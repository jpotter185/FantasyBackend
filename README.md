# Football Backend API

A Go-based REST API for managing American football games, player statistics, and team data. This backend provides endpoints for creating, reading, updating, and deleting football-related data with comprehensive statistics tracking for players, teams, and games.

## Features

- **Teams Management**: Create and manage football teams with conference and division information
- **Games Management**: Create and manage football games with scores, status, and week tracking
- **Player Management**: Add and manage player information with detailed physical attributes
- **Player Statistics**: Comprehensive football statistics including offensive, defensive, and special teams stats
- **Team Statistics**: Team-level statistics for games including yards, downs, and possession
- **Game Statistics**: Game-level information including weather, attendance, and venue details
- **RESTful API**: Clean, RESTful endpoints following proper resource-based organization
- **SQLite Database**: Lightweight, file-based database with proper foreign key relationships
- **CORS Support**: Ready for frontend integration

## Prerequisites

- Go 1.21 or higher
- Git

## Installation

1. **Clone or navigate to the project directory**:
   ```bash
   cd sports-backend
   ```

2. **Install dependencies**:
   ```bash
   go mod tidy
   ```

3. **Build the application**:
   ```bash
   go build -o sports-backend .
   ```

4. **Run the server**:
   ```bash
   ./sports-backend
   ```

The server will start on port 8080 by default. You can change this by setting the `PORT` environment variable.

## API Endpoints

### Health Check
- `GET /health` - Check if the server is running

### Teams
- `GET /api/teams` - Get all teams
- `POST /api/teams` - Create a new team
- `GET /api/teams/{id}` - Get a specific team
- `PUT /api/teams/{id}` - Update a team
- `DELETE /api/teams/{id}` - Delete a team
- `GET /api/teams/{id}/stats` - Get statistics for a specific team

### Games
- `GET /api/games` - Get all games
- `POST /api/games` - Create a new game
- `GET /api/games/{id}` - Get a specific game
- `PUT /api/games/{id}` - Update a game
- `DELETE /api/games/{id}` - Delete a game
- `GET /api/games/{id}/stats` - Get statistics for a specific game

### Players
- `GET /api/players` - Get all players
- `POST /api/players` - Create a new player
- `GET /api/players/{id}` - Get a specific player
- `PUT /api/players/{id}` - Update a player
- `DELETE /api/players/{id}` - Delete a player
- `GET /api/players/{id}/stats` - Get statistics for a specific player

### Player Statistics
- `GET /api/player-stats` - Get all player statistics
- `POST /api/player-stats` - Create new player statistics

### Team Statistics
- `GET /api/team-stats` - Get all team statistics
- `POST /api/team-stats` - Create new team statistics

### Game Statistics
- `GET /api/game-stats` - Get all game statistics
- `POST /api/game-stats` - Create new game statistics

## API Usage Examples

### Create a Game
```bash
curl -X POST http://localhost:8080/api/games \
  -H "Content-Type: application/json" \
  -d '{
    "home_team": "Lakers",
    "away_team": "Warriors",
    "sport": "basketball",
    "league": "NBA",
    "season": "2024-25",
    "game_date": "2024-12-25T20:00:00Z"
  }'
```

### Update Game Score
```bash
curl -X PUT http://localhost:8080/api/games/1 \
  -H "Content-Type: application/json" \
  -d '{
    "status": "completed",
    "home_score": 110,
    "away_score": 108
  }'
```

### Create a Player
```bash
curl -X POST http://localhost:8080/api/players \
  -H "Content-Type: application/json" \
  -d '{
    "name": "LeBron James",
    "team": "Lakers",
    "position": "Forward",
    "sport": "basketball",
    "league": "NBA"
  }'
```

### Add Player Statistics
```bash
curl -X POST http://localhost:8080/api/player-stats \
  -H "Content-Type: application/json" \
  -d '{
    "player_id": 1,
    "game_id": 1,
    "points": 25,
    "rebounds": 8,
    "assists": 10,
    "steals": 2,
    "blocks": 1,
    "turnovers": 3,
    "fouls": 2
  }'
```

## Data Models

### Game
- `id`: Unique identifier
- `home_team`: Home team name
- `away_team`: Away team name
- `sport`: Sport type (basketball, football, baseball, etc.)
- `league`: League name (NBA, NFL, MLB, etc.)
- `season`: Season identifier
- `game_date`: Date and time of the game
- `status`: Game status (scheduled, in_progress, completed, cancelled)
- `home_score`: Home team score (optional)
- `away_score`: Away team score (optional)

### Player
- `id`: Unique identifier
- `name`: Player name
- `team`: Team name
- `position`: Player position
- `sport`: Sport type
- `league`: League name

### PlayerStats
Supports statistics for multiple sports:
- **Basketball**: points, rebounds, assists, steals, blocks, turnovers, fouls
- **Football/Soccer**: goals, assists, yellow_cards, red_cards
- **Baseball**: at_bats, hits, runs, rbis, home_runs, strikeouts, walks

## Database

The application uses SQLite for data storage. The database file (`sports.db`) will be created automatically when you first run the application. Database migrations are run automatically on startup.

## Environment Variables

- `PORT`: Server port (default: 8080)
- `DB_PATH`: Database file path (default: ./sports.db)

## Project Structure

```
sports-backend/
├── main.go                 # Application entry point
├── go.mod                  # Go module file
├── go.sum                  # Go module checksums
├── models/
│   └── sports.go          # Data models
├── handlers/
│   ├── games.go           # Game-related handlers
│   ├── players.go         # Player-related handlers
│   └── stats.go           # Statistics handlers
├── database/
│   ├── connection.go      # Database connection
│   └── migrations.go      # Database migrations
└── README.md              # This file
```

## Development

### Adding New Endpoints

1. Add new handlers in the `handlers/` directory
2. Register routes in `main.go`
3. Update data models in `models/sports.go` if needed
4. Add database migrations in `database/migrations.go` if needed

### Testing

You can test the API using curl, Postman, or any HTTP client. The server includes CORS headers to allow frontend applications to connect.

## Future Enhancements

- Authentication and authorization
- Data validation middleware
- Rate limiting
- Logging middleware
- Database connection pooling
- Support for additional sports
- Real-time updates via WebSockets
- Data export/import functionality

## License

This project is open source and available under the MIT License.
