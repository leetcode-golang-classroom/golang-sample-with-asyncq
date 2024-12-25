# golang-sample-with-asyncq

This repository is demo how to use asyncq with handle complex job

## setup task on producer

```golang
package main

import (
	"log"
	"time"

	"github.com/hibiken/asynq"
	"github.com/leetcode-golang-classroom/golang-sample-with-asyncq/internal/config"
	"github.com/leetcode-golang-classroom/golang-sample-with-asyncq/internal/tasks"
)

func main() {
	redisClientOpt, err := asynq.ParseRedisURI(config.AppConfig.RedisURL)
	if err != nil {
		log.Fatalf("failed to parse redis url: %v", err)
	}
	client := asynq.NewClient(redisClientOpt)
	defer client.Close()

	task, err := tasks.NewEmailTask("test1@test.com", "Welcome!", "Thank you for signing up.")
	if err != nil {
		log.Fatalf("could not create task: %v\n", err)
	}

	info, err := client.Enqueue(task)
	if err != nil {
		log.Fatalf("could not enqueue task: %v", err)
	}
	log.Printf("enqueued task: id=%s queue=%s\n", info.ID, info.Queue)

	task, err = tasks.NewReportTask(10)
	if err != nil {
		log.Fatalf("could not create task: %v\n", err)
	}
	info, err = client.Enqueue(task, asynq.MaxRetry(10), asynq.Timeout(3*time.Minute))
	if err != nil {
		log.Fatalf("could not enqueue task: %v\n", err)
	}
	log.Printf("enqueued task: id=%s queue=%s\n", info.ID, info.Queue)

	task, err = tasks.NewImageProcessingTask("http://test.com/image")
	if err != nil {
		log.Fatalf("could not create task: %v\n", err)
	}
	info, err = client.Enqueue(task, asynq.MaxRetry(10), asynq.Timeout(3*time.Minute))
	if err != nil {
		log.Fatalf("could not enqueue task: %v\n", err)
	}
	log.Printf("enqueued task: id=%s queue=%s\n", info.ID, info.Queue)

	task, err = tasks.NewEmailTask("test3@test.com", "Welcome!", "Thank you for signing up.")
	if err != nil {
		log.Fatalf("could not create task: %v\n", err)
	}

	info, err = client.Enqueue(task, asynq.ProcessIn(10*time.Minute), asynq.Queue("critical"))
	if err != nil {
		log.Fatalf("could not enqueue task: %v", err)
	}
	log.Printf("enqueued task: id=%s queue=%s\n", info.ID, info.Queue)
}

```

## setup consumer to consumer jobs

```golang
package main

import (
	"log"

	"github.com/hibiken/asynq"
	"github.com/leetcode-golang-classroom/golang-sample-with-asyncq/internal/config"
	"github.com/leetcode-golang-classroom/golang-sample-with-asyncq/internal/tasks"
)

func main() {
	redisClientOpt, err := asynq.ParseRedisURI(config.AppConfig.RedisURL)
	if err != nil {
		log.Fatalf("failed to parse redis url: %v", err)
	}

	srv := asynq.NewServer(
		redisClientOpt,
		asynq.Config{
			Concurrency: 10,
			Queues: map[string]int{
				"critical": 6,
				"default":  3,
				"low":      1,
			},
		},
	)

	// mux maps a type to a handler
	mux := asynq.NewServeMux()
	mux.HandleFunc(tasks.TypeEmail, tasks.EmailTaskHandler)
	mux.HandleFunc(tasks.TypeReport, tasks.ReportTaskHandler)
	mux.Handle(tasks.TypeImageProcessing, tasks.NewImageProcessor())
	if err := srv.Run(mux); err != nil {
		log.Fatalf("could not run server: %v", err)
	}
}
```