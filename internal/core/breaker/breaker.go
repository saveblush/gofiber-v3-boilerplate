package breaker

import "github.com/afex/hystrix-go/hystrix"

const (
	// employees hris
	BreakerNameEmployees = "employees"
)

func Init() {
	hystrix.ConfigureCommand(BreakerNameEmployees, hystrix.CommandConfig{
		Timeout:     1500,
		SleepWindow: 2000,
	})

	/* test
	hystrix.ConfigureCommand(BreakerNameEmployees, hystrix.CommandConfig{
		Timeout:                300,
		MaxConcurrentRequests:  2,
		RequestVolumeThreshold: 1,
		ErrorPercentThreshold:  30,
	})*/
}
