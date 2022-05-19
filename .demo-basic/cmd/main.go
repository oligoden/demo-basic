package main

import (
	"github.com/oligoden/chassis/adapter"
	demobasic "github.com/oligoden/demo-basic/.demo-basic"
)

func main() {
	adapter.NewMux().
		SetURL("http://test.com:8080").
		AddRPD("swagger:8080").
		Register(demobasic.MuxFunc()).
		Serve()
}
