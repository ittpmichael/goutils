package rootsolve

import (
	//"fmt"
	"math"
	"os"
)

/*Brent method used for root solve. This function takes [a0, b0] as initial interval, f as a fuction return the value and interate util the answer is converge.*/
func Brent(a0 float64, b0 float64, f func(float64) float64) float64 {
	a, b := a0, b0
	fa, fb := f(a), f(b)
	var s, c, d float64
	delta := 1e-12
	const eps = 0.000000001
	// if f(a0)*f(b0) >=0 then exit function because the root
	// is not bracketed.
	if fa*fb >= 0 {
		os.Exit(1)
	}
	// if |f(a)| < |f(b)| the swap (a,b)
	if math.Abs(fa) < math.Abs(fb) {
		a, b = b, a
	}
	c = a
	fc := f(c)
	mflag := true
	var i int
	for i = 0; fb > eps || math.Abs(b-a) > eps; i++ {
		//fmt.Printf("====Loop %d====\n", i+1)
		if fa != fc && fb != fc {
			// inverse quadratic interpolation
			s = a*fb*fc/(fa-fb)/(fa-fc) + b*fa*fc/(fb-fa)/(fb-fc) + c*fa*fb/(fc-fa)/(fc-fb)
		} else {
			// perform secant method
			s = b - fb*(b-a)/(fb-fa)
		}
		cond1 := s > (3*a+b)/4 && b < s
		cond2 := mflag == true && math.Abs(s-b) >= math.Abs(b-c)/2
		cond3 := mflag == false && math.Abs(s-b) >= math.Abs(c-d)/2
		cond4 := mflag == true && math.Abs(b-c) < math.Abs(delta)
		cond5 := mflag == false && math.Abs(c-d) < math.Abs(delta)
		if cond1 || cond2 || cond3 || cond4 || cond5 {
			// bisection method
			s = (a + b) / 2
			mflag = true
		} else {
			mflag = false
		}
		fs := f(s)
		d, c = c, b
		if fa*fs < 0.0 {
			b = s
		} else {
			a = s
		}
		/*if math.Abs(fa) < math.Abs(fb) {
			a, b = b, a
		}*/
		//fmt.Printf("s = %0.6v\tf(%0.6v) = %0.6v\n", s, s, f(s))
		fa, fb, fc = f(a), f(b), f(c)
	}
	//fmt.Printf("Iteration: %d\tans: %.5v\n", i, s)
	return s
}
