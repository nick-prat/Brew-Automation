from __future__ import print_function

import logging

import grpc
import brewautomation_pb2
import brewautomation_pb2_grpc


def run():
    with grpc.insecure_channel("localhost:50051") as channel:
        stub = brewautomation_pb2_grpc.APIStub(channel)
        response = stub.PublishDeviceInstruction(brewautomation_pb2.DeviceInstruction(code=1, deviceId="abc"))
    print("Greeter client received:\n" + str(response))


if __name__ == "__main__":
    logging.basicConfig()
    run()
