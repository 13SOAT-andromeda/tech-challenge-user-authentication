package address

import (
	"testing"
	"tech-challenge-user-validation/internal/core/domain"
	"github.com/stretchr/testify/assert"
)

func TestAddressModel_ToDomain(t *testing.T) {
	model := &Model{
		Street:     "Main St",
		Number:     "123",
		Complement: "Apt 4B",
		City:       "Springfield",
		State:      "IL",
		ZipCode:    "62701",
	}

	domainAddress := model.ToDomain()

	assert.NotNil(t, domainAddress)
	assert.Equal(t, model.Street, domainAddress.Street)
	assert.Equal(t, model.Number, domainAddress.Number)
	assert.Equal(t, model.Complement, domainAddress.Complement)
	assert.Equal(t, model.City, domainAddress.City)
	assert.Equal(t, model.State, domainAddress.State)
	assert.Equal(t, model.ZipCode, domainAddress.ZipCode)
}

func TestAddressModel_ToDomain_Nil(t *testing.T) {
	var model *Model = nil
	domainAddress := model.ToDomain()
	assert.Nil(t, domainAddress)
}

func TestAddressModel_FromDomain(t *testing.T) {
	domainAddress := &domain.Address{
		Street:     "Main St",
		Number:     "123",
		Complement: "Apt 4B",
		City:       "Springfield",
		State:      "IL",
		ZipCode:    "62701",
	}

	model := &Model{}
	model.FromDomain(domainAddress)

	assert.Equal(t, domainAddress.Street, model.Street)
	assert.Equal(t, domainAddress.Number, model.Number)
	assert.Equal(t, domainAddress.Complement, model.Complement)
	assert.Equal(t, domainAddress.City, model.City)
	assert.Equal(t, domainAddress.State, model.State)
	assert.Equal(t, domainAddress.ZipCode, model.ZipCode)
}

func TestAddressModel_FromDomain_Nil(t *testing.T) {
	model := &Model{Street: "Old St"}
	model.FromDomain(nil)
	assert.Equal(t, "Old St", model.Street)
}
