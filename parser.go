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

	tempstr   string
	lastToken sToken
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

func (pr *Parser) ParseFunctionDefinition(script *ScriptCode, expectListPattern bool) int {
	pr.Reset()

	pr.IsParsingAppInterface = true

	pr.Script = script
	pr.Node = pr.ParseFuncDefinition()

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

func (pr *Parser) ParseDataTypeScript(script *ScriptCode, isReturnType bool) int {
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

func (pr *Parser) ParsePropertyDeclaration(script *ScriptCode) int {
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

func (pr *Parser) ParseFuncDefinition() *ScriptNode {
	node := pr.CreateNode(ASsnFunction)
	if node == nil { return nil }
	
	node.AddChildLast(pr.ParseType(true, false, false))
	if pr.IsSyntaxError { return node }
	
	node.AddChildLast(pr.ParseTypeMod(false))
	if pr.IsSyntaxError { return node }
	
	pr.ParseOptionalScope(node)
	
	node.AddChildLast(pr.ParseIdentifier())
	if pr.IsSyntaxError { return node }
	
	node.AddChildLast(pr.ParseParameterList())
	if pr.IsSyntaxError { return node }
	
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
	if node == nil { return nil }
	
	var t sToken
	
	pr.GetToken(&t)
	pr.RewindTo(&t)
	
	if t.Type == angelscript.ASttAmp {
		node.AddChildLast(pr.ParseToken(angelscript.ASttAmp))
		if pr.IsSyntaxError { return node }
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
		if pr.IsSyntaxError { return node }
	}
	
	pr.GetToken(&t)
	pr.RewindTo(&t)
	if pr.IdentifierIs(t, angelscript.ASIfHandleToken) {
		node.AddChildLast(pr.ParseToken(angelscript.ASttIdentifier))
		if pr.IsSyntaxError { return node }
	}
	
	return node
}

func (pr *Parser) ParseType(allowConst, allowVariableType, allowAuto bool) *ScriptNode {
	node := pr.CreateNode(ASsnDataType)
	if node == nil { return nil }
	
	var t sToken
	
	if allowConst {
		pr.GetToken(&t)
		pr.RewindTo(&t)
		if t.Type == angelscript.ASttConst {
			node.AddChildLast(pr.ParseToken(angelscript.ASttConst))
			if pr.IsSyntaxError { return node }
		}
	}
	
	pr.ParseOptionalScope(node)
	
	node.AddChildLast(pr.ParseDataType(allowVariableType, allowAuto))
	if pr.IsSyntaxError { return node }
	
	pr.GetToken(&t)
	pr.RewindTo(&t)
	tr := node.LastChild
	
	pr.tempstr = pr.Script[tr.TokenPosition:tr.TokenLength]
	if pr.Engine.IsTemplateType(pr.tempstr) && t.Type == angelscript.ASttLessThan {
		pr.ParseTemplTypeList(node, true)
		if pr.IsSyntaxError { return node }
	}
	
	pr.GetToken(&t)
	pr.RewindTo(&t)
	for t.Type == angelscript.ASttOpenBracket || t.Type == angelscript.ASttHandle {
		if t.Type == angelscript.ASttOpenBracket {
			node.AddChildLast(pr.ParseToken(angelscript.ASttOpenBracket))
			if pr.IsSyntaxError { return node }
			
			pr.GetToken(&t)
			if t.Type == angelscript.ASttCloseBracket {
				//TODO: ERROR (expect ])
				return node
			}
		} else {
			node.AddChildLast(pr.ParseToken(angelscript.ASttHandle))
			if pr.IsSyntaxError { return node }
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
		if reuqired {
			//TODO: ERROR (expect Lessthan)
		}
		return false
	}
	
	node.AddChildLast(pr.ParseType(true, false, false))
	if pr.IsSyntaxError { return false }
	
	pr.GetToken(&t)
	
	for t.Type == angelscript.ASttListSeparator {
		node.AddChildLast(pr.ParseType(true, false, false))
		if pr.IsSyntaxError { return false }
		pr.GetToken(&t)
	}
	
	if pr.Script.Code[t.Position:1] != ">" {
		if required {
			//TODO: ERROR (Expect GreaterThan)
		} else {
			isValid = false
		}
	} else {
		pr.SetPos(t.Position+1)
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
	if node == nil { return nil }
	
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
	if node == nil { return nil }
	
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
	if node == nil { return nil }
	
	var t1 sToken
	pr.GetToken(&t1)
	if !pr.IsDataType(t1) && !(allowVariableType && t1.Type == angelscript.ASttQuestion) && !(allowAuto && t1.Type == angelscript.ASttAuto) {
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
	if node == nil { return nil }
	
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
	if node == nil { return nil }
	
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
	if node == nil { return nil }
	
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
			if pr.IsSyntaxError { return node }
			
			node.AddChildLast(pr.ParseTypeMod(true))
			if pr.IsSyntaxError { return node }
			
			pr.GetToken(&t1)
			if t1.Type == angelscript.ASttIdentifier {
				pr.RewindTo(&t1)
				
				node.AddChildLast(pr.ParseIdentifier())
				if pr.IsSyntaxError { return node }
				
				pr.GetToken(&t1)
			}
			
			if t1.Type == angelscript.ASttAssignment {
				node.AddChildLast(pr.SuperficiallyParseExpression())
				if pr.IsSyntaxError { return node }
				
				pr.GetTokenetToken(&t1)
			}
			
			if t1.Type == angelscript.ASttCloseParanthesis {
				node.UpdateSourcePos(t1.Position, t1.Length)
			} else if  t1.Type == angelscript.ASttListSeparator {
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
	if node == nil { return nil }
	
	var start sToken
	pr.GetToken(&start)
	pr.RewindTo(&start)
	
	stack = ""
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
			if stack == "" || stack[len(stack)-1:1] == "{"  {
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
		} else if t.Type == angelscript.ASttEndStatement{
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
	if pr.lastToken.pos == pr.sourcePos {
		token = pr.lastToken
		pr.sourcePos += token.Length

		if token.Type == ASttWhiteSpace ||
			token.Type == ASttOnelineComment ||
			token.Type == ttMultilineComment {
			pr.GetToken(token)
		}
		return
	}

	sl := pr.Script.Length
	for token.Type == ASttWhiteSpace ||
		token.Type == ASttOnelineComment ||
		token.Type == ttMultilineComment {
		if pr.sourcePos >= sl {
			token.Type = ASttEnd
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
	pr.lastToken = -1
	pr.sourcePos = pos
}

func (pr *Parser) Error(text string) {

}

func (pr *Parser) Warning(text string) {

}

func (pr *Parser) Info(text string) {

}

func (pr *Parser) IsRealType(tokenType angelscript.Token) bool {
	if tokenType == angelscript.ASttVoid  ||
	   tokenType == angelscript.ASttInt   ||
	   tokenType == angelscript.ASttInt8  ||
	   tokenType == angelscript.ASttInt16 ||
	   tokenType == angelscript.ASttInt64 ||
	   tokenType == angelscript.ASttUInt  ||
	   tokenType == angelscript.ASttUInt8 ||
	   tokenType == angelscript.ASttUInt16||
	   tokenType == angelscript.ASttUInt64||
	   tokenType == angelscript.ASttFloat ||
	   tokenType == angelscript.ASttDouble||
	   tokenType == angelscript.ASttBool {
			return true
	   }
	   return false
}

func (pr *Parser) IsDataType(token *sToken) bool {
	if token.Type == angelscript.ASttIdentifier {
		if pr.CheckValidTypes {
			pr.tempstr = pr.Script.Code[token.Position:token.Length]
			if !pr.Builder.DoesTypeExist(pr.tempstr) { return false }
		}
		return true
	}
	if pr.IsRealType(token.Type) { return true }
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
	if node == nil { return nil }
	
	var t1 sToken
	pr.GetToken(&t1)
	if t1.Type != angelscript.ASttStartStatementBlock {
		//TODO: ERROR (Expected "{")
		return node
	}
	
	node.UpdateSourcePos(t1.Position, t1.Length)
	var start sToken
	
	isBeginning := true
	afterType := false
	
	for !pr.IsSyntaxError {
		pr.GetToken(&t1)
		if t1.Type == angelscript.ASttEndStatementBlock {
			if !afterType {
				//TODO: ERROR TXT_EXPECTED_DATA_TYPE
			}
			break;
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
	if t.Type != angelscript.ASttIdentifier { return false }
	return pr.Script.TokenEquals(t.Position, t.Length, str)
}


func (pr *Parser) ParseTemplateDecl(script *ScriptCode) int {
	pr.Reset()

	pr.Script = script
	pr.Node = pr.CreateNode(ASsnDataType)
	if pr.Node == nil {
		return -1
	}

	pr.Node.AddChildLast(ParseType(true))
	if pr.IsSyntaxError {
		return -1
	}

	if !pr.IsSyntaxError {
		var t sToken
		pr.GetToken(&t)
		if t.Type != ASttEnd {
			//TODO: ERROR
		}
	}

	if pr.ErrorWhileParsing {
		return -1
	}

	return 0
}

func (pr *Parser) ParseScript(script *ScriptCode) int {
	pr.Reset()

	pr.Script = script
	pr.Node = CreateNode(ASsnDataType)
	if pr.Node == nil {
		return -1
	}

	pr.Node.AddChildLast(ParseType(true))
	if pr.IsSyntaxError {
		return -1
	}

	if !pr.IsSyntaxError {
		var t sToken
		pr.GetToken(&t)
		if t.Type != ASttEnd {
			//TODO: ERROR
		}
	}

	if pr.ErrorWhileParsing {
		return -1
	}

	return 0
}

func (pr *Parser) ParseStatementBlock(script *ScriptCode, block *ScriptNode) {
	pr.Reset()

	pr.Script = script
	pr.Node = CreateNode(ASsnDataType)
	if pr.Node == nil {
		return -1
	}

	pr.Node.AddChildLast(ParseType(true))
	if pr.IsSyntaxError {
		return -1
	}

	if !pr.IsSyntaxError {
		var t sToken
		pr.GetToken(&t)
		if t.Type != ASttEnd {
			//TODO: ERROR
		}
	}

	if pr.ErrorWhileParsing {
		return -1
	}

	return 0
}

func (pr *Parser) ParseVarInit(script *ScriptCode, init *ScriptNode) int {
	pr.Reset()

	pr.Script = script
	pr.Node = CreateNode(ASsnDataType)
	if pr.Node == nil {
		return -1
	}

	pr.Node.AddChildLast(ParseType(true))
	if pr.IsSyntaxError {
		return -1
	}

	if !pr.IsSyntaxError {
		var t sToken
		pr.GetToken(&t)
		if t.Type != ASttEnd {
			//TODO: ERROR
		}
	}

	if pr.ErrorWhileParsing {
		return -1
	}

	return 0
}

func (pr *Parser) ParseExpression(script *ScriptCode) int {

	if !pr.IsSyntaxError {
		var t sToken
		pr.GetToken(&t)
		if t.Type != ASttEnd {
			//TODO: ERROR
		}
	}

	if pr.ErrorWhileParsing {
		return -1
	}

	return 0
}


//func (pr *Parser)
