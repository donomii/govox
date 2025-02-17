package main

import (
	"io/ioutil"
	"strings"

	//	"math"

	"fmt"
	"log"

	//"math"
	"math/rand"
	"runtime"

	_ "log"
	"time"

	"github.com/donomii/glim"
	"github.com/tbogdala/Voxfile"

	"github.com/donomii/govox"
	"github.com/go-gl/glfw/v3.2/glfw"
	"github.com/go-gl/mathgl/mgl32"
)

var finalLockout time.Time
var mode = "game"
var roty, rotx float32

var textOffset = 300

func init() {
	runtime.LockOSThread()
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

func voxFlipX(vox *voxfile.VoxFile) {
	for _, v := range vox.Voxels {
		v.X = 5 - v.X
	}
}

func voxFlipY(vox *voxfile.VoxFile) {
	for _, v := range vox.Voxels {
		v.Y = 5 - v.Y
	}
}

func voxFlipZ(vox *voxfile.VoxFile) {
	for _, v := range vox.Voxels {
		v.Z = 5 - v.Z
	}
}

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

		blocks[v.X+x][v.Z+y][v.Y+z].Active = true
		blocks[v.X+x][v.Z+y][v.Y+z].Color = mgl32.Vec4{float32(vox.Palette[v.Index].R) / 255, float32(vox.Palette[v.Index].G) / 255, float32(vox.Palette[v.Index].B) / 255, 1.0}
		//blocks[v.X+x][v.Z+y][v.Y+z].Color = palette[v.Index]

	}

}

func handleKeys(window *glfw.Window, maze [][]int) {
	lastInputTime := time.Now()
	lastInputTime2 := time.Now()

	for {
		if glfw.Press == 1 {
			if time.Now().Sub(lastInputTime).Nanoseconds() > 150000000 {
				lastInputTime = time.Now()
				wantPos := PlayerPos
				if mode == "game" {
					if window.GetKey(glfw.KeyUp) == glfw.Press || window.GetKey(glfw.KeyK) == glfw.Press {
						wantPos[2] = wantPos[2] - 1
					}
					if window.GetKey(glfw.KeyDown) == glfw.Press || window.GetKey(glfw.KeyJ) == glfw.Press {
						wantPos[2] = wantPos[2] + 1
					}
					if window.GetKey(glfw.KeyRight) == glfw.Press || window.GetKey(glfw.KeyL) == glfw.Press {
						wantPos[0] = wantPos[0] + 1
					}
					if window.GetKey(glfw.KeyLeft) == glfw.Press || window.GetKey(glfw.KeyH) == glfw.Press {
						wantPos[0] = wantPos[0] - 1
					}
					if moveOk(wantPos, maze) {
						PlayerPos = wantPos
					}
					monsters = handleCollision(PlayerPos, wantPos)
					if len(monsters) < 1 {
						mode = "finish"
						textOffset = 570
						finalLockout = time.Now()
					}
					//log.Println("Player position", PlayerPos)
				} else {
					if time.Now().Sub(finalLockout).Seconds() > 3 && (window.GetKey(glfw.KeyLeft) == glfw.Press || window.GetKey(glfw.KeyRight) == glfw.Press || window.GetKey(glfw.KeyDown) == glfw.Press || window.GetKey(glfw.KeyUp) == glfw.Press || window.GetKey(glfw.KeyK) == glfw.Press || window.GetKey(glfw.KeyJ) == glfw.Press || window.GetKey(glfw.KeyL) == glfw.Press || window.GetKey(glfw.KeyH) == glfw.Press) {
						InitAll(105, maze)
					}
				}
			}
			if time.Now().Sub(lastInputTime2).Nanoseconds() > 15000000 {
				lastInputTime2 = time.Now()

				// Rotate view
				if window.GetKey(glfw.KeyA) == glfw.Press {
					roty -= 0.05
				}

				if window.GetKey(glfw.KeyD) == glfw.Press {
					roty += 0.05
				}

				if window.GetKey(glfw.KeyW) == glfw.Press {
					rotx -= 0.05
				}

				if window.GetKey(glfw.KeyS) == glfw.Press {
					rotx += 0.05
				}

			}
		}
	}
}

var tiles int = 21
var tileRadius = 10
var blocksSize = 105

func main() {
	rand.Seed(time.Now().UnixNano())
	var size int = 105
	window, rv := govox.InitGraphics(size, 800, 600)
	maze := GenerateMaze(125, 125)
	markov := GenerateMarkov(blocksSize, blocksSize, blocksSize)
	markovB := GenerateMarkov(blocksSize, blocksSize, blocksSize)
	InitAll(size, maze)
	go handleKeys(window, maze)
	go BlocksWorker(size, &rv, maze, markov, markovB)
	for !window.ShouldClose() {
		govox.GlRenderer(size, &rv, window)
	}
	log.Println("Finished!")

}

func InitAll(size int, maze [][]int) {
	mode = "game"

	//FIXME use standard magicavox palette
	palette = make([]mgl32.Vec4, 2000)
	for i := 0; i < 2000; i++ {
		palette[i] = mgl32.Vec4{
			rand.Float32(),
			rand.Float32(),
			rand.Float32(),
			1.0,
		}
	}

	InitGame(size, maze)

	lifeBlocks = make([]bool, int(size*size*size))
	Actrs = []Actor{
		Actor{[3]int{25, int(size) - 1, 25}, 1},
		Actor{[3]int{35, int(size) - 1, 25}, 1},
	}

	for i, _ := range lifeBlocks {
		lifeBlocks[i] = (rand.Float32() < 0.5)
	}

	//Read map from TSV file
	raw, _ := ioutil.ReadFile("map.tsv")
	mapstr := string(raw)
	mapstrs := strings.Split(mapstr, "\n")
	xoffset := tileRadius + 1
	yoffset := tileRadius + 1
	for y, v := range mapstrs {
		cols := strings.Split(v, "	")
		for x, v := range cols {
			//log.Println(x, y)
			if v == "" {
				maze[y+yoffset][x+xoffset] = 1
				fmt.Printf("*")
			} else {
				maze[y+yoffset][x+xoffset] = 0
				fmt.Printf(" ")
			}
		}
		fmt.Printf("\n")
	}

	PlayerPos = Vec3{27, 0, 19}

}

func LoadVox(file string) *voxfile.VoxFile {
	player, err := voxfile.DecodeFile(file)
	if err != nil {
		log.Println(err)
	}
	log.Println("Loaded character with size ", player.SizeX, player.SizeY, player.SizeZ)
	return player
}

func AddText(size int, blocks voxMap) {
	message := "Markov voxel demo"
	colour := mgl32.Vec4{1.0, 1.0, 1.0, 1.0}
	if mode == "finish" {
		message = "Good.  Do it again."
		colour = mgl32.Vec4{0.0, 1.0, 0.0, 1.0}
	}
	im, _ := glim.DrawStringRGBA(20, glim.RGBA{1.0, 1.0, 1.0, 1.0}, message, "Asdfasdf")
	pic, w, h := glim.GFormatToImage(im, nil, 0, 0)

	fi := 10
	fj := 10
	fk := 10
	textOffset = textOffset + 1
	if !(textOffset < w) {
		textOffset = 0
	}
	for i := 0; i < w && fi+i < size; i = i + 1.0 {
		for j := 0; j < h && fj+j < size; j = j + 1.0 {
			xpos := i + textOffset
			if xpos > w {
				xpos = xpos - w
			}
			if pic[4*(xpos+j*w)] > 0 {
				blocks[fi+i][size-(fj+j)][fk].Active = true
				blocks[fi+i][size-(fj+j)][fk].Color = colour
			}
		}
	}
}

var lastRuleTime time.Time
var currentRule int
var rules []MarkovRule

func CopyMarkov(size int, a, b [][][]int) [][][]int {
	for i := 0; i < size; i++ {
		for j := 0; j < size; j++ {
			for k := 0; k < size; k++ {
				b[i][j][k] = a[i][j][k]
			}
		}
	}
	return b
}
func BlocksWorker(size int, rv *govox.RenderVars, maze [][]int, markov, markovB [][][]int) {
	player := LoadVox("models/chr_player.vox")
	wall := LoadVox("models/wall5.vox")
	monster := LoadVox("models/monster.vox")
	markov[4][4][4] = 1                //Move me
	markov[size-4][4][4] = 1           //Move me
	markov[size-4][size-4][4] = 1      //Move me
	markov[size-4][size-4][size-4] = 1 //Move me
	voxFlipY(monster)

	rules = []MarkovRule{
		//MarkovRule{From: []int{1, 0, 0}, To: []int{2, 2, 1}},
		MarkovRule{From: []int{1, 0, 0}, To: []int{2, 2, 1}},
		MarkovRule{From: []int{1, 2, 2}, To: []int{3, 3, 1}},
	}

	for {
		startFrame := time.Now()
		if time.Now().Sub(lastRuleTime).Milliseconds() > 1 {

			changed := ApplyRule(blocksSize, markov, markovB, rules[currentRule], false /* don't do all changes at once*/, 1)
			CopyMarkov(size, markovB, markov)
			if changed == 0 {
				currentRule = (currentRule + 1) % len(rules)
			}
			lastRuleTime = time.Now()
		}
		//ClearDisplay(size, BlocksBuffer)
		BlocksBuffer := govox.MakeBlocks(size)
		//AddText(size, BlocksBuffer)
		//ClearDisplay(size/5, LowResBlocks)
		//AddMaze(size/5, PlayerPos, wall, maze, LowResBlocks)

		//AddMaze(size, PlayerPos, wall, maze, BlocksBuffer)

		AddMarkov(blocksSize, PlayerPos, wall, markov, BlocksBuffer)
		//DrawAbstractMaze(size, maze, BlocksBuffer)
		//DrawPlayer(size, PlayerPos, BlocksBuffer)
		magica2govox(size, Vec3{size / 2, 0, size / 2}, player, BlocksBuffer)
		//for _, m := range monsters {
		//DrawAbstractMonster(size, m, BlocksBuffer)
		//AddMonster(size, m, PlayerPos, monster, BlocksBuffer)
		//}

		//AddActors(Actrs, BlocksBuffer)
		//BlocksBuffer = rise(int(size), BlocksBuffer)

		//screenshot("voxeltest.png", 4000, 2000)
		if govox.ShowTimings {
			log.Println("Calculated ", size*size*size, "blocks in", (time.Now().Sub(startFrame)).Nanoseconds()/1000000)
		}
		govox.RenderBlocks(rv, BlocksBuffer, rotx, roty, int(size), true)
		//govox.SetCam(size/5, rv.Program)
		//govox.RenderBlocks(&rv, window, &LowResBlocks, rotx, roty, int(size)/5)

		//AddFourier(size, BlocksBuffer)
		//DrawCustom(rv, window, BlocksBuffer, rotx, roty, 0, 0, 0, size)

	}
}
