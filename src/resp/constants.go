package resp

const (
	SimpleStringOp = '+'
	ErrorOp        = '-'
	IntegerOp      = ':'
	BulkStringOp   = '$'
	ArrayOp        = '*'
)

const readBufferSize = 512 // bytes

const (
	Get = "GET"
	Set = "SET"
)
