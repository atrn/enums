// Copyright (C) 2015 A.Newman
//
package main

import (
	"bufio"
	"errors"
	"io"
	"unicode"
)

var (
	// ErrUnexpectedEOF is the error returned when an unexpected
	// end-of-file occurs during parsing.  This can happen if
	// a block comment is left unterminated.
	//
	ErrUnexpectedEOF = errors.New("unexpected end of file")

	// Keywords is a mapping from an identifier to a token for a
	// specific keywords. After reading an identifier it's looked
	// up in the Keywords map to determine if it is actually a
	// keyword rather than an user-defined identifier.
	//
	Keywords = map[string]Token{
		"package": TokPackage,
		"type":    TokType,
		"enum":    TokEnum,
		"byte":    TokByte,
		"int":     TokInt,
		"uint":    TokUint,
		"int8":    TokInt8,
		"uint8":   TokUint8,
		"int16":   TokInt16,
		"uint16":  TokUint16,
		"int32":   TokInt32,
		"uint32":  TokUint32,
		"int64":   TokInt64,
		"uint64":  TokUint64,
	}
)

// Lexer is the lexer's state. It holds the source of input, a
// bufio.Reader as we use it to "unread" input, and holds the
// current line number within the input (used by the parser
// when generating error messages related to the input).
//
type Lexer struct {
	r    *bufio.Reader
	Line int
}

// NewLexer returns a new Lexer that will read from the supplied
// io.Reader.
//
func NewLexer(r io.Reader) *Lexer {
	return &Lexer{
		r:    bufio.NewReader(r),
		Line: 1,
	}
}

// Get returns the next input token.  Comments and whitespace are
// skipped and the next token read.  Punctuation and other non-specific
// input is returned with a Token of "Other" with the text holding
// the single rune of input.
//
func (l *Lexer) Get() (Token, string, error) {
	type State int
	const (
		initialState State = iota
		maybeStartComment
		commentToEOL
		blockComment
		maybeEndComment
		identifier
		tag
	)
	var tok Token
	var text string
	var state State = initialState
	success := func(t Token, s string) (Token, string, error) {
		return t, s, nil
	}
	finalResult := func(err error) (Token, string, error) {
		failure := func(err error) (Token, string, error) {
			return TokEOF, "", err
		}
		if err != io.EOF {
			return failure(err)
		}
		switch state {
		case initialState, commentToEOL:
			return success(TokEOF, "")
		case maybeStartComment, blockComment, maybeEndComment:
			return failure(ErrUnexpectedEOF)
		case identifier:
			return success(TokIdent, text)
		}
		panic(err)
	}
	getchar := func() (rune, error) {
		ch, _, err := l.r.ReadRune()
		return ch, err
	}
	ungetchar := func() {
		l.r.UnreadRune()
	}
	isSpace := func(r rune) bool {
		return r == ' ' || r == '\t' || r == '\n' || r == '\r'
	}
	isFirstIdentifierChar := func(r rune) bool {
		return unicode.IsLetter(r) || r == '_'
	}
	isIdentifierChar := func(r rune) bool {
		return isFirstIdentifierChar(r) || unicode.IsDigit(r)
	}
	for {
		ch, err := getchar()
		if err != nil {
			return finalResult(err)
		}
		if ch == '\n' {
			l.Line++
		}
		switch state {
		case initialState:
			switch {
			case isSpace(ch):
				// eat space
			case ch == '/':
				state = maybeStartComment
			case ch == '=':
				return success(TokEq, "=")
			case ch == '{':
				return success(TokLBrace, "{")
			case ch == '}':
				return success(TokRBrace, "}")
			case ch == '`':
				state, text, tok = tag, "", TokTag
			case isFirstIdentifierChar(ch):
				state, text, tok = identifier, string(ch), TokIdent
			default:
				return success(TokOther, string(ch))
			}
		case tag:
			if ch == '`' {
				return success(TokTag, text)
			} else {
				text += string(ch)
			}
		case identifier:
			if isIdentifierChar(ch) {
				text += string(ch)
			} else {
				ungetchar()
				state = initialState
				if keywordTok, found := Keywords[text]; found {
					tok = keywordTok
				}
				return success(tok, text)
			}
		case maybeStartComment:
			if ch == '/' {
				state = commentToEOL
			} else if ch == '*' {
				state = blockComment
			} else {
				ungetchar()
				state = initialState
				return success(TokOther, "/")
			}
		case commentToEOL:
			if ch == '\n' {
				state = initialState
			}
		case blockComment:
			if ch == '*' {
				state = maybeEndComment
			}
		case maybeEndComment:
			if ch == '/' {
				state = initialState
			} else {
				state = blockComment
			}
		}
	}
	return TokEOF, "", nil
}
