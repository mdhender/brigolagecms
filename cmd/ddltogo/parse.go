// brigolagecms/cmd/ddltogo/parse.go
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
	"fmt"
	"log"
)

func parse(tokens []*Token) []*Node {
	var nodes []*Node
	for i := 0; i < len(tokens); {
		if agg, l := readAggregate(tokens[i:]); agg != nil {
			nodes = append(nodes, &Node{Kind: NAggregate, Aggregate: agg})
			i += l
			continue
		}

		if c, l := readComments(tokens[i:]); c != nil {
			i = i + l
			nodes = append(nodes, &Node{Kind: NComment, Comment: c})
			continue
		}

		if fnc, l := readFunction(tokens[i:]); fnc != nil {
			nodes = append(nodes, &Node{Kind: NFunction, Function: fnc})
			i += l
			continue
		}

		if idx, l := readIndex(tokens[i:]); idx != nil {
			nodes = append(nodes, &Node{Kind: NIndex, Index: idx})
			i += l
			continue
		}

		if seq, l := readSequence(tokens[i:]); seq != nil {
			nodes = append(nodes, &Node{Kind: NSequence, Sequence: seq})
			i += l
			continue
		}

		if tbl, l := readTable(tokens[i:]); tbl != nil {
			nodes = append(nodes, &Node{Kind: NTable, Table: tbl})
			i += l
			continue
		}

		if tokens[i].Kind == TEndOfInput {
			break
		}

		for k := i - 3; 0 <= k && k < i && k < len(tokens); k++ {
			tmp := tokens[k]
			fmt.Printf("%4d:%3d: prior      token %q\n", tmp.Line, tmp.Col, string(tmp.Val))
		}
		for k := i; k < i+4 && k < len(tokens); k++ {
			tmp := tokens[k]
			fmt.Printf("%4d:%3d: unexpected token %q\n", tmp.Line, tmp.Col, string(tmp.Val))
		}
		log.Fatalln("invalid input")
	}

	return nodes
}

func acceptAggregate(tokens []*Token) bool {
	return len(tokens) != 0 && tokens[0].Kind == TAggregate
}

func acceptCheck(tokens []*Token) bool {
	return len(tokens) != 0 && tokens[0].Kind == TCheck
}

func acceptComma(tokens []*Token) bool {
	return len(tokens) != 0 && tokens[0].Kind == TComma
}

func acceptConstraint(tokens []*Token) bool {
	return len(tokens) != 0 && tokens[0].Kind == TConstraint
}

func acceptCreate(tokens []*Token) bool {
	return len(tokens) != 0 && tokens[0].Kind == TCreate
}

func acceptCreateSequence(tokens []*Token) (*Node, int) {
	var name, start *Token
	i := 0
	if !acceptCreate(tokens[i:]) {
		return nil, 0
	}
	i++
	if !acceptSequence(tokens[i:]) {
		return nil, 0
	}
	i++
	if !acceptName(tokens[i:]) {
		log.Fatalf("%d:%d: expected sequence name\n", tokens[0].Line, tokens[0].Col)
	}
	name = tokens[i]
	i++
	if acceptStart(tokens[i:]) {
		i++
		if !acceptNumber(tokens[i:]) {
			log.Fatalf("%d:%d: expected sequence start value\n", tokens[0].Line, tokens[0].Col)
		}
		start = tokens[i]
		i++
	}
	if !acceptSemiColon(tokens[i:]) {
		log.Fatalf("%d:%d: unterminated create sequence statement\n", tokens[0].Line, tokens[0].Col)
	}
	i++

	n := &Node{Kind: NCreate, Sequence: &Sequence{Name: string(name.Val), Start: "0"}}
	if start != nil {
		n.Sequence.Start = string(start.Val)
	}

	return n, i
}

func acceptCreateTable(tokens []*Token) (*Node, int) {
	var name *Token
	i := 0
	if !acceptCreate(tokens[i:]) {
		return nil, 0
	}
	i++
	if !acceptTable(tokens[i:]) {
		return nil, 0
	}
	i++
	if !acceptName(tokens[i:]) {
		log.Fatalf("%d:%d: expected table name\n", tokens[0].Line, tokens[0].Col)
	}
	name = tokens[i]
	i++
	if !acceptParenLeft(tokens[i:]) {
		log.Fatalf("%d:%d: expected table '(', found %q\n", tokens[i].Line, tokens[i].Col, string(tokens[i].Val))
	}
	i++
	if !acceptParenRight(tokens[i:]) {
		log.Fatalf("%d:%d: expected table ')', found %q\n", tokens[i].Line, tokens[i].Col, string(tokens[i].Val))
	}
	i++
	if !acceptSemiColon(tokens[i:]) {
		log.Fatalf("%d:%d: unterminated create table statement\n", tokens[0].Line, tokens[0].Col)
	}
	i++

	n := &Node{Kind: NCreate, Table: &Table{Name: string(name.Val)}}

	return n, i
}

func acceptCurrentTimestamp(tokens []*Token) bool {
	return len(tokens) != 0 && tokens[0].Kind == TCurrentTimestamp
}

func acceptDefault(tokens []*Token) bool {
	return len(tokens) != 0 && tokens[0].Kind == TDefault
}

func acceptFalse(tokens []*Token) bool {
	return len(tokens) != 0 && tokens[0].Kind == TFalse
}

func acceptFunction(tokens []*Token) bool {
	return len(tokens) != 0 && tokens[0].Kind == TFunction
}

func acceptIndex(tokens []*Token) bool {
	return len(tokens) != 0 && tokens[0].Kind == TIndex
}

func acceptKey(tokens []*Token) bool {
	return len(tokens) != 0 && tokens[0].Kind == TKey
}

func acceptName(tokens []*Token) bool {
	return len(tokens) != 0 && tokens[0].Kind == TWord
}

func acceptNextVal(tokens []*Token) bool {
	return len(tokens) != 0 && tokens[0].Kind == TNextVal
}

func acceptNot(tokens []*Token) bool {
	return len(tokens) != 0 && tokens[0].Kind == TNot
}

func acceptNull(tokens []*Token) bool {
	return len(tokens) != 0 && tokens[0].Kind == TNull
}

func acceptNumber(tokens []*Token) bool {
	return len(tokens) != 0 && tokens[0].Kind == TNumber
}

func acceptParenLeft(tokens []*Token) bool {
	return len(tokens) != 0 && tokens[0].Kind == TParenLeft
}

func acceptParenRight(tokens []*Token) bool {
	return len(tokens) != 0 && tokens[0].Kind == TParenRight
}

func acceptPrimary(tokens []*Token) bool {
	return len(tokens) != 0 && tokens[0].Kind == TPrimary
}

func acceptPrimaryKey(tokens []*Token) bool {
	return acceptPrimary(tokens) && acceptKey(tokens[1:])
}

func acceptSemiColon(tokens []*Token) bool {
	return len(tokens) != 0 && tokens[0].Kind == TSemiColon
}

func acceptSequence(tokens []*Token) bool {
	return len(tokens) != 0 && tokens[0].Kind == TSequence
}

func acceptStart(tokens []*Token) bool {
	return len(tokens) != 0 && tokens[0].Kind == TStart
}

func acceptString(tokens []*Token) bool {
	return len(tokens) != 0 && tokens[0].Kind == TString
}

func acceptTable(tokens []*Token) bool {
	return len(tokens) != 0 && tokens[0].Kind == TTable
}

func acceptTrue(tokens []*Token) bool {
	return len(tokens) != 0 && tokens[0].Kind == TTrue
}

func acceptTypeBoolean(tokens []*Token) bool {
	return len(tokens) != 0 && tokens[0].Kind == TTypeBoolean
}

func acceptTypeChar(tokens []*Token) bool {
	return len(tokens) != 0 && tokens[0].Kind == TTypeChar
}

func acceptTypeInteger(tokens []*Token) bool {
	return len(tokens) != 0 && tokens[0].Kind == TTypeInteger
}

func acceptTypeText(tokens []*Token) bool {
	return len(tokens) != 0 && tokens[0].Kind == TTypeText
}

func acceptTypeTimestamp(tokens []*Token) bool {
	return len(tokens) != 0 && tokens[0].Kind == TTypeTimestamp
}

func acceptTypeVarchar(tokens []*Token) bool {
	return len(tokens) != 0 && tokens[0].Kind == TTypeVarchar
}

func acceptUnique(tokens []*Token) bool {
	return len(tokens) != 0 && tokens[0].Kind == TUnique
}

func acceptWord(tokens []*Token) bool {
	return len(tokens) != 0 && tokens[0].Kind == TWord
}
