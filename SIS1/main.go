package main

import "fmt"

func main() {
	processor := PasswordProcessor{}
	manager := NewTaskManager(processor)

	task1, err := manager.CreateTask("Rapper123$")
	if err != nil {
		fmt.Println("Error:", err)
		return
	}
	fmt.Printf("Created Task #%d | State: %s | Completed: %v\n", task1.ID, task1.State, task1.IsCompleted())

	processedTask1, err := manager.ProcessTask(task1.ID)
	if err != nil {
		fmt.Println("Error:", err)
		return
	}
	fmt.Printf("Processed Task #%d | State: %s | Strength: %s | Length: %d | Digits: %d | Upper: %d | Special: %d\n",
		processedTask1.ID,
		processedTask1.State,
		processedTask1.Result.Strength,
		processedTask1.Result.Length,
		processedTask1.Result.Digits,
		processedTask1.Result.Upper,
		processedTask1.Result.Special,
	)

	_, err = manager.CreateTask("bad password")
	if err != nil {
		fmt.Println("Expected validation error:", err)
	}

	_, err = manager.ProcessTask(task1.ID)
	if err != nil {
		fmt.Println("Expected state error:", err)
	}
}
