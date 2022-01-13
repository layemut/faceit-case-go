# faceit-case-go
---
#### to start-up the application

``` bash
docker-compose up --build
```
---
#### contains following operations; 
+ Create User
+ Modify User
+ Remove User
+ List User
---
#### Application uses mongodb for persistence storage
---
#### It has management route for application health check
Example response:
```
{
    "mongo": {
        "status": "UP"
    },
    "app": {
        "status": "UP"
    }
}
```