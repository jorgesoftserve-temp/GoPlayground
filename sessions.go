package main

import (
	"encoding/json"
	"net/http"
)

func sessionHandler(
	w http.ResponseWriter,
	r *http.Request,
) {
	switch r.Method {

	case http.MethodGet:
		getSessions(w)

	case http.MethodPost:
		createSession(w, r)

	case http.MethodPut:
		updateSession(w, r)

	default:
		http.Error(
			w,
			"Método no permitido",
			http.StatusMethodNotAllowed,
		)
	}
}

func getSessions(
	w http.ResponseWriter,
) {
	w.Header().Set(
		"Content-Type",
		"application/json",
	)

	json.NewEncoder(w).Encode(
		sessions,
	)
}

func createSession(
	w http.ResponseWriter,
	r *http.Request,
) {
	var session Session

	err := json.NewDecoder(
		r.Body,
	).Decode(&session)

	if err != nil {
		http.Error(
			w,
			"JSON inválido",
			http.StatusBadRequest,
		)

		return
	}

	sessions = append(
		sessions,
		session,
	)

	w.Header().Set(
		"Content-Type",
		"application/json",
	)

	w.WriteHeader(
		http.StatusCreated,
	)

	json.NewEncoder(w).Encode(
		session,
	)
}

func updateSession(
	w http.ResponseWriter,
	r *http.Request,
) {
	var updatedSession Session

	err := json.NewDecoder(
		r.Body,
	).Decode(&updatedSession)

	if err != nil {
		http.Error(
			w,
			"JSON inválido",
			http.StatusBadRequest,
		)

		return
	}

	for index, session := range sessions {
		if session.AppointmentID ==
			updatedSession.AppointmentID {

			sessions[index] = updatedSession

			json.NewEncoder(w).Encode(
				updatedSession,
			)

			return
		}
	}

	http.Error(
		w,
		"Sesión no encontrada",
		http.StatusNotFound,
	)
}
