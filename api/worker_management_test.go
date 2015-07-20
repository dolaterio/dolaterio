package main

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestCreateAndFetchWorker(t *testing.T) {
	setup()
	defer clean()

	worker, res := createWorker(t, `{"docker_image":"dolaterio/dummy-worker"}`)
	assert.Equal(t, 201, res.Code)
	assert.NotEmpty(t, worker.ID)
	worker = fetchWorker(t, worker.ID)
	assert.NotEmpty(t, worker.ID)
}
