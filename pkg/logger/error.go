package logger

import "unicode/utf8"

// 定义选项
type BoxErrorOption struct {
	Label string
	Name  string
}

type BoxError struct {
	Name  string
	Code  Code
	Level LevelSupport
	Label string
}

func (c *BoxError) Error() string {
	return c.Label
}

func NewBoxError(code Code, label string) *BoxError {
	value := err[code]
	if utf8.RuneCountInString(label) > 1 {
		value.Label = label
	}
	return &value
}

type BoxFatal struct {
	Label string
}

func (c *BoxFatal) Error() string {
	return c.Label
}

func NewBoxFatal(code Code, level LevelSupport, label string) *BoxFatal {
	return &BoxFatal{
		Label: label,
	}
}

func GetDefault() BoxError {
	value := err[DEFAULT_ERROR_CODE]
	return value
}
