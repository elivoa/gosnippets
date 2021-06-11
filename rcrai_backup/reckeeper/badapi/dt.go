package badapi

import (
	"fmt"
	"time"

	"github.com/araddon/dateparse"
	"github.com/gin-gonic/gin"
	"rpkg.cc/apps/kerrigan/gwauth"
	"rpkg.cc/apps/kerrigan/pkg/reckeeper/reckcfg"
	"rpkg.cc/apps/kerrigan/pkg/reckeeper/reckhelper"
	"rpkg.cc/apps/kerrigan/pkg/reckeeper/service"
	"rpkg.cc/apps/kerrigan/rcrai/rapi"
	"rpkg.cc/ecode"
	"rpkg.cc/infras/servers/apiserver"
	"rpkg.cc/log"
)

type DtAPIs struct {
	apiserver.VersionedAPIModule
}

func NewDtAPIs(codeName, version string) *DtAPIs {
	api := &DtAPIs{}
	api.SetCodeName(codeName)
	api.SetVersion(version)
	return api
}

func NewDefaultDtAPIs() *DtAPIs {
	return NewDtAPIs(reckcfg.CodeName, reckcfg.APIVersion)
}

func (p *DtAPIs) RegisterAPIs(r *gin.Engine) {
	// r.GET(p.URL("/dt/staffs/search"), p.DtStaffsSearch)
	r.GET(p.URL("/dt/conversations/search"), p.DtConversationsSearch)
	r.PUT(p.URL("/dt/conversations/modify"), p.DtConversationModify)
	r.POST(p.URL("/dt/customer/upload"), p.DtCustomerUpload)
}

// func (p *DtAPIs) DtStaffsSearch(c *gin.Context) {
// 	authuser := gwauth.GetAuthUser(c).RequireUser()
// 	dtc := reckhelper.GetDefaultDTC(authuser)
// 	var bid = dtc.Config.Bid
// 	staffs, err := service.ListDealtapeStaffs(bid)
// 	if err != nil {
// 		rapi.SmartError(err, ecode.Internal).Apply(c)
// 		return
// 	}
// 	rapi.Resp().Set("staffs", staffs).Total(int64(len(staffs))).Apply(c)
// }

func (p *DtAPIs) DtConversationsSearch(c *gin.Context) {
	staffId := c.Query("staffId") // authed | candidate | all
	BeganTimeLine := c.Query("BeganTime")
	EndTimeLine := c.Query("EndTime")
	BeganTime, err := dateparse.ParseLocal(BeganTimeLine)
	if err != nil {
		log.Info("\n bad day: ", BeganTimeLine, err)
		BeganTime = time.Date(2000, 0, 0, 0, 0, 0, 0, time.Local)
	}
	EndTime, err := dateparse.ParseLocal(EndTimeLine)
	if err != nil {
		log.Info("\n bad day: ", EndTimeLine, err)
		EndTime = time.Date(2030, 0, 0, 0, 0, 0, 0, time.Local)
	}

	sheetName := c.Query("sheetName")
	authuser := gwauth.GetAuthUser(c).RequireUser()
	dtc := reckhelper.GetDefaultDTC(authuser)
	var bid = dtc.Config.Bid
	items, err := service.ConversationsSearch(bid, staffId, BeganTime, EndTime, sheetName)
	if err != nil {
		rapi.SmartError(err, ecode.Internal).Apply(c)
		return
	}
	rapi.Resp().Set("conversations", items).Total(int64(len(items))).Apply(c)
}

func (p *DtAPIs) DtConversationModify(c *gin.Context) {
	conversationId := c.Query("conversationId")
	customerId := c.Query("customerId")
	authuser := gwauth.GetAuthUser(c).RequireUser()
	dtc := reckhelper.GetDefaultDTC(authuser)
	var bid = dtc.Config.Bid
	items, err := service.ConversationModify(bid, conversationId, customerId)
	if err != nil {
		rapi.SmartError(err, ecode.Internal).Apply(c)
		return
	}
	rapi.Resp().Set("result", items).Apply(c)
}

func (p *DtAPIs) DtCustomerUpload(c *gin.Context) {
	source_id := c.Query("source_id")
	name := c.Query("name")
	phone := c.Query("phone")
	authuser := gwauth.GetAuthUser(c).RequireUser()
	dtc := reckhelper.GetDefaultDTC(authuser)
	var bid = dtc.Config.Bid
	err := service.CustomerUpload(bid, source_id, name, phone)
	if err != nil {
		var line = fmt.Sprintf("上传客户出错: %v", err)
		rapi.SmartErrorf(fmt.Errorf(line), ecode.InvalidArgument, line).Apply(c)
		return
	}
	rapi.Resp().Set("result", "ok").Apply(c)
}
