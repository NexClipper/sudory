package v2

type HttpRsp_ManagedChannel struct {
	ManagedChannel_tangled `json:",inline"`

	// NotifierEdge_property `json:"notifier_type,omitempty"`

	// Notifier struct {
	// 	*NotifierConsole_property  `json:"console,omitempty"`
	// 	*NotifierRabbitMq_property `json:"rabbitmq,omitempty"`
	// 	*NotifierWebhook_property  `json:"webhook,omitempty"`
	// } `json:"notifier,omitempty"`

	// Notifier struct {
	// 	NotifierEdge_option `json:",inline"`
	// } `json:"notifier,omitempty"`
}

type HttpReq_ManagedChannel_create struct {
	ManagedChannel_property `json:",inline"`

	Uuid string `json:"uuid,omitempty"` // optional
}

type HttpReq_ManagedChannel_update struct {
	ManagedChannel_property `json:",inline"`
}

type HttpRsp_ManagedChannel_NotifierEdge struct {
	NotifierEdge_option `json:",inline"`
}

type HttpReq_ManagedChannel_NotifierConsole_update struct {
	NotifierConsole_essential `json:",inline"`
}

type HttpReq_ManagedChannel_NotifierRabbitMq_update struct {
	NotifierRabbitMq_essential `json:",inline"`
}

type HttpReq_ManagedChannel_NotifierWebhook_update struct {
	NotifierWebhook_essential `json:",inline"`
}

type HttpReq_ManagedChannel_Format_update struct {
	Format_essential `json:",inline"`
}
type HttpRsq_ManagedChannel_Format struct {
	Format `json:",inline"`
}

type HttpReq_ManagedChannel_ChannelStatusOption_update struct {
	ChannelStatusOption_essential `json:",inline"`
}

type HttpRsp_ManagedChannel_ChannelStatusOption struct {
	ChannelStatusOption `json:",inline"`
}

type HttpRsp_ManagedChannel_ChannelStatus struct {
	ChannelStatus `json:",inline"`
}
