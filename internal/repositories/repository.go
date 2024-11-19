package repositories

import (
	"fmt"
	"math"
	"strings"

	"gorm.io/gorm"
	"gorm.io/gorm/clause"
	"gorm.io/gorm/schema"

	"github.com/samber/lo"

	"github.com/saveblush/gofiber-v3-boilerplate/internal/core/utils"
	"github.com/saveblush/gofiber-v3-boilerplate/internal/models"
)

// Repository common repository
type Repository struct {
	//mutex sync.Mutex
}

// NewRepository new repository
func NewRepository() Repository {
	return Repository{}
}

// Find find
func (r *Repository) Find(db *gorm.DB, i interface{}) error {
	return db.First(i).Error
}

// FindByID find by id
func (r *Repository) FindByID(db *gorm.DB, id uint, i interface{}) error {
	return r.FindByField(db, "id", id, i)
}

// FindByIDString find
func (r *Repository) FindByIDString(db *gorm.DB, field string, value string, i interface{}) error {
	return r.FindByField(db, field, value, i)
}

// FindByField find
func (r *Repository) FindByField(db *gorm.DB, field string, value interface{}, i interface{}) error {
	return db.Where(fmt.Sprintf("%s = ?", field), value).First(i).Error
}

// FindLastByField find
func (r *Repository) FindLastByField(db *gorm.DB, field string, value interface{}, i interface{}) error {
	return db.Where(fmt.Sprintf("%s = ?", field), value).Last(i).Error
}

// FindAllByIDs get all by ids
func (r *Repository) FindAllByIDs(db *gorm.DB, ids []uint, i interface{}) error {
	return db.Where("id in (?)", ids).Find(i).Error
}

// FindAllByStrings get all by strins
func (r *Repository) FindAllByStrings(db *gorm.DB, field string, values []string, i interface{}) error {
	return db.Where(fmt.Sprintf("%s in (?)", field), values).Find(i).Error
}

// FindAllByField get all by field
func (r *Repository) FindAllByField(db *gorm.DB, field string, values interface{}, i interface{}) error {
	return db.Where(fmt.Sprintf("%s in (?)", field), values).Find(i).Error
}

// FindAllByValues get all by values
func (r *Repository) FindAllByValues(db *gorm.DB, field string, values interface{}, i interface{}) error {
	return db.Where(fmt.Sprintf("%s IN (?)", field), values).Find(i).Error
}

// Create create
func (r *Repository) Create(db *gorm.DB, i interface{}) error {
	return db.Omit(clause.Associations).Create(i).Error
}

// CreateInBatch create with batch size
func (r *Repository) CreateInBatch(db *gorm.DB, i interface{}, batchSize int) error {
	return db.Omit(clause.Associations).CreateInBatches(i, batchSize).Error
}

// Update update
func (r *Repository) Update(db *gorm.DB, m, i interface{}) error {
	return db.Model(m).Omit(clause.Associations).Updates(i).Error
}

// Upsert upsert
func (r *Repository) Upsert(db *gorm.DB, uniqueKey string, columns []string, i interface{}) error {
	split := strings.Split(uniqueKey, ",")
	uniqueKeys := []clause.Column{}
	for _, v := range split {
		uniqueKeys = append(uniqueKeys, clause.Column{Name: v})
	}

	return db.Clauses(clause.OnConflict{
		Columns:   uniqueKeys,
		DoUpdates: clause.AssignmentColumns(columns),
		UpdateAll: len(columns) == 0,
	}).
		Omit(clause.Associations).
		Create(i).Error
}

// BulkUpsert bulk upsert
func (r *Repository) BulkUpsert(db *gorm.DB, uniqueKeys []string, columns []string, i interface{}, batchSize int, doNothing ...bool) error {
	uniqueColumns := []clause.Column{}
	for _, uniqueKey := range uniqueKeys {
		uniqueColumns = append(uniqueColumns, clause.Column{Name: uniqueKey})
	}

	condition := clause.OnConflict{
		Columns:   uniqueColumns,
		DoUpdates: clause.AssignmentColumns(columns),
		DoNothing: doNothing != nil && doNothing[0],
	}

	if !condition.DoNothing {
		condition.UpdateAll = len(columns) == 0
	}

	return db.Clauses(condition).
		Omit(clause.Associations).
		CreateInBatches(i, batchSize).Error
}

// Delete update stamp deleted_at
func (r *Repository) Delete(db *gorm.DB, i interface{}) error {
	return db.Omit(clause.Associations).Delete(i).Error
}

// SoftDelete soft delete
func (r *Repository) SoftDelete(db *gorm.DB, field string, value interface{}, actorID string, i interface{}) error {
	values := map[string]interface{}{
		"deleted_at": utils.Now(),
	}
	if actorID != "" {
		values["deleted_user"] = actorID
	}
	err := db.
		Model(i).
		Where(fmt.Sprintf("%s IN (?)", field), value).
		Updates(values).Error
	if err != nil {
		return err
	}

	return nil
}

// PageForm page info interface
type PageForm interface {
	GetPage() int
	GetSize() int
	GetQuery() string
	GetSort() string
	GetReverse() bool
	GetOrderBy() string
}

const (
	// DefaultPage default page in page query
	DefaultPage int = 1
	// DefaultSize default size in page query
	DefaultSize int = 25
)

// FindAllAndPageInformation get page information
func (r *Repository) FindAllAndPageInformation(db *gorm.DB, pageForm PageForm, entities interface{}, selector ...string) (*models.PageInformation, error) {
	var count int64
	stmt := gorm.Statement{DB: db}

	statementModel := entities
	if db.Statement.Model != nil && db.Statement.Table != "" {
		statementModel = db.Statement.Model
	}

	if err := stmt.Parse(statementModel); err != nil {
		return nil, err
	}

	if db.Statement.Distinct {
		selects := stmt.DB.Statement.Selects
		db = db.Distinct(selects[0]).Model(entities).Count(&count)
		db = db.Distinct(selects)
	} else {
		db.Model(entities).Count(&count)
	}

	sort := pageForm.GetSort()
	if pageForm.GetOrderBy() != "" {
		db = db.Order(pageForm.GetOrderBy())
	} else if sort != "" {
		order := sort
		if pageForm.GetReverse() {
			order = order + " DESC"
		}
		db = db.Order(order)
	} else {
		db = db.Order("id")
	}

	page := pageForm.GetPage()
	if pageForm.GetPage() < 1 {
		page = DefaultPage
	}

	limit := pageForm.GetSize()
	if pageForm.GetSize() == 0 {
		limit = DefaultSize
	}

	var offset int
	if page != 1 {
		offset = (page - 1) * limit
	}

	if len(selector) > 0 {
		db = db.Select(selector)
	}

	if err := db.
		Limit(limit).
		Offset(offset).
		Find(entities).Error; err != nil {
		return nil, err
	}

	return &models.PageInformation{
		Page:     page,
		Size:     limit,
		Count:    count,
		LastPage: int(math.Ceil(float64(count) / float64(limit))),
	}, nil
}

// SortInfo sort info
type SortInfo interface {
	GetSort() string
	GetReverse() bool
}

// Sort order by
func (r *Repository) Sort(opt SortInfo, entities interface{}) func(db *gorm.DB) *gorm.DB {
	return func(db *gorm.DB) *gorm.DB {
		stmt := gorm.Statement{DB: db}
		err := stmt.Parse(entities)
		if err != nil {
			return db
		}

		defaultSort := func() *gorm.DB {
			return db.Order(fmt.Sprintf("%s.%s",
				stmt.Schema.Table,
				stmt.Schema.PrioritizedPrimaryField.DBName,
			))
		}
		if opt.GetSort() == "" {
			return defaultSort()
		}

		_, found := lo.Find(stmt.Schema.Fields, func(field *schema.Field) bool {
			return opt.GetSort() == field.DBName
		})
		if found {
			order := opt.GetSort()
			if opt.GetReverse() {
				order = order + " DESC"
			}
			return db.Order(order)
		}

		return defaultSort()
	}
}
