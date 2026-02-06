package main

import (
	"fmt"
	"net"
)

func main() {
	ln, err := net.Listen("tcp", ":8080")
	if err != nil {
		fmt.Println("Ошибка:", err)
		return
	}
	defer ln.Close()

	fmt.Println("Сервер слушает :8080")

	for {
		conn, err := ln.Accept()
		if err != nil {
			fmt.Println("Ошибка подключения:", err)
			continue
		}

		fmt.Println("Новое подключение:", conn.RemoteAddr())
		go handleClient(conn)
	}
}

func handleClient(conn net.Conn) {
	defer conn.Close()

	conn.Write([]byte("Добро пожаловать в чат!\n"))
	buf := make([]byte, 1024)
	for {
		n, err := conn.Read(buf)
		if err != nil {
			fmt.Println("Клиент отключился:", conn.RemoteAddr())
			return
		}

		msg := string(buf[:n])
		fmt.Printf("[%s]: %s", conn.RemoteAddr(), msg)
		response := fmt.Sprintf("Вы сказали: %s", msg)
		conn.Write([]byte(response))
	}
}
