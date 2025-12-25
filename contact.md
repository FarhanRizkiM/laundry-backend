# Contact API Spec

## Create Contact API

Endpoint : POST /api/contacts

Headers :

- Authorization : token

Request Body :

```json
{
  "firstName": "Farhan Rizki",
  "lastName": "Maulana",
  "email": "farhanrizki@gmail.com",
  "phone": "08123456789"
}
```

Response Body Success :

```json
{
  "data": {
    "id": 1,
    "firstName": "Farhan Rizki",
    "lastName": "Maulana",
    "email": "farhanrizki@gmail.com",
    "phone": "08123456789"
  }
}
```

Response Body Error :

```json
{
  "errors": "email is not valid format"
}
```

## Update Contact API

Endpoint : PUT /api/contacts/:id

Headers :

- Authorization : token

Request Body :

```json
{
  "firstName": "Farhan",
  "lastName": "Maulana",
  "email": "farhanrizki@gmail.com",
  "phone": "08123456789"
}
```

Response Body Success :

```json
{
  "data": {
    "id": 1,
    "firstName": "Farhan",
    "lastName": "Maulana",
    "email": "farhanrizki@gmail.com",
    "phone": "08123456789"
  }
}
```

Response Body Error :

```json
{
  "errors": "email is not valid format"
}
```

## Get Contact API

Endpoint : GET /api/contacts/:id

Headers :

- Authorization : token

Response Body Success :

```json
{
  "data": {
    "id": 1,
    "firstName": "Farhan",
    "lastName": "Maulana",
    "email": "farhanrizki@gmail.com",
    "phone": "08123456789"
  }
}
```

Response Body Error :

```json
{
  "errors": "Contact is not found"
}
```

## Seacrh Contact API

Endpoint : GET /api/contacts

Headers :

- Authorization : token

Query Params :

- name : Search by first_name or last_name, using like query, optional
- email : Search by email using like, optional
- phone : Search by phone using like, optional
- page : number of page, default 1
- size : size per page, default 10

Response Body Success :

```json
{
  "data": [
    {
      "id": 1,
      "firstName": "Farhan",
      "lastName": "Maulana",
      "email": "farhanrizki@gmail.com",
      "phone": "08123456789"
    },
    {
      "id": 2,
      "firstName": "Farhan Rizki",
      "lastName": "Maulana",
      "email": "farhanrizki@gmail.com",
      "phone": "08123456789"
    },
  ],
  "paging": {
    "page": 1,
    "total_page": 3,
    "total_item": 30
  }
}
```

Response Body Error :

```json
{
    "errors": "data not found"
}
```

## Remove Contact API

Endpoint : DELETE /api/contacts/:id

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
  "errors": "Contact is not found"
}
```
