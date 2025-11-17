# Deployment Guide

This document provides instructions for deploying Canvus PowerToys in production environments.

## System Requirements

### Windows (Primary Platform)
- **OS**: Windows 10 or later
- **Architecture**: x64 (amd64)
- **Privileges**:
  - Standard user privileges for user config writes
  - Administrator privileges for system config writes
- **Disk Space**: ~50 MB (application + logs)
- **Network**: Required only for WebUI feature (LAN access)

### Linux/Ubuntu (Future Support)
- **OS**: Ubuntu 20.04 or later
- **Architecture**: x64 (amd64)
- **Dependencies**: GTK3 libraries (for Fyne framework)
- **Privileges**: Standard user privileges (sudo for system config)
- **Disk Space**: ~50 MB (application + logs)

## Deployment Methods

### Method 1: Standalone Executable (Recommended)

#### Windows Deployment
1. **Download the executable**
   - Download `canvus-powertoys.X.X.X.exe` from releases
   - Or build from source (see [BUILD.md](BUILD.md))

2. **Place executable**
   - Create a deployment directory (e.g., `C:\Program Files\CanvusPowerToys\`)
   - Copy executable to this directory
   - Optional: Create desktop shortcut

3. **First Run**
   - Run the executable
   - Application will auto-detect configuration files in standard locations
   - No installation or setup required

4. **Verify Installation**
   - Check that application launches successfully
   - Verify system tray icon appears
   - Confirm configuration files are detected

#### Linux Deployment
1. **Download the executable**
   - Download `canvus-powertoys-linux` from releases
   - Or build from source (see [BUILD.md](BUILD.md))

2. **Place executable**
   - Create deployment directory (e.g., `/opt/canvus-powertoys/`)
   - Copy executable to this directory
   - Make executable: `chmod +x canvus-powertoys-linux`

3. **Create desktop entry** (optional)
   ```bash
   # Create .desktop file in ~/.local/share/applications/
   cat > ~/.local/share/applications/canvus-powertoys.desktop << EOF
   [Desktop Entry]
   Name=Canvus PowerToys
   Exec=/opt/canvus-powertoys/canvus-powertoys-linux
   Icon=/opt/canvus-powertoys/icon.png
   Type=Application
   Categories=Utility;
   EOF
   ```

4. **First Run**
   - Run: `./canvus-powertoys-linux`
   - Application will auto-detect configuration files

### Method 2: Portable Deployment

For portable deployments (USB drives, network shares):

1. **Create portable directory structure**
   ```
   CanvusPowerToys/
   ├── canvus-powertoys.exe (or canvus-powertoys-linux)
   ├── config/              # Optional: pre-configured files
   │   ├── mt-canvus.ini
   │   └── screen.xml
   └── README.txt           # Usage instructions
   ```

2. **Configuration**
   - Application will use standard config locations
   - For portable configs, users can manually specify file paths in the UI

## Configuration

### Standard Configuration File Locations

#### Windows
- **User Config**: `%APPDATA%\MultiTaction\canvus\mt-canvus.ini`
  - Default: `C:\Users\<username>\AppData\Roaming\MultiTaction\canvus\mt-canvus.ini`
- **System Config**: `%ProgramData%\MultiTaction\canvus\mt-canvus.ini`
  - Default: `C:\ProgramData\MultiTaction\canvus\mt-canvus.ini`
- **Logs**: `%LOCALAPPDATA%\MultiTaction\Canvus\logs\`
  - Default: `C:\Users\<username>\AppData\Local\MultiTaction\Canvus\logs\`

#### Linux/Ubuntu
- **User Config**: `~/.config/MultiTaction/canvus/mt-canvus.ini`
- **System Config**: `/etc/MultiTaction/canvus/mt-canvus.ini`
- **Logs**: `~/.local/share/MultiTaction/Canvus/logs/`

### Auto-Detection

The application automatically detects configuration files in these locations:
1. User config location
2. System config location
3. Current directory (for portable deployments)

### Manual Configuration

Users can manually specify file paths in the UI:
- **Config Editor**: "Open File" button to browse for mt-canvus.ini
- **Screen.xml Creator**: File path can be specified before generation

## WebUI Deployment

### Enabling WebUI

1. **Open WebUI Settings Tab**
   - Launch application
   - Navigate to "WebUI Settings" tab

2. **Configure Canvus Server**
   - Enable WebUI server (toggle switch)
   - Enter Canvus Server URL (e.g., `https://canvus.example.com`)
   - Enter Private-Token (stored securely, only last 6 digits displayed)

3. **Network Configuration**
   - Default port: 8080
   - Binding: 0.0.0.0 (accessible from LAN)
   - No password/authentication (LAN-only, trusted network)

4. **Access WebUI**
   - From LAN devices: `http://<server-ip>:8080`
   - Mobile devices: Same URL (mobile-responsive interface)

### WebUI Security Considerations

- **LAN-Only Access**: WebUI binds to 0.0.0.0 (all interfaces) but is intended for LAN use
- **No Authentication**: WebUI has no password protection (trusted network assumption)
- **Token Security**: Canvus Server Private-Token is stored encrypted
- **Firewall**: Ensure port 8080 is open on LAN (block from WAN)

### Firewall Configuration

#### Windows Firewall
```powershell
# Allow inbound connections on port 8080 (LAN only)
New-NetFirewallRule -DisplayName "Canvus PowerToys WebUI" -Direction Inbound -LocalPort 8080 -Protocol TCP -Action Allow -Profile Private
```

#### Linux Firewall (ufw)
```bash
# Allow port 8080 from LAN
sudo ufw allow from 192.168.0.0/16 to any port 8080
```

## First-Time Setup

### Initial Configuration

1. **Launch Application**
   - Run executable
   - Application window appears
   - System tray icon appears

2. **Verify File Detection**
   - Check that mt-canvus.ini is detected (shown in Config Editor tab)
   - Check that screen.xml is detected (if exists, shown in Screen.xml Creator tab)

3. **Configure Features** (as needed)
   - **Screen.xml Creator**: Create screen.xml for multi-display setup
   - **Config Editor**: Modify Canvus settings
   - **CSS Options Manager**: Enable kiosk mode or CSS customizations
   - **Custom Menu Designer**: Create custom menus
   - **WebUI Settings**: Enable remote management

4. **Test Configuration**
   - Make a test change in Config Editor
   - Verify backup is created
   - Verify save succeeds
   - Restore from backup to verify backup system

### WebUI First-Time Setup

1. **Enable WebUI**
   - Open WebUI Settings tab
   - Enable WebUI server

2. **Configure Canvus Server**
   - Enter Canvus Server URL
   - Enter Private-Token
   - Verify connection (check connection status indicator)

3. **Test Access**
   - From another device on LAN, navigate to `http://<server-ip>:8080`
   - Verify main page loads
   - Verify canvas tracking works (canvas header shows installation/canvas name)

4. **Test Features**
   - Test Pages Management
   - Test Macros operations
   - Test Remote Content Upload
   - Verify real-time updates (SSE)

## Backup System

### Automatic Backups

- Backups are created automatically before all file saves
- Backup location: Same directory as source file
- Backup naming: `filename.YYYYMMDD-HHMMSS.backup`
- Backup rotation: Keeps last N backups (configurable)

### Backup Restoration

1. **Manual Restoration**
   - Locate backup file in same directory as source file
   - Rename backup file to original filename (remove `.backup` extension)
   - Or copy backup file to desired location

2. **Automatic Restoration** (if implemented)
   - Use backup manager in application (if available)

### Backup Management

- **Location**: Same directory as source file
- **Retention**: Last N backups (default: 10)
- **Cleanup**: Old backups are automatically removed when limit exceeded

## Logging

### Log Location

#### Windows
- `%LOCALAPPDATA%\MultiTaction\Canvus\logs\`
- Default: `C:\Users\<username>\AppData\Local\MultiTaction\Canvus\logs\`

#### Linux
- `~/.local/share/MultiTaction/Canvus/logs/`

### Log Files

- **Development Mode**: Comprehensive logging to console and file
- **Production Mode**: Simple error dialogs, detailed logs to file
- **Log Rotation**: Automatic (keeps last N log files)
- **Log Format**: Timestamp, level, message

### Troubleshooting with Logs

1. **Check Log Files**
   - Navigate to log directory
   - Open latest log file
   - Look for ERROR or WARN entries

2. **Common Issues**
   - File permission errors: Check write permissions
   - Configuration file not found: Verify file paths
   - WebUI connection issues: Check network and Canvus Server URL

## Updating

### Update Process

1. **Backup Current Installation**
   - Backup current executable
   - Backup configuration files (if custom)

2. **Download New Version**
   - Download latest executable from releases
   - Or build from source (see [BUILD.md](BUILD.md))

3. **Replace Executable**
   - Stop application (close window, exit from tray)
   - Replace old executable with new one
   - Launch new version

4. **Verify Update**
   - Check version in application (if available)
   - Test core features
   - Verify configuration files are still detected

### Configuration Migration

- Configuration files are forward-compatible
- Old configurations are automatically detected
- No migration required (application reads existing files)

## Uninstallation

### Windows

1. **Stop Application**
   - Close application window
   - Exit from system tray (right-click → Exit)

2. **Remove Files**
   - Delete executable
   - Optional: Delete configuration files in `%APPDATA%\MultiTaction\canvus\`
   - Optional: Delete logs in `%LOCALAPPDATA%\MultiTaction\Canvus\logs\`

3. **Remove Shortcuts**
   - Delete desktop shortcut (if created)
   - Remove from Start Menu (if added)

### Linux

1. **Stop Application**
   - Close application window
   - Exit from system tray

2. **Remove Files**
   - Delete executable
   - Optional: Delete configuration files in `~/.config/MultiTaction/canvus/`
   - Optional: Delete logs in `~/.local/share/MultiTaction/Canvus/logs/`
   - Optional: Delete desktop entry in `~/.local/share/applications/`

## Troubleshooting

### Application Won't Start

**Symptoms**: Executable doesn't launch or crashes immediately

**Solutions**:
1. Check system requirements (Windows 10+, x64 architecture)
2. Verify executable permissions (Linux: `chmod +x`)
3. Check logs in log directory
4. Run from command line to see error messages
5. Verify dependencies (GTK3 for Linux)

### Configuration Files Not Detected

**Symptoms**: Application shows "File not found" or empty configuration

**Solutions**:
1. Verify configuration files exist in standard locations
2. Check file permissions (read access required)
3. Use "Open File" button to manually specify path
4. Verify file format (valid INI/XML/YAML)

### WebUI Not Accessible

**Symptoms**: Cannot access WebUI from LAN devices

**Solutions**:
1. Verify WebUI is enabled in settings
2. Check firewall rules (port 8080)
3. Verify server IP address
4. Check network connectivity (ping server)
5. Verify Canvus Server URL and Private-Token
6. Check logs for connection errors

### Backup Issues

**Symptoms**: Backups not created or restoration fails

**Solutions**:
1. Check write permissions in config directory
2. Verify disk space available
3. Check backup directory exists
4. Review logs for backup errors

### Performance Issues

**Symptoms**: Application is slow or unresponsive

**Solutions**:
1. Check system resources (CPU, memory)
2. Review log files for errors
3. Close other applications
4. Verify configuration files are not corrupted
5. Check network connectivity (if using WebUI)

## Production Checklist

Before deploying to production:

- [ ] Application tested on target OS version
- [ ] Configuration files verified in standard locations
- [ ] Backup system tested and verified
- [ ] WebUI tested (if enabled) from LAN devices
- [ ] Firewall rules configured (if WebUI enabled)
- [ ] Logging verified and log rotation working
- [ ] System tray integration tested
- [ ] All features tested and working
- [ ] Documentation provided to users
- [ ] Support contact information available

## Support

For deployment issues:
1. Check logs in log directory
2. Review [BUILD.md](BUILD.md) for build issues
3. Review [README.md](README.md) for usage instructions
4. Check [TASKS.md](dev/TASKS.md) for known issues

## Additional Resources

- [README.md](README.md) - Project overview and features
- [BUILD.md](BUILD.md) - Build instructions
- [PRD.md](dev/PRD.md) - Product requirements and specifications
- [TECH_STACK.md](dev/TECH_STACK.md) - Technology stack and architecture

