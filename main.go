package main

import (
	"context"
	"errors"
	"fmt"
	"google.golang.org/grpc"
	"log"
	"math"
	"net"
	"sort"

	pb "github.com/Frozelo/recServcie/proto/gen/go"
)

var usersLikes = map[int64][]float64{
	1: {0, 1, 0, 1, 1, 1, 0, 1, 0, 0, 0, 0, 1, 0, 1, 0, 1, 1, 1, 0, 0, 1, 0, 0, 1, 1, 0, 1, 1, 0},
	2: {0, 1, 0, 1, 0, 1, 1, 1, 0, 1, 0, 1, 0, 1, 0, 1, 1, 1, 1, 0, 0, 1, 0, 0, 1, 0, 1, 1, 1, 1},
	3: {1, 1, 0, 1, 0, 1, 1, 1, 1, 0, 0, 1, 1, 1, 1, 1, 1, 0, 0, 0, 1, 1, 0, 1, 0, 0, 1, 1, 0, 1},
	4: {0, 1, 0, 1, 1, 1, 0, 1, 0, 1, 1, 1, 1, 0, 1, 0, 1, 1, 1, 1, 0, 1, 0, 1, 1, 1, 1, 0, 1, 1},
}

type server struct {
	pb.UnimplementedRecommendationServiceServer
}

func (s *server) GetRecommendations(ctx context.Context, req *pb.RecommendationRequest) (*pb.RecommendationResponse, error) {
	userId := req.GetUserId()
	minSimilarity := req.GetMinSimilarity()
	topN := req.GetTopN()

	topSimilarUsers, recommendations, err := ComputeRecommendations(userId, usersLikes, float64(minSimilarity), int(topN))
	if err != nil {
		log.Fatalf("Error computing recommendations: %v", err)
		return nil, err
	}

	log.Println(convertToFloat32Map(recommendations))

	return &pb.RecommendationResponse{
		SimilarUsers:    topSimilarUsers,
		Recommendations: convertToFloat32Map(recommendations),
	}, nil
}

func main() {
	grpcServer := grpc.NewServer()
	lis, err := net.Listen("tcp", ":8081")
	if err != nil {
		log.Fatalf("failed to listen: %v", err)
	}

	// Регистрируем наш gRPC сервис
	pb.RegisterRecommendationServiceServer(grpcServer, &server{})

	log.Println("Starting gRPC server on :8081...")
	if err := grpcServer.Serve(lis); err != nil {
		log.Fatalf("failed to serve: %v", err)
	}
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
	userRecs := getFilteredRecommendations(userLikes, usersLike)

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

func getFilteredRecommendations(userLikes []float64, usersLike map[int64][]float64) map[int64]float64 {
	recommendations := make(map[int64]float64)
	for otherUserID, likes := range usersLike {
		for i := 0; i < len(userLikes); i++ {
			if userLikes[i] == 0 && likes[i] > 0 {
				recommendations[otherUserID] = 1
			}
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

func convertToFloat32Map(input map[int64]float64) map[int64]float32 {
	output := make(map[int64]float32)
	for key, value := range input {
		output[key] = float32(value)
	}
	return output
}
