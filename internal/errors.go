package internal

import "errors"

// ErrUndefinedAngle is returned when the angle between two points
// p and q is undefined. e.g. when p == q
var ErrUndefinedAngle = errors.New("Can't compute between given points")
