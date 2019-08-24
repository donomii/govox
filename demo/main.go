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
	"github.com/donomii/myvox"
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
type voxMap [][][]myvox.Block

var Actrs []Actor
var PlayerPos Vec3
var palette []mgl32.Vec4

func magica2myvox(sizei int, pos Vec3, vox *voxfile.VoxFile, blocks voxMap) {
	defer func() {
		if r := recover(); r != nil {
			fmt.Printf("Recovered in magica2myvox, drawing at position %v: %v\n", pos, r)
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

func main() {
	rand.Seed(time.Now().UnixNano())
	player, err := voxfile.DecodeFile("models/chr_sword.vox")
	log.Println(err)
	log.Println("Loaded character with size ", player.SizeX, player.SizeY, player.SizeZ)
	//wall, _ := voxfile.DecodeFile("models/wall.vox")
	//eye, _ := voxfile.DecodeFile("models/eye.vox")

	var size int = 100.0
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
	window, rv := myvox.InitGraphics()
	BlocksBuffer = MakeBlocks(int(size))

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

	//maze := GenerateMaze(size)
	middle := size / 2
	PlayerPos = Vec3{middle, 0, middle}
	//var roty, rotx float32
	//lastInputTime := time.Now()

	angle := 0.0
	previousTime := glfw.GetTime()
	texture, err := myvox.NewTexture("square.png")
	if err != nil {
		log.Fatalln(err)
	}

	for !window.ShouldClose() {
		log.Println("Start loop")
		gl.Clear(gl.COLOR_BUFFER_BIT | gl.DEPTH_BUFFER_BIT)
		log.Println("Updating")
		//time.Sleep(1 * time.Second)
		// Update
		time := glfw.GetTime()
		elapsed := time - previousTime
		previousTime = time

		angle += elapsed

		model := mgl32.Ident4()
		modelUniform := gl.GetUniformLocation(rv.Program, gl.Str("model\x00"))
		gl.UniformMatrix4fv(modelUniform, 1, false, &model[0])
		model = mgl32.HomogRotate3D(float32(angle), mgl32.Vec3{0, 1, 0})
		log.Println("Rendering")
		// Render
		gl.UseProgram(rv.Program)
		gl.UniformMatrix4fv(modelUniform, 1, false, &model[0])

		gl.BindVertexArray(rv.Vao)

		gl.ActiveTexture(gl.TEXTURE0)
		gl.BindTexture(gl.TEXTURE_2D, texture)

		gl.DrawArrays(gl.TRIANGLES, 0, 6*2*3)
		log.Println("Swapping")
		// Maintenance
		window.SwapBuffers()
		glfw.PollEvents()
		log.Println("End loop")
	}

	/*
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
					if moveOk(wantPos, maze) {
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

			ClearDisplay(size, myvox.BlocksBuffer)
			AddMaze(size, PlayerPos, wall, maze, myvox.BlocksBuffer)
			AddFloor(size, maze, myvox.BlocksBuffer)
			DrawPlayer(size, PlayerPos, myvox.BlocksBuffer)
			magica2myvox(size, Vec3{2 * size / 5, 0, 2 * size / 5}, player, myvox.BlocksBuffer)
			for _, m := range monsters {
				DrawMonster(size, m, myvox.BlocksBuffer)
				AddMonster(size, m, PlayerPos, eye, myvox.BlocksBuffer)
			}

			//blocks := lifeBlocks2Blocks(int(size), lifeBlocks, nil)

			AddActors(Actrs, myvox.BlocksBuffer)
			//myvox.BlocksBuffer = rise(int(size), myvox.BlocksBuffer)
			myvox.Renderblocks(rv, window, myvox.BlocksBuffer, rotx, roty, int(size))
			//AddFourier(size, myvox.BlocksBuffer)
			//DrawCustom(rv, window, myvox.BlocksBuffer, rotx, roty, 0, 0, 0, size)

		}
	*/
	log.Println("Finished!")
}

func ClearDisplay(size int, blocks voxMap) {
	mapBlock(size, func(b myvox.Block, i, j, k int) myvox.Block {
		b.Active = false
		return b
	}, blocks)
}
func RenderBlocks(rv myvox.RenderVars, window *glfw.Window, blocks voxMap, rotx, roty float32, size int) {

	model := mgl32.Ident4()
	//modelUni := gl.GetUniformLocation(rv.Program, gl.Str("model\x00"))
	//gl.UniformMatrix4fv(modelUni, 1, false, &model[0])

	for i := 0; i < size; i++ {
		for j := 0; j < size; j++ {
			for k := 0; k < size; k++ {
				b := blocks[i][j][k]
				if !b.Active {
					continue
				}

				//gl.Uniform4fv(rv.ColUni, 1, &b.Color[0])

				fi := float32(i) - float32(size)/2
				fj := float32(j) - float32(size)/2
				fk := float32(k) - float32(size)/2

				model = mgl32.HomogRotate3DY(roty)
				model = model.Mul4(mgl32.HomogRotate3DX(rotx))
				model = model.Mul4(mgl32.Translate3D(fi, fj, fk))
				model = model.Mul4(mgl32.Scale3D(0.5, 0.5, 0.5))

				//gl.UniformMatrix4fv(modelUni, 1, false, &model[0])
				gl.BindVertexArray(rv.Vao)
				gl.DrawArrays(gl.TRIANGLES, 0, 6*2*3)
			}
		}
	}
}

func DrawCustom(rv myvox.RenderVars, window *glfw.Window, blocks voxMap, rotx, roty float32, fi, fj, fk float32, size int) {
	im, _ := glim.DrawStringRGBA(20, glim.RGBA{1.0, 1.0, 1.0, 1.0}, fmt.Sprintf("Monsters Remaining: %v", len(monsters)), "Asdfasdf")
	pic, w, h := glim.GFormatToImage(im, nil, 0, 0)
	//log.Printf("Rotx: %v, roty: %v\n", rotx, roty)
	//glim.DumpBuff(pic, uint(w), uint(h))
	/*pic := glim.RandPic(20, 20)
	w := 20
	h := 20*/
	scale := float32(0.05)

	myvox.DrawAny(rv, window, rotx, roty, fi, fj, fk, func() {
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
