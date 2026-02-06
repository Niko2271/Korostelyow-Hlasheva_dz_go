package main

import (
	"bufio"
	"fmt"
	"log"
	"net"
	"os"
	"strings"
)

func main() {
	fmt.Println("TCP-клиент")
	fmt.Println("Подключение к localhost:8080...")

	conn, err := net.Dial("tcp", "localhost:8080")
	if err != nil {
		log.Fatal("Не удалось подключиться:", err)
	}
	defer conn.Close()

	fmt.Println("Подключено к localhost:8080")
	fmt.Println("Введите сообщения. Для выхода нажмите Ctrl+C")
	fmt.Println(strings.Repeat("-", 40))

	go readFromServer(conn)
	reader := bufio.NewReader(os.Stdin)

	for {
		fmt.Print("> ")
		text, _ := reader.ReadString('\n')
		text = strings.TrimSpace(text)
		if text == "" {
			continue
		}
		_, err := conn.Write([]byte(text + "\n"))
		if err != nil {
			fmt.Println("Ошибка отправки:", err)
			break
		}
	}
}

func readFromServer(conn net.Conn) {
	buffer := make([]byte, 1024)

	for {
		n, err := conn.Read(buffer)
		if err != nil {
			fmt.Println("\nСоединение разорвано")
			os.Exit(0)
		}

		fmt.Printf("\rСервер: %s\n> ", string(buffer[:n]))
	}
}
