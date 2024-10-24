# Autocrud
Automatically generate REST APIs


# TODO V0.1.0
 - [X] parse yaml schema.
 - [X] validate yaml schema.
 - [X] create sqlite db from yaml schema with a simple resource.
        only types available will be INT, VARCHAR, TIMESTAMP, VARTEXT
 - [X] write test for codegen.
 - [X] workout absolute path for codegen templates.
 - [X] rest api to create and read a single resource:
   - [X] create resource from a single table.
   - [X] read resource from a single table.
 - [X] rest api to update and delete a single resource:
   - [X] update resource from a single table.
   - [X] delete resource from a single table.
 - [X] add DTO type.
   - [X] Config should include type info: primary key, default for insert date,
        default for update date.
   - [X] Use this info from config to create the DTO types
 - [X] getAll frontend: list view with option to add or delete entry.
 - [X] use zod to validate inputs.
 - [X] add notification system.
 - [ ] cmd: autocrud [-f --file <config.yaml>] [-t --template]

# TODO future
 - [ ] getById frontend with edit option.
 - [ ] add proper loggin.
 - [ ] parse yaml schema for relational tables.
 - [ ] search items in getAll view
 - [ ] get names for backend frontend and database directories from config.
 - [ ] make generate DB optional
 - [ ] database backend interface
   - [ ] postgress db in docker 
   - [ ] duckDB
 - [ ] dockerize API
 - [ ] dockerize frontend
 - [ ] add other django backend.
 - [ ] add other views in frontend: eg cards, complex tables, images
 - [ ] add migrations table and go script for migrations
