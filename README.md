# Olives

A Reddit clone. The back-end will be written in go while the front-end will be made using React and TypeScript. The database will MariaDB and to interface with the database I will use sqlc. Will need to use JWT library to handle authentication of users and bcrypt to hash the passwords. To style the front-end I will use Tailwind. The images will be stored locally. Each forum will have a limit of posts and when a post reaches a limit of replies, the post is deleted. When the limit of posts is reached, the oldest post will be deleted to make room for a new one. The maximum image size will 1 megabyte. I plan to learn how to use TypeScript and Tailwind in the front-end and use test driven development for the back-end. Subreddits will be called subolives.

## TODO

- create middleware for endpoints that should only be accessed by admins
- use sessions. create methods:
	- new session
	- delete session
	- get session
	- get user from session?

## Dependencies Summary

- go
	- sqlc
	- bcrypt
- mariadb
- react
	- Nextjs
	- TypeScript
	- Tailwind

## Front-end pages

- /                        home
- /subolive/:id            for subolives
- /subolive/post/:id       for a post in a subolive
- /login                   for the user to login
- /user/:id                for user info

## Back-end endpoints

   ENDPOINT                        METHOD   HANDLER  DONE   Description

- /getSubolivePosts/:id?page=x     GET       done            to get all the posts in a subolive to show in the subolive page (create version of this without including comments for performance or just be lazy and keep this one)
- /getPost/:id                     GET       done            to see the info of a post
- /deletePost/:id                  DELETE    done            to delete a post (for admins)
- /createPost                      POST      done            to create a new post
- /createComment/:postId           POST      next            to create a comment
- /deleteCooment/:postId           DELETE                    to delete a comment (for admins)
- /getUser/:userId                 GET       done            to get user info
- /createUser                      POST      done            to create a new user
- /deleteUser/:userId              DELETE    done            to delete an existent user
- /login                           POST                      to login
- /logout                          GET                       to logout

## Back-end helper functions

- bcrypt
- middlewares for auth
- Image handling (creating and deleting)
- deleting posts. need to decide whether to implement a reaper that is excecuted every ten minutes or some other time interval or to excecute a function that checks the conditions for deletion every time a request is made to the server. when:
	- a new post wants to be created and the number of maximum posts has been reached
	- the limit of comments has been reached in a post
- maybe add search function to be able to search for posts








