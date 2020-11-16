
## template

before this... install this

```sh
go install \
    github.com/grpc-ecosystem/grpc-gateway/v2/protoc-gen-grpc-gateway \
    github.com/grpc-ecosystem/grpc-gateway/v2/protoc-gen-openapiv2 \
    google.golang.org/protobuf/cmd/protoc-gen-go \
    google.golang.org/grpc/cmd/protoc-gen-go-grpc \
    github.com/grpc-ecosystem/grpc-health-probe
```

before this.. get this
```sh
go get github.com/mwitkow/go-proto-validators/protoc-gen-govalidators
```

```sh
protoc \
-I . \
-I ${GOPATH}/src \
-I `go list -m -f "{{.Dir}}" github.com/golang/protobuf` \
-I `go list -m -f "{{.Dir}}" google.golang.org/protobuf` \
-I `go list -m -f "{{.Dir}}" github.com/mwitkow/go-proto-validators` \
--go_out . \
--go_opt paths=source_relative \
--go-grpc_out . \
--go-grpc_opt paths=source_relative \
--govalidators_out . \
--grpc-gateway_out . \
--grpc-gateway_opt logtostderr=true \
--grpc-gateway_opt paths=source_relative \
--grpc-gateway_opt generate_unbound_methods=true \
--openapiv2_out . \
--openapiv2_opt logtostderr=true \
--openapiv2_opt generate_unbound_methods=true \
pkg/proto/grpcgwgokit/pb/grpcgwgokit.proto
```

https://github.com/rephus/grpc-gateway-example


## Migrations

For migrations we are using [golang-migrate](https://github.com/golang-migrate/migrate) .

You can follow this awesome original [tutorial](https://github.com/golang-migrate/migrate/blob/master/database/postgres/TUTORIAL.md) or continue reading what I did with my project.


1. install golang-migrate cli tool.

We can dowload the file and put it into `$GOPATH` or just install with `brew`

$GOPATH

* please go to the folder `$GOPATH/bin`
* download the golang migrate cli package
```bash
curl -L https://github.com/golang-migrate/migrate/releases/download/v4.13.0/migrate.darwin-amd64.tar.gz | tar xvz
```
* change its name to just `migrate`

BREW

```bash
brew install golang-migrate
```

CLI usage can be found [here](https://github.com/golang-migrate/migrate/tree/master/cmd/migrate)

2. Create database e.g. `example`

```bash
psql -h localhost -U postgres -w -c "create database example;"
```

you can add a `init.sql` file into `db` folder to create this as well. The file will contain this

```sql
CREATE DATABASE example;
GRANT ALL PRIVILEGES ON DATABASE example TO postgres;
```

3. add this env var. (remember to change the password)

```bash
export POSTGRESQL_URL='postgres://postgres:password@localhost:5432/example?sslmode=disable'
```

4. get back to the project folder.

5. let's create a table called user (this will create only the files for upgrade and downgrade)

```sh
migrate create -ext sql -dir db/migrations -seq create_users_table

%YOUR_GO_WORKSPACE%/grpcgwgokit/db/migrations/000001_create_users_table.up.sql
%YOUR_GO_WORKSPACE%/grpcgwgokit/db/migrations/000001_create_users_table.down.sql
```

you should see two EMPTY files under `db/migrations`

* 000001_create_users_table.up.sql
* 000001_create_users_table.down.sql

4. Add the expected create table syntax on `000001_create_users_table.up.sql` file.

```sql
CREATE TABLE IF NOT EXISTS users(
   user_id serial PRIMARY KEY,
   username VARCHAR (50) UNIQUE NOT NULL,
   password VARCHAR (50) NOT NULL,
   email VARCHAR (300) UNIQUE NOT NULL
);
```

5. Add the expected drop table syntax in case we need to downgrade the migration.

```sql
DROP TABLE IF EXISTS users;
```

6. Now run the migration and see what happens. Remember that we added the `${POSTGRESQL_URL}` env var in step `3` 

```bash
migrate -database ${POSTGRESQL_URL} -path db/migrations up

1/u create_users_table (14.130244ms)
```

7. Great opportunity to test how `down` works

```bash
migrate -database ${POSTGRESQL_URL} -path db/migrations down

Are you sure you want to apply all down migrations? [y/N]
y
Applying all down migrations
1/d create_users_table (10.672845ms)
```

8. let's continue adding more things to the users table.

```bash
migrate create -ext sql -dir db/migrations -seq add_mood_to_users

%YOUR_GO_WORKSPACE%/grpcgwgokit/db/migrations/000002_add_mood_to_users.up.sql
%YOUR_GO_WORKSPACE%/grpcgwgokit/db/migrations/000002_add_mood_to_users.down.sql
```

see current files in the project and check the sql syntax. Remember the `BEGIN;` and `COMMIT;` for postgresql transactions.

9. execute the migrate up again

```bash
migrate -database ${POSTGRESQL_URL} -path db/migrations up
2/u add_mood_to_users (9.609364ms)
```


