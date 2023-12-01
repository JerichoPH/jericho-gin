package tools

import "math"

// Math 数学
type Math struct{}

// NewMath 构造函数
func NewMath() *Math {
	return &Math{}
}

// Decimal 保留小数（四舍五入）
func (Math) Decimal(value float64) float64 {
	return math.Trunc(value*1e2+0.5) * 1e-2
}
