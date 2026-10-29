package main

import (
	"fmt"
)

type Task struct {
	Name string
	Run  func() error
}

func SafeRun(tasks []Task) []error {
	var errs []error
	for _, t := range tasks {
		err := func() (e error) {
			defer func() {
				if r := recover(); r != nil {
					e = fmt.Errorf("task %s: panic: %v", t.Name, r)
				}
			}()

			if err := t.Run(); err != nil {
				return fmt.Errorf("task %s: %w", t.Name, err)
			}
			return nil
		}()

		if err != nil {
			errs = append(errs, err)
		}
	}
	return errs
}
