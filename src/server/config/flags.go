package config

type flag int

const (
	singleFlag flag = 1 << iota
	multipleFlag
)

type option struct {
	name  string
	flags flag
}

func (o option) hasFlag(check flag) bool {
	return o.flags&check > 0
}
