package routes

import (
	"github.com/gofiber/fiber/v3"

	"github.com/saveblush/gofiber-v3-boilerplate/internal/core/utils/logger"
	"github.com/saveblush/gofiber-v3-boilerplate/internal/pgk/book"
)

// CronStart cron start
func (s *server) CronStart() {
	if !fiber.IsChild() {
		logger.Log.Info("Cron start")
		s.cronSchedule()
		s.cron.Start()
	}
}

// CronStop cron stop
func (s *server) CronStop() {
	if !fiber.IsChild() {
		s.cron.Stop()
	}
}

// cronSchedule cron schedule
func (s *server) cronSchedule() {
	// book service
	book := book.NewService()

	s.cron.AddFunc("*/1 * * * *", func() {
		book.Script(s.cctx)
	})
}
