//go:build windows
// +build windows

package tray

import (
	"reflect"
	"syscall"

	"fyne.io/fyne/v2"
	"golang.org/x/sys/windows"
)

var (
	user32               = windows.NewLazyDLL("user32.dll")
	procGetWindowLongPtr = user32.NewProc("GetWindowLongPtrW")
	procSetWindowLongPtr = user32.NewProc("SetWindowLongPtrW")
	procCallWindowProc   = user32.NewProc("CallWindowProcW")
	procFindWindow       = user32.NewProc("FindWindowW")
)

const (
	GWL_WNDPROC   = -4
	WM_SYSCOMMAND = 0x0112
	SC_MINIMIZE   = 0xF020
)

var (
	originalWndProc uintptr
	managerInstance *Manager
)

// setupMinimizeToTray sets up Windows message hook to intercept minimize
// and hide window to tray instead
func setupMinimizeToTray(m *Manager) error {
	managerInstance = m

	// Get HWND from Fyne window using reflection
	hwnd, err := getWindowHandle(m.window)
	if err != nil {
		// If we can't get HWND, fall back to polling approach
		return err
	}

	// Get current window procedure
	// GWL_WNDPROC is -4, use local variable to handle negative value conversion
	gwlWndProc := int32(GWL_WNDPROC)
	ret, _, _ := procGetWindowLongPtr.Call(uintptr(hwnd), uintptr(gwlWndProc))
	if ret == 0 {
		return syscall.GetLastError()
	}
	originalWndProc = ret

	// Set our custom window procedure
	newWndProc := syscall.NewCallback(windowProc)
	ret, _, _ = procSetWindowLongPtr.Call(uintptr(hwnd), uintptr(gwlWndProc), newWndProc)
	if ret == 0 {
		return syscall.GetLastError()
	}

	return nil
}

// getWindowHandle extracts the HWND from Fyne's window using reflection
func getWindowHandle(window fyne.Window) (syscall.Handle, error) {
	// Use reflection to access Fyne's internal window structure
	// Fyne's window driver stores the HWND in the window object
	val := reflect.ValueOf(window)
	if val.Kind() == reflect.Ptr {
		val = val.Elem()
	}

	// Try to find HWND field in the window structure
	// This is platform-specific and may vary with Fyne versions
	typ := val.Type()
	for i := 0; i < typ.NumField(); i++ {
		field := typ.Field(i)
		if field.Type.String() == "syscall.Handle" || field.Type.String() == "uintptr" {
			fieldVal := val.Field(i)
			if fieldVal.CanInterface() {
				if hwnd, ok := fieldVal.Interface().(syscall.Handle); ok {
					return hwnd, nil
				}
				if hwnd, ok := fieldVal.Interface().(uintptr); ok {
					return syscall.Handle(hwnd), nil
				}
			}
		}
	}

	// Alternative: Try to get from driver
	// Fyne stores window in driver, which has the HWND
	driverVal := reflect.ValueOf(window).MethodByName("Driver")
	if driverVal.IsValid() {
		driver := driverVal.Call(nil)[0]
		// Look for window or viewport field in driver
		driverTyp := driver.Type()
		for i := 0; i < driverTyp.NumField(); i++ {
			field := driverTyp.Field(i)
			if field.Name == "viewport" || field.Name == "window" {
				fieldVal := driver.Field(i)
				// Recursively search for HWND
				if hwnd, err := findHWNDInStruct(fieldVal); err == nil {
					return hwnd, nil
				}
			}
		}
	}

	return 0, syscall.EINVAL
}

// findHWNDInStruct recursively searches for HWND in a struct
func findHWNDInStruct(val reflect.Value) (syscall.Handle, error) {
	if val.Kind() == reflect.Ptr {
		if val.IsNil() {
			return 0, syscall.EINVAL
		}
		val = val.Elem()
	}

	if val.Kind() != reflect.Struct {
		return 0, syscall.EINVAL
	}

	typ := val.Type()
	for i := 0; i < typ.NumField(); i++ {
		field := typ.Field(i)
		fieldVal := val.Field(i)

		if !fieldVal.CanInterface() {
			continue
		}

		switch field.Type.Kind() {
		case reflect.Uintptr:
			if hwnd, ok := fieldVal.Interface().(uintptr); ok && hwnd != 0 {
				return syscall.Handle(hwnd), nil
			}
		case reflect.Struct, reflect.Ptr:
			if hwnd, err := findHWNDInStruct(fieldVal); err == nil {
				return hwnd, nil
			}
		}
	}

	return 0, syscall.EINVAL
}

// windowProc is the Windows window procedure that intercepts messages
func windowProc(hwnd syscall.Handle, msg uint32, wParam, lParam uintptr) uintptr {
	if msg == WM_SYSCOMMAND && wParam == SC_MINIMIZE {
		// Intercept minimize - hide window to tray instead
		if managerInstance != nil && managerInstance.window != nil {
			managerInstance.window.Hide()
		}
		return 0
	}
	// Call original window procedure for other messages
	if originalWndProc != 0 {
		ret, _, _ := procCallWindowProc.Call(originalWndProc, uintptr(hwnd), uintptr(msg), wParam, lParam)
		return ret
	}
	return 0
}

