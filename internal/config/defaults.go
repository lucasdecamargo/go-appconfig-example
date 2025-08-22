package config

import (
	"log"
	"strconv"
	"time"
)

func defaultString(val string) any {
	if val == "" {
		return nil
	}

	return val
}

func defaultBool(val string) any {
	switch val {
	case "":
		return nil
	case "true":
		return true
	case "false":
		return false
	default:
		log.Fatalf("Invalid default bool: %s", val)
	}

	return false
}

func defaultInt(val string) any {
	if val == "" {
		return nil
	}

	i, err := strconv.Atoi(val)
	if err != nil {
		log.Fatalf("Invalid default int: %s", val)
	}

	return i
}

func defaultDuration(val string) any {
	if val == "" {
		return nil
	}

	d, err := time.ParseDuration(val)
	if err != nil {
		log.Fatalf("Invalid default duration: %s", val)
	}

	return d
}
