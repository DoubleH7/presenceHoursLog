# presenceHoursLog

this is a warmup project for go development and uses basic authentication for admin actions

one can create users and monitor their working hours.
users can record their presence hours by starting or stopping their sitting sessions

## Setup requirements
You'll need to run `go build` to create your own .exe file. prior to that you'll need the following
## .env file
The directory requires a .env file 
- specifying the mongodb connection url as "DB_URI"
- specifying the server connection port as "PORT"
---

supported handlers:

- GET request to server root
> returns a string to show that server is up and running

- GET request to /admin/users/all
> returns a json of all users 

> do include basic authentication for this request

> the __admins__ collection from the __presenceLog__ database contains all valid username password combinations

- POST request to /admin/users
> do include a json body as the example suggests:

> `{
  "name" :"Hesam",
  "age": 24
}`

> do include basic authentication for this request

> the __admins__ collection from the __presenceLog__ database contains all valid username password combinations

- GET request to /admin/users/id/{USER ID HERE}
> returns the complete user info

> do include basic authentication for this request

> the __admins__ collection from the __presenceLog__ database contains all valid username password combinations

- GET request to /admin/users/name/{USER NAME HERE}
> returns the complete user info

> do include basic authentication for this request

> the __admins__ collection from the __presenceLog__ database contains all valid username password combinations

- POST request to /start/{USER ID HERE}
> starts a session for the user

- POST request to /stop/{USER ID HERE}
> stops the session for the user


