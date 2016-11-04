package config

type Config struct {
	SNMPTrapBinary  string
	SNMPTrapAddress string
	SNMPCommunity   string
	SNMPRetries     uint
	WebhookAddress  string
}
