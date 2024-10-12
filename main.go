package main

import (
	"fmt"
	"log"
	"math"
)

func main() {
	user1 := []float64{0, 1, 0, 1, 1, 1, 0, 1, 0, 1, 1, 1}
	user2 := []float64{1, 0, 0, 1, 1, 1, 1, 1, 1, 1, 1, 0}
	fmt.Println("the similarity is:", cosineSimilarity(user1, user2))
}

func cosineSimilarity(v1 []float64, v2 []float64) float64 {
	if len(v1) != len(v2) {
		log.Println("v1 and v2 are not the same length")
		return 0.0
	}
	var dotProduct, normV1, normV2 float64
	for i := 0; i < len(v1); i++ {
		dotProduct += v1[i] * v2[i]
		normV1 += math.Pow(v1[i], 2)
		normV2 += math.Pow(v2[i], 2)
	}
	if normV1 == 0 || normV2 == 0 {
		return 0.0
	}

	return dotProduct / (math.Sqrt(normV1) * math.Sqrt(normV2))
}
