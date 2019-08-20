// game.go
package main

import (
	"math/rand"

	//"time"

	"github.com/donomii/govox"
	"github.com/go-gl/mathgl/mgl32"
)

var monsters []Vec3

func InitGame(size int) {
	for i := 0; i < 5; i++ {
		monsters = append(monsters, Vec3{rand.Intn(size), 0, rand.Intn(size)})
	}
}

func handleCollision(pl, other Vec3) []Vec3 {
	out := []Vec3{}
	for _, v := range monsters {
		if v[0] == other[0] && v[1] == other[1] && v[2] == other[2] {

		} else {
			out = append(out, v)
		}
	}
	return out
}

func moveOk(pos Vec3, blocks voxMap) bool {
	return !blocks[pos[0]][pos[1]][pos[2]].Active
}

func GenerateMaze(size int) [][]int {
	grid := make([][]int, size)
	for i := 0; i < size; i++ {
		var row []int
		for j := 0; j < size; j++ {
			row = append(row, rand.Intn(5))
		}
		grid[i] = row
	}
	return grid
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
