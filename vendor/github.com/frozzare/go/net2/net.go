package net2

import (
	"net"
	"strconv"
)

// RandomPort will find a free port that can be used.
func RandomPort() (int, error) {
	l, err := net.Listen("tcp", ":0")

	if err != nil {
		return 0, err
	}

	defer l.Close()

	s := l.Addr().String()

	_, p, err := net.SplitHostPort(s)

	if err != nil {
		return 0, err
	}

	p2, err := strconv.Atoi(p)

	if err != nil {
		return 0, err
	}

	return p2, nil
}
