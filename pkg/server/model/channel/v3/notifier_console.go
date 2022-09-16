package v3

type ConsoleConfig struct{}

func (ConsoleConfig) Type() NotifierType {
	return NotifierTypeConsole
}

func (ConsoleConfig) Valid() error {
	return nil
}

type NotifierConsole_update = ConsoleConfig

type NotifierConsole_property = ConsoleConfig

type NotifierConsole struct {
	Uuid string `column:"uuid"    json:"uuid,omitempty"` // pk

	NotifierConsole_property `json:",inline"`

	// Created vanilla.NullTime `column:"created" json:"created,omitempty" swaggertype:"string"`
	// Updated vanilla.NullTime `column:"updated" json:"updated,omitempty" swaggertype:"string"`
}

func (NotifierConsole) TableName() string {
	return "managed_channel_notifier_console"
}
