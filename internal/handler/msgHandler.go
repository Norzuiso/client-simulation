package handler

import (
	"context"
	"fmt"
	"log"

	"github.com/Norzuiso/client/internal/experiment"
	pb "github.com/Norzuiso/protocol/gen/go/proto/orchestrator/v1"
	"google.golang.org/protobuf/proto"
	"google.golang.org/protobuf/types/known/anypb"
	"google.golang.org/protobuf/types/known/wrapperspb"
)

type Connection struct {
	Stream pb.Broadcast_ClientToClientMessageClient
	Outbox chan *pb.Message
	ErrCh  chan error
}

type MsgHandler struct {
	ctx             context.Context
	ClientBroadcast pb.BroadcastClient

	client     *pb.Client
	Connection *Connection

	// Experiment
	student experiment.Student
}

func NewMsgHandler(ctx context.Context, clientBroadcast pb.BroadcastClient) *MsgHandler {
	return &MsgHandler{
		ctx:             ctx,
		ClientBroadcast: clientBroadcast,
	}
}

func (m *MsgHandler) Init(servName string, initState string, id int64) {
	m.student = experiment.Student{InitState: initState}
	m.student.Initialize()

	m.client = &pb.Client{
		Id:   id,
		Name: servName,
		Seed: 0,
	}
	m.Connection = &Connection{
		ErrCh:  make(chan error, 1),
		Outbox: make(chan *pb.Message, 10),
	}
	m.ConnectClient()
	m.StartCommunication()
}

func (m *MsgHandler) StartCommunication() {
	var err error
	m.Connection.Stream, err = m.ClientBroadcast.ClientToClientMessage(m.ctx)
	if err != nil {
		log.Fatal(err)
	}
	err = m.OpenStream()
	if err != nil {
		log.Fatal(err)
	}
	latestSendMsg := &pb.Message{}
	go func() {
		for {
			msg, err := m.Connection.Stream.Recv()
			if err != nil {
				m.Connection.ErrCh <- err
				return
			}
			log.Println("Read: ", msg)

			msgType := msg.MessageType
			m.student.Seed = msg.Seed

			switch msgType {
			case pb.MessageType_MESSAGE_TYPE_DEFAULT:
				if msg.GetContent() == "End of simulation" {
					m.student.Initialize()
				}

			case pb.MessageType_MESSAGE_TYPE_REQUEST_EVENT:
				m.student.DeltInt()
				response := m.student.Lambda()
				latestSendMsg = m.BuildResEventMsg(response, msg)
				log.Println("Send: ", latestSendMsg)
				m.Connection.Stream.Send(latestSendMsg)

			case pb.MessageType_MESSAGE_TYPE_EVENT_DISPATCH:
				attr := msg.Attributes["influence"]
				result, err := getFloatFromAny(attr)
				if err != nil {
					log.Fatalf("Failed to create Any: %v", err)
				}
				m.student.CurrentState.AddInfluenceFromState(msg.Content, result)

			case pb.MessageType_MESSAGE_TYPE_ERROR_PHASE:
				log.Println("Retry Send: ", latestSendMsg)
				m.Connection.Stream.Send(latestSendMsg)
			}
		}
	}()
	log.Println(<-m.Connection.ErrCh)
}

func getFloatFromAny(anyField *anypb.Any) (float64, error) {
	var doubleMsg wrapperspb.DoubleValue

	if err := anypb.UnmarshalTo(anyField, &doubleMsg, proto.UnmarshalOptions{}); err != nil {
		return 0, fmt.Errorf("failed to unmarshal any: %w", err)
	}

	return doubleMsg.GetValue(), nil
}

func (m *MsgHandler) OpenStream() error {
	openStreamMsg := &pb.Message{
		SenderId:    m.client.GetId(),
		Content:     "",
		MessageType: pb.MessageType_MESSAGE_TYPE_OPEN_STREAM,
		Epoch:       0,
		Attributes:  nil,
		Seed:        m.client.GetSeed(),
	}
	m.Connection.Stream.Send(openStreamMsg)
	msgRes, err := m.Connection.Stream.Recv()
	if err != nil {
		log.Fatal(err)
		return err
	}
	log.Println(msgRes.String())

	if msgRes.MessageType == pb.MessageType_MESSAGE_TYPE_ERROR {
		return fmt.Errorf("%s", msgRes.GetContent())
	}

	return nil
}

func (m *MsgHandler) ConnectClient() {
	connResponse, err := m.ClientBroadcast.ConnectClient(m.ctx, &pb.ConnectionRequest{Client: m.client})
	if err != nil {
		log.Fatal(err)
	}
	m.client = connResponse.GetClient()
}

func (m *MsgHandler) ReadMsg(msg *pb.Message) error {
	msgType := msg.MessageType

	switch msgType {

	case pb.MessageType_MESSAGE_TYPE_ERROR:
		return fmt.Errorf("%s", msg.GetContent())

	case pb.MessageType_MESSAGE_TYPE_REQUEST_EVENT:
		response := m.student.Lambda()
		m.Connection.Outbox <- m.BuildResEventMsg(response, msg)

	case pb.MessageType_MESSAGE_TYPE_EVENT_DISPATCH:
		m.student.ExtEvent <- msg.GetContent()

	case pb.MessageType_MESSAGE_TYPE_REQUEST_CLIENT_STATUS:
		log.Println(msg.Content)

	case pb.MessageType_MESSAGE_TYPE_ERROR_PHASE:
		log.Println(msg.Content)
	}

	m.student.DeltExt(0)
	m.student.DeltInt()
	return nil
}

func (m *MsgHandler) SendMsg(msg *pb.Message) {
	m.Connection.Outbox <- msg
}

func (m *MsgHandler) BuildResEventMsg(content string, msg *pb.Message) *pb.Message {
	msg.Content = content
	msg.SenderId = m.client.Id
	msg.MessageType = pb.MessageType_MESSAGE_TYPE_EVENT_RESPONSE
	return msg
}
