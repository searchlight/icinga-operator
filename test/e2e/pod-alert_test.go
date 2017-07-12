package e2e_test

import (
	"fmt"

	. "github.com/onsi/ginkgo"
)

func countIcinga2Host() {
	fmt.Println("Counting Icinga2 host")
}

func countIcinga2Service() {
	fmt.Println("Counting Icinga2 service")
}

func incrementReplica() {
	fmt.Println("Incrementing deployment replica")
}

var _ = Describe("PodAlert", func() {

	Describe(`Create`, func() {
		Context(`a`, func() {
			It(`namespace`, func() {
			})
			It(`deployment`, func() {
			})

			It(`PodAlert`, func() {
			})
		})
	})

	Describe(`Check`, func() {
		Context(`Checking Icinga2`, func() {
			It(`Counting Icinga2 Host`, countIcinga2Host)
			It(`Counting Icinga2 Service`, countIcinga2Service)
		})

		Context(`Increment replica`, func() {
			It(`Incrementing replica`, incrementReplica)
		})

		Context(`Checking Icinga2 again`, func() {
			It(`Counting Icinga2 Host`, countIcinga2Host)
			It(`Counting Icinga2 Service`, countIcinga2Service)
		})
	})
	//It(`Deleting namespace`, func() {})
})
