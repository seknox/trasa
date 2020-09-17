## Development Setup

- Clone the repo  
 `git clone https://github.com/seknox/trasa.git && cd trasa`
- Download [GeoCityLite.mmdb]() and move it to `/etc/trasa/static/GeoCityLite.mmdb`
- Install dashboard dependencies   
`yarn install` inside dashboard directory
- Start dashboard in dev mode
`yarn start`
- Build backend server
`go build` inside server directory
- Run the binary
`sudo ./server`

- Edit `/etc/trasa/config/config.toml` file and change `trasa.proxyDashboard` as ` true`
- Open TRASA dashboard at http://localhost:3000

- Go through the wiki to get overview of codebase. 


# Code convensions 


## API route
Format: `/api/{version}/{entity/module}/{subentity/submodule}/{action}/{params...}`  
Examples:   
`/api/v1/service/create`  
`/api/v1/service/creds/add`  
`/api/v2/service/delete/0065d2f7-0222-4eb2-a993-c6983a0517fe`

To get a single entity `/{entity}/{entityID}`
To get all entities `/{entity}/all`

Example:  
`/api/v2/service/delete/0065d2f7-0222-4eb2-a993-c6983a0517fe`  
`/api/v2/service/all`


Read API should have GET method, all other APIs which change state of server should have POST method



## Errors and logging
1. Instead of handling errors from the lowest call stack (eg. database functions), handle them in HTTP handler functions where it makes sense.
2. If there is more than one query in a single database function, wrap errors to give a context.
3. Do not directly return err values to an HTTP response.


## Packages
All packages follow same file structure.
##### init.go
init.go contains interfaces, InitStore function and State variable
##### *hndl.go
handlers contain http handler functions. eg hndlauth.go, hndlcrud.go
##### *store.go
store contains statefull methods of Store struct. eg crudstore.go




#### Commit message syntax
```
{type}:[{scope}]  {summary}    

{description}
```

Type: One of these:  
1. fix : Bug fix
2. feat : New feature 
2. cfeat : Change feature (breaking changes)
3. ref : Refractor (no change in structure/unit tests)
4. ntc : Nothing to commit

Scope: Scope of change. example: service/creds.go, dbstore

Summary: One sentence summary of change  

Description: Long description    



## Pull Request Process

1. Ensure unit test and integration-test passes.
2. Use go fmt to format your code
2. If you have added exported functions/packages, update wiki and comments.
4. You may merge the Pull Request in once you have the sign-off of two other developers, or if you 
   do not have permission to do that, you may request the second reviewer to merge it for you.

