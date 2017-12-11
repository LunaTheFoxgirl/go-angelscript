package angelscript

import ()

type IScriptEngine interface {
	AddRef() int
	Release() int
	ShutDownAndRelease() int
}

type ScriptEngine struct {
}
