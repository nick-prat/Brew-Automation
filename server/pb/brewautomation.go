package pb

import (
	context "context"
)

type Server struct {
	UnimplementedAPIServer
}

// SayHello implements helloworld.GreeterServer
func (s *Server) CreateTempLog(_ context.Context, in *TempLogRequest) (*TempLogRequest, error) {
	return &TempLogRequest{FermentRunId: in.FermentRunId, Temperature: in.Temperature}, nil
}
