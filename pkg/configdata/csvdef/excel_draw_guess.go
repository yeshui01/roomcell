// this file generated by tools,don't edit it!!!
package csvdef
import (
	"roomcell/pkg/configdata/csvparse"
	"strings"

	"github.com/sirupsen/logrus"
)
const (
	DrawGuessColumnID = iota
	DrawGuessColumnWordType
	DrawGuessColumnWord
)
type DrawGuess struct {
	ID      	int32
	WordType	int32
	Word    	string
}
type DrawGuessCfgModule struct {
	cfgMap  map[int32]*DrawGuess
	cfgList []*DrawGuess
}

func NewDrawGuessCfgModule() *DrawGuessCfgModule {
	return &DrawGuessCfgModule{
		cfgMap: make(map[int32]*DrawGuess),
	}
}

func (cfgMod *DrawGuessCfgModule) Load(lineList []string) {
	logrus.Info("DrawGuessCfgModule begin load")
	cfgMod.cfgList=make([]*DrawGuess, len(lineList))
	for i, lstr := range lineList {
		columnList := strings.Split(lstr, ",")
		oneData := &DrawGuess{}
		oneData.ID = csvparse.ParseColumnInt(columnList[DrawGuessColumnID])
		oneData.WordType = csvparse.ParseColumnInt(columnList[DrawGuessColumnWordType])
		oneData.Word = csvparse.ParseColumnString(columnList[DrawGuessColumnWord])
		cfgMod.cfgMap[oneData.ID] = oneData
		cfgMod.cfgList[i] = oneData
	}
	logrus.Infof("DrawGuessCfgModule load finish!dataNum:%d",len(cfgMod.cfgList))
}

func (cfgMod *DrawGuessCfgModule) GetData(id int32) *DrawGuess {
	if cfg, ok := cfgMod.cfgMap[id]; ok {
		return cfg
	}
	return nil
}

func (cfgMod *DrawGuessCfgModule) GetDataList() []*DrawGuess {
	return cfgMod.cfgList
}
