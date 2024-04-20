package fileio

import (
	"bufio"
	"bytes"
	"encoding/binary"
	"fmt"
	"io"
	"os"
)

func CharacterFrequencies(filename string) (map[rune]int, error) {
	file, err := os.Open(filename)
	if err != nil {
		return nil, fmt.Errorf("failed to open file: %v", err)

	}
	defer file.Close()

	reader := bufio.NewReader(file)
	frequencies := make(map[rune]int)

	for {
		char, _, err := reader.ReadRune()
		if err != nil {
			break
		}
		frequencies[char]++
	}
	return frequencies, nil

}

func WriteHeader(filename string, freqTable map[rune]int) error {
	// Open the file for writing
	file, err := os.Create(filename)
	if err != nil {
		return err
	}

	defer file.Close()

	freqTableLength := len(freqTable)

	if err := binary.Write(file, binary.LittleEndian, int32(freqTableLength)); err != nil {
		return err
	}

	for char, freq := range freqTable {

		if err := binary.Write(file, binary.LittleEndian, int32(char)); err != nil {
			return err
		}

		if err := binary.Write(file, binary.LittleEndian, int32(freq)); err != nil {
			return err
		}

	}

	// Write a delimiter to indicate the end of the header and start of the compressed data
	delimiter := []byte("HEADER_END")

	if _, err := file.Write(delimiter); err != nil {
		return err
	}

	return nil
}

func ReadHeader(filename string) (map[rune]int, error) {

	file, err := os.Open(filename)
	if err != nil {
		return nil, fmt.Errorf("failed to open file: %v", err)
	}

	defer file.Close()

	// Skip the compressed data and find the HEADER_END delimiter
	delimiter := []byte("HEADER_END")
	offset, err := findHeaderEnd(file, delimiter)

	if err != nil {
		return nil, err
	}

	// Read the header section starting from the beginning of the file
	headerBytes := make([]byte, offset)
	_, err = file.Read(headerBytes)
	if err != nil {
		return nil, fmt.Errorf("failed to read header: %v", err)
	}

	// Create a buffer to decode binary data
	buf := headerBytes

	// Read the length of the frequency table
	var freqTableLength int32
	err = binary.Read(bytes.NewReader(buf), binary.LittleEndian, &freqTableLength)
	if err != nil {
		return nil, fmt.Errorf("failed to read frequency table length: %v", err)
	}

	// Read the frequency table
	freqTable := make(map[rune]int)
	for i := 0; i < int(freqTableLength); i++ {
		var char int32
		err = binary.Read(bytes.NewReader(buf), binary.LittleEndian, &char)
		if err != nil {
			return nil, fmt.Errorf("failed to read character from frequency table: %v", err)
		}

		var freq int32
		err = binary.Read(bytes.NewReader(buf), binary.LittleEndian, &freq)
		if err != nil {
			return nil, fmt.Errorf("failed to read frequency from frequency table: %v", err)
		}

		freqTable[rune(char)] = int(freq)
	}

	return freqTable, nil

}

func findHeaderEnd(file *os.File, delimiter []byte) (int64, error) {
	stat, err := file.Stat()

	if err != nil {
		return 0, fmt.Errorf("failed to get file information: %v", err)
	}

	offset := stat.Size()
	const chunkSize = 1024

	for offset > 0 {
		readSize := chunkSize
		if offset < chunkSize {
			readSize = int(offset)
		}

		offset -= int64(readSize)

		_, err := file.Seek(offset, io.SeekStart)
		if err != nil {
			return 0, fmt.Errorf("failed to seek file: %v", err)
		}

		chunk := make([]byte, readSize)
		_, err = file.Read(chunk)
		if err != nil {
			return 0, fmt.Errorf("failed to read file: %v", err)
		}

		if index := bytes.LastIndex(chunk, delimiter); index >= 0 {
			return offset + int64(index), nil
		}
	}

	return 0, fmt.Errorf("HEADER_END delimiter not found")

}
