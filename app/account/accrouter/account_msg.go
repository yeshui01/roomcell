/*
 * @Author: mknight(tianyh)
 * @Mail: 824338670@qq.com
 * @Date: 2022-10-08 15:15:19
 * @LastEditTime: 2022-10-14 10:24:13
 * @FilePath: \roomcell\app\account\accrouter\account_msg.go
 */
package accrouter

// 注册
type AccountRegisterReq struct {
	UserName     string `json:"user_name"`
	Nickname     string `json:"nickname"`
	Pswd         string `json:"pswd"`
	ThirdPlat    string `json:"third_plat"`
	ThirdAccount string `json:"third_account"`
	Sign         string `json:"sign"`
	CltX         int64  `json:"cltx"`
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

// 公告
type QueryNoticeReq struct {
	Sign string `json:"sign"`
	CltX int64  `json:"cltx"`
}
type QueryNoticeRsp struct {
	ID     int32
	Notice string // 公告内容
}
