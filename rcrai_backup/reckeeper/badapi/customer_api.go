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

type CustomerAPIs struct {
	apiserver.API
	apiserver.VersionedAPIModule
	service service.Service
}

func NewCustomerAPIs(codeName, version string, service service.Service) *CustomerAPIs {
	api := &CustomerAPIs{}
	api.service = service
	api.SetCodeName(codeName)
	api.SetVersion(version)
	return api
}

func NewDefaultCustomerAPIs() *CustomerAPIs {
	var service = service.Service{Pager: &services.PagerInfo{}}
	return NewCustomerAPIs(reckcfg.CodeName, reckcfg.APIVersion, service)
}

func (p *CustomerAPIs) RegisterAPIs(r *gin.Engine) {
	p.GET(r, p.URL("/customers/search"), p.SearchCustomers)
	p.POST(r, p.URL("/customer/create"), p.CreateCustomer)
	r.PUT(p.URL("/customer/modify"), p.UpdateCustomer)
	p.PUT(r, p.URL("/customer/push"), p.CustomerPush)
	p.GET(r, p.URL("/customer/find/:id"), p.GetCustomer)
}

func (p *CustomerAPIs) SearchCustomers(c *gin.Context) api.Results {
	authuser := gwauth.GetAuthUser(c).RequireUser()

	customer_id := c.Query("customer_id")
	assigned_staff := c.Query("assigned_staff")
	models, total, err := p.service.SearchCustomers(customer_id, assigned_staff,
		services.DefaultPager(c),
		services.NewContext().SetAuth(authuser),
	)
	if err != nil {
		return rapi.SmartError(err, ecode.Internal).Apply(c)
	}
	return rapi.Resp().Set("result", models).Total(total).Apply(c)
}

func (p *CustomerAPIs) CreateCustomer(c *gin.Context) api.Results {
	var model model.Customer
	if _, err := api.BindArgs(c, &model); err != nil {
		return rapi.SmartError(err, ecode.InvalidArgument).Apply(c)
	}

	authuser := gwauth.GetAuthUser(c).RequireUser()
	result, _, err := p.service.CreateCustomer(&model, services.NewContext().SetAuth(authuser))
	if err != nil {
		return rapi.SmartError(err, ecode.Internal).Apply(c)
	}
	return rapi.Resp().Set("result", result).Apply(c)
}

func (p *CustomerAPIs) UpdateCustomer(c *gin.Context) {
	var model model.Customer
	if _, err := api.BindArgs(c, &model); err != nil {
		rapi.SmartError(err, ecode.InvalidArgument).Apply(c)
		return
	}

	authuser := gwauth.GetAuthUser(c).RequireUser()
	result, _, err := p.service.UpdateCustomer(&model, services.NewContext().SetAuth(authuser))
	if err != nil {
		rapi.SmartError(err, ecode.Internal).Apply(c)
		return
	}
	rapi.Resp().Set("result", result).Apply(c)
}

// 尝试创建customer，如果有了就不创建了。这里有优化空间。
func (p *CustomerAPIs) CustomerPush(c *gin.Context) api.Results {
	var model model.Customer
	if _, err := api.BindArgs(c, &model); err != nil {
		return rapi.SmartError(err, ecode.InvalidArgument).Apply(c)
	}

	authuser := gwauth.GetAuthUser(c).RequireUser()
	result, err := p.service.CustomerPush(&model, services.NewContext().SetAuth(authuser))
	if err != nil {
		return rapi.SmartError(err, ecode.Internal).Apply(c)
	}
	return rapi.Resp().Set("result", result).Apply(c)
}

func (p *CustomerAPIs) GetCustomer(c *gin.Context) api.Results {
	id := c.Param("id")
	authuser := gwauth.GetAuthUser(c).RequireUser()
	result, err := p.service.GetCustomer(id, services.NewContext().SetAuth(authuser))
	if err != nil {
		return rapi.SmartError(err, ecode.Internal).Apply(c)
	}
	return rapi.Resp().Set("result", result).Apply(c)
}
