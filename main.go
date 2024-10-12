package main

import (
	"errors"
	"fmt"
	"log"
	"math"
	"sort"
)

var usersLikes = map[int64][]float64{
	1: {0, 1, 0, 1, 1, 1, 0, 1, 0, 0, 0, 0, 1, 0, 1, 0, 1, 1, 1, 0, 0, 1, 0, 0, 1, 1, 0, 1, 1, 0},
	2: {0, 1, 0, 1, 0, 1, 1, 1, 0, 1, 0, 1, 0, 1, 0, 1, 1, 1, 1, 0, 0, 1, 0, 0, 1, 0, 1, 1, 1, 1},
	3: {1, 1, 0, 1, 0, 1, 1, 1, 1, 0, 0, 1, 1, 1, 1, 1, 1, 0, 0, 0, 1, 1, 0, 1, 0, 0, 1, 1, 0, 1},
	4: {0, 1, 0, 1, 1, 1, 0, 1, 0, 1, 1, 1, 1, 0, 1, 0, 1, 1, 1, 1, 0, 1, 0, 1, 1, 1, 1, 0, 1, 1},
}

func main() {
	userID := int64(4)
	topSimilarUsers, recommendations, err := ComputeRecommendations(userID, usersLikes, 0.5, 2)
	if err != nil {
		log.Fatalf("Error computing recommendations: %v", err)
	}
	fmt.Printf("Top similar users: %v\n", topSimilarUsers)
	fmt.Printf("Recommendations for user %d: %v\n", userID, recommendations)
}

func ComputeRecommendations(userID int64, usersLike map[int64][]float64, minSimilarity float64, topN int) ([]int64, map[int64]float64, error) {
	userLikes, exists := usersLike[userID]
	if !exists {
		return nil, nil, fmt.Errorf("user %d not found", userID)
	}

	similarities := make(map[int64]float64)
	for otherUserID, likes := range usersLike {
		if otherUserID == userID {
			continue
		}
		similarity, err := cosineSimilarity(userLikes, likes)
		if err != nil {
			return nil, nil, fmt.Errorf("error calculating similarity for user %d: %v", otherUserID, err)
		}
		if similarity >= minSimilarity {
			similarities[otherUserID] = similarity
		}
	}

	topUsers := getTopSimilarUsers(similarities, topN)
	var userRecs map[int64]float64
	for _, similarUserID := range topUsers {
		userRecs = getFilteredRecommendations(userLikes, usersLike[similarUserID])
	}

	return topUsers, userRecs, nil
}

func cosineSimilarity(v1, v2 []float64) (float64, error) {
	if len(v1) != len(v2) {
		return 0.0, errors.New("vectors must have the same length")
	}

	var dotProduct, normV1, normV2 float64
	for i := 0; i < len(v1); i++ {
		dotProduct += v1[i] * v2[i]
		normV1 += v1[i] * v1[i]
		normV2 += v2[i] * v2[i]
	}

	if normV1 == 0 || normV2 == 0 {
		return 0.0, errors.New("one of the vectors is zero")
	}

	return dotProduct / (math.Sqrt(normV1) * math.Sqrt(normV2)), nil
}

func getFilteredRecommendations(userLikes, similarUserLikes []float64) map[int64]float64 {
	recommendations := make(map[int64]float64)
	for i := 0; i < len(userLikes); i++ {
		if userLikes[i] == 0 && similarUserLikes[i] > 0 {
			recommendations[int64(i)] = 1
		}
	}
	return recommendations
}

func getTopSimilarUsers(similarities map[int64]float64, topN int) []int64 {
	type similarityPair struct {
		userID     int64
		similarity float64
	}

	var pairs []similarityPair
	for userID, similarity := range similarities {
		pairs = append(pairs, similarityPair{userID, similarity})
	}

	sort.Slice(pairs, func(i, j int) bool {
		return pairs[i].similarity > pairs[j].similarity
	})

	var topUsers []int64
	for i := 0; i < len(pairs) && i < topN; i++ {
		topUsers = append(topUsers, pairs[i].userID)
	}

	return topUsers
}
