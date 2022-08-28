package webserve

import (
	"context"
	"net/http"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/sirupsen/logrus"
)

type WebServe struct {
	router *gin.Engine
}

func NewWebServe() *WebServe {
	w := &WebServe{
		router: gin.Default(),
	}
	return w
}

func (w *WebServe) GetRouter() *gin.Engine {
	return w.router
}

func (w *WebServe) Run(listenAddr string, stopCh chan bool, releaseMode int32) {
	if releaseMode > 0 {
		gin.SetMode(gin.ReleaseMode)
	}
	srv := &http.Server{
		Addr:         listenAddr,
		Handler:      w.router,
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
