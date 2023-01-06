package lib

import (
	"bytes"
	"crypto/md5"
	"fmt"
	"io"
)

// md5 hash a chunk of data
func calculate_strong_checksum(data []byte) string
{
	s := ""	
	reader := bytes.NewReader(data)
	h:= md5.New()
	h.Write(data)
	return h.Sum()

}