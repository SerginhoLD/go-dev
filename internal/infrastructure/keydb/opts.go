package keydb

type option interface {
	typ() string
	val() any
}

type MaxLen int64

func (opt MaxLen) typ() string {
	return "maxlen"
}

func (opt MaxLen) val() any {
	return int64(opt)
}

type Header struct {
	Key string
	Val string
}

func (opt Header) typ() string {
	return "header"
}

func (opt Header) val() any {
	return map[string]string{opt.Key: opt.Val}
}

type StartId string

func (opt StartId) typ() string {
	return "start_id"
}

func (opt StartId) val() any {
	return string(opt)
}

type DelAfterAck bool

func (opt DelAfterAck) typ() string {
	return "del_after_ack"
}

func (opt DelAfterAck) val() any {
	return bool(opt)
}
