package bitcask

import (
	"encoding/binary"
	"encoding/json"
	"hash/crc32"
)


func CRC(b []byte) []byte {
	chksum:= crc32.ChecksumIEEE(b)
	cb:= make([]byte, 4)
	binary.LittleEndian.PutUint32(cb[0:4], chksum)

	return cb
}

func Pack32(i uint) []byte {
	var b []byte = make([]byte, 4)
	binary.LittleEndian.PutUint32(b[0:4], uint32(i))
	return b
}

func Pack16(i uint) []byte {
	var b []byte = make([]byte, 2)
	binary.LittleEndian.PutUint16(b[0:2], uint16(i))
	return b
}

func PackAny(v any) []byte {
	//TODO(@Richd0tocom): Handle case where v is a pointer (users can be funny)

	var b []byte
	b, _ = json.Marshal(v)
	return b
}

func serializeHeader(epoch, keysz, valuesz []byte) []byte {
	header:= make([]byte, 12)
	
	copy(header[0:4], epoch)
	copy(header[4:8], keysz)
	copy(header[8:12], valuesz)

	return header
}

func Serialize(header []byte, key, value any) []byte{
	keyBytes:=PackAny(key)
	valueBytes:=PackAny(value)
	datasz:= len(keyBytes)+ len(valueBytes)+ len(header)

	c:= make([]byte, datasz)
	c = append(c, header...)
	c = append(c, keyBytes...)
	c = append(c, valueBytes...)
	
	data:= make([]byte, 4 + datasz)
	data = append(data, CRC(c)...)
	data = append(data, c...)

	return data
}





