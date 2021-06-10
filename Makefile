.PHONY: postgres migrate_up migrate_down

postgres: 
		docker start postgresDB

migrate_up:
		migrate -source file://db/migrations \
                        -database postgres://postgres:kjunn2000@localhost/straper_db?sslmode=disable up

migrate_down:
		migrate -source file://db/migrations \
                        -database postgres://postgres:kjunn2000@localhost/straper_db?sslmode=disable down