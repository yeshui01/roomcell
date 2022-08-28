package appconfig

import (
	"path"

	"github.com/spf13/viper"
)

// room cell config
// mysql配置
type MysqlCfg struct {
	User   string `yaml:"user"`
	Host   string `yaml:"host"`
	Port   int32  `yaml:"port"`
	Pswd   string `yaml:"pswd"`
	DbName string `yaml:"dbName"`
}

// 账号配置
type AccountCfg struct {
	ListenAddr string `yaml:"listenAddr"`
	LogLevel   int32  `yaml:"logLevel"`
}

// 网关配置
type GateCfg struct {
	ListenAddr string `yaml:"listenAddr"`
	LogLevel   int32  `yaml:"logLevel"`
}

// room配置
type RoomCfg struct {
	ListenAddr string `yaml:"listenAddr"`
	LogLevel   int32  `yaml:"logLevel"`
}

// 大区节点配置
type ZoneNode struct {
	ZoneID   int32      `yaml:"zoneID"`
	GateCfgs []*GateCfg `yaml:"gateCfgs"`
	RoomCfgs []*RoomCfg `yaml:"roomCfgs"`
}

// 账号db
type AccountLocalDB struct {
	DbFile string `yaml:"dbFile"`
}

type RoomCellConfig struct {
	ZoneNodeCfgs   *ZoneNode       `yaml:"zoneNodeCfgs"`
	AccountCfgs    []*AccountCfg   `yaml:"accountCfgs"`
	ServerID       int32           `yaml:"serverID"`
	AccountLocalDB *AccountLocalDB `yaml:"accountLocalDB"`
	AccountDB      *MysqlCfg       `yaml:"accountDB"`
}

func NewRoomCellCfg() *RoomCellConfig {
	return &RoomCellConfig{}
}

//从本地文件读取配置
func ReadRoomCellConfigFromFile(configFilePath string, defaultSetting *RoomCellConfig) error {
	fullPath := path.Join(configFilePath, "roomcell.yaml")
	viper.SetConfigFile(fullPath)
	err := viper.ReadInConfig()
	if err != nil {
		return err
	}
	return viper.Unmarshal(defaultSetting)
}
