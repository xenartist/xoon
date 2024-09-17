package utils

import "runtime"

var (
	isWindows = runtime.GOOS == "windows"
	fileName  string
	minerDir  string
)
