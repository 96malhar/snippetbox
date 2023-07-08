Table is created via - https://www.tablesgenerator.com/markdown_tables

| Method | Pattern            | Handler           | Action                                         |
|--------|--------------------|-------------------|------------------------------------------------|
| GET    | /                  | home              | Display the home page                          |
| GET    | /about             | about             | Display the about page                         |
| GET    | /snippet/view/{id} | snippetView       | Display a specific snippet                     |
| GET    | /snippet/create    | snippetCreate     | Display a HTML form for creating a new snippet |
| POST   | /snippet/create    | snippetCreatePost | Create a new snippet                           |
| GET    | /user/signup       | userSignup        | Display a HTML form for signing up a new user  |
| POST   | /user/signup       | userSignupPost    | Create a new user                              |
| GET    | /user/login        | userLogin         | Display a HTML form for logging in a user      |
| POST   | /user/login        | userLoginPost     | Authenticate and login the user                |
| POST   | /user/logout       | userLogoutPost    | Logout the user                                |
| GET    | /static/*          | http.FileServer   | Serve a specific static file                   |
| GET    | /ping              | ping              | Return a 200 OK response                       |
| GET    | /account/view      | accountView       | Returns account details of the user            |
