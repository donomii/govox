// map
package main

import (
	"math/rand"

	"github.com/donomii/myvox"
	"github.com/go-gl/mathgl/mgl32"
)

func GetBlockV(blocks voxMap, position Vec3) *myvox.Block {
	return GetBlock(blocks, position[0], position[1], position[2])
}
func GetBlock(blocks voxMap, x, y, z int) *myvox.Block {
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

	//log.Println(x, y, z)
	return &blocks[x][y][z]
}
func AddActors(Actrs []Actor, blocks voxMap) {
	for i, a := range Actrs {

		GetBlock(blocks, a.Position[0], a.Position[1], a.Position[2]).Active = true
		GetBlock(blocks, a.Position[0], a.Position[1], a.Position[2]).Color = mgl32.Vec4{1.0, 0.0, 0.0, 0.1}

		//for j, _ := range a.Position {
		//Actrs[i].Position[j] = a.Position[j] + rand.Intn(3) - 1
		//}

		Actrs[i].Position[1] = a.Position[0] + rand.Intn(3) - 1
		Actrs[i].Position[2] = a.Position[2] + rand.Intn(3) - 1
		//GetBlock(blocks, a.Position[0], a.Position[1], a.Position[2]).Active = true
		//GetBlock(blocks, Actrs[i].Position[0], Actrs[i].Position[1], Actrs[i].Position[2]).Color = mgl32.Vec4{0.0, 0.0, 0.0, 1.0}
	}

}
