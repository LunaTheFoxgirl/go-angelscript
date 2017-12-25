package angelscript

import (
	"testing"
	"fmt"
	"github.com/Member1221/go-angelscript/tokenizer"
)

func TestTokenizeAndParse(t *testing.T) {
	program := `import "test"
class MyClass {
	int myInt = 0;
}	
`
	fmt.Println("CODE:\n", program)
	code := NewScriptCode()
	code.SetCode("MyCode", program, false)
	
	engine := ScriptEngine{}
	engine.tok = angelscript.NewTokenizer()
	
	builder := ScriptBuilder{&engine}
	
	parser := NewParser(&builder)
	parser.ParseScriptX(&code)
	fmt.Println("Parse Completed")
	fmt.Println(parser.Node.ToTList())
}





