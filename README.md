# GEOGAME
>* Geogame is a backend that serves parts of an API for a location based game using Go. 
>* The backend consists of two services: ​Players ​and ​Locations
 There are two endpoints: ​Admin ​and ​Client
>* A client can ​Register​and ​Log In t​o the game with a username and password. Clients are authenticated using a ​JSON Web Token.
>* Clients with valid authentication can perform an action (eg. send their location, change their username) as a demonstration of working authentication.
>* A Location​​ consists of geographic coordinates and some metadata (eg. name, type). Admins can ​create, edit, read ​and​ delete​ locations.​
Clients can r​ead​ locations.
>* For data storage posggres is used, but also considered to make it simple to switch to a different storage solution.
Similarly, the API is served over HTTP.
>* Admin endpoints don’t have authentication, what options exist to protect them from being publicly accessible?

# Pre requisites 
- Mac(10.14)
- Goland 2019.2
- Go 1.13
- Docker  2.1.0.4
# Build & Run
> Application  is designed to run as a full stack which means all controller layer,business logic layer and persistent layers are containerised as single application.
>* Application can be build and started by Makefile.
>* Make sure to cd to project folder.
>* Run the below commands in the terminal shell.Make sure to follow the order

# How to start the app
  
    make up  
    make migrate

# How to stop the app    
	make down
	
# How to run unit test
    make test_unit

# How to run build
    make build
    
# How to create docker image
     make image
     
# Admin Endpoint info
**Create Location**
----
  Returns ok.
    
    `{"Ok":"success"}`

* **URL**

  /admin/loc/create

* **Method:**

  `POST`
  
* **Data Params**

   **Required:**
   
     `{"id":"1","geoPoint": {"longitude":19.2,"latitude":58.1},"metaData":{"locationName":"Stockholm","locationType":"city"}}`

* **Success Response:**

  * **Code:** 200 <br />
 
* **Error Response:**

  * **Code:** 401 UNAUTHORIZED <br />
    **Content:** `{ error : "You are unauthorized to make this request." }`

* **Sample Call:**

  `curl -X POST "http://localhost:8080/v1/admin/loc/create" -d '{"id":"1","geoPoint": {"longitude":19.2,"latitude":58.1},"metaData":{"locationName":"Stockholm","locationType":"city"}}'`
  
**Get Location**
----
  Returns output.
  `{"id":"1","geoPoint":{"longitude":19.2,"latitude":58.1},"metaData":{"locationName":"Paris","locationType":"city"}}`
  
* **Sample Call:**

    `curl -X GET "http://localhost:8080/v1/admin/loc/1"`

**Update Location**
----
  Returns ok.
  
  `{"Ok":"success"}`
  
* **Sample Call:**

    `curl -X PUT "http://localhost:8080/v1/admin/loc/update" -d '{"id":"1","geoPoint": {"longitude":19.2,"latitude":58.1},"metaData":{"locationName":"Paris","locationType":"city"}}'`

**Delete Location**
----
 Returns ok.
   
   `{"Ok":"success"}`
  
* **Sample Call:**

    `curl -X DELETE "http://localhost:8080/v1/admin/loc/1/delete"`
    
# Client Endpoint info

**Register client**
----

Returns ok.
  
  `{"Ok":"success"}`
  
* **Sample Call:**
                
`curl -X POST "http://localhost:8080/v1/client/register" -d '{"email":"dummy+test@gmail.com","name":"dummy fullname","password":"password"}'`

**Login client**
----

Returns token.

`{"token":"eyJhbGciOiJIUzI1NiIsInR5cCI6IkpXVCJ9.eyJ1c2VySWQiOiJkZDcxMTdjYy0zNDg4LTQzZmEtOWNjMS00NjBjNjY4ZTM4N2UiLCJleHAiOjE1OTI4MzM2OTAsImp0aSI6IjAiLCJpYXQiOjE1OTI4MzI3OTAsImlzcyI6Imdlby1nYW1lIiwibmJmIjoxNTkyODMyNzg5fQ.eoEC1LIsVxasMbsEKJHZOzmwuDTtF0ORSExNwL6FzXM"}`
  
* **Sample Call:**
                
`curl -X POST "http://localhost:8080/v1/client/login" -d '{"email":"dummy+test@gmail.com","password":"password"}'`

**Update Name**
----

Returns ok.
  
  `{"Ok":"success"}`
  
* **Sample Call:**
                
`curl -X PUT "http://localhost:8080/v1/client/update-name" -d '{"Name":"updated fullname"}' -H 'Authorization: Bearer ${Bearer token}'`

**Send location**
----

Returns ok.
  
  `{"Ok":"success"}`
  
* **Sample Call:**
                
`curl -X POST "http://localhost:8080/v1/client/loc/send" -d '{"id":"1","geoPoint": {"longitude":19.2,"latitude":58.1},"metaData":{"locationName":"Stockholm","locationType":"city"}}' -H 'Authorization: Bearer ${Bearer token}'`

**Get Location**
----
  Returns output.
  `{"id":"1","geoPoint":{"longitude":19.2,"latitude":58.1},"metaData":{"locationName":"Paris","locationType":"city"}}`
  
* **Sample Call:**

    `curl -X GET "http://localhost:8080/v1/client/loc/get" -H 'Authorization: Bearer ${Bearer token}'`

* ***NOTE:***
    Use bearer token which is part of login response
    
## Technical info
* kartoza/postgis container is used to perform GIS operation
* golang/alpine container is used
* Oauth 2(JWT) is used for client endpoint authentication


## Improvement
* Admin endpoints don't have authentication.Certificate authentication can be added to protect in public env.
* Right now location objects are stored in two different table.Instead keep the location in Locations table,make the clients to reference the location from Players table to Locations table.