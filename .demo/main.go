package main

import (
	"bufio"
	"bytes"
	"fmt"
	"log"
	"strings"

	_ "go.beyondstorage.io/services/fs/v4"
	"go.beyondstorage.io/v5/pairs"
	"go.beyondstorage.io/v5/services"
)

// func main() {
// 	var buf bytes.Buffer
// 	container := "fs:///home/ceyhun/projects/tractor/.demo"
// 	store, err := services.NewStoragerFromString(container)
// 	filename := "dd.csv"
// 	// filename = filepath.Join(container, filename)
// 	if err != nil {
// 		log.Fatal(err)
// 	}

// 	cur := int64(0)
// 	fn := func(bs []byte) {
// 		println(string(bs))
// 		cur += int64(len(bs))
// 		// log.Printf("read %d bytes already", cur)
// 	}

// 	// If IoCallback is specified, the storage will call it in every I/O operation.
// 	// User could use this feature to implement progress bar.
// 	n, err := store.Read(filename, &buf, pairs.WithIoCallback(fn))
// 	if err != nil {
// 		log.Fatalf("read %v: %v", filename, err)
// 	}

// 	log.Printf("read size: %d", n)
// 	// log.Printf("read content: %s", buf.Bytes())
// }

func main() {
	var buf bytes.Buffer
	container := "fs:///home/ceyhun/projects/tractor/.demo"
	store, err := services.NewStoragerFromString(container)
	filename := "dd.csv"
	// filename = filepath.Join(container, filename)
	if err != nil {
		log.Fatal(err)
	}
	size := int64(100)
	offset := int64(0)
	var rest []byte = nil
	var scanner *bufio.Scanner
	for {
		n, err := store.Read(filename, &buf, pairs.WithOffset(offset), pairs.WithSize(size))

		if err != nil {
			log.Fatalf("read %v: %v", filename, err)
		} else if n == 0 {
			println("read done at offset:", offset)
			break
		}
		if rest != nil {
			buf = *bytes.NewBuffer(append(rest, buf.Bytes()...))
			rest = nil
		}

		bs := buf.String()
		lines := strings.Split(bs, "\n")
		if !strings.HasSuffix(bs, "\n") {
			if len(lines) > 1 {
				lines = strings.Split(bs, "\n")
				rest = []byte(lines[len(lines)-1])
				lines = lines[:len(lines)-1]
			} else {
				rest = []byte(bs)
				lines = []string{}
			}
		}
		scanner = bufio.NewScanner(strings.NewReader(strings.Join(lines, "\n")))

		for scanner.Scan() {
			fmt.Println(scanner.Text())
		}

		// lines := strings.Split(buf.String(), "\n")

		// if len(lines) > 0 && len(lines[len(lines)-1]) != 0 {
		// 	rest = []byte(lines[len(lines)-1])
		// 	if len(lines) > 1 {
		// 		lines = lines[:len(lines)-2]
		// 	}
		// }
		// for _, line := range lines {
		// 	if len(line) > 0 {
		// 		println(string(line))
		// 		return
		// 	}
		// }

		// if len(lines) > 0 && len(lines[len(lines)-1]) != 0 {
		// 	rest = []byte(lines[len(lines)-1])
		// }
		offset += n
		buf.Reset()
		// println(buf.String())
		// buf.Reset()
	}
	fmt.Println("read done %d", offset)
}

// func main() {
// 	lines := "hello\nworld\n"
// 	tokens := strings.Split(lines, "\n")
// 	println(len(tokens[2]))
// 	if tokens[2] == "" {
// 		println("empty")
// 	}
// }
