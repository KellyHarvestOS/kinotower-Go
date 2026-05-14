package repositories

import (
	"context"
	"database/sql"
	"errors"
	"fmt"
	"strings"
	"time"

	"github.com/KellyHarvestOS/kinotower-Go/internal/models"
	"github.com/KellyHarvestOS/kinotower-Go/internal/utils"
)

type Repository struct {
	db *sql.DB
}

var ErrNotFound = sql.ErrNoRows

func New(db *sql.DB) *Repository {
	return &Repository{db: db}
}

func (r *Repository) ListFilms(ctx context.Context, filters models.FilmFilters) ([]models.Film, int, error) {
	where := []string{"f.deleted_at IS NULL"}
	args := []any{}
	if filters.Category > 0 {
		args = append(args, filters.Category)
		where = append(where, fmt.Sprintf("EXISTS (SELECT 1 FROM categories_films cf WHERE cf.film_id = f.id AND cf.category_id = $%d)", len(args)))
	}
	if filters.Country > 0 {
		args = append(args, filters.Country)
		where = append(where, fmt.Sprintf("f.country_id = $%d", len(args)))
	}
	if filters.Search != "" {
		args = append(args, "%"+filters.Search+"%")
		where = append(where, fmt.Sprintf("f.name ILIKE $%d", len(args)))
	}
	whereSQL := strings.Join(where, " AND ")

	var total int
	if err := r.db.QueryRowContext(ctx, `SELECT count(*) FROM films f WHERE `+whereSQL, args...).Scan(&total); err != nil {
		return nil, 0, err
	}

	orderBy := map[string]string{"name": "f.name", "year": "f.year_of_issue", "rating": "rating_avg"}[filters.SortBy]
	if orderBy == "" {
		orderBy = "f.name"
	}
	sortDir := "asc"
	if strings.EqualFold(filters.SortDir, "desc") {
		sortDir = "desc"
	}
	args = append(args, filters.Limit, filters.Offset)
	query := fmt.Sprintf(`
		SELECT f.id, f.name, f.duration, f.year_of_issue, f.age, f.link_img, f.link_kinopoisk,
			f.link_video, f.created_at, c.id, c.name,
			COALESCE(AVG(rt.ball), 0)::float, COUNT(DISTINCT rv.id)
		FROM films f
		JOIN countries c ON c.id = f.country_id
		LEFT JOIN ratings rt ON rt.film_id = f.id
		LEFT JOIN reviews rv ON rv.film_id = f.id AND rv.deleted_at IS NULL AND rv.is_approved = TRUE
		WHERE %s
		GROUP BY f.id, c.id
		ORDER BY %s %s
		LIMIT $%d OFFSET $%d`, whereSQL, orderBy, sortDir, len(args)-1, len(args))

	rows, err := r.db.QueryContext(ctx, query, args...)
	if err != nil {
		return nil, 0, err
	}
	defer rows.Close()

	films := []models.Film{}
	for rows.Next() {
		film, err := r.scanFilm(ctx, rows)
		if err != nil {
			return nil, 0, err
		}
		films = append(films, film)
	}
	return films, total, rows.Err()
}

func (r *Repository) GetFilm(ctx context.Context, id int) (models.Film, error) {
	row := r.db.QueryRowContext(ctx, `
		SELECT f.id, f.name, f.duration, f.year_of_issue, f.age, f.link_img, f.link_kinopoisk,
			f.link_video, f.created_at, c.id, c.name,
			COALESCE(AVG(rt.ball), 0)::float, COUNT(DISTINCT rv.id)
		FROM films f
		JOIN countries c ON c.id = f.country_id
		LEFT JOIN ratings rt ON rt.film_id = f.id
		LEFT JOIN reviews rv ON rv.film_id = f.id AND rv.deleted_at IS NULL AND rv.is_approved = TRUE
		WHERE f.id = $1 AND f.deleted_at IS NULL
		GROUP BY f.id, c.id`, id)
	return r.scanFilm(ctx, row)
}

func (r *Repository) scanFilm(ctx context.Context, scanner interface{ Scan(dest ...any) error }) (models.Film, error) {
	var film models.Film
	var img, kp sql.NullString
	var created time.Time
	if err := scanner.Scan(&film.ID, &film.Name, &film.Duration, &film.YearOfIssue, &film.Age, &img, &kp, &film.LinkVideo, &created, &film.Country.ID, &film.Country.Name, &film.RatingAvg, &film.ReviewCount); err != nil {
		return film, err
	}
	film.LinkImg = utils.NullString(img)
	film.LinkKinopoisk = utils.NullString(kp)
	film.CreatedAt = utils.TimeISO(created)
	cats, err := r.FilmCategories(ctx, film.ID)
	if err != nil {
		return film, err
	}
	film.Categories = cats
	return film, nil
}

func (r *Repository) FilmCategories(ctx context.Context, filmID int) ([]models.SimpleRef, error) {
	rows, err := r.db.QueryContext(ctx, `SELECT c.id, c.name FROM categories c JOIN categories_films cf ON cf.category_id = c.id WHERE cf.film_id = $1 AND c.deleted_at IS NULL ORDER BY c.id`, filmID)
	if err != nil {
		return nil, err
	}
	defer rows.Close()
	items := []models.SimpleRef{}
	for rows.Next() {
		var item models.SimpleRef
		if err := rows.Scan(&item.ID, &item.Name); err != nil {
			return nil, err
		}
		items = append(items, item)
	}
	return items, rows.Err()
}

func (r *Repository) ListCategories(ctx context.Context) ([]models.Category, error) {
	rows, err := r.db.QueryContext(ctx, `SELECT c.id, c.name, p.id, p.name, COUNT(f.id) FROM categories c LEFT JOIN categories p ON p.id = c.parent_id LEFT JOIN categories_films cf ON cf.category_id = c.id LEFT JOIN films f ON f.id = cf.film_id AND f.deleted_at IS NULL WHERE c.deleted_at IS NULL GROUP BY c.id, p.id ORDER BY c.id`)
	if err != nil {
		return nil, err
	}
	defer rows.Close()
	items := []models.Category{}
	for rows.Next() {
		var item models.Category
		var pid sql.NullInt64
		var pname sql.NullString
		if err := rows.Scan(&item.ID, &item.Name, &pid, &pname, &item.FilmCount); err != nil {
			return nil, err
		}
		if pid.Valid {
			item.ParentCategory = &models.SimpleRef{ID: int(pid.Int64), Name: pname.String}
		}
		items = append(items, item)
	}
	return items, rows.Err()
}

func (r *Repository) ListCountries(ctx context.Context) ([]models.Country, error) {
	rows, err := r.db.QueryContext(ctx, `SELECT c.id, c.name, COUNT(f.id) FROM countries c LEFT JOIN films f ON f.country_id = c.id AND f.deleted_at IS NULL GROUP BY c.id ORDER BY c.id`)
	if err != nil {
		return nil, err
	}
	defer rows.Close()
	items := []models.Country{}
	for rows.Next() {
		var item models.Country
		if err := rows.Scan(&item.ID, &item.Name, &item.FilmCount); err != nil {
			return nil, err
		}
		items = append(items, item)
	}
	return items, rows.Err()
}

func (r *Repository) ListGenders(ctx context.Context) ([]models.Gender, error) {
	rows, err := r.db.QueryContext(ctx, `SELECT id, name FROM genders ORDER BY id`)
	if err != nil {
		return nil, err
	}
	defer rows.Close()
	items := []models.Gender{}
	for rows.Next() {
		var item models.Gender
		if err := rows.Scan(&item.ID, &item.Name); err != nil {
			return nil, err
		}
		items = append(items, item)
	}
	return items, rows.Err()
}

func (r *Repository) CreateUser(ctx context.Context, req models.SignupRequest, passwordHash string, birthday time.Time) (models.User, error) {
	var user models.User
	err := r.db.QueryRowContext(ctx, `INSERT INTO users (fio, birthday, gender_id, email, password, created_at) VALUES ($1,$2,$3,$4,$5,now()) RETURNING id, fio, email`, req.FIO, birthday, req.GenderID, req.Email, passwordHash).Scan(&user.ID, &user.FIO, &user.Email)
	return user, err
}

func (r *Repository) FindUserByEmail(ctx context.Context, email string) (models.User, error) {
	var user models.User
	err := r.db.QueryRowContext(ctx, `SELECT id, fio, email, password FROM users WHERE email=$1 AND deleted_at IS NULL`, email).Scan(&user.ID, &user.FIO, &user.Email, &user.Password)
	return user, err
}

func (r *Repository) GetUser(ctx context.Context, id int) (models.User, error) {
	var user models.User
	err := r.db.QueryRowContext(ctx, `SELECT u.id, u.fio, u.email, u.birthday, g.id, g.name, (SELECT count(*) FROM reviews WHERE user_id=u.id AND deleted_at IS NULL), (SELECT count(*) FROM ratings WHERE user_id=u.id) FROM users u JOIN genders g ON g.id=u.gender_id WHERE u.id=$1 AND u.deleted_at IS NULL`, id).Scan(&user.ID, &user.FIO, &user.Email, &user.Birthday, &user.GenderID, &user.GenderName, &user.ReviewCount, &user.RatingCount)
	return user, err
}

func (r *Repository) UpdateUser(ctx context.Context, id int, req models.UpdateUserRequest, birthday time.Time) error {
	res, err := r.db.ExecContext(ctx, `UPDATE users SET fio=$1,email=$2,birthday=$3,gender_id=$4 WHERE id=$5 AND deleted_at IS NULL`, req.FIO, req.Email, birthday, req.GenderID, id)
	if err != nil {
		return err
	}
	if rows, _ := res.RowsAffected(); rows == 0 {
		return sql.ErrNoRows
	}
	return nil
}

func (r *Repository) DeleteUser(ctx context.Context, id int) error {
	_, err := r.db.ExecContext(ctx, `UPDATE users SET deleted_at=now() WHERE id=$1`, id)
	return err
}

func (r *Repository) Exists(ctx context.Context, query string, args ...any) bool {
	var ok bool
	_ = r.db.QueryRowContext(ctx, query, args...).Scan(&ok)
	return ok
}

func (r *Repository) CreateReview(ctx context.Context, userID int, req models.CreateReviewRequest) (models.Review, error) {
	var item models.Review
	err := r.db.QueryRowContext(ctx, `INSERT INTO reviews (film_id,user_id,message,created_at) VALUES ($1,$2,$3,now()) RETURNING id, film_id, message, is_approved, created_at`, req.FilmID, userID, req.Message).Scan(&item.ID, &item.FilmID, &item.Message, &item.IsApproved, &item.CreatedAt)
	if err != nil {
		return item, err
	}
	_ = r.db.QueryRowContext(ctx, `SELECT name FROM films WHERE id=$1`, item.FilmID).Scan(&item.FilmName)
	return item, nil
}

func (r *Repository) FilmReviews(ctx context.Context, filmID int) ([]models.Review, error) {
	rows, err := r.db.QueryContext(ctx, `SELECT r.id,u.id,u.fio,r.message,r.created_at FROM reviews r JOIN users u ON u.id=r.user_id WHERE r.film_id=$1 AND r.is_approved=TRUE AND r.deleted_at IS NULL ORDER BY r.created_at DESC`, filmID)
	if err != nil {
		return nil, err
	}
	defer rows.Close()
	items := []models.Review{}
	for rows.Next() {
		var item models.Review
		if err := rows.Scan(&item.ID, &item.UserID, &item.UserFIO, &item.Message, &item.CreatedAt); err != nil {
			return nil, err
		}
		items = append(items, item)
	}
	return items, rows.Err()
}

func (r *Repository) UserReviews(ctx context.Context, userID int) ([]models.Review, error) {
	rows, err := r.db.QueryContext(ctx, `SELECT r.id,f.id,f.name,r.message,r.is_approved,r.created_at FROM reviews r JOIN films f ON f.id=r.film_id WHERE r.user_id=$1 AND r.deleted_at IS NULL ORDER BY r.created_at DESC`, userID)
	if err != nil {
		return nil, err
	}
	defer rows.Close()
	items := []models.Review{}
	for rows.Next() {
		var item models.Review
		if err := rows.Scan(&item.ID, &item.FilmID, &item.FilmName, &item.Message, &item.IsApproved, &item.CreatedAt); err != nil {
			return nil, err
		}
		items = append(items, item)
	}
	return items, rows.Err()
}

func (r *Repository) DeleteReview(ctx context.Context, userID, id int) error {
	res, err := r.db.ExecContext(ctx, `UPDATE reviews SET deleted_at=now() WHERE id=$1 AND user_id=$2 AND deleted_at IS NULL`, id, userID)
	if err != nil {
		return err
	}
	if rows, _ := res.RowsAffected(); rows == 0 {
		return sql.ErrNoRows
	}
	return nil
}

func (r *Repository) CreateRating(ctx context.Context, userID int, req models.CreateRatingRequest) (models.Rating, error) {
	var item models.Rating
	err := r.db.QueryRowContext(ctx, `INSERT INTO ratings (film_id,user_id,ball,created_at) VALUES ($1,$2,$3,now()) RETURNING id, film_id, ball, created_at`, req.FilmID, userID, req.Ball).Scan(&item.ID, &item.FilmID, &item.Score, &item.CreatedAt)
	if err != nil {
		return item, err
	}
	_ = r.db.QueryRowContext(ctx, `SELECT name FROM films WHERE id=$1`, item.FilmID).Scan(&item.FilmName)
	return item, nil
}

func (r *Repository) UserRatings(ctx context.Context, userID int) ([]models.Rating, error) {
	rows, err := r.db.QueryContext(ctx, `SELECT rt.id,f.id,f.name,rt.ball,rt.created_at FROM ratings rt JOIN films f ON f.id=rt.film_id WHERE rt.user_id=$1 ORDER BY rt.created_at DESC`, userID)
	if err != nil {
		return nil, err
	}
	defer rows.Close()
	items := []models.Rating{}
	for rows.Next() {
		var item models.Rating
		if err := rows.Scan(&item.ID, &item.FilmID, &item.FilmName, &item.Score, &item.CreatedAt); err != nil {
			return nil, err
		}
		items = append(items, item)
	}
	return items, rows.Err()
}

func (r *Repository) DeleteRating(ctx context.Context, userID, id int) error {
	res, err := r.db.ExecContext(ctx, `DELETE FROM ratings WHERE id=$1 AND user_id=$2`, id, userID)
	if err != nil {
		return err
	}
	if rows, _ := res.RowsAffected(); rows == 0 {
		return sql.ErrNoRows
	}
	return nil
}

func IsNotFound(err error) bool { return errors.Is(err, sql.ErrNoRows) }
