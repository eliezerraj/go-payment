# go-payment

POC for test purposes.

CRUD a payment

## Diagram

## database

    CREATE TABLE payment (
        id                  SERIAL PRIMARY KEY,
        fk_account_id       integer REFERENCES account(id),
        card_number         varchar(200) NULL,
        card_type           varchar(200) NULL,
        status              varchar(200) NULL,
        currency            varchar(10) NULL,   
        amount              float8 NULL,
        mcc                 varchar(10) NULL,
        create_at           timestamptz NULL,
        update_at           timestamptz NULL,
        tenant_id           varchar(200) NULL
    );

    CREATE TABLE card (
        id                  SERIAL PRIMARY KEY,
        fk_account_id       integer REFERENCES account(id),
        card_number         varchar(200) NULL,
        card_type           varchar(200) NULL,
        card_pin            varchar(200) NULL,
        status              varchar(200) NULL,
        expire_at           timestamptz NULL,
        create_at           timestamptz NULL,
        update_at           timestamptz NULL,
        tenant_id           varchar(200) NULL
    );

## Endpoints

+ POST /payment/pay

        {
            "account_id":"ACC-1",
            "card_number":"111.222.333.444",
            "card_type":"DEBIT",
            "currency":"BRL",
            "mcc": "FOOD",
            "amount":12
        }
        
+ GET  /payment/get/6