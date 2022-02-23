package jq_test

import (
	"bytes"
	"context"
	"encoding/json"
	"testing"

	"github.com/NexClipper/sudory/pkg/client/jq"
	"github.com/NexClipper/sudory/pkg/client/k8s"
	"github.com/NexClipper/sudory/pkg/client/k8s/dynamic"
	"github.com/itchyny/gojq"
	v1 "k8s.io/api/apps/v1"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/apimachinery/pkg/runtime/schema"
	"k8s.io/client-go/util/jsonpath"
)

var testrawjson = `{"metadata":{"name":"prom","uid":"27ad4ebf-360e-4a2f-9540-7049b2755886","resourceVersion":"486439","creationTimestamp":"2022-02-14T04:40:26Z","labels":{"kubernetes.io/metadata.name":"prom"},"managedFields":[{"manager":"kubectl-create","operation":"Update","apiVersion":"v1","time":"2022-02-14T04:40:26Z","fieldsType":"FieldsV1","fieldsV1":{"f:metadata":{"f:labels":{".":{},"f:kubernetes.io/metadata.name":{}}}}}]},"spec":{"finalizers":["kubernetes"]},"status":{"phase":"Active"}}`

func TestGoJq0(t *testing.T) {
	// query, err := gojq.Parse(".metadata.name")
	// // query, err := gojq.Parse(".spec.replicas")
	// if err != nil {
	// 	t.Fatal(err)
	// }

	input := make(map[string]interface{})
	if err := json.Unmarshal([]byte(testrawjson), &input); err != nil {
		t.Fatal(err)
	}

	// iter := query.Run(input) // or query.RunWithContext
	// for {
	// 	v, ok := iter.Next()
	// 	if !ok {
	// 		break
	// 	}
	// 	if err, ok := v.(error); ok {
	// 		t.Fatal(err)
	// 	}
	// 	// t.Logf("%#v\n", v)
	// 	b, err := json.MarshalIndent(v, "", "  ")
	// 	if err != nil {
	// 		t.Fatal(err)
	// 	}
	// 	t.Logf("\n%v", string(b))
	// }

	res, err := jq.ProcessJson(input, ".metadata.uu")
	if err != nil {
		t.Fatal(err)
	}

	t.Log(res)
}

func TestGoJq(t *testing.T) {
	client, err := k8s.GetClient()
	if err != nil {
		t.Fatal(err)
	}

	outputJson, err := client.ResourceGet(schema.GroupVersion{Group: "apps", Version: "v1"}, "deployments", "sudory", "sudory-client")
	if err != nil {
		t.Fatal(err)
	}

	query, err := gojq.Parse(".spec.strategy")
	// query, err := gojq.Parse(".spec.replicas")
	if err != nil {
		t.Fatal(err)
	}

	input := make(map[string]interface{})
	if err := json.Unmarshal([]byte(outputJson), &input); err != nil {
		t.Fatal(err)
	}

	iter := query.Run(input) // or query.RunWithContext
	for {
		v, ok := iter.Next()
		if !ok {
			break
		}
		if err, ok := v.(error); ok {
			t.Fatal(err)
		}
		// t.Logf("%#v\n", v)
		b, err := json.MarshalIndent(v, "", "  ")
		if err != nil {
			t.Fatal(err)
		}
		t.Logf("\n%v", string(b))
	}
}

func TestJsonpath(t *testing.T) {
	client, err := k8s.GetClient()
	if err != nil {
		t.Fatal(err)
	}

	outputJson, err := client.ResourceGet(schema.GroupVersion{Group: "apps", Version: "v1"}, "deployments", "sudory", "sudory-client")
	if err != nil {
		t.Fatal(err)
	}
	found := new(v1.Deployment)
	if err := json.Unmarshal([]byte(outputJson), found); err != nil {
		t.Fatal(err)
	}

	// fields, err := get.RelaxedJSONPathExpression(".spec.strategy")
	// if err != nil {
	// 	t.Fatal(err)
	// }
	// t.Logf("RelaxedJSONPathExpression : %s\n", fields)

	parser := jsonpath.New("deploy")
	if err := parser.Parse("{.spec.strategy}"); err != nil {
		t.Fatal(err)
	}

	values, err := parser.FindResults(found)
	if err != nil {
		t.Fatal(err)
	}

	for i, vv := range values {

		buf := bytes.NewBuffer([]byte{})
		if err := parser.PrintResults(buf, vv); err != nil {
			t.Error(i, err)
		}
		t.Log("buffer string :", buf.String())
		for j, vvv := range vv {
			t.Logf("%d, %d : %v\n", i, j, vvv)
			ii := vvv.Interface().(v1.DeploymentStrategy)
			t.Log(ii.String())
			b, err := json.MarshalIndent(ii, "", "  ")
			if err != nil {
				t.Fatal(err)
			}
			t.Logf("\n%v", string(b))
		}
	}
}

func TestGoJq2(t *testing.T) {
	dc, err := dynamic.NewClient("")
	if err != nil {
		t.Fatal(err)
	}
	item, err := dc.Client.Resource(schema.GroupVersionResource{Group: "monitoring.coreos.com", Version: "v1", Resource: "prometheusrules"}).
		Namespace("prom").
		Get(context.TODO(), "prometheus-kube-prometheus-prometheus", metav1.GetOptions{})
	if err != nil {
		t.Fatal(err)
	}

	query, err := gojq.Parse(src_ruledeletion)
	if err != nil {
		t.Fatal(err)
	}

	iter := query.RunWithContext(context.TODO(), item.UnstructuredContent()) // or query.RunWithContext
	for {
		v, ok := iter.Next()
		if !ok {
			break
		}
		if err, ok := v.(error); ok {
			t.Fatal(err)
		}
		// t.Logf("%#v\n", v)
		b, err := json.MarshalIndent(v, "", "  ")
		if err != nil {
			t.Fatal(err)
		}
		t.Logf("result: \n%v", string(b))
	}
}
func TestGoJq3(t *testing.T) {
	dc, err := dynamic.NewClient("")
	if err != nil {
		t.Fatal(err)
	}
	items, err := dc.Client.Resource(schema.GroupVersionResource{Group: "monitoring.coreos.com", Version: "v1", Resource: "prometheuses"}).
		Namespace("prom").
		// Get(context.TODO(), "prometheus-kube-prometheus-prometheus", metav1.GetOptions{})
		List(context.TODO(), metav1.ListOptions{})
	if err != nil {
		t.Fatal(err)
	}

	query1, err := gojq.Parse(src_retentionperiod)
	if err != nil {
		t.Fatal(err)
	}
	query2, err := gojq.Parse(src_scrapinterval)
	if err != nil {
		t.Fatal(err)
	}
	query3, err := gojq.Parse(src_timeout)
	if err != nil {
		t.Fatal(err)
	}

	iter1 := query1.RunWithContext(context.TODO(), items.UnstructuredContent())
	for {
		v, ok := iter1.Next()
		if !ok {
			break
		}
		if err, ok := v.(error); ok {
			t.Fatal(err)
		}
		// t.Logf("%#v\n", v)
		b, err := json.MarshalIndent(v, "", "  ")
		if err != nil {
			t.Fatal(err)
		}
		t.Logf("result: \n%v", string(b))
	}

	iter2 := query2.RunWithContext(context.TODO(), items.UnstructuredContent()) // or query.RunWithContext
	for {
		v, ok := iter2.Next()
		if !ok {
			break
		}
		if err, ok := v.(error); ok {
			t.Fatal(err)
		}
		// t.Logf("%#v\n", v)
		b, err := json.MarshalIndent(v, "", "  ")
		if err != nil {
			t.Fatal(err)
		}
		t.Logf("result: \n%v", string(b))
	}

	iter3 := query3.RunWithContext(context.TODO(), items.UnstructuredContent()) // or query.RunWithContext
	for {
		v, ok := iter3.Next()
		if !ok {
			break
		}
		if err, ok := v.(error); ok {
			t.Fatal(err)
		}
		// t.Logf("%#v\n", v)
		b, err := json.MarshalIndent(v, "", "  ")
		if err != nil {
			t.Fatal(err)
		}
		t.Logf("result: \n%v", string(b))
	}
}

func TestGoJq4(t *testing.T) {
	dc, err := dynamic.NewClient("")
	if err != nil {
		t.Fatal(err)
	}
	item, err := dc.Client.Resource(schema.GroupVersionResource{Group: "monitoring.coreos.com", Version: "v1", Resource: "servicemonitors"}).
		Namespace("prom").
		// Get(context.TODO(), "prometheus-kube-prometheus-prometheus", metav1.GetOptions{})
		Get(context.TODO(), "prometheus-prometheus-node-exporter", metav1.GetOptions{})
	if err != nil {
		t.Fatal(err)
	}

	query, err := gojq.Parse(src_endpoint)
	if err != nil {
		t.Fatal(err)
	}

	iter := query.RunWithContext(context.TODO(), item.UnstructuredContent())
	for {
		v, ok := iter.Next()
		if !ok {
			break
		}
		if err, ok := v.(error); ok {
			t.Fatal(err)
		}
		// t.Logf("%#v\n", v)
		b, err := json.MarshalIndent(v, "", "  ")
		if err != nil {
			t.Fatal(err)
		}
		t.Logf("result: \n%v", string(b))
	}

}

// alert_rules
// - rule deletion
// kubectl get PrometheusRule prometheus-kube-prometheus-prometheus -n prom
var src_ruledeletion = `.spec.groups[0].rules| map(.alert == "Prometheus is down") | index(true)`

// p8s_config
// kubectl get prometheus -n prom -o json
var src_retentionperiod = ".items[0].spec.retention"
var src_scrapinterval = ".items[0].spec.scrapeInterval"
var src_timeout = ".items[0].spec.scrapeTimeout"

// kubectl get prometheus prometheus-kube-prometheus-prometheus -n prom -o json
var src_retentionperiod2 = ".spec.retention"
var src_scrapinterval2 = ".spec.scrapeInterval"
var src_timeout2 = ".spec.scrapeTimeout"

// service_monitor
// kubectl get servicemonitor prometheus-prometheus-node-exporter -n prom -o json
var src_endpoint = `.spec.endpoints | map(.port == "mysql-exporter2") | index(true)`
