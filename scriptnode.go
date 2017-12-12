/*
go-angelscript Copyright (c) 2017 Clipsey

This software is provided 'as-is', without any express or implied warranty. In no event will the authors be held liable for any damages arising from the use of this software.

Permission is granted to anyone to use this software for any purpose, including commercial applications, and to alter it and redistribute it freely, subject to the following restrictions:

    The origin of this software must not be misrepresented; you must not claim that you wrote the original software. If you use this software in a product, an acknowledgment in the product documentation would be appreciated but is not required.

    Altered source versions must be plainly marked as such, and must not be misrepresented as being the original software.

    This notice may not be removed or altered from any source distribution.

The original version of this library can be located at: https://www.github.com/Member1221/go-angelscript

Clipsey clipseypone@gmail.com
*/
package angelscript

import (
	"github.com/Member1221/go-angelscript/tokenizer"
)

type ScriptNodeType uint32

const (
	ASsnUndefined = ScriptNodeType(iota)
	ASsnScript
	ASsnFunction
	ASsnConstant
	ASsnDataType
	ASsnIdentifier
	ASsnParameterList
	ASsnStatementBlock
	ASsnDeclaration
	ASsnExpressionStatement
	ASsnIf
	ASsnFor
	ASsnWhile
	ASsnReturn
	ASsnExpression
	ASsnExprTerm
	ASsnFunctionCall
	ASsnConstructCall
	ASsnArgList
	ASsnExprPreOp
	ASsnExprPostOp
	ASsnExprOperator
	ASsnExprValue
	ASsnBreak
	ASsnContinue
	ASsnDoWhile
	ASsnAssignment
	ASsnCondition
	ASsnSwitch
	ASsnCase
	ASsnImport
	ASsnClass
	ASsnInitList
	ASsnInterface
	ASsnEnum
	ASsnTypedef
	ASsnCast
	ASsnVariableAccess
	ASsnFuncDef
	ASsnVirtualProperty
	ASsnNamespace
	ASsnMixin
	ASsnListPattern
	ASsnNamedArgument
	ASsnScope
)

type sToken struct {
	Type angelscript.Token
	Position int
	Length int
}

type ScriptNode struct {
	NodeType ScriptNodeType
	TokenType angelscript.Token
	TokenPosition uint32
	TokenLength uint32
	Parent *ScriptNode
	Next *ScriptNode
	Previous *ScriptNode
	FirstChild *ScriptNode
	LastChild *ScriptNode
}

func NewScriptNode(t ScriptNodeType) *ScriptNode {
	return &ScriptNode{
		NodeType: t,
		TokenType: angelscript.ASttUnrecognizedToken,
		TokenPosition: 0,
		TokenLength: 0,
	}
}

func (sn *ScriptNode) Destroy(engine *ScriptEngine) {
	n := sn.FirstChild
	
	for n != nil {
		nxt := n.Next
		n.Destroy(engine)
		n = nxt
	}
	
	//TODO: CLEAR STUFF
}

func (sn *ScriptNode) CreateCopy(engine *ScriptEngine) *ScriptNode {
	//TODO: Stub: implement CreateCopy
	return nil
}

func (sn *ScriptNode) SetToken(token *sToken) {
	sn.TokenType = token.Type
}

func (sn *ScriptNode) UpdateSourcePos(pos, l uint32) {
	if pos == 0 && l == 0 {
		return
	}
	
	if sn.TokenPosition == 0 && sn.TokenLength == 0 {
		sn.TokenPosition = pos
		sn.TokenLength = l
	} else {
		if sn.TokenPosition > pos {
			sn.TokenLength = sn.TokenPosition + sn.TokenLength - pos
			sn.TokenPosition = pos
		}
		
		if pos + l > sn.TokenPosition + sn.TokenLength {
			sn.TokenLength = pos + l - sn.TokenPosition
		}
	}
}

func (sn *ScriptNode) AddChildLast(node *ScriptNode) {
	if node == nil {
		return
	}
	
	if sn.LastChild != nil {
		sn.LastChild.Next = node
		node.Next = nil
		node.Previous = sn.LastChild
		node.Parent = sn
		sn.LastChild = node
	} else {
		sn.FirstChild = node
		sn.LastChild = node
		node.Next = nil
		node.Previous = nil
		node.Parent = sn
	}
	
	sn.UpdateSourcePos(node.TokenPosition, node.TokenLength)
}

func (sn *ScriptNode) DisconnectParent() {
	if sn.Parent != nil {
		if sn.Parent.FirstChild == sn {
			sn.Parent.FirstChild = sn.Next
		}
		if sn.Parent.LastChild == sn {
			sn.Parent.LastChild = sn.Previous
		}
	}
	
	if sn.Next != nil {
		sn.Next.Previous = sn.Previous
	}

	if sn.Previous != nil {
			sn.Previous.Next = sn.Next
	}
	
	sn.Parent = nil
	sn.Next = nil
	sn.Previous = nil
}
