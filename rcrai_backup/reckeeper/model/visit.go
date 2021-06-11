package model

import (
	"fmt"

	"gorm.io/gorm"
	"rpkg.cc/apps/kerrigan/countermgr/cmmodels/cmmc"
	"rpkg.cc/apps/kerrigan/pkg/reckeeper/reckcfg"
	"rpkg.cc/infras/models"
	"rpkg.cc/infras/models/mc"
)

type Visit struct {
	models.ID
	cmmc.BusinessIDModelIndexed

	DTSourceID string `gorm:"column:dt_id;type:varchar(128);not null;index" json:"dt_id"` // should be sales_sn
	CustomerId string `gorm:"column:customer_id;type:varchar(128);not null;index" json:"customer_id"`

	// 访问记录，和excel同步的记录。
	CustomerName string ``
	StartTime    int64  `gorm:"column:start_time;index" json:"start_time"` // 拜访时间
	Site         string `gorm:"column:site;type:varchar(128)" json:"site"`
	mc.OpTimeModel
}

func (s *Visit) TableName() string {
	return fmt.Sprintf("%s_visit", reckcfg.Schema)
}
func (m *Visit) BeforeCreate(tx *gorm.DB) (err error) {
	if m.GetID() == "" {
		m.FillNewID()
	}
	return nil
}

type VisitDetail struct {
	Visit
	Member   *Member   `json:"member"`
	Customer *Customer `json:"customer"`
}

type VisitRow struct {
	Visit
	CustomerName  string `json:"customer_name"`
	CustomerPhone string `json:"customer_phone"`
	MemberName    string `json:"member_name"`
	MemberPhone   string `json:"member_phone"`
}

func (a VisitDetail) ToVisitRow() (b *VisitRow) {
	b = &VisitRow{
		Visit: a.Visit,
	}
	if a.Customer != nil {
		b.CustomerName = a.Customer.Name
		b.CustomerPhone = a.Customer.Phone
	}
	if a.Member != nil {
		b.MemberName = a.Member.Name
		b.MemberPhone = a.Member.Phone
	}
	return b
}
