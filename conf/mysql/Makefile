# Migrate version / arguments
MIGRATE_VERSION=v4.11.0
MIGRATE_MODE?=up
MIGRATE_STEPS?=1

all: getimage

# Pull the database migration framework
getimage: 
		docker pull migrate/migrate:$(MIGRATE_VERSION)

# Run and enter the database migration framework image
enterContainer: check-db
		echo "Enter migration container"
		docker run -it -v $(CURDIR)/$(DATABASE):/migrations --network host --entrypoint /bin/sh migrate/migrate

# Migrate the database schema
# To migrate a database, you can use:
# make migrate HOST=localhost DATABASE=botdata DBUSER=root DBPASS=root
migrate: check-env check-db check-credentials
		echo "Migrating $(DATABASE) db $(MIGRATE_STEPS) version $(MIGRATE_MODE)"
		docker run -v $(CURDIR)/migrations/$(DATABASE):/migrations --network host migrate/migrate:$(MIGRATE_VERSION) -verbose -path=/migrations/ -database "mysql://$(DBUSER):$(DBPASS)@tcp($(HOST):3306)/$(DATABASE)" $(MIGRATE_MODE) $(MIGRATE_STEPS)


# Check for missing target arguments
check-env:
ifndef HOST
	$(error HOST is undefined)
endif

check-db:
ifndef DATABASE
	$(error DATABASE is undefined)
endif

check-credentials:
ifndef DBUSER
	$(error DBUSER is undefined)
endif
ifndef DBPASS
	$(error DBPASS is undefined)
endif