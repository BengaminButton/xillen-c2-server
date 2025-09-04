package main

import (
	"bufio"
	"fmt"
	"net"
	"os"
	"strings"
)

var author = "t.me/Bengamin_Button t.me/XillenAdapter"

type Client struct {
	conn   net.Conn
	addr   string
	active bool
}

type C2Server struct {
	clients []Client
	port    string
}

func (s *C2Server) handleClient(client *Client) {
	defer client.conn.Close()
	scanner := bufio.NewScanner(client.conn)
	for scanner.Scan() {
		command := strings.TrimSpace(scanner.Text())
		if command == "exit" {
			client.active = false
			break
		}
		fmt.Printf("[%s] %s\n", client.addr, command)
		s.sendResponse(client, "Command received: "+command)
	}
}

func (s *C2Server) sendResponse(client *Client, response string) {
	client.conn.Write([]byte(response + "\n"))
}

func (s *C2Server) start() {
	listener, err := net.Listen("tcp", ":"+s.port)
	if err != nil {
		fmt.Println("Ошибка запуска сервера:", err)
		return
	}
	defer listener.Close()
	fmt.Printf("C2 сервер запущен на порту %s\n", s.port)
	for {
		conn, err := listener.Accept()
		if err != nil {
			continue
		}
		client := &Client{
			conn:   conn,
			addr:   conn.RemoteAddr().String(),
			active: true,
		}
		s.clients = append(s.clients, *client)
		fmt.Printf("Новый клиент подключён: %s\n", client.addr)
		go s.handleClient(client)
	}
}

func main() {
	fmt.Println(author)
	port := "8080"
	if len(os.Args) > 1 {
		port = os.Args[1]
	}
	server := &C2Server{port: port}
	server.start()
}

