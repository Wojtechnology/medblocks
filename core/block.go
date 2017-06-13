package core

import (
	"math/big"

	"github.com/wojtechnology/glacier/meddb"
)

type Block struct {
	Transactions []*Transaction // Contains hashes of all contained transactions
	CreatedAt    *big.Int       // Time at which block was created, will be used to determine order
	Creator      []byte
	Voters       [][]byte
}

// ---------
// Block API
// ---------

func (b *Block) Hash() Hash {
	// TODO: Think about this, maybe we want to hash a subset
	return rlpHash(b)
}

func (b *Block) toDBBlock() *meddb.Block {
	var createdAt *big.Int = nil
	if b.CreatedAt != nil {
		createdAt = big.NewInt(b.CreatedAt.Int64())
	}

	var txs []*meddb.Transaction = nil
	if b.Transactions != nil {
		txs = make([]*meddb.Transaction, len(b.Transactions))
		for i, tx := range b.Transactions {
			txs[i] = tx.toDBTransaction()
		}
	}

	// TODO(wojtek): Maybe make copies here
	return &meddb.Block{
		Hash:         b.Hash().Bytes(),
		Transactions: txs,
		CreatedAt:    createdAt,
		Creator:      b.Creator,
		Voters:       b.Voters,
	}
}

func fromDBBlock(b *meddb.Block) *Block {
	var createdAt *big.Int = nil
	if b.CreatedAt != nil {
		createdAt = big.NewInt(b.CreatedAt.Int64())
	}

	var txs []*Transaction = nil
	if b.Transactions != nil {
		txs = make([]*Transaction, len(b.Transactions))
		for i, tx := range b.Transactions {
			txs[i] = fromDBTransaction(tx)
		}
	}

	// TODO(wojtek): Maybe make copies here
	return &Block{
		Transactions: txs,
		CreatedAt:    createdAt,
		Creator:      b.Creator,
		Voters:       b.Voters,
	}
}
