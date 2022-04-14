package logger

import "unicode/utf8"

type Err struct {
	Err   error
	Label string
	Code  Code
	Level LevelSupport
}

type ErrOption struct {
	Err   error
	Label string
	Code  Code
}

func (c *Err) Error() string {
	return c.Label
}

func (c *Err) GetErr() error {
	return c.Err
}

func (c *Err) GetLabel() string {
	return c.Label
}

func (c *Err) GetLevel() LevelSupport {
	return c.Level
}

func (c *Err) GetCode() Code {
	return c.Code
}

func NewErr(option ErrOption) *Err {
	err := GetDefaultErr()
	if option.Code != err.Code {
		err.Code = option.Code
	}
	if utf8.RuneCountInString(option.Label) > 0 {
		err.Label = option.Label
	}
	if option.Err != nil {
		err.Err = option.Err
	}
	err.Level = GetErr(err.Code, err.Label).GetLevel()
	return err
}

func GetErr(code Code, label string) *Err {
	value := err[code]
	if utf8.RuneCountInString(label) > 1 {
		value.Label = label
	}
	return &value
}

func GetDefaultErr() *Err {
	value := err[DEFAULT_CODE]
	return &value
}
