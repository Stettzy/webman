# Changelog

All notable changes to this project will be documented in this file.

## [1.0.0] - 2025-05-29

### Added
- Initial release of Webman
- Modern, clean interface for API testing
- Support for all HTTP methods (GET, POST, PUT, DELETE, PATCH)
- Request body support for JSON, Form Data, and Raw Text
- URL parameter management
- Headers management with default headers
- Collections for organizing requests
- Response viewer with formatting and headers inspection
- Docker support for both development and production
- Local-first approach with SQLite storage
- Comprehensive test coverage
- Multi-stage Docker builds for optimal performance
- Hot reloading in development mode
- Vue.js frontend with Tailwind CSS
- Go backend with Gin framework
- Complete documentation

### Technical Details
- Frontend: Vue.js 3.5.13 with TypeScript
- Backend: Go 1.21 with Gin framework
- Database: SQLite with GORM
- UI: Tailwind CSS 4.1.7
- Editor: Monaco Editor for request body
- Testing: Go tests and Vue Test Utils
- Container: Multi-stage Docker builds
- Development: Hot reloading for both frontend and backend 