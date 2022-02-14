# dsn2mycnf
Converts MySQL DSN connection string to a my.cnf configuration file.

### Getting started

Installation:

```
go install rafal.dev/dsn2mycnf
```

Usage:

```
$ dsn2mycnf 'user:password@tcp(host:3306)/database'
[client]
  host = "host"
  port = 3306
  user = "user"
  password = "password"
  database = "database"
  ssl-mode = "PREFERRED"
```

### Example

Generate client configuration file `database.cnf`:

```
$ dsn2mycnf -out database.cnf 'user:password@tcp(host:3306)/database'
```

Start mysql client container:

```
$ docker run -ti --rm --mount type=bind,source=$PWD/database.cnf,target=/etc/mysql/conf.d/database.cnf mysql mysql -A
```
```
Welcome to the MySQL monitor.  Commands end with ; or \g.
Your MySQL connection id is 970989
Server version: 5.7.33-log Source distribution

Copyright (c) 2000, 2022, Oracle and/or its affiliates.

Oracle is a registered trademark of Oracle Corporation and/or its
affiliates. Other names may be trademarks of their respective
owners.

Type 'help;' or '\h' for help. Type '\c' to clear the current input statement.

mysql>
```
