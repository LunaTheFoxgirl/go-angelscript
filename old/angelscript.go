// +build ignore

/*
go-angelscript
Copyright (c) 2017 Clipsey

   This software is provided 'as-is', without any express or implied 
   warranty. In no event will the authors be held liable for any 
   damages arising from the use of this software.

   Permission is granted to anyone to use this software for any 
   purpose, including commercial applications, and to alter it and 
   redistribute it freely, subject to the following restrictions:

   1. The origin of this software must not be misrepresented; you 
      must not claim that you wrote the original software. If you use
      this software in a product, an acknowledgment in the product 
      documentation would be appreciated but is not required.

   2. Altered source versions must be plainly marked as such, and 
      must not be misrepresented as being the original software.

   3. This notice may not be removed or altered from any source 
      distribution.

   The original version of this library can be located at:
   https://www.github.com/Member1221/go-angelscript

Clipsey
clipseypone@gmail.com
*/
package angelscript

/*
#cgo CXXFLAGS: -Wall -fPIC -fno-strict-aliasing -std=c++11
#include <stddef.h>
#include <stdlib.h>
#include "angelscript_c.h"
*/
import "C"
import "unsafe"
import "github.com/Member1221/go-angelscript/flags"

const ASAngelScriptVersion = 23102

//Util stuff.
type ASByte = C.asBYTE
type ASWord = C.asWORD
type ASUInt = C.asUINT
type ASDWord = C.asDWORD
type ASPWord = C.asPWORD
type ASQWORD = C.asQWORD
type ASInt64 = C.asINT64
type ASBool = C.asBOOL

func GetLibraryVersion() string {
	return C.GoString(C.asGetLibraryOptions())
}

func GetLibraryOptions() string {
	return C.GoString(C.asGetLibraryOptions())
}

//Current function being wrapped.
var wrapFuncs map[string]wrapableFunction

type wrapableFunction struct {
	fun func(...interface{})
	args ...interface{}
}

func (w *wrapableFunction) doUnwrap() {
	w.fun(w.args...)
}


// ScriptEngine is an angelscript engine/context.
type ScriptEngine struct {
	engine *C.struct_asIScriptEngine
}

func CreateScriptEngine() *ScriptEngine {
	return &ScriptEngine{
		engine: C.asCreateScriptEngine(ASDWord(ASAngelScriptVersion)),
	}
}

func CreateScriptEngineVersion(version int32) *ScriptEngine {
	return &ScriptEngine{
		engine: C.asCreateScriptEngine(ASDWord(version)),
	}
}

//export function_wrapper
func function_wrapper(f func(vars []interface{}), input C.int) {
	wrapping_function.doUnwrap()
}

func RegisterGlobalFunction(declaration string, function func(...interface{}), args ...interface{}, callingConvention flags.ASCallConvention) {
	cdeclaration := C.CString(declaration)
	defer C.free(unsafe.Pointer(cdeclaration))
	wrappingFunction = wrapableFunction{
		fun: func,
		args: args,
	}
}


func (engine *ScriptEngine) WriteMessage(section string, row, collumn int, typ flags.ASMessageType, message string) {
	// TODO: Implement this
	//C.asEngine_WriteMessage()
}
