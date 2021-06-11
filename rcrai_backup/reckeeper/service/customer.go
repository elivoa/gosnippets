package service

import (
	"fmt"

	"gorm.io/gorm"
	"rpkg.cc/apps/kerrigan/gwauth"
	"rpkg.cc/apps/kerrigan/pkg/reckeeper/model"
	"rpkg.cc/infra/utils/times"
	"rpkg.cc/infras/helper/dbs"
	"rpkg.cc/infras/services"
)

func (s *Service) CreateCustomer(model *model.Customer, ctx *services.Context) (
	result *model.Customer, created *services.CreateResult, err error) {

	authuser := gwauth.GetAuthUserFromServiceContext(ctx).RequireUser()

	var db *gorm.DB
	if _, db, err = services.GetDB(ctx); err != nil {
		return
	}

	err = db.Transaction(func(tx *gorm.DB) error {
		ctx.UseTransaction(tx)

		// 2. create one.
		model.SetBid(authuser.GetBidSafe())

		if err := tx.Create(model).Error; err != nil {
			return fmt.Errorf("error when creating %s: %w", model.TableName(), err)
		}
		result = model
		created = services.NewCreateResult(model.GetID()) // tx.RowsAffected
		return nil
	})
	return
}

func (s *Service) SearchCustomers(customer_id string, assigned_staff string, pager *services.PagerInfo, ctx *services.Context) (
	models []*model.Customer, total int64, err error) {

	authuser := gwauth.GetAuthUserFromServiceContext(ctx)

	cond := dbs.NewCond().Add("bid", authuser.GetBidSafe())
	cond.AddIf(customer_id != "", "customer_id", customer_id)
	cond.AddIf(assigned_staff != "", "assigned_staff", assigned_staff)

	var db *gorm.DB
	if _, db, err = services.GetDB(ctx); err != nil {
		return
	}

	db = db.Model(&model.Customer{}).Where(cond.Conditions()) // todo: performance, cache model.
	db.Order("customer_id asc")

	// * Count and query in helper.
	if total, err = dbs.CountAndQuery(db, pager, &models); err != nil {
		return
	}
	return
}

func (s *Service) UpdateCustomer(model *model.Customer, ctx *services.Context) (
	result *model.Customer, updated *services.UpdateResult, err error) {

	// * Do db operations. init service.Context
	ctx, db, err := services.GetDB(ctx)
	if err != nil {
		return nil, nil, err
	}

	err = db.Transaction(func(tx *gorm.DB) error {
		ctx.UseTransaction(tx)

		// 2. update branch.
		if err := tx.Model(model).Updates(model).Error; err != nil {
			return fmt.Errorf("error when updating %s: %w", model.TableName(), err)
		}
		result = model
		updated = services.NewUpdateResult(model.GetID()) // tx.RowsAffected
		return nil
	})
	return
}

// TODO 疑问：assigned staff是什么东西？
func (s *Service) CustomerPush(model *model.Customer, ctx *services.Context) (
	result *model.Customer, err error) {

	foundCustomers, _, err := s.SearchCustomers(model.CustomerId, model.AssignedStaff, nil, ctx)
	if err != nil {
		return nil, err
	}
	if len(foundCustomers) > 0 {
		dbitem := foundCustomers[0]
		id := dbitem.GetID()
		model.SetID(&id)
		model.SetBid(dbitem.GetBid())
		model.SetCreatedAt(dbitem.GetCreatedAt())
		model.SetUpdatedAt(times.NowRef())
		result, _, err = s.UpdateCustomer(model, ctx)
	} else {
		// create
		result, _, err = s.CreateCustomer(model, ctx)
	}
	return
}

func (s *Service) GetCustomer(customerID string, ctx *services.Context) (result *model.Customer, err error) {
	if foundCustomers, _, err := s.SearchCustomers(customerID, "", services.NewPager(0, 1), ctx); err != nil {
		return nil, err
	} else if len(foundCustomers) > 0 {
		return foundCustomers[0], nil
	}
	return
}
