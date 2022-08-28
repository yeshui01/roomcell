package tserver

type tserver_step_t int32

const (
	EServerRunStepCheck     tserver_step_t = 0 // 运行检测
	EServerRunStepInit      tserver_step_t = 1 // 初始化
	EServerRunStepPreRun    tserver_step_t = 2 // 准备运行
	EServerRunStepNormalRun tserver_step_t = 3 // 正常运行
	EServerRunStepStop      tserver_step_t = 4 // 停止
	EServerRunStepEnd       tserver_step_t = 5 // 结束
	EServerRunStepExit      tserver_step_t = 6 // 退出
)
