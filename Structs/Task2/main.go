package main

import "fmt"

type Person struct {
	Name  string
	Email string
}

func (p Person) Contact() string {
	return p.Name + " <" + p.Email + ">"
}

type Student struct {
	Person
	GPA float64
}

type Teacher struct {
	Person
	Department string
}

func (t Teacher) Contact() string {
	return t.Person.Contact() + " (Dept: " + t.Department + ")"
}

type TypeA struct{ ID int }
type TypeB struct{ ID int }
type Combined struct {
	TypeA
	TypeB
}

func main() {
	s := Student{Person: Person{Name: "Dias", Email: "dias@astanait.edu.kz"}, GPA: 3.9}
	t := Teacher{Person: Person{Name: "Ali", Email: "ali@astanait.edu.kz"}, Department: "CS"}

	fmt.Println(s.Contact())
	fmt.Println(t.Contact())
	fmt.Println(t.Person.Contact())
	c := Combined{TypeA: TypeA{ID: 1}, TypeB: TypeB{ID: 2}}
	fmt.Println("ID из TypeA:", c.TypeA.ID)
	fmt.Println("ID из TypeB:", c.TypeB.ID)
}
