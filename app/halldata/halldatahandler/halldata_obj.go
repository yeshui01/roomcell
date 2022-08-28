package halldatahandler

import "roomcell/app/halldata/ihalldata"

var (
	hallDataServe ihalldata.IHallData
)

func InitHallDataObj(ihdata ihalldata.IHallData) {
	hallDataServe = ihdata
}
