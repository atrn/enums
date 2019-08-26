// Copyright (C) 2015 A.Newman
//
package main

import (
	"errors"
	"fmt"
	"io"
)

var (
	errSyntax = errors.New("Syntax error")
)

type Parser struct {
	l       *Lexer
	Current Token
	Text    string
	Err     error
}

func NewParser(r io.Reader) *Parser {
	p := &Parser{l: NewLexer(r)}
	p.Advance()
	return p
}

func (p *Parser) SetError(err error) {
	if p.Err == nil {
		p.Err = err
	}
}

func (p *Parser) Broken() bool {
	return p.Err != nil
}

func (p *Parser) Advance() {
	if !p.Broken() {
		p.Current, p.Text, p.Err = p.l.Get()
	}
}

func (p *Parser) Done() bool {
	return p.Current == TokEOF || p.Broken()
}

func (p *Parser) Match(t Token) (string, bool) {
	if p.Broken() {
		return "", false
	}
	matched := t == p.Current
	text := ""
	if matched {
		text = p.Text
		p.Advance()
	}
	return text, matched
}

func (p *Parser) MatchOneOf(tokens []Token) (Token, bool) {
	if p.Broken() {
		return TokEOF, false
	}
	matched := false
	tok := TokEOF
	for _, token := range tokens {
		if p.Current == token {
			matched = true
			tok = token
			break
		}
	}
	if matched {
		p.Advance()
	}
	return tok, matched
}

func (p *Parser) MustMatch(t Token) (text string, matched bool) {
	if text, matched = p.Match(t); !matched {
		p.SetError(p.Expected(t, text))
	}
	return
}

func (p *Parser) Expected(t Token, text string) error {
	return fmt.Errorf("%d: expected %s, got %s %q", p.l.Line, t, p.Current, text)
}

func (p *Parser) Syntax(msg string) error {
	return fmt.Errorf("syntax error: %s", msg)
}
