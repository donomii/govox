// display
package main

import (
	"log"
	"math/rand"

	"github.com/donomii/govox"
	"github.com/go-gl/mathgl/mgl32"
	voxfile "github.com/tbogdala/Voxfile"
)

type MarkovRule struct {
	From []int
	To   []int
}

func DoRule(tiles int, direction int, pos []int, rule MarkovRule, maze, mazeB [][][]int) bool {
	i := pos[0]
	j := pos[1]
	k := pos[2]

	switch direction {
	case 0:
		if k < tiles-3 && maze[i][j][k] == rule.From[0] && maze[i][j][k+1] == rule.From[1] && maze[i][j][k+2] == rule.From[2] {
			mazeB[i][j][k] = rule.To[0]
			mazeB[i][j][k+1] = rule.To[1]
			mazeB[i][j][k+2] = rule.To[2]
			return true
		}
		return false
	case 1:
		if j < tiles-3 && maze[i][j][k] == rule.From[0] && maze[i][j+1][k] == rule.From[1] && maze[i][j+2][k] == rule.From[2] {
			mazeB[i][j][k] = rule.To[0]
			mazeB[i][j+1][k] = rule.To[1]
			mazeB[i][j+2][k] = rule.To[2]
			return true
		}
	case 2:
		if i < tiles-3 && maze[i][j][k] == rule.From[0] && maze[i+1][j][k] == rule.From[1] && maze[i+2][j][k] == rule.From[2] {
			mazeB[i][j][k] = rule.To[0]
			mazeB[i+1][j][k] = rule.To[1]
			mazeB[i+2][j][k] = rule.To[2]
			return true
		}
	case 3:
		if k > 3 && maze[i][j][k] == rule.From[0] && maze[i][j][k-1] == rule.From[1] && maze[i][j][k-2] == rule.From[2] {
			mazeB[i][j][k] = rule.To[0]
			mazeB[i][j][k-1] = rule.To[1]
			mazeB[i][j][k-2] = rule.To[2]
			return true
		}
	case 4:
		if j > 3 && maze[i][j][k] == rule.From[0] && maze[i][j-1][k] == rule.From[1] && maze[i][j-2][k] == rule.From[2] {
			mazeB[i][j][k] = rule.To[0]
			mazeB[i][j-1][k] = rule.To[1]
			mazeB[i][j-2][k] = rule.To[2]
			return true
		}
	case 5:
		if i > 3 && maze[i][j][k] == rule.From[0] && maze[i-1][j][k] == rule.From[1] && maze[i-2][j][k] == rule.From[2] {
			mazeB[i][j][k] = rule.To[0]
			mazeB[i-1][j][k] = rule.To[1]
			mazeB[i-2][j][k] = rule.To[2]
			return true
		}

	}
	return false
}
func ApplyRule(tiles int, maze, mazeB [][][]int, rule MarkovRule, step bool, maxChange int) int {
	changed := 0
	for i := 0; i < tiles; i++ {
		for j := 0; j < tiles; j++ {
			for k := 0; k < tiles; k++ {

				order := []int{0, 1, 2, 3, 4, 5}
				//Permute array to random order
				for ii := range order {
					jj := rand.Intn(ii + 1)
					order[ii], order[jj] = order[jj], order[ii]
				}
				//log.Printf("Order: %v", order)
				//order = order[:1]
			done:
				for _, dir := range order {
					if DoRule(tiles, dir, []int{i, j, k}, rule, maze, mazeB) {

						changed++
						break done
						if step {
							if changed > maxChange {
								return changed
							}
						}
						continue
					}

				}
			}
		}

	}
	return changed
}

func AddMarkov(edgeLength int, pos Vec3, wall *voxfile.VoxFile, maze [][][]int, blocks voxMap) {
	log.Printf("Adding markov, edgeLength: %d", edgeLength)
	for i := 0; i < edgeLength; i++ {
		for j := 0; j < edgeLength; j++ {
			for k := 0; k < edgeLength; k++ {
				x := uint8(i)
				y := uint8(j)
				z := uint8(k)
				switch maze[i][j][k] {
				case 1:
					blocks[x][y][z].Active = true
					blocks[x][y][z].Color = mgl32.Vec4{1.0, 0.0, 0.0, 1.0}
				case 2:
					blocks[x][y][z].Active = true
					blocks[x][y][z].Color = mgl32.Vec4{0.0, 1.0, 0.0, 1.0}
				case 3:
					blocks[x][y][z].Active = true
					blocks[x][y][z].Color = mgl32.Vec4{0.0, 0.0, 1.0, 1.0}
				}
			}
		}
	}
}

func AddMaze(size int, pos Vec3, wall *voxfile.VoxFile, maze [][]int, blocks voxMap) {
	imin := pos[0] - tileRadius
	kmin := pos[2] - tileRadius
	tileWidth := size / tiles
	for i := 0; i < tiles; i++ {
		for k := 0; k < tiles; k++ {
			if maze[i+imin][k+kmin] == 1 {
				magica2govox(size, Vec3{tileWidth * i, 0, tileWidth * k}, wall, blocks)
			}
		}
	}
}

func AddMonster(size int, pos, player Vec3, eye *voxfile.VoxFile, blocks voxMap) {
	x := pos[0] - player[0] + tileRadius
	y := pos[2] - player[2] + tileRadius
	if InView(player, pos) {
		magica2govox(size, Vec3{size / tiles * x, 0, size / tiles * y}, eye, blocks)
	}
}

func ClearDisplay(size int, blocks voxMap) {
	mapBlock(size, func(b govox.Block, i, j, k int) govox.Block {
		b.Active = false
		return b
	}, blocks)
}
