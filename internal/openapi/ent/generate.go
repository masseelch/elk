package ent

// Generate ent code and swagger file.
//go:generate go run -mod=mod entc.go
// Use the swagger file to create a go client for our server.
//go:generate bash -c "docker run --rm -u \"$(id -u):$(id -g)\" -v ${PWD}/../:/local swaggerapi/swagger-codegen-cli-v3 generate -i /local/ent/openapi.json -l go -o /local/client"
