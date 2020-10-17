## Development Setup

- Clone the repo  
 `git clone https://github.com/seknox/trasa.git && cd trasa`
- Download [GeoCityLite.mmdb](https://github.com/seknox/trasa/raw/master/build/etc/trasa/static/GeoLite2-City.mmdb) and move it to `/etc/trasa/static/GeoCityLite.mmdb`
- Install dashboard dependencies   
`yarn install` inside dashboard directory
- Start dashboard in dev mode
`yarn start`
- Build backend server
`go build` inside server directory
- Start database
`sudo docker run -d -p 5432:5432 --name db -e POSTGRES_PASSWORD=trasauser -e POSTGRES_USER=trasauser -e POSTGRES_DB=trasadb postgres`
-Start redis
`sudo docker run -d -p 6379:6379 --name redis redis`
- Run the binary
`sudo ./server`

- Edit `/etc/trasa/config/config.toml` file and   
    - change `trasa.proxyDashboard` to ` true`
    - change `trasa.dashboardAddr` to `http://localhost:3000`
- Open TRASA dashboard at https://localhost

> Go through the [wiki](https://github.com/seknox/trasa/wiki) to get overview of codebase. 


# Code conventions 


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
##### h*.go
handlers contain http handler functions. eg hAuth.go, hTfa.go
##### store*.go
store contains stateful methods of Store struct. eg storeSession.go




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

Scope: Scope of change. example: service/creds.go, dbstore

Summary: One sentence summary of change  

Description: Long description    



## Pull Request Process

1. Ensure unit test and integration-test passes.
2. Use go fmt to format your code
2. If you have added exported functions/packages, update wiki and comments.
4. You may merge the Pull Request in once you have the sign-off of one other developer, or if you 
   do not have permission to do that, you may request the second reviewer to merge it for you.

