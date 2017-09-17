/*******************************************************************************
*
* Copyright 2017 Stefan Majewsky <majewsky@gmx.net>
*
* This program is free software: you can redistribute it and/or modify it under
* the terms of the GNU General Public License as published by the Free Software
* Foundation, either version 3 of the License, or (at your option) any later
* version.
*
* This program is distributed in the hope that it will be useful, but WITHOUT ANY
* WARRANTY; without even the implied warranty of MERCHANTABILITY or FITNESS FOR
* A PARTICULAR PURPOSE. See the GNU General Public License for more details.
*
* You should have received a copy of the GNU General Public License along with
* this program. If not, see <http://www.gnu.org/licenses/>.
*
*******************************************************************************/

package main

import (
	"bytes"
	"fmt"
	"image"
	"image/png"
	"math/rand"
	"os"
	"time"
)

func randomUint(min, max int) int {
	return min + int(rand.Float64()*float64(max-min))
}

func squareNorm(p image.Point) int {
	return p.X*p.X + p.Y*p.Y
}

type Area struct {
	Bounds image.Rectangle
	Points []image.Point
}

func (a Area) randomPoint(padding int) image.Point {
	return image.Point{
		X: randomUint(a.Bounds.Min.X+padding, a.Bounds.Max.X-padding),
		Y: randomUint(a.Bounds.Min.Y+padding, a.Bounds.Max.Y-padding),
	}
}

func (a *Area) PlacePoint(img image.Image, tries uint) {
	if len(a.Points) == 0 {
		a.Points = append(a.Points, a.randomPoint(5))
		return
	}

	var bestCandidate image.Point
	biggestMinSqDist := 0

	for try := uint(0); try < tries; try++ {
		candidate := a.randomPoint(5)
		scale := getScalingFactor(img, candidate)

		minSqDist := squareNorm(a.Bounds.Max)
		for _, p := range a.Points {
			sqDist := int(float64(squareNorm(p.Sub(candidate))) * scale)
			if sqDist < minSqDist {
				minSqDist = sqDist
			}
			if minSqDist < biggestMinSqDist {
				break
			}
		}

		if minSqDist > biggestMinSqDist {
			biggestMinSqDist = minSqDist
			bestCandidate = candidate
		}
	}

	a.Points = append(a.Points, bestCandidate)
}

func getScalingFactor(img image.Image, point image.Point) float64 {
	ir, ig, ib, ia := img.At(point.X, point.Y).RGBA()
	r := float64(ir)
	g := float64(ig)
	b := float64(ib)
	a := float64(ia)
	return 0.1 + 1.8*(r/a+g/a+b/a)
}

func main() {
	rand.Seed(time.Now().Unix())

	if len(os.Args) != 2 {
		fmt.Fprintln(os.Stderr, "usage: "+os.Args[0]+" <png-file>")
		os.Exit(1)
	}

	file, err := os.Open(os.Args[1])
	fatalIfError(err)
	img, err := png.Decode(file)
	fatalIfError(err)
	fatalIfError(file.Close())

	a := Area{
		Bounds: img.Bounds(),
	}
	for len(a.Points) < 10000 {
		a.PlacePoint(img, uint(len(a.Points))/3)
		// a.PlacePoint(img, 100)
		if len(a.Points)%100 == 0 {
			fmt.Fprintf(os.Stderr, "Placed %d points\n", len(a.Points))
		}
	}

	width := a.Bounds.Max.X - a.Bounds.Min.X
	height := a.Bounds.Max.Y - a.Bounds.Min.Y

	pixels := make([]bool, width*height)
	for _, p := range a.Points {
		x := p.X - a.Bounds.Min.X
		y := p.Y - a.Bounds.Min.Y
		pixels[x+width*y] = true
	}

	fmt.Printf("P2 %d %d 1\n", width, height)
	var buf bytes.Buffer
	buf.Grow(len(pixels) * 2)
	for _, val := range pixels {
		if val {
			buf.WriteString("1 ")
		} else {
			buf.WriteString("0 ")
		}
	}
	buf.WriteTo(os.Stdout)
}

func fatalIfError(err error) {
	if err == nil {
		return
	}
	fmt.Fprintln(os.Stderr, err.Error())
	os.Exit(1)
}
