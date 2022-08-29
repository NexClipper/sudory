package v3

type SlackhookConfig struct {
	Url            string `column:"url,default('')"            json:"url"`
	RequestTimeout uint   `column:"request_timeout,default(0)" json:"request_timeout"` // second
}

func (SlackhookConfig) Type() NotifierType {
	return NotifierTypeSlackhook
}

func (cfg SlackhookConfig) Valid() bool {
	return true
}

type NotifierSlackhook_update = SlackhookConfig

type NotifierSlackhook_property = SlackhookConfig

type NotifierSlackhook struct {
	Uuid string `column:"uuid" json:"uuid,omitempty"` // pk

	NotifierSlackhook_property `json:",inline"`

	// Created vanilla.NullTime `column:"created" json:"created,omitempty" swaggertype:"string"`
	// Updated vanilla.NullTime `column:"updated" json:"updated,omitempty" swaggertype:"string"`
}

func (NotifierSlackhook) TableName() string {
	return "managed_channel_notifier_slackhook"
}
