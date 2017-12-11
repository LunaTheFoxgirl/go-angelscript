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
package flags

//Return codes
type ASReturnCode = int32

const (
	ASReturnSuccess                          = ASReturnCode(0)
	ASReturnError                            = ASReturnCode(-1)
	ASReturnContextActive                    = ASReturnCode(-2)
	ASReturnContextNotFinished               = ASReturnCode(-3)
	ASReturnContextNotPrepare                = ASReturnCode(-4)
	ASReturnInvalidArg                       = ASReturnCode(-5)
	ASReturnNoFunction                       = ASReturnCode(-6)
	ASReturnNotSupported                     = ASReturnCode(-7)
	ASReturnInvalidName                      = ASReturnCode(-8)
	ASReturnNameTaken                        = ASReturnCode(-9)
	ASReturnInvalidDeclaration               = ASReturnCode(-10)
	ASReturnInvalidObject                    = ASReturnCode(-11)
	ASReturnInvalidType                      = ASReturnCode(-12)
	ASReturnAlreadyRegistered                = ASReturnCode(-13)
	ASReturnMultipleFunctions                = ASReturnCode(-14)
	ASReturnNoModule                         = ASReturnCode(-15)
	ASReturnNoGlobalVar                      = ASReturnCode(-16)
	ASReturnInvalidConfiguration             = ASReturnCode(-17)
	ASReturnInvalidInterface                 = ASReturnCode(-18)
	ASReturnCantBindAllFunctions             = ASReturnCode(-19)
	ASReturnLowerArrayDimensionNotRegistered = ASReturnCode(-20)
	ASReturnWrongConfigGroup                 = ASReturnCode(-21)
	ASReturnConfigGroupIsInUse               = ASReturnCode(-22)
	ASReturnIllegalBehaviourForType          = ASReturnCode(-23)
	ASReturnWrongCallingConvention           = ASReturnCode(-24)
	ASReturnBuildInProgress                  = ASReturnCode(-25)
	ASReturnInitGlobalVarsFailed             = ASReturnCode(-26)
	ASReturnOutOfMemory                      = ASReturnCode(-27)
	ASReturnModuleIsInUse                    = ASReturnCode(-28)
)

//Engine properties

type ASEngineProperties = uint32
const (
	ASPropertyAllowUnsafeReferences = ASEngineProperties(1)
	ASPropertyOptimizeBytecode = ASEngineProperties(2)
	ASPropertyCopyScriptSections = ASEngineProperties(3)
	ASPropertyMaxStackSize = ASEngineProperties(4)
	ASPropertyUseCharacterLiterals = ASEngineProperties(5)
	ASPropertyAllowMultilineStrings = ASEngineProperties(6)
	ASPropertyAllowImplicitHandleTypes = ASEngineProperties(7)
	ASPropertyBuildWithoutLineCues = ASEngineProperties(8)
	ASPropertyInitGlobalVarsAfterBuild = ASEngineProperties(9)
	ASPropertyRequireEnumScope = ASEngineProperties(10)
	ASPropertyScriptScanner = ASEngineProperties(11)
	ASPropertyIncludeJITInstructions = ASEngineProperties(12)
	ASPropertyStringEncoding = ASEngineProperties(13)
	ASPropertyPropertyAccessorMode = ASEngineProperties(14)
	ASPropertyExpandDefArrayToTemplate = ASEngineProperties(15)
	ASPropertyAutoGarbageCollect = ASEngineProperties(16)
	ASPropertyDissallowGlobalVars = ASEngineProperties(17)
	ASPropertyAlwaysImplDefaultConstruct = ASEngineProperties(18)
	ASPropertyCompilerWarnings = ASEngineProperties(19)
	ASPropertyDisallowValueAssingForRefType = ASEngineProperties(20)
	ASPropertyAlterSyntaxNamedArgs = ASEngineProperties(21)
	ASPropertyDisableIntegerDivision = ASEngineProperties(22)
	ASPropertyDisallowEmptyListElements = ASEngineProperties(23)
	ASPropertyPrivatePropAsProtected = ASEngineProperties(24)
	ASPropertyAllowUnicodeIdentifiers = ASEngineProperties(25)
	ASPropertyHeredocTrimMode = ASEngineProperties(26)
	ASPropertyLastProperty = ASEngineProperties(27)
)

//Calling conventions
type ASCallConvention = uint32

const (
	ASCallConventionCDeclare            = ASCallConvention(0)
	ASCallConventionSTDCall             = ASCallConvention(1)
	ASCallConventionThisCallAsGlobal    = ASCallConvention(2)
	ASCallConventionThisCall            = ASCallConvention(3)
	ASCallConventionCDeclareObjectLast  = ASCallConvention(4)
	ASCallConventionCDeclareObjectFirst = ASCallConvention(5)
	ASCallConventionCallGeneric         = ASCallConvention(6)
	ASCallConventionThisCallObjectLast  = ASCallConvention(7)
	ASCallConventionThisCallObjectFirst = ASCallConvention(8)
)

//Object type flags
type ASObjectType = uint32

const (
	ASTypeReference               = ASObjectType(1 << 0)
	ASTypeValue                   = ASObjectType(1 << 1)
	ASTypeGC                      = ASObjectType(1 << 2)
	ASTypePod                     = ASObjectType(1 << 3)
	ASTypeNoHandle                = ASObjectType(1 << 4)
	ASTypeScoped                  = ASObjectType(1 << 5)
	ASTypeTemplate                = ASObjectType(1 << 6)
	ASTypeASHandle                = ASObjectType(1 << 7)
	ASTypeAppClass                = ASObjectType(1 << 8)
	ASTypeAppClassConstructor     = ASObjectType(1 << 9)
	ASTypeAppClassDestructor      = ASObjectType(1 << 10)
	ASTypeAppClassAssignment      = ASObjectType(1 << 11)
	ASTypeAppClassCopyConstructor = ASObjectType(1 << 12)
	AsTypeAppClassC               = ASObjectType(ASTypeAppClass + ASTypeAppClassConstructor)
	AsTypeAppClassCD              = ASObjectType(ASTypeAppClass + ASTypeAppClassConstructor + ASTypeAppClassDestructor)
	AsTypeAppClassCA              = ASObjectType(ASTypeAppClass + ASTypeAppClassConstructor + ASTypeAppClassAssignment)
	AsTypeAppClassCK              = ASObjectType(ASTypeAppClass + ASTypeAppClassConstructor + ASTypeAppClassCopyConstructor)
	AsTypeAppClassCDA             = ASObjectType(ASTypeAppClass + ASTypeAppClassConstructor + ASTypeAppClassDestructor + ASTypeAppClassAssignment)
	AsTypeAppClassCDK             = ASObjectType(ASTypeAppClass + ASTypeAppClassConstructor + ASTypeAppClassDestructor + ASTypeAppClassCopyConstructor)
	AsTypeAppClassCAK             = ASObjectType(ASTypeAppClass + ASTypeAppClassConstructor + ASTypeAppClassAssignment + ASTypeAppClassCopyConstructor)
	AsTypeAppClassCDAK            = ASObjectType(ASTypeAppClass + ASTypeAppClassConstructor + ASTypeAppClassDestructor + ASTypeAppClassAssignment + ASTypeAppClassCopyConstructor)
	AsTypeAppClassD               = ASObjectType(ASTypeAppClass + ASTypeAppClassDestructor)
	AsTypeAppClassDA              = ASObjectType(ASTypeAppClass + ASTypeAppClassDestructor + ASTypeAppClassAssignment)
	AsTypeAppClassDAK             = ASObjectType(ASTypeAppClass + ASTypeAppClassDestructor + ASTypeAppClassAssignment + ASTypeAppClassCopyConstructor)
	AsTypeAppClassA               = ASObjectType(ASTypeAppClass + ASTypeAppClassAssignment)
	AsTypeAppClassAK              = ASObjectType(ASTypeAppClass + ASTypeAppClassAssignment + ASTypeAppClassCopyConstructor)
	AsTypeAppClassK               = ASObjectType(ASTypeAppClass + ASTypeAppClassCopyConstructor)
	AsTypeAppPrimitive            = ASObjectType(1 << 13)
	AsTypeAppFloat                = ASObjectType(1 << 14)
	AsTypeAppArray                = ASObjectType(1 << 15)
	AsTypeAppClassAllInts         = ASObjectType(1 << 16)
	AsTypeAppClassAllFloats       = ASObjectType(1 << 17)
	AsTypeNoCount                 = ASObjectType(1 << 18)
	AsTypeAppClassAlign8          = ASObjectType(1 << 19)
	AsTypeAppImplicitHandle       = ASObjectType(1 << 20)
	AsTypeAppMaskValidFlags       = ASObjectType(0x1FFFFF)

	//Internal Flags
	AsTypeScriptObject    = ASObjectType(1 << 21)
	AsTypeShared          = ASObjectType(1 << 22)
	AsTypeNoInherit       = ASObjectType(1 << 23)
	AsTypeFuncDev         = ASObjectType(1 << 24)
	AsTypeListPattern     = ASObjectType(1 << 25)
	AsTypeEnum            = ASObjectType(1 << 26)
	AsTypeTemplateSubType = ASObjectType(1 << 27)
	AsTypeTypeDef         = ASObjectType(1 << 28)
	AsTypeAbstract        = ASObjectType(1 << 29)
	AsTypeAppAlign16      = ASObjectType(1 << 30)
)

//Behaviours
type ASBehaviour = uint32

const (
	ASBehaviourConstruct = ASBehaviour(0)
	ASBehaviourListConstruct = ASBehaviour(1)
	ASBehaviourDestruct = ASBehaviour(2)
	ASBehaviourFactory = ASBehaviour(3)
	ASBehaviourAddReference = ASBehaviour(4)
	ASBehaviourRelease = ASBehaviour(5)
	ASBehaviourGetWeakReferenceFlag = ASBehaviour(6)
	ASBehaviourTemplateCallback = ASBehaviour(7)
	ASBehaviourFirstGC = ASBehaviour(8)
	ASBehaviourGetReferenceCount = ASBehaviourFirstGC
	ASBehaviourSetGCFlag = ASBehaviour(9)
	ASBehaviourGetGCFlag = ASBehaviour(10)
	ASBehaviourEnumReferences = ASBehaviour(11)
	ASBehaviourReleaseReferences = ASBehaviour(12)
	ASBehaviourLastGC = ASBehaviourReleaseReferences
	ASBehaviourMax = ASBehaviour(13)
)

//Context States
type ASContextState = uint32

const (
	ASExecutionFinished = ASContextState(0)
	ASExecutionSuspended = ASContextState(1)
	ASExecutionAborted = ASContextState(2)
	ASExecutionException = ASContextState(3)
	ASExecutionPrepared = ASContextState(4)
	ASExecutionUninitialized = ASContextState(5)
	ASExecutionActive = ASContextState(6)
	ASExecutionError = ASContextState(7)
)

//MessageTypes
type ASMessageType = uint32
const (
	ASMsgTypeError = ASMessageType(0)
	ASMsgTypeWarning = ASMessageType(1)
	ASMsgTypeInformation = ASMessageType(2)
)

//Garbage Collector Flags
const (
	ASGCFullCycle      = 1
	ASGCOneStep        = 2
	ASGCDestroyGarbage = 4
	ASGCDetectGarbage  = 0
)

//Token classes
type ASTokenClass = uint32
const (
	ASTokenUnknown = ASTokenClass(iota)
	ASTokenKeyword
	ASTokenValue
	ASTokenTypeentifier
	ASTokenComment
	ASTokenWhitespace
)

//Type id flags
type ASTokenType = uint32
const (
	ASTypeIdVoid               = ASTokenType(0)
	ASTypeIdBool               = ASTokenType(1)
	ASTypeIdInt8               = ASTokenType(2)
	ASTypeIdInt16              = ASTokenType(3)
	ASTypeIdInt32              = ASTokenType(4)
	ASTypeIdInt64              = ASTokenType(5)
	ASTypeIdUInt8              = ASTokenType(6)
	ASTypeIdUInt16             = ASTokenType(7)
	ASTypeIdUInt32             = ASTokenType(8)
	ASTypeIdUInt64             = ASTokenType(9)
	ASTypeIdFloat              = ASTokenType(10)
	ASTypeIdDouble             = ASTokenType(11)
	ASTypeIdObjectHandle       = ASTokenType(0x40000000)
	ASTypeIdHandleToConst      = ASTokenType(0x20000000)
	ASTypeIdMaskObject         = ASTokenType(0x1C000000)
	ASTypeIdAppObject          = ASTokenType(0x04000000)
	ASTypeIdScriptObject       = ASTokenType(0x08000000)
	ASTypeIdTemplate           = ASTokenType(0x10000000)
	ASTypeIdMaskSequenceNumber = ASTokenType(0x03FFFFFF)
)

// Type modifiers
const (
	ASTypeModifierNone = iota
	ASTypeModifierInRef
	ASTypeModifierOutRef
	ASTypeModifierInoutRef
	ASTypeModifierConst
)

// GetModule flags
const (
	ASGetModuleOnlyIfExists = iota
	ASGetModuleCreateIfNotExists
	ASGetModuleAlwaysCreate
)

// Compile flags
const ASCompileAddToModule = 1

//Function types
const (
	ASFunctionDummy     = -1
	ASFunctionSystem    = 0
	ASFunctionScript    = 1
	ASFunctionInterface = 2
	ASFunctionVirtual   = 3
	ASFunctionFuncDef   = 4
	ASFunctionImported  = 5
	ASFunctionDelegate  = 6
)
