//go:build windows

package sys

import "syscall"

func GetSystemMetrics() (int, int) {
	ret1, _, _ := syscall.NewLazyDLL(`User32.dll`).NewProc(`GetSystemMetrics`).Call(uintptr(0))
	ret2, _, _ := syscall.NewLazyDLL(`User32.dll`).NewProc(`GetSystemMetrics`).Call(uintptr(1))
	return int(ret1), int(ret2)
}
