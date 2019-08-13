package customerimporter

import (
	"bufio"
	"errors"
	"os"
	"strings"
)

func findDomain(text string, emailPosition int) (string, error) {
	lineValues := strings.Split(text, ",")
	if len(lineValues) > emailPosition {
		emailValues := strings.Split(lineValues[emailPosition], "@")
		if len(emailValues) == 2 {
			return emailValues[1], nil
		}
	}
	return "", errors.New("Can't find a valid domain")
}

// ReadFile recives the file name with customers details.
// returns a sorted map of email domains along with the number
// of customers with e-mail addresses for each domain.
func ReadFile(fileName string) (map[string]int, error) {

	// open the file and prefer the scanner
	file, err := os.Open(fileName)
	if err != nil {
		return nil, err
	}
	defer file.Close()
	scanner := bufio.NewScanner(file)

	emailPosition := 2
	mapOfDomains := make(map[string]int)

	// read the file to the end
	for scanner.Scan() {
		domain, err := findDomain(scanner.Text(), emailPosition)
		if err == nil {
			mapOfDomains[domain]++
		}
	}

	if err := scanner.Err(); err != nil {
		return nil, err
	}
	return mapOfDomains, nil
}
