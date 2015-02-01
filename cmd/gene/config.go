package main

// Config holds the config parameters for gene package
type Config struct {
	// Schema holds the given schema file
	Schema string `required:"true"`

	// Target holds the target folder
	Target string `required:"true"`

	// Generators holds the generator names for processing
	Generators []string
}
