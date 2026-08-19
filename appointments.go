package main

import (
	"encoding/json"
	"net/http"
)

func appointmentsHandler(
	w http.ResponseWriter,
	r *http.Request,
) {
	switch r.Method {

	case http.MethodGet:
		getAppointments(w)

	case http.MethodPost:
		createAppointment(w, r)

	case http.MethodPut:
		updateAppointment(w, r)

	case http.MethodDelete:
		deleteAppointment(w, r)

	default:
		http.Error(
			w,
			"Método no permitido",
			http.StatusMethodNotAllowed,
		)
	}
}

func getAppointments(
	w http.ResponseWriter,
) {
	w.Header().Set(
		"Content-Type",
		"application/json",
	)

	json.NewEncoder(w).Encode(
		appointments,
	)
}

func createAppointment(
	w http.ResponseWriter,
	r *http.Request,
) {
	var appointment Appointment

	err := json.NewDecoder(
		r.Body,
	).Decode(&appointment)

	if err != nil {
		http.Error(
			w,
			"JSON inválido",
			http.StatusBadRequest,
		)

		return
	}

	appointment.ID = nextAppointmentID
	nextAppointmentID++

	appointments = append(
		appointments,
		appointment,
	)

	w.Header().Set(
		"Content-Type",
		"application/json",
	)

	w.WriteHeader(
		http.StatusCreated,
	)

	json.NewEncoder(w).Encode(
		appointment,
	)
}

func updateAppointment(
	w http.ResponseWriter,
	r *http.Request,
) {
	var updatedAppointment Appointment

	err := json.NewDecoder(
		r.Body,
	).Decode(&updatedAppointment)

	if err != nil {
		http.Error(
			w,
			"JSON inválido",
			http.StatusBadRequest,
		)

		return
	}

	for index, appointment := range appointments {
		if appointment.ID == updatedAppointment.ID {
			appointments[index] = updatedAppointment

			json.NewEncoder(w).Encode(
				updatedAppointment,
			)

			return
		}
	}

	http.Error(
		w,
		"Cita no encontrada",
		http.StatusNotFound,
	)
}

func deleteAppointment(
	w http.ResponseWriter,
	r *http.Request,
) {
	var appointmentToDelete Appointment

	err := json.NewDecoder(
		r.Body,
	).Decode(&appointmentToDelete)

	if err != nil {
		http.Error(
			w,
			"JSON inválido",
			http.StatusBadRequest,
		)

		return
	}

	for index, appointment := range appointments {
		if appointment.ID == appointmentToDelete.ID {
			appointments = append(
				appointments[:index],
				appointments[index+1:]...,
			)

			w.WriteHeader(
				http.StatusNoContent,
			)

			return
		}
	}

	http.Error(
		w,
		"Cita no encontrada",
		http.StatusNotFound,
	)
}
