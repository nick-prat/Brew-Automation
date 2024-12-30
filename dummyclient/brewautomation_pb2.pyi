from google.protobuf import empty_pb2 as _empty_pb2
from google.protobuf import descriptor as _descriptor
from google.protobuf import message as _message
from typing import ClassVar as _ClassVar, Optional as _Optional

DESCRIPTOR: _descriptor.FileDescriptor

class Command(_message.Message):
    __slots__ = ("messageCode",)
    MESSAGECODE_FIELD_NUMBER: _ClassVar[int]
    messageCode: int
    def __init__(self, messageCode: _Optional[int] = ...) -> None: ...

class TempLogRequest(_message.Message):
    __slots__ = ("temperature", "fermentRunId")
    TEMPERATURE_FIELD_NUMBER: _ClassVar[int]
    FERMENTRUNID_FIELD_NUMBER: _ClassVar[int]
    temperature: float
    fermentRunId: int
    def __init__(self, temperature: _Optional[float] = ..., fermentRunId: _Optional[int] = ...) -> None: ...

class TempLogResponse(_message.Message):
    __slots__ = ("id", "temperature", "fermentRunId", "timestamp")
    ID_FIELD_NUMBER: _ClassVar[int]
    TEMPERATURE_FIELD_NUMBER: _ClassVar[int]
    FERMENTRUNID_FIELD_NUMBER: _ClassVar[int]
    TIMESTAMP_FIELD_NUMBER: _ClassVar[int]
    id: int
    temperature: float
    fermentRunId: int
    timestamp: str
    def __init__(self, id: _Optional[int] = ..., temperature: _Optional[float] = ..., fermentRunId: _Optional[int] = ..., timestamp: _Optional[str] = ...) -> None: ...

class FermentRunGetRequest(_message.Message):
    __slots__ = ("id",)
    ID_FIELD_NUMBER: _ClassVar[int]
    id: int
    def __init__(self, id: _Optional[int] = ...) -> None: ...

class FermentRunCreateRequest(_message.Message):
    __slots__ = ("name",)
    NAME_FIELD_NUMBER: _ClassVar[int]
    name: str
    def __init__(self, name: _Optional[str] = ...) -> None: ...

class FermentRunResponse(_message.Message):
    __slots__ = ("id", "name")
    ID_FIELD_NUMBER: _ClassVar[int]
    NAME_FIELD_NUMBER: _ClassVar[int]
    id: int
    name: str
    def __init__(self, id: _Optional[int] = ..., name: _Optional[str] = ...) -> None: ...

class DeviceInstruction(_message.Message):
    __slots__ = ("code", "deviceId")
    CODE_FIELD_NUMBER: _ClassVar[int]
    DEVICEID_FIELD_NUMBER: _ClassVar[int]
    code: int
    deviceId: str
    def __init__(self, code: _Optional[int] = ..., deviceId: _Optional[str] = ...) -> None: ...
