package errors_test

import (
	"context"
	"log"
	"testing"

	"github.com/appscode/errors"
)

func TestNew(t *testing.T) {
	err := errors.New().WithMessage("hello-world").Err()
	if err == nil {
		t.Error("expected not nil, got nil")
	}
}

func TestErrOverlaps(t *testing.T) {
	err := errors.New("this-is-internal", "this-is-also-internal").Err()
	if err == nil {
		t.Error("expected not nil, got nil")
	}

	errOverlaps := errors.New().WithCause(err).Err()

	parsedErrOverlaps := errors.FromErr(errOverlaps)
	log.Println("got messages", parsedErrOverlaps.Message())
	log.Println("got error", errOverlaps.Error())
}

type fakeContext string

func (fakeContext) String() string {
	return "fake-context-values"
}

func TestErrWithContext(t *testing.T) {
	err := errors.New().WithContext(fakeContext("")).WithMessage("hello-world").Err()
	if err == nil {
		t.Error("expected not nil, got nil")
	}

	parsedErr := errors.FromErr(err)
	if val := parsedErr.Context().String(); val == "" {
		t.Error("expected value fond nil")
	}
	log.Println(parsedErr.Context().String())
}

func TestStackTrace(t *testing.T) {
	err := errors.New().WithContext(fakeContext("")).WithMessage("hello-world").Err()
	if err == nil {
		t.Error("expected not nil, got nil")
	}
	parsedErr := errors.FromErr(err)
	if parsedErr.TraceString() == "" {
		t.Error("expected values got empty")
	}
	log.Println(parsedErr.TraceString())
}

func TestGoContext(t *testing.T) {
	ctx := context.Background()
	err := errors.New().WithGoContext(ctx, func(context.Context) string { return "context" }).Err()
	if err == nil {
		t.Error("expected not nil, got nil")
	}
	parsedErr := errors.FromErr(err)
	if parsedErr.Context().String() != "context" {
		t.Error("expected values got empty")
	}
	log.Println(parsedErr.Context().String())

	ctx = context.Background()
	ctx = context.WithValue(ctx, "foo", "bar")
	err = errors.New().WithGoContext(ctx, func(c context.Context) string { return c.Value("foo").(string) }).Err()
	if err == nil {
		t.Error("expected not nil, got nil")
	}
	parsedErr = errors.FromErr(err)
	if parsedErr.Context().String() != "bar" {
		t.Error("expected values got empty")
	}
	log.Println(parsedErr.Context().String())
}

func TestMessagef(t *testing.T) {
	err := errors.Newf("foo-%s", "bar").Err()
	if err == nil {
		t.Error("expected not nil, got nil")
	}
	parsedErr := errors.FromErr(err)
	if parsedErr.Message() != "foo-bar" {
		t.Error("expected values got empty")
	}
	log.Println(parsedErr.Message())
}
