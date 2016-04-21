// Copyright 2015 PingCAP, Inc.
//
// Licensed under the Apache License, Version 2.0 (the "License");
// you may not use this file except in compliance with the License.
// You may obtain a copy of the License at
//
//     http://www.apache.org/licenses/LICENSE-2.0
//
// Unless required by applicable law or agreed to in writing, software
// distributed under the License is distributed on an "AS IS" BASIS,
// See the License for the specific language governing permissions and
// limitations under the License.

package ast

import (
	"strings"

	"github.com/pingcap/tidb/expression/builtin"
	"github.com/pingcap/tidb/util/types"
)

var (
	_ FuncNode = &FuncCallExpr{}
	_ FuncNode = &FuncExtractExpr{}
	_ FuncNode = &FuncConvertExpr{}
	_ FuncNode = &FuncCastExpr{}
	_ FuncNode = &FuncSubstringExpr{}
	_ FuncNode = &FuncLocateExpr{}
	_ FuncNode = &FuncTrimExpr{}
	_ FuncNode = &FuncDateArithExpr{}
	_ FuncNode = &AggregateFuncExpr{}
)

// UnquoteString is not quoted when printed.
type UnquoteString string

// FuncCallExpr is for function expression.
type FuncCallExpr struct {
	funcNode
	// F is the function name.
	FnName string
	// Args is the function args.
	Args []ExprNode
}

// Accept implements Node interface.
func (n *FuncCallExpr) Accept(v Visitor) (Node, bool) {
	newNod, skipChildren := v.Enter(n)
	if skipChildren {
		return v.Leave(newNod)
	}
	n = newNod.(*FuncCallExpr)
	for i, val := range n.Args {
		node, ok := val.Accept(v)
		if !ok {
			return n, false
		}
		n.Args[i] = node.(ExprNode)
	}
	return v.Leave(n)
}

// IsStatic implements the ExprNode IsStatic interface.
func (n *FuncCallExpr) IsStatic() bool {
	v := builtin.Funcs[strings.ToLower(n.FnName)]
	if v.F == nil || !v.IsStatic {
		return false
	}

	for _, v := range n.Args {
		if !v.IsStatic() {
			return false
		}
	}
	return true
}

// FuncExtractExpr is for time extract function.
// See https://dev.mysql.com/doc/refman/5.7/en/date-and-time-functions.html#function_extract
type FuncExtractExpr struct {
	funcNode

	Unit string
	Date ExprNode
}

// Accept implements Node Accept interface.
func (n *FuncExtractExpr) Accept(v Visitor) (Node, bool) {
	newNod, skipChildren := v.Enter(n)
	if skipChildren {
		return v.Leave(newNod)
	}
	n = newNod.(*FuncExtractExpr)
	node, ok := n.Date.Accept(v)
	if !ok {
		return n, false
	}
	n.Date = node.(ExprNode)
	return v.Leave(n)
}

// IsStatic implements the ExprNode IsStatic interface.
func (n *FuncExtractExpr) IsStatic() bool {
	return n.Date.IsStatic()
}

// FuncConvertExpr provides a way to convert data between different character sets.
// See: https://dev.mysql.com/doc/refman/5.7/en/cast-functions.html#function_convert
type FuncConvertExpr struct {
	funcNode
	// Expr is the expression to be converted.
	Expr ExprNode
	// Charset is the target character set to convert.
	Charset string
}

// IsStatic implements the ExprNode IsStatic interface.
func (n *FuncConvertExpr) IsStatic() bool {
	return n.Expr.IsStatic()
}

// Accept implements Node Accept interface.
func (n *FuncConvertExpr) Accept(v Visitor) (Node, bool) {
	newNod, skipChildren := v.Enter(n)
	if skipChildren {
		return v.Leave(newNod)
	}
	n = newNod.(*FuncConvertExpr)
	node, ok := n.Expr.Accept(v)
	if !ok {
		return n, false
	}
	n.Expr = node.(ExprNode)
	return v.Leave(n)
}

// CastFunctionType is the type for cast function.
type CastFunctionType int

// CastFunction types
const (
	CastFunction CastFunctionType = iota + 1
	CastConvertFunction
	CastBinaryOperator
)

// FuncCastExpr is the cast function converting value to another type, e.g, cast(expr AS signed).
// See https://dev.mysql.com/doc/refman/5.7/en/cast-functions.html
type FuncCastExpr struct {
	funcNode
	// Expr is the expression to be converted.
	Expr ExprNode
	// Tp is the conversion type.
	Tp *types.FieldType
	// Cast, Convert and Binary share this struct.
	FunctionType CastFunctionType
}

// IsStatic implements the ExprNode IsStatic interface.
func (n *FuncCastExpr) IsStatic() bool {
	return n.Expr.IsStatic()
}

// Accept implements Node Accept interface.
func (n *FuncCastExpr) Accept(v Visitor) (Node, bool) {
	newNod, skipChildren := v.Enter(n)
	if skipChildren {
		return v.Leave(newNod)
	}
	n = newNod.(*FuncCastExpr)
	node, ok := n.Expr.Accept(v)
	if !ok {
		return n, false
	}
	n.Expr = node.(ExprNode)
	return v.Leave(n)
}

// FuncSubstringExpr returns the substring as specified.
// See: https://dev.mysql.com/doc/refman/5.7/en/string-functions.html#function_substring
type FuncSubstringExpr struct {
	funcNode

	StrExpr ExprNode
	Pos     ExprNode
	Len     ExprNode
}

// Accept implements Node Accept interface.
func (n *FuncSubstringExpr) Accept(v Visitor) (Node, bool) {
	newNod, skipChildren := v.Enter(n)
	if skipChildren {
		return v.Leave(newNod)
	}
	n = newNod.(*FuncSubstringExpr)
	node, ok := n.StrExpr.Accept(v)
	if !ok {
		return n, false
	}
	n.StrExpr = node.(ExprNode)
	node, ok = n.Pos.Accept(v)
	if !ok {
		return n, false
	}
	n.Pos = node.(ExprNode)
	if n.Len != nil {
		node, ok = n.Len.Accept(v)
		if !ok {
			return n, false
		}
		n.Len = node.(ExprNode)
	}
	return v.Leave(n)
}

// IsStatic implements the ExprNode IsStatic interface.
func (n *FuncSubstringExpr) IsStatic() bool {
	return n.StrExpr.IsStatic() && n.Pos.IsStatic() && n.Len.IsStatic()
}

// FuncSubstringIndexExpr returns the substring as specified.
// See: https://dev.mysql.com/doc/refman/5.7/en/string-functions.html#function_substring-index
type FuncSubstringIndexExpr struct {
	funcNode

	StrExpr ExprNode
	Delim   ExprNode
	Count   ExprNode
}

// Accept implements Node Accept interface.
func (n *FuncSubstringIndexExpr) Accept(v Visitor) (Node, bool) {
	newNod, skipChildren := v.Enter(n)
	if skipChildren {
		return v.Leave(newNod)
	}
	n = newNod.(*FuncSubstringIndexExpr)
	node, ok := n.StrExpr.Accept(v)
	if !ok {
		return n, false
	}
	n.StrExpr = node.(ExprNode)
	node, ok = n.Delim.Accept(v)
	if !ok {
		return n, false
	}
	n.Delim = node.(ExprNode)
	node, ok = n.Count.Accept(v)
	if !ok {
		return n, false
	}
	n.Count = node.(ExprNode)
	return v.Leave(n)
}

// FuncLocateExpr returns the position of the first occurrence of substring.
// See: https://dev.mysql.com/doc/refman/5.7/en/string-functions.html#function_locate
type FuncLocateExpr struct {
	funcNode

	Str    ExprNode
	SubStr ExprNode
	Pos    ExprNode
}

// Accept implements Node Accept interface.
func (n *FuncLocateExpr) Accept(v Visitor) (Node, bool) {
	newNod, skipChildren := v.Enter(n)
	if skipChildren {
		return v.Leave(newNod)
	}
	n = newNod.(*FuncLocateExpr)
	node, ok := n.Str.Accept(v)
	if !ok {
		return n, false
	}
	n.Str = node.(ExprNode)
	node, ok = n.SubStr.Accept(v)
	if !ok {
		return n, false
	}
	n.SubStr = node.(ExprNode)
	node, ok = n.Pos.Accept(v)
	if !ok {
		return n, false
	}
	n.Pos = node.(ExprNode)
	return v.Leave(n)
}

// TrimDirectionType is the type for trim direction.
type TrimDirectionType int

const (
	// TrimBothDefault trims from both direction by default.
	TrimBothDefault TrimDirectionType = iota
	// TrimBoth trims from both direction with explicit notation.
	TrimBoth
	// TrimLeading trims from left.
	TrimLeading
	// TrimTrailing trims from right.
	TrimTrailing
)

// FuncTrimExpr remove leading/trailing/both remstr.
// See: https://dev.mysql.com/doc/refman/5.7/en/string-functions.html#function_trim
type FuncTrimExpr struct {
	funcNode

	Str       ExprNode
	RemStr    ExprNode
	Direction TrimDirectionType
}

// Accept implements Node Accept interface.
func (n *FuncTrimExpr) Accept(v Visitor) (Node, bool) {
	newNod, skipChildren := v.Enter(n)
	if skipChildren {
		return v.Leave(newNod)
	}
	n = newNod.(*FuncTrimExpr)
	node, ok := n.Str.Accept(v)
	if !ok {
		return n, false
	}
	n.Str = node.(ExprNode)
	node, ok = n.RemStr.Accept(v)
	if !ok {
		return n, false
	}
	n.RemStr = node.(ExprNode)
	return v.Leave(n)
}

// IsStatic implements the ExprNode IsStatic interface.
func (n *FuncTrimExpr) IsStatic() bool {
	return n.Str.IsStatic() && n.RemStr.IsStatic()
}

// DateArithType is type for DateArith type.
type DateArithType byte

const (
	// DateAdd is to run adddate or date_add function option.
	// See: https://dev.mysql.com/doc/refman/5.7/en/date-and-time-functions.html#function_adddate
	// See: https://dev.mysql.com/doc/refman/5.7/en/date-and-time-functions.html#function_date-add
	DateAdd DateArithType = iota + 1
	// DateSub is to run subdate or date_sub function option.
	// See: https://dev.mysql.com/doc/refman/5.7/en/date-and-time-functions.html#function_subdate
	// See: https://dev.mysql.com/doc/refman/5.7/en/date-and-time-functions.html#function_date-sub
	DateSub
	// DateArithDaysForm is to run adddate or subdate function with days form Flag.
	DateArithDaysForm
)

// DateArithInterval is the struct of DateArith interval part.
type DateArithInterval struct {
	// Form is the flag of DateArith running form.
	// The function runs with interval or days.
	Form     DateArithType
	Unit     string
	Interval ExprNode
}

// FuncDateArithExpr is the struct for date arithmetic functions.
type FuncDateArithExpr struct {
	funcNode

	// Op is used for distinguishing date_add and date_sub.
	Op   DateArithType
	Date ExprNode
	DateArithInterval
}

// Accept implements Node Accept interface.
func (n *FuncDateArithExpr) Accept(v Visitor) (Node, bool) {
	newNod, skipChildren := v.Enter(n)
	if skipChildren {
		return v.Leave(newNod)
	}
	n = newNod.(*FuncDateArithExpr)
	if n.Date != nil {
		node, ok := n.Date.Accept(v)
		if !ok {
			return n, false
		}
		n.Date = node.(ExprNode)
	}
	if n.Interval != nil {
		node, ok := n.Interval.Accept(v)
		if !ok {
			return n, false
		}
		n.Interval = node.(ExprNode)
	}
	return v.Leave(n)
}

// AggregateFuncExpr represents aggregate function expression.
type AggregateFuncExpr struct {
	funcNode
	// F is the function name.
	F string
	// Args is the function args.
	Args []ExprNode
	// If distinct is true, the function only aggregate distinct values.
	// For example, column c1 values are "1", "2", "2",  "sum(c1)" is "5",
	// but "sum(distinct c1)" is "3".
	Distinct bool
}

// Accept implements Node Accept interface.
func (n *AggregateFuncExpr) Accept(v Visitor) (Node, bool) {
	newNod, skipChildren := v.Enter(n)
	if skipChildren {
		return v.Leave(newNod)
	}
	n = newNod.(*AggregateFuncExpr)
	for i, val := range n.Args {
		node, ok := val.Accept(v)
		if !ok {
			return n, false
		}
		n.Args[i] = node.(ExprNode)
	}
	return v.Leave(n)
}
