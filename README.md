# go-payment

POC for test purposes.

CRUD a payment

## Diagram

go-payment (get:/get) == (REST) ==> go-account (service.GetAccount)
go-payment (get:/fundBalanceAccount) == (REST) ==> go-account (service.GetAccount)
go-payment ==> Repository

## database

    CREATE TABLE payment (
        id                  SERIAL PRIMARY KEY,
        fk_card_id          integer REFERENCES card(id),
        card_number         varchar(200) NULL,
        fk_terminal_id      integer REFERENCES terminal(id),
        terminal_name       varchar(200) NULL,
        card_type           varchar(200) NULL,
        card_model          varchar(200) NULL,
        payment_at          timestamptz NULL,
        mcc                 varchar(10) NULL,
        status              varchar(200) NULL,
        currency            varchar(10) NULL,   
        amount              float8 NULL,
        create_at           timestamptz NULL,
        update_at           timestamptz NULL,
        fraud            float8 NULL,
        tenant_id           varchar(200) NULL
    );

    CREATE INDEX payment_idx ON payment (card_number);

    CREATE TABLE card (
        id                  SERIAL PRIMARY KEY,
        fk_account_id       integer REFERENCES account(id),
        card_number         varchar(200) NULL,
        card_type           varchar(200) NULL,
        card_model           varchar(200) NULL,
        card_pin            varchar(200) NULL,
        status              varchar(200) NULL,
        expire_at           timestamptz NULL,
        create_at           timestamptz NULL,
        update_at           timestamptz NULL,
        tenant_id           varchar(200) NULL
    );

    CREATE INDEX card_idx ON card (card_number);

    CREATE TABLE terminal (
        id                  SERIAL PRIMARY KEY,
        terminal_name       varchar(200) NULL,
        coord_x             float8 NULL,
        coord_y             float8 NULL,
        status              varchar(200) NULL,
        create_at           timestamptz NULL,
        update_at           timestamptz NULL
    );

## Endpoints

+ POST /payment/pay

        {
        "card_number":"111.111.000.001",
        "payment_at":"2024-02-14T22:59:01.859507132-03:00",
        "terminal_name": "TERM-1",
        "card_type":"CREDIT",
        "card_style":"CHIP",
        "currency":"BRL",
        "mcc":"STORE",
        "amount":52
        }
        
+ GET  /payment/get/6

+ GET  /podGrpc

## Compile grpc proto

    protoc -I proto proto/fraud.proto --go_out=plugins=grpc:proto
