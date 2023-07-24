package config

import (
	"bufio"
	"os"
)

func ReadECRPassword() (string, error) {
	// TODO: 하드 코딩 수정
	filePath := "/app/config/ecr-password"

	file, err := os.Open(filePath)
	if err != nil {
		return "", err
	}
	defer file.Close()

	scanner := bufio.NewScanner(file)
	scanner.Scan()
	content := scanner.Text()

	if err := scanner.Err(); err != nil {
		return "", err
	}

	return content, nil
}
