# Waza - Senior Engineer Take Home Coding Assessment


### About Service
The program is written in go.

| Service Name | API Type | Description                                                                                                                                                               |
|--------------|----------|---------------------------------------------------------------------------------------------------------------------------------------------------------------------------|
| waza         | GraphQL  | This api uses GraphQL as it's input & output contract layer, <br/> It also gives you a nice playground to test the api on https://localhost:[port] - default port is 4000 | 

### Internal Dependencies

- GraphQL/GQLGEN (already setup) ( https://gqlgen.com/ )
- Sqlite for embedded db

### Environment Variables

```dotenv
ENVIRONMENT=local
PORT=4000 # This can be channged to what you desire
```

### Useful run scripts

```shell script
# To generate server code

$ make schema
```

```shell script
# To start it in docker

$ make docker ( this fails because of sqlite. but a binary has been built )
```

```shell script
# To run it manually

$ go run main.go 
```

```shell script
# For local development

$ make local
```

```shell script
# to format code

$ make fmt 
```

### Simple Documentation

| Methods               | Description                                                           |
|-----------------------|-----------------------------------------------------------------------|
| CreateUser            | This creates a user and automatically creates a transactional account |
| GetUserById           | To retrieve a user by id                                              |
| GetUserByPhone        | To retrieve a user by phone                                           |
| GetUserByEmail        | To retrieve a user by email (if provided)                             |
| TransferFunds         | To send an `amount` of funds `fromAccountId`, `toAccountId`           |
| GetAccountById        | To retrieve a transactional account by id                             |
| GetAccountByOwnerId   | To retrieve a transactional account by owner id                       |
| GetTransactionHistory | To retrieve all transaction history for `accountId`                   |

