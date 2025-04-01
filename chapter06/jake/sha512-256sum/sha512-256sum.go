package main

import (
	"crypto/sha512"
	"flag"
	"fmt"
	"io/ioutil"
	"os"
)

func init() {
	flag.Usage = func() {
		fmt.Printf("Usage: %s file...\n", os.Args[0])
		flag.PrintDefaults()
	}
}

func main() {
	flag.Parse()
	for _, file := range flag.Args() {
		// 파일 경로를 받아서 checksum 을 진행한다.
		fmt.Printf("\n%s =>\n%s\n", file, checksum(file))
	}
}

// checksum
func checksum(file string) string {
	b, err := ioutil.ReadFile(file)
	if err != nil {
		return err.Error()
	}

	return fmt.Sprintf("%x", sha512.Sum512_256(b))
}
