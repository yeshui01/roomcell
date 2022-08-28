package ihalldata

import (
	"roomcell/app/halldata/halldatamain"

	"gorm.io/gorm"
)

type IHallData interface {
	GetGameDB() *gorm.DB
	HallDataGlobal() *halldatamain.HallDataGlobal
}
