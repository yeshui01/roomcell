/*
 * @Author: mknight(tianyh)
 * @Mail: 824338670@qq.com
 * @Date: 2022-10-08 15:15:19
 * @LastEditTime: 2022-10-14 10:25:19
 * @FilePath: \roomcell\app\account\accrouter\account_model.go
 */
package accrouter

type OrmUser struct {
	UserID       int64  `gorm:"column:user_id;primaryKey"`
	UserName     string `gorm:"column:user_name"`
	Nickname     string `gorm:"column:nickname"`
	Pswd         string `gorm:"column:pswd"`
	RegisterTime int64  `gorm:"column:register_time"`
	Status       int32  `gorm:"column:status"`
	DataZone     int32  `gorm:"column:data_zone"`
	ThirdPlat    string `gorm:"column:third_plat"`
	ThirdAccount string `gorm:"column:third_account"`
}

func (tbuser *OrmUser) TableName() string {
	return "user"
}

type OrmHallList struct {
	ID        int32  `gorm:"column:id;primaryKey"`
	GateAddr  string `gorm:"column:gate_addr"`
	Recommend int32  `gorm:"column:recommend"`
}

func (tb *OrmHallList) TableName() string {
	return "hall_list"
}

type OrmCellNotice struct {
	ID      int32  `gorm:"column:id;primaryKey"`
	Content string `gorm:"column:content"`
	UpdTime int64  `gorm:"column:upd_time"`
}

func (tb *OrmCellNotice) TableName() string {
	return "cell_notice"
}
