package main

type TaskManager struct {
	tasks     map[int]*Task
	nextID    int
	processor Processor
}

func NewTaskManager(p Processor) *TaskManager {
	return &TaskManager{
		tasks:     make(map[int]*Task),
		nextID:    1,
		processor: p,
	}
}

func (m *TaskManager) CreateTask(password string) (Task, error) {
	if err := ValidatePassword(password); err != nil {
		return Task{}, err
	}

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
	task, exists := m.tasks[id]
	if !exists {
		return Task{}, ErrTaskNotFound
	}
	return *task, nil
}

func (m *TaskManager) ProcessTask(id int) (Task, error) {
	task, exists := m.tasks[id]
	if !exists {
		return Task{}, ErrTaskNotFound
	}

	if err := task.SetState(StatusProcessing); err != nil {
		return Task{}, err
	}

	result := m.processor.Process(task.Password)
	task.Result = result

	if err := task.SetState(StatusCompleted); err != nil {
		return Task{}, err
	}

	return *task, nil
}
