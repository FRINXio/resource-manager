# Built-in allocation strategies

Contains scripts to load build-in strategies into the database.

To update the strategies in resource-manager follow these steps:

1. delete the strategies in the DB (or wipe the whole DB)
2. in the **backend** folder run the following command `go generate ./pools/...`
3. start resource-manager 
