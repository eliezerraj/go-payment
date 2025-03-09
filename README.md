# go-payment

POC for test purposes.

CRUD a payment

## Diagram

go-payment (get:/get) == (REST) ==> go-account (service.GetAccount)
go-payment (get:/fundBalanceAccount) == (REST) ==> go-account (service.GetAccount)
go-payment ==> Repository

## database

See repo https://github.com/eliezerraj/go-account-migration-worker.git

## Endpoints

+ GET /header

+ GET /info

+ GET  /get/6

+ GET  /infoPodGrpc

+ POST /checkFeaturePaymentFraudGrpc

        {
            "card_number":"111.111.000.001",
            "terminal_name": "TERM-1",
            "coord_x": 90,
            "coord_y": 30,
            "card_type":"CREDIT",
            "card_model":"VIRTUAL",
            "mcc":"COMPUTE",
            "status":"OK",
            "currency":"BRL",
            "amount": 300.55,
            "payment_at":"2024-02-14T22:59:01.859507132-03:00",
            "tx_1d": 2,
            "avg_1d": 300.57,
            "tx_7d": 3,
            "avg_7d": 300.12,
            "tx_30d": 6,
            "avg_30d": 900.82,
            "time_btw_cc_tx": 60
        }

+ POST /paymentWithCheckFraud

        {
        "card_number":"111.000.000.001",
        "terminal_name": "TERM-12",
        "card_type":"CREDIT",
        "card_model":"VIRTUAL",
        "currency":"BRL",
        "mcc":"STORE",
        "amount":330.00
        }

+ POST /payment

        {
        "card_number":"111.000.000.001",
        //"payment_at":"2024-02-14T22:59:01.859507132-03:00",
        "terminal_name": "TERM-1",
        "card_type":"CREDIT",
        "card_model":"CHIP",
        "currency":"BRL",
        "mcc":"STORE",
        "amount":52,
        "fraud": 0
        }

## Compile grpc proto

    protoc -I proto proto/fraud.proto --go_out=plugins=grpc:proto
