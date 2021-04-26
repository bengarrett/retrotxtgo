package version

import "errors"

var (
	ErrCacheYaml = errors.New("set cache cannot marshal yaml")
	ErrCacheData = errors.New("set cache cannot create a data path")
	ErrCacheSave = errors.New("set cache cannot save data")
	ErrMarshal   = errors.New("version could not marshal")
)
