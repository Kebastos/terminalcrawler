package core

import (
	"sync"

	"github.com/reugn/go-quartz/quartz"
)

// Планировщик задания
func Scheduler(service *Service, wg *sync.WaitGroup) {
	sched := quartz.NewStdScheduler()
	cronTrigger, _ := quartz.NewCronTrigger(service.Cfg.CronTime)
	cronJob := Job{service, "TerminalCrawlerJob"}
	sched.Start()
	sched.ScheduleJob(&cronJob, cronTrigger)
}
