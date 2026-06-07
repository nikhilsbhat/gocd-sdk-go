package gocd

import "github.com/jinzhu/copier"

var copyClientWithOption = func(dst, src any) error {
	return copier.CopyWithOption(dst, src, copier.Option{IgnoreEmpty: true, DeepCopy: true})
}

func (conf *client) clone() *client {
	newClient := &client{}
	if err := copyClientWithOption(newClient, conf); err != nil {
		panic(err)
	}

	return newClient
}
