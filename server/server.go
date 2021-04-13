package server

import (
	"context"
	"log"
	"net/http"
	"os"
	"os/signal"
	"time"
)

/**
*@Author lyer
*@Date 2/20/21 15:22
*@Describe
**/
func RunServer(addr string, handler http.Handler) {
	srv := &http.Server{
		Addr:    addr,
		Handler: handler,
	}
	go func() {
		if err := srv.ListenAndServe(); err != nil && err != http.ErrServerClosed {
			panic(err)
		}
	}()
	quit := make(chan os.Signal)
	signal.Notify(quit, os.Interrupt)
	<-quit
	log.Println("shutdown server......")

	//设置超时时间context
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()
	//等待所有连接都关闭 如果都关闭了则结束 这里传入context是为了设置超时时间
	if err := srv.Shutdown(ctx); err != nil {
		panic(err)
	}
	log.Println("Server exiting")
}
