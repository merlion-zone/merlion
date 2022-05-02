package types

const (
	SecondsPerMinute = 60
	SecondsPerHour   = 60 * SecondsPerMinute
	SecondsPerDay    = 24 * SecondsPerHour
	SecondsPerWeek   = 7 * SecondsPerDay
	DaysPer4Years    = 365*4 + 1
	SecondsPer4Years = DaysPer4Years * SecondsPerDay
)
