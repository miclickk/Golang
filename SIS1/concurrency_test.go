package main

import (
	"errors"
	"sync"
	"testing"
)

type blockingProcessor struct {
	started chan struct{}
	release chan struct{}
}

func (p blockingProcessor) Process(password string) AnalysisResult {
	p.started <- struct{}{}
	<-p.release
	return PasswordProcessor{}.Process(password)
}

func TestWorkerPoolMultipleTasks(t *testing.T) {
	manager := NewTaskManager(PasswordProcessor{})
	manager.Start(3, 5)

	inputs := []struct {
		pwd      string
		strength string
	}{
		{"abcdef", "Weak"},
		{"Abcdef12", "Medium"},
		{"Abcdef12!", "Strong"},
		{"12345", "Weak"},
	}

	ids := make([]int, 0, len(inputs))
	for _, in := range inputs {
		task, err := manager.CreateTask(in.pwd)
		if err != nil {
			t.Fatalf("create failed: %v", err)
		}
		if err := manager.SubmitTask(task.ID); err != nil {
			t.Fatalf("submit failed: %v", err)
		}
		ids = append(ids, task.ID)
	}

	manager.Stop()

	for i, id := range ids {
		task, err := manager.GetTask(id)
		if err != nil {
			t.Fatalf("get failed: %v", err)
		}
		if task.State != StatusCompleted {
			t.Errorf("task %d expected %s, got %s", id, StatusCompleted, task.State)
		}
		if task.Result.Strength != inputs[i].strength {
			t.Errorf("task %d expected %s, got %s", id, inputs[i].strength, task.Result.Strength)
		}
	}
}

func TestConcurrentOverlappingProcessing(t *testing.T) {
	bp := blockingProcessor{
		started: make(chan struct{}, 2),
		release: make(chan struct{}),
	}

	manager := NewTaskManager(bp)
	manager.Start(2, 2)

	t1, _ := manager.CreateTask("Pass123!")
	t2, _ := manager.CreateTask("Pass456!")

	if err := manager.SubmitTask(t1.ID); err != nil {
		t.Fatalf("submit 1 failed: %v", err)
	}
	if err := manager.SubmitTask(t2.ID); err != nil {
		t.Fatalf("submit 2 failed: %v", err)
	}

	<-bp.started
	<-bp.started
	close(bp.release)

	manager.Stop()

	for _, id := range []int{t1.ID, t2.ID} {
		task, _ := manager.GetTask(id)
		if task.State != StatusCompleted {
			t.Errorf("task %d not completed", id)
		}
	}
}

func TestQueueFullHandling(t *testing.T) {
	bp := blockingProcessor{
		started: make(chan struct{}, 2),
		release: make(chan struct{}),
	}

	manager := NewTaskManager(bp)
	manager.Start(1, 1)

	t1, _ := manager.CreateTask("Pass1!")
	t2, _ := manager.CreateTask("Pass2!")
	t3, _ := manager.CreateTask("Pass3!")

	if err := manager.SubmitTask(t1.ID); err != nil {
		t.Fatalf("submit 1 failed: %v", err)
	}
	<-bp.started

	if err := manager.SubmitTask(t2.ID); err != nil {
		t.Fatalf("submit 2 failed: %v", err)
	}

	err := manager.SubmitTask(t3.ID)
	if !errors.Is(err, ErrQueueFull) {
		t.Fatalf("expected %v, got %v", ErrQueueFull, err)
	}

	fetched3, _ := manager.GetTask(t3.ID)
	if fetched3.State != StatusPending {
		t.Errorf("rejected task should remain %s, got %s", StatusPending, fetched3.State)
	}

	close(bp.release)
	manager.Stop()
}

func TestInvalidAndDuplicateSubmissions(t *testing.T) {
	bp := blockingProcessor{
		started: make(chan struct{}, 2),
		release: make(chan struct{}),
	}

	manager := NewTaskManager(bp)
	manager.Start(2, 2)

	if err := manager.SubmitTask(999); !errors.Is(err, ErrTaskNotFound) {
		t.Errorf("expected %v, got %v", ErrTaskNotFound, err)
	}

	task, _ := manager.CreateTask("Valid123!")
	if err := manager.SubmitTask(task.ID); err != nil {
		t.Fatalf("first submit failed: %v", err)
	}
	<-bp.started

	if err := manager.SubmitTask(task.ID); !errors.Is(err, ErrInvalidState) {
		t.Errorf("expected %v while processing, got %v", ErrInvalidState, err)
	}

	close(bp.release)
	manager.Stop()
}

func TestStopRejectsNewSubmissionsAndIsIdempotent(t *testing.T) {
	manager := NewTaskManager(PasswordProcessor{})
	manager.Start(2, 2)

	t1, _ := manager.CreateTask("Valid123!")
	t2, _ := manager.CreateTask("Another123!")

	if err := manager.SubmitTask(t1.ID); err != nil {
		t.Fatalf("submit failed: %v", err)
	}

	manager.Stop()
	manager.Stop()

	if err := manager.SubmitTask(t2.ID); !errors.Is(err, ErrManagerStopped) {
		t.Errorf("expected %v after stop, got %v", ErrManagerStopped, err)
	}

	fetched1, _ := manager.GetTask(t1.ID)
	if fetched1.State != StatusCompleted {
		t.Errorf("accepted task was lost during shutdown: %s", fetched1.State)
	}
}

func TestConcurrentSharedStateAccess(t *testing.T) {
	manager := NewTaskManager(PasswordProcessor{})
	manager.Start(4, 20)

	var wg sync.WaitGroup
	for i := 0; i < 15; i++ {
		wg.Add(1)
		go func() {
			defer wg.Done()
			task, err := manager.CreateTask("Concurrent1!")
			if err != nil {
				t.Errorf("create error: %v", err)
				return
			}
			_ = manager.SubmitTask(task.ID)
			_, _ = manager.GetTask(task.ID)
		}()
	}

	wg.Wait()
	manager.Stop()
}
