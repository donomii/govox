// display
package main

import (
	"github.com/donomii/govox"
	"github.com/go-gl/mathgl/mgl32"
)

func DrawPlayer(size int, pos Vec3, blocks voxMap) {

	for i := 0; i < 2; i++ {
		blocks[pos[0]][pos[1]+i][pos[2]] = govox.Block{
			Active: true,
			Color: mgl32.Vec4{
				0.0,
				1.0,
				0.0,
				1.0,
			},
		}
	}
}

func DrawMonster(size int, pos Vec3, blocks voxMap) {

	for i := 0; i < 2; i++ {
		if pos[0] > -1 {
			blocks[pos[0]][pos[1]+i][pos[2]] = govox.Block{
				Active: true,
				Color: mgl32.Vec4{
					1.0,
					0.0,
					0.0,
					1.0,
				},
			}
		}
	}
}
