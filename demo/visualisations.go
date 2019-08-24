// visualisations
package main

import (
	"github.com/chewxy/math32"

	//"time"

	"github.com/donomii/myvox"
	"github.com/go-gl/mathgl/mgl32"
)

func mapBlock(size int, f func(myvox.Block, int, int, int) myvox.Block, blocks voxMap) voxMap {

	for i := 0; i < size; i++ {

		for j := 0; j < size; j++ {

			for k := 0; k < size; k++ {
				blocks[i][j][k] = f(blocks[i][j][k], i, j, k)
			}
		}
	}
	return blocks
}

func AddFourier(size int, blocks voxMap) voxMap {
	return mapBlock(size,
		func(b myvox.Block, i, j, k int) myvox.Block {

			sizef := float32(size) / 8.0
			return myvox.Block{
				Active: fourier([3]float32{float32(i-size/2) / sizef, float32(j-size/2) / sizef, float32(k-size/2) / sizef}),
				Color: mgl32.Vec4{
					float32(i*i) / float32(size*size),
					float32(j*j) / float32(size*size),
					float32(k*k) / float32(size*size),
					1.0,
				},
			}

		}, blocks)
}

func fourier(q [3]float32) bool {
	r := float32(1.0)
	p := [3]float32{}
	for i := 0; i < 3; i++ {
		p[i] = float32(q[i])
	}
	rotangle := float32(2)

	wr := math32.Sqrt(p[0]*p[0] + p[1]*p[1] + p[2]*p[2])
	wo := math32.Acos(p[1] / wr)
	wi := math32.Atan2(p[0], p[2])

	radius := 1 + r*(0+math32.Sin((float32(1.0)+rotangle*wo))) + (math32.Sin(wi)) + (0 + math32.Cos(wi))
	if wr < radius {
		return true
	} else {
		return false
	}
}

func sphere(q [3]float32) bool {
	radius := float32(0.25)
	p := [3]float32{}
	for i := 0; i < 3; i++ {
		p[i] = float32(q[i])
	}

	wr := math32.Sqrt(p[0]*p[0] + p[1]*p[1] + p[2]*p[2])

	if wr < radius {
		return true
	} else {
		return false
	}
}
