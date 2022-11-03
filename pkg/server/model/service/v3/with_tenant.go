package service

import (
	"fmt"
	"reflect"
	"strings"

	"github.com/NexClipper/sudory/pkg/server/database/vanilla/ice_cream_maker"
	"github.com/NexClipper/sudory/pkg/server/macro/slicestrings"
	clusterv3 "github.com/NexClipper/sudory/pkg/server/model/cluster/v3"
	"github.com/NexClipper/sudory/pkg/server/model/tenants/v3"
)

var (
	TableNameWithTenant_Service       = tableNameWithTenant_Service()
	TableNameWithTenant_ServiceStep   = tableNameWithTenant_ServiceStep()
	TableNameWithTenant_ServiceResult = tableNameWithTenant_ServiceResult()
)

func tableNameWithTenant_Service() func(tenant_hash string) string {
	var S = Service{}
	var C = clusterv3.Cluster{}
	var TC = tenants.TenantClusters{}
	var T = tenants.Tenant{}

	aliasS := S.TableName()
	aliasC := C.TableName()
	aliasTC := TC.TableName()
	aliasT := T.TableName()

	columninfos := ice_cream_maker.ParseColumnTag(reflect.TypeOf(S), ice_cream_maker.ParseColumnTag_opt{})

	columns := make([]string, 0, len(columninfos))
	for i := range columninfos {
		columns = append(columns, columninfos[i].Name)
	}

	columns = slicestrings.Map(columns, func(s string, i int) string {
		return aliasS + "." + s
	})

	tables := []string{aliasS, aliasC, aliasTC, aliasT}

	conds := []string{
		fmt.Sprintf("%v.hash = '%%v'", aliasT),
		fmt.Sprintf("%v.deleted IS NULL", aliasT),
		fmt.Sprintf("%v.tenant_id = %v.id", aliasTC, aliasT),
		fmt.Sprintf("%v.id = %v.cluster_id", aliasC, aliasTC),
		fmt.Sprintf("%v.deleted IS NULL", aliasC),
		fmt.Sprintf("%v.cluster_uuid = %v.uuid", aliasS, aliasC),
	}

	format := fmt.Sprintf(`( SELECT %v FROM %v 
		INNER JOIN ( SELECT service.cluster_uuid,service.uuid,service.pdate,MAX(service.timestamp) AS timestamp FROM %v 
		WHERE %v 
		GROUP BY cluster_uuid,uuid,pdate ) stat
		ON stat.cluster_uuid = service.cluster_uuid AND stat.uuid = service.uuid AND stat.pdate = service.pdate AND stat.timestamp = service.timestamp ) X`,
		strings.Join(columns, ", "),
		aliasS,
		strings.Join(tables, ", "),
		strings.Join(conds, " AND "),
	)

	return func(tenant_hash string) string {
		return fmt.Sprintf(format, tenant_hash)
	}
}

func tableNameWithTenant_ServiceStep() func(tenant_hash string) string {
	var S = ServiceStep{}
	var C = clusterv3.Cluster{}
	var TC = tenants.TenantClusters{}
	var T = tenants.Tenant{}

	aliasS := S.TableName()
	aliasC := C.TableName()
	aliasTC := TC.TableName()
	aliasT := T.TableName()

	columninfos := ice_cream_maker.ParseColumnTag(reflect.TypeOf(S), ice_cream_maker.ParseColumnTag_opt{})

	columns := make([]string, 0, len(columninfos))
	for i := range columninfos {
		columns = append(columns, columninfos[i].Name)
	}

	columns = slicestrings.Map(columns, func(s string, i int) string {
		return aliasS + "." + s
	})

	tables := []string{aliasS, aliasC, aliasTC, aliasT}

	conds := []string{
		fmt.Sprintf("%v.hash = '%%v'", aliasT),
		fmt.Sprintf("%v.deleted IS NULL", aliasT),
		fmt.Sprintf("%v.tenant_id = %v.id", aliasTC, aliasT),
		fmt.Sprintf("%v.id = %v.cluster_id", aliasC, aliasTC),
		fmt.Sprintf("%v.deleted IS NULL", aliasC),
		fmt.Sprintf("%v.cluster_uuid = %v.uuid", aliasS, aliasC),
	}

	// format := fmt.Sprintf("( SELECT %v FROM %v WHERE %v ) X",
	// 	strings.Join(columns, ", "),
	// 	strings.Join(tables, ", "),
	// 	strings.Join(conds, " AND "),
	// )

	format := fmt.Sprintf(`( SELECT %v FROM %v 
		INNER JOIN ( SELECT service_step.cluster_uuid,service_step.uuid,service_step.seq,service_step.pdate,MAX(service_step.timestamp) AS timestamp FROM %v 
		WHERE %v 
		GROUP BY cluster_uuid,uuid,seq,pdate ) stat
		ON stat.cluster_uuid = service_step.cluster_uuid AND stat.uuid = service_step.uuid AND stat.pdate = service_step.pdate AND stat.timestamp = service_step.timestamp ) X`,
		strings.Join(columns, ", "),
		aliasS,
		strings.Join(tables, ", "),
		strings.Join(conds, " AND "),
	)

	return func(tenant_hash string) string {
		return fmt.Sprintf(format, tenant_hash)
	}
}

func tableNameWithTenant_ServiceResult() func(tenant_hash string) string {
	var R = ServiceResult{}
	var C = clusterv3.Cluster{}
	var TC = tenants.TenantClusters{}
	var T = tenants.Tenant{}

	aliasR := R.TableName()
	aliasC := C.TableName()
	aliasTC := TC.TableName()
	aliasT := T.TableName()

	columninfos := ice_cream_maker.ParseColumnTag(reflect.TypeOf(R), ice_cream_maker.ParseColumnTag_opt{})

	columns := make([]string, 0, len(columninfos))
	for i := range columninfos {
		columns = append(columns, columninfos[i].Name)
	}

	columns = slicestrings.Map(columns, func(s string, i int) string {
		return aliasR + "." + s
	})

	tables := []string{aliasR, aliasC, aliasTC, aliasT}

	conds := []string{
		fmt.Sprintf("%v.hash = '%%v'", aliasT),
		fmt.Sprintf("%v.deleted IS NULL", aliasT),
		fmt.Sprintf("%v.tenant_id = %v.id", aliasTC, aliasT),
		fmt.Sprintf("%v.id = %v.cluster_id", aliasC, aliasTC),
		fmt.Sprintf("%v.deleted IS NULL", aliasC),
		fmt.Sprintf("%v.cluster_uuid = %v.uuid", aliasR, aliasC),
	}

	// format := fmt.Sprintf("( SELECT %v FROM %v WHERE %v ) X",
	// 	strings.Join(columns, ", "),
	// 	strings.Join(tables, ", "),
	// 	strings.Join(conds, " AND "),
	// )

	format := fmt.Sprintf(`( SELECT %v FROM %v 
		INNER JOIN ( SELECT service_result.cluster_uuid,service_result.uuid,service_result.pdate,MAX(service_result.timestamp) AS timestamp FROM %v 
		WHERE %v 
		GROUP BY cluster_uuid,uuid,pdate ) stat 
		ON stat.cluster_uuid = service_result.cluster_uuid AND stat.uuid = service_result.uuid AND stat.pdate = service_result.pdate AND stat.timestamp = service_result.timestamp ) X`,
		strings.Join(columns, ", "),
		aliasR,
		strings.Join(tables, ", "),
		strings.Join(conds, " AND "),
	)

	return func(tenant_hash string) string {
		return fmt.Sprintf(format, tenant_hash)
	}
}
