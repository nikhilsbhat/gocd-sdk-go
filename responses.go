package main

import (
	"fmt"
)

var (
	goCdTryMessage       = `trying to retrieve %s information present in GoCD`
	goCdSuccessMessage   = `successfully retrieved information of %s from GoCD`
	cronMessage          = `%s cron will be scheduled for %s as specified`
	cronScheduledMessage = `%d instance of %s cron got scheduled`
	cronCompleteMessage  = `scheduled %d instance of %s cron was completed`
	symlinkMessage       = `path %s is link to %s so fetching size of destination`
	cronEnabled          = `cron is enabled for %s metric collection`
)

func apiError(component, errMsg string) error {
	return fmt.Errorf("retrieving %s information errored with: %s", component, errMsg) //nolint:goerr113
}

func apiWithCodeError(code int) error {
	return fmt.Errorf("goCd server returned code %d with message", code) //nolint:goerr113
}

func responseReadError(msg string) error {
	return fmt.Errorf("reading resposne body errored with: %s", msg) //nolint:goerr113
}

func getTryMessages(component string) string {
	return fmt.Sprintf(goCdTryMessage, component)
}

func getSuccessMessages(component string) string {
	return fmt.Sprintf(goCdSuccessMessage, component)
}

func getCronMessages(component, schedule string) string {
	return fmt.Sprintf(cronMessage, component, schedule)
}

func getCronEnbaled(component string) string {
	return fmt.Sprintf(cronEnabled, component)
}

func getCronScheduledMessage(component string, instance int) string {
	return fmt.Sprintf(cronScheduledMessage, instance, component)
}

func getCronCompleteMessage(component string, instance int) string {
	return fmt.Sprintf(cronCompleteMessage, instance, component)
}

func getLinkMessage(link, path string) string {
	return fmt.Sprintf(symlinkMessage, link, path)
}

// func getTraceInfo(component string, request *resty.Request) string {
//	out, err := json.Marshal(request.TraceInfo())
//	if err != nil {
//		log.Println(err)
//		return ""
//	}
//	return fmt.Sprintf("%s: %s", component, string(out))
// }
