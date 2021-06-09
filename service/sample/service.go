//Created by Goland
//@User: lenora
//@Date: 2021/3/8
//@Time: 2:43 下午
package sample

import "github.com/jinzhu/gorm"

type SampleService interface {

}

type sampleService struct {
	db *gorm.DB
}

func NewSampleService(db *gorm.DB) SampleService {
	u := new(sampleService)
	u.db = db
	return u
}

func (srv *sampleService) ApplicationList(page, limit uint32) (error, uint32, []Sample) {
	offset := (page - 1) * limit
	return getList(srv.db, offset, limit)
}

