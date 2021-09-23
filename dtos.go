package main

import "time"

type (
	response struct {
		Timestamp time.Time `json:"timestamp"`
		Code      string    `json:"code"`
		Message   string    `json:"message"`
	}

	search struct {
		Filter     string   `json:"filter" validate:"required"`
		Attributes []string `json:"attributes" validate:"required"`
	}

	errorResponse struct {
		Tag     string `json:"tag"`
		Name    string `json:"name"`
		Message string `json:"message"`
	}

	errorValidator struct {
		Code           string           `json:"code"`
		Message        string           `json:"message"`
		ErrorResponses []*errorResponse `json:"errors"`
	}
)
