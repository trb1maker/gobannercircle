// Одно из условий задачи: расчеты должны выполнять в отдельном модуле.
// Нет проблем реализовать расчеты в Postgres, но SQLite таких функций не
// поддерживает. Чтобы решить эти проблемы я вынес расчеты в отдельный модуль.
// Структура Stat необходимо для передачи данных.

package storage

// Stat описывает как часто был показан баннер в конкретном слоте
// для конкретной группы пользователей.
type Stat struct {
	ID     int
	Views  int
	Clicks int
	P      float64
}

type Stats []Stat

func (s Stats) Less(i, j int) bool {
	return s[i].P < s[j].P
}

func (s Stats) Len() int {
	return len(s)
}

func (s Stats) Swap(i, j int) {
	s[i], s[j] = s[j], s[i]
}
