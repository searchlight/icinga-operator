package framework

import (
	"errors"
	"fmt"
	"time"

	. "github.com/onsi/gomega"
)

func (f *Framework) EventuallyIcingaAPI() GomegaAsyncAssertion {
	return Eventually(
		func() error {
			if f.icingaClient.Check().Get(nil).Do().Status == 200 {
				PrintSeparately("Connected to icinga api")
				return nil
			}
			fmt.Println("Waiting for icinga to start")
			return errors.New("icigna is not ready")
		},
		time.Minute*10,
		time.Second*10,
	)
}
