package configeditor

// GetEmbeddedSchema returns a manually created comprehensive schema
// based on the mt-canvus.ini documentation and example file.
// This schema is manually maintained and should be updated when new settings are added.
func GetEmbeddedSchema() *ConfigSchema {
	schema := NewConfigSchema()

	// ============================================================================
	// ROOT SECTION (no section header)
	// ============================================================================

	// lock-config
	schema.AddOption(&ConfigOption{
		Section:     "",
		Key:         "lock-config",
		Default:     "false",
		Description: "Prevent Canvus client from making any changes to this configuration file. This option is primarily meant for public installations.",
		Type:        ValueTypeBoolean,
		EnumValues:  []string{"true", "false"},
	})

	// multi-user-mode-enabled
	schema.AddOption(&ConfigOption{
		Section:     "",
		Key:         "multi-user-mode-enabled",
		Default:     "auto",
		Description: "Enabling multi-user mode optimizes the application user experience for large, multi-user installations. Disable it for a single user experience on personal devices like laptops.",
		Type:        ValueTypeEnum,
		EnumValues:  []string{"true", "false", "auto"},
	})

	// virtual-keyboard-enabled
	schema.AddOption(&ConfigOption{
		Section:     "",
		Key:         "virtual-keyboard-enabled",
		Default:     "auto",
		Description: "Should Canvus use a built-in virtual keyboard. Virtual keyboard is opened automatically if true, not opened automatically if false, or used if multi-user mode is enabled or when Windows 10 Tablet mode is enabled if auto.",
		Type:        ValueTypeEnum,
		EnumValues:  []string{"true", "false", "auto"},
	})

	// virtual-keyboard-layouts
	schema.AddOption(&ConfigOption{
		Section:     "",
		Key:         "virtual-keyboard-layouts",
		Default:     "en",
		Description: "A comma separated list of enabled virtual keyboard layouts. The first item in the list is the default layout.",
		Type:        ValueTypeCommaList,
		EnumValues:  []string{"en", "fr", "ru"},
	})

	// virtual-keyboard-layout
	schema.AddOption(&ConfigOption{
		Section:     "",
		Key:         "virtual-keyboard-layout",
		Default:     "en",
		Description: "Virtual keyboard layout / language.",
		Type:        ValueTypeEnum,
		EnumValues:  []string{"en", "fr"},
	})

	// inactive-timeout
	schema.AddOption(&ConfigOption{
		Section:     "",
		Key:         "inactive-timeout",
		Default:     "14400",
		Description: "Timeout (in seconds) after which workspace is closed and returned to the welcome screen if the workspace hasn't been interacted with. Setting inactive-timeout=0 disables this feature.",
		Type:        ValueTypeNumber,
	})

	// menu-timeout
	schema.AddOption(&ConfigOption{
		Section:     "",
		Key:         "menu-timeout",
		Default:     "10",
		Description: "Timeout (in seconds) for expanded menus on the side or bottom of the canvas. Set to zero for no timeout.",
		Type:        ValueTypeNumber,
	})

	// canvus-folder
	schema.AddOption(&ConfigOption{
		Section:     "",
		Key:         "canvus-folder",
		Default:     "",
		Description: "The directory that is opened from the finger menu. DEFAULT: Home directory of the active user in single user mode or folder stored inside content/root in multi-user mode.",
		Type:        ValueTypeFilePath,
	})

	// default-canvas
	schema.AddOption(&ConfigOption{
		Section:     "",
		Key:         "default-canvas",
		Default:     "",
		Description: "URL of the default canvas that can be quickly opened from the welcome screen.",
		Type:        ValueTypeString,
	})

	// ============================================================================
	// [system] SECTION
	// ============================================================================

	// show-volumes
	schema.AddOption(&ConfigOption{
		Section:     "system",
		Key:         "show-volumes",
		Default:     "*",
		Description: "A comma separated list of drive letters, which should always be displayed on the side menu if mounted. A value of '*' means all drive letters. This setting is only supported on Windows. DEFAULT (multi-user-mode-enabled=false): * DEFAULT (multi-user-mode-enabled=true): empty",
		Type:        ValueTypeCommaList,
	})

	// codice-server
	schema.AddOption(&ConfigOption{
		Section:     "system",
		Key:         "codice-server",
		Default:     "",
		Description: "Name of the Codice server where all Codice data like the personal folders and Canvas Codice maps are stored. To use these features, configure MT Canvus multisite server by adding a new section for the server, for instance [server:My server], and then setting codice-server=My server. DEFAULT: empty, which disables all Codice server features",
		Type:        ValueTypeString,
	})

	// enable-personal-folders
	schema.AddOption(&ConfigOption{
		Section:     "system",
		Key:         "enable-personal-folders",
		Default:     "true",
		Description: "Enables personal folder -feature that is triggered with Codice markers. See also server/codice-server setting in this file.",
		Type:        ValueTypeBoolean,
		EnumValues:  []string{"true", "false"},
	})

	// ============================================================================
	// [smtp] SECTION
	// ============================================================================

	// host
	schema.AddOption(&ConfigOption{
		Section:     "smtp",
		Key:         "host",
		Default:     "",
		Description: "SMTP server hostname.",
		Type:        ValueTypeString,
	})

	// port
	schema.AddOption(&ConfigOption{
		Section:     "smtp",
		Key:         "port",
		Default:     "",
		Description: "SMTP server port.",
		Type:        ValueTypeString,
	})

	// username
	schema.AddOption(&ConfigOption{
		Section:     "smtp",
		Key:         "username",
		Default:     "",
		Description: "SMTP username.",
		Type:        ValueTypeString,
	})

	// password
	schema.AddOption(&ConfigOption{
		Section:     "smtp",
		Key:         "password",
		Default:     "",
		Description: "SMTP password.",
		Type:        ValueTypeString,
	})

	// connection-encryption
	schema.AddOption(&ConfigOption{
		Section:     "smtp",
		Key:         "connection-encryption",
		Default:     "auto",
		Description: "Encryption option for email server connection.",
		Type:        ValueTypeEnum,
		EnumValues:  []string{"auto", "none", "SSL", "TLS", "SSL-ignore-errors"},
	})

	// ignore-proxy
	schema.AddOption(&ConfigOption{
		Section:     "smtp",
		Key:         "ignore-proxy",
		Default:     "false",
		Description: "Should system proxy settings be ignored while connecting to SMTP.",
		Type:        ValueTypeBoolean,
		EnumValues:  []string{"true", "false"},
	})

	// ============================================================================
	// [mail] SECTION
	// ============================================================================

	// sender-name
	schema.AddOption(&ConfigOption{
		Section:     "mail",
		Key:         "sender-name",
		Default:     "Canvus",
		Description: "Sender name of the email",
		Type:        ValueTypeString,
	})

	// sender-address
	schema.AddOption(&ConfigOption{
		Section:     "mail",
		Key:         "sender-address",
		Default:     "noreply@multitaction.com",
		Description: "Sender address of the email",
		Type:        ValueTypeString,
	})

	// subject
	schema.AddOption(&ConfigOption{
		Section:     "mail",
		Key:         "subject",
		Default:     "Content from Canvus",
		Description: "Subject for the emails.",
		Type:        ValueTypeString,
	})

	// max-attachment-size
	schema.AddOption(&ConfigOption{
		Section:     "mail",
		Key:         "max-attachment-size",
		Default:     "32",
		Description: "Maximum email attachment size in megabytes",
		Type:        ValueTypeNumber,
	})

	// default-email-address
	schema.AddOption(&ConfigOption{
		Section:     "mail",
		Key:         "default-email-address",
		Default:     "",
		Description: "Default email address when sending email.",
		Type:        ValueTypeString,
	})

	// allow-user-email
	schema.AddOption(&ConfigOption{
		Section:     "mail",
		Key:         "allow-user-email",
		Default:     "true",
		Description: "Allow users to specify another email address when sending email. Note: Disabling this setting and setting default-email-address to an empty value is not a valid configuration. Canvus will refuse to start if you specify it.",
		Type:        ValueTypeBoolean,
		EnumValues:  []string{"true", "false"},
	})

	// ============================================================================
	// [fixed-workspace:<index>] SECTION (Compound)
	// ============================================================================
	// Mark as compound section
	if section := schema.GetSection("fixed-workspace"); section == nil {
		section = &ConfigSection{
			Name:            "fixed-workspace",
			Description:     "Fixed workspaces that cannot be resized or removed. Multiple workspaces can be defined that cannot be changed while the application is running. Each workspace is defined with the format [fixed-workspace:<index>], where the index (starting at 1) defines the left to right positional order of the workspaces.",
			Options:         []*ConfigOption{},
			IsCompound:      true,
			Pattern:         "fixed-workspace",
			CompoundEntries: make(map[string][]*ConfigOption),
		}
		schema.Sections["fixed-workspace"] = section
	}
	fixedWorkspaceSection := schema.GetSection("fixed-workspace")
	fixedWorkspaceSection.IsCompound = true
	fixedWorkspaceSection.Pattern = "fixed-workspace"

	// size (required for each fixed-workspace)
	schema.AddOption(&ConfigOption{
		Section:     "fixed-workspace",
		Key:         "size",
		Default:     "",
		Description: "The size setting must be specified for each workspace. Format: width height (e.g., 3240 1920)",
		Type:        ValueTypeString,
	})

	// name
	schema.AddOption(&ConfigOption{
		Section:     "fixed-workspace",
		Key:         "name",
		Default:     "",
		Description: "If specified, this name appears in the workspace label, otherwise the standard \"Workspace n\" naming is used.",
		Type:        ValueTypeString,
	})

	// view-location
	schema.AddOption(&ConfigOption{
		Section:     "fixed-workspace",
		Key:         "view-location",
		Default:     "0 0",
		Description: "The initial location of the view on the workspace. If 0 0, the default location will be used. NOTE: Values specifying an area outside the view size will be adjusted to ensure the viewport is valid. NOTE: setting either view-location or view-scale will override any other settings that influence the initial position and scale of the view.",
		Type:        ValueTypeString,
	})

	// view-scale
	schema.AddOption(&ConfigOption{
		Section:     "fixed-workspace",
		Key:         "view-scale",
		Default:     "0",
		Description: "The initial scale of the view on the workspace. If zero, the default scale will be used. NOTE: The maximum scale factor value allowed is 10.",
		Type:        ValueTypeNumber,
	})

	// ============================================================================
	// [annotation] SECTION
	// ============================================================================

	// eraser-marker-codes
	schema.AddOption(&ConfigOption{
		Section:     "annotation",
		Key:         "eraser-marker-codes",
		Default:     "612, 614",
		Description: "Codice eraser codes. You can assign a new value or add a comma-separated list of values. By default, Codice cards 612 and 614 are defined as Codice eraser cards.",
		Type:        ValueTypeCommaList,
	})

	// enable-canvas-codices
	schema.AddOption(&ConfigOption{
		Section:     "annotation",
		Key:         "enable-canvas-codices",
		Default:     "false",
		Description: "Allows assignment of Codice marker to be used as a shortcut for canvas selection.",
		Type:        ValueTypeBoolean,
		EnumValues:  []string{"true", "false"},
	})

	// ============================================================================
	// [canvas] SECTION
	// ============================================================================

	// pin-canvas
	schema.AddOption(&ConfigOption{
		Section:     "canvas",
		Key:         "pin-canvas",
		Default:     "false",
		Description: "If true, the canvas is automatically pinned when it is opened.",
		Type:        ValueTypeBoolean,
		EnumValues:  []string{"true", "false"},
	})

	// anchor-viewport
	schema.AddOption(&ConfigOption{
		Section:     "canvas",
		Key:         "anchor-viewport",
		Default:     "false",
		Description: "If true, the canvas is automatically moved to the first anchor position when it is opened. If there are no anchors or it is false, the zoom-viewport setting is used to determine the position.",
		Type:        ValueTypeBoolean,
		EnumValues:  []string{"true", "false"},
	})

	// zoom-viewport
	schema.AddOption(&ConfigOption{
		Section:     "canvas",
		Key:         "zoom-viewport",
		Default:     "false",
		Description: "If true, the canvas is automatically zoomed to ensure all widgets are visible when it is opened. If false, the initial position will be determined by the last saved position within the session. NOTE: If there is no saved position then the viewport will be centralized.",
		Type:        ValueTypeBoolean,
		EnumValues:  []string{"true", "false"},
	})

	// auto-pin-after
	schema.AddOption(&ConfigOption{
		Section:     "canvas",
		Key:         "auto-pin-after",
		Default:     "0",
		Description: "If set, an unpinned canvas is pinned after it has not been scaled or moved for the specified number of seconds. If zero, the canvas is never auto-pinned.",
		Type:        ValueTypeNumber,
	})

	// custom-menu
	schema.AddOption(&ConfigOption{
		Section:     "canvas",
		Key:         "custom-menu",
		Default:     "",
		Description: "Custom menu file in YAML format used to display a user-defined menu structure on the finger menu. NOTE: Backslashes in paths need special handling on Windows computers. Replace each \\ with \\\\ or use forward slashes.",
		Type:        ValueTypeFilePath,
	})

	// ============================================================================
	// [widget] SECTION
	// ============================================================================

	// auto-pin-after
	schema.AddOption(&ConfigOption{
		Section:     "widget",
		Key:         "auto-pin-after",
		Default:     "0",
		Description: "This setting only applies to local canvases. If set, an unpinned widget is pinned after it has not been scaled or moved for the specified number of seconds. If zero, the widget is never auto-pinned.",
		Type:        ValueTypeNumber,
	})

	// ============================================================================
	// [hardware] SECTION
	// ============================================================================

	// third-party-touch
	schema.AddOption(&ConfigOption{
		Section:     "hardware",
		Key:         "third-party-touch",
		Default:     "false",
		Description: "Enable or disable third party touch mode. This mode is meant to be used with devices that do not have a separate pen and/or eraser functionality and only work with touch input. When enabled, users will see extra UI that allows them to use their fingers to annotate and erase.",
		Type:        ValueTypeBoolean,
		EnumValues:  []string{"true", "false"},
	})

	// ============================================================================
	// [touch-mode-selector] SECTION
	// ============================================================================

	// timeout
	schema.AddOption(&ConfigOption{
		Section:     "touch-mode-selector",
		Key:         "timeout",
		Default:     "10",
		Description: "Idle time in seconds before the selector is returned to its home location in finger mode.",
		Type:        ValueTypeNumber,
	})

	// home-x-percent
	schema.AddOption(&ConfigOption{
		Section:     "touch-mode-selector",
		Key:         "home-x-percent",
		Default:     "90",
		Description: "Home position of the selector as a percentage of the window width.",
		Type:        ValueTypeNumber,
	})

	// home-y-percent
	schema.AddOption(&ConfigOption{
		Section:     "touch-mode-selector",
		Key:         "home-y-percent",
		Default:     "90",
		Description: "Home position of the selector as a percentage of the window height.",
		Type:        ValueTypeNumber,
	})

	// home-radius
	schema.AddOption(&ConfigOption{
		Section:     "touch-mode-selector",
		Key:         "home-radius",
		Default:     "500",
		Description: "Distance in pixels from the home position outside which the selector will return to the home position.",
		Type:        ValueTypeNumber,
	})

	// ============================================================================
	// [server:<name>] SECTION (Compound)
	// ============================================================================
	// Mark as compound section
	if section := schema.GetSection("server"); section == nil {
		section = &ConfigSection{
			Name:            "server",
			Description:     "Canvus Connect server list. Configure each server with [server:<name>] section. You can have multiple servers, but each server name must be unique.",
			Options:         []*ConfigOption{},
			IsCompound:      true,
			Pattern:         "server",
			CompoundEntries: make(map[string][]*ConfigOption),
		}
		schema.Sections["server"] = section
	}
	serverSection := schema.GetSection("server")
	serverSection.IsCompound = true
	serverSection.Pattern = "server"

	// server (hostname)
	schema.AddOption(&ConfigOption{
		Section:     "server",
		Key:         "server",
		Default:     "",
		Description: "The Canvus Connect server host to connect to. NOTE: This value must be specified.",
		Type:        ValueTypeString,
	})

	// protocol
	schema.AddOption(&ConfigOption{
		Section:     "server",
		Key:         "protocol",
		Default:     "ssl",
		Description: "Protocol used for communication with the Canvus Connect server. Options are ssl or tcp. The default protocol is determined by the chosen port below. DEFAULT (port 443): ssl DEFAULT (port 80): tcp",
		Type:        ValueTypeEnum,
		EnumValues:  []string{"ssl", "tcp"},
	})

	// port
	schema.AddOption(&ConfigOption{
		Section:     "server",
		Key:         "port",
		Default:     "443",
		Description: "Port used for communication with the Canvus Connect server. The default port is determined by the chosen protocol above. DEFAULT (ssl protocol): 443 DEFAULT (tcp protocol): 80",
		Type:        ValueTypeNumber,
	})

	// rest-api-access
	schema.AddOption(&ConfigOption{
		Section:     "server",
		Key:         "rest-api-access",
		Default:     "none",
		Description: "Control REST API access to this client from this server.",
		Type:        ValueTypeEnum,
		EnumValues:  []string{"none", "ro", "rw"},
	})

	// connection-password
	schema.AddOption(&ConfigOption{
		Section:     "server",
		Key:         "connection-password",
		Default:     "",
		Description: "Set the connection password (if required) for connecting to this server. The password cannot contain semicolon, or quotation marks.",
		Type:        ValueTypeString,
	})

	// ============================================================================
	// [identity] SECTION
	// ============================================================================

	// installation-name
	schema.AddOption(&ConfigOption{
		Section:     "identity",
		Key:         "installation-name",
		Default:     "",
		Description: "Name of the installation of this Canvus client. Is shown in the participant list for this client. If empty, the hostname of the machine is used.",
		Type:        ValueTypeString,
	})

	// ============================================================================
	// [remote-desktop:<name>] SECTION (Compound)
	// ============================================================================
	// Mark as compound section
	if section := schema.GetSection("remote-desktop"); section == nil {
		section = &ConfigSection{
			Name:            "remote-desktop",
			Description:     "Remote desktop connection configuration block. There is no limit on the number of remote-desktop configuration blocks as long as they all have unique names.",
			Options:         []*ConfigOption{},
			IsCompound:      true,
			Pattern:         "remote-desktop",
			CompoundEntries: make(map[string][]*ConfigOption),
		}
		schema.Sections["remote-desktop"] = section
	}
	remoteDesktopSection := schema.GetSection("remote-desktop")
	remoteDesktopSection.IsCompound = true
	remoteDesktopSection.Pattern = "remote-desktop"

	// hostname
	schema.AddOption(&ConfigOption{
		Section:     "remote-desktop",
		Key:         "hostname",
		Default:     "",
		Description: "Hostname or IP address of the RDP server with an optional port number. If this parameter is empty, the whole configuration block is disabled. Examples: localhost, rdp.example.com, 10.0.0.1:1234",
		Type:        ValueTypeString,
	})

	// username
	schema.AddOption(&ConfigOption{
		Section:     "remote-desktop",
		Key:         "username",
		Default:     "",
		Description: "Username to use to log in to the server.",
		Type:        ValueTypeString,
	})

	// password
	schema.AddOption(&ConfigOption{
		Section:     "remote-desktop",
		Key:         "password",
		Default:     "",
		Description: "Password for the given username. Leave empty if password is not used or if you wish to enter the password manually when connecting.",
		Type:        ValueTypeString,
	})

	// domain
	schema.AddOption(&ConfigOption{
		Section:     "remote-desktop",
		Key:         "domain",
		Default:     "",
		Description: "Domain to use for authentication.",
		Type:        ValueTypeString,
	})

	// max-connections
	schema.AddOption(&ConfigOption{
		Section:     "remote-desktop",
		Key:         "max-connections",
		Default:     "0",
		Description: "Maximum number of concurrent connections to this server. DEFAULT: 0 (unlimited)",
		Type:        ValueTypeNumber,
	})

	// show-in-server-list
	schema.AddOption(&ConfigOption{
		Section:     "remote-desktop",
		Key:         "show-in-server-list",
		Default:     "true",
		Description: "Show this server in the Remote Desktop list. If this is set to false, you can't directly connect to the server, instead this connection is just used to open files. Either show-in-server-list must be true or file-extensions must not be empty.",
		Type:        ValueTypeBoolean,
		EnumValues:  []string{"true", "false"},
	})

	// file-extensions
	schema.AddOption(&ConfigOption{
		Section:     "remote-desktop",
		Key:         "file-extensions",
		Default:     "",
		Description: "Comma-separated list of file extensions of files that are opened through this server. Either show-in-server-list must be true or file-extensions must not be empty. Example: doc,docx,ppt,pptx,xls,xlsx",
		Type:        ValueTypeCommaList,
	})

	// gateway
	schema.AddOption(&ConfigOption{
		Section:     "remote-desktop",
		Key:         "gateway",
		Default:     "",
		Description: "RDP gateway hostname with optional port number. Leave empty if gateway is not used.",
		Type:        ValueTypeString,
	})

	// gateway-username
	schema.AddOption(&ConfigOption{
		Section:     "remote-desktop",
		Key:         "gateway-username",
		Default:     "",
		Description: "Gateway username.",
		Type:        ValueTypeString,
	})

	// gateway-password
	schema.AddOption(&ConfigOption{
		Section:     "remote-desktop",
		Key:         "gateway-password",
		Default:     "",
		Description: "Gateway password.",
		Type:        ValueTypeString,
	})

	// gateway-domain
	schema.AddOption(&ConfigOption{
		Section:     "remote-desktop",
		Key:         "gateway-domain",
		Default:     "",
		Description: "Gateway Domain.",
		Type:        ValueTypeString,
	})

	// desktop-size
	schema.AddOption(&ConfigOption{
		Section:     "remote-desktop",
		Key:         "desktop-size",
		Default:     "1920 1080",
		Description: "Initial RDP desktop resolution for new RDP widgets.",
		Type:        ValueTypeString,
	})

	// app
	schema.AddOption(&ConfigOption{
		Section:     "remote-desktop",
		Key:         "app",
		Default:     "",
		Description: "RemoteApp to launch on connect. This setting is used only when connecting using the Remote Desktop list. NOTE: Backslashes in paths need special handling on Windows computers. Replace each \\ with \\\\. Example: %WINDIR%\\\\System32\\\\notepad.exe Example: ||notepad (using configured RemoteApp alias) DEFAULT: empty (show the whole desktop)",
		Type:        ValueTypeFilePath,
	})

	// app-args
	schema.AddOption(&ConfigOption{
		Section:     "remote-desktop",
		Key:         "app-args",
		Default:     "",
		Description: "Command line arguments to RemoteApp, used when \"app\" setting is used.",
		Type:        ValueTypeString,
	})

	// extra-args
	schema.AddOption(&ConfigOption{
		Section:     "remote-desktop",
		Key:         "extra-args",
		Default:     "",
		Description: "Extra FreeRDP command line arguments. Example: /drive:home,/home/multi",
		Type:        ValueTypeString,
	})

	// ============================================================================
	// [cache] SECTION
	// ============================================================================

	// max-age
	schema.AddOption(&ConfigOption{
		Section:     "cache",
		Key:         "max-age",
		Default:     "2592000",
		Description: "Maximum time in seconds that assets are allowed to be cached. DEFAULT: 2592000 (30 days). To disable caching set max-age=0.",
		Type:        ValueTypeNumber,
	})

	// ============================================================================
	// [browser] SECTION
	// ============================================================================

	// start-page
	schema.AddOption(&ConfigOption{
		Section:     "browser",
		Key:         "start-page",
		Default:     "https://www.google.com",
		Description: "Browser start page that is opened when a new web browser widget is created.",
		Type:        ValueTypeString,
	})

	// search-url
	schema.AddOption(&ConfigOption{
		Section:     "browser",
		Key:         "search-url",
		Default:     "https://www.google.com/search?q=QUERY",
		Description: "Search URL that is opened when a search term is entered in the web browser address bar. A word QUERY written in uppercase letters is replaced with the search terms.",
		Type:        ValueTypeString,
	})

	// upload-sites
	schema.AddOption(&ConfigOption{
		Section:     "browser",
		Key:         "upload-sites",
		Default:     "*",
		Description: "Specify which internet sites allow content to be uploaded by dragging onto a browser widget. if upload-sites=* or missing, allow all uploads. if upload-sites is empty, prevent all uploads. if upload-sites are specified a widget can only be dropped on the browser if its current URL contains one of the specified names, e.g. upload-sites=drive.google.com,docs.google.com,mail.google.com",
		Type:        ValueTypeCommaList,
	})

	// persistent-sessions
	schema.AddOption(&ConfigOption{
		Section:     "browser",
		Key:         "persistent-sessions",
		Default:     "true",
		Description: "Specify if browser sessions should be saved on disk so that cookies, browsing history and other session data persists between application restarts. Set to false to use in-memory sessions that are cleared when application is restarted or user logs out.",
		Type:        ValueTypeBoolean,
		EnumValues:  []string{"true", "false"},
	})

	// mouse-emulation
	schema.AddOption(&ConfigOption{
		Section:     "browser",
		Key:         "mouse-emulation",
		Default:     "false",
		Description: "Browsers can handle touches by either passing all touches to the browser or converting the first touch to mouse clicks and drags. Different web apps work better with the two alternative behaviors. The user can pick the best behavior on the browser. This setting determines the default state for new browsers. If true, then, by default, touches are converted to mouse clicks.",
		Type:        ValueTypeBoolean,
		EnumValues:  []string{"true", "false"},
	})

	return schema
}

