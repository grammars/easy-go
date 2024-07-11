package ego

func Version() string {
	return "0.0.11 (2024-07-11 18:01)"
}

func Return[T any](boolExpression bool, trueReturnValue, falseReturnValue T) T {
	if boolExpression {
		return trueReturnValue
	} else {
		return falseReturnValue
	}
}

func ReturnByFunc[T any](boolExpression bool, trueFuncForReturnValue, falseFuncForReturnValue func() T) T {
	if boolExpression {
		return trueFuncForReturnValue()
	} else {
		return falseFuncForReturnValue()
	}
}
