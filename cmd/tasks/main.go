package main

import (
	"log"

	"github.com/hibiken/asynq"
	"github.com/nneji123/ecommerce-golang/config"
	emailTasks "github.com/nneji123/ecommerce-golang/internal/email/emailsender"
)

// A list of task types
const (
	TypeEmailDelivery = "email:deliver"
)

func main() {
	cfg, err := config.LoadConfig()
	if err != nil {
		log.Fatalf("Error loading configuration: %s", err)
	}

	redisAddr := cfg.RedisAddr
	log.Println(redisAddr)

	client := asynq.NewClient(asynq.RedisClientOpt{Addr: redisAddr})
	defer client.Close()
	srv := asynq.NewServer(
		asynq.RedisClientOpt{Addr: redisAddr},
		asynq.Config{
			Concurrency: 10,
			Queues: map[string]int{
				"critical": 6,
				"default":  3,
				"low":      1,
			},
		},
	)
	mux := asynq.NewServeMux()
	mux.HandleFunc(TypeEmailDelivery, emailTasks.HandleEmailDeliveryTask)
	go func() {
		if err := srv.Run(mux); err != nil {
			log.Fatalf("could not run server: %v", err)
		}
	}()
	select {}
}
