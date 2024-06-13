package links

import (
	"context"
	"fmt"
	"time"

	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"

	"gitlab.com/robotomize/gb-golang/homework/03-01-umanager/internal/database"
)

const collection = "links"

func New(db *mongo.Database, timeout time.Duration) *Repository {
	return &Repository{db: db, timeout: timeout}
}

type Repository struct {
	db      *mongo.Database
	timeout time.Duration
}

func (r *Repository) Create(ctx context.Context, req CreateReq) (database.Link, error) {
	var l database.Link

	ctx, cancel := context.WithTimeout(ctx, r.timeout)
	defer cancel()

	req.CreatedAt = time.Now()
	req.UpdatedAt = time.Now()

	result, err := r.db.Collection(collection).InsertOne(ctx, req)
	if err != nil {
		return database.Link{}, fmt.Errorf("failed to create link: %w", err)
	}

	filter := bson.M{"_id": result.InsertedID}
	err = r.db.Collection(collection).FindOne(ctx, filter).Decode(&l)
	if err != nil {
		return database.Link{}, fmt.Errorf("failed to decode created link: %w", err)
	}

	return l, nil
}

func (r *Repository) FindByUserAndURL(ctx context.Context, link, userID string) (database.Link, error) {
	var l database.Link

	ctx, cancel := context.WithTimeout(ctx, r.timeout)
	defer cancel()

	filter := bson.M{"url": url, "userID": userID}

	err := r.db.Collection(collection).FindOne(ctx, filter).Decode(&l)
	if err != nil {
		return database.Link{}, fmt.Errorf("failed to find link by user and URL: %w", err)
	}

	return l, nil
}

func (r *Repository) FindByCriteria(ctx context.Context, criteria Criteria) ([]database.Link, error) {
	var links []database.Link

	ctx, cancel := context.WithTimeout(ctx, r.timeout)
	defer cancel()

	filter := bson.M{}
	if criteria.UserID != nil {
		filter["userID"] = *criteria.UserID
	}
	if len(criteria.Tags) > 0 {
		filter["tags"] = bson.M{"$in": criteria.Tags}
	}

	findOptions := options.Find()
	if criteria.Limit != nil {
		findOptions.SetLimit(*criteria.Limit)
	}
	if criteria.Offset != nil {
		findOptions.SetSkip(*criteria.Offset)
	}

	cursor, err := r.db.Collection(collection).Find(ctx, filter, findOptions)
	if err != nil {
		return nil, fmt.Errorf("failed to find links by criteria: %w", err)
	}
	defer cursor.Close(ctx)

	err = cursor.All(ctx, &links)
	if err != nil {
		return nil, fmt.Errorf("failed to decode links: %w", err)
	}

	return links, nil
}
