package main

import (
	"errors"
	"fmt"
)

type Celsius float64

func (c Celsius) String() string { return fmt.Sprintf("%.1f°C", float64(c)) }

func describe(v any) string {
	switch x := v.(type) {
	case nil:
		return "nil"
	case int, int64:
		return fmt.Sprintf("integer %v", x)
	case string:
		return fmt.Sprintf("string of %d bytes", len(x))
	case []int:
		return fmt.Sprintf("slice with %d ints", len(x))
	case error:
		return "error: " + x.Error()
	case fmt.Stringer:
		return "Stringer: " + x.String()
	case bool:
		return fmt.Sprintf("boolean %v", x)
	case map[string]int:
		return fmt.Sprintf("map with %d keys", len(x))
	default:
		return fmt.Sprintf("unknown type %T", x)
	}
}

func describeNoSwitch(v any) string {
	if v == nil {
		return "nil"
	}
	if x, ok := v.(int); ok {
		return fmt.Sprintf("integer %v", x)
	}
	if x, ok := v.(int64); ok {
		return fmt.Sprintf("integer %v", x)
	}
	if x, ok := v.(string); ok {
		return fmt.Sprintf("string of %d bytes", len(x))
	}
	// Остальные проверки по аналогии...
	return fmt.Sprintf("unknown type %T", v)
}

func main() {
	values := []any{
		nil,
		7,
		[]int{1, 2, 3},
		errors.New("disk full"),
		Celsius(21.5),
		int64(7),
		"привет",
		3.14,
		struct{}{},
		true,
		map[string]int{"a": 1},
	}

	for _, val := range values {
		fmt.Printf("%-25s -> %s\n", fmt.Sprintf("%#v", val), describe(val))
	}
}
