package v3

import (
	"fmt"
	"reflect"
	"strings"

	"github.com/NexClipper/sudory/pkg/server/database/vanilla/ice_cream_maker"
	"github.com/NexClipper/sudory/pkg/server/macro/slicestrings"
	tenantv3 "github.com/NexClipper/sudory/pkg/server/model/tenant/v3"
)

var (
	TableNameWithTenant_ManagedChannel = tableNameWithTenant_ManagedChannel()
	// TableNameWithTenant_ManagedChannel_scoped = tableNameWithTenant_ManagedChannel_scoped()
	TableNameWithTenant_ChannelStatusOption = tableNameWithTenant_ChannelStatusOption()
	TableNameWithTenant_Format              = tableNameWithTenant_Format()
	TableNameWithTenant_NotifierEdge        = tableNameWithTenant_NotifierEdge()
	TableNameWithTenant_NotifierConsole     = tableNameWithTenant_NotifierConsole()
	TableNameWithTenant_NotifierWebhook     = tableNameWithTenant_NotifierWebhook()
	TableNameWithTenant_NotifierRabbitMq    = tableNameWithTenant_NotifierRabbitMq()
	TableNameWithTenant_NotifierSlackhook   = tableNameWithTenant_NotifierSlackhook()
	TableNameWithTenant_ChannelStatus       = tableNameWithTenant_ChannelStatus()
)

func tableNameWithTenant_ManagedChannel() func(tenant_hash string) string {
	var C = ManagedChannel{}
	var TC = tenantv3.TenantChannels{}
	var T = tenantv3.Tenant{}

	aliasC := C.TableName()
	aliasTC := TC.TableName()
	aliasT := T.TableName()

	columninfos := ice_cream_maker.ParseColumnTag(reflect.TypeOf(C), ice_cream_maker.ParseColumnTag_opt{})

	columns := make([]string, 0, len(columninfos))
	for i := range columninfos {
		columns = append(columns, columninfos[i].Name)
	}

	columns = slicestrings.Map(columns, func(s string, i int) string {
		return aliasC + "." + s
	})

	tables := []string{aliasC, aliasTC, aliasT}

	conds := []string{
		fmt.Sprintf("%v.hash = '%%v'", aliasT),
		fmt.Sprintf("%v.deleted IS NULL", aliasT),
		fmt.Sprintf("%v.tenant_id = %v.id", aliasTC, aliasT),
		fmt.Sprintf("%v.uuid = %v.channel_uuid", aliasC, aliasTC),
		// fmt.Sprintf("%v.deleted IS NULL", aliasC),
	}

	format := fmt.Sprintf("( SELECT %v FROM %v WHERE %v ) X",
		strings.Join(columns, ", "),
		// aliasC,
		strings.Join(tables, ", "),
		strings.Join(conds, " AND "),
	)

	return func(tenant_hash string) string {
		return fmt.Sprintf(format, tenant_hash)
	}
}

// func tableNameWithTenant_ManagedChannel_scoped() func(tenant_hash string) string {
// 	var C = ManagedChannel{}
// 	var TC = tenantv3.TenantChannels{}
// 	var T = tenantv3.Tenant{}

// 	aliasC := C.TableName()
// 	aliasTC := TC.TableName()
// 	aliasT := T.TableName()

// 	columninfos := ice_cream_maker.ParseColumnTag(reflect.TypeOf(C), ice_cream_maker.ParseColumnTag_opt{})

// 	columns := make([]string, 0, len(columninfos))
// 	for i := range columninfos {
// 		columns = append(columns, columninfos[i].Name)
// 	}

// 	columns = slicestrings.Map(columns, func(s string, i int) string {
// 		return aliasC + "." + s
// 	})

// 	tables := []string{aliasC, aliasTC, aliasT}

// 	conds := []string{
// 		fmt.Sprintf("%v.hash = '%%v'", aliasT),
// 		fmt.Sprintf("%v.deleted IS NULL", aliasT),
// 		fmt.Sprintf("%v.tenant_id = %v.id", aliasTC, aliasT),
// 		fmt.Sprintf("%v.uuid = %v.channel_uuid", aliasC, aliasTC),
// 		// fmt.Sprintf("%v.deleted IS NULL", aliasC),
// 	}

// 	format := fmt.Sprintf("( SELECT %v FROM %v WHERE %v ) X",
// 		strings.Join(columns, ", "),
// 		// aliasC,
// 		strings.Join(tables, ", "),
// 		strings.Join(conds, " AND "),
// 	)

// 	return func(tenant_hash string) string {
// 		return fmt.Sprintf(format, tenant_hash)
// 	}
// }

func tableNameWithTenant_ChannelStatusOption() func(tenant_hash string) string {
	var C = ManagedChannel{}
	var SO = ChannelStatusOption{}
	var TC = tenantv3.TenantChannels{}
	var T = tenantv3.Tenant{}

	aliasC := C.TableName()
	aliasSO := SO.TableName()
	aliasTC := TC.TableName()
	aliasT := T.TableName()

	columninfos := ice_cream_maker.ParseColumnTag(reflect.TypeOf(SO), ice_cream_maker.ParseColumnTag_opt{})

	columns := make([]string, 0, len(columninfos))
	for i := range columninfos {
		columns = append(columns, columninfos[i].Name)
	}

	columns = slicestrings.Map(columns, func(s string, i int) string {
		return aliasSO + "." + s
	})

	tables := []string{aliasSO, aliasC, aliasTC, aliasT}

	conds := []string{
		fmt.Sprintf("%v.hash = '%%v'", aliasT),
		fmt.Sprintf("%v.deleted IS NULL", aliasT),
		fmt.Sprintf("%v.tenant_id = %v.id", aliasTC, aliasT),
		fmt.Sprintf("%v.uuid = %v.channel_uuid", aliasC, aliasTC),
		fmt.Sprintf("%v.deleted IS NULL", aliasC),
		fmt.Sprintf("%v.uuid = %v.uuid", aliasSO, aliasC),
	}

	format := fmt.Sprintf("( SELECT %v FROM %v WHERE %v ) X",
		strings.Join(columns, ", "),
		// aliasC,
		strings.Join(tables, ", "),
		strings.Join(conds, " AND "),
	)

	return func(tenant_hash string) string {
		return fmt.Sprintf(format, tenant_hash)
	}
}

func tableNameWithTenant_Format() func(tenant_hash string) string {
	var C = ManagedChannel{}
	var F = Format{}
	var TC = tenantv3.TenantChannels{}
	var T = tenantv3.Tenant{}

	aliasC := C.TableName()
	aliasF := F.TableName()
	aliasTC := TC.TableName()
	aliasT := T.TableName()

	columninfos := ice_cream_maker.ParseColumnTag(reflect.TypeOf(F), ice_cream_maker.ParseColumnTag_opt{})

	columns := make([]string, 0, len(columninfos))
	for i := range columninfos {
		columns = append(columns, columninfos[i].Name)
	}

	columns = slicestrings.Map(columns, func(s string, i int) string {
		return aliasF + "." + s
	})

	tables := []string{aliasF, aliasC, aliasTC, aliasT}

	conds := []string{
		fmt.Sprintf("%v.hash = '%%v'", aliasT),
		fmt.Sprintf("%v.deleted IS NULL", aliasT),
		fmt.Sprintf("%v.tenant_id = %v.id", aliasTC, aliasT),
		fmt.Sprintf("%v.uuid = %v.channel_uuid", aliasC, aliasTC),
		fmt.Sprintf("%v.deleted IS NULL", aliasC),
		fmt.Sprintf("%v.uuid = %v.uuid", aliasF, aliasC),
	}

	format := fmt.Sprintf("( SELECT %v FROM %v WHERE %v ) X",
		strings.Join(columns, ", "),
		// aliasC,
		strings.Join(tables, ", "),
		strings.Join(conds, " AND "),
	)

	return func(tenant_hash string) string {
		return fmt.Sprintf(format, tenant_hash)
	}
}

func tableNameWithTenant_NotifierEdge() func(tenant_hash string) string {
	var C = ManagedChannel{}
	var NE = NotifierEdge{}
	var TC = tenantv3.TenantChannels{}
	var T = tenantv3.Tenant{}

	aliasC := C.TableName()
	aliasNE := NE.TableName()
	aliasTC := TC.TableName()
	aliasT := T.TableName()

	columninfos := ice_cream_maker.ParseColumnTag(reflect.TypeOf(NE), ice_cream_maker.ParseColumnTag_opt{})

	columns := make([]string, 0, len(columninfos))
	for i := range columninfos {
		columns = append(columns, columninfos[i].Name)
	}

	columns = slicestrings.Map(columns, func(s string, i int) string {
		return aliasNE + "." + s
	})

	tables := []string{aliasNE, aliasC, aliasTC, aliasT}

	conds := []string{
		fmt.Sprintf("%v.hash = '%%v'", aliasT),
		fmt.Sprintf("%v.deleted IS NULL", aliasT),
		fmt.Sprintf("%v.tenant_id = %v.id", aliasTC, aliasT),
		fmt.Sprintf("%v.uuid = %v.channel_uuid", aliasC, aliasTC),
		fmt.Sprintf("%v.deleted IS NULL", aliasC),
		fmt.Sprintf("%v.uuid = %v.uuid", aliasNE, aliasC),
	}

	format := fmt.Sprintf("( SELECT %v FROM %v WHERE %v ) X",
		strings.Join(columns, ", "),
		// aliasC,
		strings.Join(tables, ", "),
		strings.Join(conds, " AND "),
	)

	return func(tenant_hash string) string {
		return fmt.Sprintf(format, tenant_hash)
	}
}

func tableNameWithTenant_NotifierConsole() func(tenant_hash string) string {
	var C = ManagedChannel{}
	var NC = NotifierConsole{}
	var TC = tenantv3.TenantChannels{}
	var T = tenantv3.Tenant{}

	aliasC := C.TableName()
	aliasNC := NC.TableName()
	aliasTC := TC.TableName()
	aliasT := T.TableName()

	columninfos := ice_cream_maker.ParseColumnTag(reflect.TypeOf(NC), ice_cream_maker.ParseColumnTag_opt{})

	columns := make([]string, 0, len(columninfos))
	for i := range columninfos {
		columns = append(columns, columninfos[i].Name)
	}

	columns = slicestrings.Map(columns, func(s string, i int) string {
		return aliasNC + "." + s
	})

	tables := []string{aliasNC, aliasC, aliasTC, aliasT}

	conds := []string{
		fmt.Sprintf("%v.hash = '%%v'", aliasT),
		fmt.Sprintf("%v.deleted IS NULL", aliasT),
		fmt.Sprintf("%v.tenant_id = %v.id", aliasTC, aliasT),
		fmt.Sprintf("%v.uuid = %v.channel_uuid", aliasC, aliasTC),
		fmt.Sprintf("%v.deleted IS NULL", aliasC),
		fmt.Sprintf("%v.uuid = %v.uuid", aliasNC, aliasC),
	}

	format := fmt.Sprintf("( SELECT %v FROM %v WHERE %v ) X",
		strings.Join(columns, ", "),
		// aliasC,
		strings.Join(tables, ", "),
		strings.Join(conds, " AND "),
	)

	return func(tenant_hash string) string {
		return fmt.Sprintf(format, tenant_hash)
	}
}

func tableNameWithTenant_NotifierWebhook() func(tenant_hash string) string {
	var C = ManagedChannel{}
	var NW = NotifierWebhook{}
	var TC = tenantv3.TenantChannels{}
	var T = tenantv3.Tenant{}

	aliasC := C.TableName()
	aliasNW := NW.TableName()
	aliasTC := TC.TableName()
	aliasT := T.TableName()

	columninfos := ice_cream_maker.ParseColumnTag(reflect.TypeOf(NW), ice_cream_maker.ParseColumnTag_opt{})

	columns := make([]string, 0, len(columninfos))
	for i := range columninfos {
		columns = append(columns, columninfos[i].Name)
	}

	columns = slicestrings.Map(columns, func(s string, i int) string {
		return aliasNW + "." + s
	})

	tables := []string{aliasNW, aliasC, aliasTC, aliasT}

	conds := []string{
		fmt.Sprintf("%v.hash = '%%v'", aliasT),
		fmt.Sprintf("%v.deleted IS NULL", aliasT),
		fmt.Sprintf("%v.tenant_id = %v.id", aliasTC, aliasT),
		fmt.Sprintf("%v.uuid = %v.channel_uuid", aliasC, aliasTC),
		fmt.Sprintf("%v.deleted IS NULL", aliasC),
		fmt.Sprintf("%v.uuid = %v.uuid", aliasNW, aliasC),
	}

	format := fmt.Sprintf("( SELECT %v FROM %v WHERE %v ) X",
		strings.Join(columns, ", "),
		// aliasC,
		strings.Join(tables, ", "),
		strings.Join(conds, " AND "),
	)

	return func(tenant_hash string) string {
		return fmt.Sprintf(format, tenant_hash)
	}
}

func tableNameWithTenant_NotifierRabbitMq() func(tenant_hash string) string {
	var C = ManagedChannel{}
	var NR = NotifierRabbitMq{}
	var TC = tenantv3.TenantChannels{}
	var T = tenantv3.Tenant{}

	aliasC := C.TableName()
	aliasNR := NR.TableName()
	aliasTC := TC.TableName()
	aliasT := T.TableName()

	columninfos := ice_cream_maker.ParseColumnTag(reflect.TypeOf(NR), ice_cream_maker.ParseColumnTag_opt{})

	columns := make([]string, 0, len(columninfos))
	for i := range columninfos {
		columns = append(columns, columninfos[i].Name)
	}

	columns = slicestrings.Map(columns, func(s string, i int) string {
		return aliasNR + "." + s
	})

	tables := []string{aliasNR, aliasC, aliasTC, aliasT}

	conds := []string{
		fmt.Sprintf("%v.hash = '%%v'", aliasT),
		fmt.Sprintf("%v.deleted IS NULL", aliasT),
		fmt.Sprintf("%v.tenant_id = %v.id", aliasTC, aliasT),
		fmt.Sprintf("%v.uuid = %v.channel_uuid", aliasC, aliasTC),
		fmt.Sprintf("%v.deleted IS NULL", aliasC),
		fmt.Sprintf("%v.uuid = %v.uuid", aliasNR, aliasC),
	}

	format := fmt.Sprintf("( SELECT %v FROM %v WHERE %v ) X",
		strings.Join(columns, ", "),
		// aliasC,
		strings.Join(tables, ", "),
		strings.Join(conds, " AND "),
	)

	return func(tenant_hash string) string {
		return fmt.Sprintf(format, tenant_hash)
	}
}

func tableNameWithTenant_NotifierSlackhook() func(tenant_hash string) string {
	var C = ManagedChannel{}
	var NS = NotifierSlackhook{}
	var TC = tenantv3.TenantChannels{}
	var T = tenantv3.Tenant{}

	aliasC := C.TableName()
	aliasNS := NS.TableName()
	aliasTC := TC.TableName()
	aliasT := T.TableName()

	columninfos := ice_cream_maker.ParseColumnTag(reflect.TypeOf(NS), ice_cream_maker.ParseColumnTag_opt{})

	columns := make([]string, 0, len(columninfos))
	for i := range columninfos {
		columns = append(columns, columninfos[i].Name)
	}

	columns = slicestrings.Map(columns, func(s string, i int) string {
		return aliasNS + "." + s
	})

	tables := []string{aliasNS, aliasC, aliasTC, aliasT}

	conds := []string{
		fmt.Sprintf("%v.hash = '%%v'", aliasT),
		fmt.Sprintf("%v.deleted IS NULL", aliasT),
		fmt.Sprintf("%v.tenant_id = %v.id", aliasTC, aliasT),
		fmt.Sprintf("%v.uuid = %v.channel_uuid", aliasC, aliasTC),
		fmt.Sprintf("%v.deleted IS NULL", aliasC),
		fmt.Sprintf("%v.uuid = %v.uuid", NS, aliasC),
	}

	format := fmt.Sprintf("( SELECT %v FROM %v WHERE %v ) X",
		strings.Join(columns, ", "),
		// aliasC,
		strings.Join(tables, ", "),
		strings.Join(conds, " AND "),
	)

	return func(tenant_hash string) string {
		return fmt.Sprintf(format, tenant_hash)
	}
}

func tableNameWithTenant_ChannelStatus() func(tenant_hash string) string {
	var C = ManagedChannel{}
	var CS = ChannelStatus{}
	var TC = tenantv3.TenantChannels{}
	var T = tenantv3.Tenant{}

	aliasC := C.TableName()
	aliasCS := CS.TableName()
	aliasTC := TC.TableName()
	aliasT := T.TableName()

	columninfos := ice_cream_maker.ParseColumnTag(reflect.TypeOf(CS), ice_cream_maker.ParseColumnTag_opt{})

	columns := make([]string, 0, len(columninfos))
	for i := range columninfos {
		columns = append(columns, columninfos[i].Name)
	}

	columns = slicestrings.Map(columns, func(s string, i int) string {
		return aliasCS + "." + s
	})

	tables := []string{aliasCS, aliasC, aliasTC, aliasT}

	conds := []string{
		fmt.Sprintf("%v.hash = '%%v'", aliasT),
		fmt.Sprintf("%v.deleted IS NULL", aliasT),
		fmt.Sprintf("%v.tenant_id = %v.id", aliasTC, aliasT),
		fmt.Sprintf("%v.uuid = %v.channel_uuid", aliasC, aliasTC),
		fmt.Sprintf("%v.deleted IS NULL", aliasC),
		fmt.Sprintf("%v.uuid = %v.uuid", aliasCS, aliasC),
	}

	format := fmt.Sprintf("( SELECT %v FROM %v WHERE %v ) X",
		strings.Join(columns, ", "),
		// aliasC,
		strings.Join(tables, ", "),
		strings.Join(conds, " AND "),
	)

	return func(tenant_hash string) string {
		return fmt.Sprintf(format, tenant_hash)
	}
}
