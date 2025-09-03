# Football Backend API

A professional Go-based REST API for managing American football teams, players, and games. Built with clean architecture principles, this backend provides comprehensive endpoints for football data management with detailed statistics tracking.

## 🏈 Features

- **Teams Management**: Create and manage NFL teams with conference and division information
- **Player Management**: Add and manage players with detailed physical attributes and positions
- **Games Management**: Create and manage football games with scores, status, and week tracking
- **Player Statistics**: Comprehensive football statistics including offensive, defensive, and special teams stats
- **Clean Architecture**: Layered architecture with handlers, services, and repositories
- **RESTful API**: Clean, RESTful endpoints following proper resource-based organization
- **SQLite Database**: Lightweight, file-based database with proper foreign key relationships
- **CORS Support**: Ready for frontend integration
- **Input Validation**: Comprehensive validation for all API endpoints
- **Error Handling**: Proper error handling with meaningful error messages

## 🚀 Prerequisites

- Go 1.21 or higher
- Git

## 📦 Installation

1. **Clone or navigate to the project directory**:
   ```bash
   cd sports-backend
   ```

2. **Install dependencies**:
   ```bash
   go mod tidy
   ```

3. **Run the server**:
   ```bash
   go run main.go
   ```

The server will start on port 8080 by default. You can change this by setting the `PORT` environment variable.

## 🔗 API Endpoints

### Health Check
- `GET /health` - Check if the server is running

### Teams
- `GET /api/teams` - Get all teams
- `POST /api/teams` - Create a new team
- `GET /api/teams/{id}` - Get a specific team
- `PUT /api/teams/{id}` - Update a team
- `DELETE /api/teams/{id}` - Delete a team
- `GET /api/teams/{id}/stats` - Get statistics for a specific team (coming soon)
- `POST /api/teams/{id}/stats` - Create team statistics (coming soon)

### Players
- `GET /api/players` - Get all players
- `POST /api/players` - Create a new player
- `GET /api/players/{id}` - Get a specific player
- `PUT /api/players/{id}` - Update a player
- `DELETE /api/players/{id}` - Delete a player
- `GET /api/players/{id}/stats` - Get statistics for a specific player (coming soon)
- `POST /api/players/{id}/stats` - Create player statistics (coming soon)
- `PUT /api/players/{id}/stats/{stats_id}` - Update player statistics (coming soon)
- `DELETE /api/players/{id}/stats/{stats_id}` - Delete player statistics (coming soon)

## 📝 API Usage Examples

### Create a Team
```bash
curl -X POST http://localhost:8080/api/teams \
  -H "Content-Type: application/json" \
  -d '{
    "name": "Chiefs",
    "city": "Kansas City",
    "conference": "AFC",
    "division": "West"
  }'
```

### Create a Player
```bash
curl -X POST http://localhost:8080/api/players \
  -H "Content-Type: application/json" \
  -d '{
    "team_id": 1,
    "first_name": "Patrick",
    "last_name": "Mahomes",
    "position": "QB",
    "jersey_number": 15,
    "height": 75,
    "weight": 230
  }'
```

### Update a Player
```bash
curl -X PUT http://localhost:8080/api/players/1 \
  -H "Content-Type: application/json" \
  -d '{
    "jersey_number": 15,
    "weight": 235
  }'
```

### Get All Teams
```bash
curl http://localhost:8080/api/teams
```

### Get All Players
```bash
curl http://localhost:8080/api/players
```

## 🏗️ Data Models

### Team
```json
{
  "id": 1,
  "name": "Chiefs",
  "city": "Kansas City",
  "conference": "AFC",
  "division": "West",
  "created_at": "2024-01-15T10:30:00Z",
  "updated_at": "2024-01-15T10:30:00Z"
}
```

### Player
```json
{
  "id": 1,
  "team_id": 1,
  "first_name": "Patrick",
  "last_name": "Mahomes",
  "position": "QB",
  "jersey_number": 15,
  "height": 75,
  "weight": 230,
  "created_at": "2024-01-15T10:30:00Z",
  "updated_at": "2024-01-15T10:30:00Z"
}
```

### PlayerStats (Coming Soon)
Comprehensive football statistics including:
- **Offensive**: Passing attempts, completions, yards, touchdowns, interceptions
- **Rushing**: Attempts, yards, touchdowns
- **Receiving**: Targets, receptions, yards, touchdowns
- **Defensive**: Tackles, sacks, interceptions, pass deflections
- **Special Teams**: Field goals, punts, kick returns, punt returns

## 🗄️ Database

The application uses SQLite for data storage. The database file (`sports.db`) will be created automatically when you first run the application. Database migrations are run automatically on startup with `CREATE TABLE IF NOT EXISTS` statements for safe re-runs.

### Database Schema
- **teams**: Team information with conference and division
- **players**: Player information with team relationships
- **games**: Game information (coming soon)
- **player_stats**: Detailed player statistics (coming soon)

## 🌍 Environment Variables

- `PORT`: Server port (default: 8080)

## 📁 Project Structure

```
sports-backend/
├── main.go                    # Application entry point with dependency injection
├── go.mod                     # Go module file
├── go.sum                     # Go module checksums
├── models/
│   ├── player.go             # Player and PlayerStats models
│   └── team.go               # Team and Game models
├── handlers/
│   ├── player_handler.go     # Player HTTP handlers
│   └── team_handler.go       # Team HTTP handlers
├── services/
│   ├── player_service.go     # Player business logic
│   └── team_service.go       # Team business logic
├── repositories/
│   ├── player_repository.go  # Player data access
│   └── team_repository.go    # Team data access
├── database/
│   └── migrations.go         # Database migrations
└── README.md                 # This file
```

## 🏛️ Architecture

This project follows clean architecture principles with clear separation of concerns:

```
HTTP Request → Handler → Service → Repository → Database
                ↓         ↓         ↓
HTTP Response ← Handler ← Service ← Repository ← Database
```

- **Handlers**: HTTP request/response handling, JSON encoding/decoding
- **Services**: Business logic, validation, data transformation
- **Repositories**: Data access, SQL queries, database operations
- **Models**: Data structures and request/response DTOs

## 🧪 Testing

You can test the API using curl, Postman, or any HTTP client. The server includes CORS headers to allow frontend applications to connect.

### Quick Test
```bash
# Check if server is running
curl http://localhost:8080/health

# Get all teams
curl http://localhost:8080/api/teams

# Get all players
curl http://localhost:8080/api/players
```

## 🚀 Quick Start with Sample Data

1. **Start the server**:
   ```bash
   go run main.go
   ```

2. **Create some teams**:
   ```bash
   curl -X POST http://localhost:8080/api/teams \
     -H "Content-Type: application/json" \
     -d '{"name": "Chiefs", "city": "Kansas City", "conference": "AFC", "division": "West"}'
   
   curl -X POST http://localhost:8080/api/teams \
     -H "Content-Type: application/json" \
     -d '{"name": "Bills", "city": "Buffalo", "conference": "AFC", "division": "East"}'
   ```

3. **Create some players**:
   ```bash
   curl -X POST http://localhost:8080/api/players \
     -H "Content-Type: application/json" \
     -d '{"team_id": 1, "first_name": "Patrick", "last_name": "Mahomes", "position": "QB", "jersey_number": 15}'
   
   curl -X POST http://localhost:8080/api/players \
     -H "Content-Type: application/json" \
     -d '{"team_id": 2, "first_name": "Josh", "last_name": "Allen", "position": "QB", "jersey_number": 17}'
   ```

## 🔮 Future Enhancements

- **Games Management**: Complete game creation and management
- **Player Statistics**: Full statistics tracking and analysis
- **Team Statistics**: Team-level performance metrics
- **Authentication**: JWT-based authentication and authorization
- **Data Validation**: Enhanced input validation middleware
- **Rate Limiting**: API rate limiting for production use
- **Logging**: Structured logging with different levels
- **Database Connection Pooling**: Optimized database connections
- **Real-time Updates**: WebSocket support for live updates
- **Data Export/Import**: CSV/JSON data export functionality
- **API Documentation**: OpenAPI/Swagger documentation
- **Unit Tests**: Comprehensive test coverage
- **Docker Support**: Containerization for easy deployment

## 📄 License

This project is open source and available under the MIT License.

## 🤝 Contributing

1. Fork the repository
2. Create a feature branch
3. Make your changes
4. Add tests if applicable
5. Submit a pull request

## 📞 Support

If you have any questions or need help, please open an issue in the repository.