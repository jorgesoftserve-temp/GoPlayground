package main

import (
	"bytes"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestCreateSession(t *testing.T) {
	sessions = nil

	request := httptest.NewRequest(
		http.MethodPost,
		"/appointments/session",
		bytes.NewBufferString(
			`{
				"appointmentId": 1,
				"diagnosis": "Infeccion leve",
				"notes": "Reposo"
			}`,
		),
	)

	response := httptest.NewRecorder()

	createSession(
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
		sessions,
		1,
	)

	assert.Equal(
		t,
		1,
		sessions[0].AppointmentID,
	)

	assert.Equal(
		t,
		"Infeccion leve",
		sessions[0].Diagnosis,
	)
}

func TestUpdateSession(t *testing.T) {
	sessions = []Session{
		{
			AppointmentID: 1,
			Diagnosis:     "Pendiente",
			Notes:         "",
		},
	}

	request := httptest.NewRequest(
		http.MethodPut,
		"/appointments/session",
		bytes.NewBufferString(
			`{
				"appointmentId": 1,
				"diagnosis": "Gastritis",
				"notes": "Cambiar dieta"
			}`,
		),
	)

	response := httptest.NewRecorder()

	updateSession(
		response,
		request,
	)

	require.Equal(
		t,
		http.StatusOK,
		response.Code,
	)

	require.Len(
		t,
		sessions,
		1,
	)

	assert.Equal(
		t,
		"Gastritis",
		sessions[0].Diagnosis,
	)

	assert.Equal(
		t,
		"Cambiar dieta",
		sessions[0].Notes,
	)
}

func TestUpdateSessionNotFound(t *testing.T) {
	sessions = nil

	request := httptest.NewRequest(
		http.MethodPut,
		"/appointments/session",
		bytes.NewBufferString(
			`{
				"appointmentId": 99,
				"diagnosis": "Nada"
			}`,
		),
	)

	response := httptest.NewRecorder()

	updateSession(
		response,
		request,
	)

	assert.Equal(
		t,
		http.StatusNotFound,
		response.Code,
	)
}
