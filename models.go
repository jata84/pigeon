package main

type Location struct {
	Status      string    `json:"status"`
	Coordinates []float64 `json:"coordinates"`
	HDOP        float64   `json:"hdop"`
	DeviceID    string    `json:"deviceId"`
	EfiUserID   int       `json:"efiuserId"`
}
