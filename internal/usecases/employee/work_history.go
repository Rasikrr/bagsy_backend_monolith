package employee

import (
	"context"
	"time"

	"github.com/Rasikrr/bagsy_backend_monolith/internal/domain/identity"
	"github.com/cockroachdb/errors"
)

// saveEmployeeWithWorkHistory сохраняет сотрудника, закрывает текущую запись истории
// и создаёт новую с указанным changeType в рамках одной транзакции.
func (u *UseCase) saveEmployeeWithWorkHistory(ctx context.Context, emp *identity.Employee, changeType identity.ChangeType) error {
	return u.txManager.Do(ctx, func(txCtx context.Context) error {
		if e := u.employeeRepo.Save(txCtx, emp); e != nil {
			return errors.Wrap(e, "save employee")
		}

		currentWH, e := u.workHistoryRepo.GetActiveByEmployeeID(txCtx, emp.ID)
		if e != nil {
			return errors.Wrap(e, "get active work history")
		}

		if currentWH != nil {
			currentWH.End(time.Now())
			if e = u.workHistoryRepo.Save(txCtx, currentWH); e != nil {
				return errors.Wrap(e, "close work history")
			}
		}

		newWH := identity.NewWorkHistory(
			emp.ID,
			emp.OrganizationID,
			emp.LocationID,
			emp.Role,
			changeType,
			nil,
		)
		if e = u.workHistoryRepo.Save(txCtx, newWH); e != nil {
			return errors.Wrap(e, "save new work history")
		}

		return nil
	})
}
