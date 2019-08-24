// display
package main

import (
	"math/rand"

	"github.com/donomii/govox"
	"github.com/go-gl/mathgl/mgl32"
	"github.com/tbogdala/Voxfile"
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

func AddMaze(size int, pos Vec3, wall *voxfile.VoxFile, maze [][]int, blocks voxMap) {
	imin := pos[0] - tileRadius
	kmin := pos[2] - tileRadius
	for i := 0; i < tiles; i++ {
		for k := 0; k < tiles; k++ {
			if maze[i+imin][k+kmin] == 1 {
				magica2govox(size, Vec3{20 * i, 0, 20 * k}, wall, blocks)
			}
		}
	}
}

func AddMonster(size int, pos, player Vec3, eye *voxfile.VoxFile, blocks voxMap) {
	x := pos[0] - player[0] + tileRadius
	y := pos[2] - player[2] + tileRadius
	if x*x < 16 && y*y < 16 {
		magica2govox(size, Vec3{size / tiles * x, 0, size / tiles * y}, eye, blocks)
	}
}

func AddFloor(size int, maze [][]int, blocks voxMap) {
	rs := rand.NewSource(9384598375)
	r := rand.New(rs)

	for i := 0; i < size; i++ {

		for k := 0; k < size; k++ {
			tweak := r.Float32() / 2.0
			for j := 0; j < 2; j++ {
				if maze[i][k] == 1 {

					blocks[i][j][k] = govox.Block{
						Active: true,
						Color: mgl32.Vec4{
							0.5 + tweak,
							0.5 + tweak,
							0.5 + tweak,
							1.0,
						},
					}
				} else {
					blocks[i][j][k] = govox.Block{
						Active: false,
						Color: mgl32.Vec4{
							0.5,
							0.5,
							0.5,
							1.0,
						},
					}

				}
			}
		}
	}
	//rand.Seed(time.Now().Unix())
}
