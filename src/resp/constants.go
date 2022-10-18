package resp

const (
	SimpleStringOp = '+'
	ErrorOp        = '-'
	IntegersOp     = ':'
	BulkStringOp   = '$'
	ArrayOp        = '*'
)

const readBufferSize = 512 // bytes

const (
	Get = "GET"
	Set = "SET"
)
