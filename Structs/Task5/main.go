package main

import "testing"

type Report struct {
	Title  string
	Values [1024]float64
}

func (r Report) FirstValue() float64 { return r.Values[0] }

func (r *Report) FirstPointer() float64 { return r.Values[0] }

type Point struct{ X, Y float64 }

func (p Point) FirstValue() float64 { return p.X }

func (p *Point) FirstPointer() float64 { return p.X }

func BenchmarkReportValue(b *testing.B) {
	r := Report{}
	for i := 0; i < b.N; i++ {
		r.FirstValue()
	}
}

func BenchmarkReportPointer(b *testing.B) {
	r := Report{}
	for i := 0; i < b.N; i++ {
		r.FirstPointer()
	}
}

func BenchmarkPointValue(b *testing.B) {
	p := Point{}
	for i := 0; i < b.N; i++ {
		p.FirstValue()
	}
}

func BenchmarkPointPointer(b *testing.B) {
	p := Point{}
	for i := 0; i < b.N; i++ {
		p.FirstPointer()
	}
}
