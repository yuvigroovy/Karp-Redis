package main

import (
	"bufio"
	"os"
	"sync"
	"time"
)

type Aof struct {
	file   *os.File
	reader *bufio.Reader
	mutex  sync.Mutex
}

func NewAof(path string) (*Aof, error) {
	file, err := os.OpenFile(path, os.O_CREATE|os.O_RDWR, 0666)

	if err != nil {
		return nil, err
	}

	aof := &Aof{
		file:   file,
		reader: bufio.NewReader(file),
	}

	go func() {
		for {
			aof.mutex.Lock()

			aof.file.Sync()

			aof.mutex.Unlock()

			time.Sleep(time.Second)
		}
	}()

	return aof, nil
}

func (aof *Aof) close() error {
	aof.mutex.Lock()
	defer aof.mutex.Unlock()

	return aof.file.Close()
}

func (aof *Aof) write(value Value) error {
	aof.mutex.Lock()
	defer aof.mutex.Unlock()

	_, err := aof.file.Write(value.marshal())

	if err != nil {
		return err
	}

	return nil
}

func (aof *Aof) read(callback func(value Value)) error {
	aof.mutex.Lock()
	defer aof.mutex.Unlock()

	resp := newResp(aof.file)

	for {
		value, err := read(resp)

		if err != nil { //loop until we get an error (EOF counts as an error)
			return err
		}

		callback(value)
	}
}
