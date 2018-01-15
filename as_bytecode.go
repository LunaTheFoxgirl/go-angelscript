package angelscript

import (
	"strconv"
)

// InstructionType is the type of an bytecode instruction.
type InstructionType byte

// Instruction type sizes
var ASBCTypeSize [21]int = [21]int{
	0, // asBCTYPE_INFO
	1, // asBCTYPE_NO_ARG
	1, // asBCTYPE_W_ARG
	1, // asBCTYPE_wW_ARG
	2, // asBCTYPE_DW_ARG
	2, // asBCTYPE_rW_DW_ARG
	3, // asBCTYPE_QW_ARG
	3, // asBCTYPE_DW_DW_ARG
	2, // asBCTYPE_wW_rW_rW_ARG
	3, // asBCTYPE_wW_QW_ARG
	2, // asBCTYPE_wW_rW_ARG
	1, // asBCTYPE_rW_ARG
	2, // asBCTYPE_wW_DW_ARG
	3, // asBCTYPE_wW_rW_DW_ARG
	2, // asBCTYPE_rW_rW_ARG
	2, // asBCTYPE_wW_W_ARG
	4, // asBCTYPE_QW_DW_ARG
	3, // asBCTYPE_rW_QW_ARG
	2, // asBCTYPE_W_DW_ARG
	3, // asBCTYPE_rW_W_DW_ARG
	3, // asBCTYPE_rW_DW_DW_ARG
}

const (
	ASBCTYPE_INFO = InstructionType(iota)
	ASBCTYPE_NO_ARG
	ASBCTYPE_W_ARG
	ASBCTYPE_wW_ARG
	ASBCTYPE_DW_ARG
	ASBCTYPE_rW_DW_ARG
	ASBCTYPE_QW_ARG
	ASBCTYPE_DW_DW_ARG
	ASBCTYPE_wW_rW_rW_ARG
	ASBCTYPE_wW_QW_ARG
	ASBCTYPE_wW_rW_ARG
	ASBCTYPE_rW_ARG
	ASBCTYPE_wW_DW_ARG
	ASBCTYPE_wW_rW_DW_ARG
	ASBCTYPE_rW_rW_ARG
	ASBCTYPE_wW_W_ARG
	ASBCTYPE_QW_DW_ARG
	ASBCTYPE_rW_QW_ARG
	ASBCTYPE_W_DW_ARG
	ASBCTYPE_rW_W_DW_ARG
	ASBCTYPE_rW_DW_DW_ARG

	ASBCTYPE_PTR_ARG    = ASBCTYPE_DW_ARG
	ASBCTYPE_PTR_DW_ARG = ASBCTYPE_DW_DW_ARG
	ASBCTYPE_wW_PTR_ARG = ASBCTYPE_wW_DW_ARG
	ASBCTYPE_rW_PTR_ARG = ASBCTYPE_rW_DW_ARG
)

// EBCInstruction is an bytecode instruction.
// See all ASBC_ consts for bytecode instructions.
type ByteCodeInstruction byte

// Byte code instructions
const (
	ASBC_PopPtr = ByteCodeInstruction(iota)
	ASBC_PshGPtr
	ASBC_PshC4
	ASBC_PshV4
	ASBC_PSF
	ASBC_SwapPtr
	ASBC_NOT
	ASBC_PshG4
	ASBC_LdGRdR4
	ASBC_CALL
	ASBC_RET
	ASBC_JMP
	ASBC_JZ
	ASBC_JNZ
	ASBC_JS
	ASBC_JNS
	ASBC_JP
	ASBC_JNP
	ASBC_TZ
	ASBC_TNZ
	ASBC_TS
	ASBC_TNS
	ASBC_TP
	ASBC_TNP
	ASBC_NEGi
	ASBC_NEGf
	ASBC_NEGd
	ASBC_INCi16
	ASBC_INCi8
	ASBC_DECi16
	ASBC_DECi8
	ASBC_INCi
	ASBC_DECi
	ASBC_INCf
	ASBC_DECf
	ASBC_INCd
	ASBC_DECd
	ASBC_IncVi
	ASBC_DecVi
	ASBC_BNOT
	ASBC_BAND
	ASBC_BOR
	ASBC_BXOR
	ASBC_BSLL
	ASBC_BSRL
	ASBC_BSRA
	ASBC_COPY
	ASBC_PshC8
	ASBC_PshVPtr
	ASBC_RDSPtr
	ASBC_CMPd
	ASBC_CMPu
	ASBC_CMPf
	ASBC_CMPi
	ASBC_CMPIi
	ASBC_CMPIf
	ASBC_CMPIu
	ASBC_JMPP
	ASBC_PopRPtr
	ASBC_PshRPtr
	ASBC_STR
	ASBC_CALLSYS
	ASBC_CALLBND
	ASBC_SUSPEND
	ASBC_ALLOC
	ASBC_FREE
	ASBC_LOADOBJ
	ASBC_STOREOBJ
	ASBC_GETOBJ
	ASBC_REFCPY
	ASBC_CHKREF
	ASBC_GETOBJREF
	ASBC_GETREF
	ASBC_PshNull
	ASBC_ClrVPtr
	ASBC_OBJTYPE
	ASBC_TYPEID
	ASBC_SetV4
	ASBC_SetV8
	ASBC_ADDSi
	ASBC_CpyVtoV4
	ASBC_CpyVtoV8
	ASBC_CpyVtoR4
	ASBC_CpyVtoR8
	ASBC_CpyVtoG4
	ASBC_CpyRtoV4
	ASBC_CpyRtoV8
	ASBC_CpyGtoV4
	ASBC_WRTV1
	ASBC_WRTV2
	ASBC_WRTV4
	ASBC_WRTV8
	ASBC_RDR1
	ASBC_RDR2
	ASBC_RDR4
	ASBC_RDR8
	ASBC_LDG
	ASBC_LDV
	ASBC_PGA
	ASBC_CmpPtr
	ASBC_VAR
	ASBC_iTOf
	ASBC_fTOi
	ASBC_uTOf
	ASBC_fTOu
	ASBC_sbTOi
	ASBC_swTOi
	ASBC_ubTOi
	ASBC_uwTOi
	ASBC_dTOi
	ASBC_dTOu
	ASBC_dTOf
	ASBC_iTOd
	ASBC_uTOd
	ASBC_fTOd
	ASBC_ADDi
	ASBC_SUBi
	ASBC_MULi
	ASBC_DIVi
	ASBC_MODi
	ASBC_ADDf
	ASBC_SUBf
	ASBC_MULf
	ASBC_DIVf
	ASBC_MODf
	ASBC_ADDd
	ASBC_SUBd
	ASBC_MULd
	ASBC_DIVd
	ASBC_MODd
	ASBC_ADDIi
	ASBC_SUBIi
	ASBC_MULIi
	ASBC_ADDIf
	ASBC_SUBIf
	ASBC_MULIf
	ASBC_SetG4
	ASBC_ChkRefS
	ASBC_ChkNullV
	ASBC_CALLINTF
	ASBC_iTOb
	ASBC_iTOw
	ASBC_SetV1
	ASBC_SetV2
	ASBC_Cast
	ASBC_i64TOi
	ASBC_uTOi64
	ASBC_iTOi64
	ASBC_fTOi64
	ASBC_dTOi64
	ASBC_fTOu64
	ASBC_dTOu64
	ASBC_i64TOf
	ASBC_u64TOf
	ASBC_i64TOd
	ASBC_u64TOd
	ASBC_NEGi64
	ASBC_INCi64
	ASBC_DECi64
	ASBC_BNOT64
	ASBC_ADDi64
	ASBC_SUBi64
	ASBC_MULi64
	ASBC_DIVi64
	ASBC_MODi64
	ASBC_BAND64
	ASBC_BOR64
	ASBC_BXOR64
	ASBC_BSLL64
	ASBC_BSRL64
	ASBC_BSRA64
	ASBC_CMPi64
	ASBC_CMPu64
	ASBC_ChkNullS
	ASBC_ClrHi
	ASBC_JitEntry
	ASBC_CallPtr
	ASBC_FuncPtr
	ASBC_LoadThisR
	ASBC_PshV8
	ASBC_DIVu
	ASBC_MODu
	ASBC_DIVu64
	ASBC_MODu64
	ASBC_LoadRObjR
	ASBC_LoadVObjR
	ASBC_RefCpyV
	ASBC_JLowZ
	ASBC_JLowNZ
	ASBC_AllocMem
	ASBC_SetListSize
	ASBC_PshListElmnt
	ASBC_SetListType
	ASBC_POWi
	ASBC_POWu
	ASBC_POWf
	ASBC_POWd
	ASBC_POWdi
	ASBC_POWi64
	ASBC_POWu64
	ASBC_Thiscall1
	ASBC_MAXBYTECODE

	// Temporary tokens. Can't be output to the final program
	ASBC_VarDecl = ByteCodeInstruction(251)
	ASBC_Block   = ByteCodeInstruction(iota)
	ASBC_ObjInfo
	ASBC_LINE
	ASBC_LABEL
)

// ByteCodeInfo contains relevant information about a byte code instruction.
type ByteCodeInfo struct {
	// Instruction is the instruction itself.
	Instruction ByteCodeInstruction
	// Type is the type of instruction
	Type InstructionType
	// StackIncrement is the size of the instruction.
	StackIncrement int
	// Name is the name of the instruction
	Name string
}

// helper function to keep the definition clean.
func asbcinfo(inst ByteCodeInstruction, t InstructionType, s int, name string) ByteCodeInfo {
	return ByteCodeInfo{inst, t, s, name}
}

// helper function to define dummy bytecode entries.
func asbcdummy(id int) ByteCodeInfo {
	return asbcinfo(ASBC_MAXBYTECODE, ASBCTYPE_INFO, 0, "<bc DUMMY "+strconv.Itoa(id)+">")
}

// helper function to define dummy bytecode entries.
func asbcdummyrange(from, to int) []ByteCodeInfo {
	out := make([]ByteCodeInfo, 0)
	for i := from; i <= to; i++ {
		out = append(out, asbcdummy(i))
	}
	return out
}

//TODO: Finish off human names
// ASBCInfo is a list of info about ASBC instructions.
var ASBCInfo []ByteCodeInfo = []ByteCodeInfo{
	asbcinfo(ASBC_PopPtr, ASBCTYPE_NO_ARG, -int(ASPointerSize), "<pop pointer>"),
	asbcinfo(ASBC_PshGPtr, ASBCTYPE_PTR_ARG, int(ASPointerSize), "<push pointer>"),
	asbcinfo(ASBC_PshC4, ASBCTYPE_DW_ARG, 1, "<push C4>"),
	asbcinfo(ASBC_PshV4, ASBCTYPE_rW_ARG, 1, "<push V4>"),
	asbcinfo(ASBC_PSF, ASBCTYPE_rW_ARG, int(ASPointerSize), "<psf>"),
	asbcinfo(ASBC_SwapPtr, ASBCTYPE_NO_ARG, 0, "<swap pointer>"),
	asbcinfo(ASBC_NOT, ASBCTYPE_rW_ARG, 0, "<not>"),
	asbcinfo(ASBC_PshG4, ASBCTYPE_PTR_ARG, 1, "<push G4>"),
	asbcinfo(ASBC_LdGRdR4, ASBCTYPE_wW_PTR_ARG, 0, "<LdGRdR4>"),
	asbcinfo(ASBC_CALL, ASBCTYPE_DW_ARG, 0xFFFF, "<call>"),
	asbcinfo(ASBC_RET, ASBCTYPE_W_ARG, 0xFFFF, "<return>"),

	//jumps
	asbcinfo(ASBC_JMP, ASBCTYPE_DW_ARG, 0, "<jump>"),
	asbcinfo(ASBC_JZ, ASBCTYPE_DW_ARG, 0, "<jump if 0>"),
	asbcinfo(ASBC_JNZ, ASBCTYPE_DW_ARG, 0, "<jump if not 0>"),
	asbcinfo(ASBC_JS, ASBCTYPE_DW_ARG, 0, "<jump if sign>"),
	asbcinfo(ASBC_JNS, ASBCTYPE_DW_ARG, 0, "<jump if not sign>"),
	asbcinfo(ASBC_JP, ASBCTYPE_DW_ARG, 0, ""),
	asbcinfo(ASBC_JNP, ASBCTYPE_DW_ARG, 0, ""),
	asbcinfo(ASBC_TZ, ASBCTYPE_NO_ARG, 0, ""),
	asbcinfo(ASBC_TNZ, ASBCTYPE_NO_ARG, 0, ""),
	asbcinfo(ASBC_TS, ASBCTYPE_NO_ARG, 0, ""),
	asbcinfo(ASBC_TNS, ASBCTYPE_NO_ARG, 0, ""),
	asbcinfo(ASBC_TP, ASBCTYPE_NO_ARG, 0, ""),
	asbcinfo(ASBC_TNP, ASBCTYPE_NO_ARG, 0, ""),

	// NEGs
	asbcinfo(ASBC_NEGi, ASBCTYPE_rW_ARG, 0, ""),
	asbcinfo(ASBC_NEGf, ASBCTYPE_rW_ARG, 0, ""),
	asbcinfo(ASBC_NEGd, ASBCTYPE_rW_ARG, 0, ""),

	// Inc and Dec
	asbcinfo(ASBC_INCi16, ASBCTYPE_NO_ARG, 0, ""),
	asbcinfo(ASBC_INCi8, ASBCTYPE_NO_ARG, 0, ""),
	asbcinfo(ASBC_DECi16, ASBCTYPE_NO_ARG, 0, ""),
	asbcinfo(ASBC_DECi8, ASBCTYPE_NO_ARG, 0, ""),
	asbcinfo(ASBC_INCi, ASBCTYPE_NO_ARG, 0, ""),
	asbcinfo(ASBC_DECi, ASBCTYPE_NO_ARG, 0, ""),
	asbcinfo(ASBC_INCf, ASBCTYPE_NO_ARG, 0, ""),
	asbcinfo(ASBC_DECf, ASBCTYPE_NO_ARG, 0, ""),
	asbcinfo(ASBC_INCd, ASBCTYPE_NO_ARG, 0, ""),
	asbcinfo(ASBC_DECd, ASBCTYPE_NO_ARG, 0, ""),

	//Special cases?
	asbcinfo(ASBC_IncVi, ASBCTYPE_rW_ARG, 0, ""),
	asbcinfo(ASBC_DecVi, ASBCTYPE_rW_ARG, 0, ""),
	asbcinfo(ASBC_BNOT, ASBCTYPE_rW_ARG, 0, ""),

	//B Operators
	asbcinfo(ASBC_BAND, ASBCTYPE_wW_rW_rW_ARG, 0, ""),
	asbcinfo(ASBC_BOR, ASBCTYPE_wW_rW_rW_ARG, 0, ""),
	asbcinfo(ASBC_BSLL, ASBCTYPE_wW_rW_rW_ARG, 0, ""),
	asbcinfo(ASBC_BSRL, ASBCTYPE_wW_rW_rW_ARG, 0, ""),
	asbcinfo(ASBC_BSRA, ASBCTYPE_wW_rW_rW_ARG, 0, ""),

	//Copy
	asbcinfo(ASBC_COPY, ASBCTYPE_DW_ARG, -int(ASPointerSize), "<copy>"),

	//More pointer stuff
	asbcinfo(ASBC_PshC8, ASBCTYPE_QW_ARG, 2, ""),
	asbcinfo(ASBC_PshVPtr, ASBCTYPE_rW_ARG, int(ASPointerSize), ""),
	asbcinfo(ASBC_RDSPtr, ASBCTYPE_NO_ARG, 0, ""),

	//Compare
	asbcinfo(ASBC_CMPd, ASBCTYPE_rW_rW_ARG, 0, ""),
	asbcinfo(ASBC_CMPu, ASBCTYPE_rW_rW_ARG, 0, ""),
	asbcinfo(ASBC_CMPf, ASBCTYPE_rW_rW_ARG, 0, ""),
	asbcinfo(ASBC_CMPi, ASBCTYPE_rW_rW_ARG, 0, ""),
	asbcinfo(ASBC_CMPIi, ASBCTYPE_DW_ARG, 0, ""),
	asbcinfo(ASBC_CMPIf, ASBCTYPE_DW_ARG, 0, ""),
	asbcinfo(ASBC_CMPIu, ASBCTYPE_DW_ARG, 0, ""),

	//EVEN MORE pointer stuff, and jump?
	asbcinfo(ASBC_JMPP, ASBCTYPE_rW_ARG, 0, "<jump pointer>"),
	asbcinfo(ASBC_PopRPtr, ASBCTYPE_NO_ARG, -int(ASPointerSize), "<pop ref pointer>"),
	asbcinfo(ASBC_PshRPtr, ASBCTYPE_NO_ARG, int(ASPointerSize), "<push ref pointer>"),

	//string and essential stuff
	asbcinfo(ASBC_STR, ASBCTYPE_W_ARG, 1+int(ASPointerSize), "<str>"),
	asbcinfo(ASBC_CALLSYS, ASBCTYPE_DW_ARG, 0xFFFF, "<callsys>"),
	asbcinfo(ASBC_CALLBND, ASBCTYPE_DW_ARG, 0xFFFF, "<callbind>"),
	asbcinfo(ASBC_SUSPEND, ASBCTYPE_NO_ARG, 0, "<suspend>"),
	asbcinfo(ASBC_ALLOC, ASBCTYPE_PTR_DW_ARG, 0xFFFF, "<allocate>"),
	asbcinfo(ASBC_FREE, ASBCTYPE_wW_PTR_ARG, 0, "<free>"),
	asbcinfo(ASBC_LOADOBJ, ASBCTYPE_rW_ARG, 0, "<load object>"),
	asbcinfo(ASBC_STOREOBJ, ASBCTYPE_wW_ARG, 0, "<store object>"),
	asbcinfo(ASBC_GETOBJ, ASBCTYPE_W_ARG, 0, "<get object>"),
	asbcinfo(ASBC_REFCPY, ASBCTYPE_W_ARG, -int(ASPointerSize), "<ref copy>"),
	asbcinfo(ASBC_CHKREF, ASBCTYPE_NO_ARG, 0, "<check ref>"),
	asbcinfo(ASBC_GETOBJREF, ASBCTYPE_W_ARG, 0, "<get object ref>"),
	asbcinfo(ASBC_GETREF, ASBCTYPE_W_ARG, 0, "<get ref>"),
	asbcinfo(ASBC_PshNull, ASBCTYPE_NO_ARG, int(ASPointerSize), "<push null>"),
	asbcinfo(ASBC_ClrVPtr, ASBCTYPE_wW_ARG, 0, "<clear virtual pointer>"),
	asbcinfo(ASBC_OBJTYPE, ASBCTYPE_PTR_ARG, int(ASPointerSize), "<objtype>"),
	asbcinfo(ASBC_TYPEID, ASBCTYPE_DW_ARG, 1, "<typeid>"),

	asbcinfo(ASBC_SetV4, ASBCTYPE_wW_DW_ARG, 0, ""),
	asbcinfo(ASBC_SetV8, ASBCTYPE_wW_QW_ARG, 0, ""),
	asbcinfo(ASBC_ADDSi, ASBCTYPE_W_DW_ARG, 0, ""),

	//Copy X to Y(4/8)
	asbcinfo(ASBC_CpyVtoV4, ASBCTYPE_wW_rW_ARG, 0, ""),
	asbcinfo(ASBC_CpyVtoV8, ASBCTYPE_wW_rW_ARG, 0, ""),
	asbcinfo(ASBC_CpyVtoR4, ASBCTYPE_rW_ARG, 0, ""),
	asbcinfo(ASBC_CpyVtoR8, ASBCTYPE_rW_ARG, 0, ""),
	asbcinfo(ASBC_CpyVtoG4, ASBCTYPE_rW_PTR_ARG, 0, ""),
	asbcinfo(ASBC_CpyRtoV4, ASBCTYPE_rW_ARG, 0, ""),
	asbcinfo(ASBC_CpyRtoV8, ASBCTYPE_rW_ARG, 0, ""),
	asbcinfo(ASBC_CpyGtoV4, ASBCTYPE_rW_PTR_ARG, 0, ""),

	asbcinfo(ASBC_WRTV1, ASBCTYPE_rW_ARG, 0, ""),
	asbcinfo(ASBC_WRTV2, ASBCTYPE_rW_ARG, 0, ""),
	asbcinfo(ASBC_WRTV4, ASBCTYPE_rW_ARG, 0, ""),
	asbcinfo(ASBC_WRTV8, ASBCTYPE_rW_ARG, 0, ""),
	asbcinfo(ASBC_RDR1, ASBCTYPE_rW_ARG, 0, ""),
	asbcinfo(ASBC_RDR2, ASBCTYPE_rW_ARG, 0, ""),
	asbcinfo(ASBC_RDR4, ASBCTYPE_rW_ARG, 0, ""),
	asbcinfo(ASBC_RDR8, ASBCTYPE_rW_ARG, 0, ""),

	asbcinfo(ASBC_LDG, ASBCTYPE_PTR_ARG, 0, ""),
	asbcinfo(ASBC_LDV, ASBCTYPE_rW_ARG, 0, ""),
	asbcinfo(ASBC_PGA, ASBCTYPE_PTR_ARG, int(ASPointerSize), ""),
	asbcinfo(ASBC_CmpPtr, ASBCTYPE_rW_rW_ARG, 0, ""),
	asbcinfo(ASBC_VAR, ASBCTYPE_rW_rW_ARG, int(ASPointerSize), "<var>"),
	asbcinfo(ASBC_iTOf, ASBCTYPE_rW_ARG, 0, ""),
	asbcinfo(ASBC_fTOi, ASBCTYPE_rW_ARG, 0, ""),
	asbcinfo(ASBC_uTOf, ASBCTYPE_rW_ARG, 0, ""),
	asbcinfo(ASBC_fTOu, ASBCTYPE_rW_ARG, 0, ""),
	asbcinfo(ASBC_sbTOi, ASBCTYPE_rW_ARG, 0, ""),
	asbcinfo(ASBC_swTOi, ASBCTYPE_rW_ARG, 0, ""),
	asbcinfo(ASBC_ubTOi, ASBCTYPE_rW_ARG, 0, ""),
	asbcinfo(ASBC_uwTOi, ASBCTYPE_rW_ARG, 0, ""),
	asbcinfo(ASBC_dTOi, ASBCTYPE_wW_rW_ARG, 0, ""),
	asbcinfo(ASBC_dTOu, ASBCTYPE_wW_rW_ARG, 0, ""),
	asbcinfo(ASBC_dTOf, ASBCTYPE_wW_rW_ARG, 0, ""),
	asbcinfo(ASBC_iTOd, ASBCTYPE_wW_rW_ARG, 0, ""),
	asbcinfo(ASBC_uTOd, ASBCTYPE_wW_rW_ARG, 0, ""),
	asbcinfo(ASBC_fTOd, ASBCTYPE_wW_rW_ARG, 0, ""),
	asbcinfo(ASBC_ADDi, ASBCTYPE_wW_rW_rW_ARG, 0, ""),
	asbcinfo(ASBC_SUBi, ASBCTYPE_wW_rW_rW_ARG, 0, ""),
	asbcinfo(ASBC_MULi, ASBCTYPE_wW_rW_rW_ARG, 0, ""),
	asbcinfo(ASBC_DIVi, ASBCTYPE_wW_rW_rW_ARG, 0, ""),
	asbcinfo(ASBC_MODi, ASBCTYPE_wW_rW_rW_ARG, 0, ""),
	asbcinfo(ASBC_ADDf, ASBCTYPE_wW_rW_rW_ARG, 0, ""),
	asbcinfo(ASBC_SUBf, ASBCTYPE_wW_rW_rW_ARG, 0, ""),
	asbcinfo(ASBC_MULf, ASBCTYPE_wW_rW_rW_ARG, 0, ""),
	asbcinfo(ASBC_DIVf, ASBCTYPE_wW_rW_rW_ARG, 0, ""),
	asbcinfo(ASBC_MODf, ASBCTYPE_wW_rW_rW_ARG, 0, ""),
	asbcinfo(ASBC_ADDd, ASBCTYPE_wW_rW_rW_ARG, 0, ""),
	asbcinfo(ASBC_SUBd, ASBCTYPE_wW_rW_rW_ARG, 0, ""),
	asbcinfo(ASBC_MULd, ASBCTYPE_wW_rW_rW_ARG, 0, ""),
	asbcinfo(ASBC_DIVd, ASBCTYPE_wW_rW_rW_ARG, 0, ""),
	asbcinfo(ASBC_MODd, ASBCTYPE_wW_rW_rW_ARG, 0, ""),
	asbcinfo(ASBC_ADDIi, ASBCTYPE_wW_rW_rW_ARG, 0, ""),
	asbcinfo(ASBC_SUBIi, ASBCTYPE_wW_rW_rW_ARG, 0, ""),
	asbcinfo(ASBC_MULIi, ASBCTYPE_wW_rW_rW_ARG, 0, ""),
	asbcinfo(ASBC_ADDIf, ASBCTYPE_wW_rW_rW_ARG, 0, ""),
	asbcinfo(ASBC_SUBIf, ASBCTYPE_wW_rW_rW_ARG, 0, ""),
	asbcinfo(ASBC_MULIf, ASBCTYPE_wW_rW_rW_ARG, 0, ""),
	asbcinfo(ASBC_SetG4, ASBCTYPE_PTR_DW_ARG, 0, ""),
	asbcinfo(ASBC_ChkRefS, ASBCTYPE_NO_ARG, 0, ""),
	asbcinfo(ASBC_ChkNullV, ASBCTYPE_rW_ARG, 0, ""),
	asbcinfo(ASBC_CALLINTF, ASBCTYPE_DW_ARG, 0xFFFF, "<call intf>"),
	asbcinfo(ASBC_iTOb, ASBCTYPE_rW_ARG, 0, ""),
	asbcinfo(ASBC_iTOw, ASBCTYPE_rW_ARG, 0, ""),
	asbcinfo(ASBC_SetV1, ASBCTYPE_wW_DW_ARG, 0, ""),
	asbcinfo(ASBC_SetV2, ASBCTYPE_wW_DW_ARG, 0, ""),
	asbcinfo(ASBC_Cast, ASBCTYPE_DW_ARG, -int(ASPointerSize), "<CASt>"),
	asbcinfo(ASBC_i64TOi, ASBCTYPE_wW_rW_ARG, 0, ""),
	asbcinfo(ASBC_uTOi64, ASBCTYPE_wW_rW_ARG, 0, ""),
	asbcinfo(ASBC_iTOi64, ASBCTYPE_wW_rW_ARG, 0, ""),
	asbcinfo(ASBC_fTOi64, ASBCTYPE_wW_rW_ARG, 0, ""),
	asbcinfo(ASBC_dTOi64, ASBCTYPE_rW_ARG, 0, ""),
	asbcinfo(ASBC_fTOu64, ASBCTYPE_wW_rW_ARG, 0, ""),
	asbcinfo(ASBC_dTOu64, ASBCTYPE_rW_ARG, 0, ""),
	asbcinfo(ASBC_fTOu64, ASBCTYPE_wW_rW_ARG, 0, ""),
	asbcinfo(ASBC_dTOu64, ASBCTYPE_rW_ARG, 0, ""),
	asbcinfo(ASBC_i64TOf, ASBCTYPE_wW_rW_ARG, 0, ""),
	asbcinfo(ASBC_u64TOf, ASBCTYPE_wW_rW_ARG, 0, ""),
	asbcinfo(ASBC_i64TOd, ASBCTYPE_rW_ARG, 0, ""),
	asbcinfo(ASBC_u64TOd, ASBCTYPE_rW_ARG, 0, ""),
	asbcinfo(ASBC_NEGi64, ASBCTYPE_rW_ARG, 0, ""),
	asbcinfo(ASBC_INCi64, ASBCTYPE_NO_ARG, 0, ""),
	asbcinfo(ASBC_DECi64, ASBCTYPE_NO_ARG, 0, ""),
	asbcinfo(ASBC_BNOT64, ASBCTYPE_rW_ARG, 0, ""),
	asbcinfo(ASBC_ADDi64, ASBCTYPE_wW_rW_rW_ARG, 0, ""),
	asbcinfo(ASBC_SUBi64, ASBCTYPE_wW_rW_rW_ARG, 0, ""),
	asbcinfo(ASBC_MULi64, ASBCTYPE_wW_rW_rW_ARG, 0, ""),
	asbcinfo(ASBC_DIVi64, ASBCTYPE_wW_rW_rW_ARG, 0, ""),
	asbcinfo(ASBC_MODi64, ASBCTYPE_wW_rW_rW_ARG, 0, ""),
	asbcinfo(ASBC_BAND64, ASBCTYPE_wW_rW_rW_ARG, 0, ""),
	asbcinfo(ASBC_BOR64, ASBCTYPE_wW_rW_rW_ARG, 0, ""),
	asbcinfo(ASBC_BXOR64, ASBCTYPE_wW_rW_rW_ARG, 0, ""),
	asbcinfo(ASBC_BSLL64, ASBCTYPE_wW_rW_rW_ARG, 0, ""),
	asbcinfo(ASBC_BSRL64, ASBCTYPE_wW_rW_rW_ARG, 0, ""),
	asbcinfo(ASBC_BSRA64, ASBCTYPE_wW_rW_rW_ARG, 0, ""),
	asbcinfo(ASBC_CMPi64, ASBCTYPE_rW_rW_ARG, 0, ""),
	asbcinfo(ASBC_CMPu64, ASBCTYPE_rW_rW_ARG, 0, ""),
	asbcinfo(ASBC_ChkNullS, ASBCTYPE_W_ARG, 0, "<check null S>"),
	asbcinfo(ASBC_ClrHi, ASBCTYPE_NO_ARG, 0, "<clr hi>"),
	asbcinfo(ASBC_JitEntry, ASBCTYPE_PTR_ARG, 0, "<JIT entry>"),
	asbcinfo(ASBC_CallPtr, ASBCTYPE_rW_ARG, 0xFFFF, "<call pointer>"),
	asbcinfo(ASBC_FuncPtr, ASBCTYPE_PTR_ARG, int(ASPointerSize), "<func pointer>"),
	asbcinfo(ASBC_LoadThisR, ASBCTYPE_W_DW_ARG, 0, ""),
	asbcinfo(ASBC_PshV8, ASBCTYPE_rW_ARG, 2, ""),
	asbcinfo(ASBC_DIVu, ASBCTYPE_wW_rW_rW_ARG, 0, ""),
	asbcinfo(ASBC_MODu, ASBCTYPE_wW_rW_rW_ARG, 0, ""),
	asbcinfo(ASBC_DIVu64, ASBCTYPE_wW_rW_rW_ARG, 0, ""),
	asbcinfo(ASBC_MODu64, ASBCTYPE_wW_rW_rW_ARG, 0, ""),
	asbcinfo(ASBC_LoadRObjR, ASBCTYPE_rW_W_DW_ARG, 0, ""),
	asbcinfo(ASBC_LoadVObjR, ASBCTYPE_rW_W_DW_ARG, 0, ""),
	asbcinfo(ASBC_RefCpyV, ASBCTYPE_wW_PTR_ARG, 0, ""),
	asbcinfo(ASBC_JLowZ, ASBCTYPE_DW_ARG, 0, ""),
	asbcinfo(ASBC_JLowNZ, ASBCTYPE_DW_ARG, 0, ""),
	asbcinfo(ASBC_AllocMem, ASBCTYPE_wW_DW_ARG, 0, "<alloc mem>"),
	asbcinfo(ASBC_SetListSize, ASBCTYPE_rW_DW_DW_ARG, 0, "<set list size>"),
	asbcinfo(ASBC_PshListElmnt, ASBCTYPE_rW_DW_ARG, int(ASPointerSize), "<push list element>"),
	asbcinfo(ASBC_SetListType, ASBCTYPE_rW_DW_DW_ARG, 0, "<set list type>"),
	asbcinfo(ASBC_POWi, ASBCTYPE_wW_rW_rW_ARG, 0, ""),
	asbcinfo(ASBC_POWu, ASBCTYPE_wW_rW_rW_ARG, 0, ""),
	asbcinfo(ASBC_POWf, ASBCTYPE_wW_rW_rW_ARG, 0, ""),
	asbcinfo(ASBC_POWd, ASBCTYPE_wW_rW_rW_ARG, 0, ""),
	asbcinfo(ASBC_POWdi, ASBCTYPE_wW_rW_rW_ARG, 0, ""),
	asbcinfo(ASBC_POWi64, ASBCTYPE_wW_rW_rW_ARG, 0, ""),
	asbcinfo(ASBC_POWu64, ASBCTYPE_wW_rW_rW_ARG, 0, ""),
	asbcinfo(ASBC_Thiscall1, ASBCTYPE_DW_ARG, -int(ASPointerSize)-1, "<thiscall1>"),

	//Dummies
	asbcdummy(201),
	asbcdummy(202),
	asbcdummy(203),
	asbcdummy(204),
	asbcdummy(205),
	asbcdummy(206),
	asbcdummy(207),
	asbcdummy(208),
	asbcdummy(209),
	asbcdummy(210),
	asbcdummy(211),
	asbcdummy(212),
	asbcdummy(213),
	asbcdummy(214),
	asbcdummy(215),
	asbcdummy(216),
	asbcdummy(217),
	asbcdummy(218),
	asbcdummy(219),
	asbcdummy(220),
	asbcdummy(221),
	asbcdummy(222),
	asbcdummy(223),
	asbcdummy(224),
	asbcdummy(225),
	asbcdummy(226),
	asbcdummy(227),
	asbcdummy(228),
	asbcdummy(229),
	asbcdummy(230),
	asbcdummy(231),
	asbcdummy(232),
	asbcdummy(233),
	asbcdummy(234),
	asbcdummy(235),
	asbcdummy(236),
	asbcdummy(237),
	asbcdummy(238),
	asbcdummy(239),
	asbcdummy(240),
	asbcdummy(241),
	asbcdummy(242),
	asbcdummy(243),
	asbcdummy(244),
	asbcdummy(245),
	asbcdummy(246),
	asbcdummy(247),
	asbcdummy(248),
	asbcdummy(249),
	asbcdummy(250),

	asbcinfo(ASBC_VarDecl, ASBCTYPE_W_ARG, 0, "<var decl>"),
	asbcinfo(ASBC_Block, ASBCTYPE_INFO, 0, "<block info>"),
	asbcinfo(ASBC_ObjInfo, ASBCTYPE_rW_DW_ARG, 0, "<obj info>"),
	asbcinfo(ASBC_LINE, ASBCTYPE_INFO, 0, "<line info>"),
	asbcinfo(ASBC_LABEL, ASBCTYPE_INFO, 0, "<label info>"),
}

// ByteInstruction is an instruction in the bytecode.
type ByteInstruction struct {
	OPCode   ByteCodeInstruction
	Arg      uint64
	WArg     [3]int
	size     int
	stackInc int

	Marked    bool
	stackSize int

	Next     *ByteInstruction
	Previous *ByteInstruction
}

// NewByteInstruction creates a new instance of ByteInstruction.
func NewByteInstruction() *ByteInstruction {
	inst := ByteInstruction{}
	inst.Next = nil
	inst.Previous = nil

	inst.OPCode = ASBC_LABEL

	inst.Arg = 0
	inst.WArg = [3]int{0, 0, 0}
	inst.size = 0
	inst.stackInc = 0
	inst.Marked = false
	inst.stackSize = 0

	return &inst
}

// AddAfter adds an instruction after this instruction.
func (bi *ByteInstruction) AddAfter(next *ByteInstruction) {
	if bi.Next != nil {
		bi.Previous = next
	}

	next.Next = bi.Next
	next.Previous = bi
	bi.Next = next
}

// AddBefore adds an instruction before this instruction.
func (bi *ByteInstruction) AddBefore(prev *ByteInstruction) {
	if bi.Previous != nil {
		bi.Next = prev
	}

	prev.Previous = bi.Next
	prev.Next = bi
	bi.Previous = prev
}

// GetSize gets the size of the instruction.
func (bi *ByteInstruction) GetSize() int {
	return bi.size
}

// GetStackIncrease gets the stack increase.
func (bi *ByteInstruction) GetStackIncrease() int {
	return bi.stackInc
}

// Remove removes this instruction from the instruction list.
func (bi *ByteInstruction) Remove() {
	if bi.Previous != nil {
		bi.Previous.Next = bi.Next
	}
	if bi.Next != nil {
		bi.Next.Previous = bi.Previous
	}
	bi.Next = nil
	bi.Previous = nil
}

type ByteCode struct {
	LineNumbers  []int
	SectionIdxs  []int
	LargestStack int

	engine   *ScriptEngine
	first    *ByteInstruction
	last     *ByteInstruction
	tempVars []int
}

func (bc *ByteCode) TempVarsExists(item int) bool {
	for _, i := range bc.tempVars {
		if i == item {
			return true
		}
	}
}

// NewBytecode creates a new bytecode instance.
func NewBytecode(engine *ScriptEngine) ByteCode {
	return ByteCode{first: nil, last: nil, LargestStack: -1, tempVars: make([]int, 0), engine: engine}
}

// Finalize finalizes the bytecode.
func (bc *ByteCode) Finalize(tempVarOffsets []int) {
	bc.tempVars = tempVarOffsets

	bc.PostProcess()

	bc.Optimize()

	bc.ResolveJumpAddrs()

	bc.ExtractLineNumbers()
}

// ClearAll clears all instructions from this bytecode instance, and cleans up.
func (bc *ByteCode) ClearAll() {
	del := bc.first

	for del != nil {
		bc.first = del.Next
		// If no references to this byte instruction is left, it should be removed by the gc.
		del = nil
		del = bc.first
	}

	bc.first = nil
	bc.last = nil

	bc.LineNumbers = make([]int, 0)
	bc.LargestStack = -1
}

// InsertIfNotExists inserts v into vars if it does not exist in vars.
func (bc *ByteCode) InsertIfNotExists(vars []int, v int) {
	for _, va := range vars {
		if va == v {
			return
		}
	}
	vars = append(vars, v)
}

// GetVarsUsed gets the variables used.
func (bc *ByteCode) GetVarsUsed(vars []int) {
	curr := bc.first
	for curr != nil {
		if ASBCInfo[curr.OPCode].Type == ASBCTYPE_wW_rW_rW_ARG {
			bc.InsertIfNotExists(vars, curr.WArg[0])
			bc.InsertIfNotExists(vars, curr.WArg[1])
			bc.InsertIfNotExists(vars, curr.WArg[2])
		} else if ASBCInfo[curr.OPCode].Type == ASBCTYPE_rW_ARG ||
			ASBCInfo[curr.OPCode].Type == ASBCTYPE_wW_ARG ||
			ASBCInfo[curr.OPCode].Type == ASBCTYPE_wW_W_ARG ||
			ASBCInfo[curr.OPCode].Type == ASBCTYPE_rW_DW_ARG ||
			ASBCInfo[curr.OPCode].Type == ASBCTYPE_wW_DW_ARG ||
			ASBCInfo[curr.OPCode].Type == ASBCTYPE_wW_QW_ARG ||
			ASBCInfo[curr.OPCode].Type == ASBCTYPE_rW_W_DW_ARG ||
			ASBCInfo[curr.OPCode].Type == ASBCTYPE_rW_DW_DW_ARG {
			bc.InsertIfNotExists(vars, curr.WArg[0])
		} else if ASBCInfo[curr.OPCode].Type == ASBCTYPE_wW_rW_ARG ||
			ASBCInfo[curr.OPCode].Type == ASBCTYPE_rW_rW_ARG ||
			ASBCInfo[curr.OPCode].Type == ASBCTYPE_wW_rW_DW_ARG {
			bc.InsertIfNotExists(vars, curr.WArg[0])
			bc.InsertIfNotExists(vars, curr.WArg[1])
		} else if curr.OPCode == ASBC_LoadThisR {
			bc.InsertIfNotExists(vars, curr.WArg[0])
		}

		curr = curr.Next
	}
}

// IsVarUsed return true if the variable at the offset is used.
func (bc *ByteCode) IsVarUsed(offset int) bool {
	curr := bc.first
	for curr != nil {
		if ASBCInfo[curr.OPCode].Type == ASBCTYPE_wW_rW_rW_ARG {
			if curr.WArg[0] == offset || curr.WArg[1] == offset || curr.WArg[2] == offset {
				return true
			}
		} else if ASBCInfo[curr.OPCode].Type == ASBCTYPE_rW_ARG ||
			ASBCInfo[curr.OPCode].Type == ASBCTYPE_wW_ARG ||
			ASBCInfo[curr.OPCode].Type == ASBCTYPE_wW_W_ARG ||
			ASBCInfo[curr.OPCode].Type == ASBCTYPE_rW_DW_ARG ||
			ASBCInfo[curr.OPCode].Type == ASBCTYPE_wW_DW_ARG ||
			ASBCInfo[curr.OPCode].Type == ASBCTYPE_wW_QW_ARG ||
			ASBCInfo[curr.OPCode].Type == ASBCTYPE_rW_W_DW_ARG ||
			ASBCInfo[curr.OPCode].Type == ASBCTYPE_rW_DW_DW_ARG {
			if curr.WArg[0] == offset {
				return true
			}
		} else if ASBCInfo[curr.OPCode].Type == ASBCTYPE_wW_rW_ARG ||
			ASBCInfo[curr.OPCode].Type == ASBCTYPE_rW_rW_ARG ||
			ASBCInfo[curr.OPCode].Type == ASBCTYPE_wW_rW_DW_ARG {
			if curr.WArg[0] == offset || curr.WArg[1] == offset {
				return true
			}
		} else if curr.OPCode == ASBC_LoadThisR {
			if offset == 0 {
				return true
			}
		}

		curr = curr.Next
	}
	return false
}

// ExchangeVar exchanges a variable offset.
func (bc *ByteCode) ExchangeVar(oldoffset, newoffset int) {
	curr := bc.first
	for curr != nil {
		if ASBCInfo[curr.OPCode].Type == ASBCTYPE_wW_rW_rW_ARG {
			if curr.WArg[0] == oldoffset {
				curr.WArg[0] = newoffset
			}
			if curr.WArg[1] == oldoffset {
				curr.WArg[1] = newoffset
			}
			if curr.WArg[2] == oldoffset {
				curr.WArg[2] = newoffset
			}
		} else if ASBCInfo[curr.OPCode].Type == ASBCTYPE_rW_ARG ||
			ASBCInfo[curr.OPCode].Type == ASBCTYPE_wW_ARG ||
			ASBCInfo[curr.OPCode].Type == ASBCTYPE_wW_W_ARG ||
			ASBCInfo[curr.OPCode].Type == ASBCTYPE_rW_DW_ARG ||
			ASBCInfo[curr.OPCode].Type == ASBCTYPE_wW_DW_ARG ||
			ASBCInfo[curr.OPCode].Type == ASBCTYPE_wW_QW_ARG ||
			ASBCInfo[curr.OPCode].Type == ASBCTYPE_rW_W_DW_ARG ||
			ASBCInfo[curr.OPCode].Type == ASBCTYPE_rW_DW_DW_ARG {
			if curr.WArg[0] == oldoffset {
				curr.WArg[0] = newoffset
			}
		} else if ASBCInfo[curr.OPCode].Type == ASBCTYPE_wW_rW_ARG ||
			ASBCInfo[curr.OPCode].Type == ASBCTYPE_rW_rW_ARG ||
			ASBCInfo[curr.OPCode].Type == ASBCTYPE_wW_rW_DW_ARG {
			if curr.WArg[0] == oldoffset {
				curr.WArg[0] = newoffset
			}
			if curr.WArg[1] == oldoffset {
				curr.WArg[1] = newoffset
			}
		}

		curr = curr.Next
	}
}

func (bc *ByteCode) AddPath(paths []ByteInstruction, instr ByteInstruction, stackSize int) {
	if instr.Marked {
		if instr.stackSize != stackSize {
			panic("Instruction stack size invalid!")
		}
		instr.Marked = true
		instr.stackSize = stackSize
		paths = append(paths, instr)
	}
}

func (bc *ByteCode) ChangeFirstDeleteNext(curr *ByteInstruction, next ByteCodeInstruction) *ByteInstruction {
	curr.OPCode = next
	if curr.Next != nil {
		curr.DeleteInstruction(curr.Next)
	}

	if curr.Previous != nil {
		return curr.Previous
	}
	return curr
}

func (bc *ByteCode) DeleteFirstChangeNext(curr *ByteInstruction, next ByteCodeInstruction) *ByteInstruction {
	if curr.Next == nil {
		panic("curr.Next was nil!")
	}
	instr := curr.Next
	instr.OPCode = next

	if instr.Previous != nil {
		return instr.Previous
	}
	return instr
}

func (bc *ByteCode) InsertBefore(before *ByteInstruction, instr *ByteInstruction) {
	if instr.Next != nil {
		panic("instruction is already in a bytecode context!")
	}
	if instr.Previous != nil {
		panic("instruction is already in a bytecode context!")
	}

	if before.Previous != nil {
		before.Previous.Next = instr
	}
	instr.Previous = before.Previous
	before.Previous = instr
	instr.Next = before

	if bc.first == before {
		first = instr
	}
}

func (bc *ByteCode) RemoveInstruction(instr *ByteInstruction) {
	if instr == bc.first {
		bc.first = bc.first.Next
	}
	if instr == bc.last {
		bc.last = bc.last.Previous
	}

	if instr.Previous != nil {
		instr.Previous.Next = instr.Next
	}
	if instr.Next != nil {
		instr.Next.Previous = instr.Previous
	}

	instr.Next = nil
	instr.Previous = nil
}

func (bc *ByteCode) CanBeSwapped(instr *ByteInstruction) bool {
	if instr.OPCode != ASBC_SwapPtr {
		panic("Pointer was not swap")
	}

	if instr.Previous != nil || instr.Previous.Previous != nil {
		return false
	}

	b := instr.Previous
	a := b.Previous

	if a.OPCode != ASBC_PshNull && a.OPCode != ASBC_PshVPtr && a.OPCode != ASBC_PSF {
		return false
	}
	if b.OPCode != ASBC_PshNull && b.OPCode != ASBC_PshVPtr && b.OPCode != ASBC_PSF {
		return false
	}

	return true
}

func (bc *ByteCode) GoBack(instr *ByteInstruction) *ByteInstruction {
	if instr == nil {
		return nil
	}
	if instr.Previous != nil {
		instr = instr.Previous
	}
	if instr.Previous != nil {
		instr = instr.Previous
	}

	return instr
}

func (bc *ByteCode) GoForward(instr *ByteInstruction) *ByteInstruction {
	if instr == nil {
		return nil
	}
	if instr.Next != nil {
		instr = instr.Next
	}
	if instr.Next != nil {
		instr = instr.Next
	}

	return instr
}

func (bc *ByteCode) PostponeInitOfTemp(curr *ByteInstruction, next **ByteInstruction) bool {
	if curr.OPCode != ASBC_SetV4 && curr.OPCode != ASBC_SetV8 || !bc.IsTemporary(curr.WArg[0]) {
		return false
	}

	use := curr.Next
	for use != nil {

		if bc.IsTempVarReadByInstr(use, curr.WArg[0]) {
			break
		}

		if bc.IsTempVarOverwrittenByInstr(use, curr.WArg[0]) {
			return false
		}

		if bc.IsInstrJmpOrLabel(usr) {
			return false
		}

		use = use.Next
	}

	if use != nil && use.Previous != curr {
		orig := curr.Next

		bc.RemoveInstruction(curr)
		bc.InsertBefore(use, curr)

		if bc.RemoveUnusedValue(curr, 0) {
			*next = orig
			return true
		}

		bc.RemoveInstruction(curr)
		bc.InsertBefore(orig, curr)
	}

	return false
}

func (bc *ByteCode) RemoveUnusedValue(curr *ByteInstruction, next **ByteInstruction) bool {
	var dummy *ByteInstruction
	if next == nil {
		next = &dummy
	}

	if curr.OPCode != ASBC_FREE &&
		(ASBCInfo[curr.OPCode].Type == ASBCTYPE_wW_rW_rW_ARG ||
			ASBCInfo[curr.OPCode].Type == ASBCTYPE_wW_rW_ARG ||
			ASBCInfo[curr.OPCode].Type == ASBCTYPE_wW_rW_DW_ARG ||
			ASBCInfo[curr.OPCode].Type == ASBCTYPE_wW_ARG ||
			ASBCInfo[curr.OPCode].Type == ASBCTYPE_wW_DW_ARG ||
			ASBCInfo[curr.OPCode].Type == ASBCTYPE_wW_QW_ARG) &&
		bc.IsTemporary(curr.WArg[0]) && !bc.IsTempVarRead(curr, curr.WArg[0]) {
		if curr.OPCode == ASBC_LdGRdR4 && bc.IsTempRegUsed(curr) {
			curr.OPCode = ASBC_LDG
			*next = bc.GoForward(curr)
			return true
		}
		*next = bc.GoForward(bc.DeleteInstruction(curr))
		return true
	}

	if curr.OPCode == ASBC_SetV4 && curr.Next != nil {
		if (curr.Next.OPCode == ASBC_CMPi ||
			curr.Next.OPCode == ASBC_CMPf ||
			curr.Next.OPCode == ASBC_CMPu) &&
			curr.WArg[0] == curr.Next.WArg[1] && bc.IsTemporary(curr.WArg[0]) && !bc.IsTempVarRead(curr.Next, curr.WArg[0]) {
			if curr.Next.OPCode == ASBC_CMPi {
				curr.Next.OPCode = ASBC_CMPIi
			}
			if curr.Next.OPCode == ASBC_CMPf {
				curr.Next.OPCode = ASBC_CMPIf
			}
			if curr.Next.OPCode == ASBC_CMPu {
				curr.Next.OPCode = ASBC_CMPIu
			}
			curr.Next.size = ASBCTypeSize[ASBCInfo[ASBC_CMPIi].Type]
			curr.Next.Arg = curr.Arg
			*next = bc.GoForward(bc.DeleteInstruction(curr))
			return true
		}

		if (curr.Next.OPCode == ASBC_ADDi ||
			curr.Next.OPCode == ASBC_SUBi ||
			curr.Next.OPCode == ASBC_MULi ||
			curr.Next.OPCode == ASBC_ADDf ||
			curr.Next.OPCode == ASBC_SUBf ||
			curr.Next.OPCode == ASBC_MULf) &&
			curr.WArg[0] == curr.Next.WArg[2] && (curr.Next.WArg[0] == curr.WArg[0] || (bc.IsTemporary(curr.WArg[0]) && !bc.IsTempVarRead(curr.Next, curr.WArg[0]))) {
			if curr.Next.OPCode == ASBC_ADDi {
				curr.Next.OPCode = ASBC_ADDIi
			}
			if curr.Next.OPCode == ASBC_SUBi {
				curr.Next.OPCode = ASBC_SUBIi
			}
			if curr.Next.OPCode == ASBC_MULi {
				curr.Next.OPCode = ASBC_MULIi
			}
			if curr.Next.OPCode == ASBC_ADDf {
				curr.Next.OPCode = ASBC_ADDIf
			}
			if curr.Next.OPCode == ASBC_SUBf {
				curr.Next.OPCode = ASBC_SUBIf
			}
			if curr.Next.OPCode == ASBC_MULf {
				curr.Next.OPCode = ASBC_MULIf
			}
			curr.Next.size = ASBCTypeSize[ASBCInfo[ASBC_ADDIi].Type]
			curr.Next.Arg = curr.Arg
			*next = bc.GoForward(bc.DeleteInstruction(curr))
			return true
		}

		if (curr.Next.OPCode == ASBC_ADDi ||
			curr.Next.OPCode == ASBC_MULi ||
			curr.Next.OPCode == ASBC_ADDf ||
			curr.Next.OPCode == ASBC_MULf) &&
			curr.WArg[0] == curr.Next.WArg[2] && (curr.Next.WArg[0] == curr.WArg[0] || (bc.IsTemporary(curr.WArg[0]) && !bc.IsTempVarRead(curr.Next, curr.WArg[0]))) {
			if curr.Next.OPCode == ASBC_ADDi {
				curr.Next.OPCode = ASBC_ADDIi
			}
			if curr.Next.OPCode == ASBC_MULi {
				curr.Next.OPCode = ASBC_SUBIi
			}
			if curr.Next.OPCode == ASBC_ADDf {
				curr.Next.OPCode = ASBC_ADDIf
			}
			if curr.Next.OPCode == ASBC_MULf {
				curr.Next.OPCode = ASBC_SUBIf
			}
			curr.Next.size = ASBCTypeSize[ASBCInfo[ASBC_ADDIi].Type]
			curr.Next.Arg = curr.Arg

			curr.Next.WArg[1] = curr.Next.WArg[2]
			*next = bc.GoForward(bc.DeleteInstruction(curr))
			return true
		}

		// The constant value is immediately moved to another variable and then not used again
		if curr.Next.OPCode == ASBC_CpyVtoV4 &&
			curr.WArg[0] == curr.Next.WArg[1] &&
			bc.IsTemporary(curr.WArg[0]) &&
			!bc.IsTempVarRead(curr.Next, curr.WArg[0]) {
			curr.WArg[0] = curr.Next.WArg[0]
			*next = bc.GoForward(bc.DeleteInstruction(curr.Next))
			return true
		}

		// The constant is copied to a temp and then immediately pushed on the stack
		if curr.Next.OPCode == ASBC_PshV4 &&
			curr.WArg[0] == curr.Next.WArg[0] &&
			bc.IsTemporary(curr.WArg[0]) &&
			!bc.IsTempVarRead(curr.Next, curr.WArg[0]) {
			curr.OPCode = ASBC_PshC4
			curr.stackInc = ASBCInfo[ASBC_PshC4].StackIncrement
			*next = bc.GoForward(bc.DeleteInstruction(curr.Next))
			return true
		}

		// The constant is copied to a global variable and then never used again
		if curr.Next.OPCode == ASBC_CpyVtoG4 &&
			curr.WArg[0] == curr.Next.WArg[0] &&
			bc.IsTemporary(curr.WArg[0]) &&
			!bc.IsTempVarRead(curr.Next, curr.WArg[0]) {
			curr.OPCode = ASBC_SetG4
			curr.size = ASBCTypeSize[ASBCInfo[ASBC_SetG4].Type]
			curr.Arg = curr.Next.Arg
			*next = bc.GoForward(bc.DeleteInstruction(curr.Next))
			return true
		}
	}

	// The value is immediately moved to another variable and then not used again
	if (ASBCInfo[curr.OPCode].Type == ASBCTYPE_wW_rW_rW_ARG ||
		ASBCInfo[curr.OPCode].Type == ASBCTYPE_wW_rW_DW_ARG) &&
		curr.Next != nil &&
		curr.Next.OPCode == ASBC_CpyVtoV4 &&
		curr.WArg[0] == curr.Next.WArg[1] &&
		bc.IsTemporary(curr.WArg[0]) &&
		!bc.IsTempVarRead(curr.Next, curr.WArg[0]) {
		curr.WArg[0] = curr.Next.WArg[0]
		*next = bc.GoForward(bc.DeleteInstruction(curr.Next))
		return true
	}

	// The register is copied to a temp variable and then back to the register again without being used afterwards
	if curr.OPCode == ASBC_CpyRtoV4 && curr.Next != nil && curr.Next.OPCode == ASBC_CpyVtoR4 &&
		curr.WArg[0] == curr.Next.WArg[0] &&
		bc.IsTemporary(curr.WArg[0]) &&
		!bc.IsTempVarRead(curr.Next, curr.WArg[0]) {
		bc.DeleteInstruction(curr.Next)
		*next = bc.GoForward(bc.DeleteInstruction(curr))
		return true
	}

	// The global value is copied to a temp and then immediately pushed on the stack
	if curr.OPCode == ASBC_CpyGtoV4 && curr.Next != nil && curr.Next.OPCode == ASBC_PshV4 &&
		curr.WArg[0] == curr.Next.WArg[0] &&
		bc.IsTemporary(curr.WArg[0]) &&
		!bc.IsTempVarRead(curr.Next, curr.WArg[0]) {
		curr.OPCode = ASBC_PshG4
		curr.size = ASBCTypeSize[ASBCInfo[ASBC_PshG4].Type]
		curr.stackInc = ASBCInfo[ASBC_PshG4].StackIncrement
		*next = bc.GoForward(bc.DeleteInstruction(curr.Next))
		return true
	}

	// The constant is assigned to a variable, then the value of the variable
	// pushed on the stack, and then the variable is never used again
	if curr.OPCode == ASBC_SetV8 && curr.Next != nil && curr.Next.OPCode == ASBC_PshV8 &&
		curr.WArg[0] == curr.Next.WArg[0] &&
		bc.IsTemporary(curr.WArg[0]) &&
		!bc.IsTempVarRead(curr.Next, curr.WArg[0]) {
		curr.OPCode = ASBC_PshC8
		curr.stackInc = ASBCInfo[ASBC_PshC8].StackIncrement
		*next = bc.GoForward(bc.DeleteInstruction(curr.Next))
		return true
	}

	return false
}

func (bc *ByteCode) IsTemporary(offset int) bool {
	return bc.TempVarsExists(offset)
}

func (bc *ByteCode) OptimizeLocally(tempVarOffsets []int) {
	bc.tempVars = tempVarOffsets

	instr := bc.last
	for instr != nil {
		curr := instr
		instr = instr.Previous

		// Remove instructions when the result is not used anywhere
		// This will return true if the instruction is deleted, and
		// false if it is not deleted. Observe that the instruction
		// can be modified.
		if bc.RemoveUnusedValue(curr, &instr) {
			continue
		}

		// Postpone initializations so that they may be combined in the second pass.
		// If the initialization is postponed, then the optimizations should continue
		// from where the value was used, so instr will be updated to point to that.
		if bc.PostponeInitOfTemp(curr, &instr) {
			continue
		}

		currOp := curr.OPCode
		if currOp == ASBC_SwapPtr {
			if bc.CanBeSwapped(curr) {
				bc.DeleteInstruction(curr)

				a := instr.Previous
				bc.RemoveInstruction(instr)
				bc.InsertBefore(a, instr)

				instr = bc.GoForward(a)
				continue
			}
		} else if currOp == ASBC_ClrHi {
			if instr != nil &&
				(instr.OPCode == ASBC_TZ ||
					instr.OPCode == ASBC_TNZ ||
					instr.OPCode == ASBC_TS ||
					instr.OPCode == ASBC_TNS ||
					instr.OPCode == ASBC_TP ||
					instr.OPCode == ASBC_TNP) {
				instr = bc.GoForward(bc.DeleteInstruction(curr))
				continue
			}

			if curr.Next != nil && curr.Next.OPCode == ASBC_JZ {
				curr.Next.OPCode = ASBC_JLowZ
				instr = bc.GoForward(bc.DeleteInstruction(curr))
				continue
			}

			if curr.Next != nil && curr.Next.OPCode == ASBC_JNZ {
				curr.Next.OPCode = ASBC_JLowNZ
				instr = bc.GoForward(bc.DeleteInstruction(curr))
				continue
			}
		} else if currOp == ASBC_LDV && curr.Next != nil {
			if curr.Next.OPCode == ASBC_INCi && !bc.isTempRegUsed(curr.Next) {
				curr.OPCode = ASBC_IncVi
				bc.DeleteInstruction(curr.Next)
				instr = bc.GoForward(curr)
			}

			if curr.Next.OPCode == ASBC_DECi && !bc.isTempRegUsed(curr.Next) {
				curr.OPCode = ASBC_DecVi
				bc.DeleteInstruction(curr.Next)
				instr = bc.GoForward(curr)
			}
		} else if currOp == ASBC_LDG && curr.Next != nil {
			if curr.Next.OPCode == ASBC_WRTV4 && !bc.IsTempRegUsed(curr.Next) {
				curr.OPCode = ASBC_CpyVtoG4
				curr.size = ASBCTypeSize[ASBCInfo[ASBC_CpyVtoG4].Type]
				curr.WArg[0] = curr.Next.WArg[0]
				bc.DeleteInstruction(curr.Next)
				instr = bc.GoForward(curr)
			}

			else if curr.Next.OPCode == ASBC_RDR4 {
				if !bc.IsTempRegUsed(curr.Next) {
					curr.OPCode = ASBC_CpyGtoV4
				} else {
					curr.OPCode = ASBC_LdGRdR4
				}
				curr.size = ASBCTypeSize[ASBCInfo[ASBC_CpyGtoV4].Type]
				curr.WArg[0] = curr.Next.WArg[0]
				bc.DeleteInstruction(curr.Next)
				instr = bc.GoForward(curr)
			}
		} else if currOp == ASBC_CHKREF {

		}
	}
}


func (bc *ByteCode) DeleteInstruction(instr *ByteInstruction) *ByteInstruction {
	if instr == nil { return nil }
	ret := instr.Previous
	if instr.Previous == nil {
		ret = instr.Next
	}
	bc.RemoveInstruction(instr)

	return ret
}

/*
func (bc *ByteCode) () {

}
*/
