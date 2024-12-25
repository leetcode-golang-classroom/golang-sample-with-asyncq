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
