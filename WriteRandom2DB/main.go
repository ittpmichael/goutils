package main

import (
	"database/sql"
	"log"
	"math"
	"os"
	"rng/SphereRejectSampling"
	"rootsolve"
	"strconv"

	_ "github.com/mattn/go-sqlite3"
	"gonum.org/v1/gonum/integrate/quad"
)

const (
	//sample int = 100
	sample int = 1215856128
	// R = 700 Mpc
	R float64 = 700
	// OmegaM (Matter) = 0.315
	OmegaM float64 = 0.315
	// OmegaL (Lambda) = 0.685
	OmegaL float64 = 0.685
	// H0 Hubble's parameter
	H0 float64 = 67.3
	// c light speed in km/s
	c float64 = 299792.458
	// w = -1
	w float64 = -1
)

func main() {
	//define file name and remove existed file
	filename := os.Args[1]
	os.Remove(filename)
	// Open File for write, and defer close file
	db, err := sql.Open("sqlite3", filename)
	check(err)
	defer db.Close()

	// prepare sql command to execute
	sqlStmt := `CREATE TABLE RDLightcone ( RA float, Dec float, RARad float, DecRad float, R float, Rh float, Redshift float);`
	_, err = db.Exec(sqlStmt)
	check(err)

	tx, err := db.Begin()
	check(err)

	var StrRadii, StrDec, StrRA, StrRARad, StrDecRad, StrRadiih, StrRedshift string

	stmt, err := tx.Prepare("INSERT INTO RDLightcone(RA, Dec, RARad, DecRad, R, Rh, Redshift) VALUES (?,?,?,?,?,?,?)")
	check(err)
	defer stmt.Close()
	for i := 1; i <= sample; i++ {
		radx, rady, radz := rng.SphereRejectSampling()
		X := radx * R
		Y := rady * R
		Z := radz * R

		//Comoving Radial Distance
		Radii := math.Sqrt(X*X + Y*Y + Z*Z)
		Radiih := Radii * H0 / 100.0

		//Phi span across RA [0,2pi]
		RARad := math.Atan2(Y, X)
		RA := RadToDeg(RARad) + 180.0

		//compute Theta(Dec) with domain [0,pi],
		// then change the domain to [-pi/2,pi/2]
		DecRad := math.Acos(Z/Radii) - math.Pi
		Dec := RadToDeg(DecRad)

		Approx := Radii / 4375.0 /* numerical approximation z = r/70*0.016 */
		a0 := Approx * 0.9
		b0 := Approx * 1.1
		//fmt.Printf("i: %d\tRadii: %v\ta0: %0.5v\tb0: %0.5v\n", i, Radii, a0, b0)
		Redshift := rootsolve.Brent(a0, b0, func(x float64) float64 {
			ev := quad.Fixed(df, 0, x, 10, nil, 0) - Radii
			return ev
		})

		//convert number to string
		StrRadii = strconv.FormatFloat(Radii, 'f', 6, 64)
		StrRadiih = strconv.FormatFloat(Radiih, 'f', 6, 64)
		StrRA = strconv.FormatFloat(RA, 'f', 6, 64)
		StrRARad = strconv.FormatFloat(RARad, 'f', 6, 64)
		StrDec = strconv.FormatFloat(Dec, 'f', 6, 64)
		StrDecRad = strconv.FormatFloat(DecRad, 'f', 6, 64)
		StrRedshift = strconv.FormatFloat(Redshift, 'f', 9, 64)

		_, err = stmt.Exec(
			StrRA, StrDec, StrDecRad,
			StrRARad, StrRadii, StrRadiih,
			StrRedshift,
		)
		check(err)
	}
	tx.Commit()
}

func check(e error) {
	if e != nil {
		log.Fatal(e)
	}
}

// RadToDeg converts angle unit from radian to degree
func RadToDeg(rad float64) float64 {
	return 180.0 * rad / math.Pi
}

func df(z float64) float64 {
	Omega := OmegaM + OmegaL
	return c / H0 / math.Sqrt((1-Omega)*(1+z)*(1+z)+OmegaL+OmegaM*math.Pow((1+z), 3))
}
