package utils

import (
	pb "github.com/Norzuiso/protocol/gen/go/proto/orchestrator/v1"
)

func BuildPhaseErrorMsg(msg *pb.Message, err error) *pb.Message {
	return msg
}
