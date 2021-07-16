package middleware

import (
	"bytes"
	"compress/gzip"
	"io/ioutil"
)

// RpcDecompress RPC数据解密
func RpcDecompress(data []byte, err error) ([]byte, error) {
	rdata := bytes.NewReader(data)
	reader, err := gzip.NewReader(rdata)
	if err != nil {
		return data, err
	}
	defer func(reader *gzip.Reader) {
		_ = reader.Close()
	}(reader)
	dataNow, errs := ioutil.ReadAll(reader)
	if errs != nil {
		return data, errs
	}
	return dataNow, err
}

// RpcCompress RPC数据加密
func RpcCompress(data []byte, err error) ([]byte, error) {
	var b bytes.Buffer
	gz := gzip.NewWriter(&b)
	if _, err := gz.Write(data); err != nil {
		return data, err
	}
	if err := gz.Flush(); err != nil {
		return data, err
	}
	if err := gz.Close(); err != nil {
		return data, err
	}
	return b.Bytes(), err
}
