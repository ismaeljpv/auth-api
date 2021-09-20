-- auth_api.users definition

DROP TABLE IF EXISTS `users`;
CREATE TABLE `users` (
  `id` int NOT NULL AUTO_INCREMENT,
  `firstname` varchar(100) NOT NULL,
  `lastname` varchar(100) NOT NULL,
  `email` varchar(100) NOT NULL,
  `password` varchar(256) NOT NULL,
  `status` varchar(10) NOT NULL,
  `createdOn` mediumtext,
  PRIMARY KEY (`id`)
) ENGINE=InnoDB AUTO_INCREMENT=2 DEFAULT CHARSET=utf8mb3;

INSERT INTO auth_api.users (firstname, lastname, email, password, status, createdOn) VALUES('Ismael', 'Pena', 'ip@gmail.com', '$2a$10$jR2qQUJx9b/CiC81KtPbXu5/1rwOA7AjQzTxDXNnzv4/IweV14mj6', 'ACTIVE', '1631906445');
INSERT INTO auth_api.users (firstname, lastname, email, password, status, createdOn) VALUES('Cristian', 'Vargas', 'cjpv@gmail.com', '$2a$10$jR2qQUJx9b/CiC81KtPbXu5/1rwOA7AjQzTxDXNnzv4/IweV14mj6', 'ACTIVE', '1631906445');