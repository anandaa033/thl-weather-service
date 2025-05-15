package state

type Type string

const (
	LOCAL Type = "LOCAL"
	DEV   Type = "DEV"
	SIT   Type = "SIT"
	UAT   Type = "UAT"
	PROD  Type = "PROD"
)

func (o *Type) String() string {
	return string(*o)
}
