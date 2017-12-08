package helpers

import (
	"bufio"
	"fmt"
	"os"

	log "github.com/sirupsen/logrus"
)

func readDataFromScanner(scanner *bufio.Scanner) []string {
	result := []string{}
	for scanner.Scan() { // internally, it advances token based on sperator
		line := scanner.Text()
		result = append(result, line)
	}
	return result
}

func ReadData(filename string) ([]string, error) {
	fi, err := os.Stdin.Stat()
	if err != nil {
		panic(err)
	}

	scanner := bufio.NewScanner(os.Stdin)
	if fi.Mode()&os.ModeNamedPipe == 0 {
		if filename == "" {
			log.Debug("no pipe, no file :(")
		} else {
			log.Debug("read data from file")
			fd, err := os.Open(filename)
			if err != nil {
				return nil, fmt.Errorf("open file: reading error")
			}
			defer fd.Close()

			scanner = bufio.NewScanner(fd)
		}
	} else {
		if filename == "" {
			log.Debug("pipe!")
		} else {
			log.Debug("pipe, file skipped")
		}
	}

	result := readDataFromScanner(scanner)
	if err != nil {
		return nil, fmt.Errorf("readDataFromScanner: reading error")
	}
	return result, nil
}
