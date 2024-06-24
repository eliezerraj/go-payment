#!/bin/bash

var_cc=0
genCC(){
    var_cc=$(($RANDOM%($max_cc-$min_cc+1)+$min_cc))
}

var_tx=0
genTXDay(){
    var_tx=$(($RANDOM%($max_tx-$min_tx+1)+$min_tx))
}

var_min=0
genMinutes(){
    var_min=$(($RANDOM%($max_minutes-$min_minutes+1)+$min_minutes))
}

var_skip=0
genSkip(){
    var_skip=$(($RANDOM%($max_skip-$min_skip+1)+$min_skip))
}

var_type_mcc=0
genMcc(){
    var_type_mcc=$(($RANDOM%($max_mcc-$min_mcc+1)+$min_mcc))
}

# Normal ditribuition
var_amount=0
genAmount(){
    a=$(($RANDOM%($max_amount-$min_amount+1)+$min_amount))
    b=$(($RANDOM%($max_amount-$min_amount+1)+$min_amount))
    c=$(($RANDOM%($max_amount-$min_amount+1)+$min_amount))
    var_amount=$(( (a+b+c)/3 ))
}

var_term=0
genTerm(){
    var_term=$(($RANDOM%($max_term-$min_term+1)+$min_term))
}

var_model_card=0
genModelCard(){
    var_model_card=$(($RANDOM%($max_model-$min_model+1)+$min_model))
}

var_type_card=0
genTypeCard(){
    var_type_card=$(($RANDOM%($max_tcc-$min_tcc+1)+$min_tcc))
}

echo "-------------------------------------"
echo "-----------STARTING DAY---------------"
echo "-------------------------------------"
#token=eyJhbGciOiJIUzI1NiIsInR5cCI6IkpXVCJ9.eyJ1c2VybmFtZSI6ImFkbWluIiwic2NvcGUiOlsiYWRtaW4iXSwiZXhwIjoxNzEzNDAwODM0fQ.iUebcV3URus1K1r13X9E-EhfYfEdFZQbEBWoCKD9N8A
#domain=https://go-api-global.architecture.caradhras.io/payment/payment/pay
domain=http://localhost:5007/payment/pay

declare -a arr_mcc
arr_mcc=(PARKING BEVERAGE FOOD LAUNDRY CINEMA BOOK GIFT CASH GAS PET DRUG_STORE COSMETIC GYM STORE SPORTING COMPUTER MOTOR)

declare -a arr_model_card
arr_model_card=(VIRTUAL CHIP)

declare -a arr_type_card
arr_type_card=(CREDIT DEBIT)

cc_pan=111000000001

min_model=0
max_model=1
min_tcc=0
max_tcc=1
min_mcc=0  # start idx arr-mcc
max_mcc=16 # final idx arr-mcc
min_cc=1 # credit card
max_cc=2500 # credit card
min_tx=0  #min transaction
max_tx=2  #max transaction
min_minutes=900 # min qtd between transaction (15min)
max_minutes=3600 # max qtd between transaction (60min)
min_amount=20 # min amount transaction
max_amount=800 # max amount transaction
min_term=1
max_term=100
var_fraud=0 # NO FRAUD
fraud_rate=1 # NO FRAUD RATE

for c in {0..30..1}
do
    genCC

    tmp_cc_pan=$(($cc_pan+$var_cc))       
    cc_final="${tmp_cc_pan:0:3}"."${tmp_cc_pan:3:3}"."${tmp_cc_pan:6:3}"."${tmp_cc_pan:9:3}"

    echo "==============================="
    echo "=====> " $cc_final "<====="
    echo "==============================="
    
    min_skip=0
    max_skip=7
    genSkip

    #for d in {0..30..1} 
    for d in $(eval echo "{0..30..$var_skip}")
    do
        echo "day => "$d
        
        min_skip=6
        max_skip=17
        genSkip

        #for h in {0..22..1} # Hour
        for h in $(eval echo "{0..22..$var_skip}")
        do
            echo "---> Hour => "$h
            genTXDay
            for (( z=var_tx; z>0; z-- ))
            do
                genMinutes
                start_dt=`date '+%Y-%m-%dT%T.%9N%:z' -d "2024-03-01T09:00:00.000-03:00 +$d days +$h hours +$var_min seconds"`
                genMcc
                if [ $var_type_mcc -lt 2 ]
                then
                    min_amount=20
                    max_amount=80
                elif [ $var_type_mcc -lt 7 ]
                then
                    min_amount=70
                    max_amount=170
                elif [ $var_type_mcc -lt 12 ]
                then
                    min_amount=165
                    max_amount=340
                elif [ $var_type_mcc -lt 15 ]
                then
                    min_amount=335
                    max_amount=530
                elif [ $var_type_mcc -lt 17 ]
                then
                    min_amount=500
                    max_amount=880
                else
                    min_amount=20
                    max_amount=200
                fi

                genModelCard
                genAmount
                genTerm

                var_amount=$((var_amount * fraud_rate))
                echo '---> {"terminal_name":"TERM-'$var_term'","card_number":"'$cc_final'","payment_at":"'${start_dt}'","card_type":"'${arr_type_card[var_type_card]}'","card_model":"'${arr_model_card[var_model_card]}'","currency":"BRL","mcc":"'${arr_mcc[var_type_mcc]}'","amount":'$var_amount',"fraud":'$var_fraud'}'
                curl -X POST $domain --header "Authorization: Bearer $token" --header 'Content-Type: application/json' -d '{"terminal_name":"TERM-'$var_term'","card_number":"'$cc_final'","payment_at":"'${start_dt}'","card_type":"'${arr_type_card[var_type_card]}'","card_model":"'${arr_model_card[var_model_card]}'","currency":"BRL","mcc":"'${arr_mcc[var_type_mcc]}'","amount":'$var_amount',"fraud":'$var_fraud'}'
            done
        done
    done
done