from google.protobuf import descriptor as _descriptor
from google.protobuf import message as _message
from typing import ClassVar as _ClassVar, Optional as _Optional

DESCRIPTOR: _descriptor.FileDescriptor

class TempLogRequest(_message.Message):
    __slots__ = ("temperature", "fermentRunId")
    TEMPERATURE_FIELD_NUMBER: _ClassVar[int]
    FERMENTRUNID_FIELD_NUMBER: _ClassVar[int]
    temperature: float
    fermentRunId: int
    def __init__(self, temperature: _Optional[float] = ..., fermentRunId: _Optional[int] = ...) -> None: ...
