package main

import (
	"bytes"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestCreateAppointment(t *testing.T) {
	appointments = nil
	nextAppointmentID = 1

	request := httptest.NewRequest(
		http.MethodPost,
		"/appointments",
		bytes.NewBufferString(
			`{
				"petId": 1,
				"reason": "Consulta general"
			}`,
		),
	)

	response := httptest.NewRecorder()

	createAppointment(
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
		appointments,
		1,
	)

	assert.Equal(
		t,
		1,
		appointments[0].ID,
	)

	assert.Equal(
		t,
		1,
		appointments[0].PetID,
	)

	assert.Equal(
		t,
		"Consulta general",
		appointments[0].Reason,
	)
}

func TestUpdateAppointment(t *testing.T) {
	appointments = []Appointment{
		{
			ID:     1,
			PetID:  1,
			Reason: "Consulta",
		},
	}

	request := httptest.NewRequest(
		http.MethodPut,
		"/appointments",
		bytes.NewBufferString(
			`{
				"id": 1,
				"petId": 1,
				"reason": "Vacunacion"
			}`,
		),
	)

	response := httptest.NewRecorder()

	updateAppointment(
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
		appointments,
		1,
	)

	assert.Equal(
		t,
		"Vacunacion",
		appointments[0].Reason,
	)
}

func TestDeleteAppointmentNotFound(t *testing.T) {
	appointments = nil

	request := httptest.NewRequest(
		http.MethodDelete,
		"/appointments",
		bytes.NewBufferString(
			`{"id":99}`,
		),
	)

	response := httptest.NewRecorder()

	deleteAppointment(
		response,
		request,
	)

	assert.Equal(
		t,
		http.StatusNotFound,
		response.Code,
	)
}
