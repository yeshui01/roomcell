package tserver

const (
	EUserDataTypeNetInfo = 1
)

type UserData struct {
	dataType  int32
	NodeType  int32
	NodeIndex int32
	ZoneID    int32
}
