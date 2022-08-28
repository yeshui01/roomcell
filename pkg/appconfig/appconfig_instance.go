package appconfig

import (
	"path"

	"github.com/spf13/viper"
)

var roomcellCfg *RoomCellConfig

func Instance() *RoomCellConfig {
	if roomcellCfg == nil {
		roomcellCfg = NewRoomCellCfg()
	}
	return roomcellCfg
}

//从本地文件读取配置
func (cfg *RoomCellConfig) Load(configFilePath string) error {
	fullPath := path.Join(configFilePath, "roomcell.yaml")
	viper.SetConfigFile(fullPath)
	err := viper.ReadInConfig()
	if err != nil {
		return err
	}
	err = viper.Unmarshal(cfg)
	return err
}
