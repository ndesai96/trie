package trie

import (
	"fmt"
	"testing"

	"github.com/stretchr/testify/assert"
)

var keys = []string{
	"foo",
	"bar",
	"foobar",
	"baz",
	"qux",
	"fizz",
}

func TestTrie_Insert(t *testing.T) {
	trie := NewTrie()

	for idx, key := range keys {
		testName := fmt.Sprintf("Insert %s", key)
		t.Run(testName, func(t *testing.T) {
			trie.Insert(key)

			currentNode := trie.root
			for _, char := range key {
				if node, ok := currentNode.children[char]; ok {
					currentNode = node
				} else {
					t.Fatal("Node path not found for key:", key)
				}
			}

			if !currentNode.isKey {
				t.Fatal("Last node does not have isKey set for key:", key)
			}

			assert.Equal(t, idx+1, trie.Len())
		})
	}
}

func initializeTrie() *Trie {
	trie := NewTrie()

	for _, key := range keys {
		trie.Insert(key)
	}

	return trie
}

func TestTrie_Search(t *testing.T) {
	trie := initializeTrie()

	type expectation struct {
		key   string
		found bool
	}

	var expectations []expectation

	for _, k := range keys {
		expectations = append(expectations, expectation{key: k, found: true})
	}

	expectations = append(expectations, expectation{key: "none", found: false})

	for _, e := range expectations {
		testName := fmt.Sprintf("Key %s", e.key)
		t.Run(testName, func(t *testing.T) {
			found := trie.Search(e.key)

			assert.Equal(t, e.found, found)
		})
	}
}

func TestTrie_FindAllWithPrefix(t *testing.T) {
	trie := initializeTrie()

	type expectation struct {
		prefix       string
		expectedKeys []string
	}

	expectations := []expectation{
		{prefix: "f", expectedKeys: []string{"foo", "foobar", "fizz"}},
		{prefix: "fo", expectedKeys: []string{"foo", "foobar"}},
		{prefix: "foo", expectedKeys: []string{"foo", "foobar"}},
		{prefix: "foob", expectedKeys: []string{"foobar"}},
		{prefix: "ba", expectedKeys: []string{"bar", "baz"}},
		{prefix: "bar", expectedKeys: []string{"bar"}},
		{prefix: "q", expectedKeys: []string{"qux"}},
		{prefix: "qux", expectedKeys: []string{"qux"}},
		{prefix: "none", expectedKeys: []string{}},
	}

	for _, e := range expectations {
		testName := fmt.Sprintf("Prefix %s", e.prefix)
		t.Run(testName, func(t *testing.T) {
			keysWithPrefix := trie.FindAllWithPrefix(e.prefix)

			assert.ElementsMatch(t, e.expectedKeys, keysWithPrefix)
		})
	}
}

func TestTrie_Serialize(t *testing.T) {
	trie := initializeTrie()

	var serializedTrie string
	trie.Serialize(&serializedTrie)

	assert.Equal(t, 39, len(serializedTrie))
}

func TestTrie_Len(t *testing.T) {
	trie := initializeTrie()

	assert.Equal(t, len(keys), trie.Len())

	t.Run("Duplicate key inserted", func(t *testing.T) {
		trie.Insert(keys[0])

		assert.Equal(t, len(keys), trie.Len())
	})
}

func TestDeserialize(t *testing.T) {
	serializedTrie := "foo]bar]>>>>>izz]>>>>bar]>z]>>>qux]>>>>"

	trie := NewTrie()
	Deserialize(trie, serializedTrie)

	for _, key := range keys {
		assert.True(t, trie.Search(key))
	}

	assert.Equal(t, len(keys), trie.Len())
}
