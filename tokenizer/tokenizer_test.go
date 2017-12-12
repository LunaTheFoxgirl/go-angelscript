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
	"testing"
	_ "github.com/Member1221/go-angelscript"
	"fmt"
)

func TestTokenizeProgram(t *testing.T) {
	program := `
//import base stuff
import "std"
"""MyClass is a test class, this is a test docstring."""
class MyClass
{
	/*
		Properties
	*/
	int property = 22;
	
	/*
		Constructor
	*/
	MyClass()
	{
		
	}
	
	/*
		Methods/functions
	*/
	
	//This is a function.
	void MyFunction()
	{
		std::print("Hello, world!");
	}
}`
	end := uint32(len(program))
	cur := program
	i := uint32(0)
	tk := NewTokenizer()
	for i < end {
		cur = program[i:]
		_, l, token := tk.ParseToken(cur);
		if token == 0 {
			return
		}
		if int(l) < len(cur) && l > 0 {
			def := GetDefinition(token)
			if def != "" {
				if def == "[Comment: 1 liner]" {
					fmt.Print(def + "\n")
				} else {
					fmt.Print(def + " ")
				}
			} else {
				if cur[:1] == "\n" {
					fmt.Print("\n")
				}
			}
		}
		i += l
	}
	
}