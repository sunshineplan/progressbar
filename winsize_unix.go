//go:build unix

package progressbar

import "golang.org/x/sys/unix"

func GetWinsize() int {
	ws, err := unix.IoctlGetWinsize(0, unix.TIOCGWINSZ)
	if err != nil {
		return 0
	}
	return int(ws.Col)
}
