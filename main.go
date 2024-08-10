// original sources: https://abdus.dev/projects/dumpclip/ and https://abdus.dev/posts/monitor-clipboard/
package main

import (
	"fmt"
	"log"
	"syscall"
	"unsafe"

	"github.com/gonutz/w32/v2"
)

const (
	WM_CLIPBOARDUPDATE = 0x031D
)

func main() {
	// Create a window to receive messages
	className, err := syscall.UTF16PtrFromString("ClipboardMonitorClass")
	if err != nil {
		log.Fatalf("Failed to create class name: %v", err)
	}

	wc := w32.WNDCLASSEX{
		WndProc:   syscall.NewCallback(wndProc),
		ClassName: className,
	}
	if atom := w32.RegisterClassEx(&wc); atom == 0 {
		log.Fatal("Failed to register window class")
	}

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
	if hwnd == 0 {
		log.Fatal("Failed to create window")
	}

	// Add the window as a clipboard viewer
	if !w32.AddClipboardFormatListener(hwnd) {
		log.Fatal("Failed to add clipboard format listener")
	}

	fmt.Println("Monitoring clipboard for changes. Press Ctrl+C to exit.")

	// Message loop
	var msg w32.MSG
	for {
		result := w32.GetMessage(&msg, 0, 0, 0)
		if result == 0 {
			break // WM_QUIT
		}
		if result == -1 {
			log.Println("Error getting message")
			continue
		}
		w32.TranslateMessage(&msg)
		w32.DispatchMessage(&msg)
	}
}

func wndProc(hwnd w32.HWND, msg uint32, wParam, lParam uintptr) uintptr {
	switch msg {
	case WM_CLIPBOARDUPDATE:
		go func() {
			if !w32.OpenClipboard(hwnd) {
				log.Println("Failed to open clipboard")
				return
			}
			defer w32.CloseClipboard()

			hClipData := w32.GetClipboardData(w32.CF_UNICODETEXT)
			if hClipData == 0 {
				log.Println("Failed to get clipboard data")
				return
			}

			pszText := w32.GlobalLock(w32.HGLOBAL(hClipData))
			if pszText == nil {
				log.Println("Failed to lock global memory")
				return
			}
			defer w32.GlobalUnlock(w32.HGLOBAL(hClipData))

			clipboardText := w32.UTF16PtrToString((*uint16)(unsafe.Pointer(pszText)))
			fmt.Printf("New clipboard content: %s\n", clipboardText)
		}()
	}
	return w32.DefWindowProc(hwnd, msg, wParam, lParam)
}
