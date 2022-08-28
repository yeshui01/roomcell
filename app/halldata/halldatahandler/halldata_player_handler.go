package halldatahandler

import (
	"roomcell/app/halldata/halldatamain"
	"roomcell/pkg/loghlp"
	"roomcell/pkg/ormdef"
	"roomcell/pkg/pb/pbserver"
	"roomcell/pkg/protocol"
	"roomcell/pkg/tbobj"
	"roomcell/pkg/trframe"
	"roomcell/pkg/trframe/iframe"
	"time"
)

func HandlePlayerLoadRoleData(tmsgCtx *iframe.TMsgContext) (isok int32, retData interface{}, rt iframe.IHandleResultType) {
	req := &pbserver.ESMsgPlayerLoadRoleReq{}
	if !trframe.DecodePBMessage(tmsgCtx.NetMessage, req) {
		return protocol.ECodePBDecodeError, nil, iframe.EHandleContent
	}
	trframe.LogMsgInfo(tmsgCtx.NetMessage, req)
	rep := &pbserver.ESMsgPlayerLoadRoleRep{}
	//
	dbGlobal := hallDataServe.HallDataGlobal()
	gameDB := dbGlobal.GetGameDB()

	dataPlayer := dbGlobal.FindDataPlayer(req.RoleID)
	if dataPlayer == nil {
		tbRoleBase := tbobj.NewTbRoleBase()
		tbRoleBase.SetRoleID(req.RoleID)
		err := gameDB.Model(tbRoleBase.GetOrmMeta()).First(tbRoleBase.GetOrmMeta()).Error
		if err != nil {
			// 不存在，创建
			if !req.WithCreate {
				return protocol.ECodeRoleNotExisted, rep, iframe.EHandleContent
			}
			// 创建
			tbRoleBase.SetRoleID(req.RoleID)
			tbRoleBase.SetUserID(req.CreateParam.UserID)
			tbRoleBase.SetCreateTime(time.Now().Unix())
			tbRoleBase.SetRoleName(req.CreateParam.Account)
			tbRoleBase.SetMoney(0)
			tbRoleBase.SetLevel(req.CreateParam.Level)
			err = gameDB.Model(tbRoleBase.GetOrmMeta()).Create(tbRoleBase.GetOrmMeta()).Error
			if err != nil {
				loghlp.Errorf("create role dberror:%s", err.Error())
				return protocol.ECodeDBError, rep, iframe.EHandleContent
			} else {
				loghlp.Infof("create role succ,%+v", tbRoleBase.GetOrmMeta())
			}
		}
		dataPlayer = halldatamain.NewHalldataPlayer()
		dataPlayer.RoleID = req.RoleID
		dataPlayer.DataTbRoleBase = tbRoleBase
		dbGlobal.AddDataPlayer(req.RoleID, dataPlayer)
	}

	// 返回数据
	rep.RoleDetailData = &pbserver.HallRoleData{}
	oneTable := &pbserver.DbTableData{
		TableID: ormdef.ETableRoleBase,
	}
	oneTable.Data, _ = dataPlayer.DataTbRoleBase.ToBytes()
	rep.RoleDetailData.RoleTables = append(rep.RoleDetailData.RoleTables, oneTable)
	return protocol.ECodeSuccess, rep, iframe.EHandleContent
}

func HandlePlayerSaveRoleData(tmsgCtx *iframe.TMsgContext) (isok int32, retData interface{}, rt iframe.IHandleResultType) {
	req := &pbserver.ESMsgPlayerSaveRoleReq{}
	if !trframe.DecodePBMessage(tmsgCtx.NetMessage, req) {
		return protocol.ECodePBDecodeError, nil, iframe.EHandleContent
	}
	trframe.LogMsgInfo(tmsgCtx.NetMessage, req)
	gameDB := hallDataServe.GetGameDB()
	dbGlobal := hallDataServe.HallDataGlobal()
	dataPlayer := dbGlobal.FindDataPlayer(req.RoleID)
	if dataPlayer != nil {
		// 更新到缓存
		for _, v := range req.RoleTables {
			switch v.TableID {
			case ormdef.ETableRoleBase:
				{
					tbRoleBase := tbobj.NewTbRoleBase()
					tbRoleBase.FromBytes(v.Data)
					dataPlayer.DataTbRoleBase.FromBytes(v.Data)
					// 发送到db线程更新
					dbJob := func() {
						gameDB.Model(tbRoleBase.GetOrmMeta()).Select("*").Updates(tbRoleBase.GetOrmMeta())
					}
					dbGlobal.PostDBJob(&halldatamain.HallDataDBJob{
						DoJob: dbJob,
					})
					break
				}
			default:
				{
					loghlp.Errorf("known table id:%d", v.TableID)
				}
			}
		}
	}
	// for _, v := range req.RoleTables {
	// 	switch v.TableID {
	// 	case ormdef.ETableRoleBase:
	// 		{
	// 			tbRoleBase := tbobj.NewTbRoleBase()
	// 			tbRoleBase.FromBytes(v.Data)
	// 			gameDB.Model(tbRoleBase.GetOrmMeta()).Select("*").Updates(tbRoleBase.GetOrmMeta())
	// 			break
	// 		}
	// 	default:
	// 		{
	// 			loghlp.Errorf("known table id:%d", v.TableID)
	// 		}
	// 	}
	// }
	rep := &pbserver.ESMsgPlayerSaveRoleRep{}
	return protocol.ECodeSuccess, rep, iframe.EHandleContent
}
