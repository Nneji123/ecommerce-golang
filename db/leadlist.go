package db

import (
	"gorm.io/gorm"
	"log"
)

type LeadList struct {
	ID     uint    `gorm:"primaryKey"`
	Name   string  // Name of the lead list
	UserID uint    // User ID associated with the lead list
	Leads  []*Lead `gorm:"foreignKey:LeadListID"` // Relationship with leads
}

// CreateLeadList creates a new lead list in the database
func CreateLeadList(leadList *LeadList) error {
	result := DB.Create(leadList)
	if result.Error != nil {
		log.Printf("Error creating lead list: %v", result.Error)
		return result.Error
	}
	log.Println("Lead list created successfully with ID:", leadList.ID)
	return nil
}

// GetLeadListByID retrieves a lead list from the database by ID
func GetLeadListByID(id uint) (*LeadList, error) {
	leadList := &LeadList{}
	result := DB.Preload("Leads").First(leadList, id)
	if result.Error != nil {
		if result.Error == gorm.ErrRecordNotFound {
			log.Println("Lead list not found")
			return nil, nil // Return nil if no record found
		}
		log.Printf("Error retrieving lead list: %v", result.Error)
		return nil, result.Error
	}
	log.Println("Retrieved lead list:", leadList)
	return leadList, nil
}

// GetAllLeadLists retrieves all lead lists from the database
func GetAllLeadLists() ([]*LeadList, error) {
	var leadLists []*LeadList
	result := DB.Find(&leadLists)
	if result.Error != nil {
		return nil, result.Error
	}
	return leadLists, nil
}

// UpdateLeadList updates a lead list in the database
func UpdateLeadList(leadList *LeadList) error {
	result := DB.Save(leadList)
	if result.Error != nil {
		log.Printf("Error updating lead list: %v", result.Error)
		return result.Error
	}
	log.Println("Lead list updated successfully")
	return nil
}

// DeleteLeadList deletes a lead list from the database
func DeleteLeadList(id uint) error {
	result := DB.Delete(&LeadList{}, id)
	if result.Error != nil {
		log.Printf("Error deleting lead list: %v", result.Error)
		return result.Error
	}
	log.Println("Lead list deleted successfully")
	return nil
}
