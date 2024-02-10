//go:build !windows

package main

import (
	_ "embed"
	"errors"
	"fmt"
	"os/exec"
	"runtime"
)

//go:embed bin/wget_linux_any.txt
var linuxHelp string

//go:embed bin/wget_macos_any.txt
var macosHelp string

func InitializeBinary() error {
	cmd := exec.Command("wget", "--version")
	_, err := cmd.CombinedOutput()
	if err != nil {
		fmt.Println("wget is not installed! follow these instructions to install on most operating systems:")
		if //goland:noinspection GoBoolExpressions
		runtime.GOOS == "darwin" {
			fmt.Println(macosHelp)
		} else {
			fmt.Println(linuxHelp)
		}

		return errors.New("wget not installed")
	}

	return nil
}
