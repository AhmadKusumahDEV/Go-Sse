package main

import "fmt"

// Shape interface
type Shape interface {
	Area() float64
	Perimeter() float64
}

// Circle struct implementing Shape
type Circle struct {
	radius float64
}

func (c Circle) Area() float64 {
	return 3.14 * c.radius * c.radius
}

func (c Circle) Perimeter() float64 {
	return 2 * 3.14 * c.radius
}

// Rectangle struct implementing Shape
type Rectangle struct {
	width, height float64
}

func (r Rectangle) Area() float64 {
	return r.width * r.height
}

func (r Rectangle) Perimeter() float64 {
	return 2 * (r.width + r.height)
}

func main() {
	c := Circle{radius: 5}
	r := Rectangle{width: 3, height: 4}

	shapes := []Shape{c, r}

	for _, shape := range shapes {
		fmt.Printf("Area: %f, Perimeter: %f\n", shape.Area(), shape.Perimeter())
	}
}
