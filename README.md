```
app-name
├─ .VSCodeCounter
│  ├─ 2026-03-20_06-40-16
│  │  ├─ details.md
│  │  ├─ diff-details.md
│  │  ├─ diff.csv
│  │  ├─ diff.md
│  │  ├─ diff.txt
│  │  ├─ results.csv
│  │  ├─ results.json
│  │  ├─ results.md
│  │  └─ results.txt
│  └─ 2026-03-22_16-19-42
│     ├─ details.md
│     ├─ diff-details.md
│     ├─ diff.csv
│     ├─ diff.md
│     ├─ diff.txt
│     ├─ results.csv
│     ├─ results.json
│     ├─ results.md
│     └─ results.txt
├─ Dockerfile
├─ README.md
├─ cmd
│  ├─ api
│  │  └─ main.go
│  └─ bootsrap
│     └─ bootstrap.go
├─ config
│  └─ config.go
├─ docker-compose.yaml
├─ docs
│  └─ openapi.yaml
├─ go.mod
├─ go.sum
├─ internal
│  ├─ app
│  │  ├─ applications
│  │  │  ├─ contract
│  │  │  │  ├─ repository.go
│  │  │  │  └─ service.go
│  │  │  ├─ dto
│  │  │  │  ├─ request.go
│  │  │  │  └─ response.go
│  │  │  ├─ entity
│  │  │  │  └─ applications.go
│  │  │  ├─ handler
│  │  │  │  ├─ handler.go
│  │  │  │  └─ routes.go
│  │  │  ├─ repository
│  │  │  │  └─ repository.go
│  │  │  └─ service
│  │  │     └─ service.go
│  │  ├─ career_mapping
│  │  │  ├─ contract
│  │  │  │  ├─ repository.go
│  │  │  │  └─ service.go
│  │  │  ├─ dto
│  │  │  │  ├─ request.go
│  │  │  │  └─ response.go
│  │  │  ├─ entity
│  │  │  │  └─ career_mapping.go
│  │  │  ├─ handler
│  │  │  │  ├─ handler.go
│  │  │  │  └─ routes.go
│  │  │  ├─ repository
│  │  │  │  └─ repository.go
│  │  │  └─ service
│  │  │     └─ service.go
│  │  ├─ company
│  │  │  ├─ contract
│  │  │  │  ├─ repository.go
│  │  │  │  └─ service.go
│  │  │  ├─ dto
│  │  │  │  └─ response.go
│  │  │  ├─ entity
│  │  │  │  └─ company.go
│  │  │  ├─ handler
│  │  │  │  ├─ handler.go
│  │  │  │  └─ routes.go
│  │  │  ├─ repository
│  │  │  │  └─ repository.go
│  │  │  └─ service
│  │  │     └─ service.go
│  │  ├─ gemini
│  │  │  ├─ contract
│  │  │  │  ├─ repository.go
│  │  │  │  └─ service.go
│  │  │  ├─ dto
│  │  │  │  ├─ request.go
│  │  │  │  └─ response.go
│  │  │  ├─ entity
│  │  │  │  └─ cv.go
│  │  │  ├─ handler
│  │  │  │  ├─ handler.go
│  │  │  │  └─ routes.go
│  │  │  ├─ repository
│  │  │  │  └─ repository.go
│  │  │  └─ service
│  │  │     └─ service.go
│  │  ├─ home
│  │  │  ├─ dto
│  │  │  │  └─ response.go
│  │  │  ├─ handler
│  │  │  │  ├─ handler.go
│  │  │  │  └─ routes.go
│  │  │  ├─ repository
│  │  │  └─ service
│  │  │     └─ service.go
│  │  ├─ job_board
│  │  │  ├─ contract
│  │  │  │  ├─ repository.go
│  │  │  │  └─ service.go
│  │  │  ├─ dto
│  │  │  │  ├─ request.go
│  │  │  │  └─ response.go
│  │  │  ├─ entity
│  │  │  │  ├─ job_listing.go
│  │  │  │  └─ saved_jobs.go
│  │  │  ├─ handler
│  │  │  │  ├─ handler.go
│  │  │  │  └─ routes.go
│  │  │  ├─ repository
│  │  │  │  └─ repository.go
│  │  │  └─ service
│  │  │     └─ service.go
│  │  ├─ onboarding
│  │  │  ├─ contract
│  │  │  │  ├─ repository.go
│  │  │  │  └─ service.go
│  │  │  ├─ dto
│  │  │  │  ├─ request.go
│  │  │  │  └─ response.go
│  │  │  ├─ handler
│  │  │  │  ├─ handler.go
│  │  │  │  └─ routes.go
│  │  │  ├─ repository
│  │  │  │  └─ repository.go
│  │  │  └─ service
│  │  │     └─ service.go
│  │  ├─ payment
│  │  │  ├─ contract
│  │  │  │  ├─ repository.go
│  │  │  │  └─ service.go
│  │  │  ├─ dto
│  │  │  │  ├─ request.go
│  │  │  │  └─ response.go
│  │  │  ├─ entity
│  │  │  │  └─ order.go
│  │  │  ├─ handler
│  │  │  │  ├─ handler.go
│  │  │  │  └─ routes.go
│  │  │  ├─ repository
│  │  │  │  └─ repository.go
│  │  │  └─ service
│  │  │     └─ service.go
│  │  ├─ smart_profile
│  │  │  ├─ contract
│  │  │  │  └─ service.go
│  │  │  ├─ dto
│  │  │  │  └─ response.go
│  │  │  ├─ handler
│  │  │  │  ├─ handler.go
│  │  │  │  └─ routes.go
│  │  │  └─ service
│  │  │     └─ service.go
│  │  └─ user
│  │     ├─ contract
│  │     │  ├─ repository.go
│  │     │  └─ service.go
│  │     ├─ dto
│  │     │  ├─ request.go
│  │     │  └─ response.go
│  │     ├─ entity
│  │     │  ├─ refresh_token.go
│  │     │  ├─ user.go
│  │     │  └─ verification_token.go
│  │     ├─ handler
│  │     │  ├─ auth_handler.go
│  │     │  ├─ routes.go
│  │     │  └─ user_handler.go
│  │     ├─ repository
│  │     │  ├─ refresh_token_repositroy.go
│  │     │  ├─ repository.go
│  │     │  └─ verification_token_repository.go
│  │     └─ service
│  │        └─ service.go
│  ├─ infra
│  │  └─ database
│  │     ├─ connection.go
│  │     ├─ migration.go
│  │     └─ seed.go
│  └─ middleware
│     └─ jwt.go
└─ pkg
   ├─ email
   │  ├─ email.go
   │  └─ template.go
   ├─ gemini
   │  └─ gemini.go
   ├─ jwt
   │  └─ jwt.go
   ├─ response
   │  ├─ error.go
   │  └─ response.go
   ├─ storage
   │  └─ storage.go
   └─ xendit
      └─ xendit.go

```
