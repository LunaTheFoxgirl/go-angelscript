package angelscript

import (
	_ "encoding/binary"
	_ "reflect"
	"github.com/Member1221/go-angelscript/flags"
)
//Token definition

type Token = uint32

const (
	ttUnrecognizedToken = Token(iota)

	ttEnd // End of file

	// White space and comments
	ttWhiteSpace       // ' ', '\t', '\r', '\n', UTF8 byte-order-mark
	ttOnelineComment   // // \n
	ttMultilineComment // /* */

	// Atoms
	ttIdentifier                  // abc123
	ttIntConstant                 // 1234
	ttFloatConstant               // 12.34e56f
	ttDoubleConstant              // 12.34e56
	ttStringConstant              // "123"
	ttMultilineStringConstant     //
	ttHeredocStringConstant       // """text"""
	ttNonTerminatedStringConstant // "123
	ttBitsConstant                // 0xFFFF

	// Math operators
	ttPlus     // +
	ttMinus    // -
	ttStar     // *
	ttSlash    // /
	ttPercent  // %
	ttStarStar // **

	ttHandle // @

	ttAddAssign // +=
	ttSubAssign // -=
	ttMulAssign // *=
	ttDivAssign // /=
	ttModAssign // %=
	ttPowAssign // **=

	ttOrAssign          // |=
	ttAndAssign         // &=
	ttXorAssign         // ^=
	ttShiftLeftAssign   // <<=
	ttShiftRightLAssign // >>=
	ttShiftRightAAssign // >>>=

	ttInc // ++
	ttDec // --

	ttDot   // .
	ttScope // ::

	// Statement tokens
	ttAssignment          // =
	ttEndStatement        // ;
	ttListSeparator       // ,
	ttStartStatementBlock // {
	ttEndStatementBlock   // }
	ttOpenParanthesis     // (
	ttCloseParanthesis    // )
	ttOpenBracket         // [
	ttCloseBracket        // ]
	ttAmp                 // &

	// Bitwise operators
	ttBitOr              // |
	ttBitNot             // ~
	ttBitXor             // ^
	ttBitShiftLeft       // <<
	ttBitShiftRight      // >>     // TODO: In Java this is the arithmetical shift
	ttBitShiftRightArith // >>>    // TODO: In Java this is the logical shift

	// Compare operators
	ttEqual              // ==
	ttNotEqual           // !=
	ttLessThan           // <
	ttGreaterThan        // >
	ttLessThanOrEqual    // <=
	ttGreaterThanOrEqual // >=

	ttQuestion // ?
	ttColon    // :

	// Reserved keywords
	ttIf        // if
	ttElse      // else
	ttFor       // for
	ttWhile     // while
	ttBool      // bool
	ttFuncDef   // funcdef
	ttImport    // import
	ttInt       // int
	ttInt8      // int8
	ttInt16     // int16
	ttInt64     // int64
	ttInterface // interface
	ttIs        // is
	ttNotIs     // !is
	ttUInt      // uint
	ttUInt8     // uint8
	ttUInt16    // uint16
	ttUInt64    // uint64
	ttFloat     // float
	ttVoid      // void
	ttTrue      // true
	ttFalse     // false
	ttReturn    // return
	ttNot       // not
	ttAnd       // and, &&
	ttOr        // or, ||
	ttXor       // xor, ^^
	ttBreak     // break
	ttContinue  // continue
	ttConst     // const
	ttDo        // do
	ttDouble    // double
	ttSwitch    // switch
	ttCase      // case
	ttDefault   // default
	ttIn        // in
	ttOut       // out
	ttInOut     // inout
	ttNull      // null
	ttClass     // class
	ttTypedef   // typedef
	ttEnum      // enum
	ttCast      // cast
	ttPrivate   // private
	ttProtected // protected
	ttNamespace // namespace
	ttMixin     // mixin
	ttAuto      // auto
)

const (
	ASWhitespaceToken = " \t\r\n"
	ASThisToken = "this"
	ASFromToken = "from"
	ASSuperToken = "super"
	ASSharedToken = "shared"
	ASFinalToken = "final"
	ASOverrideToken = "override"
	ASGetToken = "get"
	ASSetToken = "set"
	ASAbstractToken = "abstract"
	ASFunctionToken = "function"
	ASIfHandleToken = "if_handle_then_const"
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
	asTokenDef("+", ttPlus),
	asTokenDef("+=", ttAddAssign),
	asTokenDef("++", ttInc),
	asTokenDef("-", ttMinus),
	asTokenDef("-=", ttSubAssign),
	asTokenDef("--", ttDec),
	asTokenDef("*", ttStar),
	asTokenDef("*=", ttMulAssign),
	asTokenDef("/", ttSlash),
	asTokenDef("/=", ttDivAssign),
	asTokenDef("%", ttPercent),
	asTokenDef("%=", ttModAssign),
	asTokenDef("**", ttStarStar),
	asTokenDef("**=", ttPowAssign),
	asTokenDef("=", ttAssignment),
	asTokenDef("==", ttEqual),
	asTokenDef(".", ttDot),
	asTokenDef("|", ttBitOr),
	asTokenDef("|=", ttOrAssign),
	asTokenDef("||", ttOr),
	asTokenDef("&", ttAmp),
	asTokenDef("&=", ttAndAssign),
	asTokenDef("&&", ttAnd),
	asTokenDef("^", ttBitXor),
	asTokenDef("^=", ttXorAssign),
	asTokenDef("^^", ttXor),
	asTokenDef("<", ttLessThan),
	asTokenDef("<=", ttLessThanOrEqual),
	asTokenDef("<<", ttBitShiftLeft),
	asTokenDef("<<=", ttShiftLeftAssign),
	asTokenDef(">", ttGreaterThan),
	asTokenDef(">=", ttGreaterThanOrEqual),
	asTokenDef(">>", ttBitShiftRight),
	asTokenDef(">>=", ttShiftRightLAssign),
	asTokenDef(">>>", ttBitShiftRightArith),
	asTokenDef(">>>=", ttShiftRightAAssign),
	asTokenDef("~", ttBitNot),
	asTokenDef(";", ttEndStatement),
	asTokenDef(",", ttListSeparator),
	asTokenDef("{", ttStartStatementBlock),
	asTokenDef("}", ttEndStatementBlock),
	asTokenDef("(", ttOpenParanthesis),
	asTokenDef(")", ttCloseParanthesis),
	asTokenDef("[", ttOpenBracket),
	asTokenDef("]", ttCloseBracket),
	asTokenDef("?", ttQuestion),
	asTokenDef(":", ttColon),
	asTokenDef("::", ttScope),
	asTokenDef("!", ttNot),
	asTokenDef("!=", ttNotEqual),
	asTokenDef("!is", ttNotIs),
	asTokenDef("@", ttHandle),
	asTokenDef("and", ttAnd),
	asTokenDef("auto", ttAuto),
	asTokenDef("bool", ttBool),
	asTokenDef("break", ttBreak),
	asTokenDef("case", ttCase),
	asTokenDef("cast", ttCast),
	asTokenDef("class", ttClass),
	asTokenDef("const", ttConst),
	asTokenDef("continue", ttContinue),
	asTokenDef("default", ttDefault),
	asTokenDef("do", ttDo),
	asTokenDef("double", ttDouble),
	asTokenDef("else", ttElse),
	asTokenDef("enum", ttEnum),
	asTokenDef("false", ttFalse),
	asTokenDef("float", ttFloat),
	asTokenDef("for", ttFor),
	asTokenDef("funcdef", ttFuncDef),
	asTokenDef("if", ttIf),
	asTokenDef("import", ttImport),
	asTokenDef("in", ttIn),
	asTokenDef("inout", ttInOut),
	asTokenDef("int", ttInt),
	asTokenDef("int8", ttInt8),
	asTokenDef("int16", ttInt16),
	asTokenDef("int32", ttInt),
	asTokenDef("int64", ttInt64),
	asTokenDef("interface", ttInterface),
	asTokenDef("is", ttIs),
	asTokenDef("mixin", ttMixin),
	asTokenDef("namespace", ttNamespace),
	asTokenDef("not", ttNot),
	asTokenDef("null", ttNull),
	asTokenDef("or", ttOr),
	asTokenDef("out", ttOut),
	asTokenDef("private", ttPrivate),
	asTokenDef("protected", ttProtected),
	asTokenDef("return", ttReturn),
	asTokenDef("switch", ttSwitch),
	asTokenDef("true", ttTrue),
	asTokenDef("typedef", ttTypedef),
	asTokenDef("uint", ttUInt),
	asTokenDef("uint8", ttUInt8),
	asTokenDef("uint16", ttUInt16),
	asTokenDef("uint32", ttUInt),
	asTokenDef("uint64", ttUInt64),
	asTokenDef("void", ttVoid),
	asTokenDef("while", ttWhile),
	asTokenDef("xor", ttXor),
}

type Tokenizer struct {
	engine int
	keywordTable map[string][]TokenWord
}

func NewTokenizer() *Tokenizer {
	tk := Tokenizer{}
	tk.engine = 0
	for i := 0; i < len(tokenWords); i++ {
		current := tokenWords[i]
		start := current.Word[0:1]
		
		if tk.keywordTable[start] == nil {
			tk.keywordTable[start] = make([]TokenWord, 32)
		}
		
		insert, index := 0, 0
		for _, w := range tk.keywordTable[start] {
			if w.Length >= current.Length {
				insert++
			}
			index++
		}
		for index > insert {
			tk.keywordTable[start][index] = tk.keywordTable[start][index - 1]
			index--
		}
		tk.keywordTable[start][insert] = current
		
	}
	return &tk
}

func (tk *Tokenizer) Cleanup() {
	tk.keywordTable = nil
}

func (tk *Tokenizer) IsDigitinRadix(ch rune, radix int) bool {
	if ch >= '0' && ch <= '9' {
		return int(ch - '0') < radix;
	}
	if ch >= 'A' && ch <= 'Z' {
		return int(ch - 'A'-10) < radix;
	}
	if ch >= 'a' && ch <= 'z' {
		return int(ch - 'a'-10) < radix;
	}
	return false
}

func (tk *Tokenizer) GetToken(source string) (flags.ASTokenType, uint32, flags.ASTokenClass) {
	t, tokenLength, tokenType := tk.ParseToken(source)
	return t, tokenLength, tokenType
}

func (tk *Tokenizer) ParseToken(source string) (flags.ASTokenClass, uint32, Token) {
	//Whitespace token
	if ok, l, token := tk.IsWhiteSpace(source); ok == true {
		return flags.ASTokenWhitespace, l, token
	}
	
	return flags.ASTokenUnknown, 1, ttUnrecognizedToken
}

func (tk *Tokenizer) IsWhiteSpace(source string) (bool, uint32, Token) {
	
	if 	len(source) >= 3 && 
		uint32(source[0]) == 0xEF && 
		uint32(source[1]) == 0xBB && 
		uint32(source[2]) == 0xBF {
		return true, 3, ttWhiteSpace
	}
	
	n := 0
	numWsChars := len(ASWhitespaceToken)
	for n = 0; n < len(source); n++ {
		isWhitespace := false
		for w := 0; w < numWsChars; w++ {
			if source[n] == ASWhitespaceToken[w] {
				isWhitespace = true
				break; 
			}
		}
		if !isWhitespace {
			break;
		}
	}
	
	if n > 0 {
		return true, uint32(n), ttWhiteSpace
	}
	
	return false, 0, ttUnrecognizedToken
}

func (tk *Tokenizer) IsComment(source string) (bool, uint32, Token) {
	
	src := []rune(source)
	
	//Definately not a comment
	if len(src) < 2 {
		return false, 0, ttUnrecognizedToken
	}
	//Not a comment
	if src[0] != '/' {
		return false, 0, ttUnrecognizedToken
	}
	
	//Oneliner comment
	if src[1] == '/' {
		n := 0
		for n = 2; n < len(src); n++ {
			if src[n] == '\n' {
				break;
			}
		}
		
		tlen := n
		if n < sourceLength {
			tlen = n + 1
		}
		
		return true, tlen, ttOnelineComment
	}
	
	//Multiliner comment
	if src[1] == '*' {
		n := 0
		for n = 2; n < len(src)-1; n++ {
			if src[n] == '*' && src[n+1] == '*' {
				break;
			}
		}
		n++
		
		return true, n+1, ttOnelineComment
	}
	
	return false, 0, ttUnrecognizedToken
}
