package product

// import (
// 	"errors"

// 	"gorm.io/gorm"
// )

// type Repository interface {
// 	Create(product *Product, categories []string) error
// 	Update(product *Product, categories []string) error
// 	Delete(id uint) error
// 	FindByID(id uint) (*Product, error)
// 	List(page, limit int, category string) ([]Product, int, error)
// 	FindByIDs(ids []uint) ([]Product, error)
// 	UpdateStock(id uint, quantity int) error
// }

// type repository struct {
// 	db *gorm.DB
// }

// func NewRepository(db *gorm.DB) Repository {
// 	return &repository{db: db}
// }

// func (r *repository) Create(product *Product, categoryNames []string) error {
// 	return r.db.Transaction(func(tx *gorm.DB) error {
// 		// Create or get categories
// 		var categories []Category
// 		for _, name := range categoryNames {
// 			var category Category
// 			err := tx.Where("name = ?", name).FirstOrCreate(&category, Category{Name: name}).Error
// 			if err != nil {
// 				return err
// 			}
// 			categories = append(categories, category)
// 		}

// 		// Create product
// 		if err := tx.Create(product).Error; err != nil {
// 			return err
// 		}

// 		// Associate categories
// 		return tx.Model(product).Association("Categories").Replace(categories)
// 	})
// }

// func (r *repository) Update(product *Product, categoryNames []string) error {
// 	return r.db.Transaction(func(tx *gorm.DB) error {
// 		// Update product fields
// 		if err := tx.Model(product).Updates(map[string]interface{}{
// 			"name":        product.Name,
// 			"description": product.Description,
// 			"price":       product.Price,
// 			"stock":       product.Stock,
// 		}).Error; err != nil {
// 			return err
// 		}

// 		// Update categories if provided
// 		if categoryNames != nil {
// 			var categories []Category
// 			for _, name := range categoryNames {
// 				var category Category
// 				err := tx.Where("name = ?", name).FirstOrCreate(&category, Category{Name: name}).Error
// 				if err != nil {
// 					return err
// 				}
// 				categories = append(categories, category)
// 			}
// 			if err := tx.Model(product).Association("Categories").Replace(categories); err != nil {
// 				return err
// 			}
// 		}

// 		return nil
// 	})
// }

// func (r *repository) Delete(id uint) error {
// 	return r.db.Transaction(func(tx *gorm.DB) error {
// 		// Check if product exists
// 		var product Product
// 		if err := tx.First(&product, id).Error; err != nil {
// 			if errors.Is(err, gorm.ErrRecordNotFound) {
// 				return nil
// 			}
// 			return err
// 		}

// 		// Soft delete the product
// 		return tx.Delete(&product).Error
// 	})
// }

// func (r *repository) FindByID(id uint) (*Product, error) {
// 	var product Product
// 	err := r.db.Preload("Categories").First(&product, id).Error
// 	if err != nil {
// 		if errors.Is(err, gorm.ErrRecordNotFound) {
// 			return nil, nil
// 		}
// 		return nil, err
// 	}
// 	return &product, nil
// }

// func (r *repository) List(page, limit int, category string) ([]Product, int, error) {
// 	var products []Product
// 	var total int64
// 	query := r.db.Model(&Product{})

// 	if category != "" {
// 		query = query.Joins("JOIN product_categories ON products.id = product_categories.product_id").
// 			Joins("JOIN categories ON product_categories.category_id = categories.id").
// 			Where("categories.name = ?", category)
// 	}

// 	// Count total records
// 	if err := query.Count(&total).Error; err != nil {
// 		return nil, 0, err
// 	}

// 	// Get paginated records
// 	err := query.Offset((page - 1) * limit).
// 		Limit(limit).
// 		Preload("Categories").
// 		Find(&products).Error

// 	return products, int(total), err
// }

// func (r *repository) FindByIDs(ids []uint) ([]Product, error) {
// 	var products []Product
// 	err := r.db.Preload("Categories").Where("id IN ?", ids).Find(&products).Error
// 	return products, err
// }

// func (r *repository) UpdateStock(id uint, quantity int) error {
// 	return r.db.Transaction(func(tx *gorm.DB) error {
// 		var product Product
// 		if err := tx.Lock("FOR UPDATE").First(&product, id).Error; err != nil {
// 			return err
// 		}

// 		if product.Stock+quantity < 0 {
// 			return errors.New("insufficient stock")
// 		}

// 		return tx.Model(&product).Update("stock", gorm.Expr("stock + ?", quantity)).Error
// 	})
// }
