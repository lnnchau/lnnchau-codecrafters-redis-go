package main

import (
	"fmt"
	"io"
	"strconv"

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
	
	resp, _ := ParseRESP(buf[:length]).GetArray()
	
	command, _ := resp[0].GetString()
	args := resp[1:]

	switch command {
	case "echo":
		echoMsg, _ := args[0].GetString()
		conn.Write([]byte(fmt.Sprintf("$%d\r\n%s\r\n", len(echoMsg), echoMsg)))
	case "ping":
		conn.Write([]byte("+PONG\r\n"))
	case "set":
		key, _ := args[0].GetString()
		value, _ := args[1].GetString()

		var expiryCmd string
		var expiryValue int64

		if len(args) > 2 {
			expiryCmd, _ = args[2].GetString()

			expiryValueInStr, _ := args[3].GetString()
			expiryValue, _ = strconv.ParseInt(expiryValueInStr, 10, 64)
		}

		storage.Set(key, value, expiryCmd, expiryValue)
		conn.Write([]byte("+OK\r\n"))
	case "get":
		key, _ := args[0].GetString()
		value, ok := storage.Get(key)

		if (ok) {
			conn.Write([]byte(fmt.Sprintf("$%d\r\n%s\r\n", len(value), value)))
		} else {
			conn.Write([]byte("$-1\r\n"))
		}

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
