package progressbar

import "golang.org/x/sys/windows"

func getWinsize() int {
	var info windows.ConsoleScreenBufferInfo
	if err := windows.GetConsoleScreenBufferInfo(windows.Handle(windows.Stdout), &info); err != nil {
		return 0
	}
	return int(info.Size.X)
}
