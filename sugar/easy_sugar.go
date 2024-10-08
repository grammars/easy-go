package sugar

type number interface {
	int | int8 | int16 | int32 | int64 | float32 | float64
}

// EnsureNotBlank 确保不为龙， 如果空就返回defaultStr
func EnsureNotBlank(str, defaultStr string) string {
	if len(str) == 0 {
		return defaultStr
	}
	return str
}

// EnsurePositive 确保是正数(>0) 如果不是就返回defaultNum
func EnsurePositive[T number](num, defaultNum T) T {
	if num <= 0 {
		return defaultNum
	}
	return num
}

// ReturnIf 三元表达式
func ReturnIf[T any](boolExpression bool, trueReturnValue, falseReturnValue T) T {
	if boolExpression {
		return trueReturnValue
	} else {
		return falseReturnValue
	}
}
