package terrors

type TError int

const (
	// common
	NotFound            TError = iota + 10000
	DatabaseError              // 数据库错误
	ParametersAreWrong         // 参数错误
	CidToHashFiled             // cid转hash出错
	EncodingError              // 编码错误
	DecodingError              // 解码错误
	NotFoundUser               // 找不到用户
	InsufficientBalance        // 余额不足
	NotAdministrator           // 不是管理员
	NotFoundSignCode           // 找不到签名码
	SignError                  // 签名错误
	NotFoundAddress            // 没有可用地址

	Success = 0
	Unknown = -1
)

func (e TError) Int() int {
	return int(e)
}

func (e TError) String() string {
	switch e {
	case DatabaseError:
		return "database error"
	case NotAdministrator:
		return "not a administrator"
	case InsufficientBalance:
		return "user insufficient balance"
	case NotFoundSignCode:
		return "not found code"
	case NotFoundUser:
		return "not found user"
	case SignError:
		return "sign error"
	case NotFoundAddress:
		return "not found address"
	default:
		return ""
	}
}
