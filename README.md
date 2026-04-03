# WorkAble Backend API

WorkAble Backend API adalah _core system_ yang menggerakkan platform WorkAble. Proyek ini dibangun menggunakan arsitektur berbasis **Domain-Driven Design (DDD)** yang dimodifikasi untuk memastikan modularitas, skalabilitas, dan kemudahan pemeliharaan (_maintainability_). API ini menangani seluruh proses bisnis utama, mulai dari autentikasi pengguna, pemrosesan profil cerdas menggunakan AI, manajemen lowongan pekerjaan, hingga integrasi gerbang pembayaran (_payment gateway_).

## 🛠️ Tech Stack & Tools

Proyek ini dikembangkan menggunakan teknologi berikut:

- **Language & Framework:** Golang + Go Fiber (Dipilih untuk performa tinggi dan _concurrency_).
- **Database:** PostgreSQL (Relational Database Management System utama).
- **Caching:** Redis (Diimplementasikan untuk optimasi _response time_ pada _endpoint_ dengan beban baca tinggi, seperti _Home Summary_).
- **AI Engine:** Google Gemini API (Digunakan untuk _Career Mapping_ dan analisis profil pengguna).
- **Payment Gateway:** Xendit (Untuk pemrosesan transaksi yang aman dan otomatis).
- **Infrastructure:** Docker & Docker Compose (Untuk konsistensi _environment_ pengembangan dan _deployment_).

---

## 📁 Struktur Anatomi (Folder Structure)

Proyek ini mengadopsi standar tata letak proyek Go yang terstruktur berdasarkan domain bisnis:

```text
.
├── cmd/                # Entry point utama dari aplikasi (main.go).
├── config/             # Manajemen konfigurasi dan pemuatan environment variables (.env).
├── internal/           # Kode aplikasi internal (business logic) yang tidak dapat diakses oleh package luar.
│   ├── app/            # Kumpulan domain/modul fitur (User, Job Board, Home, dll).
│   │   ├── handler/    # Presentation layer (Menerima HTTP request dan mengembalikan HTTP response).
│   │   ├── service/    # Business logic layer (Memproses data dan aturan bisnis).
│   │   ├── repository/ # Data access layer (Berinteraksi langsung dengan Database/Redis).
│   │   └── dto/        # Data Transfer Objects (Struktur data untuk request/response).
│   ├── infra/          # Konfigurasi infrastruktur (Koneksi DB, Migrasi, Seeding).
│   └── middleware/     # Interceptor untuk HTTP requests (Validasi JWT token, dll).
└── pkg/                # Reusable packages/libraries yang dapat digunakan oleh modul lain (Email, Xendit, Gemini, Storage).
```
