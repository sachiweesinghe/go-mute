package main

import (
	"runtime"
	"unsafe"

	"github.com/getlantern/systray"
	"github.com/go-ole/go-ole"
	"github.com/moutend/go-wca/pkg/wca"
	"golang.org/x/sys/windows"
)

const (
	MOD_CONTROL = 0x0002
	MOD_SHIFT   = 0x0004
	VK_F8       = 0x77
	WM_HOTKEY   = 0x0312
)

var (
	user32           = windows.NewLazySystemDLL("user32.dll")
	registerHotKey   = user32.NewProc("RegisterHotKey")
	unregisterHotKey = user32.NewProc("UnregisterHotKey")
	getMessage       = user32.NewProc("GetMessageW")
	isMuted          = false
)

func updateIcon() {
	if isMuted {
		systray.SetIcon(iconMicOff)
		systray.SetTooltip("Microphone: MUTED (Press F8)")
	} else {
		systray.SetIcon(iconMicOn)
		systray.SetTooltip("Microphone: ON (Press F8)")
	}
}

func toggleMuteWithCOM() {
	// Get device enumerator
	var deviceEnumerator *wca.IMMDeviceEnumerator
	if err := wca.CoCreateInstance(
		wca.CLSID_MMDeviceEnumerator,
		0,
		wca.CLSCTX_ALL,
		wca.IID_IMMDeviceEnumerator,
		&deviceEnumerator,
	); err != nil {
		return
	}
	defer deviceEnumerator.Release()

	// Get default microphone (capture device)
	var device *wca.IMMDevice
	if err := deviceEnumerator.GetDefaultAudioEndpoint(wca.ECapture, wca.EConsole, &device); err != nil {
		return
	}
	defer device.Release()

	// Activate the audio endpoint volume interface
	var endpointVolume *wca.IAudioEndpointVolume
	if err := device.Activate(
		wca.IID_IAudioEndpointVolume,
		wca.CLSCTX_ALL,
		nil,
		&endpointVolume,
	); err != nil {
		return
	}
	defer endpointVolume.Release()

	// Get current mute state
	var currentMute bool
	if err := endpointVolume.GetMute(&currentMute); err != nil {
		return
	}

	// Toggle mute state
	newMute := !currentMute
	eventContext := ole.GUID{}
	endpointVolume.SetMute(newMute, &eventContext)

	// Update global state and icon
	isMuted = newMute
	updateIcon()
}

func onReady() {
	systray.SetTitle("Mic Mute")
	updateIcon()

	mQuit := systray.AddMenuItem("Quit", "Exit the application")

	go func() {
		<-mQuit.ClickedCh
		systray.Quit()
	}()
}

func onExit() {
	// Cleanup
}

func startHotkeyListener() {
	// Lock the OS thread to ensure COM operations happen on the same thread
	runtime.LockOSThread()
	defer runtime.UnlockOSThread()

	// Initialize COM for the main thread
	if err := ole.CoInitializeEx(0, ole.COINIT_APARTMENTTHREADED); err != nil {
		return
	}
	defer ole.CoUninitialize()

	// Register hotkey: F8
	hotKeyID := 1
	ret, _, _ := registerHotKey.Call(
		0,                 // NULL window handle
		uintptr(hotKeyID), // hotkey ID
		0,                 // no modifiers (just F8)
		VK_F8,             // virtual key code
	)
	if ret == 0 {
		return
	}
	defer unregisterHotKey.Call(0, uintptr(hotKeyID))

	// Message loop to listen for hotkey events
	var msg struct {
		hwnd    uintptr
		message uint32
		wParam  uintptr
		lParam  uintptr
		time    uint32
		pt      struct{ x, y int32 }
	}

	for {
		ret, _, _ := getMessage.Call(
			uintptr(unsafe.Pointer(&msg)),
			0,
			0,
			0,
		)
		if ret == 0 {
			break
		}

		if msg.message == WM_HOTKEY {
			toggleMuteWithCOM()
		}
	}
}

func main() {
	go startHotkeyListener()
	systray.Run(onReady, onExit)
}
