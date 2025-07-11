DATABASE="postgres://admin:admin@localhost:5432/service?sslmode=disable"
PATH_MIGRATE="./migrations"

mig-up:
	migrate \
	-path=$(PATH_MIGRATE) \
	-database=$(DATABASE) \
	up

mig-down:
	migrate \
	-path=$(PATH_MIGRATE) \
	-database=$(DATABASE) \
	down

mig-force:
	migrate \
	-path=$(PATH_MIGRATE) \
	-database=$(DATABASE) \
	force 1