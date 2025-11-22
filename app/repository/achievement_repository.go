package repository

import (
	"context"

	"backend/app/models"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo"
)

type IAchievementRepository interface {
	Create(ctx context.Context, achievement *models.Achievement) (string, error)
	FindByID(ctx context.Context, id string) (*models.Achievement, error)
	Update(ctx context.Context, achievement *models.Achievement) error
	Delete(ctx context.Context, id string) error
}

type AchievementRepository struct {
	Collection *mongo.Collection
}

func NewAchievementRepository(db *mongo.Database) IAchievementRepository {
	return &AchievementRepository{
		Collection: db.Collection("achievements"),
	}
}

func (r *AchievementRepository) Create(ctx context.Context, achievement *models.Achievement) (string, error) {
	res, err := r.Collection.InsertOne(ctx, achievement)
	if err != nil {
		return "", err
	}
	return res.InsertedID.(interface{}).(string), nil
}

func (r *AchievementRepository) FindByID(ctx context.Context, id string) (*models.Achievement, error) {
	var a models.Achievement
	err := r.Collection.FindOne(ctx, bson.M{"_id": id}).Decode(&a)
	if err != nil {
		return nil, err
	}
	return &a, nil
}

func (r *AchievementRepository) Update(ctx context.Context, achievement *models.Achievement) error {
	_, err := r.Collection.UpdateOne(
		ctx,
		bson.M{"_id": achievement.ID},
		bson.M{"$set": achievement},
	)
	return err
}

func (r *AchievementRepository) Delete(ctx context.Context, id string) error {
	_, err := r.Collection.DeleteOne(ctx, bson.M{"_id": id})
	return err
}
