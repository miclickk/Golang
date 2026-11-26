package main

import "fmt"

func Letter(score int) (string, error) {
	switch {
	case score < 0 || score > 100:
		return "", fmt.Errorf("score %d out of range 0..100", score)
	case score >= 90:
		return "A", nil
	case score >= 80:
		return "B", nil
	case score >= 70:
		return "C", nil
	case score >= 60:
		return "D", nil
	default:
		return "F", nil
	}
}

func Average(scores []int) float64 {
	if len(scores) == 0 {
		return 0
	}
	sum := 0
	for _, s := range scores {
		sum += s
	}
	return float64(sum) / float64(len(scores))
}

func main() {
	g, err := Letter(90)
	fmt.Println(g, err)
	fmt.Println(Average([]int{90, 80, 70}))
}
