import grpc
import hello_pb2_grpc as pb2_grpc
import hello_pb2 as pb2

def run():
    conn = grpc.insecure_channel("localhost:50051")
    client = pb2_grpc.GreeterStub(channel=conn)
    response = client.SayHello(pb2.HelloRequest(name="Alice", age=30))
    print("Greeter client received: " + response.message)

if __name__ == "__main__":
    run()