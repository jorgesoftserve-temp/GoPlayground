package main

import (
	"encoding/json"
	"net/http"
)

func petsHandler(
	w http.ResponseWriter,
	r *http.Request,
) {
	switch r.Method {

	case http.MethodGet:
		getPets(w)

	case http.MethodPost:
		createPet(w, r)

	case http.MethodPut:
		updatePet(w, r)

	case http.MethodDelete:
		deletePet(w, r)

	default:
		http.Error(
			w,
			"Método no permitido",
			http.StatusMethodNotAllowed,
		)
	}
}

func getPets(
	w http.ResponseWriter,
) {
	w.Header().Set(
		"Content-Type",
		"application/json",
	)

	json.NewEncoder(w).Encode(pets)
}

func createPet(
	w http.ResponseWriter,
	r *http.Request,
) {
	var pet Pet

	err := json.NewDecoder(
		r.Body,
	).Decode(&pet)

	if err != nil {
		http.Error(
			w,
			"JSON inválido",
			http.StatusBadRequest,
		)

		return
	}

	pet.ID = nextPetID
	nextPetID++

	pets = append(
		pets,
		pet,
	)

	w.Header().Set(
		"Content-Type",
		"application/json",
	)

	w.WriteHeader(
		http.StatusCreated,
	)

	json.NewEncoder(w).Encode(pet)
}

func updatePet(
	w http.ResponseWriter,
	r *http.Request,
) {
	var updatedPet Pet

	err := json.NewDecoder(
		r.Body,
	).Decode(&updatedPet)

	if err != nil {
		http.Error(
			w,
			"JSON inválido",
			http.StatusBadRequest,
		)

		return
	}

	for index, pet := range pets {
		if pet.ID == updatedPet.ID {
			pets[index] = updatedPet

			json.NewEncoder(w).Encode(
				updatedPet,
			)

			return
		}
	}

	http.Error(
		w,
		"Mascota no encontrada",
		http.StatusNotFound,
	)
}

func deletePet(
	w http.ResponseWriter,
	r *http.Request,
) {
	var petToDelete Pet

	err := json.NewDecoder(
		r.Body,
	).Decode(&petToDelete)

	if err != nil {
		http.Error(
			w,
			"JSON inválido",
			http.StatusBadRequest,
		)

		return
	}

	for index, pet := range pets {
		if pet.ID == petToDelete.ID {
			pets = append(
				pets[:index],
				pets[index+1:]...,
			)

			w.WriteHeader(
				http.StatusNoContent,
			)

			return
		}
	}

	http.Error(
		w,
		"Mascota no encontrada",
		http.StatusNotFound,
	)
}
