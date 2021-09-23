package main

import (
	"github.com/gofiber/fiber/v2/middleware/compress"
	"github.com/gofiber/fiber/v2/middleware/cors"
	"github.com/gofiber/fiber/v2/middleware/etag"
	"github.com/gofiber/fiber/v2/middleware/logger"
	"github.com/gofiber/fiber/v2/middleware/recover"
	"github.com/gofiber/fiber/v2/middleware/requestid"
	"log"
	"strings"
	"time"

	"github.com/go-playground/validator/v10"
	"github.com/gofiber/fiber/v2"
)

type Server struct {
	fiber *fiber.App
	c     *client
}

const value = "go-ldap/0.0.1"

var val = validator.New()

var Config = fiber.Config{
	ServerHeader: value,
	BodyLimit:    10 * 1024 * 1024, // 10 MB
	ReadTimeout:  30 * time.Second,
	WriteTimeout: 30 * time.Second,
	IdleTimeout:  10 * time.Second,
	ErrorHandler: func(c *fiber.Ctx, err error) error {
		var code = fiber.StatusInternalServerError
		var res = &response{
			Timestamp: time.Now(),
			Code:      "SERVER-ERROR",
			Message:   "Internal Server Error",
		}
		if e, ok := err.(*fiber.Error); ok {
			code = e.Code
			res.Message = e.Message
		}

		return c.Status(code).JSON(res)
	},
}

func (s *Server) validate(i interface{}) []*errorResponse {
	var errors []*errorResponse
	if err := val.Struct(i); err != nil {
		for _, e := range err.(validator.ValidationErrors) {
			var element errorResponse
			element.Tag = e.Tag()
			element.Name = strings.ToLower(e.Field())
			switch e.Tag() {
			case "gte":
				element.Message = e.Field() + " value must be greater than " + e.Param()
			case "lte":
				element.Message = e.Field() + " value must be lower than " + e.Param()
			default:
				element.Message = e.Field() + " is required type: " + e.Type().Name()
			}
			errors = append(errors, &element)
		}
	}
	return errors
}

func (s *Server) use() {
	s.fiber.Use(
		cors.New(cors.Config{MaxAge: 3600}),
		recover.New(),
		logger.New(),
		requestid.New(),
		etag.New(),
		compress.New(compress.Config{
			Level: compress.LevelBestSpeed, // 1
		}),
		func(c *fiber.Ctx) error {
			c.Response().Header.Set(fiber.HeaderAccessControlAllowOrigin, "*")
			c.Response().Header.Set(fiber.HeaderContentType, fiber.MIMEApplicationJSONCharsetUTF8)
			c.Response().Header.Set(fiber.HeaderServer, value)
			return c.Next()
		},
	)
}

func (s *Server) start() {
	defer recovered()
	log.Printf("[%s]: %s", value, "starting...")

	if err := s.fiber.Listen("127.0.0.1:3000"); err != nil {
		panic(err)
	}
}

func (s *Server) stop() {
	defer recovered()

	log.Printf("[%s]: %s", value, "stopping...")
	if err := s.fiber.Shutdown(); err != nil {
		panic(err)
	}
}
