# Documentation Structure

This document explains the consolidated documentation structure for the Atomizer project.

## Main Documentation Files

### 1. README.md
- **Purpose**: Quick start guide and command reference
- **Audience**: New users and developers
- **Contents**: 
  - Project overview
  - Installation instructions
  - Basic usage examples
  - Command reference
  - Links to detailed documentation

### 2. ARCHITECTURE.md
- **Purpose**: Technical design and implementation details
- **Audience**: Developers and contributors
- **Contents**:
  - System architecture diagrams
  - Component descriptions
  - Package structure
  - Exchange integration patterns
  - Data flow
  - Extension points

### 3. OPERATIONS.md
- **Purpose**: Deployment, configuration, and operational guide
- **Audience**: System operators and traders
- **Contents**:
  - Configuration options
  - Deployment strategies
  - Trading operations
  - Hedging strategies
  - Monitoring and analysis
  - Troubleshooting
  - Performance tuning

## Archived Documentation

Historical and migration-related documentation has been moved to `docs/archive/`:

- Migration and refactoring documentation
- Old architecture documents
- Legacy component READMEs
- Historical design decisions

## Reference Documentation

External API documentation is stored in `docs/reference/`:

- `derive.md` - Complete Derive API reference

## Component-Specific Documentation

Some directories maintain their own documentation for specific tools:

- `scripts/market_analysis/README.md` - Market analysis tools documentation
- `sdk/README.md` - SDK submodule documentation

## Consolidation Benefits

1. **Reduced Duplication**: Information is now in one authoritative location
2. **Better Organization**: Clear separation between architecture, operations, and quick start
3. **Easier Maintenance**: Three main files instead of 20+ scattered documents
4. **Improved Discovery**: Users know exactly where to look for information
5. **Progressive Detail**: README → Architecture → Operations flow