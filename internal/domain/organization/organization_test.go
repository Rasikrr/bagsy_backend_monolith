package organization

import (
	"testing"

	"github.com/google/uuid"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestNewOrganization(t *testing.T) {
	ownerID := uuid.New()
	org, err := NewStubOrganization(ownerID)

	require.NoError(t, err)
	assert.Equal(t, ownerID, org.OwnerID)
	assert.NotEqual(t, uuid.Nil, org.ID)
	assert.True(t, org.Active)
	assert.False(t, org.CreatedAt.IsZero())
}

func TestOrganization_UpdateInfo(t *testing.T) {
	ownerID := uuid.New()
	org, _ := NewStubOrganization(ownerID)

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

func TestOrganization_ChangeOwnership(t *testing.T) {
	ownerID := uuid.New()
	org, _ := NewStubOrganization(ownerID)
	_ = org.SetupProfile("Test Org", nil)

	newOwnerID := uuid.New()

	t.Run("successful ownership change", func(t *testing.T) {
		err := org.ChangeOwnership(newOwnerID)
		assert.NoError(t, err)
		assert.Equal(t, newOwnerID, org.OwnerID)
	})

	t.Run("change to same owner", func(t *testing.T) {
		err := org.ChangeOwnership(newOwnerID)
		assert.ErrorIs(t, err, ErrSameOwner)
	})

	t.Run("change on inactive organization", func(t *testing.T) {
		err := org.Deactivate()
		require.NoError(t, err)
		err = org.ChangeOwnership(uuid.New())
		assert.ErrorIs(t, err, ErrOrganizationInactive)
		_ = org.Activate()
	})

	t.Run("change on deleted organization", func(t *testing.T) {
		err := org.Delete()
		require.NoError(t, err)
		err = org.ChangeOwnership(uuid.New())
		assert.ErrorIs(t, err, ErrOrganizationDeleted)
	})
}

func TestOrganization_Activation(t *testing.T) {
	ownerID := uuid.New()
	org, _ := NewStubOrganization(ownerID)

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
	ownerID := uuid.New()
	org, _ := NewStubOrganization(ownerID)

	err := org.Delete()
	assert.NoError(t, err)
	assert.True(t, org.IsDeleted())
	assert.False(t, org.Active)
	assert.NotNil(t, org.DeletedAt)

	// Double delete
	err = org.Delete()
	assert.NoError(t, err)
}

func TestOrganization_IsOwnedBy(t *testing.T) {
	ownerID := uuid.New()
	org, _ := NewStubOrganization(ownerID)

	assert.True(t, org.IsOwnedBy(ownerID))
	assert.False(t, org.IsOwnedBy(uuid.New()))
}

func TestOrganization_CanOperate(t *testing.T) {
	ownerID := uuid.New()
	org, _ := NewStubOrganization(ownerID)

	assert.True(t, org.CanOperate())

	_ = org.Deactivate()
	assert.False(t, org.CanOperate())

	_ = org.Activate()
	_ = org.Delete()
	assert.False(t, org.CanOperate())
}

func TestSlug(t *testing.T) {
	t.Run("NewSlug with value", func(t *testing.T) {
		s, err := NewSlug("Test Name")
		require.NoError(t, err)
		assert.Equal(t, "test_name", s.String())
		assert.False(t, s.IsEmpty())
	})

	t.Run("NewSlug with empty value", func(t *testing.T) {
		_, err := NewSlug("")
		assert.ErrorIs(t, err, ErrEmptySlug)
	})
}
