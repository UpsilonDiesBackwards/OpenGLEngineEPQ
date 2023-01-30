package primitives

import (
	"fmt"
	"os"
	"strconv"
	"strings"
)

var vertices float32
var indices uint32

func CreateNewOBJ(filePath string) (vertices []float32, indices []uint32) {
	// TODO: Make this it's own function

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
			vertices = append(vertices, float32(x), float32(y), float32(z))
		case "vt":
			fmt.Println("Texture coordinate found:", line)
		case "vn":
			fmt.Println("Vertex normal found:", line)
		case "f":
			for _, field := range fields[1:] {
				vertexIndices := strings.Split(field, "/")
				vertexIndex, _ := strconv.Atoi(vertexIndices[0])
				indices = append(indices, uint32(vertexIndex))

				//// Use vertexIndex here
				//textureIndex, _ := strconv.Atoi(vertexIndices[1])
				//// Use textureIndex here
				//normalIndex, _ := strconv.Atoi(vertexIndices[2])
				//// Use normalIndex here
			}
		}
	}

	return vertices, indices
}
