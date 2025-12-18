package repository

import (
	"context"
	"database/sql"
	"time"

	"backend/app/models"

	"github.com/google/uuid"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo"
)

type ReportRepository interface {
	GetAchievementCountByType(ctx context.Context) ([]models.AchievementTypeStat, error)
	GetAchievementCountByPeriod(ctx context.Context, start time.Time, end time.Time) ([]models.AchievementPeriodStat, error)
	GetTopStudents(ctx context.Context, limit int) ([]models.TopStudentStat, error)
	GetCompetitionLevelDistribution(ctx context.Context) ([]models.CompetitionLevelStat, error)
	GetStudentTotalPoint(ctx context.Context, studentID uuid.UUID) (float64, error)
	// Buat Filtered Version untuk Dosen Wali (FR-011)
	GetCountByTypeFiltered(ctx context.Context, studentIDs []uuid.UUID) ([]models.AchievementTypeStat, error)
	GetLevelDistributionFiltered(ctx context.Context, studentIDs []uuid.UUID) ([]models.CompetitionLevelStat, error)
}

type reportRepository struct {
	db         *sql.DB           // Untuk Join Status & User
	collection *mongo.Collection // Untuk Aggregation Data Dinamis
}

func NewReportRepository(db *sql.DB, mongoDb *mongo.Database) ReportRepository {
	return &reportRepository{
		db:         db,
		collection: mongoDb.Collection("achievements"),
	}
}

// 1. Total prestasi per tipe (FR-011)
func (r *reportRepository) GetAchievementCountByType(ctx context.Context) ([]models.AchievementTypeStat, error) {
	pipeline := mongo.Pipeline{
		{{Key: "$match", Value: bson.M{"is_deleted": bson.M{"$ne": true}}}},
		{{Key: "$group", Value: bson.M{"_id": "$achievementType", "total": bson.M{"$sum": 1}}}},
		{{Key: "$project", Value: bson.M{"achievementType": "$_id", "total": 1, "_id": 0}}},
	}
	cursor, err := r.collection.Aggregate(ctx, pipeline)
	if err != nil {
		return nil, err
	}
	var result []models.AchievementTypeStat
	err = cursor.All(ctx, &result)
	return result, err
}

// 2. Total prestasi per periode (FR-011)
func (r *reportRepository) GetAchievementCountByPeriod(ctx context.Context, start time.Time, end time.Time) ([]models.AchievementPeriodStat, error) {
	pipeline := mongo.Pipeline{
		{{Key: "$match", Value: bson.M{
			"createdAt":  bson.M{"$gte": start, "$lte": end},
			"is_deleted": bson.M{"$ne": true},
		}}},
		{{Key: "$group", Value: bson.M{
			"_id":   bson.M{"$dateToString": bson.M{"format": "%Y-%m", "date": "$createdAt"}},
			"total": bson.M{"$sum": 1},
		}}},
		{{Key: "$sort", Value: bson.M{"_id": 1}}},
		{{Key: "$project", Value: bson.M{"period": "$_id", "total": 1, "_id": 0}}},
	}
	cursor, err := r.collection.Aggregate(ctx, pipeline)
	if err != nil {
		return nil, err
	}
	var result []models.AchievementPeriodStat
	err = cursor.All(ctx, &result)
	return result, err
}

// 3. Top mahasiswa berprestasi (FR-011) - Kawinin data Mongo & Postgres
func (r *reportRepository) GetTopStudents(ctx context.Context, limit int) ([]models.TopStudentStat, error) {
	pipeline := mongo.Pipeline{
		{{Key: "$match", Value: bson.M{"is_deleted": bson.M{"$ne": true}}}},
		{{Key: "$group", Value: bson.M{"_id": "$studentId", "totalPoint": bson.M{"$sum": "$points"}}}},
		{{Key: "$sort", Value: bson.M{"totalPoint": -1}}},
		{{Key: "$limit", Value: limit}},
	}
	cursor, err := r.collection.Aggregate(ctx, pipeline)
	if err != nil {
		return nil, err
	}

	var results []models.TopStudentStat
	for cursor.Next(ctx) {
		var temp struct {
			ID         string  `bson:"_id"`
			TotalPoint float64 `bson:"totalPoint"`
		}
		cursor.Decode(&temp)

		// Ambil Nama dari Postgres
		var fullName string
		query := `SELECT u.full_name FROM students s JOIN users u ON s.user_id = u.id WHERE s.id = $1`
		r.db.QueryRowContext(ctx, query, temp.ID).Scan(&fullName)

		results = append(results, models.TopStudentStat{
			StudentID:  uuid.MustParse(temp.ID),
			FullName:   fullName,
			TotalPoint: temp.TotalPoint,
		})
	}
	return results, nil
}

// 4. Distribusi tingkat kompetisi (FR-011)
func (r *reportRepository) GetCompetitionLevelDistribution(ctx context.Context) ([]models.CompetitionLevelStat, error) {
	pipeline := mongo.Pipeline{
		{{Key: "$match", Value: bson.M{"achievementType": "competition", "is_deleted": bson.M{"$ne": true}}}},
		{{Key: "$group", Value: bson.M{"_id": "$details.competitionLevel", "total": bson.M{"$sum": 1}}}},
		{{Key: "$project", Value: bson.M{"level": "$_id", "total": 1, "_id": 0}}},
	}
	cursor, err := r.collection.Aggregate(ctx, pipeline)
	if err != nil {
		return nil, err
	}
	var result []models.CompetitionLevelStat
	err = cursor.All(ctx, &result)
	return result, err
}

// 5. Total Poin per Mahasiswa
func (r *reportRepository) GetStudentTotalPoint(ctx context.Context, studentID uuid.UUID) (float64, error) {
	pipeline := mongo.Pipeline{
		{{Key: "$match", Value: bson.M{"studentId": studentID.String(), "is_deleted": bson.M{"$ne": true}}}},
		{{Key: "$group", Value: bson.M{"_id": nil, "total": bson.M{"$sum": "$points"}}}},
	}
	cursor, err := r.collection.Aggregate(ctx, pipeline)
	if err != nil {
		return 0, err
	}
	var res []struct{ Total float64 `bson:"total"` }
	cursor.All(ctx, &res)
	if len(res) > 0 {
		return res[0].Total, nil
	}
	return 0, nil
}

// Filtered Version untuk Dosen Wali (Bimbingan)
func (r *reportRepository) GetCountByTypeFiltered(ctx context.Context, studentIDs []uuid.UUID) ([]models.AchievementTypeStat, error) {
	ids := make([]string, len(studentIDs))
	for i, id := range studentIDs { ids[i] = id.String() }

	pipeline := mongo.Pipeline{
		{{Key: "$match", Value: bson.M{"studentId": bson.M{"$in": ids}, "is_deleted": bson.M{"$ne": true}}}},
		{{Key: "$group", Value: bson.M{"_id": "$achievementType", "total": bson.M{"$sum": 1}}}},
		{{Key: "$project", Value: bson.M{"achievementType": "$_id", "total": 1, "_id": 0}}},
	}
	cursor, err := r.collection.Aggregate(ctx, pipeline)
	if err != nil { return nil, err }
	var result []models.AchievementTypeStat
	cursor.All(ctx, &result)
	return result, nil
}

func (r *reportRepository) GetLevelDistributionFiltered(ctx context.Context, studentIDs []uuid.UUID) ([]models.CompetitionLevelStat, error) {
	ids := make([]string, len(studentIDs))
	for i, id := range studentIDs { ids[i] = id.String() }

	pipeline := mongo.Pipeline{
		{{Key: "$match", Value: bson.M{"studentId": bson.M{"$in": ids}, "achievementType": "competition", "is_deleted": bson.M{"$ne": true}}}},
		{{Key: "$group", Value: bson.M{"_id": "$details.competitionLevel", "total": bson.M{"$sum": 1}}}},
		{{Key: "$project", Value: bson.M{"level": "$_id", "total": 1, "_id": 0}}},
	}
	cursor, err := r.collection.Aggregate(ctx, pipeline)
	if err != nil { return nil, err }
	var result []models.CompetitionLevelStat
	cursor.All(ctx, &result)
	return result, nil
}