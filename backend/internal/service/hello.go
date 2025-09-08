package service

type HelloService interface {
	Hello() string
}

type helloServiceImpl struct{}

func NewHelloService() HelloService {
	return &helloServiceImpl{}
}

func (s *helloServiceImpl) Hello() string {
	return "Hello, Golang!"
}
