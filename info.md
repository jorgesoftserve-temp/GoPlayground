# Pruebas manuales de la API Veterinaria con cURL

## 1. Levantar la API

``` bash
go run .
```

------------------------------------------------------------------------

## 2. Crear una mascota

``` bash
curl -X POST http://localhost:8080/pets \
  -H "Content-Type: application/json" \
  -d '{
    "name": "Max",
    "ownerName": "Jorge"
  }'
```

Respuesta:

``` json
{"id":1,"name":"Max","ownerName":"Jorge"}
```

------------------------------------------------------------------------

## 3. Consultar mascotas

``` bash
curl http://localhost:8080/pets
```

Respuesta:

``` json
[
  {
    "id": 1,
    "name": "Max",
    "ownerName": "Jorge"
  }
]
```

------------------------------------------------------------------------

## 4. Actualizar una mascota

``` bash
curl -X PUT http://localhost:8080/pets \
  -H "Content-Type: application/json" \
  -d '{
    "id": 1,
    "name": "Maximo",
    "ownerName": "Jorge"
  }'
```

Consultar:

``` bash
curl http://localhost:8080/pets
```

------------------------------------------------------------------------

## 5. Crear una cita

``` bash
curl -X POST http://localhost:8080/appointments \
  -H "Content-Type: application/json" \
  -d '{
    "petId": 1,
    "reason": "Consulta general"
  }'
```

Respuesta:

``` json
{"id":1,"petId":1,"reason":"Consulta general"}
```

Consultar:

``` bash
curl http://localhost:8080/appointments
```

------------------------------------------------------------------------

## 6. Actualizar una cita

``` bash
curl -X PUT http://localhost:8080/appointments \
  -H "Content-Type: application/json" \
  -d '{
    "id": 1,
    "petId": 1,
    "reason": "Vacunacion"
  }'
```

------------------------------------------------------------------------

## 7. Crear una sesión activa de la cita

``` bash
curl -X POST http://localhost:8080/appointments/session \
  -H "Content-Type: application/json" \
  -d '{
    "appointmentId": 1,
    "diagnosis": "Infeccion leve",
    "notes": "Reposo y seguimiento"
  }'
```

Consultar:

``` bash
curl http://localhost:8080/appointments/session
```

------------------------------------------------------------------------

## 8. Actualizar la sesión

``` bash
curl -X PUT http://localhost:8080/appointments/session \
  -H "Content-Type: application/json" \
  -d '{
    "appointmentId": 1,
    "diagnosis": "Paciente estable",
    "notes": "Continuar tratamiento por 5 dias"
  }'
```

------------------------------------------------------------------------

## 9. Eliminar una cita

``` bash
curl -i -X DELETE http://localhost:8080/appointments \
  -H "Content-Type: application/json" \
  -d '{"id":1}'
```

------------------------------------------------------------------------

## 10. Eliminar una mascota

``` bash
curl -i -X DELETE http://localhost:8080/pets \
  -H "Content-Type: application/json" \
  -d '{"id":1}'
```

Consultar:

``` bash
curl http://localhost:8080/pets
curl http://localhost:8080/appointments
```

------------------------------------------------------------------------
