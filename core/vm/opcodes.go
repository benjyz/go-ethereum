// Copyright 2014 The go-ethereum Authors
// This file is part of the go-ethereum library.
//
// The go-ethereum library is free software: you can redistribute it and/or modify
// it under the terms of the GNU Lesser General Public License as published by
// the Free Software Foundation, either version 3 of the License, or
// (at your option) any later version.
//
// The go-ethereum library is distributed in the hope that it will be useful,
// but WITHOUT ANY WARRANTY; without even the implied warranty of
// MERCHANTABILITY or FITNESS FOR A PARTICULAR PURPOSE. See the
// GNU Lesser General Public License for more details.
//
// You should have received a copy of the GNU Lesser General Public License
// along with the go-ethereum library. If not, see <http://www.gnu.org/licenses/>.

package vm

import (
	"fmt"
)

/*
--- OP code list ---

0s: Stop and Arithmetic Operations

0x00    STOP        Halts execution
0x01    ADD         Addition operation
0x02    MUL         Multiplication operation
0x03    SUB         Subtraction operation
0x04    DIV         Integer division operation
0x05    SDIV        Signed integer
0x06    MOD         Modulo
0x07    SMOD        Signed modulo
0x08    ADDMOD      Modulo
0x09    MULMOD      Modulo
0x0a    EXP         Exponential operation
0x0b    SIGNEXTEND  Extend length of two's complement signed integer
10s: Comparison & Bitwise Logic Operations

0x10    LT      Lesser-than comparison
0x11    GT      Greater-than comparison
0x12    SLT     Signed less-than comparison
0x13    SGT     Signed greater-than comparison
0x14    EQ      Equality  comparison
0x15    ISZERO  Simple not operator
0x16    AND     Bitwise AND operation
0x17    OR      Bitwise OR operation
0x18    XOR     Bitwise XOR operation
0x19    NOT     Bitwise NOT operation
0x1a    BYTE    Retrieve single byte from word
20s: SHA3

0x20    SHA3    Compute Keccak-256 hash

30s: Environmental Information

0x30    ADDRESS         Get address of currently executing account
0x31    BALANCE         Get balance of the given account
0x32    ORIGIN          Get execution origination address
0x33    CALLER          Get caller address. This is the address of the account that is directly responsible for this execution
0x34    CALLVALUE       Get deposited value by the instruction/transaction responsible for this execution
0x35    CALLDATALOAD    Get input data of current environment
0x36    CALLDATASIZE    Get size of input data in current environment
0x37    CALLDATACOPY    Copy input data in current environment to memory This pertains to the input data passed with the message call instruction or transaction
0x38    CODESIZE        Get size of code running in current environment
0x39    CODECOPY        Copy code running in current environment to memory
0x3a    GASPRICE        Get price of gas in current environment
0x3b    EXTCODESIZE     Get size of an account's code
0x3c    EXTCODECOPY     Copy an account's code to memory

40s: Block Information

0x40    BLOCKHASH   Get the hash of one of the 256 most recent complete blocks
0x41    COINBASE    Get the block's beneficiary address
0x42    TIMESTAMP   Get the block's timestamp
0x43    NUMBER      Get the block's number
0x44    DIFFICULTY  Get the block's difficulty
0x45    GASLIMIT    Get the block's gas limit

50s Stack, Memory, Storage and Flow Operations

0x50    POP         Remove item from stack
0x51    MLOAD       Load word from memory
0x52    MSTORE      Save word to memory
0x53    MSTORE8     Save byte to memory
0x54    SLOAD       Load word from storage
0x55    SSTORE      Save word to storage
0x56    JUMP        Alter the program counter
0x57    JUMPI       Conditionally alter the program counter
0x58    PC          Get the value of the program counter prior to the increment
0x59    MSIZE       Get the size of active memory in bytes
0x5a    GAS         Get the amount of available gas, including the corresponding reduction
0x5b    JUMPDEST    Mark a valid destination for jumps
60s & 70s: Push Operations

0x60    PUSH1   Place 1 byte item on stack
0x61    PUSH2   Place 2-byte item on stack
…
0x7f    PUSH32  Place 32-byte (full word) item on stack

80s: Duplication Operations

0x80    DUP1    Duplicate 1st stack item
0x81    DUP2    Duplicate 2nd stack item
…
0x8f    DUP16   Duplicate 16th stack item

90s: Exchange Operations

0x90    SWAP1   Exchange 1st and 2nd stack items
0x91    SWAP2   Exchange 1st and 3rd stack items
…   …
0x9f    SWAP16  Exchange 1st and 17th stack items

a0s: Logging Operations

0xa0    LOG0    Append log record with no topics
0xa1    LOG1    Append log record with one topic
…   …
0xa4    LOG4    Append log record with four topics

f0s: System operations

0xf0    CREATE          Create a new account with associated code
0xf1    CALL            Message-call into an account
0xf2    CALLCODE        Message-call into this account with alternative account's code
0xf3    RETURN          Halt execution returning output data
0xf4    DELEGATECALL    Message-call into this account with an alternative account's code, but persisting the current values for `sender` and `value`
Halt Execution, Mark for deletion

0xff    SELFDESTRUCT    Halt execution and register account for later deletion
*/

// OpCode is an EVM opcode
type OpCode byte

func (op OpCode) IsPush() bool {
	switch op {
	case PUSH1, PUSH2, PUSH3, PUSH4, PUSH5, PUSH6, PUSH7, PUSH8, PUSH9, PUSH10, PUSH11, PUSH12, PUSH13, PUSH14, PUSH15, PUSH16, PUSH17, PUSH18, PUSH19, PUSH20, PUSH21, PUSH22, PUSH23, PUSH24, PUSH25, PUSH26, PUSH27, PUSH28, PUSH29, PUSH30, PUSH31, PUSH32:
		return true
	}
	return false
}

func (op OpCode) IsStaticJump() bool {
	return op == JUMP
}

const (
	// 0x0 range - arithmetic ops
	STOP OpCode = iota
	ADD
	MUL
	SUB
	DIV
	SDIV
	MOD
	SMOD
	ADDMOD
	MULMOD
	EXP
	SIGNEXTEND
)

const (
	LT OpCode = iota + 0x10
	GT
	SLT
	SGT
	EQ
	ISZERO
	AND
	OR
	XOR
	NOT
	BYTE
	SHL
	SHR
	SAR

	SHA3 = 0x20
)

const (
	// 0x30 range - closure state
	ADDRESS OpCode = 0x30 + iota
	BALANCE
	ORIGIN
	CALLER
	CALLVALUE
	CALLDATALOAD
	CALLDATASIZE
	CALLDATACOPY
	CODESIZE
	CODECOPY
	GASPRICE
	EXTCODESIZE
	EXTCODECOPY
	RETURNDATASIZE
	RETURNDATACOPY
)

const (
	// 0x40 range - block operations
	BLOCKHASH OpCode = 0x40 + iota
	COINBASE
	TIMESTAMP
	NUMBER
	DIFFICULTY
	GASLIMIT
)

const (
	// 0x50 range - 'storage' and execution
	POP OpCode = 0x50 + iota
	MLOAD
	MSTORE
	MSTORE8
	SLOAD
	SSTORE
	JUMP
	JUMPI
	PC
	MSIZE
	GAS
	JUMPDEST
)

const (
	// 0x60 range
	PUSH1 OpCode = 0x60 + iota
	PUSH2
	PUSH3
	PUSH4
	PUSH5
	PUSH6
	PUSH7
	PUSH8
	PUSH9
	PUSH10
	PUSH11
	PUSH12
	PUSH13
	PUSH14
	PUSH15
	PUSH16
	PUSH17
	PUSH18
	PUSH19
	PUSH20
	PUSH21
	PUSH22
	PUSH23
	PUSH24
	PUSH25
	PUSH26
	PUSH27
	PUSH28
	PUSH29
	PUSH30
	PUSH31
	PUSH32
	DUP1
	DUP2
	DUP3
	DUP4
	DUP5
	DUP6
	DUP7
	DUP8
	DUP9
	DUP10
	DUP11
	DUP12
	DUP13
	DUP14
	DUP15
	DUP16
	SWAP1
	SWAP2
	SWAP3
	SWAP4
	SWAP5
	SWAP6
	SWAP7
	SWAP8
	SWAP9
	SWAP10
	SWAP11
	SWAP12
	SWAP13
	SWAP14
	SWAP15
	SWAP16
)

const (
	LOG0 OpCode = 0xa0 + iota
	LOG1
	LOG2
	LOG3
	LOG4
)

// unofficial opcodes used for parsing
const (
	PUSH OpCode = 0xb0 + iota
	DUP
	SWAP
)

const (
	// 0xf0 range - closures
	CREATE OpCode = 0xf0 + iota
	CALL
	CALLCODE
	RETURN
	DELEGATECALL
	STATICCALL = 0xfa

	REVERT       = 0xfd
	SELFDESTRUCT = 0xff
)

// Since the opcodes aren't all in order we can't use a regular slice
var opCodeToString = map[OpCode]string{
	// 0x0 range - arithmetic ops
	STOP:       "STOP",
	ADD:        "ADD",
	MUL:        "MUL",
	SUB:        "SUB",
	DIV:        "DIV",
	SDIV:       "SDIV",
	MOD:        "MOD",
	SMOD:       "SMOD",
	EXP:        "EXP",
	NOT:        "NOT",
	LT:         "LT",
	GT:         "GT",
	SLT:        "SLT",
	SGT:        "SGT",
	EQ:         "EQ",
	ISZERO:     "ISZERO",
	SIGNEXTEND: "SIGNEXTEND",

	// 0x10 range - bit ops
	AND:    "AND",
	OR:     "OR",
	XOR:    "XOR",
	BYTE:   "BYTE",
	SHL:    "SHL",
	SHR:    "SHR",
	SAR:    "SAR",
	ADDMOD: "ADDMOD",
	MULMOD: "MULMOD",

	// 0x20 range - crypto
	SHA3: "SHA3",

	// 0x30 range - closure state
	ADDRESS:        "ADDRESS",
	BALANCE:        "BALANCE",
	ORIGIN:         "ORIGIN",
	CALLER:         "CALLER",
	CALLVALUE:      "CALLVALUE",
	CALLDATALOAD:   "CALLDATALOAD",
	CALLDATASIZE:   "CALLDATASIZE",
	CALLDATACOPY:   "CALLDATACOPY",
	CODESIZE:       "CODESIZE",
	CODECOPY:       "CODECOPY",
	GASPRICE:       "GASPRICE",
	EXTCODESIZE:    "EXTCODESIZE",
	EXTCODECOPY:    "EXTCODECOPY",
	RETURNDATASIZE: "RETURNDATASIZE",
	RETURNDATACOPY: "RETURNDATACOPY",

	// 0x40 range - block operations
	BLOCKHASH:  "BLOCKHASH",
	COINBASE:   "COINBASE",
	TIMESTAMP:  "TIMESTAMP",
	NUMBER:     "NUMBER",
	DIFFICULTY: "DIFFICULTY",
	GASLIMIT:   "GASLIMIT",

	// 0x50 range - 'storage' and execution
	POP: "POP",
	//DUP:     "DUP",
	//SWAP:    "SWAP",
	MLOAD:    "MLOAD",
	MSTORE:   "MSTORE",
	MSTORE8:  "MSTORE8",
	SLOAD:    "SLOAD",
	SSTORE:   "SSTORE",
	JUMP:     "JUMP",
	JUMPI:    "JUMPI",
	PC:       "PC",
	MSIZE:    "MSIZE",
	GAS:      "GAS",
	JUMPDEST: "JUMPDEST",

	// 0x60 range - push
	PUSH1:  "PUSH1",
	PUSH2:  "PUSH2",
	PUSH3:  "PUSH3",
	PUSH4:  "PUSH4",
	PUSH5:  "PUSH5",
	PUSH6:  "PUSH6",
	PUSH7:  "PUSH7",
	PUSH8:  "PUSH8",
	PUSH9:  "PUSH9",
	PUSH10: "PUSH10",
	PUSH11: "PUSH11",
	PUSH12: "PUSH12",
	PUSH13: "PUSH13",
	PUSH14: "PUSH14",
	PUSH15: "PUSH15",
	PUSH16: "PUSH16",
	PUSH17: "PUSH17",
	PUSH18: "PUSH18",
	PUSH19: "PUSH19",
	PUSH20: "PUSH20",
	PUSH21: "PUSH21",
	PUSH22: "PUSH22",
	PUSH23: "PUSH23",
	PUSH24: "PUSH24",
	PUSH25: "PUSH25",
	PUSH26: "PUSH26",
	PUSH27: "PUSH27",
	PUSH28: "PUSH28",
	PUSH29: "PUSH29",
	PUSH30: "PUSH30",
	PUSH31: "PUSH31",
	PUSH32: "PUSH32",

	DUP1:  "DUP1",
	DUP2:  "DUP2",
	DUP3:  "DUP3",
	DUP4:  "DUP4",
	DUP5:  "DUP5",
	DUP6:  "DUP6",
	DUP7:  "DUP7",
	DUP8:  "DUP8",
	DUP9:  "DUP9",
	DUP10: "DUP10",
	DUP11: "DUP11",
	DUP12: "DUP12",
	DUP13: "DUP13",
	DUP14: "DUP14",
	DUP15: "DUP15",
	DUP16: "DUP16",

	SWAP1:  "SWAP1",
	SWAP2:  "SWAP2",
	SWAP3:  "SWAP3",
	SWAP4:  "SWAP4",
	SWAP5:  "SWAP5",
	SWAP6:  "SWAP6",
	SWAP7:  "SWAP7",
	SWAP8:  "SWAP8",
	SWAP9:  "SWAP9",
	SWAP10: "SWAP10",
	SWAP11: "SWAP11",
	SWAP12: "SWAP12",
	SWAP13: "SWAP13",
	SWAP14: "SWAP14",
	SWAP15: "SWAP15",
	SWAP16: "SWAP16",
	LOG0:   "LOG0",
	LOG1:   "LOG1",
	LOG2:   "LOG2",
	LOG3:   "LOG3",
	LOG4:   "LOG4",

	// 0xf0 range
	CREATE:       "CREATE",
	CALL:         "CALL",
	RETURN:       "RETURN",
	CALLCODE:     "CALLCODE",
	DELEGATECALL: "DELEGATECALL",
	STATICCALL:   "STATICCALL",
	REVERT:       "REVERT",
	SELFDESTRUCT: "SELFDESTRUCT",

	PUSH: "PUSH",
	DUP:  "DUP",
	SWAP: "SWAP",
}

func (o OpCode) String() string {
	str := opCodeToString[o]
	if len(str) == 0 {
		return fmt.Sprintf("Missing opcode 0x%x", int(o))
	}

	return str
}

var stringToOp = map[string]OpCode{
	"STOP":           STOP,
	"ADD":            ADD,
	"MUL":            MUL,
	"SUB":            SUB,
	"DIV":            DIV,
	"SDIV":           SDIV,
	"MOD":            MOD,
	"SMOD":           SMOD,
	"EXP":            EXP,
	"NOT":            NOT,
	"LT":             LT,
	"GT":             GT,
	"SLT":            SLT,
	"SGT":            SGT,
	"EQ":             EQ,
	"ISZERO":         ISZERO,
	"SIGNEXTEND":     SIGNEXTEND,
	"AND":            AND,
	"OR":             OR,
	"XOR":            XOR,
	"BYTE":           BYTE,
	"SHL":            SHL,
	"SHR":            SHR,
	"SAR":            SAR,
	"ADDMOD":         ADDMOD,
	"MULMOD":         MULMOD,
	"SHA3":           SHA3,
	"ADDRESS":        ADDRESS,
	"BALANCE":        BALANCE,
	"ORIGIN":         ORIGIN,
	"CALLER":         CALLER,
	"CALLVALUE":      CALLVALUE,
	"CALLDATALOAD":   CALLDATALOAD,
	"CALLDATASIZE":   CALLDATASIZE,
	"CALLDATACOPY":   CALLDATACOPY,
	"DELEGATECALL":   DELEGATECALL,
	"STATICCALL":     STATICCALL,
	"CODESIZE":       CODESIZE,
	"CODECOPY":       CODECOPY,
	"GASPRICE":       GASPRICE,
	"EXTCODESIZE":    EXTCODESIZE,
	"EXTCODECOPY":    EXTCODECOPY,
	"RETURNDATASIZE": RETURNDATASIZE,
	"RETURNDATACOPY": RETURNDATACOPY,
	"BLOCKHASH":      BLOCKHASH,
	"COINBASE":       COINBASE,
	"TIMESTAMP":      TIMESTAMP,
	"NUMBER":         NUMBER,
	"DIFFICULTY":     DIFFICULTY,
	"GASLIMIT":       GASLIMIT,
	"POP":            POP,
	"MLOAD":          MLOAD,
	"MSTORE":         MSTORE,
	"MSTORE8":        MSTORE8,
	"SLOAD":          SLOAD,
	"SSTORE":         SSTORE,
	"JUMP":           JUMP,
	"JUMPI":          JUMPI,
	"PC":             PC,
	"MSIZE":          MSIZE,
	"GAS":            GAS,
	"JUMPDEST":       JUMPDEST,
	"PUSH1":          PUSH1,
	"PUSH2":          PUSH2,
	"PUSH3":          PUSH3,
	"PUSH4":          PUSH4,
	"PUSH5":          PUSH5,
	"PUSH6":          PUSH6,
	"PUSH7":          PUSH7,
	"PUSH8":          PUSH8,
	"PUSH9":          PUSH9,
	"PUSH10":         PUSH10,
	"PUSH11":         PUSH11,
	"PUSH12":         PUSH12,
	"PUSH13":         PUSH13,
	"PUSH14":         PUSH14,
	"PUSH15":         PUSH15,
	"PUSH16":         PUSH16,
	"PUSH17":         PUSH17,
	"PUSH18":         PUSH18,
	"PUSH19":         PUSH19,
	"PUSH20":         PUSH20,
	"PUSH21":         PUSH21,
	"PUSH22":         PUSH22,
	"PUSH23":         PUSH23,
	"PUSH24":         PUSH24,
	"PUSH25":         PUSH25,
	"PUSH26":         PUSH26,
	"PUSH27":         PUSH27,
	"PUSH28":         PUSH28,
	"PUSH29":         PUSH29,
	"PUSH30":         PUSH30,
	"PUSH31":         PUSH31,
	"PUSH32":         PUSH32,
	"DUP1":           DUP1,
	"DUP2":           DUP2,
	"DUP3":           DUP3,
	"DUP4":           DUP4,
	"DUP5":           DUP5,
	"DUP6":           DUP6,
	"DUP7":           DUP7,
	"DUP8":           DUP8,
	"DUP9":           DUP9,
	"DUP10":          DUP10,
	"DUP11":          DUP11,
	"DUP12":          DUP12,
	"DUP13":          DUP13,
	"DUP14":          DUP14,
	"DUP15":          DUP15,
	"DUP16":          DUP16,
	"SWAP1":          SWAP1,
	"SWAP2":          SWAP2,
	"SWAP3":          SWAP3,
	"SWAP4":          SWAP4,
	"SWAP5":          SWAP5,
	"SWAP6":          SWAP6,
	"SWAP7":          SWAP7,
	"SWAP8":          SWAP8,
	"SWAP9":          SWAP9,
	"SWAP10":         SWAP10,
	"SWAP11":         SWAP11,
	"SWAP12":         SWAP12,
	"SWAP13":         SWAP13,
	"SWAP14":         SWAP14,
	"SWAP15":         SWAP15,
	"SWAP16":         SWAP16,
	"LOG0":           LOG0,
	"LOG1":           LOG1,
	"LOG2":           LOG2,
	"LOG3":           LOG3,
	"LOG4":           LOG4,
	"CREATE":         CREATE,
	"CALL":           CALL,
	"RETURN":         RETURN,
	"CALLCODE":       CALLCODE,
	"REVERT":         REVERT,
	"SELFDESTRUCT":   SELFDESTRUCT,
}

func StringToOp(str string) OpCode {
	return stringToOp[str]
}
