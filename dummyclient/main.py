from __future__ import print_function

import logging

import grpc
import brewautomation_pb2
import brewautomation_pb2_grpc


def run():
    with grpc.insecure_channel("localhost:50051") as channel:
        stub = brewautomation_pb2_grpc.APIStub(channel)
        response = stub.CreateTempLog(brewautomation_pb2.TempLogRequest(temperature=10.5, fermentRunId=0))
    print("Greeter client received: " + str(response.temperature))


if __name__ == "__main__":
    logging.basicConfig()
    run()
