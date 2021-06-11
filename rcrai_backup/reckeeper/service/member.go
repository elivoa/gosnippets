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

// ListStaff - support: filters.Type
func (s *Service) SearchMembers(dt_id string, site string, pager *services.PagerInfo, ctx *services.Context) (
	models []*model.Member, total int64, err error) {

	authuser := gwauth.GetAuthUserFromServiceContext(ctx)

	cond := dbs.NewCond().Add("bid", authuser.GetBidSafe())
	cond.AddIf(dt_id != "", "dt_id", dt_id)
	cond.AddIf(site != "", "site", site)

	var db *gorm.DB
	if _, db, err = services.GetDB(ctx); err != nil {
		return
	}

	db = db.Model(&model.Member{}).Where(cond.Conditions()) // todo: performance, cache model.
	db.Order("dt_id asc")

	// * Count and query in helper.
	if total, err = dbs.CountAndQuery(db, pager, &models); err != nil {
		return
	}
	return
}

func (s *Service) CreateMember(model *model.Member, ctx *services.Context) (
	result *model.Member, created *services.CreateResult, err error) {

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

func (s *Service) UpdateMember(model *model.Member, ctx *services.Context) (
	result *model.Member, updated *services.UpdateResult, err error) {

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

func (s *Service) GetMember(id string, ctx *services.Context) (result *model.Member, err error) {
	if models, _, err := s.SearchMembers(id, "", services.NewPager(0, 1), ctx); err != nil {
		return nil, err
	} else if len(models) > 0 {
		return models[0], nil
	}
	return
}

// TODO 这里判断一下，如果不需要更新，为何还要更新数据库？只更新一个时间么？
func (s *Service) MemberPush(model *model.Member, ctx *services.Context) (
	result *model.Member, err error) {

	items, _, err := s.SearchMembers(model.DtId, model.Site, nil, ctx)
	if err == nil && len(items) > 0 {
		dbitem := items[0]
		id := dbitem.GetID()
		model.SetID(&id)
		model.SetBid(dbitem.GetBid())
		model.SetCreatedAt(dbitem.GetCreatedAt())
		model.SetUpdatedAt(times.NowRef())
		result, _, err = s.UpdateMember(model, ctx)
	} else {
		result, _, err = s.CreateMember(model, ctx)
	}
	if err != nil {
		return
	}
	return result, err
}
