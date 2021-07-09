package core

import (
	"github.com/reugn/go-quartz/quartz"
)

type Job struct {
	service *Service
	desc    string
}

// Возвращает описание задания
func (j Job) Description() string {
	return j.desc
}

// Возвращает уникальный ключ задания
func (j Job) Key() int {
	return quartz.HashCode(j.Description())
}

// Вызов выполнения задания
func (j Job) Execute() {
	l, err := j.service.MsSqlConn.SelectEnvInfo(j.service.Cfg.GetEnvInfoQuery)
	if err != nil {
		j.service.Errlog.Fatal("Ошибка возникла!", err)
	}

	for i := 0; i < len(l); i++ {
		if l[i].Mrc != "Тест" && l[i].Mrc != "Синк" && l[i].Mrc != "АУП" {
			err = j.service.MsSqlConn.MergeQuery(j.service.Cfg.CollectTerminalQuery, l[i])
			if err != nil {
				j.service.Errlog.Printf("Ошибка обработки информации с региона %v. \n Ошибка: \n %v", l[i].Caption, err)
			}
		}
	}
}
