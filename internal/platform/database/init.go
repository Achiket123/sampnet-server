package database

import (
	"fmt"
	"log"
	"os"
	"server/internal/platform/database/models"

	"gorm.io/driver/postgres"
	"gorm.io/gorm"
	"gorm.io/gorm/logger"
)

var DB *gorm.DB

func Init() {
	DB = initDB()
}

func initDB() *gorm.DB {
	host := os.Getenv("DB_HOST")
	user := os.Getenv("DB_USER")
	password := os.Getenv("DB_PASSWORD")
	dbname := os.Getenv("DB_NAME")
	port := os.Getenv("DB_PORT")

	if host == "" || user == "" || password == "" || dbname == "" || port == "" {
		log.Fatalf("Database configuration missing. Please check environment variables.")
	}

	dsn := fmt.Sprintf("host=%s user=%s password=%s dbname=%s port=%s sslmode=disable", host, user, password, dbname, port)

	db, err := gorm.Open(postgres.Open(dsn), &gorm.Config{
		Logger: logger.Default.LogMode(logger.Info),
	})
	if err != nil {
		panic("failed to connect database")
	}
	return db
}

func PerformMigrations() error {
	err := DB.AutoMigrate(
		&models.Organisation{},
		&models.UserModel{},
		&models.Task{},
		&models.File{},
		&models.Employee{},
		&models.Manager{},
		&models.Boss{},
		&models.Attendance{},
		&models.Project{},
		&models.Team{},
		&models.TeamMember{},
		&models.Notification{},
		&models.Chat{},
		&models.ChatMessage{},
		&models.CallState{},
		&models.TaskComment{},
		&models.TaskAttachment{},
		&models.OnboardingProgress{},
		&models.Leave{},
		&models.Milestone{},
		&models.LeavePolicy{},
		&models.AttendancePolicy{},
		&models.RolePermissions{},
		&models.TaskType{},
		&models.TaskActivity{},
		&models.WorkSchedule{},
		&models.AuditLog{},
		&models.Invite{},
		&models.ResourceCollection{},
		&models.ResourceRecord{},
		&models.ResourceRecordHistory{},
		&models.ResourceRecordAttachment{},
		&models.ResearchEntry{},
		&models.ResearchFolder{},
		&models.ResearchDocument{},
		&models.ResearchDocumentVersion{},
		&models.ResearchFile{},
		&models.ResearchReference{},
		&models.ResearchCollaborator{},
		&models.ResearchActivity{},
		&models.PeopleContact{},
		&models.PeopleInteraction{},
		&models.PeoplePipelineStage{},
		&models.PeopleList{},
		&models.PeopleListContact{},
	)
	if err != nil {
		return err
	}

	// Create GIN index and search helper indexes for Resource Collection and Records
	DB.Exec(`CREATE INDEX IF NOT EXISTS idx_resource_records_data ON resource_records USING GIN (data)`)
	DB.Exec(`CREATE INDEX IF NOT EXISTS idx_resource_records_collection_org ON resource_records (collection_id, organisation_id)`)
	DB.Exec(`CREATE INDEX IF NOT EXISTS idx_resource_collections_org ON resource_collections (organisation_id)`)

	return nil
}
