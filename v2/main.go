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
	"fmt"
	"math"
	"math/rand"
	"time"
)

func makePolygon(color uint8, x1, y1, x2, y2, x3, y3 float64, stroke bool) string {
	colorCode := fmt.Sprintf("#%02x%02x%02x", color, color, color)
	strokeSpec := ""
	if stroke {
		strokeSpec = fmt.Sprintf(`stroke="%s" stroke-width="1" `, colorCode)
	}
	return fmt.Sprintf(`<path d="M %.3f %.3f L %.3f %.3f L %.3f %.3f Z" fill="%s" %s/>`,
		x1, y1, x2, y2, x3, y3, colorCode, strokeSpec,
	)
}

func main() {
	rand.Seed(time.Now().Unix())

	scale := 120
	width := 16 * scale
	height := 9 * scale
	fmt.Printf(`<svg version="1.1" baseProfile="full" xmlns="http://www.w3.org/2000/svg" width="%d" height="%d">`,
		width, height,
	)

	//size of triangles
	stride := 25.0 //TODO: provisional value
	yMax := int(float64(height)/stride) + 1
	xMax := 2 * (int(float64(width)/stride) + 2)
	logo := makeLogo(xMax, yMax)

	//initialize noise
	noise := make([]float64, xMax*yMax)
	for x := 0; x < xMax; x++ {
		for y := 0; y < yMax; y++ {
			noise[x+xMax*y] = rand.Float64()
		}
	}
	//smoothen noise
	for x := 0; x < xMax; x++ {
		for y := 0; y < yMax; y++ {
			getNoiseTorus := func(dx, dy int) float64 {
				idx := (x + dx) + xMax*(y+dy)
				if idx < 0 {
					idx += len(noise)
				}
				if idx >= len(noise) {
					idx -= len(noise)
				}
				return noise[idx]
			}
			sum := getNoiseTorus(-1, -1)
			sum += getNoiseTorus(+0, -1)
			sum += getNoiseTorus(+1, -1)
			sum += getNoiseTorus(-1, +0)
			sum += getNoiseTorus(+0, +0) * 1
			sum += getNoiseTorus(+1, +0)
			sum += getNoiseTorus(-1, +1)
			sum += getNoiseTorus(+0, +1)
			sum += getNoiseTorus(+1, +1)
			noise[x+xMax*y] = sum / 9
		}
	}

	//place triangles
	var (
		foreground string
		background string
	)
	for yIdx := 0; yIdx < yMax; yIdx++ {
		yTop := float64(yIdx) * stride
		yBottom := yTop + stride

		for xIdx := 0; xIdx < xMax; xIdx++ {
			xBase := math.Floor(float64(xIdx)/2) * stride
			switch {
			case yIdx%2 == 0:
				xBase -= stride / 2
			case xIdx%2 == 0:
				xBase -= stride
			}
			if xBase > float64(width) {
				break
			}

			color := uint8(noise[xIdx+xMax*yIdx] * 25)
			if logo.Contains(xIdx, yIdx) {
				color = 255 - uint8(noise[xIdx+xMax*yIdx]*120)
			}

			facingDown := (xIdx+yIdx)%2 == 0
			if facingDown {
				foreground += makePolygon(
					color,
					xBase, yTop,
					xBase+stride, yTop,
					xBase+stride/2, yBottom,
					false,
				)
				background += makePolygon(
					color,
					xBase, yTop,
					xBase+stride, yTop,
					xBase+stride/2, yBottom,
					true,
				)
			} else {
				foreground += makePolygon(
					color,
					xBase+stride, yTop,
					xBase+stride/2, yBottom,
					xBase+stride*3/2, yBottom,
					false,
				)
				background += makePolygon(
					color,
					xBase+stride, yTop,
					xBase+stride/2, yBottom,
					xBase+stride*3/2, yBottom,
					true,
				)
			}
		}
	}

	fmt.Println(background + foreground + `</svg>`)
}

type structure interface {
	Contains(x, y int) bool
}

type union []structure

func (u union) Contains(x, y int) bool {
	for _, s := range u {
		if s.Contains(x, y) {
			return true
		}
	}
	return false
}

type forwardSlab struct {
	//left+bottom must be odd, right+top must be even
	left, right int
	top, bottom int
}

func (s forwardSlab) Contains(x, y int) bool {
	if y < s.top || y > s.bottom {
		return false
	}
	if x < s.left+(s.bottom-y) {
		return false
	}
	if x > s.right-(y-s.top) {
		return false
	}
	return true
}

type backwardSlab struct {
	//left+top must be even, right+bottom must be odd
	left, right int
	top, bottom int
}

func (s backwardSlab) Contains(x, y int) bool {
	if y < s.top || y > s.bottom {
		return false
	}
	if x < s.left+(y-s.top) {
		return false
	}
	if x > s.right-(s.bottom-y) {
		return false
	}
	return true
}

func makeLogo(xMax, yMax int) structure {
	y0 := yMax/2 - 1
	x0 := xMax/2 - 26
	if (x0+y0)%2 == 0 {
		x0++
	}
	return union{
		//first letter C
		forwardSlab{
			left:   x0,
			right:  x0 + 6,
			top:    y0 - 3,
			bottom: y0,
		},
		backwardSlab{
			left:   x0,
			right:  x0 + 6,
			top:    y0 + 1,
			bottom: y0 + 4,
		},
		//second letter C
		forwardSlab{
			left:   x0 + 6,
			right:  x0 + 12,
			top:    y0 - 3,
			bottom: y0,
		},
		backwardSlab{
			left:   x0 + 6,
			right:  x0 + 12,
			top:    y0 + 1,
			bottom: y0 + 4,
		},
		//third letter C
		forwardSlab{
			left:   x0 + 12,
			right:  x0 + 18,
			top:    y0 - 3,
			bottom: y0,
		},
		backwardSlab{
			left:   x0 + 12,
			right:  x0 + 18,
			top:    y0 + 1,
			bottom: y0 + 4,
		},
		//slash
		forwardSlab{
			left:   x0 + 22,
			right:  x0 + 32,
			top:    y0 - 3,
			bottom: y0 + 4,
		},
		//first letter D
		backwardSlab{
			left:   x0 + 36,
			right:  x0 + 42,
			top:    y0 - 3,
			bottom: y0,
		},
		forwardSlab{
			left:   x0 + 36,
			right:  x0 + 42,
			top:    y0 + 1,
			bottom: y0 + 4,
		},
		//second letter D
		backwardSlab{
			left:   x0 + 42,
			right:  x0 + 48,
			top:    y0 - 3,
			bottom: y0,
		},
		forwardSlab{
			left:   x0 + 42,
			right:  x0 + 48,
			top:    y0 + 1,
			bottom: y0 + 4,
		},
	}
}
