package context

import (
	"context"
	"net/http"
)

type contextKey string

// ключи для ResponseContext
const (
	ResponseDataKey contextKey = "response.data"
)

// ResponseContext контекст для ответа
type ResponseContext struct {
	context.Context
	err error
	key interface{}
	val interface{}
}

// Err возврашает ошибку если не nil иначе ошибку родителя
func (rc *ResponseContext) Err() error {
	if rc.err != nil {
		return rc.err
	}
	return rc.Context.Err()
}

// Value возврашает значение контекста по ключу
// если в данном контексте нет то ищет у родителя
func (rc *ResponseContext) Value(key interface{}) interface{} {
	if key == rc.key {
		return rc.val
	}
	return rc.Context.Value(key)
}

// WithResponseContext возвращает копию родителя,
// если ошибка не nil то вернет ошибку родителя,
// возврашает поле data по ключу key
func WithResponseContext(key, val interface{}, err error) context.Context {
	return &ResponseContext{
		Context: context.Background(),
		err:     err,
		key:     key,
		val:     val,
	}
}

// SetError устанавливает в контекст ошибку
func SetError(r *http.Request, err error) {
	*r = *r.WithContext(
		WithResponseContext(
			ResponseDataKey,
			nil,
			err,
		),
	)
}

// SetResponse устанавливает в контекст json ответ
func SetResponse(r *http.Request, data interface{}) {
	*r = *r.WithContext(
		WithResponseContext(
			ResponseDataKey,
			data,
			nil,
		),
	)
}

// SetErrorOrResponse устанавливает в контекст ошибку, или дату если ошибка равна nil
func SetErrorOrResponse(r *http.Request, data interface{}, err error) {
	if err != nil {
		SetError(r, err)
	} else {
		SetResponse(r, data)
	}
}
