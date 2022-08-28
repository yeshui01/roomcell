package hallroomhandler

import "roomcell/app/hallroom/ihallroom"

var (
	roomServe ihallroom.IHallRoom
)

func InitRoomObj(s ihallroom.IHallRoom) {
	roomServe = s
}
