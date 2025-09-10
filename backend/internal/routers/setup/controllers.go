package setup

import "civ/internal/controller/hello"

type Controllers struct {
	HelloController hello.HelloController
}

func NewControllers() *Controllers {

	HelloController := hello.NewHelloController()
	return &Controllers{
		HelloController: *HelloController,
	}
}
