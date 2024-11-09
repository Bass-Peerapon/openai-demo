package postgres

import (
	"context"

	"github.com/Bass-Peerapon/openai-demo/domain"
	"github.com/Bass-Peerapon/openai-demo/internal/repository/resty"
	"github.com/jmoiron/sqlx"
	"github.com/pgvector/pgvector-go"
)

type ProductRepository struct {
	client     *sqlx.DB
	openaiRepo *resty.OpenaiService
}

func (r *ProductRepository) SearchProduct(ctx context.Context, query string) ([]domain.Product, error) {
	sql := "SELECT * FROM products ORDER BY embedding <-> $1 LIMIT 5"
	embedding, err := r.openaiRepo.GetEmbedding(ctx, query)
	if err != nil {
		return nil, err
	}
	rows, err := r.client.QueryxContext(ctx, sql, pgvector.NewVector(embedding[0]))
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	products := make([]domain.Product, 0)
	for rows.Next() {
		product := domain.Product{}
		if err := rows.StructScan(&product); err != nil {
			return nil, err
		}
		products = append(products, product)
	}
	return products, nil
}

func (r *ProductRepository) MigrateData() error {
	products := []*domain.Product{
		{Title: "Headphones", Content: "Advanced functionality at an affordable price."},
		{Title: "Webcam", Content: "Designed for comfort and performance."},
		{Title: "Mouse", Content: "Reliable and durable, ideal for both work and leisure."},
		{Title: "Monitor", Content: "Stable and secure, ensuring uninterrupted connectivity."},
		{Title: "Printer", Content: "High quality product with great features."},
		{Title: "Speakers", Content: "Crafted to offer the best audio/visual experience."},
		{Title: "Laptop", Content: "Compact and efficient, perfect for daily use."},
		{Title: "Router", Content: "Stable and secure, ensuring uninterrupted connectivity."},
		{Title: "Smartphone", Content: "A state-of-the-art device with the latest technology."},
		{Title: "Keyboard", Content: "Reliable and durable, ideal for both work and leisure."},
	}

	for _, p := range products {
		embedding, err := r.openaiRepo.GetEmbedding(context.Background(), p.Title)
		if err != nil {
			return err
		}
		if len(embedding) == 0 {
			continue
		}
		p.Embedding = pgvector.NewVector(embedding[0])
	}

	tx, err := r.client.BeginTxx(context.Background(), nil)
	if err != nil {
		return err
	}
	defer tx.Rollback()

	for _, p := range products {
		_, err := tx.NamedExec(`INSERT INTO products (title, content, embedding) VALUES (:title, :content, :embedding)`, p)
		if err != nil {
			return err
		}
	}
	return tx.Commit()
}

func NewProductRepository(client *sqlx.DB, openaiRepo *resty.OpenaiService) *ProductRepository {
	return &ProductRepository{
		client:     client,
		openaiRepo: openaiRepo,
	}
}
