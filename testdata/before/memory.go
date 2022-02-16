package main

import "fmt"

// Inmemory database
type Memory struct {
	cars []Car
}

func (me *Memory) Insert(v interface{}) error {
	switch v := v.(type) {
	case Car:
		me.cars = append(me.cars, v)
	default:
		return fmt.Errorf("Insert: cannot handle %T", v)
	}
	return nil
}

func (me *Memory) Select(v interface{}) error {
	switch v := v.(type) {
	case []Car:
		copy(v, me.cars)
	default:
		return fmt.Errorf("Select: cannot handle %T", v)
	}
	return nil

}
