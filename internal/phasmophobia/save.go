package phasmophobia

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"log"
	"os"
	"time"

	cp "github.com/nmrshll/go-cp"
	"github.com/pkg/errors"
)

// ReadSave is the function that will decode the save file and load the json in-memory
func ReadSave(path string) (*Save, error) {
	save := &Save{Path: path}

	file, err := os.Open(path)
	if err != nil {
		return nil, err
	}
	defer file.Close()

	data, err := ioutil.ReadAll(file)
	if err != nil {
		return nil, err
	}

	// Decode save file, this may need rewriting in future releases ( works with  0.176.39 and prior )
	decodedSave := xor(data, salt)

	// Load JSON
	err = json.Unmarshal([]byte(decodedSave), save)
	if err != nil {
		log.Printf("Unable to decode json from %s ( showing 100 first characters )", decodedSave[:100])
		return nil, errors.Wrapf(err, "invalid json in file %s", path)
	}

	return save, nil
}

// Save overrides the current save game but creates a backup file first
func (s Save) Save() error {
	// Backup previous save
	err := cp.CopyFile(s.Path, s.Path+fmt.Sprintf("-%d", time.Now().UnixNano()))
	if err != nil {
		log.Fatal(err)
	}

	// Now we can override our save
	file, err := os.Create(s.Path)
	if err != nil {
		return err
	}
	defer file.Close()

	// Marshal our save
	bytes, err := json.Marshal(s)
	if err != nil {
		return err
	}

	// Override with xor
	_, err = file.Write(xor(bytes, salt))

	return err
}

// Simple XOR with passphrase as salt
func xor(data, salt []byte) []byte {
	decoded := []byte{}
	for i, c := range data {
		decoded = append(decoded, c^salt[i%len(salt)])
	}
	return decoded
}
