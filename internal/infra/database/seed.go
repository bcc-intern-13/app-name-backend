package database

import (
	"log"

	careermapping "github.com/bcc-intern-13/app-name-backend/internal/career_mapping/entity"
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
		log.Println("✅ Questions already seeded, skipping...")
		return
	}

	questions := []careermapping.CareerMappingQuestion{
		{
			Nomor:      1,
			Pertanyaan: "Saat punya waktu bebas, kamu lebih sering?",
			PilihanA:   "Coret-coret ide, gambar, atau eksplorasi referensi visual",
			PilihanB:   "Coba aplikasi baru, nonton tutorial tech, atau otak-atik perangkat",
			PilihanC:   "Ngobrol panjang dengan teman atau gabung komunitas online",
			PilihanD:   "Rapikan file, buat to-do list, atau beresin sesuatu sampai tuntas",
			SkorA:      []byte(`{"KR": 2, "TK": 1}`),
			SkorB:      []byte(`{"TK": 2, "KO": 1}`),
			SkorC:      []byte(`{"KO": 2, "ED": 1}`),
			SkorD:      []byte(`{"AD": 2, "OP": 1}`),
		},
		{
			Nomor:      2,
			Pertanyaan: "Kalau harus menjelaskan sesuatu ke orang lain, kamu lebih suka?",
			PilihanA:   "Bikin infografis atau sketsa visual",
			PilihanB:   "Buat diagram alur atau presentasi data",
			PilihanC:   "Ngobrol langsung dan tanya balik untuk pastikan mereka paham",
			PilihanD:   "Buat dokumen terstruktur atau tunjukkan langsung step by step",
			SkorA:      []byte(`{"KR": 2, "TK": 1}`),
			SkorB:      []byte(`{"TK": 2, "AD": 1}`),
			SkorC:      []byte(`{"KO": 2, "ED": 1}`),
			SkorD:      []byte(`{"AD": 2, "OP": 1}`),
		},
		{
			Nomor:      3,
			Pertanyaan: "Kamu paling mudah masuk ke kondisi fokus total saat?",
			PilihanA:   "Membuat sesuatu yang punya nilai estetika",
			PilihanB:   "Memecahkan masalah teknis yang kompleks",
			PilihanC:   "Berdiskusi atau brainstorming dengan orang lain",
			PilihanD:   "Mengorganisir data atau mengerjakan tugas berulang dengan teliti",
			SkorA:      []byte(`{"KR": 2, "TK": 1}`),
			SkorB:      []byte(`{"TK": 2, "AD": 1}`),
			SkorC:      []byte(`{"KO": 2, "ED": 1}`),
			SkorD:      []byte(`{"AD": 2, "OP": 1}`),
		},
		{
			Nomor:      4,
			Pertanyaan: "Di antara kalimat ini, mana yang paling menggambarkan kamu?",
			PilihanA:   "Saya selalu punya ide visual di kepala",
			PilihanB:   "Saya suka tahu cara kerja sesuatu di balik layar",
			PilihanC:   "Saya mudah akrab dengan orang baru dan senang bantu orang berkembang",
			PilihanD:   "Saya tidak tenang kalau data belum rapi atau pekerjaan belum sesuai prosedur",
			SkorA:      []byte(`{"KR": 2, "TK": 1}`),
			SkorB:      []byte(`{"TK": 2, "KR": 1}`),
			SkorC:      []byte(`{"KO": 2, "ED": 1}`),
			SkorD:      []byte(`{"AD": 2, "OP": 1}`),
		},
		{
			Nomor:      5,
			Pertanyaan: "Di tim baru, peran apa yang paling natural kamu ambil?",
			PilihanA:   "Yang bikin materi presentasi atau desain visualnya",
			PilihanB:   "Yang urus tools, sistem, atau infrastruktur digitalnya",
			PilihanC:   "Yang jadi penghubung antar anggota tim",
			PilihanD:   "Yang kelola jadwal dan dokumentasi atau pastikan semua tugas selesai tepat waktu",
			SkorA:      []byte(`{"KR": 2, "TK": 1}`),
			SkorB:      []byte(`{"TK": 2, "AD": 1}`),
			SkorC:      []byte(`{"KO": 2, "ED": 1}`),
			SkorD:      []byte(`{"AD": 2, "OP": 1}`),
		},
		{
			Nomor:      6,
			Pertanyaan: "Kata orang-orang di sekitarmu, kamu dikenal sebagai?",
			PilihanA:   "Orang yang kreatif dan punya selera estetika tinggi",
			PilihanB:   "Orang yang jago teknologi atau problem solver",
			PilihanC:   "Orang yang enak diajak ngobrol dan sabar menjelaskan",
			PilihanD:   "Orang yang terorganisir, detail, dan tidak suka setengah-setengah",
			SkorA:      []byte(`{"KR": 2, "TK": 1}`),
			SkorB:      []byte(`{"TK": 2, "KR": 1}`),
			SkorC:      []byte(`{"KO": 2, "ED": 1}`),
			SkorD:      []byte(`{"AD": 2, "OP": 1}`),
		},
		{
			Nomor:      7,
			Pertanyaan: "Kalau kamu gagal dalam sesuatu, reaksi pertamamu?",
			PilihanA:   "Coba pendekatan visual atau kreatif yang berbeda",
			PilihanB:   "Analisis di mana sistem atau logikanya yang salah",
			PilihanC:   "Cari feedback dari orang lain",
			PilihanD:   "Review ulang catatan atau ulangi dari awal lebih hati-hati",
			SkorA:      []byte(`{"KR": 2, "TK": 1}`),
			SkorB:      []byte(`{"TK": 2, "AD": 1}`),
			SkorC:      []byte(`{"KO": 2, "ED": 1}`),
			SkorD:      []byte(`{"AD": 2, "OP": 1}`),
		},
		{
			Nomor:      8,
			Pertanyaan: "Jenis feedback yang paling berarti bagimu?",
			PilihanA:   "Desain atau karya kamu selalu punya ciri khas",
			PilihanB:   "Solusi yang kamu buat efisien dan cerdas",
			PilihanC:   "Kamu cairkan suasana dan penjelasanmu mudah dipahami",
			PilihanD:   "Laporan kamu rapi dan kamu bisa diandalkan sampai selesai",
			SkorA:      []byte(`{"KR": 2, "TK": 1}`),
			SkorB:      []byte(`{"TK": 2, "AD": 1}`),
			SkorC:      []byte(`{"KO": 2, "ED": 1}`),
			SkorD:      []byte(`{"AD": 2, "OP": 1}`),
		},
		{
			Nomor:      9,
			Pertanyaan: "Kalau ada proyek besar, bagian mana yang paling kamu tunggu?",
			PilihanA:   "Tahap ideasi dan membuat konsep visual",
			PilihanB:   "Tahap build yaitu membangun sistemnya",
			PilihanC:   "Tahap presentasi dan pitching ke stakeholder",
			PilihanD:   "Tahap perencanaan roadmap atau eksekusi dan delivery",
			SkorA:      []byte(`{"KR": 2, "TK": 1}`),
			SkorB:      []byte(`{"TK": 2, "KR": 1}`),
			SkorC:      []byte(`{"KO": 2, "ED": 1}`),
			SkorD:      []byte(`{"AD": 2, "OP": 1}`),
		},
		{
			Nomor:      10,
			Pertanyaan: "Lingkungan kerja seperti apa yang bikin kamu paling produktif?",
			PilihanA:   "Bebas bereksperimen dan tidak terlalu kaku",
			PilihanB:   "Punya tools yang lengkap dan sistem yang jelas",
			PilihanC:   "Banyak kolaborasi dan komunikasi antar tim",
			PilihanD:   "Terstruktur ada standar jelas dan target nyata terukur setiap hari",
			SkorA:      []byte(`{"KR": 2, "TK": 1}`),
			SkorB:      []byte(`{"TK": 2, "AD": 1}`),
			SkorC:      []byte(`{"KO": 2, "ED": 1}`),
			SkorD:      []byte(`{"AD": 2, "OP": 1}`),
		},
		{
			Nomor:      11,
			Pertanyaan: "Saat kamu harus belajar skill baru, kamu lebih suka?",
			PilihanA:   "Langsung praktik dan lihat hasilnya secara visual",
			PilihanB:   "Pahami konsep dasarnya dulu baru eksekusi",
			PilihanC:   "Belajar bareng orang lain atau join komunitas",
			PilihanD:   "Ikuti panduan step by step atau langsung kerjakan task nyata",
			SkorA:      []byte(`{"KR": 2, "TK": 1}`),
			SkorB:      []byte(`{"TK": 2, "AD": 1}`),
			SkorC:      []byte(`{"KO": 2, "ED": 1}`),
			SkorD:      []byte(`{"AD": 2, "OP": 1}`),
		},
		{
			Nomor:      12,
			Pertanyaan: "Kamu merasa paling tidak nyaman saat?",
			PilihanA:   "Diminta kerja tanpa ada ruang untuk ekspresi",
			PilihanB:   "Harus kerja tanpa data atau logika yang jelas",
			PilihanC:   "Tidak ada komunikasi dan harus kerja sendiri terus",
			PilihanD:   "Sistem berantakan atau pekerjaan tidak selesai dan hasilnya tidak terukur",
			SkorA:      []byte(`{"KR": 2, "TK": 1}`),
			SkorB:      []byte(`{"TK": 2, "AD": 1}`),
			SkorC:      []byte(`{"KO": 2, "ED": 1}`),
			SkorD:      []byte(`{"AD": 2, "OP": 1}`),
		},
		{
			Nomor:      13,
			Pertanyaan: "Kalau kamu punya satu hari bebas penuh tanpa gadget, kamu akan?",
			PilihanA:   "Melukis, foto-foto, atau buat kerajinan tangan",
			PilihanB:   "Bongkar-pasang sesuatu untuk tahu cara kerjanya",
			PilihanC:   "Kumpul bareng teman atau keluarga",
			PilihanD:   "Rapikan rumah atau selesaikan satu proyek kecil dari awal sampai akhir",
			SkorA:      []byte(`{"KR": 2, "TK": 1}`),
			SkorB:      []byte(`{"TK": 2, "KR": 1}`),
			SkorC:      []byte(`{"KO": 2, "ED": 1}`),
			SkorD:      []byte(`{"AD": 2, "OP": 1}`),
		},
		{
			Nomor:      14,
			Pertanyaan: "Apa yang paling sering kamu lakukan secara spontan?",
			PilihanA:   "Screenshot atau simpan referensi visual yang menarik",
			PilihanB:   "Cari tahu cara kerja teknologi atau aplikasi baru",
			PilihanC:   "Hubungi teman lama atau catat hal menarik yang baru dipelajari",
			PilihanD:   "Buat list atau kategorisasi atau selesaikan tugas kecil yang tertunda",
			SkorA:      []byte(`{"KR": 2, "TK": 1}`),
			SkorB:      []byte(`{"TK": 2, "KR": 1}`),
			SkorC:      []byte(`{"KO": 2, "ED": 1}`),
			SkorD:      []byte(`{"AD": 2, "OP": 1}`),
		},
		{
			Nomor:      15,
			Pertanyaan: "Kalau diminta memilih satu kata untuk menggambarkan cara kerjamu?",
			PilihanA:   "Imajinatif",
			PilihanB:   "Analitis",
			PilihanC:   "Kolaboratif atau Reflektif",
			PilihanD:   "Sistematis atau Eksekutif",
			SkorA:      []byte(`{"KR": 2, "TK": 1}`),
			SkorB:      []byte(`{"TK": 2, "AD": 1}`),
			SkorC:      []byte(`{"KO": 2, "ED": 1}`),
			SkorD:      []byte(`{"AD": 2, "OP": 1}`),
		},
		{
			Nomor:      16,
			Pertanyaan: "Saat kamu melihat produk atau layanan yang buruk, reaksimu?",
			PilihanA:   "Langsung kepikiran bagaimana tampilannya bisa lebih baik",
			PilihanB:   "Penasaran apa masalah teknis atau sistemnya",
			PilihanC:   "Ingin kasih feedback atau kepikiran bagaimana edukasi penggunanya diperbaiki",
			PilihanD:   "Perhatikan alur proses atau cari cara paling praktis untuk solve-nya",
			SkorA:      []byte(`{"KR": 2, "TK": 1}`),
			SkorB:      []byte(`{"TK": 2, "AD": 1}`),
			SkorC:      []byte(`{"KO": 2, "ED": 1}`),
			SkorD:      []byte(`{"AD": 2, "OP": 1}`),
		},
		{
			Nomor:      17,
			Pertanyaan: "Tipe pencapaian mana yang paling membuatmu bangga?",
			PilihanA:   "Karya yang diakui dan diapresiasi banyak orang",
			PilihanB:   "Sistem atau solusi yang kamu bangun benar-benar berfungsi",
			PilihanC:   "Seseorang bilang kamu sangat membantu atau mengubah cara pandang mereka",
			PilihanD:   "Proyek selesai tepat waktu rapi atau target tercapai dan hasilnya langsung dirasakan",
			SkorA:      []byte(`{"KR": 2, "TK": 1}`),
			SkorB:      []byte(`{"TK": 2, "KR": 1}`),
			SkorC:      []byte(`{"KO": 2, "ED": 1}`),
			SkorD:      []byte(`{"AD": 2, "OP": 1}`),
		},
		{
			Nomor:      18,
			Pertanyaan: "Kalau kamu harus memilih satu kalimat tentang karir idealmu?",
			PilihanA:   "Saya ingin berkarya dan meninggalkan jejak kreatif",
			PilihanB:   "Saya ingin membangun sesuatu yang berdampak lewat teknologi",
			PilihanC:   "Saya ingin pekerjaan yang melibatkan interaksi manusia dan bermanfaat bagi masyarakat",
			PilihanD:   "Saya ingin pekerjaan terstruktur dengan target nyata yang bisa dicapai setiap hari",
			SkorA:      []byte(`{"KR": 2, "TK": 1}`),
			SkorB:      []byte(`{"TK": 2, "KR": 1}`),
			SkorC:      []byte(`{"KO": 2, "ED": 1}`),
			SkorD:      []byte(`{"AD": 2, "OP": 1}`),
		},
		{
			Nomor:      19,
			Pertanyaan: "Kalau ada konflik dalam tim, kamu cenderung?",
			PilihanA:   "Buat visualisasi masalah supaya semua bisa lihat gambaran besarnya",
			PilihanB:   "Analisis akar masalahnya secara objektif berdasarkan data",
			PilihanC:   "Jadi mediator atau ajak semua refleksi dan belajar dari situasi ini",
			PilihanD:   "Dokumentasikan masalah buat SOP atau fokus ke solusi praktis yang bisa dieksekusi",
			SkorA:      []byte(`{"KR": 2, "TK": 1}`),
			SkorB:      []byte(`{"TK": 2, "AD": 1}`),
			SkorC:      []byte(`{"KO": 2, "ED": 1}`),
			SkorD:      []byte(`{"AD": 2, "OP": 1}`),
		},
		{
			Nomor:      20,
			Pertanyaan: "Kalimat mana yang paling menggambarkan motivasi kerjamu?",
			PilihanA:   "Saya bekerja karena ingin mengekspresikan diri dan berkreasi",
			PilihanB:   "Saya bekerja karena ingin terus berinovasi dan memecahkan masalah",
			PilihanC:   "Saya bekerja karena koneksi dengan orang lain dan ingin memberi dampak positif",
			PilihanD:   "Saya bekerja karena suka keteraturan dan kepuasan terbesar adalah menyelesaikan sesuatu",
			SkorA:      []byte(`{"KR": 2, "TK": 1}`),
			SkorB:      []byte(`{"TK": 2, "KR": 1}`),
			SkorC:      []byte(`{"KO": 2, "ED": 1}`),
			SkorD:      []byte(`{"AD": 2, "OP": 1}`),
		},
	}

	if err := db.Create(&questions).Error; err != nil {
		log.Fatal("❌ Failed to seed questions:", err)
	}
	log.Println("✅ Career mapping questions seeded successfully!")
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
		log.Fatal("❌ Failed to seed categories:", err)
	}
	log.Println("✅ Career categories seeded successfully!")
}
