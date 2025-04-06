[![Go Report Card](https://goreportcard.com/badge/github.com/furmanp/relaise)](https://goreportcard.com/report/github.com/furmanp/relaise)
![Latest Release](https://img.shields.io/github/v/release/furmanp/relaise)
# Relaise

Relaise is a tool to automatically generate release changelogs or release notes based on commit messages from the latest Git tag.

## Features
- Extract commit messages since the latest Git tag.
- Categorize commits (e.g., features, fixes, etc.).
- Generate a formatted changelog.

## Installation
1. Clone the repository:
   ```bash
   git clone https://github.com/your-username/relaise.git
   cd relaise
   ```
2. Build the project:
   ```bash
   go build
   ```

## Usage
Run the tool in a Git repository:
```bash
./relaise
```


## License
This project is licensed under the MIT License, which allows for free, unrestricted use, copying, modification, and distribution with attribution.
