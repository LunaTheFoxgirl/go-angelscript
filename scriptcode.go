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
)

type ScriptCode struct {
	Name      string
	Code      string
	Offset    int
	Length    int
	Shared    bool
	Positions []uint32
	IDX       int
}

func NewScriptCode() ScriptCode {
	return ScriptCode{
		Offset:     0,
		Code:       "",
		Length:		0,
		Shared:     false,
	}
}

func (sc *ScriptCode) SetCode(name, code string, makecopy bool) int {
	
	if !sc.Shared && sc.Code != "" {
		sc.Code = ""
	}
	
	if makecopy {
		sc.Shared = false
		sc.Code = code
	} else {
		sc.Shared = true
		sc.Code = code
	}
	
	coder := []rune(code)
	sc.Positions = append(sc.Positions, 0)
	for n := 0; n < len(coder); n++ {
		if coder[n] == '\n' {
			sc.Positions = append(sc.Positions, uint32(n+1))
		}
	}
	sc.Positions = append(sc.Positions, uint32(len(coder)))
	sc.Length = len(coder)
	//ASSuccess
	return 0
}

func (sc *ScriptCode) PosToRowCol(pos int) (int, int) {
	if len(sc.Positions) == 0 {
		return sc.Offset, 1
	}
	
	max := len(sc.Positions)-1
	min := 0
	i := max/2
	for {
		if sc.Positions[i] < uint32(pos) {
			if min == i {
				break
			}
			
			min = i
			i = (max+min/2)
		}else if sc.Positions[i] > uint32(pos) {
			if max == i {
				break
			}
			
			max = i
			i = (max+min/2)
		} else {
			break
		}
	}
	
	return i+1+sc.Offset, (pos-int(sc.Positions[i]))+1
}

func (sc *ScriptCode) TokenEquals(pos, l int, str string) bool {
	if pos + l > sc.Length {
		return false
	}
	s := sc.Code[pos:pos+l]
	if s == str {
		return true
	}
	return false
}
