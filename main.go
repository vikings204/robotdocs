package main

import (
	"fmt"
	"github.com/gofiber/fiber/v2"
)

func main() {
	fmt.Println("binary error:", InitializeBinary())

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

	app.Static("/", "./htmltest", fiber.Static{
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

	if err := app.Listen(":3000"); err != nil {
		panic(err)
	}
}
