/*
 * MinIO Go Library for Amazon S3 Compatible Cloud Storage
 * Copyright 2015-2020 MinIO, Inc.
 *
 * Licensed under the Apache License, Version 2.0 (the "License");
 * you may not use this file except in compliance with the License.
 * You may obtain a copy of the License at
 *
 *     http://www.apache.org/licenses/LICENSE-2.0
 *
 * Unless required by applicable law or agreed to in writing, software
 * distributed under the License is distributed on an "AS IS" BASIS,
 * WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
 * See the License for the specific language governing permissions and
 * limitations under the License.
 */

package minio

import (
	"net/http"
	"net/url"
)

// PromptObjectOptions provides options to PromptObject call.
type PromptObjectOptions struct {
	lambdaArn string
	kwargs    map[string]string
	headers   map[string]string
	reqParams url.Values
}

// Header returns the http.Header representation of the GET options.
func (o PromptObjectOptions) Header() http.Header {
	headers := make(http.Header, len(o.headers))
	for k, v := range o.headers {
		headers.Set(k, v)
	}
	return headers
}

// Add a key value pair to the kwargs where the key is a string but the value can be any type.
func (o *PromptObjectOptions) AddKwarg(key string, value string) {
	if o.kwargs == nil {
		o.kwargs = make(map[string]string)
	}
	o.kwargs[key] = value
}

// Set adds a key value pair to the options. The
// key-value pair will be part of the HTTP GET request
// headers.
func (o *PromptObjectOptions) SetHeader(key, value string) {
	if o.headers == nil {
		o.headers = make(map[string]string)
	}
	o.headers[http.CanonicalHeaderKey(key)] = value
}

// SetReqParam - set request query string parameter
// supported key: see supportedQueryValues and allowedCustomQueryPrefix.
// If an unsupported key is passed in, it will be ignored and nothing will be done.
func (o *PromptObjectOptions) SetReqParam(key, value string) {
	if !isCustomQueryValue(key) && !isStandardQueryValue(key) {
		// do nothing
		return
	}
	if o.reqParams == nil {
		o.reqParams = make(url.Values)
	}
	o.reqParams.Set(key, value)
}

// AddReqParam - add request query string parameter
// supported key: see supportedQueryValues and allowedCustomQueryPrefix.
// If an unsupported key is passed in, it will be ignored and nothing will be done.
func (o *PromptObjectOptions) AddReqParam(key, value string) {
	if !isCustomQueryValue(key) && !isStandardQueryValue(key) {
		// do nothing
		return
	}
	if o.reqParams == nil {
		o.reqParams = make(url.Values)
	}
	o.reqParams.Add(key, value)
}

// toQueryValues - Convert the reqParams in Options to query string parameters.
func (o *PromptObjectOptions) toQueryValues() url.Values {
	urlValues := make(url.Values)
	if o.reqParams != nil {
		for key, values := range o.reqParams {
			for _, value := range values {
				urlValues.Add(key, value)
			}
		}
	}

	return urlValues
}
