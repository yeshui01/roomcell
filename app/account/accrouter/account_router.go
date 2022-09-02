package accrouter

import (
	"net/http"
	"roomcell/pkg/crossdef"
	"roomcell/pkg/loghlp"
	"roomcell/pkg/protocol"
	"roomcell/pkg/sconst"
	"strings"
	"time"

	"github.com/gin-gonic/gin"
	"gorm.io/gorm"
)

type IAccApp interface {
	GetRouter() *gin.Engine
	GetAccountDB() *gorm.DB
}

var accApp IAccApp

func SetupAccountApi(app IAccApp) {
	accApp = app
	accApi := accApp.GetRouter().Group("/account")
	accApi.POST("/register", registerAccount)
	accApi.POST("/login", loginAccount)

	noticeApi := accApp.GetRouter().Group("/notice")
	noticeApi.POST("/query", queryNotice)
}

// 注册
func registerAccount(c *gin.Context) {
	var req AccountRegisterReq
	if err := c.BindJSON(&req); err != nil {
		loghlp.Errorf("req bind json fail!!!")
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}
	accDB := accApp.GetAccountDB()
	var userAccount OrmUser = OrmUser{
		UserName: req.UserName,
	}
	if len(req.UserName) < 5 || len(req.UserName) > 64 {
		c.JSON(http.StatusOK, gin.H{
			"code": protocol.ECodeParamError,
			"msg":  "账号名字格式错误",
		})
		return
	}
	errdb := accDB.Model(userAccount).Where("user_name=?", req.UserName).First(&userAccount).Error
	if errdb == nil {
		// c.JSON(http.StatusOK, gin.H{
		// 	"code": protocol.ECodeDBError,
		// 	"msg":  "db error",
		// })
		// return
		c.JSON(http.StatusOK, gin.H{
			"code": protocol.ECodeAccNameHasExisted,
			"msg":  "账号已经存在",
		})
		return
	}
	if userAccount.UserID > 0 {
		// 账号已经存在
		c.JSON(http.StatusOK, gin.H{
			"code": protocol.ECodeAccNameHasExisted,
			"msg":  "db error",
		})
		return
	}
	if len(req.Pswd) < 6 || len(req.Pswd) > 12 {
		// 密码长度不符合
		c.JSON(http.StatusOK, gin.H{
			"code": protocol.ECodeParamError,
			"msg":  "param error",
		})
		return
	}
	// 注册
	userAccount.UserName = req.UserName
	userAccount.RegisterTime = time.Now().Unix()
	userAccount.Pswd = req.Pswd
	userAccount.Status = 0
	var hallInfo OrmHallList = OrmHallList{
		Recommend: 1,
	}
	errHall := accDB.Model(hallInfo).Where("recommend=?", hallInfo.Recommend).First(&hallInfo).Error
	if errHall != nil {
		c.JSON(http.StatusOK, gin.H{
			"code": protocol.ECodeAccNotExisted,
			"msg":  "db error",
		})
		//return
	}
	userAccount.DataZone = hallInfo.ID // 目前默认1
	if userAccount.DataZone == 0 {
		userAccount.DataZone = 1
	}
	errdb = accDB.Create(&userAccount).Error
	if errdb != nil {
		// 账号已经存在
		c.JSON(http.StatusOK, gin.H{
			"code": protocol.ECodeAccNameHasExisted,
			"msg":  "db error",
		})
		return
	}

	// 注册成功
	var rep AccountRegisterRsp
	rep.UserID = userAccount.UserID

	c.JSON(http.StatusOK, gin.H{
		"code": 0,
		"msg":  "register success",
		"data": rep,
	})
}

// 登录
func loginAccount(c *gin.Context) {
	var req AccountLoginReq
	if err := c.BindJSON(&req); err != nil {
		loghlp.Errorf("req bind json fail!!!")
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}
	accDB := accApp.GetAccountDB()
	var userAccount OrmUser = OrmUser{
		UserName: req.UserName,
	}
	errdb := accDB.Model(userAccount).Where("user_name=?", req.UserName).First(&userAccount).Error
	if errdb != nil {
		c.JSON(http.StatusOK, gin.H{
			"code": protocol.ECodeAccNotExisted,
			"msg":  "db error",
		})
		return
	}
	// 密码验证
	if req.Pswd != userAccount.Pswd {
		loghlp.Errorf("check user(%s) password fail,pswd:%s, input pswd:%s", req.UserName, userAccount.Pswd, req.Pswd)
		c.JSON(http.StatusOK, gin.H{
			"code": protocol.ECodeAccPasswordError,
			"msg":  "密码错误",
		})
		return
	}
	nowTime := time.Now().Unix()
	if nowTime-userAccount.RegisterTime >= sconst.AccountCertificationTime {
		c.JSON(http.StatusOK, gin.H{
			"code": protocol.ECodeAccCertificationTimeOut,
			"msg":  "db error",
		})
		return
	}
	// 生成token
	token, errToken := genJwtToken(userAccount.UserID, userAccount.UserName, userAccount.DataZone)
	if errToken != nil {
		// 系统错误
		c.JSON(http.StatusOK, gin.H{
			"code": protocol.ECodeSysError,
			"msg":  "sys error",
		})
		return
	}
	// 当前大厅
	if strings.HasPrefix(userAccount.UserName, "debug_") {
		// 临时修改zone为2测试
		userAccount.DataZone = 2
		loghlp.Warnf("find debug_ prefix, dynamic modify the data zone to 2 for test")
	}
	var hallInfo OrmHallList = OrmHallList{}
	errHall := accDB.Model(hallInfo).Where("id=?", userAccount.DataZone).First(&hallInfo).Error
	if errHall != nil {
		c.JSON(http.StatusOK, gin.H{
			"code": protocol.ECodeAccNotExisted,
			"msg":  "db error",
		})
		return
	}
	var rep AccountLoginRsp = AccountLoginRsp{
		Token: token,
		// HallAddr: "localhost:7200", // 测试数据
		HallAddr: hallInfo.GateAddr,
		RestTime: int32((userAccount.RegisterTime + int64(sconst.AccountCertificationTime)) - nowTime),
	}

	// 解析token
	jwtObj := crossdef.NewJWT()
	jwtObj.SetSignKey(crossdef.SignKey)
	claimData, errJwt := jwtObj.ParseToken(token)
	if errJwt == nil {
		loghlp.Infof("parse token success:%+v", claimData)
	} else {
		loghlp.Errorf("parse token fail:%s", errJwt.Error())
	}
	c.JSON(http.StatusOK, gin.H{
		"code": 0,
		"msg":  "login success",
		"data": rep,
	})
}

// 查询公告
func queryNotice(c *gin.Context) {
	var req QueryNoticeReq
	if err := c.BindJSON(&req); err != nil {
		loghlp.Errorf("req bind json fail!!!")
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}
	accDB := accApp.GetAccountDB()
	var cellNotice OrmCellNotice = OrmCellNotice{}
	errdb := accDB.Model(cellNotice).First(&cellNotice).Error
	if errdb != nil {
		c.JSON(http.StatusOK, gin.H{
			"code": protocol.ECodeNotFindNotice,
			"msg":  "not find notice",
		})
		return
	}
	var rep QueryNoticeRsp = QueryNoticeRsp{
		ID:     cellNotice.ID,
		Notice: cellNotice.Content,
	}
	c.JSON(http.StatusOK, gin.H{
		"code": 0,
		"msg":  "login success",
		"data": rep,
	})
}
