package main

type Database interface {
	Insert(interface{}) error
	Select(interface{}) error
}

