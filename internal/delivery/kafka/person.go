package kafka

import (
	"encoding/json"
	"fio/internal/domain"

	"github.com/IBM/sarama"
)

func (h *MessageHandler) SendPersonErrorReponseToKafka(person personErrorResponse, topic string) error {
	personBytes, err := json.Marshal(person)
	if err != nil {
		return err
	}

	msg := &sarama.ProducerMessage{
		Topic: topic,
		Value: sarama.ByteEncoder(personBytes),
	}

	_, _, err = h.syncProducer.SendMessage(msg)
	return err
}

func (h *MessageHandler) handleAddPersonMessage(message []byte) (domain.Person, error) {
	var person domain.Person
	err := json.Unmarshal(message, &person)
	if err != nil {
		return domain.Person{}, err
	}

	err = h.validator.Struct(person)
	if err != nil {
		return domain.Person{}, err
	}

	id, err := h.services.AddPerson(person)
	if err != nil {
		return domain.Person{}, err
	}
	h.logger.Infof("Person with id %d was added", id)

	return person, nil
}
