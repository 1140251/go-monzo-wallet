package components

type Radius struct {
	TopLeft     int
	TopRight    int
	BottomRight int
	BottomLeft  int
}

func NewRadius(radius int) Radius {
	return Radius{
		TopLeft:     radius,
		TopRight:    radius,
		BottomRight: radius,
		BottomLeft:  radius,
	}
}

func TopRadius(radius int) Radius {
	return Radius{
		TopLeft:  radius,
		TopRight: radius,
	}
}

func BottomRadius(radius int) Radius {
	return Radius{
		BottomRight: radius,
		BottomLeft:  radius,
	}
}

const (
	defaultRadius = 14
)
