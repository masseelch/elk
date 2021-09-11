package ent

// Generate ent code and swagger file.
//go:generate go run -mod=mod entc.go
// Example of how to use the OAS file to create a go client for our server with the swagger-codegen docker image.
// go:generate bash -c "docker run --rm -u \"$(id -u):$(id -g)\" -v ${PWD}/../:/local swaggerapi/swagger-codegen-cli-v3 generate -i /local/ent/openapi.json -l go -o /local/client"
