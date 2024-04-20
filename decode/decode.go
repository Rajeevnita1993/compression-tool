package decode

import (
	"bufio"
	"fmt"
	"os"

	"github.com/Rajeevnita1993/compression-tool/encode"
	"github.com/Rajeevnita1993/compression-tool/fileio"
)

func DecodeFile(encodedFilename string, outputFilename string) error {
	// Open the encoded file for reading
	encodedFile, err := os.Open(encodedFilename)
	if err != nil {
		return fmt.Errorf("failed to open encoded file: %v", err)
	}
	defer encodedFile.Close()

	// Read the header to reconstruct the prefix code table
	frequencies, err := fileio.ReadHeader(encodedFilename)
	fmt.Println("frquencies: ", frequencies)
	if err != nil {
		return fmt.Errorf("failed to read header: %v", err)
	}

	// Build the Huffman tree from the frequencies
	root := encode.BuildHuffmanTree(frequencies)

	// Generate the prefix code table from the Huffman tree
	prefixCodes := encode.GeneratePrefixCodeTable(root)
	fmt.Println("prefixCodes: ", prefixCodes)

	// Open the output file for writing
	outputFile, err := os.Create(outputFilename)
	if err != nil {
		return fmt.Errorf("failed to create output file: %v", err)
	}
	defer outputFile.Close()

	// Decode the remainder of the encoded file using the prefix code table
	if err := DecodeAndWrite(encodedFile, outputFile, prefixCodes); err != nil {
		return fmt.Errorf("failed to decode and write: %v", err)
	}

	return nil
}

func DecodeAndWrite(encodedFile, outputFile *os.File, prefixCodes encode.PrefixCodeTable) error {
	reader := bufio.NewReader(encodedFile)
	var bits []byte
	for {
		bit, err := reader.ReadByte()
		if err != nil {
			break // End of file
		}
		bits = append(bits, bit)
	}
	fmt.Println("bits: ", bits)

	// Decode the bits using the prefix code table
	decodedText, err := DecodeBits(bits, prefixCodes)
	if err != nil {
		fmt.Println("Error: ", err)
	}

	fmt.Println("decodedText: ", decodedText)

	//Write the decoded text to the output file
	_, err = outputFile.WriteString(decodedText)
	if err != nil {
		return fmt.Errorf("failed to write decoded text to output file: %v", err)
	}

	return nil
}

func DecodeBits(bits []byte, prefixCodes encode.PrefixCodeTable) (string, error) {
	var decodedText string
	var currentCode string

	for _, bit := range bits {
		if bit != 0 && bit != 1 {
			return "", fmt.Errorf("invalid bit encountered: %d", bit)
		}

		// Append the current bit to the current code
		currentCode += fmt.Sprintf("%d", bit)

		// Check if the current code matches any prefix code
		for char, code := range prefixCodes {
			if code == currentCode {
				// Found a match, append the corresponding character to the decoded text
				decodedText += string(char)
				currentCode = "" // Reset the current code
				break
			}
		}
	}

	return decodedText, nil
}
