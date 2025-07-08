# Proyek Toko Klontong

Aplikasi web CRUD (Create, Read, Update, Delete) sederhana untuk mengelola daftar produk di sebuah toko kelontong. Dibangun menggunakan bahasa Go dengan arsitektur berlapis untuk memisahkan antara logika bisnis, akses data, dan presentasi.

## Fitur

*   **Daftar Produk:** Menampilkan semua produk yang ada di dalam database.
*   **Tambah Produk:** Menambahkan produk baru ke dalam database melalui form.

## Teknologi yang Digunakan

*   **Bahasa:** Go
*   **Web Framework:** `net/http` (standard library) dengan `julienschmidt/httprouter` untuk routing.
*   **Database:** SQLite
*   **ORM:** GORM
*   **Templating:** `html/template` (standard library)
*   **Testing:** `testing` (standard library) & `stretchr/testify` (untuk assertion dan mock).
*   **Frontend:** Tailwind CSS (via CDN)

## Cara Menjalankan

1.  **Pastikan Go terinstall** di sistem Anda.
2.  **Clone repositori ini.**
3.  **Buka terminal** dan masuk ke direktori proyek.
4.  **Jalankan perintah untuk mengunduh dependensi:**
    ```bash
    go mod tidy
    ```
5.  **Jalankan aplikasi:**
    ```bash
    go run cmd/main.go
    ```
6.  **Buka browser** dan akses:
    *   **Halaman utama:** `http://localhost:8080/`
    *   **Daftar Produk:** `http://localhost:8080/produk`

## Cara Menjalankan Test

Untuk menjalankan semua unit test, gunakan perintah:

```bash
go test ./... -v
```

---

## Arsitektur & Detail Implementasi

Proyek ini menggunakan arsitektur berlapis (layered architecture) untuk memisahkan tanggung jawab.

```
/
├── cmd/                # Entrypoint & Database
├── config/             # Konfigurasi Database
├── controllers/        # Lapisan Presentasi (HTTP Handler)
├── models/             # Struktur Data (Entitas)
├── repositories/       # Lapisan Akses Data (Database)
├── routes/             # Definisi Rute URL
├── services/           # Lapisan Logika Bisnis
├── templates/          # File HTML
├── tests/              # Unit Tests & Mocks
└── utils/              # Fungsi Bantuan
```

### `cmd/main.go`

Titik masuk utama aplikasi.

*   `func main()`:
    1.  `config.ConnectDatabase()`: Menginisialisasi koneksi ke database SQLite.
    2.  `utils.ParseTemplates()`: Memuat dan mem-parsing semua template HTML agar siap digunakan.
    3.  **Dependency Injection**: Menginisialisasi semua lapisan dari bawah ke atas (Repository -> Service -> Controller). Ini memastikan bahwa setiap lapisan menerima dependensi yang dibutuhkannya melalui constructor, bukan membuatnya sendiri.
    4.  `routes.SetupRouter()`: Mengatur semua rute URL aplikasi.
    5.  `http.ListenAndServe()`: Menjalankan server web pada port `8080`.

### `config/database.go`

*   `func ConnectDatabase()`:
    *   Membuka koneksi ke file database `toko.db` menggunakan GORM dengan driver SQLite.
    *   `db.AutoMigrate(&models.Produk{})`: Secara otomatis membuat atau memperbarui skema tabel `produks` di database berdasarkan definisi `models.Produk`.

### `models/produk.go`

*   `struct Produk`:
    *   Mendefinisikan struktur data untuk sebuah produk.
    *   Ini adalah representasi dari tabel `produks` di database.
    *   `gorm.Model`: Menyematkan struct ini secara otomatis menambahkan field `ID`, `CreatedAt`, `UpdatedAt`, dan `DeletedAt`, sesuai dengan konvensi GORM.

### `repositories/produk_repository.go`

Lapisan ini bertanggung jawab penuh untuk komunikasi dengan database.

*   `interface ProdukRepository`: Mendefinisikan "kontrak" atau fungsi apa saja yang harus disediakan oleh sebuah produk repository (`FindAll`, `Save`).
*   `struct produkRepositoryImpl`: Implementasi konkret dari interface di atas. Menyimpan instance koneksi database (`*gorm.DB`).
*   `func NewProdukRepository(...)`: Constructor untuk membuat instance baru dari `produkRepositoryImpl`.
*   `func (r *produkRepositoryImpl) FindAll()`: Mengambil semua data produk dari database menggunakan `r.db.Find()`.
*   `func (r *produkRepositoryImpl) Save(...)`: Menyimpan satu data produk baru ke database menggunakan `r.db.Create()`.

### `services/produk_service.go`

Lapisan ini berisi logika bisnis aplikasi.

*   `interface ProdukService`: Mendefinisikan "kontrak" untuk logika bisnis produk.
*   `struct produkServiceImpl`: Implementasi dari interface service. Menyimpan instance dari `ProdukRepository`.
*   `func NewProdukService(...)`: Constructor untuk `produkServiceImpl`.
*   `func (s *produkServiceImpl) FindAll()`: Logika untuk mencari semua produk. Saat ini, ia hanya meneruskan panggilan ke repository. Di masa depan, bisa ditambahkan logika caching, validasi, dll.
*   `func (s *produkServiceImpl) Create(...)`: Logika untuk membuat produk baru. Meneruskan panggilan ke repository untuk menyimpan data.

### `controllers/produk_controller.go`

Lapisan ini menjembatani antara permintaan HTTP dari pengguna dan logika bisnis.

*   `interface ProdukController`: Mendefinisikan "kontrak" untuk handler HTTP yang terkait dengan produk.
*   `struct produkcotrollerimpl`: Implementasi dari interface controller. Menyimpan instance `ProdukService` dan template yang sudah di-parsing.
*   `func NewProdukController(...)`: Constructor untuk `produkcotrollerimpl`.
*   `func (c *produkcotrollerimpl) Index(...)`: Handler untuk `GET /produk`. Memanggil `produkService.FindAll()` untuk mendapatkan data, lalu merender template `list.html` dengan data tersebut.
*   `func (c *produkcotrollerimpl) Form(...)`: Handler untuk `GET /produk/tambah`. Hanya merender template `form.html`.
*   `func (c *produkcotrollerimpl) Store(...)`: Handler untuk `POST /produk/tambah`. Mem-parsing data dari form, membuat struct `models.Produk`, memanggil `produkService.Create()`, dan mengalihkan pengguna kembali ke halaman daftar produk.

### `routes/routes.go`

*   `func SetupRouter(...)`:
    *   Membuat instance baru dari `httprouter.Router`.
    *   Mendefinisikan setiap rute (misal: `GET /produk`), metode HTTP-nya, dan fungsi controller mana yang harus menanganinya.
    *   Mengembalikan router yang sudah dikonfigurasi.

### `utils/template.go`

*   `func ParseTemplates()`:
    *   Fungsi bantuan yang sangat penting untuk performa.
    *   `template.ParseGlob(...)`: Membaca semua file `.html` dari direktori `templates/produk` dan `templates/partials` sekali saja saat aplikasi pertama kali berjalan.
    *   Hasilnya disimpan dalam satu objek `*template.Template` dan digunakan kembali setiap kali ada permintaan untuk merender halaman, menghindari pembacaan file dari disk berulang kali.

---

## Log Proses Pembuatan Dokumentasi

Berikut adalah langkah-langkah yang saya (Gemini) lakukan untuk menganalisis proyek dan membuat file `README.md` ini:

1.  **Permintaan Awal**: Pengguna meminta saya untuk membaca seluruh proyek dan membuat dokumentasi, serta melacak semua langkah yang dilakukan.
2.  **Membaca `README.md` Awal**: Saya memulai dengan membaca file `README.md` yang sudah ada untuk memahami konteks awal. Isinya sangat dasar.
3.  **Analisis Proyek**: Saya menggunakan tool `read_many_files` untuk membaca semua file dengan ekstensi `.go`, `.html`, dan `.md`. Ini memberi saya pemahaman lengkap tentang kode sumber, struktur, tampilan, dan pengujian.
4.  **Menyusun Draf Awal**: Berdasarkan analisis, saya menyusun draf pertama `README.md` yang mencakup deskripsi umum, fitur, struktur direktori, dan cara menjalankan proyek.
5.  **Revisi Berdasarkan Feedback**: Pengguna meminta detail lebih lanjut, termasuk penjelasan per fungsi dan pencatatan (log) dari proses yang saya lakukan.
6.  **Analisis Mendalam**: Saya kembali "membaca" kode di setiap file `.go` untuk memahami tujuan spesifik dari setiap fungsi, struct, dan interface.
7.  **Menulis Dokumentasi Final**: Saya menggabungkan draf awal dengan analisis mendalam per fungsi dan menambahkan bagian "Log Proses Pembuatan Dokumentasi" ini untuk memenuhi semua permintaan pengguna.
8.  **Menulis ke File**: Saya menulis konten final ini ke dalam file `README.md` di proyek Anda.

---