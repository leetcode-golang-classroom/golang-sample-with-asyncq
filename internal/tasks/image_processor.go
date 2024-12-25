package tasks

import (
	"context"
	"encoding/json"
	"fmt"
	"log"
	"time"

	"github.com/hibiken/asynq"
)

type ImageProcessor struct {
}

type ImageProcessingPayload struct {
	ImageURL string
}

func NewImageProcessor() *ImageProcessor {
	return &ImageProcessor{}
}

func (processor *ImageProcessor) ProcessTask(ctx context.Context, t *asynq.Task) error {
	var p ImageProcessingPayload
	if err := json.Unmarshal(t.Payload(), &p); err != nil {
		return fmt.Errorf("json.Unmarshal failed: %v: %w", err, asynq.SkipRetry)
	}
	log.Printf("Resizing image: URL=%s", p.ImageURL)
	time.Sleep(5 * time.Second)
	return nil
}

func NewImageProcessingTask(imageURL string) (*asynq.Task, error) {
	payload, err := json.Marshal(ImageProcessingPayload{
		ImageURL: imageURL,
	})
	if err != nil {
		return nil, err
	}
	return asynq.NewTask(TypeImageProcessing, payload), nil
}
