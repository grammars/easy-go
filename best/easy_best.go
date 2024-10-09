package best

func Echo() {
	println("最佳实践 V 2024-10-09 11:14")
}

func SuccessResult[D any](message string, data D) *EcResult[D] {
	er := EcResult[D]{Message: message, Data: data}
	er.Success()
	return &er
}

func FailResult(message string) *EcResult[any] {
	er := EcResult[any]{Message: message}
	er.Fail()
	return &er
}

type EcResult[D any] struct {
	Code    int    `json:"code"`
	Message string `json:"message"`
	Data    D      `json:"data"`
}

func (er *EcResult[D]) Success() *EcResult[D] {
	er.Code = 0
	return er
}

func (er *EcResult[D]) Fail() *EcResult[D] {
	er.Code = 1
	return er
}

func (er *EcResult[D]) FailError(err error) *EcResult[D] {
	er.Code = 1
	er.Message = err.Error()
	return er
}

func (er *EcResult[D]) SetCode(code int) *EcResult[D] {
	er.Code = code
	return er
}

func (er *EcResult[D]) SetMessage(message string) *EcResult[D] {
	er.Message = message
	return er
}

func (er *EcResult[D]) SetData(data D) *EcResult[D] {
	er.Data = data
	return er
}
