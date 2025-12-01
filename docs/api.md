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



#  **Task 模块 API** 

所有任务接口都需要携带 JWT Token
 **Header:**

```
Authorization: <你的 token>
```

## 1. 创建任务（Create Task）

**URL：**

```
POST /tasks
```

**Headers：**

```
Authorization: <JWT Token>
Content-Type: application/json
```

**请求示例：**

```
{
  "title": "学习 Go",
  "desc": "完成 Todo 项目",
  "deadline": "2025-12-31T23:59:59Z"
}
```

**返回示例：**

```
{
  "task": {
    "id": 1,
    "title": "学习 Go",
    "desc": "完成 Todo 项目",
    "deadline": "2025-12-31T23:59:59Z",
    "done": false,
    "user_id": 1
  }
}
```

------

##  2. 获取任务列表（Get Tasks）

**URL：**

```
GET /tasks
```

**Headers：**

```
Authorization: <JWT Token>
```

**返回示例：**

```
{
  "tasks": [
    {
      "id": 1,
      "title": "学习 Go",
      "desc": "完成 Todo 项目",
      "deadline": "2025-12-31T23:59:59Z",
      "done": false,
      "user_id": 1
    },
    {
      "id": 2,
      "title": "健身",
      "desc": "一天 20 分钟",
      "deadline": "2025-12-10T10:00:00Z",
      "done": true,
      "user_id": 1
    }
  ]
}
```

------

## 3. 更新任务（Update Task）

**URL：**

```
PUT /tasks/:id
```

示例：

```
PUT /tasks/1
```

**Headers：**

```
Authorization: <JWT Token>
Content-Type: application/json
```

**请求示例：**

```
{
  "title": "学习 Go（更新）",
  "desc": "每天写代码 2 小时",
  "deadline": "2025-12-25T12:00:00Z",
  "done": true
}
```

**返回示例：**

```
{
  "task": {
    "id": 1,
    "title": "学习 Go（更新）",
    "desc": "每天写代码 2 小时",
    "deadline": "2025-12-25T12:00:00Z",
    "done": true,
    "user_id": 1
  }
}
```

------

4. 删除任务（Delete Task）

**URL：**

```
DELETE /tasks/:id
```

示例：

```
DELETE /tasks/1
```

**Headers：**

```
Authorization: <JWT Token>
```

**返回示例：**

```
{
  "message": "Task deleted"
}
```

------

## 错误响应示例

**Token 无效：**

```
{
  "error": "Invalid token"
}
```

**任务不存在：**

```
{
  "error": "Task not found"
}
```