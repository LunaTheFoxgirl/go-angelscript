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
	"fmt"
)

type ScriptNodeType string

const (
	ASsnUndefined = ScriptNodeType("UNDEFINED")
	ASsnScript = ScriptNodeType("<script>")
	ASsnFunction = ScriptNodeType("<function>")
	ASsnConstant = ScriptNodeType("<constant>")
	ASsnDataType = ScriptNodeType("<data type>")
	ASsnIdentifier = ScriptNodeType("<identifier>")
	ASsnParameterList = ScriptNodeType("<param list>")
	ASsnStatementBlock = ScriptNodeType("<statement block>")
	ASsnDeclaration = ScriptNodeType("<declaration>")
	ASsnExpressionStatement = ScriptNodeType("<expr statement>")
	ASsnIf = ScriptNodeType("<if>")
	ASsnFor = ScriptNodeType("<for>")
	ASsnWhile = ScriptNodeType("<while>")
	ASsnReturn = ScriptNodeType("<return>")
	ASsnExpression = ScriptNodeType("<expr>")
	ASsnExprTerm = ScriptNodeType("<expr term>")
	ASsnFunctionCall = ScriptNodeType("<function call>")
	ASsnConstructCall = ScriptNodeType("<construct call>")
	ASsnArgList = ScriptNodeType("<arg list>")
	ASsnExprPreOp = ScriptNodeType("<expr pre op>")
	ASsnExprPostOp = ScriptNodeType("<expr post op>")
	ASsnExprOperator = ScriptNodeType("<expr operator>")
	ASsnExprValue = ScriptNodeType("<expr value>")
	ASsnBreak = ScriptNodeType("<break>")
	ASsnContinue = ScriptNodeType("<continue>")
	ASsnDoWhile = ScriptNodeType("<do-while>")
	ASsnAssignment = ScriptNodeType("<assignment>")
	ASsnCondition = ScriptNodeType("<condition>")
	ASsnSwitch = ScriptNodeType("<switch>")
	ASsnCase = ScriptNodeType("<case>")
	ASsnImport = ScriptNodeType("<import>")
	ASsnClass = ScriptNodeType("<class>")
	ASsnInitList = ScriptNodeType("<init list>")
	ASsnInterface = ScriptNodeType("<interface>")
	ASsnEnum = ScriptNodeType("<enum>")
	ASsnTypedef = ScriptNodeType("<type def>")
	ASsnCast = ScriptNodeType("<cast>")
	ASsnVariableAccess = ScriptNodeType("<var access>")
	ASsnFuncDef = ScriptNodeType("<function definition>")
	ASsnVirtualProperty = ScriptNodeType("<virtual property>")
	ASsnNamespace = ScriptNodeType("<namespace>")
	ASsnMixin = ScriptNodeType("<mixin>")
	ASsnListPattern = ScriptNodeType("<list pattern>")
	ASsnNamedArgument = ScriptNodeType("<named arg>")
	ASsnScope = ScriptNodeType("<scope>")
)

type sToken struct {
	Type tokens.Token
	Position int
	Length int
}

// ScriptNode is a node in the script containing a range of tokens for compilation.
type ScriptNode struct {
	NodeType ScriptNodeType
	TokenType tokens.Token
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
		TokenType: tokens.ASttUnrecognizedToken,
		TokenPosition: 0,
		TokenLength: 0,
	}
}

func (sn *ScriptNode) ToTList() string {
	f := sn.FirstChild
	//fmt.Println("Going down a level -->")
	var o string = ""
	o = string(sn.NodeType)
	if sn.NodeType == ASsnDataType && sn.TokenType != tokens.ASttUnrecognizedToken { o = tokens.GetDefinition(sn.TokenType)}
	if f != nil {
		for f != nil {
			if f.TokenLength <= 1 { f = f.Next; continue; }
			if f.Previous == nil {
				if f.Next == nil {
					o += " { " + f.ToTList() + " }"
					f = f.Next
					continue
				}
				o += " { " + f.ToTList() + " }"
			} else if f.Next != nil {
				if f.Next.Next == nil {
					o += " & { " + f.ToTList() + " }"
					f = f.Next.Next
					continue
				}
				o += ", { " + f.ToTList() + " }"
			} else {
				o += ", { " + f.ToTList() + " }"
			}
			f = f.Next
		}
	}
	return o
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

func (sn *ScriptNode) UpdateSourcePos(pos, l int) {
	if pos == 0 && l == 0 {
		return
	}
	
	if sn.TokenPosition == 0 && sn.TokenLength == 0 {
		sn.TokenPosition = uint32(pos)
		sn.TokenLength = uint32(l)
	} else {
		if sn.TokenPosition > uint32(pos) {
			sn.TokenLength = sn.TokenPosition + sn.TokenLength - uint32(pos)
			sn.TokenPosition = uint32(pos)
		}
		
		if uint32(pos + l) > sn.TokenPosition + sn.TokenLength {
			sn.TokenLength = uint32(pos) + uint32(l) - sn.TokenPosition
		}
	}
}

func (sn *ScriptNode) AddChildLast(node *ScriptNode, err error) {
	if err != nil {
		fmt.Println("[Angelscript Error]", err)
		return
	}
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
	
	sn.UpdateSourcePos(int(node.TokenPosition), int(node.TokenLength))
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
