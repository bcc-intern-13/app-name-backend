# WorkAble Backend API

> "Mencari kerja itu seperti mencari jodoh. Kadang butuh algoritma, kadang butuh keajaiban. Kami menyediakan algoritmanya."

WorkAble Backend adalah sistem di balik layar yang menggerakkan platform WorkAble. Dibangun dengan arsitektur yang menjunjung tinggi modularitas (karena hidup udah cukup berantakan, kodingan jangan), API ini menangani segalanya: mulai dari autentikasi _user_, pencocokan karir menggunakan AI, pengelolaan loker, sampai urusan duit via _payment gateway_.

## Tech Stack & Tools

Ibarat bangun rumah, ini bahan bangunan yang kita pakai:

- **Bahasa & Framework:** Golang 🐹 + Go Fiber (Cepat, ringan, dan nggak rewel).
- **Database:** PostgreSQL (Buat nyimpen data yang butuh komitmen jangka panjang).
- **Caching:** Redis (Si etalase cepat yang sempat bikin drama _connection refused_ jam 3 pagi).
- **AI Engine:** Google Gemini API (Buat ngasih rekomendasi karir dan _screening_ CV).
- **Payment Gateway:** Xendit (Karena ngurus pembayaran manual itu capek).
- **Infrastructure:** Docker & Docker Compose.

---

## 📁 Struktur Anatomi (Folder Structure)

Proyek ini menggunakan pendekatan **Domain-Driven Design (DDD)** yang dimodifikasi. Biar gampang nyari letak _bug_ kalau lagi _error_:

```text
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
