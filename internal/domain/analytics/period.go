package analytics

import "time"

// Period — закрытый интервал дат [From..To] (включительно), нормализованный до даты (без времени).
type Period struct {
	From time.Time
	To   time.Time
}

// NewPeriod создаёт период, нормализуя границы до даты. Возвращает ErrInvalidPeriod при From > To.
func NewPeriod(from, to time.Time) (Period, error) {
	f, t := dateOnly(from), dateOnly(to)
	if f.After(t) {
		return Period{}, ErrInvalidPeriod
	}
	return Period{From: f, To: t}, nil
}

// Days возвращает количество дней в периоде включительно.
func (p Period) Days() int {
	return daysBetween(p.From, p.To) + 1
}

// Dates возвращает список всех дат периода по возрастанию.
func (p Period) Dates() []time.Time {
	n := p.Days()
	dates := make([]time.Time, 0, n)
	for i := range n {
		dates = append(dates, p.From.AddDate(0, 0, i))
	}
	return dates
}

// Previous вычисляет период сравнения по правилам ТЗ:
//   - From == To (один день)            → вчера;
//   - MTD (с 1-го числа по конец/сегодня) → тот же промежуток прошлого месяца (с клампом дня);
//   - иначе                              → equal-length: N дней до From.
func (p Period) Previous(now time.Time) Period {
	now = dateOnly(now)

	// Один день → вчера.
	if p.From.Equal(p.To) {
		y := p.From.AddDate(0, 0, -1)
		return Period{From: y, To: y}
	}

	// MTD → тот же промежуток прошлого месяца.
	if p.isMTD(now) {
		prevFrom := p.From.AddDate(0, -1, 0) // From — 1-е число, сдвиг безопасен.
		day := p.To.Day()
		if last := lastDayOfMonth(prevFrom.Year(), prevFrom.Month()); day > last {
			day = last
		}
		prevTo := time.Date(prevFrom.Year(), prevFrom.Month(), day, 0, 0, 0, 0, p.From.Location())
		return Period{From: prevFrom, To: prevTo}
	}

	// Equal-length: N дней до From.
	n := p.Days()
	prevTo := p.From.AddDate(0, 0, -1)
	prevFrom := prevTo.AddDate(0, 0, -(n - 1))
	return Period{From: prevFrom, To: prevTo}
}

// isMTD: From — 1-е число месяца И To — последний день того же месяца ИЛИ сегодня (в том же месяце).
func (p Period) isMTD(now time.Time) bool {
	if p.From.Day() != 1 {
		return false
	}
	sameMonth := p.To.Year() == p.From.Year() && p.To.Month() == p.From.Month()
	if !sameMonth {
		return false
	}
	lastDay := lastDayOfMonth(p.From.Year(), p.From.Month())
	endOfMonth := time.Date(p.From.Year(), p.From.Month(), lastDay, 0, 0, 0, 0, p.From.Location())
	return p.To.Equal(endOfMonth) || p.To.Equal(now)
}

// dateOnly обнуляет время, оставляя только дату.
func dateOnly(t time.Time) time.Time {
	return time.Date(t.Year(), t.Month(), t.Day(), 0, 0, 0, 0, t.Location())
}

// daysBetween возвращает количество полных дней между from и to (to - from).
func daysBetween(from, to time.Time) int {
	return int(dateOnly(to).Sub(dateOnly(from)).Hours() / 24)
}

// lastDayOfMonth возвращает номер последнего дня месяца.
func lastDayOfMonth(year int, month time.Month) int {
	return time.Date(year, month+1, 0, 0, 0, 0, 0, time.UTC).Day()
}
