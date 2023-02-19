package windowing

import (
	"fmt"
	"github.com/UpsilonDiesBackwards/behngine_epq/camera"
	"github.com/UpsilonDiesBackwards/behngine_epq/input"
	"github.com/UpsilonDiesBackwards/behngine_epq/shaders"
	"github.com/go-gl/gl/all-core/gl"
	"github.com/go-gl/glfw/v3.2/glfw"
	"github.com/go-gl/mathgl/mgl32"
	"log"
)

type GLWindow struct {
	width  int
	height int
	glfw   *glfw.Window

	inputManager *input.UserInput
}

var (
	square = []float32{
		// first triangle
		0.5, 0.5, 0.0,
		0.5, -0.5, 0.0,
		-0.5, 0.5, 0.0,

		0.5, -0.5, 0.0,
		-0.5, -0.5, 0.0,
		-0.5, 0.5, 0.0,

		// second triangle
		-0.5, 0.5, 0,
		-0.5, -0.5, 0,
		0.5, -0.5, 0,

		-0.5, 0.5, 0,
		0.5, 0.5, 0,
		0.5, -0.5, 0,
	}
	square_i = []uint32{
		0, 1, 2, // first triangle of square
		2, 3, 0, // second triangle of square
	}
)

type PerspectiveBlock struct {
	project *mgl32.Mat4
	camera  *mgl32.Mat4
	world   *mgl32.Mat4
}

var fov = float32(60.0)
var projectionTransform = mgl32.Perspective(mgl32.DegToRad(fov),
	float32(800/600), 0.1, 2000)

var UBO uint32

func (w *GLWindow) InputManager() *input.UserInput {
	return w.inputManager
}

var inputMan = *input.CreateInput_Manager()
var program = initGL()
var VAO uint32

func NewWindow(width, height int, title string) *GLWindow {
	fmt.Printf("Initializing GLFW... ")
	if err := glfw.Init(); err != nil {
		panic(err)
	} else {
		fmt.Println("Success!")
	}

	// Establish window hints
	glfw.WindowHint(glfw.Resizable, glfw.False)
	glfw.WindowHint(glfw.ContextVersionMajor, 4) // OR 2
	glfw.WindowHint(glfw.ContextVersionMinor, 1)
	glfw.WindowHint(glfw.OpenGLProfile, glfw.OpenGLCoreProfile)
	glfw.WindowHint(glfw.OpenGLForwardCompatible, glfw.True)

	// Create new appWindow
	appWindow, err := glfw.CreateWindow(
		width,
		height,
		title,
		nil, nil)
	if err != nil {
		panic(err)
	}
	appWindow.MakeContextCurrent()
	appWindow.SetInputMode(glfw.CursorMode, glfw.CursorDisabled)

	VAO := createVAO(square, square_i)
	if VAO == 0 {
		fmt.Println("VAO is 0! Failed to create VAO")
	}

	return &GLWindow{
		width:  width,
		height: height,
		glfw:   appWindow,
	}
}

func (appWindow *GLWindow) Execute(camera *camera.Camera) {
	appWindow.glfw.SetKeyCallback(inputMan.KeyCallback)
	appWindow.glfw.SetCursorPosCallback(inputMan.MouseCallBack)

	ViewportTransform := camera.GetTransform()
	world := createWorldMatrix() // Create world matrix

	block := createPerspectiveBlock(&projectionTransform, &ViewportTransform, &world)
	UBO = createUBO(block)

	drawWindowContent(VAO, program)
}

func initGL() uint32 {
	fmt.Printf("Initializing OpenGL... ")
	if err := gl.Init(); err != nil {
		panic(err)
	} else {
		fmt.Println("Success!")
	}

	_GLVersion := gl.GoStr(gl.GetString(gl.VERSION)) // Get current OpenGL version
	log.Println("Using OpenGL version:", _GLVersion)

	// Create vShader and fShader
	vShader, err := shaders.ShaderCompiler(shaders.SHADER_DEFAULT_V, gl.VERTEX_SHADER)
	if err != nil {
		panic(err)
	}

	fShader, err := shaders.ShaderCompiler(shaders.SHADER_DEFAULT_F, gl.FRAGMENT_SHADER)
	if err != nil {
		panic(err)
	}

	// Create a new shader program, then attach shaders, then link and use it.
	glProgram := gl.CreateProgram()
	gl.AttachShader(glProgram, vShader)
	gl.AttachShader(glProgram, fShader)
	gl.LinkProgram(glProgram)
	gl.UseProgram(glProgram)
	return glProgram
}

func createVAO(vertices []float32, indices []uint32) uint32 {
	var VAO uint32
	gl.GenVertexArrays(1, &VAO)
	gl.BindVertexArray(VAO)

	var VBO uint32
	gl.GenBuffers(1, &VBO)
	gl.BindBuffer(gl.ARRAY_BUFFER, VBO)
	gl.BufferData(gl.ARRAY_BUFFER, len(vertices)*4, gl.Ptr(vertices), gl.STATIC_DRAW)

	var EBO uint32
	gl.GenBuffers(1, &EBO)
	gl.BindBuffer(gl.ELEMENT_ARRAY_BUFFER, EBO)
	gl.BufferData(gl.ELEMENT_ARRAY_BUFFER, len(indices)*4, gl.Ptr(indices), gl.STATIC_DRAW)

	//gl.VertexAttribPointer(0, 3, gl.FLOAT, false, 3*4, gl.PtrOffset(0))
	gl.VertexAttribPointer(0, 3, gl.FLOAT, false, 0, gl.PtrOffset(0))
	gl.EnableVertexAttribArray(0)

	gl.BindBuffer(gl.ARRAY_BUFFER, 0)
	gl.BindVertexArray(0)

	return VAO
}

func createUBO(block PerspectiveBlock) uint32 {
	var ubo uint32
	gl.GenBuffers(1, &ubo)
	gl.BindBuffer(gl.UNIFORM_BUFFER, ubo)
	gl.BufferData(gl.UNIFORM_BUFFER, 3*16*4, nil, gl.STREAM_DRAW)
	gl.BindBufferBase(gl.UNIFORM_BUFFER, 1, ubo)

	gl.BindBuffer(gl.UNIFORM_BUFFER, ubo)
	gl.BufferSubData(gl.UNIFORM_BUFFER, 0, 16*4, gl.Ptr(&block.project[0]))
	gl.BufferSubData(gl.UNIFORM_BUFFER, 16*4, 16*4, gl.Ptr(&block.camera[0]))
	gl.BufferSubData(gl.UNIFORM_BUFFER, 32*4, 16*4, gl.Ptr(&block.world[0]))

	return ubo
}

func drawWindowContent(VAO uint32, program uint32) {
	gl.Clear(gl.COLOR_BUFFER_BIT | gl.DEPTH_BUFFER_BIT) // Clear the color and depth buffer bits
	gl.ClearColor(0.52, 0.80, 0.96, 1.0)

	gl.BindVertexArray(VAO)
	if err := gl.GetError(); err != 0 {
		log.Printf("error binding vertex array object: %v\n", err)
	}

	gl.UseProgram(program)
	if err := gl.GetError(); err != 0 {
		log.Printf("error using program: %v\n", err)
	}

	gl.DrawArrays(gl.TRIANGLES, 0, int32(len(square)/3))
	//gl.DrawElements(gl.TRIANGLES, int32(len(object.Indices)), gl.UNSIGNED_INT, nil)

	gl.Enable(gl.DEPTH_TEST)
}

func (w *GLWindow) Width() int {
	return w.width
}

func (w *GLWindow) Height() int {
	return w.height
}

func (w *GLWindow) ShouldClose() bool {
	return w.glfw.ShouldClose()
}

func (w *GLWindow) OnFrameStart() {
	w.glfw.SwapBuffers()
	glfw.PollEvents()
}

func createWorldMatrix() mgl32.Mat4 {
	translate := mgl32.Vec3{0, 0, -3}
	angle := float32(1 / 4)
	axis := mgl32.Vec3{0, 1, 0}
	scale := mgl32.Vec3{2, 2, 2}

	translateMatrix := mgl32.Translate3D(translate.X(), translate.Y(), translate.Z())
	rotateMatrix := mgl32.HomogRotate3D(angle, axis)
	scaleMatrix := mgl32.Scale3D(scale.X(), scale.Y(), scale.Z())

	return translateMatrix.Mul4(rotateMatrix.Mul4(scaleMatrix))
}

func createPerspectiveBlock(projection *mgl32.Mat4, viewport *mgl32.Mat4, world *mgl32.Mat4) PerspectiveBlock {
	return PerspectiveBlock{
		project: projection,
		camera:  viewport,
		world:   world,
	}
}
