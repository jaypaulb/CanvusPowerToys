package custommenu

import (
	"fmt"
	"os"
	"path/filepath"
	"strconv"
	"strings"

	"fyne.io/fyne/v2"
	"fyne.io/fyne/v2/container"
	"fyne.io/fyne/v2/dialog"
	"fyne.io/fyne/v2/widget"

	"github.com/jaypaulb/CanvusPowerToys/internal/atoms/backup"
	"github.com/jaypaulb/CanvusPowerToys/internal/atoms/config"
	"github.com/jaypaulb/CanvusPowerToys/internal/organisms/services"
)

// MenuItem represents a menu item in the YAML structure.
type MenuItem struct {
	Tooltip  string     `yaml:"tooltip"`
	Icon     string     `yaml:"icon,omitempty"`
	IconBack string     `yaml:"icon-back,omitempty"`
	Items    []MenuItem `yaml:"items,omitempty"`
	Actions  []Action   `yaml:"actions,omitempty"`
}

// Action represents a menu item action.
type Action struct {
	Name       string                 `yaml:"name"` // "create" or "open-folder"
	Parameters map[string]interface{} `yaml:"parameters,omitempty"`
}

// Designer handles custom menu design and YAML generation.
type Designer struct {
	fileService   *services.FileService
	yamlHandler   *config.YAMLHandler
	iniParser     *config.INIParser
	backupManager *backup.Manager
	menuTree      *widget.Tree
	menuData      []MenuItem
	formContainer *container.Scroll
}

// NewDesigner creates a new Custom Menu Designer.
func NewDesigner(fileService *services.FileService) (*Designer, error) {
	return &Designer{
		fileService:   fileService,
		yamlHandler:   config.NewYAMLHandler(),
		iniParser:     config.NewINIParser(),
		backupManager: backup.NewManager(""),
		menuData:      []MenuItem{},
	}, nil
}

// CreateUI creates the UI for the Custom Menu Designer tab.
func (d *Designer) CreateUI(window fyne.Window) fyne.CanvasObject {
	title := widget.NewLabel("Custom Menu Designer")
	title.TextStyle = fyne.TextStyle{Bold: true}

	instructions := widget.NewRichTextFromMarkdown(`
**Custom Menu Designer**

Create and edit custom menus for Canvus. Menus are saved as menu.yml files.

**Features:**
- Hierarchical menu structure with unlimited nesting
- Icon support (937x937 PNG)
- Actions: create (note, pdf, video, image, browser) and open-folder
`)

	// Menu tree view
	d.menuTree = widget.NewTree(
		func(id widget.TreeNodeID) []widget.TreeNodeID {
			return d.getChildNodes(id)
		},
		func(id widget.TreeNodeID) bool {
			return d.hasChildren(id)
		},
		func(branch bool) fyne.CanvasObject {
			return widget.NewLabel("Menu Item")
		},
		func(id widget.TreeNodeID, branch bool, obj fyne.CanvasObject) {
			label := obj.(*widget.Label)
			item := d.getItemByID(id)
			if item != nil {
				label.SetText(item.Tooltip)
			}
		},
	)

	d.menuTree.OnSelected = func(id widget.TreeNodeID) {
		d.showItemForm(id)
	}

	// Form container
	d.formContainer = container.NewScroll(widget.NewLabel("Select a menu item to edit"))
	d.formContainer.SetMinSize(fyne.NewSize(400, 0))

	// Buttons
	addItemBtn := widget.NewButton("Add Menu Item", func() {
		d.addMenuItem(window)
	})

	importBtn := widget.NewButton("Import menu.yml", func() {
		d.importMenu(window)
	})

	saveBtn := widget.NewButton("Save menu.yml", func() {
		d.saveMenu(window)
	})

	// Layout
	leftPanel := container.NewBorder(
		container.NewVBox(
			title,
			instructions,
			widget.NewSeparator(),
			container.NewHBox(addItemBtn, importBtn),
		),
		nil, nil, nil,
		d.menuTree,
	)

	split := container.NewHSplit(leftPanel, d.formContainer)
	split.SetOffset(0.4)

	return container.NewBorder(
		container.NewHBox(saveBtn),
		nil, nil, nil,
		split,
	)
}

// getChildNodes returns child node IDs for a given parent.
func (d *Designer) getChildNodes(id widget.TreeNodeID) []widget.TreeNodeID {
	if id == "" {
		// Root level - return top-level items
		var ids []widget.TreeNodeID
		for i := range d.menuData {
			ids = append(ids, fmt.Sprintf("%d", i))
		}
		return ids
	}

	// Get item and return its children
	item := d.getItemByID(id)
	if item != nil && len(item.Items) > 0 {
		var ids []widget.TreeNodeID
		for i := range item.Items {
			ids = append(ids, fmt.Sprintf("%s.%d", id, i))
		}
		return ids
	}
	return nil
}

// hasChildren checks if a node has children.
func (d *Designer) hasChildren(id widget.TreeNodeID) bool {
	item := d.getItemByID(id)
	return item != nil && len(item.Items) > 0
}

// getItemByID gets a menu item by its ID.
func (d *Designer) getItemByID(id widget.TreeNodeID) *MenuItem {
	if id == "" {
		return nil
	}

	// Parse ID (format: "0", "0.1", "0.1.2", etc.)
	parts := d.parseID(id)
	if len(parts) == 0 {
		return nil
	}

	item := &d.menuData[parts[0]]
	for i := 1; i < len(parts); i++ {
		if parts[i] >= len(item.Items) {
			return nil
		}
		item = &item.Items[parts[i]]
	}
	return item
}

// parseID parses a node ID into indices.
func (d *Designer) parseID(id widget.TreeNodeID) []int {
	if id == "" {
		return nil
	}
	parts := strings.Split(string(id), ".")
	indices := make([]int, len(parts))
	for i, part := range parts {
		idx, err := strconv.Atoi(part)
		if err != nil {
			return nil
		}
		indices[i] = idx
	}
	return indices
}

// showItemForm shows the form for editing a menu item.
func (d *Designer) showItemForm(id widget.TreeNodeID) {
	item := d.getItemByID(id)
	if item == nil {
		return
	}

	form := d.createItemForm(item, id)
	d.formContainer.Content = form
	d.formContainer.Refresh()
}

// createItemForm creates a form for editing a menu item.
func (d *Designer) createItemForm(item *MenuItem, id widget.TreeNodeID) fyne.CanvasObject {
	tooltipEntry := widget.NewEntry()
	tooltipEntry.SetText(item.Tooltip)

	iconEntry := widget.NewEntry()
	iconEntry.SetText(item.Icon)
	iconBrowseBtn := widget.NewButton("Browse...", func() {
		d.browseIcon()
	})

	actionTypeSelect := widget.NewSelect([]string{"create", "open-folder"}, nil)
	if len(item.Actions) > 0 {
		actionTypeSelect.SetSelected(item.Actions[0].Name)
	}

	contentTypeSelect := widget.NewSelect([]string{"note", "pdf", "video", "image", "browser"}, nil)
	if len(item.Actions) > 0 && item.Actions[0].Name == "create" {
		if params := item.Actions[0].Parameters; params != nil {
			if t, ok := params["type"].(string); ok {
				contentTypeSelect.SetSelected(t)
			}
		}
	}

	saveBtn := widget.NewButton("Save", func() {
		item.Tooltip = tooltipEntry.Text
		item.Icon = iconEntry.Text
		if actionTypeSelect.Selected != "" {
			action := Action{Name: actionTypeSelect.Selected}
			if actionTypeSelect.Selected == "create" && contentTypeSelect.Selected != "" {
				action.Parameters = map[string]interface{}{
					"type": contentTypeSelect.Selected,
				}
			}
			item.Actions = []Action{action}
		}
		d.menuTree.Refresh()
	})

	addSubItemBtn := widget.NewButton("Add Sub-Item", func() {
		item.Items = append(item.Items, MenuItem{Tooltip: "New Sub-Item"})
		d.menuTree.Refresh()
	})

	return container.NewVBox(
		widget.NewLabel("Menu Item Editor"),
		widget.NewSeparator(),
		widget.NewLabel("Tooltip:"),
		tooltipEntry,
		widget.NewLabel("Icon:"),
		container.NewHBox(iconEntry, iconBrowseBtn),
		widget.NewLabel("Action Type:"),
		actionTypeSelect,
		widget.NewLabel("Content Type (for create):"),
		contentTypeSelect,
		widget.NewSeparator(),
		saveBtn,
		addSubItemBtn,
	)
}

// addMenuItem adds a new top-level menu item.
func (d *Designer) addMenuItem(window fyne.Window) {
	d.menuData = append(d.menuData, MenuItem{
		Tooltip: "New Menu Item",
	})
	d.menuTree.Refresh()
}

// browseIcon opens a file browser for icon selection.
func (d *Designer) browseIcon() {
	// TODO: Implement file browser dialog
}

// importMenu imports an existing menu.yml file.
func (d *Designer) importMenu(window fyne.Window) {
	dialog.ShowFileOpen(func(reader fyne.URIReadCloser, err error) {
		if err != nil || reader == nil {
			return
		}
		defer reader.Close()

		// Read YAML file - structure is root item with items array
		var rootItem MenuItem
		if err := d.yamlHandler.Read(reader.URI().Path(), &rootItem); err != nil {
			dialog.ShowError(fmt.Errorf("failed to import menu.yml: %w", err), window)
			return
		}

		d.menuData = rootItem.Items
		d.menuTree.Refresh()
		dialog.ShowInformation("Imported", "menu.yml imported successfully", window)
	}, window)
}

// saveMenu saves the menu to menu.yml and updates mt-canvus.ini.
func (d *Designer) saveMenu(window fyne.Window) {
	configDir := d.fileService.GetUserConfigPath()
	menuPath := filepath.Join(configDir, "menu.yml")

	// Ensure directory exists
	if err := os.MkdirAll(configDir, 0755); err != nil {
		dialog.ShowError(fmt.Errorf("failed to create directory: %w", err), window)
		return
	}

	// Create backup if file exists
	if _, err := os.Stat(menuPath); err == nil {
		if _, err := d.backupManager.CreateBackup(menuPath); err != nil {
			fmt.Printf("Warning: Failed to create backup: %v\n", err)
		}
	}

	// Create root structure with items
	rootItem := MenuItem{
		Tooltip: "Custom Menu",
		Icon:    "icons/custom-menu.png",
		Items:   d.menuData,
	}

	// Save YAML
	if err := d.yamlHandler.Write(menuPath, rootItem); err != nil {
		dialog.ShowError(fmt.Errorf("failed to save menu.yml: %w", err), window)
		return
	}

	// Update mt-canvus.ini
	iniPath := d.fileService.DetectMtCanvusIni()
	if iniPath != "" {
		if err := d.updateCustomMenuIni(iniPath, menuPath); err != nil {
			dialog.ShowError(fmt.Errorf("failed to update mt-canvus.ini: %w", err), window)
			return
		}
	}

	dialog.ShowInformation("Saved", fmt.Sprintf("menu.yml saved to:\n%s", menuPath), window)
}

// updateCustomMenuIni updates the custom-menu entry in mt-canvus.ini.
func (d *Designer) updateCustomMenuIni(iniPath, menuPath string) error {
	iniFile, err := d.iniParser.Read(iniPath)
	if err != nil {
		return fmt.Errorf("failed to read mt-canvus.ini: %w", err)
	}

	// Get or create [canvas] section
	section, err := iniFile.GetSection("canvas")
	if err != nil {
		section, _ = iniFile.NewSection("canvas")
	}

	// Set custom-menu to relative path
	relPath, err := filepath.Rel(filepath.Dir(iniPath), menuPath)
	if err != nil {
		relPath = menuPath // Fallback to absolute path
	}

	section.Key("custom-menu").SetValue(relPath)

	// Create backup before updating
	if _, err := os.Stat(iniPath); err == nil {
		if _, err := d.backupManager.CreateBackup(iniPath); err != nil {
			fmt.Printf("Warning: Failed to create backup: %v\n", err)
		}
	}

	return d.iniParser.Write(iniFile, iniPath)
}
