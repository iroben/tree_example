CREATE TABLE `group`
(
    `id`   int(11) unsigned NOT NULL AUTO_INCREMENT,
    `name` varchar(32) NOT NULL,
    `pid`  int(11) NOT NULL DEFAULT '0',
    `lval` int(11) unsigned NOT NULL DEFAULT '0',
    `rval` int(11) unsigned NOT NULL DEFAULT '0',
    PRIMARY KEY (`id`)
) ENGINE = InnoDB
  DEFAULT CHARSET = utf8;