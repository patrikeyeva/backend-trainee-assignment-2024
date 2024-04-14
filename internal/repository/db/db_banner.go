package db

import (
	"avito-banner/internal/domain"
	"context"
	"fmt"
	"log"
	"time"

	"github.com/jackc/pgx/v5/pgxpool"
)

type Repo struct {
	conn *pgxpool.Pool
}

func NewRepository(conn *pgxpool.Pool) *Repo {
	return &Repo{conn: conn}

}

func (r *Repo) GetUserBanner(ctx context.Context, tagId, featureId int) (ResponseUserBanner, error) {
	log.Println("Start Get from db")
	defer log.Println("End Get from db")

	var ret ResponseUserBanner
	query := `SELECT b.title, b.text, b.url, b.is_active
			  FROM Banner b
			  JOIN BannerTag bt ON b.banner_id = bt.banner_id
			  WHERE b.feature_id = $2
	  		  AND bt.tag_id = $1`
	row := r.conn.QueryRow(ctx, query, tagId, featureId)
	err := row.Scan(&ret.Title, &ret.Text, &ret.Url, &ret.IsActive)

	return ret, err
}

func (r *Repo) GetBanner(ctx context.Context, tagId, featureId, limit, offset int) ([]domain.ResponseBanner, error) {
	var ret []domain.ResponseBanner
	// Начальный SQL-запрос без условий
	query := `
	 SELECT b.banner_id, array_agg(bt.tag_id) as tag_ids, b.feature_id, b.title, b.text, b.url, b.is_active, b.created_at, b.updated_at
	 FROM Banner b
	 JOIN BannerTag bt ON b.banner_id = bt.banner_id
	 `

	// Сборка условий WHERE
	whereClause := "WHERE 1=1"

	cntr := 1
	if featureId != -1 {
		whereClause += fmt.Sprintf(" AND b.feature_id = $%d", cntr)
		cntr += 1
	}

	if tagId != -1 {
		whereClause += fmt.Sprintf(" AND bt.tag_id = $%d", cntr)
		cntr += 1
	}

	// Добавление условий WHERE к основному запросу
	query += whereClause

	// Группировка, сортировка
	query += `
		GROUP BY b.banner_id
		ORDER BY b.banner_id
	`

	if limit != -1 {
		query += fmt.Sprintf(" LIMIT $%d", cntr)
		cntr += 1
	}

	if offset != -1 {
		query += fmt.Sprintf(" OFFSET $%d", cntr)
		cntr += 1
	}

	rows, err := r.conn.Query(ctx, query, tagId, featureId)
	if err != nil {
		return ret, err
	}
	for rows.Next() {
		var row domain.ResponseBanner
		if err := rows.Scan(&row.BannerID,
			&row.TagIds,
			&row.FeatureID,
			&row.Content.Title,
			&row.Content.Text,
			&row.Content.Url,
			&row.IsActive,
			&row.CreatedAt,
			&row.UpdatedAt); err != nil {
			return nil, err
		}
		ret = append(ret, row)
	}
	if err := rows.Err(); err != nil {
		return nil, err
	}

	return ret, err
}

func (r *Repo) PostBanner(ctx context.Context, TagIds []int, featureId int, content domain.ResponseUserBanner, isActive bool) (int, error) {
	var bannerID int
	insertBannerQuery := `
        INSERT INTO Banner (feature_id, title, text, url, is_active)
        VALUES ($1, $2, $3, $4, $5)
        RETURNING banner_id
    `

	// Выполнение запроса и получение идентификатора созданной записи в таблице Banner
	err := r.conn.QueryRow(ctx, insertBannerQuery, featureId, content.Title, content.Text, content.Url, isActive).Scan(&bannerID)
	if err != nil {
		return 0, err
	}

	insertBannerTagQuery := `
        INSERT INTO BannerTag (banner_id, tag_id)
        VALUES ($1, $2)
    `

	// Выполнение запросов для каждого тега из списка TagIds
	for _, tagID := range TagIds {
		_, err := r.conn.Exec(ctx, insertBannerTagQuery, bannerID, tagID)
		if err != nil {
			return 0, err
		}
	}

	return bannerID, nil
}

func (r *Repo) PatchBanner(ctx context.Context, bannerId int, TagIds []int, featureId int, content domain.ResponseUserBanner, isActive bool) error {
	updateBannerQuery := `
		UPDATE Banner
		SET feature_id = $1, title = $2, text = $3, url = $4, is_active = $5, updated_at = $6
		WHERE banner_id = $7
	`

	_, err := r.conn.Exec(ctx, updateBannerQuery, featureId, content.Title, content.Text, content.Url, isActive, time.Now(), bannerId)
	if err != nil {
		return err
	}

	deleteBannerTagQuery := `
		DELETE FROM BannerTag
		WHERE banner_id = $1
	`

	_, err = r.conn.Exec(ctx, deleteBannerTagQuery, bannerId)
	if err != nil {
		return err
	}

	insertBannerTagQuery := `
		INSERT INTO BannerTag (banner_id, tag_id)
		VALUES ($1, $2)
	`

	for _, tagID := range TagIds {
		_, err := r.conn.Exec(ctx, insertBannerTagQuery, bannerId, tagID)
		if err != nil {
			return err
		}
	}

	return nil
}

func (r *Repo) DeleteBanner(ctx context.Context, bannerId int) error {
	deleteBannerTagQuery := `
		DELETE FROM BannerTag
		WHERE banner_id = $1
	`

	_, err := r.conn.Exec(ctx, deleteBannerTagQuery, bannerId)
	if err != nil {
		return err
	}

	deleteBannerQuery := `
		DELETE FROM Banner
		WHERE banner_id = $1
	`

	_, err = r.conn.Exec(ctx, deleteBannerQuery, bannerId)
	if err != nil {
		return err
	}

	return nil
}
