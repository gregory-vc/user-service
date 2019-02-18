package main

var (
	returningColumns = []string{
		"id",
		"created_at",
		"updated_at",
		"first_name",
		"last_name",
		"email",
		"service",
	}

	returningPrivateColumns = []string{
		"id",
		"created_at",
		"updated_at",
		"first_name",
		"last_name",
		"email",
		"service",
		"`password`",
	}

	updateColumns = []string{
		"first_name=$first_name",
		"last_name=$last_name",
		"email=$email",
		"updated_at=$updated_at",
	}
)
