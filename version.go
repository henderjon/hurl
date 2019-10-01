package main

// these vars are built at compile time, DO NOT ALTER
var (
	// Version adds build information
	buildVersion string
	// BuildTimestamp adds build information
	buildTimestamp string
	// CompiledBy adds the make/model that was used to compile
	compiledBy string
)

// GetBuildVersion returns the build version set when the binary was built
func getBuildVersion() string {
	return buildVersion
}

// GetBuildTimestamp returns the build timestamp set when the binary was built
func getBuildTimestamp() string {
	return buildTimestamp
}

// GetBuildTimestamp returns the build timestamp set when the binary was built
func getCompiledBy() string {
	return compiledBy
}
