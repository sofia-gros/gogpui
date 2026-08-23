//go:build windows

package gogpui

import (
	"syscall"
	"unsafe"
)

var (
	user32            = syscall.NewLazyDLL("user32.dll")
	procFindWindowW   = user32.NewProc("FindWindowW")
	procIsIconic      = user32.NewProc("IsIconic")
	procGetClientRect = user32.NewProc("GetClientRect")
)

type winRect struct {
	left, top, right, bottom int32
}

// checkWindowMinimized は Win32 API を直接呼び出してウィンドウが最小化されているか確認する。
// 最小化されている場合は (true, 0, 0) を返し、通常表示の場合は (false, width, height) を返す。
func checkWindowMinimized(title string) (bool, int, int) {
	titleUTF16, err := syscall.UTF16PtrFromString(title)
	if err != nil {
		return false, 0, 0
	}
	classUTF16, err := syscall.UTF16PtrFromString("GoGPUWindow")
	if err != nil {
		return false, 0, 0
	}

	hwnd, _, _ := procFindWindowW.Call(
		uintptr(unsafe.Pointer(classUTF16)),
		uintptr(unsafe.Pointer(titleUTF16)),
	)
	if hwnd == 0 {
		return false, 0, 0
	}

	isIconic, _, _ := procIsIconic.Call(hwnd)
	if isIconic != 0 {
		return true, 0, 0
	}

	var r winRect
	procGetClientRect.Call(hwnd, uintptr(unsafe.Pointer(&r)))
	w := int(r.right - r.left)
	h := int(r.bottom - r.top)
	if w <= 0 || h <= 0 {
		return true, 0, 0
	}

	return false, w, h
}
