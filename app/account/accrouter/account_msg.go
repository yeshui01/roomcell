package accrouter

// 注册
type AccountRegisterReq struct {
	UserName string `json:"user_name"`
	Pswd     string `json:"pswd"`
	Sign     string `json:"sign"`
	CltX     int64  `json:"cltx"`
}

type AccountRegisterRsp struct {
	UserID int64 `json:"user_id"`
}

// 登录
type AccountLoginReq struct {
	UserName string `json:"user_name"`
	Pswd     string `json:"pswd"`
	Sign     string `json:"sign"`
	CltX     int64  `json:"cltx"`
}
type AccountLoginRsp struct {
	HallAddr string `json:"hall_addr"` // 大厅地址
	Token    string `json:"token"`     // 返回token信息
	Statu    int32  `json:"status"`    // 当前账号状态 0-未验证 1-认证通过
	RestTime int32  `json:"rest_time"` // 剩余的认证时间,为0则不可用
}
