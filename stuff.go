package govox

import (
	"errors"
	"time"

	"log"
	"math/rand"
	"strings"

	"github.com/donomii/glim"
	"github.com/go-gl/gl/v3.3-core/gl"
	"github.com/go-gl/glfw/v3.2/glfw"
	"github.com/go-gl/mathgl/mgl32"
)

type Block struct {
	Active bool
	Color  mgl32.Vec4
}
type RenderVars struct {
	Col         mgl32.Vec4
	ColUni      int32
	Vao         uint32
	Vbo         uint32
	VertAttrib  uint32
	Vaoc        uint32
	Vboc        uint32
	VertAttribc uint32
	Program     uint32
}

var BlocksBuffer [][][]Block

var saveCount int = 1

func Screenshot(filename string, width, height int) {
	gl.ReadBuffer(gl.BACK_LEFT)
	data := make([]byte, width*height*4)
	//data[0], data[1], data[2] = 123, 213, 132 // Test if it's overwritten
	gl.ReadPixels(0, 0, int32(width), int32(height), gl.RGBA, gl.UNSIGNED_BYTE, gl.Ptr(data))
	//fmt.Println("Read at", 0, 0, data)
	saveCount += 1
	glim.SaveBuff(int(width), int(height), data, filename)

}

func ScreenshotBuff(width, height int) []byte {
	gl.ReadBuffer(gl.BACK_LEFT)
	data := make([]byte, width*height*4)
	//data[0], data[1], data[2] = 123, 213, 132 // Test if it's overwritten
	gl.ReadPixels(0, 0, int32(width), int32(height), gl.RGBA, gl.UNSIGNED_BYTE, gl.Ptr(data))
	//fmt.Println("Read at", 0, 0, data)
	saveCount += 1
	return data

}

func SetCam(size int, p uint32) {
	cam := mgl32.LookAtV(mgl32.Vec3{float32(size) * 1.8, float32(size) * 1.5, float32(size) * 2}, mgl32.Vec3{0, 0, 0}, mgl32.Vec3{0, 1, 0})
	camUni := gl.GetUniformLocation(p, gl.Str("camera\x00"))
	gl.UniformMatrix4fv(camUni, 1, false, &cam[0])
}
func InitGraphics(size int, width, height int) (*glfw.Window, RenderVars) {
	if err := glfw.Init(); err != nil {
		log.Fatal(err)
	}

	glfw.WindowHint(glfw.Resizable, glfw.False)
	glfw.WindowHint(glfw.ContextVersionMajor, 3)
	glfw.WindowHint(glfw.ContextVersionMinor, 3)
	glfw.WindowHint(glfw.OpenGLProfile, glfw.OpenGLCoreProfile)
	glfw.WindowHint(glfw.OpenGLForwardCompatible, glfw.True)

	window, err := glfw.CreateWindow(width, height, "govox", nil, nil)
	if err != nil {
		log.Fatal(err)
	}
	window.MakeContextCurrent()
	if err := gl.Init(); err != nil {
		log.Fatal(err)
	}

	log.Printf("OpenGL version %s", gl.GoStr(gl.GetString(gl.VERSION)))

	// Set up the program
	p, err := NewProgram(vertexShader, fragmentShader)
	if err != nil {
		log.Fatal(err)
	}
	gl.UseProgram(p)

	proj := mgl32.Perspective(mgl32.DegToRad(45.0), float32(width)/float32(height), 0.1, 1000.0)
	projUni := gl.GetUniformLocation(p, gl.Str("projection\x00"))
	gl.UniformMatrix4fv(projUni, 1, false, &proj[0])

	SetCam(size, p)
	col := mgl32.Vec4{0, 0, 0, 1}
	colUni := gl.GetUniformLocation(p, gl.Str("col\x00"))
	gl.Uniform4fv(colUni, 1, &col[0])

	// Vertex data
	var vao uint32
	gl.GenVertexArrays(1, &vao)
	gl.BindVertexArray(vao)

	var vbo uint32
	gl.GenBuffers(1, &vbo)
	gl.BindBuffer(gl.ARRAY_BUFFER, vbo)
	gl.BufferData(gl.ARRAY_BUFFER, len(cubeVerts)*4, gl.Ptr(cubeVerts), gl.STATIC_DRAW)

	vertAttrib := uint32(gl.GetAttribLocation(p, gl.Str("vert\x00")))
	gl.EnableVertexAttribArray(vertAttrib)
	gl.VertexAttribPointer(vertAttrib, 3, gl.FLOAT, false, 3*4, gl.PtrOffset(0))

	// Colour data
	var vaoc uint32
	gl.GenVertexArrays(1, &vaoc)
	gl.BindVertexArray(vaoc)

	var vboc uint32
	gl.GenBuffers(1, &vboc)
	gl.BindBuffer(gl.ARRAY_BUFFER, vboc)
	gl.BufferData(gl.ARRAY_BUFFER, len(cubeVerts)*4, gl.Ptr(cubeVerts), gl.STATIC_DRAW)

	vertAttribc := uint32(gl.GetAttribLocation(p, gl.Str("colour\x00")))
	gl.EnableVertexAttribArray(vertAttribc)
	gl.VertexAttribPointer(vertAttribc, 4, gl.FLOAT, false, 4*4, gl.PtrOffset(0))

	rv := RenderVars{col, colUni, vao, vbo, vertAttrib, vaoc, vboc, vertAttribc, p}

	gl.BindFragDataLocation(p, 0, gl.Str("outputColor\x00"))

	readyCh = make(chan *RenderData, 2)
	cycleCh = make(chan *RenderData, 2)
	finishCh = make(chan *RenderData, 2)

	rd := RenderData{}
	rd.Points = make([]float32, 200000)
	rd.PointsLength = 0
	rd.Colours = make([]float32, 200000)
	rd.ColoursLength = 0
	//rd.Blocks = blocks
	finishCh <- &rd

	rd = RenderData{}
	rd.Points = make([]float32, 200000)
	rd.PointsLength = 0
	rd.Colours = make([]float32, 200000)
	rd.ColoursLength = 0
	//rd.Blocks = blocks
	finishCh <- &rd

	go RenderPrepWorker(size, cycleCh, readyCh)
	return window, rv
}

func ShutdownGraphics() {
	defer glfw.Terminate()
}

func MakeBlocks(size int) [][][]Block {
	// initialize blocks
	blocks := make([][][]Block, size)
	for i := 0; i < size; i++ {
		blocks[i] = make([][]Block, size)
		for j := 0; j < size; j++ {
			blocks[i][j] = make([]Block, size)
			for k := 0; k < size; k++ {
				blocks[i][j][k] = Block{
					Active: false,
					Color: mgl32.Vec4{
						0.0,
						0.0,
						0.0,
						0.0,
					},
				}
			}
		}
	}
	return blocks
}

func MakeRandomBlocks(size int) [][][]Block {
	// initialize blocks
	blocks := make([][][]Block, size)
	for i := 0; i < size; i++ {
		blocks[i] = make([][]Block, size)
		for j := 0; j < size; j++ {
			blocks[i][j] = make([]Block, size)
			for k := 0; k < size; k++ {
				blocks[i][j][k] = Block{
					Active: (rand.Float32() < 0.1),
					Color: mgl32.Vec4{
						rand.Float32(),
						rand.Float32(),
						rand.Float32(),
						1.0,
					},
				}
			}
		}
	}
	return blocks
}

func StartRender() {
	// globals
	gl.Enable(gl.DEPTH_TEST)
	gl.DepthFunc(gl.LESS)
	//gl.ClearColor(0.8, 0.8, 1.0, 1.0)
	gl.ClearColor(0.0, 0.0, 0.0, 1.0)

	//screenshot("voxeltest.png", 4000, 2000)

	gl.Clear(gl.COLOR_BUFFER_BIT | gl.DEPTH_BUFFER_BIT)

}

func FinishRender(window *glfw.Window) {
	window.SwapBuffers()

}

var startFrame time.Time

type RenderData struct {
	Points        []float32
	Colours       []float32
	ColoursLength int
	PointsLength  int
	Blocks        [][][]Block
	rotx, roty    float32
}

func RenderPrepWorker(size int, cycleCh, readyCh chan *RenderData) {
	go func() {
		for {
			startFrame := time.Now()
			rd := <-cycleCh
			if rd.Blocks != nil {
				points := rd.Points
				colours := rd.Colours
				pointsi := 0
				coloursi := 0
				blocks := rd.Blocks
				offset := float32(size / 2)
				for i := 0; i < size; i++ {
					for j := 0; j < size; j++ {
						for k := 0; k < size; k++ {
							b := blocks[i][j][k]
							if !b.Active {
								continue
							}

							fi := float32(i) - float32(size)/2 - offset
							fj := float32(j) - float32(size)/2 - offset
							fk := float32(k) - float32(size)/2 - offset

							//Copy voxel locations into VBO
							points[pointsi] = fi
							pointsi = pointsi + 1
							points[pointsi] = fj
							pointsi = pointsi + 1
							points[pointsi] = fk
							pointsi = pointsi + 1

							//Copy colours into VBO buffer
							colours[coloursi] = b.Color.X()
							coloursi = coloursi + 1
							colours[coloursi] = b.Color.Y()
							coloursi = coloursi + 1
							colours[coloursi] = b.Color.Z()
							coloursi = coloursi + 1
							colours[coloursi] = b.Color.W()
							coloursi = coloursi + 1

						}
					}
				}

				rd.Points = points
				rd.Colours = colours
				rd.PointsLength = pointsi
				rd.ColoursLength = coloursi
				readyCh <- rd
			}
			log.Println("Prepared", size*size*size, "blocks in", (time.Now().Sub(startFrame)).Nanoseconds()/1000000)

		}
	}()
}

var readyCh chan *RenderData
var cycleCh chan *RenderData
var finishCh chan *RenderData

func RenderBlocks(rv *RenderVars, blocks [][][]Block, rotx, roty float32, size int) {

	rd := <-finishCh
	rd.Blocks = blocks
	rd.roty = roty
	rd.rotx = rotx
	cycleCh <- rd

}

func GlRenderer(size int, rv *RenderVars, window *glfw.Window) {
	select {
	case rd1 := <-readyCh:
		startFrame = time.Now()

		StartRender()
		SetCam(size/2, rv.Program)

		model := mgl32.Ident4()
		modelUni := gl.GetUniformLocation(rv.Program, gl.Str("model\x00"))

		model = mgl32.HomogRotate3DY(rd1.roty)
		model = model.Mul4(mgl32.HomogRotate3DX(rd1.rotx))
		gl.UniformMatrix4fv(modelUni, 1, false, &model[0])
		gl.PointSize(8)

		gl.BindBuffer(gl.ARRAY_BUFFER, rv.Vbo)
		gl.BufferData(gl.ARRAY_BUFFER, (rd1.PointsLength)*4, gl.Ptr(rd1.Points), gl.STATIC_DRAW)

		gl.EnableVertexAttribArray(rv.VertAttrib)
		gl.VertexAttribPointer(rv.VertAttrib, 3, gl.FLOAT, false, 3*4, gl.PtrOffset(0))

		gl.BindBuffer(gl.ARRAY_BUFFER, rv.Vboc)
		gl.BufferData(gl.ARRAY_BUFFER, rd1.ColoursLength*4, gl.Ptr(rd1.Colours), gl.STATIC_DRAW)
		gl.EnableVertexAttribArray(rv.VertAttribc)
		gl.VertexAttribPointer(rv.VertAttribc, 4, gl.FLOAT, false, 4*4, gl.PtrOffset(0))

		gl.DrawArrays(gl.POINTS, 0, int32(rd1.PointsLength))

		FinishRender(window)

		log.Println("Drew", rd1.PointsLength, "points in", (time.Now().Sub(startFrame)).Nanoseconds()/1000000)

		finishCh <- rd1
	default:
		glfw.PollEvents()
	}
	glfw.PollEvents()

}

func NewProgram(vSource, fSource string) (uint32, error) {
	vShader, err := CompileShader(vSource, gl.VERTEX_SHADER)
	if err != nil {
		return 0, err
	}

	fShader, err := CompileShader(fSource, gl.FRAGMENT_SHADER)
	if err != nil {
		return 0, err
	}

	p := gl.CreateProgram()

	gl.AttachShader(p, vShader)
	gl.AttachShader(p, fShader)
	gl.LinkProgram(p)

	var status int32
	gl.GetProgramiv(p, gl.LINK_STATUS, &status)
	if status == gl.FALSE {
		var ll int32
		gl.GetProgramiv(p, gl.INFO_LOG_LENGTH, &ll)

		l := strings.Repeat("\x00", int(ll+1))
		gl.GetProgramInfoLog(p, ll, nil, gl.Str(l))

		return 0, errors.New(l)
	}

	gl.DeleteShader(vShader)
	gl.DeleteShader(fShader)

	return p, nil
}

func CompileShader(s string, t uint32) (uint32, error) {
	shader := gl.CreateShader(t)

	cs, free := gl.Strs(s)
	gl.ShaderSource(shader, 1, cs, nil)
	free()
	gl.CompileShader(shader)

	var status int32
	gl.GetShaderiv(shader, gl.COMPILE_STATUS, &status)
	if status == gl.FALSE {
		var ll int32
		gl.GetShaderiv(shader, gl.INFO_LOG_LENGTH, &ll)

		l := strings.Repeat("\x00", int(ll+1))
		gl.GetShaderInfoLog(shader, ll, nil, gl.Str(l))

		return 0, errors.New(l)
	}

	return shader, nil
}
