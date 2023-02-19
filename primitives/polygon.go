package primitives

import (
	"fmt"
	"github.com/UpsilonDiesBackwards/behngine_epq/camera"
	"github.com/go-gl/gl/all-core/gl"
	"github.com/go-gl/glfw/v3.2/glfw"
	"github.com/go-gl/mathgl/mgl32"
	"log"
	"os"
	"strconv"
	"strings"
)

var vertices []float32
var normals []float32
var textureCoords []float32
var indices []uint32

// Object
type ObjectPrimitive struct {
	Vertices      []float32
	Normals       []float32
	TextureCoords []float32
	Indices       []uint32
}

type Primitive struct {
	*ObjectPrimitive

	VAO uint32
	VBO uint32
	EBO uint32
	UBO uint32

	ShaderProgram uint32
}

type PerspectiveBlock struct {
	project *mgl32.Mat4
	camera  *mgl32.Mat4
	world   *mgl32.Mat4
}

var ViewportTransform mgl32.Mat4

var object *Primitive
var UBO uint32

var fov = float32(60.0)
var projectionTransform = mgl32.Perspective(mgl32.DegToRad(fov),
	float32(800/600), 0.1, 2000)

func CreateOBJObject(filepath, ObjectID string, appWindow *glfw.Window, c *camera.Camera) {
	object = &Primitive{}
	objPrim := <-loadOBJFromFile(filepath)
	object.ObjectPrimitive = &objPrim

	world := createWorldMatrix() // Create world matrix

	VAO := createVAO(object.ObjectPrimitive.Vertices, object.ObjectPrimitive.Indices, object.VAO, object.VBO, object.EBO)
	if VAO == 0 {
		log.Fatalln("Error creating VAO for object: ", ObjectID)
	}

	for !appWindow.ShouldClose() {
		ViewportTransform = c.GetTransform()

		block := createPerspectiveBlock(&projectionTransform, &ViewportTransform, &world)
		UBO = createUBO(block)

		fmt.Println(object.Vertices)
		gl.DrawArrays(gl.TRIANGLES, 0, int32(len(object.Vertices)/3))
	}

}

func loadOBJFromFile(filePath string) <-chan ObjectPrimitive {
	result := make(chan ObjectPrimitive)

	go func() {
		// Read file
		fileBytes, err := os.ReadFile(filePath)
		if err != nil {
			fmt.Printf("Failed to load OBJ file: %s", filePath)
		}

		lines := strings.Split(string(fileBytes), "\n")
		for _, line := range lines {
			if len(line) < 2 {
				continue
			}

			fields := strings.Fields(line)
			switch fields[0] {
			case "v":
				x, _ := strconv.ParseFloat(fields[1], 32)
				y, _ := strconv.ParseFloat(fields[2], 32)
				z, _ := strconv.ParseFloat(fields[3], 32)
				vertices = append(vertices, float32(x), float32(y), float32(z))
			case "vt":
				x, _ := strconv.ParseFloat(fields[1], 32)
				y, _ := strconv.ParseFloat(fields[2], 32)
				textureCoords = append(textureCoords, float32(x), float32(y))
			case "vn":
				x, _ := strconv.ParseFloat(fields[1], 32)
				y, _ := strconv.ParseFloat(fields[2], 32)
				z, _ := strconv.ParseFloat(fields[3], 32)
				normals = append(normals, float32(x), float32(y), float32(z))
			case "f":
				for _, field := range fields[1:] {
					vertexIndices := strings.Split(field, "/")
					if len(vertexIndices) > 0 {
						vertexIndex, _ := strconv.Atoi(vertexIndices[0])
						indices = append(indices, uint32(vertexIndex-1))
					}
				}
			}
		}

		result <- ObjectPrimitive{
			Vertices:      vertices,
			Indices:       indices,
			Normals:       normals,
			TextureCoords: textureCoords,
		}
		close(result)
	}()
	return result
}

func triangulate(vertices []float32) []float32 {
	triangulatedVertices := make([]float32, 0)
	for i := 0; i < len(vertices); i += 3 {
		triangulatedVertices = append(triangulatedVertices, vertices[i], vertices[i+1], vertices[i+2])
	}
	return triangulatedVertices
}

func createVAO(vertices []float32, indices []uint32, VAO, VBO, EBO uint32) uint32 {
	gl.GenVertexArrays(1, &VAO)
	gl.BindVertexArray(VAO)

	gl.GenBuffers(1, &VBO)
	gl.BindBuffer(gl.ARRAY_BUFFER, VBO)
	gl.BufferData(gl.ARRAY_BUFFER, len(vertices)*4, gl.Ptr(vertices), gl.STATIC_DRAW)

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

func Render(VAO uint32) {
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
