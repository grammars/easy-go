package ego

import (
	"github.com/grammars/easy-go/file"
	"github.com/grammars/easy-go/practice"
	"github.com/grammars/easy-go/socket"
)

func Version() string {
	return "0.0.9"
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

var File *file.Utils = &file.Utils{}
var Practice *practice.Utils = &practice.Utils{}
var Socket *socket.Utils = &socket.Utils{}
