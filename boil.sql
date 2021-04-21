CREATE TABLE boil_user
(
    id        INT PRIMARY KEY AUTO_INCREMENT NOT NULL,
    username  VARCHAR(100)                   NOT NULL,
    password  BINARY(32)                     NOT NULL,
    avatar_id INT                            NOT NULL,
    salt      CHAR(5)                        NOT NULL,
    bio       VARCHAR(100)                   NOT NULL
);

CREATE TABLE boil_user_follow_user
(
    id          INT PRIMARY KEY AUTO_INCREMENT NOT NULL,
    follower_id INT                            NOT NULL,
    user_id     INT                            NOT NULL,
    FOREIGN KEY (follower_id) REFERENCES boil_user (id) ON DELETE CASCADE,
    FOREIGN KEY (user_id) REFERENCES boil_user (id) ON DELETE CASCADE
);


CREATE TABLE boil_tag
(
    id    INT PRIMARY KEY AUTO_INCREMENT NOT NULL,
    title VARCHAR(10)                    NOT NULL
);

CREATE TABLE boil_boil
(
    id          INT PRIMARY KEY AUTO_INCREMENT NOT NULL,
    tag_id      INT                            NOT NULL,
    user_id     INT                            NOT NULL,
    create_time DATETIME                       NOT NULL,
    content     VARCHAR(300)                   NOT NULL,
    FOREIGN KEY (tag_id) REFERENCES boil_tag (id) ON DELETE CASCADE,
    FOREIGN KEY (user_id) REFERENCES boil_user (id) ON DELETE CASCADE
);


CREATE TABLE boil_comment
(
    id          INT PRIMARY KEY AUTO_INCREMENT NOT NULL,
    boil_id     INT                            NOT NULL,
    user_id     INT                            NOT NULL,
    create_time DATETIME                       NOT NULL,
    content     VARCHAR(100)                   NOT NULL,
    FOREIGN KEY (boil_id) REFERENCES boil_boil (id) ON DELETE CASCADE,
    FOREIGN KEY (user_id) REFERENCES boil_user (id) ON DELETE CASCADE
);

CREATE TABLE boil_user_like_boil
(
    id      INT PRIMARY KEY AUTO_INCREMENT NOT NULL,
    user_id INT                            NOT NULL,
    boil_id INT                            NOT NULL,
    FOREIGN KEY (user_id) REFERENCES boil_user (id) ON DELETE CASCADE,
    FOREIGN KEY (boil_id) REFERENCES boil_boil (id) ON DELETE CASCADE
);