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
	Ping   = "PING"
	Config = "CONFIG"
	Get    = "GET"
	Set    = "SET"
	DBSize = "DBSIZE"
	Del    = "DEL"
)
