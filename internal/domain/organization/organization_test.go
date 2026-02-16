package organization

import (
	"testing"

	"github.com/google/uuid"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestNewOrganization(t *testing.T) {
	org, err := NewStubOrganization()

	require.NoError(t, err)
	assert.NotEqual(t, uuid.Nil, org.ID)
	assert.True(t, org.Active)
	assert.False(t, org.CreatedAt.IsZero())
}

func TestOrganization_UpdateInfo(t *testing.T) {
	org, _ := NewStubOrganization()

	name := "Test Org"
	desc := "Test Description"
	err := org.SetupProfile(name, &desc)

	require.NoError(t, err)
	assert.Equal(t, name, org.Name)
	require.NotNil(t, org.Description)
	assert.Equal(t, desc, *org.Description)
	assert.False(t, org.Slug.IsEmpty())
	assert.NotNil(t, org.UpdatedAt)

	// Test update on deleted organization
	err = org.Delete()
	require.NoError(t, err)
	err = org.SetupProfile("New Name", nil)
	assert.ErrorIs(t, err, ErrOrganizationDeleted)
}

func TestOrganization_Activation(t *testing.T) {
	org, _ := NewStubOrganization()

	t.Run("deactivate active organization", func(t *testing.T) {
		err := org.Deactivate()
		assert.NoError(t, err)
		assert.False(t, org.Active)
	})

	t.Run("deactivate already inactive organization", func(t *testing.T) {
		err := org.Deactivate()
		assert.NoError(t, err)
	})

	t.Run("activate inactive organization", func(t *testing.T) {
		err := org.Activate()
		assert.NoError(t, err)
		assert.True(t, org.Active)
	})

	t.Run("activate on deleted organization", func(t *testing.T) {
		_ = org.Delete()
		err := org.Activate()
		assert.ErrorIs(t, err, ErrOrganizationDeleted)
	})
}

func TestOrganization_Delete(t *testing.T) {
	org, _ := NewStubOrganization()

	err := org.Delete()
	assert.NoError(t, err)
	assert.True(t, org.IsDeleted())
	assert.False(t, org.Active)
	assert.NotNil(t, org.DeletedAt)

	// Double delete
	err = org.Delete()
	assert.NoError(t, err)
}

func TestOrganization_CanOperate(t *testing.T) {
	org, _ := NewStubOrganization()

	assert.True(t, org.CanOperate())

	_ = org.Deactivate()
	assert.False(t, org.CanOperate())

	_ = org.Activate()
	_ = org.Delete()
	assert.False(t, org.CanOperate())
}
