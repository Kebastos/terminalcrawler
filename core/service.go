package core

import (
	"log"
	"sync"

	"github.com/kardianos/service"
	"gitlab.tools.russianpost.ru/eas-support/monitoring-eas-ops/microservices/terminalcrawler/config"
)

type Service struct {
	MsSqlConn      *MsSqlConn
	Stdlog, Errlog *log.Logger
	Cfg            config.Config
}

// Запуск сервиса
func (s *Service) Start(svc service.Service) (err error) {
	s.Stdlog.Printf("Инициализация подключения к базе данных мониторинга([%v].[%v]).", s.Cfg.MonitoringDb.Server, s.Cfg.MonitoringDb.DbName)
	s.MsSqlConn = NewMsSqlConn(s)
	s.Stdlog.Print("Инициализация шедулера")
	wg := new(sync.WaitGroup)
	wg.Add(1)
	go Scheduler(s, wg)

	wg.Wait()

	return nil
}

// Остановка сервиса
func (s *Service) Stop(svc service.Service) (err error) {
	s.Stdlog.Println("Производится остановка сервиса.")
	s.MsSqlConn.SqlConn.Close()

	return
}
