package angelscript

import (
	"strconv"
)

// InstructionType is the type of an bytecode instruction.
type InstructionType byte

const 
(
	ASBCTYPE_INFO         = itoa(InstructionType)
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
)

// EBCInstruction is an bytecode instruction.
// See all ASBC_ consts for bytecode instructions.
type ByteCodeInstruction byte

// Byte code instructions
const (
	ASBC_PopPtr			= ByteCodeInstruction(iota)
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
	ASBC_CASt			
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
	ASBC_VarDecl		
	ASBC_Block			
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
func asbcinfo(inst ByteInstruction, t InstructionType, s int, name string) ByteCodeInfo {
	return ByteCodeInfo{inst, t, s, name}
}

// helper function to define dummy bytecode entries.
func asbcdummy(id int) ByteCodeInfo {
	return asbcinfo(ASBC_MAXBYTECODE, ASBCTYPE_INFO, 0, "<bc DUMMY " + strconv.Itoa(id) + ">")
}
// helper function to define dummy bytecode entries.
func asbcdummyrange(from, to int) []ByteCodeInfo {
	out := make([]ByteCodeInfo, 0)
	for i := from; i <= to; i++ {
		out = append(out, asbcdummy(i))
	}
	return out
}

const ASBCInfo []ByteCodeInfo = []ByteCodeInfo {
	asbcinfo(ASBC_PopPtr, ASBCTYPE_NO_ARG, -ASPointerSize, "<pop pointer>"),
	asbcinfo(ASBC_PshGPtr, ASBCTYPE_PTR_ARG, ASPointerSize, "<push pointer>"),
	asbcinfo(ASBC_PshC4, ASBCTYPE_DW_ARG, 1, "<push C4>"),
	asbcinfo(ASBC_PshV4, ASBCTYPE_rW_ARG, 1, "<push V4>"),
	asbcinfo(ASBC_PSF, ASBCTYPE_rW_ARG, ASPointerSize, "<psf>"),
	asbcinfo(ASBC_SwapPtr, ASBCTYPE_NO_ARG, 0, "<swap pointer>"),
	asbcinfo(ASBC_NOT, ASBCTYPE_rW_ARG, 0, "<not>"),
	asbcinfo(ASBC_PshG4, ASBCTYPE_PTR_ARG, 1, "<push G4>"),
	asbcinfo(ASBC_LdGRdR4, ASBCTYPE_wW_PTR_ARG, 0, "<LdGRdR4>"),
	asbcinfo(ASBC_CALL, ASBCTYPE_ARG, 0xFFFF, "<call>"),
	asbcinfo(ASBC_RET, ASBCTYPE_W_ARG, 0xFFFF, "<return>"),
	
	//jumps
	asbcinfo(ASBC_JMP, ASBCTYPE_DW_ARG, 0, "<jump>"),
	asbcinfo(ASBC_JZ, ASBCTYPE_DW_ARG, 0, "<jump if 0>"),
	asbcinfo(ASBC_JNZ, ASBCTYPE_DW_ARG, 0, "<jump if not 0>"),
	asbcinfo(ASBC_JS, ASBCTYPE_DW_ARG, 0, "<jump if sign>"),
	asbcinfo(ASBC_JNS, ASBCTYPE_DW_ARG, 0, "<jump if not sign>"),
	asbcinfo(ASBC_JP, ASBCTYPE_DW_ARG, , ""),
	asbcinfo(ASBC_JNP, ASBCTYPE_DW_ARG, , ""),
	asbcinfo(ASBC_TZ, ASBCTYPE_NO_ARG, , ""),
	asbcinfo(ASBC_TNZ, ASBCTYPE_NO_ARG, , ""),
	asbcinfo(ASBC_TS, ASBCTYPE_NO_ARG, , ""),
	asbcinfo(ASBC_TNS, ASBCTYPE_NO_ARG, , ""),
	asbcinfo(ASBC_TP, ASBCTYPE_NO_ARG, , ""),
	asbcinfo(ASBC_TNP, ASBCTYPE_NO_ARG, , ""),
	
	// NEGs
	asbcinfo(ASBC_NEGi, ASBCTYPE_rW_ARG, , ""),
	asbcinfo(ASBC_NEGf, ASBCTYPE_rW_ARG, , ""),
	asbcinfo(ASBC_NEGd, ASBCTYPE_rW_ARG, , ""),
	
	// Inc and Dec
	asbcinfo(ASBC_INCi16, ASBCTYPE_NO_ARG, , ""),
	asbcinfo(ASBC_INCi8, ASBCTYPE_NO_ARG, , ""),
	asbcinfo(ASBC_DECi16, ASBCTYPE_NO_ARG, , ""),
	asbcinfo(ASBC_DECi8, ASBCTYPE_NO_ARG, , ""),
	asbcinfo(ASBC_INCi, ASBCTYPE_NO_ARG, , ""),
	asbcinfo(ASBC_DECi, ASBCTYPE_NO_ARG, , ""),
	asbcinfo(ASBC_INCf, ASBCTYPE_NO_ARG, , ""),
	asbcinfo(ASBC_DECf, ASBCTYPE_NO_ARG, , ""),
	asbcinfo(ASBC_INCd, ASBCTYPE_NO_ARG, , ""),
	asbcinfo(ASBC_DECd, ASBCTYPE_NO_ARG, , ""),
	
	//Special cases?
	asbcinfo(ASBC_IncVi, ASBCTYPE_rW_ARG, , ""),
	asbcinfo(ASBC_DecVi, ASBCTYPE_rW_ARG, , ""),
	asbcinfo(ASBC_BNOT, ASBCTYPE_rW_ARG, , ""),
	
	//B Operators
	asbcinfo(ASBC_BAND, ASBCTYPE_wW_rW_rW_ARG, , ""),
	asbcinfo(ASBC_BOR, ASBCTYPE_wW_rW_rW_ARG, , ""),
	asbcinfo(ASBC_BSLL, ASBCTYPE_wW_rW_rW_ARG, , ""),
	asbcinfo(ASBC_BSRL, ASBCTYPE_wW_rW_rW_ARG, , ""),
	asbcinfo(ASBC_BSRA, ASBCTYPE_wW_rW_rW_ARG, , ""),
	
	//Copy
	asbcinfo(ASBC_COPY, ASBCTYPE_DW_ARG, , ""),
	
	//More pointer stuff
	asbcinfo(ASBC_PshC8, ASBCTYPE_QW_ARG, , ""),
	asbcinfo(ASBC_PshVPtr, ASBCTYPE_rW_ARG, , ""),
	asbcinfo(ASBC_RDSPtr, ASBCTYPE_NO_ARG, , ""),
	
	//Compare
	asbcinfo(ASBC_CMPd, ASBCTYPE_rW_rW_ARG, , ""),
	asbcinfo(ASBC_CMPu, ASBCTYPE_rW_rW_ARG, , ""),
	asbcinfo(ASBC_CMPf, ASBCTYPE_rW_rW_ARG, , ""),
	asbcinfo(ASBC_CMPi, ASBCTYPE_rW_rW_ARG, , ""),
	asbcinfo(ASBC_CMPIi, ASBCTYPE_DW_ARG, , ""),
	asbcinfo(ASBC_CMPIf, ASBCTYPE_DW_ARG, , ""),
	asbcinfo(ASBC_CMPIu, ASBCTYPE_DW_ARG, , ""),
	
	//EVEN MORE pointer stuff, and jump?
	asbcinfo(ASBC_JMPP, ASBCTYPE_rW_ARG, , ""),
	asbcinfo(ASBC_PopRPtr, ASBCTYPE_NO_ARG, , ""),
	asbcinfo(ASBC_PshRPtr, ASBCTYPE_NO_ARG, , ""),
	
	//string and essential stuff
	asbcinfo(ASBC_STR, ASBCTYPE_W_ARG, , ""),
	asbcinfo(ASBC_CALLSYS, ASBCTYPE_DW_ARG, , ""),
	asbcinfo(ASBC_CALLBND, ASBCTYPE_DW_ARG, , ""),
	asbcinfo(ASBC_SUSPEND, ASBCTYPE_NO_ARG, , ""),
	asbcinfo(ASBC_ALLOC, ASBCTYPE_PTR_DW_ARG, , ""),
	asbcinfo(ASBC_FREE, ASBCTYPE_wW_PTR_ARG, , ""),
	asbcinfo(ASBC_LOADOBJ, ASBCTYPE_rW_ARG, , ""),
	asbcinfo(ASBC_STOREOBJ, ASBCTYPE_wW_ARG, , ""),
	asbcinfo(ASBC_GETOBJ, ASBCTYPE_W_ARG, , ""),
	asbcinfo(ASBC_REFCPY, ASBCTYPE_W_ARG, , ""),
	asbcinfo(ASBC_CHKREF, ASBCTYPE_NO_ARG, , ""),
	asbcinfo(ASBC_GETOBJREF, ASBCTYPE_W_ARG, , ""),
	asbcinfo(ASBC_GETREF, ASBCTYPE_W_ARG, , ""),
	asbcinfo(ASBC_PshNull, ASBCTYPE_NO_ARG, , ""),
	asbcinfo(ASBC_ClrVptr, ASBCTYPE_wW_ARG, , ""),
	asbcinfo(ASBC_OBJTYPE, ASBCTYPE_PTR_ARG, , ""),
	asbcinfo(ASBC_TYPEID, ASBCTYPE_DW_ARG, , ""),
	
	asbcinfo(ASBC_SetV4, ASBCTYPE_wW_DW_ARG, , ""),
	asbcinfo(ASBC_SetV8, ASBCTYPE_wW_QW_ARG, , ""),
	asbcinfo(ASBC_ADDSi, ASBCTYPE_W_DW_ARG, , ""),
	
	//Copy X to Y(4/8)
	asbcinfo(ASBC_CpyVtoV4, ASBCTYPE_wW_RW_ARG, , ""),
	asbcinfo(ASBC_CpyVtoV8, ASBCTYPE_wW_RW_ARG, , ""),
	asbcinfo(ASBC_CpyVtoR4, ASBCTYPE_rW_ARG, , ""),
	asbcinfo(ASBC_CpyVtoR8, ASBCTYPE_rW_ARG, , ""),
	asbcinfo(ASBC_CpyVtoG4, ASBCTYPE_rW_PTR_ARG, , ""),
	asbcinfo(ASBC_CpyRtoV4, ASBCTYPE_rW_ARG, , ""),
	asbcinfo(ASBC_CpyRtoV8, ASBCTYPE_rW_ARG, , ""),
	asbcinfo(ASBC_CpyGtoV4, ASBCTYPE_rW_PTR_ARG, , ""),
	
	asbcinfo(ASBC_WRTV1, ASBCTYPE_rW_ARG, , ""),
	asbcinfo(ASBC_WRTV2, ASBCTYPE_rW_ARG, , ""),
	asbcinfo(ASBC_WRTV4, ASBCTYPE_rW_ARG, , ""),
	asbcinfo(ASBC_WRTV8, ASBCTYPE_rW_ARG, , ""),
	asbcinfo(ASBC_RDR1, ASBCTYPE_rW_ARG, , ""),
	asbcinfo(ASBC_RDR2, ASBCTYPE_rW_ARG, , ""),
	asbcinfo(ASBC_RDR4, ASBCTYPE_rW_ARG, , ""),
	asbcinfo(ASBC_RDR8, ASBCTYPE_rW_ARG, , ""),
	
	asbcinfo(ASBC_LDG, ASBCTYPE_, , ""),
	asbcinfo(ASBC_LDV, ASBCTYPE_, , ""),
	asbcinfo(ASBC_PGA, ASBCTYPE_, , ""),
	asbcinfo(ASBC_CmpPtr, ASBCTYPE_, , ""),
	asbcinfo(ASBC_VAR, ASBCTYPE_, , ""),
	asbcinfo(ASBC_iTOf, ASBCTYPE_, , ""),
	asbcinfo(ASBC_fTOi, ASBCTYPE_, , ""),
	asbcinfo(ASBC_uTOf, ASBCTYPE_, , ""),
	asbcinfo(ASBC_fTOu, ASBCTYPE_, , ""),
	asbcinfo(ASBC_sbTOi, ASBCTYPE_, , ""),
	asbcinfo(ASBC_swTOi, ASBCTYPE_, , ""),
	asbcinfo(ASBC_ubTOi, ASBCTYPE_, , ""),
	asbcinfo(ASBC_uwTOi, ASBCTYPE_, , ""),
	asbcinfo(ASBC_dTOi, ASBCTYPE_, , ""),
	asbcinfo(ASBC_dTOu, ASBCTYPE_, , ""),
	asbcinfo(ASBC_dTOf, ASBCTYPE_, , ""),
	asbcinfo(ASBC_iTOd, ASBCTYPE_, , ""),
	asbcinfo(ASBC_uTOd, ASBCTYPE_, , ""),
	asbcinfo(ASBC_fTOd, ASBCTYPE_, , ""),
	asbcinfo(ASBC_ADDi, ASBCTYPE_, , ""),
	asbcinfo(ASBC_SUBi, ASBCTYPE_, , ""),
	asbcinfo(ASBC_MULi, ASBCTYPE_, , ""),
	asbcinfo(ASBC_DIVi, ASBCTYPE_, , ""),
	asbcinfo(ASBC_MODi, ASBCTYPE_, , ""),
	asbcinfo(ASBC_ADDf, ASBCTYPE_, , ""),
	asbcinfo(ASBC_SUBf, ASBCTYPE_, , ""),
	asbcinfo(ASBC_MULf, ASBCTYPE_, , ""),
	asbcinfo(ASBC_DIVf, ASBCTYPE_, , ""),
	asbcinfo(ASBC_MODf, ASBCTYPE_, , ""),
	asbcinfo(ASBC_ADDd, ASBCTYPE_, , ""),
	asbcinfo(ASBC_SUBd, ASBCTYPE_, , ""),
	asbcinfo(ASBC_MULd, ASBCTYPE_, , ""),
	asbcinfo(ASBC_DIVd, ASBCTYPE_, , ""),
	asbcinfo(ASBC_MODd, ASBCTYPE_, , ""),
	asbcinfo(ASBC_ADDIi, ASBCTYPE_, , ""),
	asbcinfo(ASBC_SUBIi, ASBCTYPE_, , ""),
	asbcinfo(ASBC_MULIi, ASBCTYPE_, , ""),
	asbcinfo(ASBC_ADDIf, ASBCTYPE_, , ""),
	asbcinfo(ASBC_SUBIf, ASBCTYPE_, , ""),
	asbcinfo(ASBC_MULIf, ASBCTYPE_, , ""),
	asbcinfo(ASBC_SetG4, ASBCTYPE_, , ""),
	asbcinfo(ASBC_ChkRefS, ASBCTYPE_, , ""),
	asbcinfo(ASBC_ChkNullV, ASBCTYPE_, , ""),
	asbcinfo(ASBC_CALLINTF, ASBCTYPE_, , ""),
	asbcinfo(ASBC_iTOb, ASBCTYPE_, , ""),
	asbcinfo(ASBC_iTOw, ASBCTYPE_, , ""),
	asbcinfo(ASBC_SetV1, ASBCTYPE_, , ""),
	asbcinfo(ASBC_SetV2, ASBCTYPE_, , ""),
	asbcinfo(ASBC_Cast, ASBCTYPE_, , ""),
	asbcinfo(ASBC_i64Toi, ASBCTYPE_, , ""),
	asbcinfo(ASBC_uTOi64, ASBCTYPE_, , ""),
	asbcinfo(ASBC_iTOi64, ASBCTYPE_, , ""),
	asbcinfo(ASBC_fTOi64, ASBCTYPE_, , ""),
	asbcinfo(ASBC_dTOi64, ASBCTYPE_, , ""),
	asbcinfo(ASBC_uTOu64, ASBCTYPE_, , ""),
	asbcinfo(ASBC_iTOu64, ASBCTYPE_, , ""),
	asbcinfo(ASBC_fTOu64, ASBCTYPE_, , ""),
	asbcinfo(ASBC_dTOu64, ASBCTYPE_, , ""),
	asbcinfo(ASBC_i64TOf, ASBCTYPE_, , ""),
	asbcinfo(ASBC_u64TOf, ASBCTYPE_, , ""),
	asbcinfo(ASBC_i64TOd, ASBCTYPE_, , ""),
	asbcinfo(ASBC_u64TOd, ASBCTYPE_, , ""),
	asbcinfo(ASBC_NEGi64, ASBCTYPE_, , ""),
	asbcinfo(ASBC_INCi64, ASBCTYPE_, , ""),
	asbcinfo(ASBC_DECi64, ASBCTYPE_, , ""),
	asbcinfo(ASBC_BNOT64, ASBCTYPE_, , ""),
	asbcinfo(ASBC_ADDi64, ASBCTYPE_, , ""),
	asbcinfo(ASBC_SUBi64, ASBCTYPE_, , ""),
	asbcinfo(ASBC_MULi64, ASBCTYPE_, , ""),
	asbcinfo(ASBC_DIVi64, ASBCTYPE_, , ""),
	asbcinfo(ASBC_MODi64, ASBCTYPE_, , ""),
	asbcinfo(ASBC_BAND64, ASBCTYPE_, , ""),
	asbcinfo(ASBC_BOR64, ASBCTYPE_, , ""),
	asbcinfo(ASBC_BXOR64, ASBCTYPE_, , ""),
	asbcinfo(ASBC_BSLL64, ASBCTYPE_, , ""),
	asbcinfo(ASBC_BSRL64, ASBCTYPE_, , ""),
	asbcinfo(ASBC_BSRA64, ASBCTYPE_, , ""),
	asbcinfo(ASBC_CMPi64, ASBCTYPE_, , ""),
	asbcinfo(ASBC_CMPu64, ASBCTYPE_, , ""),
	asbcinfo(ASBC_ChkNullS, ASBCTYPE_, , ""),
	asbcinfo(ASBC_ClrHi, ASBCTYPE_, , ""),
	asbcinfo(ASBC_JitEntry, ASBCTYPE_, , ""),
	asbcinfo(ASBC_CallPtr, ASBCTYPE_, , ""),
	asbcinfo(ASBC_FuncPtr, ASBCTYPE_, , ""),
	asbcinfo(ASBC_LoadThisR, ASBCTYPE_, , ""),
	asbcinfo(ASBC_PshV8, ASBCTYPE_, , ""),
	asbcinfo(ASBC_DIVu, ASBCTYPE_, , ""),
	asbcinfo(ASBC_MODu, ASBCTYPE_, , ""),
	asbcinfo(ASBC_DIVu64, ASBCTYPE_, , ""),
	asbcinfo(ASBC_MODu64, ASBCTYPE_, , ""),
	asbcinfo(ASBC_LoadRObjR, ASBCTYPE_, , ""),
	asbcinfo(ASBC_LoadVObjR, ASBCTYPE_, , ""),
	asbcinfo(ASBC_RefCpyV, ASBCTYPE_, , ""),
	asbcinfo(ASBC_JLowZ, ASBCTYPE_, , ""),
	asbcinfo(ASBC_JLowNZ, ASBCTYPE_, , ""),
	asbcinfo(ASBC_AllocMem, ASBCTYPE_, , ""),
	asbcinfo(ASBC_SetListSize, ASBCTYPE_, , ""),
	asbcinfo(ASBC_PshListElmnt, ASBCTYPE_, , ""),
	asbcinfo(ASBC_SetListType, ASBCTYPE_, , ""),
	asbcinfo(ASBC_POWi, ASBCTYPE_, , ""),
	asbcinfo(ASBC_POWu, ASBCTYPE_, , ""),
	asbcinfo(ASBC_POWf, ASBCTYPE_, , ""),
	asbcinfo(ASBC_POWd, ASBCTYPE_, , ""),
	asbcinfo(ASBC_POWdi, ASBCTYPE_, , ""),
	asbcinfo(ASBC_POWi64, ASBCTYPE_, , ""),
	asbcinfo(ASBC_POWu64, ASBCTYPE_, , ""),
	asbcinfo(ASBC_Thiscall1, ASBCTYPE_, , ""),
	
	asbcdummyrange(201, 250),
	
	asbcinfo(ASBC_VarDecl, ASBCTYPE_W_ARG, 0, "<var decl>"),
	asbcinfo(ASBC_Block, ASBCTYPE_INFO, 0, "<block info>"),
	asbcinfo(ASBC_ObjInfo, ASBCTYPE_rW_DW_ARG, 0, "<obj info>"),
	asbcinfo(ASBC_LINE, ASBCTYPE_INFO, 0, "<line info>"),
	asbcinfo(ASBC_LABEL, ASBCTYPE_INFO, 0, "<label info>"),
	
}

// ByteInstruction is an instruction in the bytecode.
type ByteInstruction struct {
	OPCode EBCInstruction
	Arg uint64
	WArg [3]uint16
	size int
	stackInc int
	
	Marked bool
	stackSize int
	
	Next *ByteInstruction
	Previous *ByteInstruction
}

// NewByteInstruction creates a new instance of ByteInstruction.
func NewByteInstruction() *ByteInstruction {
	inst := ByteInstruction{}
	inst.Next = nil
	inst.Previous = nil
	
	inst.OPCode = ASBC_LABEL
	
	inst.Arg = 0
	inst.WArg = [3]uint16{0,0,0}
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
	if bi.Previous != nil { bi.Previous.Next = bi.Next }
	if bi.Next != nil { bi.Next.Previous = bi.Previous }
	bi.Next = nil
	bi.Previous = nil
}

type ByteCode struct {
	LineNumbers  []int
	SectionIdxs  []int
	LargestStack int
	
	engine *ScriptEngine
	first *ByteInstruction
	last *ByteInstruction
	tempVars []int 
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
func (bc *ByteCode) InsertIfNotExists(vars *[]int, v int) {
	for _, va := range vars {
		if va == v { return }
	}
	vars = append(vars, v)
}
// GetVarsUsed gets the variables used.
func (bc *ByteCode) GetVarsUsed(vars *[]int) {
	curr := bc.first
	for curr != nil {
		if bc.
	}
}