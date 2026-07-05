# API Reference

## Base URL
```
http://localhost:8080/api
```

## Health Endpoints

### Check Application Health
```
GET /health
```

**Response:**
```json
{
  "status": "ok"
}
```

### Check Database Connection
```
GET /health/db
```

**Response:**
```json
{
  "status": "ok"
}
```

---

## Scenario Endpoints

### Create Scenario
```
POST /scenario
```

**Request Body:**
```json
{
  "title": "Смена паспорта",
  "description": "Процесс смены паспорта РФ",
  "category": "документы",
  "start_node_id": 1
}
```

**Response (201 Created):**
```json
{
  "id": 1
}
```

### Get Scenario by ID
```
GET /scenario?id={id}
```

**Query Parameters:**
- `id` (required): Scenario ID

**Response:**
```json
{
  "id": 1,
  "title": "Смена паспорта",
  "description": "Процесс смены паспорта РФ",
  "category": "документы",
  "start_node_id": 1,
  "created_at": "2026-07-05T15:00:00Z"
}
```

### List All Scenarios
```
GET /scenarios
```

**Response:**
```json
[
  {
    "id": 1,
    "title": "Смена паспорта",
    "description": "Процесс смены паспорта РФ",
    "category": "документы",
    "start_node_id": 1,
    "created_at": "2026-07-05T15:00:00Z"
  }
]
```

### Update Scenario
```
PUT /scenario?id={id}
```

**Query Parameters:**
- `id` (required): Scenario ID

**Request Body:**
```json
{
  "title": "Смена паспорта (обновлено)",
  "description": "Процесс смены паспорта РФ",
  "category": "документы",
  "start_node_id": 1
}
```

**Response:**
```json
{
  "status": "ok"
}
```

### Delete Scenario
```
DELETE /scenario?id={id}
```

**Query Parameters:**
- `id` (required): Scenario ID

**Response:**
```json
{
  "status": "ok"
}
```

---

## Node Endpoints

### Create Node
```
POST /node
```

**Request Body:**
```json
{
  "scenario_id": 1,
  "type": "question",
  "text": "Вы потеряли паспорт?",
  "order": 1
}
```

**For Action Nodes:**
```json
{
  "scenario_id": 1,
  "type": "action",
  "text": "Уведомить паспортный стол",
  "action_type": "notification",
  "organization": "ФМС России",
  "deadline": "30 дней",
  "consequences": "Штраф до 2500 рублей",
  "reference_links": "[\"https://example.com\"]",
  "order": 2
}
```

**Response:**
```json
{
  "id": 1
}
```

### Get Node by ID
```
GET /node?id={id}
```

**Response:**
```json
{
  "id": 1,
  "scenario_id": 1,
  "type": "question",
  "text": "Вы потеряли паспорт?",
  "order": 1,
  "created_at": "2026-07-05T15:00:00Z"
}
```

### Get Scenario Nodes
```
GET /nodes?scenario_id={id}
```

**Response:**
```json
[
  {
    "id": 1,
    "scenario_id": 1,
    "type": "question",
    "text": "Вы потеряли паспорт?",
    "order": 1,
    "created_at": "2026-07-05T15:00:00Z"
  }
]
```

### Update Node
```
PUT /node?id={id}
```

**Response:**
```json
{
  "status": "ok"
}
```

### Delete Node
```
DELETE /node?id={id}
```

**Response:**
```json
{
  "status": "ok"
}
```

---

## Edge Endpoints

### Create Edge
```
POST /edge
```

**Request Body:**
```json
{
  "from_node": 1,
  "to_node": 2,
  "condition": "yes",
  "logic": null
}
```

**Response:**
```json
{
  "id": 1
}
```

### Get Edge by ID
```
GET /edge?id={id}
```

### Get Scenario Edges
```
GET /edges?scenario_id={id}
```

**Response:**
```json
[
  {
    "id": 1,
    "from_node": 1,
    "to_node": 2,
    "condition": "yes"
  }
]
```

### Update Edge
```
PUT /edge?id={id}
```

### Delete Edge
```
DELETE /edge?id={id}
```

---

## User Scenario Flow

### Start Scenario
```
POST /user/scenario/start
```

**Request Body:**
```json
{
  "user_id": 1,
  "scenario_id": 1
}
```

**Response:**
```json
{
  "state_id": 1,
  "current_node": {
    "id": 1,
    "type": "question",
    "text": "Вы потеряли паспорт?"
  },
  "status": "in_progress"
}
```

### Process Answer (Step)
```
POST /user/scenario/step
```

**Request Body:**
```json
{
  "user_state_id": 1,
  "answer": "yes"
}
```

**Response (Next Node):**
```json
{
  "current_node": {
    "id": 2,
    "type": "action",
    "text": "Уведомить паспортный стол",
    "organization": "ФМС России",
    "deadline": "30 дней",
    "consequences": "Штраф до 2500 рублей"
  },
  "is_completed": false
}
```

**Response (Scenario Complete):**
```json
{
  "current_node": null,
  "is_completed": true,
  "message": "Scenario completed!"
}
```

### Get User Progress
```
GET /user/scenario/status?user_id={id}&scenario_id={id}
```

**Response:**
```json
{
  "state": {
    "id": 1,
    "user_id": 1,
    "scenario_id": 1,
    "current_node_id": 2,
    "status": "in_progress",
    "completed_steps": [1, 2],
    "created_at": "2026-07-05T15:00:00Z",
    "updated_at": "2026-07-05T15:05:00Z"
  },
  "actions": [
    {
      "id": 1,
      "node_id": 2,
      "action": "notification",
      "organization": "ФМС России",
      "deadline": "2026-08-04T00:00:00Z",
      "completed": false,
      "created_at": "2026-07-05T15:05:00Z"
    }
  ]
}
```

---

## Error Responses

### 400 Bad Request
```json
{
  "error": "invalid request"
}
```

### 404 Not Found
```json
{
  "error": "scenario not found"
}
```

### 500 Internal Server Error
```json
{
  "error": "failed to create scenario"
}
```
