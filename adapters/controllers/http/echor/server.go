// Package echor Пакет управления HTTP/HTTPS
// Основной объект - Server
package echor

import (
	"context"
	"crypto/tls"
	"github.com/labstack/echo/v4"
	"github.com/labstack/echo/v4/middleware"
	"github.com/pkg/errors"
	"log/slog"
	"prjs/itemsd/adapters/controllers/http/echor/handler"
	"prjs/itemsd/model/usecase"
)

// Server структура управления HTTPS сервером
type Server struct {
	cfg        Config      // Конфигурация
	certFile   interface{} // серфикат
	keyFile    interface{} // ключ
	errorsChan chan error  // канал ошибок
	server     *echo.Echo  //  переменная для управления echo сервером
}

// NewServer Конструктор для создания нового экзмепляра Server
func NewServer(cfg Config, certFile interface{}, keyFile interface{}, uc usecase.Usecase) (s *Server, err error) {
	defer func() {
		err = errors.Wrap(err, "echor NewServer()")
	}()

	err = cfg.Validate()
	if err != nil {
		return nil, err
	}

	h := handler.NewHandler(uc)

	es := echo.New()

	es.Use(middleware.CORS())
	es.Use(middleware.Recover())

	setRoutes(es, h)
	return &Server{cfg: cfg, certFile: certFile, keyFile: keyFile, errorsChan: make(chan error, 10), server: es}, nil
}

// Start Запуск HTTP/HTTPS сервера
func (s *Server) Start() <-chan error {
	go func() {
		if s.cfg.TLSEnable {
			s.server.Server.TLSConfig = &tls.Config{
				MinVersion: tls.VersionTLS12,
			}
			slog.Info("Starting HTTPS server")
			err := s.server.StartTLS(s.cfg.Address, s.certFile, s.keyFile)
			if err != nil {
				err = errors.Wrap(err, "s.server.StartTLS(s.cfg.Address, s.certFile, s.keyFile)")
				s.errorsChan <- err
				return
			}
		} else {
			slog.Info("Starting HTTP server")
			err := s.server.Start(s.cfg.Address)
			if err != nil {
				err = errors.Wrap(err, "s.server.Start(s.cfg.Address)")
				s.errorsChan <- err
				return
			}
		}
	}()
	return s.errorsChan
}

// Stop Завершение работы сервера
func (s *Server) Stop() (err error) {
	defer func() {
		err = errors.Wrap(err, "echor (s *Server) Stop()")
	}()
	err = s.server.Shutdown(context.Background())
	if err != nil {
		err = errors.Wrap(err, "s.server.Shutdown(context.Background())")
		return err
	}
	return nil
}
