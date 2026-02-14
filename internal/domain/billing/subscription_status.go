package billing

type SubscriptionStatus string

const (
	SubscriptionStatusTrial     SubscriptionStatus = "trial"
	SubscriptionStatusActive    SubscriptionStatus = "active"
	SubscriptionStatusPastDue   SubscriptionStatus = "past_due"
	SubscriptionStatusSuspended SubscriptionStatus = "suspended"
	SubscriptionStatusCanceled  SubscriptionStatus = "canceled"
)

func (s SubscriptionStatus) IsValid() bool {
	switch s {
	case SubscriptionStatusTrial,
		SubscriptionStatusActive,
		SubscriptionStatusPastDue,
		SubscriptionStatusSuspended,
		SubscriptionStatusCanceled:
		return true
	}
	return false
}

func (s SubscriptionStatus) String() string {
	return string(s)
}

// CanTransitionTo проверяет допустимость перехода между статусами.
//
//	trial     → active, past_due
//	active    → active (продление), past_due
//	past_due  → active, suspended
//	suspended → active, canceled
//	canceled  → (финальный)
func (s SubscriptionStatus) CanTransitionTo(target SubscriptionStatus) bool {
	transitions := map[SubscriptionStatus][]SubscriptionStatus{
		SubscriptionStatusTrial:     {SubscriptionStatusActive, SubscriptionStatusPastDue},
		SubscriptionStatusActive:    {SubscriptionStatusActive, SubscriptionStatusPastDue},
		SubscriptionStatusPastDue:   {SubscriptionStatusActive, SubscriptionStatusSuspended},
		SubscriptionStatusSuspended: {SubscriptionStatusActive, SubscriptionStatusCanceled},
		SubscriptionStatusCanceled:  {},
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

// IsFinal возвращает true если статус финальный (canceled).
func (s SubscriptionStatus) IsFinal() bool {
	return s == SubscriptionStatusCanceled
}

// CanOperate возвращает true если организация может полноценно работать.
// trial, active, past_due — полный функционал.
func (s SubscriptionStatus) CanOperate() bool {
	switch s {
	case SubscriptionStatusTrial, SubscriptionStatusActive, SubscriptionStatusPastDue:
		return true
	}
	return false
}

// CanRead возвращает true если допустим read-only доступ.
// suspended — можно просматривать данные, но не создавать новые.
func (s SubscriptionStatus) CanRead() bool {
	switch s {
	case SubscriptionStatusTrial, SubscriptionStatusActive,
		SubscriptionStatusPastDue, SubscriptionStatusSuspended:
		return true
	}
	return false
}

func ParseSubscriptionStatus(s string) (SubscriptionStatus, error) {
	status := SubscriptionStatus(s)
	if !status.IsValid() {
		return "", ErrInvalidSubscriptionStatus
	}
	return status, nil
}
