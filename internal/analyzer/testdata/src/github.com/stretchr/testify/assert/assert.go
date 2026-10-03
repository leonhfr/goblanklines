// Package assert is a minimal stand-in for github.com/stretchr/testify/assert,
// exposing only the exported API surface the testify rule's test fixtures
// need to type-check under analysistest (which type-checks but never runs
// this code).
package assert

type TestingT interface {
	Errorf(format string, args ...any)
}

func Equal(t TestingT, expected, actual any, msgAndArgs ...any) bool    { return true }
func NotEqual(t TestingT, expected, actual any, msgAndArgs ...any) bool { return true }
func True(t TestingT, value bool, msgAndArgs ...any) bool               { return true }
func False(t TestingT, value bool, msgAndArgs ...any) bool              { return true }
func Nil(t TestingT, object any, msgAndArgs ...any) bool                { return true }
func NotNil(t TestingT, object any, msgAndArgs ...any) bool             { return true }
func NoError(t TestingT, err error, msgAndArgs ...any) bool             { return true }
func Error(t TestingT, err error, msgAndArgs ...any) bool               { return true }
func Len(t TestingT, object any, length int, msgAndArgs ...any) bool    { return true }
func Contains(t TestingT, s, contains any, msgAndArgs ...any) bool      { return true }

type Assertions struct {
	t TestingT
}

func New(t TestingT) *Assertions { return &Assertions{t: t} }

func (a *Assertions) Equal(expected, actual any, msgAndArgs ...any) bool {
	return Equal(a.t, expected, actual, msgAndArgs...)
}

func (a *Assertions) NotEqual(expected, actual any, msgAndArgs ...any) bool {
	return NotEqual(a.t, expected, actual, msgAndArgs...)
}

func (a *Assertions) True(value bool, msgAndArgs ...any) bool { return True(a.t, value, msgAndArgs...) }

func (a *Assertions) False(value bool, msgAndArgs ...any) bool {
	return False(a.t, value, msgAndArgs...)
}

func (a *Assertions) Nil(object any, msgAndArgs ...any) bool { return Nil(a.t, object, msgAndArgs...) }

func (a *Assertions) NotNil(object any, msgAndArgs ...any) bool {
	return NotNil(a.t, object, msgAndArgs...)
}

func (a *Assertions) NoError(err error, msgAndArgs ...any) bool {
	return NoError(a.t, err, msgAndArgs...)
}

func (a *Assertions) Error(err error, msgAndArgs ...any) bool { return Error(a.t, err, msgAndArgs...) }

func (a *Assertions) Len(object any, length int, msgAndArgs ...any) bool {
	return Len(a.t, object, length, msgAndArgs...)
}

func (a *Assertions) Contains(s, contains any, msgAndArgs ...any) bool {
	return Contains(a.t, s, contains, msgAndArgs...)
}
