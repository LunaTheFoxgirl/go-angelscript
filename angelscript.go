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

//Return codes
type ASCode = int

const (
	ASSuccess                          = ASCode(0)
	ASError                            = ASCode(-1)
	ASContextActive                    = ASCode(-2)
	ASContextNotFinished               = ASCode(-3)
	ASContextNotPrepare                = ASCode(-4)
	ASInvalidArg                       = ASCode(-5)
	ASNoFunction                       = ASCode(-6)
	ASNotSupported                     = ASCode(-7)
	ASInvalidName                      = ASCode(-8)
	ASNameTaken                        = ASCode(-9)
	ASInvalidDeclaration               = ASCode(-10)
	ASInvalidObject                    = ASCode(-11)
	ASInvalidType                      = ASCode(-12)
	ASAlreadyRegistered                = ASCode(-13)
	ASMultipleFunctions                = ASCode(-14)
	ASNoModule                         = ASCode(-15)
	ASNoGlobalVar                      = ASCode(-16)
	ASInvalidConfiguration             = ASCode(-17)
	ASInvalidInterface                 = ASCode(-18)
	ASCantBindAllFunctions             = ASCode(-19)
	ASLowerArrayDimensionNotRegistered = ASCode(-20)
	ASWrongConfigGroup                 = ASCode(-21)
	ASConfigGroupIsInUse               = ASCode(-22)
	ASIllegalBehaviourForType          = ASCode(-23)
	ASWrongCallingConvention           = ASCode(-24)
	ASBuildInProgress                  = ASCode(-25)
	ASInitGlobalVarsFailed             = ASCode(-26)
	ASOutOfMemory                      = ASCode(-27)
	ASModuleIsInUse                    = ASCode(-28)
)

//Util stuff.
type ASBYTE = byte
type ASWORD = uint16
type ASUINT = uint
type ASDWORD = uint32
type ASQWORD = uint64
type ASPWORD = uint
type ASINT64 = int64
type ASBOOL = bool

// PointerSize is the size of a pointer, default 8 bytes.
// Changing this is not advised, but possible.
var ASPointerSize uint8 = 8

type UserData struct {
	Type ASPWORD
	Data interface{}
}

type IScriptEngine interface {
	AddRef() int
	Release() int
	ShutDownAndRelease() int
}

type IThreadManager interface {
	
}

type IScriptModule interface {
	
}

type IScriptContext interface {
	
}

type IScriptGeneric interface {
	
}

type IScriptObject interface {
	
}

type ITypeInfo interface {
	GetEngine() *ScriptEngine
	GetConfigGroup() string
	GetAccessMask() ASDWORD
	GetModule() *Module
	
	AddRef()
	Release()
	
	GetName() string
	GetNamespace() string
	GetBaseType() *ITypeInfo
	DerivesFrom(Type *ITypeInfo) bool
	GetFlags() ASDWORD
	GetSize() ASUINT
	GetTypeId() int
	GetSubTypeId(index uint) int
	GetSubType(index uint) *ITypeInfo
	GetSubTypeCount() uint
	
	GetInterfaceCount() uint
	GetInterface(index uint) *ITypeInfo
	Implements(Type *ITypeInfo) bool
	
	GetFactoryCount() uint
	GetFactoryByIndex(index uint) *IScriptFunction
	GetFactoryByDecl(decl string) *IScriptFunction
	
	GetMethodCount() uint
	GetMethodByIndex(index uint, getVirtual bool) *IScriptFunction
	GetMethodByName(name string, getVirtual bool) *IScriptFunction
	GetMethodByDecl(decl string, getVirtual bool) *IScriptFunction
	
	GetPropertyCount() uint
	GetProperty(index uint) (name string, typeId int, isPrivate, isProtected bool, offset int, isReference bool, accessMask uint16, err error)
	GetPropertyDeclaration(index uint, includeNamespace bool) string
	
	GetBehaviourCount() uint
	GetBehaviourByIndex(index uint) (IScriptFunction /*TODO: EBehaviour outBehaviour*/)
	
	GetChildFuncdefCount() uint
	GetChildFuncdef(index uint) *ITypeInfo
	GetParentType() *ITypeInfo
	
	GetEnumValueCount() uint
	GetEnumValueByIndex(index uint) (a string, outVal int)
	
	GetTypedefTypeId() int
	
	GetFuncdefSignature() *IScriptFunction
	
	SetUserData(data UserData) *UserData
	GetUserData(id ASPWORD) *UserData
}

type IScriptFunction interface {
	
}

type IBinaryStream interface {
	
}

type ILockableSharedBool interface {
	
}