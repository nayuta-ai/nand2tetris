package main

import (
	"bytes"
	"encoding/binary"
	"os"
)

// create Hack file converting from [][]byte
func createHack(dataset [][]byte, filepath string) error {
	file, _ := os.Create(filepath + ".hack")
	defer file.Close()
	for _, data := range dataset {
		buf := new(bytes.Buffer)
		err := binary.Write(buf, binary.BigEndian, data)
		if err != nil {
			return err
		}
		file.Write(buf.Bytes())
	}
	return nil
}
