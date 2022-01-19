package control

// // CreateToken
// // @Description Create a Token
// // @Accept json
// // @Produce json
// // @Tags server
// // @Router /server/cluster/{id}/token [post]
// // @Param id path string true "cluster id"
// // @Param token body model.ReqToken true "Token의 정보"
// // @Success 200
// func (c *Control) CreateToken(ctx echo.Context) error {
// 	v := view.NewCreateToken(operator.NewToken(c.db))
// 	return v.Request(ctx)
// }
