package objLoader

import (
	"fmt"
	"github.com/go-gl/gl/all-core/gl"
	"github.com/go-gl/mathgl/mgl32"
	"os"
	"strconv"
	"strings"
)

var vertices float32
var normals []float32
var textureCoords []float32
var indices uint32

// Object
type Object struct {
	UUID string // UUID does nothing atm, will be used if I implement imgui or an alternative.

	ModelMatrix mgl32.Mat4
	Translation mgl32.Vec3
	Rotation    mgl32.Vec3
	Scale       mgl32.Vec3

	VAO           uint32
	Vertices      []float32
	Indices       []uint32
	Normals       []float32
	TextureCoords []float32
}

var object = &Object{}
var Objects []*Object // Arrary of objects in scene

func CreateObject(filepath, uuid string, translation, rotation, scale mgl32.Vec3) *Object {
	object.Vertices, object.Normals, object.TextureCoords, object.Indices = loadOBJFromFile("resources/obj/" + filepath)

	vao := createVAO(object.Vertices, object.Indices)

	obj := &Object{
		UUID:        uuid,
		Translation: translation,
		Rotation:    rotation,
		Scale:       scale,

		VAO:           vao,
		Vertices:      object.Vertices,
		Indices:       object.Indices,
		Normals:       object.Normals,
		TextureCoords: object.TextureCoords,
	}
	Objects = append(Objects, obj) // add object to object array

	return obj
}

func loadOBJFromFile(filePath string) (vertices, normals, textureCoords []float32, indices []uint32) {
	// Read file
	fileBytes, err := os.ReadFile(filePath)
	if err != nil {
		fmt.Printf("Failed to load OBJ file!")
		return
	}

	lines := strings.Split(string(fileBytes), "\n")
	for _, line := range lines {
		if len(line) < 2 {
			continue
		}

		fields := strings.Fields(line)
		switch fields[0] {
		case "v": // process vertex data
			x, _ := strconv.ParseFloat(fields[1], 32)
			y, _ := strconv.ParseFloat(fields[2], 32)
			z, _ := strconv.ParseFloat(fields[3], 32)
			vertices = append(triangulate(vertices), float32(x), float32(-y), float32(z))
		case "vt": // process texture coord data
			x, _ := strconv.ParseFloat(fields[1], 32)
			y, _ := strconv.ParseFloat(fields[2], 32)
			textureCoords = append(textureCoords, float32(x), float32(y))
		case "vn": // process vertex normal data
			x, _ := strconv.ParseFloat(fields[1], 32)
			y, _ := strconv.ParseFloat(fields[2], 32)
			z, _ := strconv.ParseFloat(fields[3], 32)
			normals = append(normals, float32(x), float32(y), float32(z))
		case "f": // process index
			for _, field := range fields[1:] {
				vertexIndices := strings.Split(field, "/")
				vertexIndex, _ := strconv.Atoi(vertexIndices[0])
				indices = append(indices, uint32(vertexIndex-1))
			}
		}
	}

	return vertices, normals, textureCoords, indices
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

	gl.VertexAttribPointer(0, 3, gl.FLOAT, false, 0, gl.PtrOffset(0))
	gl.EnableVertexAttribArray(0)

	gl.BindBuffer(gl.ARRAY_BUFFER, 0)
	gl.BindVertexArray(0)

	return VAO
}

func (obj *Object) CreateModelMatrix() {
	translation := obj.Translation
	rotation := obj.Rotation
	scale := obj.Scale

	translateMatrix := mgl32.Translate3D(translation.X(), translation.Y(), translation.Z())
	var rotateMatrix mgl32.Mat4
	if rotation.Len() > 0 {
		rotateMatrix = mgl32.HomogRotate3D(rotation.Len(), rotation.Normalize())
	} else {
		rotateMatrix = mgl32.Ident4()
	}
	scaleMatrix := mgl32.Scale3D(scale.X(), scale.Y(), scale.Z())

	obj.ModelMatrix = translateMatrix.Mul4(rotateMatrix.Mul4(scaleMatrix))
}

func triangulate(vertices []float32) []float32 {
	// Triangulate *.obj mesh. I have no idea if it works or not, but it gets angry if I don't use it. :O
	triangulatedVertices := make([]float32, 0)
	for i := 0; i < len(vertices); i += 3 {
		triangulatedVertices = append(triangulatedVertices, vertices[i], vertices[i+1], vertices[i+2])
	}
	return triangulatedVertices
}
