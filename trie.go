package trie

import (
	"strings"

	"github.com/golang-collections/collections/stack"
)

type Trie struct {
	root *node
	size int
}

func NewTrie() *Trie {
	return &Trie{
		root: newNode(),
		size: 0,
	}
}

type node struct {
	isKey    bool
	children map[rune]*node
}

func newNode() *node {
	return &node{
		isKey:    false,
		children: make(map[rune]*node),
	}
}

func (t *Trie) Len() int {
	return t.size
}

func (t *Trie) Insert(key string) {
	currentNode := t.root
	for _, char := range strings.ToLower(key) {
		if childNode, ok := currentNode.children[char]; ok {
			currentNode = childNode
		} else {
			newNode := newNode()
			currentNode.children[char] = newNode
			currentNode = newNode
		}
	}

	// Check if key already exists
	if !currentNode.isKey {
		currentNode.isKey = true
		t.size++
	}
}

func (t *Trie) searchNode(s string) *node {
	currentNode := t.root
	for _, char := range strings.ToLower(s) {
		if childNode, ok := currentNode.children[char]; ok {
			currentNode = childNode
		} else {
			return nil
		}
	}

	return currentNode
}

func (t *Trie) Search(key string) bool {
	node := t.searchNode(key)
	return node != nil && node.isKey
}

type nodePrefix struct {
	node   *node
	prefix string
}

func newNodePrefix(node *node, prefix string) *nodePrefix {
	return &nodePrefix{
		node:   node,
		prefix: prefix,
	}
}

func (t *Trie) FindAllWithPrefix(prefix string) []string {
	var keys []string

	node := t.searchNode(prefix)
	if node == nil {
		return []string{}
	}

	nodePrefixStack := stack.New()
	nodePrefixStack.Push(newNodePrefix(node, prefix))

	for nodePrefixStack.Len() > 0 {
		nodePrefix := nodePrefixStack.Pop().(*nodePrefix)
		if nodePrefix.node.isKey {
			keys = append(keys, nodePrefix.prefix)
		}

		if len(nodePrefix.node.children) > 0 {
			for char, childNode := range nodePrefix.node.children {
				newPrefix := nodePrefix.prefix + string(char)
				nodePrefixStack.Push(newNodePrefix(childNode, newPrefix))
			}
		}
	}

	return keys
}

func (t *Trie) Serialize(serializedTrie *string) {
	*serializedTrie = ""

	var serialize func(node *node)

	serialize = func(node *node) {
		if node.isKey {
			*serializedTrie += "]"
		}
		for char, childNode := range node.children {
			*serializedTrie += string(char)
			serialize(childNode)
		}
		*serializedTrie += ">"
	}

	serialize(t.root)
}

func Deserialize(trie *Trie, serializedTrie string) {
	nodeStack := stack.New()
	nodeStack.Push(trie.root)

	for _, char := range serializedTrie {
		if char == ']' {
			nodeStack.Peek().(*node).isKey = true
			trie.size++
		} else if char == '>' {
			nodeStack.Pop()
		} else {
			newNode := newNode()
			nodeStack.Peek().(*node).children[char] = newNode
			nodeStack.Push(newNode)
		}
	}
}
