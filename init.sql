DROP TABLE IF EXISTS `gifs`;
DROP TABLE IF EXISTS `users`;
DROP TABLE IF EXISTS `wxUsers`;
DROP TABLE IF EXISTS `tgUsers`;
DROP TABLE IF EXISTS `gifGroups`;


CREATE TABLE `gifs`
(
    `id` INTEGER NOT NULL AUTO_INCREMENT,
    `GroupID` INTEGER  NULL DEFAULT NULL,
    `FileID` CHAR(40) NOT NULL,
    `UserID` INTEGER NOT NULL,
    PRIMARY KEY (`id`)
);

CREATE TABLE `wxUsers`
(
    `id` INTEGER NOT NULL AUTO_INCREMENT,
    `openID` CHAR(40) NOT NULL,
    `nickName` CHAR(40) NOT NULL DEFAULT "",
    PRIMARY KEY (`id`),
    UNIQUE (`openID`)
);

CREATE TABLE `tgUsers`
(
    `id` INTEGER NOT NULL DEFAULT 0,
    PRIMARY KEY (`id`)
);

CREATE TABLE `users`
(
    `id` INTEGER NOT NULL AUTO_INCREMENT,
    `tgUserID` INTEGER NOT NULL DEFAULT 0,
    `wxUserID` INTEGER  NOT NULL,
    PRIMARY KEY (`id`)
);


CREATE TABLE `gifGroups`
(
    `id` INTEGER NOT NULL AUTO_INCREMENT,
    `name` CHAR(40) NOT NULL DEFAULT "",
    PRIMARY KEY (`id`)
);

ALTER TABLE `gifs` ADD FOREIGN KEY (GroupID) REFERENCES `gifGroups` (`id`);
ALTER TABLE `gifs` ADD FOREIGN KEY (UserID) REFERENCES `users` (`id`);
ALTER TABLE `users` ADD FOREIGN KEY (tgUserID) REFERENCES `tgUsers` (`id`);
ALTER TABLE `users` ADD FOREIGN KEY (wxUserID) REFERENCES `wxUsers` (`id`);


INSERT INTO tgUsers (id) VALUES (0);
INSERT INTO 