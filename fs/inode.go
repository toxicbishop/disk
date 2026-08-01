package fs

import (
	"encoding/binary"
	"fmt"
)

const (
	FileTypeRegular = 0
	FileTypeDir     = 1
	InodeSize       = 64
	DirectBlocks    = 10
)

// Inode represents a file or directory's metadata
type Inode struct {
	Mode   uint32   // Unix-style permissions + file type (dir/reg)
	Size   uint32   // Size in bytes
	Blocks [DirectBlocks]uint32 // Direct block pointers
	// (Keeping it simple without indirect blocks for this scale)
}

func (in *Inode) Serialize() []byte {
	b := make([]byte, InodeSize)
	binary.LittleEndian.PutUint32(b[0:4], in.Mode)
	binary.LittleEndian.PutUint32(b[4:8], in.Size)
	for i := 0; i < DirectBlocks; i++ {
		binary.LittleEndian.PutUint32(b[8+i*4:12+i*4], in.Blocks[i])
	}
	return b
}

func DeserializeInode(b []byte) *Inode {
	if len(b) < InodeSize {
		return nil
	}
	in := &Inode{}
	in.Mode = binary.LittleEndian.Uint32(b[0:4])
	in.Size = binary.LittleEndian.Uint32(b[4:8])
	for i := 0; i < DirectBlocks; i++ {
		in.Blocks[i] = binary.LittleEndian.Uint32(b[8+i*4 : 12+i*4])
	}
	return in
}

// Unix-style permission helpers
func (in *Inode) IsDir() bool {
	return in.Mode&(1<<31) != 0 // use highest bit for dir flag
}

func (in *Inode) SetDir(isDir bool) {
	if isDir {
		in.Mode |= (1 << 31)
	} else {
		in.Mode &= ^uint32(1 << 31)
	}
}

func (in *Inode) Permissions() uint16 {
	return uint16(in.Mode & 0x0FFF) // bottom 12 bits for rwx (user, group, other)
}

func (in *Inode) SetPermissions(perms uint16) {
	in.Mode = (in.Mode & 0xFFFFF000) | uint32(perms&0x0FFF)
}

func (in *Inode) String() string {
	typeStr := "-"
	if in.IsDir() {
		typeStr = "d"
	}
	// Simple rwxrwxrwx string
	p := in.Permissions()
	rwx := func(v uint16, shift int) string {
		s := ""
		if v&(4<<shift) != 0 { s += "r" } else { s += "-" }
		if v&(2<<shift) != 0 { s += "w" } else { s += "-" }
		if v&(1<<shift) != 0 { s += "x" } else { s += "-" }
		return s
	}
	return fmt.Sprintf("%s%s%s%s %d", typeStr, rwx(p, 6), rwx(p, 3), rwx(p, 0), in.Size)
}
