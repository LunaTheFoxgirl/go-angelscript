package angelscript

import (
	"sync/atomic"
)

// DAtomic provides functions for handling ASDWORDs threadsafetly.
type DAtomic struct {
	val ASDWORD
}

func (da *DAtomic) Get() ASDWORD {
	return atomic.LoadUint32(&da.val)
}

func (da *DAtomic) Set(val ASDWORD) {
	atomic.StoreUint32(&da.val, val)
}

func (da *DAtomic) Inc() ASDWORD {
	da.Set(da.Get()+1)
	return da.Get()
}

func (da *DAtomic) Dec() ASDWORD {
	da.Set(da.Get()-1)
	return da.Get()
}