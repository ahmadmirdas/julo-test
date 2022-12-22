## Guideline Running Code

1. Build service using `make start`
2. Run Migration DB `make migration-up` (you can using library migration go like `goose`)
3. Access several API has been provide with PreffixUrl `/api/v1`
4. You can access database using adminer, to access them please [here](http://localhost:8080/?pgsql=postgres&username=postgres&db=julotest&ns=public)
