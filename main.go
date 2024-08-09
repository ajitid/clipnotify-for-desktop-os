// original sources: https://abdus.dev/projects/dumpclip/ and https://abdus.dev/posts/monitor-clipboard/
package main

import (
	"os"
	"syscall"

	"github.com/gonutz/w32/v2"
)

const (
	WM_CLIPBOARDUPDATE = 0x031D
)

var exitChan = make(chan struct{})

func main() {
	// Create a window to receive messages
	className, _ := syscall.UTF16PtrFromString("ClipboardMonitorClass")
	wc := w32.WNDCLASSEX{
		WndProc:   syscall.NewCallback(wndProc),
		ClassName: className,
	}
	w32.RegisterClassEx(&wc)

	hwnd := w32.CreateWindowEx(
		0,
		className,
		nil,
		0,
		0, 0, 0, 0,
		0,
		0,
		0,
		nil,
	)

	// Add the window as a clipboard viewer
	w32.AddClipboardFormatListener(hwnd)

	// Start a goroutine to handle the exit signal
	go func() {
		<-exitChan
		w32.PostQuitMessage(0)
		os.Exit(0)
	}()

	// Message loop
	var msg w32.MSG
	for w32.GetMessage(&msg, 0, 0, 0) != 0 {
		w32.TranslateMessage(&msg)
		w32.DispatchMessage(&msg)
	}
}

func wndProc(hwnd w32.HWND, msg uint32, wParam, lParam uintptr) uintptr {
	switch msg {
	case WM_CLIPBOARDUPDATE:
		if w32.OpenClipboard(hwnd) {
			defer w32.CloseClipboard()

			hClipData := w32.GetClipboardData(w32.CF_UNICODETEXT)
			if hClipData != 0 {
				pszText := w32.GlobalLock(w32.HGLOBAL(hClipData))
				if pszText != nil {
					defer w32.GlobalUnlock(w32.HGLOBAL(hClipData))

					// clipboardText := w32.UTF16PtrToString((*uint16)(unsafe.Pointer(pszText)))
					// fmt.Printf("New clipboard content: %s\n", clipboardText)

					// Signal to exit the program
					close(exitChan)
				}
			}
		}
	}
	return w32.DefWindowProc(hwnd, msg, wParam, lParam)
}
