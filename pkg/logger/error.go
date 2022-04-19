package logger

import (
	"unicode/utf8"
)

type BaseError struct {
	err   error
	label string
	code  Code
	level LevelSupport
}

type ErrOption struct {
	Err   error
	Label string
	Code  Code
}

func NewError(code Code, label string, err error) error {
	ers := baseErrorIns[code]
	if utf8.RuneCountInString(label) == 0 {
		label = ers.label
	}
	return &BaseError{
		label: label,
		code:  code,
		level: ers.level,
		err:   err,
	}
}

func (c *BaseError) Error() string {
	if utf8.RuneCountInString(c.label) > 0 {
		return c.label
	}
	if c.err != nil {
		return c.err.Error()
	}
	return c.label
}

func (c *BaseError) UnWrap() error {
	return c.err
}

func (c *BaseError) Label() string {
	return c.label
}

func (c *BaseError) Code() Code {
	return c.code
}

func (c *BaseError) Level() LevelSupport {
	return c.level
}
