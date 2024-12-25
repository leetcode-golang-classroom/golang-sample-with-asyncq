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
