from __future__ import print_function

import logging

import grpc
import brewautomation_pb2
import brewautomation_pb2_grpc
import google.protobuf.empty_pb2


def run():
    with grpc.insecure_channel("localhost:50051") as channel:
        stub = brewautomation_pb2_grpc.APIStub(channel)
        for di in stub.SubscribeDeviceInstruction(google.protobuf.empty_pb2.Empty()):
            print("Greeter client received:\n" + str(di))


if __name__ == "__main__":
    logging.basicConfig()
    run()
