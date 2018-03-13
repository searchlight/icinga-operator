package check_influx_query

import (
	"encoding/json"
	"errors"
	"fmt"
	"net/url"
	"os"
	"strconv"

	"github.com/Knetic/govaluate"
	"github.com/appscode/go/flags"
	"github.com/appscode/searchlight/pkg/icinga"
	"github.com/influxdata/influxdb/client"
	"github.com/spf13/cobra"
)

type Request struct {
	masterURL      string
	kubeconfigPath string

	Host          string
	A, B, C, D, E string
	R             string
	Warning       string
	Critical      string
	SecretName    string
	Namespace     string
}

func trunc(val float64) interface{} {
	intData := int64(val * 1000)
	return float64(intData) / 1000.0
}

func getInfluxDBClient(authData *AuthInfo) (*client.Client, error) {
	config := &client.Config{
		URL: url.URL{
			Scheme: "http",
			Host:   authData.Host,
		},
		Username: authData.Username,
		Password: authData.Password,
	}
	client, err := client.NewClient(*config)
	if err != nil {
		return nil, err
	}
	return client, nil
}

func getInfluxdbData(con *client.Client, command, db, queryName string) (float64, error) {
	q := client.Query{
		Command:  command,
		Database: db,
	}
	res, err := con.Query(q)
	if err != nil {
		return 0.0, err
	}

	if len(res.Results[0].Series) == 0 {
		return 0.0, errors.New(fmt.Sprint("Value not found for query: ", queryName))
	}
	data, err := strconv.ParseFloat((string(res.Results[0].Series[0].Values[0][1].(json.Number))), 64)
	if err != nil {
		return 0.0, err
	}
	return data, nil
}

func getValue(con *client.Client, db string, req *Request) (map[string]interface{}, error) {
	valMap := make(map[string]interface{})

	defer func() {
		if e := recover(); e != nil {
			fmt.Fprintln(os.Stdout, icinga.Warning, e)
			os.Exit(3)
		}
	}()

	if req.A != "" {
		data, err := getInfluxdbData(con, req.A, db, "A")
		if err != nil {
			return nil, err
		}
		valMap["A"] = data
	}

	if req.B != "" {
		data, err := getInfluxdbData(con, req.B, db, "B")
		if err != nil {
			return nil, err
		}
		valMap["B"] = data
	}

	if req.C != "" {
		data, err := getInfluxdbData(con, req.C, db, "C")
		if err != nil {
			return nil, err
		}
		valMap["C"] = data
	}

	if req.D != "" {
		data, err := getInfluxdbData(con, req.D, db, "D")
		if err != nil {
			return nil, err
		}
		valMap["D"] = data
	}

	if req.E != "" {
		data, err := getInfluxdbData(con, req.E, db, "E")
		if err != nil {
			return nil, err
		}
		valMap["E"] = data
	}

	return valMap, nil
}

func checkResult(checkQuery string, valueMap map[string]interface{}) (bool, error) {
	expr, err := govaluate.NewEvaluableExpression(checkQuery)
	if err != nil {
		return false, err
	}

	res, err := expr.Evaluate(valueMap)
	if err != nil {
		return false, err
	}

	if res.(bool) {
		return true, nil
	}
	return false, nil
}

func CheckInfluxQuery(req *Request) (icinga.State, interface{}) {
	authData, err := GetInfluxDBSecretData(req, req.SecretName, req.Namespace)
	if err != nil {
		return icinga.Unknown, err
	}

	if req.Host != "" {
		authData.Host = req.Host
	}
	if authData.Host == "" {
		return icinga.Unknown, "No InfluxDB host found"
	}
	client, err := getInfluxDBClient(authData)
	if err != nil {
		return icinga.Unknown, err
	}

	valMap, err := getValue(client, authData.Database, req)
	if err != nil {
		return icinga.Unknown, err
	}

	expression, err := govaluate.NewEvaluableExpression(req.R)
	if err != nil {
		return icinga.Unknown, err
	}

	if valMap["R"], err = expression.Evaluate(valMap); err != nil {
		return icinga.Unknown, err
	}
	valMap["R"] = trunc(valMap["R"].(float64))

	if req.Critical != "" {
		isCritical, err := checkResult(req.Critical, valMap)
		if err != nil {
			return icinga.Unknown, err.Error()
		}
		if isCritical {
			return icinga.Critical, nil
		}
	}

	if req.Warning != "" {
		isWarning, err := checkResult(req.Warning, valMap)
		if err != nil {
			return icinga.Unknown, err
		}
		if isWarning {
			return icinga.Warning, nil
		}
	}

	return icinga.OK, "Fine"
}

func NewCmd() *cobra.Command {
	var req Request
	var icingaHost string

	c := &cobra.Command{
		Use:     "check_influx_query",
		Short:   "Check InfluxDB Query Data",
		Example: "",

		Run: func(cmd *cobra.Command, args []string) {
			flags.EnsureRequiredFlags(cmd, "host")

			host, err := icinga.ParseHost(icingaHost)
			if err != nil {
				fmt.Fprintln(os.Stdout, icinga.Warning, "Invalid icinga host.name")
				os.Exit(3)
			}
			req.Namespace = host.AlertNamespace

			flags.EnsureAlterableFlags(cmd, "A", "B", "C", "D", "E")
			flags.EnsureAlterableFlags(cmd, "warning", "critical")
			icinga.Output(CheckInfluxQuery(&req))
		},
	}

	c.Flags().StringVar(&req.masterURL, "master", req.masterURL, "The address of the Kubernetes API server (overrides any value in kubeconfig)")
	c.Flags().StringVar(&req.kubeconfigPath, "kubeconfig", req.kubeconfigPath, "Path to kubeconfig file with authorization information (the master location is set by the master flag).")

	c.Flags().StringVarP(&icingaHost, "host", "H", "", "Icinga host name")
	c.Flags().StringVar(&req.Host, "influxHost", "", "URL of InfluxDB host to query")
	c.Flags().StringVarP(&req.SecretName, "secretName", "s", "", `Kubernetes secret name`)
	c.Flags().StringVar(&req.A, "A", "", "InfluxDB query A")
	c.Flags().StringVar(&req.B, "B", "", "InfluxDB query B")
	c.Flags().StringVar(&req.C, "C", "", "InfluxDB query C")
	c.Flags().StringVar(&req.D, "D", "", "InfluxDB query D")
	c.Flags().StringVar(&req.E, "E", "", "InfluxDB query E")
	c.Flags().StringVar(&req.R, "R", "", `Equation to evaluate result`)
	c.Flags().StringVarP(&req.Warning, "warning", "w", "", `Warning query which returns [true/false]`)
	c.Flags().StringVarP(&req.Critical, "critical", "c", "", `Critical query which returns [true/false]`)
	return c
}
