//go:build windows

package main

import (
	_ "embed"
	"os"
)

//go:embed bin/wget_win_64.exe
var binary []byte

func InitializeBinary() error {
	os.Chdir(os.UserConfigDir())
	os.Mkdir("RobotDocs", os.ModeDir)
	return os.WriteFile("wget.exe", binary, os.ModePerm)
}
