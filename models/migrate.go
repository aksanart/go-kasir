package models

type Model struct {
	Model interface{}
}

func MigrateModels() []Model {
	return []Model{
		{Model: Order{}},
		{Model: Payment{}},
		{Model: Cashier{}},
		{Model: Category{}},
		{Model: Discount{}},
		{Model: Product{}},
	}
}
