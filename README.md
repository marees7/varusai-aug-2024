# varusai-aug-2024

# shopping_site


REST API for shopping_site using golang, fiber framework & GORM ORM.


## Overview
- There are 3 different roles, User,Merchant and Admin
- General operations (anyone can view all and specific products)
- User operations (can create and cancel his orders,can view his orders and specific order,can update his details)
- Merchant operations (can create and delete product,can view only his products and specific product,can update his details,product details and order status,can view orders from customers for him and specific order)
- Admin Operations (can create brands and category )


## Features
- User authentication and authorization using JSON Web Tokens (JWT)
- General view for all product without signup
- CRUD operations for products and orders
- Role based validation for every oprations
- Filtering products based on keywords
- Pagination and sorting of products and orders
- Error handling and response formatting
- Input validation and data sanitization
- Database integration using PostgreSQL


## Requirements
- Golang 
- Postgres


## Run Locally

Clone the project

```bash
  git clone https://github.com/marees7/varusai-aug-2024
```

Go to the project directory
go to the cmd folder and main.go file.
change the credentials of postgres db in the internals.

```bash
  go run main.go
```


## API Endpoints

The following endpoints are available in the API:

## AUTH API

| Method | 	Endpoint | 	Description |
| ---- | -------- | -------- |
| POST |	/signup	| Register a new user |
| POST |	/login	| Log in and obtain JWT |

## General API

| Method | 	Endpoint | 	Description |
| ---- | -------- | -------- |
| GET  |	v1/common/product	| Get all products |
| GET  |	v1/common/product/:product_id	| Get a specific product |

## USER API

| Method | 	Endpoint | 	Description |
| ---- | -------- | -------- |
| POST |	v1/user/order	| Create order |
| GET  |	v1/user/order	| Get all orders |
| GET  |	v1/user/order/:order_id	| Get a specific order |
| PATCH |	v1/user/order/:order_id	| Cancel order |
| PATCH |	v1/user/order/:order_id	| Update user details |

## MERCHANT API

| Method | 	Endpoint | 	Description |
| ---- | -------- | -------- |
| POST |	v1/merchant/product	| Create product |
| GET  |	v1/merchant/product	| Get all products from his listing |
| GET  |	v1/merchant/product/:product_id	| Get a specific product from his listing |
| GET  |	v1/merchant/order | Get all orders |
| GET  |	v1/merchant/order/:order_id	| Get a specific order |
| PATCH |	v1/merchant/ | Update merchant details |
| PATCH |	v1/merchant/product	| Update product details |
| PATCH |	v1/merchant/order/:order_id	| Update order status |
| DELETE |	v1/merchant/product/:product_id	| Delete product |

## ADMIN API

| Method | 	Endpoint | 	Description |
| ---- | -------- | -------- |
| POST |	v1/admin/category	| Create category |
| POST |	v1/admin/brand	| Create brand |

## Database Schema

The application uses a PostgreSQL database with the following schema:

```sql
CREATE TABLE IF NOT EXISTS public.users
(
    user_id uuid NOT NULL,
    first_name text COLLATE pg_catalog."default" NOT NULL,
    last_name text COLLATE pg_catalog."default" NOT NULL,
    email text COLLATE pg_catalog."default" NOT NULL,
    phone text COLLATE pg_catalog."default" NOT NULL,
    password text COLLATE pg_catalog."default" NOT NULL,
    role text COLLATE pg_catalog."default" NOT NULL,
    is_verified boolean NOT NULL,
    created_at timestamp with time zone,
    updated_at timestamp with time zone,
    deleted_at timestamp with time zone,
    CONSTRAINT users_pkey PRIMARY KEY (user_id),
    CONSTRAINT uni_users_email UNIQUE (email),
    CONSTRAINT uni_users_phone UNIQUE (phone),
    CONSTRAINT chk_users_role CHECK (role = 'user'::text OR role = 'merchant'::text OR role = 'admin'::text)
)

CREATE TABLE IF NOT EXISTS public.categories
(
    category_id uuid NOT NULL,
    category_name text COLLATE pg_catalog."default" NOT NULL,
    CONSTRAINT categories_pkey PRIMARY KEY (category_id)
)

CREATE TABLE IF NOT EXISTS public.brands
(
    brand_id uuid NOT NULL,
    brand_name text COLLATE pg_catalog."default" NOT NULL,
    CONSTRAINT brands_pkey PRIMARY KEY (brand_id)
)

CREATE TABLE IF NOT EXISTS public.products
(
    product_id uuid NOT NULL,
    product_name text COLLATE pg_catalog."default" NOT NULL,
    category_id uuid NOT NULL,
    brand_id uuid NOT NULL,
    user_id uuid NOT NULL,
    price numeric NOT NULL,
    rating numeric NOT NULL,
    is_approved boolean NOT NULL,
    CONSTRAINT products_pkey PRIMARY KEY (product_id),
    CONSTRAINT fk_brands_product FOREIGN KEY (brand_id)
        REFERENCES public.brands (brand_id) MATCH SIMPLE
        ON UPDATE CASCADE
        ON DELETE SET NULL,
    CONSTRAINT fk_categories_product FOREIGN KEY (category_id)
        REFERENCES public.categories (category_id) MATCH SIMPLE
        ON UPDATE CASCADE
        ON DELETE SET NULL,
    CONSTRAINT fk_users_product FOREIGN KEY (user_id)
        REFERENCES public.users (user_id) MATCH SIMPLE
        ON UPDATE CASCADE
        ON DELETE SET NULL
)

CREATE TABLE IF NOT EXISTS public.orders
(
    order_id uuid NOT NULL,
    user_id uuid NOT NULL,
    address_id uuid NOT NULL,
    name text COLLATE pg_catalog."default" NOT NULL,
    email text COLLATE pg_catalog."default" NOT NULL,
    phone text COLLATE pg_catalog."default" NOT NULL,
    total_amount numeric,
    status text COLLATE pg_catalog."default",
    created_at timestamp with time zone,
    CONSTRAINT orders_pkey PRIMARY KEY (order_id),
    CONSTRAINT fk_addresses_order FOREIGN KEY (address_id)
        REFERENCES public.addresses (address_id) MATCH SIMPLE
        ON UPDATE CASCADE
        ON DELETE SET NULL,
    CONSTRAINT fk_users_order FOREIGN KEY (user_id)
        REFERENCES public.users (user_id) MATCH SIMPLE
        ON UPDATE CASCADE
        ON DELETE SET NULL
)

CREATE TABLE IF NOT EXISTS public.ordered_items
(
    ordered_items_id uuid NOT NULL,
    product_id text COLLATE pg_catalog."default" NOT NULL,
    merchant_id text COLLATE pg_catalog."default" NOT NULL,
    product_name text COLLATE pg_catalog."default" NOT NULL,
    quantity bigint NOT NULL,
    price numeric NOT NULL,
    status text COLLATE pg_catalog."default",
    order_id uuid,
    customer_id text COLLATE pg_catalog."default",
    address_id text COLLATE pg_catalog."default",
    created_at timestamp with time zone,
    CONSTRAINT ordered_items_pkey PRIMARY KEY (ordered_items_id),
    CONSTRAINT fk_orders_products FOREIGN KEY (order_id)
        REFERENCES public.orders (order_id) MATCH SIMPLE
        ON UPDATE CASCADE
        ON DELETE SET NULL
)
```

## Sample API Requests and Responses


## Auth API

#### POST /signup

sample request:
```json
       {
        "first_name" : "varusai",
        "last_name": "mohamed",
        "email":"varusaimohamed@gmail.com",
        "phone":"9715538624",
        "password":"varusai@123",
        "role":"merchant",
        "address":[
            {
                "door_no": "15/144",
                "street": "bigcar",
                "city": "madurai",
                "state": "tamilnadu",
                "zip_code" : 625005
            }
        ]
    }
```

sample response:
```json
{
    "message": "user created successfully"
}
```

##### POST /login

sample request:
```json
     {
        "email":"varusaimohamed@gmail.com",
        "password":"varusai@123"
    }
```

sample response:
```json
{
    "message": "logged in successfully"
}
```

## General API

#### GET v1/common/product?category_name=phone&brand_name=xiomi&limit=10&offset=1&price=300&rating=3

sample response:
```json
{
    "data": [
        {
            "product_id": "de0f1048-6624-4bb4-97b4-9ad705c70caa",
            "product_name": "redmi_4",
            "category_id": "5eb1d5c6-e7b0-49e1-a715-3ea1b1650835",
            "brand_id": "dd12b1f2-6d8a-4fd6-b335-40dba068ca1c",
            "user_id": "1ba51d85-bde6-4f12-b153-92e9233a0492",
            "price": 9000,
            "rating": 4.6
        },
        {
            "product_id": "6a7ef81c-60a4-4b63-838c-af7415c3134e",
            "product_name": "xiomi_notepat_7",
            "category_id": "5eb1d5c6-e7b0-49e1-a715-3ea1b1650835",
            "brand_id": "dd12b1f2-6d8a-4fd6-b335-40dba068ca1c",
            "user_id": "3efb583d-cc54-41ce-954f-440ea67fbae8",
            "price": 32000,
            "rating": 4.6
        }
    ],
    "total_records": 2,
    "limit": 10
}
```

#### GET v1/common/product/de0f1048-6624-4bb4-97b4-9ad705c70caa

sample response:
```json
{
    "data": {
        "product_id": "de0f1048-6624-4bb4-97b4-9ad705c70caa",
        "product_name": "redmi_4",
        "category_id": "5eb1d5c6-e7b0-49e1-a715-3ea1b1650835",
        "brand_id": "dd12b1f2-6d8a-4fd6-b335-40dba068ca1c",
        "user_id": "1ba51d85-bde6-4f12-b153-92e9233a0492",
        "price": 9000,
        "rating": 4.6
    }
}
```

## Merchant API

#### POST v1/merchant/product

sample request:
```json
{
    "product_name":"s1_pro",
    "category_id":"5eb1d5c6-e7b0-49e1-a715-3ea1b1650835",
    "brand_id":"fe4ad93a-d17e-4240-a6ab-d421713c49cc",
    "price":58000,
    "rating":4.2
}
```

sample response:
```json
{
    "message": "product created successfully",
    "data": {
        "product_id": "81b196ce-bb2d-4c47-836b-7bbab8a1ab7e",
        "product_name": "s1_pro",
        "category_id": "5eb1d5c6-e7b0-49e1-a715-3ea1b1650835",
        "brand_id": "fe4ad93a-d17e-4240-a6ab-d421713c49cc",
        "user_id": "3efb583d-cc54-41ce-954f-440ea67fbae8",
        "price": 58000,
       "rating":4.2
    }
}
```

#### GET v1/merchant/product?category_name=phone&brand_name=xiomi&limit=3&offset=1

sample response:
```json
{
    "data": [
        {
            "product_id": "6a7ef81c-60a4-4b63-838c-af7415c3134e",
            "product_name": "xiomi_notepat_7",
            "category_id": "5eb1d5c6-e7b0-49e1-a715-3ea1b1650835",
            "brand_id": "dd12b1f2-6d8a-4fd6-b335-40dba068ca1c",
            "user_id": "3efb583d-cc54-41ce-954f-440ea67fbae8",
            "price": 32000,
            "rating": 4.6
        }
    ],
    "total_records": 1,
    "limit": 3
}
```

#### GET v1/merchant/product/4fac47bc-beed-4f17-810e-c1a9c1aece83

sample response:
```json
{
    "data": {
        "product_id": "4fac47bc-beed-4f17-810e-c1a9c1aece83",
        "product_name": "s1_pro",
        "category_id": "5eb1d5c6-e7b0-49e1-a715-3ea1b1650835",
        "brand_id": "fe4ad93a-d17e-4240-a6ab-d421713c49cc",
        "user_id": "3efb583d-cc54-41ce-954f-440ea67fbae8",
        "price": 58000
    }
}
```

#### GET v1/merchant/order?from_date=2024-12-30&to_date=2024-12-31&limit=3&offset=1

sample response:
```json
{
    "data": [
        {
            "ordered_items_id": "115a4a62-6540-4bcf-8ee0-e814ebe07146",
            "product_id": "21649249-0fb8-43ca-85b1-ca788709315b",
            "merchant_id": "3efb583d-cc54-41ce-954f-440ea67fbae8",
            "product_name": "walk_smooth_002",
            "quantity": 1,
            "price": 700,
            "status": "placed",
            "order_id": "d80d9a24-2130-4f03-900a-2e474ab72d43",
            "customer_id": "f96d2ba9-0d6f-4d70-b53b-709f4833889d",
            "address_id": "719e5b93-d27c-4ab0-b21d-f2c459b5a078",
            "created_at": "2024-12-30T17:37:28.441892+05:30"
        },
        {
            "ordered_items_id": "605912da-d958-4353-bf6b-1a6d97da8492",
            "product_id": "4fac47bc-beed-4f17-810e-c1a9c1aece83",
            "merchant_id": "3efb583d-cc54-41ce-954f-440ea67fbae8",
            "product_name": "s1_pro",
            "quantity": 1,
            "price": 58000,
            "status": "placed",
            "order_id": "d80d9a24-2130-4f03-900a-2e474ab72d43",
            "customer_id": "f96d2ba9-0d6f-4d70-b53b-709f4833889d",
            "address_id": "719e5b93-d27c-4ab0-b21d-f2c459b5a078",
            "created_at": "2024-12-30T17:37:28.441892+05:30"
        },
        {
            "ordered_items_id": "c52354af-efd3-447c-a152-8e4160ae627e",
            "product_id": "21649249-0fb8-43ca-85b1-ca788709315b",
            "merchant_id": "3efb583d-cc54-41ce-954f-440ea67fbae8",
            "product_name": "walk_smooth_002",
            "quantity": 4,
            "price": 700,
            "status": "placed",
            "order_id": "01976146-8844-4f18-8e2d-88ea4a01b3e3",
            "customer_id": "f96d2ba9-0d6f-4d70-b53b-709f4833889d",
            "address_id": "719e5b93-d27c-4ab0-b21d-f2c459b5a078",
            "created_at": "2024-12-30T17:37:35.460546+05:30"
        }
    ],
    "total_records": 7,
    "limit": 3
}
```

#### GET v1/merchant/order/04ee0160-0ad9-42af-a682-041f9334846c

sample response:
```json
{
    "data": [
        {
            "ordered_items_id": "113157e7-74f0-486a-b4ea-54d0ec033052",
            "product_id": "21649249-0fb8-43ca-85b1-ca788709315b",
            "merchant_id": "3efb583d-cc54-41ce-954f-440ea67fbae8",
            "product_name": "walk_smooth_002",
            "quantity": 1,
            "price": 700,
            "status": "placed",
            "order_id": "04ee0160-0ad9-42af-a682-041f9334846c",
            "customer_id": "f96d2ba9-0d6f-4d70-b53b-709f4833889d",
            "address_id": "719e5b93-d27c-4ab0-b21d-f2c459b5a078",
            "created_at": "2024-12-30T17:37:57.092258+05:30"
        },
        {
            "ordered_items_id": "fc0b1f1d-255b-4fb8-887a-c8a35d06f532",
            "product_id": "4fac47bc-beed-4f17-810e-c1a9c1aece83",
            "merchant_id": "3efb583d-cc54-41ce-954f-440ea67fbae8",
            "product_name": "s1_pro",
            "quantity": 2,
            "price": 58000,
            "status": "placed",
            "order_id": "04ee0160-0ad9-42af-a682-041f9334846c",
            "customer_id": "f96d2ba9-0d6f-4d70-b53b-709f4833889d",
            "address_id": "719e5b93-d27c-4ab0-b21d-f2c459b5a078",
            "created_at": "2024-12-30T17:37:57.092258+05:30"
        }
    ]
}
```

#### PATCH v1/merchant/product

sample request:
```json
{
    "product_id": "81b196ce-bb2d-4c47-836b-7bbab8a1ab7e",
    "product_name":"s1_pro_ultra",
    "price":66000
}
```

sample response:
```json
{
    "message": "product updated successfully",
    "data": {
        "product_id": "81b196ce-bb2d-4c47-836b-7bbab8a1ab7e"
    }
}
```

#### PATCH v1/merchant

sample request:
```json
{
   "first_name":"varusai",
   "last_name":"s",
   "email":"varusai652000@gmail.com",
   "phone":"8220902268",
   "password":"varusai@1234",
    "address":[
            {
                "address_id":"79c19b9c-6cf0-4ef2-a79c-be23489eb704",
                "door_no": "12/22",
                "street": "sanathi",
                "city": "ulundhoorpettai",
                "state": "tamilnadu",
                "zip_code" : 625006
            }
        ]
}
```

sample response:
```json
{
    "message": "merchant details updated successfully",
    "data": {
        "user_id": "3efb583d-cc54-41ce-954f-440ea67fbae8"
    }
}
```

#### PATCH v1/merchant/order/91fe8d8a-2a5e-421b-8370-3837e3054389?order_status=shipped

sample response:
```json
{
    "message": "order_item status updated successfully",
    "data": {
        "order_item_id": "91fe8d8a-2a5e-421b-8370-3837e3054389"
    }
}
```

#### DELETE v1/merchant/product/productId

sample response:
```json
{
    "message": "product deleted successfully"
}
```

## User API

#### POST v1/user/order

sample request:
```json
{
"address_id":"719e5b93-d27c-4ab0-b21d-f2c459b5a078",
"products":[{
    "product_Id":"21649249-0fb8-43ca-85b1-ca788709315b",
    "quantity":1
},
{
    "product_Id":"4fac47bc-beed-4f17-810e-c1a9c1aece83",
    "quantity":2
}]
}
```

sample response:
```json
{
    "message": "order placed successfully",
    "data": {
        "ordered_id": "eb395ef7-c6d1-4c46-a73c-f8ee5bb13983",
        "user_id": "f96d2ba9-0d6f-4d70-b53b-709f4833889d",
        "address_id": "719e5b93-d27c-4ab0-b21d-f2c459b5a078",
        "name": "mohaideen v",
        "first_name": "mohaideen@gmail.com",
        "phone": "9033115452",
        "products": [
            {
                "ordered_items_id": "1288a400-51c6-4388-97b7-88785f0a24e0",
                "product_id": "21649249-0fb8-43ca-85b1-ca788709315b",
                "merchant_id": "3efb583d-cc54-41ce-954f-440ea67fbae8",
                "product_name": "walk_smooth_002",
                "quantity": 1,
                "price": 700,
                "order_id": "eb395ef7-c6d1-4c46-a73c-f8ee5bb13983",
                "customer_id": "f96d2ba9-0d6f-4d70-b53b-709f4833889d",
                "address_id": "719e5b93-d27c-4ab0-b21d-f2c459b5a078",
                "created_at": "2024-12-30T21:02:19.503+05:30"
            },
            {
                "ordered_items_id": "00377539-549a-4dd3-969e-29bbb7cf69ee",
                "product_id": "4fac47bc-beed-4f17-810e-c1a9c1aece83",
                "merchant_id": "3efb583d-cc54-41ce-954f-440ea67fbae8",
                "product_name": "s1_pro",
                "quantity": 2,
                "price": 58000,
                "order_id": "eb395ef7-c6d1-4c46-a73c-f8ee5bb13983",
                "customer_id": "f96d2ba9-0d6f-4d70-b53b-709f4833889d",
                "address_id": "719e5b93-d27c-4ab0-b21d-f2c459b5a078",
                "created_at": "2024-12-30T21:02:19.503+05:30"
            }
        ],
        "total_amount": 116700,
        "status": "inprogress",
        "created_at": "2024-12-30T21:02:19.5009559+05:30"
    }
}
```

#### GET v1/user/order?from_date=2024-12-31&to_date=2025-12-31&limit=10&offset=1

sample response:
```json
{
    "data": [
        {
            "ordered_id": "a3c86b72-5270-44ab-b8c2-83c6a1e0e855",
            "user_id": "f96d2ba9-0d6f-4d70-b53b-709f4833889d",
            "address_id": "719e5b93-d27c-4ab0-b21d-f2c459b5a078",
            "name": "mohamed mohaideen",
            "first_name": "mohaideen@gmail.com",
            "phone": "9092576225",
            "products": [
                {
                    "ordered_items_id": "0f028fea-ac3a-4de8-8b53-99261058fd49",
                    "product_id": "21649249-0fb8-43ca-85b1-ca788709315b",
                    "merchant_id": "3efb583d-cc54-41ce-954f-440ea67fbae8",
                    "product_name": "walk_smooth_002",
                    "quantity": 1,
                    "price": 700,
                    "status": "placed",
                    "order_id": "a3c86b72-5270-44ab-b8c2-83c6a1e0e855",
                    "customer_id": "f96d2ba9-0d6f-4d70-b53b-709f4833889d",
                    "address_id": "719e5b93-d27c-4ab0-b21d-f2c459b5a078",
                    "created_at": "2024-12-31T11:14:48.228936+05:30"
                },
                {
                    "ordered_items_id": "48d75d54-eb89-4d43-9aca-a04ace14dd4c",
                    "product_id": "4fac47bc-beed-4f17-810e-c1a9c1aece83",
                    "merchant_id": "3efb583d-cc54-41ce-954f-440ea67fbae8",
                    "product_name": "s1_pro",
                    "quantity": 2,
                    "price": 58000,
                    "status": "placed",
                    "order_id": "a3c86b72-5270-44ab-b8c2-83c6a1e0e855",
                    "customer_id": "f96d2ba9-0d6f-4d70-b53b-709f4833889d",
                    "address_id": "719e5b93-d27c-4ab0-b21d-f2c459b5a078",
                    "created_at": "2024-12-31T11:14:48.228936+05:30"
                }
            ],
            "total_amount": 116700,
            "status": "inprogress",
            "created_at": "2024-12-31T11:14:48.228936+05:30"
        },
        {
            "ordered_id": "7ca449ba-ae0a-4cb5-9a57-d1c9fe1d3f6c",
            "user_id": "f96d2ba9-0d6f-4d70-b53b-709f4833889d",
            "address_id": "719e5b93-d27c-4ab0-b21d-f2c459b5a078",
            "name": "mohamed mohaideen",
            "first_name": "mohaideen@gmail.com",
            "phone": "9092576225",
            "products": [
                {
                    "ordered_items_id": "a14f6bbf-2fa2-466a-9e8b-4c0820984fec",
                    "product_id": "21649249-0fb8-43ca-85b1-ca788709315b",
                    "merchant_id": "3efb583d-cc54-41ce-954f-440ea67fbae8",
                    "product_name": "walk_smooth_002",
                    "quantity": 1,
                    "price": 700,
                    "status": "placed",
                    "order_id": "7ca449ba-ae0a-4cb5-9a57-d1c9fe1d3f6c",
                    "customer_id": "f96d2ba9-0d6f-4d70-b53b-709f4833889d",
                    "address_id": "719e5b93-d27c-4ab0-b21d-f2c459b5a078",
                    "created_at": "2024-12-31T11:37:19.043075+05:30"
                },
                {
                    "ordered_items_id": "32bcdf5f-82ee-46bc-a35d-0eb990772e88",
                    "product_id": "4fac47bc-beed-4f17-810e-c1a9c1aece83",
                    "merchant_id": "3efb583d-cc54-41ce-954f-440ea67fbae8",
                    "product_name": "s1_pro",
                    "quantity": 2,
                    "price": 58000,
                    "status": "placed",
                    "order_id": "7ca449ba-ae0a-4cb5-9a57-d1c9fe1d3f6c",
                    "customer_id": "f96d2ba9-0d6f-4d70-b53b-709f4833889d",
                    "address_id": "719e5b93-d27c-4ab0-b21d-f2c459b5a078",
                    "created_at": "2024-12-31T11:37:19.043075+05:30"
                }
            ],
            "total_amount": 116700,
            "status": "inprogress",
            "created_at": "2024-12-31T11:37:19.043075+05:30"
        }
    ],
    "total_records": 2,
    "limit": 10
}
```

#### GET v1/user/order/7d287d9c-aee1-4b27-ac34-45513a69e977

sample response:
```json
{
    "data": {
        "ordered_id": "7d287d9c-aee1-4b27-ac34-45513a69e977",
        "user_id": "f96d2ba9-0d6f-4d70-b53b-709f4833889d",
        "address_id": "719e5b93-d27c-4ab0-b21d-f2c459b5a078",
        "name": "mohaideen v",
        "first_name": "mohaideen@gmail.com",
        "phone": "9033115452",
        "products": [
            {
                "ordered_items_id": "97657edc-773e-4703-89fe-e48d52a583d7",
                "product_id": "21649249-0fb8-43ca-85b1-ca788709315b",
                "merchant_id": "3efb583d-cc54-41ce-954f-440ea67fbae8",
                "product_name": "walk_smooth_002",
                "quantity": 1,
                "price": 700,
                "status": "placed",
                "order_id": "7d287d9c-aee1-4b27-ac34-45513a69e977",
                "customer_id": "f96d2ba9-0d6f-4d70-b53b-709f4833889d",
                "address_id": "719e5b93-d27c-4ab0-b21d-f2c459b5a078",
                "created_at": "2024-12-30T17:38:17.189898+05:30"
            }
        ],
        "total_amount": 700,
        "status": "inprogress",
        "created_at": "2024-12-30T17:38:17.189898+05:30"
    }
}
```

#### PATCH v1/user/order/eb395ef7-c6d1-4c46-a73c-f8ee5bb13983

sample response:
```json
{
    "message": "order cancelled successfully"
}
```

#### PATCH v1/user

sample request:
```json
{
   "first_name":"mohamed",
   "last_name":"mohaideen",
   "email":"mohaideen@gmail.com",
   "phone":"9092576225",
   "password":"mohaideen@123",
    "address":[
            {
                "address_id":"719e5b93-d27c-4ab0-b21d-f2c459b5a079",
                "door_no": "15/122",
                "street": "periyaratham",
                "city": "madurai",
                "state": "tamilnadu",
                "zip_code" : 625005
            }
        ]
}
```

sample response:
```json
{
    "message": "user details updated successfully",
    "data": {
        "user_id": "f96d2ba9-0d6f-4d70-b53b-709f4833889d"
    }
}
```