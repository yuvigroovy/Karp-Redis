package main

import "io"

type Writer struct {
	writer io.Writer
}

func NewWriter(writer io.Writer) *Writer {
	return &Writer{writer: writer}
}

func (writer *Writer) write(value Value) error {
	var bytes = value.marshal()

	_, err := writer.writer.Write(bytes)
	if err != nil {
		return err
	}

	return nil
}
