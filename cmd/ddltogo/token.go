// brigolagecms/cmd/ddltogo/token.go
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

// Token is a literal from the buffer.
type Token struct {
	Kind      TokenKind
	Line, Col int
	Val       []byte
}

// TokenKind is the type of a literal.
type TokenKind int

// enums for tokens
const (
	TEndOfInput TokenKind = iota
	TAggregate
	TCComment
	TCheck
	TComma
	TComment
	TConstraint
	TCreate
	TCreateIndex
	TCreateTable
	TCurrentTimestamp
	TDefault
	TFalse
	TFunction
	TIndex
	TKey
	TNextVal
	TNot
	TNull
	TNumber
	TParenLeft
	TParenRight
	TPrimary
	TSemiColon
	TSequence
	TSpace
	TStart
	TString
	TTable
	TTrue
	TTypeBoolean
	TTypeChar
	TTypeInteger
	TTypeText
	TTypeTimestamp
	TTypeVarchar
	TUnique
	TWord
)

// Tokenize loads a file and returns a slice of tokens.
func Tokenize(filename string) []*Token {
	b := NewBuffer(filename)

	var tokens []*Token
	for {
		// read spaces
		if t := b.GetSpaces(); t != nil {
			// log.Printf("[spaces ] %3d %3d %s\n", t.Line, t.Col, string(t.Val))
			continue
		}

		// read comments
		if t := b.GetComment(); t != nil {
			// log.Printf("[comment] %3d %3d %s\n", t.Line, t.Col, string(t.Val))
			tokens = append(tokens, t)
			continue
		} else if t := b.GetCComment(); t != nil {
			// log.Printf("[comment] %3d %3d %s\n", t.Line, t.Col, string(t.Val))
			tokens = append(tokens, t)
			continue
		}

		// read strings
		if t := b.GetString(); t != nil {
			// log.Printf("[string ] %3d %3d %s\n", t.Line, t.Col, string(t.Val))
			tokens = append(tokens, t)
			continue
		}

		// read text
		if t := b.GetWord(); t != nil {
			// log.Printf("[word   ] %3d %3d %s\n", t.Line, t.Col, string(t.Val))
			tokens = append(tokens, t)
			continue
		}

		return append(tokens, &Token{Kind: TEndOfInput})
	}
}
