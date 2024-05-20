package main

import (
	"fmt"
	"math"
)

type Shape struct {
	Name string
}

func (s Shape) GetName() string {
	return s.Name
}

func (s Shape) Area() float64 {
	return 0.0
}

type Rectangle struct {
	Shape
	Width  float64
	Height float64
}

func (r Rectangle) Area() float64 {
	return r.Width * r.Height
}

type Circle struct {
	Shape
	Radius float64
}

func (c Circle) Area() float64 {
	return math.Pi * c.Radius * c.Radius
}

func main() {
	circle := Circle{Shape{"Circle"}, 5}
	rectangle := Rectangle{Shape{"Rectangle"}, 4, 6}

	fmt.Printf("%s: Area = %.2f\n",
		circle.GetName(), circle.Area())
	fmt.Printf("%s: Area = %.2f\n",
		rectangle.GetName(), rectangle.Area())
}
