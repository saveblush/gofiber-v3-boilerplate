package models

import "github.com/saveblush/gofiber-v3-boilerplate/internal/core/generic"

// PageInformation page information
type PageInformation struct {
	Page     int   `json:"page,omitempty"`
	Size     int   `json:"size,omitempty"`
	Count    int64 `json:"count,omitempty"`
	LastPage int   `json:"last_page,omitempty"`
}

// Page page model
type Page struct {
	PageInformation *PageInformation `json:"page_information,omitempty"`
	Entities        interface{}      `json:"entities,omitempty"`
}

// NewPage new page
func NewPage(pageInfo *PageInformation, entities interface{}) *Page {
	v := &Page{}
	v.Entities = entities

	if !generic.IsEmpty(pageInfo) {
		v.PageInformation = pageInfo
	}

	return v
}

// GetEntities get entities
func (p *Page) GetEntities() interface{} {
	return p.Entities
}

// PageForm page form
type PageForm struct {
	Page    int    `json:"page,omitempty" form:"page" query:"page"`
	Size    int    `json:"size,omitempty" form:"size" query:"size"`
	Query   string `json:"query,omitempty" form:"query" query:"query"`
	Sort    string `json:"sort,omitempty" form:"sort" query:"sort"`
	Reverse bool   `json:"reverse,omitempty" form:"reverse" query:"reverse"`
	OrderBy string `json:"-" form:"-"`
}

// GetPage get page
func (f *PageForm) GetPage() int {
	return f.Page
}

// GetSize get size
func (f *PageForm) GetSize() int {
	return f.Size
}

// GetQuery get query
func (f *PageForm) GetQuery() string {
	return f.Query
}

// GetSort get sort
func (f *PageForm) GetSort() string {
	return f.Sort
}

// GetReverse get reverse
func (f *PageForm) GetReverse() bool {
	return f.Reverse
}

// GetOrderBy get order by
func (f *PageForm) GetOrderBy() string {
	return f.OrderBy
}
