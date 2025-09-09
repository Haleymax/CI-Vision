package service

type HelloService interface {
	Hello() (string, error)
}

type helloServiceImpl struct{}

func NewHelloService() HelloService {
	return &helloServiceImpl{}
}

func (s *helloServiceImpl) Hello() (string, error) {
	return "Hello, Gin!", nil
}
