package dsn

import (
	"net/url"
	"strings"
)

type DSN struct {
	scheme   string
	host     string
	port     string
	user     string
	password string
	path     string
	query    url.Values
}

func (dsn *DSN) Host() string {
	return dsn.host
}

func (dsn *DSN) Port() string {
	return dsn.port
}

func (dsn *DSN) Short() string {
	return dsn.scheme + "://" + dsn.host + ":" + dsn.port
}

func (dsn *DSN) Query(key string) string {
	return dsn.query[key][0]
}

func (dsn *DSN) Socket() string {
	return dsn.host + ":" + dsn.port
}

func (dsn *DSN) User() string {
	return dsn.user
}

func (dsn *DSN) Password() string {
	return dsn.password
}

func Decode(dsn string) (d DSN) {

	var connect string

	parts := strings.SplitN(dsn, "?", 2)
	if len(parts) == 2 {
		connect = parts[0]
		d.query, _ = url.ParseQuery(parts[1])
	} else {
		connect = dsn
	}

	parts = strings.SplitN(connect, "://", 2)
	if len(parts) == 0 {
		d.scheme = ""
		connect = dsn
	} else if len(parts) == 1 {
		d.scheme = ""
		connect = parts[0]
	} else {
		d.scheme = parts[0]
		connect = parts[1]
	}
	parts = strings.SplitN(connect, "@", 2)
	var socket, credential string
	if len(parts) <= 1 {
		credential = ""
		socket = connect
	} else {
		credential = parts[0]
		socket = parts[1]
	}
	if credential == "" {
		d.user = ""
		d.password = ""
	} else {
		parts = strings.SplitN(credential, ":", 2)
		if len(parts) <= 1 {
			d.user = credential
			d.password = ""
		} else {
			d.user = parts[0]
			d.password = parts[1]
		}
	}
	if socket == "" {
		d.host = ""
		d.port = ""
	} else {
		parts = strings.SplitN(socket, ":", 2)
		if len(parts) <= 1 {
			d.host = socket
			d.port = ""
		} else {
			d.host = parts[0]
			d.port = parts[1]
			if strings.Contains(d.port, "/") {
				d.port, d.path, _ = strings.Cut(d.port, "/")
			}
		}
	}
	return
}
