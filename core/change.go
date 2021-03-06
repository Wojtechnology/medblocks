package core

import "github.com/wojtechnology/glacier/meddb"

// -------------------
// Transaction Changes
// -------------------

type TransactionChange struct {
	NewTransaction *Transaction
	OldTransaction *Transaction
}

// Wrapper around a meddb transaction that maps meddb transactions to core transactions.
type TransactionChangeCursor struct {
	changefeed meddb.TransactionChangefeed
}

func (cursor *TransactionChangeCursor) Next(change *TransactionChange) bool {
	var res meddb.TransactionChangefeedRes

	changed := cursor.changefeed.Next(&res)
	if changed {
		if res.NewVal != nil {
			change.NewTransaction = fromDBTransaction(res.NewVal)
		} else {
			change.NewTransaction = nil
		}
		if res.OldVal != nil {
			change.OldTransaction = fromDBTransaction(res.OldVal)
		} else {
			change.OldTransaction = nil
		}
	}

	return changed
}

type BlockChange struct {
	NewBlock *Block
	OldBlock *Block
}

// Wrapper around a meddb block that maps meddb blocks to core blocks.
type BlockChangeCursor struct {
	changefeed meddb.BlockChangefeed
}

func (cursor *BlockChangeCursor) Next(change *BlockChange) bool {
	var res meddb.BlockChangefeedRes

	changed := cursor.changefeed.Next(&res)
	if changed {
		if res.NewVal != nil {
			change.NewBlock = fromDBBlock(res.NewVal)
		} else {
			change.NewBlock = nil
		}
		if res.OldVal != nil {
			change.OldBlock = fromDBBlock(res.OldVal)
		} else {
			change.OldBlock = nil
		}
	}

	return changed
}

type VoteChange struct {
	NewVote *Vote
	OldVote *Vote
}

// Wrapper around a meddb vote that maps votes to core votes
type VoteChangeCursor struct {
	changefeed meddb.VoteChangefeed
}

func (cursor *VoteChangeCursor) Next(change *VoteChange) bool {
	var res meddb.VoteChangefeedRes

	changed := cursor.changefeed.Next(&res)
	if changed {
		if res.NewVal != nil {
			change.NewVote = fromDBVote(res.NewVal)
		} else {
			change.NewVote = nil
		}
		if res.OldVal != nil {
			change.OldVote = fromDBVote(res.OldVal)
		} else {
			change.OldVote = nil
		}
	}

	return changed
}
