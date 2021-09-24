package aws_utils

import (
	"fmt"
	"reflect"
	"time"

	"github.com/aws/aws-sdk-go/aws/session"
)

func IsSlice(v interface{}) bool {
	return reflect.TypeOf(v).Kind() == reflect.Slice
}

func GetEnvSession(profile string) *session.Session {
	opts := session.Options{
		SharedConfigState: session.SharedConfigEnable,
	}

	if profile != "" {
		opts.Profile = profile
	}

	sess := session.Must(session.NewSessionWithOptions(opts))
	return sess
}

func GetPrettyUptime(lunch time.Time) string {
	var timeStr string

	timeDiff := time.Since(lunch)
	hoursDiff := timeDiff.Hours()
	minutesDiff := timeDiff.Minutes()

	// if uptime is less than 1 day
	if hoursDiff < 24 {

		if minutesDiff < 60 {
			timeStr = fmt.Sprintf("%.0f Minutes", hoursDiff)
		} else {
			timeStr = fmt.Sprintf("%.2f Hours", hoursDiff)
		}
	} else {
		daysDiff := hoursDiff / 24
		timeStr = fmt.Sprintf("%.0f Days", daysDiff)
	}
	return timeStr
}
