package appconfig

import (
	"path"

	"github.com/spf13/viper"
)

type CellClientCfg struct {
	AccountAddr  string `yaml:"accountAddr"`
	LogLevel     int32  `yaml:"logLevel"`
	WindowWidth  int32  `yaml:"windowWidth"`
	WindowHeight int32  `yaml:"windowWHeight"`
}

func NewCellClientCfg() *CellClientCfg {
	return &CellClientCfg{
		AccountAddr: "",
		LogLevel:    5,
	}
}

//从本地文件读取配置
func ReadCellClientConfigFromFile(configFilePath string) (*CellClientCfg, error) {
	defaultSetting := NewCellClientCfg()
	fullPath := path.Join(configFilePath, "cell_client_config.yaml")
	viper.SetConfigFile(fullPath)
	err := viper.ReadInConfig()
	if err != nil {
		return nil, err
	}
	err = viper.Unmarshal(defaultSetting)
	return defaultSetting, err
}
