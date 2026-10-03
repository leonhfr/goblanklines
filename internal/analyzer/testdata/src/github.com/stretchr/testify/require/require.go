// Package require is a minimal stand-in for github.com/stretchr/testify/require,
// exposing only the exported API surface the testify rule's test fixtures
// need to type-check under analysistest (which type-checks but never runs
// this code).
package require

type TestingT interface {
	Errorf(format string, args ...any)
	FailNow()
}

func Equal(t TestingT, expected, actual any, msgAndArgs ...any)    {}
func NotEqual(t TestingT, expected, actual any, msgAndArgs ...any) {}
func True(t TestingT, value bool, msgAndArgs ...any)               {}
func False(t TestingT, value bool, msgAndArgs ...any)              {}
func Nil(t TestingT, object any, msgAndArgs ...any)                {}
func NotNil(t TestingT, object any, msgAndArgs ...any)             {}
func NoError(t TestingT, err error, msgAndArgs ...any)             {}
func Error(t TestingT, err error, msgAndArgs ...any)               {}
func Len(t TestingT, object any, length int, msgAndArgs ...any)    {}
func Contains(t TestingT, s, contains any, msgAndArgs ...any)      {}

type Assertions struct {
	t TestingT
}

func New(t TestingT) *Assertions { return &Assertions{t: t} }

func (a *Assertions) Equal(expected, actual any, msgAndArgs ...any) { Equal(a.t, expected, actual, msgAndArgs...) }

func (a *Assertions) NotEqual(expected, actual any, msgAndArgs ...any) {
	NotEqual(a.t, expected, actual, msgAndArgs...)
}

func (a *Assertions) True(value bool, msgAndArgs ...any)  { True(a.t, value, msgAndArgs...) }
func (a *Assertions) False(value bool, msgAndArgs ...any) { False(a.t, value, msgAndArgs...) }
func (a *Assertions) Nil(object any, msgAndArgs ...any)   { Nil(a.t, object, msgAndArgs...) }
func (a *Assertions) NotNil(object any, msgAndArgs ...any) {
	NotNil(a.t, object, msgAndArgs...)
}

func (a *Assertions) NoError(err error, msgAndArgs ...any) { NoError(a.t, err, msgAndArgs...) }
func (a *Assertions) Error(err error, msgAndArgs ...any)   { Error(a.t, err, msgAndArgs...) }

func (a *Assertions) Len(object any, length int, msgAndArgs ...any) {
	Len(a.t, object, length, msgAndArgs...)
}

func (a *Assertions) Contains(s, contains any, msgAndArgs ...any) {
	Contains(a.t, s, contains, msgAndArgs...)
}
