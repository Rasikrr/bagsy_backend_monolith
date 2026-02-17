package auth

type VerifyRegistrationUseCase struct {
	pending       pendingRegistrationStore
	employees     employeeRepository
	organizations organizationRepository
	plans         planRepository
	subscriptions subscriptionRepository
	workHistory   workHistoryRepository
	tokens        tokenService
	tx            txManager
}

func NewVerifyRegistrationUseCase(
	pending pendingRegistrationStore,
	employees employeeRepository,
	organizations organizationRepository,
	plans planRepository,
	subscriptions subscriptionRepository,
	workHistory workHistoryRepository,
	tokens tokenService,
	tx txManager,
) *VerifyRegistrationUseCase {
	return &VerifyRegistrationUseCase{
		pending:       pending,
		employees:     employees,
		organizations: organizations,
		plans:         plans,
		subscriptions: subscriptions,
		workHistory:   workHistory,
		tokens:        tokens,
		tx:            tx,
	}
}
