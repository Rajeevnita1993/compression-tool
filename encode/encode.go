package encode

import (
	"container/heap"
	"fmt"
	"os"

	"github.com/Rajeevnita1993/compression-tool/fileio"
)

type Node struct {
	char  rune
	freq  int
	left  *Node
	right *Node
}

type PriorityQueue []*Node

func (pq PriorityQueue) Len() int {
	return len(pq)
}

func (pq PriorityQueue) Less(i, j int) bool {
	return pq[i].freq < pq[j].freq
}

func (pq PriorityQueue) Swap(i, j int) {
	pq[i], pq[j] = pq[j], pq[i]
}

func (pq *PriorityQueue) Push(x interface{}) {
	item := x.(*Node)
	*pq = append(*pq, item)
}

func (pq *PriorityQueue) Pop() interface{} {
	old := *pq
	n := len(old)
	item := old[n-1]
	*pq = old[0 : n-1]
	return item
}

func BuildHuffmanTree(frequencies map[rune]int) *Node {
	// Create leaf nodes for each character & its frequency
	var nodes PriorityQueue
	for char, freq := range frequencies {

		nodes = append(nodes, &Node{char: char, freq: freq})
	}

	heap.Init(&nodes)

	// Build the huffman tree
	for len(nodes) > 1 {
		// deque the two nodes with lowest frequencies
		left := heap.Pop(&nodes).(*Node)
		right := heap.Pop(&nodes).(*Node)

		// Create a new internam node
		internal := &Node{char: ' ', freq: left.freq + right.freq, left: left, right: right}

		// Enqueue the new internam node into priority queue
		heap.Push(&nodes, internal)
	}

	return nodes[0]
}

func EncodeFile(inputFile, outputFile string) error {

	frequencies, err := fileio.CharacterFrequencies(inputFile)

	if err != nil {
		fmt.Println("Error: ", err)
		os.Exit(1)
	}

	node := BuildHuffmanTree(frequencies)
	fmt.Println("node: ", node)

	prefixCodes := GeneratePrefixCodeTable(node)
	fmt.Println("prefixCodes: ", prefixCodes)

	fileio.WriteHeader(outputFile, frequencies)

	data, err := os.ReadFile(inputFile)

	if err != nil {
		fmt.Println("Error: ", err)
		os.Exit(2)
	}

	EncodeAndWrite(outputFile, string(data), prefixCodes)

	return nil

}
