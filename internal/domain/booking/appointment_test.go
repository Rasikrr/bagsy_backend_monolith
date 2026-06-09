package booking

import (
	"testing"

	"github.com/google/uuid"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func appointmentWithStatus(status Status) *Appointment {
	return &Appointment{
		ID:            uuid.New(),
		Status:        status,
		StatusHistory: []StatusHistoryEntry{{ID: uuid.New(), ToStatus: status}},
	}
}

func TestAppointment_AutoComplete(t *testing.T) {
	t.Run("from confirmed → completed", func(t *testing.T) {
		a := appointmentWithStatus(StatusConfirmed)

		err := a.AutoComplete()

		require.NoError(t, err)
		assert.Equal(t, StatusCompleted, a.Status)
		assert.NotNil(t, a.UpdatedAt)

		last := a.StatusHistory[len(a.StatusHistory)-1]
		assert.Equal(t, StatusCompleted, last.ToStatus)
		require.NotNil(t, last.FromStatus)
		assert.Equal(t, StatusConfirmed, *last.FromStatus)
		assert.Nil(t, last.ChangedBy) // системное действие
		require.NotNil(t, last.Reason)
		assert.Equal(t, reasonAutoCompleted, *last.Reason)
	})

	t.Run("from in_progress → completed", func(t *testing.T) {
		a := appointmentWithStatus(StatusInProgress)

		err := a.AutoComplete()

		require.NoError(t, err)
		assert.Equal(t, StatusCompleted, a.Status)
	})

	t.Run("invalid from non-active statuses", func(t *testing.T) {
		for _, status := range []Status{StatusPending, StatusCompleted, StatusCancelled} {
			a := appointmentWithStatus(status)

			err := a.AutoComplete()

			require.ErrorIs(t, err, ErrAppointmentInvalidStatusTransition)
			assert.Equal(t, status, a.Status, "status must not change on invalid transition")
		}
	})
}
