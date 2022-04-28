package config

type (
	CategoryParam struct {
		Limit int `form:"limit" binding:"required,numeric,min=0"`
		Skip  int `form:"skip" binding:"numeric"`
	}

	CategoryParamJson struct {
		Name string `json:"name" binding:"required"`
	}

	ProductParam struct {
		Limit      int    `form:"limit" binding:"required,numeric,min=0"`
		Skip       int    `form:"skip" binding:"numeric"`
		CategoryID int    `form:"categoryId" binding:"numeric"`
		Q          string `form:"q"`
	}

	ProductParamJson struct {
		CategoryID int    `json:"categoryId" binding:"required,numeric,min=0"`
		Name       string `json:"name" binding:"required"`
		Image      string `json:"image" binding:"required"`
		Price      int    `json:"price" binding:"required,numeric,min=0"`
		Stock      int    `json:"stock" binding:"required,numeric,min=0"`
		Discount   struct {
			Qty       int    `json:"qty"`
			Type      string `json:"type"`
			Result    int    `json:"result"`
			ExpiredAt int    `json:"expiredAt"`
		} `json:"discount"`
	}

	ProductParamJsonUpdate struct {
		CategoryID int    `json:"categoryId,omitempty"`
		Name       string `json:"name,omitempty"`
		Image      string `json:"image,omitempty"`
		Price      int    `json:"price,omitempty"`
		Stock      int    `json:"stock,omitempty"`
	}

	CashierParamJson struct {
		Name     string `json:"name,omitempty" binding:"required"`
		Passcode string `json:"passcode,omitempty" binding:"required"`
	}

	CashierParamUpdate struct {
		Name     string `json:"name,omitempty"`
		Passcode string `json:"passcode,omitempty"`
	}

	PaymentParam struct {
		Limit    int `form:"limit" binding:"required,numeric,min=0"`
		Skip     int `form:"skip" binding:"numeric"`
		SubTotal int `form:"subtotal" binding:"numeric"`
	}

	PaymentParamJson struct {
		Name string `json:"name,omitempty" binding:"required"`
		Type string `json:"type,omitempty" binding:"required"`
		Logo string `json:"logo,omitempty"`
	}

	PaymentParamUpdate struct {
		Name string `json:"name,omitempty"`
		Type string `json:"type,omitempty"`
		Logo string `json:"logo,omitempty"`
	}

	OrderParam struct {
		ProductId int `json:"productId,omitempty" binding:"required,numeric"`
		Qty       int `json:"qty,omitempty" binding:"required"`
	}

	OrderParamJson struct {
		PaymentId int          `json:"paymentId,omitempty" binding:"required"`
		TotalPaid int64        `json:"totalPaid,omitempty" binding:"required"`
		Products  []OrderParam `json:"products"`
		CashierId int          `json:"cashierId"`
	}
)
