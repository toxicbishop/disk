package fs

import (
	"bytes"
	"encoding/binary"
)

// DirEntry is a simple mapping from name to Inode number
type DirEntry struct {
	InodeNum uint32
	Name     string // Fixed to 28 bytes in disk
}

func SerializeDirEntry(e DirEntry) []byte {
	b := make([]byte, 32)
	binary.LittleEndian.PutUint32(b[0:4], e.InodeNum)
	copy(b[4:], e.Name)
	return b
}

func DeserializeDirEntry(b []byte) DirEntry {
	num := binary.LittleEndian.Uint32(b[0:4])
	nameBytes := b[4:32]
	name := string(bytes.TrimRight(nameBytes, "\x00"))
	return DirEntry{InodeNum: num, Name: name}
}
