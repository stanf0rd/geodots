package models

import (
	"fmt"
	"math"
)

type Dot struct {
	Name string  `json:"name"` // имя точки
	Lat  float64 `json:"lat"`  // широта
	Lon  float64 `json:"lon"`  // долгота
}

type Area struct {
	Center Dot     `json:"center"`
	Radius float64 `json:"radius"`
}

type Target struct {
	Center   Dot     `json:"dot"`
	Distance float64 `json:"distance"`
}

func AllDots() ([]*Dot, error) {
	rows, err := db.Query("SELECT * FROM dots")
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	dots := make([]*Dot, 0)
	for rows.Next() {
		dot := new(Dot)
		err := rows.Scan(&dot.Name, &dot.Lat, &dot.Lon)
		if err != nil {
			return nil, err
		}
		dots = append(dots, dot)
	}
	if err = rows.Err(); err != nil {
		return nil, err
	}

	return dots, nil
}

func AddDots(dots []Dot) {
	fmt.Println(dots)

	sqlStr := "INSERT INTO dots (name, lat, lon) VALUES "
	for _, row := range dots {
		sqlStr += fmt.Sprintf(
			"('%s', %f, %f),",
			row.Name, row.Lat, row.Lon,
		)
	}
	sqlStr = sqlStr[0 : len(sqlStr)-1]

	_, err := db.Query(sqlStr)
	if err != nil {
		panic(err)
	}
}

func AreaDots(area Area) []Target {
	dots, err := AllDots()
	if err != nil {
		panic(err)
	}
	var areaDots []Target
	for _, dot := range dots {
		distance := Distance(&area.Center, dot)
		if distance < area.Radius {
			target := Target{*dot, distance}
			areaDots = append(areaDots, target)
		}
	}
	return areaDots
}

func hsin(theta float64) float64 {
	return math.Pow(math.Sin(theta/2), 2)
}

// Distance function returns the distance (in km) between two points of
//     a given longitude and latitude relatively accurately (using a spherical
//     approximation of the Earth) through the Haversin Distance Formula for
//     great arc distance on a sphere with accuracy for small distances
//
// Point coordinates are supplied in degrees and converted into rad.
//
// Distance returned is KM
// http://en.wikipedia.org/wiki/Haversine_formula
func Distance(first, second *Dot) float64 {
	var la1, lo1, la2, lo2, r float64
	la1 = first.Lat * math.Pi / 180
	lo1 = first.Lon * math.Pi / 180
	la2 = second.Lat * math.Pi / 180
	lo2 = second.Lon * math.Pi / 180

	r = 6378 // Earth radius in KM

	// calculate
	h := hsin(la2-la1) + math.Cos(la1)*math.Cos(la2)*hsin(lo2-lo1)

	return 2 * r * math.Asin(math.Sqrt(h))
}
