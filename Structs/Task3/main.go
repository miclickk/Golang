package main

type Counter struct{ n int }

func (c Counter) Get() int { return c.n }
func (c *Counter) Inc()    { c.n++ }
func NewCounter() Counter  { return Counter{} }

func main() {
	var c Counter
	c.Inc()
	(&c).Inc()
	Counter{}.Get()

	m := map[string]Counter{"a": {}}
	m["a"].Get()

	s := []Counter{{}}
	s[0].Inc()
	pm := map[string]*Counter{"a": {}}
	pm["a"].Inc()
}
