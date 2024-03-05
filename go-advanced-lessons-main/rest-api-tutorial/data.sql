-- Таблица покупатели
CREATE TABLE public.customers
(
    customer_id   UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    entity_type   VARCHAR(50) NOT NULL, -- "физ" или "юр"
    contact_name  VARCHAR(100) NOT NULL, -- ФИО или наименование организации
    address       VARCHAR(200),
    phone         VARCHAR(20)
);

-- Таблица поставщики
CREATE TABLE public.suppliers
(
    supplier_id   UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    organization  VARCHAR(100) NOT NULL,
    address       VARCHAR(200),
    phone         VARCHAR(20)
);

-- Таблица виды изделия
CREATE TABLE public.product_types
(
    type_id       UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    type_name     VARCHAR(100) NOT NULL
);

-- Таблица о каталоге изделия
CREATE TABLE public.products
(
    product_id    UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    type_id       UUID NOT NULL,
    product_name  VARCHAR(100) NOT NULL,
    weight        DECIMAL(10, 2) NOT NULL, -- Вес изделия в кг
    unit          VARCHAR(20) NOT NULL, -- Единица измерения
    description   TEXT,
    price_pickup  DECIMAL(10, 2) NOT NULL, -- Цена самовывоза в рублях
    price_delivery DECIMAL(10, 2) NOT NULL, -- Цена доставки в рублях

    CONSTRAINT product_type_fk FOREIGN KEY (type_id) REFERENCES public.product_types (type_id)
);

-- Таблица приходной накладной
CREATE TABLE public.incoming_invoice
(
    invoice_id    UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    product_id    UUID NOT NULL,
    quantity      INT NOT NULL,
    price         DECIMAL(10, 2) NOT NULL,
    transaction_date TIMESTAMP NOT NULL,
    payment_method VARCHAR(50),
    supplier_id   UUID NOT NULL,

    CONSTRAINT invoice_product_fk FOREIGN KEY (product_id) REFERENCES public.products (product_id),
    CONSTRAINT invoice_supplier_fk FOREIGN KEY (supplier_id) REFERENCES public.suppliers (supplier_id)
);

-- Таблица Расходной накладной
CREATE TABLE public.outgoing_invoice
(
    invoice_id    UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    product_id    UUID NOT NULL,
    quantity      INT NOT NULL,
    price         DECIMAL(10, 2) NOT NULL,
    transaction_date TIMESTAMP NOT NULL,
    payment_method VARCHAR(50),
    customer_id   UUID NOT NULL,

    CONSTRAINT invoice_product_fk FOREIGN KEY (product_id) REFERENCES public.products (product_id),
    CONSTRAINT invoice_customer_fk FOREIGN KEY (customer_id) REFERENCES public.customers (customer_id)
);

-- Таблица о наличии товара
CREATE TABLE public.inventory
(
    product_id          UUID PRIMARY KEY,
    type_id             UUID NOT NULL,
    quantity           INT NOT NULL,
    unit               VARCHAR(20) NOT NULL,
    last_receipt_date  TIMESTAMP NOT NULL,
    availability       BOOLEAN NOT NULL,

    CONSTRAINT inventory_product_fk FOREIGN KEY (product_id) REFERENCES public.products (product_id),
    CONSTRAINT inventory_type_fk FOREIGN KEY (type_id) REFERENCES public.product_types (type_id)
);

INSERT INTO product_types (type_name) VALUES ('Колодцы унифицированные'); -- 8287af71-d66d-4071-a373-abe6046a9f42
INSERT INTO product_types (type_name) VALUES ('Кабельные колодцы (ККС)'); -- 976b7521-5ff5-4542-8869-d56a012f963a
INSERT INTO product_types (type_name) VALUES ('Кольца стеновые с четвертью (КС, КЦ, К)'); -- 142a9a5e-7d47-4b96-9124-6d6c6f848a31
INSERT INTO product_types (type_name) VALUES ('Кольца стеновые с днищем(КСД, КЦД)'); --74d15225-287c-4d1f-b749-f119b4a8d125

INSERT INTO products (type_id, product_name, weight, unit, description, price_pickup, price_delivery) VALUES ('8287af71-d66d-4071-a373-abe6046a9f42', 'ВГ 15', 2820, 'шт.', 'Объём - 1,128 м^3. Унифицированный', 9520, 12340); -- 94b2b277-2468-4348-b0e8-fb775ce48370
INSERT INTO products (type_id, product_name, weight, unit, description, price_pickup, price_delivery) VALUES ('8287af71-d66d-4071-a373-abe6046a9f42', 'ВГ 12', 1930, 'шт.', 'Объём - 0,772 м^3. Унифицированный', 7840, 9770); -- 097d1ed7-a839-432b-9916-e7a3f234f2f0

ALTER TABLE public.customers
    ADD COLUMN user_priority VARCHAR(50) DEFAULT 'покупатель' NOT NULL;

ALTER TABLE public.customers
    ADD COLUMN login VARCHAR(100) NOT NULL,
    ADD COLUMN password VARCHAR(100) NOT NULL;

UPDATE public.customers
SET address = 'none'
WHERE address IS NULL;

UPDATE public.customers
SET user_priority = 'покупатель'
WHERE user_priority IS NULL;

ALTER TABLE public.customers
    ADD COLUMN email VARCHAR(50);