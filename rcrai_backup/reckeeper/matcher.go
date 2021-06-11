package uploader

import (
	"fmt"
	"time"

	"rpkg.cc/apps/kerrigan/countermgr/cmservice/edgedb"
	"rpkg.cc/apps/kerrigan/gwauth"
	"rpkg.cc/apps/kerrigan/pkg/audiocollector/aclog"
	"rpkg.cc/dev/slacker"
	"rpkg.cc/infra/utils/refs"
	"rpkg.cc/infras/servers"
	ss "rpkg.cc/infras/services"
	"rpkg.cc/log"

	"rpkg.cc/apps/kerrigan/pkg/audiocollector/acconfig"
	"rpkg.cc/apps/kerrigan/pkg/reckeeper/model"
	"rpkg.cc/apps/kerrigan/pkg/reckeeper/service"
)

// ! 注意，本文件中的代码和逻辑已经吸收并在reckeeper中实现，此文件不再使用。有问题找gaobo

type Matcher struct {
	servers.IntervalServer // enable IntervalServer
	service.Service
}

type UnmatchedAudio struct {
	LogID     int
	StaffID   string
	Timestamp int64
	Site      string
}

func NewMatcher() *Matcher {
	core := &Matcher{}
	core.SetInterval(3000).AddCore(core)
	return core
}

func (p *Matcher) OnStart() bool {
	log.Info("Start Matcher ...")
	return true
}

func (p *Matcher) DoInterval() bool {
	defer func() {
		if err := recover(); err != nil {
			slacker.ProcessPanicBottomLine("[Matcher.DoInterval]", err)
		}
	}()

	// TODO 考虑这了是否用数据库中的edgedevices，还是用内存中的？
	edgedevices, err := edgedb.FindAllEdgeDevicesDB("")
	if err != nil {
		return false
	}

	timeEnd := time.Now()
	tiemStart := timeEnd.Add(-time.Hour * 24 * 7)

	for _, edgedevice := range edgedevices {
		// fake a bid. // TODO need admin?
		fackAuth := &gwauth.AuthUser{}
		fackAuth.SetBid(edgedevice.GetBid())
		ctx := ss.NewContext().SetAuth(fackAuth)

		// * Main logic
		audioInfos, err := getUnmatched(&[]string{edgedevice.GetID()}, tiemStart, timeEnd)
		handleErr(err)
		uploaderC := acconfig.GetConfig().AudioCollector.Uploader
		for _, audioInfo := range audioInfos {
			// TODO 这里需要build一个有BID的场景；
			visits, total, err := p.Service.VisitsSearch(audioInfo.Site, audioInfo.StaffID, "",
				audioInfo.Timestamp-uploaderC.BeganTime,
				audioInfo.Timestamp+uploaderC.EndTime,
				false, nil, ctx,
			)
			handleErr(err)
			if total == 0 {
				continue
			}
			visit, err := findMatchedVisited(audioInfo.Timestamp, visits)
			handleErr(err)
			handleErr(matchCustomerId(&audioInfo, visit.CustomerId))
		}
	}

	return true
}

func (p *Matcher) OnShutdown() bool { return true }

func getUnmatched(deviceIDs *[]string, from time.Time, to time.Time) ([]UnmatchedAudio, error) {
	filter := &aclog.UploadLogFilter{
		TimeRangeFilter: ss.NewTimeRangeFilter().From(from).To(to),
		Match:           "unmatched",
	}
	models, _, err := aclog.GetUploadLogs(
		filter,
		nil,
		deviceIDs,
		"t_event desc",
		ss.Pager().Init(nil, nil, nil, refs.IntRef(1000)).Set(0, 1000),
		nil,
	)
	if err != nil {
		return nil, err
	}
	ret := make([]UnmatchedAudio, len(models))
	for i := range models {
		ret[i] = UnmatchedAudio{
			StaffID:   refs.StrVal(models[i].SPDID),
			Timestamp: models[i].EventTime.Unix(),
			LogID:     int(*models[i].ID),
			Site:      refs.StrVal(models[i].EdgeName),
		}
	}
	return ret, nil
}

func FindConversationId(staffId string, startTime time.Time) (string, error) {
	var bid = acconfig.GetConfig().AudioCollector.Uploader.Dealtape.Bid
	conversations, err := service.ConversationsSearch(bid, staffId, startTime.Add(-time.Minute*10), startTime.Add(-time.Minute*10), "all")
	if err != nil {
		return "", err
	}
	if len(conversations) == 0 {
		return "", fmt.Errorf("conversation not found")
	}
	var gap int64 = -1
	var conversationId = ""
	for i := range conversations {
		g := Abs(startTime.Unix(), conversations[i].StartTime.Unix())
		if gap < 0 || g <= gap {
			gap = g
			conversationId = conversations[i].ConversationId
		}
	}
	return conversationId, nil
}

func matchCustomerId(audioInfo *UnmatchedAudio, customerId string) error {
	conversationId, err := FindConversationId(audioInfo.StaffID, time.Unix(audioInfo.Timestamp, 0))
	if err != nil {
		return err
	}
	if acconfig.GetConfig().AudioCollector.Uploader.Debug {
		log.Infof("[matcher debug info]: dealtape： 将未匹配录音 (conversationId) %v 匹配到付记录 (customerId) %v", conversationId, customerId)
		log.Infof("[matcher debug info]: aclog: 标记记录 %v (audioInfo.LogID)为匹配 (customerId) %v ", audioInfo.LogID, customerId)
		return nil
	}
	var bid = acconfig.GetConfig().AudioCollector.Uploader.Dealtape.Bid
	_, err = service.ConversationModify(bid, conversationId, customerId)
	if err != nil {
		return err
	}
	err = aclog.MatchUploadRecord(audioInfo.LogID, customerId, nil)
	if err != nil {
		return err
	}
	return nil
}

func Abs(a, b int64) int64 {
	if a >= b {
		return a - b
	}
	return Abs(b, a)
}

func findMatchedVisited(timestamp int64, visits []*model.VisitDetail) (*model.Visit, error) {
	if len(visits) == 0 {
		return nil, fmt.Errorf("visit not found")
	}
	if len(visits) == 1 {
		return &visits[0].Visit, nil
	}
	var beforeAnchorTime int64 = -1
	var recentAnchorTime int64 = -1
	var beforeV *model.Visit = nil
	var recentV *model.Visit = nil

	for i := range visits {
		if visits[i].StartTime < timestamp && visits[i].StartTime > beforeAnchorTime {
			beforeAnchorTime = visits[i].StartTime
			beforeV = &visits[i].Visit
		}
		gap := Abs(visits[i].StartTime, timestamp)
		if gap < recentAnchorTime || recentAnchorTime < 0 {
			recentAnchorTime = gap
			recentV = &visits[i].Visit
		}
	}
	if beforeAnchorTime > 0 && beforeV != nil && (timestamp-beforeAnchorTime) < 3600 {
		return beforeV, nil
	}
	if recentV == nil {
		return nil, fmt.Errorf("visit not found")
	}
	return recentV, nil
}

func handleErr(e error) {
	if e != nil {
		panic(e)
	}
}
