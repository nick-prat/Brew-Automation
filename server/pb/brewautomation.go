package pb

import (
	context "context"
	"fmt"
	"raspberrysour/dao"
	"sync"

	"github.com/google/uuid"
	"github.com/jmoiron/sqlx"
	"google.golang.org/protobuf/types/known/emptypb"
)

type ChanKey = uuid.UUID

type Server struct {
	DB *sqlx.DB
	UnimplementedAPIServer
}

type Topic[T any] struct {
	Count    int64
	Channels map[ChanKey]chan *T
}

type TopicHandler struct {
	DeviceInstructionTopic Topic[DeviceInstruction]
}

var topicHandlerLock = &sync.Mutex{}
var topicHandler *TopicHandler

func getTopicHandler() *TopicHandler {
	if topicHandler == nil {
		topicHandlerLock.Lock()
		defer topicHandlerLock.Unlock()
		topicHandler = &TopicHandler{}
	}

	return topicHandler
}

func (topicHandler *TopicHandler) subscribeDeviceInstruction() (ChanKey, chan *DeviceInstruction) {
	topicHandlerLock.Lock()
	defer topicHandlerLock.Unlock()
	dit := &topicHandler.DeviceInstructionTopic

	if dit.Channels == nil {
		dit.Channels = make(map[uuid.UUID]chan *DeviceInstruction)
	}

	id := uuid.New()
	c := make(chan *DeviceInstruction)
	dit.Channels[id] = c
	dit.Count += 1

	return id, c
}

func (topicHandler *TopicHandler) unsubscribeDeviceInstruction(key ChanKey) error {
	topicHandlerLock.Lock()
	defer topicHandlerLock.Unlock()
	dit := &topicHandler.DeviceInstructionTopic

	if dit.Channels == nil {
		return fmt.Errorf("channel map has not been initiated")
	}

	c, ok := dit.Channels[key]
	if !ok {
		return fmt.Errorf("key doesn't exist in channels")
	}
	close(c)
	delete(dit.Channels, key)

	return nil
}

func (topicHandler *TopicHandler) publishDeviceInstruction(msg *DeviceInstruction) {
	dit := &topicHandler.DeviceInstructionTopic

	for _, c := range dit.Channels {
		c <- msg
	}
}

func (s *Server) CreateTempLog(_ context.Context, in *TempLogRequest) (*TempLogResponse, error) {
	tempLogDAO := dao.NewTempLogDAO(s.DB)

	pk, err := tempLogDAO.Create(&dao.TempLog{
		Temp:       in.Temperature,
		FermentRun: in.FermentRunId,
	})
	if err != nil {
		return nil, err
	}

	tempLog, err := tempLogDAO.Get(pk)
	if err != nil {
		return nil, err
	}

	return &TempLogResponse{
		FermentRunId: tempLog.FermentRun,
		Temperature:  tempLog.Temp,
		Id:           tempLog.Id,
		Timestamp:    tempLog.TimeStamp.String(),
	}, nil
}

func getFermentRun(fermentRunDAO *dao.FermentRunDAO, id int32) (*FermentRunResponse, error) {
	fermentRun, err := fermentRunDAO.Get(id)
	if err != nil {
		return nil, err
	}

	return &FermentRunResponse{
		Id:   fermentRun.Id,
		Name: fermentRun.Name,
	}, nil
}

func (s *Server) GetFermentRun(_ context.Context, in *FermentRunGetRequest) (*FermentRunResponse, error) {
	return getFermentRun(dao.NewFermentRunDAO(s.DB), in.Id)
}

func (s *Server) CreateFermentRun(_ context.Context, in *FermentRunCreateRequest) (*FermentRunResponse, error) {
	fermentRunDAO := dao.NewFermentRunDAO(s.DB)

	pk, err := fermentRunDAO.Create(&dao.FermentRun{
		Name: in.Name,
	})
	if err != nil {
		return nil, err
	}

	return getFermentRun(fermentRunDAO, pk)
}

func (s *Server) SubscribeDeviceInstruction(_ *emptypb.Empty, stream API_SubscribeDeviceInstructionServer) error {
	key, c := getTopicHandler().subscribeDeviceInstruction()
	defer getTopicHandler().unsubscribeDeviceInstruction(key)
	defer func() {
		fmt.Println("Client Disconnected")
	}()

	fmt.Println("Client Connected")

	for deviceInstruction := range c {
		if err := stream.Send(deviceInstruction); err != nil {
			return err
		}
	}

	return nil
}

func (s *Server) PublishDeviceInstruction(_ context.Context, in *DeviceInstruction) (*emptypb.Empty, error) {
	getTopicHandler().publishDeviceInstruction(in)
	fmt.Printf("Published device instruction %d: %s\n", in.Code, in.DeviceId)
	return nil, nil
}
