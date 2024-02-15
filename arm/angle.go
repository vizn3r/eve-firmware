package arm

import "math"

// in radians
type Angle float64

func (a Angle) Radians() float64 {
	return float64(a)
}

func (a Angle) Degrees() float64 {
	return math.Acos(float64(a)) * (180 / math.Pi)
}

func (a Angle) Float64() float64 {
	return float64(a)
}
