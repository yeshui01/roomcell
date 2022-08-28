package tframeconfig

import (
	"path"

	"github.com/spf13/viper"
)

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
	ListenAddr string    `yaml:"listenAddr"`
	LogLevel   int32     `yaml:"logLevel"`
	LogPath    string    `yaml:"logPath"`
	AccountDb  *MysqlCfg `yaml:"accountDb"`
}

// 网关配置
type GateCfg struct {
	ListenAddr string `yaml:"listenAddr"`
	LogLevel   int32  `yaml:"logLevel"`
	ListenMode string `yaml:"listenMode"`
	LogPath    string `yaml:"logPath"`
}

// root配置
type RootCfg struct {
	ListenAddr string `yaml:"listenAddr"`
	LogLevel   int32  `yaml:"logLevel"`
	ListenMode string `yaml:"listenMode"`
	LogPath    string `yaml:"logPath"`
}

// hallgate配置
type HallGateCfg struct {
	ListenAddr   string `yaml:"listenAddr"`
	LogLevel     int32  `yaml:"logLevel"`
	ListenMode   string `yaml:"listenMode"`
	WsListenAddr string `yaml:"wsListenAddr"`
	LogPath      string `yaml:"logPath"`
}

// hallserver
type HallServerCfg struct {
	ListenAddr string `yaml:"listenAddr"`
	LogLevel   int32  `yaml:"logLevel"`
	ListenMode string `yaml:"listenMode"`
	LogPath    string `yaml:"logPath"`
}

// halldata
type HallDataCfg struct {
	ListenAddr  string    `yaml:"listenAddr"`
	LogLevel    int32     `yaml:"logLevel"`
	ListenMode  string    `yaml:"listenMode"`
	LocalDBFile string    `yaml:"localDBFile"`
	LogPath     string    `yaml:"logPath"`
	GameDb      *MysqlCfg `yaml:"gameDb"`
}

// hallroom
type HallRoomCfg struct {
	ListenAddr string `yaml:"listenAddr"`
	LogLevel   int32  `yaml:"logLevel"`
	ListenMode string `yaml:"listenMode"`
	LogPath    string `yaml:"logPath"`
	CsvPath    string `yaml:"csvPath"`
}

// hallroommgr
type HallRoomMgrCfg struct {
	ListenAddr string `yaml:"listenAddr"`
	LogLevel   int32  `yaml:"logLevel"`
	ListenMode string `yaml:"listenMode"`
	LogPath    string `yaml:"logPath"`
}

type FrameConfig struct {
	AccountCfgs     []*AccountCfg     `yaml:"accountCfgs"`
	ZoneID          int32             `yaml:"zoneID"`
	GateCfgs        []*GateCfg        `yaml:"gateCfgs"`
	RootCfgs        []*RootCfg        `yaml:"rootCfgs"`
	HallGateCfgs    []*HallGateCfg    `yaml:"hallgateCfgs"`
	HallServerCfgs  []*HallServerCfg  `yaml:"hallserverCfgs"`
	HallDataCfgs    []*HallDataCfg    `yaml:"halldataCfgs"`
	HallRoomCfgs    []*HallRoomCfg    `yaml:"hallroomCfgs"`
	HallRoomMgrCfgs []*HallRoomMgrCfg `yaml:"hallroomMgrCfgs"`
}

func NewFrameConfig() *FrameConfig {
	return &FrameConfig{}
}

//从本地文件读取配置
func ReadFrameConfigFromFile(configFilePath string, defaultSetting *FrameConfig) error {
	fullPath := path.Join(configFilePath, "trframe.yaml")
	viper.SetConfigFile(fullPath)
	err := viper.ReadInConfig()
	if err != nil {
		return err
	}
	return viper.Unmarshal(defaultSetting)
}
