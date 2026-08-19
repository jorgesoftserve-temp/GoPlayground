package main

type Pet struct {
	ID        int    `json:"id"`
	Name      string `json:"name"`
	OwnerName string `json:"ownerName"`
}

type Appointment struct {
	ID     int    `json:"id"`
	PetID  int    `json:"petId"`
	Reason string `json:"reason"`
}

type Session struct {
	AppointmentID int    `json:"appointmentId"`
	Diagnosis     string `json:"diagnosis"`
	Notes         string `json:"notes"`
}

var pets []Pet

var appointments []Appointment

var sessions []Session

var nextPetID = 1

var nextAppointmentID = 1
