package badapi

import (
	"github.com/gin-gonic/gin"
	"rpkg.cc/apps/kerrigan/gwauth"
	"rpkg.cc/apps/kerrigan/pkg/reckeeper/model"
	"rpkg.cc/apps/kerrigan/pkg/reckeeper/reckcfg"
	"rpkg.cc/apps/kerrigan/pkg/reckeeper/service"
	"rpkg.cc/apps/kerrigan/rcrai/rapi"
	"rpkg.cc/ecode"
	"rpkg.cc/infras/servers/api"
	"rpkg.cc/infras/servers/apiserver"
	"rpkg.cc/infras/services"
)

type MemberAPIs struct {
	apiserver.API
	apiserver.VersionedAPIModule
	service service.Service
}

func NewMemberAPIs(codeName, version string, service service.Service) *MemberAPIs {
	api := &MemberAPIs{}
	api.service = service
	api.SetCodeName(codeName)
	api.SetVersion(version)
	return api
}

func NewDefaultMemberAPIs() *MemberAPIs {
	var service = service.Service{Pager: &services.PagerInfo{}}
	return NewMemberAPIs(reckcfg.CodeName, reckcfg.APIVersion, service)
}

func (p *MemberAPIs) RegisterAPIs(r *gin.Engine) {
	p.GET(r, p.URL("/members/search"), p.SearchMembers)
	p.POST(r, p.URL("/member/create"), p.CreateMember)
	p.PUT(r, p.URL("/member/modify"), p.UpdateMember)
	p.PUT(r, p.URL("/member/push"), p.UpdateMember)
	p.GET(r, p.URL("/member/find/:id"), p.GetMember)
}

func (p *MemberAPIs) SearchMembers(c *gin.Context) api.Results {
	dt_id := c.Query("dt_id")
	site := c.Query("site")
	authuser := gwauth.GetAuthUser(c).RequireUser()
	models, total, err := p.service.SearchMembers(dt_id, site,
		services.DefaultPager(c),
		services.NewContext().SetAuth(authuser),
	)
	if err != nil {
		return rapi.SmartError(err, ecode.Internal).Apply(c)

	}
	return rapi.Resp().Set("result", models).Total(total).Apply(c)
}

func (p *MemberAPIs) CreateMember(c *gin.Context) api.Results {
	var model model.Member
	if _, err := api.BindArgs(c, &model); err != nil {
		return rapi.SmartError(err, ecode.InvalidArgument).Apply(c)
	}

	authuser := gwauth.GetAuthUser(c).RequireUser()
	result, _, err := p.service.CreateMember(&model, services.NewContext().SetAuth(authuser))
	if err != nil {
		return rapi.SmartError(err, ecode.Internal).Apply(c)
	}
	return rapi.Resp().Set("result", result).Apply(c)
}

func (p *MemberAPIs) UpdateMember(c *gin.Context) api.Results {
	var model model.Member
	if _, err := api.BindArgs(c, &model); err != nil {
		return rapi.SmartError(err, ecode.InvalidArgument).Apply(c)
	}

	authuser := gwauth.GetAuthUser(c).RequireUser()
	result, _, err := p.service.UpdateMember(&model, services.NewContext().SetAuth(authuser))
	if err != nil {
		return rapi.SmartError(err, ecode.Internal).Apply(c)

	}
	return rapi.Resp().Set("result", result).Apply(c)
}

func (p *MemberAPIs) GetMember(c *gin.Context) api.Results {
	id := c.Param("id")
	authuser := gwauth.GetAuthUser(c).RequireUser()
	result, err := p.service.GetMember(id, services.NewContext().SetAuth(authuser))
	if err != nil {
		return rapi.SmartError(err, ecode.Internal).Apply(c)
	}
	return rapi.Resp().Set("result", result).Apply(c)
}

func (p *MemberAPIs) MemberPush(c *gin.Context) api.Results {
	var model model.Member
	if _, err := api.BindArgs(c, &model); err != nil {
		return rapi.SmartError(err, ecode.InvalidArgument).Apply(c)
	}
	authuser := gwauth.GetAuthUser(c).RequireUser()
	result, err := p.service.MemberPush(&model, services.NewContext().SetAuth(authuser))
	if err != nil {
		return rapi.SmartError(err, ecode.Internal).Apply(c)
	}
	return rapi.Resp().Set("result", result).Apply(c)
}
