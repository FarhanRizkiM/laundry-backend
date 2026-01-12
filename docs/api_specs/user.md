# User API Spec

## Register User API

Endpoint : POST /api/users

Request Body : 

```json
{
    "username": "rizkiganteng",
    "password": "rizki123",
    "name": "Farhan Rizki"
}
```

Response Body Success :

```json
{
    "data"{
        "username": "rizkiganteng",
        "name": "Farhan Rizki"
    }
}
```

Response Body Error :

```json
{
    "errors": "Username already registered"
}
```

## Login User API

Endpoint : POST /api/users/login

Request Body : 

```json
{
    "username": "rizkiganteng",
    "password": "rizki123"
}
```

Response Body Success : 

```json
{
    "data": {
        "token": "unique-token"
    }
}
```

Response Body Error : 

```json
{
    "errors": "Username or password wrong"
}
```

## Update User API

Endpoint : PATCH /api/users/current

Headers : 
- Authorization : token

Request Body : 

```json
{
    "name": "Farhan Rizki Maulana", //optional
    "password": "maulana123" //optional
}
```

Response Body Success : 

```json
{
    "data": {
        "username": "rizkiganteng", 
        "name": "Farhan Rizki Maulana"
    }
}
```

Response Body Error : 

```json
{
    "errors": "Name length max 100"
}
```

## Get User API

Endpoint : GET /api/users/current

Headers : 
- Authorization : token

Response Body Success : 

```json
{
    "data": {
        "username": "rizkiganteng",
        "name": "Farhan Rizki Maulana"
    }
}
```

Response Body Error :

```json
{
    "errors": "Unauthorized"
}
```

## Logout User API

Endpoint : DELETE /api/users/logout

Headers : 
- Authorization : token

Response Body Success : 

```json
{
    "data": "OK"
}
```

Response Body Error : 

```json
{
    "errors": "Unauthorized"
}
```