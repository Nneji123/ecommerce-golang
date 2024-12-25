package db

import (
	"gorm.io/datatypes"
	"gorm.io/gorm"
	"log"
)

type Lead struct {
	ID                      uint `gorm:"primaryKey"`
	LeadListID              uint // Foreign key for the lead list
	Name                    string
	Phone                   uint
	Location                string
	Website                 string
	LinkedinURL             string
	EmailAddress            string
	ScrapedDataFromLinkedin datatypes.JSON
}

// CreateLead creates a new lead record in the database
func CreateLead(lead *Lead) error {
	result := DB.Create(lead)
	if result.Error != nil {
		log.Printf("Error creating lead: %v", result.Error)
		return result.Error
	}
	log.Println("Lead created successfully with ID:", lead.ID)
	return nil
}

// GetLead retrieves a lead record from the database by ID
func GetLead(id uint) (*Lead, error) {
	lead := &Lead{}
	result := DB.First(lead, id)
	if result.Error != nil {
		if result.Error == gorm.ErrRecordNotFound {
			log.Println("Lead not found")
			return nil, nil // Return nil if no record found
		}
		log.Printf("Error retrieving lead: %v", result.Error)
		return nil, result.Error
	}
	log.Println("Retrieved lead:", lead)
	return lead, nil
}

// GetAllLeads retrieves all leads from the database
func GetAllLeads() ([]*Lead, error) {
	var leads []*Lead
	result := DB.Find(&leads)
	if result.Error != nil {
		return nil, result.Error
	}
	return leads, nil
}

// UpdateLead updates a lead record in the database
func UpdateLead(lead *Lead) error {
	result := DB.Save(lead)
	if result.Error != nil {
		log.Printf("Error updating lead: %v", result.Error)
		return result.Error
	}
	log.Println("Lead updated successfully")
	return nil
}

// DeleteLead deletes a lead record from the database
func DeleteLead(id uint) error {
	result := DB.Delete(&Lead{}, id)
	if result.Error != nil {
		log.Printf("Error deleting lead: %v", result.Error)
		return result.Error
	}
	log.Println("Lead deleted successfully")
	return nil
}
