// This file prevents Go from trying to parse template files as Go modules
// Template files contain template variables like {{.Config.GoModule}} which are not valid Go
module templates

go 1.21

// This module is intentionally empty and serves only to isolate template files
