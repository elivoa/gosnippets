package service

import (
	"fmt"
	"time"

	"rpkg.cc/apps/kerrigan/confs"
	"rpkg.cc/apps/kerrigan/server/dealtapeclient"
	"rpkg.cc/endpoint/zhijian"
	"rpkg.cc/log"
)

type StaffFilter struct {
	BranchID string
	Type     string
	Role     string
}

func NewStaffFilter(branchID, stafftype, role string) StaffFilter {
	return StaffFilter{
		BranchID: branchID,
		Type:     stafftype,
		Role:     role,
	}
}

func ListDealtapeStaffs(bid string) (staffs []*zhijian.Staff, err error) {
	dtConfig, ok := confs.GetGlobalConfig().KerriganConfig.DealtapeAccounts[bid]
	if !ok || dtConfig == nil {
		err = fmt.Errorf("found no dt config", bid)
		return
	}

	client := dealtapeclient.NewZhijianClient(dtConfig)
	log.Info("DealtapeClientConfig", dtConfig)
	staffs, err = client.ListStaffs()
	return
}

func CustomersUpload(bid string, customers []*zhijian.Customer) (err error) {
	dtConfig, ok := confs.GetGlobalConfig().KerriganConfig.DealtapeAccounts[bid]
	if !ok || dtConfig == nil {
		err = fmt.Errorf("found no dt config", bid)
		return
	}
	client := dealtapeclient.NewZhijianClient(dtConfig)
	err = client.CustomersUpload(customers)
	return
}

func CustomerUpload(bid, source_id, name, phone string) (err error) {
	dtConfig, ok := confs.GetGlobalConfig().KerriganConfig.DealtapeAccounts[bid]
	if !ok || dtConfig == nil {
		err = fmt.Errorf("found no dt config", bid)
		return
	}
	client := dealtapeclient.NewZhijianClient(dtConfig)
	var customer = zhijian.Customer{SourceId: source_id, Name: name, Phone: phone}
	err = client.CustomerUpload(&customer)
	return
}

func ConversationsSearch(bid, staffId string, BeganTime time.Time, EndTime time.Time, sheetName string) (result []*zhijian.Conversation, err error) {
	dtConfig, ok := confs.GetGlobalConfig().KerriganConfig.DealtapeAccounts[bid]
	if !ok || dtConfig == nil {
		err = fmt.Errorf("found no dt config", bid)
		return
	}
	client := dealtapeclient.NewZhijianClient(dtConfig)

	result, err = client.ConversationsSearch(staffId, BeganTime, EndTime, sheetName)
	return
}

func ConversationModify(bid, conversationId, customerId string) (result bool, err error) {
	dtConfig, ok := confs.GetGlobalConfig().KerriganConfig.DealtapeAccounts[bid]
	if !ok || dtConfig == nil {
		err = fmt.Errorf("found no dt config", bid)
		return
	}
	client := dealtapeclient.NewZhijianClient(dtConfig)
	result, err = client.ConversationModify(conversationId, customerId)
	return
}
