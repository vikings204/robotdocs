package main

import (
	"fmt"
	"github.com/gofiber/fiber/v2"
	"io/fs"
	"os"
	"path/filepath"
	"strings"
	"sync"
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

	// tests DO NOT USE
	//fmt.Println(fixExternalRefs("og.html", "htmltest\\baller", "htmltest\\baller\\_external"))
	//fmt.Println(filepath.WalkDir("htmltest\\rev", func(path string, di fs.DirEntry, err error) error {
	//	if strings.Contains(path, ".htm") {
	//		fmt.Println(path)
	//		return fixExternalRefs(path, "", "htmltest\\rev\\_external")
	//	} else {
	//		return nil
	//	}
	//}))

	// multithreaded entire website
	var paths []string
	fmt.Println(filepath.WalkDir("htmltest\\rev", func(path string, di fs.DirEntry, err error) error {
		if strings.Contains(path, ".htm") {
			fmt.Println(path)
			paths = append(paths, path)
		}
		return nil
	}))
	var wg sync.WaitGroup
	wg.Add(len(paths))
	fmt.Println("going")
	for pn := range paths {
		go func(p string) {
			defer wg.Done()
			fmt.Println(fixExternalRefs(p, "", "htmltest\\rev\\_external"))
		}(paths[pn])
	}

	//wg.Wait()

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
		Index:          "NO_INDEX_PLEASE",
		CacheDuration:  0,
		MaxAge:         0,
		ModifyResponse: nil,
		Next:           nil,
	})

	panic(app.Listen(":3000"))
}
