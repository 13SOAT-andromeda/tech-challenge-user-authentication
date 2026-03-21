package address

import "tech-challenge-user-validation/internal/core/domain"

type Model struct {
	Street     string `gorm:"column:street"`
	Number     string `gorm:"column:number"`
	Complement string `gorm:"column:complement"`
	City       string `gorm:"column:city"`
	State      string `gorm:"column:state"`
	ZipCode    string `gorm:"column:zip_code"`
}

func (m *Model) ToDomain() *domain.Address {
	if m == nil {
		return nil
	}
	return &domain.Address{
		Street:     m.Street,
		Number:     m.Number,
		Complement: m.Complement,
		City:       m.City,
		State:      m.State,
		ZipCode:    m.ZipCode,
	}
}

func (m *Model) FromDomain(d *domain.Address) {
	if d == nil {
		return
	}
	m.Street = d.Street
	m.Number = d.Number
	m.Complement = d.Complement
	m.City = d.City
	m.State = d.State
	m.ZipCode = d.ZipCode
}
