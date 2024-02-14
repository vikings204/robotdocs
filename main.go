package main

import (
	"fmt"
	"github.com/gofiber/fiber/v2"
	"io/fs"
	"os"
	"path/filepath"
	"strings"
	"time"
)

func main() {
	dir, err := os.UserConfigDir()
	if err != nil {
		panic(err)
	}
	err = os.Chdir(dir)
	if err != nil {
		panic(err)
	}
	err = os.Chdir("RobotDocs")
	if err != nil {
		err = os.Mkdir("RobotDocs", os.ModeDir)
		err = os.Chdir("RobotDocs")
		if err != nil {
			panic(err)
		}
	}

	//panic(fixExternalRefs("htmltest\\baller\\og.html", "htmltest\\baller\\_external"))
	fmt.Println(fixExternalRefs("og.html", "htmltest\\baller", "htmltest\\baller\\_external"))
	//fmt.Println(fixExternalRefs("app.gitbook.com_public_app_public-PO7REWUB.min.js-v-10.9.593-82f7df415a98c42348302ad74af126b6adfbf759-7828934982.js", "htmltest\\baller\\_external", "htmltest\\baller\\_external"))

	for i := 0; i < 1; i++ {
		fmt.Println(filepath.WalkDir("htmltest\\baller\\_external", func(path string, di fs.DirEntry, err error) error {
			if strings.Contains(path, ".js") {
				fmt.Println(path)
				return fixExternalRefs(path, "", "htmltest\\baller\\_external")
			} else {
				return nil
			}
		}))
	}
	fmt.Println(filepath.WalkDir("htmltest\\baller\\_external", func(path string, di fs.DirEntry, err error) error {
		if strings.Contains(path, ".js") {
			fmt.Println(path)
			bytes, er := os.ReadFile(path)
			if er != nil {
				return er
			}
			str := string(bytes)

			str = strings.Replace(str, "https://app.gitbook.com/public/app/chunks/", "http://127.0.0.1:3000/baller/_external/app.gitbook.com_public_app_chunks_", -1)

			return os.WriteFile(path, []byte(str), os.ModePerm)
		} else {
			return nil
		}
	}))

	err = InitializeBinary()
	if err != nil {
		panic(err)
	} else {
		fmt.Println("binary error:", err)
	}

	app := fiber.New(fiber.Config{
		Prefork:           false,
		ServerHeader:      "RobotDocs",
		GETOnly:           true,
		DisableKeepalive:  true,
		AppName:           "RobotDocs",
		StreamRequestBody: true,
		ReduceMemoryUsage: true,
		Network:           "tcp4",
		EnablePrintRoutes: true,
	})

	app.Get("/stop", func(c *fiber.Ctx) error {
		fmt.Println("stopping server in 5s")
		go func() {
			time.Sleep(5 * time.Second)
			fmt.Println("stopped")
			os.Exit(0)
		}()
		return c.SendString("RobotDocs server stopping in 5 seconds")
	})

	app.Static("/", "htmltest", fiber.Static{
		Compress:       false,
		ByteRange:      false,
		Browse:         true,
		Download:       false,
		Index:          "",
		CacheDuration:  0,
		MaxAge:         0,
		ModifyResponse: nil,
		Next:           nil,
	})

	panic(app.Listen(":3000"))
}
