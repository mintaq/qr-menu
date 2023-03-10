package models

import (
	"errors"
	"log"
	"math"
	"regexp"
	"strconv"
	"strings"

	"github.com/gofiber/fiber/v2"
	"gorm.io/gorm"
)

type Pagination struct {
	Limit      int   `json:"limit" gorm:"omitempty;query:limit"`
	Page       int   `json:"page" gorm:"omitempty;query:page"`
	TotalRows  int64 `json:"total_rows"`
	TotalPages int64 `json:"total_pages"`
}

type PaginationQueryParams struct {
	Limit  int    `query:"limit"`
	Page   int    `query:"page"`
	SortBy string `query:"sort_by"`
}

func (p *Pagination) GetOffset() int {
	return (p.GetPage() - 1) * p.GetLimit()
}

func (p *Pagination) GetLimit() int {
	if p.Limit == 0 {
		p.Limit = 10
	}
	return p.Limit
}

func (p *Pagination) GetPage() int {
	if p.Page == 0 {
		p.Page = 1
	}
	return p.Page
}

// Paginate performs pagination on a GORM query by extracting pagination parameters from
// a Fiber context object and applying them to the query. The function returns a Pagination
// object containing information about the current page and the total number of pages, as well
// as a function that can be used to apply the pagination to a GORM query. The pagination parameters
// are extracted from the Fiber context object as follows:
//
//   - The "page" query parameter specifies the current page number. If this parameter is not present
//     or its value is invalid, the function defaults to page 1.
//   - The "limit" query parameter specifies the maximum number of rows per page. If this parameter is
//     not present or its value is invalid, the function defaults to 10.
//   - The "sort_by" query parameter specifies the column to sort by and the sort order (ascending or
//     descending). If this parameter is not present or its value is invalid, the function defaults to
//     sorting by the "id" column in ascending order.
//
// The function also takes a GORM query object and a model object as parameters. The query object
// should already contain any necessary filters for the query, and the model object should be a
// pointer to the model that the query operates on. The function applies the pagination to the
// query by adding an ORDER BY clause with the specified sort order, an OFFSET clause with the
// calculated offset, and a LIMIT clause with the specified limit. The function also calculates
// the total number of pages by dividing the total number of rows by the limit, and rounding up
// to the nearest integer.
//
// Example usage:
//
//	var users []User
//	query := db.Model(&users).Where("status = ?", "active")
//	pagination, applyPagination := Paginate(&users, c, query)
//	applyPagination(query).Find(&users)
//
// This example applies pagination to a query that retrieves all active users, using pagination
// parameters extracted from a Fiber context object "c". The `applyPagination` function can be used
// to apply the pagination to the query, and the `pagination` object contains information about
// the current page and the total number of pages.
func Paginate(model interface{}, c *fiber.Ctx, query *gorm.DB) (p *Pagination, f func(db *gorm.DB) *gorm.DB) {
	var totalPages, totalRows int64
	page, _ := strconv.Atoi(c.Query("page", "1"))
	if page <= 0 {
		page = 1
	}

	limit, _ := strconv.Atoi(c.Query("limit", "10"))
	switch {
	case limit > 100:
		limit = 100
	case limit <= 0:
		limit = 1
	}

	sortBy := c.Query("sort_by", "id:asc")
	if match, _ := regexp.MatchString(`^[a-zA-Z]+:(?i)(asc|desc)$`, sortBy); !match {
		sortBy = "id:asc"
	}
	sortBy = strings.ReplaceAll(sortBy, ":", " ")

	offset := (page - 1) * limit

	if tx := query.Select("id").Count(&totalPages); tx.Error != nil {
		log.Println(tx.Error.Error())
		return nil, nil
	}

	if tx := query.Select("id").Limit(limit).Offset(offset).Count(&totalRows); tx.Error != nil {
		log.Println(tx.Error.Error())
		return nil, nil
	}

	pagination := Pagination{
		Limit:      limit,
		Page:       page,
		TotalPages: int64(math.Round(float64(totalPages) / float64(limit))),
		TotalRows:  totalRows,
	}

	return &pagination, func(db *gorm.DB) *gorm.DB {
		return db.Where(query).Limit(limit).Offset(offset).Order(sortBy)
	}
}

func Paginate2(c *fiber.Ctx, query *gorm.DB) (*Pagination, *gorm.DB, error) {
	// Initialize pagination variables
	var totalPages, totalRows int64
	page, err := strconv.Atoi(c.Query("page", "1"))
	if err != nil || page <= 0 {
		return nil, nil, errors.New("invalid page number")
	}

	limit, err := strconv.Atoi(c.Query("limit", "10"))
	if err != nil || limit <= 0 || limit > 100 {
		return nil, nil, errors.New("invalid limit value")
	}

	sortBy := c.Query("sort_by", "id:asc")
	if match, _ := regexp.MatchString(`^[a-zA-Z]+:(?i)(asc|desc)$`, sortBy); !match {
		sortBy = "id:asc"
	}
	sortBy = strings.ReplaceAll(sortBy, ":", " ")

	offset := (page - 1) * limit

	// Count the total number of pages and rows
	if err := query.Count(&totalPages).Error; err != nil {
		return nil, nil, err
	}

	if err := query.Limit(limit).Offset(offset).Order(sortBy).Count(&totalRows).Error; err != nil {
		return nil, nil, err
	}

	db := query.Session(&gorm.Session{NewDB: true})
	sql := db.ToSQL(func(tx *gorm.DB) *gorm.DB {
		return tx.Model(Product{}).Limit(limit).Offset(offset).Order(sortBy).Count(&totalRows)
	})

	log.Println(sql)

	// Create pagination object
	pagination := Pagination{
		Limit:      limit,
		Page:       page,
		TotalPages: int64(math.Round(float64(totalPages) / float64(limit))),
		TotalRows:  totalRows,
	}

	paginatedQuery := query.Limit(limit).Offset(offset).Order(sortBy)

	// Return the pagination object and a function that applies pagination to a GORM query
	return &pagination, paginatedQuery, nil
}

func PaginateByScope(model interface{}, c *fiber.Ctx, scope func(db *gorm.DB) *gorm.DB, db *gorm.DB) (p *Pagination, f func(db *gorm.DB) *gorm.DB) {
	var totalPages, totalRows int64
	page, _ := strconv.Atoi(c.Query("page", "1"))
	if page <= 0 {
		page = 1
	}

	limit, _ := strconv.Atoi(c.Query("limit", "10"))
	switch {
	case limit > 100:
		limit = 100
	case limit <= 0:
		limit = 1
	}

	sortBy := c.Query("sort_by", "id:asc")
	if match, _ := regexp.MatchString(`^[a-zA-Z]+:(?i)(asc|desc)$`, sortBy); !match {
		sortBy = "id:asc"
	}
	sortBy = strings.ReplaceAll(sortBy, ":", " ")

	offset := (page - 1) * limit

	db.Model(&model).Scopes(scope).Count(&totalPages)
	db.Model(&model).Scopes(scope).Offset(offset).Limit(limit).Count(&totalRows)

	pagination := Pagination{
		Limit:      limit,
		Page:       page,
		TotalPages: int64(math.Round(float64(totalPages) / float64(limit))),
		TotalRows:  totalRows,
	}

	return &pagination, func(db *gorm.DB) *gorm.DB {
		return db.Scopes(scope).Offset(offset).Limit(limit).Order(sortBy)
	}
}
