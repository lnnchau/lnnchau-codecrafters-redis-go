package main

import (
	"fmt"
	"io"

	"net"
	"os"
)


func processConn(conn net.Conn, storage Storage) {
	buf := make([]byte, 1024)
	length, err := conn.Read(buf)
	if err != nil {
		if err == io.EOF {
			return
		}

		fmt.Println("error reading from client: ", err.Error())
		os.Exit(1)
	}

	fmt.Printf("Message received:\n%s\n", string(buf[:length]))
	
	resp := ParseRESP(buf[:length])
	
	command := resp[0]
	args := resp[1:]

	switch command {
	case "echo":
		echoMsg := args[0]
		conn.Write([]byte(fmt.Sprintf("$%d\r\n%s\r\n", len(echoMsg), echoMsg)))
	case "ping":
		conn.Write([]byte("+PONG\r\n"))
	case "set":
		key := args[0]
		value := args[1]

		storage.Set(key, value)

		conn.Write([]byte("+OK\r\n"))
	case "get":
		key := args[0]
		value := storage.Get(key)

		conn.Write([]byte(fmt.Sprintf("$%d\r\n%s\r\n", len(value), value)))
	default:
		conn.Write([]byte("+OK\r\n"))
	}
}

func handleConn(conn net.Conn, storage Storage) {
	defer conn.Close()
	for {
		processConn(conn, storage)
	}
}

func main() {
	// You can use print statements as follows for debugging, they'll be visible when running tests.
	fmt.Println("Logs from your program will appear here!")

	// Uncomment this block to pass the first stage
	
	l, err := net.Listen("tcp", "0.0.0.0:6379")
	if err != nil {
		fmt.Println("Failed to bind to port 6379")
		os.Exit(1)
	}

	for {
		conn, err := l.Accept()
		if err != nil {
			fmt.Println("Error accepting connection: ", err.Error())
			os.Exit(1)
		}

		storage := NewStorage()
	
		go handleConn(conn, storage)
	}
}
