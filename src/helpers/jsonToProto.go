package helpers

import (
	"google.golang.org/protobuf/encoding/protojson"
	"google.golang.org/protobuf/proto"
)

func ProtobufToJSON(message proto.Message) (string, error) {
	b, err := protojson.MarshalOptions{
		Indent:          "  ",
		UseProtoNames:   true,
		EmitUnpopulated: true,
	}.Marshal(message)
	if err != nil {
		return "", err
	}
	return string(b), err
}

func JSONToProtobuf(message proto.Message, json string) error {
	return protojson.UnmarshalOptions{
		AllowPartial: true,
	}.Unmarshal([]byte(json), message)
}
