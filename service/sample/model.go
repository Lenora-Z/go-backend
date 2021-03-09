//Created by Goland
//@User: lenora
//@Date: 2021/3/8
//@Time: 2:42 下午
package sample

import (
	"github.com/jinzhu/gorm"
	"time"
)

type Sample struct {
	ID        int       `gorm:"primary_key;column:id;type:int(11) unsigned;not null" json:"id"`
	CreatedAt time.Time `gorm:"column:created_at;type:datetime;not null" json:"created_at"`
	UpdatedAt time.Time `gorm:"column:updated_at;type:datetime" json:"updated_at"`
}

// TableName get sql table name.获取数据库表名
func (m *Sample) TableName() string {
	return "sample"
}

var tablename = (&Sample{}).TableName()

func create(db *gorm.DB, m *Sample) (error, *Sample) {
	err := db.Table(tablename).Create(m).Error
	return err, m
}

func detail(db *gorm.DB, id uint64) (bool, *Sample) {
	var item Sample
	status := db.Table(tablename).Where("id = ?", id).First(&item).RecordNotFound()
	return status, &item
}

func detailByColumn(db *gorm.DB, group map[string]interface{}) (bool, *Sample) {
	var item Sample
	query := db.Table(tablename)
	for i, x := range group {
		query = query.Where(i+" = ?", x)
	}
	status := query.First(&item).RecordNotFound()
	return status, &item
}

func update(db *gorm.DB, m *Sample) (error, *Sample) {
	m.UpdatedAt = time.Now()
	err := db.Table(tablename).Where("id = ?", m.ID).Update(m).Error
	return err, m
}

func getList(db *gorm.DB, offset, limit uint32) (error, uint32, []Sample) {
	var count uint32
	list := make([]Sample, 0)
	query := db.Table(tablename)
	err := query.Count(&count).
		Offset(offset).
		Limit(limit).
		Find(&list).
		Error
	return err, count, list
}



