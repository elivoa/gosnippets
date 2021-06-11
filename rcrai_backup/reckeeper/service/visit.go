package service

import (
	"fmt"

	"gorm.io/gorm"
	"rpkg.cc/apps/kerrigan/gwauth"
	"rpkg.cc/apps/kerrigan/pkg/reckeeper/model"
	"rpkg.cc/ecode"
	"rpkg.cc/infra/utils/times"
	"rpkg.cc/infras/driver/dbdriver"
	"rpkg.cc/infras/helper/dbs"
	"rpkg.cc/infras/services"
	"rpkg.cc/log"
)

func (s *Service) CreateVisit(model *model.Visit, ctx *services.Context) (
	result *model.Visit, created *services.CreateResult, err error) {

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

func (s *Service) UpdateVisit(model *model.Visit, ctx *services.Context) (
	result *model.Visit, updated *services.UpdateResult, err error) {

	// * Do db operations. init service.Context
	ctx, db, err := services.GetDB(ctx)
	if err != nil {
		return nil, nil, err
	}

	err = db.Transaction(func(tx *gorm.DB) error {
		ctx.UseTransaction(tx)

		// 2. update branch.
		if err := tx.Save(model).Error; err != nil {
			return fmt.Errorf("error when updating %s: %w", model.TableName(), err)
		}
		result = model
		updated = services.NewUpdateResult(model.GetID()) // tx.RowsAffected
		return nil
	})
	return
}

func (s *Service) VisitPush(model *model.Visit, ctx *services.Context) (
	result *model.Visit, err error) {

	items, _, err := s.VisitsSearch(
		model.Site, model.DTSourceID, model.CustomerId, model.StartTime-300, model.StartTime+300,
		false, nil, ctx,
	)
	if err == nil && len(items) > 0 {
		dbitem := items[0]
		id := dbitem.GetID()
		model.SetID(&id)
		model.SetBid(dbitem.GetBid())
		model.SetCreatedAt(dbitem.GetCreatedAt())
		model.SetUpdatedAt(times.NowRef())
		result, _, err = s.UpdateVisit(model, ctx)
	} else {
		result, _, err = s.CreateVisit(model, ctx)
	}
	if err != nil {
		return
	}
	return result, err
}

func (s *Service) VisitFind(id string) (result *model.Visit, err error) {
	if id == "" {
		return nil, ecode.InvalidArgument
	}

	var db *gorm.DB
	if db, err = dbdriver.QuickGetDB(); err != nil {
		return
	}

	if err = db.First(&result, id).Error; err != nil {
		return
	}
	return
}

// ! --------- gaobo's refactoring line ----------------

func (s *Service) VisitsSearch(site string, dt_id string, customer_id string, began_time int64, end_time int64,
	withDetail bool, pager *services.PagerInfo, ctx *services.Context) (
	models []*model.VisitDetail, total int64, err error) {

	authuser := gwauth.GetAuthUserFromServiceContext(ctx)

	cond := dbs.NewCond().Add("bid", authuser.GetBidSafe())
	cond.AddIf(site != "", "site", site)
	cond.AddIf(dt_id != "", "dt_id", dt_id)
	cond.AddIf(customer_id != "", "customer_id", customer_id)
	cond.AddwcWhen(began_time > 0 && end_time > 0, "start_time BETWEEN ? AND ?", began_time, end_time)

	var db *gorm.DB
	if _, db, err = services.GetDB(ctx); err != nil {
		return
	}

	db = db.Model(&model.Visit{}).Where(cond.Conditions())
	db.Order("start_time asc")

	var visits []*model.Visit

	// * Count and query in helper.
	if total, err = dbs.CountAndQuery(db, pager, &visits); err != nil {
		return
	}
	log.Info("VisitsSearch", site, dt_id, customer_id, began_time, end_time, "results", total)

	// TODO kill entropy！ 这后面应该是个fill的过程，需要改。
	for _, visit := range visits {
		var v = model.VisitDetail{Visit: *visit}
		if withDetail {
			member, err := s.GetMember(visit.DTSourceID, ctx)
			if err != nil {
				log.Errorf("%w", err)
			}
			customer, err := s.GetCustomer(visit.CustomerId, ctx)
			if err != nil {
				log.Errorf("%w", err)
			}
			v.Member = member
			v.Customer = customer
		}
		models = append(models, &v)
	}
	return
}
