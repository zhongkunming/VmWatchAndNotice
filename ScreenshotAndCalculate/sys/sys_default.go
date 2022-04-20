//go:build !windows

package sys

import (
	"github.com/go-vgo/robotgo"
)

func GetSystemMetrics() (int, int) {

	event := robotgo.AddEvent("mleft")
	if event {
		x, y := robotgo.GetMousePos()
		return x + 1, y + 1
	}
	panic("")
}
