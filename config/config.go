package config

import (
	"bufio"
	"log"
	"os"
	"strings"
)

type Config struct {
	SMTPHost              string
	SMTPPort              string
	SMTPUser              string
	SMTPPass              string
	SMTPFrom              string
	AppName               string
	WaitlistReceiverEmail string
	Port                  string
}

func Load() Config {

	file, err := os.Open(".env")
	if err == nil {
		defer file.Close()
		scanner := bufio.NewScanner(file)
		for scanner.Scan() {
			line := scanner.Text()
			if strings.TrimSpace(line) == "" || strings.HasPrefix(line, "#") {
				continue
			}
			parts := strings.SplitN(line, "=", 2)
			if len(parts) == 2 {
				key := strings.TrimSpace(parts[0])
				val := strings.TrimSpace(parts[1])
				os.Setenv(key, val)
			}
		}
	}

	cfg := Config{
		SMTPHost:              getEnv("SMTP_HOST", true),
		SMTPPort:              getEnv("SMTP_PORT", true),
		SMTPUser:              getEnv("SMTP_USER", true),
		SMTPPass:              getEnv("SMTP_PASS", true),
		SMTPFrom:              getEnv("SMTP_FROM", true),
		AppName:               getEnv("APP_NAME", true),
		WaitlistReceiverEmail: getEnv("WAITLIST_RECEIVER_EMAIL", true),
		Port:                  getEnv("PORT", false, "8080"),
	}
	return cfg
}

func getEnv(key string, required bool, defaults ...string) string {
	val := os.Getenv(key)
	if val == "" && len(defaults) > 0 {
		return defaults[0]
	}
	if val == "" && required {
		log.Fatalf("Missing required env var: %s", key)
	}
	return val
}
