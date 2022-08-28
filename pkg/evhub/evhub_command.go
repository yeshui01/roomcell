/*
 * @Author: mknight(tianyh)
 * @Mail: 824338670@qq.com
 * @Date: 2022-03-02 15:00
 * @LastEditTime: 2022-03-02 15:00
 */
package evhub

const (
	HubCmdFrame    = 1 // 框架命令
	HubCmdNetEvent = 2 // 网络事件

	HubCmdUserBase = 1000 // 用户命令起始
)

// 	HubCmdNetEvent = 2 // 网络事件
const (
	NetEventConnected = 1 // 连接
	NetEventClose     = 2 // 关闭
	NetEventMessage   = 3 // 收到消息
	NetEventExitWrite = 4 // 退出写线程
)

type HubCommand struct {
	cmdClass int32 // 主命令
	cmdType  int32 // 子命令类型
	cmdData  interface{}
}

func (hcmd *HubCommand) GetCmdClass() int32 {
	return hcmd.cmdClass
}
func (hcmd *HubCommand) GetCmdType() int32 {
	return hcmd.cmdType
}

func (hcmd *HubCommand) GetCmdData() interface{} {
	return hcmd.cmdData
}
func NewHubCommand(cmd int32, cType int32, data interface{}) *HubCommand {
	return &HubCommand{
		cmdClass: cmd,
		cmdType:  cType,
		cmdData:  data,
	}
}
