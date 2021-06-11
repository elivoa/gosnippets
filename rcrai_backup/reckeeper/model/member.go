package model

import (
	"fmt"

	"gorm.io/gorm"
	"rpkg.cc/apps/kerrigan/countermgr/cmmodels/cmmc"
	"rpkg.cc/apps/kerrigan/pkg/reckeeper/reckcfg"
	"rpkg.cc/infras/models"
	"rpkg.cc/infras/models/mc"
)

// 员工 Later
type Member struct {
	models.ID
	cmmc.BusinessIDModelIndexed

	DtId  string `json:"dt_id,omitempty" gorm:"type:varchar(24);index;"` // dt.sourceId
	Phone string `gorm:"column:phone;type:varchar(128);not null" json:"phone"`
	Name  string `gorm:"column:name;type:varchar(1024)" json:"name"`
	Site  string `gorm:"column:site;type:varchar(128)" json:"site"`
	// create time is authedtimes
	mc.OpTimeModel
}

func (s *Member) TableName() string {
	return fmt.Sprintf("%s_member", reckcfg.Schema)
}

// gorm: hooks to create ID
func (m *Member) BeforeCreate(tx *gorm.DB) (err error) {
	if m.GetID() == "" {
		m.FillNewID()
	}
	return nil
}
