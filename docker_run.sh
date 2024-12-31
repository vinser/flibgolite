#!/usr/bin/env bash

# the script runs flibgolite docker container

# run this script from the root of the app folder

# put your books in the "books" folder
books_dir=$(pwd)/books
mkdir -p $books_dir

# here you can view logs
logs_dir=$(pwd)/logs
mkdir -p $logs_dir

# here is your books index database. Back up it if you want
dbdata_dir=$(pwd)/dbdata
mkdir -p $dbdata_dir


docker run -d \
--name=flibgolite \
-p 8085:8085 \
--mount type=bind,source=$books_dir,target=/flibgolite/books \
--mount type=bind,source=$dbdata_dir,target=/flibgolite/dbdata \
--mount type=bind,source=$logs_dir,target=/flibgolite/logs \
vinser/flibgolite:latest