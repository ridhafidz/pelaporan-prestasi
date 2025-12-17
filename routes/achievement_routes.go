package routes

import (
	"fmt"
	"path/filepath"
	"time"

	"github.com/gofiber/fiber/v2"
	"github.com/google/uuid"

	"backend/app/models"
	"backend/app/services"
)

func processCreateAchievement(s service.AchievementService) fiber.Handler {
	return func(c *fiber.Ctx) error {
		userIDVal := c.Locals("userID")
		if userIDVal == nil {
			return c.Status(fiber.StatusUnauthorized).JSON(fiber.Map{"status": "error", "message": "Unauthorized"})
		}
		studentID := userIDVal.(uuid.UUID)

		req := new(models.CreateAchievementReferenceRequest)
		if err := c.BodyParser(req); err != nil {
			return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{"status": "error", "message": "Invalid request body"})
		}

		req.StudentID = studentID

		if err := validate.Struct(req); err != nil {
			return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{"status": "error", "message": err.Error()})
		}

		resp, err := s.Create(c.Context(), req)
		if err != nil {
			return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{"status": "error", "message": err.Error()})
		}

		return c.Status(fiber.StatusCreated).JSON(fiber.Map{"status": "success", "data": resp})
	}
}

func processGetAchievements(s service.AchievementService) fiber.Handler {
	return func(c *fiber.Ctx) error {
		userIDVal := c.Locals("userID")
		roleVal := c.Locals("role")
		
		if userIDVal == nil {
			return c.Status(fiber.StatusUnauthorized).JSON(fiber.Map{"status": "error", "message": "Unauthorized"})
		}
		userID := userIDVal.(uuid.UUID)
		role := roleVal.(string)

		page := c.QueryInt("page", 1)
		limit := c.QueryInt("limit", 10)

		var achievements []models.Achievement
		var err error

		// Logic Filter by Role
		if role == "Mahasiswa" {
			// Mahasiswa hanya melihat miliknya sendiri
			achievements, err = s.GetByStudentID(c.Context(), userID, page, limit)
		} else {
			// Dosen/Admin melihat semua (atau logic specific Dosen Wali bisa ditambahkan di sini)
			// Untuk sekarang kita pakai GetByStudentID dulu jika dosen ingin lihat spesifik mahasiswa, 
			// atau Anda perlu buat method `GetAll(ctx, page, limit)` di service untuk Admin.
			// Contoh Mock response untuk non-mahasiswa sementara:
			return c.Status(fiber.StatusNotImplemented).JSON(fiber.Map{"message": "View for Lecturer/Admin not implemented yet in service"})
		}

		if err != nil {
			return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{"status": "error", "message": err.Error()})
		}

		return c.Status(fiber.StatusOK).JSON(fiber.Map{
			"status": "success",
			"data":   achievements,
			"meta":   fiber.Map{"page": page, "limit": limit},
		})
	}
}

func processGetAchievementByID(s service.AchievementReferenceService) fiber.Handler {
	return func(c *fiber.Ctx) error {
		idParam := c.Params("id")
		id, err := uuid.Parse(idParam)
		if err != nil {
			return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{"status": "error", "message": "Invalid UUID"})
		}

		achievement, err := s.GetByID(c.Context(), id)
		if err != nil {
			return c.Status(fiber.StatusNotFound).JSON(fiber.Map{"status": "fail", "message": "Achievement not found"})
		}

		return c.Status(fiber.StatusOK).JSON(fiber.Map{"status": "success", "data": achievement})
	}
}

func processSubmitAchievement(s service.AchievementService) fiber.Handler {
	return func(c *fiber.Ctx) error {
		idParam := c.Params("id")
		id, err := uuid.Parse(idParam)
		if err != nil {
			return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{"status": "error", "message": "Invalid UUID"})
		}

		if err := s.Submit(c.Context(), id); err != nil {
			return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{"status": "error", "message": err.Error()})
		}

		return c.Status(fiber.StatusOK).JSON(fiber.Map{"status": "success", "message": "Achievement submitted successfully"})
	}
}

func processVerifyAchievement(s service.AchievementService) fiber.Handler {
	return func(c *fiber.Ctx) error {
		// Verifier (Dosen) ID
		userIDVal := c.Locals("userID")
		verifierID := userIDVal.(uuid.UUID)

		// Target Achievement ID
		idParam := c.Params("id")
		id, err := uuid.Parse(idParam)
		if err != nil {
			return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{"status": "error", "message": "Invalid UUID"})
		}

		if err := s.Verify(c.Context(), id, verifierID); err != nil {
			return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{"status": "error", "message": err.Error()})
		}

		return c.Status(fiber.StatusOK).JSON(fiber.Map{"status": "success", "message": "Achievement verified successfully"})
	}
}

func processRejectAchievement(s service.AchievementService) fiber.Handler {
	return func(c *fiber.Ctx) error {
		idParam := c.Params("id")
		id, err := uuid.Parse(idParam)
		if err != nil {
			return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{"status": "error", "message": "Invalid UUID"})
		}

		req := new(models.RejectAchievementRequest)
		if err := c.BodyParser(req); err != nil {
			return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{"status": "error", "message": "Invalid body"})
		}

		if err := s.Reject(c.Context(), id, req); err != nil {
			return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{"status": "error", "message": err.Error()})
		}

		return c.Status(fiber.StatusOK).JSON(fiber.Map{"status": "success", "message": "Achievement rejected"})
	}
}

func processAddAttachment(s service.AchievementService) fiber.Handler {
	return func(c *fiber.Ctx) error {
		idParam := c.Params("id")
		id, err := uuid.Parse(idParam)
		if err != nil {
			return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{"status": "error", "message": "Invalid UUID"})
		}

		// 1. Handle File Upload (Fiber)
		file, err := c.FormFile("file")
		if err != nil {
			return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{"status": "error", "message": "File is required"})
		}

		// 2. Simpan File (Contoh sederhana: simpan lokal)
		// Di production, upload ke S3/GCS lalu dapatkan URL-nya
		filename := fmt.Sprintf("%d_%s", time.Now().Unix(), file.Filename)
		filePath := filepath.Join("./uploads", filename)
		
		if err := c.SaveFile(file, filePath); err != nil {
			return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{"status": "error", "message": "Failed to save file"})
		}

		// 3. Buat Object Attachment
		attachment := models.Attachment{
			FileName:   file.Filename,
			FileURL:    "/uploads/" + filename, // URL akses publik
			FileType:   filepath.Ext(filename),
			UploadedAt: time.Now(),
		}

		// 4. Call Service
		if err := s.AddAttachment(c.Context(), id, attachment); err != nil {
			return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{"status": "error", "message": err.Error()})
		}

		return c.Status(fiber.StatusOK).JSON(fiber.Map{"status": "success", "data": attachment})
	}
}

// --- SETUP ROUTES ---

// Tambahkan authService jika middleware membutuhkannya, atau cukup UserService untuk profile
func SetupAchievementRoutes(api fiber.Router, achievementService service.AchievementService) {
	
	achievements := api.Group("/achievements")

	// Middleware (Auth Check) sudah harus terpasang di level 'api' atau group induknya di routes.go
	// Jika belum, tambahkan: achievements.Use(middleware.JWTMiddleware())

	// Endpoint List sesuai Gambar
	achievements.Get("/", processGetAchievements(achievementService))
	achievements.Post("/", processCreateAchievement(achievementService))
	achievements.Get("/:id", processGetAchievementByID(achievementService))
	
	// Workflow Endpoints
	achievements.Post("/:id/submit", processSubmitAchievement(achievementService))
	achievements.Post("/:id/verify", processVerifyAchievement(achievementService)) // Sebaiknya tambah middleware OnlyLecturer
	achievements.Post("/:id/reject", processRejectAchievement(achievementService)) // Sebaiknya tambah middleware OnlyLecturer
	
	achievements.Post("/:id/attachments", processAddAttachment(achievementService))

	// Note: Endpoint Update (PUT) dan Delete (DELETE) bisa ditambahkan serupa dengan pola di atas
	// achievements.Put("/:id", processUpdateAchievement(achievementService))
	// achievements.Delete("/:id", processDeleteAchievement(achievementService))
}