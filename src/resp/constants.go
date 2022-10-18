package resp

const (
	simpleStringOp = '+'
	errorOp        = '-'
	integersOp     = ':'
	bulkStringOp   = '$'
	arrayOp        = '*'
)

const readBufferSize = 512 // bytes
