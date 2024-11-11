package service

import (
	"encoding/json"
	"fmt"
	"github.com/GrzegorzManiak/NoiseBackend/config/structs"
)

var Prefix = "/service"

type Announcement struct {
	Id      string                `json:"id"`
	Port    string                `json:"port"`
	Host    string                `json:"host"`
	Service structs.ServiceConfig `json:"service"`
}

func (s Announcement) String() string {
	return fmt.Sprintf("ServiceAnnouncement{Id: %s, Port: %s, Host: %s, Service: %s}", s.Id, s.Port, s.Host, s.Service.Prefix)
}

func (s Announcement) Marshal() (string, error) {
	s.Service.Password = ""
	data, err := json.Marshal(s)
	return string(data), err
}

func (s Announcement) BuildID() string {
	return fmt.Sprintf("%s/%s%s", Prefix, s.Id, s.Service.Prefix)
}

func UnmarshalServiceAnnouncement(data []byte) (Announcement, error) {
	var serviceAnnouncement Announcement
	err := json.Unmarshal(data, &serviceAnnouncement)
	return serviceAnnouncement, err
}
