package fs

import (
	"fmt"
	"os"

	"github.com/edsrzf/mmap-go"
)

const (
	BlockSize = 4096            // 4KB blocks
	DiskSize  = 10 * 1024 * 1024 // 10 MB default disk size
	NumBlocks = DiskSize / BlockSize
)

// Disk manages the virtual disk mapped to memory
type Disk struct {
	file *os.File
	mmap mmap.MMap
}

// OpenDisk opens or creates a memory-mapped virtual disk
func OpenDisk(path string) (*Disk, error) {
	// Open file (create if not exists)
	f, err := os.OpenFile(path, os.O_RDWR|os.O_CREATE, 0644)
	if err != nil {
		return nil, err
	}

	// Ensure file is the correct size
	info, err := f.Stat()
	if err != nil {
		return nil, err
	}

	if info.Size() < DiskSize {
		if err := f.Truncate(DiskSize); err != nil {
			return nil, err
		}
	}

	// Map the file into memory
	m, err := mmap.Map(f, mmap.RDWR, 0)
	if err != nil {
		f.Close()
		return nil, fmt.Errorf("failed to mmap file: %w", err)
	}

	return &Disk{
		file: f,
		mmap: m,
	}, nil
}

// Close flushes and unmaps the disk
func (d *Disk) Close() error {
	var err1, err2, err3 error
	if d.mmap != nil {
		err1 = d.mmap.Flush()
		err2 = d.mmap.Unmap()
	}
	if d.file != nil {
		err3 = d.file.Close()
	}
	if err1 != nil {
		return err1
	}
	if err2 != nil {
		return err2
	}
	return err3
}

// ReadBlock reads a block of data from the disk
func (d *Disk) ReadBlock(blockIndex int) ([]byte, error) {
	if blockIndex < 0 || blockIndex >= NumBlocks {
		return nil, fmt.Errorf("block index %d out of bounds", blockIndex)
	}
	start := blockIndex * BlockSize
	end := start + BlockSize
	
	// Return a copy so the caller doesn't accidentally modify the mmap directly without knowing
	data := make([]byte, BlockSize)
	copy(data, d.mmap[start:end])
	return data, nil
}

// WriteBlock writes a block of data to the disk
func (d *Disk) WriteBlock(blockIndex int, data []byte) error {
	if blockIndex < 0 || blockIndex >= NumBlocks {
		return fmt.Errorf("block index %d out of bounds", blockIndex)
	}
	if len(data) > BlockSize {
		return fmt.Errorf("data exceeds block size")
	}
	
	start := blockIndex * BlockSize
	end := start + BlockSize
	
	// Zero out the block and then copy
	for i := start; i < end; i++ {
		d.mmap[i] = 0
	}
	copy(d.mmap[start:end], data)
	return nil
}
