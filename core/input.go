package core

import (
	"errors"
	"fmt"
	"math/big"

	"github.com/wojtechnology/glacier/meddb"
)

type InputType int

const (
	INPUT_TYPE_ADMIN      InputType = iota // ADMIN      = 0
	INPUT_TYPE_WRITER                      // WRITER     = 1
	INPUT_TYPE_ROW_WRITER                  // ROW_WRITER = 2
)

type Input interface {
	OutputHash() Hash
	Type() InputType
	Data() []byte
	FromData([]byte) error
}

// ---------
// InputLink
// ---------

// Forms a link from an input to an output
type InputLink struct {
	LinksTo Hash // The hash of the output that this input links to.
}

func (link *InputLink) OutputHash() Hash {
	return link.LinksTo
}

// --------------------------------
// AdminInput implementation
//
// Allows a particular user to update a table
// --------------------------------

type AdminInput struct {
	InputLink
	Sig []byte
}

func (in *AdminInput) Type() InputType {
	return INPUT_TYPE_ADMIN
}

func (in *AdminInput) Data() []byte {
	return in.Sig
}

func (in *AdminInput) FromData(data []byte) error {
	in.Sig = data
	return nil
}

// --------------------------------
// WriterInput implementation
//
// Allows a particular user to write to a table
// --------------------------------

type WriterInput struct {
	InputLink
	Sig []byte
}

func (in *WriterInput) Type() InputType {
	return INPUT_TYPE_WRITER
}

func (in *WriterInput) Data() []byte {
	return in.Sig
}

func (in *WriterInput) FromData(data []byte) error {
	in.Sig = data
	return nil
}

// --------------------------------
// RowWriterInput implementation
//
// Allows a particular user to write to a row in a table
// --------------------------------

type RowWriterInput struct {
	InputLink
	Sig []byte
}

func (in *RowWriterInput) Type() InputType {
	return INPUT_TYPE_ROW_WRITER
}

func (in *RowWriterInput) Data() []byte {
	return in.Sig
}

func (in *RowWriterInput) FromData(data []byte) error {
	in.Sig = data
	return nil
}

// -------
// Helpers
// -------

type inputHashObject struct {
	Type       *big.Int
	OutputHash []byte
	Data       []byte
}

func hashInput(in Input) Hash {
	return rlpHash(&inputHashObject{
		Type:       intToBigInt(int(in.Type())),
		OutputHash: in.OutputHash().Bytes(),
		Data:       in.Data(),
	})
}

func toDBInput(in Input) *meddb.Input {
	return &meddb.Input{
		Type:       int(in.Type()),
		OutputHash: in.OutputHash().Bytes(),
		Data:       in.Data(),
	}
}

func NewInput(inputType int, outputHash, data []byte) (Input, error) {
	var coreInput Input

	switch InputType(inputType) {
	case INPUT_TYPE_ADMIN:
		coreInput = &AdminInput{InputLink: InputLink{BytesToHash(outputHash)}}
	case INPUT_TYPE_WRITER:
		coreInput = &WriterInput{InputLink: InputLink{BytesToHash(outputHash)}}
	case INPUT_TYPE_ROW_WRITER:
		coreInput = &RowWriterInput{InputLink: InputLink{BytesToHash(outputHash)}}
	default:
		return nil, errors.New(fmt.Sprint("Invalid input type: %d\n", inputType))
	}

	if err := coreInput.FromData(data); err != nil {
		return nil, err
	}

	return coreInput, nil
}
