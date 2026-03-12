package service

import (
	"encoding/json"
	"errors"
	"sort"

	"github.com/bcc-intern-13/app-name-backend/internal/career_mapping/dto"
	"github.com/bcc-intern-13/app-name-backend/internal/career_mapping/entity"
	"github.com/google/uuid"
	"gorm.io/datatypes"
)

type careerMappingService struct {
	repo dto.CareerMappingRepository
}

func NewCareerMappingService(repo dto.CareerMappingRepository) dto.CareerMappingService {
	return &careerMappingService{repo: repo}
}

func (s *careerMappingService) GetQuestions() ([]entity.CareerMappingQuestion, error) {
	return s.repo.GetAllQuestions()
}

func (s *careerMappingService) Submit(userID uuid.UUID, req *dto.SubmitCareerMappingRequest) (*dto.CareerMappingResponse, error) {
	questions, err := s.repo.GetAllQuestions()
	if err != nil {
		return nil, errors.New("failed to get questions")
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
			skorJSON = q.SkorA
		case "B":
			skorJSON = q.SkorB
		case "C":
			skorJSON = q.SkorC
		case "D":
			skorJSON = q.SkorD
		}

		var skor map[string]int
		if err := json.Unmarshal(skorJSON, &skor); err != nil {
			return nil, errors.New("failed to parse score")
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
		// tiebreaker: priority order
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
		if err != nil || category == nil {
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
		return nil, errors.New("failed to count attempts")
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
		return nil, errors.New("failed to save result")
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

func (s *careerMappingService) GetLatestResult(userID uuid.UUID) (*dto.CareerMappingResponse, error) {
	result, err := s.repo.FindLatestResultByUserID(userID)
	if err != nil {
		return nil, err
	}
	if result == nil {
		return nil, errors.New("result not found")
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
