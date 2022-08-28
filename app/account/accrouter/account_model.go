package accrouter

type OrmUser struct {
	UserID       int64  `gorm:"column:user_id;primaryKey"`
	UserName     string `gorm:"column:user_name"`
	Pswd         string `gorm:"column:pswd"`
	RegisterTime int64  `gorm:"column:register_time"`
	Status       int32  `gorm:"column:status"`
	DataZone     int32  `gorm:"column:data_zone"`
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
