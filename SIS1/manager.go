package main

import "sync"

type TaskManager struct {
	mu        sync.RWMutex
	tasks     map[int]*Task
	nextID    int
	processor Processor
	queue     chan int
	wg        sync.WaitGroup
	started   bool
	stopped   bool
}

func NewTaskManager(p Processor) *TaskManager {
	return &TaskManager{
		tasks:     make(map[int]*Task),
		nextID:    1,
		processor: p,
	}
}

func (m *TaskManager) Start(workers, queueSize int) {
	m.mu.Lock()
	defer m.mu.Unlock()

	if m.started {
		return
	}
	if workers < 2 {
		workers = 2
	}
	if queueSize < 1 {
		queueSize = 1
	}

	m.queue = make(chan int, queueSize)
	m.started = true
	m.stopped = false

	for i := 0; i < workers; i++ {
		m.wg.Add(1)
		go m.worker()
	}
}

func (m *TaskManager) Stop() {
	m.mu.Lock()
	if !m.started || m.stopped {
		m.mu.Unlock()
		return
	}
	m.stopped = true
	close(m.queue)
	m.mu.Unlock()

	m.wg.Wait()
}

func (m *TaskManager) CreateTask(password string) (Task, error) {
	if err := ValidatePassword(password); err != nil {
		return Task{}, err
	}

	m.mu.Lock()
	defer m.mu.Unlock()

	task := &Task{
		ID:       m.nextID,
		Password: password,
		State:    StatusPending,
	}

	m.tasks[task.ID] = task
	m.nextID++

	return *task, nil
}

func (m *TaskManager) GetTask(id int) (Task, error) {
	m.mu.RLock()
	defer m.mu.RUnlock()

	task, exists := m.tasks[id]
	if !exists {
		return Task{}, ErrTaskNotFound
	}
	return *task, nil
}

func (m *TaskManager) SubmitTask(id int) error {
	m.mu.Lock()
	defer m.mu.Unlock()

	if !m.started || m.stopped {
		return ErrManagerStopped
	}

	task, exists := m.tasks[id]
	if !exists {
		return ErrTaskNotFound
	}
	if task.State != StatusPending {
		return ErrInvalidState
	}

	select {
	case m.queue <- id:
		_ = task.SetState(StatusQueued)
		return nil
	default:
		return ErrQueueFull
	}
}
func (m *TaskManager) ProcessTask(id int) (Task, error) {
	m.mu.Lock()
	task, exists := m.tasks[id]
	if !exists {
		m.mu.Unlock()
		return Task{}, ErrTaskNotFound
	}
	if err := task.SetState(StatusProcessing); err != nil {
		m.mu.Unlock()
		return Task{}, err
	}
	password := task.Password
	m.mu.Unlock()

	result := m.processor.Process(password)

	m.mu.Lock()
	defer m.mu.Unlock()

	task.Result = result
	_ = task.SetState(StatusCompleted)

	return *task, nil
}

func (m *TaskManager) worker() {
	defer m.wg.Done()

	for id := range m.queue {
		m.mu.Lock()
		task, exists := m.tasks[id]
		if !exists {
			m.mu.Unlock()
			continue
		}
		if err := task.SetState(StatusProcessing); err != nil {
			m.mu.Unlock()
			continue
		}
		password := task.Password
		m.mu.Unlock()

		result := m.processor.Process(password)

		m.mu.Lock()
		task.Result = result
		_ = task.SetState(StatusCompleted)
		m.mu.Unlock()
	}
}
