package project_common

import (
	"context"
	"github.com/gin-gonic/gin"
	"log"
	"net/http"
	"os"
	"os/signal"
	"syscall"
	"time"
)

func Run(r *gin.Engine, srvName string, addr string, stop func()) {
	srv := &http.Server{
		Addr:    addr,
		Handler: r,
	}
	// 保证下面的优雅启停
	go func() {
		log.Printf("web server running in %s \n", srv.Addr)
		if err := srv.ListenAndServe(); err != nil && err != http.ErrServerClosed {
			log.Fatalln(err)
		}
	}()
	quit := make(chan os.Signal)
	// SIGINT 用户发送INTR字符(Ctrl+C)触发 kill -2
	// SIGTERM 结束程序(可以被捕获，阻塞或者忽略) 正常结束程序
	// Notify函数让signal包将输入信号转发到c。如果没有列出要传递的信号，会将所有输入信号传递到c；否则只传递列出的输入信号。
	// 如果发送时c阻塞了，signal包会直接放弃
	signal.Notify(quit, syscall.SIGTERM, syscall.SIGINT)
	// 代码是在这里进行阻塞的
	<-quit
	log.Printf("Shutting Down project %s... \n", srvName)
	ctx, cancel := context.WithTimeout(context.Background(), 2*time.Second)
	defer cancel()
	if stop != nil {
		stop()
	}

	if err := srv.Shutdown(ctx); err != nil {
		log.Fatalf("%s Shutdown, cause by : %v", srvName, err)
	}
	select {
	case <-ctx.Done():
		log.Println("wait timeout...")
	}
	log.Printf("%s stop success...", srvName)
}
