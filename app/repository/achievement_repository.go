package repository

import (
	"context"
	"time"

	"backend/app/models"

	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"
)

type AchievementRepository interface {
	Create(ctx context.Context, achievement *models.Achievement) (string, error)
	FindByID(ctx context.Context, id string) (*models.Achievement, error)
	AddAttachment(ctx context.Context, id string, attachment models.Attachment) error
	UpdatePoints(ctx context.Context, id string, points float64) error
	Update(ctx context.Context, id string, achievement *models.Achievement) error
	SoftDelete(ctx context.Context, id string) error
}

type achievementRepository struct {
	collection *mongo.Collection
}

func NewAchievementRepository(db *mongo.Database) AchievementRepository {
	return &achievementRepository{
		collection: db.Collection("achievements"),
	}
}

func (r *achievementRepository) Create(ctx context.Context, achievement *models.Achievement) (string, error) {
	achievement.ID = primitive.NewObjectID()
	achievement.CreatedAt = time.Now()
	achievement.UpdatedAt = time.Now()

	if achievement.Points == 0 {
		achievement.Points = 0
	}

	_, err := r.collection.InsertOne(ctx, achievement)
	if err != nil {
		return "", err
	}
	return achievement.ID.Hex(), nil
}

func (r *achievementRepository) FindByID(ctx context.Context, id string) (*models.Achievement, error) {
	objID, _ := primitive.ObjectIDFromHex(id)
	filter := bson.M{
		"_id":        objID,
		"is_deleted": bson.M{"$ne": true},
	}

	var achievement models.Achievement
	err := r.collection.FindOne(ctx, filter).Decode(&achievement)
	return &achievement, err
}

func (r *achievementRepository) AddAttachment(ctx context.Context, id string, attachment models.Attachment) error {
	objID, _ := primitive.ObjectIDFromHex(id)
	filter := bson.M{"_id": objID}
	update := bson.M{
		"$push": bson.M{"attachments": attachment},
		"$set":  bson.M{"updatedAt": time.Now()},
	}
	_, err := r.collection.UpdateOne(ctx, filter, update)
	return err
}

func (r *achievementRepository) UpdatePoints(ctx context.Context, id string, points float64) error {
	objID, err := primitive.ObjectIDFromHex(id)
	if err != nil {
		return err
	}

	update := bson.M{
		"$set": bson.M{
			"points":    points,
			"updatedAt": time.Now(),
		},
	}

	_, err = r.collection.UpdateOne(ctx, bson.M{"_id": objID}, update)
	return err
}

func (r *achievementRepository) Update(
	ctx context.Context,
	id string,
	achievement *models.Achievement,
) error {

	objID, err := primitive.ObjectIDFromHex(id)
	if err != nil {
		return err
	}

	filter := bson.M{"_id": objID}

	update := bson.M{
		"$set": bson.M{
			"title":       achievement.Title,
			"description": achievement.Description,
			"details":     achievement.Details,
			"tags":        achievement.Tags,
			"points":      achievement.Points,
			"updatedAt":   achievement.UpdatedAt,
		},
	}

	_, err = r.collection.UpdateOne(ctx, filter, update)
	return err
}

func (r *achievementRepository) SoftDelete(
	ctx context.Context,
	id string,
) error {

	objID, err := primitive.ObjectIDFromHex(id)
	if err != nil {
		return err
	}

	now := time.Now()

	update := bson.M{
		"$set": bson.M{
			"is_deleted": true,
			"deletedAt":  now,
			"updatedAt":  now,
		},
	}

	_, err = r.collection.UpdateOne(ctx, bson.M{"_id": objID}, update)
	return err
}
