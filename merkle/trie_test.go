package merkle

import (
	"testing"

	"github.com/wojtechnology/medblocks/test"
)

func buildTrie() *MerkleTrie {
	trie := &MerkleTrie{}

	leaf := &MerkleLeafNode{
		key: []byte{0, 1, 2, 3, 4},
		val: "someValue",
	}
	innerLeaf := &MerkleLeafNode{
		key: []byte{0, 1, 2},
		val: "someValueInner",
	}

	innerBranch := &MerkleBranchNode{
		keyPrefix: []byte{0, 1, 2},
		innerLeaf: innerLeaf,
	}
	innerBranch.children[3] = leaf
	branch := &MerkleBranchNode{
		keyPrefix: []byte{0},
	}
	branch.children[1] = innerBranch
	trie.root = branch

	return trie
}

func buildRootLeafTrie() *MerkleTrie {
	trie := &MerkleTrie{}

	leaf := &MerkleLeafNode{
		key: []byte{0, 1, 2},
		val: "someValue",
	}
	trie.root = leaf

	return trie
}

func buildEmptyTrie() *MerkleTrie {
	return &MerkleTrie{}
}

// -------------
// Test Contains
// -------------

func TestDoesNotContainEmpty(t *testing.T) {
	trie := buildTrie()
	test.AssertEqual(t, false, trie.Contains([]byte{}))
}

func TestDoesNotContainInnerLeaf(t *testing.T) {
	trie := buildTrie()
	test.AssertEqual(t, false, trie.Contains([]byte{0}))
}

func TestContainsInnerLeaf(t *testing.T) {
	trie := buildTrie()
	test.AssertEqual(t, true, trie.Contains([]byte{0, 1, 2}))
}

func TestContainsLeaf(t *testing.T) {
	trie := buildTrie()
	test.AssertEqual(t, true, trie.Contains([]byte{0, 1, 2, 3, 4}))
}

func TestDoesNotContainLong(t *testing.T) {
	trie := buildTrie()
	test.AssertEqual(t, false, trie.Contains([]byte{0, 1, 2, 3, 4, 5}))
}

func TestDoesNotContainMissingChild(t *testing.T) {
	trie := buildTrie()
	test.AssertEqual(t, false, trie.Contains([]byte{0, 1, 4}))
}

func TestDoesNotContainRootLeaf(t *testing.T) {
	trie := buildRootLeafTrie()
	test.AssertEqual(t, false, trie.Contains([]byte{0, 1}))
}

func TestContainsRootLeaf(t *testing.T) {
	trie := buildRootLeafTrie()
	test.AssertEqual(t, true, trie.Contains([]byte{0, 1, 2}))
}

func TestContainsEmptyTrie(t *testing.T) {
	trie := buildEmptyTrie()
	test.AssertEqual(t, false, trie.Contains([]byte{}))
	test.AssertEqual(t, false, trie.Contains([]byte{1}))
}

// ----------------
// Test Add and Get
// ----------------

func TestAddToEmpty(t *testing.T) {
	trie := buildEmptyTrie()
	key := []byte{1}
	val := "someValue"

	test.AssertEqual(t, true, trie.Add(key, val))
	test.AssertEqual(t, val, trie.Get(key))
}

func TestAddBranchToEmpty(t *testing.T) {
	trie := buildEmptyTrie()
	key1 := []byte{1, 2}
	val1 := "someValue1"
	key2 := []byte{1, 4}
	val2 := "someValue2"

	test.AssertEqual(t, true, trie.Add(key1, val1))
	test.AssertEqual(t, true, trie.Add(key2, val2))
	test.AssertEqual(t, val1, trie.Get(key1))
	test.AssertEqual(t, val2, trie.Get(key2))
}

func testAddAndGet(t *testing.T, trie *MerkleTrie, key []byte) {
	val := "someValue"
	test.AssertEqual(t, true, trie.Add(key, val))
	test.AssertEqual(t, val, trie.Get(key))
}

func TestAddInnerLeaf(t *testing.T) {
	testAddAndGet(t, buildTrie(), []byte{0})
}

func TestAddToNewBranch(t *testing.T) {
	testAddAndGet(t, buildTrie(), []byte{0, 1, 3, 4})
}

// New node becomes inner leaf
func TestAddToNewBranchInnerLeaf1(t *testing.T) {
	testAddAndGet(t, buildTrie(), []byte{0, 1})
}

// Existing node becomes inner leaf
func TestAddToNewBranchInnerLeaf2(t *testing.T) {
	testAddAndGet(t, buildTrie(), []byte{0, 1, 2, 3, 4, 5})
}

func TestAddToExistingBranch(t *testing.T) {
	testAddAndGet(t, buildTrie(), []byte{0, 2, 3, 4})
}

func testAddAlreadyExists(t *testing.T, trie *MerkleTrie, key []byte) {
	val := "someValue"
	test.AssertEqual(t, false, trie.Add(key, val))
}

func TestAddAlreadyExistsRootLeaf(t *testing.T) {
	testAddAlreadyExists(t, buildRootLeafTrie(), []byte{0, 1, 2})
}

func TestAddAlreadyExistsLeaf(t *testing.T) {
	testAddAlreadyExists(t, buildTrie(), []byte{0, 1, 2, 3, 4})
}

func TestAddAlreadyInnerLeaf(t *testing.T) {
	testAddAlreadyExists(t, buildTrie(), []byte{0, 1, 2})
}