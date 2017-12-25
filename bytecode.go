package angelscript

import ()

type EBCInstruction byte

// Byte code instructions
const (
	ASBC_PopPtr			= EBCInstruction(iota)
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

func (bi *ByteInstruction) AddAfter(next *ByteInstruction) {
	if bi.Next != nil {
		bi.Previous = next
	}
	
	next.Next = bi.Next
	next.Previous = bi
	bi.Next = next
}

func (bi *ByteInstruction) AddBefore(prev *ByteInstruction) {
	if bi.Previous != nil {
		bi.Next = prev
	}
	
	prev.Previous = bi.Next
	prev.Next = bi
	bi.Previous = prev
}

func (bi *ByteInstruction) GetSize() int {
	return bi.size
}

func (bi *ByteInstruction) GetStackIncrease() int {
	return bi.stackInc
}

func (bi *ByteInstruction) Remove() {
	if bi.Previous != nil { bi.Previous.Next = bi.Next }
	if bi.Next != nil { bi.Next.Previous = bi.Previous }
	bi.Next = nil
	bi.Previous = nil
}

type ByteCode struct {
	lineNumbers  []int
	sectionIdxs  []int
	largestStack int
}
