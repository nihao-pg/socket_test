package main

import (
	"flag"
	"fmt"
	"io"
	"net"
)

var (
	ServerIp   string
	ServerPort int
)

func main() {
	flag.StringVar(&ServerIp, "serverIp", "127.0.0.1", "服务端Ip") //命令行参数
	flag.IntVar(&ServerPort, "serverPort", 8080, "服务端端口")
	flag.Parse()
	//拼接地址
	serverAddr := fmt.Sprintf("%s:%d", ServerIp, ServerPort)
	fmt.Printf("连接至服务器%s\n", serverAddr)
	conn, err := net.Dial("tcp", fmt.Sprintf("%s", serverAddr))
	if err != nil {
		fmt.Println("连接服务器失败!!!")
		return
	}
	go func() {
		buf := make([]byte, 1024)
		for {
			n, err := conn.Read(buf)
			if n == 0 {
				fmt.Println("失去连接")
				return
			}

			if err != nil && err != io.EOF {
				fmt.Println("读取异常:", err)
			}
			fmt.Println(string(buf[:n-1]))

		}
	}()
	fmt.Println("请输入要发送的内容(输入exit退出):")
	var msg string
	_, err = fmt.Scan(&msg)
	if err != nil {
		fmt.Println("输入异常:", err)
	}
	for msg != "exit" {
		if err != nil {
			fmt.Println("输入异常:", err)
			continue
		}
		if msg == "" {
			continue
		}
		conn.Write([]byte(msg))
		fmt.Println("请输入要发送的内容(输入exit退出):")
		_, err = fmt.Scan(&msg)
	}
}
