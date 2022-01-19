package control

// // CreateCluster
// // @Description Create a Cluster
// // @Accept json
// // @Produce json
// // @Tags server
// // @Router /server/cluster [post]
// // @Param cluster body model.ReqCluster true "Cluster의 정보"
// // @Success 200
// func (c *Control) CreateCluster(ctx echo.Context) error {
// 	v := view.NewCreateCluster(operator.NewCluster(c.db))
// 	return v.Request(ctx)
// }

// // GetCluster
// // @Description Get a Cluster
// // @Accept json
// // @Produce json
// // @Tags server
// // @Router /server/cluster/{id} [get]
// // @Param id path string true "Cluster의 ID"
// // @Success 200 {object} model.Cluster
// func (c *Control) GetCluster(ctx echo.Context) error {
// 	v := view.NewGetCluster(operator.NewCluster(c.db))
// 	return v.Request(ctx)
// }
