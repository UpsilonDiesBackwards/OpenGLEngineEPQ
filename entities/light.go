package entities

import (
	"github.com/go-gl/gl/all-core/gl"
	"github.com/go-gl/mathgl/mgl32"
	"math"
)

type Light struct {
	VAO uint32

	ModelMatrix mgl32.Mat4
	Position    mgl32.Vec3
	LightColor  mgl32.Vec3

	Vertices []float32
	Indices  []uint32
}

var light = &Light{}
var Lights []*Light

func createSphere(radius float32, slices int, stacks int) ([]float32, []uint32) {
	var vertices []float32
	var indices []uint32

	// Generate vertices
	for i := 0; i <= stacks; i++ {
		phi := math.Pi/2 - float64(i)*math.Pi/float64(stacks)
		for j := 0; j <= slices; j++ {
			theta := float64(j) * 2 * math.Pi / float64(slices)

			x := radius * float32(math.Cos(theta)*math.Sin(phi))
			y := radius * float32(math.Cos(phi))
			z := radius * float32(math.Sin(theta)*math.Sin(phi))

			vertices = append(vertices, x, y, z)
		}
	}

	// Generate indices
	for i := 0; i < stacks; i++ {
		for j := 0; j < slices; j++ {
			index := uint32(i*(slices+1) + j)
			indices = append(indices, index, index+uint32(slices)+1, uint32(index)+1)
			indices = append(indices, uint32(index)+1, index+uint32(slices)+1, index+uint32(slices)+2)
		}
	}

	return vertices, indices
}

func CreateLight(position, lColor mgl32.Vec3) *Light {
	lightVert, lightInd := createSphere(0.5, 6, 6)

	vao := createLightVAO(lightVert, lightInd)

	light := &Light{
		VAO: vao,

		Position:   position,
		LightColor: lColor,

		Vertices: lightVert,
		Indices:  lightInd,
	}
	Lights = append(Lights, light)

	return light
}

func createLightVAO(vertices []float32, indices []uint32) uint32 {
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

	gl.VertexAttribPointer(0, 3, gl.FLOAT, false, 0, gl.PtrOffset(0))
	gl.EnableVertexAttribArray(0)

	gl.BindBuffer(gl.ARRAY_BUFFER, 0)
	gl.BindVertexArray(0)

	return VAO
}

func (l *Light) CreateLightModelMatrix() {
	position := l.Position

	translateMatrix := mgl32.Translate3D(position.X(), position.Y(), position.Z())
	scaleMatrix := mgl32.Ident4()
	rotationMatrix := mgl32.Ident4()

	l.ModelMatrix = translateMatrix.Mul4(scaleMatrix).Mul4(rotationMatrix)
}
