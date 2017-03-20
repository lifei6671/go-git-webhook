package goclient

import "time"

type ClientConfig struct {
	User string
	Password string
	Timeout time.Duration
}

type Client struct {

}

func Dial(addr string,port int,config *ClientConfig) (*Client,error){

	return &Client{},nil
}