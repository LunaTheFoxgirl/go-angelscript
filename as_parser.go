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
	"errors"
	"fmt"

	"github.com/Member1221/go-angelscript/tokenizer"
)

var ex_debug = false

type Parser struct {
	ErrorWhileParsing     bool
	IsSyntaxError         bool
	CheckValidTypes       bool
	IsParsingAppInterface bool

	Engine  *ScriptEngine
	Script  *ScriptCode
	Node    *ScriptNode
	Builder *ScriptBuilder

	tempstr   string
	lastToken *sToken
	sourcePos int
}

func NewParser(builder *ScriptBuilder) *Parser {
	p := Parser{}
	p.Builder = builder
	p.Engine = builder.Engine
	p.Script = nil
	p.Node = nil
	p.CheckValidTypes = false
	p.IsParsingAppInterface = false
	return &p
}

func (pr *Parser) Reset() {
	pr.ErrorWhileParsing = false
	pr.IsSyntaxError = false
	pr.CheckValidTypes = false
	pr.IsParsingAppInterface = false

	pr.sourcePos = 0

	if pr.Node != nil {
		pr.Node.Destroy(pr.Engine)
	}

	pr.Node = nil
	pr.Script = nil

	pr.lastToken = &sToken{0, 0, 0}
	pr.lastToken.Position = -1
}

func (pr *Parser) GetScriptNode() *ScriptNode {
	return pr.Node
}

func (pr *Parser) CreateNode(t ScriptNodeType) *ScriptNode {
	//TODO: create proper node layout
	/*
		::::C++::::
		void *ptr = engine->memoryMgr.AllocScriptNode();
		if( ptr == 0 )
		{
			// Out of memory
			errorWhileParsing = true;
			return 0;
		}

		return new(ptr) asCScriptNode(type);
	*/
	return NewScriptNode(t)
}

func (pr *Parser) ParseFunctionDefinitionX(script *ScriptCode, expectListPattern bool) int {
	pr.Reset()

	pr.IsParsingAppInterface = true

	pr.Script = script
	var err error
	pr.Node, err = pr.ParseFunctionDefinition()
	if err != nil {
		fmt.Println("[Angelscript Error]", err)
		return -1
	}
	if expectListPattern {
		pr.Node.AddChildLast(pr.ParseListPattern())
	}

	if !pr.IsSyntaxError {
		var t sToken
		pr.GetToken(&t)
		if t.Type != tokens.ASttEnd {
			//TODO: ERROR
		}
	}

	if pr.ErrorWhileParsing {
		return -1
	}

	return 0
}

func (pr *Parser) ParseScriptX(script *ScriptCode) int {
	pr.Reset()

	pr.Script = script
	var err error
	pr.Node, err = pr.ParseScript(false)
	if err != nil {
		fmt.Println("[Angelscript Error]", err)
		return -1
	}

	if pr.ErrorWhileParsing {
		return -1
	}

	return 0
}

func (pr *Parser) ParseStatementBlockX(script *ScriptCode, block *ScriptNode) int {
	pr.Reset()

	//TODO: implement TimeIt function (?)

	pr.CheckValidTypes = true
	pr.Script = script
	pr.sourcePos = int(block.TokenPosition)
	var err error
	pr.Node, err = pr.ParseStatementBlock()
	if err != nil {
		fmt.Println("[Angelscript Error]", err)
		return -1
	}

	if pr.IsSyntaxError || pr.ErrorWhileParsing {
		return -1
	}
	return 0
}

func (pr *Parser) ParseVarInitX(script *ScriptCode, init *ScriptNode) int {
	pr.Reset()

	pr.CheckValidTypes = true

	pr.Script = script
	pr.sourcePos = int(init.TokenPosition)

	var err error
	var t sToken
	pr.GetToken(&t)
	if t.Type == tokens.ASttAssignment {
		pr.GetToken(&t)
		pr.RewindTo(&t)
		if t.Type == tokens.ASttStartStatementBlock {
			pr.Node, err = pr.ParseInitList()
		} else {
			pr.Node, err = pr.ParseAssignment()
		}
	} else if t.Type == tokens.ASttOpenParanthesis {
		pr.RewindTo(&t)
		pr.Node, err = pr.ParseArgList(true)
	} else {
		//TODO: error expected assignment or open paranthesis
	}
	if err != nil {
		fmt.Println("[Angelscript Error]", err)
		return -1
	}
	pr.GetToken(&t)
	if t.Type != tokens.ASttEnd && t.Type != tokens.ASttEndStatement && t.Type != tokens.ASttListSeparator && t.Type != tokens.ASttEndStatementBlock {
		//TODO: error TXT_UNEXPECTED_TOKEN
	}

	if pr.IsSyntaxError || pr.ErrorWhileParsing {
		return -1
	}

	return 0
}

func (pr *Parser) ParseExpressionX(script *ScriptCode) int {
	pr.Reset()
	pr.Script = script
	pr.CheckValidTypes = true

	var err error
	pr.Node, err = pr.ParseExpression()
	if err != nil {
		fmt.Println("[Angelscript Error]", err)
		return -1
	}
	if pr.ErrorWhileParsing {
		return -1
	}
	return 0
}

func (pr *Parser) ParseDataTypeX(script *ScriptCode, isReturnType bool) int {
	pr.Reset()

	pr.Script = script
	pr.Node = pr.CreateNode(ASsnDataType)
	if pr.Node == nil {
		return -1
	}

	pr.Node.AddChildLast(pr.ParseType(true, false, false))
	if pr.IsSyntaxError {
		return -1
	}

	if isReturnType {
		pr.Node.AddChildLast(pr.ParseTypeMod(false))
		if pr.IsSyntaxError {
			return -1
		}
	}

	if !pr.IsSyntaxError {
		var t sToken
		pr.GetToken(&t)
		if t.Type != tokens.ASttEnd {
			//TODO: ERROR
		}
	}

	if pr.ErrorWhileParsing {
		return -1
	}

	return 0
}

func (pr *Parser) ParsePropertyDeclarationX(script *ScriptCode) int {
	pr.Reset()

	pr.Script = script
	pr.Node = pr.CreateNode(ASsnDataType)
	if pr.Node == nil {
		return -1
	}

	pr.Node.AddChildLast(pr.ParseType(true, false, false))
	if pr.IsSyntaxError {
		return -1
	}

	var t sToken
	pr.GetToken(&t)
	pr.RewindTo(&t)
	if t.Type == tokens.ASttAmp {
		pr.Node.AddChildLast(pr.ParseToken(tokens.ASttAmp))
	}

	if !pr.IsSyntaxError {
		if t.Type != tokens.ASttEnd {
			//TODO: ERROR
		}
	}

	if pr.ErrorWhileParsing {
		return -1
	}

	return 0
}

func (pr *Parser) ParseTemplateDecl(script *ScriptCode) int {
	pr.Reset()

	pr.Script = script
	pr.Node = pr.CreateNode(ASsnDataType)
	if pr.Node == nil {
		return -1
	}

	pr.Node.AddChildLast(pr.ParseType(true, false, false))
	if pr.IsSyntaxError {
		return -1
	}

	if !pr.IsSyntaxError {
		var t sToken
		pr.GetToken(&t)
		if t.Type != tokens.ASttEnd {
			//TODO: ERROR
		}
	}

	if pr.ErrorWhileParsing {
		return -1
	}

	return 0
}

func (pr *Parser) ParseOptionalScope(script *ScriptNode) {
	scope := pr.CreateNode(ASsnScope)
	var t1 sToken
	var t2 sToken
	pr.GetToken(&t1)
	pr.GetToken(&t2)
	if t1.Type == tokens.ASttScope {
		pr.RewindTo(&t1)
		scope.AddChildLast(pr.ParseToken(tokens.ASttScope))
		pr.GetToken(&t1)
		pr.GetToken(&t2)
	}
	for t1.Type == tokens.ASttIdentifier && t2.Type == tokens.ASttScope {
		pr.RewindTo(&t1)
		scope.AddChildLast(pr.ParseIdentifier())
		scope.AddChildLast(pr.ParseToken(tokens.ASttScope))
		pr.GetToken(&t1)
		pr.GetToken(&t2)
	}

	if t1.Type == tokens.ASttIdentifier && t2.Type == tokens.ASttLessThan {
		pr.tempstr = pr.Script.Code[t1.Position : t1.Position+t1.Length]
		if pr.Engine.IsTemplateType(pr.tempstr) {

			pr.RewindTo(&t1)
			restore := scope.LastChild
			scope.AddChildLast(pr.ParseIdentifier())
			t, err := pr.ParseTemplTypeList(scope, false)
			if err != nil {
				panic(err)
			}
			if t {
				pr.GetToken(&t2)
				if t2.Type == tokens.ASttScope {
					pr.Node.AddChildLast(scope, nil)
					return
				} else {
					pr.RewindTo(&t1)

					for scope.LastChild != restore {
						last := scope.LastChild
						last.DisconnectParent()
						last.Destroy(pr.Engine)
					}
					if scope.LastChild != nil {
						pr.Node.AddChildLast(scope, nil)
					} else {
						scope.Destroy(pr.Engine)
					}
				}
				return
			}
		}
	}

	pr.RewindTo(&t1)

	if scope.LastChild != nil {
		pr.Node.AddChildLast(scope, nil)
	} else {
		scope.Destroy(pr.Engine)
	}
}

func (pr *Parser) ParseFunctionDefinition() (*ScriptNode, error) {
	node := pr.CreateNode(ASsnFunction)
	if node == nil {
		return nil, nil
	}

	node.AddChildLast(pr.ParseType(true, false, false))
	if pr.IsSyntaxError {
		return node, nil
	}

	node.AddChildLast(pr.ParseTypeMod(false))
	if pr.IsSyntaxError {
		return node, nil
	}

	pr.ParseOptionalScope(node)

	node.AddChildLast(pr.ParseIdentifier())
	if pr.IsSyntaxError {
		return node, nil
	}

	node.AddChildLast(pr.ParseParameterList())
	if pr.IsSyntaxError {
		return node, nil
	}

	var t1 sToken
	pr.GetToken(&t1)
	pr.RewindTo(&t1)
	if t1.Type == tokens.ASttConst {
		node.AddChildLast(pr.ParseToken(tokens.ASttConst))
	}
	return node, nil
}

func (pr *Parser) ParseTypeMod(isParam bool) (*ScriptNode, error) {
	node := pr.CreateNode(ASsnDataType)
	if node == nil {
		return nil, nil
	}

	var t sToken

	pr.GetToken(&t)
	pr.RewindTo(&t)

	if t.Type == tokens.ASttAmp {
		node.AddChildLast(pr.ParseToken(tokens.ASttAmp))
		if pr.IsSyntaxError {
			return node, nil
		}
		if isParam {
			pr.GetToken(&t)
			pr.RewindTo(&t)

			if t.Type == tokens.ASttIn || t.Type == tokens.ASttOut || t.Type == tokens.ASttInOut {
				tokens := []tokens.Token{tokens.ASttIn, tokens.ASttOut, tokens.ASttInOut}
				node.AddChildLast(pr.ParseOneOf(tokens))
			}
		}
	}

	pr.GetToken(&t)
	pr.RewindTo(&t)
	if t.Type == tokens.ASttPlus {
		node.AddChildLast(pr.ParseToken(tokens.ASttPlus))
		if pr.IsSyntaxError {
			return node, nil
		}
	}

	pr.GetToken(&t)
	pr.RewindTo(&t)
	if pr.IdentifierIs(t, tokens.ASIfHandleToken) {
		node.AddChildLast(pr.ParseToken(tokens.ASttIdentifier))
		if pr.IsSyntaxError {
			return node, nil
		}
	}

	return node, nil
}

func (pr *Parser) ParseType(allowConst, allowVariableType, allowAuto bool) (*ScriptNode, error) {
	node := pr.CreateNode(ASsnDataType)
	if node == nil {
		return nil, nil
	}

	var t sToken

	if allowConst {
		pr.GetToken(&t)
		pr.RewindTo(&t)
		if t.Type == tokens.ASttConst {
			node.AddChildLast(pr.ParseToken(tokens.ASttConst))
			if pr.IsSyntaxError {
				return node, nil
			}
		}
	}

	pr.ParseOptionalScope(node)

	node.AddChildLast(pr.ParseDataType(allowVariableType, allowAuto))
	if pr.IsSyntaxError {
		return node, nil
	}

	pr.GetToken(&t)
	pr.RewindTo(&t)
	tr := node.LastChild

	pr.tempstr = pr.Script.Code[tr.TokenPosition : tr.TokenPosition+tr.TokenLength]
	if pr.Engine.IsTemplateType(pr.tempstr) && t.Type == tokens.ASttLessThan {
		pr.ParseTemplTypeList(node, true)
		if pr.IsSyntaxError {
			return node, nil
		}
	}

	pr.GetToken(&t)
	pr.RewindTo(&t)
	for t.Type == tokens.ASttOpenBracket || t.Type == tokens.ASttHandle {
		if t.Type == tokens.ASttOpenBracket {
			node.AddChildLast(pr.ParseToken(tokens.ASttOpenBracket))
			if pr.IsSyntaxError {
				return node, nil
			}

			pr.GetToken(&t)
			if t.Type == tokens.ASttCloseBracket {
				//TODO: ERROR (expect ])
				return node, nil
			}
		} else {
			node.AddChildLast(pr.ParseToken(tokens.ASttHandle))
			if pr.IsSyntaxError {
				return node, nil
			}
		}

		pr.GetToken(&t)
		pr.RewindTo(&t)
	}
	return node, nil
}

func (pr *Parser) ParseTemplTypeList(node *ScriptNode, required bool) (bool, error) {
	var t sToken
	isValid := true

	last := node.LastChild

	pr.GetToken(&t)

	if t.Type != tokens.ASttLessThan {
		if required {
			//TODO: ERROR (expect Lessthan)
			return false, errors.New("Expected " + tokens.GetDefinition(tokens.ASttLessThan) + ", got " + tokens.GetDefinition(t.Type))
		}
		return false, nil
	}

	node.AddChildLast(pr.ParseType(true, false, false))
	if pr.IsSyntaxError {
		return false, nil
	}

	pr.GetToken(&t)

	for t.Type == tokens.ASttListSeparator {
		node.AddChildLast(pr.ParseType(true, false, false))
		if pr.IsSyntaxError {
			return false, nil
		}
		pr.GetToken(&t)
	}

	if pr.Script.Code[t.Position:1] != ">" {
		if required {
			return false, errors.New("Expected" + tokens.GetDefinition(tokens.ASttGreaterThan) + ", got " + tokens.GetDefinition(t.Type))
		} else {
			isValid = false
		}
	} else {
		pr.SetPos(t.Position + 1)
	}

	if !required && !isValid {
		for node.LastChild != last {
			n := node.LastChild
			n.DisconnectParent()
			n.Destroy(pr.Engine)
		}
		return false, nil
	}

	return true, nil

}

func (pr *Parser) ParseToken(token tokens.Token) (*ScriptNode, error) {
	node := pr.CreateNode(ASsnUndefined)
	if node == nil {
		return nil, nil
	}

	var t1 sToken

	pr.GetToken(&t1)

	if t1.Type != token {
		//TODO: ERROR (Expect TOKEN)
		return nil, errors.New("Expected " + tokens.GetDefinition(token) + ", got " + tokens.GetDefinition(t1.Type))
	}

	node.SetToken(&t1)
	node.UpdateSourcePos(t1.Position, t1.Length)
	return node, nil
}

func (pr *Parser) ParseOneOf(toks []tokens.Token) (*ScriptNode, error) {
	node := pr.CreateNode(ASsnUndefined)
	if node == nil {
		return nil, nil
	}

	var t1 sToken
	pr.GetToken(&t1)

	n := 0
	for n = 0; n < len(toks); n++ {
		if toks[n] == t1.Type {
			break
		}
	}

	if n == len(toks) {
		//TODO: ERROR (Expect tokens/count, got t1)
		return node, errors.New("Expected " + tokens.GetDefinitionOrList(toks) + ", got " + tokens.GetDefinition(t1.Type))
	}

	node.SetToken(&t1)
	node.UpdateSourcePos(t1.Position, t1.Length)

	return node, nil
}

func (pr *Parser) ParseDataType(allowVariableType, allowAuto bool) (*ScriptNode, error) {
	node := pr.CreateNode(ASsnDataType)
	if node == nil {
		return nil, nil
	}

	var t1 sToken
	pr.GetToken(&t1)
	if !pr.IsDataType(&t1) && !(allowVariableType && t1.Type == tokens.ASttQuestion) && !(allowAuto && t1.Type == tokens.ASttAuto) {
		if t1.Type == tokens.ASttIdentifier {
			//TODO: FATAL ERROR:
			/*
				asCString errMsg;
				tempString.Assign(&script->code[t1.pos], t1.length);
				errMsg.Format(TXT_IDENTIFIER_s_NOT_DATA_TYPE, tempString.AddressOf());
				Error(errMsg, &t1);
			*/
			return node, errors.New("Identifier(s) not a data type, got " + tokens.GetDefinition(t1.Type))
		} else if t1.Type == tokens.ASttAuto {
			//TODO: ERROR TXT_AUTO_NOT_ALLOWED
			return node, errors.New("Auto not allowed, got " + tokens.GetDefinition(t1.Type))
		} else {
			//TODO: ERROR TXT_EXPECTED_DATA_TYPE
			return node, errors.New("Expected data type, got " + tokens.GetDefinition(t1.Type))
		}
	}

	node.SetToken(&t1)
	node.UpdateSourcePos(t1.Position, t1.Length)
	return node, nil
}

func (pr *Parser) ParseRealType() (*ScriptNode, error) {
	node := pr.CreateNode(ASsnDataType)
	if node == nil {
		return nil, nil
	}

	var t1 sToken
	pr.GetToken(&t1)
	if !pr.IsRealType(t1.Type) {
		//TODO: ERROR TXT_EXPECTED_DATATYPE
		return node, errors.New("Expected data type, got " + tokens.GetDefinition(t1.Type))
	}

	node.SetToken(&t1)
	node.UpdateSourcePos(t1.Position, t1.Length)
	return node, nil
}

func (pr *Parser) ParseIdentifier() (*ScriptNode, error) {
	node := pr.CreateNode(ASsnIdentifier)
	if node == nil {
		return nil, nil
	}

	var t1 sToken
	pr.GetToken(&t1)
	if t1.Type != tokens.ASttIdentifier {
		//TODO: ERROR TXT_EXPECTED_DATATYPE
		return node, errors.New("Expected data type, got " + tokens.GetDefinition(t1.Type))
	}

	node.SetToken(&t1)
	node.UpdateSourcePos(t1.Position, t1.Length)
	return node, nil
}

func (pr *Parser) ParseParameterList() (*ScriptNode, error) {
	node := pr.CreateNode(ASsnParameterList)
	if node == nil {
		return nil, nil
	}

	var t1 sToken
	pr.GetToken(&t1)
	if t1.Type != tokens.ASttOpenParanthesis {
		//TODO: ERROR (Expected "(")
		return nil, errors.New("Expected " + tokens.GetDefinitionOrList([]tokens.Token{tokens.ASttOpenParanthesis}) + ", got " + tokens.GetDefinition(t1.Type))
	}

	node.UpdateSourcePos(t1.Position, t1.Length)
	pr.GetToken(&t1)
	if t1.Type == tokens.ASttCloseParanthesis {
		node.UpdateSourcePos(t1.Position, t1.Length)
		return node, nil
	} else {
		if t1.Type == tokens.ASttVoid {
			var t2 sToken
			pr.GetToken(&t2)
			if t2.Type == tokens.ASttCloseParanthesis {
				node.UpdateSourcePos(t2.Position, t2.Length)
				return node, nil
			}
		}

		pr.RewindTo(&t1)

		for {
			node.AddChildLast(pr.ParseType(true, pr.IsParsingAppInterface, false))
			if pr.IsSyntaxError {
				return nil, nil
			}

			node.AddChildLast(pr.ParseTypeMod(true))
			if pr.IsSyntaxError {
				return nil, nil
			}

			pr.GetToken(&t1)
			if t1.Type == tokens.ASttIdentifier {
				pr.RewindTo(&t1)

				node.AddChildLast(pr.ParseIdentifier())
				if pr.IsSyntaxError {
					return node, nil
				}

				pr.GetToken(&t1)
			}

			if t1.Type == tokens.ASttAssignment {
				node.AddChildLast(pr.SuperficiallyParseExpression())
				if pr.IsSyntaxError {
					return node, nil
				}

				pr.GetToken(&t1)
			}

			if t1.Type == tokens.ASttCloseParanthesis {
				node.UpdateSourcePos(t1.Position, t1.Length)
			} else if t1.Type == tokens.ASttListSeparator {
				continue
			} else {
				//TODO: Error (Expected Tokens: ")", ",")
				return nil, errors.New("Expected " + tokens.GetDefinitionOrList([]tokens.Token{tokens.ASttCloseParanthesis, tokens.ASttListSeparator}) + ", got " + tokens.GetDefinition(t1.Type))
			}
		}
	}
	return node, nil
}

func (pr *Parser) SuperficiallyParseExpression() (*ScriptNode, error) {
	node := pr.CreateNode(ASsnParameterList)
	if node == nil {
		return nil, nil
	}

	var start sToken
	pr.GetToken(&start)
	pr.RewindTo(&start)

	stack := ""
	var t sToken
	for {
		pr.GetToken(&t)

		if t.Type == tokens.ASttOpenParanthesis {
			stack += "("
		} else if t.Type == tokens.ASttCloseParanthesis {
			if stack == "" {
				pr.RewindTo(&t)
				break
			} else if stack[len(stack)-1:1] == "(" {
				stack = stack[:len(stack)-1]
			} else {
				pr.RewindTo(&t)
				//TODO: FATAL ERROR
				/*
					asCString str;
					str.Format(TXT_UNEXPECTED_TOKEN_s, ")");
					Error(str, &t);
				*/
				return node, nil
			}
		} else if t.Type == tokens.ASttListSeparator {
			if stack == "" {
				pr.RewindTo(&t)
				break
			}
		} else if t.Type == tokens.ASttStartStatementBlock {
			stack += "{"
		} else if t.Type == tokens.ASttEndStatementBlock {
			if stack == "" || stack[len(stack)-1:1] == "{" {
				pr.RewindTo(&t)
				/*
					asCString str;
					str.Format(TXT_UNEXPECTED_TOKEN_s, "}");
					Error(str, &t);
				*/
				return node, nil
			} else {
				stack = stack[:len(stack)-1]
			}
		} else if t.Type == tokens.ASttEndStatement {
			pr.RewindTo(&t)
			/*
				asCString str;
				str.Format(TXT_UNEXPECTED_TOKEN_s, ";");
				Error(str, &t);
			*/
			return node, nil
		} else if t.Type == tokens.ASttNonTerminatedStringConstant {
			pr.RewindTo(&t)
			//TODO: ERROR (TXT_NONTERMINATED_STRING)
			return node, nil
		} else if t.Type == tokens.ASttEnd {
			pr.RewindTo(&t)
			//TODO: ERROR (TXT_UNEXPECTED_END_OF_FILE)
			return node, nil
		}

		node.UpdateSourcePos(t.Position, t.Length)
	}

	return node, nil
}

func (pr *Parser) GetToken(token *sToken) {
	if pr.lastToken.Position == pr.sourcePos {
		token.Length = pr.lastToken.Length
		token.Position = pr.lastToken.Position
		token.Type = pr.lastToken.Type
		if ex_debug {
			if token.Type != tokens.ASttWhiteSpace {
				fmt.Println(tokens.GetDefinition(token.Type), "of length", token.Length, "at position", token.Position)
			}
		}
		pr.sourcePos += token.Length

		if token.Type == tokens.ASttWhiteSpace ||
			token.Type == tokens.ASttOnelineComment ||
			token.Type == tokens.ASttMultilineComment {
			pr.GetToken(token)
		}
		return
	}

	if pr.sourcePos >= len(pr.Script.Code) {
		token.Type = tokens.ASttEnd
		token.Length = 0
	} else {
		l, t := pr.Engine.tok.GetToken(pr.Script.Code[pr.sourcePos:])
		token.Type = t
		token.Length = int(l)
		if ex_debug {
			if token.Type != tokens.ASttWhiteSpace {
				fmt.Println(tokens.GetDefinition(token.Type), "of length", token.Length, "at position", token.Position)
			}
		}
	}
	token.Position = pr.sourcePos
	pr.sourcePos += token.Length
	for token.Type == tokens.ASttWhiteSpace ||
		token.Type == tokens.ASttOnelineComment ||
		token.Type == tokens.ASttMultilineComment {
		if pr.sourcePos >= len(pr.Script.Code) {
			token.Type = tokens.ASttEnd
			token.Length = 0
		} else {
			l, t := pr.Engine.tok.GetToken(pr.Script.Code[pr.sourcePos:])
			token.Type = t
			token.Length = int(l)
			if ex_debug {
				if token.Type != tokens.ASttWhiteSpace {
					fmt.Println(tokens.GetDefinition(token.Type), "of length", token.Length, "at position", token.Position)
				}
			}
		}

		token.Position = pr.sourcePos
		pr.sourcePos += token.Length
	}

}

func (pr *Parser) RewindTo(token *sToken) {
	if ex_debug {
		fmt.Println("Rewinding to", token.Position, "with type", tokens.GetDefinition(token.Type), "of length", token.Length, "...")
	}
	pr.lastToken.Position = token.Position
	pr.lastToken.Length = token.Length
	pr.lastToken.Type = token.Type
	pr.sourcePos = token.Position
}

func (pr *Parser) SetPos(pos int) {
	pr.lastToken = nil
	pr.sourcePos = pos
}

func (pr *Parser) Error(text string) {

}

func (pr *Parser) Warning(text string) {

}

func (pr *Parser) Info(text string) {

}

func (pr *Parser) IsRealType(tokenType tokens.Token) bool {
	if tokenType == tokens.ASttVoid ||
		tokenType == tokens.ASttInt ||
		tokenType == tokens.ASttInt8 ||
		tokenType == tokens.ASttInt16 ||
		tokenType == tokens.ASttInt64 ||
		tokenType == tokens.ASttUInt ||
		tokenType == tokens.ASttUInt8 ||
		tokenType == tokens.ASttUInt16 ||
		tokenType == tokens.ASttUInt64 ||
		tokenType == tokens.ASttFloat ||
		tokenType == tokens.ASttDouble ||
		tokenType == tokens.ASttBool {
		return true
	}
	return false
}

func (pr *Parser) IsDataType(token *sToken) bool {
	if token.Type == tokens.ASttIdentifier {
		if pr.CheckValidTypes {
			pr.tempstr = pr.Script.Code[token.Position : token.Position+token.Length]
			if !pr.Builder.DoesTypeExist(pr.tempstr) {
				return false
			}
		}
		return true
	}
	if pr.IsRealType(token.Type) {
		return true
	}
	return false
}

func (pr *Parser) ExpectedToken(token string) string {
	return ""
}

func (pr *Parser) ExpectedTokens(tokena, tokenb string) string {
	return ""
}

func (pr *Parser) ExpectedOneOf(tokens []string) string {
	return ""
}

func (pr *Parser) ExpectedOneOfMap(tokens map[string][]tokens.Token) string {
	return ""
}

func (pr *Parser) InsteadFound(token sToken) string {
	return ""
}

func (pr *Parser) ParseListPattern() (*ScriptNode, error) {
	node := pr.CreateNode(ASsnParameterList)
	if node == nil {
		return nil, nil
	}

	var t1 sToken
	pr.GetToken(&t1)
	if t1.Type != tokens.ASttStartStatementBlock {
		//TODO: ERROR (Expected "{")
		return node, nil
	}

	node.UpdateSourcePos(t1.Position, t1.Length)
	//var start sToken

	isBeginning := true
	afterType := false

	for !pr.IsSyntaxError {
		pr.GetToken(&t1)
		if t1.Type == tokens.ASttEndStatementBlock {
			if !afterType {
				//TODO: ERROR TXT_EXPECTED_DATA_TYPE
			}
			break
		} else if t1.Type == tokens.ASttStartStatementBlock {
			if afterType {
				//TODO: ERROR (Expected ",", "}")
			}
			pr.RewindTo(&t1)
			node.AddChildLast(pr.ParseListPattern())
			afterType = true
		} else if t1.Type == tokens.ASttIdentifier && (pr.IdentifierIs(t1, "repeat") || pr.IdentifierIs(t1, "repeat_same")) {
			if !isBeginning {
				/*
					asCString msg;
					asCString token(&script->code[t1.pos], t1.length);
					msg.Format(TXT_UNEXPECTED_TOKEN_s, token.AddressOf());
					Error(msg.AddressOf(), &t1);
				*/
			}
			pr.RewindTo(&t1)
			node.AddChildLast(pr.ParseIdentifier())
		} else if t1.Type == tokens.ASttEnd {
			//TODO: ERROR TXT_UNEXPECTED_END_OF_FILE
			break
		} else if t1.Type == tokens.ASttListSeparator {
			if !afterType {
				//TODO:: ERROR TXT_EXPECTED_DATA_TYPE
			}
			afterType = false
		} else {
			if afterType {
				//TODO: ERROR (Expected ",", "}")
			}
			pr.RewindTo(&t1)
			node.AddChildLast(pr.ParseType(true, true, false))
			afterType = true
		}
		isBeginning = false
	}

	node.UpdateSourcePos(t1.Position, t1.Length)

	return node, nil
}

func (pr *Parser) IdentifierIs(t sToken, str string) bool {
	if t.Type != tokens.ASttIdentifier {
		return false
	}
	return pr.Script.TokenEquals(t.Position, t.Length, str)
}

func (pr *Parser) CheckTemplateType(t sToken) bool {
	pr.tempstr = pr.Script.Code[t.Position : t.Position+t.Length]
	if pr.Engine.IsTemplateType(pr.tempstr) {
		var t1 sToken
		pr.GetToken(&t1)
		if t1.Type == tokens.ASttLessThan {
			pr.RewindTo(&t1)
			return true
		}

		for {
			pr.GetToken(&t1)
			if t1.Type == tokens.ASttScope {
				pr.GetToken(&t1)
			}

			var t2 sToken
			pr.GetToken(&t2)
			for t1.Type == tokens.ASttIdentifier && t2.Type == tokens.ASttScope {
				pr.GetToken(&t1)
				pr.GetToken(&t2)
			}
			pr.RewindTo(&t2)

			if !pr.IsDataType(&t1) {
				return false
			}
			if !pr.CheckTemplateType(t1) {
				return false
			}

			pr.GetToken(&t1)

			for t1.Type == tokens.ASttHandle || t1.Type == tokens.ASttOpenBracket {
				if t1.Type == tokens.ASttOpenBracket {
					pr.GetToken(&t1)
					if t1.Type != tokens.ASttCloseBracket {
						return false
					}
				}

				pr.GetToken(&t1)
			}

			if t1.Type == tokens.ASttListSeparator {
				break
			}
		}

		if pr.Script.Code[t1.Position:1] == ">" {
			return false
		} else if t1.Length != 1 {
			pr.SetPos(t1.Position + 1)
		}
	}
	return true
}

func (pr *Parser) ParseCast() (*ScriptNode, error) {
	node := pr.CreateNode(ASsnCast)
	if node == nil {
		return nil, nil
	}

	var t1 sToken
	pr.GetToken(&t1)
	if t1.Type == tokens.ASttCast {
		//TODO: error expected cast
		return node, nil
	}

	node.UpdateSourcePos(t1.Position, t1.Length)

	pr.GetToken(&t1)
	if t1.Type != tokens.ASttLessThan {
		//TODO: error expected <
		return node, nil
	}

	node.AddChildLast(pr.ParseType(true, false, false))
	if pr.IsSyntaxError {
		return node, nil
	}

	pr.GetToken(&t1)
	if t1.Type != tokens.ASttGreaterThan {
		//TODO: error expected >
	}

	pr.GetToken(&t1)
	if t1.Type != tokens.ASttOpenParanthesis {
		errors.New("Expected " + tokens.GetDefinitionOrList([]tokens.Token{tokens.ASttOpenParanthesis}) + ", got " + tokens.GetDefinition(t1.Type))
		return node, nil
	}

	node.AddChildLast(pr.ParseAssignment())
	if pr.IsSyntaxError {
		return node, nil
	}

	pr.GetToken(&t1)
	if t1.Type != tokens.ASttCloseParanthesis {
		errors.New("Expected " + tokens.GetDefinitionOrList([]tokens.Token{tokens.ASttCloseParanthesis}) + ", got " + tokens.GetDefinition(t1.Type))
		return node, nil
	}

	node.UpdateSourcePos(t1.Position, t1.Length)

	return node, nil
}

func (pr *Parser) ParseExprValue() (*ScriptNode, error) {
	node := pr.CreateNode(ASsnExprValue)
	if node == nil {
		return nil, nil
	}

	var t1 sToken
	var t2 sToken
	pr.GetToken(&t1)
	pr.GetToken(&t2)
	pr.RewindTo(&t1)

	if t1.Type == tokens.ASttVoid {
		node.AddChildLast(pr.ParseToken(tokens.ASttVoid))
	} else if pr.IsRealType(t1.Type) {
		node.AddChildLast(pr.ParseConstructCall())
	} else if t1.Type == tokens.ASttIdentifier || t1.Type == tokens.ASttScope {
		if pr.IsLambda() {
			node.AddChildLast(pr.ParseLambda())
		} else {
			var t sToken
			if t1.Type == tokens.ASttScope {
				t = t2
			} else {
				t = t1
			}
			pr.RewindTo(&t)
			pr.GetToken(&t2)
			for t.Type == tokens.ASttIdentifier {
				t2 = t
				pr.GetToken(&t)
				if t.Type == tokens.ASttScope {
					pr.GetToken(&t)
				} else {
					break
				}
			}

			isDataType := pr.IsDataType(&t2)
			isTemplateType := false

			if isDataType {
				pr.tempstr = pr.Script.Code[t2.Position : t2.Position+t2.Length]
				if pr.Engine.IsTemplateType(pr.tempstr) {
					isTemplateType = true
				}
			}

			pr.GetToken(&t2)

			pr.RewindTo(&t1)

			if isDataType && (t.Type == tokens.ASttOpenParanthesis || (t.Type == tokens.ASttOpenBracket && t2.Type == tokens.ASttCloseBracket)) {
				node.AddChildLast(pr.ParseConstructCall())
			} else if isTemplateType && t.Type == tokens.ASttLessThan {
				node.AddChildLast(pr.ParseConstructCall())
			} else if pr.IsFunctionCall() {
				node.AddChildLast(pr.ParseFunctionCall())
			} else {
				node.AddChildLast(pr.ParseVariableAccess())
			}
		}
	} else if t1.Type == tokens.ASttCast {
		node.AddChildLast(pr.ParseCast())
	} else if pr.IsConstant(t1.Type) {
		node.AddChildLast(pr.ParseConstant())
	} else if t1.Type == tokens.ASttOpenParanthesis {
		pr.GetToken(&t1)
		node.UpdateSourcePos(t1.Position, t1.Length)

		node.AddChildLast(pr.ParseAssignment())
		if pr.IsSyntaxError {
			return node, nil
		}

		pr.GetToken(&t1)
		if t1.Type != tokens.ASttCloseParanthesis {
			return node, errors.New("Expected " + tokens.GetDefinitionOrList([]tokens.Token{tokens.ASttCloseParanthesis}) + ", got " + tokens.GetDefinition(t1.Type))
		}

		node.UpdateSourcePos(t1.Position, t1.Length)
	} else {
		//TODO: error TXT_EXPECTED_EXPRESSION_VALUE
	}

	return node, nil
}

func (pr *Parser) ParseConstant() (*ScriptNode, error) {
	node := pr.CreateNode(ASsnConstant)
	if node == nil {
		return nil, nil
	}

	var t sToken
	pr.GetToken(&t)
	if !(pr.IsConstant(t.Type)) {
		//TODO: error TXT_EXPECTED_CONSTANT
		return node, nil
	}

	node.SetToken(&t)
	node.UpdateSourcePos(t.Position, t.Length)

	if t.Type == tokens.ASttStringConstant || t.Type == tokens.ASttMultilineStringConstant || t.Type == tokens.ASttHeredocStringConstant {
		pr.RewindTo(&t)
	}

	for t.Type == tokens.ASttStringConstant || t.Type == tokens.ASttMultilineStringConstant || t.Type == tokens.ASttHeredocStringConstant {
		node.AddChildLast(pr.ParseStringConstant())

		pr.GetToken(&t)
		pr.RewindTo(&t)
	}

	return node, nil
}

func (pr *Parser) IsLambda() bool {
	isLambda := false
	var t sToken
	pr.GetToken(&t)
	if t.Type == tokens.ASttIdentifier && pr.IdentifierIs(t, tokens.ASFunctionToken) {
		var t2 sToken
		pr.GetToken(&t2)
		if t2.Type == tokens.ASttOpenParanthesis {
			for t2.Type != tokens.ASttCloseParanthesis && t2.Type != tokens.ASttEnd {
				pr.GetToken(&t2)
			}

			pr.GetToken(&t2)

			if t2.Type == tokens.ASttStartStatementBlock {
				isLambda = true
			}
		}
	}

	pr.RewindTo(&t)
	return isLambda
}

func (pr *Parser) ParseLambda() (*ScriptNode, error) {
	node := pr.CreateNode(ASsnFunction)
	if node == nil {
		return nil, nil
	}

	var t sToken
	pr.GetToken(&t)

	if t.Type != tokens.ASttIdentifier || !pr.IdentifierIs(t, tokens.ASFunctionToken) {
		//TODO: error expected function token
		return node, nil
	}

	pr.GetToken(&t)
	if t.Type != tokens.ASttOpenParanthesis {
		errors.New("Expected " + tokens.GetDefinitionOrList([]tokens.Token{tokens.ASttOpenParanthesis}) + ", got " + tokens.GetDefinition(t.Type))
		return node, nil
	}

	pr.GetToken(&t)
	if t.Type == tokens.ASttIdentifier {
		pr.RewindTo(&t)
		node.AddChildLast(pr.ParseIdentifier())

		pr.GetToken(&t)
		for t.Type == tokens.ASttListSeparator {
			node.AddChildLast(pr.ParseIdentifier())
			if pr.IsSyntaxError {
				return node, nil
			}

			pr.GetToken(&t)
		}
	}

	if t.Type != tokens.ASttCloseParanthesis {
		errors.New("Expected " + tokens.GetDefinitionOrList([]tokens.Token{tokens.ASttCloseParanthesis}) + ", got " + tokens.GetDefinition(t.Type))
		return node, nil
	}

	node.AddChildLast(pr.SuperficiallyParseStatementBlock())

	return node, nil
}

func (pr *Parser) ParseStringConstant() (*ScriptNode, error) {
	node := pr.CreateNode(ASsnConstant)
	if node == nil {
		return nil, nil
	}

	var t sToken
	pr.GetToken(&t)
	if t.Type != tokens.ASttStringConstant && t.Type != tokens.ASttMultilineComment && t.Type != tokens.ASttHeredocStringConstant {
		//TODO: error TXT_EXPECTED_STRING
		return node, nil
	}

	node.SetToken(&t)
	node.UpdateSourcePos(t.Position, t.Length)

	return node, nil
}

func (pr *Parser) ParseFunctionCall() (*ScriptNode, error) {
	node := pr.CreateNode(ASsnFunctionCall)
	if node == nil {
		return nil, nil
	}

	pr.ParseOptionalScope(node)

	node.AddChildLast(pr.ParseIdentifier())
	if pr.IsSyntaxError {
		return node, nil
	}

	node.AddChildLast(pr.ParseArgList(true))

	return node, nil
}

func (pr *Parser) ParseVariableAccess() (*ScriptNode, error) {
	node := pr.CreateNode(ASsnVariableAccess)
	if node == nil {
		return nil, nil
	}

	pr.ParseOptionalScope(node)

	node.AddChildLast(pr.ParseIdentifier())

	return node, nil
}

func (pr *Parser) ParseConstructCall() (*ScriptNode, error) {
	node := pr.CreateNode(ASsnConstructCall)
	if node == nil {
		return nil, nil
	}

	node.AddChildLast(pr.ParseType(false, false, false))
	if pr.IsSyntaxError {
		return node, nil
	}

	node.AddChildLast(pr.ParseArgList(true))

	return node, nil
}

func (pr *Parser) ParseArgList(withParenthesis bool) (*ScriptNode, error) {
	node := pr.CreateNode(ASsnArgList)
	if node == nil {
		return nil, nil
	}

	var t1 sToken
	if withParenthesis {
		pr.GetToken(&t1)
		if t1.Type != tokens.ASttOpenParanthesis {
			errors.New("Expected " + tokens.GetDefinitionOrList([]tokens.Token{tokens.ASttOpenParanthesis}) + ", got " + tokens.GetDefinition(t1.Type))
			return node, nil
		}
		node.UpdateSourcePos(t1.Position, t1.Length)
	}

	pr.GetToken(&t1)
	if t1.Type == tokens.ASttCloseParanthesis || t1.Type == tokens.ASttCloseBracket {
		if withParenthesis {
			if t1.Type == tokens.ASttCloseParanthesis {
				node.UpdateSourcePos(t1.Position, t1.Length)
			} else {
				/*
					asCString str;
					str.Format(TXT_UNEXPECTED_TOKEN_s, asCTokenizer::GetDefinition(ttCloseBracket));

					Error(str.AddressOf(), &t1);
				*/
			}
		} else {
			pr.RewindTo(&t1)
		}

		return node, nil
	} else {
		pr.RewindTo(&t1)

		for {
			var tl sToken
			var t2 sToken
			pr.GetToken(&tl)
			pr.GetToken(&t2)
			pr.RewindTo(&tl)

			if tl.Type == tokens.ASttIdentifier && t2.Type == tokens.ASttColon {
				named := pr.CreateNode(ASsnNamedArgument)
				if named == nil {
					return nil, nil
				}

				node.AddChildLast(named, nil)

				named.AddChildLast(pr.ParseIdentifier())

				pr.GetToken(&t2)

				named.AddChildLast(pr.ParseAssignment())
			} else {
				node.AddChildLast(pr.ParseAssignment())
			}

			if pr.IsSyntaxError {
				return node, nil
			}

			pr.GetToken(&t1)
			if t1.Type == tokens.ASttListSeparator {
				continue
			} else {
				if withParenthesis {
					if t1.Type == tokens.ASttCloseParanthesis {
						node.UpdateSourcePos(t1.Position, t1.Length)
					} else {
						return node, errors.New("Expected " + tokens.GetDefinitionOrList([]tokens.Token{tokens.ASttCloseParanthesis}) + ", got " + tokens.GetDefinition(t1.Type))
					}
				} else {
					pr.RewindTo(&t1)
				}

				return node, nil
			}
		}
	}
}

func (pr *Parser) IsFunctionCall() bool {
	var s sToken
	var t1 sToken
	var t2 sToken

	pr.GetToken(&s)
	t1 = s

	if t1.Type == tokens.ASttScope {
		pr.GetToken(&t1)
	}
	pr.GetToken(&t2)

	for t1.Type == tokens.ASttIdentifier && t2.Type == tokens.ASttScope {
		pr.GetToken(&t1)
		pr.GetToken(&t2)
	}

	if t1.Type == tokens.ASttIdentifier || pr.IsDataType(&t1) {
		pr.RewindTo(&s)
		return false
	}

	if t2.Type == tokens.ASttOpenParanthesis {
		pr.RewindTo(&s)
		return true
	}

	pr.RewindTo(&s)
	return false
}

func (pr *Parser) ParseAssignment() (*ScriptNode, error) {
	node := pr.CreateNode(ASsnAssignment)
	if node == nil {
		return nil, nil
	}

	node.AddChildLast(pr.ParseCondition())
	if pr.IsSyntaxError {
		return node, nil
	}

	var t sToken
	pr.GetToken(&t)
	pr.RewindTo(&t)

	if pr.IsAssignOperator(t.Type) {
		node.AddChildLast(pr.ParseAssignOperator())
		if pr.IsSyntaxError {
			return node, nil
		}

		node.AddChildLast(pr.ParseAssignment())
		if pr.IsSyntaxError {
			return node, nil
		}
	}

	return node, nil
}

func (pr *Parser) ParseCondition() (*ScriptNode, error) {
	node := pr.CreateNode(ASsnCondition)
	if node == nil {
		return nil, nil
	}

	node.AddChildLast(pr.ParseExpression())
	if pr.IsSyntaxError {
		return nil, nil
	}

	var t sToken
	pr.GetToken(&t)
	if t.Type == tokens.ASttQuestion {
		node.AddChildLast(pr.ParseAssignment())
		if pr.IsSyntaxError {
			return node, nil
		}

		pr.GetToken(&t)

		if t.Type != tokens.ASttColon {
			//TODO: error expect :
			return node, nil
		}

		node.AddChildLast(pr.ParseAssignment())
		if pr.IsSyntaxError {
			return node, nil
		}
	} else {
		pr.RewindTo(&t)
	}

	return node, nil
}

func (pr *Parser) ParseExpression() (*ScriptNode, error) {
	node := pr.CreateNode(ASsnExpression)
	if node == nil {
		return nil, nil
	}

	node.AddChildLast(pr.ParseExprTerm())
	if pr.IsSyntaxError {
		return node, nil
	}

	for {
		var t sToken
		pr.GetToken(&t)
		pr.RewindTo(&t)

		if !pr.IsOperator(t.Type) {
			return node, nil
		}

		node.AddChildLast(pr.ParseExprOperator())
		if pr.IsSyntaxError {
			return node, nil
		}

		node.AddChildLast(pr.ParseExprTerm())
		if pr.IsSyntaxError {
			return node, nil
		}
	}
	return node, nil
}

func (pr *Parser) ParseExprTerm() (*ScriptNode, error) {
	node := pr.CreateNode(ASsnExprTerm)
	if node == nil {
		return nil, nil
	}

	var t sToken
	pr.GetToken(&t)
	var t2 sToken = t
	var t3 sToken

	if pr.IsDataType(&t2) && pr.CheckTemplateType(t2) {
		pr.GetToken(&t2)
		pr.GetToken(&t3)
		if t2.Type == tokens.ASttAssignment && t3.Type == tokens.ASttStartStatementBlock {
			pr.RewindTo(&t)
			node.AddChildLast(pr.ParseType(false, false, false))

			pr.GetToken(&t2)
			node.AddChildLast(pr.ParseInitList())
			return node, nil
		}
	}

	pr.RewindTo(&t)

	for {
		pr.GetToken(&t)
		pr.RewindTo(&t)

		if !pr.IsPreOperator(t.Type) {
			break
		}

		node.AddChildLast(pr.ParseExprPreOp())
		if pr.IsSyntaxError {
			return node, nil
		}
	}

	node.AddChildLast(pr.ParseExprValue())
	if pr.IsSyntaxError {
		return node, nil
	}

	for {
		pr.GetToken(&t)
		pr.RewindTo(&t)
		if !pr.IsPostOperator(t.Type) {
			return node, nil
		}

		node.AddChildLast(pr.ParseExprPostOp())
		if pr.IsSyntaxError {
			return node, nil
		}
	}
	return node, nil
}

func (pr *Parser) ParseExprPreOp() (*ScriptNode, error) {
	node := pr.CreateNode(ASsnExprPreOp)
	if node == nil {
		return nil, nil
	}

	var t sToken
	pr.GetToken(&t)
	if !pr.IsPreOperator(t.Type) {
		//TODO: error TXT_EXPECTED_PRE_OPERATOR
		return node, nil
	}

	node.SetToken(&t)
	node.UpdateSourcePos(t.Position, t.Length)

	return node, nil
}

func (pr *Parser) ParseExprPostOp() (*ScriptNode, error) {
	node := pr.CreateNode(ASsnExprPostOp)
	if node == nil {
		return nil, nil
	}

	var t sToken
	pr.GetToken(&t)
	if !pr.IsPostOperator(t.Type) {
		//TODO: error TXT_EXPECTED_POST_OPERATOR
		return node, nil
	}

	node.SetToken(&t)
	node.UpdateSourcePos(t.Position, t.Length)

	if t.Type == tokens.ASttDot {
		var t1 sToken
		var t2 sToken

		pr.GetToken(&t1)
		pr.GetToken(&t2)
		pr.RewindTo(&t1)
		if t2.Type == tokens.ASttOpenParanthesis {
			node.AddChildLast(pr.ParseFunctionCall())
		} else {
			node.AddChildLast(pr.ParseIdentifier())
		}
	} else if t.Type == tokens.ASttOpenBracket {
		node.AddChildLast(pr.ParseArgList(false))

		pr.GetToken(&t)
		if t.Type != tokens.ASttCloseBracket {
			//TODO: error expected ]
			return node, nil
		}

		node.UpdateSourcePos(t.Position, t.Length)
	} else if t.Type == tokens.ASttOpenParanthesis {
		pr.RewindTo(&t)
		node.AddChildLast(pr.ParseArgList(true))
	}

	return node, nil
}

func (pr *Parser) ParseExprOperator() (*ScriptNode, error) {
	node := pr.CreateNode(ASsnExprOperator)
	if node == nil {
		return node, nil
	}

	var t sToken
	pr.GetToken(&t)
	if !pr.IsOperator(t.Type) {
		//TODO: error TXT_EXPECTED_OPERATOR
		return node, nil
	}

	node.SetToken(&t)
	node.UpdateSourcePos(t.Position, t.Length)

	return node, nil
}

func (pr *Parser) ParseAssignOperator() (*ScriptNode, error) {
	node := pr.CreateNode(ASsnExprOperator)
	if node == nil {
		return nil, nil
	}

	var t sToken
	pr.GetToken(&t)

	if !pr.IsAssignOperator(t.Type) {
		//TODO: error TXT_EXPECTED_OPERATOR
		return node, nil
	}

	node.SetToken(&t)
	node.UpdateSourcePos(t.Position, t.Length)

	return node, nil
}

func (pr *Parser) IsOperator(tokenType tokens.Token) bool {
	if tokenType == tokens.ASttPlus ||
		tokenType == tokens.ASttMinus ||
		tokenType == tokens.ASttStar ||
		tokenType == tokens.ASttSlash ||
		tokenType == tokens.ASttPercent ||
		tokenType == tokens.ASttStarStar ||
		tokenType == tokens.ASttAnd ||
		tokenType == tokens.ASttOr ||
		tokenType == tokens.ASttXor ||
		tokenType == tokens.ASttEqual ||
		tokenType == tokens.ASttNotEqual ||
		tokenType == tokens.ASttLessThan ||
		tokenType == tokens.ASttLessThanOrEqual ||
		tokenType == tokens.ASttGreaterThan ||
		tokenType == tokens.ASttGreaterThanOrEqual ||
		tokenType == tokens.ASttAmp ||
		tokenType == tokens.ASttBitOr ||
		tokenType == tokens.ASttBitXor ||
		tokenType == tokens.ASttBitShiftLeft ||
		tokenType == tokens.ASttBitShiftRight ||
		tokenType == tokens.ASttBitShiftRightArith ||
		tokenType == tokens.ASttIs ||
		tokenType == tokens.ASttNotIs {
		return true
	}
	return false
}

func (pr *Parser) IsAssignOperator(tokenType tokens.Token) bool {
	if tokenType == tokens.ASttAssignment ||
		tokenType == tokens.ASttAddAssign ||
		tokenType == tokens.ASttSubAssign ||
		tokenType == tokens.ASttMulAssign ||
		tokenType == tokens.ASttDivAssign ||
		tokenType == tokens.ASttModAssign ||
		tokenType == tokens.ASttPowAssign ||
		tokenType == tokens.ASttAndAssign ||
		tokenType == tokens.ASttOrAssign ||
		tokenType == tokens.ASttXorAssign ||
		tokenType == tokens.ASttShiftLeftAssign ||
		tokenType == tokens.ASttShiftRightLAssign ||
		tokenType == tokens.ASttShiftRightAAssign {
		return true
	}
	return false
}

func (pr *Parser) IsPreOperator(tokenType tokens.Token) bool {
	if tokenType == tokens.ASttMinus ||
		tokenType == tokens.ASttPlus ||
		tokenType == tokens.ASttNot ||
		tokenType == tokens.ASttInc ||
		tokenType == tokens.ASttDec ||
		tokenType == tokens.ASttBitNot ||
		tokenType == tokens.ASttHandle {
		return true
	}
	return false
}

func (pr *Parser) IsPostOperator(tokenType tokens.Token) bool {
	if tokenType == tokens.ASttInc || // post increment
		tokenType == tokens.ASttDec || // post decrement
		tokenType == tokens.ASttDot || // member access
		tokenType == tokens.ASttOpenBracket || // index operator
		tokenType == tokens.ASttOpenParanthesis { // argument list for call on function pointer
		return true
	}
	return false
}

func (pr *Parser) IsConstant(tokenType tokens.Token) bool {
	if tokenType == tokens.ASttIntConstant ||
		tokenType == tokens.ASttFloatConstant ||
		tokenType == tokens.ASttDoubleConstant ||
		tokenType == tokens.ASttStringConstant ||
		tokenType == tokens.ASttMultilineStringConstant ||
		tokenType == tokens.ASttHeredocStringConstant ||
		tokenType == tokens.ASttTrue ||
		tokenType == tokens.ASttFalse ||
		tokenType == tokens.ASttBitsConstant ||
		tokenType == tokens.ASttNull {
		return true
	}
	return false
}

func (pr *Parser) ParseImport() (*ScriptNode, error) {
	node := pr.CreateNode(ASsnImport)
	if node == nil {
		return nil, nil
	}

	var t sToken
	pr.GetToken(&t)
	if t.Type != tokens.ASttImport {
		//TODO: error expected import
		return node, nil
	}

	node.SetToken(&t)
	node.UpdateSourcePos(t.Position, t.Length)

	node.AddChildLast(pr.ParseFunctionDefinition())
	if pr.IsSyntaxError {
		return node, nil
	}

	pr.GetToken(&t)
	if t.Type != tokens.ASttIdentifier {
		//TODO: error expected from
		return node, nil
	}

	pr.tempstr = pr.Script.Code[t.Position : t.Position+t.Length]
	if pr.tempstr != tokens.ASFromToken {
		//TODO: error expected from
		return node, nil
	}

	node.UpdateSourcePos(t.Position, t.Length)

	pr.GetToken(&t)
	if t.Type != tokens.ASttStringConstant {
		//TODO: error TXT_EXPECTED_STRING
		return node, nil
	}

	mod := pr.CreateNode(ASsnConstant)
	if mod == nil {
		return nil, nil
	}

	node.AddChildLast(mod, nil)

	mod.SetToken(&t)
	mod.UpdateSourcePos(t.Position, t.Length)

	pr.GetToken(&t)
	if t.Type != tokens.ASttEndStatement {
		//TODO: error expected end statement
		return node, nil
	}

	node.UpdateSourcePos(t.Position, t.Length)
	return node, nil
}

func (pr *Parser) ParseScript(inBlock bool) (*ScriptNode, error) {
	node := pr.CreateNode(ASsnScript)
	if node == nil {
		return nil, nil
	}

	var t1 sToken
	var t2 sToken

	//TODO/FIXME: Look this through, some of the if statements might be wrong.
	for {
		for !pr.IsSyntaxError {
			pr.GetToken(&t1)
			pr.GetToken(&t2)
			pr.RewindTo(&t1)

			if t1.Type == tokens.ASttImport {
				node.AddChildLast(pr.ParseImport())
			} else if t1.Type == tokens.ASttEnum || (pr.IdentifierIs(t1, tokens.ASSharedToken) && t2.Type == tokens.ASttEnum) {
				node.AddChildLast(pr.ParseEnumeration())
			} else if t1.Type == tokens.ASttTypedef {
				node.AddChildLast(pr.ParseTypedef())
			} else if t1.Type == tokens.ASttClass ||
				((pr.IdentifierIs(t1, tokens.ASSharedToken) || pr.IdentifierIs(t1, tokens.ASFinalToken) || pr.IdentifierIs(t1, tokens.ASAbstractToken)) && t2.Type == tokens.ASttClass) ||
				(pr.IdentifierIs(t1, tokens.ASSharedToken) && (pr.IdentifierIs(t1, tokens.ASFinalToken) || pr.IdentifierIs(t1, tokens.ASAbstractToken))) {
				node.AddChildLast(pr.ParseClass())
			} else if t1.Type == tokens.ASttMixin {
				node.AddChildLast(pr.ParseMixin())
			} else if t1.Type == tokens.ASttInterface || (t1.Type == tokens.ASttIdentifier && t2.Type == tokens.ASttInterface) {
				node.AddChildLast(pr.ParseInterface())
			} else if t1.Type == tokens.ASttFuncDef {
				node.AddChildLast(pr.ParseFuncDef())
			} else if t1.Type == tokens.ASttConst || t1.Type == tokens.ASttScope || t1.Type == tokens.ASttAuto || pr.IsDataType(&t1) {
				if pr.IsVirtualPropertyDecl() {
					node.AddChildLast(pr.ParseVirtualPropertyDecl(false, false))
				} else if pr.IsVarDecl() {
					node.AddChildLast(pr.ParseDeclaration(false, true))
				} else {
					node.AddChildLast(pr.ParseFunction(false))
				}
			} else if t1.Type == tokens.ASttEndStatement {
				pr.GetToken(&t1)
			} else if t1.Type == tokens.ASttNamespace {
				node.AddChildLast(pr.ParseNamespace())
			} else if t1.Type == tokens.ASttEnd {
				return node, nil
			} else if inBlock && t1.Type == tokens.ASttEndStatementBlock {
				return node, nil
			} else {
				/*
					asCString str;
					const char *t = asCTokenizer::GetDefinition(t1.type);
					if( t == 0 ) t = "<unknown token>";

					str.Format(TXT_UNEXPECTED_TOKEN_s, t);

					Error(str, &t1);
				*/
			}
		}

		if pr.IsSyntaxError {
			pr.GetToken(&t1)
			for t1.Type != tokens.ASttEndStatement && t1.Type != tokens.ASttEnd && t1.Type != tokens.ASttStartStatementBlock {
				pr.GetToken(&t1)
			}

			if t1.Type == tokens.ASttStartStatementBlock {
				level := 1
				for level > 0 {
					pr.GetToken(&t1)
					if t1.Type == tokens.ASttStartStatementBlock {
						level++
					}
					if t1.Type == tokens.ASttEndStatementBlock {
						level--
					}
					if t1.Type == tokens.ASttEnd {
						break
					}
				}
			}
			pr.IsSyntaxError = false
		}
	}
	return nil, nil
}

func (pr *Parser) ParseNamespace() (*ScriptNode, error) {
	node := pr.CreateNode(ASsnNamespace)
	if node == nil {
		return nil, nil
	}

	var t1 sToken

	pr.GetToken(&t1)
	if t1.Type == tokens.ASttNamespace {
		node.UpdateSourcePos(t1.Position, t1.Length)
	} else {
		//TODO: error expected namespace
		return node, nil
	}

	node.AddChildLast(pr.ParseIdentifier())
	if pr.IsSyntaxError {
		return node, nil
	}

	pr.GetToken(&t1)
	if t1.Type == tokens.ASttStartStatementBlock {
		node.UpdateSourcePos(t1.Position, t1.Length)
	} else {
		//TODO: error expected start statement block
		return node, nil
	}
	//start := t1

	node.AddChildLast(pr.ParseScript(true))
	if !pr.IsSyntaxError {
		pr.GetToken(&t1)
		if t1.Type == tokens.ASttEndStatementBlock {
			node.UpdateSourcePos(t1.Position, t1.Length)
		} else {
			if t1.Type == tokens.ASttEnd {
				//TODO: error TXT_UNEXPECTED_END_OF_FILE
			} else {
				//TODO: error expected end statement block
			}
			//TODO: info TXT_WHITE_PARSING_NAMESPACE
			return node, nil
		}
	}

	return node, nil
}

func (pr *Parser) ParseEnumeration() (*ScriptNode, error) {
	var ident *ScriptNode
	var dataType *ScriptNode
	node := pr.CreateNode(ASsnEnum)
	if node == nil {
		return nil, nil
	}

	var token sToken

	pr.GetToken(&token)
	if pr.IdentifierIs(token, tokens.ASSharedToken) {
		pr.RewindTo(&token)
		node.AddChildLast(pr.ParseIdentifier())
		if pr.IsSyntaxError {
			return node, nil
		}

		pr.GetToken(&token)
	}

	if token.Type != tokens.ASttEnd {
		//TODO: error expected token enum
		return node, nil
	}

	node.SetToken(&token)
	node.UpdateSourcePos(token.Position, token.Length)

	pr.GetToken(&token)
	if tokens.ASttIdentifier != token.Type {
		//TODO: error TXT_EXPECTED_IDENTIFIER
		return node, nil
	}

	dataType = pr.CreateNode(ASsnDataType)
	if dataType == nil {
		return nil, nil
	}

	node.AddChildLast(dataType, nil)

	ident = pr.CreateNode(ASsnIdentifier)
	if ident == nil {
		return nil, nil
	}

	ident.SetToken(&token)
	ident.UpdateSourcePos(token.Position, token.Length)
	dataType.AddChildLast(ident, nil)

	pr.GetToken(&token)
	if token.Type != tokens.ASttStartStatementBlock {
		pr.RewindTo(&token)
		//TODO: error expected token.type
		return node, nil
	}

	for tokens.ASttEnd != token.Type {
		pr.GetToken(&token)

		if tokens.ASttEndStatement == token.Type {
			pr.RewindTo(&token)
			break
		}

		if tokens.ASttIdentifier != token.Type {
			//TODO: error TXT_EXPECTED_IDENTIFIER
			return node, nil
		}

		ident = pr.CreateNode(ASsnIdentifier)
		if ident == nil {
			return nil, nil
		}

		ident.SetToken(&token)
		ident.UpdateSourcePos(token.Position, token.Length)
		node.AddChildLast(ident, nil)

		pr.GetToken(&token)

		if token.Type == tokens.ASttAssignment {
			pr.RewindTo(&token)
			tmp, err := pr.SuperficiallyParseVarInit()
			_ = err
			node.AddChildLast(tmp, nil)
			if pr.IsSyntaxError {
				return node, nil
			}
			pr.GetToken(&token)
		}

		if tokens.ASttListSeparator != token.Type {
			pr.RewindTo(&token)
			break
		}
	}

	pr.GetToken(&token)
	if token.Type != tokens.ASttEndStatementBlock {
		pr.RewindTo(&token)
		//TODO: error expected }
		return node, nil
	}

	return node, nil
}

func (pr *Parser) IsVarDecl() bool {
	var t sToken
	pr.GetToken(&t)
	pr.RewindTo(&t)

	var t1 sToken
	pr.GetToken(&t1)
	if t1.Type != tokens.ASttPrivate && t1.Type != tokens.ASttProtected {
		pr.RewindTo(&t1)
	}

	pr.GetToken(&t1)
	if t1.Type == tokens.ASttConst {
		pr.GetToken(&t1)
	}

	var t2 sToken
	if t1.Type != tokens.ASttAuto {
		if t1.Type == tokens.ASttScope {
			pr.GetToken(&t1)
		}

		pr.GetToken(&t2)
		for t1.Type == tokens.ASttIdentifier {
			if t2.Type == tokens.ASttScope {
				pr.GetToken(&t1)
				pr.GetToken(&t2)
				continue
			} else if t2.Type == tokens.ASttLessThan {
				pr.RewindTo(&t2)
				if pr.CheckTemplateType(t1) {
					var t3 sToken
					pr.GetToken(&t3)
					if t3.Type == tokens.ASttScope {
						pr.GetToken(&t1)
						pr.GetToken(&t2)
						continue
					}
				}
			}

			break
		}
		pr.RewindTo(&t2)
	}

	if !pr.IsRealType(t1.Type) && t1.Type != tokens.ASttIdentifier && t1.Type != tokens.ASttAuto {
		pr.RewindTo(&t)
		return false
	}

	if !pr.CheckTemplateType(t1) {
		pr.RewindTo(&t)
		return false
	}

	pr.GetToken(&t2)
	for t2.Type == tokens.ASttHandle || t2.Type == tokens.ASttAmp || t2.Type == tokens.ASttOpenBracket {
		if t2.Type == tokens.ASttOpenBracket {
			pr.GetToken(&t2)
			if t2.Type != tokens.ASttCloseBracket {
				pr.RewindTo(&t)
				return false
			}
		}
		pr.GetToken(&t2)
	}

	if t2.Type != tokens.ASttIdentifier {
		pr.RewindTo(&t)
		return false
	}

	pr.GetToken(&t2)
	if t2.Type == tokens.ASttEnd || t2.Type == tokens.ASttAssignment || t2.Type == tokens.ASttListSeparator {
		pr.RewindTo(&t)
		return true
	}
	if t2.Type == tokens.ASttOpenParanthesis {
		nest := 0
		for t2.Type != tokens.ASttEnd {
			if t2.Type == tokens.ASttOpenParanthesis {
				nest++
			} else if t2.Type == tokens.ASttCloseParanthesis {
				nest--
				if nest <= 0 {
					break
				}
			}
			pr.GetToken(&t2)
		}

		if t2.Type == tokens.ASttEnd {
			return false
		} else {
			pr.GetToken(&t1)
			pr.RewindTo(&t)
			if t1.Type == tokens.ASttStartStatementBlock || t1.Type == tokens.ASttEnd {
				return false
			}
		}

		pr.RewindTo(&t)
		return true
	}

	pr.RewindTo(&t)
	return true
}

func (pr *Parser) IsVirtualPropertyDecl() bool {
	var t sToken
	pr.GetToken(&t)
	pr.RewindTo(&t)

	var t1 sToken
	pr.GetToken(&t1)
	if t1.Type != tokens.ASttPrivate && t1.Type != tokens.ASttProtected {
		pr.RewindTo(&t1)
	}

	pr.GetToken(&t1)
	if t1.Type == tokens.ASttConst {
		pr.GetToken(&t1)
	}

	if t1.Type == tokens.ASttScope {
		pr.GetToken(&t1)
	}

	if t1.Type == tokens.ASttIdentifier {
		var t2 sToken
		pr.GetToken(&t2)
		for t1.Type == tokens.ASttIdentifier && t2.Type == tokens.ASttScope {
			pr.GetToken(&t1)
			pr.GetToken(&t2)
		}

		pr.RewindTo(&t2)
	} else if !pr.IsRealType(t1.Type) {
		pr.RewindTo(&t)
		return false
	}

	if !pr.CheckTemplateType(t1) {
		pr.RewindTo(&t)
		return false
	}

	var t2 sToken
	pr.GetToken(&t2)
	for t2.Type == tokens.ASttHandle || t2.Type == tokens.ASttOpenBracket {
		if t2.Type == tokens.ASttOpenBracket {
			pr.GetToken(&t2)
			if t2.Type != tokens.ASttCloseBracket {
				pr.RewindTo(&t)
				return false
			}
		}

		pr.GetToken(&t2)
	}

	if t2.Type != tokens.ASttIdentifier {
		pr.RewindTo(&t)
		return false
	}

	pr.GetToken(&t2)
	if t2.Type == tokens.ASttStartStatementBlock {
		pr.RewindTo(&t)
		return true
	}

	pr.RewindTo(&t)
	return false
}

func (pr *Parser) IsFuncDecl(isMethod bool) bool {
	var t sToken
	pr.GetToken(&t)
	pr.RewindTo(&t)

	if isMethod {
		var t1 sToken
		var t2 sToken
		pr.GetToken(&t1)
		if t1.Type != tokens.ASttPrivate && t1.Type != tokens.ASttProtected {
			pr.RewindTo(&t1)
		}

		pr.GetToken(&t1)
		pr.GetToken(&t2)
		pr.RewindTo(&t1)

		if (t1.Type == tokens.ASttIdentifier && t2.Type == tokens.ASttOpenParanthesis) || t1.Type == tokens.ASttBitNot {
			pr.RewindTo(&t)
			return true
		}
	}

	var t1 sToken
	pr.GetToken(&t1)
	if t1.Type == tokens.ASttConst {
		pr.GetToken(&t1)
	}

	if t1.Type == tokens.ASttScope {
		pr.GetToken(&t1)
	}
	for t1.Type == tokens.ASttIdentifier {
		var t2 sToken
		pr.GetToken(&t2)
		if t2.Type == tokens.ASttScope {
			pr.GetToken(&t1)
		} else {
			pr.RewindTo(&t2)
			break
		}
	}

	if !pr.IsDataType(&t1) {
		pr.RewindTo(&t)
		return false
	}

	if !pr.CheckTemplateType(t1) {
		pr.RewindTo(&t)
		return false
	}

	var t2 sToken
	pr.GetToken(&t2)
	for t2.Type == tokens.ASttHandle || t2.Type == tokens.ASttOpenBracket {
		if t2.Type == tokens.ASttOpenBracket {
			pr.GetToken(&t2)
			if t2.Type != tokens.ASttCloseBracket {
				pr.RewindTo(&t)
				return false
			}
		}

		pr.GetToken(&t2)
	}

	if t2.Type == tokens.ASttAmp {
		pr.RewindTo(&t)
		return false
	}

	if t2.Type != tokens.ASttIdentifier {
		pr.RewindTo(&t)
		return false
	}

	pr.GetToken(&t2)
	if t2.Type == tokens.ASttOpenParanthesis {
		nest := 0
		pr.GetToken(&t2)
		for (nest >= 1 || t2.Type != tokens.ASttCloseParanthesis) && t2.Type != tokens.ASttEnd {
			if t2.Type == tokens.ASttOpenParanthesis {
				nest++
			}
			if t2.Type == tokens.ASttCloseParanthesis {
				nest--
			}

			pr.GetToken(&t2)
		}

		if t2.Type == tokens.ASttEnd {
			return false
		} else {
			if isMethod {

				pr.GetToken(&t1)
				if t1.Type != tokens.ASttConst {
					pr.RewindTo(&t1)
				}

				for {
					pr.GetToken(&t2)
					if !pr.IdentifierIs(t2, tokens.ASFinalToken) && !pr.IdentifierIs(t2, tokens.ASOverrideToken) {
						pr.RewindTo(&t2)
						break
					}
				}
			}

			pr.GetToken(&t1)
			pr.RewindTo(&t)
			if t1.Type == tokens.ASttStartStatementBlock {
				return true
			}
		}

		pr.RewindTo(&t)
		return false
	}

	pr.RewindTo(&t)
	return false
}

func (pr *Parser) ParseFuncDef() (*ScriptNode, error) {
	node := pr.CreateNode(ASsnFuncDef)
	if node == nil {
		return nil, nil
	}

	var t1 sToken
	pr.GetToken(&t1)
	if t1.Type != tokens.ASttFuncDef {
		//TODO: error ttfuncdef???
		return node, nil
	}

	node.SetToken(&t1)

	node.AddChildLast(pr.ParseType(true, false, false))
	if pr.IsSyntaxError {
		return node, nil
	}

	node.AddChildLast(pr.ParseTypeMod(false))
	if pr.IsSyntaxError {
		return node, nil
	}

	node.AddChildLast(pr.ParseIdentifier())
	if pr.IsSyntaxError {
		return node, nil
	}

	node.AddChildLast(pr.ParseParameterList())
	if pr.IsSyntaxError {
		return node, nil
	}

	pr.GetToken(&t1)
	if t1.Type != tokens.ASttEndStatement {
		//TODO: error expected end statement
		return node, nil
	}

	node.UpdateSourcePos(t1.Position, t1.Length)
	return node, nil
}

func (pr *Parser) ParseFunction(isMethod bool) (*ScriptNode, error) {
	node := pr.CreateNode(ASsnFunction)
	if node == nil {
		return nil, nil
	}

	var t1 sToken
	var t2 sToken
	pr.GetToken(&t1)
	pr.GetToken(&t2)
	pr.RewindTo(&t1)

	fmt.Println(tokens.GetDefinition(t1.Type), tokens.GetDefinition(t2.Type))

	if isMethod && t1.Type == tokens.ASttPrivate {
		node.AddChildLast(pr.ParseToken(tokens.ASttPrivate))
	} else if isMethod && t1.Type == tokens.ASttProtected {
		node.AddChildLast(pr.ParseToken(tokens.ASttProtected))
	}
	if pr.IsSyntaxError {
		return node, nil
	}

	if !isMethod && pr.IdentifierIs(t1, tokens.ASSharedToken) {
		node.AddChildLast(pr.ParseIdentifier())
		if pr.IsSyntaxError {
			return node, nil
		}
	}

	if !isMethod || (t1.Type != tokens.ASttBitNot && t2.Type != tokens.ASttOpenParanthesis) {
		node.AddChildLast(pr.ParseType(true, false, false))
		if pr.IsSyntaxError {
			return node, nil
		}

		node.AddChildLast(pr.ParseTypeMod(false))
		if pr.IsSyntaxError {
			return node, nil
		}
	}

	if isMethod && t1.Type == tokens.ASttBitNot {
		node.AddChildLast(pr.ParseToken(tokens.ASttBitNot))
		if pr.IsSyntaxError {
			return node, nil
		}
	}

	node.AddChildLast(pr.ParseIdentifier())
	if pr.IsSyntaxError {
		return node, nil
	}

	node.AddChildLast(pr.ParseParameterList())
	if pr.IsSyntaxError {
		return node, nil
	}

	if isMethod {
		pr.GetToken(&t1)
		pr.RewindTo(&t1)

		if t1.Type == tokens.ASttConst {
			node.AddChildLast(pr.ParseToken(tokens.ASttConst))
		}

		pr.ParseMethodOverrideBehaviors(node)
		if pr.IsSyntaxError {
			return node, nil
		}
	}

	node.AddChildLast(pr.SuperficiallyParseStatementBlock())

	return node, nil
}

func (pr *Parser) ParseInterfaceMethod() (*ScriptNode, error) {
	node := pr.CreateNode(ASsnFunction)
	if node == nil {
		return nil, nil
	}

	node.AddChildLast(pr.ParseType(true, false, false))
	if pr.IsSyntaxError {
		return node, nil
	}

	node.AddChildLast(pr.ParseTypeMod(false))
	if pr.IsSyntaxError {
		return node, nil
	}

	node.AddChildLast(pr.ParseIdentifier())
	if pr.IsSyntaxError {
		return node, nil
	}

	node.AddChildLast(pr.ParseParameterList())
	if pr.IsSyntaxError {
		return node, nil
	}

	var t1 sToken
	pr.GetToken(&t1)
	pr.RewindTo(&t1)
	if t1.Type == tokens.ASttConst {
		node.AddChildLast(pr.ParseToken(tokens.ASttConst))
	}

	pr.GetToken(&t1)
	if t1.Type != tokens.ASttEndStatement {
		errors.New("Expected " + tokens.GetDefinitionOrList([]tokens.Token{tokens.ASttEndStatement}) + ", got " + tokens.GetDefinition(t1.Type))
		return node, nil
	}

	node.UpdateSourcePos(t1.Position, t1.Length)
	return node, nil
}

func (pr *Parser) ParseVirtualPropertyDecl(isMethod, isInterface bool) (*ScriptNode, error) {
	node := pr.CreateNode(ASsnVirtualProperty)
	if node == nil {
		return nil, nil
	}

	var t1 sToken
	var t2 sToken
	pr.GetToken(&t1)
	pr.GetToken(&t2)
	pr.RewindTo(&t1)

	if isMethod && t1.Type == tokens.ASttPrivate {
		node.AddChildLast(pr.ParseToken(tokens.ASttPrivate))
	} else if isMethod && t1.Type == tokens.ASttProtected {
		node.AddChildLast(pr.ParseToken(tokens.ASttProtected))
	}
	if pr.IsSyntaxError {
		return node, nil
	}

	node.AddChildLast(pr.ParseType(true, false, false))
	if pr.IsSyntaxError {
		return node, nil
	}

	node.AddChildLast(pr.ParseTypeMod(false))
	if pr.IsSyntaxError {
		return node, nil
	}

	node.AddChildLast(pr.ParseIdentifier())
	if pr.IsSyntaxError {
		return node, nil
	}

	pr.GetToken(&t1)
	if t1.Type != tokens.ASttStartStatementBlock {
		return nil, errors.New("Expected " + tokens.GetDefinitionOrList([]tokens.Token{tokens.ASttStartStatementBlock}) + ", got " + tokens.GetDefinition(t1.Type))
		return node, nil
	}

	for {
		pr.GetToken(&t1)
		var aNode *ScriptNode

		if pr.IdentifierIs(t1, tokens.ASGetToken) || pr.IdentifierIs(t1, tokens.ASSetToken) {
			aNode = pr.CreateNode(ASsnVirtualProperty)
			if aNode == nil {
				return nil, nil
			}

			node.AddChildLast(aNode, nil)

			pr.RewindTo(&t1)
			aNode.AddChildLast(pr.ParseIdentifier())

			if isMethod {
				pr.GetToken(&t1)
				pr.RewindTo(&t1)
				if t1.Type == tokens.ASttConst {
					aNode.AddChildLast(pr.ParseToken(tokens.ASttConst))
				}

				if !isInterface {
					pr.ParseMethodOverrideBehaviors(aNode)
					if pr.IsSyntaxError {
						return node, nil
					}
				}
			}

			if !isInterface {
				pr.GetToken(&t1)
				if t1.Type == tokens.ASttStartStatementBlock {
					pr.RewindTo(&t1)
					aNode.AddChildLast(pr.SuperficiallyParseStatementBlock())
					if pr.IsSyntaxError {
						return node, nil
					}
				} else if t1.Type != tokens.ASttEndStatement {
					return node, errors.New("Expected " + tokens.GetDefinitionOrList([]tokens.Token{tokens.ASttEndStatement, tokens.ASttOpenBracket}) + ", got " + tokens.GetDefinition(t1.Type))
				}
			} else {
				pr.GetToken(&t1)
				if t1.Type != tokens.ASttEndStatement {
					return node, errors.New("Expected " + tokens.GetDefinitionOrList([]tokens.Token{tokens.ASttEndStatement}) + ", got " + tokens.GetDefinition(t1.Type))
				}
			}
		} else if t1.Type == tokens.ASttEndStatementBlock {
			break
		} else {
			//TODO: error expected get, set or END_STATEMENT_BLOCK
			return node, nil
		}
	}

	return node, nil
}

func (pr *Parser) ParseInterface() (*ScriptNode, error) {
	node := pr.CreateNode(ASsnInterface)
	if node == nil {
		return nil, nil
	}

	var t sToken
	pr.GetToken(&t)

	if t.Type == tokens.ASttIdentifier {
		pr.tempstr = pr.Script.Code[t.Position : t.Position+t.Length]
		if pr.tempstr != tokens.ASSharedToken {
			//TODO: error expected shared_token
			return node, nil
		}

		pr.RewindTo(&t)
		node.AddChildLast(pr.ParseIdentifier())
		pr.GetToken(&t)
	}
	if t.Type != tokens.ASttInterface {
		//TODO: error expected interface
		return node, nil
	}

	node.SetToken(&t)
	node.AddChildLast(pr.ParseIdentifier())

	pr.GetToken(&t)
	if t.Type == tokens.ASttColon {
		inherit := pr.CreateNode(ASsnIdentifier)
		node.AddChildLast(inherit, nil)

		pr.ParseOptionalScope(inherit)
		inherit.AddChildLast(pr.ParseIdentifier())
		pr.GetToken(&t)
		for t.Type == tokens.ASttListSeparator {
			inherit = pr.CreateNode(ASsnIdentifier)
			node.AddChildLast(inherit, nil)

			pr.ParseOptionalScope(inherit)
			inherit.AddChildLast(pr.ParseIdentifier())
			pr.GetToken(&t)
		}
	}

	if t.Type != tokens.ASttStartStatementBlock {
		return nil, errors.New("Expected " + tokens.GetDefinitionOrList([]tokens.Token{tokens.ASttStartStatementBlock}) + ", got " + tokens.GetDefinition(t.Type))
		return node, nil
	}

	pr.GetToken(&t)
	pr.RewindTo(&t)
	for t.Type != tokens.ASttEndStatementBlock && t.Type != tokens.ASttEnd {
		if pr.IsVirtualPropertyDecl() {
			node.AddChildLast(pr.ParseVirtualPropertyDecl(true, true))
		} else if t.Type == tokens.ASttEndStatement {
			pr.GetToken(&t)
		} else {
			node.AddChildLast(pr.ParseInterfaceMethod())
		}

		if pr.IsSyntaxError {
			return node, nil
		}

		pr.GetToken(&t)
		pr.RewindTo(&t)
	}

	pr.GetToken(&t)
	if t.Type != tokens.ASttEndStatementBlock {
		//TODO: error expected }
		return node, nil
	}

	node.UpdateSourcePos(t.Position, t.Length)
	return node, nil
}

func (pr *Parser) ParseMixin() (*ScriptNode, error) {
	node := pr.CreateNode(ASsnMixin)
	if node == nil {
		return nil, nil
	}

	var t sToken
	pr.GetToken(&t)

	if t.Type != tokens.ASttMixin {
		//TODO: error expected mixin
		return node, nil
	}

	node.SetToken(&t)

	node.AddChildLast(pr.ParseClass())
	return node, nil
}

func (pr *Parser) ParseClass() (*ScriptNode, error) {
	node := pr.CreateNode(ASsnClass)
	if node == nil {
		return nil, nil
	}

	var t sToken
	pr.GetToken(&t)

	for pr.IdentifierIs(t, tokens.ASSharedToken) ||
		pr.IdentifierIs(t, tokens.ASAbstractToken) ||
		pr.IdentifierIs(t, tokens.ASFinalToken) {
		pr.RewindTo(&t)
		node.AddChildLast(pr.ParseIdentifier())
		pr.GetToken(&t)
	}

	if t.Type != tokens.ASttClass {
		return node, errors.New("Expected " + tokens.GetDefinitionOrList([]tokens.Token{tokens.ASttClass}) + ", got " + tokens.GetDefinition(t.Type))
	}

	node.SetToken(&t)

	/*
		TODO: Add implicit handle class?
		GetToken(&t);

		if ( t.type == ttHandle )
			node->SetToken(&t);
		else
			RewindTo(&t);
	*/

	node.AddChildLast(pr.ParseIdentifier())

	pr.GetToken(&t)

	if t.Type == tokens.ASttColon {
		inherit := pr.CreateNode(ASsnIdentifier)
		node.AddChildLast(inherit, nil)

		pr.ParseOptionalScope(inherit)
		inherit.AddChildLast(pr.ParseIdentifier())
		pr.GetToken(&t)
		for t.Type == tokens.ASttListSeparator {
			inherit = pr.CreateNode(ASsnIdentifier)
			node.AddChildLast(inherit, nil)

			pr.ParseOptionalScope(inherit)
			inherit.AddChildLast(pr.ParseIdentifier())
			pr.GetToken(&t)
		}
	}

	if t.Type != tokens.ASttStartStatementBlock {
		return nil, errors.New("Expected " + tokens.GetDefinitionOrList([]tokens.Token{tokens.ASttStartStatementBlock}) + ", got " + tokens.GetDefinition(t.Type))
	}

	pr.GetToken(&t)
	pr.RewindTo(&t)
	for t.Type != tokens.ASttEndStatementBlock && t.Type != tokens.ASttEnd {
		if t.Type == tokens.ASttFuncDef {
			node.AddChildLast(pr.ParseFuncDef())
		} else if pr.IsFuncDecl(true) {
			node.AddChildLast(pr.ParseFunction(true))
		} else if pr.IsVirtualPropertyDecl() {
			node.AddChildLast(pr.ParseVirtualPropertyDecl(true, false))
		} else if pr.IsVarDecl() {
			node.AddChildLast(pr.ParseDeclaration(true, false))
		} else if t.Type == tokens.ASttEndStatement {
			pr.GetToken(&t)
		} else {
			//TODO: error TXT_EXPECTED_METHOD_OR_PROPERTY
			return node, errors.New("Expected method or property, got " + tokens.GetDefinition(t.Type))
		}

		if pr.IsSyntaxError {
			return node, nil
		}

		pr.GetToken(&t)
		pr.RewindTo(&t)
	}

	pr.GetToken(&t)
	if t.Type != tokens.ASttEndStatementBlock {
		//TODO: error expected }
		return node, nil
	}

	node.UpdateSourcePos(t.Position, t.Length)
	return node, nil
}

func (pr *Parser) SuperficiallyParseVarInit() (*ScriptNode, error) {
	node := pr.CreateNode(ASsnAssignment)
	if node == nil {
		return nil, nil
	}

	var t sToken
	pr.GetToken(&t)
	node.UpdateSourcePos(t.Position, t.Length)

	if t.Type == tokens.ASttAssignment {
		pr.GetToken(&t)
		//start := t

		indentParan := 0
		indentBrace := 0

		for indentParan >= 1 || indentBrace >= 1 || (t.Type != tokens.ASttListSeparator && t.Type != tokens.ASttEndStatement && t.Type != tokens.ASttEndStatementBlock) {
			if t.Type == tokens.ASttOpenParanthesis {
				indentParan++
			} else if t.Type == tokens.ASttCloseParanthesis {
				indentParan--
			} else if t.Type == tokens.ASttStartStatementBlock {
				indentBrace++
			} else if t.Type == tokens.ASttEndStatementBlock {
				indentBrace--
			} else if t.Type == tokens.ASttNonTerminatedStringConstant {
				//TODO: error TXT_NONTERMINATED_STRING
				break
			} else if t.Type == tokens.ASttEnd {
				//TODO: error TXT_UNEXPECTED_END_OF_FILE
				break
			}
			pr.GetToken(&t)
		}
		pr.RewindTo(&t)
	} else if t.Type == tokens.ASttOpenParanthesis {
		//start := t

		indent := 1
		for indent >= 1 {
			pr.GetToken(&t)
			if t.Type == tokens.ASttOpenParanthesis {
				indent++
			} else if t.Type == tokens.ASttCloseParanthesis {
				indent--
			} else if t.Type == tokens.ASttNonTerminatedStringConstant {
				//TODO: error TXT_NONTERMINATED_STRING
				break
			} else if t.Type == tokens.ASttEnd {
				//TODO: error TXT_UNEXPECTED_END_OF_FILE
				break
			}
		}
	} else {
		//TODO: error expected assignment or open paranthesis
	}
	return node, nil
}

func (pr *Parser) SuperficiallyParseStatementBlock() (*ScriptNode, error) {
	node := pr.CreateNode(ASsnStatementBlock)
	if node == nil {
		return nil, nil
	}

	var t sToken

	pr.GetToken(&t)
	if t.Type != tokens.ASttStartStatementBlock {
		return nil, errors.New("Expected " + tokens.GetDefinitionOrList([]tokens.Token{tokens.ASttStartStatementBlock}) + ", got " + tokens.GetDefinition(t.Type))
		return node, nil
	}

	node.UpdateSourcePos(t.Position, t.Length)

	//start := t1
	level := 1
	for level > 0 && !pr.IsSyntaxError {
		pr.GetToken(&t)
		if t.Type == tokens.ASttEndStatementBlock {
			level--
		} else if t.Type == tokens.ASttStartStatementBlock {
			level++
		} else if t.Type == tokens.ASttNonTerminatedStringConstant {
			//TODO: error TXT_NONTERMINATED_STRING
			break
		} else if t.Type == tokens.ASttEnd {
			//TODO: error TXT_UNEXPECTED_END_OF_FILE
			break
		}
	}

	node.UpdateSourcePos(t.Position, t.Length)
	return node, nil
}

func (pr *Parser) ParseStatementBlock() (*ScriptNode, error) {
	node := pr.CreateNode(ASsnStatementBlock)
	if node == nil {
		return nil, nil
	}

	var t1 sToken

	pr.GetToken(&t1)
	if t1.Type != tokens.ASttStartStatementBlock {
		return nil, errors.New("Expected " + tokens.GetDefinitionOrList([]tokens.Token{tokens.ASttStartStatementBlock}) + ", got " + tokens.GetDefinition(t1.Type))
		return node, nil
	}

	//start := t1
	for {
		for !pr.IsSyntaxError {
			pr.GetToken(&t1)
			if t1.Type == tokens.ASttEndStatementBlock {
				node.UpdateSourcePos(t1.Position, t1.Length)

				return node, nil
			} else {
				pr.RewindTo(&t1)

				if pr.IsVarDecl() {
					node.AddChildLast(pr.ParseDeclaration(false, false))
				} else {
					node.AddChildLast(pr.ParseStatement())
				}
			}
		}

		if pr.IsSyntaxError {
			pr.GetToken(&t1)
			for t1.Type != tokens.ASttEndStatement && t1.Type != tokens.ASttEnd &&
				t1.Type != tokens.ASttStartStatementBlock && t1.Type != tokens.ASttEndStatementBlock {
				pr.GetToken(&t1)
			}

			if t1.Type == tokens.ASttStartStatementBlock {
				level := 1
				for level > 0 {
					if t1.Type == tokens.ASttStartStatementBlock {
						level++
					}
					if t1.Type == tokens.ASttEndStatementBlock {
						level--
					}
					if t1.Type == tokens.ASttEnd {
						break
					}
				}
			} else if t1.Type == tokens.ASttEndStatementBlock {
				pr.RewindTo(&t1)
			} else if t1.Type == tokens.ASttEnd {
				//TODO: error TXT_UNEXPECTED_END_OF_FILE
				return node, nil
			}
			pr.IsSyntaxError = false
		}
	}
	return node, nil
}

func (pr *Parser) ParseInitList() (*ScriptNode, error) {
	node := pr.CreateNode(ASsnInitList)
	if node == nil {
		return nil, nil
	}

	var t1 sToken
	pr.GetToken(&t1)
	if t1.Type != tokens.ASttStartStatementBlock {
		return nil, errors.New("Expected " + tokens.GetDefinitionOrList([]tokens.Token{tokens.ASttStartStatementBlock}) + ", got " + tokens.GetDefinition(t1.Type))
		return node, nil
	}

	node.UpdateSourcePos(t1.Position, t1.Length)

	pr.GetToken(&t1)
	if t1.Type == tokens.ASttEndStatementBlock {
		node.UpdateSourcePos(t1.Position, t1.Length)

		return node, nil
	} else {
		pr.RewindTo(&t1)
		for {
			pr.GetToken(&t1)

			if t1.Type == tokens.ASttListSeparator {
				node.AddChildLast(pr.CreateNode(ASsnUndefined), nil)
				node.LastChild.UpdateSourcePos(t1.Position, 1)

				pr.GetToken(&t1)
				if t1.Type == tokens.ASttEndStatementBlock {
					node.AddChildLast(pr.CreateNode(ASsnUndefined), nil)
					node.LastChild.UpdateSourcePos(t1.Position, 1)
					node.UpdateSourcePos(t1.Position, t1.Length)

					return node, nil
				}
				pr.RewindTo(&t1)
			} else if t1.Type == tokens.ASttEndStatementBlock {
				node.AddChildLast(pr.CreateNode(ASsnUndefined), nil)
				node.LastChild.UpdateSourcePos(t1.Position, 1)
				node.UpdateSourcePos(t1.Position, t1.Length)

				return node, nil
			} else if t1.Type == tokens.ASttStartStatementBlock {
				pr.RewindTo(&t1)
				node.AddChildLast(pr.ParseInitList())
				if pr.IsSyntaxError {
					return node, nil
				}

				pr.GetToken(&t1)
				if t1.Type == tokens.ASttListSeparator {
					continue
				} else if t1.Type == tokens.ASttEndStatementBlock {
					node.UpdateSourcePos(t1.Position, t1.Length)

					return node, nil
				} else {
					//TODO: error expected } or ,
					return node, nil
				}
			} else {
				pr.RewindTo(&t1)
				node.AddChildLast(pr.ParseAssignment())
				if pr.IsSyntaxError {
					return node, nil
				}

				pr.GetToken(&t1)
				if t1.Type == tokens.ASttEndStatementBlock {
					node.UpdateSourcePos(t1.Position, t1.Length)

					return node, nil
				} else {
					//TODO: error expected } or ,
				}
			}
		}
	}
	return node, nil
}

func (pr *Parser) ParseDeclaration(isClassProp bool, isGlobal bool) (*ScriptNode, error) {
	node := pr.CreateNode(ASsnDeclaration)
	if node == nil {
		return nil, nil
	}

	var t sToken
	pr.GetToken(&t)
	pr.RewindTo(&t)

	if t.Type == tokens.ASttPrivate && isClassProp {
		node.AddChildLast(pr.ParseToken(tokens.ASttPrivate))
	} else if t.Type == tokens.ASttProtected && isClassProp {
		node.AddChildLast(pr.ParseToken(tokens.ASttProtected))
	}

	node.AddChildLast(pr.ParseType(true, false, !isClassProp))
	if pr.IsSyntaxError {
		return node, nil
	}

	for {
		node.AddChildLast(pr.ParseIdentifier())
		if pr.IsSyntaxError {
			return node, nil
		}

		if isClassProp || isGlobal {
			pr.GetToken(&t)
			pr.RewindTo(&t)
			if t.Type == tokens.ASttAssignment || t.Type == tokens.ASttOpenParanthesis {
				node.AddChildLast(pr.SuperficiallyParseVarInit())
				if pr.IsSyntaxError {
					return node, nil
				}
			}
		} else {
			pr.GetToken(&t)
			if t.Type == tokens.ASttOpenParanthesis {
				pr.RewindTo(&t)
				node.AddChildLast(pr.ParseArgList(true))
				if pr.IsSyntaxError {
					return node, nil
				}
			} else if t.Type == tokens.ASttAssignment {
				pr.GetToken(&t)
				pr.RewindTo(&t)
				if t.Type == tokens.ASttStartStatementBlock {
					node.AddChildLast(pr.ParseInitList())
					if pr.IsSyntaxError {
						return node, nil
					}
				} else {
					node.AddChildLast(pr.ParseAssignment())
					if pr.IsSyntaxError {
						return node, nil
					}
				}
			} else {
				pr.RewindTo(&t)
			}
		}

		pr.GetToken(&t)
		if t.Type == tokens.ASttListSeparator {
			continue
		} else if t.Type == tokens.ASttEndStatement {
			node.UpdateSourcePos(t.Position, t.Length)

			return node, nil
		} else {
			//TODO: error expected , or ;
			return node, errors.New("Expected " + tokens.GetDefinitionOrList([]tokens.Token{tokens.ASttListSeparator, tokens.ASttEndStatement}) + " got " + tokens.GetDefinition(t.Type))
		}
	}
	return node, nil
}

func (pr *Parser) ParseStatement() (*ScriptNode, error) {
	var t1 sToken
	pr.GetToken(&t1)
	pr.RewindTo(&t1)

	if t1.Type == tokens.ASttIf {
		return pr.ParseIf()
	} else if t1.Type == tokens.ASttFor {
		return pr.ParseFor()
	} else if t1.Type == tokens.ASttWhile {
		return pr.ParseWhile()
	} else if t1.Type == tokens.ASttReturn {
		return pr.ParseReturn()
	} else if t1.Type == tokens.ASttStartStatementBlock {
		return pr.ParseStatementBlock()
	} else if t1.Type == tokens.ASttBreak {
		return pr.ParseBreak()
	} else if t1.Type == tokens.ASttContinue {
		return pr.ParseContinue()
	} else if t1.Type == tokens.ASttDo {
		return pr.ParseDoWhile()
	} else if t1.Type == tokens.ASttSwitch {
		return pr.ParseSwitch()
	} else {
		if pr.IsVarDecl() {
			//TODO: error TXT_UNEXPECTED_VAR_DECL
		}
		return pr.ParseExpressionStatement()
	}
}

func (pr *Parser) ParseExpressionStatement() (*ScriptNode, error) {
	node := pr.CreateNode(ASsnExpressionStatement)
	if node == nil {
		return nil, nil
	}

	var t sToken
	pr.GetToken(&t)
	if t.Type == tokens.ASttEndStatement {
		node.UpdateSourcePos(t.Position, t.Length)
		return node, nil
	}

	pr.RewindTo(&t)
	node.AddChildLast(pr.ParseAssignment())
	if pr.IsSyntaxError {
		return node, nil
	}

	pr.GetToken(&t)
	if t.Type != tokens.ASttEndStatement {
		return node, errors.New("Expected " + tokens.GetDefinitionOrList([]tokens.Token{tokens.ASttEndStatement}) + ", got " + tokens.GetDefinition(t.Type))
	}

	node.UpdateSourcePos(t.Position, t.Length)
	return node, nil
}

func (pr *Parser) ParseSwitch() (*ScriptNode, error) {
	node := pr.CreateNode(ASsnSwitch)
	if node == nil {
		return nil, nil
	}

	var t sToken
	pr.GetToken(&t)
	if t.Type != tokens.ASttSwitch {
		//TODO: error expected switch
		return node, nil
	}

	node.UpdateSourcePos(t.Position, t.Length)

	pr.GetToken(&t)
	if t.Type != tokens.ASttOpenParanthesis {
		errors.New("Expected " + tokens.GetDefinitionOrList([]tokens.Token{tokens.ASttOpenParanthesis}) + ", got " + tokens.GetDefinition(t.Type))
		return node, nil
	}

	node.AddChildLast(pr.ParseAssignment())
	if pr.IsSyntaxError {
		return node, nil
	}

	pr.GetToken(&t)
	if t.Type != tokens.ASttCloseParanthesis {
		errors.New("Expected " + tokens.GetDefinitionOrList([]tokens.Token{tokens.ASttCloseParanthesis}) + ", got " + tokens.GetDefinition(t.Type))
		return node, nil
	}

	pr.GetToken(&t)
	if t.Type != tokens.ASttStartStatementBlock {
		return nil, errors.New("Expected " + tokens.GetDefinitionOrList([]tokens.Token{tokens.ASttStartStatementBlock}) + ", got " + tokens.GetDefinition(t.Type))
		return node, nil
	}

	for !pr.IsSyntaxError {
		pr.GetToken(&t)
		if t.Type == tokens.ASttEndStatementBlock {
			break
		}

		pr.RewindTo(&t)

		if t.Type != tokens.ASttCase && t.Type != tokens.ASttDefault {
			return nil, errors.New("Expected " + tokens.GetDefinitionOrList([]tokens.Token{tokens.ASttCase, tokens.ASttDefault}) + ", got " + tokens.GetDefinition(t.Type))
			return node, nil
		}

		node.AddChildLast(pr.ParseCase())
		if pr.IsSyntaxError {
			return node, nil
		}
	}

	if t.Type != tokens.ASttEndStatementBlock {
		//TODO: error expected }
		return node, nil
	}
	return node, nil
}

func (pr *Parser) ParseCase() (*ScriptNode, error) {
	node := pr.CreateNode(ASsnCase)
	if node == nil {
		return nil, nil
	}

	var t sToken
	pr.GetToken(&t)
	if t.Type != tokens.ASttCase && t.Type != tokens.ASttDefault {
		return nil, errors.New("Expected " + tokens.GetDefinitionOrList([]tokens.Token{tokens.ASttCase, tokens.ASttDefault}) + ", got " + tokens.GetDefinition(t.Type))
		return node, nil
	}

	node.UpdateSourcePos(t.Position, t.Length)

	if t.Type == tokens.ASttCase {
		node.AddChildLast(pr.ParseExpression())
	}

	pr.GetToken(&t)
	if t.Type != tokens.ASttColon {
		//TODO: error expected :
		return node, nil
	}

	pr.GetToken(&t)
	pr.RewindTo(&t)
	for t.Type != tokens.ASttCase &&
		t.Type != tokens.ASttDefault &&
		t.Type != tokens.ASttEndStatementBlock &&
		t.Type != tokens.ASttBreak {
		if pr.IsVarDecl() {
			node.AddChildLast(pr.ParseDeclaration(false, false))
		} else {
			node.AddChildLast(pr.ParseStatement())
		}
		if pr.IsSyntaxError {
			return node, nil
		}

		pr.GetToken(&t)
		pr.RewindTo(&t)
	}

	if t.Type == tokens.ASttBreak {
		node.AddChildLast(pr.ParseBreak())
	}
	return node, nil
}

func (pr *Parser) ParseIf() (*ScriptNode, error) {
	node := pr.CreateNode(ASsnIf)

	var t sToken
	pr.GetToken(&t)
	if t.Type != tokens.ASttIf {
		//TODO: error expected if
		return node, nil
	}

	node.UpdateSourcePos(t.Position, t.Length)

	pr.GetToken(&t)
	if t.Type != tokens.ASttOpenParanthesis {
		errors.New("Expected " + tokens.GetDefinitionOrList([]tokens.Token{tokens.ASttOpenParanthesis}) + ", got " + tokens.GetDefinition(t.Type))
		return node, nil
	}

	node.AddChildLast(pr.ParseAssignment())
	if pr.IsSyntaxError {
		return node, nil
	}

	pr.GetToken(&t)
	if t.Type != tokens.ASttCloseParanthesis {
		errors.New("Expected " + tokens.GetDefinitionOrList([]tokens.Token{tokens.ASttCloseParanthesis}) + ", got " + tokens.GetDefinition(t.Type))
		return node, nil
	}

	node.AddChildLast(pr.ParseStatement())
	if pr.IsSyntaxError {
		return node, nil
	}

	pr.GetToken(&t)
	if t.Type != tokens.ASttElse {
		pr.RewindTo(&t)
		return node, nil
	}

	node.AddChildLast(pr.ParseStatement())

	return node, nil
}

func (pr *Parser) ParseFor() (*ScriptNode, error) {
	node := pr.CreateNode(ASsnFor)
	if node == nil {
		return nil, nil
	}

	var t sToken
	pr.GetToken(&t)
	if t.Type != tokens.ASttFor {
		//TODO: error expected for
		return node, nil
	}

	node.UpdateSourcePos(t.Position, t.Length)

	pr.GetToken(&t)
	if t.Type != tokens.ASttOpenParanthesis {
		errors.New("Expected " + tokens.GetDefinitionOrList([]tokens.Token{tokens.ASttOpenParanthesis}) + ", got " + tokens.GetDefinition(t.Type))
		return node, nil
	}

	if pr.IsVarDecl() {
		node.AddChildLast(pr.ParseDeclaration(false, false))
	} else {
		node.AddChildLast(pr.ParseExpressionStatement())
	}
	if pr.IsSyntaxError {
		return node, nil
	}

	node.AddChildLast(pr.ParseExpressionStatement())
	if pr.IsSyntaxError {
		return node, nil
	}

	pr.GetToken(&t)
	if t.Type != tokens.ASttCloseParanthesis {
		pr.RewindTo(&t)

		for {
			n := pr.CreateNode(ASsnExpressionStatement)
			if n == nil {
				return nil, nil
			}
			node.AddChildLast(n, nil)
			n.AddChildLast(pr.ParseAssignment())
			if pr.IsSyntaxError {
				return node, nil
			}

			pr.GetToken(&t)
			if t.Type == tokens.ASttListSeparator {
				continue
			} else if t.Type == tokens.ASttCloseParanthesis {
				break
			} else {
				//TODO: error expected , or )
				return node, nil
			}
		}
	}

	node.AddChildLast(pr.ParseStatement())

	return node, nil
}

func (pr *Parser) ParseWhile() (*ScriptNode, error) {
	node := pr.CreateNode(ASsnWhile)
	if node == nil {
		return nil, nil
	}

	var t sToken
	pr.GetToken(&t)
	if t.Type != tokens.ASttWhile {
		//TODO: error expected while
		return node, nil
	}

	node.UpdateSourcePos(t.Position, t.Length)

	pr.GetToken(&t)
	if t.Type != tokens.ASttOpenParanthesis {
		errors.New("Expected " + tokens.GetDefinitionOrList([]tokens.Token{tokens.ASttOpenParanthesis}) + ", got " + tokens.GetDefinition(t.Type))
		return node, nil
	}

	node.AddChildLast(pr.ParseAssignment())
	if pr.IsSyntaxError {
		return node, nil
	}

	pr.GetToken(&t)
	if t.Type != tokens.ASttCloseParanthesis {
		errors.New("Expected " + tokens.GetDefinitionOrList([]tokens.Token{tokens.ASttCloseParanthesis}) + ", got " + tokens.GetDefinition(t.Type))
		return node, nil
	}

	node.AddChildLast(pr.ParseStatement())

	return node, nil
}

func (pr *Parser) ParseDoWhile() (*ScriptNode, error) {
	node := pr.CreateNode(ASsnDoWhile)
	if node == nil {
		return nil, nil
	}

	var t sToken
	pr.GetToken(&t)
	if t.Type != tokens.ASttDo {
		//TODO: error expected do
		return node, nil
	}

	node.UpdateSourcePos(t.Position, t.Length)

	node.AddChildLast(pr.ParseStatement())
	if pr.IsSyntaxError {
		return node, nil
	}

	pr.GetToken(&t)
	if t.Type == tokens.ASttWhile {
		//TODO: error expected while
		return node, nil
	}

	pr.GetToken(&t)
	if t.Type != tokens.ASttOpenParanthesis {
		errors.New("Expected " + tokens.GetDefinitionOrList([]tokens.Token{tokens.ASttOpenParanthesis}) + ", got " + tokens.GetDefinition(t.Type))
		return node, nil
	}

	node.AddChildLast(pr.ParseAssignment())
	if pr.IsSyntaxError {
		return node, nil
	}

	pr.GetToken(&t)
	if t.Type != tokens.ASttCloseParanthesis {
		errors.New("Expected " + tokens.GetDefinitionOrList([]tokens.Token{tokens.ASttCloseParanthesis}) + ", got " + tokens.GetDefinition(t.Type))
		return node, nil
	}

	pr.GetToken(&t)
	if t.Type != tokens.ASttEndStatement {
		errors.New("Expected " + tokens.GetDefinitionOrList([]tokens.Token{tokens.ASttEndStatement}) + ", got " + tokens.GetDefinition(t.Type))
		return node, nil
	}

	node.UpdateSourcePos(t.Position, t.Length)

	return node, nil

}

func (pr *Parser) ParseReturn() (*ScriptNode, error) {
	node := pr.CreateNode(ASsnReturn)
	if node == nil {
		return nil, nil
	}

	var t sToken
	pr.GetToken(&t)
	if t.Type != tokens.ASttReturn {
		//TODO: error expected return
		return node, nil
	}

	node.UpdateSourcePos(t.Position, t.Length)

	pr.GetToken(&t)
	if t.Type == tokens.ASttEndStatement {
		node.UpdateSourcePos(t.Position, t.Length)
		return node, nil
	}
	pr.RewindTo(&t)

	node.AddChildLast(pr.ParseAssignment())
	if pr.IsSyntaxError {
		return node, nil
	}

	pr.GetToken(&t)
	if t.Type != tokens.ASttEndStatement {
		errors.New("Expected " + tokens.GetDefinitionOrList([]tokens.Token{tokens.ASttEndStatement}) + ", got " + tokens.GetDefinition(t.Type))
		return node, nil
	}

	node.UpdateSourcePos(t.Position, t.Length)
	return node, nil
}

func (pr *Parser) ParseBreak() (*ScriptNode, error) {
	node := pr.CreateNode(ASsnBreak)
	if node == nil {
		return nil, nil
	}

	var t sToken
	pr.GetToken(&t)
	if t.Type != tokens.ASttBreak {
		errors.New("Expected " + tokens.GetDefinitionOrList([]tokens.Token{tokens.ASttBreak}) + ", got " + tokens.GetDefinition(t.Type))
		return node, nil
	}

	node.UpdateSourcePos(t.Position, t.Length)

	pr.GetToken(&t)
	if t.Type != tokens.ASttEndStatement {
		errors.New("Expected " + tokens.GetDefinitionOrList([]tokens.Token{tokens.ASttEndStatement}) + ", got " + tokens.GetDefinition(t.Type))
	}

	node.UpdateSourcePos(t.Position, t.Length)
	return node, nil
}

func (pr *Parser) ParseContinue() (*ScriptNode, error) {
	node := pr.CreateNode(ASsnContinue)
	if node == nil {
		return nil, nil
	}

	var t sToken
	pr.GetToken(&t)
	if t.Type != tokens.ASttContinue {
		errors.New("Expected " + tokens.GetDefinitionOrList([]tokens.Token{tokens.ASttContinue}) + ", got " + tokens.GetDefinition(t.Type))
		return node, nil
	}

	node.UpdateSourcePos(t.Position, t.Length)

	pr.GetToken(&t)
	if t.Type != tokens.ASttEndStatement {
		//TODO: error expected token ;
	}

	node.UpdateSourcePos(t.Position, t.Length)
	return node, nil
}

func (pr *Parser) ParseTypedef() (*ScriptNode, error) {
	node := pr.CreateNode(ASsnTypedef)
	if node == nil {
		return nil, nil
	}

	var token sToken

	pr.GetToken(&token)
	if token.Type != tokens.ASttTypedef {
		errors.New("Expected " + tokens.GetDefinitionOrList([]tokens.Token{tokens.ASttTypedef}) + ", got " + tokens.GetDefinition(token.Type))
		return node, nil
	}

	node.SetToken(&token)
	node.UpdateSourcePos(token.Position, token.Length)

	pr.GetToken(&token)
	pr.RewindTo(&token)

	if !pr.IsRealType(token.Type) || token.Type == tokens.ASttVoid {
		errors.New("Unexpected token: " + tokens.GetDefinition(token.Type))
		return node, nil
	}

	node.AddChildLast(pr.ParseRealType())
	node.AddChildLast(pr.ParseIdentifier())

	pr.GetToken(&token)
	if token.Type != tokens.ASttEndStatement {
		pr.RewindTo(&token)
		errors.New("Expected " + tokens.GetDefinitionOrList([]tokens.Token{tokens.ASttEndStatement}) + ", got " + tokens.GetDefinition(token.Type))
	}

	return node, nil
}

func (pr *Parser) ParseMethodOverrideBehaviors(node *ScriptNode) {
	var t sToken
	for {
		pr.GetToken(&t)
		pr.RewindTo(&t)
		if pr.IdentifierIs(t, tokens.ASFinalToken) || pr.IdentifierIs(t, tokens.ASOverrideToken) {
			node.AddChildLast(pr.ParseIdentifier())
		} else {
			break
		}
	}
}
