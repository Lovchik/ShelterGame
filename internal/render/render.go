package render

import (
	"ShelterGame/internal/database/sqlite"
	"fmt"
	"log/slog"
	"math/rand"
	"os"
	"strings"
	"time"
)

var order = map[string]Pair{
	"-1-": {
		tableName: "accident",
		rollKind:  "roll",
	},
	"-2-": {
		tableName: "shelter_info",
		rollKind:  "roll",
	},
	"-3-": {
		tableName: "profession",
		rollKind:  "roll",
	},
	"-4-": {
		tableName: "health",
		rollKind:  "rollWithChance",
	},
	"-5-": {
		tableName: "",
		rollKind:  "hardcode",
	},
	"-6-": {
		tableName: "hobby",
		rollKind:  "roll",
	},
	"-7-": {
		tableName: "trait",
		rollKind:  "roll",
	},
	"-8-": {
		tableName: "phobia",
		rollKind:  "rollWithChance",
	},
	"-9-": {
		tableName: "add_info",
		rollKind:  "roll",
	},
	"-10-": {
		tableName: "baggage",
		rollKind:  "roll",
	},
}

type Pair struct {
	tableName string
	rollKind  string
}

func Render(memberNumber string) string {
	file, err := os.ReadFile("/Users/vitali.louchy/Desktop/ShelterGame/sample")
	if err != nil {
		panic(err)
	}
	info := string(file)
	info = strings.Replace(info, "-0-", memberNumber, -1)

	for k, v := range order {
		info = strings.Replace(info, k, decide(v), -1)
	}
	slog.Info("info:", info)
	return info
}

func rollWithChance(tableName string, firstChance int) string {
	rand.Seed(time.Now().UnixNano())
	randomValue := rand.Intn(100) + 1
	var result string
	if randomValue <= firstChance {
		query := fmt.Sprintf("SELECT name FROM %s where id=1", tableName)
		sqlite.GetDB().Raw(query).Scan(&result)
	} else {
		return roll(tableName)
	}
	return result
}

func roll(tableName string) string {
	var count int64
	var result string
	query := fmt.Sprintf("SELECT Count(id) FROM %s", tableName)
	sqlite.GetDB().Raw(query).Debug().Count(&count)
	randomValue := rand.Intn(int(count)) + 1
	query = fmt.Sprintf("SELECT name FROM %s where id=?", tableName)
	sqlite.GetDB().Raw(query, randomValue).Debug().Scan(&result)
	return result
}

func decide(pair Pair) string {
	switch pair.rollKind {
	case "roll":
		{
			return roll(pair.tableName)
		}
	case "rollWithChance":
		{
			return rollWithChance(pair.tableName, 30)
		}
	case "hardcode":
		{
			return "Не готово еще"
		}

	}
	return ""
}
