
CREATE TABLE IF NOT EXISTS orders (
    order_uid          VARCHAR(255) PRIMARY KEY,
    track_number       VARCHAR(255),
    entry              VARCHAR(255),
    locale             VARCHAR(10),
    internal_signature VARCHAR(255),
    customer_id        VARCHAR(255),
    delivery_service   VARCHAR(255),
    shardkey           VARCHAR(10),
    sm_id              INTEGER,
    date_created       TIMESTAMP WITH TIME ZONE,
    oof_shard          VARCHAR(10)
);


CREATE TABLE IF NOT EXISTS delivery (
    order_uid VARCHAR(255) PRIMARY KEY REFERENCES orders(order_uid) ON DELETE CASCADE,
    name      VARCHAR(255),
    phone     VARCHAR(20),
    zip       VARCHAR(20),
    city      VARCHAR(100),
    address   VARCHAR(255),
    region    VARCHAR(100),
    email     VARCHAR(100)
);


CREATE TABLE IF NOT EXISTS payment (
    order_uid      VARCHAR(255) PRIMARY KEY REFERENCES orders(order_uid) ON DELETE CASCADE,
    transaction    VARCHAR(255),
    request_id     VARCHAR(255),
    currency       VARCHAR(10),
    provider       VARCHAR(50),
    amount         INTEGER,
    payment_dt     BIGINT,
    bank           VARCHAR(100),
    delivery_cost  INTEGER,
    goods_total    INTEGER,
    custom_fee     INTEGER
);


CREATE TABLE IF NOT EXISTS items (
    id           SERIAL PRIMARY KEY,
    order_uid    VARCHAR(255) REFERENCES orders(order_uid) ON DELETE CASCADE,
    chrt_id      INTEGER,
    track_number VARCHAR(255),
    price        INTEGER,
    rid          VARCHAR(255),
    name         VARCHAR(255),
    sale         INTEGER,
    size         VARCHAR(10),
    total_price  INTEGER,
    nm_id        INTEGER,
    brand        VARCHAR(100),
    status       INTEGER
);


CREATE INDEX IF NOT EXISTS idx_items_order_uid ON items(order_uid);
CREATE INDEX IF NOT EXISTS idx_orders_date_created ON orders(date_created);