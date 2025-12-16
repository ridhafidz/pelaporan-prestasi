package repository

import (
	"context"
	"errors"
	"time"

	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"

	"backend/app/models"
)

// AchievementRepository defines operations for achievements stored in MongoDB
type AchievementRepository interface {
	Create(ctx context.Context, a *models.Achievement) error
	GetByID(ctx context.Context, id string) (*models.Achievement, error)
	GetByStudentID(ctx context.Context, studentID string) ([]models.Achievement, error)
	Update(ctx context.Context, id string, update bson.M) error
	SoftDelete(ctx context.Context, id string) error
	AddAttachment(ctx context.Context, id string, att models.Attachment) error
	RemoveAttachmentByFileName(ctx context.Context, id string, fileName string) error
	FindByTag(ctx context.Context, tag string, limit int64) ([]models.Achievement, error)
}

type achievementRepo struct {
	coll *mongo.Collection
}

// NewAchievementRepository creates repository using a mongo client and database name
func NewAchievementRepository(client *mongo.Client, dbName string) AchievementRepository {
	coll := client.Database(dbName).Collection("achievements")
	return &achievementRepo{coll: coll}
}

func (r *achievementRepo) Create(ctx context.Context, a *models.Achievement) error {
	if a == nil {
		return errors.New("achievement is nil")
	}
	now := time.Now()
	if a.ID.IsZero() {
		a.ID = primitive.NewObjectID()
	}
	if a.CreatedAt.IsZero() {
		a.CreatedAt = now
	}
	a.UpdatedAt = now

	_, err := r.coll.InsertOne(ctx, a)
	return err
}

func (r *achievementRepo) GetByID(ctx context.Context, id string) (*models.Achievement, error) {
	oid, err := primitive.ObjectIDFromHex(id)
	if err != nil {
		return nil, err
	}
	var a models.Achievement
	err = r.coll.FindOne(ctx, bson.M{"_id": oid}).Decode(&a)
	if err != nil {
		if err == mongo.ErrNoDocuments {
			return nil, nil
		}
		return nil, err
	}
	return &a, nil
}

func (r *achievementRepo) GetByStudentID(ctx context.Context, studentID string) ([]models.Achievement, error) {
	filter := bson.M{"studentId": studentID}
	opts := options.Find().SetSort(bson.D{{"createdAt", -1}})
	cur, err := r.coll.Find(ctx, filter, opts)
	if err != nil {
		return nil, err
	}
	defer cur.Close(ctx)

	var results []models.Achievement
	for cur.Next(ctx) {
		var a models.Achievement
		if err := cur.Decode(&a); err != nil {
			return nil, err
		}
		results = append(results, a)
	}
	if err := cur.Err(); err != nil {
		return nil, err
	}
	return results, nil
}

func (r *achievementRepo) Update(ctx context.Context, id string, update bson.M) error {
	oid, err := primitive.ObjectIDFromHex(id)
	if err != nil {
		return err
	}
	if update == nil {
		return errors.New("update is empty")
	}
	// always update updatedAt
	if update["$set"] == nil {
		update["$set"] = bson.M{"updatedAt": time.Now()}
	} else {
		update["$set"].(bson.M)["updatedAt"] = time.Now()
	}

	res, err := r.coll.UpdateByID(ctx, oid, update)
	if err != nil {
		return err
	}
	if res.MatchedCount == 0 {
		return mongo.ErrNoDocuments
	}
	return nil
}

func (r *achievementRepo) SoftDelete(ctx context.Context, id string) error {
	oid, err := primitive.ObjectIDFromHex(id)
	if err != nil {
		return err
	}
	now := time.Now()
	res, err := r.coll.UpdateByID(ctx, oid, bson.M{"$set": bson.M{"deletedAt": now, "updatedAt": now}})
	if err != nil {
		return err
	}
	if res.MatchedCount == 0 {
		return mongo.ErrNoDocuments
	}
	return nil
}

func (r *achievementRepo) AddAttachment(ctx context.Context, id string, att models.Attachment) error {
	oid, err := primitive.ObjectIDFromHex(id)
	if err != nil {
		return err
	}
	att.UploadedAt = time.Now()
	res, err := r.coll.UpdateByID(ctx, oid, bson.M{"$push": bson.M{"attachments": att}, "$set": bson.M{"updatedAt": time.Now()}})
	if err != nil {
		return err
	}
	if res.MatchedCount == 0 {
		return mongo.ErrNoDocuments
	}
	return nil
}

func (r *achievementRepo) RemoveAttachmentByFileName(ctx context.Context, id string, fileName string) error {
	oid, err := primitive.ObjectIDFromHex(id)
	if err != nil {
		return err
	}
	res, err := r.coll.UpdateByID(ctx, oid, bson.M{"$pull": bson.M{"attachments": bson.M{"fileName": fileName}}, "$set": bson.M{"updatedAt": time.Now()}})
	if err != nil {
		return err
	}
	if res.MatchedCount == 0 {
		return mongo.ErrNoDocuments
	}
	return nil
}

func (r *achievementRepo) FindByTag(ctx context.Context, tag string, limit int64) ([]models.Achievement, error) {
	filter := bson.M{"tags": tag}
	opts := options.Find().SetSort(bson.D{{"createdAt", -1}})
	if limit > 0 {
		opts.SetLimit(limit)
	}
	cur, err := r.coll.Find(ctx, filter, opts)
	if err != nil {
		return nil, err
	}
	defer cur.Close(ctx)

	var res []models.Achievement
	for cur.Next(ctx) {
		var a models.Achievement
		if err := cur.Decode(&a); err != nil {
			return nil, err
		}
		res = append(res, a)
	}
	if err := cur.Err(); err != nil {
		return nil, err
	}
	return res, nil
}
