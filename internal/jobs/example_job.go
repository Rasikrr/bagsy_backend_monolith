// nolint
package jobs

import (
	"context"

	"github.com/Rasikrr/core/log"
)

// Интерфейс джобы
type Job interface {
	Name() string
	Schedule() string
	Run()
}

type ExampleJob struct {
	name     string
	schedule string

	// Твои зависимости(redis, репо, клиенты т.д)
}

func NewExampleJob(name, schedule string) *ExampleJob {
	return &ExampleJob{
		name:     name,
		schedule: schedule,
	}
}

func (j *ExampleJob) Name() string {
	return j.name
}

func (j *ExampleJob) Schedule() string {
	return j.schedule
}

func (j *ExampleJob) Run() {
	log.Info(
		context.Background(),
		"job working",
		log.String("name", j.Name()),
		log.String("schedule", j.Schedule()),
	)
}
