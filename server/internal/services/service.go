package services

import (
	"encoding/json"
	"fmt"
	"github.com/GrzegorzManiak/NoiseBackend/config/structs"
)

type ServiceAnnouncement struct {
	Id      string                `json:"id"`
	Port    string                `json:"port"`
	Host    string                `json:"host"`
	Service structs.ServiceConfig `json:"service"`
}

func (s ServiceAnnouncement) String() string {
	return fmt.Sprintf("ServiceAnnouncement{Id: %s, Port: %d, Host: %s, Service: %s}", s.Id, s.Port, s.Host, s.Service.Prefix)
}

func (s ServiceAnnouncement) Marshal() (string, error) {
	data, err := json.Marshal(s)
	return string(data), err
}

func (s ServiceAnnouncement) BuildID() string {
	return fmt.Sprintf("%s%s", s.Id, s.Service.Prefix)
}

func UnmarshalServiceAnnouncement(data []byte) (ServiceAnnouncement, error) {
	var s ServiceAnnouncement
	err := json.Unmarshal(data, &s)
	return s, err
}
