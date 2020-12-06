// brigolagecms/pkg/http/jsonapi/jsonapi.go
//
// Copyright (c) 2020, Michael D Henderson.
// All rights reserved.
//
// Redistribution and use in source and binary forms, with or without
// modification, are permitted provided that the following conditions are met:
//
//  1. Redistributions of source code must retain the above copyright notice, this
//     list of conditions and the following disclaimer.
//
//  2. Redistributions in binary form must reproduce the above copyright notice,
//     this list of conditions and the following disclaimer in the documentation
//     and/or other materials provided with the distribution.
//
//  3. Neither the name of the copyright holder nor the names of its
//     contributors may be used to endorse or promote products derived from
//     this software without specific prior written permission.
//
// THIS SOFTWARE IS PROVIDED BY THE COPYRIGHT HOLDERS AND CONTRIBUTORS "AS IS"
// AND ANY EXPRESS OR IMPLIED WARRANTIES, INCLUDING, BUT NOT LIMITED TO, THE
// IMPLIED WARRANTIES OF MERCHANTABILITY AND FITNESS FOR A PARTICULAR PURPOSE ARE
// DISCLAIMED. IN NO EVENT SHALL THE COPYRIGHT HOLDER OR CONTRIBUTORS BE LIABLE
// FOR ANY DIRECT, INDIRECT, INCIDENTAL, SPECIAL, EXEMPLARY, OR CONSEQUENTIAL
// DAMAGES (INCLUDING, BUT NOT LIMITED TO, PROCUREMENT OF SUBSTITUTE GOODS OR
// SERVICES; LOSS OF USE, DATA, OR PROFITS; OR BUSINESS INTERRUPTION) HOWEVER
// CAUSED AND ON ANY THEORY OF LIABILITY, WHETHER IN CONTRACT, STRICT LIABILITY,
// OR TORT (INCLUDING NEGLIGENCE OR OTHERWISE) ARISING IN ANY WAY OUT OF THE USE
// OF THIS SOFTWARE, EVEN IF ADVISED OF THE POSSIBILITY OF SUCH DAMAGE.

// Package jsonapi tries to implement JSONAPI.
package jsonapi

import (
	"encoding/json"
	"fmt"
	"log"
	"net/http"
)

type ErrorObject struct {
	Id     string       `json:"id,omitempty"`
	Status string       `json:"status,omitempty"`
	Code   string       `json:"code,omitempty"`
	Title  string       `json:"title,omitempty"`
	Detail string       `json:"detail,omitempty"`
	Source *ErrorSource `json:"source,omitempty"`
}
type ErrorSource struct {
	Pointer   string `json:"pointer,omitempty"`
	Parameter string `json:"parameter,omitempty"`
}

func Error(w http.ResponseWriter, r *http.Request, status int, text string, details ...string) {
	var failure struct {
		Status string        `json:"status"`
		Errors []ErrorObject `json:"errors"`
		Links  struct {
			Self string `json:"self"`
		} `json:"links"`
	}
	failure.Status = http.StatusText(status)
	failure.Links.Self = r.URL.Path

	// the first error, by convention, is always the http status being returned
	if text == "" {
		text = http.StatusText(status)
	}
	failure.Errors = append(failure.Errors, ErrorObject{
		Status: fmt.Sprintf("%d", status),
		Detail: text,
		Source: &ErrorSource{Parameter: r.URL.Path},
	})

	// then append any error details that the user provided
	for _, detail := range details {
		failure.Errors = append(failure.Errors, ErrorObject{
			Detail: detail,
		})
	}

	w.Header().Set("Content-Type", "application/vnd.api+json")
	w.WriteHeader(status)

	if err := json.NewEncoder(w).Encode(failure); err != nil {
		log.Printf("[http] error writing response: %+v\n", err)
	}
}

func Errors(w http.ResponseWriter, r *http.Request, status int, errors ...ErrorObject) {
	var failure struct {
		Status string        `json:"status"`
		Errors []ErrorObject `json:"errors"`
		Links  struct {
			Self string `json:"self"`
		} `json:"links"`
	}
	failure.Status = http.StatusText(status)
	failure.Links.Self = r.URL.Path

	// then append any error details that the user provided
	for _, err := range errors {
		failure.Errors = append(failure.Errors, err)
	}

	w.Header().Set("Content-Type", "application/vnd.api+json")
	w.WriteHeader(status)

	if err := json.NewEncoder(w).Encode(failure); err != nil {
		log.Printf("[http] error writing response: %+v\n", err)
	}
}

func NoContent(w http.ResponseWriter, r *http.Request) {
	w.WriteHeader(http.StatusNoContent)
}

func Ok(w http.ResponseWriter, r *http.Request, status int, data interface{}) {
	w.Header().Set("Content-Type", "application/vnd.api+json")
	w.WriteHeader(status)

	var success struct {
		Status string      `json:"status"`
		Data   interface{} `json:"data"`
		Links  struct {
			Self string `json:"self"`
		} `json:"links"`
	}
	success.Status = "ok"
	success.Data = data
	success.Links.Self = r.URL.Path

	if err := json.NewEncoder(w).Encode(success); err != nil {
		log.Printf("[http] error writing response: %+v\n", err)
	}
}
