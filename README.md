# Qonto

## Exercise

At Qonto, we have a wide variety of customers.
Some of them are big organizations that need to perform a large number of transfers.
For example, paying the salaries of hundreds of employees at the end of the month.
Performing those transfers one by one would be painful and time consuming for our customers,
and for that reason we allow them to perform transfers in bulk.
Your mission, should you choose to accept it, is to write a web service to receive
bulk transfer requests from a single Qonto account and:
1. Verify the validity of the request:
* whether the Qonto customer has enough funds for all the transfers in the request.
* If the customer does not have enough funds, the entire request should be denied;
2. If the request must be denied, return a 422 HTTP response;
3. Otherwise, add the transfers to the database, update the customer's balance, and return a 201 HTTP response.

## Answer

### Build and Run

This code is integrated with github ci and docker hub.
If you trigger the pipeline it will compile the code, generate a new docker image and publish it in the docker hub.
The image is already published with the name adrianoccosta/qonto:latest, but fell free to use the toolkit.

To run the application (with docker compose):
1.	clone the code to your local environment:
>  git clone https://github.com/adrianoccosta/exercise-qonto.git
2.	Run:
*  docker-compose up -d
*  docker logs -f qonto-service (to check the loading process)
3.	The application should be up and running. Check the health endpoint: http://127.0.0.1:8080/qonto/api/health

### Test

This application expose a **swagger web page**, where all the available web endpoints can be found and trigger:
http://127.0.0.1:8080/qonto/api/swagger/index.html

Aditionally a postman collection is provided within the project:
>  {base_project}/test/qonto.postman_collection.json


#### Endpoints to run the exercises:

**Bank Account Endpoints**

1. register new bank account
> curl -X POST 'http://127.0.0.1:8080/qonto/api/v1/bank-account' -H 'accept: application/json' -H 'Content-Type: application/json' -d '{ "name": "ACME Corp", "balance": "100000", "iban": "FR10474608000002006107XXXXX", "bic": "OIVUSCLQXXX"}'
 
2. find bank account by its iban
> curl -X GET 'http://127.0.0.1:8080/qonto/api/v1/bank-account/iban/FR10474608000002006107XXXXX' -H 'accept: application/json'

3. update bank account
> curl -X PUT 'http://127.0.0.1:8080/qonto/api/v1/bank-account' -H 'accept: application/json' -H 'Content-Type: application/json' -d '{"name": "ACME Corp", "balance": "123456", "iban": "FR10474608000002006107XXXXX", "bic": "OIVUSCLQXXX"}'

4. delete bank account by its iban
> curl -X DELETE 'http://127.0.0.1:8080/qonto/api/v1/bank-account/iban/FR10474608000002006107XXXXX' -H 'accept: application/json'

**Transaction Endpoints**

1. Get transactions
> curl -X GET 'http://127.0.0.1:8080/qonto/api/v1/transaction' -H 'accept: application/json'

> curl -X GET 'http://127.0.0.1:8080/qonto/api/v1/transaction?counterparty_iban=FR0010009380540930414023042' -H 'accept: application/json'

**Subscriber Information Endpoints**

1. Bulk transfer operation
> curl -X POST 'http://127.0.0.1:8080/qonto/api/v1/transfer/bulk' -H 'accept: application/json' -H 'Content-Type: application/json' -d '{...}'

**Health Endpoints**

1. get metrics
> curl -X POST 'http://127.0.0.1:8080/qonto/api/health' -H 'accept: application/json'

**Metrics Endpoints**

1. get metrics
> curl -X POST 'http://127.0.0.1:8080/qonto/api/metrics' -H 'accept: application/json'

### Notes

This project is implemented with:
* go 1.18
* SQlite
* Swagger
* GORM was not used because is hard to manage complex queries with it
