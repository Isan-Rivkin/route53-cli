package aws_utils

import (
	"github.com/aws/aws-sdk-go/aws/session"
)

func GetEnvSession() *session.Session {
	sess := session.Must(session.NewSessionWithOptions(session.Options{
		SharedConfigState: session.SharedConfigEnable,
	}))
	return sess
}
