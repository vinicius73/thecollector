# TheCollector - Database Backup Tool

A tool to make and schedule postgres backups.
All bakups are sent to S3 Buckets with checksums

Check [`thecollector.yml`](thecollector.yml) about configs and environment variables.

You also can define a custom config with all your datasources.

## Dependencies

- [Go v1.21](https://go.dev/dl/)
- [task](https://taskfile.dev/installation/)

### Requeriments

- [pg_dump](https://www.postgresql.org/docs/current/app-pgdump.html)
  - [postgresql-libs - Archlinux](https://archlinux.org/packages/extra/x86_64/postgresql-libs/)
  - [postgresql15-client - Alpine](https://pkgs.alpinelinux.org/package/edge/main/x86/postgresql15-client)

## Other projects

Check my other Go projects

- [gear-feed - Bot and Scrapper system](https://github.com/vinicius73/gear-feed)
- [peristera - Telegram Bot for ORGs](https://github.com/comunidade-shallom/peristera)
