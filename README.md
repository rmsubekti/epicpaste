# Just Another Simple Note Taking App
EPIC PASTE is just like a note taking app but for snippets.

### Try epicpaste using docker compose
```yml
version: "3.9"

services:
  db:
    image: postgres:latest
    restart: always
    environment:
      POSTGRES_DB: epicpaste
      POSTGRES_PASSWORD: EpicPaste*password*

  epicpaste:
    image: rmsubekti/epicpaste:latest
    restart: unless-stopped
    environment:
      POSTGRES_DB: epicpaste # Same as postgres environment
      POSTGRES_PASSWORD: EpicPaste*password* # Same as postgres environment
      POSTGRES_HOSTNAME: db
      EPIC_HOSTNAME: localhost #192.168.192.199 or domainname.com
    ports:
      - 80:80
    depends_on:
      - db
```

### All environment variables avalaible
```toml
POSTGRES_USER= dbuser
POSTGRES_PASSWORD= dbpassword
POSTGRES_DB= dbname
POSTGRES_HOSTNAME=localhost
POSTGRES_SSLMODE=disable
POSTGRES_HOST_AUTH_METHOD=trust
POSTGRES_PORT=5432

EPIC_HOSTNAME=localhost 
EPIC_PORT="3030"
EPIC_JWT_SECRET_KEY= saltKey
EPIC_COOKIE_SECRET_KEY= saltKey
```

### API DOCUMENTATION
This is only a backend app, you need to build the frontend yourself. See the documentation at:
```
localhost:80/docs/index.html
```