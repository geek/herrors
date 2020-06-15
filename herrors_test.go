package herrors

import (
	"errors"
	"fmt"
	"net/http"
	"net/http/httptest"
	"testing"
)

func TestNew(t *testing.T) {
	tests := []struct {
		code uint
		want error
	}{
		{400, fmt.Errorf(http.StatusText(400))},
		{404, fmt.Errorf(http.StatusText(404))},
		{401, ErrUnauthorized},
		{0, ErrInternalServer},
	}

	for _, tt := range tests {
		got := New(tt.code)
		if got.Error() != tt.want.Error() {
			t.Errorf("New.Error(): got: %q, want %q", got, tt.want)
		}
	}
}

func TestErrorIs(t *testing.T) {
	tests := []struct {
		err        error
		wantErr    error
		shouldFail bool
	}{
		{ErrBadRequest, &ErrHttp{code: 400}, false},
		{ErrUnauthorized, &ErrHttp{code: 401}, false},
		{ErrUnauthorized, ErrUnauthorized, false},
		{ErrUnauthorized, ErrBadRequest, true},
		{ErrUnauthorized, errors.New("foo"), true},
	}

	for _, tt := range tests {
		if !errors.Is(tt.err, tt.wantErr) && !tt.shouldFail {
			t.Errorf("Is: got: %v, want %v", tt.err, tt.wantErr)
		}
	}
}

func TestErrorWrapUnwrap(t *testing.T) {
	tests := []struct {
		err  error
		code uint
	}{
		{errors.New("foo"), 400},
		{ErrBadRequest, 401},
	}

	for _, tt := range tests {
		w := Wrap(tt.err, tt.code)
		u := errors.Unwrap(w)

		type hasCode interface {
			Code() uint
		}

		var ce hasCode
		if !errors.As(w, &ce) || ce.Code() != tt.code {
			t.Errorf("Wrap: got code: %d, want code %d", ce.Code(), tt.code)
		}

		if !errors.Is(tt.err, u) {
			t.Errorf("Unwrap: got: %v, want %v", u, tt.err)
		}
	}
}

func TestFmtErrorfWrap(t *testing.T) {
	w := fmt.Errorf("foo %w", ErrBadRequest)

	type hasCode interface {
		Code() uint
	}

	var ce hasCode
	if !errors.As(w, &ce) || ce.Code() != http.StatusBadRequest {
		t.Errorf("Wrap: got code: %d, want code %d", ce.Code(), http.StatusBadRequest)
	}

	if !errors.Is(w, ErrBadRequest) {
		t.Errorf("Unwrap: got: %v, want %v", w, ErrBadRequest)
	}
}

func TestWrite(t *testing.T) {
	tests := []struct {
		err  error
		code int
	}{
		{ErrUnauthorized, 401},
		{ErrBadRequest, 400},
		{nil, 200},
		{errors.New("foo"), 500},
	}

	for _, tt := range tests {
		handlerFunc := func(w http.ResponseWriter, r *http.Request) {
			Write(w, tt.err)
		}

		req, err := http.NewRequest("GET", "/", nil)
		if err != nil {
			t.Fatal(err)
		}

		rr := httptest.NewRecorder()
		handler := http.HandlerFunc(handlerFunc)
		handler.ServeHTTP(rr, req)

		if code := rr.Code; code != tt.code {
			t.Errorf("Write: got %v want %v", code, tt.code)
		}
	}
}
