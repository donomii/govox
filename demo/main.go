package main

import (
	//	"log"
	"fmt"
	//"math"
	"math/rand"
	"runtime"

	_ "log"
	"time"

	"github.com/go-gl/gl/v3.3-core/gl"

	"github.com/donomii/glim"
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
	Position Vec3
	Type     int
}

type Vec3 [3]int
type voxMap [][][]govox.Block

var Actrs []Actor
var PlayerPos Vec3

func main() {
	rand.Seed(time.Now().UnixNano())

	var size int = 20.0
	InitGame(size)
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
	/*
		go func() {
			for {
				lifeBlocks = cycle(int(size), lifeBlocks, true)
				time.Sleep(1 * time.Second)
			}
		}()
	*/

	maze := GenerateMaze(size)
	middle := size / 2
	PlayerPos = Vec3{middle, 0, middle}
	var roty, rotx float32
	lastInputTime := time.Now()

	for !window.ShouldClose() {

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
				if moveOk(wantPos, govox.BlocksBuffer) {
					PlayerPos = wantPos
				} else {
					monsters = handleCollision(PlayerPos, wantPos)
				}

			}
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
		AddFloor(size, maze, govox.BlocksBuffer)
		DrawPlayer(size, PlayerPos, govox.BlocksBuffer)
		for _, m := range monsters {
			DrawMonster(size, m, govox.BlocksBuffer)
		}

		//blocks := lifeBlocks2Blocks(int(size), lifeBlocks, nil)

		AddActors(Actrs, govox.BlocksBuffer)
		//govox.BlocksBuffer = rise(int(size), govox.BlocksBuffer)
		//govox.Renderblocks(rv, window, govox.BlocksBuffer, rotx, roty, int(size))

		DrawCustom(rv, window, govox.BlocksBuffer, rotx, roty, 0, 0, 0, size)

	}
}

func RenderBlocks(rv govox.RenderVars, window *glfw.Window, blocks voxMap, rotx, roty float32, size int) {

	model := mgl32.Ident4()
	modelUni := gl.GetUniformLocation(rv.Program, gl.Str("model\x00"))
	gl.UniformMatrix4fv(modelUni, 1, false, &model[0])

	for i := 0; i < size; i++ {
		for j := 0; j < size; j++ {
			for k := 0; k < size; k++ {
				b := blocks[i][j][k]
				if !b.Active {
					continue
				}

				gl.Uniform4fv(rv.ColUni, 1, &b.Color[0])

				fi := float32(i) - float32(size)/2
				fj := float32(j) - float32(size)/2
				fk := float32(k) - float32(size)/2

				model = mgl32.HomogRotate3DY(roty)
				model = model.Mul4(mgl32.HomogRotate3DX(rotx))
				model = model.Mul4(mgl32.Translate3D(fi, fj, fk))
				model = model.Mul4(mgl32.Scale3D(0.5, 0.5, 0.5))

				gl.UniformMatrix4fv(modelUni, 1, false, &model[0])
				gl.BindVertexArray(rv.Vao)
				gl.DrawArrays(gl.TRIANGLES, 0, 6*2*3)
			}
		}
	}
}

func DrawCustom(rv govox.RenderVars, window *glfw.Window, blocks voxMap, rotx, roty float32, fi, fj, fk float32, size int) {
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
		RenderBlocks(rv, window, blocks, rotx, roty, size)
		rotx = -3.449998
		roty = 0.7500001
		for i := float32(0.0); i < float32(w); i = i + 1.0 {
			for j := float32(0.0); j < float32(h); j = j + 1.0 {
				if pic[4*(int(i)+int(j)*int(w))] > 0 {

					model := mgl32.Ident4()
					modelUni := gl.GetUniformLocation(rv.Program, gl.Str("model\x00"))
					gl.UniformMatrix4fv(modelUni, 1, false, &model[0])

					//screenshot("voxeltest.png", 4000, 2000)

					Color := mgl32.Vec4{1.0, 0.0, 0.0, 1.0}
					gl.Uniform4fv(rv.ColUni, 1, &Color[0])
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
