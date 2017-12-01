package repot

import (
	"bufio"
	"crypto/rand"
	"errors"
	"fmt"
	"io"
	"os"
	"sort"
	"strconv"
	"strings"

	log "github.com/sirupsen/logrus"
)

// newUUID generates a random UUID according to RFC 4122
func UUID() (string, error) {
	uuid := make([]byte, 16)
	n, err := io.ReadFull(rand.Reader, uuid)
	if n != len(uuid) || err != nil {
		return "", err
	}
	// variant bits; see section 4.1.1
	uuid[8] = uuid[8]&^0xc0 | 0x80
	// version 4 (pseudo-random); see section 4.1.3
	uuid[6] = uuid[6]&^0xf0 | 0x40
	return fmt.Sprintf("%x-%x-%x-%x-%x", uuid[0:4], uuid[4:6], uuid[6:8], uuid[8:10], uuid[10:]), nil
}

func keys(m map[int]bool) []int {
	keys := make([]int, len(m))

	i := 0
	for k := range m {
		keys[i] = k
		i++
	}
	sort.Ints(keys)
	return keys
}

func ParseRangesString(ranges string) ([]int, error) {
	repos := map[int]bool{}

	ranges = strings.TrimSpace(ranges)
	if len(ranges) == 0 {
		return []int{}, errors.New("empty input")
	}
	segments := strings.SplitN(ranges, ",", -1)
	for _, r := range segments {
		if rr := strings.SplitN(r, "-", -1); len(rr) > 1 {
			start, start_err := strconv.Atoi(rr[0])
			finish, finish_err := strconv.Atoi(rr[1])
			if start_err != nil || finish_err != nil {
				return []int{}, errors.New("range error")
			}

			if finish <= start {
				continue
			}
			for i := start; i <= finish; i++ {
				repos[i] = true
			}
		} else {
			if val, err := strconv.Atoi(r); err != nil {
				return []int{}, err
			} else {
				repos[val] = true
			}
		}
	}
	return keys(repos), nil
}

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
