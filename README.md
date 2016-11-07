# SaleStock Technical Assignment

## How to run
1. Install MongoDB, make sure it listen on port 27017 on localhost
2. `cd` to the root of repository
3. Run `./seeder/seed`
4. Run `./main`
5. It will listen on port 8080 by default
6. Login using username `administrator` with password `administrator` or username `customer` with password `customer`

## Routes and Params
- /login POST
```
Form:
username
password

Return:
token - to be used as auth, you can put it as form "token", query "token" or auth bearer header
```

- /logout POST
```
Form:
token - your returning token

Return:
success message
```

- /user GET
```
Query:
order_by
page
per_page
id
username
role

Return:
User objects
```

- /user POST
```
Form:
username
role
password

Return:
User object
```

- /user/:id PUT
```
Form:
username
role
password

Return:
User Object
```

- /user/:id DELETE
```
Return:
Deleted message
```

- /product GET
```
Query:
page
per_page
order_by
id
name
quantity_less_than_equal
quantity_more_than_equal
price_less_than_equal
price_more_than_equal

Return:
Product objects
```

- /product POST
```
Form:
name
quantity
price

Return:
Product objects
```

- /product/:id PUT
```
Form:
name
quantity
price

Return:
Product objects
```

- /product/:id DELETE
```
Return:
Deleted message
```

- /coupon GET
```
Query:
page
per_page
order_by
id
valid_until_before
valid_until_after
quantity_less_than_equal
quantity_more_than_equal
discount_less_than_equal
discount_more_than_equal
discount_type

Return:
Coupon objects
```

- /coupon POST
```
Form:
valid_until
quantity
discount
discount_type

Return:
Coupon object
```

- /coupon/:id PUT
```
Form:
valid_until
quantity
discount
discount_type

Return:
Coupon object
```

- /coupon/:id DELETE
```
Return:
Deleted message
```


- /transaction GET
```
Query:
page
per_page
order_by
id

Return:
Transaction objects
```

- /transaction POST
```
Form:
coupon_id
customer_id
products
address
name
email
phone_number

Return:
Transaction object
```

- /transaction/:id PUT
```
Form:
order_status
shipment_id
shipment_status

Return:
Transaction object
```

- /transaction/:id DELETE
```
Return:
Deleted message
```

#### Disclaimer
To be honest, I haven't tested them all thoroughly, so bear with it if it spew some (many) errors.
There are also some questionable practices inside, but the time is too limited to do a thorough checking, so oh well.