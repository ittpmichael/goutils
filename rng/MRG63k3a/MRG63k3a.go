package rng
/*
   63-bits Random number generator U(0,1): MRG63k3a
   Author: Pierre L'Ecuyer,
   Source: Good Parameter Sets for Combined Multiple Recursive Random
           Number Generators,
           Shorter version in Operations Research,
           47, 1 (1999), 159--164.
   ---------------------------------------------------------
*/

const norm float64 = 1.0842021724855052e-19
const m1 = 9223372036854769163
const m2 = 9223372036854754679
const a12 =  1754669720
const q12 = 5256471877
const r12 = 251304723
const a13n = 3182104042
const q13 = 2898513661
const r13 = 394451401
const a21 = 31387477935
const q21 = 293855150
const r21 = 143639429
const a23n = 6199136374
const q23 = 1487847900
const r23 = 985240079


/*** State variables s10, s11, s12, s20, s21, s22 must be 64-bits integers.

The seeds for s10, s11, s12 must be integers in [0, m1 - 1] and not all 0. 
The seeds for s20, s21, s22 must be integers in [0, m2 - 1] and not all 0. 
***/

const SEED = 123456789

var s10, s11, s12, s20, s21, s22 int64 = SEED, SEED, SEED, SEED, SEED, SEED
			
func MRG63k3a() float64 {
	var h, p12, p13, p21, p23 int64
	/* Component 1 */
	h = s10 / q13
	p13 = a13n * (s10 - h * q13) - h * r13
	h = s11 / q12
	p12 = a12 * (s11 - h * q12) - h * r12
	if p13 < 0 {
	   p13 += m1
	}
	if p12 < 0 {
	   p12 += m1 - p13
	} else {
	   p12 -= p13
	}
	if p12 < 0 {
	   p12 += m1
	}
	s10 = s11
	s11 = s12
	s12 = p12
 
	/* Component 2 */
	h = s20 / q23
	p23 = a23n * (s20 - h * q23) - h * r23
	h = s22 / q21
	p21 = a21 * (s22 - h * q21) - h * r21
	if p23 < 0 {
	   p23 += m2 
	}
	if p21 < 0 {
	   p21 += m2 - p23
	} else {
	   p21 -= p23
	}
	if p21 < 0 {
	   p21 += m2
	}
	s20 = s21
	s21 = s22
	s22 = p21
 
	/* Combination */
	if p12 > p21 {
	   return (float64(p12 - p21) * norm)
	} else {
	   return (float64(p12 - p21 + m1) * norm)
	}
}
