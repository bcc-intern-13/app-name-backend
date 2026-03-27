package database

import (
	"log"

	careermapping "github.com/bcc-intern-13/WorkAble-backend/internal/app/career_mapping/entity"
	"gorm.io/gorm"
)

func Seed(db *gorm.DB) {
	seedQuestions(db)
	seedCategories(db)
}

func seedQuestions(db *gorm.DB) {
	var count int64
	db.Model(&careermapping.CareerMappingQuestion{}).Count(&count)
	if count > 0 {
		log.Println("Questions already seeded, skipping...")
		return
	}

	questions := []careermapping.CareerMappingQuestion{
		{
			Number:   1,
			Question: "Saat punya waktu bebas, kamu lebih sering?",
			OptionA:  "Coret-coret ide, gambar, atau eksplorasi referensi visual",
			OptionB:  "Coba aplikasi baru, nonton tutorial tech, atau otak-atik perangkat",
			OptionC:  "Ngobrol panjang dengan teman atau gabung komunitas online",
			OptionD:  "Rapikan file, buat to-do list, atau beresin sesuatu sampai tuntas",
			ScoreA:   []byte(`{"KR": 2, "TK": 1}`),
			ScoreB:   []byte(`{"TK": 2, "KO": 1}`),
			ScoreC:   []byte(`{"KO": 2, "ED": 1}`),
			ScoreD:   []byte(`{"AD": 2, "OP": 1}`),
		},
		{
			Number:   2,
			Question: "Kalau harus menjelaskan sesuatu ke orang lain, kamu lebih suka?",
			OptionA:  "Bikin infografis atau sketsa visual",
			OptionB:  "Buat diagram alur atau presentasi data",
			OptionC:  "Ngobrol langsung dan tanya balik untuk pastikan mereka paham",
			OptionD:  "Buat dokumen terstruktur atau tunjukkan langsung step by step",
			ScoreA:   []byte(`{"KR": 2, "TK": 1}`),
			ScoreB:   []byte(`{"TK": 2, "AD": 1}`),
			ScoreC:   []byte(`{"KO": 2, "ED": 1}`),
			ScoreD:   []byte(`{"AD": 2, "OP": 1}`),
		},
		{
			Number:   3,
			Question: "Kamu paling mudah masuk ke kondisi fokus total saat?",
			OptionA:  "Membuat sesuatu yang punya nilai estetika",
			OptionB:  "Memecahkan masalah teknis yang kompleks",
			OptionC:  "Berdiskusi atau brainstorming dengan orang lain",
			OptionD:  "Mengorganisir data atau mengerjakan tugas berulang dengan teliti",
			ScoreA:   []byte(`{"KR": 2, "TK": 1}`),
			ScoreB:   []byte(`{"TK": 2, "AD": 1}`),
			ScoreC:   []byte(`{"KO": 2, "ED": 1}`),
			ScoreD:   []byte(`{"AD": 2, "OP": 1}`),
		},
		{
			Number:   4,
			Question: "Di antara kalimat ini, mana yang paling menggambarkan kamu?",
			OptionA:  "Saya selalu punya ide visual di kepala",
			OptionB:  "Saya suka tahu cara kerja sesuatu di balik layar",
			OptionC:  "Saya mudah akrab dengan orang baru dan senang bantu orang berkembang",
			OptionD:  "Saya tidak tenang kalau data belum rapi atau pekerjaan belum sesuai prosedur",
			ScoreA:   []byte(`{"KR": 2, "TK": 1}`),
			ScoreB:   []byte(`{"TK": 2, "KR": 1}`),
			ScoreC:   []byte(`{"KO": 2, "ED": 1}`),
			ScoreD:   []byte(`{"AD": 2, "OP": 1}`),
		},
		{
			Number:   5,
			Question: "Di tim baru, peran apa yang paling natural kamu ambil?",
			OptionA:  "Yang bikin materi presentasi atau desain visualnya",
			OptionB:  "Yang urus tools, sistem, atau infrastruktur digitalnya",
			OptionC:  "Yang jadi penghubung antar anggota tim",
			OptionD:  "Yang kelola jadwal dan dokumentasi atau pastikan semua tugas selesai tepat waktu",
			ScoreA:   []byte(`{"KR": 2, "TK": 1}`),
			ScoreB:   []byte(`{"TK": 2, "AD": 1}`),
			ScoreC:   []byte(`{"KO": 2, "ED": 1}`),
			ScoreD:   []byte(`{"AD": 2, "OP": 1}`),
		},
		{
			Number:   6,
			Question: "Kata orang-orang di sekitarmu, kamu dikenal sebagai?",
			OptionA:  "Orang yang kreatif dan punya selera estetika tinggi",
			OptionB:  "Orang yang jago teknologi atau problem solver",
			OptionC:  "Orang yang enak diajak ngobrol dan sabar menjelaskan",
			OptionD:  "Orang yang terorganisir, detail, dan tidak suka setengah-setengah",
			ScoreA:   []byte(`{"KR": 2, "TK": 1}`),
			ScoreB:   []byte(`{"TK": 2, "KR": 1}`),
			ScoreC:   []byte(`{"KO": 2, "ED": 1}`),
			ScoreD:   []byte(`{"AD": 2, "OP": 1}`),
		},
		{
			Number:   7,
			Question: "Kalau kamu gagal dalam sesuatu, reaksi pertamamu?",
			OptionA:  "Coba pendekatan visual atau kreatif yang berbeda",
			OptionB:  "Analisis di mana sistem atau logikanya yang salah",
			OptionC:  "Cari feedback dari orang lain",
			OptionD:  "Review ulang catatan atau ulangi dari awal lebih hati-hati",
			ScoreA:   []byte(`{"KR": 2, "TK": 1}`),
			ScoreB:   []byte(`{"TK": 2, "AD": 1}`),
			ScoreC:   []byte(`{"KO": 2, "ED": 1}`),
			ScoreD:   []byte(`{"AD": 2, "OP": 1}`),
		},
		{
			Number:   8,
			Question: "Jenis feedback yang paling berarti bagimu?",
			OptionA:  "Desain atau karya kamu selalu punya ciri khas",
			OptionB:  "Solusi yang kamu buat efisien dan cerdas",
			OptionC:  "Kamu cairkan suasana dan penjelasanmu mudah dipahami",
			OptionD:  "Laporan kamu rapi dan kamu bisa diandalkan sampai selesai",
			ScoreA:   []byte(`{"KR": 2, "TK": 1}`),
			ScoreB:   []byte(`{"TK": 2, "AD": 1}`),
			ScoreC:   []byte(`{"KO": 2, "ED": 1}`),
			ScoreD:   []byte(`{"AD": 2, "OP": 1}`),
		},
		{
			Number:   9,
			Question: "Kalau ada proyek besar, bagian mana yang paling kamu tunggu?",
			OptionA:  "Tahap ideasi dan membuat konsep visual",
			OptionB:  "Tahap build yaitu membangun sistemnya",
			OptionC:  "Tahap presentasi dan pitching ke stakeholder",
			OptionD:  "Tahap perencanaan roadmap atau eksekusi dan delivery",
			ScoreA:   []byte(`{"KR": 2, "TK": 1}`),
			ScoreB:   []byte(`{"TK": 2, "KR": 1}`),
			ScoreC:   []byte(`{"KO": 2, "ED": 1}`),
			ScoreD:   []byte(`{"AD": 2, "OP": 1}`),
		},
		{
			Number:   10,
			Question: "Lingkungan kerja seperti apa yang bikin kamu paling produktif?",
			OptionA:  "Bebas bereksperimen dan tidak terlalu kaku",
			OptionB:  "Punya tools yang lengkap dan sistem yang jelas",
			OptionC:  "Banyak kolaborasi dan komunikasi antar tim",
			OptionD:  "Terstruktur ada standar jelas dan target nyata terukur setiap hari",
			ScoreA:   []byte(`{"KR": 2, "TK": 1}`),
			ScoreB:   []byte(`{"TK": 2, "AD": 1}`),
			ScoreC:   []byte(`{"KO": 2, "ED": 1}`),
			ScoreD:   []byte(`{"AD": 2, "OP": 1}`),
		},
		{
			Number:   11,
			Question: "Saat kamu harus belajar skill baru, kamu lebih suka?",
			OptionA:  "Langsung praktik dan lihat hasilnya secara visual",
			OptionB:  "Pahami konsep dasarnya dulu baru eksekusi",
			OptionC:  "Belajar bareng orang lain atau join komunitas",
			OptionD:  "Ikuti panduan step by step atau langsung kerjakan task nyata",
			ScoreA:   []byte(`{"KR": 2, "TK": 1}`),
			ScoreB:   []byte(`{"TK": 2, "AD": 1}`),
			ScoreC:   []byte(`{"KO": 2, "ED": 1}`),
			ScoreD:   []byte(`{"AD": 2, "OP": 1}`),
		},
		{
			Number:   12,
			Question: "Kamu merasa paling tidak nyaman saat?",
			OptionA:  "Diminta kerja tanpa ada ruang untuk ekspresi",
			OptionB:  "Harus kerja tanpa data atau logika yang jelas",
			OptionC:  "Tidak ada komunikasi dan harus kerja sendiri terus",
			OptionD:  "Sistem berantakan atau pekerjaan tidak selesai dan hasilnya tidak terukur",
			ScoreA:   []byte(`{"KR": 2, "TK": 1}`),
			ScoreB:   []byte(`{"TK": 2, "AD": 1}`),
			ScoreC:   []byte(`{"KO": 2, "ED": 1}`),
			ScoreD:   []byte(`{"AD": 2, "OP": 1}`),
		},
		{
			Number:   13,
			Question: "Kalau kamu punya satu hari bebas penuh tanpa gadget, kamu akan?",
			OptionA:  "Melukis, foto-foto, atau buat kerajinan tangan",
			OptionB:  "Bongkar-pasang sesuatu untuk tahu cara kerjanya",
			OptionC:  "Kumpul bareng teman atau keluarga",
			OptionD:  "Rapikan rumah atau selesaikan satu proyek kecil dari awal sampai akhir",
			ScoreA:   []byte(`{"KR": 2, "TK": 1}`),
			ScoreB:   []byte(`{"TK": 2, "KR": 1}`),
			ScoreC:   []byte(`{"KO": 2, "ED": 1}`),
			ScoreD:   []byte(`{"AD": 2, "OP": 1}`),
		},
		{
			Number:   14,
			Question: "Apa yang paling sering kamu lakukan secara spontan?",
			OptionA:  "Screenshot atau simpan referensi visual yang menarik",
			OptionB:  "Cari tahu cara kerja teknologi atau aplikasi baru",
			OptionC:  "Hubungi teman lama atau catat hal menarik yang baru dipelajari",
			OptionD:  "Buat list atau kategorisasi atau selesaikan tugas kecil yang tertunda",
			ScoreA:   []byte(`{"KR": 2, "TK": 1}`),
			ScoreB:   []byte(`{"TK": 2, "KR": 1}`),
			ScoreC:   []byte(`{"KO": 2, "ED": 1}`),
			ScoreD:   []byte(`{"AD": 2, "OP": 1}`),
		},
		{
			Number:   15,
			Question: "Kalau diminta memilih satu kata untuk menggambarkan cara kerjamu?",
			OptionA:  "Imajinatif",
			OptionB:  "Analitis",
			OptionC:  "Kolaboratif atau Reflektif",
			OptionD:  "Sistematis atau Eksekutif",
			ScoreA:   []byte(`{"KR": 2, "TK": 1}`),
			ScoreB:   []byte(`{"TK": 2, "AD": 1}`),
			ScoreC:   []byte(`{"KO": 2, "ED": 1}`),
			ScoreD:   []byte(`{"AD": 2, "OP": 1}`),
		},
		{
			Number:   16,
			Question: "Saat kamu melihat produk atau layanan yang buruk, reaksimu?",
			OptionA:  "Langsung kepikiran bagaimana tampilannya bisa lebih baik",
			OptionB:  "Penasaran apa masalah teknis atau sistemnya",
			OptionC:  "Ingin kasih feedback atau kepikiran bagaimana edukasi penggunanya diperbaiki",
			OptionD:  "Perhatikan alur proses atau cari cara paling praktis untuk solve-nya",
			ScoreA:   []byte(`{"KR": 2, "TK": 1}`),
			ScoreB:   []byte(`{"TK": 2, "AD": 1}`),
			ScoreC:   []byte(`{"KO": 2, "ED": 1}`),
			ScoreD:   []byte(`{"AD": 2, "OP": 1}`),
		},
		{
			Number:   17,
			Question: "Tipe pencapaian mana yang paling membuatmu bangga?",
			OptionA:  "Karya yang diakui dan diapresiasi banyak orang",
			OptionB:  "Sistem atau solusi yang kamu bangun benar-benar berfungsi",
			OptionC:  "Seseorang bilang kamu sangat membantu atau mengubah cara pandang mereka",
			OptionD:  "Proyek selesai tepat waktu rapi atau target tercapai dan hasilnya langsung dirasakan",
			ScoreA:   []byte(`{"KR": 2, "TK": 1}`),
			ScoreB:   []byte(`{"TK": 2, "KR": 1}`),
			ScoreC:   []byte(`{"KO": 2, "ED": 1}`),
			ScoreD:   []byte(`{"AD": 2, "OP": 1}`),
		},
		{
			Number:   18,
			Question: "Kalau kamu harus memilih satu kalimat tentang karir idealmu?",
			OptionA:  "Saya ingin berkarya dan meninggalkan jejak kreatif",
			OptionB:  "Saya ingin membangun sesuatu yang berdampak lewat teknologi",
			OptionC:  "Saya ingin pekerjaan yang melibatkan interaksi manusia dan bermanfaat bagi masyarakat",
			OptionD:  "Saya ingin pekerjaan terstruktur dengan target nyata yang bisa dicapai setiap hari",
			ScoreA:   []byte(`{"KR": 2, "TK": 1}`),
			ScoreB:   []byte(`{"TK": 2, "KR": 1}`),
			ScoreC:   []byte(`{"KO": 2, "ED": 1}`),
			ScoreD:   []byte(`{"AD": 2, "OP": 1}`),
		},
		{
			Number:   19,
			Question: "Kalau ada konflik dalam tim, kamu cenderung?",
			OptionA:  "Buat visualisasi masalah supaya semua bisa lihat gambaran besarnya",
			OptionB:  "Analisis akar masalahnya secara objektif berdasarkan data",
			OptionC:  "Jadi mediator atau ajak semua refleksi dan belajar dari situasi ini",
			OptionD:  "Dokumentasikan masalah buat SOP atau fokus ke solusi praktis yang bisa dieksekusi",
			ScoreA:   []byte(`{"KR": 2, "TK": 1}`),
			ScoreB:   []byte(`{"TK": 2, "AD": 1}`),
			ScoreC:   []byte(`{"KO": 2, "ED": 1}`),
			ScoreD:   []byte(`{"AD": 2, "OP": 1}`),
		},
		{
			Number:   20,
			Question: "Kalimat mana yang paling menggambarkan motivasi kerjamu?",
			OptionA:  "Saya bekerja karena ingin mengekspresikan diri dan berkreasi",
			OptionB:  "Saya bekerja karena ingin terus berinovasi dan memecahkan masalah",
			OptionC:  "Saya bekerja karena koneksi dengan orang lain dan ingin memberi dampak positif",
			OptionD:  "Saya bekerja karena suka keteraturan dan kepuasan terbesar adalah menyelesaikan sesuatu",
			ScoreA:   []byte(`{"KR": 2, "TK": 1}`),
			ScoreB:   []byte(`{"TK": 2, "KR": 1}`),
			ScoreC:   []byte(`{"KO": 2, "ED": 1}`),
			ScoreD:   []byte(`{"AD": 2, "OP": 1}`),
		},
	}

	if err := db.Create(&questions).Error; err != nil {
		log.Fatal(" Failed to seed questions:", err)
	}
	log.Println(" Career mapping questions seeded successfully!")
}

func seedCategories(db *gorm.DB) {
	var count int64
	db.Model(&careermapping.CareerCategory{}).Count(&count)
	if count > 0 {
		log.Println("✅ Categories already seeded, skipping...")
		return
	}

	categories := []careermapping.CareerCategory{
		{
			ID:          "KR",
			Name:        "Kreatif dan Seni",
			Description: "Bekerja dengan visual, seni, desain, dan ekspresi kreatif. Cocok untuk individu yang berpikir secara visual, imajinatif, dan senang mengekspresikan diri melalui karya.",
			FormalJobs:  []byte(`["Desainer Grafis","Illustrator","Fotografer","Videografer","Animator","UI/UX Designer","Motion Designer","Content Creator Visual"]`),
			SideJobs:    []byte(`["Jasa desain lepas","Jual produk handmade","Jasa foto dan video","Jual karya seni digital","Toko merchandise custom"]`),
		},
		{
			ID:          "TK",
			Name:        "Teknologi dan Digital",
			Description: "Bekerja dengan tools digital, sistem, dan teknologi. Cocok untuk individu yang analitis, suka memecahkan masalah teknis, dan tertarik pada cara kerja di balik teknologi.",
			FormalJobs:  []byte(`["Web Developer","UI/UX Designer","Data Analyst","Digital Marketer","IT Support","Content Creator","Social Media Specialist","Programmer"]`),
			SideJobs:    []byte(`["Freelance developer","Jasa pembuatan website","Jualan template digital","Jasa IT support","YouTuber atau Podcaster tech"]`),
		},
		{
			ID:          "KO",
			Name:        "Komunikasi dan Orang",
			Description: "Bekerja dengan interaksi manusia, komunikasi, dan relasi. Cocok untuk individu yang mudah bergaul, empatik, dan mendapat energi dari berinteraksi dengan orang lain.",
			FormalJobs:  []byte(`["Customer Service","HRD / Rekrutmen","Public Relations","Sales","Account Manager","Community Manager","Event Coordinator","Resepsionis"]`),
			SideJobs:    []byte(`["Reseller / Dropshipper","Jasa event organizer kecil","Agen properti independen","Wedding organizer"]`),
		},
		{
			ID:          "ED",
			Name:        "Edukasi dan Sosial",
			Description: "Bekerja dengan pengajaran, pemberdayaan, dan dampak sosial. Cocok untuk individu yang sabar, suka berbagi ilmu, dan termotivasi melihat orang lain berkembang.",
			FormalJobs:  []byte(`["Guru / Pengajar","Mentor","Social Worker","Konselor","Trainer / Fasilitator","Penulis Konten Edukasi","Peneliti Sosial"]`),
			SideJobs:    []byte(`["Les privat","Kursus online","Jual modul belajar digital","Pelatihan keterampilan komunitas"]`),
		},
		{
			ID:          "AD",
			Name:        "Administrasi dan Data",
			Description: "Bekerja dengan data, dokumen, angka, dan struktur. Cocok untuk individu yang teliti, terorganisir, suka keteraturan, dan tidak nyaman dengan hal-hal yang ambigu.",
			FormalJobs:  []byte(`["Admin Kantor","Data Entry","Analis Data","Sekretaris","Staf Keuangan","Akuntansi","Operator Input Data"]`),
			SideJobs:    []byte(`["Virtual assistant lepas","Jasa pengetikan dan transkripsi","Jasa pembukuan UMKM","Pengelola toko online"]`),
		},
		{
			ID:          "OP",
			Name:        "Operasional dan Detail",
			Description: "Bekerja dengan prosedur, detail, dan eksekusi terstruktur. Cocok untuk individu yang konsisten, praktis, berorientasi hasil nyata, dan puas saat pekerjaan selesai tuntas.",
			FormalJobs:  []byte(`["Quality Control","Logistik dan Gudang","Operator Produksi","Teknisi","Staff Purchasing","Petugas Kasir","Staff Inventori"]`),
			SideJobs:    []byte(`["Catering rumahan","Laundry rumahan","Toko kelontong","Bengkel kecil","Jasa packing produk"]`),
		},
	}

	if err := db.Create(&categories).Error; err != nil {
		log.Fatal("Failed to seed categories:", err)
	}
	log.Println("Career categories seeded successfully!")
}
