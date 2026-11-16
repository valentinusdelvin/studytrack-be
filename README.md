# StudyTrack Backend

Go/Fiber backend for StudyTrack authentication, courses, and tasks. The schema is managed from GORM models using `AutoMigrate` when the database connection is opened.

## Setup

1. Copy `config/.env.example` to `config/.env`, then enter your PostgreSQL credentials and `JWT_SECRET_KEY`.
2. Make sure PostgreSQL is running.
3. Apply the schema from the repository root (optional; the schema is also applied automatically when the server or seeder runs):

   ```sh
   go run ./cmd/migrate
   ```

4. Load the demo data:

   ```sh
   go run ./cmd/seed
   ```

   The seeder idempotently populates all demo data using the email `demo@studytrack.test` and password `demo1234`. Passwords are stored as bcrypt hashes.

5. Run the server:

   ```sh
   go run ./cmd/app
   ```

`go run ./cmd/migrate` runs GORM `AutoMigrate` once and then exits. `go run ./cmd/app` and `go run ./cmd/seed` also ensure the latest schema is applied when opening a database connection. Schema changes should be made to the models in `domain/entity/models.go`.

## API

Successful responses contain JSON as specified by each endpoint. Error responses include an `error` field in the JSON body.

| Method | Endpoint | Auth | Main payload |
| --- | --- | --- | --- |
| POST | `/api/auth/register` | No | `{ nama, email, password }` |
| POST | `/api/auth/login` | No | `{ email, password }` |
| GET | `/api/auth/me` | Yes | Token profile |
| GET, POST | `/api/courses` | Yes | `{ name, kode?, dosen? }` |
| PUT, DELETE | `/api/courses/:id` | Yes | Course fields |
| GET, POST | `/api/tasks` | Yes | `{ mata_kuliah_id, judul, deadline, deskripsi? }` |
| GET, PUT, DELETE | `/api/tasks/:id` | Yes | Task fields; status is `belum` or `selesai` |

Authentication uses `Authorization: Bearer <token>`.

`GET /api/tasks` returns the user's tasks across all courses. Each task includes `mata_kuliah` and `status_deadline`, which is calculated when the request is made.

| `status_deadline` | Condition |
| --- | --- |
| `due_soon` | A task with status `belum` and 0 to less than 24 hours remaining |
| `overdue` | A task with status `belum` whose deadline has passed |
| `normal` | A task with status `selesai`, or a task with status `belum` and at least 24 hours remaining |

`due_soon` tasks appear first, ordered by the nearest deadline. Other tasks are ordered by most recent `created_at`. The `created_at` order is also used when two `due_soon` tasks have the same deadline.

Deadlines without an offset use WIB (`Asia/Jakarta`). RFC3339 deadlines with an offset are converted to WIB. Responses return deadlines in WIB using the `YYYY-MM-DDTHH:mm:ss` format.
