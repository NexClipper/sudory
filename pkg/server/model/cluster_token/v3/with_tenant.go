package v3

import (
	"fmt"
	"reflect"
	"strings"

	"github.com/NexClipper/sudory/pkg/server/database/vanilla/ice_cream_maker"
	"github.com/NexClipper/sudory/pkg/server/macro/slicestrings"
	clusterv3 "github.com/NexClipper/sudory/pkg/server/model/cluster/v3"
	tenantv3 "github.com/NexClipper/sudory/pkg/server/model/tenant/v3"
)

var (
	TableNameWithTenant = tableNameWithTenant()
)

func tableNameWithTenant() func(tenant_hash string) string {
	var CT = ClusterToken{}
	var C = clusterv3.Cluster{}
	var TC = tenantv3.TenantClusters{}
	var T = tenantv3.Tenant{}

	aliasCT := CT.TableName()
	aliasC := C.TableName()
	aliasTC := TC.TableName()
	aliasT := T.TableName()

	columninfos := ice_cream_maker.ParseColumnTag(reflect.TypeOf(CT), ice_cream_maker.ParseColumnTag_opt{})

	columns := make([]string, 0, len(columninfos))
	for i := range columninfos {
		columns = append(columns, columninfos[i].Name)
	}

	columns = slicestrings.Map(columns, func(s string, i int) string {
		return aliasCT + "." + s
	})

	tables := []string{aliasCT, aliasC, aliasTC, aliasT}

	conds := []string{
		fmt.Sprintf("%v.hash = '%%v'", aliasT),
		fmt.Sprintf("%v.deleted IS NULL", aliasT),
		fmt.Sprintf("%v.tenant_id = %v.id", aliasTC, aliasT),
		fmt.Sprintf("%v.id = %v.cluster_id", aliasC, aliasTC),
		fmt.Sprintf("%v.deleted IS NULL", aliasC),
		fmt.Sprintf("%v.cluster_uuid = %v.uuid", aliasCT, aliasC),
	}

	format := fmt.Sprintf("( SELECT %v FROM %v WHERE %v ) X",
		strings.Join(columns, ", "),
		strings.Join(tables, ", "),
		strings.Join(conds, " AND "),
	)

	return func(tenant_hash string) string {
		return fmt.Sprintf(format, tenant_hash)
	}
}
