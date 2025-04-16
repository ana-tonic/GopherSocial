package store

import (
	"context"
	"database/sql"
	"errors"
	"strconv"

	"github.com/lib/pq"
)

type Post struct {
	ID        int64     `json:"id"`
	Content   string    `json:"content"`
	Title     string    `json:"title"`
	UserID    int64     `json:"user_id"`
	Tags      []string  `json:"tags"`
	CreatedAt string    `json:"created_at"`
	UpdatedAt string    `json:"updated_at"`
	Version   int       `json:"version"`
	Comments  []Comment `json:"comments"`
	User      User      `json:"user"`
}

type PostWithMetadata struct {
	Post
	CommentsCount int `json:"comments_count"`
}

type PostStore struct {
	db *sql.DB
}

func (s *PostStore) Create(ctx context.Context, post *Post) error {
	query := `INSERT INTO posts (content, title, user_id, tags)
	VALUES ($1, $2, $3, $4) RETURNING id, created_at, updated_at`

	ctx, cancel := context.WithTimeout(ctx, QueryTimeoutDuration)
	defer cancel()

	err := s.db.QueryRowContext(
		ctx,
		query,
		post.Content,
		post.Title,
		post.UserID,
		pq.Array(post.Tags),
	).Scan(
		&post.ID,
		&post.CreatedAt,
		&post.UpdatedAt,
	)
	if err != nil {
		return err
	}

	return nil
}

func (s *PostStore) GetByID(ctx context.Context, id int64) (*Post, error) {
	query := `
		SELECT id, content, title, user_id, tags, created_at, updated_at, version 
		FROM posts 
		WHERE id = $1
		`

	ctx, cancel := context.WithTimeout(ctx, QueryTimeoutDuration)
	defer cancel()

	var post Post

	err := s.db.QueryRowContext(ctx, query, id).Scan(
		&post.ID,
		&post.Content,
		&post.Title,
		&post.UserID,
		pq.Array(&post.Tags),
		&post.CreatedAt,
		&post.UpdatedAt,
		&post.Version,
	)
	if err != nil {
		switch {
		case errors.Is(err, sql.ErrNoRows):
			return nil, ErrNotFound
		default:
			return nil, err
		}
	}

	return &post, nil
}

func (s *PostStore) Delete(ctx context.Context, id int64) error {
	query := `DELETE FROM posts WHERE id = $1`

	ctx, cancel := context.WithTimeout(ctx, QueryTimeoutDuration)
	defer cancel()

	res, err := s.db.ExecContext(ctx, query, id)
	if err != nil {
		return err
	}

	rows, err := res.RowsAffected()
	if err != nil {
		return err
	}

	if rows == 0 {
		return ErrNotFound
	}

	return nil
}

func (s *PostStore) Update(ctx context.Context, post *Post) error {
	query := `UPDATE posts 
	SET content = $1, title = $2, version = version + 1
	WHERE id = $3 AND version = $4
	RETURNING version
	`

	err := s.db.QueryRowContext(
		ctx,
		query,
		post.Content,
		post.Title,
		post.ID,
		post.Version,
	).Scan(&post.Version)
	if err != nil {
		switch {
		case errors.Is(err, sql.ErrNoRows):
			return ErrNotFound
		default:
			return err
		}
	}

	return nil
}

func (s *PostStore) GetUserFeed(ctx context.Context, userID int64, fq PaginatedFeedQuery) ([]PostWithMetadata, error) {
	orderBy := "DESC"
	if fq.Sort == "asc" {
		orderBy = "ASC"
	}

	query := `
		SELECT 
			p.id, p.user_id, p.title, p.content, p.created_at, p.version, p.tags,
			u.username, 
			COUNT(c.id) as comments_count
		FROM posts p
		LEFT JOIN comments c ON c.post_id = p.id
		LEFT JOIN users u ON p.user_id = u.id
		LEFT JOIN followers f ON f.user_id = $1 AND f.follower_id = p.user_id
		WHERE (p.user_id = $1 OR f.user_id IS NOT NULL)
	`

	args := []interface{}{userID}
	argPosition := 2

	if fq.Search != "" {
		query += ` AND (p.title ILIKE $` + strconv.Itoa(argPosition) +
			` OR p.content ILIKE $` + strconv.Itoa(argPosition) + `)`
		args = append(args, "%"+fq.Search+"%")
		argPosition++
	}

	if len(fq.Tags) > 0 {
		query += ` AND p.tags && $` + strconv.Itoa(argPosition)
		args = append(args, pq.Array(fq.Tags))
		argPosition++
	}

	if fq.Since != "" {
		query += ` AND p.created_at >= $` + strconv.Itoa(argPosition)
		args = append(args, fq.Since)
		argPosition++
	}

	if fq.Until != "" {
		query += ` AND p.created_at <= $` + strconv.Itoa(argPosition)
		args = append(args, fq.Until)
		argPosition++
	}

	query += `
		GROUP BY p.id, u.username
		ORDER BY p.created_at ` + orderBy + `
		LIMIT $` + strconv.Itoa(argPosition) + ` OFFSET $` + strconv.Itoa(argPosition+1) + `
	`

	args = append(args, fq.Limit, fq.Offset)

	ctx, cancel := context.WithTimeout(ctx, QueryTimeoutDuration)
	defer cancel()

	rows, err := s.db.QueryContext(ctx, query, args...)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var feed []PostWithMetadata

	for rows.Next() {
		var post PostWithMetadata
		if err := rows.Scan(
			&post.ID,
			&post.UserID,
			&post.Title,
			&post.Content,
			&post.CreatedAt,
			&post.Version,
			pq.Array(&post.Tags),
			&post.User.Username,
			&post.CommentsCount,
		); err != nil {
			return nil, err
		}
		feed = append(feed, post)
	}

	if err = rows.Err(); err != nil {
		return nil, err
	}

	return feed, nil
}
