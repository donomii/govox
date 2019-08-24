// animate
package main

import (
	"math"

	"github.com/donomii/myvox"
	"github.com/go-gl/mathgl/mgl32"
)

func rise(size int, blocks voxMap) voxMap {
	// initialize blocks

	for i := 0; i < size; i++ {

		for j := 1; j < size; j++ {

			for k := 0; k < size; k++ {

				blocks[i][j-1][k] = blocks[i][j][k]
				blocks[i][j-1][k].Color[0] = blocks[i][j-1][k].Color.X() * float32(0.95)
			}
		}
	}
	for i := 0; i < size; i++ {

		for j := size - 1; j < size; j++ {

			for k := 0; k < size; k++ {

				blocks[i][j][k] = myvox.Block{
					Active: false,
					Color: mgl32.Vec4{
						float32(math.Mod(float64(i*size*size+j*size+k), 256)) / 256,
						0.0,
						0.0,
						1.0,
					},
				}
			}
		}
	}

	return blocks
}
