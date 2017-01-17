package model

type DyDanmuColor int

func (c DyDanmuColor) IsNormal() bool {
	return c == 0
}

func (c DyDanmuColor) IsRed() bool {
	return c == 1
}

func (c DyDanmuColor) IsBlue() bool {
	return c == 2
}

func (c DyDanmuColor) IsGreen() bool {
	return c == 3
}
func (c DyDanmuColor) IsOrange() bool {
	return c == 4
}
func (c DyDanmuColor) IsPurple() bool {
	return c == 5
}
func (c DyDanmuColor) IsPink() bool {
	return c == 6
}
