# BeesInTheTrap

A turn-based command-line game where you battle against a hive of bees. Destroy the hive before the bees sting you to death!

## Game Rules
- Players start with 100 HP
- The hive consists of three types of bees:
  - 1 Queen Bee (100 HP, deals 10 damage)
  - 5 Worker Bees (75 HP each, deal 5 damage)
  - 25 Drone Bees (60 HP each, deal 1 damage)
- Enter "hit" during your turn to attack a random bee
- After your turn, the bees will attack you
- When the Queen Bee dies, all remaining bees die too
- The game ends when either all bees are dead, or you die
- Both you and the bees have a chance to miss attacks

## Installation

### Pre-compiled Binaries

Download the latest release for your platform from the [Releases](https://github.com/lewwolfe/beesinthetrap/releases) page - v1.0.0.

### Build from Source

1. Clone the repository
```
git clone https://github.com/lewwolfe/beesinthetrap.git
cd beesinthetrap
```

2. Build the project
```
make build
```

## Usage

Run the game (pick with OS you are running) :
```
./beesinthetrap-[windows/linux]-amd64
```

Or using make:
```
make run
```

## Development

### Prerequisites

- Go 1.20+
- golangci-lint (for linting)

### Available Commands

- `make all`: Run linting, tests, and build
- `make build`: Build the application
- `make test`: Run all tests
- `make lint`: Run code linting
- `make clean`: Remove build artifacts
- `make run`: Build and run the application

### CI/CD

This project uses GitHub Actions for continuous integration:

- Automatic build for both Windows and Linux
- Automated tests on push
- Release binaries generation for both platforms

## Project Structure

```
beesinthetrap/
├── cmd/
│   └── beesinthetrap/
│       └── main.go       # Application entry point
├── cli/
│   ├── cli.go            # Handle all the inputs and output logic
│   └── messages.go       # Display message to CLI
├── config/
│   └── config.go         # Game configuration
├── game/
│   ├── bee.go            # Bee types and behavior
│   ├── game.go           # Game engine
│   └── player.go         # Player implementation
├── .github/
│   └── workflows/        # GitHub Actions workflows
├── Makefile              # Build automation
└── README.md             # This file
```

## License

[MIT](LICENSE)
