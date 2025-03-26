package protocol

import (
	"bufio"
	"fmt"
	"io"
	"strconv"
)

//Protocol Specification: https://redis.io/docs/latest/develop/reference/protocol-spec/
const (
	//string
	S = '+'

	//Array
	A = '*'

	//Integer
	I = ':'

	//Bulk
	B = '$'

	//Error
	E =  '-'
)

type RESP struct {
	reader *bufio.Reader
}

func NewRESP(r io.Reader) *RESP {
	return &RESP{
		reader: bufio.NewReader(r),
	}
}

//reads each line
func (r *RESP) ReadLine() (line []byte, n int, e error){

	for {
		b, err:= r.reader.ReadByte()

		if err != nil {
			return nil, 0, err
		}

		n+=1
		line = append(line, b)

		if len(line) >= 2 && line[len(line)-2] == '\r' {
			break
		}
	}

	//exclude the CRLF chars \r\n
	return line[:len(line)-2], n, nil
}

// like Readline but reads and parses an integer 
func (r *RESP) readInteger() (x int, n int, err error) {
	line, n, err := r.ReadLine() 
	if err != nil {
		return 0, 0, err
	}
	i64, err := strconv.ParseInt(string(line), 10, 64)
	if err != nil {
		return 0, n, err
	}
	return int(i64), n, nil
}

func (r *RESP) readArray() (Value,  error) {
	v:= Value{}
	v.Type = "array"

	length,_, err:= r.readInteger()

	if err != nil {
		return v, err
	}

	v.Array = make([]Value, length)

	for i:=range length {
		val, err:= r.Read()

		if err != nil {
			return Value{}, err
		}

		v.Array[i] = val
	}




	return v, nil
}

func (r *RESP) Read() (Value, error) {
	typ, err:= r.reader.ReadByte()

	fmt.Println("BYTE READ",string(typ))

	if err != nil {
		return Value{}, err
	}

	switch typ {
		case A:
			return r.readArray()
		case S:
			return r.readSimpleString()
		case B:
			return r.readBulk()
		//implement reading raw integers
		default:
			fmt.Printf("Unknown type: %v", string(typ))
			return Value{}, nil
	}
}

func (r *RESP) readBulk() (Value, error) {
	v:= Value{}

	v.Type = "bulk"

	length, _ , err:= r.readInteger()

	if err != nil {
		return Value{}, err
	}

	bulk:= make([]byte, length)

	r.reader.Read(bulk)

	v.Bulk = string(bulk)

	r.ReadLine()

	return v, nil
}

func (r *RESP) readSimpleString() (Value, error) {
	return Value{}, nil
}

func (r *RESP) readRawInteger() (Value,  error) {
	return Value{}, nil
}