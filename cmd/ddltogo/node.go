// brigolagecms/cmd/ddltogo/node.go
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

// Tree is a syntrax tree.
type Tree struct {
	Statements  []*Node
	CurrentNode *Node
}

// Node is a node in the syntax tree
type Node struct {
	Kind       NodeKind
	Aggregate  *Aggregate
	Comment    *Comment
	Function   *Function
	Index      *Index
	Name       string
	Sequence   *Sequence
	StartValue string
	Table      *Table
}

// NodeKind is the type of a node.
type NodeKind int

// enums for nodes
const (
	NRoot NodeKind = iota
	NAggregate
	NColumnBoolean
	NColumnChar
	NColumnInteger
	NColumnText
	NColumnTimestamp
	NColumnVarchar
	NComment
	NCreate
	NFunction
	NIndex
	NName
	NSequence
	NStartValue
	NTable
)

// Aggregate represents an aggregator.
type Aggregate struct {
	Name string
}

// Column represents a column.
type Column struct {
	Name        string
	Kind        NodeKind
	Nullable    bool
	Width       string
	Constraints []*Constraint
	Default     *Default
}

// Columns is a node for columns.
type Columns struct {
	Val []*Column
}

// Comment is a node for comments.
type Comment struct {
	Val []*Token
}

// Constraint is a node for constraints
type Constraint struct {
	Name  string
	Check struct {
		Is      bool
		Actions []*Token
	}
	PrimaryKey struct {
		Is   bool
		Keys []string // names of the key fields (in order)
	}
}

// Default gives the column a default value
type Default struct {
	Kind string
	Val  string
}

// Function is a node for a function.
type Function struct {
	Name string
}

// Index is a node for an index.
type Index struct {
	Name   string
	Unique bool
}

// Sequence represents a sequence.
type Sequence struct {
	Name  string
	Incr  string
	Start string
}

// Table represents a table.
type Table struct {
	Name        string
	Columns     []*Column
	Constraints []*Constraint
}

func readAggregate(tokens []*Token) (*Aggregate, int) {
	agg, i := Aggregate{}, 0
	if !acceptCreate(tokens[i:]) {
		return nil, 0
	}
	i++
	if !acceptAggregate(tokens[i:]) {
		return nil, 0
	}
	i++
	if !acceptName(tokens[i:]) {
		log.Fatalf("%d:%d: expected aggregate name\n", tokens[0].Line, tokens[0].Col)
	}
	agg.Name = string(tokens[i].Val)
	i++
	for i < len(tokens) && !acceptSemiColon(tokens[i:]) {
		i++
	}
	if !acceptSemiColon(tokens[i:]) {
		log.Fatalf("%d:%d: unterminated create aggregate statement\n", tokens[0].Line, tokens[0].Col)
	}
	i++

	return &agg, i
}

func readBoolean(tokens []*Token) (*Column, int) {
	c := Column{Kind: NColumnBoolean}
	i := 0
	if !acceptTypeBoolean(tokens[i:]) {
		return nil, 0
	}
	i++
	return &c, i
}

func readChar(tokens []*Token) (*Column, int) {
	c := Column{Kind: NColumnChar}
	i := 0
	if !acceptTypeChar(tokens[i:]) {
		return nil, 0
	}
	i++
	if !acceptParenLeft(tokens[i:]) {
		log.Fatalf("%d:%d: expected char '(', found %q\n", tokens[i].Line, tokens[i].Col, string(tokens[i].Val))
	}
	i++
	if !acceptNumber(tokens[i:]) {
		log.Fatalf("%d:%d: expected char width, found %q\n", tokens[i].Line, tokens[i].Col, string(tokens[i].Val))
	}
	c.Width = string(tokens[i].Val)
	i++
	if !acceptParenRight(tokens[i:]) {
		log.Fatalf("%d:%d: unterminated char, found %q\n", tokens[i].Line, tokens[i].Col, string(tokens[i].Val))
	}
	i++
	return &c, i
}

func readCheck(tokens []*Token) ([]*Token, int) {
	i := 0
	if !acceptCheck(tokens[i:]) {
		return nil, 0
	}
	i++
	if !acceptParenLeft(tokens[i:]) {
		log.Fatalf("%d:%d: expected check '(', found %q\n", tokens[i].Line, tokens[i].Col, string(tokens[i].Val))
	}
	i++
	var actions []*Token
	parenLevel := 1
	for parenLevel != 0 {
		if tokens[i].Kind == TEndOfInput {
			log.Fatalf("%d:%d: unterminated check\n", tokens[0].Line, tokens[0].Col)
		} else if acceptParenLeft(tokens[i:]) {
			parenLevel++
		} else if acceptParenRight(tokens[i:]) {
			parenLevel--
		} else {
			actions = append(actions, tokens[i])
		}
		i++
	}
	return actions, i
}

func readColumn(tokens []*Token) (*Column, int) {
	i := 0
	if !acceptName(tokens[i:]) {
		return nil, 0
	}
	col := Column{
		Name:     string(tokens[i].Val),
		Nullable: true,
	}
	i++
	if c, l := readBoolean(tokens[i:]); c != nil {
		col.Kind = c.Kind
		i += l
	} else if c, l := readChar(tokens[i:]); c != nil {
		col.Kind = c.Kind
		col.Width = c.Width
		i += l
	} else if c, l := readInteger(tokens[i:]); c != nil {
		col.Kind = c.Kind
		i += l
	} else if c, l := readText(tokens[i:]); c != nil {
		col.Kind = c.Kind
		i += l
	} else if c, l := readTimestamp(tokens[i:]); c != nil {
		col.Kind = c.Kind
		i += l
	} else if c, l := readVarchar(tokens[i:]); c != nil {
		col.Kind = c.Kind
		col.Width = c.Width
		i += l
	} else {
		log.Fatalf("%d:%d: expected column type, found %q\n", tokens[i].Line, tokens[i].Col, string(tokens[i].Val))
	}
	if acceptNot(tokens[i:]) {
		i++
		if !acceptNull(tokens[i:]) {
			log.Fatalf("%d:%d: expected not null, found %q\n", tokens[i].Line, tokens[i].Col, string(tokens[i].Val))
		}
		col.Nullable = false
		i++
	}
	if dflt, l := readDefault(tokens[i:]); dflt != nil {
		col.Default = dflt
		i += l
	}
	if c, l := readConstraint(tokens[i:]); c != nil {
		col.Constraints = append(col.Constraints, c)
		i += l
	}
	return &col, i
}

func readComments(tokens []*Token) (*Comment, int) {
	c := Comment{}
	for i := 0; i < len(tokens) && (tokens[i].Kind == TComment || tokens[i].Kind == TCComment); i = i + 1 {
		c.Val = append(c.Val, tokens[i])
	}
	if len(c.Val) == 0 {
		return nil, 0
	}
	return &c, len(c.Val)
}

// CONSTRAINT ck_media_id               CHECK       (alias_id != id)
// CONSTRAINT ck_template___tplate_type CHECK       (tplate_type IN (1, 2, 3))
// CONSTRAINT pk_element_type__id       PRIMARY KEY (id)
func readConstraint(tokens []*Token) (*Constraint, int) {
	c := Constraint{}
	i := 0
	if !acceptConstraint(tokens[i:]) {
		return nil, 0
	}
	i++
	if !acceptName(tokens[i:]) {
		log.Fatalf("%d:%d: expected constraint name, found %q\n", tokens[i].Line, tokens[i].Col, string(tokens[i].Val))
	}
	c.Name = string(tokens[i].Val)
	i++
	if chk, l := readCheck(tokens[i:]); chk != nil {
		c.Check.Is = true
		c.Check.Actions = chk
		i += l
		return &c, i
	} else if acceptPrimaryKey(tokens[i:]) {
		c.PrimaryKey.Is = true
		i += 2
		list, l := readListOfNames(tokens[i:])
		c.PrimaryKey.Keys = list
		i += l
		return &c, i
	}
	for k := 0; k < i && k < len(tokens); k++ {
		tmp := tokens[k]
		fmt.Printf("%4d:%3d:      prior token %q\n", tmp.Line, tmp.Col, string(tmp.Val))
	}
	for k := i; k < i+4 && k < len(tokens); k++ {
		tmp := tokens[k]
		fmt.Printf("%4d:%3d: unexpected token %q\n", tmp.Line, tmp.Col, string(tmp.Val))
	}
	log.Fatalf("%d:%d: expected constraint, found %q\n", tokens[i].Line, tokens[i].Col, string(tokens[i].Val))
	// NOT REACHED
	return nil, 0
}

// DEFAULT 0
// DEFAULT CURRENT_TIMESTAMP
// DEFAULT FALSE
// DEFAULT 'text'
func readDefault(tokens []*Token) (*Default, int) {
	i := 0
	if !acceptDefault(tokens[i:]) {
		return nil, 0
	}
	i++
	if acceptCurrentTimestamp(tokens[i:]) {
		return &Default{Kind: "timestamp", Val: string(tokens[i].Val)}, i + 1
	} else if c, l := readNextVal(tokens[i:]); c != nil {
		return &Default{Kind: "generator", Val: c.Name}, i + l
	} else if acceptFalse(tokens[i:]) {
		return &Default{Kind: "boolean", Val: string(tokens[i].Val)}, i + 1
	} else if acceptNumber(tokens[i:]) {
		return &Default{Kind: "number", Val: string(tokens[i].Val)}, i + 1
	} else if acceptString(tokens[i:]) {
		return &Default{Kind: "text", Val: string(tokens[i].Val)}, i + 1
	} else if acceptTrue(tokens[i:]) {
		return &Default{Kind: "boolean", Val: string(tokens[i].Val)}, i + 1
	}
	log.Fatalf("%d:%d: expected default value, found %q\n", tokens[i].Line, tokens[i].Col, string(tokens[i].Val))
	// NOT REACHED
	return nil, 0
}

func readFunction(tokens []*Token) (*Function, int) {
	fnc, i := Function{}, 0
	if !acceptCreate(tokens[i:]) {
		return nil, 0
	}
	i++
	if !acceptFunction(tokens[i:]) {
		return nil, 0
	}
	i++
	if !acceptName(tokens[i:]) {
		log.Fatalf("%d:%d: expected function name\n", tokens[0].Line, tokens[0].Col)
	}
	fnc.Name = string(tokens[i].Val)
	i++
	for i < len(tokens) && !acceptSemiColon(tokens[i:]) {
		i++
	}
	if !acceptSemiColon(tokens[i:]) {
		log.Fatalf("%d:%d: unterminated create function statement\n", tokens[0].Line, tokens[0].Col)
	}
	i++

	return &fnc, i
}

func readIndex(tokens []*Token) (*Index, int) {
	idx, i := Index{}, 0
	if !acceptCreate(tokens[i:]) {
		return nil, 0
	}
	i++
	if acceptUnique(tokens[i:]) {
		idx.Unique = true
		i++
	}
	if !acceptIndex(tokens[i:]) {
		return nil, 0
	}
	i++
	if !acceptName(tokens[i:]) {
		log.Fatalf("%d:%d: expected index name\n", tokens[0].Line, tokens[0].Col)
	}
	idx.Name = string(tokens[i].Val)
	i++
	for i < len(tokens) && !acceptSemiColon(tokens[i:]) {
		i++
	}
	if !acceptSemiColon(tokens[i:]) {
		log.Fatalf("%d:%d: unterminated create index statement\n", tokens[0].Line, tokens[0].Col)
	}
	i++

	return &idx, i
}

func readInteger(tokens []*Token) (*Column, int) {
	c := Column{Kind: NColumnInteger}
	i := 0
	if !acceptTypeInteger(tokens[i:]) {
		return nil, 0
	}
	i++
	return &c, i
}

func readListOfNames(tokens []*Token) ([]string, int) {
	i := 0
	if !acceptParenLeft(tokens[i:]) {
		return nil, 0
	}
	i++
	var list []string
	for acceptName(tokens[i:]) {
		list = append(list, string(tokens[i].Val))
		i++
		if !acceptComma(tokens[i:]) {
			break
		}
		i++
	}
	if !acceptParenRight(tokens[i:]) {
		for k := i; k < i+4 && k < len(tokens); k++ {
			tmp := tokens[k]
			fmt.Printf("%4d:%3d: unexpected token %q\n", tmp.Line, tmp.Col, string(tmp.Val))
		}
		log.Fatalf("%d:%d: expected list ')', found %q\n", tokens[i].Line, tokens[i].Col, string(tokens[i].Val))
	}
	i++
	return list, i
}

func readNextVal(tokens []*Token) (*Sequence, int) {
	seq := Sequence{Start: "0", Incr: "1"}
	i := 0
	if !acceptNextVal(tokens[i:]) {
		return nil, 0
	}
	i++
	if !acceptParenLeft(tokens[i:]) {
		log.Fatalf("%d:%d: expected nextval '(', found %q\n", tokens[i].Line, tokens[i].Col, string(tokens[i].Val))
	}
	i++
	if !acceptString(tokens[i:]) {
		log.Fatalf("%d:%d: expected nextval generator, found %q\n", tokens[i].Line, tokens[i].Col, string(tokens[i].Val))
	}
	seq.Name = string(tokens[i].Val)
	i++
	if !acceptParenRight(tokens[i:]) {
		for k := i; k < i+4 && k < len(tokens); k++ {
			tmp := tokens[k]
			fmt.Printf("%4d:%3d: unexpected token %q\n", tmp.Line, tmp.Col, string(tmp.Val))
		}
		log.Fatalf("%d:%d: expected nextval ')', found %q\n", tokens[i].Line, tokens[i].Col, string(tokens[i].Val))
	}
	i++
	return &seq, i
}

func readSequence(tokens []*Token) (*Sequence, int) {
	s, i := Sequence{Start: "0", Incr: "1"}, 0
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
	s.Name = string(tokens[i].Val)
	i++
	if v, l := readStartValue(tokens[i:]); l != 0 {
		s.Start = v
		i = i + l
	}
	if !acceptSemiColon(tokens[i:]) {
		log.Fatalf("%d:%d: unterminated create sequence statement\n", tokens[0].Line, tokens[0].Col)
	}
	i++

	return &s, i
}

func readStartValue(tokens []*Token) (string, int) {
	if !acceptStart(tokens) || !acceptNumber(tokens[1:]) {
		return "", 0
	}
	return string(tokens[1].Val), 2
}

func readTable(tokens []*Token) (*Table, int) {
	tbl, i := Table{}, 0
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
	tbl.Name = string(tokens[i].Val)
	i++

	if !acceptParenLeft(tokens[i:]) {
		log.Fatalf("%d:%d: expected columns '(', found %q\n", tokens[i].Line, tokens[i].Col, string(tokens[i].Val))
	}
	i++

	for {
		if col, l := readColumn(tokens[i:]); col != nil {
			tbl.Columns = append(tbl.Columns, col)
			i += l
		} else if con, l := readConstraint(tokens[i:]); con != nil {
			tbl.Constraints = append(tbl.Constraints, con)
			i += l
		} else {
			log.Fatalf("%d:%d: expected column or constraint, found %q\n", tokens[i].Line, tokens[i].Col, string(tokens[i].Val))
		}
		if !acceptComma(tokens[i:]) {
			break
		}
		i++
	}
	if len(tbl.Constraints) == 0 {
		log.Fatalf("%d:%d: expected list of columns, found %q\n", tokens[i].Line, tokens[i].Col, string(tokens[i].Val))
	}

	if !acceptParenRight(tokens[i:]) {
		for k := i; k < i+4 && k < len(tokens); k++ {
			tmp := tokens[k]
			fmt.Printf("%4d:%3d: unexpected token %q\n", tmp.Line, tmp.Col, string(tmp.Val))
		}
		log.Fatalf("%d:%d: expected columns ')', found %q\n", tokens[i].Line, tokens[i].Col, string(tokens[i].Val))
	}
	i++

	if !acceptSemiColon(tokens[i:]) {
		log.Fatalf("%d:%d: unterminated create table statement\n", tokens[0].Line, tokens[0].Col)
	}
	i++

	return &tbl, i
}

func readText(tokens []*Token) (*Column, int) {
	c := Column{Kind: NColumnText}
	i := 0
	if !acceptTypeText(tokens[i:]) {
		return nil, 0
	}
	i++
	return &c, i
}

func readTimestamp(tokens []*Token) (*Column, int) {
	c := Column{Kind: NColumnTimestamp}
	i := 0
	if !acceptTypeTimestamp(tokens[i:]) {
		return nil, 0
	}
	i++
	return &c, i
}

func readVarchar(tokens []*Token) (*Column, int) {
	c := Column{Kind: NColumnVarchar}
	i := 0
	if !acceptTypeVarchar(tokens[i:]) {
		return nil, 0
	}
	i++
	if !acceptParenLeft(tokens[i:]) {
		log.Fatalf("%d:%d: expected varchar '(', found %q\n", tokens[i].Line, tokens[i].Col, string(tokens[i].Val))
	}
	i++
	if !acceptNumber(tokens[i:]) {
		log.Fatalf("%d:%d: expected varchar width, found %q\n", tokens[i].Line, tokens[i].Col, string(tokens[i].Val))
	}
	c.Width = string(tokens[i].Val)
	i++
	if !acceptParenRight(tokens[i:]) {
		for k := i; k < i+4 && k < len(tokens); k++ {
			tmp := tokens[k]
			fmt.Printf("%4d:%3d: unexpected token %q\n", tmp.Line, tmp.Col, string(tmp.Val))
		}
		log.Fatalf("%d:%d: expected varchar ')', found %q\n", tokens[i].Line, tokens[i].Col, string(tokens[i].Val))
	}
	i++
	return &c, i

}
