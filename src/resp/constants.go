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
	Ping = "PING"
	Get  = "GET"
	Set  = "SET"
)
