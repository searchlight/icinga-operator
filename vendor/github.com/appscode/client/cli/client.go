package cli

import (
	"errors"

	"github.com/appscode/client"
	term "github.com/appscode/go-term"
	"google.golang.org/grpc/grpclog"
)

func Client(ua string) (*client.Client, error) {
	rc, err := LoadApprc()
	if err != nil {
		return nil, err
	}
	a := rc.GetAuth()
	if a == nil {
		return nil, errors.New("Command requires authentication, please run `appctl login`")
	}
	options := client.NewOption(a.ApiServer)
	options.BearerAuth(a.TeamId, a.Token)
	if ua != "" {
		options.UserAgent(ua)
	}
	grpclog.SetLogger(&term.NullLogger{})
	c, err := client.New(options)
	if err != nil {
		return nil, err
	}
	return c, nil
}

func ClientOrDie(ua string) *client.Client {
	c, err := Client(ua)
	term.ExitOnError(err)
	return c
}
