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

	active, err := manager.GetTask(t1.ID)
	if err != nil || active.State != StatusProcessing {
		t.Errorf("expected task 1 to be %s while running, got %s (err=%v)", StatusProcessing, active.State, err)
	}

	close(bp.release)
	manager.Stop()

	for _, id := range []int{t1.ID, t2.ID} {
		task, _ := manager.GetTask(id)
		if task.State != StatusCompleted || task.Result.Strength != "Strong" {
			t.Errorf("task %d invalid final state: %+v", id, task)
		}
	}
}

func TestQueueFullHandling(t *testing.T) {
	bp := blockingProcessor{
		started: make(chan struct{}, 2),
		release: make(chan struct{}),
	}

	manager := NewTaskManager(bp)
	manager.Start(2, 1)

	t1, _ := manager.CreateTask("Pass123!")
	t2, _ := manager.CreateTask("Pass456!")
	t3, _ := manager.CreateTask("Pass789!")
	t4, _ := manager.CreateTask("Pass000!")

	if err := manager.SubmitTask(t1.ID); err != nil {
		t.Fatalf("submit 1 failed: %v", err)
	}
	<-bp.started
	if err := manager.SubmitTask(t2.ID); err != nil {
		t.Fatalf("submit 2 failed: %v", err)
	}

	<-bp.started

	if err := manager.SubmitTask(t3.ID); err != nil {
		t.Fatalf("submit 3 failed: %v", err)
	}

	err := manager.SubmitTask(t4.ID)
	if !errors.Is(err, ErrQueueFull) {
		t.Fatalf("expected %v, got %v", ErrQueueFull, err)
	}

	fetched4, _ := manager.GetTask(t4.ID)
	if fetched4.State != StatusPending {
		t.Errorf("rejected task should remain %s, got %s", StatusPending, fetched4.State)
	}

	close(bp.release)
	manager.Stop()
}

func TestInvalidAndDuplicateSubmissions(t *testing.T) {
	bp := blockingProcessor{
		started: make(chan struct{}, 3),
		release: make(chan struct{}),
	}

	manager := NewTaskManager(bp)
	manager.Start(2, 2)

	if err := manager.SubmitTask(999); !errors.Is(err, ErrTaskNotFound) {
		t.Errorf("expected %v, got %v", ErrTaskNotFound, err)
	}

	t1, _ := manager.CreateTask("Valid123!")
	t2, _ := manager.CreateTask("Valid456!")
	t3, _ := manager.CreateTask("Valid789!")

	_ = manager.SubmitTask(t1.ID)
	_ = manager.SubmitTask(t2.ID)
	<-bp.started
	<-bp.started

	if err := manager.SubmitTask(t1.ID); !errors.Is(err, ErrInvalidState) {
		t.Errorf("expected %v while processing, got %v", ErrInvalidState, err)
	}

	_ = manager.SubmitTask(t3.ID)
	if err := manager.SubmitTask(t3.ID); !errors.Is(err, ErrInvalidState) {
		t.Errorf("expected %v while queued, got %v", ErrInvalidState, err)
	}

	close(bp.release)
	manager.Stop()

	manager2 := NewTaskManager(PasswordProcessor{})
	manager2.Start(2, 2)
	doneTask, _ := manager2.CreateTask("Completed1!")
	if _, err := manager2.ProcessTask(doneTask.ID); err != nil {
		t.Fatalf("process failed: %v", err)
	}
	if err := manager2.SubmitTask(doneTask.ID); !errors.Is(err, ErrInvalidState) {
		t.Errorf("expected %v for completed task, got %v", ErrInvalidState, err)
	}
	manager2.Stop()
}

func TestStopConcurrentWithSubmissionsAndIdempotent(t *testing.T) {
	manager := NewTaskManager(PasswordProcessor{})
	manager.Start(3, 20)

	var accepted []int
	var mu sync.Mutex
	var wg sync.WaitGroup

	for i := 0; i < 15; i++ {
		wg.Add(1)
		go func() {
			defer wg.Done()
			task, err := manager.CreateTask("Valid123!")
			if err != nil {
				t.Errorf("create failed: %v", err)
				return
			}
			err = manager.SubmitTask(task.ID)
			if err == nil {
				mu.Lock()
				accepted = append(accepted, task.ID)
				mu.Unlock()
			} else if !errors.Is(err, ErrManagerStopped) && !errors.Is(err, ErrQueueFull) {
				t.Errorf("unexpected submit error: %v", err)
			}
		}()
	}

	manager.Stop()
	manager.Stop()
	wg.Wait()

	extra, _ := manager.CreateTask("AfterStop1!")
	if err := manager.SubmitTask(extra.ID); !errors.Is(err, ErrManagerStopped) {
		t.Errorf("expected %v after stop, got %v", ErrManagerStopped, err)
	}

	for _, id := range accepted {
		task, err := manager.GetTask(id)
		if err != nil {
			t.Errorf("get failed: %v", err)
			continue
		}
		if task.State != StatusCompleted {
			t.Errorf("accepted task %d lost during shutdown: %s", id, task.State)
		}
	}
}

func TestConcurrentSharedStateAccess(t *testing.T) {
	manager := NewTaskManager(PasswordProcessor{})
	manager.Start(4, 20)

	ids := make(chan int, 15)
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
			if err := manager.SubmitTask(task.ID); err != nil {
				t.Errorf("submit error: %v", err)
				return
			}
			if _, err := manager.GetTask(task.ID); err != nil {
				t.Errorf("get error: %v", err)
				return
			}
			ids <- task.ID
		}()
	}

	wg.Wait()
	close(ids)
	manager.Stop()

	for id := range ids {
		task, err := manager.GetTask(id)
		if err != nil {
			t.Fatalf("final get failed for %d: %v", id, err)
		}
		if task.State != StatusCompleted || task.Result.Strength != "Strong" {
			t.Errorf("task %d unexpected result: %+v", id, task)
		}
	}
}
