package utils_test

import (
	"github.com/stretchr/testify/assert"
	"testing"
	"ticken-ticket-service/utils"
)

func Test_IsDev_ReturnsCorrectly(t *testing.T) {
	config := new(utils.TickenConfig)
	config.Env.TickenEnv = "dev"
	assert.True(t, config.IsDev())
	assert.False(t, config.IsTest())
	assert.False(t, config.IsProd())
}

func Test_IsTest_ReturnsCorrectly(t *testing.T) {
	config := new(utils.TickenConfig)
	config.Env.TickenEnv = "test"
	assert.False(t, config.IsDev())
	assert.True(t, config.IsTest())
	assert.False(t, config.IsProd())
}

func Test_IsProd_ReturnsCorrectly(t *testing.T) {
	config := new(utils.TickenConfig)
	config.Env.TickenEnv = "prod"
	assert.False(t, config.IsDev())
	assert.False(t, config.IsTest())
	assert.True(t, config.IsProd())
}
