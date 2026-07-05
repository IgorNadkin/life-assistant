# Architecture Documentation

## Clean Architecture Overview

This project implements Clean Architecture to maintain separation of concerns and ensure testability.

### Layer Diagram

```
┌─────────────────────────────────────┐
│         HTTP/REST (chi)             │
├─────────────────────────────────────┤
│         Handler Layer               │
│  (HTTP Request/Response handling)   │
├─────────────────────────────────────┤
│         Service Layer               │
│  (Business Logic, Use Cases)        │
├─────────────────────────────────────┤
│         Repository Layer            │
│  (Data Access Abstraction)          │
├─────────────────────────────────────┤
│         PostgreSQL Database         │
└─────────────────────────────────────┘
```

## Key Principles

### 1. Dependency Inversion
- Services depend on Repository interfaces, not implementations
- Handlers depend on Service interfaces
- Easy to swap implementations (PostgreSQL → MongoDB)

### 2. Separation of Concerns
- **Handler**: Only HTTP concerns (parsing, serialization)
- **Service**: Only business logic (validation, decisions)
- **Repository**: Only database queries
- **Domain**: Pure business entities, no technical details

### 3. Independence
- Business logic doesn't know about HTTP
- Database queries don't know about handlers
- Easy to add new interfaces (gRPC, events)

## Directory Structure

### `cmd/api/`
Application entry point. Minimal code - just bootstraps the app.

### `internal/domain/`
Pure business entities with no external dependencies.

```go
// Example: Scenario doesn't know it's stored in PostgreSQL
type Scenario struct {
    ID          int64
    Title       string
    Description string
    StartNodeID int64
}
```

### `internal/repository/`
Data access abstraction.

```go
// Interface (in interfaces.go)
type ScenarioRepository interface {
    Create(s *Scenario) (int64, error)
    Get(id int64) (*Scenario, error)
    Update(s *Scenario) error
    Delete(id int64) error
}

// Implementation (in postgres/)
type ScenarioRepo struct {
    db *sqlx.DB
}

func (r *ScenarioRepo) Create(s *Scenario) (int64, error) {
    // PostgreSQL-specific code here
}
```

### `internal/service/`
Business logic layer.

```go
type ScenarioService struct {
    scenarioRepo repository.ScenarioRepository  // Depends on interface
    stateRepo    repository.UserStateRepository
}

func (s *ScenarioService) StartScenario(userID, scenarioID int64) (*UserState, error) {
    // Business logic here
    // Doesn't care if data is in PostgreSQL, MongoDB, etc.
}
```

### `internal/api/handler/`
HTTP request handling.

```go
type ScenarioHandler struct {
    service *service.ScenarioService  // Depends on service
}

func (h *ScenarioHandler) Create(w http.ResponseWriter, r *http.Request) {
    // Only HTTP concerns here
    // Parse request, call service, serialize response
}
```

## Data Flow Examples

### Creating a Scenario

```
1. HTTP POST /api/scenario
   ↓
2. ScenarioHandler.Create receives request
   - Parses JSON
   - Validates input
   ↓
3. Calls ScenarioService.CreateScenario(scenario)
   - Applies business rules
   - Calls repository method
   ↓
4. ScenarioRepo.Create executes SQL
   - Inserts into database
   - Returns ID
   ↓
5. Service returns to Handler
   ↓
6. Handler serializes response
   - Returns 201 Created with ID
```

### Processing Scenario Step

```
1. HTTP POST /api/user/scenario/step
   ↓
2. EngineHandler.Step receives answer
   ↓
3. Calls ScenarioFlow.ProcessAnswer
   - Gets current node from engine
   - Determines next node based on answer
   - Updates user state
   ↓
4. Updates database via UserStateRepository
   ↓
5. If node type is "action", creates UserAction
   - Calls UserActionRepo.Create
   ↓
6. Returns next node to handler
   ↓
7. Handler serializes and returns to client
```

## Adding New Features

### Example: Add "Delete Scenario" Feature

1. **Domain** (already exists, no changes)
   ```go
   // internal/domain/scenario/scenario.go
   type Scenario struct { ... }
   ```

2. **Repository Interface**
   ```go
   // internal/repository/interfaces.go
   type ScenarioRepository interface {
       Delete(id int64) error  // Add this
   }
   ```

3. **Repository Implementation**
   ```go
   // internal/repository/postgres/scenario_repo.go
   func (r *ScenarioRepo) Delete(id int64) error {
       _, err := r.db.Exec(`DELETE FROM scenario WHERE id=$1`, id)
       return err
   }
   ```

4. **Service**
   ```go
   // internal/service/scenario_service.go
   func (s *ScenarioService) DeleteScenario(id int64) error {
       return s.scenarioRepo.Delete(id)
   }
   ```

5. **Handler**
   ```go
   // internal/api/handler/scenario.go
   func (h *ScenarioHandler) Delete(w http.ResponseWriter, r *http.Request) {
       idStr := r.URL.Query().Get("id")
       id, _ := strconv.ParseInt(idStr, 10, 64)
       if err := h.scenarioService.DeleteScenario(id); err != nil {
           http.Error(w, "failed", 500)
           return
       }
       json.NewEncoder(w).Encode(map[string]string{"status": "ok"})
   }
   ```

6. **Routes**
   ```go
   // internal/app/router.go
   r.Delete("/scenario", scenarioHandler.Delete)
   ```

## Dependency Injection

We use constructor injection for all dependencies:

```go
// Good: Dependencies are explicit
func NewScenarioService(
    scenarioRepo repository.ScenarioRepository,
    stateRepo repository.UserStateRepository,
) *ScenarioService {
    return &ScenarioService{
        scenarioRepo: scenarioRepo,
        stateRepo:    stateRepo,
    }
}

// Bad: Hidden dependencies (we don't do this)
var globalDB *sql.DB
func NewScenarioService() *ScenarioService {
    return &ScenarioService{
        db: globalDB,  // Where did this come from?
    }
}
```

## Testing Strategy

With this architecture, testing is straightforward:

```go
// Mock repository for testing
type MockScenarioRepo struct {
    scenarios map[int64]*Scenario
}

func (m *MockScenarioRepo) Create(s *Scenario) (int64, error) {
    m.scenarios[s.ID] = s
    return s.ID, nil
}

// Test service without database
func TestStartScenario(t *testing.T) {
    mockRepo := &MockScenarioRepo{scenarios: make(map[int64]*Scenario)}
    service := NewScenarioService(mockRepo, nil)
    
    state, err := service.StartScenario(1, 1)
    assert.NoError(t, err)
    assert.NotNil(t, state)
}
```

## Common Mistakes to Avoid

❌ **Don't**: Import `net/http` in service
```go
// Bad!
func (s *Service) Process() (http.StatusOK, error) { }
```

❌ **Don't**: Use database queries in handlers
```go
// Bad!
func (h *Handler) Get(w http.ResponseWriter, r *http.Request) {
    rows, _ := db.Query("SELECT ...")  // Should use repository!
}
```

❌ **Don't**: Add business logic to domain models
```go
// Bad!
func (s *Scenario) Save() error {  // Save is infrastructure concern!
    return db.Insert(s)
}
```

✅ **Do**: Keep layers separated
```go
// Good!
scenarioService.CreateScenario(scenario)  // Service handles logic
scenarioRepo.Create(scenario)               // Repo handles persistence
scenarioHandler.Create(w, r)                // Handler handles HTTP
```

## References

- [Clean Code by Robert C. Martin](https://www.amazon.com/Clean-Code-Handbook-Software-Craftsmanship/dp/0132350882)
- [Clean Architecture](https://blog.cleancoder.com/uncle-bob/2012/08/13/the-clean-architecture.html)
- [Go Project Layout](https://github.com/golang-standards/project-layout)
