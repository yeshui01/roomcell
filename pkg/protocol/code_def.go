package protocol

const (
	ECodeSuccess     = 0 // 正常
	ECodeSysError    = 1 // 系统错误
	ECodeParamError  = 2 // 参数错误
	ECodeAsyncHandle = 3 // 异步处理,当前不处理
	// 系统错误定义
	ECodeDBError           = 100 // db错误
	ECodePBDecodeError     = 101 // pb反序列化错误
	ECodeTokenError        = 102 // token解析失败
	ECodeTokenExpire       = 103 // token过期
	ECodeInvalideOperation = 104 // 无效操作
	// 错误码定义
	ECodeAccNameHasExisted       = 1000 // 账号已经存在
	ECodeAccNotExisted           = 1001 // 账号不存在
	ECodeAccCertificationTimeOut = 1002 // 账号超出最终认证时间
	ECodeAccPasswordError        = 1003 // 账号密码错误
	ECodeRoleNotExisted          = 1004 // 角色不存在
	ECodeRoleHasOnline           = 1005 // 玩家已经在线
	ECodeNotFindNotice           = 1106 // 未找到公告
	ECodeRoomMaxPlayerNumLimit   = 1107 // 房间最大人数限制
	// 房间
	ECodeRoomCreateFail                  = 1100 // 房间创建失败
	ECodeRoomPlayerHasInRoom             = 1101 // 玩家已经在房间
	ECodeRoomNotExisted                  = 1102 // 房间不存在
	ECodeRoomPlayerNotFound              = 1103 // 找不到房间玩家
	ECodeRoomPlayerNotInRoom             = 1104 // 玩家不在房间
	ECodeRoomDrawSettingAuthError        = 1105 // 你画我猜房间设定-权限不足,只有房间管理员可以设定
	ECodeRoomDrawReapeatedGuess          = 1106 // 你画我猜-不能重复猜词
	ECodeRoomDrawInvalideOption          = 1107 // 你画我猜-无效的操作
	ECodeRoomCantjoin                    = 1108 // 房间游戏中,此时不能加入房间
	ECodeRoomUndercoverInvalideOption    = 1109 // 谁是卧底-无效的操作
	ECodeRoomUndercoverVoted             = 1110 // 谁是卧底-已经投票过了
	ECodeRoomUndercoverOutPlayerCantTalk = 1111 // 谁是卧底-已经出局的玩家不能发言

	ECodeRoomNumberbombInvalideOption = 1112 // 数字炸弹-无效的操作
	ECodeRoomNumberbombHasTalked      = 1113 // 数字炸弹-已经发言了
)
