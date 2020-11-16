# Resource-Manager Built-In Strategies

This subproject contains build-in strategies in the **src** folder.

Before importing (into resource-manager) we need to process individual strategies
 (resolve their dependencies, remove code used in tests etc). After processing the resulting strategies
are copied into the **generated** folder. 
 
 
The following commands are possible:

- run `yarn test` to test the strategies
- run `yarn generate:all` to re-process the strategies and copy them into the **generated** folder
