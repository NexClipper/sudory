package view

// type CreateClient struct {
// 	opr *operator.Client
// }

// func NewCreateClient(o operator.Operator) Viewer {
// 	return &CreateClient{opr: o.(*operator.Client)}
// }

// func (v *CreateClient) fromModel(m *model.ReqClient) {
// 	v.opr.MachineID = m.MachineID
// 	v.opr.ClusterID = m.ClusterID
// 	v.opr.IP = m.IP
// 	v.opr.Port = m.Port
// 	v.opr.Response = v.Response
// }

// func (v *CreateClient) Request(ctx echo.Context) error {
// 	reqModel := &model.ReqClient{}
// 	if err := ctx.Bind(reqModel); err != nil {
// 		return ctx.JSON(http.StatusBadRequest, nil)
// 	}

// 	v.fromModel(reqModel)
// 	if err := v.opr.Create(ctx); err != nil {
// 		return ctx.JSON(http.StatusInternalServerError, nil)
// 	}

// 	return nil
// }

// func (v *CreateClient) Response(ctx echo.Context, m model.Modeler) error {
// 	if err := ctx.JSON(http.StatusOK, nil); err != nil {
// 		return err
// 	}

// 	return nil
// }
