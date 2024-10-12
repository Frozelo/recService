package main

import (
	"errors"
	"fmt"
	"log"
	"math"
)

var usersLikes = map[int64][]float64{
	1: {0, 1, 0, 1, 1, 1, 0, 1, 0, 0, 0, 0},
	2: {0, 1, 0, 1, 0, 1, 1, 1, 0, 1, 0, 1},
	3: {0, 1, 0, 1, 0, 1, 1, 1, 0, 1, 0, 1},
}

func main() {
	user1 := []float64{0, 1, 0, 1, 1, 1, 0, 1, 0, 1, 1, 1}
	fmt.Println(ComputeForUser(user1, usersLikes))
}

func ComputeForUser(userLikes []float64, usersLike map[int64][]float64) float64 {
	var err error
	var similarity, maxSimilarity float64
	var similarUserId int64
	for userId, likes := range usersLike {
		similarity, err = cosineSimilarity(userLikes, likes)
		if err != nil {
			log.Fatal(err)
		}
		if similarity > maxSimilarity {
			maxSimilarity = similarity
			similarUserId = userId
		}
	}
	log.Fatalf("the maximum similarity is %f\nand the similar user is %v", maxSimilarity, similarUserId)
	return maxSimilarity
}

func cosineSimilarity(v1 []float64, v2 []float64) (float64, error) {
	if len(v1) != len(v2) {
		return 0.0, errors.New("the length of v1 and v2 are not equal")
	}
	var dotProduct, normV1, normV2 float64
	for i := 0; i < len(v1); i++ {
		dotProduct += v1[i] * v2[i]
		normV1 += math.Pow(v1[i], 2)
		normV2 += math.Pow(v2[i], 2)
	}
	if normV1 == 0 || normV2 == 0 {
		return 0.0, errors.New("the sum of all coefficients is zero")
	}

	return dotProduct / (math.Sqrt(normV1) * math.Sqrt(normV2)), nil
}
