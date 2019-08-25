package main

import (
	//	"math"

	"fmt"
	"log"

	//"math"
	"math/rand"
	"runtime"

	_ "log"
	"time"

	"github.com/go-gl/gl/v4.1-core/gl"
	"github.com/tbogdala/Voxfile"

	"github.com/donomii/glim"
	"github.com/donomii/govox"
	"github.com/go-gl/glfw/v3.2/glfw"
	"github.com/go-gl/mathgl/mgl32"
)

type RenderVars struct {
	//Col        mgl32.Vec4
	//ColUni     int32
	Vao        uint32
	Vbo        uint32
	VertAttrib uint32
	Program    uint32
}

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
	Position Vec3
	Type     int
}

type Vec3 [3]int
type voxMap [][][]govox.Block

var Actrs []Actor
var PlayerPos Vec3
var palette []mgl32.Vec4

func magica2govox(sizei int, pos Vec3, vox *voxfile.VoxFile, blocks voxMap) {
	defer func() {
		if r := recover(); r != nil {
			fmt.Printf("Recovered in magica2govox, drawing at position %v: %v\n", pos, r)
		}
	}()
	x := uint8(pos[0])
	y := uint8(pos[1])
	z := uint8(pos[2])
	for _, v := range vox.Voxels {
		//log.Printf("%+v,%v,%v,%v\n", v, x, y, z)
		blocks[v.X+x][v.Z+y][v.Y+z].Active = true
		blocks[v.X+x][v.Z+y][v.Y+z].Color = palette[v.Index]

	}
	/*
		size := uint32(sizei)
		for i := uint32(0); i < vox.SizeX && i < size; i++ {
			for j := uint32(0); j < vox.SizeY && j < size; j++ {
				for k := uint32(0); k < vox.SizeZ && k < size; k++ {
					log.Println(i, j, k)
					blocks[i][j][k].Active = vox.Voxels[k*vox.SizeX*vox.SizeY+j*vox.SizeX+i].Index > 0
				}
			}
		}
	*/
}

var roty, rotx float32

func handleKeys(window *glfw.Window, maze [][]int) {
	lastInputTime := time.Now()
	lastInputTime2 := time.Now()

	for {
		if glfw.Press == 1 {
			if time.Now().Sub(lastInputTime).Nanoseconds() > 150000000 {
				lastInputTime = time.Now()
				wantPos := PlayerPos
				if window.GetKey(glfw.KeyW) == glfw.Press {
					wantPos[2] = wantPos[2] - 1
				}
				if window.GetKey(glfw.KeyS) == glfw.Press {
					wantPos[2] = wantPos[2] + 1
				}
				if window.GetKey(glfw.KeyD) == glfw.Press {
					wantPos[0] = wantPos[0] + 1
				}
				if window.GetKey(glfw.KeyA) == glfw.Press {
					wantPos[0] = wantPos[0] - 1
				}
				if moveOk(wantPos, maze) {
					PlayerPos = wantPos
				} else {
					log.Println("Move not ok")
					monsters = handleCollision(PlayerPos, wantPos)
				}
				log.Printf("Player: %v\n", PlayerPos)

			}
			if time.Now().Sub(lastInputTime2).Nanoseconds() > 15000000 {
				lastInputTime2 = time.Now()

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

				if window.GetKey(glfw.KeyDown) == glfw.Press {
					rotx += 0.05
				}

			}
		}
	}
}

var tiles int = 21
var tileRadius = 10

func main() {
	rand.Seed(time.Now().UnixNano())
	player, err := voxfile.DecodeFile("models/chr_sword.vox")
	log.Println(err)
	log.Println("Loaded character with size ", player.SizeX, player.SizeY, player.SizeZ)
	wall, _ := voxfile.DecodeFile("models/wall5.vox")
	eye, _ := voxfile.DecodeFile("models/eye.vox")

	var size int = 105

	palette = make([]mgl32.Vec4, 2000)
	for i := 0; i < 2000; i++ {
		palette[i] = mgl32.Vec4{
			rand.Float32(),
			rand.Float32(),
			rand.Float32(),
			1.0,
		}
	}
	InitGame(size)
	window, rv := govox.InitGraphics(size, 1900, 1000)
	//BlocksBuffer := govox.MakeBlocks(int(size))

	//	blocks := makeBlocks(int(size))

	lifeBlocks = make([]bool, int(size*size*size))
	Actrs = []Actor{
		Actor{[3]int{25, int(size) - 1, 25}, 1},
		Actor{[3]int{35, int(size) - 1, 25}, 1},
	}

	for i, _ := range lifeBlocks {
		lifeBlocks[i] = (rand.Float32() < 0.5)
	}

	BlocksBuffer := govox.MakeBlocks(int(size))
	LowResBlocks := govox.MakeBlocks(int(size) / 5)
	/*
		go func() {
			for {
				lifeBlocks = cycle(int(size), lifeBlocks, true)
				time.Sleep(1 * time.Second)
			}
		}()
	*/

	maze := GenerateMaze(300)
	middle := 150
	PlayerPos = Vec3{middle, 0, middle}

	//	angle := 0.0
	//previousTime := glfw.GetTime()
	//	texture, err := govox.NewTexture("square.png")
	//if err != nil {
	///	log.Fatalln(err)
	//}
	go handleKeys(window, maze)
	for !window.ShouldClose() {

		ClearDisplay(size, BlocksBuffer)
		ClearDisplay(size/5, LowResBlocks)
		//AddMaze(size/5, PlayerPos, wall, maze, LowResBlocks)
		AddMaze(size, PlayerPos, wall, maze, BlocksBuffer)
		AddFloor(size, maze, BlocksBuffer)
		DrawPlayer(size, PlayerPos, BlocksBuffer)
		magica2govox(size, Vec3{2 * size / 5, 0, 2 * size / 5}, player, BlocksBuffer)
		for _, m := range monsters {
			DrawMonster(size, m, BlocksBuffer)
			AddMonster(size, m, PlayerPos, eye, BlocksBuffer)
		}

		//blocks := lifeBlocks2Blocks(int(size), lifeBlocks, nil)

		AddActors(Actrs, BlocksBuffer)
		//BlocksBuffer = rise(int(size), BlocksBuffer)

		// globals
		//		gl.Enable(gl.DEPTH_TEST)
		//		gl.DepthFunc(gl.LESS)
		//gl.ClearColor(0.8, 0.8, 1.0, 1.0)
		//		gl.ClearColor(0.0, 0.0, 0.0, 1.0)

		//screenshot("voxeltest.png", 4000, 2000)
		govox.StartRender()
		govox.SetCam(size, rv.Program)
		govox.RenderBlocks(&rv, window, &BlocksBuffer, rotx, roty, int(size))
		govox.SetCam(size/5, rv.Program)
		//govox.RenderBlocks(&rv, window, &LowResBlocks, rotx, roty, int(size)/5)
		govox.FinishRender(window)

		//AddFourier(size, BlocksBuffer)
		//DrawCustom(rv, window, BlocksBuffer, rotx, roty, 0, 0, 0, size)

	}

	log.Println("Finished!")
}

func ClearDisplay(size int, blocks voxMap) {
	mapBlock(size, func(b govox.Block, i, j, k int) govox.Block {
		b.Active = false
		return b
	}, blocks)
}

func DrawCustom(rv *govox.RenderVars, window *glfw.Window, blocks voxMap, rotx, roty float32, fi, fj, fk float32, size int) {
	im, _ := glim.DrawStringRGBA(20, glim.RGBA{1.0, 1.0, 1.0, 1.0}, fmt.Sprintf("Monsters Remaining: %v", len(monsters)), "Asdfasdf")
	pic, w, h := glim.GFormatToImage(im, nil, 0, 0)
	//log.Printf("Rotx: %v, roty: %v\n", rotx, roty)
	//glim.DumpBuff(pic, uint(w), uint(h))
	/*pic := glim.RandPic(20, 20)
	w := 20
	h := 20*/
	scale := float32(0.05)

	govox.DrawAny(rv, window, rotx, roty, fi, fj, fk, func() {
		gl.Clear(gl.COLOR_BUFFER_BIT | gl.DEPTH_BUFFER_BIT)
		//		RenderBlocks(rv, window, blocks, rotx, roty, size)
		rotx = -3.449998
		roty = 0.7500001
		for i := float32(0.0); i < float32(w); i = i + 1.0 {
			for j := float32(0.0); j < float32(h); j = j + 1.0 {
				if pic[4*(int(i)+int(j)*int(w))] > 0 {

					model := mgl32.Ident4()
					modelUni := gl.GetUniformLocation(rv.Program, gl.Str("model\x00"))
					gl.UniformMatrix4fv(modelUni, 1, false, &model[0])

					//screenshot("voxeltest.png", 4000, 2000)

					//Color := mgl32.Vec4{1.0, 0.0, 0.0, 1.0}
					//gl.Uniform4fv(rv.ColUni, 1, &Color[0])
					model = mgl32.HomogRotate3DY(roty)
					model = model.Mul4(mgl32.HomogRotate3DX(rotx))
					model = model.Mul4(mgl32.Scale3D(scale, scale, scale))
					model = model.Mul4(mgl32.Translate3D(fi+i, fj+j, fk))

					gl.UniformMatrix4fv(modelUni, 1, false, &model[0])
					gl.BindVertexArray(rv.Vao)
					gl.DrawArrays(gl.TRIANGLES, 0, 6*2*3)
				}
			}
		}
	})
}
