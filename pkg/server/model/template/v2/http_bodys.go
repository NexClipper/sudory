package v2

type HttpRsp_Template struct {
	Template `json:",inline"`
	Commands []TemplateCommand `json:"commands,omitempty"`
}

type HttpRsp_TemplateCommand TemplateCommand
