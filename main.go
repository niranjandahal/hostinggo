package main

import (
	"log"
	"os"

	"github.com/gofiber/fiber/v2"
)

func main (){
	app := fiber.New()

	app.Get("/" , func(c *fiber.Ctx) error {
		return c.SendString("first sample page")

	})

	app.Get("/env", func(c *fiber.Ctx) error {
		return c.SendString("TEst ENV" + os.Getenv("Test_env"))
	})

	port := os.Getenv("PORT")

	if port == ""{
		port = "3000"
	}

	log.Fatal(app.Listen("0.0.0.0:"+ port))
}