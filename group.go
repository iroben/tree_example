package main

import (
	"fmt"
	"gorm.io/driver/mysql"
	"gorm.io/gorm"
	"gorm.io/gorm/logger"
	"log"
)

var (
	DB *gorm.DB
)

func init() {
	mainDsn := fmt.Sprintf(
		"%s:%s@tcp(%s)/%s?charset=utf8&parseTime=True&loc=Local",
		"root",
		"root",
		"localhost:3306",
		"test")
	db, err := gorm.Open(mysql.Open(mainDsn), &gorm.Config{
		Logger: logger.Default.LogMode(logger.Info),
	})
	if err != nil {
		log.Println("数据库连接创建失败: ", err.Error())
		return
	}
	DB = db
}

type Group struct {
	Id       int      `gorm:"column:id;primary_key:true"`
	Name     string   `gorm:"column:name"`
	Pid      int      `gorm:"column:pid"`
	Lval     int      `gorm:"column:lval"`
	Rval     int      `gorm:"column:rval"`
	Children []*Group `gorm:"-"`
}

func (Group) TableName() string {
	return "group"
}
func (m *Group) Get() (*Group, error) {
	var retVal Group
	if err := DB.Model(&Group{}).
		Where("id=?", m.Id).
		First(&retVal).Error; err != nil {
		return nil, err
	}
	return &retVal, nil
}
func (m *Group) AddRoot() (bool, error) {
	var max int
	err := DB.Model(&Group{}).Select("IFNULL(MAX(rval), 0)").Row().Scan(&max)
	if err != nil {
		return false, err
	}
	m.Lval = max + 1
	m.Rval = max + 2
	return m.Save()
}

func (m *Group) Add() (bool, error) {
	if m.Pid == 0 {
		return m.AddRoot()
	}
	return (&Group{Id: m.Pid}).AddChildren(m)
}

func (m *Group) AddChildren(g *Group) (bool, error) {
	parentGroup, err := m.Get()
	if parentGroup == nil {
		return false, err
	}
	result, err := parentGroup.updateLeft()
	if !result {
		return false, err
	}
	result, err = parentGroup.updateRight()
	if !result {
		return false, err
	}
	g.Pid = parentGroup.Id
	g.Lval = parentGroup.Rval
	g.Rval = parentGroup.Rval + 1
	return g.Save()
}

func (m *Group) updateLeft() (bool, error) {
	err := DB.Model(&Group{}).
		Where("lval>?", m.Rval).
		UpdateColumn("lval", gorm.Expr("lval + 2")).Error
	return err == nil, err
}

func (m *Group) updateRight() (bool, error) {
	err := DB.Model(&Group{}).
		Where("rval>=?", m.Rval).
		UpdateColumn("rval", gorm.Expr("rval + 2")).Error
	return err == nil, err
}

func (m *Group) Save() (bool, error) {
	err := DB.Create(m).Error
	if err != nil {
		return false, err
	}
	return true, nil
}

func (m *Group) FindAll() ([]*Group, error) {
	var retVal []*Group
	if err := DB.Model(&Group{}).Find(&retVal).Error; err != nil {
		return nil, err
	}
	g := &Group{}
	g.formatChild(retVal)
	return g.Children, nil
}

func (m *Group) formatChild(groups []*Group) {
	children := make([]*Group, 0)
	for _, g := range groups {
		if g.Pid == m.Id {
			g.formatChild(groups)
			children = append(children, g)
		}
	}
	m.Children = children
}
