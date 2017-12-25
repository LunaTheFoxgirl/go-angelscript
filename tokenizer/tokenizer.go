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
package tokens

import (
	_ "encoding/binary"
	"github.com/Member1221/go-angelscript/flags"
	"strconv"
)

//Token definition

type Token uint32

const (
	ASttUnrecognizedToken = Token(iota)

	ASttEnd // End of file

	// White space and comments
	ASttWhiteSpace       // ' ', '\t', '\r', '\n', UTF8 byte-order-mark
	ASttOnelineComment   // // \n
	ASttMultilineComment // /* */

	// Atoms
	ASttIdentifier                  // abc123
	ASttIntConstant                 // 1234
	ASttFloatConstant               // 12.34e56f
	ASttDoubleConstant              // 12.34e56
	ASttStringConstant              // "123"
	ASttMultilineStringConstant     //
	ASttHeredocStringConstant       // """text"""
	ASttNonTerminatedStringConstant // "123
	ASttBitsConstant                // 0xFFFF

	// Math operators
	ASttPlus     // +
	ASttMinus    // -
	ASttStar     // *
	ASttSlash    // /
	ASttPercent  // %
	ASttStarStar // **

	ASttHandle // @

	ASttAddAssign // +=
	ASttSubAssign // -=
	ASttMulAssign // *=
	ASttDivAssign // /=
	ASttModAssign // %=
	ASttPowAssign // **=

	ASttOrAssign          // |=
	ASttAndAssign         // &=
	ASttXorAssign         // ^=
	ASttShiftLeftAssign   // <<=
	ASttShiftRightLAssign // >>=
	ASttShiftRightAAssign // >>>=

	ASttInc // ++
	ASttDec // --

	ASttDot   // .
	ASttScope // ::

	// Statement tokens
	ASttAssignment          // =
	ASttEndStatement        // ;
	ASttListSeparator       // ,
	ASttStartStatementBlock // {
	ASttEndStatementBlock   // }
	ASttOpenParanthesis     // (
	ASttCloseParanthesis    // )
	ASttOpenBracket         // [
	ASttCloseBracket        // ]
	ASttAmp                 // &

	// Bitwise operators
	ASttBitOr              // |
	ASttBitNot             // ~
	ASttBitXor             // ^
	ASttBitShiftLeft       // <<
	ASttBitShiftRight      // >>     // TODO: In Java this is the arithmetical shift
	ASttBitShiftRightArith // >>>    // TODO: In Java this is the logical shift

	// Compare operators
	ASttEqual              // ==
	ASttNotEqual           // !=
	ASttLessThan           // <
	ASttGreaterThan        // >
	ASttLessThanOrEqual    // <=
	ASttGreaterThanOrEqual // >=

	ASttQuestion // ?
	ASttColon    // :

	// Reserved keywords
	ASttIf        // if
	ASttElse      // else
	ASttFor       // for
	ASttWhile     // while
	ASttBool      // bool
	ASttFuncDef   // funcdef
	ASttImport    // import
	ASttInt       // int
	ASttInt8      // int8
	ASttInt16     // int16
	ASttInt64     // int64
	ASttInterface // interface
	ASttIs        // is
	ASttNotIs     // !is
	ASttUInt      // uint
	ASttUInt8     // uint8
	ASttUInt16    // uint16
	ASttUInt64    // uint64
	ASttFloat     // float
	ASttVoid      // void
	ASttTrue      // true
	ASttFalse     // false
	ASttReturn    // return
	ASttNot       // not
	ASttAnd       // and, &&
	ASttOr        // or, ||
	ASttXor       // xor, ^^
	ASttBreak     // break
	ASttContinue  // continue
	ASttConst     // const
	ASttDo        // do
	ASttDouble    // double
	ASttSwitch    // switch
	ASttCase      // case
	ASttDefault   // default
	ASttIn        // in
	ASttOut       // out
	ASttInOut     // inout
	ASttNull      // null
	ASttClass     // class
	ASttTypedef   // typedef
	ASttEnum      // enum
	ASttCast      // cast
	ASttPrivate   // private
	ASttProtected // protected
	ASttNamespace // namespace
	ASttMixin     // mixin
	ASttAuto      // auto
)

const (
	ASWhitespaceToken = " \t\r\n"
	ASThisToken       = "this"
	ASFromToken       = "from"
	ASSuperToken      = "super"
	ASSharedToken     = "shared"
	ASFinalToken      = "final"
	ASOverrideToken   = "override"
	ASGetToken        = "get"
	ASSetToken        = "set"
	ASAbstractToken   = "abstract"
	ASFunctionToken   = "function"
	ASIfHandleToken   = "if_handle_then_const"
)

type TokenWord struct {
	Word   string
	Length int
	Type   Token
}

func asTokenDef(str string, tok Token) TokenWord {
	return TokenWord{
		Word:   str,
		Length: len(str),
		Type:   tok,
	}
}

var tokenWords []TokenWord = []TokenWord{
	asTokenDef("+", ASttPlus),
	asTokenDef("+=", ASttAddAssign),
	asTokenDef("++", ASttInc),
	asTokenDef("-", ASttMinus),
	asTokenDef("-=", ASttSubAssign),
	asTokenDef("--", ASttDec),
	asTokenDef("*", ASttStar),
	asTokenDef("*=", ASttMulAssign),
	asTokenDef("/", ASttSlash),
	asTokenDef("/=", ASttDivAssign),
	asTokenDef("%", ASttPercent),
	asTokenDef("%=", ASttModAssign),
	asTokenDef("**", ASttStarStar),
	asTokenDef("**=", ASttPowAssign),
	asTokenDef("=", ASttAssignment),
	asTokenDef("==", ASttEqual),
	asTokenDef(".", ASttDot),
	asTokenDef("|", ASttBitOr),
	asTokenDef("|=", ASttOrAssign),
	asTokenDef("||", ASttOr),
	asTokenDef("&", ASttAmp),
	asTokenDef("&=", ASttAndAssign),
	asTokenDef("&&", ASttAnd),
	asTokenDef("^", ASttBitXor),
	asTokenDef("^=", ASttXorAssign),
	asTokenDef("^^", ASttXor),
	asTokenDef("<", ASttLessThan),
	asTokenDef("<=", ASttLessThanOrEqual),
	asTokenDef("<<", ASttBitShiftLeft),
	asTokenDef("<<=", ASttShiftLeftAssign),
	asTokenDef(">", ASttGreaterThan),
	asTokenDef(">=", ASttGreaterThanOrEqual),
	asTokenDef(">>", ASttBitShiftRight),
	asTokenDef(">>=", ASttShiftRightLAssign),
	asTokenDef(">>>", ASttBitShiftRightArith),
	asTokenDef(">>>=", ASttShiftRightAAssign),
	asTokenDef("~", ASttBitNot),
	asTokenDef(";", ASttEndStatement),
	asTokenDef(",", ASttListSeparator),
	asTokenDef("{", ASttStartStatementBlock),
	asTokenDef("}", ASttEndStatementBlock),
	asTokenDef("(", ASttOpenParanthesis),
	asTokenDef(")", ASttCloseParanthesis),
	asTokenDef("[", ASttOpenBracket),
	asTokenDef("]", ASttCloseBracket),
	asTokenDef("?", ASttQuestion),
	asTokenDef(":", ASttColon),
	asTokenDef("::", ASttScope),
	asTokenDef("!", ASttNot),
	asTokenDef("!=", ASttNotEqual),
	asTokenDef("!is", ASttNotIs),
	asTokenDef("@", ASttHandle),
	asTokenDef("and", ASttAnd),
	asTokenDef("auto", ASttAuto),
	asTokenDef("bool", ASttBool),
	asTokenDef("break", ASttBreak),
	asTokenDef("case", ASttCase),
	asTokenDef("cast", ASttCast),
	asTokenDef("class", ASttClass),
	asTokenDef("const", ASttConst),
	asTokenDef("continue", ASttContinue),
	asTokenDef("default", ASttDefault),
	asTokenDef("do", ASttDo),
	asTokenDef("double", ASttDouble),
	asTokenDef("else", ASttElse),
	asTokenDef("enum", ASttEnum),
	asTokenDef("false", ASttFalse),
	asTokenDef("float", ASttFloat),
	asTokenDef("for", ASttFor),
	asTokenDef("funcdef", ASttFuncDef),
	asTokenDef("if", ASttIf),
	asTokenDef("import", ASttImport),
	asTokenDef("in", ASttIn),
	asTokenDef("inout", ASttInOut),
	asTokenDef("int", ASttInt),
	asTokenDef("int8", ASttInt8),
	asTokenDef("int16", ASttInt16),
	asTokenDef("int32", ASttInt),
	asTokenDef("int64", ASttInt64),
	asTokenDef("interface", ASttInterface),
	asTokenDef("is", ASttIs),
	asTokenDef("mixin", ASttMixin),
	asTokenDef("namespace", ASttNamespace),
	asTokenDef("not", ASttNot),
	asTokenDef("null", ASttNull),
	asTokenDef("or", ASttOr),
	asTokenDef("out", ASttOut),
	asTokenDef("private", ASttPrivate),
	asTokenDef("protected", ASttProtected),
	asTokenDef("return", ASttReturn),
	asTokenDef("switch", ASttSwitch),
	asTokenDef("true", ASttTrue),
	asTokenDef("typedef", ASttTypedef),
	asTokenDef("uint", ASttUInt),
	asTokenDef("uint8", ASttUInt8),
	asTokenDef("uint16", ASttUInt16),
	asTokenDef("uint32", ASttUInt),
	asTokenDef("uint64", ASttUInt64),
	asTokenDef("void", ASttVoid),
	asTokenDef("while", ASttWhile),
	asTokenDef("xor", ASttXor),
}

func GetDefinitionOrList(toks []Token) string {
	if len(toks) == 1 {
		return GetDefinition(toks[0])
	}
	o := ""
	for i, t := range toks {
		if i == 0 {
			if len(toks) == 2 {
				o += GetDefinition(t) + " or "
				continue
			}
			o += GetDefinition(t) + ", "
		} else if i == len(toks)-1 {
			o += GetDefinition(t)
		} else {
			if i == len(toks) - 2 {
				o += GetDefinition(t) + " or "
				continue
			}
			o += GetDefinition(t) + ", "
		}
	}
	return o
}

func GetDefinition(tok Token) string {
	for _, t := range tokenWords {
		if t.Type == tok {
			return t.Word
		}
	}

	if tok == ASttWhiteSpace {
		return ""
	}
	if tok == ASttIdentifier {
		return "<identifier>"
	}

	if tok == ASttEnd {
		return "<EOF>"
	}

	if tok == ASttOnelineComment {
		return "<comment>"
	}

	if tok == ASttMultilineComment {
		return "<multiline comment>"
	}

	if tok == ASttClass {
		return "class"
	}
	
	if tok == ASttIntConstant {
		return "int"
	}

	if tok == ASttFloatConstant {
		return "float"
	}

	if tok == ASttDoubleConstant {
		return "double"
	}

	if tok == ASttStringConstant {
		return "<string>"
	}

	if tok == ASttMultilineStringConstant {
		return "<multiline string>"
	}

	if tok == ASttNonTerminatedStringConstant {
		return "<nonterminated string>"
	}

	if tok == ASttBitsConstant {
		return "bits"
	}

	if tok == ASttHeredocStringConstant {
		return "<heredoc comment>"
	}

	return "unknown token (" + strconv.Itoa(int(tok)) + ")"
}

type Tokenizer struct {
	engine       int
	keywordTable map[string][]TokenWord
}

func NewTokenizer() *Tokenizer {
	tk := Tokenizer{}
	tk.engine = 0
	tk.keywordTable = make(map[string][]TokenWord)
	for i := 0; i < len(tokenWords); i++ {
		current := tokenWords[i]
		start := current.Word[0:1]

		if tk.keywordTable[start] == nil {
			tk.keywordTable[start] = make([]TokenWord, 0)
		}

		insert, index := 0, 0
		for _, w := range tk.keywordTable[start] {
			if w.Length >= current.Length {
				insert++
			}
			index++
		}
		if len(tk.keywordTable[start]) <= i {
			tk.keywordTable[start] = append(tk.keywordTable[start], TokenWord{})
		}
		for index > insert {
			tk.keywordTable[start][index] = tk.keywordTable[start][index-1]
			index--
		}
		tk.keywordTable[start][insert] = current

	}
	return &tk
}

func (tk *Tokenizer) Cleanup() {
	tk.keywordTable = nil
}

func (tk *Tokenizer) IsDigitInRadix(ch rune, radix int) bool {
	if ch >= '0' && ch <= '9' {
		return int(ch-'0') < radix
	}
	if ch >= 'A' && ch <= 'Z' {
		return int(ch-'A'-10) < radix
	}
	if ch >= 'a' && ch <= 'z' {
		return int(ch-'a'-10) < radix
	}
	return false
}

func (tk *Tokenizer) GetToken(source string) (uint32, Token) {
	_, tokenLength, t := tk.ParseToken(source)
	return tokenLength, t
}

func (tk *Tokenizer) ParseToken(source string) (flags.ASTokenClass, uint32, Token) {
	//Whitespace token
	if ok, l, token := tk.IsWhiteSpace(source); ok == true {
		return flags.ASTokenWhitespace, l, token
	}
	if ok, l, token := tk.IsComment(source); ok == true {
		return flags.ASTokenComment, l, token
	}
	if ok, l, token := tk.IsConstant(source); ok == true {
		return flags.ASTokenValue, l, token
	}
	if ok, l, token := tk.IsIdentifier(source); ok == true {
		return flags.ASTokenIdentifier, l, token
	}
	if ok, l, token := tk.IsKeyword(source); ok == true {
		return flags.ASTokenKeyword, l, token
	}
	return flags.ASTokenUnknown, 1, ASttUnrecognizedToken
}

func (tk *Tokenizer) IsWhiteSpace(source string) (bool, uint32, Token) {

	wSpace := []rune(ASWhitespaceToken)
	src := []rune(source)

	if len(src) >= 3 &&
		uint32(src[0]) == 0xEF &&
		uint32(src[1]) == 0xBB &&
		uint32(src[2]) == 0xBF {
		return true, 3, ASttWhiteSpace
	}

	n := 0
	numWsChars := len(wSpace)
	for n = 0; n < len(src); n++ {
		isWhitespace := false
		for w := 0; w < numWsChars; w++ {
			if src[n] == wSpace[w] {
				isWhitespace = true
				break
			}
		}
		if !isWhitespace {
			break
		}
	}

	if n > 0 {
		return true, uint32(n), ASttWhiteSpace
	}

	return false, 0, ASttUnrecognizedToken
}

func (tk *Tokenizer) IsComment(source string) (bool, uint32, Token) {

	src := []rune(source)

	//Definately not a comment
	if len(src) < 2 {
		return false, 0, ASttUnrecognizedToken
	}
	//Not a comment
	if src[0] != '/' {
		return false, 0, ASttUnrecognizedToken
	}

	//Oneliner comment
	if src[1] == '/' {
		n := 0
		for n = 2; n < len(src); n++ {
			if src[n] == '\n' {
				break
			}
		}

		tlen := n
		if n < len(src) {
			tlen = n + 1
		}

		return true, uint32(tlen), ASttOnelineComment
	}

	//Multiliner comment
	if src[1] == '*' {
		n := 0
		for n = 2; n < len(src)-1; n++ {
			if src[n] == '*' && src[n+1] == '/' {
				break
			}
		}
		n++

		return true, uint32(n + 1), ASttMultilineComment
	}

	return false, 0, ASttUnrecognizedToken
}

func (tk *Tokenizer) IsConstant(source string) (bool, uint32, Token) {
	src := []rune(source)

	if (src[0] >= '0' && src[0] <= '9') || (src[0] == '.' && len(src) > 1 && src[1] >= '0' && src[1] <= '9') {
		if src[0] == '0' && len(src) > 1 {
			radix := 0
			switch src[1] {
			case 'b':
			case 'B':
				radix = 2
				break
			case 'o':
			case 'O':
				radix = 8
				break
			case 'd':
			case 'D':
				radix = 10
				break
			case 'x':
			case 'X':
				radix = 16
				break
			}
			if radix != 0 {
				n := 0
				for n = 2; n < len(src); n++ {
					if !tk.IsDigitInRadix(src[n], radix) {
						break
					}
				}
				return true, uint32(n), ASttBitsConstant
			}
		}
		n := 0
		for n = 0; n < len(src); n++ {
			if src[n] < '0' || src[n] > '9' {
				break
			}
		}

		if n < len(src) && (src[n] == '.' || src[n] == 'e' || src[n] == 'E') {
			if src[n] == '.' {
				n++
				for ; n < len(src); n++ {
					if src[n] < '0' || src[n] > '9' {
						break
					}
				}
			}

			if n < len(src) && (src[n] == 'e' || src[n] == 'E') {
				n++
				if n < len(src) && (src[n] == '-' || src[n] == '+') {
					n++
				}
				for ; n < len(src); n++ {
					if src[n] < '0' || src[n] > '9' {
						break
					}
				}
			}
			if n < len(src) && (src[n] == 'f' || src[n] == 'F') {
				return true, uint32(n + 1), ASttFloatConstant
			} else {
				return true, uint32(n), ASttDoubleConstant
			}
		}
		return true, uint32(n), ASttIntConstant
	}
	l := 0
	t := ASttUnrecognizedToken
	//String constants
	if src[0] == '"' || src[0] == '\'' {

		if len(src) >= 6 && src[0] == '"' && src[1] == '"' && src[2] == '"' {
			//Heredoc string constant
			n := 0
			for n = 3; n < len(src)-2; n++ {
				if src[n] == '"' && src[n+1] == '"' && src[n+2] == '"' {
					break
				}
			}
			return true, uint32(n + 3), ASttHeredocStringConstant
		} else {
			//Normal string constant
			ASttype := ASttStringConstant
			quote := src[0]
			evenSlashes := true
			n := 0
			for n = 1; n < len(src); n++ {
				if src[n] == '\n' {
					ASttype = ASttMultilineStringConstant
				}
				if src[n] == quote && evenSlashes {
					return true, uint32(n + 1), ASttype
				}
				if src[n] == '\\' {
					evenSlashes = !evenSlashes
				} else {
					evenSlashes = true
				}
			}
			t = ASttNonTerminatedStringConstant
			l = n
		}
		return true, uint32(l),t
	}

	return false, 0, ASttUnrecognizedToken
}

func (tk *Tokenizer) IsIdentifier(source string) (bool, uint32, Token) {
	src := []rune(source)
	c := src[0]

	if (c >= 'a' && c <= 'z') || (c >= 'A' && c <= 'Z') || c == '_' || c < 0 {
		AStt := ASttIdentifier
		tl := 1
		for n := 1; n < len(src); n++ {
			c = src[n]
			if (c >= 'a' && c <= 'z') || (c >= 'A' && c <= 'Z') || c == '_' || c < 0 {
				tl++
			} else {
				break
			}
		}

		//Check keyword
		if ok, _, _ := tk.IsKeyword(source); ok {
			return false, 0, ASttUnrecognizedToken
		}

		return true, uint32(tl), AStt
	}
	return false, 0, ASttUnrecognizedToken
}

func (tk *Tokenizer) IsKeyword(source string) (bool, uint32, Token) {
	src := []rune(source)
	start := string(src[0])
	tokenWord := tk.keywordTable[start]

	if tokenWord == nil {
		return false, 0, ASttUnrecognizedToken
	}

	for _, ptr := range tokenWord {
		wlen := ptr.Length
		if len(src) >= wlen && source[0:wlen] == ptr.Word {
			if wlen < len(src) &&
				((src[wlen-1] >= 'a' && src[wlen-1] <= 'z') ||
					(src[wlen-1] >= 'A' && src[wlen-1] <= 'Z') ||
					(src[wlen-1] >= '0' && src[wlen-1] <= '9')) &&
				((src[wlen] >= 'a' && src[wlen] <= 'z') ||
					(src[wlen] >= 'A' && src[wlen] <= 'Z') ||
					(src[wlen] >= '0' && src[wlen] <= '9') ||
					(src[wlen] == '_')) {
				// The token doesn't really match, even though
				// the start of the source matches the token
				continue
			}

			return true, uint32(wlen), ptr.Type
		}
	}
	return false, 0, ASttUnrecognizedToken
}
