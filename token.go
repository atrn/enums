// Copyright (C) 2015 A.Newman
//
package main

import "fmt"

// Token is the type of an input token read by the lexer.
//
type Token int

const (
	TokBOF Token = iota
	TokEOF
	TokByte
	TokEnum
	TokEq
	TokIdent
	TokInt
	TokInt16
	TokInt32
	TokInt64
	TokInt8
	TokLBrace
	TokOther
	TokPackage
	TokRBrace
	TokTag
	TokType
	TokUint
	TokUint16
	TokUint32
	TokUint64
	TokUint8
)

func (t Token) String() string {
	switch t {
	case TokBOF:
		return "BOF"
	case TokEOF:
		return "EOF"
	case TokByte:
		return "byte"
	case TokEnum:
		return "Enum"
	case TokEq:
		return "Eq"
	case TokIdent:
		return "Ident"
	case TokInt16:
		return "int16"
	case TokInt32:
		return "int32"
	case TokInt64:
		return "int64"
	case TokInt8:
		return "int8"
	case TokInt:
		return "int"
	case TokLBrace:
		return "LBrace"
	case TokOther:
		return "Other"
	case TokPackage:
		return "Package"
	case TokRBrace:
		return "RBrace"
	case TokTag:
		return "Tag"
	case TokType:
		return "Type"
	case TokUint16:
		return "uint16"
	case TokUint32:
		return "uint32"
	case TokUint64:
		return "uint64"
	case TokUint8:
		return "uint8"
	case TokUint:
		return "uint"
	}
	panic(fmt.Errorf("Bad Token: %x", int(t)))
}
