// --- this file generated by tools, don't edit it!!!! ---
package ormdef

type RoleBase struct {
	RoleID      int64  `gorm:"primary_key"`
	UserID      int64  `gorm:"column:user_id"`
	RoleName    string `gorm:"column:role_name"`
	CreateTime  int64  `gorm:"column:create_time"`
	Level       int32  `gorm:"column:level"`
	LoginTime   int64  `gorm:"column:login_time"`
	OfflineTime int64  `gorm:"column:offline_time"`
	Money       int64  `gorm:"column:money"`
}

func (t *RoleBase) TableName() string {
	return "role_base"
}
