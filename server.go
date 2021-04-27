package main

import (
	"net"
	"flag"
	"fmt"
	"sync"
	"io"
	"encoding/json"
)


type Server struct {
	Port		int
	OnlineMap	map[string]*User
	RequestNum	int
	MapLock		sync.RWMutex
	NumLock		sync.RWMutex
}

func NewServer(port int) *Server {
	return &Server{
		Port:		port,
		OnlineMap:	make(map[string]*User),
		RequestNum:	0,
	}
}

type User struct {
	Ip	string
	conn	net.Conn
	server	*Server
}

func NewUser(conn net.Conn, server *Server) *User {
	return &User{
		conn:	conn,
		Ip:	conn.RemoteAddr().String(),
		server:	server,
	}
}

func (user *User)SendMessage(msg string) {
	//向指定客户端发送消息
	_, err := user.conn.Write([]byte(msg + "\n"))
	if err != nil{
		fmt.Println("写入错误", err)
		return
	}
}

func (user *User)Online() {
	//上线,将新用户加入在线列表中
	user.server.MapLock.Lock()
	user.server.OnlineMap[user.Ip] = user
	user.server.MapLock.Unlock()
	fmt.Println(user.Ip,"已登录")
}

func (user *User)Offline() {
	//下线,将用户从在线列表中删除
        user.server.MapLock.Lock()
        delete(user.server.OnlineMap, user.Ip)
        user.server.MapLock.Unlock()
        fmt.Println(user.Ip,"已下线")
}


func main(){
	//可用启动参数配置服务其端口
	port := 8080
	flag.IntVar(&port, "port", 8080, "服务端启动端口")
	flag.Parse()
	server := NewServer(port)
	//启动服务
	server.Start()
}

func (s *Server) Start(){
	add := fmt.Sprintf("127.0.0.1:%d", s.Port)
	//开启端口监听
	listener, err := net.Listen("tcp", add)
	defer listener.Close()
	if err != nil {
		fmt.Println("端口监听异常:", err)
	}
	fmt.Println(s.Port, "端口开始监听...")
	for {
		conn, err := listener.Accept()
		if err != nil {
			fmt.Println("接受异常:", err)
			continue
		}
		//收到新连接，开启处理事务
		go s.Handler(conn)

	}
}

//回复请求时用的结构体
type ReplyMsg struct {
	Message		string	//请求消息内容
	OnlineNum	int	//在线用户数量
	RequestNum	int	//请求数量
	ClientAddr	string
}

func (s *Server)Handler(conn net.Conn) {
	defer conn.Close()
	//上线
	user := NewUser(conn, s)
	user.Online()
	buf := make([]byte, 2048)
	for{
		n, err := conn.Read(buf)
		if n == 0 {
			user.Offline()
			return
		}
		if err != nil && err != io.EOF {
			fmt.Println("信息读入异常:", err)
		}
		msg := string(buf[:n - 1])
		s.NumLock.Lock()
		s.RequestNum ++
		reply := ReplyMsg {
			Message:	msg,
			OnlineNum:	len(s.OnlineMap),
			RequestNum:	s.RequestNum,
			ClientAddr: user.Ip,
		}
		s.NumLock.Unlock()
		bt, err := json.Marshal(reply)
		if err != nil {
			fmt.Println("json转换错误:", err)
			continue
		}
		user.SendMessage(string(bt))
	}
}
