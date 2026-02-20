package auth

type ActionTokenPurpose string

const (
	ActionTokenPurposeUnknown         ActionTokenPurpose = "unknown"
	ActionTokenPurposeStaffInvitation ActionTokenPurpose = "staff_invitation"
	ActionTokenPurposePasswordReset   ActionTokenPurpose = "password_reset"
)

func (t ActionTokenPurpose) String() string {
	return string(t)
}

func ParseActionTokenPurpose(s string) (ActionTokenPurpose, error) {
	p := ActionTokenPurpose(s)
	switch p {
	case ActionTokenPurposeStaffInvitation, ActionTokenPurposePasswordReset:
		return p, nil
	default:
		return ActionTokenPurposeUnknown, ErrUnknownTokenPurpose
	}
}
