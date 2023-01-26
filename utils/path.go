package utils

import (
	"path/filepath"
	"runtime"
)

func GetServiceRootPath() (string, error) {
	_, thisFile, _, _ := runtime.Caller(0)
	thisFilePath := filepath.Dir(thisFile)

	// we know that this file is currently inside
	// the path /ticken-event-service/utils
	// so the root is ../{this-path}

	// Dir is going to remove the last element of
	// the slash and after that we are going to
	// keep the root value
	serviceRootPath := filepath.Dir(thisFilePath)

	// Here it is! Probable one the most hardcoded
	// things I did, but i'm proud of this logic :)
	return serviceRootPath, nil
}
