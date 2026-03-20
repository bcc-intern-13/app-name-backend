package service

import (
	"encoding/json"
	"log/slog"
	"sort"

	"github.com/bcc-intern-13/app-name-backend/internal/app/career_mapping/contract"
	"github.com/bcc-intern-13/app-name-backend/internal/app/career_mapping/dto"
	"github.com/bcc-intern-13/app-name-backend/internal/app/career_mapping/entity"
	"github.com/bcc-intern-13/app-name-backend/pkg/response"
	"github.com/google/uuid"
	"gorm.io/datatypes"
)

type careerMappingService struct {
	repo contract.CareerMappingRepository
}

func NewCareerMappingService(repo contract.CareerMappingRepository) contract.CareerMappingService {
	return &careerMappingService{repo: repo}
}

func (s *careerMappingService) GetQuestions() ([]entity.CareerMappingQuestion, *response.APIError) {
	questions, err := s.repo.GetAllQuestions()
	if err != nil {
		slog.Error("failed to get questions", "error", err)
		return nil, response.ErrInternal("failed to get questions")
	}
	return questions, nil
}

func (s *careerMappingService) Submit(userID uuid.UUID, req *dto.SubmitCareerMappingRequest) (*dto.CareerMappingResponse, *response.APIError) {
	questions, err := s.repo.GetAllQuestions()
	if err != nil {
		slog.Error("failed to get questions", "error", err)
		return nil, response.ErrInternal("failed to get questions")
	}

	scores := map[string]int{
		"KR": 0, "TK": 0, "KO": 0,
		"ED": 0, "AD": 0, "OP": 0,
	}

	for i, answer := range req.Answers {
		q := questions[i]
		var skorJSON []byte
		switch answer {
		case "A":
			skorJSON = q.ScoreA
		case "B":
			skorJSON = q.ScoreB
		case "C":
			skorJSON = q.ScoreC
		case "D":
			skorJSON = q.ScoreD
		}

		var skor map[string]int
		if err := json.Unmarshal(skorJSON, &skor); err != nil {
			slog.Error("failed to parse score", "error", err, "answer", answer)
			return nil, response.ErrInternal("failed to parse score")
		}
		for cat, pts := range skor {
			scores[cat] += pts
		}
	}

	type catScore struct {
		Code  string
		Score int
	}
	priority := []string{"KR", "TK", "KO", "ED", "AD", "OP"}

	var sorted []catScore
	for _, code := range priority {
		sorted = append(sorted, catScore{Code: code, Score: scores[code]})
	}
	sort.SliceStable(sorted, func(i, j int) bool {
		if sorted[i].Score != sorted[j].Score {
			return sorted[i].Score > sorted[j].Score
		}
		pi, pj := 0, 0
		for idx, p := range priority {
			if p == sorted[i].Code {
				pi = idx
			}
			if p == sorted[j].Code {
				pj = idx
			}
		}
		return pi < pj
	})
	top3 := sorted[:3]

	var topCategories []dto.CategoryScore
	for rank, cat := range top3 {
		category, err := s.repo.GetCategoryByID(cat.Code)
		if err != nil {
			slog.Error("failed to get category", "error", err, "code", cat.Code)
			continue
		}
		if category == nil {
			continue
		}
		topCategories = append(topCategories, dto.CategoryScore{
			Rank:        rank + 1,
			Code:        cat.Code,
			Name:        category.Name,
			Score:       cat.Score,
			Description: category.Description,
			FormalJobs:  category.FormalJobs,
			SideJobs:    category.SideJobs,
		})
	}

	count, err := s.repo.CountByUserID(userID)
	if err != nil {
		slog.Error("failed to count attempts", "error", err, "userID", userID)
		return nil, response.ErrInternal("failed to count attempts")
	}
	attemptNumber := int(count) + 1

	answersJSON, _ := json.Marshal(req.Answers)
	scoresJSON, _ := json.Marshal(scores)
	topCatCodes := []string{top3[0].Code, top3[1].Code, top3[2].Code}
	topCatJSON, _ := json.Marshal(topCatCodes)

	result := &entity.CareerMappingResult{
		ID:            uuid.New(),
		UserID:        userID,
		Answers:       datatypes.JSON(answersJSON),
		Scores:        datatypes.JSON(scoresJSON),
		TopCategories: datatypes.JSON(topCatJSON),
		AttemptNumber: attemptNumber,
	}

	if err := s.repo.CreateResult(result); err != nil {
		slog.Error("failed to save result", "error", err, "userID", userID)
		return nil, response.ErrInternal("failed to save result")
	}

	var allScores []dto.CategoryScore
	for _, cat := range sorted {
		category, _ := s.repo.GetCategoryByID(cat.Code)
		if category != nil {
			allScores = append(allScores, dto.CategoryScore{
				Code:  cat.Code,
				Name:  category.Name,
				Score: cat.Score,
			})
		}
	}

	return &dto.CareerMappingResponse{
		TopCategories: topCategories,
		AllScores:     allScores,
		AttemptNumber: attemptNumber,
	}, nil
}

func (s *careerMappingService) GetLatestResult(userID uuid.UUID) (*dto.CareerMappingResponse, *response.APIError) {
	result, err := s.repo.FindLatestResultByUserID(userID)
	if err != nil {
		slog.Error("failed to get latest result", "error", err, "userID", userID)
		return nil, response.ErrInternal("failed to get latest result")
	}
	if result == nil {
		return nil, response.ErrNotFound("result not found")
	}

	var topCatCodes []string
	json.Unmarshal(result.TopCategories, &topCatCodes)

	var topCategories []dto.CategoryScore
	for rank, code := range topCatCodes {
		category, _ := s.repo.GetCategoryByID(code)
		var scores map[string]int
		json.Unmarshal(result.Scores, &scores)
		if category != nil {
			topCategories = append(topCategories, dto.CategoryScore{
				Rank:        rank + 1,
				Code:        code,
				Name:        category.Name,
				Score:       scores[code],
				Description: category.Description,
				FormalJobs:  category.FormalJobs,
				SideJobs:    category.SideJobs,
			})
		}
	}

	return &dto.CareerMappingResponse{
		TopCategories: topCategories,
		AttemptNumber: result.AttemptNumber,
	}, nil
}
func (s *careerMappingService) GetHistory(userID uuid.UUID) ([]dto.CareerMappingResponse, *response.APIError) {
	results, err := s.repo.FindAllResultsByUserID(userID)
	if err != nil {
		slog.Error("failed to get career mapping history", "error", err, "userID", userID)
		return nil, response.ErrInternal("failed to get history")
	}
	if len(results) == 0 {
		return []dto.CareerMappingResponse{}, nil
	}

	var responses []dto.CareerMappingResponse
	for _, result := range results {
		var topCatCodes []string
		if err := json.Unmarshal(result.TopCategories, &topCatCodes); err != nil {
			slog.Error("failed to parse top categories", "error", err)
			continue
		}

		var topCategories []dto.CategoryScore
		for rank, code := range topCatCodes {
			category, _ := s.repo.GetCategoryByID(code)
			var scores map[string]int
			json.Unmarshal(result.Scores, &scores)
			if category != nil {
				topCategories = append(topCategories, dto.CategoryScore{
					Rank:        rank + 1,
					Code:        code,
					Name:        category.Name,
					Score:       scores[code],
					Description: category.Description,
					FormalJobs:  category.FormalJobs,
					SideJobs:    category.SideJobs,
				})
			}
		}

		responses = append(responses, dto.CareerMappingResponse{
			TopCategories: topCategories,
			AttemptNumber: result.AttemptNumber,
			CreatedAt:     result.CreatedAt,
		})
	}

	return responses, nil
}
