package cli

import (
	"fmt"
	"regexp"
	"strconv"
	"time"
)

func parseDue(d string) (time.Duration, error) {
	durationExpression := regexp.MustCompile(`^(?<hours>\d+h)?(?<days>\d+d)?(?<weeks>\d+w)?(?<months>\d+m)?(?<years>\d+y)?$`)
	var duration time.Duration
	match := durationExpression.FindStringSubmatch(d)
	for i, name := range durationExpression.SubexpNames() {
		if i != 0 && name != "" && len(match[i]) != 0 {
			num, err := strconv.Atoi(match[i][:len(match[i])-1])
			if err != nil {
				return 0, fmt.Errorf("could not parse number: %w", err)
			}
			switch name {
			case "hours":
				duration += time.Duration(num) * time.Hour
			case "days":
				duration += time.Duration(24*num) * time.Hour
			case "weeks":
				duration += time.Duration(7*24*num) * time.Hour
			case "months":
				duration += time.Duration(30*24*num) * time.Hour
			case "years":
				duration += time.Duration(365*24*num) * time.Hour
			}
		}
	}
	return duration, nil
}
