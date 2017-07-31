package check_influx_query

import (
	"testing"
	"fmt"
)

func TestCheckInfluxQuery(t *testing.T) {
	req := &Request{
		Host: "a8682d825f99811e6899912f236046fb-1008332040.us-east-1.elb.amazonaws.com:8086",
		Namespace: "kube-system",
		A: "select value from \"memory/limit\" where pod_name='nginx-jqk30' and pod_namespace='kube-system';",
		R: "A",
		Warning: "R >= 0",
		SecretName: "appscode-influx",
	}

	fmt.Println(CheckInfluxQuery(req))
}
