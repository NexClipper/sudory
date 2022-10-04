package v3

type HttpRsp_ManagedChannel = ManagedChannel_tangled

type HttpReq_ManagedChannel_create = ManagedChannel_create

type HttpReq_ManagedChannel_update = ManagedChannel_update

type HttpRsp_ManagedChannel_NotifierEdge = NotifierEdge_option

type HttpReq_ManagedChannel_NotifierConsole_update = NotifierConsole_update

type HttpReq_ManagedChannel_NotifierRabbitMq_update = NotifierRabbitMq_update

type HttpReq_ManagedChannel_NotifierWebhook_update = NotifierWebhook_update

type HttpReq_ManagedChannel_NotifierSlackhook_update = NotifierSlackhook_update

type HttpReq_ManagedChannel_Format_update = Format_update

type HttpRsq_ManagedChannel_Format = Format

type HttpReq_ManagedChannel_ChannelStatusOption_update = ChannelStatusOption_update

type HttpRsp_ManagedChannel_ChannelStatusOption = ChannelStatusOption

type HttpRsp_ManagedChannel_ChannelStatus = ChannelStatus

type ManagedChannel_tangled struct {
	ManagedChannel `json:",inline"`

	StatusOption ChannelStatusOption_property `json:"status_option,omitempty"`
	Format       Format_property              `json:"format,omitempty"`

	Notifiers struct {
		NotifierEdge_property `json:",inline"`
		Console               *NotifierConsole_property   `json:"console,omitempty"`
		Webhook               *NotifierWebhook_property   `json:"webhook,omitempty"`
		RabbitMq              *NotifierRabbitMq_property  `json:"rabbitmq,omitempty"`
		Slackhook             *NotifierSlackhook_property `json:"slackhook,omitempty"`
	} `json:"notifiers,omitempty"`
}

type NotifierEdge_option struct {
	NotifierEdge `json:",inline"`

	Console   *NotifierConsole_property   `json:"console,omitempty"`
	Webhook   *NotifierWebhook_property   `json:"webhook,omitempty"`
	RabbitMq  *NotifierRabbitMq_property  `json:"rabbitmq,omitempty"`
	Slackhook *NotifierSlackhook_property `json:"slackhook,omitempty"`
}
