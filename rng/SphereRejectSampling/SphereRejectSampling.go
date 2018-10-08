package rng

import (
	"rng/MRG63k3a"
)

func SphereRejectSampling() (float64, float64, float64) {
	/* SphereRejectSampling return uniformly distribution in Cartesian
	coor. x, y, z respectively.
	*/
	var x, y, z float64
	for {
		x = 1.0 - 2.0*rng.MRG63k3a()
		y = 1.0 - 2.0*rng.MRG63k3a()
		z = 1.0 - 2.0*rng.MRG63k3a()
		r2 := x*x + y*y + z*z

		if r2 < 1.00 {
			break
		}
	}
	return x, y, z
}
