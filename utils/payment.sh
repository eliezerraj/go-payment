#!/bin/bash

var_acc=0
genAcc(){
    var_acc=$(($RANDOM%($max-$min+1)+$min))
}

var_amount=0
genAmount(){
    var_amount=$(($RANDOM%($max_amount-$min_amount+1)+$min_amount))
}

var_type_mcc=0
genMcc(){
    var_type_mcc=$(($RANDOM%($max_mcc-$min_mcc+1)+$min_mcc))
}

var_style_card=0
genStyleCard(){
    var_style_card=$(($RANDOM%($max_scc-$min_scc+1)+$min_scc))
}

var_type_card=0
genTypeCard(){
    var_type_card=$(($RANDOM%($max_tcc-$min_tcc+1)+$min_tcc))
}

declare -a arr_style_card
arr_style_card=(FISICO VIRTUAL)

declare -a arr_type_card
arr_type_card=(CREDIT DEBIT)

declare -a arr_mcc
arr_mcc=(FOOD GAS HOTEL AIRLINE TRANSPORT EDUCATION STORE GYM CINEMA PARKING CAR_RENTAL DRUG_STORE BEVERAGE )
# -----------------------------------------------------

arr_mcc=(HOTEL AIRLINE CAR_RENTAL)

min=1
max=500

min_amount=1000
max_amount=4000

min_mcc=0
max_mcc=2

min_scc=0
max_scc=1

min_tcc=0
max_tcc=1

dayfreq=1

sum_day=21

#domain=http://localhost:5007/payment/pay
domain=https://97x38r33ag.execute-api.us-east-2.amazonaws.com/Live/payment/pay

for (( x=100; x>0; x-- ))
do
    numdays=$((numdays + sum_day))
    new_dt=`date '+%Y-%m-%dT%T.%9N%:z' -d "+$numdays days"`
    echo new_dt
    genAcc
    for (( y=2; y>0; y-- ))
    do
        genStyleCard
        genTypeCard
            for (( z=3; z>0; z-- ))
            do
                genAmount
                genMcc
                echo curl -X POST $domain -H 'Content-Type: application/json' -d '{"account_id":"ACC-'$var_acc'","card_number":"111.222.333.'$var_acc'","payment_at":"'${new_dt}'","card_type":"'${arr_type_card[var_type_card]}'","card_style":"'${arr_style_card[var_style_card]}'","currency":"BRL","mcc":"'${arr_mcc[var_type_mcc]}'","amount":'$var_amount'}'
                     curl -X POST $domain -H 'Content-Type: application/json' -d '{"account_id":"ACC-'$var_acc'","card_number":"111.222.333.'$var_acc'","payment_at":"'${new_dt}'","card_type":"'${arr_type_card[var_type_card]}'","card_style":"'${arr_style_card[var_style_card]}'","currency":"BRL","mcc":"'${arr_mcc[var_type_mcc]}'","amount":'$var_amount'}'  
            done
    done
done

arr_mcc=(FOOD STORE CINEMA PARKING BEVERAGE)

min=1
max=500
min_amount=20
max_amount=120
min_mcc=0
max_mcc=4
sum_day=3

for (( x=100; x>0; x-- ))
do
    numdays=$((numdays + sum_day))
    new_dt=`date '+%Y-%m-%dT%T.%9N%:z' -d "+$numdays days"`
    genAcc
    for (( y=3; y>0; y-- ))
    do
        genStyleCard
        genTypeCard
            for (( z=5; z>0; z-- ))
            do
                genAmount
                genMcc
                echo curl -X POST $domain -H 'Content-Type: application/json' -d '{"account_id":"ACC-'$var_acc'","card_number":"111.222.333.'$var_acc'","payment_at":"'${new_dt}'","card_type":"'${arr_type_card[var_type_card]}'","card_style":"'${arr_style_card[var_style_card]}'","currency":"BRL","mcc":"'${arr_mcc[var_type_mcc]}'","amount":'$var_amount'}' 
                     curl -X POST $domain -H 'Content-Type: application/json' -d '{"account_id":"ACC-'$var_acc'","card_number":"111.222.333.'$var_acc'","payment_at":"'${new_dt}'","card_type":"'${arr_type_card[var_type_card]}'","card_style":"'${arr_style_card[var_style_card]}'","currency":"BRL","mcc":"'${arr_mcc[var_type_mcc]}'","amount":'$var_amount'}' 
            done
    done
done