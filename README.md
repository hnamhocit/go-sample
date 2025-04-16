# GO SAMPLE

-   Clean folder structure for beginners
-   SQLC to generate sql queries, migrations and models
-   Go GIN framework
-   Load env, jwt auth, argon2

-   Create migrations file: `
goose -dir ./sql/migrations mysql "root:@hn71LM5@tcp(127.0.0.1:3306)/sample" create create_users_table sql`

-   Sync migrations files: `
goose -dir ./sql/migrations mysql "root:@hn71LM5@tcp(127.0.0.1:3306)/sample" up`

-   Roll back migrations files: `
goose -dir ./sql/migrations mysql "root:@hn71LM5@tcp(127.0.0.1:3306)/sample" down`
