// entities - пакет с сущностями для БД
package entities

import (
	"math"
	"math/rand"
	"pora/internal/errors"
)

// MemberColor обозначает цвет пользователя в конкретной семье
type MemberColor string

const (
	ColorRed        MemberColor = "0xFFE57373"
	ColorPink       MemberColor = "0xFFF06292"
	ColorPurple     MemberColor = "0xFFBA68C8"
	ColorViolet     MemberColor = "0xFF9575CD"
	ColorBlue       MemberColor = "0xFF64B5F6"
	ColorCyan       MemberColor = "0xFF4DD0E1"
	ColorGreen      MemberColor = "0xFF81C784"
	ColorLightGreen MemberColor = "0xFFAED581"
	ColorYellow     MemberColor = "0xFFFFD54F"
	ColorOrange     MemberColor = "0xFFFFB74D"
	ColorCoral      MemberColor = "0xFFFF8A65"
	ColorBrown      MemberColor = "0xFFA1887F"
	ColorGray       MemberColor = "0xFF90A4AE"
	ColorTurquoise  MemberColor = "0xFF4DB6AC"
	ColorIndigo     MemberColor = "0xFF7986CB"
	ColorLightBlue  MemberColor = "0xFF4FC3F7"
	ColorOrchid     MemberColor = "0xFFCE93D8"
)

// AvailableColors - массив со всеми возможными
// цветами для членов семьи
var AvailableColors = []MemberColor{
	ColorRed,
	ColorPink,
	ColorPurple,
	ColorViolet,
	ColorBlue,
	ColorCyan,
	ColorGreen,
	ColorLightGreen,
	ColorYellow,
	ColorOrange,
	ColorCoral,
	ColorBrown,
	ColorGray,
	ColorTurquoise,
	ColorIndigo,
	ColorLightBlue,
	ColorOrchid,
}

// GetFreeColor получает случайный цвет с наименьшим
// количеством повторений внутри семьи
func (f *Family) GetFreeColor() (MemberColor, error) {
	counts := make(map[MemberColor]int)

	for _, color := range AvailableColors {
		counts[color] = 0
	}

	for _, member := range f.Members {
		counts[member.Color]++
	}

	minCount := math.MaxInt

	for _, count := range counts {
		if count < minCount {
			minCount = count
		}
	}

	var freeColors []MemberColor

	for color, count := range counts {
		if count == minCount {
			freeColors = append(freeColors, color)
		}
	}

	if len(freeColors) == 0 {
		return "", errors.ErrorInternal
	}

	randomIndex := rand.Intn(len(freeColors))
	return freeColors[randomIndex], nil
}
