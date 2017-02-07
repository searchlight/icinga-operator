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
	"github.com/appscode/searchlight/pkg/client/influxdb"
	"github.com/appscode/searchlight/util"
	"github.com/influxdata/influxdb/client"
	"github.com/spf13/cobra"
)

type request struct {
	host          string
	a, b, c, d, e string
	r             string
	warning       string
	critical      string
	secret        string
}

func trunc(val float64) interface{} {
	intData := int64(val * 1000)
	return float64(intData) / 1000.0
}

func getInfluxDBClient(authData *influxdb.AuthInfo) (*client.Client, error) {
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

func getValue(con *client.Client, db string, req *request) (map[string]interface{}, error) {

	valMap := make(map[string]interface{})

	defer func() {
		if e := recover(); e != nil {
			fmt.Fprintln(os.Stdout, util.State[3], e)
			os.Exit(3)
		}
	}()

	if req.a != "" {
		data, err := getInfluxdbData(con, req.a, db, "A")
		if err != nil {
			return nil, err
		}
		valMap["A"] = data
	}

	if req.b != "" {
		data, err := getInfluxdbData(con, req.b, db, "B")
		if err != nil {
			return nil, err
		}
		valMap["B"] = data
	}

	if req.c != "" {
		data, err := getInfluxdbData(con, req.c, db, "C")
		if err != nil {
			return nil, err
		}
		valMap["C"] = data
	}

	if req.d != "" {
		data, err := getInfluxdbData(con, req.d, db, "D")
		if err != nil {
			return nil, err
		}
		valMap["D"] = data
	}

	if req.e != "" {
		data, err := getInfluxdbData(con, req.e, db, "E")
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

func CheckInfluxQuery(req *request) (util.IcingaState, interface{}) {
	authData, err := influxdb.GetInfluxDBSecretData(req.secret)
	if err != nil {
		return util.Unknown, err
	}

	if req.host != "" {
		authData.Host = req.host
	}
	if authData.Host == "" {
		return util.Unknown, "No InfluxDB host found"
	}
	client, err := getInfluxDBClient(authData)
	if err != nil {
		return util.Unknown, err
	}

	valMap, err := getValue(client, authData.Database, req)
	if err != nil {
		return util.Unknown, err
	}

	expression, err := govaluate.NewEvaluableExpression(req.r)
	if err != nil {
		return util.Unknown, err
	}

	if valMap["R"], err = expression.Evaluate(valMap); err != nil {
		return util.Unknown, err
	}
	valMap["R"] = trunc(valMap["R"].(float64))

	if req.critical != "" {
		isCritical, err := checkResult(req.critical, valMap)
		if err != nil {
			return util.Unknown, err.Error()
		}
		if isCritical {
			return util.Critical, nil
		}
	}

	if req.warning != "" {
		isWarning, err := checkResult(req.warning, valMap)
		if err != nil {
			return util.Unknown, err
		}
		if isWarning {
			return util.Warning, nil
		}
	}

	return util.Ok, "Fine"
}

func NewCmd() *cobra.Command {
	var req request

	c := &cobra.Command{
		Use:     "check_influx_query",
		Short:   "Check InfluxDB Query Data",
		Example: "",

		Run: func(cmd *cobra.Command, args []string) {
			flags.EnsureRequiredFlags(cmd, "secret", "R")
			flags.EnsureAlterableFlags(cmd, "A", "B", "C", "D", "E")
			flags.EnsureAlterableFlags(cmd, "warning", "critical")
			util.Output(CheckInfluxQuery(&req))
		},
	}

	c.Flags().StringVarP(&req.host, "influx_host", "H", "", "URL of InfluxDB host to query")
	c.Flags().StringVarP(&req.secret, "secret", "s", "", `Kubernetes secret name`)
	c.Flags().StringVar(&req.a, "A", "", "InfluxDB query A")
	c.Flags().StringVar(&req.b, "B", "", "InfluxDB query B")
	c.Flags().StringVar(&req.c, "C", "", "InfluxDB query C")
	c.Flags().StringVar(&req.d, "D", "", "InfluxDB query D")
	c.Flags().StringVar(&req.e, "E", "", "InfluxDB query E")
	c.Flags().StringVar(&req.r, "R", "", `Equation to evaluate result`)
	c.Flags().StringVarP(&req.warning, "warning", "w", "", `Warning query which returns [true/false]`)
	c.Flags().StringVarP(&req.critical, "critical", "c", "", `Critical query which returns [true/false]`)
	return c
}
