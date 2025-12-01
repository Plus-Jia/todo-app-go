## 🔖 Auth API

------

###  **POST /register**

用户注册

#### Request Body

```
{
  "email": "string",
  "password": "string"
}
```

#### Response

**200 OK**

```
{
  "message": "User registered successfully"
}
```

**400 Bad Request**

```
{
  "error": "Email is already registered"
}
```

------

###  **POST /login**

用户登录并返回 JWT

#### Request Body

```
{
  "email": "string",
  "password": "string"
}
```

#### Response

**200 OK**

```
{
  "token": "jwt-string"
}
```

**401 Unauthorized**

```
{
  "error": "Invalid email or password"
}
```

------

### ### **GET /protected**

需要 Token

#### Headers

```
Authorization: Bearer <token>
```

#### Response

```
{
  "message": "This is a protected route"
}
```