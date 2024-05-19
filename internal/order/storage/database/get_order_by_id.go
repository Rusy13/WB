package storage

import (
	"context"
	"errors"

	"WB/internal/order/model"
	"WB/internal/order/storage"
	"WB/internal/order/storage/database/dto"
	"github.com/jackc/pgx/v4"
)

func (s *OrderStorageDB) GetOrderByID(ctx context.Context, orderUID string) (*model.Order, error) {
	var orderDB dto.OrderDB

	// Fetch order details
	err := s.db.QueryRow(ctx,
		`SELECT order_uid, track_number, entry, locale, internal_signature, customer_id, delivery_service, shardkey, sm_id, date_created, oof_shard
		 FROM orders WHERE order_uid = $1`, orderUID).Scan(
		&orderDB.OrderUID, &orderDB.TrackNumber, &orderDB.Entry, &orderDB.Locale, &orderDB.InternalSignature,
		&orderDB.CustomerID, &orderDB.DeliveryService, &orderDB.ShardKey, &orderDB.SmID, &orderDB.DateCreated, &orderDB.OofShard)
	if err != nil {
		if errors.Is(err, pgx.ErrNoRows) {
			return nil, storage.ErrOrderNotFound
		}
		return nil, err
	}

	// Fetch delivery details
	err = s.db.QueryRow(ctx,
		`SELECT name, phone, zip, city, address, region, email 
		 FROM delivery WHERE order_uid = $1`, orderUID).Scan(
		&orderDB.Delivery.Name, &orderDB.Delivery.Phone, &orderDB.Delivery.Zip, &orderDB.Delivery.City,
		&orderDB.Delivery.Address, &orderDB.Delivery.Region, &orderDB.Delivery.Email)
	if err != nil {
		return nil, err
	}

	// Fetch payment details
	err = s.db.QueryRow(ctx,
		`SELECT transaction, request_id, currency, provider, amount, payment_dt, bank, delivery_cost, goods_total, custom_fee
		 FROM payment WHERE order_uid = $1`, orderUID).Scan(
		&orderDB.Payment.Transaction, &orderDB.Payment.RequestID, &orderDB.Payment.Currency, &orderDB.Payment.Provider,
		&orderDB.Payment.Amount, &orderDB.Payment.PaymentDt, &orderDB.Payment.Bank, &orderDB.Payment.DeliveryCost,
		&orderDB.Payment.GoodsTotal, &orderDB.Payment.CustomFee)
	if err != nil {
		return nil, err
	}

	// Fetch items details
	rows, err := s.db.Query(ctx,
		`SELECT chrt_id, track_number, price, rid, name, sale, size, total_price, nm_id, brand, status
		 FROM items WHERE order_uid = $1`, orderUID)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	for rows.Next() {
		var item dto.ItemDB
		err := rows.Scan(
			&item.ChrtID, &item.TrackNumber, &item.Price, &item.Rid, &item.Name,
			&item.Sale, &item.Size, &item.TotalPrice, &item.NmID, &item.Brand, &item.Status)
		if err != nil {
			return nil, err
		}
		orderDB.Items = append(orderDB.Items, item)
	}

	if err = rows.Err(); err != nil {
		return nil, err
	}

	order := dto.ConvertToOrder(orderDB)
	return &order, nil
}
