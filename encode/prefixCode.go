package encode

import (
	"fmt"
	"os"
)

type PrefixCodeTable map[rune]string

func GeneratePrefixCodeTable(root *Node) PrefixCodeTable {
	prefixCodes := make(PrefixCodeTable)
	traverse(root, "", prefixCodes)
	return prefixCodes
}

func traverse(node *Node, code string, prefixCodes PrefixCodeTable) {

	if node == nil {
		return
	}

	// if node is leaf node
	if node.left == nil && node.right == nil {
		prefixCodes[node.char] = code
		return
	}

	traverse(node.left, code+"0", prefixCodes)
	traverse(node.right, code+"1", prefixCodes)

}

func EncodeAndWrite(filename string, text string, codeTable PrefixCodeTable) error {
	// Open the file for writing in append mode
	file, err := os.OpenFile(filename, os.O_APPEND|os.O_WRONLY, os.ModeAppend)
	if err != nil {
		return err
	}

	defer file.Close()
	var encodedBytes []byte

	for _, char := range text {
		code, ok := codeTable[char]

		if !ok {
			return fmt.Errorf("character '%c' not found in code table", char)
		}

		for _, bit := range code {
			if bit == '0' {
				encodedBytes = append(encodedBytes, 0)
			} else if bit == '1' {
				encodedBytes = append(encodedBytes, 1)
			} else {
				return fmt.Errorf("invalid bit in code: %c", bit)
			}
		}
	}

	// Pack the bits into bytes
	packedBytes := packBits(encodedBytes)

	// Write the packed bytes into file
	_, err = file.Write(packedBytes)
	if err != nil {
		return err
	}

	return nil

}

func packBits(bits []byte) []byte {
	packed := make([]byte, (len(bits)+7)/8)
	for i := range bits {
		if bits[i] == 1 {
			packed[i/8] |= 1 << uint(7-i%8)
		} else if bits[i] != 0 {
			panic(fmt.Errorf("invalid bit encountered: %d", bits[i]))
		}
	}
	return packed
}
