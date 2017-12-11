package angelscript

//Token definition

const (
	ttUnrecognizedToken = iota

	ttEnd // End of file

	// White space and comments
	ttWhiteSpace       = " ;\t;\r;\n" // ' ', '\t', '\r', '\n', UTF8 byte-order-mark
	ttOnelineComment   = "//;\n"      // // \n
	ttMultilineComment = "/*;*/"      // /* */

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
	ttPlus     = "+"  // +
	ttMinus    = "-"  // -
	ttStar     = "*"  // *
	ttSlash    = "//" // /
	ttPercent  = "%"  // %
	ttStarStar = "**" // **

	ttHandle = "@" // @

	ttAddAssign = "+="  // +=
	ttSubAssign = "-="  // -=
	ttMulAssign = "*="  // *=
	ttDivAssign = "/="  // /=
	ttModAssign = "%="  // %=
	ttPowAssign = "**=" // **=

	ttOrAssign          = "|="   // |=
	ttAndAssign         = "&="   // &=
	ttXorAssign         = "^="   // ^=
	ttShiftLeftAssign   = "<<="  // <<=
	ttShiftRightLAssign = ">>="  // >>=
	ttShiftRightAAssign = ">>>=" // >>>=

	ttInc = "++" // ++
	ttDec = "--" // --

	ttDot   = "."  // .
	ttScope = "::" // ::

	// Statement tokens
	ttAssignment          = "=" // =
	ttEndStatement        = ";" // ;
	ttListSeparator       = "," // ,
	ttStartStatementBlock = "{" // {
	ttEndStatementBlock   = "}" // }
	ttOpenParanthesis     = "(" // (
	ttCloseParanthesis    = ")" // )
	ttOpenBracket         = "[" // [
	ttCloseBracket        = "]" // ]
	ttAmp                 = "&" // &

	// Bitwise operators
	ttBitOr              = "|"   // |
	ttBitNot             = "~"   // ~
	ttBitXor             = "^"   // ^
	ttBitShiftLeft       = "<<"  // <<
	ttBitShiftRight      = ">>"  // >>     // TODO: In Java this is the arithmetical shift
	ttBitShiftRightArith = ">>>" // >>>    // TODO: In Java this is the logical shift

	// Compare operators
	ttEqual              = "==" // ==
	ttNotEqual           = "!=" // !=
	ttLessThan           = "<"  // <
	ttGreaterThan        = ">"  // >
	ttLessThanOrEqual    = "<=" // <=
	ttGreaterThanOrEqual = ">=" // >=

	ttQuestion = "?" // ?
	ttColon    = ":" // :

	// Reserved keywords
	ttIf        = "if"        // if
	ttElse      = "else"      // else
	ttFor       = "for"       // for
	ttWhile     = "while"     // while
	ttBool      = "bool"      // bool
	ttFuncDef   = "funcdef"   // funcdef
	ttImport    = "import"    // import
	ttInt       = "int"       // int
	ttInt8      = "int8"      // int8
	ttInt16     = "int16"     // int16
	ttInt64     = "int64"     // int64
	ttInterface = "interface" // interface
	ttIs        = "is"        // is
	ttNotIs     = "!is"       // !is
	ttUInt      = "uint"      // uint
	ttUInt8     = "uint8"     // uint8
	ttUInt16    = "uint16"    // uint16
	ttUInt64    = "uint64"    // uint64
	ttFloat     = "float"     // float
	ttVoid      = "void"      // void
	ttTrue      = "true"      // true
	ttFalse     = "false"     // false
	ttReturn    = "return"    // return
	ttNot       = "not"       // not
	ttAnd       = "and;&&"    // and, &&
	ttOr        = "or;||"     // or, ||
	ttXor       = "xor;^^"    // xor, ^^
	ttBreak     = "break"     // break
	ttContinue  = "continue"  // continue
	ttConst     = "const"     // const
	ttDo        = "do"        // do
	ttDouble    = "double"    // double
	ttSwitch    = "switch"    // switch
	ttCase      = "case"      // case
	ttDefault   = "default"   // default
	ttIn        = "in"        // in
	ttOut       = "out"       // out
	ttInOut     = "inout"     // inout
	ttNull      = "null"      // null
	ttClass     = "class"     // class
	ttTypedef   = "typedef"   // typedef
	ttEnum      = "enum"      // enum
	ttCast      = "cast"      // cast
	ttPrivate   = "private"   // private
	ttProtected = "protected" // protected
	ttNamespace = "namespace" // namespace
	ttMixin     = "mixin"     // mixin
	ttAuto      = "auto"      // auto
)
