package util

import (
	"log"
	"os"
)

func AppendToFile(fileName string, content []byte) error {
	f, err := os.OpenFile(fileName, os.O_WRONLY|os.O_APPEND, 0644)
	if err != nil {
		log.Println("file created or append failed! err:" + err.Error())
	}
	_, err = f.Write(content)

	defer f.Close()

	return err
}
