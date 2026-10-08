package main

import "fmt"

type Student struct {
	Name    string
	ID      int
	GPA     float64
	Courses []string
}

func addCourseValue(s Student, course string) {
	s.Courses = append(s.Courses, course)
}

func addCoursePointer(s *Student, course string) {
	s.Courses = append(s.Courses, course)
}

func main() {
	var s1 Student
	s2 := Student{Name: "Maxim", ID: 1, GPA: 3.8, Courses: []string{"Math"}}
	s3 := Student{"Nikita", 2, 3.5, []string{"Physics"}}
	s4 := &Student{Name: "Kolyan", ID: 3, GPA: 4.0}

	fmt.Printf("s1: %+v\n", s1)
	fmt.Printf("s2: %+v\n", s2)
	fmt.Printf("s3: %+v\n", s3)
	fmt.Printf("s4: %+v\n\n", s4)

	addCourseValue(s2, "History")
	fmt.Printf("После addCourseValue: %+v\n", s2)

	addCoursePointer(&s2, "History")
	fmt.Printf("После addCoursePointer: %+v\n", s2)
}
