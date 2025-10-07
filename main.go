package main

import (
	"flag"
	"fmt"
	"log"
	"os"
	"os/signal"
	"syscall"

	"probe/internal/capture"
	"probe/internal/server"
	"probe/internal/storage"
)

func main() {
	var (
		interfaceName = flag.String("i", "", "网络接口名称 (例如: eth0, en0)")
		port          = flag.String("p", "8080", "Web服务器端口")
		//verbose       = flag.Bool("v", false, "详细输出")
	)
	flag.Parse()

	if *interfaceName == "" {
		fmt.Println("请指定网络接口名称")
		fmt.Println("使用方法: go run main.go -i <接口名称>")
		fmt.Println("例如: go run main.go -i en0")
		os.Exit(1)
	}

	// 初始化存储
	storage := storage.NewMemoryStorage()

	// 初始化抓包器
	capturer, err := capture.NewCapturer(*interfaceName, storage)
	if err != nil {
		log.Fatalf("创建抓包器失败: %v", err)
	}

	// 启动抓包
	go func() {
		if err := capturer.Start(); err != nil {
			log.Fatalf("启动抓包失败: %v", err)
		}
	}()

	// 启动Web服务器
	webServer := server.NewServer(storage, *port)
	go func() {
		fmt.Printf("Web界面启动在: http://localhost:%s\n", *port)
		if err := webServer.Start(); err != nil {
			log.Fatalf("启动Web服务器失败: %v", err)
		}
	}()

	// 等待中断信号
	sigChan := make(chan os.Signal, 1)
	signal.Notify(sigChan, syscall.SIGINT, syscall.SIGTERM)
	<-sigChan

	fmt.Println("\n正在停止抓包工具...")
	capturer.Stop()
	webServer.Stop()
	fmt.Println("已停止")
}
