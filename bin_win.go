//go:build windows

package main

import (
	_ "embed"
	"os"
)

//go:embed bin/wget_win_64.exe
var binary []byte

func InitializeBinary() error {
	return os.WriteFile("wget.exe", binary, os.ModePerm)
}
