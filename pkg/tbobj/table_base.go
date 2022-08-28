/* ====================================================================
 * Author           : tianyh(mknight)
 * Email            : 824338670@qq.com
 * Last modified    : 2022-04-08 10:33
 * Filename         : table_base.go
 * Description      :
 * ====================================================================*/
package tbobj

const (
	DbStatusNone   int8 = 0
	DbStatusInsert int8 = 1
	DbStatusUpdate int8 = 2
	DbStatusDelete int8 = 3
)

type TableObjOpt struct {
	dbStatus int8
}

func (tb *TableObjOpt) SetDbStatus(s int8) {
	switch tb.dbStatus {
	case DbStatusNone:
		tb.dbStatus = s
	case DbStatusInsert:
		if s == DbStatusDelete {
			tb.dbStatus = s
		}
	case DbStatusUpdate:
		if s == DbStatusDelete {
			tb.dbStatus = s
		}
	}
}

func (tb *TableObjOpt) GetDbStatus() int8 {
	return tb.dbStatus
}

func (tb *TableObjOpt) ClearDbStatus() {
	tb.dbStatus = DbStatusNone
}
