package main

import (
	"fmt"
	"log"
	"os"
	"path"
	"strconv"
	"strings"

	"github.com/disintegration/imaging"
)

func main() {
	filename := os.Args[1]
	width, _ := strconv.Atoi(os.Args[2])
	height, _ := strconv.Atoi(os.Args[3])

	fmt.Println(filename)
	src, err := imaging.Open(filename)
	if err != nil {
		log.Fatalf("Open failed: %v", err)
	}

	src = imaging.Resize(src, width, height, imaging.Lanczos)

	ext := path.Ext(filename)

	err = imaging.Save(src, fmt.Sprintf("%v.%v%v", strings.TrimSuffix(filename, ext), "_s", ext))
	if err != nil {
		log.Fatalf("Save failed: %v", err)

	}
}
