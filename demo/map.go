// map
package main

import (
	"math/rand"

	"github.com/donomii/govox"
	"github.com/go-gl/mathgl/mgl32"
)

func mapBlock(size int, f func(govox.Block, int, int, int) govox.Block, blocks voxMap) voxMap {
	for i := 0; i < size; i++ {
		for j := 0; j < size; j++ {
			for k := 0; k < size; k++ {
				blocks[i][j][k] = f(blocks[i][j][k], i, j, k)
			}
		}
	}
	return blocks
}

func GetBlockV(blocks voxMap, position Vec3) *govox.Block {
	return GetBlock(blocks, position[0], position[1], position[2])
}
func GetBlock(blocks voxMap, x, y, z int) *govox.Block {
	size := len(blocks)
	if x >= size {
		return GetBlock(blocks, x-size, y, z)
	}
	if x < 0 {
		return GetBlock(blocks, x+size, y, z)
	}

	if y >= size {
		return GetBlock(blocks, x, y-size, z)
	}
	if y < 0 {
		return GetBlock(blocks, x, y+size, z)
	}

	if z >= size {
		return GetBlock(blocks, x, y, z-size)
	}
	if z < 0 {
		return GetBlock(blocks, x, y, z+size)
	}

	return &blocks[x][y][z]
}
func AddActors(Actrs []Actor, blocks voxMap) {
	for i, a := range Actrs {

		GetBlock(blocks, a.Position[0], a.Position[1], a.Position[2]).Active = true
		GetBlock(blocks, a.Position[0], a.Position[1], a.Position[2]).Color = mgl32.Vec4{1.0, 0.0, 0.0, 0.1}

		Actrs[i].Position[1] = a.Position[0] + rand.Intn(3) - 1
		Actrs[i].Position[2] = a.Position[2] + rand.Intn(3) - 1
	}

}
