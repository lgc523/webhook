package resp

const (
	SUCCESS = 200

	ILLEGAL = 400

	NotExist = 800

	DEFAULT = -1
)

var MsgMapping = map[int]string{
	SUCCESS:  "OK",
	ILLEGAL:  "请求不合法",
	NotExist: "资源不存在",
	DEFAULT:  "阿里嘎多",
}

func DefaultMsg() string {
	return Msg(DEFAULT)
}

func Msg(code int) string {
	if msg, ok := MsgMapping[code]; ok {
		return msg
	}
	return DefaultMsg()
}
