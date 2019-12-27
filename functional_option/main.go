package main

import (
	"log"
)

type config struct {
	host     string
	port     int
	username string
}

type Option func(*config) error

func OptionHost(host string) Option {
	return func(conf *config) error {
		conf.host = host
		return nil
	}
}

func OptionPort(port int) Option {
	return func(conf *config) error {
		conf.port = port
		return nil
	}
}

func OptionUsername(username string) Option {
	return func(conf *config) error {
		conf.username = username
		return nil
	}
}

type Connection struct {
	conf config
}

func NewConnection(options ...Option) (*Connection, error) {
	var ret Connection
	for _, opt := range options {
		err := opt(&ret.conf)
		if err != nil {
			return nil, err
		}
	}

	return &ret, nil
}

func main() {
	conn, err := NewConnection(OptionHost("nice"), OptionPort(11001))
	if err != nil {
		log.Fatal(err)
	}

	log.Println(conn)
}
