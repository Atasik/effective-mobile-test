package kafka

import (
	"log"
	"os"

	"github.com/IBM/sarama"
)

func InitKafkaConsumerGroup(endpoints []string, groupID string, ClientID string, tlsEnable, returnSuccesses bool, reqAcks sarama.RequiredAcks) (sarama.ConsumerGroup, error) {
	kafkaBrokers := endpoints
	sarama.Logger = log.New(os.Stdout, "[sarama]", log.LstdFlags)
	config := sarama.NewConfig()
	config.ClientID = "fio-app"
	config.Net.TLS.Enable = tlsEnable
	config.Producer.RequiredAcks = reqAcks
	config.Producer.Return.Successes = returnSuccesses

	return sarama.NewConsumerGroup(kafkaBrokers, groupID, config)
}

func InitKafkaSyncProducer(endpoints []string, ClientID string, tlsEnable, returnSuccesses bool, reqAcks sarama.RequiredAcks) (sarama.SyncProducer, error) {
	kafkaBrokers := endpoints
	sarama.Logger = log.New(os.Stdout, "[sarama]", log.LstdFlags)
	config := sarama.NewConfig()
	config.ClientID = "fio-app"
	config.Net.TLS.Enable = tlsEnable
	config.Producer.RequiredAcks = reqAcks
	config.Producer.Return.Successes = returnSuccesses

	return sarama.NewSyncProducer(kafkaBrokers, config)
}
