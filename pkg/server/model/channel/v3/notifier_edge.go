package v3

type NotifierEdge_property struct {
	// enums:"NaV(0), console(1), webhook(2), rabbitmq(3), slackhook(4)"
	NotifierType NotifierType `column:"notifier_type,default(0)" json:"notifier_type,omitempty" swaggertype:"string" enums:"0,1,2,3,4"`
}

type NotifierEdge struct {
	Uuid string `column:"uuid" json:"uuid,omitempty"` // pk

	NotifierEdge_property `json:",inline"`

	// Created vanilla.NullTime `column:"created"                  json:"created,omitempty"       swaggertype:"string"`
	// Updated vanilla.NullTime `column:"updated"                  json:"updated,omitempty"       swaggertype:"string"`
}

func (NotifierEdge) TableName() string {
	return "managed_channel_notifier_edge"
}

// func (record NotifierEdge_option) TableName() string {
// 	q := `(
// 		SELECT %v /**columns**/
// 	      FROM %v EG /**managed_channel_notifier_edge**/
// 		  LEFT JOIN %v CS /**managed_channel_notifier_console**/
// 			     ON EG.uuid = CS.uuid
// 		  LEFT JOIN %v WH /**managed_channel_notifier_webhook**/
// 			     ON EG.uuid = WH.uuid
// 		  LEFT JOIN %v RQ /**managed_channel_notifier_rabbitmq**/
// 				 ON EG.uuid = RQ.uuid
// 		  LEFT JOIN %v SH /**managed_channel_notifier_slackhook**/
// 				 ON EG.uuid = SH.uuid
// 	) X
// `

// 	EG := record.NotifierEdge.TableName()
// 	CS := record.Console.TableName()
// 	WH := record.Webhook.TableName()
// 	RQ := record.RabbitMq.TableName()
// 	SH := record.Slackhook.TableName()
// 	return fmt.Sprintf(q, strings.Join(record.ColumnNamesWithAlias(), ", "), EG, CS, WH, RQ, SH)
// }
