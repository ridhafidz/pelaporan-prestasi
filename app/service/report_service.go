package service

import (
	"context"
	"errors"
	"time"

	"github.com/google/uuid"

	"backend/app/models"
	"backend/app/repository"
)

type ReportService interface {
	GetStatistics(ctx context.Context, role string, userID uuid.UUID, start time.Time, end time.Time) (*models.AchievementStatisticsResponse, error)
	GetStudentReport(ctx context.Context, requesterRole string, requesterID uuid.UUID, studentID uuid.UUID) (*models.StudentReportResponse, error)
}

type reportService struct {
	reportRepo  repository.ReportRepository
	studentRepo repository.StudentLecturerRepository
}

func NewReportService(
	reportRepo repository.ReportRepository,
	studentRepo repository.StudentLecturerRepository,
) ReportService {
	return &reportService{
		reportRepo:  reportRepo,
		studentRepo: studentRepo,
	}
}

func (s *reportService) GetStatistics(
	ctx context.Context,
	role string,
	userID uuid.UUID,
	start time.Time,
	end time.Time,
) (*models.AchievementStatisticsResponse, error) { // Tambah prefix models.

	// FR-011: Admin → Full Access [cite: 30, 253]
	if role == "Admin" {
		return s.getGlobalStatistics(ctx, start, end)
	}

	// FR-011: Mahasiswa → Only Own [cite: 30, 253]
	if role == "Mahasiswa" {
		// Pastikan method ini ada di repository.StudentLecturerRepository
		student, err := s.studentRepo.GetStudentByID(ctx, userID)
		if err != nil {
			return nil, err
		}

		point, err := s.reportRepo.GetStudentTotalPoint(ctx, student.ID)
		if err != nil {
			return nil, err
		}

		return &models.AchievementStatisticsResponse{
			TopStudents: []models.TopStudentStat{
				{
					StudentID:  student.ID,
					FullName:   student.FullName,
					TotalPoint: point,
				},
			},
		}, nil
	}

	// FR-011: Dosen Wali → Advisees [cite: 30, 253]
	// FR-011: Dosen Wali → Advisees [cite: 253]
	if role == "Dosen Wali" {
		// 1. Kita harus tau dulu lecturerID si dosen ini dari userID-nya
		// Pastikan lo punya fungsi GetLecturerByUserID di repository
		lecturer, err := s.studentRepo.GetLecturerByUserID(ctx, userID)
		if err != nil {
			return nil, err
		}

		// 2. Baru panggil GetLecturerAdvisees pake lecturer.ID [cite: 208]
		students, err := s.studentRepo.GetLecturerAdvisees(ctx, lecturer.ID)
		if err != nil {
			return nil, err
		}

		var top []models.TopStudentStat
		var studentIDs []uuid.UUID
		for _, st := range students {
			studentIDs = append(studentIDs, st.ID)
			point, _ := s.reportRepo.GetStudentTotalPoint(ctx, st.ID)
			top = append(top, models.TopStudentStat{
				StudentID:  st.ID,
				FullName:   st.FullName,
				TotalPoint: point,
			})
		}

		// 3. Ambil statistik ter-filter [cite: 254]
		byType, _ := s.reportRepo.GetCountByTypeFiltered(ctx, studentIDs)
		levels, _ := s.reportRepo.GetLevelDistributionFiltered(ctx, studentIDs)

		return &models.AchievementStatisticsResponse{
			ByType:            byType,
			TopStudents:       top,
			CompetitionLevels: levels,
		}, nil
	}

	return nil, errors.New("invalid role")
}

func (s *reportService) getGlobalStatistics(
	ctx context.Context,
	start time.Time,
	end time.Time,
) (*models.AchievementStatisticsResponse, error) { // Tambah prefix models.

	byType, err := s.reportRepo.GetAchievementCountByType(ctx)
	if err != nil {
		return nil, err
	}

	// Sesuai FR-011 Output: Total per periode & Distribusi level [cite: 254]
	byPeriod, err := s.reportRepo.GetAchievementCountByPeriod(ctx, start, end)
	if err != nil {
		return nil, err
	}

	topStudents, err := s.reportRepo.GetTopStudents(ctx, 10)
	if err != nil {
		return nil, err
	}

	levels, err := s.reportRepo.GetCompetitionLevelDistribution(ctx)
	if err != nil {
		return nil, err
	}

	return &models.AchievementStatisticsResponse{
		ByType:            byType,
		ByPeriod:          byPeriod,
		TopStudents:       topStudents,
		CompetitionLevels: levels,
	}, nil
}

func (s *reportService) GetStudentReport(
	ctx context.Context,
	requesterRole string,
	requesterID uuid.UUID,
	studentID uuid.UUID,
) (*models.StudentReportResponse, error) { // Tambah prefix models.

	// Mahasiswa hanya boleh lihat milik sendiri [cite: 253]
	if requesterRole == "Mahasiswa" {
		student, err := s.studentRepo.GetStudentByID(ctx, requesterID)
		if err != nil || student.ID != studentID {
			return nil, errors.New("access denied")
		}
	}

	// Dosen wali harus pembimbingnya [cite: 208, 253]
	if requesterRole == "Dosen Wali" {
		ok, err := s.studentRepo.IsAdvisorOf(ctx, requesterID, studentID)
		if err != nil || !ok {
			return nil, errors.New("access denied")
		}
	}

	point, err := s.reportRepo.GetStudentTotalPoint(ctx, studentID)
	if err != nil {
		return nil, err
	}

	return &models.StudentReportResponse{
		StudentID:  studentID,
		TotalPoint: point,
	}, nil
}
