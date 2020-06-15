# herrors
--
    import "github.com/geek/herrors"

## Examples

### Basic Bad Request

```go
func handler(w http.ResponseWriter, r *http.Request) {
  // perform some validation
  if !isValid(r) {
    herrors.Write(w, herrors.ErrBadRequest)
    return
  }
}
```

### Bad Request with Wrapping

```go
func validate(r *http.Request) error {
  // ... do some validation and return an error on failure
}

func handler(w http.ResponseWriter, r *http.Request) {
  if err := validate(r); err != nil {
    err = herrors.Wrap(err, http.StatusBadRequest)
    herrors.Write(w, err)
    return
  }
}
```

## Usage

#### func  New

```go
func New(statusCode uint) error
```
New constructs an ErrHttp error with the code and message populated

#### func  Wrap

```go
func Wrap(err error, statusCode uint) error
```
Wrap will create a new ErrHttp and place the provided error inside of it. There
is a complementary errors.Unwrap to retrieve the wrapped error.

#### func  Write

```go
func Write(w http.ResponseWriter, err error)
```
Write sets the response status code to the provided error code. When unable to
find an ErrHttp error in the error chain then a 500 internal error is output.

#### type ErrHttp

```go
type ErrHttp struct {
}
```


#### func (*ErrHttp) Code

```go
func (e *ErrHttp) Code() uint
```
Code returns http status code that the error represents

#### func (*ErrHttp) Error

```go
func (e *ErrHttp) Error() string
```
Error returns http friendly error message

#### func (*ErrHttp) Is

```go
func (e *ErrHttp) Is(target error) bool
```
Is used by errors.Is for comparing ErrHttp

#### func (*ErrHttp) Unwrap

```go
func (e *ErrHttp) Unwrap() error
```
Unwrap will return the first wrapped error if there is one
