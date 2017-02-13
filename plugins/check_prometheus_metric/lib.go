package check_prometheus_metric

import (
	"context"
	"errors"
	"fmt"
	"os"
	"strconv"
	"time"

	"github.com/appscode/go/flags"
	"github.com/appscode/searchlight/util"
	"github.com/prometheus/client_golang/api/prometheus"
	"github.com/prometheus/common/model"
	"github.com/spf13/cobra"
)

type Request struct {
	Host       string
	Query      string
	MetricName string
	Method     string
	Warning    int64
	Critical   int64
	AcceptNan  bool
}

func checkResult(method string, valueToCheck, result int64) bool {
	switch method {
	case "gt":
		if result > valueToCheck {
			return true
		}
	case "ge":
		if result >= valueToCheck {
			return true
		}
	case "lt":
		if result < valueToCheck {
			return true
		}
	case "le":
		if result <= valueToCheck {
			return true
		}
	case "eq":
		if result == valueToCheck {
			return true
		}
	case "ne":
		if result != valueToCheck {
			return true
		}
	}
	return false
}

func CheckPrometheusMetric(req *Request) (util.IcingaState, interface{}) {
	config := prometheus.Config{
		Address: req.Host,
	}
	client, err := prometheus.New(config)
	if err != nil {
		return util.Unknown, err
	}
	queryApi := prometheus.NewQueryAPI(client)

	result, err := queryApi.Query(context.Background(), req.Query, time.Now())
	if err != nil {
		return util.Unknown, err
	}

	vector := result.(model.Vector)

	if len(vector) == 0 {
		if req.AcceptNan {
			return util.Ok, errors.New("NaN")
		}
		return util.Unknown, errors.New("No data found")
	} else if len(vector) > 1 {
		return util.Unknown, errors.New("Invalid query.\nQuery should return single float or int")
	}

	value := vector[0].Value.String()
	val, err := strconv.ParseFloat(value, 64)
	if err != nil {
		return util.Unknown, err
	}

	valInt64 := int64(val)
	outputStr := fmt.Sprintf("%v", valInt64)
	if req.MetricName != "" {
		outputStr = fmt.Sprintf("%v is %v", req.MetricName, valInt64)
	}

	if checkResult(req.Method, req.Critical, valInt64) {
		return util.Critical, outputStr
	}
	if checkResult(req.Method, req.Warning, valInt64) {
		return util.Warning, outputStr
	}

	return util.Ok, outputStr
}

func NewCmd() *cobra.Command {
	var req Request

	c := &cobra.Command{
		Use:     "check_prometheus_metric",
		Short:   "Check prometheus metric",
		Example: "",

		Run: func(cmd *cobra.Command, args []string) {
			if req.Host == "" {
				fmt.Fprintln(os.Stdout, util.State[3], errors.New("No prometheus host found"))
				os.Exit(3)
			}
			flags.EnsureRequiredFlags(cmd, "query", "warning", "critical")
			util.Output(CheckPrometheusMetric(&req))
		},
	}

	c.Flags().StringVarP(&req.Host, "prom_host", "H", "", "URL of Prometheus host to query")
	c.Flags().StringVarP(&req.Query, "query", "q", "", "Prometheus query that returns a float or int")
	c.Flags().StringVarP(&req.MetricName, "metric_name", "n", "", "A name for the metric being checked")
	c.Flags().StringVarP(&req.Method, "method", "m", "ge", `Comparison method, one of gt, ge, lt, le, eq, ne
	(defaults to ge unless otherwise specified)`)
	c.Flags().BoolVarP(&req.AcceptNan, "accept_nan", "O", false, `Accept NaN as an "OK" result`)
	c.Flags().Int64VarP(&req.Warning, "warning", "w", 0, "Warning level value (must be zero or positive)")
	c.Flags().Int64VarP(&req.Critical, "critical", "c", 0, "Critical level value (must be zero or positive)")
	return c
}
