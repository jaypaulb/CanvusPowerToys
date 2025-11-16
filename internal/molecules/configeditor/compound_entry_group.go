package configeditor

import (
	"fmt"

	"fyne.io/fyne/v2"
	"fyne.io/fyne/v2/container"
	"fyne.io/fyne/v2/dialog"
	"fyne.io/fyne/v2/widget"
	"gopkg.in/ini.v1"
)

// CompoundEntryGroup represents a UI group for compound entries (e.g., [server:name], [remote-desktop:name]).
type CompoundEntryGroup struct {
	pattern          string
	section          *ConfigSection
	iniFile          *ini.File
	window           fyne.Window
	entries          map[string]*CompoundEntry // Map of entry name to entry
	entriesContainer *fyne.Container           // Container for entries (stored for dynamic addition)
	onValueChange    func(section, key, value string)
	onEntryAdd       func(pattern, name string)
	onEntryRemove    func(pattern, name string)
}

// CompoundEntry represents a single compound entry instance.
type CompoundEntry struct {
	name         string
	pattern      string
	options      []*ConfigOption
	formControls map[string]*FormControl
	container    *fyne.Container
}

// NewCompoundEntryGroup creates a new compound entry group.
func NewCompoundEntryGroup(pattern string, section *ConfigSection, iniFile *ini.File, window fyne.Window, onValueChange func(section, key, value string), onEntryAdd, onEntryRemove func(pattern, name string)) *CompoundEntryGroup {
	return &CompoundEntryGroup{
		pattern:       pattern,
		section:       section,
		iniFile:       iniFile,
		window:        window,
		entries:       make(map[string]*CompoundEntry),
		onValueChange: onValueChange,
		onEntryAdd:    onEntryAdd,
		onEntryRemove: onEntryRemove,
	}
}

// CreateUI creates the UI for the compound entry group.
func (ceg *CompoundEntryGroup) CreateUI() fyne.CanvasObject {
	// Title
	title := widget.NewLabel(fmt.Sprintf("%s Entries", ceg.pattern))
	title.TextStyle.Bold = true

	// Add button
	addBtn := widget.NewButton(fmt.Sprintf("Add New %s", ceg.pattern), func() {
		ceg.showAddEntryDialog()
	})

	// Container for entries
	entriesContainer := container.NewVBox()
	ceg.entriesContainer = entriesContainer // Store reference for dynamic addition

	// Load existing entries from INI file
	if ceg.iniFile != nil {
		for _, section := range ceg.iniFile.Sections() {
			sectionName := section.Name()
			// Check if this section matches the pattern (e.g., "server:name" or "remote-desktop:name")
			if ceg.matchesPattern(sectionName) {
				entryName := ceg.extractEntryName(sectionName)
				if entryName != "" {
					ceg.addEntry(entryName, entriesContainer)
				}
			}
		}
	}

	// If no entries, add at least one empty entry
	if len(ceg.entries) == 0 {
		// For fixed-workspace, start with index 1
		initialName := ""
		if ceg.pattern == "fixed-workspace" {
			initialName = "1"
		}
		ceg.addEntry(initialName, entriesContainer)
	}

	scroll := container.NewScroll(entriesContainer)
	scroll.SetMinSize(fyne.NewSize(0, 300))

	return container.NewVBox(
		container.NewHBox(title, addBtn),
		scroll,
	)
}

// matchesPattern checks if a section name matches the compound pattern.
func (ceg *CompoundEntryGroup) matchesPattern(sectionName string) bool {
	// Check if section name starts with pattern:
	// "server:name" -> pattern "server"
	// "remote-desktop:name" -> pattern "remote-desktop"
	return len(sectionName) > len(ceg.pattern) &&
		   sectionName[:len(ceg.pattern)] == ceg.pattern &&
		   sectionName[len(ceg.pattern)] == ':'
}

// extractEntryName extracts the entry name from a section name.
func (ceg *CompoundEntryGroup) extractEntryName(sectionName string) string {
	// "server:My Server" -> "My Server"
	// "remote-desktop:Connection 1" -> "Connection 1"
	prefix := ceg.pattern + ":"
	if len(sectionName) > len(prefix) {
		return sectionName[len(prefix):]
	}
	return ""
}

// addEntry adds a new compound entry to the UI.
func (ceg *CompoundEntryGroup) addEntry(name string, parent *fyne.Container) {
	entry := &CompoundEntry{
		name:         name,
		pattern:      ceg.pattern,
		options:      ceg.getCompoundOptions(),
		formControls: make(map[string]*FormControl),
	}

	// Create name entry
	nameEntry := widget.NewEntry()
	nameEntry.SetText(name)
	nameEntry.SetPlaceHolder(fmt.Sprintf("Enter %s name", ceg.pattern))

	// Create form for options
	form := container.NewVBox()

	// Add name field
	form.Add(widget.NewLabel("Name:"))
	form.Add(nameEntry)
	form.Add(widget.NewSeparator())

	// Add all option fields
	for _, option := range entry.options {
		// Get current value from INI
		currentValue := ceg.getCurrentValue(option, name)

		formControl := CreateFormControl(option, ceg.window, currentValue)
		entry.formControls[option.Key] = formControl

		label := widget.NewRichTextFromMarkdown(fmt.Sprintf("**%s**\n%s", option.Key, option.Description))
		label.Wrapping = fyne.TextWrapWord

		form.Add(label)
		form.Add(formControl.Control)
		form.Add(widget.NewSeparator())
	}

	// Remove button
	removeBtn := widget.NewButton("Remove", func() {
		if ceg.onEntryRemove != nil {
			ceg.onEntryRemove(ceg.pattern, entry.name)
		}
		parent.Remove(entry.container)
		delete(ceg.entries, entry.name)
	})

	// Container for this entry
	entryContainer := container.NewBorder(
		container.NewHBox(
			widget.NewLabel(fmt.Sprintf("%s: %s", ceg.pattern, nameEntry.Text)),
			removeBtn,
		),
		nil, nil, nil,
		container.NewScroll(form),
	)

	entry.container = entryContainer
	ceg.entries[name] = entry

	// Update name when changed
	nameEntry.OnChanged = func(text string) {
		if text != name {
			// Update entry name
			oldName := entry.name
			entry.name = text
			delete(ceg.entries, oldName)
			ceg.entries[text] = entry
		}
	}

	parent.Add(entryContainer)
	parent.Add(widget.NewSeparator())
}

// getCompoundOptions gets the options for this compound entry type.
func (ceg *CompoundEntryGroup) getCompoundOptions() []*ConfigOption {
	var options []*ConfigOption
	for _, option := range ceg.section.Options {
		if option.IsCompound && option.Pattern == ceg.pattern {
			options = append(options, option)
		}
	}
	return options
}

// getCurrentValue gets the current value for an option from a compound entry.
func (ceg *CompoundEntryGroup) getCurrentValue(option *ConfigOption, entryName string) string {
	if ceg.iniFile == nil || entryName == "" {
		return option.Default
	}

	sectionName := fmt.Sprintf("%s:%s", ceg.pattern, entryName)
	section := ceg.iniFile.Section(sectionName)
	if section == nil {
		return option.Default
	}

	key := section.Key(option.Key)
	if key == nil || key.String() == "" {
		return option.Default
	}

	return key.String()
}

// showAddEntryDialog shows a dialog to add a new compound entry.
func (ceg *CompoundEntryGroup) showAddEntryDialog() {
	var entryName string

	// For fixed-workspace, auto-generate the next index number
	if ceg.pattern == "fixed-workspace" {
		// Find the highest existing index from both entries map and INI file
		maxIndex := 0

		// Check entries map
		for name := range ceg.entries {
			var index int
			if _, err := fmt.Sscanf(name, "%d", &index); err == nil {
				if index > maxIndex {
					maxIndex = index
				}
			}
		}

		// Check INI file sections
		if ceg.iniFile != nil {
			for _, section := range ceg.iniFile.Sections() {
				sectionName := section.Name()
				if ceg.matchesPattern(sectionName) {
					extractedName := ceg.extractEntryName(sectionName)
					var index int
					if _, err := fmt.Sscanf(extractedName, "%d", &index); err == nil {
						if index > maxIndex {
							maxIndex = index
						}
					}
				}
			}
		}

		// Use next index
		entryName = fmt.Sprintf("%d", maxIndex+1)
	} else {
		// For other patterns, show dialog to enter name
		nameEntry := widget.NewEntry()
		nameEntry.SetPlaceHolder(fmt.Sprintf("Enter %s name", ceg.pattern))

		content := container.NewVBox(
			widget.NewLabel(fmt.Sprintf("Add New %s", ceg.pattern)),
			widget.NewLabel("Name:"),
			nameEntry,
		)

		dialog.ShowCustomConfirm(
			fmt.Sprintf("Add %s", ceg.pattern),
			"Add",
			"Cancel",
			content,
			func(confirmed bool) {
				if confirmed && nameEntry.Text != "" {
					entryName = nameEntry.Text
					// Add to INI file via callback
					if ceg.onEntryAdd != nil {
						ceg.onEntryAdd(ceg.pattern, entryName)
					}
					// Add to UI immediately
					if ceg.entriesContainer != nil {
						ceg.addEntry(entryName, ceg.entriesContainer)
					}
				}
			},
			ceg.window,
		)
		return
	}

	// For fixed-workspace, add immediately without dialog
	if entryName != "" {
		// Add to INI file via callback
		if ceg.onEntryAdd != nil {
			ceg.onEntryAdd(ceg.pattern, entryName)
		}
		// Add to UI immediately
		if ceg.entriesContainer != nil {
			ceg.addEntry(entryName, ceg.entriesContainer)
		}
	}
}

// GetEntries returns all compound entries.
func (ceg *CompoundEntryGroup) GetEntries() map[string]map[string]string {
	result := make(map[string]map[string]string)
	for name, entry := range ceg.entries {
		values := make(map[string]string)
		for key, formControl := range entry.formControls {
			values[key] = formControl.GetValue()
		}
		result[name] = values
	}
	return result
}

