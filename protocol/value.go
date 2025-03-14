package protocol

import (
	"io"
	"strconv"
)


type Value struct {
	//type of the value
	_type string

	//holds the value of the string received from the simple strings.
	str string

	//num holds the value of the integer received from the integers.
	num int

	//bulk is used to store the string received from the bulk strings
	//
	//SEE: https://redis.io/docs/reference/protocol-spec/#resp-bulk-strings.
	bulk string

	//array holds all the values received from the
	array []Value
}
func (v Value) Marshal() []byte {

	switch v._type {
	case "array":
		return v.marshalArray()
	case "bulk":
		return v.marshalBulk()
	case "string":
		return v.marshalString()
	case "null":
		return v.marshalNull()
	case "error":
		return v.marshalError()
	default:
		return []byte{}
	}
}


func (v Value) marshalString() []byte {
	bytes:= []byte{}

	bytes = append(bytes, S)
	bytes = append(bytes, v.str...)
	bytes = append(bytes, "\r\n"...)

	return bytes

}

func (v Value) marshalBulk() []byte {
	bytes:= []byte{}

	bytes = append(bytes, B)
	bytes = append(bytes, strconv.Itoa(len(v.bulk))...)
	bytes = append(bytes, "\r\n"...)
	bytes = append(bytes, v.bulk...)
	bytes = append(bytes, "\r\n"...)

	return bytes

}


func (v Value) marshalArray() []byte {
	bytes:= []byte{}

	bytes = append(bytes, A)
	bytes = append(bytes, strconv.Itoa(len(v.array))...)
	bytes = append(bytes, "\r\n"...)
	
	for i := range v.array {
		bytes = append(bytes, v.array[i].Marshal()...)
	} 

	return bytes
}


func (v Value) marshalError() []byte {
	bytes:= []byte{}

	bytes = append(bytes, E)
	bytes = append(bytes, "\r\n"...)

	return bytes
}

func (v Value) marshalNull() []byte {
	return []byte("$-1\r\n")
}

type Writer struct {
	writer io.Writer
}

func NewRESPWriter(w io.Writer) *Writer{
	return &Writer{
		writer: w,
	}
}