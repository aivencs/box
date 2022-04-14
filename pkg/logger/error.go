package logger

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

func NewBoxError(code Code, level LevelSupport, label string) *BoxError {
	return &BoxError{
		Code:  code,
		Level: level,
		Label: label,
	}
}

type BoxFatal struct {
	Label string
}

func (c *BoxFatal) Error() string {
	return c.Label
}
