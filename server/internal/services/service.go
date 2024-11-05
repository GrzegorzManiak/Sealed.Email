package services

import (
	"encoding/json"
	"fmt"
	"github.com/GrzegorzManiak/NoiseBackend/config/structs"
)

var ServicePrefix = "/service"

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
	s.Service.Password = ""
	data, err := json.Marshal(s)
	return string(data), err
}

func (s ServiceAnnouncement) BuildID() string {
	return fmt.Sprintf("%s/%s%s", ServicePrefix, s.Id, s.Service.Prefix)
}

func UnmarshalServiceAnnouncement(data string) (ServiceAnnouncement, error) {
	var serviceAnnouncement ServiceAnnouncement
	err := json.Unmarshal([]byte(data), &serviceAnnouncement)
	return serviceAnnouncement, err
}
