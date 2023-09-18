package kafka

import (
	"context"
	"fio/internal/service"

	"github.com/IBM/sarama"
	"github.com/go-playground/validator"
	"go.uber.org/zap"
)

const (
	FioTopic       = "FIO"
	FioFailedTopic = "FIO_FAILED"
)

type MessageHandler struct {
	services  *service.Service
	validator *validator.Validate
	logger    *zap.SugaredLogger

	syncProducer sarama.SyncProducer
}

func NewMessageHander(services *service.Service, validator *validator.Validate, logger *zap.SugaredLogger, syncProducer sarama.SyncProducer) *MessageHandler {
	return &MessageHandler{
		services:     services,
		validator:    validator,
		logger:       logger,
		syncProducer: syncProducer,
	}
}

func (h *MessageHandler) ConsumeLoop(topics []string, ctx context.Context, consumerGroup sarama.ConsumerGroup) {
	for {
		err := consumerGroup.Consume(ctx, topics, h)
		if err != nil {
			h.logger.Panicf("Error from consumer: %v", err)
		}

		if ctx.Err() != nil {
			return
		}
	}
}

func (h *MessageHandler) Setup(_ sarama.ConsumerGroupSession) error { return nil }

func (h *MessageHandler) Cleanup(_ sarama.ConsumerGroupSession) error { return nil }

func (h *MessageHandler) ConsumeClaim(sess sarama.ConsumerGroupSession, claim sarama.ConsumerGroupClaim) error {
	for msg := range claim.Messages() {
		select {
		case <-claim.Messages():
			h.handleConsumerMessages(msg)
		case <-sess.Context().Done():
			return nil
		}
		h.logger.Infof("Message topic:%q partition:%d offset:%d\n", msg.Topic, msg.Partition, msg.Offset)

		sess.MarkMessage(msg, "")
	}
	return nil
}

// TODO: refactor
func (h *MessageHandler) handleConsumerMessages(msg *sarama.ConsumerMessage) {
	if msg.Topic == FioTopic {
		person, err := h.handleAddPersonMessage(msg.Value)
		if err != nil {
			err = h.SendPersonErrorReponseToKafka(personErrorResponse{person, errorResponse{err.Error()}}, FioFailedTopic)
			if err != nil {
				h.logger.Infow("failed to send message to topic", FioFailedTopic)
			}
		}
	}
}
