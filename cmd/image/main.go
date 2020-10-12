package main

import (
	"fmt"
	"image"
	"os"

	_ "image/gif"
	_ "image/jpeg"
	_ "image/png"
)

func main() {
	if len(os.Args) < 2 {
		fmt.Println("usage: image <src>")
		os.Exit(1)
	}

	src := os.Args[1]
	f, err := os.Open(src)
	if err != nil {
		fmt.Println("ERROR: open source:", err)
		os.Exit(2)
	}
	defer f.Close()
	cfg, format, err := image.DecodeConfig(f)
	if err != nil {
		fmt.Println("ERROR: decode config:", err)
		os.Exit(2)
	}
	fmt.Printf("Source: %q\n", src)
	fmt.Printf("Format: %q\n", format)
	fmt.Printf("Dim   : w=%d, h=%d\n", cfg.Width, cfg.Height)
}
