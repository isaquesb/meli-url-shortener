package config

func getAmqpURI() string {
	host := GetEnv("RABBITMQ_HOST", "rabbitmq:5672")
	user := GetEnv("RABBITMQ_USER", "guest")
	password := GetEnv("RABBITMQ_PASS", "guest")
	return "amqp://" + user + ":" + password + "@" + host
}
