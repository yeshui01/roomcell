package wsserve

import (
	"context"
	"net/http"
	"time"

	"roomcell/pkg/loghlp"

	"github.com/gin-gonic/gin"
	"github.com/gorilla/websocket"
	"github.com/sirupsen/logrus"
)

var upgrader = websocket.Upgrader{
	//跨域
	CheckOrigin: func(r *http.Request) bool {
		return true
	},
	ReadBufferSize:  1024,
	WriteBufferSize: 2048,
}

type WSConnectedCallback func(wsConn *websocket.Conn, err error)
type WSServer struct {
	engine       *gin.Engine
	connCallback WSConnectedCallback
	wsRouterPath string
}

func NewWSServer() *WSServer {
	return &WSServer{
		engine: gin.Default(),
	}
}

func (serv *WSServer) Run(listenAddr string, stopCh chan bool, releaseMode int32) {
	if releaseMode > 0 {
		gin.SetMode(gin.ReleaseMode)
	}
	loghlp.Infof("wsserver run:%s", listenAddr)
	srv := &http.Server{
		Addr:         listenAddr,
		Handler:      serv.engine,
		ReadTimeout:  30 * time.Second,
		WriteTimeout: 30 * time.Second,
	}
	go func() {
		if err := srv.ListenAndServe(); err != nil && err != http.ErrServerClosed {
			logrus.Fatalf("listen web: %s\n", err)
		}
	}()

	<-stopCh
	logrus.Info("Shutdown Server ...\n")
	ctx, cancel := context.WithTimeout(context.Background(), time.Second*5)
	defer cancel()
	if err := srv.Shutdown(ctx); err != nil {
		logrus.Fatal("Server Shutdown:", err)
	}
	select {
	case <-ctx.Done():
	}
	logrus.Info("Server exiting")
}

func (serv *WSServer) wsHandler(c *gin.Context) {
	ws, err := upgrader.Upgrade(c.Writer, c.Request, nil)
	if serv.connCallback != nil {
		if err == nil {
			loghlp.Infof("onWebsocket connected")
			serv.connCallback(ws, err)
		} else {
			loghlp.Errorf("upgrade err:", err)
			return
		}
	}
}
func (serv *WSServer) Stop() {
}

// 设置连接路由
func (serv *WSServer) SetupWSRouter(wsPath string, connHandle WSConnectedCallback) {
	serv.connCallback = connHandle
	serv.wsRouterPath = wsPath
	if serv.wsRouterPath == "" {
		serv.wsRouterPath = "/ws"
	}
	serv.engine.Use(Cors()) // 跨域
	// ws router
	serv.engine.GET(serv.wsRouterPath, serv.wsHandler)
}
