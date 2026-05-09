package adminnotify

import "github.com/openimsdk/tools/errs"

// SendID/Sender 是产品冻结字段。
// 这类“管理员后台补发通知”必须统一使用这里的固定值，业务服务侧不要再本地重定义。
const (
	SendID = "10001"
	Sender = "系统管理员"
)

const (
	EventManualNotification        = "admin_manual_notification"
	EventLoginOtherDevice          = "login_other_device"
	EventAdminPublishVersion       = "admin_publish_version"
	EventFriendApply               = "friend_apply"
	EventFriendApplyRejected       = "friend_apply_rejected"
	EventFriendApplyApproved       = "friend_apply_approved"
	EventFriendApplyHandleApproved = "friend_apply_handle_approved"
	EventFriendApplyHandleRejected = "friend_apply_handle_rejected"
	EventGroupApply                = "group_apply"
	EventGroupMemberInvited        = "group_member_invited"
	EventGroupApplyRejected        = "group_apply_rejected"
	EventGroupApplyApproved        = "group_apply_approved"
	EventGroupApplyHandleApproved  = "group_apply_handle_approved"
	EventGroupApplyHandleRejected  = "group_apply_handle_rejected"
	EventPrivacyPhoneApply         = "privacy_phone_apply"
	EventPrivacyEmailApply         = "privacy_email_apply"
	EventPrivacyHandleApproved     = "privacy_handle_approved"
	EventPrivacyHandleRejected     = "privacy_handle_rejected"
	EventPrivacyApplyApproved      = "privacy_apply_approved"
	EventPrivacyApplyRejected      = "privacy_apply_rejected"
)

const (
	CodeManualNotification        int32 = 10001
	CodeLoginOtherDevice          int32 = 10002
	CodeAdminPublishVersion       int32 = 10003
	CodeFriendApply               int32 = 10101
	CodeFriendApplyRejected       int32 = 10102
	CodeFriendApplyApproved       int32 = 10103
	CodeFriendApplyHandleApproved int32 = 10104
	CodeFriendApplyHandleRejected int32 = 10105
	CodeGroupApply                int32 = 10201
	CodeGroupMemberInvited        int32 = 10202
	CodeGroupApplyRejected        int32 = 10203
	CodeGroupApplyApproved        int32 = 10204
	CodeGroupApplyHandleApproved  int32 = 10205
	CodeGroupApplyHandleRejected  int32 = 10206
	CodePrivacyPhoneApply         int32 = 10301
	CodePrivacyEmailApply         int32 = 10302
	CodePrivacyHandleApproved     int32 = 10303
	CodePrivacyHandleRejected     int32 = 10304
	CodePrivacyApplyApproved      int32 = 10305
	CodePrivacyApplyRejected      int32 = 10306
)

const (
	TitleManualNotification  = "系统通知"
	TitleLoginOtherDevice    = "系统安全提醒"
	TitleAdminPublishVersion = "功能更新公告"
	TitleFriendApply         = "好友申请"
	TitleFriendApplyRejected = "好友申请已拒绝"
	TitleFriendApplyApproved = "好友申请已通过"
	TitleGroupApply          = "群聊申请"
	TitleGroupMemberInvited  = "你被邀请加入群聊"
	TitleGroupApplyRejected  = "入群申请已拒绝"
	TitleGroupApplyApproved  = "你被邀请加入群聊"
	TitlePrivacyApply        = "隐私查看请求"
	TitlePrivacyApplyReject  = "隐私查看请求已拒绝"
)

const (
	ContextLoginOtherDeviceFallback     = "检测到您的账号在新设备登录，若非本人操作请立即修改密码并开启两步验证。"
	ContextAdminPublishVersionFallback  = "全新的「派友派对」设计系统已上线，立即体验更多互动玩法。"
	ContextTemplateFriendApply          = "%s 请求添加你为好友：%s"
	ContextTemplateFriendApplyApproved  = "%s 已同意了您的好友申请，现在可以开始聊天了"
	ContextTemplateFriendApplyRejected  = "%s 已拒绝了您的好友申请"
	ContextTemplateFriendHandleApproved = "%s 请求添加你为好友：%s"
	ContextTemplateFriendHandleRejected = "%s 请求添加你为好友：%s"
	ContextTemplateGroupApply           = "%s 请求加入 %s，%s"
	ContextTemplateGroupInvited         = "你已被邀请加入 %s"
	ContextTemplateGroupApplyRejected   = "管理员已拒绝您加入 %s 的申请"
	ContextTemplateGroupHandleApproved  = "%s 请求加入 %s"
	ContextTemplateGroupHandleRejected  = "%s 请求加入 %s"
	ContextTemplatePrivacyApply         = "%s 请求查看您的%s信息"
	ContextTemplatePrivacySelfPassed    = "%s 请求查看您的%s信息"
	ContextTemplatePrivacySelfRejected  = "%s 请求查看您的%s信息"
	ContextPrivacyApplyApproved         = "对方通过了您的手机号/邮箱查看申请"
	ContextPrivacyApplyRejected         = "对方已拒绝了您的手机号/邮箱查看申请"
)

const (
	StatusPending  = "pending"
	StatusApproved = "approved"
	StatusRejected = "rejected"
	StatusSuccess  = "success"
)

const (
	CategorySystem = "system"
	CategoryFriend = "friend"
	CategoryGroup  = "group"
)

type EventSpec struct {
	Category      string
	Code          int32
	DefaultTitle  string
	DefaultStatus string
}

// eventSpecs 是 event -> code/title/status 的唯一真源。
// chat-admin 和 open-im-server 都必须直接复用这张表，不能各自维护副本，
// 否则协议漂移只是时间问题。
var eventSpecs = map[string]EventSpec{
	EventManualNotification:        {Category: CategorySystem, Code: CodeManualNotification, DefaultTitle: TitleManualNotification, DefaultStatus: StatusSuccess},
	EventLoginOtherDevice:          {Category: CategorySystem, Code: CodeLoginOtherDevice, DefaultTitle: TitleLoginOtherDevice, DefaultStatus: StatusSuccess},
	EventAdminPublishVersion:       {Category: CategorySystem, Code: CodeAdminPublishVersion, DefaultTitle: TitleAdminPublishVersion, DefaultStatus: StatusSuccess},
	EventFriendApply:               {Category: CategoryFriend, Code: CodeFriendApply, DefaultTitle: TitleFriendApply, DefaultStatus: StatusPending},
	EventFriendApplyRejected:       {Category: CategoryFriend, Code: CodeFriendApplyRejected, DefaultTitle: TitleFriendApplyRejected, DefaultStatus: StatusRejected},
	EventFriendApplyApproved:       {Category: CategoryFriend, Code: CodeFriendApplyApproved, DefaultTitle: TitleFriendApplyApproved, DefaultStatus: StatusApproved},
	EventFriendApplyHandleApproved: {Category: CategoryFriend, Code: CodeFriendApplyHandleApproved, DefaultTitle: TitleFriendApply, DefaultStatus: StatusApproved},
	EventFriendApplyHandleRejected: {Category: CategoryFriend, Code: CodeFriendApplyHandleRejected, DefaultTitle: TitleFriendApply, DefaultStatus: StatusRejected},
	EventGroupApply:                {Category: CategoryGroup, Code: CodeGroupApply, DefaultTitle: TitleGroupApply, DefaultStatus: StatusPending},
	EventGroupMemberInvited:        {Category: CategoryGroup, Code: CodeGroupMemberInvited, DefaultTitle: TitleGroupMemberInvited, DefaultStatus: StatusApproved},
	EventGroupApplyRejected:        {Category: CategoryGroup, Code: CodeGroupApplyRejected, DefaultTitle: TitleGroupApplyRejected, DefaultStatus: StatusRejected},
	EventGroupApplyApproved:        {Category: CategoryGroup, Code: CodeGroupApplyApproved, DefaultTitle: TitleGroupApplyApproved, DefaultStatus: StatusApproved},
	EventGroupApplyHandleApproved:  {Category: CategoryGroup, Code: CodeGroupApplyHandleApproved, DefaultTitle: TitleGroupApply, DefaultStatus: StatusApproved},
	EventGroupApplyHandleRejected:  {Category: CategoryGroup, Code: CodeGroupApplyHandleRejected, DefaultTitle: TitleGroupApply, DefaultStatus: StatusRejected},
	EventPrivacyPhoneApply:         {Category: CategorySystem, Code: CodePrivacyPhoneApply, DefaultTitle: TitlePrivacyApply, DefaultStatus: StatusPending},
	EventPrivacyEmailApply:         {Category: CategorySystem, Code: CodePrivacyEmailApply, DefaultTitle: TitlePrivacyApply, DefaultStatus: StatusPending},
	EventPrivacyHandleApproved:     {Category: CategorySystem, Code: CodePrivacyHandleApproved, DefaultTitle: TitlePrivacyApply, DefaultStatus: StatusApproved},
	EventPrivacyHandleRejected:     {Category: CategorySystem, Code: CodePrivacyHandleRejected, DefaultTitle: TitlePrivacyApply, DefaultStatus: StatusRejected},
	EventPrivacyApplyApproved:      {Category: CategorySystem, Code: CodePrivacyApplyApproved, DefaultTitle: TitlePrivacyApply, DefaultStatus: StatusApproved},
	EventPrivacyApplyRejected:      {Category: CategorySystem, Code: CodePrivacyApplyRejected, DefaultTitle: TitlePrivacyApplyReject, DefaultStatus: StatusRejected},
}

type Message struct {
	Title    string         `json:"title"`
	Content  string         `json:"content"`
	Code     int32          `json:"code"`
	Status   string         `json:"status"`
	Category string         `json:"category"`
	Ex       map[string]any `json:"ex"`
	Sender   string         `json:"sender"`
}

func SpecByEvent(event string) (EventSpec, bool) {
	spec, ok := eventSpecs[event]
	return spec, ok
}

func CodeByEvent(event string) (int32, bool) {
	spec, ok := SpecByEvent(event)
	if !ok {
		return 0, false
	}
	return spec.Code, true
}

func DefaultTitleByEvent(event string) (string, bool) {
	spec, ok := SpecByEvent(event)
	if !ok {
		return "", false
	}
	return spec.DefaultTitle, true
}

func DefaultStatusByEvent(event string) (string, bool) {
	spec, ok := SpecByEvent(event)
	if !ok {
		return "", false
	}
	return spec.DefaultStatus, true
}

func CategoryByEvent(event string) (string, bool) {
	spec, ok := SpecByEvent(event)
	if !ok {
		return "", false
	}
	return spec.Category, true
}

// NewMessage 用来构造最终下发给 App 的 JSON 协议体。
// title/code/status 都从共享协议表里派生，调用方只负责提供动态 content，
// 这样可以避免业务层把协议结构再写散。
func NewMessage(event, content string) (*Message, error) {
	spec, ok := SpecByEvent(event)
	if !ok {
		return nil, errs.ErrArgs.WrapMsg("unknown admin background event", "event", event)
	}
	return &Message{
		Title:    spec.DefaultTitle,
		Content:  content,
		Code:     spec.Code,
		Status:   spec.DefaultStatus,
		Category: spec.Category,
		Ex:       map[string]any{},
		Sender:   Sender,
	}, nil
}

// NewMessageWithTitle 只给少数需要临时覆盖标题的场景使用。
// 即使允许覆盖 title，code/status/sender 仍然必须锚定到统一协议定义上。
func NewMessageWithTitle(event, title, content string) (*Message, error) {
	spec, ok := SpecByEvent(event)
	if !ok {
		return nil, errs.ErrArgs.WrapMsg("unknown admin background event", "event", event)
	}
	if title == "" {
		title = spec.DefaultTitle
	}
	return &Message{
		Title:    title,
		Content:  content,
		Code:     spec.Code,
		Status:   spec.DefaultStatus,
		Category: spec.Category,
		Ex:       map[string]any{},
		Sender:   Sender,
	}, nil
}
