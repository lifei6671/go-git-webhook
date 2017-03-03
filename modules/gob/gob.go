package gob

import (
	"bytes"
	"encoding/gob"
)

func Decode(value string,r interface{}) (error) {

	network := bytes.NewBuffer([]byte(value));

	dec := gob.NewDecoder(network)

	return dec.Decode(r);
}

//
func Encode(value interface{}) (string,error) {
	network:= bytes.NewBuffer(nil);

	enc := gob.NewEncoder(network)

	err := enc.Encode(value);
	if err != nil {
		return  "",err;
	}

	return  network.String(),nil;
}