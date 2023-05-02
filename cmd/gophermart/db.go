package main

type DbInterface interface {
	Connect()
	Save(interface{}) bool
	Get(interface{}) (interface{}, error)
	Update()
	SaveSession(string, interface{})
	GetOrder(interface{}) (interface{}, error)
	SaveOrder(interface{})
}

type Cursor struct {
	DbInterface
}
