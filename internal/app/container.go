package app

import (
	"database/sql"
	"github.com/aprianimmanuel/rangkaiedu-backend/internal/config"
	"github.com/aprianimmanuel/rangkaiedu-backend/internal/repositories"
	"github.com/aprianimmanuel/rangkaiedu-backend/internal/services"
)

// Container holds all the dependencies for the application
type Container struct {
	// Configuration
	Config *config.Config
	
	// Database
	DB *sql.DB
	
	// Repositories
	UserRepository             repositories.UserRepository
	ClassRepository            repositories.ClassRepository
	MaterialRepository         repositories.MaterialRepository
	StudentEnrollmentRepository repositories.StudentEnrollmentRepository
	SubjectRepository          repositories.SubjectRepository
	
	// Services
	UserService services.UserService
}

// NewContainer creates a new dependency injection container
func NewContainer(cfg *config.Config, db *sql.DB) *Container {
	// Create repositories
	userRepo := repositories.NewUserRepository(db)
	classRepo := repositories.NewClassRepository(db)
	materialRepo := repositories.NewMaterialRepository(db)
	studentEnrollmentRepo := repositories.NewStudentEnrollmentRepository(db)
	subjectRepo := repositories.NewSubjectRepository(db)
	
	// Create services
	userService := services.NewUserService(userRepo)
	
	return &Container{
		Config: cfg,
		DB:     db,
		
		// Repositories
		UserRepository:             userRepo,
		ClassRepository:            classRepo,
		MaterialRepository:         materialRepo,
		StudentEnrollmentRepository: studentEnrollmentRepo,
		SubjectRepository:          subjectRepo,
		
		// Services
		UserService: userService,
	}
}