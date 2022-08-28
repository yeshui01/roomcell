package pbtools

import (
	"roomcell/pkg/loghlp"

	"google.golang.org/protobuf/proto"
	"google.golang.org/protobuf/reflect/protoreflect"
	"google.golang.org/protobuf/reflect/protoregistry"
)

// 获取消息的全名称
func GetFullNameByMessage(msg proto.Message) string {
	// reflectPB := proto.MessageReflect(msg)
	// descripterPB := reflectPB.Descripter()
	return string(proto.MessageName(msg))
}

// 根据消息的名称获取对象
func GetNewMessageObjByName(messageName string) proto.Message {
	pbMsgType, err := protoregistry.GlobalTypes.FindMessageByName(protoreflect.FullName(messageName))
	if err != nil {
		loghlp.Errorf("not find pbmessage name:%s", messageName)
		return nil
	}
	newMsg := pbMsgType.New().Interface()
	return proto.Message(newMsg)
}
