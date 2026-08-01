package main

import (
	"bufio"
	"fmt"
	"os"
	"strings"

	"disk/fs" // assuming module name is disk
)

func main() {
	fmt.Println("Starting Advanced File System Simulation...")
	
	// Create or mount the filesystem
	var myFS *fs.FileSystem
	var err error
	if _, err = os.Stat("virtual_disk.bin"); os.IsNotExist(err) {
		fmt.Println("Formatting new virtual disk...")
		myFS, err = fs.FormatFS("virtual_disk.bin")
	} else {
		fmt.Println("Mounting existing virtual disk...")
		myFS, err = fs.MountFS("virtual_disk.bin")
	}

	if err != nil {
		fmt.Printf("Failed to initialize FS: %v\n", err)
		return
	}
	defer myFS.Close()

	reader := bufio.NewReader(os.Stdin)
	for {
		fmt.Print("fs> ")
		input, err := reader.ReadString('\n')
		if err != nil {
			break
		}
		
		input = strings.TrimSpace(input)
		if input == "" {
			continue
		}
		
		args := strings.Split(input, " ")
		cmd := args[0]
		
		switch cmd {
		case "exit", "quit":
			fmt.Println("Unmounting disk and exiting...")
			return
		case "info":
			fmt.Println("Disk: virtual_disk.bin (10 MB mapped to memory)")
			fmt.Printf("Block size: %d bytes\n", fs.BlockSize)
		case "ls":
			entries := myFS.Ls()
			if len(entries) == 0 {
				fmt.Println("(empty)")
			} else {
				for _, e := range entries {
					fmt.Println(e)
				}
			}
		case "mkdir":
			if len(args) < 2 {
				fmt.Println("Usage: mkdir <dirname>")
				continue
			}
			err := myFS.Mkdir(args[1])
			if err != nil {
				fmt.Printf("Error: %v\n", err)
			} else {
				fmt.Printf("Created directory '%s'\n", args[1])
			}
		case "rm":
			if len(args) < 2 {
				fmt.Println("Usage: rm [-r] <name>")
				continue
			}
			recursive := false
			target := args[1]
			if args[1] == "-r" {
				if len(args) < 3 {
					fmt.Println("Usage: rm -r <name>")
					continue
				}
				recursive = true
				target = args[2]
			}
			
			err := myFS.Rm(target, recursive)
			if err != nil {
				fmt.Printf("Error: %v\n", err)
			} else {
				fmt.Printf("Removed '%s'\n", target)
			}
		default:
			fmt.Printf("Unknown command: %s. (Available: info, ls, mkdir, rm, exit)\n", cmd)
		}
	}
}
