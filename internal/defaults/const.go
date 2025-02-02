package defaults

import (
	"github.com/beetlewar010785/pow-task/internal/adapter"
	"github.com/beetlewar010785/pow-task/internal/domain"
	"time"
)

const ServerPort = "8080"
const LogLevel = adapter.LogLevelDebug
const ChallengeDifficulty = 4
const ChallengeLength = 16
const VerificationTimeout = 10 * time.Second

var WordOfWisdomQuotes = []domain.Quote{
	"Cease to be idle; cease to be unclean; cease to find fault one with another.",
	"A man is saved no faster than he gains knowledge.",
	"Our thoughts determine our actions, our actions determine our habits, our habits determine our character, and our character determines our destiny.",
	"When we put God first, all other things fall into their proper place or drop out of our lives.",
	"If you don’t stand for something, you will fall for anything.",
}
