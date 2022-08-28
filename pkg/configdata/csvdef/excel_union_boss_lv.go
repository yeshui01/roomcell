// this file generated by tools,don't edit it!!!
package csvdef

import (
	"roomcell/pkg/configdata/csvparse"
	"strings"

	"github.com/sirupsen/logrus"
)

const (
	UnionBossLvColumnID = iota
	UnionBossLvColumnBossReward
	UnionBossLvColumnRewardNum
	UnionBossLvColumnGrowupNum
	UnionBossLvColumnRankReward
)

type UnionBossLv struct {
	ID         int32
	BossReward []int32
	RewardNum  []int32
	GrowupNum  int32
	RankReward int32
}
type UnionBossLvCfgModule struct {
	cfgMap  map[int32]*UnionBossLv
	cfgList []*UnionBossLv
}

func NewUnionBossLvCfgModule() *UnionBossLvCfgModule {
	return &UnionBossLvCfgModule{
		cfgMap: make(map[int32]*UnionBossLv),
	}
}

func (cfgMod *UnionBossLvCfgModule) Load(lineList []string) {
	logrus.Info("UnionBossLvCfgModule begin load")
	cfgMod.cfgList = make([]*UnionBossLv, len(lineList))
	for i, lstr := range lineList {
		columnList := strings.Split(lstr, ",")
		oneData := &UnionBossLv{}
		oneData.ID = csvparse.ParseColumnInt(columnList[UnionBossLvColumnID])
		oneData.BossReward = csvparse.ParseColumnArrayInt(columnList[UnionBossLvColumnBossReward])
		oneData.RewardNum = csvparse.ParseColumnArrayInt(columnList[UnionBossLvColumnRewardNum])
		oneData.GrowupNum = csvparse.ParseColumnInt(columnList[UnionBossLvColumnGrowupNum])
		oneData.RankReward = csvparse.ParseColumnInt(columnList[UnionBossLvColumnRankReward])
		cfgMod.cfgMap[oneData.ID] = oneData
		cfgMod.cfgList[i] = oneData
	}
	logrus.Infof("UnionBossLvCfgModule load finish!dataNum:%d", len(cfgMod.cfgList))
}

func (cfgMod *UnionBossLvCfgModule) GetData(id int32) *UnionBossLv {
	if cfg, ok := cfgMod.cfgMap[id]; ok {
		return cfg
	}
	return nil
}

func (cfgMod *UnionBossLvCfgModule) GetDataList() []*UnionBossLv {
	return cfgMod.cfgList
}