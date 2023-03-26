CREATE TABLE subolives (
	subolive_id INT AUTO_INCREMENT PRIMARY KEY,
	name VARCHAR(100) NOT NULL
);

CREATE TABLE images (
	image_id INT AUTO_INCREMENT PRIMARY KEY,
	file_path VARCHAR(255) NOT NULL
);

CREATE TABLE users (
	user_id INT AUTO_INCREMENT PRIMARY KEY,
	email VARCHAR(255) NOT NULL UNIQUE,
	username VARCHAR(255) NOT NULL UNIQUE,
	password VARCHAR(255) NOT NULL,
	admin BOOL NOT NULL
);

CREATE TABLE posts (
	post_id INT AUTO_INCREMENT PRIMARY KEY,
	title VARCHAR(255) NOT NULL,
	text VARCHAR(1275) NOT NULL, 
	created_at TIMESTAMP NOT NULL DEFAULT CURRENT_TIMESTAMP,
	subolive_id INT NOT NULL,
	user_id INT NOT NULL,
	image_id INT,
	CONSTRAINT fk_post_subolive
		FOREIGN KEY (subolive_id)
		REFERENCES subolives(subolive_id)
		ON UPDATE CASCADE
		ON DELETE CASCADE,
	CONSTRAINT fk_post_user
		FOREIGN KEY (user_id)
		REFERENCES users(user_id)
		ON UPDATE CASCADE
		ON DELETE CASCADE,
	CONSTRAINT fk_post_image
		FOREIGN KEY (image_id)
		REFERENCES images(image_id)
		ON UPDATE CASCADE
		ON DELETE CASCADE
);

CREATE TABLE comments (
	comment_id INT AUTO_INCREMENT PRIMARY KEY,
	text VARCHAR(1275) NOT NULL, 
	created_at TIMESTAMP NOT NULL DEFAULT CURRENT_TIMESTAMP,
	user_id INT NOT NULL,
	image_id INT,
	post_id INT NOT NULL,
	CONSTRAINT fk_comment_user
		FOREIGN KEY (user_id)
		REFERENCES users(user_id)
		ON UPDATE CASCADE
		ON DELETE CASCADE,
	CONSTRAINT fk_comment_image
		FOREIGN KEY (image_id)
		REFERENCES images(image_id)
		ON UPDATE CASCADE
		ON DELETE CASCADE,
	CONSTRAINT fk_comment_post
		FOREIGN KEY (post_id)
		REFERENCES posts(post_id)
		ON UPDATE CASCADE
		ON DELETE CASCADE
);

CREATE TABLE sessions (
	session_id VARCHAR(100) PRIMARY KEY,
	last_access TIMESTAMP NOT NULL DEFAULT CURRENT_TIMESTAMP,
	user_id INT NOT NULL,
	CONSTRAINT fk_session_user
		FOREIGN KEY (user_id)
		REFERENCES users(user_id)
		ON UPDATE CASCADE
		ON DELETE CASCADE
)












