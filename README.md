# StudyTrack backend

Backend Go/Fiber untuk autentikasi, mata kuliah, dan tugas StudyTrack. Skema dikelola dari model GORM menggunakan `AutoMigrate` saat koneksi database dibuka.

## Setup

1. Salin `config/.env.example` ke `config/.env`, lalu isi kredensial PostgreSQL dan `JWT_SECRET_KEY`.
2. Pastikan PostgreSQL berjalan.
3. Terapkan skema dari root repository (opsional, skema juga otomatis diterapkan saat server atau seed dijalankan):

   ```sh
   go run ./cmd/migrate
   ```

4. Muat data demo:

   ```sh
   go run ./cmd/seed
   ```

   Seeder mengisi seluruh data demo secara idempoten, menggunakan email `demo@studytrack.test` dan password `demo1234`. Password disimpan sebagai hash bcrypt.

5. Jalankan server:

   ```sh
   go run ./cmd/app
   ```

`go run ./cmd/migrate` menjalankan GORM `AutoMigrate` satu kali lalu keluar. `go run ./cmd/app` dan `go run ./cmd/seed` juga memastikan skema terbaru saat membuka koneksi database. Perubahan skema dibuat pada model di `domain/entity/models.go`.

## API

Semua response menggunakan envelope `{ "request_id": "...", "payload": ... }`. Error memakai `payload.error` dari helper response aplikasi.

| Method | Endpoint | Auth | Payload inti |
| --- | --- | --- | --- |
| POST | `/api/auth/register` | Tidak | `{ nama, email, password }` |
| POST | `/api/auth/login` | Tidak | `{ email, password }` |
| GET | `/api/auth/me` | Ya | Profil token |
| GET, POST | `/api/courses` | Ya | `{ name, kode?, dosen? }` |
| PUT, DELETE | `/api/courses/:id` | Ya | Field mata kuliah |
| GET, POST | `/api/tasks` | Ya | `{ mata_kuliah_id, judul, deadline, deskripsi? }` |
| GET, PUT, DELETE | `/api/tasks/:id` | Ya | Field tugas; status `belum` atau `selesai` |

Autentikasi memakai `Authorization: Bearer <token>`. Daftar tugas menyertakan nama mata kuliah dalam `mata_kuliah` dan disusun dari `created_at` terbaru.
