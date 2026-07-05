# Life Assistant - Digital Notification Helper

## Project Overview

Life Assistant is a web application that helps Russian citizens determine which organizations need to be notified when their personal data changes (passport, address, phone, etc.), and what deadlines and consequences they must be aware of.

## Features

- **Scenario-based guidance**: Navigate through life events using interactive question-answer flows
- **Organization checklist**: Automatically generated list of organizations to notify
- **Deadline tracking**: Clear deadlines for each notification
- **Consequences info**: Information about penalties for non-compliance
- **Legal references**: Links to relevant regulatory documents
- **User progress tracking**: Resume scenarios at any time

## Architecture

The project follows **Clean Architecture** principles:

```
HTTP Request
    ↓
Handler (API Layer)
    ↓
Service (Business Logic)
    ↓
Repository (Data Access)
    ↓
PostgreSQL Database
```

### Key Components

- **Domain**: Core business entities (Scenario, Graph Nodes/Edges, User, UserState)
- **Repository**: Data access layer with PostgreSQL implementation
- **Service**: Business logic layer (ScenarioService, NodeService, EdgeService, UserService)
- **API Handler**: HTTP request handlers for REST endpoints
- **Graph Engine**: Decision logic for processing scenario flow

## Technology Stack

### Backend
- **Language**: Go 1.26.4
- **Framework**: Chi (lightweight router)
- **Database**: PostgreSQL 15
- **Driver**: sqlx (SQL library without ORM)
- **Container**: Docker & Docker Compose

### Frontend (Coming Soon)
- React with TypeScript
- Tailwind CSS
- TanStack Query
- React Hook Form
- shadcn/ui

## Project Structure

```
backend/
├── cmd/api/                    # Application entry point
│   └── main.go
├── internal/
│   ├── api/
│   │   ├── handler/           # HTTP handlers
│   │   └── response/          # Response utilities
│   ├── app/                   # Application setup
│   ├── config/                # Configuration management
│   ├── domain/                # Business entities
│   │   ├── graph/             # Graph nodes and edges
│   │   ├── scenario/          # Scenario entity
│   │   └── user/              # User and state entities
│   ├── repository/            # Data access layer
│   │   └── postgres/          # PostgreSQL implementation
│   └── service/               # Business logic
│       └── engine/            # Graph engine & scenario flow
├── migrations/                # Database migrations
├── pkg/
│   ├── database/              # Database connection
│   └── logger/                # Logging utilities
├── docker-compose.yml
├── Dockerfile
├── go.mod
└── go.sum
```

## Getting Started

### Prerequisites
- Go 1.26.4 or higher
- Docker & Docker Compose
- PostgreSQL 15 (if running locally)

### Quick Start

1. **Clone repository**
   ```bash
   git clone <repo-url>
   cd life-assistant/backend
   ```

2. **Start PostgreSQL with Docker Compose**
   ```bash
   docker-compose up -d
   ```

3. **Create .env file**
   ```bash
   cp .env.example .env
   ```

4. **Run migrations**
   ```bash
   # Migrations are run automatically on app startup
   # Or manually using your migration tool
   ```

5. **Start the server**
   ```bash
   go run cmd/api/main.go
   ```

Server will be available at `http://localhost:8080`

## API Endpoints

### Health Checks
- `GET /api/health` - Application health
- `GET /api/health/db` - Database connection health

### Scenarios
- `POST /api/scenario` - Create scenario
- `GET /api/scenario?id={id}` - Get scenario by ID
- `GET /api/scenarios` - List all scenarios
- `PUT /api/scenario?id={id}` - Update scenario
- `DELETE /api/scenario?id={id}` - Delete scenario

### Graph Nodes
- `POST /api/node` - Create node
- `GET /api/node?id={id}` - Get node by ID
- `GET /api/nodes?scenario_id={id}` - Get scenario nodes
- `PUT /api/node?id={id}` - Update node
- `DELETE /api/node?id={id}` - Delete node

### Graph Edges
- `POST /api/edge` - Create edge
- `GET /api/edge?id={id}` - Get edge by ID
- `GET /api/edges?scenario_id={id}` - Get scenario edges
- `PUT /api/edge?id={id}` - Update edge
- `DELETE /api/edge?id={id}` - Delete edge

### Scenario Engine (User Flow)
- `POST /api/user/scenario/start` - Start new scenario session
- `POST /api/user/scenario/step` - Process user answer and move to next node
- `GET /api/user/scenario/status?user_id={id}&scenario_id={id}` - Get user progress

## Domain Models

### Scenario
```go
type Scenario struct {
    ID          int64
    Title       string      // "Смена паспорта"
    Description string
    Category    *string     // "документы", "работа"
    StartNodeID int64
    CreatedAt   time.Time
}
```

### Node
```go
type Node struct {
    ID              int64
    ScenarioID      int64
    Type            NodeType    // "question", "info", "action", "result"
    Text            string
    ActionType      *string     // Тип действия
    Organization    *string     // Организация для уведомления
    Deadline        *string     // Сроки в днях
    ReferenceLinks  *string     // JSON: ссылки на НПА
    Consequences    *string     // Последствия нарушения
    Order           int
    CreatedAt       time.Time
}
```

### Edge
```go
type Edge struct {
    ID        int64
    FromNode  int64
    ToNode    int64
    Condition string      // "yes", "no", "continue"
    Logic     *string     // Дополнительная логика
}
```

### UserState
```go
type UserState struct {
    ID              int64
    UserID          int64
    ScenarioID      int64
    CurrentNodeID   int64
    Status          StateStatus // "in_progress", "completed"
    CompletedSteps  []int64
    CreatedAt       time.Time
    UpdatedAt       time.Time
}
```

## Database Schema

See `migrations/` directory for:
- `001_init.sql` - Initial scenario table
- `002_graph.sql` - Graph nodes and edges
- `003_user_state.sql` - User state table
- `004_scenarios.sql` - Scenarios table
- `005_extend_graph.sql` - Extended node/edge fields
- `006_users.sql` - Users table
- `007_user_states.sql` - Extended user states
- `008_user_actions.sql` - User actions tracking

## Development Guidelines

### Code Organization
- Keep business logic in Services, not Handlers
- Don't let Repository know about HTTP
- Services don't know about database details
- Use repository interfaces for abstraction

### Adding New Features
1. Define domain model in `internal/domain/`
2. Create repository interface in `internal/repository/interfaces.go`
3. Implement repository in `internal/repository/postgres/`
4. Create service in `internal/service/`
5. Add HTTP handler in `internal/api/handler/`
6. Register routes in `internal/app/router.go`
7. Add database migrations in `migrations/`

## Testing

### Using Postman
Import `postman_collection.json` to test all endpoints.

### Manual Testing
```bash
# Health check
curl http://localhost:8080/api/health

# Create scenario
curl -X POST http://localhost:8080/api/scenario \
  -H "Content-Type: application/json" \
  -d '{
    "title": "Смена паспорта",
    "description": "Процесс смены паспорта",
    "start_node_id": 1
  }'
```

## Roadmap

### Phase 1 (Current)
- [x] Domain models
- [x] Database schema
- [x] Repository layer
- [x] Service layer
- [x] API handlers
- [x] CRUD operations
- [ ] Input validation
- [ ] Error handling

### Phase 2
- [ ] Frontend (React)
- [ ] User authentication
- [ ] Calendar integration
- [ ] Push notifications

### Phase 3
- [ ] Messenger bot (VK/Telegram)
- [ ] ESID integration (pending accreditation)
- [ ] Advanced rule engine
- [ ] Admin panel

## Contributing

When contributing:
1. Follow Clean Architecture principles
2. Keep code modular and testable
3. Write clear commit messages
4. Don't use ORMs (sqlx only)
5. Maintain backward compatibility

## License

MIT

## Support

For questions or issues, please create a GitHub issue.
