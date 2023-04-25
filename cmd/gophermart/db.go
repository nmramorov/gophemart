package main

type DbInterface interface {
	Connect()
	Save(interface{}) bool
	Get(interface{}) (interface{}, error)
	Update()
}

type Cursor struct {
	DbInterface
}
