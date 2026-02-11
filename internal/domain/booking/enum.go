package booking

type Status string

const (
	StatusPending    Status = "pending"
	StatusConfirmed  Status = "confirmed"
	StatusInProgress Status = "in_progress"
	StatusCompleted  Status = "completed"
	StatusCancelled  Status = "cancelled"
)

func (s Status) IsValid() bool {
	switch s {
	case StatusPending, StatusConfirmed, StatusInProgress, StatusCompleted, StatusCancelled:
		return true
	}
	return false
}

func (s Status) String() string {
	return string(s)
}

func (s Status) IsFinal() bool {
	return s == StatusCompleted || s == StatusCancelled
}

func (s Status) CanTransitionTo(target Status) bool {
	transitions := map[Status][]Status{
		StatusPending:    {StatusConfirmed, StatusCancelled},
		StatusConfirmed:  {StatusInProgress, StatusCancelled},
		StatusInProgress: {StatusCompleted, StatusCancelled},
		StatusCompleted:  {},
		StatusCancelled:  {},
	}

	allowed, ok := transitions[s]
	if !ok {
		return false
	}

	for _, t := range allowed {
		if t == target {
			return true
		}
	}
	return false
}

func ParseStatus(s string) (Status, error) {
	status := Status(s)
	if !status.IsValid() {
		return "", ErrAppointmentInvalidStatus
	}
	return status, nil
}
