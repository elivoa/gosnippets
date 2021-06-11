package badapi

import (
	"fmt"
	"io/ioutil"
	"math"
	"path/filepath"
	"regexp"
	"strconv"
	"strings"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/tealeg/xlsx/v3"
	"rpkg.cc/apps/kerrigan/confs"
	"rpkg.cc/apps/kerrigan/gwauth"
	"rpkg.cc/apps/kerrigan/pkg/reckeeper/model"
	"rpkg.cc/apps/kerrigan/pkg/reckeeper/reckcfg"
	"rpkg.cc/apps/kerrigan/pkg/reckeeper/service"
	"rpkg.cc/apps/kerrigan/rcrai/rapi"
	"rpkg.cc/ecode"
	"rpkg.cc/endpoint/zhijian"
	"rpkg.cc/infras/servers/api"
	"rpkg.cc/infras/servers/apiserver"
	"rpkg.cc/infras/services"
	"rpkg.cc/log"
)

type VisitAPIs struct {
	apiserver.API
	apiserver.VersionedAPIModule
	service service.Service
	config  confs.Config
}

func NewVisitAPIs(codeName, version string, service service.Service) *VisitAPIs {
	api := &VisitAPIs{}
	api.service = service
	api.config = *confs.GetGlobalConfig()

	api.SetCodeName(codeName)
	api.SetVersion(version)
	return api
}

func NewDefaultVisitAPIs() *VisitAPIs {
	var service = service.Service{Pager: &services.PagerInfo{}}
	return NewVisitAPIs(reckcfg.CodeName, reckcfg.APIVersion, service)
}

func (p *VisitAPIs) RegisterAPIs(r *gin.Engine) {
	p.GET(r, p.URL("/visits/search"), p.GetVisitSearch(false, false))
	p.GET(r, p.URL("/visits/search/detail"), p.GetVisitSearch(true, false))
	p.GET(r, p.URL("/visits/search/row"), p.GetVisitSearch(true, true))
	p.POST(r, p.URL("/visit/create"), p.CreateVisit)
	p.PUT(r, p.URL("/visit/modify"), p.UpdateVisit)
	p.PUT(r, p.URL("/visit/push"), p.VisitPush)
	r.GET(p.URL("/visit/find/:id"), p.VisitFind)
	r.POST(p.URL("/visits/upload"), p.VisitsUpload)
}

func (p *VisitAPIs) CreateVisit(c *gin.Context) api.Results {
	var model model.Visit
	if _, err := api.BindArgs(c, &model); err != nil {
		return rapi.SmartError(err, ecode.InvalidArgument).Apply(c)
	}

	authuser := gwauth.GetAuthUser(c).RequireUser()
	result, _, err := p.service.CreateVisit(&model, services.NewContext().SetAuth(authuser))
	if err != nil {
		return rapi.SmartError(err, ecode.Internal).Apply(c)
	}
	return rapi.Resp().Set("result", result).Apply(c)
}

func (p *VisitAPIs) UpdateVisit(c *gin.Context) api.Results {
	var model model.Visit
	if _, err := api.BindArgs(c, &model); err != nil {
		return rapi.SmartError(err, ecode.InvalidArgument).Apply(c)

	}
	authuser := gwauth.GetAuthUser(c).RequireUser()
	result, _, err := p.service.UpdateVisit(&model, services.NewContext().SetAuth(authuser))
	if err != nil {
		return rapi.SmartError(err, ecode.Internal).Apply(c)
	}
	return rapi.Resp().Set("result", result).Apply(c)
}

func (p *VisitAPIs) VisitPush(c *gin.Context) api.Results {
	var model model.Visit
	if _, err := api.BindArgs(c, &model); err != nil {
		return rapi.SmartError(err, ecode.InvalidArgument).Apply(c)
	}

	authuser := gwauth.GetAuthUser(c).RequireUser()
	result, err := p.service.VisitPush(&model, services.NewContext().SetAuth(authuser))
	if err != nil {
		return rapi.SmartError(err, ecode.Internal).Apply(c)
	}
	return rapi.Resp().Set("result", result).Apply(c)
}

func (p *VisitAPIs) GetVisitSearch(withDetail bool, transfer2row bool) func(c *gin.Context) api.Results {
	return func(c *gin.Context) api.Results {
		site := c.Query("site")
		dt_id := c.Query("dt_id")
		customer_id := c.Query("customer_id")
		began_time, err := strconv.ParseInt(c.Query("began_time"), 10, 64)
		if err != nil {
			began_time = 0
		}
		end_time, err := strconv.ParseInt(c.Query("end_time"), 10, 64)
		if err != nil {
			end_time = 1000000000000
		}

		authuser := gwauth.GetAuthUser(c).RequireUser()

		models, total, err := p.service.VisitsSearch(
			site, dt_id, customer_id, began_time, end_time, withDetail,
			services.DefaultPager(c),
			services.NewContext().SetAuth(authuser),
		)
		if err != nil {
			return rapi.SmartError(err, ecode.Internal).Apply(c)
		}

		if !transfer2row {
			return rapi.Resp().Set("result", &models).Total(total).Apply(c)
		} else {
			var items []*model.VisitRow
			for _, model := range models {
				items = append(items, model.ToVisitRow())
			}
			return rapi.Resp().Set("result", items).Total(total).Apply(c)
		}
	}
}

func (p *VisitAPIs) VisitFind(c *gin.Context) {
	id := c.Param("id")

	result, err := p.service.VisitFind(id)
	if err != nil {
		rapi.SmartError(err, ecode.Internal).Apply(c)
		return
	}
	rapi.Resp().Set("result", result).Apply(c)
}

// !----------------------------------------------------------------

// VisitsUpload 是最主要的上传接口。
func (p *VisitAPIs) VisitsUpload(c *gin.Context) {

	// 只能使用登录的bid和token么？
	authuser := gwauth.GetAuthUser(c).RequireUser()
	ctxs := services.NewContext().SetAuth(authuser)
	// dtc := reckhelper.GetDefaultDTC(authuser)
	var bid = authuser.GetBidSafe() // dtc.Config.Bid //

	// save xls
	file, err := c.FormFile("file")
	if err != nil {
		var line = fmt.Sprintf("接收上传拜访记录xls出错: %v", err)
		rapi.SmartErrorf(fmt.Errorf(line), ecode.InvalidArgument, line).Apply(c)
		return
	}
	extension := filepath.Ext(file.Filename)
	if extension != ".xls" && extension != ".xlsx" {
		var line = fmt.Sprintf("xls文件格式 %v 错误: %v", extension, err)
		rapi.SmartErrorf(fmt.Errorf(line), ecode.InvalidArgument, line).Apply(c)
		return
	}
	if file.Size > 1024*1024*100 {
		var line = fmt.Sprintf("文件超过100m: %v", err)
		rapi.SmartErrorf(fmt.Errorf(line), ecode.InvalidArgument, line).Apply(c)
		return
	}

	// var dst = p.config.KerriganConfig.Reckeeper.VistiDir + "/" + file.Filename
	var dst = filepath.Join(confs.GetReckeeperCfg().VistiBackupDir, file.Filename)
	c.SaveUploadedFile(file, dst)
	log.Info("'%s' 上传成功!,保存在", file.Filename)

	total, succeed, results, err := syncXls(bid, p, dst, ctxs)
	if err != nil {
		succeed = 0
		rapi.Resp().Set("total", total).Set("succeed", succeed).Set("warning", err.Error()).Apply(c)
		return
	}

	files, err := ioutil.ReadDir(p.config.KerriganConfig.Reckeeper.VistiBackupDir)
	if err != nil {
		log.Fatal(err)
		succeed = 0
		rapi.Resp().Set("total", total).Set("succeed", succeed).Set("warning", "无法读取到访记录文件夹："+p.config.KerriganConfig.Reckeeper.VistiBackupDir).Apply(c)
		return
	} else {
		for i := 0; i < len(files); i += 1 {
			var name = files[i].Name()
			if time.Now().Unix()-files[i].ModTime().Unix() < 3600*3 && (strings.HasSuffix(name, ".xls") || strings.HasSuffix(name, ".xlsx")) && name != file.Filename {
				total1, succeed1, results1, err1 := syncXls(bid, p, dst, ctxs)
				if err1 != nil {
					succeed1 = 0
					log.Error("解析旧的上传记录失败", name, total1, succeed1, err1, results1)
				}
			}
		}
	}

	rapi.Resp().Set("total", total).Set("succeed", succeed).Set("result", results).Apply(c)
}

func syncXls(bid string, p *VisitAPIs, path string, ctx *services.Context) (
	total, succeed int, results []model.VisitRow, err error) {

	total, succeed, visits, customers, members, times, err := parseXls(path)
	if err != nil {
		return
	}
	staffs, existStaffs, err := fetchDtStaffs(bid)
	if err != nil {
		return
	}
	_, results, err = syncReckeeper(p, visits, customers, members, ctx)
	if err != nil {
		return
	}
	_, err = syncDt(bid, p, staffs, existStaffs, visits, customers, members, times, ctx)
	if err != nil {
		return
	}
	return
}

func fetchDtStaffs(bid string) (staffs []*zhijian.Staff, existStaffs map[string]bool, err error) {
	staffs0, err := service.ListDealtapeStaffs(bid)
	if err != nil {
		err = fmt.Errorf("ListDealtapeStaffs 失败 %s", err)
		return
	}
	//var staffs = []zhijian.Staff{}
	for i := 0; i < len(staffs0); i += 1 {
		if staffs0[i].SourceId != "" {
			staffs = append(staffs, staffs0[i])
		}
	}
	existStaffs = map[string]bool{}
	for i := 0; i < len(staffs); i += 1 {
		a := staffs[i]
		b := a.SourceId
		existStaffs[b] = true
	}
	log.Info(" 找到dt员工 ", len(staffs0), "有sourceId", len(staffs), "个员工", staffs)
	return
}

func syncDt(bid string, p *VisitAPIs, staffs []*zhijian.Staff, existStaffs map[string]bool,
	visits []*model.Visit, customers []*model.Customer, members []*model.Member, times []time.Time,
	ctx *services.Context) (succeed int, err error) {

	// 更新dt
	total := len(visits)
	for i := 0; i < total; i += 1 {
		var visit = visits[i]
		var customer = customers[i]
		var member = members[i]
		var startTime = times[i]

		// dt员工必需存在
		_, ok := existStaffs[member.DtId]
		if !ok {
			err = fmt.Errorf(" %v 员工不存在 ", *member)
			return
		}
		// 推送客户  必须逐行传
		e := service.CustomerUpload(bid, customer.CustomerId, customer.Name, customer.Phone)
		if e != nil {
			err = fmt.Errorf("推送dt客户 %v 出错: %v", *customer, e)
			return
		}
		// make sure visit uploaded to reckeeper, and get a candicate range for conversations
		vis, _, e := p.service.VisitsSearch("", visit.DTSourceID, "", startTime.Add(-time.Hour*8).Unix(), startTime.Add(time.Hour*8).Unix(), false, nil, ctx)
		if e != nil {
			err = fmt.Errorf("搜索到访记录 出错: %v", e)
			return
		}
		var staffId = "" //get StaffId from dt
		for l := 0; l < len(staffs); l += 1 {
			if staffs[l].SourceId == visit.DTSourceID {
				staffId = staffs[l].StaffId
				break
			}
		}
		if staffId == "" {
			err = fmt.Errorf("未找到dt匹配员工 出错: %v", visit)
			return
		}
		conversations, e := service.ConversationsSearch(bid, staffId, startTime.Add(-time.Hour*8), startTime.Add(time.Hour*8), "all")
		if e != nil {
			err = fmt.Errorf("搜索dt会话 出错: %v", e)
			return
		}
		log.Info(" 找到附近拜访记录 VisitsSearch ", len(vis), " 找到附近会话 ConversationsSearch ", len(conversations), conversations)
		for _, conv := range conversations {
			var near = 0
			var diff = math.Abs(float64((conv.StartTime.Unix() - vis[0].StartTime)))
			for j := 1; j < len(vis); j += 1 {
				var d = math.Abs(float64((conv.StartTime.Unix() - vis[j].StartTime)))
				if d < diff {
					diff = d
					near = j
				}
			}
			var customerId = vis[near].CustomerId
			log.Info("匹配到最近拜访记录", near, vis[near], conv)
			if conv.CustomerId != customerId {
				log.Info("===会话安排客户  ", conv, near, vis[near])
				var re, e = service.ConversationModify(bid, conv.ConversationId, customerId)
				if e != nil || !re {
					err = fmt.Errorf("会话更新客户出错: %v", e)
					return
				}
			}
		}
		succeed += 1
	}
	return
}

func syncReckeeper(p *VisitAPIs, visits []*model.Visit,
	customers []*model.Customer, members []*model.Member, ctx *services.Context) (
	succeed int, results []model.VisitRow, err error) {

	total := len(visits)

	// * 这里拿到一个文件的所有上传记录，先对member和customer去重，然后尝试新建。
	// membersmap := map[string]*model.Member{}
	// customermap := map[string]*model.Customer{}

	for i := 0; i < total; i += 1 {
		// 推送reckeeper
		m, e := p.service.MemberPush(members[i], ctx)
		if e != nil {
			log.Error("reckeeper 推送员工 出错:", *members[i], e)
			err = e
			continue
		}
		cm, e := p.service.CustomerPush(customers[i], ctx)
		if e != nil {
			log.Error("reckeeper 推送客户 出错: ", *customers[i], e)
			err = e
			continue
		}

		// * 判断上传记录是否重复，再进行上传操作。

		v, e := p.service.VisitPush(visits[i], ctx)
		if e != nil {
			err = fmt.Errorf("reckeeper 推送拜访 出错: %v %w", visits[i], e)
			return
		}
		var result = model.VisitRow{Visit: *v, CustomerName: cm.Name, CustomerPhone: cm.Phone, MemberName: m.Name, MemberPhone: m.Phone}
		results = append(results, result)
	}
	succeed = len(results)

	return
}

func parseXls(path string) (total, succeed int, visits []*model.Visit, customers []*model.Customer, members []*model.Member, times []time.Time, err error) {
	total = -1
	wb, e := xlsx.OpenFile(path)
	if e != nil || len(wb.Sheets) == 0 || wb.Sheets[0].MaxRow == 0 {
		err = fmt.Errorf("解析上传拜访记录xls出错: %v", e)
		return
	}
	var rows = []*xlsx.Row{}
	sh := wb.Sheets[0]
	log.Info("正在解析表格", path, sh.Name)
	row, e := sh.Row(0)
	if e != nil {
		err = fmt.Errorf("解析上传拜访记录xls行 header 出错: %v", e)
		return
	}
	indexes, e := parseHead(row)
	if e != nil {
		err = fmt.Errorf("解析上传拜访记录xls行 header 出错: %v", e)
		return
	}

	// non-blank rows
	for j := 1; j < sh.MaxRow; j += 1 {
		row, e := sh.Row(j)
		if e != nil {
			err = fmt.Errorf("解析上传拜访记录xls行 %d 出错: %v", j, e)
			return
		}
		var length = 0
		for k := 0; k < sh.MaxCol; k += 1 {
			length += len(row.GetCell(k).String())
		}
		if length > 0 {
			rows = append(rows, row)
		}

	}
	total = len(rows)

	for j := 0; j < total; j += 1 {
		customer, member, visit, startTime, e := rowVisitor(rows[j], indexes)
		if e != nil {
			err = fmt.Errorf("解析上传拜访记录xls行结构 %d 出错: %v", j+1, e)
			return
		}
		visits = append(visits, visit)
		customers = append(customers, customer)
		members = append(members, member)
		times = append(times, startTime)
		succeed = len(visits)
	}
	return
}

func parseHead(row *xlsx.Row) ([]int, error) {
	var headerIndex = map[string]int{}
	for i := 0; i < row.Sheet.MaxCol; i += 1 {
		headerIndex[row.GetCell(i).String()] = i
	}
	var err error = nil
	var indexes = []int{}
	var fields = []string{"客户姓名", "客户电话", "拜访时间", "销售姓名", "销售电话", "所属项目"}
	for i := 0; i < len(fields); i += 1 {
		idx, ok := headerIndex[fields[i]]
		if !ok {
			err = fmt.Errorf("header not found %v", fields[i])
			return indexes, err
		}
		indexes = append(indexes, idx)
	}
	return indexes, err
}

//客户姓名	客户电话	拜访时间	销售姓名	销售电话	所属项目
//张先生	131****0702	2021/1/13 19:50	张三	138****8254	翡翠西湖
func rowVisitor(row *xlsx.Row, indexes []int) (*model.Customer, *model.Member, *model.Visit, time.Time, error) {
	var customer model.Customer
	var member model.Member
	var e0 error
	customer.Name = row.GetCell(indexes[0]).String()
	customer.Phone = row.GetCell(indexes[1]).String()
	customer.CustomerId = customer.Phone

	member.Name = row.GetCell(indexes[3]).String()
	member.Phone = row.GetCell(indexes[4]).String()
	member.Site = row.GetCell(indexes[5]).String()
	member.DtId = member.Site + "_" + member.Phone
	if !validName(customer.Name) || !validPhone(customer.Phone) || !validName(member.Name) || !validPhone(member.Phone) || !validName(member.Site) {
		e0 = fmt.Errorf("cell format not allowed")
	}
	start_time, e := row.GetCell(indexes[2]).GetTime(false)
	if e0 != nil || e != nil {
		err := fmt.Errorf("推送拜访记录xls行 %v 出错: %v  %v", &row, e0, e)
		return nil, nil, nil, time.Now(), err
	}
	var visit = model.Visit{DTSourceID: member.DtId, CustomerId: customer.CustomerId, StartTime: start_time.Unix(), Site: member.Site}
	return &customer, &member, &visit, start_time, nil
}

func validPhone(s string) bool {
	if len(s) == 0 {
		return false
	}
	var reg = regexp.MustCompile(`^[0-9_\*]+$`)
	return reg.MatchString(s)
}

func validName(s string) bool {
	if len(s) == 0 {
		return false
	}
	// var reg = regexp.MustCompile(`^[\p{Han}A-Za-z0-9_\*]+$`)
	var reg = regexp.MustCompile(`^[\p{Han}\x20-\x7E\t\w\n\r\-·=【】、；’‘，。~！￥…（）—「」|：“”《》？]+$`)
	return reg.MatchString(s)
}
