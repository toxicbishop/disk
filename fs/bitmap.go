package fs

import "fmt"

// Bitmap manages a simple bit array to track free blocks
type Bitmap struct {
	data []byte
}

func NewBitmap(data []byte) *Bitmap {
	return &Bitmap{data: data}
}

// Allocate finds the first free bit, sets it to 1, and returns its index.
func (b *Bitmap) Allocate() (int, error) {
	for i := 0; i < len(b.data); i++ {
		if b.data[i] != 0xFF { // Not fully used
			for bit := 0; bit < 8; bit++ {
				mask := byte(1 << bit)
				if b.data[i]&mask == 0 {
					b.data[i] |= mask
					return i*8 + bit, nil
				}
			}
		}
	}
	return -1, fmt.Errorf("no free space")
}

// Free sets the bit at the given index to 0.
func (b *Bitmap) Free(index int) {
	if index < 0 || index >= len(b.data)*8 {
		return
	}
	byteIdx := index / 8
	bitIdx := index % 8
	b.data[byteIdx] &= ^byte(1 << bitIdx)
}

// IsAllocated checks if the bit at the given index is 1.
func (b *Bitmap) IsAllocated(index int) bool {
	if index < 0 || index >= len(b.data)*8 {
		return false
	}
	byteIdx := index / 8
	bitIdx := index % 8
	return b.data[byteIdx]&(1<<bitIdx) != 0
}
