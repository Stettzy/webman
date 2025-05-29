package services

import (
	"errors"
	"time"

	"webman/pkg/database"
	"webman/pkg/models"

	"github.com/google/uuid"
	"gorm.io/gorm"
)

type CollectionService struct {
	db *gorm.DB
}

func NewCollectionService() *CollectionService {
	return &CollectionService{
		db: database.DB,
	}
}

func (s *CollectionService) CreateCollection(name, description string) (*models.Collection, error) {
	collection := models.Collection{
		ID:          uuid.New().String(),
		Name:        name,
		Description: description,
		CreatedAt:   time.Now(),
		UpdatedAt:   time.Now(),
	}

	if err := s.db.Create(&collection).Error; err != nil {
		return nil, err
	}

	return &collection, nil
}

func (s *CollectionService) GetCollection(id string) (*models.Collection, error) {
	var collection models.Collection
	if err := s.db.Preload("Requests").First(&collection, "id = ?", id).Error; err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, errors.New("collection not found")
		}
		return nil, err
	}
	return &collection, nil
}

func (s *CollectionService) ListCollections() ([]models.Collection, error) {
	var collections []models.Collection
	if err := s.db.Preload("Requests").Find(&collections).Error; err != nil {
		return nil, err
	}
	return collections, nil
}

func (s *CollectionService) AddRequest(collectionID string, request models.Request) error {
	// Check if collection exists
	var collection models.Collection
	if err := s.db.First(&collection, "id = ?", collectionID).Error; err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return errors.New("collection not found")
		}
		return err
	}

	request.ID = uuid.New().String()
	request.CollectionID = collectionID
	request.CreatedAt = time.Now()
	request.UpdatedAt = time.Now()

	return s.db.Create(&request).Error
}

func (s *CollectionService) DeleteCollection(id string) error {
	// Delete all requests in the collection first
	if err := s.db.Delete(&models.Request{}, "collection_id = ?", id).Error; err != nil {
		return err
	}

	// Then delete the collection
	if err := s.db.Delete(&models.Collection{}, "id = ?", id).Error; err != nil {
		return err
	}

	return nil
}

func (s *CollectionService) UpdateCollection(id string, name, description string) (*models.Collection, error) {
	var collection models.Collection
	if err := s.db.First(&collection, "id = ?", id).Error; err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, errors.New("collection not found")
		}
		return nil, err
	}

	collection.Name = name
	collection.Description = description
	collection.UpdatedAt = time.Now()

	if err := s.db.Save(&collection).Error; err != nil {
		return nil, err
	}

	return &collection, nil
}

func (s *CollectionService) DeleteRequest(collectionID, requestID string) error {
	result := s.db.Delete(&models.Request{}, "id = ? AND collection_id = ?", requestID, collectionID)
	if result.Error != nil {
		return result.Error
	}
	if result.RowsAffected == 0 {
		return errors.New("request not found")
	}
	return nil
}

func (s *CollectionService) UpdateRequest(collectionID string, request models.Request) error {
	var existingRequest models.Request
	if err := s.db.First(&existingRequest, "id = ? AND collection_id = ?", request.ID, collectionID).Error; err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return errors.New("request not found")
		}
		return err
	}

	request.CollectionID = collectionID
	request.UpdatedAt = time.Now()

	return s.db.Save(&request).Error
}
