// game.go
package main

import (
	"math/rand"

//"time"
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

func moveOk(pos Vec3, maze [][]int) bool {
	return !(maze[pos[0]][pos[2]] == 1)
	//return !blocks[pos[0]][pos[1]][pos[2]].Active
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
