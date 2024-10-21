package services

import (
	"encoding/json"
	"fmt"
)

type ServiceAnnouncement struct {
	Id   string `json:"id"`
	Port int    `json:"port"`
	Host string `json:"host"`
}

func (s ServiceAnnouncement) String() string {
	return fmt.Sprintf("ServiceAnnouncement{Id: %s, Port: %d, Host: %s}", s.Id, s.Port, s.Host)
}

func (s ServiceAnnouncement) Marshal() (string, error) {
	data, err := json.Marshal(s)
	return string(data), err
}

func UnmarshalServiceAnnouncement(data []byte) (ServiceAnnouncement, error) {
	var s ServiceAnnouncement
	err := json.Unmarshal(data, &s)
	return s, err
}
