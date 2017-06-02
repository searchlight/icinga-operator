package influxdb

import (
	"errors"
	"fmt"
	"strings"

	"github.com/appscode/searchlight/pkg/client/k8s"
	ini "github.com/vaughan0/go-ini"
)

const (
	admin string = ".admin"

	influxDBHost     string = "INFLUX_HOST"
	influxDBDatabase string = "INFLUX_DB"
	influxDBReadUser string = "INFLUX_READ_USER"
	influxDBReadPass string = "INFLUX_READ_PASSWORD"

	influxDBDefaultDatabase string = "k8s"
	influxDBHostPort               = 8086
)

type AuthInfo struct {
	Host     string
	Username string
	Password string
	Database string
}

func GetInfluxDBSecretData(secretName, namespace string) (*AuthInfo, error) {
	kubeClient, err := k8s.NewClient()
	if err != nil {
		return nil, err
	}

	secret, err := kubeClient.Client.Core().Secrets(namespace).Get(secretName)
	if err != nil {
		return nil, err
	}

	authData := new(AuthInfo)
	if data, found := secret.Data[admin]; found {
		dataReader := strings.NewReader(string(data))
		secretData, err := ini.Load(dataReader)
		if err != nil {
			return nil, err
		}

		if host, found := secretData.Get("", influxDBHost); found {
			authData.Host = fmt.Sprintf("%s:%d", host, influxDBHostPort)
		}

		if authData.Database, found = secretData.Get("", influxDBDatabase); !found {
			authData.Database = influxDBDefaultDatabase
		}
		if authData.Username, found = secretData.Get("", influxDBReadUser); !found {
			return nil, errors.New("No InfluxDB read user found")
		}
		if authData.Password, found = secretData.Get("", influxDBReadPass); !found {
			return nil, errors.New("No InfluxDB read password found")
		}
		return authData, nil
	}
	return nil, errors.New("Invalid InfluxDB secret")
}
