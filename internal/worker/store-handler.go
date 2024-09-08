package worker

import (
	"fmt"
	"github.com/ThreeDotsLabs/watermill/message"
	pkgDynamoDb "github.com/isaquesb/meli-url-shortener/pkg/dynamoDb"
	"github.com/isaquesb/meli-url-shortener/pkg/logger"
)

func StoreDatabase(m *message.Message) error {

	strPayload := string(m.Payload)
	short := strPayload[0:6]
	url := strPayload[6:]

	logger.Debug(fmt.Sprintf("short: %s, url: %s", short, url))

	pkgDynamoDb.Write(short, url)

	return nil
}
