package models

type QueueEmailId struct {
	EmailId string
}

func (ie QueueEmailId) Marshal() (string, error) {
	return ie.EmailId, nil
}

func UnmarshalQueueEmailId(data string) (QueueEmailId, error) {
	return QueueEmailId{EmailId: data}, nil
}
