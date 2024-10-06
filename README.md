# Autocrud
Automatically generate REST APIs


# TODO V0.1.0
 - [X] parse yaml schema.
 - [X] validate yaml schema.
 - [X] create sqlite db from yaml schema with a simple resource.
        only types available will be INT, VARCHAR, TIMESTAMP, VARTEXT
 - [X] write test for codegen.
 - [X] workout absolute path for codegen templates.
 - [ ] rest api to create and read a single resource:
   - [ ] create resource from a single table.
   - [ ] read resource from a single table.
 - [ ] rest api to update and delete a single resource:
   - [ ] update resource from a single table.
   - [ ] delete resource from a single table.
 - [ ] parse yaml schema for relational tables.
 - [ ] getAll frontend: list view with option to add or delete entry.
 - [ ] getById frontend with edit and delete option.
 
# TODO future
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
