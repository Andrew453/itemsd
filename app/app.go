package app

import (
	"context"
	"github.com/pkg/errors"
	"log/slog"
	"os"
	"path/filepath"
	"prjs/itemsd/adapters/controllers/http/echor"
	"prjs/itemsd/adapters/gates/csvdb"
	"prjs/itemsd/config"
	"prjs/itemsd/usecase"
)

// App Общая структура приложения. Надстройка над всеми внутренними компонентами
type App struct {
	config     *config.Config
	db         *csvdb.CSVDB
	server     *echor.Server
	uc         *usecase.Usecase
	errorsChan chan error
	ctx        context.Context
}

// New Создание нового приложения
func New(ctx context.Context) (a *App, err error) {
	defer func() {
		err = errors.Wrap(err, "app (a *App) New()")
	}()
	a = &App{
		ctx: ctx,
	}
	err = a.manage("./config.hcl")
	if err != nil {
		return nil, err
	}
	a.errorsChan = make(chan error, 10)
	return a, nil
}

func (a *App) Run() <-chan error {

	go func() {
		var err error
		a.db, err = csvdb.NewCSVDB(csvdb.Config{FilePath: a.config.DBPath})
		if err != nil {
			a.errorsChan <- err
			return
		}

		a.uc, err = usecase.NewUsecase(a.db, usecase.Config{MaxItemsIDs: a.config.MaxItems})
		if err != nil {
			a.errorsChan <- err
			return
		}
		a.server, err = echor.NewServer(echor.Config{
			TLSEnable: a.config.TLSEnabled,
			Address:   a.config.Address,
		}, a.config.CertFile, a.config.KeyFile, a.uc)
		if err != nil {
			a.errorsChan <- err
			return
		}

		errs := a.server.Start()
		for {
			select {
			case <-a.ctx.Done():
				return
			case err = <-errs:
				slog.Error(err.Error())
				a.errorsChan <- err
				return
			}
		}
	}()
	return a.errorsChan
}

func (a *App) Stop() (err error) {
	defer func() {
		err = errors.Wrap(err, "app (a *App) Stop()")
	}()

	if a.server != nil {
		err = a.server.Stop()
		if err != nil {
			slog.Error(err.Error())
			return
		}
	}
	return
}

// manage Загрузка конфигурационного файла и настройка логгера
func (a *App) manage(confPath string) (err error) {
	defer func() {
		err = errors.Wrap(err, "app (a *App) manage()")
	}()
	cPath, err := filepath.Abs(confPath)
	if err != nil {
		err = errors.Wrap(err, "filepath.Abs(confPath)")
		return
	}
	a.config = &config.Config{}
	err = config.LoadConfigHCLFromFile(cPath, a.config)
	if err != nil {
		return
	}

	err = a.config.Validate()
	if err != nil {
		return err
	}

	if a.config.LogFile == "" {
		a.config.LogFile = "./error.log"
	}

	ePath, err := filepath.Abs(a.config.LogFile)
	if err != nil {
		err = errors.Wrap(err, "filepath.Abs(a.conf.Logs.EPath)")
		return err
	}
	eFile, err := os.OpenFile(ePath, os.O_APPEND|os.O_RDWR|os.O_CREATE|os.O_SYNC, 0666)
	if err != nil {
		err = errors.Wrap(err, "os.OpenFile(ePath, os.O_APPEND|os.O_RDWR|os.O_CREATE|os.O_SYNC, 0666)")
		return err
	}

	var level slog.Level
	switch a.config.LogLevel {
	case 1:
		level = slog.LevelError
	case 2:
		level = slog.LevelInfo
	case 3:
		level = slog.LevelDebug
	}

	lg := slog.New(slog.NewJSONHandler(eFile, &slog.HandlerOptions{
		Level: level,
	}))

	slog.SetDefault(lg)
	return nil
}
