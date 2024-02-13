package main

import (
	"fmt"
	"github.com/gofiber/fiber/v2"
	"os"
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
	fmt.Println(fixExternalRefs("htmltest\\baller\\og.html", "htmltest\\baller\\_external"))

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
