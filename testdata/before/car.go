package main

import "fmt"

type Car struct {
	Brand string
	Model string
	Year  int
}

func (me *Car) String() string {
	return fmt.Sprintf("%s %s %v", me.Brand, me.Model, me.Year)
}
