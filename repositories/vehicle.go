package repositories

import (
	"fmt"
	"time"

	"gorm.io/gorm"
)

type Vehicle struct {
	ID        string    `json:"id"`
	UserID    uint      `json:"user_id"`           // Foreign key field
	User      User      `gorm:"foreignKey:UserID"` // Association
	Type      string    `json:"type"`
	Plat      string    `json:"plat"`
	Status    string    `json:"status"`
	Contact   string    `json:"contact"`
	CreatedAt time.Time `json:"created_at"`
	UpdatedAt time.Time `json:"updated_at"`
}

type CreateVehicleRequest struct {
	UID     string `form:"username"`
	Type    string `form:"type" binding:"required"`
	Plat    string `form:"plat" binding:"required"`
	Contact string `form:"contact"`
	Status  string `form:"status" binding:"required"`
}

type VehicleRepository struct {
	db *gorm.DB
}

func NewVehicleRepository(db *gorm.DB) *VehicleRepository {
	return &VehicleRepository{db: db}
}

func (r *VehicleRepository) Create(vehicle *CreateVehicleRequest) (string, error) {
	// Get count of user's vehicles to generate ID
	var count int64
	if err := r.db.Model(&Vehicle{}).Where("user_id = ?", vehicle.UID).Count(&count).Error; err != nil {
		return "", err
	}

	// Fetch the user's username from the User model
	var user User
	if err := r.db.Where("id = ?", vehicle.UID).First(&user).Error; err != nil {
		return "", err
	}
	id := fmt.Sprintf("%s-%d", user.Username, count+1)

	// Create a new vehicle instance with ID format (username-vehiclecount)
	newVehicle := Vehicle{
		ID:        id,
		UserID:    user.ID, // Set the UserID foreign key
		Type:      vehicle.Type,
		Plat:      vehicle.Plat,
		Contact:   vehicle.Contact,
		Status:    vehicle.Status,
		CreatedAt: time.Now(),
		UpdatedAt: time.Now(),
	}
	return id, r.db.Create(&newVehicle).Error
}

func (r *VehicleRepository) FindByStatus(status string) ([]Vehicle, error) {
	var vehicles []Vehicle
	err := r.db.Where("status = ?", status).Preload("User").Find(&vehicles).Error
	return vehicles, err
}

func (r *VehicleRepository) FindByUsername(username string) ([]Vehicle, error) {
	var vehicles []Vehicle
	err := r.db.
		Joins("JOIN users ON users.id = vehicles.user_id").
		Where("users.username = ?", username).
		Preload("User").
		Find(&vehicles).Error
	return vehicles, err
}

func (r *VehicleRepository) FindByID(id string) (*Vehicle, error) {
	var vehicle Vehicle
	err := r.db.Where("id = ?", id).Preload("User").First(&vehicle).Error
	return &vehicle, err
}
func (r *VehicleRepository) Update(id string, vehicle *CreateVehicleRequest) error {
	var existingVehicle Vehicle
	if err := r.db.Where("id = ?", id).First(&existingVehicle).Error; err != nil {
		return err
	}

	// Update the fields of the existing vehicle
	existingVehicle.Type = vehicle.Type
	existingVehicle.Status = vehicle.Status
	existingVehicle.Plat = vehicle.Plat
	existingVehicle.Contact = vehicle.Contact

	return r.db.Save(&existingVehicle).Error
}

func (r *VehicleRepository) Delete(id string) error {
	return r.db.Where("id = ?", id).Delete(&Vehicle{}).Error
}
