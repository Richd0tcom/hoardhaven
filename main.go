package main

import (
	// "bufio"
	// "bytes"
	"fmt"
	"io"
	"net"
	"os"
	"strings"

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

		if value.Type != "array" {
			fmt.Println("Invalid request, expected array")
			continue
		}

		if len(value.Array) == 0 {
			fmt.Println("Invalid request, expected array length > 0")
			continue
		}

		cmd:= strings.ToUpper(value.Array[0].Bulk)

		args:= value.Array[1:]

		writer:= protocol.NewRESPWriter(conn)

		handler, ok:= protocol.Handlers[cmd]

		if !ok {
			fmt.Println("Invalid command: ", cmd)
			writer.Write(protocol.Value{Type: "string", Str: ""})
			continue
		}

		result:= handler(args)

		writer.Write(result)
		
	}
}