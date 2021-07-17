package aws_utils

import (
	"github.com/aws/aws-sdk-go/aws/session"
)

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
