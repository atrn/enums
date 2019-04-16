// Copyright (C) 2015 A.Newman
//
package main

import (
	"io"
	"time"
)

type Enums struct {
	Package  string
	Imports  []string
	Enums    []EnumDefinition
	Filename string
	Time     time.Time
}

type Enumerator struct {
	Enum string
	Tag  string
}

type EnumDefinition struct {
	TypeName    string
	BaseType    Token
	Enumerators []Enumerator
}

func ParseToEnd(r io.ReadCloser, e *Enums) error {
	err := Parse(r, e)
	r.Close()
	return err
}

var (
	intTypes = []Token{
		TokByte,
		TokInt,
		TokUint,
		TokInt8,
		TokUint8,
		TokInt16,
		TokUint16,
		TokInt32,
		TokUint32,
		TokInt64,
		TokUint64,
	}
)

func Parse(r io.Reader, e *Enums) error {
	inEnum := false
	var def EnumDefinition
	p := NewParser(r)
	p.MustMatch(TokPackage)
	var ok bool
	e.Package, ok = p.MustMatch(TokIdent)
	if !ok {
		return p.Err
	}
	for !p.Done() {
		if inEnum {
			if _, match := p.Match(TokRBrace); match {
				e.Enums = append(e.Enums, def)
				inEnum = false
			} else {
				tok, ok := p.MustMatch(TokIdent)
				if !ok {
					break
				}
				enum := Enumerator{tok, tok}
				if tag, hasTag := p.Match(TokTag); hasTag {
					enum.Tag = tag
				}
				def.Enumerators = append(def.Enumerators, enum)
			}
		} else {
			_, ok = p.MustMatch(TokType)
			if !ok {
				break
			}
			def.TypeName, ok = p.MustMatch(TokIdent)
			if !ok {
				break
			}
			_, ok = p.MustMatch(TokEnum)
			if !ok {
				break
			}
			if _, matched := p.Match(TokLBrace); matched {
				def.BaseType = TokInt
			} else if baseType, matched := p.MatchOneOf(intTypes); matched {
				def.BaseType = baseType
				_, ok = p.MustMatch(TokLBrace)
				if !ok {
					break
				}
			} else {
				return p.Syntax("expected '{' or <int-type>")
			}
			inEnum = true
			def.Enumerators = nil
		}
	}
	return p.Err
}
