// Package v1alpha1 provides primitives to interact with the openapi HTTP API.
//
// Code generated by github.com/oapi-codegen/oapi-codegen/v2 version v2.3.0 DO NOT EDIT.
package v1alpha1

import (
	"bytes"
	"compress/gzip"
	"encoding/base64"
	"fmt"
	"net/url"
	"path"
	"strings"

	"github.com/getkin/kin-openapi/openapi3"
)

// Base64 encoded, gzipped, json marshaled Swagger object
var swaggerSpec = []string{

	"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",
	"DgAA//8nCwvalzAAAA==",
}

// GetSwagger returns the content of the embedded swagger specification file
// or error if failed to decode
func decodeSpec() ([]byte, error) {
	zipped, err := base64.StdEncoding.DecodeString(strings.Join(swaggerSpec, ""))
	if err != nil {
		return nil, fmt.Errorf("error base64 decoding spec: %w", err)
	}
	zr, err := gzip.NewReader(bytes.NewReader(zipped))
	if err != nil {
		return nil, fmt.Errorf("error decompressing spec: %w", err)
	}
	var buf bytes.Buffer
	_, err = buf.ReadFrom(zr)
	if err != nil {
		return nil, fmt.Errorf("error decompressing spec: %w", err)
	}

	return buf.Bytes(), nil
}

var rawSpec = decodeSpecCached()

// a naive cached of a decoded swagger spec
func decodeSpecCached() func() ([]byte, error) {
	data, err := decodeSpec()
	return func() ([]byte, error) {
		return data, err
	}
}

// Constructs a synthetic filesystem for resolving external references when loading openapi specifications.
func PathToRawSpec(pathToFile string) map[string]func() ([]byte, error) {
	res := make(map[string]func() ([]byte, error))
	if len(pathToFile) > 0 {
		res[pathToFile] = rawSpec
	}

	return res
}

// GetSwagger returns the Swagger specification corresponding to the generated code
// in this file. The external references of Swagger specification are resolved.
// The logic of resolving external references is tightly connected to "import-mapping" feature.
// Externally referenced files must be embedded in the corresponding golang packages.
// Urls can be supported but this task was out of the scope.
func GetSwagger() (swagger *openapi3.T, err error) {
	resolvePath := PathToRawSpec("")

	loader := openapi3.NewLoader()
	loader.IsExternalRefsAllowed = true
	loader.ReadFromURIFunc = func(loader *openapi3.Loader, url *url.URL) ([]byte, error) {
		pathToFile := url.String()
		pathToFile = path.Clean(pathToFile)
		getSpec, ok := resolvePath[pathToFile]
		if !ok {
			err1 := fmt.Errorf("path not found: %s", pathToFile)
			return nil, err1
		}
		return getSpec()
	}
	var specData []byte
	specData, err = rawSpec()
	if err != nil {
		return
	}
	swagger, err = loader.LoadFromData(specData)
	if err != nil {
		return
	}
	return
}
