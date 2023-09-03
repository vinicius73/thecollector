package checksum

import (
	"fmt"
	"os"
	"path/filepath"

	"github.com/codingsince1985/checksum"
)

const sumFileName = "SHA256SUMS"

func GenerateChecksum(fileName string) (string, error) {
	sha256, err := checksum.SHA256sum(fileName)
	if err != nil {
		return sha256, err
	}

	line := fmt.Sprintf("%s *%s\n", sha256, filepath.Base(fileName))

	sumFile := filepath.Join(filepath.Dir(fileName), sumFileName)

	file, err := os.OpenFile(sumFile, os.O_APPEND|os.O_CREATE|os.O_WRONLY, 0o644)
	if err != nil {
		return sha256, err
	}

	defer file.Close()

	_, err = file.WriteString(line)

	return sha256, err
}
