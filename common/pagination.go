package common

import (
	"github.com/go-xorm/xorm"
	"strconv"
)

type PageParam struct {
	Page int `form:"page"`
	PageSize int `form:"pagesize"`
}

type Pagination struct {
	Page int
	PageSize int
	TotalCount int
}

func (p *Pagination) getLimit() int {
	if p.PageSize > 1 {
		return p.PageSize
	} else {
		return -1
	}
}

func (p *Pagination) getOffset() int {
	if p.Page < 1 {
		return 0
	} else {
		return (p.Page-1) * p.PageSize
	}
}

func (p *Pagination) getPageTotal() int {
	return (p.TotalCount + p.PageSize - 1) / p.PageSize
}

func FindPaginationData(query *xorm.Session, rowsSlicePtr interface{},pageParam PageParam, bean ...interface{}) (*Pagination, error) {
	p := &Pagination{
		Page: pageParam.Page,
		PageSize:pageParam.PageSize,
	}
	//set default value
	if p.PageSize == 0 {
		p.PageSize = 10
	}
	queryCount := *query;
	if c, err := queryCount.Select("count(*)").Count(bean...); err != nil {
		return nil, err
	} else {
		cs :=strconv.FormatInt(c,10)
		p.TotalCount, _ = strconv.Atoi(cs)
	}
	if err := query.Limit(p.getLimit(), p.getOffset()).Find(rowsSlicePtr); err != nil {
		return nil, err
	} else {
		return p, nil
	}
}