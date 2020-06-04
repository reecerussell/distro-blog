CREATE SCHEMA `distro-blog-test` ;

CREATE TABLE `distro-blog-test`.`table-one` (
  `id` INT NOT NULL AUTO_INCREMENT,
  `name` VARCHAR(10) NOT NULL,
  `age` INT(8) NOT NULL,
  PRIMARY KEY (`id`),
  UNIQUE INDEX `name_UNIQUE` (`name` ASC) VISIBLE);
