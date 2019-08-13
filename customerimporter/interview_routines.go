package customerimporter

import (
	"bufio"
	"bytes"
	"os"
	"strings"
	"sync"
)

func findDomainWithGoroutines(buf <-chan string, results chan<- string, wg *sync.WaitGroup, emailPosition int) {
	defer wg.Done()

	for lines := range buf {
		scanner := bufio.NewScanner(strings.NewReader(lines))
		for scanner.Scan() {
			lineValues := strings.Split(scanner.Text(), ",")
			if len(lineValues) > emailPosition {
				emailValues := strings.Split(lineValues[emailPosition], "@")
				if len(emailValues) == 2 {
					results <- emailValues[1]
				}
			}
		}
	}
}

// https://stackoverflow.com/a/57232670/106669
func splitAt(substring string) func(data []byte, atEOF bool) (advance int, token []byte, err error) {
	searchBytes := []byte(substring)
	searchLen := len(searchBytes)
	return func(data []byte, atEOF bool) (advance int, token []byte, err error) {
		dataLen := len(data)

		// Return nothing if at end of file and no data passed
		if atEOF && dataLen == 0 {
			return 0, nil, nil
		}

		// Find next separator and return token
		if i := bytes.Index(data, searchBytes); i >= 0 {
			return i + searchLen, data[0:i], nil
		}

		// If we're at EOF, we have a final, non-terminated line. Return it.
		if atEOF {
			return dataLen, data, nil
		}

		// Request more data.
		return 0, nil, nil
	}
}

func calcTotalOfGoroutines(size int64) int {
	totalOfGoroutines := size / 50
	if totalOfGoroutines > 1000 {
		return int(totalOfGoroutines / 1000)
	}
	return 2
}

// ReadFileWithGoroutines recives the file name with customers details.
// returns a sorted map of email domains along with the number
// of customers with e-mail addresses for each domain.
func ReadFileWithGoroutines(fileName string) (map[string]int, error) {

	// open the file and prefer the scanner
	file, err := os.Open(fileName)
	if err != nil {
		return nil, err
	}
	defer file.Close()
	scanner := bufio.NewScanner(file)
	scanner.Split(splitAt("\nfirst_name,last_name,email,gender,ip_address"))

	// check the file' size to calculte the total of routines
	fileinfo, err := file.Stat()
	if err != nil {
		return nil, err
	}
	totalOfGoroutines := calcTotalOfGoroutines(fileinfo.Size())

	emailPosition := 2
	mapOfDomains := make(map[string]int)
	buf := make(chan string)
	results := make(chan string)

	// call the find func using routines and waitGroup
	wg := new(sync.WaitGroup)
	for i := 0; i < totalOfGoroutines; i++ {
		wg.Add(1)
		go findDomainWithGoroutines(buf, results, wg, emailPosition)
	}

	// read the file to the end
	go func() {
		for scanner.Scan() {
			buf <- scanner.Text()
		}
		close(buf)
	}()

	go func() {
		wg.Wait()
		close(results)
	}()

	for domain := range results {
		mapOfDomains[domain]++
	}

	if err := scanner.Err(); err != nil {
		return nil, err
	}

	return mapOfDomains, nil
}
