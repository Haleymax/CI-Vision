package setup

import "civ/internal/controller/hello"

type Controllers struct {
	HelloController hello.HelloController
}

// NewControllers creates and returns a Controllers instance with its HelloController
// field initialized.
//
// It instantiates a hello.HelloController via hello.NewHelloController and returns
// a pointer to the populated Controllers struct.
func NewControllers() *Controllers {

	HelloController := hello.NewHelloController()
	return &Controllers{
		HelloController: *HelloController,
	}
}
