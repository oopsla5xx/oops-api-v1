# Development Seeds

SQL seed files for the development environment.

## Convention

- File naming: `NNN_<description>.sql` (e.g. `001_users.sql`)
- Each file must be idempotent (`INSERT ... ON CONFLICT DO NOTHING`)
- Files are executed in alphabetical order

## Running Seeds

```bash
make seed ENV=development
```
