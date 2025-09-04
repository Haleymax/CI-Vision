import grpc
import hello_pb2_grpc as pb2_grpc
import hello_pb2 as pb2
from concurrent import futures
from time import sleep


class Greeter(pb2_grpc.GreeterServicer):
    def SayHello(self, request, context):
        name = request.name
        age = request.age
        message = f"Hello, {name}! You are {age} years old."
        return pb2.HelloReply(message=message, age=age)
    

def serve():
    grpc_server = grpc.server(
        futures.ThreadPoolExecutor(max_workers=10)
    )
    pb2_grpc.add_GreeterServicer_to_server(Greeter(), grpc_server)
    grpc_server.add_insecure_port("0.0.0.0:50051")
    print("gRPC server is running on port 50051...")
    grpc_server.start()
    grpc_server.wait_for_termination()


if __name__ == "__main__":
    serve()