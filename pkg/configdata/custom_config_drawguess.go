package configdata

import "roomcell/pkg/configdata/csvdef"

type DrawTypeWordsCfg struct {
	WordsList []*csvdef.DrawGuess
}

func (cfg *ConfigData) GetTypeWords(wordType int32) *DrawTypeWordsCfg {
	if dt, ok := cfg.typeDrawWordsConfig[wordType]; ok {
		return dt
	}
	return nil
}
func (cfg *ConfigData) GetWordsTypeList() []*DrawTypeWordsCfg {
	retList := make([]*DrawTypeWordsCfg, 0)
	for _, v := range cfg.typeDrawWordsConfig {
		retList = append(retList, v)
	}
	return retList
}
