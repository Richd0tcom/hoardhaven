package protocol

import (
	"io"
	"strconv"
)


type Value struct {
	//type of the value
	Type string

	//holds the value of the string received from the simple strings.
	Str string

	//num holds the value of the integer received from the integers.
	Num int

	//bulk is used to store the string received from the bulk strings
	//
	//SEE: https://redis.io/docs/reference/protocol-spec/#resp-bulk-strings.
	Bulk string

	//array holds all the values received from the
	Array []Value
}
func (v Value) Marshal() []byte {

	switch v.Type {
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
	bytes = append(bytes, v.Str...)
	bytes = append(bytes, "\r\n"...)

	return bytes

}

func (v Value) marshalBulk() []byte {
	bytes:= []byte{}

	bytes = append(bytes, B)
	bytes = append(bytes, strconv.Itoa(len(v.Bulk))...)
	bytes = append(bytes, "\r\n"...)
	bytes = append(bytes, v.Bulk...)
	bytes = append(bytes, "\r\n"...)

	return bytes

}


func (v Value) marshalArray() []byte {
	bytes:= []byte{}

	bytes = append(bytes, A)
	bytes = append(bytes, strconv.Itoa(len(v.Array))...)
	bytes = append(bytes, "\r\n"...)
	
	for i := range v.Array {
		bytes = append(bytes, v.Array[i].Marshal()...)
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

func (w *Writer) Write(v Value) (int, error) {
	byts:= v.Marshal()
	
	return w.writer.Write(byts)
}

func NewRESPWriter(w io.Writer) *Writer{
	return &Writer{
		writer: w,
	}
}