package main

import (
	"log"
	"os"
	"os/signal"
	"syscall"
	"time"

	"github.com/gofiber/fiber/v2"
)

func recovered() {
	if r := recover(); r != nil {
		log.Fatalf("[%s]: recovered in -> %v", value, r)
	}
}

func main() {
	defer recovered()

	con, err := connect()
	if err != nil {
		panic(err)
	}
	defer con.Close()

	s := &Server{fiber: fiber.New(Config), c: con}
	s.use()

	s.fiber.Get("/", func(c *fiber.Ctx) error {
		return c.JSON(&response{Timestamp: time.Now(), Code: "OK", Message: "Hello, World 👋!"})
	})
	s.fiber.Post("/search", func(c *fiber.Ctx) error {
		var dto *search
		if err = c.BodyParser(&dto); err != nil {
			return c.Status(fiber.StatusBadRequest).
				JSON(&response{Timestamp: time.Now(), Code: "BODY-PARSER", Message: fiber.ErrServiceUnavailable.Error()})
		}
		if errs := s.validate(*dto); len(errs) > 0 {
			return c.Status(fiber.StatusNotAcceptable).JSONP(&errorValidator{
				Code:           "VALIDATOR",
				Message:        syscall.Errno(01).Error(),
				ErrorResponses: errs,
			})
		}

		entries, err := s.c.search(dto)
		if err != nil {
			return err
		}

		return c.JSON(entries)
	})
	s.fiber.Use(
		func(c *fiber.Ctx) error {
			return c.Status(fiber.StatusNotFound).
				JSON(&response{Timestamp: time.Now(), Code: "NOTFOUND", Message: "not found route: " + c.Path()})
		},
	)

	go s.start()

	c := make(chan os.Signal, 4)
	signal.Notify(c, os.Interrupt, syscall.SIGINT, syscall.SIGTERM, syscall.SIGKILL)
	<-c

	s.stop()
}
