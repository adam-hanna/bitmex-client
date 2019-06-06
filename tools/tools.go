package tools

import (
	"encoding/gob"
	"log"
	"os"
)

// WriteGob ...
func WriteGob(filePath string, object interface{}) error {
	file, err := os.Create(filePath)
	if err != nil {
		log.Printf("err creating file %s:\n%v", filePath, err)
		return err
	}

	encoder := gob.NewEncoder(file)
	if err = encoder.Encode(object); err != nil {
		log.Printf("err encoding object:\n%v", err)
		return err
	}

	return file.Close()
}

// ReadGob ...
func ReadGob(filePath string, object interface{}) error {
	file, err := os.Open(filePath)
	if err != nil {
		log.Printf("err opening file %s:\n%v", filePath, err)
		return err
	}

	decoder := gob.NewDecoder(file)
	if err = decoder.Decode(object); err != nil {
		log.Printf("err decoding:\n%v", err)
		return err
	}

	return file.Close()
}
