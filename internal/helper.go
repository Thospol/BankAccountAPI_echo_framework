package internal

func LogFormat() string {
	// Refer to https://github.com/tkuchiki/alp
	var format string
	format += "host:${host}\t"
	format += "uri:${uri}\t"
	format += "method:${method}\t"
	format += "status:${status}\t"
	return format
}
