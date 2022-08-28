package robotcore

// 谁是卧底-机器人
type CellRobotUndercover struct {
	*CellRobot
}

func NewRobotUndercover(robotName string, targetRoomID int64) *CellRobotUndercover {
	r := &CellRobotUndercover{
		CellRobot: NewCellRobot(robotName),
	}
	r.TargetRoomID = targetRoomID
	r.UpdateFunc = func(curTime int64) {
		r.Update(curTime)
	}
	r.CellRobot.RealRobotIns = r
	return r
}

func (robotObj *CellRobotUndercover) Update(curTime int64) {
	//robotObj.Debugf("undercover robot(%s) sec update,robot status(%d) step(%d)", robotObj.RobotName, robotObj.RobotStatus, robotObj.StatusStep)
	switch robotObj.RobotStatus {
	case ERobotStatusNone:
		{
			robotObj.LoginHall()
			break
		}
	case ERobotStatusLoginHall:
		{
			if robotObj.StatusStep == ERobotStatusStepIng {
				break
			}
			// 登录大厅完成
			if robotObj.TargetRoomID > 0 {
				robotObj.SendEnterRoom(robotObj.TargetRoomID)
			}
			break
		}
	case ERobotStatusEnterRoom:
		{
			if robotObj.StatusStep == ERobotStatusStepIng {
				break
			}
			if !robotObj.IsReady {
				robotObj.DoActGameReady(1)
			}
			break
		}
	}
}

// 玩家进入
