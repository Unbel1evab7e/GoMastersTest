# GoMastersTest
This repository is implementation of simple api
# Enpoints
## GetAllUsers

Summary      Get All Users
 
Description  Return All Users
 
Tags         Users
 
Produce      json
 
Success      200  {object}  []entity.User
 
Failure      400  {object}  httputil.HTTPError
 
Failure      404  {object}  httputil.HTTPError
 
 Router /users [get]

## GetUserByID
Summary      Get User By id

Description  Return User By ID

Tags         Users

Accept       json

Produce      json

Param        id path string  true  "User Id"

Success      200  {object}  entity.User

Failure      400  {object}  httputil.HTTPError

Failure      404  {object}  httputil.HTTPError

Router /users/{id} [get]

## CreateUser 
Summary CreateUser

Description  Return Created User

Tags         Users

Accept       json

Produce      json

Param        User  body     DTOs.User  true  "Add User"

Success      201  {object}  string

Failure      422  {object}  httputil.HTTPError

Failure      404  {object}  httputil.HTTPError

Router /users [post]

## UpdateUser 
 Summary      UpdateUser
 
 Description  UpdateUser
 
 Tags         Users
 
 Accept       json
 
 Produce      json
 
 Param        id    path     string  true  "User ID"
 
 Param        User  body     DTOs.User  true  "Update user"
 
 Success      200  {object}  entity.User
 
 Failure      400  {object}  httputil.HTTPError
 
 Failure      422  {object}  httputil.HTTPError
 
 Router       /users/{id} [put]
 
 ## DeleteUser 
 Summary 	 DeleteUser
 
 Description  DeleteUser
 
 Tags         Users
 
 Accept       json
 
 Produce      json
 
 Param        id   path      string  true  "User ID"
 
 Success      204  {object}  string
 
 Failure      400  {object}  httputil.HTTPError
 
 Failure      404  {object}  httputil.HTTPError
 
 Router /users/{id} [delete]
