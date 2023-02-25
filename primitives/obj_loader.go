package primitives

import (
	"fmt"
	"os"
	"strconv"
	"strings"
)

var vertices float32
var normals []float32
var textureCoords []float32
var indices uint32

// Object
type ObjectPrimitive struct {
	ObjectV  []float32
	ObjectI  []uint32
	ObjectN  []float32
	ObjectTN []float32
}

var object = &ObjectPrimitive{}

func CreateObject(filepath string, object *ObjectPrimitive) {
	object.ObjectV, object.ObjectN, object.ObjectTN, object.ObjectI = loadOBJFromFile("resources/obj/" + filepath)
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
		case "v":
			x, _ := strconv.ParseFloat(fields[1], 32)
			y, _ := strconv.ParseFloat(fields[2], 32)
			z, _ := strconv.ParseFloat(fields[3], 32)
			vertices = append((triangulate(vertices)), float32(x), float32(-y), float32(z))
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
				vertexIndex, _ := strconv.Atoi(vertexIndices[0])
				indices = append(indices, uint32(vertexIndex-1))
			}
		}
	}

	return vertices, normals, textureCoords, indices
}

func triangulate(vertices []float32) []float32 {
	triangulatedVertices := make([]float32, 0)
	for i := 0; i < len(vertices); i += 3 {
		triangulatedVertices = append(triangulatedVertices, vertices[i], vertices[i+1], vertices[i+2])
	}
	return triangulatedVertices
}
