package services

import (
	"webman/pkg/models"
)

// DefaultHeaders provides a list of commonly used HTTP headers
var DefaultHeaders = []models.DefaultHeader{
	{
		Name:        "Accept",
		Value:       "application/json",
		Description: "Indicates that the client expects JSON response",
	},
	{
		Name:        "Content-Type",
		Value:       "application/json",
		Description: "Indicates that the request body is in JSON format",
	},
	{
		Name:        "Authorization",
		Value:       "Bearer ",
		Description: "Bearer token authentication",
	},
	{
		Name:        "Cache-Control",
		Value:       "no-cache",
		Description: "Controls caching behavior",
	},
	{
		Name:        "User-Agent",
		Value:       "Webman/1.0.0",
		Description: "Identifies the client application",
	},
	{
		Name:        "Accept-Language",
		Value:       "en-US,en;q=0.9",
		Description: "Preferred language for response",
	},
	{
		Name:        "X-Requested-With",
		Value:       "XMLHttpRequest",
		Description: "Indicates an AJAX request",
	},
}

type HeaderService struct{}

func NewHeaderService() *HeaderService {
	return &HeaderService{}
}

func (s *HeaderService) GetDefaultHeaders() []models.DefaultHeader {
	return DefaultHeaders
}

func (s *HeaderService) GetHeaderByName(name string) *models.DefaultHeader {
	for _, header := range DefaultHeaders {
		if header.Name == name {
			return &header
		}
	}
	return nil
}
