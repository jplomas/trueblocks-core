package parser

import (
	"fmt"
	"math/big"
	"strconv"

	"github.com/TrueBlocks/trueblocks-core/src/apps/chifra/pkg/base"
	"github.com/ethereum/go-ethereum/accounts/abi"
)

// ArgAddress is a type alias to capture bool values correctly
type ArgAddress base.Address

func (b *ArgAddress) Capture(values []string) error {
	debug("ArgAddress::Capture", values)
	*b = ArgAddress(base.HexToAddress(values[0]))
	return nil
}

// ArgBool is a type alias to capture bool values correctly
type ArgBool bool

func (b *ArgBool) Capture(values []string) error {
	debug("ArgBool::Capture", values)
	*b = values[0] == "true"
	return nil
}

// ArgHex represents anything that starts with 0x. If the value is a valid
// address, then it's capture into `base.Address` type, `string` otherwise.
type ArgHex struct {
	Address *base.Address
	String  *string
}

func (h *ArgHex) Capture(values []string) error {
	debug("ArgHex::Capture", values)
	hexLiteral := values[0]

	if valid, _ := base.IsValidHex("", hexLiteral, 20); !valid {
		h.String = &hexLiteral
		return nil
	}

	address := base.HexToAddress(hexLiteral)
	h.Address = &address
	return nil
}

// ArgNumber represents a number
type ArgNumber struct {
	Int  *int64
	Uint *uint64
	Big  *big.Int
}

func (n *ArgNumber) Capture(values []string) error {
	debug("ArgNumber::Capture", values)
	literal := values[0]

	// Atoi parses into `int` type, which is used by go-ethereum
	// to construct solidity int types.
	asInt, err := strconv.ParseInt(literal, 10, 64)
	if err == nil {
		n.Int = &asInt
		return nil
	}

	asUint, err := strconv.ParseUint(literal, 10, 64)
	if err == nil {
		n.Uint = &asUint
		return nil
	}

	// If we are here, the number is bigger than int64
	asBig := big.NewInt(0)
	_, ok := asBig.SetString(literal, 10)
	if ok {
		n.Big = asBig
		return nil
	}

	return fmt.Errorf("cannot parse %s as number", literal)
}

// Interface returns Number value as any
func (n *ArgNumber) Interface() any {
	if n.Int != nil {
		return *n.Int
	}
	if n.Uint != nil {
		return *n.Uint
	}
	return n.Big
}

func (n *ArgNumber) Convert(abiType *abi.Type) (any, error) {
	if abiType.Size > 64 {
		return n.Big, nil
	}

	if abiType.T == abi.UintTy {
		switch abiType.Size {
		case 8:
			return uint8(*n.Uint), nil
		case 16:
			return uint16(*n.Uint), nil
		case 32:
			return uint32(*n.Uint), nil
		case 64:
			return uint64(*n.Uint), nil
		}
	} else if abiType.T == abi.IntTy {
		switch abiType.Size {
		case 8:
			return int8(*n.Int), nil
		case 16:
			return int16(*n.Int), nil
		case 32:
			return int32(*n.Int), nil
		case 64:
			return int64(*n.Int), nil
		}
	}

	return nil, fmt.Errorf("cannot convert %v to number", n.Interface())
}

func debug(name string, values []string) {
	// fmt.Printf("%s%s: %v%s\n", colors.Green, name, values, colors.Off)
}
