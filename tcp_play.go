package main

import (
	"flag"
	"fmt"
	"net"
	"os"
	"time"
)

var server = flag.Bool("server", false, "run as a server")
var addr = flag.String("addr", "localhost", "server address")

func main() {
	flag.Parse()
	if *server {
		runServer()
	} else {
		runClient()
	}
}

func runClient() {
	conn, err := net.Dial("tcp", fmt.Sprintf("%s:6789", *addr))
	if err != nil {
		panic("failed creating a client")
	}

	fmt.Println("attempting to write to server")
	for i := 0; i < 10; i++ {
		var bytesWritten int
		bytesWritten, err = conn.Write([]byte("Hello Server"))
		if err != nil {
			panic("failed writing to server")
		}
		fmt.Printf("client bytes written %d", bytesWritten)
		time.Sleep(time.Second)
	}

	buf := make([]byte, 1024)
	// Read the incoming connection into the buffer.
	bytesRead, err := conn.Read(buf)
	if err != nil {
		fmt.Println("Error reading:", err.Error())
	}
	fmt.Printf("server bytes read %d", bytesRead)
}

func runServer() {
	listConn, err := net.Listen("tcp", ":6789")
	if err != nil {
		panic(err.Error())
	}
	defer listConn.Close()
	for {
		// Listen for an incoming connection.
		conn, err := listConn.Accept()
		if err != nil {
			fmt.Println("Error accepting: ", err.Error())
			os.Exit(1)
		}
		go func() {
			time.Sleep(10 * time.Second)
			fmt.Println("attempting to write to client")
			bytesWritten, err := conn.Write([]byte("Hello Client"))
			if err != nil {
				panic("failed writing to client")
			}

			fmt.Printf("server bytes written %d\n", bytesWritten)
			buf := make([]byte, 1024)
			bytesRead, err := conn.Read(buf)
			if err != nil {
				fmt.Println("Error reading:", err.Error())
			}
			fmt.Printf("client bytes read %d", bytesRead)
		}()
	}
}
