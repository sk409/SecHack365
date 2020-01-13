package main

type facade interface {
	Public() (interface{}, error)
}
