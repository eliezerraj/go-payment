#!/bin/bash

var_cc=0
genCC(){
    var_cc=$(($RANDOM%($max_cc-$min_cc+1)+$min_cc))
}

var_term=0
genTerm(){
    var_term=$(($RANDOM%($max_term-$min_term+1)+$min_term))
}

# Normal ditribuition
var_amount=0
genAmount(){
    a=$(($RANDOM%($max_amount-$min_amount+1)+$min_amount))
    b=$(($RANDOM%($max_amount-$min_amount+1)+$min_amount))
    c=$(($RANDOM%($max_amount-$min_amount+1)+$min_amount))
    var_amount=$(( (a+b+c)/3 ))
}

var_type_mcc=0
genMcc(){
    var_type_mcc=$(($RANDOM%($max_mcc-$min_mcc+1)+$min_mcc))
}

var_model_card=0
genModelCard(){
    var_model_card=$(($RANDOM%($max_model-$min_model+1)+$min_model))
}

var_type_card=0
genTypeCard(){
    var_type_card=$(($RANDOM%($max_tcc-$min_tcc+1)+$min_tcc))
}

var_tx_per_day=0
genTXDay(){
    var_tx_per_day=$(($RANDOM%($max_tx_day-$min_tx_day+1)+$min_tx_day))
}

var_min=0
genMinutes(){
    var_min=$(($RANDOM%($max_minutes-$min_minutes+1)+$min_minutes))
}

declare -a arr_model_card
arr_model_card=(CHIP VIRTUAL)

declare -a arr_type_card
arr_type_card=(CREDIT DEBIT)

declare -a arr_mcc
arr_mcc=(PARKING BEVERAGE FOOD LAUNDRY CINEMA BOOK GIFT CASH GAS PET DRUG_STORE COSMETIC GYM STORE SPORTING COMPUTER MOTOR)
# -----------------------------------------------------

#domain=http://localhost:5007/payment/pay
#domain=https://97x38r33ag.execute-api.us-east-2.amazonaws.com/Live/payment/pay

token=eyJhbGciOiJIUzI1NiIsInR5cCI6IkpXVCJ9.eyJ1c2VybmFtZSI6ImFkbWluIiwic2NvcGUiOlsiYWRtaW4iXSwiZXhwIjoxNzA5MDgzNDE2fQ.TnJ9WrmbIy3rVKCq9TJ7-rstl9-1Uza2wSUXth13EWk
domain=https://go-api-global.architecture.caradhras.io/payment/payment/pay

min_model=0
max_model=1

min_tcc=0
max_tcc=1

min_cc=1
max_cc=20

min_amount=20
max_amount=800

min_mcc=0
max_mcc=16

min_term=1
max_term=100

min_tx_day=0
max_tx_day=4

min_minutes=1
max_minutes=600

var_pan=111111000001

echo "-------------------------------------"
echo "-----------STARTING DAY---------------"
echo "-------------------------------------"

arr_mcc=(PARKING BEVERAGE FOOD LAUNDRY CINEMA BOOK GIFT CASH GAS PET DRUG_STORE COSMETIC GYM STORE SPORTING COMPUTER MOTOR)
var_fraud=0
fraud_rate=1
min_tx_day=0
max_tx_day=4

for (( d=0; d<30; d++ )) # Day
do
    echo "***********> New Day ****************"

    for (( w=0; w<30; w++ )) #Qtd tx per minutes
    do
        echo "* * * * New CC * * * *"
        #genCC  #Generate CC

        #genTypeCard #Generate type card (CREDIT or DEBIT)           
        #if [[ ${arr_type_card[var_type_card]} == "DEBIT" ]]
        #then
        #    var_model_card=0
        #    var_pan=333111000000
        #else
        #    genModelCard
        #    if [[ ${arr_model_card[var_model_card]} == "CHIP" ]]
        #        var_pan=1111111000000
        #    then
        #        var_pan=222111000000
        #    fi
        #fi
        #cc=$(($var_pan+$var_cc))
        cc=$(($var_pan+$w))
        cc_final="${cc:0:3}"."${cc:3:3}"."${cc:6:3}"."${cc:9:3}"

        genTXDay #Generate tx per day
        for (( z=var_tx_per_day; z>0; z-- ))
        do
            genMinutes
            start_dt=`date '+%Y-%m-%dT%T.%9N%:z' -d "2024-02-01T09:00:00.000-03:00 +$d days +$var_min minutes"`

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

            genAmount
            genTerm
            
            var_amount=$((var_amount * fraud_rate))
            
            #echo '{"terminal_name":"TERM-'$var_term'","card_number":"'$cc_final'","payment_at":"'${start_dt}'","card_type":"'${arr_type_card[var_type_card]}'","card_model":"'${arr_model_card[var_model_card]}'","currency":"BRL","mcc":"'${arr_mcc[var_type_mcc]}'","amount":'$var_amount',"fraud":'$var_fraud'}'
            #echo  curl -X POST $domain --header "Authorization: Bearer $token" --header 'Content-Type: application/json' -d '{"terminal_name":"TERM-'$var_term'","card_number":"'$cc_final'","payment_at":"'${start_dt}'","card_type":"'${arr_type_card[var_type_card]}'","card_model":"'${arr_model_card[var_model_card]}'","currency":"BRL","mcc":"'${arr_mcc[var_type_mcc]}'","amount":'$var_amount'}'
            curl -X POST $domain --header "Authorization: Bearer $token" --header 'Content-Type: application/json' -d '{"terminal_name":"TERM-'$var_term'","card_number":"'$cc_final'","payment_at":"'${start_dt}'","card_type":"'${arr_type_card[var_type_card]}'","card_model":"'${arr_model_card[var_model_card]}'","currency":"BRL","mcc":"'${arr_mcc[var_type_mcc]}'","amount":'$var_amount',"fraud":'$var_fraud'}'
        done
    done
done

echo "-------------------------------------"
echo "-----------STARTING Nigth---------------"
echo "-------------------------------------"

arr_mcc=(PARKING BEVERAGE FOOD CINEMA STORE)
min_mcc=0
max_mcc=4

var_fraud=0
fraud_rate=1
min_tx_day=0
max_tx_day=2
#max_tx_day=2

for (( d=0; d<30; d++ )) # Day
do
    echo "***********> New Day ****************"

    for (( w=0; w<30; w++ )) #Qtd de CC
    do
        echo "* * * * New CC * * * *"
        #genCC  #Generate CC

        #genTypeCard #Generate type card (CREDIT or DEBIT)           
        #if [[ ${arr_type_card[var_type_card]} == "DEBIT" ]]
        #then
        #    var_model_card=0
        #    var_pan=333111000000
        #else
        #    genModelCard
        #    if [[ ${arr_model_card[var_model_card]} == "CHIP" ]]
        #        var_pan=1111111000000
        #    then
        #        var_pan=222111000000
        #    fi
        #fi
        #cc=$(($var_pan+$var_cc))
        cc=$(($var_pan+$w))
        cc_final="${cc:0:3}"."${cc:3:3}"."${cc:6:3}"."${cc:9:3}"

        genTXDay #Generate tx per day
        for (( z=var_tx_per_day; z>0; z-- ))
        do
            genMinutes
            start_dt=`date '+%Y-%m-%dT%T.%9N%:z' -d "2024-02-01T20:00:00.000-03:00 +$d days +$var_min minutes"`

            genMcc
            if [ $var_type_mcc -lt 1 ]
            then
                min_amount=20
                max_amount=45
            elif [ $var_type_mcc -lt 2 ]
            then
                min_amount=40
                max_amount=100
            elif [ $var_type_mcc -lt 4 ]
            then
                min_amount=100
                max_amount=250
            elif [ $var_type_mcc -lt 5 ]
            then
                min_amount=200
                max_amount=500
            else
                min_amount=20
                max_amount=200
            fi

            genAmount
            genTerm

            var_amount=$((var_amount * fraud_rate))

            #echo '{"terminal_name":"TERM-'$var_term'","card_number":"'$cc_final'","payment_at":"'${start_dt}'","card_type":"'${arr_type_card[var_type_card]}'","card_model":"'${arr_model_card[var_model_card]}'","currency":"BRL","mcc":"'${arr_mcc[var_type_mcc]}'","amount":'$var_amount',"fraud":'$var_fraud'}'
            #echo  curl -X POST $domain --header "Authorization: Bearer $token" --header 'Content-Type: application/json' -d '{"terminal_name":"TERM-'$var_term'","card_number":"'$cc_final'","payment_at":"'${start_dt}'","card_type":"'${arr_type_card[var_type_card]}'","card_model":"'${arr_model_card[var_model_card]}'","currency":"BRL","mcc":"'${arr_mcc[var_type_mcc]}'","amount":'$var_amount'}'
            curl -X POST $domain --header "Authorization: Bearer $token" --header 'Content-Type: application/json' -d '{"terminal_name":"TERM-'$var_term'","card_number":"'$cc_final'","payment_at":"'${start_dt}'","card_type":"'${arr_type_card[var_type_card]}'","card_model":"'${arr_model_card[var_model_card]}'","currency":"BRL","mcc":"'${arr_mcc[var_type_mcc]}'","amount":'$var_amount',"fraud":'$var_fraud'}'
        done
    done
done