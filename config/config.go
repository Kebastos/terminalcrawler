package config

import (
	"encoding/json"
	"log"
	"os"
)

// Описание конфигурационного файла сервиса
type Config struct {
	ServiceName          string              `json:"service_name"`
	DisplayName          string              `json:"display_name"`
	ServiceDescription   string              `json:"service_description"`
	MonitoringDb         MonitoringDbOptions `json:"monitoringDb"`
	CollectTerminalQuery string              `json:"collectTerminalQuery"`
	CronTime             string              `json:"cronTime"`
	GetEnvInfoQuery      string              `json:"getEnvInfoQuery"`
}

// Описание подключения к БД Мониторинга
type MonitoringDbOptions struct {
	Server string `json:"server"`
	DbName string `json:"dbName"`
}

// Инициализация файла конфигурации
func LoadConfiguration(file string) (Config, error) {
	var config Config
	configFile, err := os.Open(file)
	if err != nil {
		log.Println(err.Error())
		return Config{}, err
	}
	jsonParser := json.NewDecoder(configFile)
	err = jsonParser.Decode(&config)
	_ = configFile.Close()
	return config, err
}
