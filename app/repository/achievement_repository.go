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
}

type achievementRepository struct {
	collection *mongo.Collection
}

func NewAchievementRepository(db *mongo.Database) AchievementRepository {
	return &achievementRepository{
		collection: db.Collection("achievements"),
	}
}

// Create: Menyimpan data detail ke Mongo dan mengembalikan ID string untuk disimpan di Postgres
func (r *achievementRepository) Create(ctx context.Context, achievement *models.Achievement) (string, error) {
	achievement.ID = primitive.NewObjectID()
	achievement.CreatedAt = time.Now()
	achievement.UpdatedAt = time.Now()

	// Poin default 0 jika tidak diisi (sesuai SRS field points: Number) [cite: 154]
	if achievement.Points == 0 {
		achievement.Points = 0
	}

	_, err := r.collection.InsertOne(ctx, achievement)
	if err != nil {
		return "", err
	}
	// Return Hex String ID untuk disimpan di tabel Postgres 'mongo_achievement_id'
	return achievement.ID.Hex(), nil
}

// FindByID: Mengambil detail data
func (r *achievementRepository) FindByID(ctx context.Context, id string) (*models.Achievement, error) {
	objID, _ := primitive.ObjectIDFromHex(id)
	filter := bson.M{"_id": objID}

	var achievement models.Achievement
	err := r.collection.FindOne(ctx, filter).Decode(&achievement)
	return &achievement, err
}

// AddAttachment: Upload files [cite: 281]
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
