// brigolagecms/cmd/ddltogo/buffer.go
//
// Copyright (c) 2020 Michael Henderson
//
// Permission is hereby granted, free of charge, to any person obtaining a copy
// of this software and associated documentation files (the "Software"), to deal
// in the Software without restriction, including without limitation the rights
// to use, copy, modify, merge, publish, distribute, sublicense, and/or sell
// copies of the Software, and to permit persons to whom the Software is
// furnished to do so, subject to the following conditions:
//
// The above copyright notice and this permission notice shall be included in all
// copies or substantial portions of the Software.
//
//  THE SOFTWARE IS PROVIDED "AS IS", WITHOUT WARRANTY OF ANY KIND, EXPRESS OR
//  IMPLIED, INCLUDING BUT NOT LIMITED TO THE WARRANTIES OF MERCHANTABILITY,
//  FITNESS FOR A PARTICULAR PURPOSE AND NONINFRINGEMENT. IN NO EVENT SHALL THE
//  AUTHORS OR COPYRIGHT HOLDERS BE LIABLE FOR ANY CLAIM, DAMAGES OR OTHER
//  LIABILITY, WHETHER IN AN ACTION OF CONTRACT, TORT OR OTHERWISE, ARISING FROM,
//  OUT OF OR IN CONNECTION WITH THE SOFTWARE OR THE USE OR OTHER DEALINGS IN THE
//  SOFTWARE.

package main

import (
	"bytes"
	"io/ioutil"
	"log"
	"unicode"
	"unicode/utf8"
)

// Buffer is used for reading sql files.
type Buffer struct {
	Line, Col int
	offset    int
	buf       []byte
}

// NewBuffer reads a file into a buffer.
func NewBuffer(filename string) *Buffer {
	buf, err := ioutil.ReadFile(filename)
	if err != nil {
		log.Fatal(err)
	}
	// normalize cr+lf to just lf, then normalize cr to lf
	return &Buffer{
		Line: 1,
		Col:  1,
		buf:  bytes.Replace(bytes.Replace(buf, []byte{13, 10}, []byte{10}, -1), []byte{13}, []byte{10}, -1),
	}
}

// Next returns the next rune
func (b *Buffer) Next() rune {
	if len(b.buf) <= b.offset {
		return utf8.RuneError
	}

	r, w := utf8.DecodeRune(b.buf[b.offset:])
	if r == utf8.RuneError {
		log.Fatalf("%d: %d: invalid utf-8\n", b.Line, b.Col)
	}
	b.offset += w

	if r == '\n' {
		b.Line, b.Col = b.Line+1, 0
	}
	b.Col = b.Col + 1

	return r
}

// Peek returns the next rune without advancing the input.
func (b *Buffer) Peek() rune {
	if len(b.buf) <= b.offset {
		return utf8.RuneError
	}
	r, _ := utf8.DecodeRune(b.buf[b.offset:])
	return r
}

// GetCComment returns a comment token
func (b *Buffer) GetCComment() *Token {
	start, t := b.offset, Token{
		Kind: TCComment,
		Line: b.Line,
		Col:  b.Col,
	}

	isComment := b.Next() == '/' && b.Next() == '*'
	if !isComment {
		b.Line, b.Col, b.offset = t.Line, t.Col, start
		return nil
	}

	for b.Peek() != utf8.RuneError {
		if b.Peek() == '*' {
			b.Next()
			if b.Peek() == '/' {
				b.Next()
				break
			}
		}
		b.Next()
	}

	t.Val = append(t.Val, b.buf[start:b.offset]...)

	return &t
}

// GetComment returns a comment token
func (b *Buffer) GetComment() *Token {
	start, t := b.offset, Token{
		Kind: TComment,
		Line: b.Line,
		Col:  b.Col,
	}

	isComment := b.Next() == '-' && b.Next() == '-'
	if !isComment {
		b.Line, b.Col, b.offset = t.Line, t.Col, start
		return nil
	}

	for r := b.Peek(); r != '\n' && r != utf8.RuneError; r = b.Peek() {
		b.Next()
	}

	t.Val = append(t.Val, b.buf[start+2:b.offset]...)

	return &t
}

// GetSpaces returns a spaces token
func (b *Buffer) GetSpaces() *Token {
	start, t := b.offset, Token{
		Kind: TSpace,
		Line: b.Line,
		Col:  b.Col,
	}

	if !unicode.IsSpace(b.Next()) {
		b.Line, b.Col, b.offset = t.Line, t.Col, start
		return nil
	}

	for r := b.Peek(); unicode.IsSpace(r); r = b.Peek() {
		b.Next()
	}

	t.Val = append(t.Val, b.buf[start:b.offset]...)

	return &t
}

// GetString returns a string
func (b *Buffer) GetString() *Token {
	start, t := b.offset, Token{
		Kind: TString,
		Line: b.Line,
		Col:  b.Col,
	}

	quote := b.Next()
	if quote != '\'' {
		b.Line, b.Col, b.offset = t.Line, t.Col, start
		return nil
	}

	r := b.Next()
	for !(r == quote || r == utf8.RuneError) {
		r = b.Next()
	}

	if r != quote {
		log.Fatalf("%d: %d: unterminated string\n", t.Line, t.Col)
	}

	t.Val = append(t.Val, b.buf[start:b.offset]...)

	return &t
}

// GetWord returns a word token
func (b *Buffer) GetWord() *Token {
	start, t := b.offset, Token{
		Kind: TWord,
		Line: b.Line,
		Col:  b.Col,
	}

	r := b.Next()
	if r == utf8.RuneError || unicode.IsSpace(r) {
		b.Line, b.Col, b.offset = t.Line, t.Col, start
		return nil
	}

	switch r {
	case ',':
		t.Kind = TComma
	case ';':
		t.Kind = TSemiColon
	case '(':
		t.Kind = TParenLeft
	case ')':
		t.Kind = TParenRight
	default:
		if unicode.IsNumber(r) {
			t.Kind = TNumber
			for r := b.Peek(); unicode.IsNumber(r); r = b.Peek() {
				b.Next()
			}
		} else if unicode.IsLetter(r) {
			t.Kind = TWord
			for r := b.Peek(); r == '_' || unicode.IsLetter(r) || unicode.IsNumber(r); r = b.Peek() {
				b.Next()
			}
		}
	}

	if t.Kind == TWord {
		t.Val = bytes.ToLower(b.buf[start:b.offset])
		switch string(t.Val) {
		case "aggregate":
			t.Kind = TAggregate
		case "boolean":
			t.Kind = TTypeBoolean
		case ",":
			t.Kind = TComma
		case "char":
			t.Kind = TTypeChar
		case "check":
			t.Kind = TCheck
		case "constraint":
			t.Kind = TConstraint
		case "create":
			t.Kind = TCreate
		case "current_timestamp":
			t.Kind = TCurrentTimestamp
		case "default":
			t.Kind = TDefault
		case "false":
			t.Kind = TFalse
		case "function":
			t.Kind = TFunction
		case "index":
			t.Kind = TIndex
		case "int2":
			t.Kind = TTypeInteger
		case "integer":
			t.Kind = TTypeInteger
		case "key":
			t.Kind = TKey
		case "nextval":
			t.Kind = TNextVal
		case "not":
			t.Kind = TNot
		case "null":
			t.Kind = TNull
		case "primary":
			t.Kind = TPrimary
		case "sequence":
			t.Kind = TSequence
		case "smallint":
			t.Kind = TTypeInteger
		case "start":
			t.Kind = TStart
		case "table":
			t.Kind = TTable
		case "text":
			t.Kind = TTypeText
		case "timestamp":
			t.Kind = TTypeTimestamp
		case "true":
			t.Kind = TTrue
		case "unique":
			t.Kind = TUnique
		case "varchar":
			t.Kind = TTypeVarchar
		}
	} else {
		t.Val = append(t.Val, b.buf[start:b.offset]...)
	}

	return &t
}
