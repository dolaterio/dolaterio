package main

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestCreateAndFetchWorker(t *testing.T) {
	setup()
	defer clean()

	worker := createWorker(t, `{"docker_image":"dolaterio/dummy-worker"}`)
	assert.NotEmpty(t, worker.ID)
	worker = fetchWorker(t, worker.ID)
	assert.NotEmpty(t, worker.ID)
}
