package database

import (
	"database/sql"
	"fmt"
	"log"
	"order-service/internal/model"
	"os"
)

type Database struct {
	db *sql.DB
}

func NewPostrgesDB() (*Database, error){
	db, err := sql.Open("postgres", fmt.Sprintf("postgres://%s:%s@%s:%s/%s?sslmode=disable",
	os.Getenv("POSTGRES_USER"),
	os.Getenv("POSTGRES_PASSWORD"),
	os.Getenv("POSTGRES_HOST"),
	os.Getenv("POSTGRES_PORT"),
	os.Getenv("POSTGRES_DB"),
	))
	if err != nil {
		log.Fatalf("Не удалось подключить базу данных: %v", err)
		return nil, fmt.Errorf("Не удалось подключить базу данных: %v", err)
	}

	if err = db.Ping(); err != nil {
		return nil, fmt.Errorf("Не удалось пингануть базу данных: %v", err)
	}

	log.Println("Удалось подключиться к базе данных")

	return &Database{db: db}, nil
}

func (d *Database) Close() error {
	return d.db.Close()
}

func (d *Database) SaveOrder(order model.Order) error {
	tx, err := d.db.Begin()
	if err != nil {
		return fmt.Errorf("Не удалось начать транзакцию: %v", err)
	}
	defer tx.Rollback()

	_, err = tx.Exec(`
		INSERT INTO orders (order_uid, track_number, entry, locale, internal_signature, 
			customer_id, delivery_service, shardkey, sm_id, date_created, oof_shard)
		VALUES ($1, $2, $3, $4, $5, $6, $7, $8, $9, $10, $11)
		ON CONFLICT (order_uid) DO UPDATE SET
			track_number = EXCLUDED.track_number,
			entry = EXCLUDED.entry,
			locale = EXCLUDED.locale,
			internal_signature = EXCLUDED.internal_signature,
			customer_id = EXCLUDED.customer_id,
			delivery_service = EXCLUDED.delivery_service,
			shardkey = EXCLUDED.shardkey,
			sm_id = EXCLUDED.sm_id,
			date_created = EXCLUDED.date_created,
			oof_shard = EXCLUDED.oof_shard
	`, order.OrderUID, order.TrackNumber, order.Entry, order.Locale, order.InternalSignature,
		order.CustomerID, order.DeliveryService, order.Shardkey, order.SmID, order.DateCreated, order.OofShard)
	if err != nil {
		return fmt.Errorf("Не удалось добавить заказ: %v", err)
	}

	_, err = tx.Exec(`
		INSERT INTO delivery (order_uid, name, phone, zip, city, address, region, email)
		VALUES ($1, $2, $3, $4, $5, $6, $7, $8)
		ON CONFLICT (order_uid) DO UPDATE SET
			name = EXCLUDED.name,
			phone = EXCLUDED.phone,
			zip = EXCLUDED.zip,
			city = EXCLUDED.city,
			address = EXCLUDED.address,
			region = EXCLUDED.region,
			email = EXCLUDED.email
	`, order.OrderUID, order.Delivery.Name, order.Delivery.Phone, order.Delivery.Zip,
		order.Delivery.City, order.Delivery.Address, order.Delivery.Region, order.Delivery.Email)
	if err != nil {
		return fmt.Errorf("Не удалось добавить доставку: %v", err)
	}

	_, err = tx.Exec(`
		INSERT INTO payment (order_uid, transaction, request_id, currency, provider, 
			amount, payment_dt, bank, delivery_cost, goods_total, custom_fee)
		VALUES ($1, $2, $3, $4, $5, $6, $7, $8, $9, $10, $11)
		ON CONFLICT (order_uid) DO UPDATE SET
			transaction = EXCLUDED.transaction,
			request_id = EXCLUDED.request_id,
			currency = EXCLUDED.currency,
			provider = EXCLUDED.provider,
			amount = EXCLUDED.amount,
			payment_dt = EXCLUDED.payment_dt,
			bank = EXCLUDED.bank,
			delivery_cost = EXCLUDED.delivery_cost,
			goods_total = EXCLUDED.goods_total,
			custom_fee = EXCLUDED.custom_fee
	`, order.OrderUID, order.Payment.Transaction, order.Payment.RequestID, order.Payment.Currency,
		order.Payment.Provider, order.Payment.Amount, order.Payment.PaymentDt, order.Payment.Bank,
		order.Payment.DeliveryCost, order.Payment.GoodsTotal, order.Payment.CustomFee)

	if err != nil {
		return fmt.Errorf("Не удалось сохранить информацию об оплате: %v", err)
	}

	_, err = tx.Exec("DELETE FROM items WHERE order_uid = $1", order.OrderUID)

	if err != nil {
		return fmt.Errorf("Не удалось удалить старые товары: %v", err)
	}


	for _, item := range order.Items {
		_, err = tx.Exec(`
			INSERT INTO items (order_uid, chrt_id, track_number, price, rid, name, 
				sale, size, total_price, nm_id, brand, status)
			VALUES ($1, $2, $3, $4, $5, $6, $7, $8, $9, $10, $11, $12)
		`, order.OrderUID, item.ChrtID, item.TrackNumber, item.Price, item.Rid, item.Name,
			item.Sale, item.Size, item.TotalPrice, item.NmID, item.Brand, item.Status)
		if err != nil {
			return fmt.Errorf("Не удалось добавить товары: %v", err)
		}
	}

	return tx.Commit()

}



func (d *Database) GetOrder(uid string) (model.Order, error) {
	var order model.Order
	
	err := d.db.QueryRow(`
		SELECT order_uid, track_number, entry, locale, internal_signature, 
			customer_id, delivery_service, shardkey, sm_id, date_created, oof_shard
		FROM orders WHERE order_uid = $1
	`, uid).Scan(
		&order.OrderUID, &order.TrackNumber, &order.Entry, &order.Locale, &order.InternalSignature,
		&order.CustomerID, &order.DeliveryService, &order.Shardkey, &order.SmID, &order.DateCreated, &order.OofShard,
	)
	if err != nil {
		if err == sql.ErrNoRows {
			return order, fmt.Errorf("Заказ не найден")
		}
		return order, fmt.Errorf("Не удалось получить заказ: %v", err)
	}


	err = d.db.QueryRow(`
		SELECT name, phone, zip, city, address, region, email
		FROM delivery WHERE order_uid = $1
	`, uid).Scan(
		&order.Delivery.Name, &order.Delivery.Phone, &order.Delivery.Zip, &order.Delivery.City,
		&order.Delivery.Address, &order.Delivery.Region, &order.Delivery.Email,
	)
	if err != nil {
		return order, fmt.Errorf("Не удалось получить информацию о доставке: %v", err)
	}
	order.Delivery.OrderUID = uid


	err = d.db.QueryRow(`
		SELECT transaction, request_id, currency, provider, amount, payment_dt, 
			bank, delivery_cost, goods_total, custom_fee
		FROM payment WHERE order_uid = $1
	`, uid).Scan(
		&order.Payment.Transaction, &order.Payment.RequestID, &order.Payment.Currency, &order.Payment.Provider,
		&order.Payment.Amount, &order.Payment.PaymentDt, &order.Payment.Bank, &order.Payment.DeliveryCost,
		&order.Payment.GoodsTotal, &order.Payment.CustomFee,
	)
	if err != nil {
		return order, fmt.Errorf("Не удалось получить информацию об оплате: %v", err)
	}
	order.Payment.OrderUID = uid


	rows, err := d.db.Query(`
		SELECT chrt_id, track_number, price, rid, name, sale, size, total_price, nm_id, brand, status
		FROM items WHERE order_uid = $1
	`, uid)
	if err != nil {
		return order, fmt.Errorf("Не удалось получить товар: %v", err)
	}
	defer rows.Close()

	for rows.Next() {
		var item model.Item
		err = rows.Scan(
			&item.ChrtID, &item.TrackNumber, &item.Price, &item.Rid, &item.Name,
			&item.Sale, &item.Size, &item.TotalPrice, &item.NmID, &item.Brand, &item.Status,
		)
		if err != nil {
			return order, fmt.Errorf("Не удалось засканить товар: %v", err)
		}
		item.OrderUID = uid
		order.Items = append(order.Items, item)
	}

	return order, nil
}


func (d *Database) GetAllOrders() (map[string]model.Order, error) {
	orders := make(map[string]model.Order)

	rows, err := d.db.Query(`
		SELECT o.order_uid, o.track_number, o.entry, o.locale, o.internal_signature, 
			o.customer_id, o.delivery_service, o.shardkey, o.sm_id, o.date_created, o.oof_shard,
			d.name, d.phone, d.zip, d.city, d.address, d.region, d.email,
			p.transaction, p.request_id, p.currency, p.provider, p.amount, p.payment_dt, 
			p.bank, p.delivery_cost, p.goods_total, p.custom_fee
		FROM orders o
		LEFT JOIN delivery d ON o.order_uid = d.order_uid
		LEFT JOIN payment p ON o.order_uid = p.order_uid
	`)
	if err != nil {
		return nil, fmt.Errorf("Не удалось получить заказы: %v", err)
	}
	defer rows.Close()

	for rows.Next() {
		var order model.Order
		err := rows.Scan(
			&order.OrderUID, &order.TrackNumber, &order.Entry, &order.Locale, &order.InternalSignature,
			&order.CustomerID, &order.DeliveryService, &order.Shardkey, &order.SmID, &order.DateCreated, &order.OofShard,
			&order.Delivery.Name, &order.Delivery.Phone, &order.Delivery.Zip, &order.Delivery.City,
			&order.Delivery.Address, &order.Delivery.Region, &order.Delivery.Email,
			&order.Payment.Transaction, &order.Payment.RequestID, &order.Payment.Currency, &order.Payment.Provider,
			&order.Payment.Amount, &order.Payment.PaymentDt, &order.Payment.Bank, &order.Payment.DeliveryCost,
			&order.Payment.GoodsTotal, &order.Payment.CustomFee,
		)
		if err != nil {
			log.Printf("Не удалось обработать заказ: %v", err)
			continue
		}
		

		itemRows, err := d.db.Query(`
			SELECT chrt_id, track_number, price, rid, name, sale, size, total_price, nm_id, brand, status
			FROM items WHERE order_uid = $1
		`, order.OrderUID)
		if err != nil {
			log.Printf("Не удалось получить товары для заказа %s: %v", order.OrderUID, err)
			continue
		}
		
		order.Items = make([]model.Item, 0)
		for itemRows.Next() {
			var item model.Item
			err = itemRows.Scan(
				&item.ChrtID, &item.TrackNumber, &item.Price, &item.Rid, &item.Name,
				&item.Sale, &item.Size, &item.TotalPrice, &item.NmID, &item.Brand, &item.Status,
			)
			if err != nil {
				log.Printf("Не удалось обработать товар для заказа %s: %v", order.OrderUID, err)
				continue
			}
			order.Items = append(order.Items, item)
		}
		itemRows.Close()
		
		orders[order.OrderUID] = order
	}

	if err = rows.Err(); err != nil {
		return nil, fmt.Errorf("Ошибка при обработке результатов заказов: %v", err)
	}

	return orders, nil
}

