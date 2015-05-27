# conntrack-logger

## What is this?

A very basic tool for collecting new conntrack entries and storing them in the database.

As it was written for solving private particular task, lot of stuff is hardcoded and non-configurable, sorry.

That was my first expirience with Go, so don't expect it to be coded in a good way.

### Internals

It directly uses libnetlink/libnetfilter for conntrack events processing. As those libraries
do not provide any ways to directly parse nf events, events are converted into xml form by library call
and the parsed on the go side.

## Configuring

conntrack-logger is configured with a json file you specify with the -config flag:

`conntrack-logger -config yourstuff.json`

Sample config file is ditributed with source code, see conntrack-logger.cfg.

It it almost self descriptive, as only database connection details are configured. 
The CommitCount parameter controls size of a single batch.

## Building it

1. Install [go](http://golang.org/doc/install)

2. Install "lib/pq" go get -u github.com/lib/pq

3. Install libnetfilter-conntrack-dev and libnfnetlink-dev

4. Compile conntrack-logger

        git clone git://github.com/akashihi/conntrack-logger.git
        cd conntrack-logger
        go build .

## Packaging it (optional)

To build the packages, you will need ruby and fpm installed.

    gem install fpm

Now you can build a package:

    make deb

Package installs to `/usr/bin`. 

You'll need to install libnetfilter-conntrack and libnfnetlink manually before using it.

## Running it

Generally:

    conntrack-logger -config conntrack-logger.cfg

Option '-config' can be omited, defaul config file location is '/etc/conntrack-logger.cfg'. 

You will also need to install PostgreSQL server, create database, user and table. Table definition 
 is in database.sql file.

## License 

See LICENSE file.

Copyright 2015 Denis V Chapligin <akashihi@gmail.com>
