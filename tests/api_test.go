package tests

import (
	"math/rand/v2"
	"testing"
	"time"

	"github.com/hibiken/asynq"
)

func TestAsynqLoop(t *testing.T) {
	var (
		client    *asynq.Client = asynq.NewClient(asynq.RedisClientOpt{Addr: "127.0.0.1:6379"})
		taskTypes []string      = []string{"test:task1", "test:task2", "test:task3"}
		err       error
	)

	defer client.Close()

	for range 128 {
		var (
			task *asynq.Task = asynq.NewTask(taskTypes[rand.IntN(len(taskTypes))], randomBytes(rand.IntN(128)+128), asynq.ProcessIn(time.Millisecond*time.Duration(rand.IntN(10000))))
			info *asynq.TaskInfo
		)

		if info, err = client.Enqueue(task); err != nil {
			t.Fatalf("failed to enqueue task: %v", err)
			return
		}

		t.Logf("enqueued task: type=%s, id=%s", task.Type(), info.ID)
	}

	// Consumer
	
}
