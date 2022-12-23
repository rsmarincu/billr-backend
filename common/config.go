package common

import (
	"os"
)

type Config interface {
	GetEnvVariable(envVarName, defaultValue string) string
	Environment() string
	ServiceName() string
}

func NewConfig(serviceName string) Config {
	env := os.Getenv("env")
	if env == "" {
		env = "development"
	}
	return &config{
		serviceName: serviceName,
		environment: env,
	}
}

type config struct {
	serviceName string
	environment string
}

func (c *config) ServiceName() string {
	return c.serviceName
}

func (c *config) Environment() string {
	return c.environment
}

func (c *config) GetEnvVariable(envVarName, defaultValue string) string {
	envVar := os.Getenv(envVarName)
	if envVar == "" {
		return defaultValue
	}
	return envVar
}
