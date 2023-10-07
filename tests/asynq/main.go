package main

import (
	"context"
	"github.com/hibiken/asynq"
	"log"
	"time"
)

func main() {
	go Client()
	go Server(1)
	go Server(2)
	for {
		time.Sleep(time.Hour)
	}
}

func Client() {
	scheduler := asynq.NewScheduler(
		asynq.RedisClientOpt{Addr: ":6379"},
		&asynq.SchedulerOpts{Location: time.Local},
	)

	task := asynq.NewTask("email:welcome", nil)

	// You can use cron spec string to specify the schedule.
	entryID, err := scheduler.Register("@every 2s", task, asynq.Queue("welcome"))
	if err != nil {
		log.Fatal(err)
	}
	log.Printf("registered an entry: %q\n", entryID)

	// You can use "@every <duration>" to specify the interval.
	entryID, err = scheduler.Register("@every 2s", asynq.NewTask("email:reminder", nil),
		asynq.Queue("reminder"))
	if err != nil {
		log.Fatal(err)
	}
	scheduler.Unregister(entryID)

	log.Printf("registered an entry: %q\n", entryID)

	// You can also pass options.
	entryID, err = scheduler.Register("@every 24h", task, asynq.Queue("myqueue"))
	if err != nil {
		log.Fatal(err)
	}
	log.Printf("registered an entry: %q\n", entryID)

	// ... Register tasks

	if err := scheduler.Run(); err != nil {
		log.Fatal(err)
	}
}

func Server(id int) {
	queue := "welcome"
	if id != 1 {
		queue = "reminder"
	}
	srv := asynq.NewServer(
		asynq.RedisClientOpt{Addr: "localhost:6379"},
		asynq.Config{Concurrency: 10, Queues: map[string]int{queue: 1}},
	)

	mux := asynq.NewServeMux()
	mux.HandleFunc("email:welcome", func(ctx context.Context, task *asynq.Task) error {
		log.Printf("id:%v email:welcome handler: type:%v payload:%v", id, task.Type(), string(task.Payload()))
		return nil
	})
	mux.HandleFunc("email:reminder", func(ctx context.Context, task *asynq.Task) error {
		log.Printf("id:%v email:reminder handler: type:%v payload:%v", id, task.Type(), string(task.Payload()))
		return nil
	})

	if err := srv.Run(mux); err != nil {
		log.Fatal(err)
	}
}
