package terrors

type TError int

const (
	// common
	NotFound            TError = iota + 10000
	DatabaseError              // 数据库错误
	ParametersWrong            // 参数错误
	CidToHashFiled             // cid转hash出错
	EncodingError              // 编码错误
	DecodingError              // 解码错误
	NotFoundUser               // 找不到用户
	InsufficientBalance        // 余额不足
	NotAdministrator           // 不是管理员
	NotFoundSignCode           // 找不到签名码
	SignError                  // 签名错误
	NotFoundAddress            // 没有可用地址
	NotFoundOrder              // 找不到订单
	StateMachinesError         // 订单状态机出错
	DescribePriceError         // 询价出错
	UserMismatch               // 用户不匹配

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
	case ParametersWrong:
		return "parameters are wrong"
	case CidToHashFiled:
		return "error converting cid to hash"
	case EncodingError:
		return "encoding error"
	case DecodingError:
		return "decoding error"
	case NotFoundUser:
		return "user not found"
	case InsufficientBalance:
		return "user insufficient balance"
	case NotAdministrator:
		return "not an administrator"
	case NotFoundSignCode:
		return "signature code not found"
	case SignError:
		return "signature error"
	case NotFoundAddress:
		return "no address available"
	case NotFoundOrder:
		return "order not found"
	case StateMachinesError:
		return "order state machine error"
	case DescribePriceError:
		return "describe price error"
	case UserMismatch:
		return "user mismatch"
	default:
		return ""
	}
}
