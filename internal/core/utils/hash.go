package utils

import "github.com/bytedance/sonic"

// Marshal marshal
func Marshal(val interface{}) ([]byte, error) {
	d, err := sonic.Marshal(&val)
	if err != nil {
		return nil, err
	}

	return d, nil
}

// Unmarshal unmarshal
func Unmarshal(buf []byte, val interface{}) error {
	err := sonic.Unmarshal(buf, val)
	if err != nil {
		return err
	}

	return nil
}
