package error_view

import (
	"encoding/json"
	"io"
	"net/http"

	"github.com/pkg/errors"
	"github.com/toolsparty/mvc"
	"github.com/valyala/fasthttp"
)

type ErrorView struct {
	*mvc.BaseView
}

func (ErrorView) Name() (string, error) {
	return "error", nil
}

func (view *ErrorView) Render(w io.Writer, tpl string, params mvc.ViewParams) error {
	e := &errResponse{
		Code:    http.StatusInternalServerError,
		Message: "Internal Server Error",
		Details: "Error Not Found",
	}

	if wr, ok := w.(*fasthttp.RequestCtx); ok {
		defer func() {
			wr.SetContentType("application/json")
			wr.SetStatusCode(e.Code)
		}()
	}

	if wr, ok := w.(http.ResponseWriter); ok {
		defer func() {
			wr.WriteHeader(e.Code)
		}()
	}

	err, ok := params["error"].(error)
	if !ok {
		if err := view.toJSON(w, &response{Error: e}); err != nil {
			return errors.Wrap(err, "encoding to json failed")
		}

		return errors.New("error not found")
	}

	he, ok := err.(HTTPError)
	if !ok {
		e.Details = err.Error()
		return view.toJSON(w, &response{
			Error: e,
		})
	}

	e.Code = he.Code()
	e.Message = he.Message()
	e.Details = he.Error()

	return view.toJSON(w, &response{
		Error: e,
	})
}

func (ErrorView) toJSON(w io.Writer, b interface{}) error {
	enc := json.NewEncoder(w)

	err := enc.Encode(b)
	if err != nil {
		return errors.Wrap(err, "encoding json failed")
	}

	return nil
}
