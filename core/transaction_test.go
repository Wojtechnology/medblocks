package core

import (
	"math/big"
	"testing"

	"github.com/wojtechnology/glacier/meddb"
	"github.com/wojtechnology/glacier/test"
)

func TestToDBTransaction(t *testing.T) {
	tx := &Transaction{
		CellAddress: &CellAddress{
			TableName: []byte{42},
			RowId:     []byte{32},
			ColId:     []byte{43},
			VerId:     big.NewInt(4),
		},
		Data: []byte{69},
	}
	var (
		assignedTo   = []byte{12}
		lastAssigned = big.NewInt(420)
		hash         = rlpHash(tx)
	)

	expected := &meddb.Transaction{
		Hash: hash.Bytes(),
		CellAddress: &meddb.CellAddress{
			TableName: []byte{42},
			RowId:     []byte{32},
			ColId:     []byte{43},
			VerId:     big.NewInt(4),
		},
		AssignedTo:   assignedTo,
		LastAssigned: lastAssigned,
		Data:         []byte{69},
	}
	actual := tx.toDBTransaction(assignedTo, lastAssigned)

	test.AssertEqual(t, expected, actual)
}

func TestToDBTransactionEmpty(t *testing.T) {
	tx := &Transaction{
		CellAddress: &CellAddress{},
	}
	var (
		assignedTo   []byte   = nil
		lastAssigned *big.Int = nil
		hash                  = rlpHash(tx)
	)

	expected := &meddb.Transaction{
		Hash:         hash.Bytes(),
		CellAddress:  &meddb.CellAddress{},
		AssignedTo:   assignedTo,
		LastAssigned: lastAssigned,
	}
	actual := tx.toDBTransaction(assignedTo, lastAssigned)

	test.AssertEqual(t, expected, actual)
}
