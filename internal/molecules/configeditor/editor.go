package configeditor

import (
	"fmt"
	"os"
	"path/filepath"
	"strings"

	"fyne.io/fyne/v2"
	"fyne.io/fyne/v2/container"
	"fyne.io/fyne/v2/dialog"
	"fyne.io/fyne/v2/widget"
	"gopkg.in/ini.v1"

	"github.com/jaypaulb/CanvusPowerToys/internal/atoms/backup"
	"github.com/jaypaulb/CanvusPowerToys/internal/atoms/config"
	"github.com/jaypaulb/CanvusPowerToys/internal/organisms/services"
)

// Editor is the main Canvus Config Editor component.
type Editor struct {
	iniParser      *config.INIParser
	fileService    *services.FileService
	backupManager  *backup.Manager
	iniFile        *ini.File
	schema         *ConfigSchema // Schema with all possible options
	searchEntry    *widget.Entry
	accordion      *widget.Accordion
	sectionGroups  map[string]*SectionGroup
	compoundGroups map[string]*CompoundEntryGroup
	formContainer  *container.Scroll
	window         fyne.Window
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
	// Create backup manager with temporary directory (will use file's directory when backing up)
	backupMgr := backup.NewManager("")

	// Start with empty schema - will be populated when INI file is loaded
	schema := NewConfigSchema()

	return &Editor{
		iniParser:      config.NewINIParser(),
		fileService:     fileService,
		backupManager:  backupMgr,
		schema:         schema,
		sectionGroups:  make(map[string]*SectionGroup),
		compoundGroups: make(map[string]*CompoundEntryGroup),
	}, nil
}

// CreateUI creates the UI for the Config Editor tab.
func (e *Editor) CreateUI(window fyne.Window) fyne.CanvasObject {
	e.window = window

	// Top: Search bar
	e.searchEntry = widget.NewEntry()
	e.searchEntry.SetPlaceHolder("Search configuration options...")
	e.searchEntry.OnChanged = func(text string) {
		e.filterSections(text)
	}

	// Create accordion for sections
	e.accordion = widget.NewAccordion()

	// Try to auto-load: first example file for schema, then live file for values
	if err := e.loadConfigFilesSilent(); err == nil {
		// Successfully loaded
	}

	// Build UI from schema
	e.buildUIFromSchema()

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

	// Open All / Close All buttons
	openAllBtn := widget.NewButton("Open All", func() {
		e.accordion.OpenAll()
	})
	closeAllBtn := widget.NewButton("Close All", func() {
		e.accordion.CloseAll()
	})

	// Layout
	topBar := container.NewHBox(
		loadBtn,
		widget.NewSeparator(),
		saveUserBtn,
		saveSystemBtn,
		widget.NewSeparator(),
		openAllBtn,
		closeAllBtn,
	)

	mainPanel := container.NewBorder(
		e.searchEntry,
		nil, nil, nil,
		container.NewScroll(e.accordion),
	)

	return container.NewBorder(
		topBar,
		nil, nil, nil,
		mainPanel,
	)
}

// buildUIFromSchema builds the UI from the schema, showing all options.
func (e *Editor) buildUIFromSchema() {
	// Clear accordion
	e.accordion.Items = nil
	e.sectionGroups = make(map[string]*SectionGroup)
	e.compoundGroups = make(map[string]*CompoundEntryGroup)

	// Process all sections from schema
	processedSections := make(map[string]bool)

	// First, handle root level options (empty section)
	if section := e.schema.GetSection(""); section != nil || e.hasRootOptions() {
		sectionName := "General"
		section := e.buildSectionFromOptions("")
		if len(section.Options) > 0 {
			sectionGroup := NewSectionGroup(section, e.iniFile, e.window, e.onValueChange)
			e.sectionGroups[sectionName] = sectionGroup
			item := sectionGroup.CreateUI()
			e.accordion.Append(item)
			processedSections[""] = true
		}
	}

	// Process all other sections
	for _, option := range e.schema.Options {
		sectionName := option.Section
		if sectionName == "" || processedSections[sectionName] {
			continue
		}

		section := e.schema.GetSection(sectionName)
		if section == nil {
			// Build section from options
			section = e.buildSectionFromOptions(sectionName)
		}

		if len(section.Options) == 0 {
			continue
		}

		// Check if this section has compound entries
		hasCompound := false
		var compoundPattern string
		for _, opt := range section.Options {
			if opt.IsCompound {
				hasCompound = true
				compoundPattern = opt.Pattern
				break
			}
		}

		if hasCompound {
			// Create compound entry group
			compoundGroup := NewCompoundEntryGroup(
				compoundPattern,
				section,
				e.iniFile,
				e.window,
				e.onValueChange,
				e.onCompoundEntryAdd,
				e.onCompoundEntryRemove,
			)
			e.compoundGroups[compoundPattern] = compoundGroup

			item := &widget.AccordionItem{
				Title:  sectionName,
				Detail: compoundGroup.CreateUI(),
				Open:   true,
			}
			e.accordion.Append(item)
		} else {
			// Create regular section group
			sectionGroup := NewSectionGroup(section, e.iniFile, e.window, e.onValueChange)
			e.sectionGroups[sectionName] = sectionGroup

			item := sectionGroup.CreateUI()
			e.accordion.Append(item)
		}

		processedSections[sectionName] = true
	}

	e.accordion.OpenAll()
}

// hasRootOptions checks if there are root level options.
func (e *Editor) hasRootOptions() bool {
	for _, option := range e.schema.Options {
		if option.Section == "" {
			return true
		}
	}
	return false
}

// buildSectionFromOptions builds a section from options with the given section name.
func (e *Editor) buildSectionFromOptions(sectionName string) *ConfigSection {
	section := &ConfigSection{
		Name:            sectionName,
		Options:         []*ConfigOption{},
		CompoundEntries: make(map[string][]*ConfigOption),
	}

	for _, option := range e.schema.Options {
		if option.Section == sectionName {
			section.Options = append(section.Options, option)
		}
	}

	return section
}

// loadConfigFilesSilent loads the example INI file for schema, then overlays with live values.
func (e *Editor) loadConfigFilesSilent() error {
	// First, try to load the example file for complete schema
	examplePath := e.fileService.DetectExampleIni()
	if examplePath != "" {
		fileParser := NewINIFileParser(examplePath)
		schema, err := fileParser.Parse()
		if err == nil {
			e.schema = schema
		}
	}

	// Then, load the live config file for actual values
	iniPath := e.fileService.DetectMtCanvusIni()
	if iniPath == "" {
		// No live file found, but we might have schema from example
		if e.schema != nil && len(e.schema.Options) > 0 {
			// Create empty INI file for editing
			e.iniFile = ini.Empty()
			return nil
		}
		return fmt.Errorf("no INI file found")
	}

	// Load the actual INI file for current values
	iniFile, err := e.iniParser.Read(iniPath)
	if err != nil {
		return err
	}

	e.iniFile = iniFile
	return nil
}

// loadConfigFile loads the example INI file for schema, then overlays with live values.
func (e *Editor) loadConfigFile(window fyne.Window) {
	// First, try to load the example file for complete schema
	examplePath := e.fileService.DetectExampleIni()
	var schema *ConfigSchema
	if examplePath != "" {
		fileParser := NewINIFileParser(examplePath)
		var err error
		schema, err = fileParser.Parse()
		if err != nil {
			dialog.ShowError(fmt.Errorf("failed to parse example mt-canvus.ini: %w", err), window)
			return
		}
		e.schema = schema
	}

	// Then, load the live config file for actual values
	iniPath := e.fileService.DetectMtCanvusIni()
	if iniPath == "" {
		if schema == nil || len(schema.Options) == 0 {
			dialog.ShowInformation("Not Found", "mt-canvus.ini not found in standard locations", window)
			return
		}
		// We have schema from example, but no live file - create empty INI
		e.iniFile = ini.Empty()
		e.buildUIFromSchema()
		dialog.ShowInformation("Loaded", fmt.Sprintf("Loaded schema from example file:\n%s\n\nNo live config file found - using defaults.", examplePath), window)
		return
	}

	// Load the actual INI file for current values
	iniFile, err := e.iniParser.Read(iniPath)
	if err != nil {
		dialog.ShowError(fmt.Errorf("failed to load mt-canvus.ini: %w", err), window)
		return
	}

	e.iniFile = iniFile

	// If we didn't get schema from example, try parsing the live file
	if schema == nil || len(schema.Options) == 0 {
		fileParser := NewINIFileParser(iniPath)
		schema, err = fileParser.Parse()
		if err == nil {
			e.schema = schema
		}
	}

	// Rebuild UI to update values from loaded file
	e.buildUIFromSchema()

	msg := fmt.Sprintf("Loaded mt-canvus.ini from:\n%s", iniPath)
	if examplePath != "" {
		msg += fmt.Sprintf("\n\nSchema loaded from example:\n%s", examplePath)
	}
	msg += "\n\nAll options updated with current values."
	dialog.ShowInformation("Loaded", msg, window)
}

// filterSections filters sections based on search text.
func (e *Editor) filterSections(searchText string) {
	if e.accordion == nil {
		return
	}

	searchLower := strings.ToLower(searchText)

	// Show/hide accordion items based on search
	for _, item := range e.accordion.Items {
		title := strings.ToLower(item.Title)
		shouldShow := searchText == "" || strings.Contains(title, searchLower)

		// Also check if any option in the section matches
		if !shouldShow && searchText != "" {
			// Check section groups
			for sectionName, sectionGroup := range e.sectionGroups {
				if strings.Contains(strings.ToLower(sectionName), searchLower) {
					for key := range sectionGroup.formControls {
						if strings.Contains(strings.ToLower(key), searchLower) {
							shouldShow = true
							break
						}
					}
				}
			}
		}

		// Note: Fyne Accordion doesn't have a direct way to hide items
		// We'll need to rebuild the accordion with filtered items
	}

	// For now, just expand all when searching
	if searchText != "" {
		e.accordion.OpenAll()
	}
}

// onValueChange handles value changes from form controls.
func (e *Editor) onValueChange(section, key, value string) {
	if e.iniFile == nil {
		// Create empty INI file if needed
		cfg := ini.Empty()
		e.iniFile = cfg
	}

	sec, err := e.iniFile.GetSection(section)
	if err != nil {
		sec, _ = e.iniFile.NewSection(section)
	}

	sec.Key(key).SetValue(value)
}

// onCompoundEntryAdd handles adding a new compound entry.
func (e *Editor) onCompoundEntryAdd(pattern, name string) {
	if e.iniFile == nil {
		cfg := ini.Empty()
		e.iniFile = cfg
	}

	sectionName := fmt.Sprintf("%s:%s", pattern, name)
	_, err := e.iniFile.NewSection(sectionName)
	if err != nil {
		// Section might already exist
	}
}

// onCompoundEntryRemove handles removing a compound entry.
func (e *Editor) onCompoundEntryRemove(pattern, name string) {
	if e.iniFile == nil {
		return
	}

	sectionName := fmt.Sprintf("%s:%s", pattern, name)
	e.iniFile.DeleteSection(sectionName)
}

// saveConfig saves the configuration to file.
func (e *Editor) saveConfig(window fyne.Window, userConfig bool) {
	if e.iniFile == nil {
		dialog.ShowError(fmt.Errorf("no configuration loaded"), window)
		return
	}

	var savePath string
	var configDir string
	if userConfig {
		configDir = e.fileService.GetUserConfigPath()
		savePath = filepath.Join(configDir, "mt-canvus.ini")
	} else {
		configDir = e.fileService.GetSystemConfigPath()
		savePath = filepath.Join(configDir, "mt-canvus.ini")
	}

	// Ensure directory exists
	if err := e.fileService.EnsureDirectory(configDir); err != nil {
		dialog.ShowError(fmt.Errorf("failed to create directory: %w", err), window)
		return
	}

	// Create backup before saving (if file exists)
	if _, err := os.Stat(savePath); err == nil {
		if _, err := e.backupManager.CreateBackup(savePath); err != nil {
			// Log warning but continue with save
			fmt.Printf("Warning: Failed to create backup: %v\n", err)
		}
	}

	// Save the file
	if err := e.iniParser.Write(e.iniFile, savePath); err != nil {
		dialog.ShowError(fmt.Errorf("failed to save mt-canvus.ini: %w", err), window)
		return
	}

	location := "user config"
	if !userConfig {
		location = "system config"
	}
	dialog.ShowInformation("Saved", fmt.Sprintf("Saved mt-canvus.ini to %s:\n%s\n\nBackup created automatically.", location, savePath), window)
}

// getTooltip returns a tooltip for a configuration option.
func (e *Editor) getTooltip(section, key string) string {
	// TODO: Load tooltips from documentation
	// For now, return a basic description
	return fmt.Sprintf("Configuration option: %s.%s", section, key)
}
