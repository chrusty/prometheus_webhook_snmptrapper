package config

type Config struct {
	SNMPTrapAddress string
	SNMPCommunity   string
	SNMPRetries     uint
	WebhookAddress  string
}
