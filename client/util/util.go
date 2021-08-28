package util

import (
	"encoding/binary"
)

// func BufChunkToByteArray(buf []byte, leftOffset int, length int) []byte {
// 	return buf[leftOffset : leftOffset+length]
// }

func ByteArrayToUint64(buf []byte) uint64 {
	return binary.BigEndian.Uint64(buf)
}

// func InsertBufChunkInBuf(buf []byte, insert_buf []byte, offset int) {
// 	for i := 0; i < len(insert_buf); i++ {
// 		buf[i+offset] = insert_buf[i]
// 	}
// }

func Uint64ToByteArray(integer uint64) []byte {
	b := make([]byte, 8)
	binary.LittleEndian.PutUint64(b, integer)
	return b
}
