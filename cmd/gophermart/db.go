package main

type DbInterface interface {
	Connect()
	Save(interface{}) bool
	Get()
	Update()
}

type Cursor struct {
	DbInterface
}
