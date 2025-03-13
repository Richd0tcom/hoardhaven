package main

import (
	// "bufio"
	// "bytes"
	"fmt"
	"io"
	"net"
	"os"

	"github.com/richd0tcom/hoardhaven/protocol"
	// "strings"
)

func main() {
    fmt.Println("Hello World")
	fmt.Println("Listening on port :6379")

	// m:= make([]byte, 4)
	// 	s:= "\r\nHello\r\n"

	// 	bf:= bytes.NewBufferString(s)

	// 	bf.Read(m)

	// 	fmt.Println("bytes",m)

	// 	ss:= string(m)

	// 	fmt.Println("string", ss )

	// 	fmt.Println([]byte(ss))

	l, err := net.Listen("tcp", ":6379")
	if err != nil {
	    fmt.Println(err)
    	return
	}

	conn, err:=l.Accept()
	if err != nil {
	    fmt.Println(err)
    	return
	}

	defer conn.Close()

	for {


		resp:=protocol.NewRESP(conn)
		value, err := resp.Read()

		if err != nil {
			if err == io.EOF {
				break
			}
			fmt.Println("error reading from client: ", err.Error())
			os.Exit(1)
		}

		fmt.Println(value)
		
		

		
		// r := bufio.NewReader(strings.NewReader(""))

		// byt, err:= r.ReadByte()

		
		
		
		conn.Write([]byte("+OK\r\n"))
	}
}