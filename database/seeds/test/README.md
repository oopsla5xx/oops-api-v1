# Test Seeds

SQL seed files for the test environment.

## Convention

- Data must be deterministic — same input produces same output on every run
- File naming: `NNN_<description>.sql`
- Each file must be idempotent (`INSERT ... ON CONFLICT DO NOTHING`)

## Running Seeds

```bash
make seed ENV=test
```
