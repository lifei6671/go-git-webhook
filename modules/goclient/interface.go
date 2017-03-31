package goclient

import "net/url"

type ClientInterface interface {
	Command (host url.URL,account,password ,shell string,channel chan <-[]byte)
}
