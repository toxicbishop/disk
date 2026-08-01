# Advanced File System Simulation in Go

This project is an advanced simulation of a virtual file system written in Go. It demonstrates key concepts in operating system design, including memory-mapped storage, free space management, inodes, and hierarchical directories.

## Features

- **Memory-Mapped Persistence:** Uses `mmap-go` to map a 10MB virtual disk (`virtual_disk.bin`) into memory, ensuring real-time persistence of changes.
- **Free Space Bitmap:** Precisely tracks which 4KB blocks are used and which are free to allow for block reuse.
- **Inodes:** Implements an inode table with Unix-style permissions (`rwxrwxrwx`) and Indexed Allocation pointers.
- **Directories:** Simulates hierarchical directories using `DirEntry` mappings.
- **Interactive Shell:** A custom CLI to interact with the virtual disk.

## Installation

Ensure you have Go installed, then clone the repository and build the shell:

```bash
go mod tidy
go build -o fs-shell ./cmd/shell
```

## Usage

Run the interactive shell:

```bash
./fs-shell
```

### Available Commands

- `info`: Display information about the virtual disk.
- `ls`: List contents of the current directory.
- `mkdir <dirname>`: Create a new directory.
- `rm [-r] <name>`: Remove a file or directory. Use `-r` to forcefully remove a directory.
- `exit`: Unmount the disk and close the shell.

## Legacy Examples

The older, simplified simulations are still available for reference:
- **example1/**: Contiguous file system simulation.
- **example2/**: Paged file system simulation.
- **chunks/**: Chunk calculation test.

## License

This project is licensed under the terms of the MIT License.
