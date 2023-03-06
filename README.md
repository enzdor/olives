# Olives

A Reddit clone. The back-end will be written in go while the front-end will be made using React and TypeScript. The database will MariaDB and to interface with the database I will use sqlc. Will need to use JWT library to handle authentication of users and bcrypt to hash the passwords. To style the front-end I will use Tailwind. The images will be stored locally. Each forum will have a limit of posts and when a post reaches a limit of replies, the post is deleted. When the limit of posts is reached, the oldest post will be deleted to make room for a new one. The maximum image size will 1 megabyte. I plan to learn how to use TypeScript and Tailwind in the front-end and use test driven development for the back-end. Subreddits will be called subolives.

## TODO

maybe try and abstract some of the error handling?

## Dependencies Summary

- go
	- sqlc
	- JWT
	- bcrypt
- mariadb
- react
	- TypeScript
	- Tailwind

## Front-end pages

- /                        home
- /subolive/:id            for subolives
- /subolive/post/:id       for a post in a subolive
- /login                   for the user to login
- /user/:id                for user info

## Back-end endpoints

   ENDPOINT                        METHOD   DONE   Description

- /getSubolivePosts/:id?page=x     GET             to get all the posts in a subolive to show in the subolive page (create version of this without including comments for performance or just be lazy and keep this one)
- /getPost/:id                     GET             to see the info of a post
- /createPost/:suboliveId          POST            to create a new post
- /getComments/:postId             GET             to get the comments when displaying a post (i don't think this will be necessary)
- /createComment/:postId           POST            to create a comment
- /getUser/:userId                 GET             to get user info
- /createUser                      POST            to create a new user
- /deleteUser/:userId              POST            to delete an existent user
- /login                           POST            to login
- /logout                          GET             to logout

## Back-end helper functions

- JWT
- bcrypt
- middlewares for auth
- Image handling (creating and deleting)
- deleting posts when:
	- a new post wants to be created and the number of maximum posts has been reached
	- the limit of comments has been reached in a post
- maybe add search function to be able to search for posts








