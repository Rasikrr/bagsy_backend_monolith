package analytics

import "github.com/google/uuid"

// StaffWeekdayInput — сырое количество записей мастера в день недели (Weekday уже в формате ТЗ, 0=Пн).
type StaffWeekdayInput struct {
	EmployeeID uuid.UUID
	Weekday    int
	Count      int
}

// NormalizeStaffWeekday нормирует нагрузку по дням недели внутри каждого мастера (макс мастера → 1).
func NormalizeStaffWeekday(inputs []StaffWeekdayInput) []StaffWeekdayCell {
	maxByEmp := make(map[uuid.UUID]int)
	for _, in := range inputs {
		if in.Count > maxByEmp[in.EmployeeID] {
			maxByEmp[in.EmployeeID] = in.Count
		}
	}

	cells := make([]StaffWeekdayCell, 0, len(inputs))
	for _, in := range inputs {
		var v float64
		if m := maxByEmp[in.EmployeeID]; m > 0 {
			v = round2(float64(in.Count) / float64(m))
		}
		cells = append(cells, StaffWeekdayCell{
			EmployeeID: in.EmployeeID,
			Weekday:    in.Weekday,
			Value:      v,
		})
	}
	return cells
}
