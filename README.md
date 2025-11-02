# GoMute - Microphone Mute Toggle

A lightweight Windows system tray application to toggle microphone mute with a keyboard shortcut.

## Features

- Toggle microphone mute with **F8** key
- System tray icon shows current mic status
- Lightweight (~2-5 MB RAM)
- Runs silently in background
- Auto-start on Windows login (if you want)

## Use Cases

- Quick mute/unmute during video calls (Teams, Zoom, Discord, etc.)
- Gaming communication control
- Privacy protection when working from home
- Alternative to hardware mute buttons

## Build

```bash
# First time setup (to add executable icon)
go install github.com/tc-hib/go-winres@latest
go-winres simply --icon icons/app.ico

# Build
go build -ldflags -H=windowsgui -o go-mute.exe
```

## Installation

1. Build the executable or use the pre-built `go-mute.exe`
2. (Optional) Copy to startup folder for auto-start:
   ```
   C:\Users\<YourUsername>\AppData\Roaming\Microsoft\Windows\Start Menu\Programs\Startup
   ```

## Usage

- **Run**: Double-click `go-mute.exe`
- **Toggle Mute**: Press **F8**
- **Exit**: Right-click system tray icon → Quit

## Customization

### Change Keybind

To use a different key instead of F8:

1. Open `main.go`
2. Find the line `VK_F8 = 0x77`
3. Replace `0x77` with your desired key code:
   - F1-F12: `0x70` to `0x7B`
   - A-Z: `0x41` to `0x5A`
   - 0-9: `0x30` to `0x39`
   - [Full list of virtual key codes](https://learn.microsoft.com/en-us/windows/win32/inputdev/virtual-key-codes)
4. Rebuild: `go build -ldflags -H=windowsgui -o go-mute.exe`

## Requirements

- Windows OS
- Go 1.16+ (for building)

## Supported Devices

Controls the default microphone via Windows Core Audio API. Works with any microphone recognized by Windows:
- USB microphones
- Headset microphones
- Built-in laptop microphones
- Audio interface inputs
- Bluetooth microphones
- Any device listed as a recording device in Windows Sound settings

**Note:** Does not require hardware mute support. Controls system-level mute independently of hardware buttons.
