package main

import (
	"math"
	"math/rand"
	"runtime"

	"time"

	_ "log"

	"github.com/donomii/govox"
	"github.com/go-gl/glfw/v3.2/glfw"
	"github.com/go-gl/mathgl/mgl32"
)

func init() {
	runtime.LockOSThread()
}

type block struct {
	active bool
	color  mgl32.Vec4
}

var saveCount int = 1

var lifeBlocks []bool

type Actor struct {
	Position [3]int
	Type     int
}

var Actrs []Actor

func main() {
	rand.Seed(time.Now().UnixNano())
	var size float32 = 50.0
	window, rv := govox.InitGraphics(size, 1000, 1000)

	//	blocks := makeBlocks(int(size))

	lifeBlocks = make([]bool, int(size*size*size))
	Actrs = []Actor{
		Actor{[3]int{25, int(size) - 1, 25}, 1},
		Actor{[3]int{35, int(size) - 1, 25}, 1},
	}

	for i, _ := range lifeBlocks {
		lifeBlocks[i] = (rand.Float32() < 0.5)
	}

	go func() {
		for {
			lifeBlocks = cycle(int(size), lifeBlocks, true)
			time.Sleep(1 * time.Second)
		}
	}()

	var roty, rotx float32
	for !window.ShouldClose() {

		// check inputs
		if window.GetKey(glfw.KeyLeft) == glfw.Press {
			roty -= 0.05
		}

		if window.GetKey(glfw.KeyRight) == glfw.Press {
			roty += 0.05
		}

		if window.GetKey(glfw.KeyUp) == glfw.Press {
			rotx -= 0.05
		}

		if window.GetKey(glfw.KeyDown) == glfw.Press {
			rotx += 0.05
		}

		//gl.Viewport(0, 0, 100, 100)
		//blocks := lifeBlocks2Blocks(int(size), lifeBlocks, nil)
		AddActors(Actrs, govox.BlocksBuffer)
		govox.Renderblocks(rv, window, rv.Program, rise(int(size), govox.BlocksBuffer), rotx, roty, int(size))
	}
}

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

func GetBlock(blocks [][][]govox.Block, x, y, z int) *govox.Block {
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
func AddActors(Actrs []Actor, blocks [][][]govox.Block) {
	for i, a := range Actrs {

		GetBlock(blocks, a.Position[0], a.Position[1], a.Position[2]).Active = true
		GetBlock(blocks, a.Position[0], a.Position[1], a.Position[2]).Color = mgl32.Vec4{1.0, 0.0, 0.0, 0.1}

		//for j, _ := range a.Position {
		//Actrs[i].Position[j] = a.Position[j] + rand.Intn(3) - 1
		//}

		Actrs[i].Position[0] = a.Position[0] + rand.Intn(3) - 1
		Actrs[i].Position[2] = a.Position[2] + rand.Intn(3) - 1
		//GetBlock(blocks, a.Position[0], a.Position[1], a.Position[2]).Active = true
		//GetBlock(blocks, Actrs[i].Position[0], Actrs[i].Position[1], Actrs[i].Position[2]).Color = mgl32.Vec4{0.0, 0.0, 0.0, 1.0}
	}

}

func lifeBlocks2Blocks(size int, lifeBlocks []bool, inblocks [][][]govox.Block) [][][]govox.Block {
	// initialize blocks
	blocks := inblocks
	if inblocks == nil {
		blocks = make([][][]govox.Block, size)
	}
	for i := 0; i < size; i++ {
		if inblocks == nil {
			blocks[i] = make([][]govox.Block, size)
		}
		for j := 0; j < size; j++ {
			if inblocks == nil {
				blocks[i][j] = make([]govox.Block, size)
			}
			for k := 0; k < size; k++ {

				blocks[i][j][k] = govox.Block{
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

func rise(size int, blocks [][][]govox.Block) [][][]govox.Block {
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

				blocks[i][j][k] = govox.Block{
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
