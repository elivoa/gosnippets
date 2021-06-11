package model

import (
	"fmt"

	"gorm.io/gorm"
	"rpkg.cc/apps/kerrigan/countermgr/cmmodels/cmmc"
	"rpkg.cc/apps/kerrigan/pkg/reckeeper/reckcfg"
	"rpkg.cc/infras/models"
	"rpkg.cc/infras/models/mc"
)

type Customer struct {
	models.ID
	cmmc.BusinessIDModelCUIndex
	CustomerId    string `gorm:"column:customer_id;type:varchar(128);not null;index:bid_comp_idx,unique" json:"customer_id"`
	Name          string `gorm:"column:name;type:varchar(1024);not null" json:"name"`
	Phone         string `gorm:"column:phone;type:varchar(128);not null" json:"phone"`
	AssignedStaff string `gorm:"column:assigned_staff;type:varchar(128)" json:"assigned_staff"`

	mc.OpTimeModel
}

func (s *Customer) TableName() string {
	return fmt.Sprintf("%s_customer", reckcfg.Schema)
}

func (m *Customer) BeforeCreate(tx *gorm.DB) (err error) {
	if m.GetID() == "" {
		m.FillNewID()
	}
	return nil
}
