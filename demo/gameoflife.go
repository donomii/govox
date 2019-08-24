// gameoflife.go
package main

import (
	"math"

	"github.com/donomii/myvox"
	"github.com/go-gl/mathgl/mgl32"
)

func countNeighbours(pos, size int, lifeBlocks []bool) int {
	var out int
	for i := -1; i < 2; i = i + 1 {
		for j := -1; j < 2; j = j + 1 {
			for k := -1; k < 2; k = k + 1 {
				if lifeBlocks[pos+i*size*size+j*size+k] && !(i == j && j == k) {
					out = out + 1
				}
			}
		}
	}
	return out
}

func cycle(size int, lifeBlocks []bool, wrapEdges bool) []bool {
	// initialize blocks
	blocks := make([]bool, size*size*size)
	/*
		for i := 0; i < size; i++ {
			for j := 0; j < size; j++ {
				for k := 0; k < size; k++ {
					ii := i
					jj := j
					kk := k

					if wrapEdges {
						if i == 0 {
							ii = size - 2
						}
						if i == size-1 {
							ii = 1
						}
						if j == 0 {
							jj = size - 2
						}
						if j == size-1 {
							jj = 1
						}
						if k == 0 {
							kk = size - 2
						}
						if k == size-1 {
							kk = 1
						}
					}

					blocks[i*size*size+j*size+k] = (countNeighbours(ii*size*size+jj*size+kk, size, lifeBlocks) > 3 && countNeighbours(ii*size*size+jj*size+kk, size, lifeBlocks) < 7)
				}
			}
		}*/

	return blocks
}
func lifeBlocks2Blocks(size int, lifeBlocks []bool, inblocks voxMap) voxMap {
	// initialize blocks
	blocks := inblocks
	if inblocks == nil {
		blocks = make(voxMap, size)
	}
	for i := 0; i < size; i++ {
		if inblocks == nil {
			blocks[i] = make([][]myvox.Block, size)
		}
		for j := 0; j < size; j++ {
			if inblocks == nil {
				blocks[i][j] = make([]myvox.Block, size)
			}
			for k := 0; k < size; k++ {

				blocks[i][j][k] = myvox.Block{
					Active: lifeBlocks[i*size*size+j*size+k],
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
