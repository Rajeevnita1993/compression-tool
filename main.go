package main

import (
	"flag"
	"fmt"
	"os"

	"github.com/Rajeevnita1993/compression-tool/decode"
	"github.com/Rajeevnita1993/compression-tool/encode"
)

func main() {

	// Define command line flags
	encodeCmd := flag.NewFlagSet("encode", flag.ExitOnError)
	decodeCmd := flag.NewFlagSet("decode", flag.ExitOnError)

	// Define flags for encode command
	encodeInput := encodeCmd.String("i", "", "input filename")
	encodeOutput := encodeCmd.String("o", "encoded_file", "output filename for encoding")

	// Define flags for decode command
	decodeInput := decodeCmd.String("i", "", "input filename")
	decodeOutput := decodeCmd.String("o", "decoded_file", "output filename for decoding")

	// Parse the command-line arguments
	flag.Parse()

	if len(os.Args) < 2 {
		fmt.Println("Usage: compression-tool <filename>")
		os.Exit(1)
	}

	// Check which command was provided
	switch os.Args[1] {
	case "encode":
		encodeCmd.Parse(os.Args[2:])
		if *encodeInput == "" {
			fmt.Println("Error: Input filename required for encode.")
			os.Exit(1)
		}
		fmt.Println("Encoding...")
		fmt.Println("Input filename:", *encodeInput)
		fmt.Println("Output filename:", *encodeOutput)
		// Call encode function with *input and *encodeOutput
		err := encode.EncodeFile(*encodeInput, *encodeOutput)
		if err != nil {
			fmt.Println("Error encoding:", err)
			os.Exit(1)
		}
	case "decode":
		decodeCmd.Parse(os.Args[2:])
		if *decodeInput == "" {
			fmt.Println("Error: Input filename required for decode.")
			os.Exit(1)
		}
		fmt.Println("Decoding...")
		fmt.Println("Input filename:", *decodeInput)
		fmt.Println("Output filename:", *decodeOutput)
		// Call decode function with *input and *decodeOutput
		err := decode.DecodeFile(*decodeInput, *decodeOutput)
		if err != nil {
			fmt.Println("Error decoding:", err)
			os.Exit(1)
		}
	default:
		fmt.Println("Unknown command:", os.Args[1])
		os.Exit(1)
	}

}
