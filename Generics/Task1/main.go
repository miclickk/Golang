package main

import (
	"errors"
	"fmt"
	"io"
	"os"
	"strings"
)

type Notifier interface {
	Notify(msg string) error
}

type Email struct{ Address string }

func (e Email) Notify(msg string) error {
	fmt.Printf("Email sent to %s: %s\n", e.Address, msg)
	return nil
}

type SMS struct{ Phone string }

func (s SMS) Notify(msg string) error {
	if len(msg) > 20 {
		return errors.New("message exceeds 20 characters")
	}
	fmt.Printf("SMS sent to %s: %s\n", s.Phone, msg)
	return nil
}

type Telegram struct {
	Chat     string
	CountMsg int
}

func (t *Telegram) Notify(msg string) error {
	t.CountMsg++
	fmt.Printf("Telegram sent to %s (total %d): %s\n", t.Chat, t.CountMsg, msg)
	return nil
}

var _ Notifier = Email{}
var _ Notifier = SMS{}
var _ Notifier = (*Telegram)(nil)

func Broadcast(ns []Notifier, msg string) []error {
	var errs []error
	for _, n := range ns {
		if err := n.Notify(msg); err != nil {
			errs = append(errs, err)
		}
	}
	return errs
}

func writeReport(w io.Writer, name string, score int) {
	fmt.Fprintf(w, "%s: %d points\n", name, score)
}

func main() {
	ns := []Notifier{
		Email{Address: "test@mail.com"},
		SMS{Phone: "+123456789"},
		&Telegram{Chat: "@durov"},
	}

	errs := Broadcast(ns, "Hello!")
	fmt.Println("Broadcast errors:", errs)

	fmt.Println("\n--- Grades ---")
	writeReport(os.Stdout, "Aruzhan", 95)

	var sb strings.Builder
	writeReport(&sb, "Dias", 80)
	fmt.Print("From builder: ", sb.String())
}
