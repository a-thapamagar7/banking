package domain

type CustomerRepositoryStub struct {
	customers []Customer
}

func (s CustomerRepositoryStub) FindAll() ([]Customer, error) {
	return s.customers, nil
}

func NewCustomerRespositoryStub() CustomerRepositoryStub {
	customers := []Customer{
		{"1001", "Ashish", "New Delhi", "110011", "2001-01-01", "1"},
		{"1002", "Rob", "New Delhi", "110011", "2001-01-01", "1"},
	}
	return CustomerRepositoryStub{customers}
}
