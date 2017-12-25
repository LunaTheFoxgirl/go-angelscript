package angelscript

import (

)

//TODO: Implement builder
type ScriptBuilder struct {
	Engine *ScriptEngine
}

func (sb *ScriptBuilder) DoesTypeExist(typ string) bool {
	return true
}