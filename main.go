package main

import (
	"bufio"
	"bytes"
	"encoding/json"
	"fmt"
	"log"
	"os"
	"strings"
)

func main() {
	dictionary := NewDictionary("./dictionary.txt")
	scheme := [][]string{
		[]string{"A", "C", "B"},
		[]string{"L", "T", "M"},
		[]string{"M", "N", "O"},
	}
	board := NewBoard(scheme, dictionary)

	board.FindWords()
	fmt.Println(board.Result())
}

type Board struct {
	scheme     [][]string
	dictionary Dictionary
	rows       int
	cols       int
	results    []string
	score      int
	min        int
}

func NewBoard(scheme [][]string, dictionary *Dictionary) *Board {
	var cols int
	for i, row := range scheme {
		for j, val := range row {
			if j > cols {
				cols = j
			}
			scheme[i][j] = strings.ToUpper(val)
		}
	}
	return &Board{
		scheme:     scheme,
		dictionary: *dictionary,
		rows:       len(scheme),
		cols:       cols + 1,
		min:        2,
	}
}

type Coordinates struct {
	x int
	y int
}

type QueueValues struct {
	OldRow  int
	OldCol  int
	OldVal  string
	OldNode *Node
	Passed  []Coordinates
}

func (b *Board) FindWords() {
	queue := b.InitQueue()
	var current QueueValues

	n := [][]int{
		[]int{0, 1},
		[]int{0, -1},
		[]int{1, 0},
		[]int{-1, 0},
		[]int{1, 1},
		[]int{-1, -1},
		[]int{1, -1},
		[]int{-1, 1},
	}

	type neighbour struct {
		position Coordinates
		checked  bool
		value    string
	}

	var neighbours []neighbour

	passed := func(elements []Coordinates, nr int, nc int) bool {
		for _, element := range elements {
			if element.x == nr && element.y == nc {
				return true
			}
		}
		return false
	}

	isNeighbour := func(current QueueValues, potential Coordinates) bool {
		coordinates := current.Passed[len(current.Passed)-1]
		for i := 0; i < len(n); i++ {
			nr := coordinates.x + n[i][0]
			nc := coordinates.y + n[i][1]

			if nr >= potential.x && nc >= potential.y && nr < b.rows && nc < b.cols {
				return !passed(current.Passed, potential.x, potential.y)
			}
		}

		return false
	}

	var buffer bytes.Buffer

	for len(queue.values) > 0 {
		current = queue.Dequeue()
		for i := 0; i < len(n); i++ {
			nr := current.OldRow + n[i][0]
			nc := current.OldCol + n[i][1]

			if nr >= 0 && nc >= 0 && nr < b.rows && nc < b.cols {
				neighbours = append(neighbours, neighbour{position: Coordinates{x: nr, y: nc}, value: b.scheme[nr][nc]})
			}
		}

		for i, cn := range neighbours {
			child := []rune(cn.value)[0]
			if current.OldNode.children[child] != nil && isNeighbour(current, cn.position) {
				buffer.WriteString(current.OldVal)
				buffer.WriteString(cn.value)
				newVal := buffer.String()

				newSeen := append(current.Passed, Coordinates{x: cn.position.x, y: cn.position.y})
				queue.Enqueue(QueueValues{
					OldRow:  cn.position.x,
					OldCol:  cn.position.y,
					OldVal:  newVal,
					OldNode: current.OldNode.children[child],
					Passed:  newSeen,
				})
				if b.dictionary.trie.isWord(newVal) && len(newVal) > b.min {
					b.score += len(newVal) - 2
					b.results = append(b.results, newVal)
				}
				buffer.Reset()
			}
			neighbours[i].checked = true
		}
		neighbours = nil
	}
}

type FormattedResult struct {
	Score int      `json:"score"`
	Words []string `json:"words"`
}

func (b *Board) Result() string {
	result, err := json.Marshal(FormattedResult{b.score, b.results})
	if err != nil {
		return ""
	}

	return string(result)
}

func (b Board) InitQueue() *Queue {
	var queue *Queue
	queue = new(Queue)
	var node *Node

	for r := 0; r < b.rows; r++ {
		for c := 0; c < b.cols; c++ {
			node = findNode(b.dictionary.trie.Root(), []rune(b.scheme[r][c]))
			if node != nil {
				queue.Enqueue(QueueValues{
					OldRow:  r,
					OldCol:  c,
					OldVal:  b.scheme[r][c],
					OldNode: node,
					Passed:  []Coordinates{Coordinates{x: r, y: c}},
				})
			}
		}
	}
	return queue
}

// ============ Dictionary ============

type Dictionary struct {
	trie *Trie
}

// A construction function for Dictionary struct
func NewDictionary(path string) *Dictionary {
	file, err := os.Open(path)
	if err != nil {
		log.Fatal(err)
	}
	defer file.Close()

	trie := NewTrie()
	scanner := bufio.NewScanner(file)

	for scanner.Scan() {
		trie.Add(strings.ToUpper(scanner.Text()))
	}

	return &Dictionary{trie: trie}
}

func (d *Dictionary) WordPresent(word string) bool {
	return d.trie.isWord(word)
}

// ============ Trie definition ============

const nul = 0x0

type Trie struct {
	root *Node
	size int
}

func NewTrie() *Trie {
	return &Trie{
		root: &Node{children: make(map[rune]*Node)},
		size: 0,
	}
}

func (t *Trie) Root() *Node {
	return t.root
}

func (t *Trie) Add(word string) *Node {
	t.size++
	runes := []rune(word)
	node := t.root
	for i := range runes {
		r := runes[i]
		if n, ok := node.children[r]; ok {
			node = n
		} else {
			node = node.AddChild(r, false)
		}
	}
	node = node.AddChild(nul, true)
	return node
}

func (t *Trie) isWord(word string) bool {
	runes := []rune(word)
	node := findNode(t.Root(), runes)

	if node == nil {
		return false
	}

	// return node.children[].leaf
	for _, child := range node.children {
		if child.leaf {
			return true
		}
	}
	return false
}

// ============ Node definition ============

type Node struct {
	val      rune
	leaf     bool
	children map[rune]*Node
}

func (n *Node) AddChild(val rune, leaf bool) *Node {
	node := &Node{
		val:      val,
		leaf:     leaf,
		children: make(map[rune]*Node),
	}
	n.children[val] = node
	return node
}

func (n Node) Children() map[rune]*Node {
	return n.children
}

func findNode(node *Node, runes []rune) *Node {
	if node == nil {
		return nil
	}

	if len(runes) == 0 {
		return node
	}

	n, ok := node.Children()[runes[0]]
	if !ok {
		return nil
	}

	var nrunes []rune
	if len(runes) > 1 {
		nrunes = runes[1:]
	} else {
		nrunes = runes[0:0]
	}

	return findNode(n, nrunes)
}

// ==== Queue ====

type Queue struct {
	values []QueueValues
}

func (q *Queue) Dequeue() QueueValues {
	var value QueueValues
	value, q.values = q.values[0], q.values[1:]
	return value
}

func (q *Queue) Enqueue(value QueueValues) []QueueValues {
	q.values = append(q.values, value)
	return q.values
}

func (q *Queue) isEmpty() bool {
	return len(q.values) == 0
}
