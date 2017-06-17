package meddb

import (
	"bytes"
	"sort"
	"sync"
)

// In-memory blockchain db mainly meant for testing
type MemoryBlockchainDB struct {
	backlogTable map[string]*Transaction
	backlogLock  sync.RWMutex
	blockTable   map[string]*Block
	blockLock    sync.RWMutex
	voteTable    map[string]*Vote
	voteLock     sync.RWMutex
}

// ----------------------
// MemoryBlockchainDB API
// ----------------------

func NewMemoryBlockchainDB() (*MemoryBlockchainDB, error) {
	return &MemoryBlockchainDB{
		backlogTable: make(map[string]*Transaction),
		blockTable:   make(map[string]*Block),
		voteTable:    make(map[string]*Vote),
	}, nil
}

func (db *MemoryBlockchainDB) SetupTables() error {
	return nil
}

func (db *MemoryBlockchainDB) WriteTransaction(tx *Transaction) error {
	db.backlogLock.Lock()
	defer db.backlogLock.Unlock()

	db.backlogTable[string(tx.Hash)] = tx.Clone()
	return nil
}

// Note: This is not performant, do not use in prod
func (db *MemoryBlockchainDB) GetAssignedTransactions(pubKey []byte) ([]*Transaction, error) {
	db.backlogLock.Lock()
	defer db.backlogLock.Unlock()

	txs := make([]*Transaction, 0)
	for _, tx := range db.backlogTable {
		if bytes.Equal(tx.AssignedTo, pubKey) {
			txs = append(txs, tx.Clone())
		}
	}

	return txs, nil
}

// Note: This is not performant, do not use in prod
func (db *MemoryBlockchainDB) GetStaleTransactions(before int64) ([]*Transaction, error) {

	db.backlogLock.Lock()
	defer db.backlogLock.Unlock()

	txs := make([]*Transaction, 0)
	for _, tx := range db.backlogTable {
		if tx.AssignedAt != nil && tx.AssignedAt.Int64() <= before {
			txs = append(txs, tx.Clone())
		}
	}

	return txs, nil
}

func (db *MemoryBlockchainDB) DeleteTransactions(txs []*Transaction) error {
	db.backlogLock.Lock()
	defer db.backlogLock.Unlock()

	for _, tx := range txs {
		delete(db.backlogTable, string(tx.Hash))
	}

	return nil
}

func (db *MemoryBlockchainDB) WriteBlock(b *Block) error {
	db.blockLock.Lock()
	defer db.blockLock.Unlock()

	db.blockTable[string(b.Hash)] = b.Clone()
	return nil
}

func (db *MemoryBlockchainDB) GetOldestBlocks(start int64, limit int) ([]*Block, error) {
	db.blockLock.Lock()
	defer db.blockLock.Unlock()

	candidates := make([]*Block, 0)
	for _, b := range db.blockTable {
		if b.CreatedAt != nil && b.CreatedAt.Int64() >= start {
			candidates = append(candidates, b.Clone())
		}
	}
	sort.Slice(candidates, func(i, j int) bool {
		// None of the CreatedAt will be nil
		return candidates[i].CreatedAt.Int64() < candidates[j].CreatedAt.Int64()
	})

	return candidates[:limit], nil
}

func (db *MemoryBlockchainDB) WriteVote(v *Vote) error {
	db.voteLock.Lock()
	defer db.voteLock.Unlock()

	db.voteTable[string(v.Hash)] = v.Clone()
	return nil
}

func (db *MemoryBlockchainDB) GetVotes(pubKey []byte, votedAt int64) ([]*Vote, error) {
	db.voteLock.Lock()
	defer db.voteLock.Unlock()

	vs := make([]*Vote, 0)
	for _, v := range db.voteTable {
		if bytes.Equal(v.Voter, pubKey) && v.VotedAt != nil && v.VotedAt.Int64() == votedAt {
			vs = append(vs, v.Clone())
		}
	}

	return vs, nil
}

func (db *MemoryBlockchainDB) GetRecentVotes(pubKey []byte, limit int) ([]*Vote, error) {
	db.voteLock.Lock()
	defer db.voteLock.Unlock()

	candidates := make([]*Vote, 0)
	for _, v := range db.voteTable {
		if bytes.Equal(v.Voter, pubKey) {
			candidates = append(candidates, v.Clone())
		}
	}
	sort.Slice(candidates, func(i, j int) bool {
		// Greater than since we want a reverse sort
		if candidates[i].VotedAt == nil {
			return false
		}
		if candidates[j].VotedAt == nil {
			return true
		}
		return candidates[i].VotedAt.Int64() > candidates[j].VotedAt.Int64()
	})

	return candidates[:limit], nil
}
