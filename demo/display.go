// display
package main

import (
	"github.com/donomii/govox"
	"github.com/tbogdala/Voxfile"
)

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
