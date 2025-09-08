package hello

import (
	context "context"
	"log"
	"time"

	grpc "google.golang.org/grpc"
	"google.golang.org/grpc/credentials/insecure"
)

func demo() {
	conn, err := grpc.Dial("localhost:50051", grpc.WithTransportCredentials(insecure.NewCredentials()))
	if err != nil {
		log.Fatalf("did not connect: %v", err)
	}
	defer conn.Close()

	c := pb.NewGreeterClient(conn)

	ctx, cancel := context.WithTimeout(context.Background(), time.Second)
	defer cancel()

	r, err := c.SayHello(ctx, &pb.HelloRequest{Name: "Go Client", Age: 30})
	if err != nil {
		log.Fatalf("could not greet: %v", err)
	}
	log.Printf("Response: %s", r.Message)
	log.Printf("Original age: 30, Server processed age: %d", r.Age)
}
