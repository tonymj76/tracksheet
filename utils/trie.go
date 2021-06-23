package main // utils

import (
	"fmt"
	"reflect"
	"strings"
)

type TrieNode struct {
	key            string
	numOfOccurence int
	children       map[string]*TrieNode
	isWord         bool
	numOfWords     int
}

type TrieChildren map[string]*TrieNode

func newNode(key string) *TrieNode {
	return &TrieNode{key: key, children: map[string]*TrieNode{}}
}

func PrefixTree() *TrieNode {
	return newNode("root")
}

func (trie *TrieNode) Add(values []string) {
	trie.numOfOccurence += 1

	//if no spereator, set word value to true
	if len(values) == 0 {
		trie.isWord = true
		return
	}

	//get parent and subsequent child words
	value, values := values[0], values[1:]

	//if parent key word exist, append children words to it
	if trie.key == value {
		trie.Add(values)

		return
	}

	//append children to parent node if it does not exist
	if _, ok := trie.children[value]; !ok {
		trie.children[value] = newNode(value)
	}

	trie.children[value].Add(values)
}

func merge(m1, m2 map[string]*TrieNode) map[string]*TrieNode {
	for m, v := range m2 {
		m1[m] = v
	}

	return m1
}

func (trie *TrieNode) Paternify() {
	common := []*TrieNode{}

	for _, current := range trie.children {
		//if node has an occurance of 1 or less, append to the common slice
		if current.numOfOccurence <= 1 {
			common = append(common, current)
		}
	}

	if len(common) > 1 {
		//pick the first item from the common as key
		pattern := newNode(common[0].key)

		for _, node := range common {
			pattern.numOfOccurence += 1
			pattern.children = merge(pattern.children, node.children)

			pattern.numOfWords += 1
			if !pattern.isWord {
				pattern.isWord = node.isWord
			}

			delete(trie.children, node.key)
		}
		trie.children[common[0].key] = pattern
	}

	for _, current := range trie.children {
		current.Paternify()
	}
}

//get words after transversing the tree
func (trie *TrieNode) getWords(words []string, word string, sep string) []string {
	if trie.isWord {
		words = append(words, word)
	}

	for _, c := range trie.children {
		words = c.getWords(words, word+c.key+sep, sep)
	}
	return words
}

func (trie *TrieNode) GetWords(sep string) []string {
	trie.Paternify()
	return trie.getWords([]string{}, "", sep)
}

func (trie *TrieNode) String() string {
	children := reflect.ValueOf(trie.children).MapKeys()
	fmt.Println(trie.key, children, trie.numOfOccurence, trie.isWord)
	for _, t := range trie.children {
		t.String()
	}

	return ""
}

func GetUnique(strs []string, sep string) []string {
	fmt.Println(len(strs))

	trie := PrefixTree()

	for _, str := range strs {
		trie.Add(strings.Split(str, sep))
	}

	return trie.GetWords(sep)
}

func main() {
	urls := []string{
		"https://amazon.com/product/123333",
		"https://amazon.com/product/456899",
		"https://amazon.com/gp/5565565",
		"https://amazon.com/gp/556556",
		"https://amazon.com/gp/55655",
		"https://amazon.com/gp/5565565/abc",
		"https://amazon.com/gp/5565565/abc/ref=ab",
	}
	trie := PrefixTree()

	for _, url := range urls {
		trie.Add(strings.Split(url, "/"))
	}
	fmt.Println(trie.GetWords("/"))
}
