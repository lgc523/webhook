package resp

type Resp struct {
	Code int
	Msg  string
	Data any
}

func Success(data any) *Resp {
	return &Resp{
		Code: SUCCESS,
		Msg:  Msg(SUCCESS),
		Data: data,
	}
}

func ClientFail(msg string) *Resp {
	return &Resp{
		Code: DEFAULT,
		Msg:  msg,
	}
}

func ClientFailWithCode(code int, msg string) *Resp {
	return &Resp{
		Code: code,
		Msg:  msg,
	}
}
func ServerFail(msg string) *Resp {
	return &Resp{
		Code: DEFAULT,
		Msg:  msg,
	}
}
func ServerFailWithCode(code int, msg string) *Resp {
	return &Resp{
		Code: code,
		Msg:  msg,
	}
}
