package main

import "encoding/binary"

func BufChunkToByteArray(buf []byte, leftOffset int, length int) []byte {
	var temp []byte

	for i := leftOffset; i < leftOffset+length; i++ {
		temp = append(temp, buf[i])
	}

	return temp
}

func BufChunkToUint64(buf []byte, leftOffset int, length int) uint64 {
	return binary.BigEndian.Uint64(BufChunkToByteArray(buf, leftOffset, length))
}

func InsertBufChunkInBuf(buf []byte, insert_buf []byte, offset int) {
	for i := 0; i < len(insert_buf); i++ {
		buf[i+offset] = insert_buf[i]
	}
}
