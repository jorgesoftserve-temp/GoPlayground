package main

import (
	"bytes"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestPetsHandler(t *testing.T) {
	tests := []struct {
		name           string
		method         string
		body           string
		expectedStatus int
	}{
		{
			name:           "crea una mascota",
			method:         http.MethodPost,
			body:           `{"name":"Max","ownerName":"Jorge"}`,
			expectedStatus: http.StatusCreated,
		},
		{
			name:           "rechaza JSON invalido",
			method:         http.MethodPost,
			body:           `{`,
			expectedStatus: http.StatusBadRequest,
		},
		{
			name:           "rechaza metodo no permitido",
			method:         http.MethodPatch,
			expectedStatus: http.StatusMethodNotAllowed,
		},
	}

	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			pets = nil
			nextPetID = 1

			request := httptest.NewRequest(
				test.method,
				"/pets",
				bytes.NewBufferString(test.body),
			)

			response := httptest.NewRecorder()

			petsHandler(
				response,
				request,
			)

			assert.Equal(
				t,
				test.expectedStatus,
				response.Code,
			)
		})
	}
}

func TestCreatePet(t *testing.T) {
	pets = nil
	nextPetID = 1

	request := httptest.NewRequest(
		http.MethodPost,
		"/pets",
		bytes.NewBufferString(
			`{"name":"Max","ownerName":"Jorge"}`,
		),
	)

	response := httptest.NewRecorder()

	createPet(
		response,
		request,
	)

	require.Equal(
		t,
		http.StatusCreated,
		response.Code,
	)

	require.Len(
		t,
		pets,
		1,
	)

	assert.Equal(
		t,
		1,
		pets[0].ID,
	)

	assert.Equal(
		t,
		"Max",
		pets[0].Name,
	)

	assert.Equal(
		t,
		"Jorge",
		pets[0].OwnerName,
	)
}

func TestDeletePet(t *testing.T) {
	pets = []Pet{
		{
			ID:        1,
			Name:      "Max",
			OwnerName: "Jorge",
		},
	}

	request := httptest.NewRequest(
		http.MethodDelete,
		"/pets",
		bytes.NewBufferString(
			`{"id":1}`,
		),
	)

	response := httptest.NewRecorder()

	deletePet(
		response,
		request,
	)

	require.Equal(
		t,
		http.StatusNoContent,
		response.Code,
	)

	assert.Len(
		t,
		pets,
		0,
	)
}
