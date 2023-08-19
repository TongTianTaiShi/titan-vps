package terrors

const (
	// common
	NotFound            = iota + 10000
	DatabaseError       // 数据库错误
	ParametersAreWrong  // 参数错误
	CidToHashFiled      // cid转hash出错
	EncodingError       // 编码错误
	DecodingError       // 解码错误
	UserNotFound        // 找不到用户
	InsufficientBalance // 余额不足

	Success = 0
	Unknown = -1
)
