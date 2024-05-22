package fixtures

import (
	"WB/internal/order/model"
	"WB/tests/states"
)

type PaymentBuilder struct {
	instance *model.Payment
}

func Payment() *PaymentBuilder {
	return &PaymentBuilder{instance: &model.Payment{}}
}

func (b *PaymentBuilder) Transaction(v string) *PaymentBuilder {
	b.instance.Transaction = v
	return b
}

func (b *PaymentBuilder) RequestID(v string) *PaymentBuilder {
	b.instance.RequestID = v
	return b
}

func (b *PaymentBuilder) Currency(v string) *PaymentBuilder {
	b.instance.Currency = v
	return b
}

func (b *PaymentBuilder) Provider(v string) *PaymentBuilder {
	b.instance.Provider = v
	return b
}

func (b *PaymentBuilder) Amount(v int) *PaymentBuilder {
	b.instance.Amount = v
	return b
}

func (b *PaymentBuilder) PaymentDT(v int) *PaymentBuilder {
	b.instance.PaymentDt = v
	return b
}

func (b *PaymentBuilder) Bank(v string) *PaymentBuilder {
	b.instance.Bank = v
	return b
}

func (b *PaymentBuilder) DeliveryCost(v int) *PaymentBuilder {
	b.instance.DeliveryCost = v
	return b
}

func (b *PaymentBuilder) GoodsTotal(v int) *PaymentBuilder {
	b.instance.GoodsTotal = v
	return b
}

func (b *PaymentBuilder) CustomFee(v int) *PaymentBuilder {
	b.instance.CustomFee = v
	return b
}

func (b *PaymentBuilder) OrderUID(v string) *PaymentBuilder {
	b.instance.OrderUID = v
	return b
}

func (b *PaymentBuilder) Valid1() *PaymentBuilder {
	return Payment().OrderUID(states.OrderUID1).Transaction(states.Payment1.Transaction).RequestID(states.Payment1.RequestID).Currency(states.Payment1.Currency).Provider(states.Payment1.Provider).Amount(states.Payment1.Amount).PaymentDT(states.Payment1.PaymentDt).Bank(states.Payment1.Bank).DeliveryCost(states.Payment1.DeliveryCost).GoodsTotal(states.Payment1.GoodsTotal).CustomFee(states.Payment1.CustomFee)
}

func (b *PaymentBuilder) Ptr() *model.Payment {
	return b.instance
}

func (b *PaymentBuilder) Val() model.Payment {
	return *b.instance
}
