package go_utils

import (
	"crypto/rand"
	"io"
	"encoding/base64"
)

/*
 uuid utils
 @author Tony Tian
 @date 2018-04-14
 @version 1.0.0
*/

/*
  return eg: 4725f5ae6a350b1c45687c9934456e6f
*/
func GUID() (string, error) {
	b := make([]byte, 48)
	_, err := io.ReadFull(rand.Reader, b)
	if err != nil {
		return "", err
	}
	guid := Md5(base64.URLEncoding.EncodeToString(b))
	return guid, nil
}

/*
 return eg: e7486845-9f24-c3d8-0db1-fe61e25c88a2
*/
func UUID() (string, error) {
	guid, err := GUID()
	if err != nil {
		return "", err
	}
	str := NewString(guid)
	builder := NewStringBuilder()
	builder.Append(str.Substring(0, 8).ToString())
	builder.Append(HLINE).Append(str.Substring(8, 12).ToString())
	builder.Append(HLINE).Append(str.Substring(12, 16).ToString())
	builder.Append(HLINE).Append(str.Substring(16, 20).ToString())
	builder.Append(HLINE).Append(str.SubstringBegin(20).ToString())
	return builder.ToString(), nil
}
