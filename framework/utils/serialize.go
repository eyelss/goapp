package utils

import (
	"google.golang.org/protobuf/proto"
)

type Serializer interface {
	Encode(interface{}) ([]byte, error)
	Decode([]byte, interface{}) error
	ContentType() string
}

type ProtoSerializer struct{}
type ErrNotProtobufMessage struct{}

func (m *ErrNotProtobufMessage) Error() string {
	return "message is not protobuf message"
}

func (p ProtoSerializer) Encode(v interface{}) ([]byte, error) {
	if msg, ok := v.(proto.Message); ok {
		return proto.Marshal(msg)
	}

	return nil, &ErrNotProtobufMessage{}
}

func (p ProtoSerializer) Decode(data []byte, v interface{}) error {
	if msg, ok := v.(proto.Message); ok {
		return proto.Unmarshal(data, msg)
	}
	return &ErrNotProtobufMessage{}
}

func (p ProtoSerializer) ContentType() string { return "application/protobuf" }
