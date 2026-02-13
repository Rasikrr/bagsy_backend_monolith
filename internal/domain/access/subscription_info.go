package access

// ─────────────────────────────────────────────────────────────────
// SubscriptionInfo
//
// Все поля pre-computed в middleware из billing.Subscription.
// Логика определения статуса остаётся в billing домене,
// здесь — только результат.
// ─────────────────────────────────────────────────────────────────

type SubscriptionInfo struct {
	Active bool
}
