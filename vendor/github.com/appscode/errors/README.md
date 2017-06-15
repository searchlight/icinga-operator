[![Go Report Card](https://goreportcard.com/badge/github.com/appscode/errors)](https://goreportcard.com/report/github.com/appscode/errors)

[Website](https://appscode.com) • [Slack](https://slack.appscode.com) • [Forum](https://discuss.appscode.com) • [Twitter](https://twitter.com/AppsCodeHQ)

# errors
Package errors provides detailed error types, and context support on top of go error.

Simple Usage
============

```go
type Context string

func (c Context) String() string {
	return "context: "+ string(c)
}
errors.New().WithMessage("error message").
	WithCause(err).
	WithContext(Context("custom context")).
	Failed()
```
######To Build Error
```go
err := errors.New().Internal() //Create external error
err := errors.New("Custom error message").
	WithCause(err).
	WithContext(Context("Custom context")).
	BadRequest() //Create BadRequest error with custom message, existing error and custom context
err := errors.New().
	WithMessage("Custom error message").
	External() // Create External error with custom message
```
######To Parse Error
```go
aError := errors.Parse(extError)
```
######To Extract Error Element
```go
c := aError.Context()
m :=  aError.Message()
code := aError.Code()
trace := aError.TraceString()
```
######To Test Error Type
```go
if errors.IsBadRequest(badReqError) {
	//Code
}

if errors.Is(badReqError, errors.BadRequest) {
    //Code
}
```

Available Error Types
=====================

```go
 External
 Failed
 Internal
 InvalidData
 InvalidPaymentInformation
 InvalidQuota
 NotFound
 PaymentInformationUnavailable
 PermissionDenied
 QuotaLimitExceed
 Unauthorized
 Unimplemented
 Unknown
```

Full documentation is available on [godoc](https://godoc.org/github.com/appscode/errors)

#License
See [LICENSE ](https://github.com/appscode/errors/blob/master/LICENSE)
