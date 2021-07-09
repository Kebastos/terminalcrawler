package main

import (
	"fmt"
	"log"
	"os"

	"github.com/kardianos/service"
	"gitlab.tools.russianpost.ru/eas-support/monitoring-eas-ops/microservices/terminalcrawler/config"
	"gitlab.tools.russianpost.ru/eas-support/monitoring-eas-ops/microservices/terminalcrawler/core"
)

// Начальная инициализация сервиса и всех его компонентов
func main() {
	fmt.Println("TerminalCrawler running!")

	svc := &core.Service{
		Stdlog: log.New(os.Stdout, "Info ", log.Ldate|log.Ltime),
		Errlog: log.New(os.Stderr, "Error ", log.Ldate|log.Ltime|log.Lshortfile),
	}

	var err error

	if svc.Cfg, err = config.LoadConfiguration("config.json"); err != nil {
		svc.Errlog.Fatalln("Ошибка при загрузке Json-конфигурации: ", err)
	}

	cfg := &service.Config{
		Name:        svc.Cfg.ServiceName,
		DisplayName: svc.Cfg.DisplayName,
		Description: svc.Cfg.ServiceDescription,
	}

	s, err := service.New(svc, cfg)
	if err != nil {
		svc.Errlog.Fatalln("Ошибка при создании сервиса:", err)
	}

	syslogger, err := s.Logger(nil)
	if err != nil {
		svc.Errlog.Fatalln("Ошибка получения логера:", err)
	}

	if err = s.Run(); err != nil {
		err = syslogger.Errorf("сервис завершил работу с неожиданной ошибкой:", err)
		if err != nil {
			panic("Критическая ошибка сервиса. Сервис остановлен!")
		}
	}
}
