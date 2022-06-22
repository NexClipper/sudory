package v2

type HttpRsp_Template struct {
	Template `json:",inline"`
	Commands []TemplateCommand `json:"commands,omitempty"`
}

type HttpReq_TemplateCreate struct {
	Template_essential `json:",inline"`
	Commands           []TemplateCommand_essential `json:"commands,omitempty"`
}

type HttpReq_TemplateUpdate struct {
	Template_essential `json:",inline"`
	Commands           []TemplateCommand_essential `json:"commands,omitempty"`
}
