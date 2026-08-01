package fs
import (
	"fmt"
)


const (
	SuperBlockIndex = 0
	InodeTableStart = 1
	InodeTableSize  = 5
	DataBlockStart  = 6
	
	MaxInodes = (BlockSize / InodeSize) * InodeTableSize
)

type FileSystem struct {
	disk   *Disk
	bitmap *Bitmap
}

func FormatFS(path string) (*FileSystem, error) {
	d, err := OpenDisk(path)
	if err != nil {
		return nil, err
	}

	// Read block 0 to setup bitmap (it should be all zeros since it's a new file)
	blk0, err := d.ReadBlock(SuperBlockIndex)
	if err != nil {
		return nil, err
	}
	
	bitmap := NewBitmap(blk0)
	
	// Mark Superblock and Inode tables as allocated
	for i := 0; i < DataBlockStart; i++ {
		bitmap.Allocate() 
	}
	
	fs := &FileSystem{disk: d, bitmap: bitmap}
	
	// Create root directory at inode 0
	rootInode := &Inode{}
	rootInode.SetDir(true)
	rootInode.SetPermissions(0755)
	
	// Allocate 1 data block for root dir entries
	blk, err := fs.allocateBlock()
	if err == nil {
		rootInode.Blocks[0] = uint32(blk)
	}
	
	fs.writeInode(0, rootInode)
	fs.syncBitmap()
	
	return fs, nil
}

func MountFS(path string) (*FileSystem, error) {
	d, err := OpenDisk(path)
	if err != nil {
		return nil, err
	}

	blk0, err := d.ReadBlock(SuperBlockIndex)
	if err != nil {
		return nil, err
	}

	bitmap := NewBitmap(blk0)
	return &FileSystem{disk: d, bitmap: bitmap}, nil
}

func (fs *FileSystem) syncBitmap() {
	fs.disk.WriteBlock(SuperBlockIndex, fs.bitmap.data)
}

func (fs *FileSystem) allocateBlock() (int, error) {
	idx, err := fs.bitmap.Allocate()
	if err != nil {
		return -1, err
	}
	fs.syncBitmap()
	return idx, nil
}

func (fs *FileSystem) freeBlock(idx int) {
	fs.bitmap.Free(idx)
	fs.syncBitmap()
}

func (fs *FileSystem) readInode(idx uint32) *Inode {
	if idx >= MaxInodes {
		return nil
	}
	blockIdx := InodeTableStart + (idx / (BlockSize / InodeSize))
	offset := (idx % (BlockSize / InodeSize)) * InodeSize
	
	blk, _ := fs.disk.ReadBlock(int(blockIdx))
	return DeserializeInode(blk[offset : offset+InodeSize])
}

func (fs *FileSystem) writeInode(idx uint32, in *Inode) {
	if idx >= MaxInodes {
		return
	}
	blockIdx := InodeTableStart + (idx / (BlockSize / InodeSize))
	offset := (idx % (BlockSize / InodeSize)) * InodeSize
	
	blk, _ := fs.disk.ReadBlock(int(blockIdx))
	copy(blk[offset:offset+InodeSize], in.Serialize())
	fs.disk.WriteBlock(int(blockIdx), blk)
}

func (fs *FileSystem) Close() error {
	fs.syncBitmap()
	return fs.disk.Close()
}

// Ls returns a list of directory entries in the root directory (simplified to root only for now)
func (fs *FileSystem) Ls() []string {
	root := fs.readInode(0)
	if root == nil || !root.IsDir() {
		return nil
	}

	var names []string
	// Read the first direct block of the root directory
	blkIdx := root.Blocks[0]
	if blkIdx == 0 {
		return names
	}

	blk, _ := fs.disk.ReadBlock(int(blkIdx))
	// Parse DirEntries (32 bytes each)
	for offset := 0; offset < BlockSize; offset += 32 {
		entry := DeserializeDirEntry(blk[offset : offset+32])
		if entry.InodeNum != 0 || entry.Name != "" {
			names = append(names, entry.Name)
		}
	}
	return names
}

// Mkdir creates a new directory in the root directory
func (fs *FileSystem) Mkdir(name string) error {
	root := fs.readInode(0)
	if root == nil || !root.IsDir() {
		return fmt.Errorf("root directory not found")
	}

	// 1. Find a free Inode (start from 1 since 0 is root)
	var newInodeNum uint32 = 0
	for i := uint32(1); i < MaxInodes; i++ {
		in := fs.readInode(i)
		if in != nil && in.Mode == 0 { // Empty inode
			newInodeNum = i
			break
		}
	}
	if newInodeNum == 0 {
		return fmt.Errorf("no free inodes")
	}

	// 2. Initialize new directory Inode
	newDir := &Inode{}
	newDir.SetDir(true)
	newDir.SetPermissions(0755)
	
	// Allocate a block for the new directory's contents
	blkIdx, err := fs.allocateBlock()
	if err != nil {
		return err
	}
	newDir.Blocks[0] = uint32(blkIdx)
	fs.writeInode(newInodeNum, newDir)

	// 3. Add entry to root directory
	rootBlkIdx := root.Blocks[0]
	rootBlk, _ := fs.disk.ReadBlock(int(rootBlkIdx))
	
	entryAdded := false
	for offset := 0; offset < BlockSize; offset += 32 {
		entry := DeserializeDirEntry(rootBlk[offset : offset+32])
		if entry.InodeNum == 0 && entry.Name == "" { // Free slot
			newEntry := DirEntry{InodeNum: newInodeNum, Name: name}
			copy(rootBlk[offset:offset+32], SerializeDirEntry(newEntry))
			entryAdded = true
			break
		}
	}
	
	if !entryAdded {
		return fmt.Errorf("root directory full")
	}
	
	fs.disk.WriteBlock(int(rootBlkIdx), rootBlk)
	return nil
}
