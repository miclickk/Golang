package main

import "fmt"

func main() {
	manager := NewTaskManager(PasswordProcessor{})
	manager.Start(3, 5)

	passwords := []string{"weak", "Medium12", "Sup3r$ecret", "Another#99"}
	var ids []int

	for _, pwd := range passwords {
		task, err := manager.CreateTask(pwd)
		if err != nil {
			fmt.Println("Create error:", err)
			continue
		}
		if err := manager.SubmitTask(task.ID); err != nil {
			fmt.Println("Submit error:", err)
			continue
		}
		ids = append(ids, task.ID)
	}

	manager.Stop()

	for _, id := range ids {
		task, err := manager.GetTask(id)
		if err != nil {
			fmt.Println("Get error:", err)
			continue
		}
		fmt.Printf("Task #%d | State: %s | Strength: %s | Len: %d\n",
			task.ID, task.State, task.Result.Strength, task.Result.Length)
	}
}
