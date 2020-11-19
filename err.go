package go_utils


type ParamError struct {
	TipErr error // 提示错误
	Code   int   // 错误码
	Err    error // 内部错误
}

func (m *ParamError) Error() string {
	return m.TipErr.Error()
}

func NewParamError(code int, tipError error, err error) *ParamError {
	return &ParamError{
		TipErr: tipError,
		Code:   code,
		Err:    err,
	}
}

type InternalError struct {
	TipErr error // 提示错误
	Code   int   // 错误码
	Err    error // 内部错误
}

func (m *InternalError) Error() string {
	return m.TipErr.Error()
}

func NewInternalError(code int, tipError error, err error) *InternalError {
	return &InternalError{
		TipErr: tipError,
		Code:   code,
		Err:    err,
	}
}
