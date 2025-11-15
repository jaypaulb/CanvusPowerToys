package configeditor

import (
	"fmt"
	"strings"

	"fyne.io/fyne/v2"
	"fyne.io/fyne/v2/container"
	"fyne.io/fyne/v2/dialog"
	"fyne.io/fyne/v2/widget"
	"gopkg.in/ini.v1"

	"github.com/jaypaulb/CanvusPowerToys/internal/atoms/config"
	"github.com/jaypaulb/CanvusPowerToys/internal/organisms/services"
)

// Editor is the main Canvus Config Editor component.
type Editor struct {
	iniParser      *config.INIParser
	fileService    *services.FileService
	iniFile        *ini.File
	searchEntry    *widget.Entry
	optionsList    *widget.List
	optionsData    []OptionItem
	formContainer  *container.Scroll
	currentSection string
}

// OptionItem represents a configuration option.
type OptionItem struct {
	Section string
	Key     string
	Value   string
	Tooltip string
}

// NewEditor creates a new Config Editor instance.
func NewEditor(fileService *services.FileService) (*Editor, error) {
	return &Editor{
		iniParser:   config.NewINIParser(),
		fileService: fileService,
		optionsData: []OptionItem{},
	}, nil
}

// CreateUI creates the UI for the Config Editor tab.
func (e *Editor) CreateUI(window fyne.Window) fyne.CanvasObject {
	// Top: Search bar
	e.searchEntry = widget.NewEntry()
	e.searchEntry.SetPlaceHolder("Search configuration options...")
	e.searchEntry.OnChanged = func(text string) {
		e.filterOptions(text)
	}

	// Left: Options list
	e.optionsList = widget.NewList(
		func() int {
			return len(e.optionsData)
		},
		func() fyne.CanvasObject {
			return container.NewHBox(
				widget.NewLabel("Section.Key"),
				widget.NewLabel("Value"),
			)
		},
		func(id widget.ListItemID, obj fyne.CanvasObject) {
			item := e.optionsData[id]
			box := obj.(*fyne.Container)
			labels := box.Objects
			labels[0].(*widget.Label).SetText(fmt.Sprintf("%s.%s", item.Section, item.Key))
			labels[1].(*widget.Label).SetText(item.Value)
		},
	)
	e.optionsList.OnSelected = func(id widget.ListItemID) {
		e.showOptionForm(id)
	}

	// Right: Form container
	e.formContainer = container.NewScroll(widget.NewLabel("Select an option to edit"))
	e.formContainer.SetMinSize(fyne.NewSize(400, 0))

	// Load button
	loadBtn := widget.NewButton("Load mt-canvus.ini", func() {
		e.loadConfigFile(window)
	})

	// Save buttons
	saveUserBtn := widget.NewButton("Save to User Config", func() {
		e.saveConfig(window, true)
	})
	saveSystemBtn := widget.NewButton("Save to System Config", func() {
		e.saveConfig(window, false)
	})

	// Layout
	topBar := container.NewHBox(
		loadBtn,
		widget.NewSeparator(),
		saveUserBtn,
		saveSystemBtn,
	)

	leftPanel := container.NewBorder(
		e.searchEntry,
		nil, nil, nil,
		e.optionsList,
	)

	split := container.NewHSplit(leftPanel, e.formContainer)
	split.SetOffset(0.4)

	return container.NewBorder(
		topBar,
		nil, nil, nil,
		split,
	)
}

// loadConfigFile loads mt-canvus.ini from detected location.
func (e *Editor) loadConfigFile(window fyne.Window) {
	iniPath := e.fileService.DetectMtCanvusIni()
	if iniPath == "" {
		dialog.ShowInformation("Not Found", "mt-canvus.ini not found in standard locations", window)
		return
	}

	iniFile, err := e.iniParser.Read(iniPath)
	if err != nil {
		dialog.ShowError(fmt.Errorf("failed to load mt-canvus.ini: %w", err), window)
		return
	}

	e.iniFile = iniFile
	e.populateOptionsList()
	dialog.ShowInformation("Loaded", fmt.Sprintf("Loaded mt-canvus.ini from:\n%s", iniPath), window)
}

// populateOptionsList populates the options list from the INI file.
func (e *Editor) populateOptionsList() {
	// Store original data before filtering
	originalData := []OptionItem{}

	if e.iniFile == nil {
		e.optionsData = originalData
		e.optionsList.Refresh()
		return
	}

	// Iterate through all sections
	for _, section := range e.iniFile.Sections() {
		sectionName := section.Name()
		if sectionName == "DEFAULT" {
			continue
		}

		// Get all keys in this section
		for _, key := range section.Keys() {
			value := key.Value()
			tooltip := e.getTooltip(sectionName, key.Name())

			originalData = append(originalData, OptionItem{
				Section: sectionName,
				Key:     key.Name(),
				Value:   value,
				Tooltip: tooltip,
			})
		}
	}

	// Apply current filter if any
	searchText := e.searchEntry.Text
	if searchText == "" {
		e.optionsData = originalData
	} else {
		e.filterOptions(searchText)
		return
	}

	e.optionsList.Refresh()
}

// filterOptions filters the options list based on search text.
func (e *Editor) filterOptions(searchText string) {
	if e.iniFile == nil {
		return
	}

	searchLower := strings.ToLower(searchText)
	filtered := []OptionItem{}

	// Rebuild from INI file
	for _, section := range e.iniFile.Sections() {
		sectionName := section.Name()
		if sectionName == "DEFAULT" {
			continue
		}

		for _, key := range section.Keys() {
			item := OptionItem{
				Section: sectionName,
				Key:     key.Name(),
				Value:   key.Value(),
				Tooltip: e.getTooltip(sectionName, key.Name()),
			}

			if strings.Contains(strings.ToLower(item.Section), searchLower) ||
				strings.Contains(strings.ToLower(item.Key), searchLower) ||
				strings.Contains(strings.ToLower(item.Value), searchLower) {
				filtered = append(filtered, item)
			}
		}
	}

	e.optionsData = filtered
	e.optionsList.Refresh()
}

// showOptionForm shows the form for editing an option.
func (e *Editor) showOptionForm(id widget.ListItemID) {
	if id < 0 || id >= len(e.optionsData) {
		return
	}

	item := e.optionsData[id]
	form := e.createOptionForm(item)
	e.formContainer.Content = form
	e.formContainer.Refresh()
}

// createOptionForm creates a form for editing an option.
func (e *Editor) createOptionForm(item OptionItem) fyne.CanvasObject {
	sectionLabel := widget.NewLabel(fmt.Sprintf("Section: %s", item.Section))
	keyLabel := widget.NewLabel(fmt.Sprintf("Key: %s", item.Key))

	valueEntry := widget.NewEntry()
	valueEntry.SetText(item.Value)

	tooltipLabel := widget.NewLabel(item.Tooltip)
	tooltipLabel.Wrapping = fyne.TextWrapWord

	saveBtn := widget.NewButton("Save", func() {
		e.saveOptionValue(item.Section, item.Key, valueEntry.Text)
	})

	return container.NewVBox(
		sectionLabel,
		keyLabel,
		widget.NewSeparator(),
		widget.NewLabel("Value:"),
		valueEntry,
		widget.NewSeparator(),
		widget.NewLabel("Description:"),
		tooltipLabel,
		widget.NewSeparator(),
		saveBtn,
	)
}

// saveOptionValue saves an option value to the INI file in memory.
func (e *Editor) saveOptionValue(section, key, value string) {
	if e.iniFile == nil {
		return
	}

	sec, err := e.iniFile.GetSection(section)
	if err != nil {
		sec, _ = e.iniFile.NewSection(section)
	}

	sec.Key(key).SetValue(value)
	e.populateOptionsList()
}

// saveConfig saves the configuration to file.
func (e *Editor) saveConfig(window fyne.Window, userConfig bool) {
	if e.iniFile == nil {
		dialog.ShowError(fmt.Errorf("no configuration loaded"), window)
		return
	}

	var savePath string
	if userConfig {
		savePath = e.fileService.GetUserConfigPath() + "/mt-canvus.ini"
	} else {
		savePath = e.fileService.GetSystemConfigPath() + "/mt-canvus.ini"
	}

	// Ensure directory exists
	if err := e.fileService.EnsureDirectory(e.fileService.GetUserConfigPath()); err != nil {
		dialog.ShowError(fmt.Errorf("failed to create directory: %w", err), window)
		return
	}

	if err := e.iniParser.Write(e.iniFile, savePath); err != nil {
		dialog.ShowError(fmt.Errorf("failed to save mt-canvus.ini: %w", err), window)
		return
	}

	location := "user config"
	if !userConfig {
		location = "system config"
	}
	dialog.ShowInformation("Saved", fmt.Sprintf("Saved mt-canvus.ini to %s:\n%s", location, savePath), window)
}

// getTooltip returns a tooltip for a configuration option.
func (e *Editor) getTooltip(section, key string) string {
	// TODO: Load tooltips from documentation
	// For now, return a basic description
	return fmt.Sprintf("Configuration option: %s.%s", section, key)
}
