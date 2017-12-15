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
	_ "fmt"
	"github.com/Member1221/go-angelscript/tokenizer"
)

type Parser struct {
	ErrorWhileParsing     bool
	IsSyntaxError         bool
	CheckValidTypes       bool
	IsParsingAppInterface bool

	Engine *ScriptEngine
	Script *ScriptCode
	Node   *ScriptNode
	Builder *ScriptBuilder

	tempstr   string
	lastToken *sToken
	sourcePos int
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
	pr.Node = pr.ParseFunctionDefinition()

	if expectListPattern {
		pr.Node.AddChildLast(pr.ParseListPattern())
	}

	if !pr.IsSyntaxError {
		var t sToken
		pr.GetToken(&t)
		if t.Type != angelscript.ASttEnd {
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

	pr.Node = pr.ParseScript(false)

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
	pr.Node = pr.ParseStatementBlock()

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

	var t sToken
	pr.GetToken(&t)
	if t.Type == angelscript.ASttAssignment {
		pr.GetToken(&t)
		pr.RewindTo(&t)
		if t.Type == angelscript.ASttStartStatementBlock {
			pr.Node = pr.ParseInitList()
		} else {
			pr.Node = pr.ParseAssignment()
		}
	} else if t.Type == angelscript.ASttOpenParanthesis {
		pr.RewindTo(&t)
		pr.Node = pr.ParseArgList(true)
	} else {
		//TODO: error expected assignment or open paranthesis
	}

	pr.GetToken(&t)
	if t.Type != angelscript.ASttEnd && t.Type != angelscript.ASttEndStatement && t.Type != angelscript.ASttListSeparator && t.Type != angelscript.ASttEndStatementBlock {
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

	pr.Node = pr.ParseExpression()
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
		if t.Type != angelscript.ASttEnd {
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
	if t.Type == angelscript.ASttAmp {
		pr.Node.AddChildLast(pr.ParseToken(angelscript.ASttAmp))
	}

	if !pr.IsSyntaxError {
		if t.Type != angelscript.ASttEnd {
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
		if t.Type != angelscript.ASttEnd {
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
	if t1.Type == angelscript.ASttScope {
		pr.RewindTo(&t1)
		scope.AddChildLast(pr.ParseToken(angelscript.ASttScope))
		pr.GetToken(&t1)
		pr.GetToken(&t2)
	}
	for t1.Type == angelscript.ASttIdentifier && t2.Type == angelscript.ASttScope {
		pr.RewindTo(&t1)
		scope.AddChildLast(pr.ParseIdentifier())
		scope.AddChildLast(pr.ParseToken(angelscript.ASttScope))
		pr.GetToken(&t1)
		pr.GetToken(&t2)
	}

	if t1.Type == angelscript.ASttIdentifier && t2.Type == angelscript.ASttLessThan {
		pr.tempstr = pr.Script.Code[t1.Position:t1.Length]
		if pr.Engine.IsTemplateType(pr.tempstr) {

			pr.RewindTo(&t1)
			restore := scope.LastChild
			scope.AddChildLast(pr.ParseIdentifier())
			if pr.ParseTemplTypeList(scope, false) {
				pr.GetToken(&t2)
				if t2.Type == angelscript.ASttScope {
					pr.Node.AddChildLast(scope)
					return
				} else {
					pr.RewindTo(&t1)

					for scope.LastChild != restore {
						last := scope.LastChild
						last.DisconnectParent()
						last.Destroy(pr.Engine)
					}
					if scope.LastChild != nil {
						pr.Node.AddChildLast(scope)
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
		pr.Node.AddChildLast(scope)
	} else {
		scope.Destroy(pr.Engine)
	}
}

func (pr *Parser) ParseFunctionDefinition() *ScriptNode {
	node := pr.CreateNode(ASsnFunction)
	if node == nil {
		return nil
	}

	node.AddChildLast(pr.ParseType(true, false, false))
	if pr.IsSyntaxError {
		return node
	}

	node.AddChildLast(pr.ParseTypeMod(false))
	if pr.IsSyntaxError {
		return node
	}

	pr.ParseOptionalScope(node)

	node.AddChildLast(pr.ParseIdentifier())
	if pr.IsSyntaxError {
		return node
	}

	node.AddChildLast(pr.ParseParameterList())
	if pr.IsSyntaxError {
		return node
	}

	var t1 sToken
	pr.GetToken(&t1)
	pr.RewindTo(&t1)
	if t1.Type == angelscript.ASttConst {
		node.AddChildLast(pr.ParseToken(angelscript.ASttConst))
	}
	return node
}

func (pr *Parser) ParseTypeMod(isParam bool) *ScriptNode {
	node := pr.CreateNode(ASsnDataType)
	if node == nil {
		return nil
	}

	var t sToken

	pr.GetToken(&t)
	pr.RewindTo(&t)

	if t.Type == angelscript.ASttAmp {
		node.AddChildLast(pr.ParseToken(angelscript.ASttAmp))
		if pr.IsSyntaxError {
			return node
		}
		if isParam {
			pr.GetToken(&t)
			pr.RewindTo(&t)

			if t.Type == angelscript.ASttIn || t.Type == angelscript.ASttOut || t.Type == angelscript.ASttInOut {
				tokens := []angelscript.Token{angelscript.ASttIn, angelscript.ASttOut, angelscript.ASttInOut}
				node.AddChildLast(pr.ParseOneOf(tokens))
			}
		}
	}

	pr.GetToken(&t)
	pr.RewindTo(&t)
	if t.Type == angelscript.ASttPlus {
		node.AddChildLast(pr.ParseToken(angelscript.ASttPlus))
		if pr.IsSyntaxError {
			return node
		}
	}

	pr.GetToken(&t)
	pr.RewindTo(&t)
	if pr.IdentifierIs(t, angelscript.ASIfHandleToken) {
		node.AddChildLast(pr.ParseToken(angelscript.ASttIdentifier))
		if pr.IsSyntaxError {
			return node
		}
	}

	return node
}

func (pr *Parser) ParseType(allowConst, allowVariableType, allowAuto bool) *ScriptNode {
	node := pr.CreateNode(ASsnDataType)
	if node == nil {
		return nil
	}

	var t sToken

	if allowConst {
		pr.GetToken(&t)
		pr.RewindTo(&t)
		if t.Type == angelscript.ASttConst {
			node.AddChildLast(pr.ParseToken(angelscript.ASttConst))
			if pr.IsSyntaxError {
				return node
			}
		}
	}

	pr.ParseOptionalScope(node)

	node.AddChildLast(pr.ParseDataType(allowVariableType, allowAuto))
	if pr.IsSyntaxError {
		return node
	}

	pr.GetToken(&t)
	pr.RewindTo(&t)
	tr := node.LastChild

	pr.tempstr = pr.Script.Code[tr.TokenPosition:tr.TokenLength]
	if pr.Engine.IsTemplateType(pr.tempstr) && t.Type == angelscript.ASttLessThan {
		pr.ParseTemplTypeList(node, true)
		if pr.IsSyntaxError {
			return node
		}
	}

	pr.GetToken(&t)
	pr.RewindTo(&t)
	for t.Type == angelscript.ASttOpenBracket || t.Type == angelscript.ASttHandle {
		if t.Type == angelscript.ASttOpenBracket {
			node.AddChildLast(pr.ParseToken(angelscript.ASttOpenBracket))
			if pr.IsSyntaxError {
				return node
			}

			pr.GetToken(&t)
			if t.Type == angelscript.ASttCloseBracket {
				//TODO: ERROR (expect ])
				return node
			}
		} else {
			node.AddChildLast(pr.ParseToken(angelscript.ASttHandle))
			if pr.IsSyntaxError {
				return node
			}
		}

		pr.GetToken(&t)
		pr.RewindTo(&t)
	}
	return node
}

func (pr *Parser) ParseTemplTypeList(node *ScriptNode, required bool) bool {
	var t sToken
	isValid := true

	last := node.LastChild

	pr.GetToken(&t)

	if t.Type != angelscript.ASttLessThan {
		if required {
			//TODO: ERROR (expect Lessthan)
		}
		return false
	}

	node.AddChildLast(pr.ParseType(true, false, false))
	if pr.IsSyntaxError {
		return false
	}

	pr.GetToken(&t)

	for t.Type == angelscript.ASttListSeparator {
		node.AddChildLast(pr.ParseType(true, false, false))
		if pr.IsSyntaxError {
			return false
		}
		pr.GetToken(&t)
	}

	if pr.Script.Code[t.Position:1] != ">" {
		if required {
			//TODO: ERROR (Expect GreaterThan)
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
		return false
	}

	return true

}

func (pr *Parser) ParseToken(token angelscript.Token) *ScriptNode {
	node := pr.CreateNode(ASsnUndefined)
	if node == nil {
		return nil
	}

	var t1 sToken

	pr.GetToken(&t1)

	if t1.Type != token {
		//TODO: ERROR (Expect TOKEN)
		return node
	}

	node.SetToken(&t1)
	node.UpdateSourcePos(t1.Position, t1.Length)
	return node
}

func (pr *Parser) ParseOneOf(tokens []angelscript.Token) *ScriptNode {
	node := pr.CreateNode(ASsnUndefined)
	if node == nil {
		return nil
	}

	var t1 sToken
	pr.GetToken(&t1)

	n := 0
	for n = 0; n < len(tokens); n++ {
		if tokens[n] == t1.Type {
			break
		}
	}

	if n == len(tokens) {
		//TODO: ERROR (Expect tokens/count, got t1)
		return node
	}

	node.SetToken(&t1)
	node.UpdateSourcePos(t1.Position, t1.Length)

	return node
}

func (pr *Parser) ParseDataType(allowVariableType, allowAuto bool) *ScriptNode {
	node := pr.CreateNode(ASsnDataType)
	if node == nil {
		return nil
	}

	var t1 sToken
	pr.GetToken(&t1)
	if !pr.IsDataType(&t1) && !(allowVariableType && t1.Type == angelscript.ASttQuestion) && !(allowAuto && t1.Type == angelscript.ASttAuto) {
		if t1.Type == angelscript.ASttIdentifier {
			//TODO: FATAL ERROR:
			/*
				asCString errMsg;
				tempString.Assign(&script->code[t1.pos], t1.length);
				errMsg.Format(TXT_IDENTIFIER_s_NOT_DATA_TYPE, tempString.AddressOf());
				Error(errMsg, &t1);
			*/
		} else if t1.Type == angelscript.ASttAuto {
			//TODO: ERROR TXT_AUTO_NOT_ALLOWED
		} else {
			//TODO: ERROR TXT_EXPECTED_DATA_TYPE
		}
		return node
	}

	node.SetToken(&t1)
	node.UpdateSourcePos(t1.Position, t1.Length)
	return node
}

func (pr *Parser) ParseRealType() *ScriptNode {
	node := pr.CreateNode(ASsnDataType)
	if node == nil {
		return nil
	}

	var t1 sToken
	pr.GetToken(&t1)
	if !pr.IsRealType(t1.Type) {
		//TODO: ERROR TXT_EXPECTED_DATATYPE
		return node
	}

	node.SetToken(&t1)
	node.UpdateSourcePos(t1.Position, t1.Length)
	return node
}

func (pr *Parser) ParseIdentifier() *ScriptNode {
	node := pr.CreateNode(ASsnIdentifier)
	if node == nil {
		return nil
	}

	var t1 sToken
	pr.GetToken(&t1)
	if t1.Type != angelscript.ASttIdentifier {
		//TODO: ERROR TXT_EXPECTED_DATATYPE
		return node
	}

	node.SetToken(&t1)
	node.UpdateSourcePos(t1.Position, t1.Length)
	return node
}

func (pr *Parser) ParseParameterList() *ScriptNode {
	node := pr.CreateNode(ASsnParameterList)
	if node == nil {
		return nil
	}

	var t1 sToken
	pr.GetToken(&t1)
	if t1.Type != angelscript.ASttOpenParanthesis {
		//TODO: ERROR (Expected "(")
		return node
	}

	node.UpdateSourcePos(t1.Position, t1.Length)
	pr.GetToken(&t1)
	if t1.Type == angelscript.ASttCloseParanthesis {
		node.UpdateSourcePos(t1.Position, t1.Length)
		return node
	} else {
		if t1.Type == angelscript.ASttVoid {
			var t2 sToken
			pr.GetToken(&t2)
			if t2.Type == angelscript.ASttCloseParanthesis {
				node.UpdateSourcePos(t2.Position, t2.Length)
				return node
			}
		}

		pr.RewindTo(&t1)

		for {
			node.AddChildLast(pr.ParseType(true, pr.IsParsingAppInterface, false))
			if pr.IsSyntaxError {
				return node
			}

			node.AddChildLast(pr.ParseTypeMod(true))
			if pr.IsSyntaxError {
				return node
			}

			pr.GetToken(&t1)
			if t1.Type == angelscript.ASttIdentifier {
				pr.RewindTo(&t1)

				node.AddChildLast(pr.ParseIdentifier())
				if pr.IsSyntaxError {
					return node
				}

				pr.GetToken(&t1)
			}

			if t1.Type == angelscript.ASttAssignment {
				node.AddChildLast(pr.SuperficiallyParseExpression())
				if pr.IsSyntaxError {
					return node
				}

				pr.GetToken(&t1)
			}

			if t1.Type == angelscript.ASttCloseParanthesis {
				node.UpdateSourcePos(t1.Position, t1.Length)
			} else if t1.Type == angelscript.ASttListSeparator {
				continue
			} else {
				//TODO: Error (Expected Tokens: ")", ",")
				return node
			}
		}
	}
	return node
}

func (pr *Parser) SuperficiallyParseExpression() *ScriptNode {
	node := pr.CreateNode(ASsnParameterList)
	if node == nil {
		return nil
	}

	var start sToken
	pr.GetToken(&start)
	pr.RewindTo(&start)

	stack := ""
	var t sToken
	for {
		pr.GetToken(&t)

		if t.Type == angelscript.ASttOpenParanthesis {
			stack += "("
		} else if t.Type == angelscript.ASttCloseParanthesis {
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
				return node
			}
		} else if t.Type == angelscript.ASttListSeparator {
			if stack == "" {
				pr.RewindTo(&t)
				break
			}
		} else if t.Type == angelscript.ASttStartStatementBlock {
			stack += "{"
		} else if t.Type == angelscript.ASttEndStatementBlock {
			if stack == "" || stack[len(stack)-1:1] == "{" {
				pr.RewindTo(&t)
				/*
					asCString str;
					str.Format(TXT_UNEXPECTED_TOKEN_s, "}");
					Error(str, &t);
				*/
				return node
			} else {
				stack = stack[:len(stack)-1]
			}
		} else if t.Type == angelscript.ASttEndStatement {
			pr.RewindTo(&t)
			/*
				asCString str;
				str.Format(TXT_UNEXPECTED_TOKEN_s, ";");
				Error(str, &t);
			*/
			return node
		} else if t.Type == angelscript.ASttNonTerminatedStringConstant {
			pr.RewindTo(&t)
			//TODO: ERROR (TXT_NONTERMINATED_STRING)
			return node
		} else if t.Type == angelscript.ASttEnd {
			pr.RewindTo(&t)
			//TODO: ERROR (TXT_UNEXPECTED_END_OF_FILE)
			return node
		}

		node.UpdateSourcePos(t.Position, t.Length)
	}

	return node
}

func (pr *Parser) GetToken(token *sToken) {
	if pr.lastToken.Position == pr.sourcePos {
		token = pr.lastToken
		pr.sourcePos += token.Length

		if token.Type == angelscript.ASttWhiteSpace ||
			token.Type == angelscript.ASttOnelineComment ||
			token.Type == angelscript.ASttMultilineComment {
			pr.GetToken(token)
		}
		return
	}

	sl := pr.Script.Length
	for token.Type == angelscript.ASttWhiteSpace ||
		token.Type == angelscript.ASttOnelineComment ||
		token.Type == angelscript.ASttMultilineComment {
		if pr.sourcePos >= sl {
			token.Type = angelscript.ASttEnd
			token.Length = 0
		} else {
			//TODO: token->type = engine->tok.GetToken(&script->code[sourcePos], sourceLength - sourcePos, &token->length);
		}

		token.Position = pr.sourcePos
		pr.sourcePos += token.Length
	}

}

func (pr *Parser) RewindTo(token *sToken) {
	pr.lastToken = token
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

func (pr *Parser) IsRealType(tokenType angelscript.Token) bool {
	if tokenType == angelscript.ASttVoid ||
		tokenType == angelscript.ASttInt ||
		tokenType == angelscript.ASttInt8 ||
		tokenType == angelscript.ASttInt16 ||
		tokenType == angelscript.ASttInt64 ||
		tokenType == angelscript.ASttUInt ||
		tokenType == angelscript.ASttUInt8 ||
		tokenType == angelscript.ASttUInt16 ||
		tokenType == angelscript.ASttUInt64 ||
		tokenType == angelscript.ASttFloat ||
		tokenType == angelscript.ASttDouble ||
		tokenType == angelscript.ASttBool {
		return true
	}
	return false
}

func (pr *Parser) IsDataType(token *sToken) bool {
	if token.Type == angelscript.ASttIdentifier {
		if pr.CheckValidTypes {
			pr.tempstr = pr.Script.Code[token.Position:token.Length]
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

func (pr *Parser) ExpectedOneOfMap(tokens map[string][]angelscript.Token) string {
	return ""
}

func (pr *Parser) InsteadFound(token sToken) string {
	return ""
}

func (pr *Parser) ParseListPattern() *ScriptNode {
	node := pr.CreateNode(ASsnParameterList)
	if node == nil {
		return nil
	}

	var t1 sToken
	pr.GetToken(&t1)
	if t1.Type != angelscript.ASttStartStatementBlock {
		//TODO: ERROR (Expected "{")
		return node
	}

	node.UpdateSourcePos(t1.Position, t1.Length)
	//var start sToken

	isBeginning := true
	afterType := false

	for !pr.IsSyntaxError {
		pr.GetToken(&t1)
		if t1.Type == angelscript.ASttEndStatementBlock {
			if !afterType {
				//TODO: ERROR TXT_EXPECTED_DATA_TYPE
			}
			break
		} else if t1.Type == angelscript.ASttStartStatementBlock {
			if afterType {
				//TODO: ERROR (Expected ",", "}")
			}
			pr.RewindTo(&t1)
			node.AddChildLast(pr.ParseListPattern())
			afterType = true
		} else if t1.Type == angelscript.ASttIdentifier && (pr.IdentifierIs(t1, "repeat") || pr.IdentifierIs(t1, "repeat_same")) {
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
		} else if t1.Type == angelscript.ASttEnd {
			//TODO: ERROR TXT_UNEXPECTED_END_OF_FILE
			break
		} else if t1.Type == angelscript.ASttListSeparator {
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

	return node
}

func (pr *Parser) IdentifierIs(t sToken, str string) bool {
	if t.Type != angelscript.ASttIdentifier {
		return false
	}
	return pr.Script.TokenEquals(t.Position, t.Length, str)
}

func (pr *Parser) CheckTemplateType(t sToken) bool {
	pr.tempstr = pr.Script.Code[t.Position:t.Length]
	if pr.Engine.IsTemplateType(pr.tempstr) {
		var t1 sToken
		pr.GetToken(&t1)
		if t1.Type == angelscript.ASttLessThan {
			pr.RewindTo(&t1)
			return true
		}

		for {
			pr.GetToken(&t1)
			if t1.Type == angelscript.ASttScope {
				pr.GetToken(&t1)
			}

			var t2 sToken
			pr.GetToken(&t2)
			for t1.Type == angelscript.ASttIdentifier && t2.Type == angelscript.ASttScope {
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

			for t1.Type == angelscript.ASttHandle || t1.Type == angelscript.ASttOpenBracket {
				if t1.Type == angelscript.ASttOpenBracket {
					pr.GetToken(&t1)
					if t1.Type != angelscript.ASttCloseBracket {
						return false
					}
				}

				pr.GetToken(&t1)
			}

			if t1.Type == angelscript.ASttListSeparator {
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

func (pr *Parser) ParseCast() *ScriptNode {
	node := pr.CreateNode(ASsnCast)
	if node == nil {
		return nil
	}

	var t1 sToken
	pr.GetToken(&t1)
	if t1.Type == angelscript.ASttCast {
		//TODO: error expected cast
		return node
	}

	node.UpdateSourcePos(t1.Position, t1.Length)

	pr.GetToken(&t1)
	if t1.Type != angelscript.ASttLessThan {
		//TODO: error expected <
		return node
	}

	node.AddChildLast(pr.ParseType(true, false, false))
	if pr.IsSyntaxError {
		return node
	}

	pr.GetToken(&t1)
	if t1.Type != angelscript.ASttGreaterThan {
		//TODO: error expected >
	}

	pr.GetToken(&t1)
	if t1.Type != angelscript.ASttOpenParanthesis {
		//TODO: error expected (
		return node
	}

	node.AddChildLast(pr.ParseAssignment())
	if pr.IsSyntaxError {
		return node
	}

	pr.GetToken(&t1)
	if t1.Type != angelscript.ASttCloseParanthesis {
		//TODO: error expected )
		return node
	}

	node.UpdateSourcePos(t1.Position, t1.Length)

	return node
}

func (pr *Parser) ParseExprValue() *ScriptNode {
	node := pr.CreateNode(ASsnExprValue)
	if node == nil {
		return nil
	}

	var t1 sToken
	var t2 sToken
	pr.GetToken(&t1)
	pr.GetToken(&t2)
	pr.RewindTo(&t1)

	if t1.Type == angelscript.ASttVoid {
		node.AddChildLast(pr.ParseToken(angelscript.ASttVoid))
	} else if pr.IsRealType(t1.Type) {
		node.AddChildLast(pr.ParseConstructCall())
	} else if t1.Type == angelscript.ASttIdentifier || t1.Type == angelscript.ASttScope {
		if pr.IsLambda() {
			node.AddChildLast(pr.ParseLambda())
		} else {
			var t sToken
			if t1.Type == angelscript.ASttScope {
				t = t2
			} else {
				t = t1
			}
			pr.RewindTo(&t)
			pr.GetToken(&t2)
			for t.Type == angelscript.ASttIdentifier {
				t2 = t
				pr.GetToken(&t)
				if t.Type == angelscript.ASttScope {
					pr.GetToken(&t)
				} else {
					break
				}
			}

			isDataType := pr.IsDataType(&t2)
			isTemplateType := false

			if isDataType {
				pr.tempstr = pr.Script.Code[t2.Position:t2.Length]
				if pr.Engine.IsTemplateType(pr.tempstr) {
					isTemplateType = true
				}
			}

			pr.GetToken(&t2)

			pr.RewindTo(&t1)

			if isDataType && (t.Type == angelscript.ASttOpenParanthesis || (t.Type == angelscript.ASttOpenBracket && t2.Type == angelscript.ASttCloseBracket)) {
				node.AddChildLast(pr.ParseConstructCall())
			} else if isTemplateType && t.Type == angelscript.ASttLessThan {
				node.AddChildLast(pr.ParseConstructCall())
			} else if pr.IsFunctionCall() {
				node.AddChildLast(pr.ParseFunctionCall())
			} else {
				node.AddChildLast(pr.ParseVariableAccess())
			}
		}
	} else if t1.Type == angelscript.ASttCast {
		node.AddChildLast(pr.ParseCast())
	} else if pr.IsConstant(t1.Type) {
		node.AddChildLast(pr.ParseConstant())
	} else if t1.Type == angelscript.ASttOpenParanthesis {
		pr.GetToken(&t1)
		node.UpdateSourcePos(t1.Position, t1.Length)

		node.AddChildLast(pr.ParseAssignment())
		if pr.IsSyntaxError {
			return node
		}

		pr.GetToken(&t1)
		if t1.Type != angelscript.ASttCloseParanthesis {
			//TODO: error expected ) (NO RETURN NODE)
		}

		node.UpdateSourcePos(t1.Position, t1.Length)
	} else {
		//TODO: error TXT_EXPECTED_EXPRESSION_VALUE
	}

	return node
}

func (pr *Parser) ParseConstant() *ScriptNode {
	node := pr.CreateNode(ASsnConstant)
	if node == nil {
		return nil
	}

	var t sToken
	pr.GetToken(&t)
	if !(pr.IsConstant(t.Type)) {
		//TODO: error TXT_EXPECTED_CONSTANT
		return node
	}

	node.SetToken(&t)
	node.UpdateSourcePos(t.Position, t.Length)

	if t.Type == angelscript.ASttStringConstant || t.Type == angelscript.ASttMultilineStringConstant || t.Type == angelscript.ASttHeredocStringConstant {
		pr.RewindTo(&t)
	}

	for t.Type == angelscript.ASttStringConstant || t.Type == angelscript.ASttMultilineStringConstant || t.Type == angelscript.ASttHeredocStringConstant {
		node.AddChildLast(pr.ParseStringConstant())

		pr.GetToken(&t)
		pr.RewindTo(&t)
	}

	return node
}

func (pr *Parser) IsLambda() bool {
	isLambda := false
	var t sToken
	pr.GetToken(&t)
	if t.Type == angelscript.ASttIdentifier && pr.IdentifierIs(t, angelscript.ASFunctionToken) {
		var t2 sToken
		pr.GetToken(&t2)
		if t2.Type == angelscript.ASttOpenParanthesis {
			for t2.Type != angelscript.ASttCloseParanthesis && t2.Type != angelscript.ASttEnd {
				pr.GetToken(&t2)
			}

			pr.GetToken(&t2)

			if t2.Type == angelscript.ASttStartStatementBlock {
				isLambda = true
			}
		}
	}

	pr.RewindTo(&t)
	return isLambda
}

func (pr *Parser) ParseLambda() *ScriptNode {
	node := pr.CreateNode(ASsnFunction)
	if node == nil {
		return nil
	}

	var t sToken
	pr.GetToken(&t)

	if t.Type != angelscript.ASttIdentifier || !pr.IdentifierIs(t, angelscript.ASFunctionToken) {
		//TODO: error expected function token
		return node
	}

	pr.GetToken(&t)
	if t.Type != angelscript.ASttOpenParanthesis {
		//TODO: error expected (
		return node
	}

	pr.GetToken(&t)
	if t.Type == angelscript.ASttIdentifier {
		pr.RewindTo(&t)
		node.AddChildLast(pr.ParseIdentifier())

		pr.GetToken(&t)
		for t.Type == angelscript.ASttListSeparator {
			node.AddChildLast(pr.ParseIdentifier())
			if pr.IsSyntaxError {
				return node
			}

			pr.GetToken(&t)
		}
	}

	if t.Type != angelscript.ASttCloseParanthesis {
		//TODO: error expected )
		return node
	}

	node.AddChildLast(pr.SuperficiallyParseStatementBlock())

	return node
}

func (pr *Parser) ParseStringConstant() *ScriptNode {
	node := pr.CreateNode(ASsnConstant)
	if node == nil {
		return nil
	}

	var t sToken
	pr.GetToken(&t)
	if t.Type != angelscript.ASttStringConstant && t.Type != angelscript.ASttMultilineComment && t.Type != angelscript.ASttHeredocStringConstant {
		//TODO: error TXT_EXPECTED_STRING
		return node
	}

	node.SetToken(&t)
	node.UpdateSourcePos(t.Position, t.Length)

	return node
}

func (pr *Parser) ParseFunctionCall() *ScriptNode {
	node := pr.CreateNode(ASsnFunctionCall)
	if node == nil {
		return nil
	}

	pr.ParseOptionalScope(node)

	node.AddChildLast(pr.ParseIdentifier())
	if pr.IsSyntaxError {
		return node
	}

	node.AddChildLast(pr.ParseArgList(true))

	return node
}

func (pr *Parser) ParseVariableAccess() *ScriptNode {
	node := pr.CreateNode(ASsnVariableAccess)
	if node == nil {
		return nil
	}

	pr.ParseOptionalScope(node)

	node.AddChildLast(pr.ParseIdentifier())

	return node
}

func (pr *Parser) ParseConstructCall() *ScriptNode {
	node := pr.CreateNode(ASsnConstructCall)
	if node == nil {
		return nil
	}

	node.AddChildLast(pr.ParseType(false, false, false))
	if pr.IsSyntaxError {
		return node
	}

	node.AddChildLast(pr.ParseArgList(true))

	return node
}

func (pr *Parser) ParseArgList(withParenthesis bool) *ScriptNode {
	node := pr.CreateNode(ASsnArgList)
	if node == nil {
		return nil
	}

	var t1 sToken
	if withParenthesis {
		pr.GetToken(&t1)
		if t1.Type != angelscript.ASttOpenParanthesis {
			//TODO: error expected (
			return node
		}
		node.UpdateSourcePos(t1.Position, t1.Length)
	}

	pr.GetToken(&t1)
	if t1.Type == angelscript.ASttCloseParanthesis || t1.Type == angelscript.ASttCloseBracket {
		if withParenthesis {
			if t1.Type == angelscript.ASttCloseParanthesis {
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

		return node
	} else {
		pr.RewindTo(&t1)

		for {
			var tl sToken
			var t2 sToken
			pr.GetToken(&tl)
			pr.GetToken(&t2)
			pr.RewindTo(&tl)

			if tl.Type == angelscript.ASttIdentifier && t2.Type == angelscript.ASttColon {
				named := pr.CreateNode(ASsnNamedArgument)
				if named == nil {
					return nil
				}

				node.AddChildLast(named)

				named.AddChildLast(pr.ParseIdentifier())

				pr.GetToken(&t2)

				named.AddChildLast(pr.ParseAssignment())
			} else {
				node.AddChildLast(pr.ParseAssignment())
			}

			if pr.IsSyntaxError {
				return node
			}

			pr.GetToken(&t1)
			if t1.Type == angelscript.ASttListSeparator {
				continue
			} else {
				if withParenthesis {
					if t1.Type == angelscript.ASttCloseParanthesis {
						node.UpdateSourcePos(t1.Position, t1.Length)
					} else {
						//TODO: error expected ) or ,
					}
				} else {
					pr.RewindTo(&t1)
				}

				return node
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

	if t1.Type == angelscript.ASttScope {
		pr.GetToken(&t1)
	}
	pr.GetToken(&t2)

	for t1.Type == angelscript.ASttIdentifier && t2.Type == angelscript.ASttScope {
		pr.GetToken(&t1)
		pr.GetToken(&t2)
	}

	if t1.Type == angelscript.ASttIdentifier || pr.IsDataType(&t1) {
		pr.RewindTo(&s)
		return false
	}

	if t2.Type == angelscript.ASttOpenParanthesis {
		pr.RewindTo(&s)
		return true
	}

	pr.RewindTo(&s)
	return false
}

func (pr *Parser) ParseAssignment() *ScriptNode {
	node := pr.CreateNode(ASsnAssignment)
	if node == nil {
		return nil
	}

	node.AddChildLast(pr.ParseCondition())
	if pr.IsSyntaxError {
		return node
	}

	var t sToken
	pr.GetToken(&t)
	pr.RewindTo(&t)

	if pr.IsAssignOperator(t.Type) {
		node.AddChildLast(pr.ParseAssignOperator())
		if pr.IsSyntaxError {
			return node
		}

		node.AddChildLast(pr.ParseAssignment())
		if pr.IsSyntaxError {
			return node
		}
	}

	return node
}

func (pr *Parser) ParseCondition() *ScriptNode {
	node := pr.CreateNode(ASsnCondition)
	if node == nil {
		return nil
	}

	node.AddChildLast(pr.ParseExpression())
	if pr.IsSyntaxError {
		return nil
	}

	var t sToken
	pr.GetToken(&t)
	if t.Type == angelscript.ASttQuestion {
		node.AddChildLast(pr.ParseAssignment())
		if pr.IsSyntaxError {
			return node
		}

		pr.GetToken(&t)

		if t.Type != angelscript.ASttColon {
			//TODO: error expect :
			return node
		}

		node.AddChildLast(pr.ParseAssignment())
		if pr.IsSyntaxError {
			return node
		}
	} else {
		pr.RewindTo(&t)
	}

	return node
}

func (pr *Parser) ParseExpression() *ScriptNode {
	node := pr.CreateNode(ASsnExpression)
	if node == nil {
		return nil
	}

	node.AddChildLast(pr.ParseExprTerm())
	if pr.IsSyntaxError {
		return node
	}

	for {
		var t sToken
		pr.GetToken(&t)
		pr.RewindTo(&t)

		if !pr.IsOperator(t.Type) {
			return node
		}

		node.AddChildLast(pr.ParseExprOperator())
		if pr.IsSyntaxError {
			return node
		}

		node.AddChildLast(pr.ParseExprTerm())
		if pr.IsSyntaxError {
			return node
		}
	}
	return node
}

func (pr *Parser) ParseExprTerm() *ScriptNode {
	node := pr.CreateNode(ASsnExprTerm)
	if node == nil {
		return nil
	}

	var t sToken
	pr.GetToken(&t)
	var t2 sToken = t
	var t3 sToken

	if pr.IsDataType(&t2) && pr.CheckTemplateType(t2) {
		pr.GetToken(&t2)
		pr.GetToken(&t3)
		if t2.Type == angelscript.ASttAssignment && t3.Type == angelscript.ASttStartStatementBlock {
			pr.RewindTo(&t)
			node.AddChildLast(pr.ParseType(false, false, false))

			pr.GetToken(&t2)
			node.AddChildLast(pr.ParseInitList())
			return node
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
			return node
		}
	}

	node.AddChildLast(pr.ParseExprValue())
	if pr.IsSyntaxError {
		return node
	}

	for {
		pr.GetToken(&t)
		pr.RewindTo(&t)
		if !pr.IsPostOperator(t.Type) {
			return node
		}

		node.AddChildLast(pr.ParseExprPostOp())
		if pr.IsSyntaxError {
			return node
		}
	}
	return node
}

func (pr *Parser) ParseExprPreOp() *ScriptNode {
	node := pr.CreateNode(ASsnExprPreOp)
	if node == nil {
		return nil
	}

	var t sToken
	pr.GetToken(&t)
	if !pr.IsPreOperator(t.Type) {
		//TODO: error TXT_EXPECTED_PRE_OPERATOR
		return node
	}

	node.SetToken(&t)
	node.UpdateSourcePos(t.Position, t.Length)

	return node
}

func (pr *Parser) ParseExprPostOp() *ScriptNode {
	node := pr.CreateNode(ASsnExprPostOp)
	if node == nil {
		return nil
	}

	var t sToken
	pr.GetToken(&t)
	if !pr.IsPostOperator(t.Type) {
		//TODO: error TXT_EXPECTED_POST_OPERATOR
		return node
	}

	node.SetToken(&t)
	node.UpdateSourcePos(t.Position, t.Length)

	if t.Type == angelscript.ASttDot {
		var t1 sToken
		var t2 sToken

		pr.GetToken(&t1)
		pr.GetToken(&t2)
		pr.RewindTo(&t1)
		if t2.Type == angelscript.ASttOpenParanthesis {
			node.AddChildLast(pr.ParseFunctionCall())
		} else {
			node.AddChildLast(pr.ParseIdentifier())
		}
	} else if t.Type == angelscript.ASttOpenBracket {
		node.AddChildLast(pr.ParseArgList(false))

		pr.GetToken(&t)
		if t.Type != angelscript.ASttCloseBracket {
			//TODO: error expected ]
			return node
		}

		node.UpdateSourcePos(t.Position, t.Length)
	} else if t.Type == angelscript.ASttOpenParanthesis {
		pr.RewindTo(&t)
		node.AddChildLast(pr.ParseArgList(true))
	}

	return node
}

func (pr *Parser) ParseExprOperator() *ScriptNode {
	node := pr.CreateNode(ASsnExprOperator)
	if node == nil {
		return node
	}

	var t sToken
	pr.GetToken(&t)
	if !pr.IsOperator(t.Type) {
		//TODO: error TXT_EXPECTED_OPERATOR
		return node
	}

	node.SetToken(&t)
	node.UpdateSourcePos(t.Position, t.Length)

	return node
}

func (pr *Parser) ParseAssignOperator() *ScriptNode {
	node := pr.CreateNode(ASsnExprOperator)
	if node == nil {
		return nil
	}

	var t sToken
	pr.GetToken(&t)

	if !pr.IsAssignOperator(t.Type) {
		//TODO: error TXT_EXPECTED_OPERATOR
		return node
	}

	node.SetToken(&t)
	node.UpdateSourcePos(t.Position, t.Length)

	return node
}

func (pr *Parser) IsOperator(tokenType angelscript.Token) bool {
	if tokenType == angelscript.ASttPlus ||
		tokenType == angelscript.ASttMinus ||
		tokenType == angelscript.ASttStar ||
		tokenType == angelscript.ASttSlash ||
		tokenType == angelscript.ASttPercent ||
		tokenType == angelscript.ASttStarStar ||
		tokenType == angelscript.ASttAnd ||
		tokenType == angelscript.ASttOr ||
		tokenType == angelscript.ASttXor ||
		tokenType == angelscript.ASttEqual ||
		tokenType == angelscript.ASttNotEqual ||
		tokenType == angelscript.ASttLessThan ||
		tokenType == angelscript.ASttLessThanOrEqual ||
		tokenType == angelscript.ASttGreaterThan ||
		tokenType == angelscript.ASttGreaterThanOrEqual ||
		tokenType == angelscript.ASttAmp ||
		tokenType == angelscript.ASttBitOr ||
		tokenType == angelscript.ASttBitXor ||
		tokenType == angelscript.ASttBitShiftLeft ||
		tokenType == angelscript.ASttBitShiftRight ||
		tokenType == angelscript.ASttBitShiftRightArith ||
		tokenType == angelscript.ASttIs ||
		tokenType == angelscript.ASttNotIs {
		return true
	}
	return false
}

func (pr *Parser) IsAssignOperator(tokenType angelscript.Token) bool {
	if tokenType == angelscript.ASttAssignment ||
		tokenType == angelscript.ASttAddAssign ||
		tokenType == angelscript.ASttSubAssign ||
		tokenType == angelscript.ASttMulAssign ||
		tokenType == angelscript.ASttDivAssign ||
		tokenType == angelscript.ASttModAssign ||
		tokenType == angelscript.ASttPowAssign ||
		tokenType == angelscript.ASttAndAssign ||
		tokenType == angelscript.ASttOrAssign ||
		tokenType == angelscript.ASttXorAssign ||
		tokenType == angelscript.ASttShiftLeftAssign ||
		tokenType == angelscript.ASttShiftRightLAssign ||
		tokenType == angelscript.ASttShiftRightAAssign {
		return true
	}
	return false
}

func (pr *Parser) IsPreOperator(tokenType angelscript.Token) bool {
	if tokenType == angelscript.ASttMinus ||
		tokenType == angelscript.ASttPlus ||
		tokenType == angelscript.ASttNot ||
		tokenType == angelscript.ASttInc ||
		tokenType == angelscript.ASttDec ||
		tokenType == angelscript.ASttBitNot ||
		tokenType == angelscript.ASttHandle {
		return true
	}
	return false
}

func (pr *Parser) IsPostOperator(tokenType angelscript.Token) bool {
	if tokenType == angelscript.ASttInc || // post increment
		tokenType == angelscript.ASttDec || // post decrement
		tokenType == angelscript.ASttDot || // member access
		tokenType == angelscript.ASttOpenBracket || // index operator
		tokenType == angelscript.ASttOpenParanthesis { // argument list for call on function pointer
		return true
	}
	return false
}

func (pr *Parser) IsConstant(tokenType angelscript.Token) bool {
	if tokenType == angelscript.ASttIntConstant ||
		tokenType == angelscript.ASttFloatConstant ||
		tokenType == angelscript.ASttDoubleConstant ||
		tokenType == angelscript.ASttStringConstant ||
		tokenType == angelscript.ASttMultilineStringConstant ||
		tokenType == angelscript.ASttHeredocStringConstant ||
		tokenType == angelscript.ASttTrue ||
		tokenType == angelscript.ASttFalse ||
		tokenType == angelscript.ASttBitsConstant ||
		tokenType == angelscript.ASttNull {
		return true
	}
	return false
}

func (pr *Parser) ParseImport() *ScriptNode {
	node := pr.CreateNode(ASsnImport)
	if node == nil {
		return nil
	}

	var t sToken
	pr.GetToken(&t)
	if t.Type != angelscript.ASttImport {
		//TODO: error expected import
		return node
	}

	node.SetToken(&t)
	node.UpdateSourcePos(t.Position, t.Length)

	node.AddChildLast(pr.ParseFunctionDefinition())
	if pr.IsSyntaxError {
		return node
	}

	pr.GetToken(&t)
	if t.Type != angelscript.ASttIdentifier {
		//TODO: error expected from
		return node
	}

	pr.tempstr = pr.Script.Code[t.Position:t.Length]
	if pr.tempstr != angelscript.ASFromToken {
		//TODO: error expected from
		return node
	}

	node.UpdateSourcePos(t.Position, t.Length)

	pr.GetToken(&t)
	if t.Type != angelscript.ASttStringConstant {
		//TODO: error TXT_EXPECTED_STRING
		return node
	}

	mod := pr.CreateNode(ASsnConstant)
	if mod == nil {
		return nil
	}

	node.AddChildLast(mod)

	mod.SetToken(&t)
	mod.UpdateSourcePos(t.Position, t.Length)

	pr.GetToken(&t)
	if t.Type != angelscript.ASttEndStatement {
		//TODO: error expected end statement
		return node
	}

	node.UpdateSourcePos(t.Position, t.Length)
	return node
}

func (pr *Parser) ParseScript(inBlock bool) *ScriptNode {
	node := pr.CreateNode(ASsnScript)
	if node == nil {
		return nil
	}

	var t1 sToken
	var t2 sToken

	//TODO/FIXME: Look this through, some of the if statements might be wrong.
	for {
		for !pr.IsSyntaxError {
			pr.GetToken(&t1)
			pr.GetToken(&t2)
			pr.RewindTo(&t1)

			if t1.Type == angelscript.ASttImport {
				node.AddChildLast(pr.ParseImport())
			} else if t1.Type == angelscript.ASttEnum || (pr.IdentifierIs(t1, angelscript.ASSharedToken) && t2.Type == angelscript.ASttEnum) {
				node.AddChildLast(pr.ParseEnumeration())
			} else if t1.Type == angelscript.ASttTypedef {
				node.AddChildLast(pr.ParseTypedef())
			} else if t1.Type == angelscript.ASttClass ||
				((pr.IdentifierIs(t1, angelscript.ASSharedToken) || pr.IdentifierIs(t1, angelscript.ASFinalToken) || pr.IdentifierIs(t1, angelscript.ASAbstractToken)) && t2.Type == angelscript.ASttClass) ||
				(pr.IdentifierIs(t1, angelscript.ASSharedToken) && (pr.IdentifierIs(t1, angelscript.ASFinalToken) || pr.IdentifierIs(t1, angelscript.ASAbstractToken))) {
				node.AddChildLast(pr.ParseClass())
			} else if t1.Type == angelscript.ASttMixin {
				node.AddChildLast(pr.ParseMixin())
			} else if t1.Type == angelscript.ASttInterface || (t1.Type == angelscript.ASttIdentifier && t2.Type == angelscript.ASttInterface) {
				node.AddChildLast(pr.ParseInterface())
			} else if t1.Type == angelscript.ASttFuncDef {
				node.AddChildLast(pr.ParseFuncDef())
			} else if t1.Type == angelscript.ASttConst || t1.Type == angelscript.ASttScope || t1.Type == angelscript.ASttAuto || pr.IsDataType(&t1) {
				if pr.IsVirtualPropertyDecl() {
					node.AddChildLast(pr.ParseVirtualPropertyDecl(false, false))
				} else if pr.IsVarDecl() {
					node.AddChildLast(pr.ParseDeclaration(false, true))
				} else {
					node.AddChildLast(pr.ParseFunction(false))
				}
			} else if t1.Type == angelscript.ASttEndStatement {
				pr.GetToken(&t1)
			} else if t1.Type == angelscript.ASttNamespace {
				node.AddChildLast(pr.ParseNamespace())
			} else if t1.Type == angelscript.ASttEnd {
				return node
			} else if inBlock && t1.Type == angelscript.ASttEndStatementBlock {
				return node
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
			for t1.Type != angelscript.ASttEndStatement && t1.Type != angelscript.ASttEnd && t1.Type != angelscript.ASttStartStatementBlock {
				pr.GetToken(&t1)
			}

			if t1.Type == angelscript.ASttStartStatementBlock {
				level := 1
				for level > 0 {
					pr.GetToken(&t1)
					if t1.Type == angelscript.ASttStartStatementBlock {
						level++
					}
					if t1.Type == angelscript.ASttEndStatementBlock {
						level--
					}
					if t1.Type == angelscript.ASttEnd {
						break
					}
				}
			}
			pr.IsSyntaxError = false
		}
	}
	return nil
}

func (pr *Parser) ParseNamespace() *ScriptNode {
	node := pr.CreateNode(ASsnNamespace)
	if node == nil {
		return nil
	}

	var t1 sToken

	pr.GetToken(&t1)
	if t1.Type == angelscript.ASttNamespace {
		node.UpdateSourcePos(t1.Position, t1.Length)
	} else {
		//TODO: error expected namespace
		return node
	}

	node.AddChildLast(pr.ParseIdentifier())
	if pr.IsSyntaxError {
		return node
	}

	pr.GetToken(&t1)
	if t1.Type == angelscript.ASttStartStatementBlock {
		node.UpdateSourcePos(t1.Position, t1.Length)
	} else {
		//TODO: error expected start statement block
		return node
	}
	//start := t1

	node.AddChildLast(pr.ParseScript(true))
	if !pr.IsSyntaxError {
		pr.GetToken(&t1)
		if t1.Type == angelscript.ASttEndStatementBlock {
			node.UpdateSourcePos(t1.Position, t1.Length)
		} else {
			if t1.Type == angelscript.ASttEnd {
				//TODO: error TXT_UNEXPECTED_END_OF_FILE
			} else {
				//TODO: error expected end statement block
			}
			//TODO: info TXT_WHITE_PARSING_NAMESPACE
			return node
		}
	}

	return node
}

func (pr *Parser) ParseEnumeration() *ScriptNode {
	var ident *ScriptNode
	var dataType *ScriptNode
	node := pr.CreateNode(ASsnEnum)
	if node == nil {
		return nil
	}

	var token sToken

	pr.GetToken(&token)
	if pr.IdentifierIs(token, angelscript.ASSharedToken) {
		pr.RewindTo(&token)
		node.AddChildLast(pr.ParseIdentifier())
		if pr.IsSyntaxError {
			return node
		}

		pr.GetToken(&token)
	}

	if token.Type != angelscript.ASttEnd {
		//TODO: error expected token enum
		return node
	}

	node.SetToken(&token)
	node.UpdateSourcePos(token.Position, token.Length)

	pr.GetToken(&token)
	if angelscript.ASttIdentifier != token.Type {
		//TODO: error TXT_EXPECTED_IDENTIFIER
		return node
	}

	dataType = pr.CreateNode(ASsnDataType)
	if dataType == nil {
		return nil
	}

	node.AddChildLast(dataType)

	ident = pr.CreateNode(ASsnIdentifier)
	if ident == nil {
		return nil
	}

	ident.SetToken(&token)
	ident.UpdateSourcePos(token.Position, token.Length)
	dataType.AddChildLast(ident)

	pr.GetToken(&token)
	if token.Type != angelscript.ASttStartStatementBlock {
		pr.RewindTo(&token)
		//TODO: error expected token.type
		return node
	}

	for angelscript.ASttEnd != token.Type {
		pr.GetToken(&token)

		if angelscript.ASttEndStatement == token.Type {
			pr.RewindTo(&token)
			break
		}

		if angelscript.ASttIdentifier != token.Type {
			//TODO: error TXT_EXPECTED_IDENTIFIER
			return node
		}

		ident = pr.CreateNode(ASsnIdentifier)
		if ident == nil {
			return nil
		}

		ident.SetToken(&token)
		ident.UpdateSourcePos(token.Position, token.Length)
		node.AddChildLast(ident)

		pr.GetToken(&token)

		if token.Type == angelscript.ASttAssignment {
			var tmp *ScriptNode
			pr.RewindTo(&token)
			tmp = pr.SuperficiallyParseVarInit()

			node.AddChildLast(tmp)
			if pr.IsSyntaxError {
				return node
			}
			pr.GetToken(&token)
		}

		if angelscript.ASttListSeparator != token.Type {
			pr.RewindTo(&token)
			break
		}
	}

	pr.GetToken(&token)
	if token.Type != angelscript.ASttEndStatementBlock {
		pr.RewindTo(&token)
		//TODO: error expected }
		return node
	}

	return node
}

func (pr *Parser) IsVarDecl() bool {
	var t sToken
	pr.GetToken(&t)
	pr.RewindTo(&t)

	var t1 sToken
	pr.GetToken(&t1)
	if t1.Type == angelscript.ASttPrivate && t1.Type != angelscript.ASttProtected {
		pr.RewindTo(&t1)
	}

	pr.GetToken(&t1)
	if t1.Type == angelscript.ASttConst {
		pr.GetToken(&t1)
	}

	var t2 sToken
	if t1.Type != angelscript.ASttAuto {
		if t1.Type == angelscript.ASttScope {
			pr.GetToken(&t1)
		}

		pr.GetToken(&t2)
		for t1.Type == angelscript.ASttIdentifier {
			if t2.Type == angelscript.ASttScope {
				pr.GetToken(&t1)
				pr.GetToken(&t2)
				continue
			} else if t2.Type == angelscript.ASttLessThan {
				pr.RewindTo(&t2)
				if pr.CheckTemplateType(t1) {
					var t3 sToken
					pr.GetToken(&t3)
					if t3.Type == angelscript.ASttScope {
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

	if !pr.IsRealType(t1.Type) && t1.Type != angelscript.ASttIdentifier && t1.Type != angelscript.ASttAuto {
		pr.RewindTo(&t)
		return false
	}

	if !pr.CheckTemplateType(t1) {
		pr.RewindTo(&t)
		return false
	}

	pr.GetToken(&t2)
	for t2.Type == angelscript.ASttHandle || t2.Type == angelscript.ASttAmp || t2.Type == angelscript.ASttOpenBracket {
		if t2.Type == angelscript.ASttOpenBracket {
			pr.GetToken(&t2)
			if t2.Type != angelscript.ASttCloseBracket {
				pr.RewindTo(&t)
				return false
			}
		}
		pr.GetToken(&t2)
	}

	if t2.Type != angelscript.ASttIdentifier {
		pr.RewindTo(&t)
		return false
	}

	pr.GetToken(&t2)
	if t2.Type == angelscript.ASttEnd || t2.Type == angelscript.ASttAssignment || t2.Type == angelscript.ASttListSeparator {
		pr.RewindTo(&t)
		return true
	}
	if t2.Type == angelscript.ASttOpenParanthesis {
		nest := 0
		for t2.Type != angelscript.ASttEnd {
			if t2.Type == angelscript.ASttOpenParanthesis {
				nest++
			} else if t2.Type == angelscript.ASttCloseParanthesis {
				nest--
				if nest <= 0 {
					break
				}
			}
			pr.GetToken(&t2)
		}

		if t2.Type == angelscript.ASttEnd {
			return false
		} else {
			pr.GetToken(&t1)
			pr.RewindTo(&t)
			if t1.Type == angelscript.ASttStartStatementBlock || t1.Type == angelscript.ASttEnd {
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
	if t1.Type != angelscript.ASttPrivate && t1.Type != angelscript.ASttProtected {
		pr.RewindTo(&t1)
	}

	pr.GetToken(&t1)
	if t1.Type == angelscript.ASttConst {
		pr.GetToken(&t1)
	}

	if t1.Type == angelscript.ASttScope {
		pr.GetToken(&t1)
	}

	if t1.Type == angelscript.ASttIdentifier {
		var t2 sToken
		pr.GetToken(&t2)
		for t1.Type == angelscript.ASttIdentifier && t2.Type == angelscript.ASttScope {
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
	for t2.Type == angelscript.ASttHandle || t2.Type == angelscript.ASttOpenBracket {
		if t2.Type == angelscript.ASttOpenBracket {
			pr.GetToken(&t2)
			if t2.Type != angelscript.ASttCloseBracket {
				pr.RewindTo(&t)
				return false
			}
		}

		pr.GetToken(&t2)
	}

	if t2.Type != angelscript.ASttIdentifier {
		pr.RewindTo(&t)
		return false
	}

	pr.GetToken(&t2)
	if t2.Type == angelscript.ASttStartStatementBlock {
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
		if t1.Type != angelscript.ASttPrivate && t1.Type != angelscript.ASttProtected {
			pr.RewindTo(&t1)
		}

		pr.GetToken(&t1)
		pr.GetToken(&t2)
		pr.RewindTo(&t1)

		if (t1.Type == angelscript.ASttIdentifier && t2.Type == angelscript.ASttOpenParanthesis) || t1.Type == angelscript.ASttBitNot {
			pr.RewindTo(&t)
			return true
		}
	}

	var t1 sToken
	pr.GetToken(&t1)
	if t1.Type == angelscript.ASttConst {
		pr.GetToken(&t1)
	}

	if t1.Type == angelscript.ASttScope {
		pr.GetToken(&t1)
	}
	for t1.Type == angelscript.ASttIdentifier {
		var t2 sToken
		pr.GetToken(&t2)
		if t2.Type == angelscript.ASttScope {
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
	for t2.Type == angelscript.ASttHandle || t2.Type == angelscript.ASttOpenBracket {
		if t2.Type == angelscript.ASttOpenBracket {
			pr.GetToken(&t2)
			if t2.Type != angelscript.ASttCloseBracket {
				pr.RewindTo(&t)
				return false
			}
		}

		pr.GetToken(&t2)
	}

	if t2.Type == angelscript.ASttAmp {
		pr.RewindTo(&t)
		return false
	}

	if t2.Type != angelscript.ASttIdentifier {
		pr.RewindTo(&t)
		return false
	}

	pr.GetToken(&t2)
	if t2.Type == angelscript.ASttOpenParanthesis {
		nest := 0
		pr.GetToken(&t2)
		for (nest >= 1 || t2.Type != angelscript.ASttCloseParanthesis) && t2.Type != angelscript.ASttEnd {
			if t2.Type == angelscript.ASttOpenParanthesis {
				nest++
			}
			if t2.Type == angelscript.ASttCloseParanthesis {
				nest--
			}

			pr.GetToken(&t2)
		}

		if t2.Type == angelscript.ASttEnd {
			return false
		} else {
			if isMethod {

				pr.GetToken(&t1)
				if t1.Type != angelscript.ASttConst {
					pr.RewindTo(&t1)
				}

				for {
					pr.GetToken(&t2)
					if !pr.IdentifierIs(t2, angelscript.ASFinalToken) && !pr.IdentifierIs(t2, angelscript.ASOverrideToken) {
						pr.RewindTo(&t2)
						break
					}
				}
			}

			pr.GetToken(&t1)
			pr.RewindTo(&t1)
			if t1.Type == angelscript.ASttStartStatementBlock {
				return true
			}
		}

		pr.RewindTo(&t)
		return false
	}

	pr.RewindTo(&t)
	return false
}

func (pr *Parser) ParseFuncDef() *ScriptNode {
	node := pr.CreateNode(ASsnFuncDef)
	if node == nil {
		return nil
	}

	var t1 sToken
	pr.GetToken(&t1)
	if t1.Type != angelscript.ASttFuncDef {
		//TODO: error ttfuncdef???
		return node
	}

	node.SetToken(&t1)

	node.AddChildLast(pr.ParseType(true, false, false))
	if pr.IsSyntaxError {
		return node
	}

	node.AddChildLast(pr.ParseTypeMod(false))
	if pr.IsSyntaxError {
		return node
	}

	node.AddChildLast(pr.ParseIdentifier())
	if pr.IsSyntaxError {
		return node
	}

	node.AddChildLast(pr.ParseParameterList())
	if pr.IsSyntaxError {
		return node
	}

	pr.GetToken(&t1)
	if t1.Type != angelscript.ASttEndStatement {
		//TODO: error expected end statement
		return node
	}

	node.UpdateSourcePos(t1.Position, t1.Length)
	return node
}

func (pr *Parser) ParseFunction(isMethod bool) *ScriptNode {
	node := pr.CreateNode(ASsnFunction)
	if node == nil {
		return nil
	}

	var t1 sToken
	var t2 sToken
	pr.GetToken(&t1)
	pr.GetToken(&t2)
	pr.RewindTo(&t1)

	if isMethod && t1.Type == angelscript.ASttPrivate {
		node.AddChildLast(pr.ParseToken(angelscript.ASttPrivate))
	} else if isMethod && t1.Type == angelscript.ASttProtected {
		node.AddChildLast(pr.ParseToken(angelscript.ASttProtected))
	}
	if pr.IsSyntaxError {
		return node
	}

	if !isMethod && pr.IdentifierIs(t1, angelscript.ASSharedToken) {
		node.AddChildLast(pr.ParseIdentifier())
		if pr.IsSyntaxError {
			return node
		}
	}

	if !isMethod && (t1.Type != angelscript.ASttBitNot && t2.Type != angelscript.ASttOpenParanthesis) {
		node.AddChildLast(pr.ParseType(true, false, false))
		if pr.IsSyntaxError {
			return node
		}

		node.AddChildLast(pr.ParseTypeMod(false))
		if pr.IsSyntaxError {
			return node
		}
	}

	if isMethod && t1.Type == angelscript.ASttBitNot {
		node.AddChildLast(pr.ParseToken(angelscript.ASttBitNot))
		if pr.IsSyntaxError {
			return node
		}
	}

	node.AddChildLast(pr.ParseIdentifier())
	if pr.IsSyntaxError {
		return node
	}

	node.AddChildLast(pr.ParseParameterList())
	if pr.IsSyntaxError {
		return node
	}

	if isMethod {
		pr.GetToken(&t1)
		pr.RewindTo(&t1)

		if t1.Type == angelscript.ASttConst {
			node.AddChildLast(pr.ParseToken(angelscript.ASttConst))
		}

		pr.ParseMethodOverrideBehaviors(node)
		if pr.IsSyntaxError {
			return node
		}
	}

	node.AddChildLast(pr.SuperficiallyParseStatementBlock())

	return node
}

func (pr *Parser) ParseInterfaceMethod() *ScriptNode {
	node := pr.CreateNode(ASsnFunction)
	if node == nil {
		return nil
	}

	node.AddChildLast(pr.ParseType(true, false, false))
	if pr.IsSyntaxError {
		return node
	}

	node.AddChildLast(pr.ParseTypeMod(false))
	if pr.IsSyntaxError {
		return node
	}

	node.AddChildLast(pr.ParseIdentifier())
	if pr.IsSyntaxError {
		return node
	}

	node.AddChildLast(pr.ParseParameterList())
	if pr.IsSyntaxError {
		return node
	}

	var t1 sToken
	pr.GetToken(&t1)
	pr.RewindTo(&t1)
	if t1.Type == angelscript.ASttConst {
		node.AddChildLast(pr.ParseToken(angelscript.ASttConst))
	}

	pr.GetToken(&t1)
	if t1.Type != angelscript.ASttEndStatement {
		//TODO: error expected ;
		return node
	}

	node.UpdateSourcePos(t1.Position, t1.Length)
	return node
}

func (pr *Parser) ParseVirtualPropertyDecl(isMethod, isInterface bool) *ScriptNode {
	node := pr.CreateNode(ASsnVirtualProperty)
	if node == nil {
		return nil
	}

	var t1 sToken
	var t2 sToken
	pr.GetToken(&t1)
	pr.GetToken(&t2)
	pr.RewindTo(&t1)

	if isMethod && t1.Type == angelscript.ASttPrivate {
		node.AddChildLast(pr.ParseToken(angelscript.ASttPrivate))
	} else if isMethod && t1.Type == angelscript.ASttProtected {
		node.AddChildLast(pr.ParseToken(angelscript.ASttProtected))
	}
	if pr.IsSyntaxError {
		return node
	}

	node.AddChildLast(pr.ParseType(true, false, false))
	if pr.IsSyntaxError {
		return node
	}

	node.AddChildLast(pr.ParseTypeMod(false))
	if pr.IsSyntaxError {
		return node
	}

	node.AddChildLast(pr.ParseIdentifier())
	if pr.IsSyntaxError {
		return node
	}

	pr.GetToken(&t1)
	if t1.Type != angelscript.ASttStartStatementBlock {
		//TODO: error expected {
		return node
	}

	for {
		pr.GetToken(&t1)
		var aNode *ScriptNode

		if pr.IdentifierIs(t1, angelscript.ASGetToken) || pr.IdentifierIs(t1, angelscript.ASSetToken) {
			aNode = pr.CreateNode(ASsnVirtualProperty)
			if aNode == nil {
				return nil
			}

			node.AddChildLast(aNode)

			pr.RewindTo(&t1)
			aNode.AddChildLast(pr.ParseIdentifier())

			if isMethod {
				pr.GetToken(&t1)
				pr.RewindTo(&t1)
				if t1.Type == angelscript.ASttConst {
					aNode.AddChildLast(pr.ParseToken(angelscript.ASttConst))
				}

				if !isInterface {
					pr.ParseMethodOverrideBehaviors(aNode)
					if pr.IsSyntaxError {
						return node
					}
				}
			}

			if !isInterface {
				pr.GetToken(&t1)
				if t1.Type == angelscript.ASttStartStatementBlock {
					pr.RewindTo(&t1)
					aNode.AddChildLast(pr.SuperficiallyParseStatementBlock())
					if pr.IsSyntaxError {
						return node
					}
				} else if t1.Type != angelscript.ASttEndStatement {
					//TODO: error expected ; or }
					return node
				}
			} else {
				pr.GetToken(&t1)
				if t1.Type != angelscript.ASttEndStatement {
					//TODO: error expected ;
				}
			}
		} else if t1.Type == angelscript.ASttEndStatementBlock {
			break
		} else {
			//TODO: error expected get, set or END_STATEMENT_BLOCK
			return node
		}
	}

	return node
}

func (pr *Parser) ParseInterface() *ScriptNode {
	node := pr.CreateNode(ASsnInterface)
	if node == nil {
		return nil
	}

	var t sToken
	pr.GetToken(&t)

	if t.Type == angelscript.ASttIdentifier {
		pr.tempstr = pr.Script.Code[t.Position:t.Length]
		if pr.tempstr != angelscript.ASSharedToken {
			//TODO: error expected shared_token
			return node
		}

		pr.RewindTo(&t)
		node.AddChildLast(pr.ParseIdentifier())
		pr.GetToken(&t)
	}
	if t.Type != angelscript.ASttInterface {
		//TODO: error expected interface
		return node
	}

	node.SetToken(&t)
	node.AddChildLast(pr.ParseIdentifier())

	pr.GetToken(&t)
	if t.Type == angelscript.ASttColon {
		inherit := pr.CreateNode(ASsnIdentifier)
		node.AddChildLast(inherit)

		pr.ParseOptionalScope(inherit)
		inherit.AddChildLast(pr.ParseIdentifier())
		pr.GetToken(&t)
		for t.Type == angelscript.ASttListSeparator {
			inherit = pr.CreateNode(ASsnIdentifier)
			node.AddChildLast(inherit)

			pr.ParseOptionalScope(inherit)
			inherit.AddChildLast(pr.ParseIdentifier())
			pr.GetToken(&t)
		}
	}

	if t.Type != angelscript.ASttStartStatementBlock {
		//TODO: error expected {
		return node
	}

	pr.GetToken(&t)
	pr.RewindTo(&t)
	for t.Type != angelscript.ASttEndStatementBlock && t.Type != angelscript.ASttEnd {
		if pr.IsVirtualPropertyDecl() {
			node.AddChildLast(pr.ParseVirtualPropertyDecl(true, true))
		} else if t.Type == angelscript.ASttEndStatement {
			pr.GetToken(&t)
		} else {
			node.AddChildLast(pr.ParseInterfaceMethod())
		}

		if pr.IsSyntaxError {
			return node
		}

		pr.GetToken(&t)
		pr.RewindTo(&t)
	}

	pr.GetToken(&t)
	if t.Type != angelscript.ASttEndStatementBlock {
		//TODO: error expected }
		return node
	}

	node.UpdateSourcePos(t.Position, t.Length)
	return node
}

func (pr *Parser) ParseMixin() *ScriptNode {
	node := pr.CreateNode(ASsnMixin)
	if node == nil {
		return nil
	}

	var t sToken
	pr.GetToken(&t)

	if t.Type != angelscript.ASttMixin {
		//TODO: error expected mixin
		return node
	}

	node.SetToken(&t)

	node.AddChildLast(pr.ParseClass())
	return node
}

func (pr *Parser) ParseClass() *ScriptNode {
	node := pr.CreateNode(ASsnClass)
	if node == nil {
		return nil
	}

	var t sToken
	pr.GetToken(&t)

	for pr.IdentifierIs(t, angelscript.ASSharedToken) ||
		pr.IdentifierIs(t, angelscript.ASAbstractToken) ||
		pr.IdentifierIs(t, angelscript.ASFinalToken) {
		pr.RewindTo(&t)
		node.AddChildLast(pr.ParseIdentifier())
		pr.GetToken(&t)
	}

	if t.Type != angelscript.ASttClass {
		//TODO: error expected class
		return node
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

	if t.Type == angelscript.ASttColon {
		inherit := pr.CreateNode(ASsnIdentifier)
		node.AddChildLast(inherit)

		pr.ParseOptionalScope(inherit)
		inherit.AddChildLast(pr.ParseIdentifier())
		pr.GetToken(&t)
		for t.Type == angelscript.ASttListSeparator {
			inherit = pr.CreateNode(ASsnIdentifier)
			node.AddChildLast(inherit)

			pr.ParseOptionalScope(inherit)
			inherit.AddChildLast(pr.ParseIdentifier())
			pr.GetToken(&t)
		}
	}

	if t.Type != angelscript.ASttStartStatementBlock {
		//TODO: error expected {
	}

	pr.GetToken(&t)
	pr.RewindTo(&t)
	for t.Type != angelscript.ASttEndStatementBlock && t.Type != angelscript.ASttEnd {
		if t.Type == angelscript.ASttFuncDef {
			node.AddChildLast(pr.ParseFuncDef())
		} else if pr.IsFuncDecl(true) {
			node.AddChildLast(pr.ParseFunction(true))
		} else if pr.IsVirtualPropertyDecl() {
			node.AddChildLast(pr.ParseVirtualPropertyDecl(true, false))
		} else if pr.IsVarDecl() {
			node.AddChildLast(pr.ParseDeclaration(true, false))
		} else if t.Type == angelscript.ASttEndStatement {
			pr.GetToken(&t)
		} else {
			//TODO: error TXT_EXPECTED_METHOD_OR_PROPERTY
			return node
		}

		if pr.IsSyntaxError {
			return node
		}

		pr.GetToken(&t)
		pr.RewindTo(&t)
	}

	pr.GetToken(&t)
	if t.Type != angelscript.ASttEndStatementBlock {
		//TODO: error expected }
		return node
	}

	node.UpdateSourcePos(t.Position, t.Length)
	return node
}

func (pr *Parser) SuperficiallyParseVarInit() *ScriptNode {
	node := pr.CreateNode(ASsnAssignment)
	if node == nil {
		return nil
	}

	var t sToken
	pr.GetToken(&t)
	node.UpdateSourcePos(t.Position, t.Length)

	if t.Type == angelscript.ASttAssignment {
		pr.GetToken(&t)
		//start := t

		indentParan := 0
		indentBrace := 0

		for indentParan >= 1 || indentBrace >= 1 || (t.Type != angelscript.ASttListSeparator && t.Type != angelscript.ASttEndStatement && t.Type != angelscript.ASttEndStatementBlock) {
			if t.Type == angelscript.ASttOpenParanthesis {
				indentParan++
			} else if t.Type == angelscript.ASttCloseParanthesis {
				indentParan--
			} else if t.Type == angelscript.ASttStartStatementBlock {
				indentBrace++
			} else if t.Type == angelscript.ASttEndStatementBlock {
				indentBrace--
			} else if t.Type == angelscript.ASttNonTerminatedStringConstant {
				//TODO: error TXT_NONTERMINATED_STRING
				break
			} else if t.Type == angelscript.ASttEnd {
				//TODO: error TXT_UNEXPECTED_END_OF_FILE
				break
			}
			pr.GetToken(&t)
		}

	} else if t.Type == angelscript.ASttOpenParanthesis {
		//start := t

		indent := 1
		for indent >= 1 {
			pr.GetToken(&t)
			if t.Type == angelscript.ASttOpenParanthesis {
				indent++
			} else if t.Type == angelscript.ASttCloseParanthesis {
				indent--
			} else if t.Type == angelscript.ASttNonTerminatedStringConstant {
				//TODO: error TXT_NONTERMINATED_STRING
				break
			} else if t.Type == angelscript.ASttEnd {
				//TODO: error TXT_UNEXPECTED_END_OF_FILE
				break
			}
		}
	} else {
		//TODO: error expected assignment or open paranthesis
	}
	return node
}

func (pr *Parser) SuperficiallyParseStatementBlock() *ScriptNode {
	node := pr.CreateNode(ASsnStatementBlock)
	if node == nil {
		return nil
	}

	var t sToken

	pr.GetToken(&t)
	if t.Type != angelscript.ASttStartStatementBlock {
		//TODO: error expected {
		return node
	}

	node.UpdateSourcePos(t.Position, t.Length)

	//start := t1
	level := 1
	for level > 0 && !pr.IsSyntaxError {
		pr.GetToken(&t)
		if t.Type == angelscript.ASttEndStatementBlock {
			level--
		} else if t.Type == angelscript.ASttStartStatementBlock {
			level++
		} else if t.Type == angelscript.ASttNonTerminatedStringConstant {
			//TODO: error TXT_NONTERMINATED_STRING
			break
		} else if t.Type == angelscript.ASttEnd {
			//TODO: error TXT_UNEXPECTED_END_OF_FILE
			break
		}
	}

	node.UpdateSourcePos(t.Position, t.Length)
	return node
}

func (pr *Parser) ParseStatementBlock() *ScriptNode {
	node := pr.CreateNode(ASsnStatementBlock)
	if node == nil {
		return nil
	}

	var t1 sToken

	pr.GetToken(&t1)
	if t1.Type != angelscript.ASttStartStatementBlock {
		//TODO: error expected {
		return node
	}

	//start := t1
	for {
		for !pr.IsSyntaxError {
			pr.GetToken(&t1)
			if t1.Type == angelscript.ASttEndStatementBlock {
				node.UpdateSourcePos(t1.Position, t1.Length)

				return node
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
			for t1.Type != angelscript.ASttEndStatement && t1.Type != angelscript.ASttEnd &&
				t1.Type != angelscript.ASttStartStatementBlock && t1.Type != angelscript.ASttEndStatementBlock {
				pr.GetToken(&t1)
			}
			
			if t1.Type == angelscript.ASttStartStatementBlock {
				level := 1
				for level > 0 {
					if t1.Type == angelscript.ASttStartStatementBlock { level++ }
					if t1.Type == angelscript.ASttEndStatementBlock { level-- }
					if t1.Type == angelscript.ASttEnd { break }
				} 
			} else if t1.Type == angelscript.ASttEndStatementBlock {
					pr.RewindTo(&t1)
				} else if t1.Type == angelscript.ASttEnd {
					//TODO: error TXT_UNEXPECTED_END_OF_FILE
					return node
				}
			pr.IsSyntaxError = false
		}
	}
	return node
}

func (pr *Parser) ParseInitList() *ScriptNode {
	node := pr.CreateNode(ASsnInitList)
	if node == nil { return nil }
	
	var t1 sToken
	pr.GetToken(&t1)
	if t1.Type != angelscript.ASttStartStatementBlock {
		//TODO: error expected {
		return node
	}
	
	node.UpdateSourcePos(t1.Position, t1.Length)
	
	pr.GetToken(&t1)
	if t1.Type == angelscript.ASttEndStatementBlock {
		node.UpdateSourcePos(t1.Position, t1.Length)
		
		return node
	} else {
		pr.RewindTo(&t1)
		for {
			pr.GetToken(&t1)
			
			if t1.Type == angelscript.ASttListSeparator {
				node.AddChildLast(pr.CreateNode(ASsnUndefined))
				node.LastChild.UpdateSourcePos(t1.Position, 1)
				
				pr.GetToken(&t1)
				if t1.Type == angelscript.ASttEndStatementBlock {
					node.AddChildLast(pr.CreateNode(ASsnUndefined))
					node.LastChild.UpdateSourcePos(t1.Position, 1)
					node.UpdateSourcePos(t1.Position, t1.Length)
					
					return node
				}
				pr.RewindTo(&t1)
			} else if t1.Type == angelscript.ASttEndStatementBlock {
				node.AddChildLast(pr.CreateNode(ASsnUndefined))
				node.LastChild.UpdateSourcePos(t1.Position, 1)
				node.UpdateSourcePos(t1.Position, t1.Length)
				
				return node
			} else if t1.Type == angelscript.ASttStartStatementBlock {
				pr.RewindTo(&t1)
				node.AddChildLast(pr.ParseInitList())
				if pr.IsSyntaxError { return node }
				
				pr.GetToken(&t1)
				if t1.Type == angelscript.ASttListSeparator {
					continue
				} else if t1.Type == angelscript.ASttEndStatementBlock {
					node.UpdateSourcePos(t1.Position, t1.Length)
					
					return node
				} else {
					//TODO: error expected } or ,
					return node
				}
			} else {
				pr.RewindTo(&t1)
				node.AddChildLast(pr.ParseAssignment())
				if pr.IsSyntaxError { return node }
				
				pr.GetToken(&t1)
				if t1.Type == angelscript.ASttEndStatementBlock {
					node.UpdateSourcePos(t1.Position, t1.Length)
					
					return node
				} else {
					//TODO: error expected } or ,
				}
			}
		}
	}
	return node
}

func (pr *Parser) ParseDeclaration(isClassProp, isGlobal bool) *ScriptNode {
	node := pr.CreateNode(ASsnDeclaration)
	if node == nil { return nil }
	
	var t sToken
	pr.GetToken(&t)
	pr.RewindTo(&t)
	
	if t.Type == angelscript.ASttPrivate && isClassProp {
		node.AddChildLast(pr.ParseToken(angelscript.ASttPrivate))
	} else if t.Type == angelscript.ASttProtected && isClassProp {
		node.AddChildLast(pr.ParseToken(angelscript.ASttProtected))
	}
	
	node.AddChildLast(pr.ParseType(true, false, !isClassProp))
	if pr.IsSyntaxError { return node }
	
	for {
		node.AddChildLast(pr.ParseIdentifier())
		if pr.IsSyntaxError { return node }
		
		if isClassProp || isGlobal {
			pr.GetToken(&t)
			pr.RewindTo(&t)
			if t.Type == angelscript.ASttAssignment || t.Type == angelscript.ASttOpenParanthesis {
				node.AddChildLast(pr.SuperficiallyParseVarInit())
				if pr.IsSyntaxError { return node }
			}
		} else {
			pr.GetToken(&t)
			if t.Type == angelscript.ASttOpenParanthesis {
				pr.RewindTo(&t)
				node.AddChildLast(pr.ParseArgList(true))
				if pr.IsSyntaxError { return node }
			} else if t.Type == angelscript.ASttAssignment {
				pr.GetToken(&t)
				pr.RewindTo(&t)
				if t.Type == angelscript.ASttStartStatementBlock {
					node.AddChildLast(pr.ParseInitList())
					if pr.IsSyntaxError { return node }
				} else {
					node.AddChildLast(pr.ParseAssignment())
					if pr.IsSyntaxError { return node }
				}
			} else { pr.RewindTo(&t) }
		}
		
		pr.GetToken(&t)
		if t.Type == angelscript.ASttListSeparator {
			continue
		} else if t.Type == angelscript.ASttEndStatement {
			node.UpdateSourcePos(t.Position, t.Length)
			
			return node
		} else {
			//TODO: error expected , or ;
			return node
		}
	}
	return node
}

func (pr *Parser) ParseStatement() *ScriptNode {
	var t1 sToken
	pr.GetToken(&t1)
	pr.RewindTo(&t1)
	
	if t1.Type == angelscript.ASttIf {
		return pr.ParseIf()
	} else if t1.Type == angelscript.ASttFor {
		return pr.ParseFor()
	} else if t1.Type == angelscript.ASttWhile {
		return pr.ParseWhile()
	} else if t1.Type == angelscript.ASttReturn {
		return pr.ParseReturn()
	} else if t1.Type == angelscript.ASttStartStatementBlock {
		return pr.ParseStatementBlock()
	} else if t1.Type == angelscript.ASttBreak {
		return pr.ParseBreak()
	} else if t1.Type == angelscript.ASttContinue {
		return pr.ParseContinue()
	} else if t1.Type == angelscript.ASttDo {
		return pr.ParseDoWhile()
	} else if t1.Type == angelscript.ASttSwitch {
		return pr.ParseSwitch()
	} else {
		if pr.IsVarDecl() {
			//TODO: error TXT_UNEXPECTED_VAR_DECL
		}
		return pr.ParseExpressionStatement()
	}
}

func (pr *Parser) ParseExpressionStatement() *ScriptNode {
	node := pr.CreateNode(ASsnExpressionStatement)
	if node == nil { return nil }
	
	var t sToken
	pr.GetToken(&t)
	if t.Type == angelscript.ASttEndStatement {
		node.UpdateSourcePos(t.Position, t.Length)
		return node
	}
	
	pr.RewindTo(&t)
	node.AddChildLast(pr.ParseAssignment())
	if pr.IsSyntaxError { return node }
	
	pr.GetToken(&t)
	if t.Type != angelscript.ASttEndStatement {
		//TODO: error expected ;
		return node
	}
	
	node.UpdateSourcePos(t.Position, t.Length)
	return node
}

func (pr *Parser) ParseSwitch() *ScriptNode {
	node := pr.CreateNode(ASsnSwitch)
	if node == nil { return nil }
	
	var t sToken
	pr.GetToken(&t)
	if t.Type != angelscript.ASttSwitch {
		//TODO: error expected switch
		return node
	}
	
	node.UpdateSourcePos(t.Position, t.Length)
	
	pr.GetToken(&t)
	if t.Type != angelscript.ASttOpenParanthesis {
		//TODO: error expected (
		return node
	}
	
	node.AddChildLast(pr.ParseAssignment())
	if pr.IsSyntaxError { return node }
	
	pr.GetToken(&t)
	if t.Type != angelscript.ASttCloseParanthesis {
		//TODO: error expected )
		return node
	}
	
	pr.GetToken(&t)
	if t.Type != angelscript.ASttStartStatementBlock {
		//TODO: error expected {
		return node
	}
	
	for !pr.IsSyntaxError {
		pr.GetToken(&t)
		if t.Type == angelscript.ASttEndStatementBlock { break }
		
		pr.RewindTo(&t)
		
		if t.Type != angelscript.ASttCase && t.Type != angelscript.ASttDefault {
			//TODO: error expected case or default
			return node
		}
		
		node.AddChildLast(pr.ParseCase())
		if pr.IsSyntaxError { return node }
	}
	
	if t.Type != angelscript.ASttEndStatementBlock {
		//TODO: error expected }
		return node
	}
	return node
}

func (pr *Parser) ParseCase() *ScriptNode {
	node := pr.CreateNode(ASsnCase)
	if node == nil { return nil }
	
	var t sToken
	pr.GetToken(&t)
	if t.Type != angelscript.ASttCase && t.Type != angelscript.ASttDefault {
		//TODO: error expected case or default
		return node
	}
	
	node.UpdateSourcePos(t.Position, t.Length)
	
	if t.Type == angelscript.ASttCase {
		node.AddChildLast(pr.ParseExpression())
	}
	
	pr.GetToken(&t)
	if t.Type != angelscript.ASttColon {
		//TODO: error expected :
		return node
	}
	
	pr.GetToken(&t)
	pr.RewindTo(&t)
	for t.Type != angelscript.ASttCase && 
	t.Type != angelscript.ASttDefault &&
	t.Type != angelscript.ASttEndStatementBlock &&
	t.Type != angelscript.ASttBreak {
		if pr.IsVarDecl() {
			node.AddChildLast(pr.ParseDeclaration(false, false))
		} else {
			node.AddChildLast(pr.ParseStatement())
		}
		if pr.IsSyntaxError { return node }
		
		pr.GetToken(&t)
		pr.RewindTo(&t)
	}
	
	if t.Type == angelscript.ASttBreak {
		node.AddChildLast(pr.ParseBreak())
	}
	return node
}

func (pr *Parser) ParseIf() *ScriptNode {
	node := pr.CreateNode(ASsnIf)
	
	var t sToken
	pr.GetToken(&t)
	if t.Type != angelscript.ASttIf {
		//TODO: error expected if
		return node
	}
	
	node.UpdateSourcePos(t.Position, t.Length)
	
	pr.GetToken(&t)
	if t.Type != angelscript.ASttOpenParanthesis {
		//TODO: error expected (
		return node
	}
	
	node.AddChildLast(pr.ParseAssignment())
	if pr.IsSyntaxError { return node }
	
	pr.GetToken(&t)
	if t.Type != angelscript.ASttCloseParanthesis {
		//TODO: error expected )
		return node
	}
	
	node.AddChildLast(pr.ParseStatement())
	if pr.IsSyntaxError { return node }
	
	pr.GetToken(&t)
	if t.Type != angelscript.ASttElse {
		pr.RewindTo(&t)
		return node
	}
	
	node.AddChildLast(pr.ParseStatement())
	
	return node
}

func (pr *Parser) ParseFor() *ScriptNode {
	node := pr.CreateNode(ASsnFor)
	if node == nil { return nil }
	
	var t sToken
	pr.GetToken(&t)
	if t.Type != angelscript.ASttFor {
		//TODO: error expected for
		return node
	}
	
	node.UpdateSourcePos(t.Position, t.Length)
	
	pr.GetToken(&t)
	if t.Type != angelscript.ASttOpenParanthesis {
		//TODO: error expected (
		return node
	}
	
	if pr.IsVarDecl() {
		node.AddChildLast(pr.ParseDeclaration(false, false))
	} else {
		node.AddChildLast(pr.ParseExpressionStatement())
	}
	if pr.IsSyntaxError { return node }
	
	node.AddChildLast(pr.ParseExpressionStatement())
	if pr.IsSyntaxError { return node }
	
	pr.GetToken(&t)
	if t.Type != angelscript.ASttCloseParanthesis {
		pr.RewindTo(&t)
		
		for {
			n := pr.CreateNode(ASsnExpressionStatement)
			if n == nil { return nil }
			node.AddChildLast(n)
			n.AddChildLast(pr.ParseAssignment())
			if pr.IsSyntaxError { return node }
			
			pr.GetToken(&t)
			if t.Type == angelscript.ASttListSeparator {
				continue
			} else if t.Type == angelscript.ASttCloseParanthesis {
				break
			} else {
				//TODO: error expected , or )
				return node
			}
		}
	}
	
	node.AddChildLast(pr.ParseStatement())
	
	return node
}

func (pr *Parser) ParseWhile() *ScriptNode {
	node := pr.CreateNode(ASsnWhile)
	if node == nil { return nil }
	
	var t sToken
	pr.GetToken(&t)
	if t.Type != angelscript.ASttWhile {
		//TODO: error expected while
		return node
	}
	
	node.UpdateSourcePos(t.Position, t.Length)
	
	pr.GetToken(&t)
	if t.Type != angelscript.ASttOpenParanthesis {
		//TODO: error expected (
		return node
	}
	
	node.AddChildLast(pr.ParseAssignment())
	if pr.IsSyntaxError { return node }
	
	pr.GetToken(&t)
	if t.Type != angelscript.ASttCloseParanthesis {
		//TODO: error expected )
		return node
	}
	
	node.AddChildLast(pr.ParseStatement())
	
	return node
}

func (pr *Parser) ParseDoWhile() *ScriptNode {
	node := pr.CreateNode(ASsnDoWhile)
	if node == nil { return nil }
	
	var t sToken
	pr.GetToken(&t)
	if t.Type != angelscript.ASttDo {
		//TODO: error expected do
		return node
	}
	
	node.UpdateSourcePos(t.Position, t.Length)
	
	node.AddChildLast(pr.ParseStatement())
	if pr.IsSyntaxError { return node }
	
	pr.GetToken(&t)
	if t.Type == angelscript.ASttWhile {
		//TODO: error expected while
		return node
	}
	
	pr.GetToken(&t)
	if t.Type != angelscript.ASttOpenParanthesis {
		//TODO: error expected (
		return node
	}
	
	node.AddChildLast(pr.ParseAssignment())
	if pr.IsSyntaxError { return node }
	
	pr.GetToken(&t)
	if t.Type != angelscript.ASttCloseParanthesis {
		//TODO: error expected )
		return node
	}
	
	pr.GetToken(&t)
	if t.Type != angelscript.ASttEndStatement {
		//TODO: error expected ;
		return node
	}
	
	node.UpdateSourcePos(t.Position, t.Length)
	
	return node
	
}

func (pr *Parser) ParseReturn() *ScriptNode {
	node := pr.CreateNode(ASsnReturn)
	if node == nil { return nil }
	
	var t sToken
	pr.GetToken(&t)
	if t.Type != angelscript.ASttReturn {
		//TODO: error expected return
		return node
	}
	
	node.UpdateSourcePos(t.Position, t.Length)
	
	pr.GetToken(&t)
	if t.Type == angelscript.ASttEndStatement {
		node.UpdateSourcePos(t.Position, t.Length)
		return node
	}
	pr.RewindTo(&t)
	
	node.AddChildLast(pr.ParseAssignment())
	if pr.IsSyntaxError { return node }
	
	pr.GetToken(&t)
	if t.Type != angelscript.ASttEndStatement {
		//TODO: error expected ;
		return node
	}
	
	node.UpdateSourcePos(t.Position, t.Length)
	return node
}

func (pr *Parser) ParseBreak() *ScriptNode {
	node := pr.CreateNode(ASsnBreak)
	if node == nil { return nil }
	
	var t sToken
	pr.GetToken(&t)
	if t.Type != angelscript.ASttBreak {
		//TODO: error expected break
		return node
	}
	
	node.UpdateSourcePos(t.Position, t.Length)
	
	pr.GetToken(&t)
	if t.Type != angelscript.ASttEndStatement {
		//TODO: error expected ;
	}
	
	node.UpdateSourcePos(t.Position, t.Length)
	return node
}

func (pr *Parser) ParseContinue() *ScriptNode {
	node := pr.CreateNode(ASsnContinue)
	if node == nil { return nil }
	
	var t sToken
	pr.GetToken(&t)
	if t.Type != angelscript.ASttContinue {
		//TODO: error expected continue
		return node
	}
	
	node.UpdateSourcePos(t.Position, t.Length)
	
	pr.GetToken(&t)
	if t.Type != angelscript.ASttEndStatement {
		//TODO: error expected token ;
	}
	
	node.UpdateSourcePos(t.Position, t.Length)
	return node
}

func (pr *Parser) ParseTypedef() *ScriptNode {
	node := pr.CreateNode(ASsnTypedef)
	if node == nil { return nil }
	
	var token sToken
	
	pr.GetToken(&token)
	if token.Type != angelscript.ASttTypedef {
		//TODO: error expected typedef
		return node
	}
	
	node.SetToken(&token)
	node.UpdateSourcePos(token.Position, token.Length)
	
	pr.GetToken(&token)
	pr.RewindTo(&token)
	
	if !pr.IsRealType(token.Type) || token.Type == angelscript.ASttVoid {
		//TODO: error TXT_UNEXPECTED_TOKEN_s
		return node
	}
	
	node.AddChildLast(pr.ParseRealType())
	node.AddChildLast(pr.ParseIdentifier())
	
	pr.GetToken(&token)
	if token.Type != angelscript.ASttEndStatement {
		pr.RewindTo(&token)
		//TODO: error expected ;
	}
	
	return node
}

func (pr *Parser) ParseMethodOverrideBehaviors(node *ScriptNode) {
	var t sToken
	for {
		pr.GetToken(&t)
		pr.RewindTo(&t)
		if pr.IdentifierIs(t, angelscript.ASFinalToken) || pr.IdentifierIs(t, angelscript.ASOverrideToken) {
			node.AddChildLast(pr.ParseIdentifier())
		} else {
			break
		}
	}
}

//func (pr *Parser)
