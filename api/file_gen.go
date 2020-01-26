// Package api :: This is auto generated file, do not edit manually
package api

import "encoding/json"

// Conf receives custom request config definition, e.g. custom headers, custom OData mod
func (file *File) Conf(config *RequestConfig) *File {
	file.config = config
	return file
}

// Select adds $select OData modifier
func (file *File) Select(oDataSelect string) *File {
	file.modifiers.AddSelect(oDataSelect)
	return file
}

// Expand adds $expand OData modifier
func (file *File) Expand(oDataExpand string) *File {
	file.modifiers.AddExpand(oDataExpand)
	return file
}

/* Response helpers */

// Data response helper
func (fileResp *FileResp) Data() *FileInfo {
	data := NormalizeODataItem(*fileResp)
	res := &FileInfo{}
	json.Unmarshal(data, &res)
	return res
}

// Normalized returns normalized body
func (fileResp *FileResp) Normalized() []byte {
	return NormalizeODataItem(*fileResp)
}