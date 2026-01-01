package request

import (
	"encoding/json"
	"io"
	"net/http"

	domainErr "github.com/Rasikrr/bagsy_backend_monolith/internal/domain/errors"
)

type QueryParametersGetter interface {
	GetQueryParameters(r *http.Request) error
}

type ParametersGetter interface {
	GetParameters(r *http.Request) error
}

type Validatable interface {
	Validate() error
}

func GetAndValidateData(r *http.Request, data interface{}) error {
	err := fillFromParameters(r, data)
	if err != nil {
		return err
	}

	if r.Body != nil {
		if err := fillFromBody(r, data); err != nil {
			return err
		}
	}

	if v, ok := data.(Validatable); ok {
		return v.Validate()
	}

	return nil
}

func fillFromParameters(r *http.Request, data interface{}) error {
	var err error
	switch v := data.(type) {
	case QueryParametersGetter:
		err = v.GetQueryParameters(r)
	case ParametersGetter:
		err = v.GetParameters(r)
	}
	if err != nil {
		return domainErr.NewInvalidInputError("failed to parse request parameters", err)
	}
	return nil
}

func fillFromBody(r *http.Request, data interface{}) error {
	defer r.Body.Close()

	bb, err := io.ReadAll(r.Body)
	if err != nil {
		return domainErr.NewInvalidInputError("failed to read request body", err)
	}
	if len(bb) == 0 {
		return nil
	}

	if unmarshaler, ok := data.(json.Unmarshaler); ok {
		if err := unmarshaler.UnmarshalJSON(bb); err != nil {
			return domainErr.NewInvalidInputError("failed to unmarshal request body", err)
		}
		return nil
	}

	if err := json.Unmarshal(bb, data); err != nil {
		return domainErr.NewInvalidInputError("failed to unmarshal request body", err)
	}
	return nil
}
