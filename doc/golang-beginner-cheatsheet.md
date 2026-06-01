# Cheatsheet for Golang beginners

## 1 - Create startup repository

See [**[first demo application in go]**](https://github.com/RobertTC32/example-demo_hello#)
  
### 1.1 - Setup tools locally
  
- install git and configure global git settings for developer:  
  pull with rebase ensures that local changes and commits,  
  can only be added and never change remote commits;  
``` 
git config --global user.email "robertthecoder32@gmail.com"  
git config --global user.name "RobertTC32"
git config --global init.defaultbranch "main"  
git config --global pull.ff "only"  
# For windows workstations, add:  
git config --global core.autocrlf "true"  
git config --global core.editor "notepad"  
# For linux workstations, add:  
#   git config --global core.autocrlf "false"  
#   git config --global core.editor "nano"  
```
- install ssh and generate ssh public + private key pair for developer:  
``` 
ssh-keygen -t ed25519 -C "robertthecoder32@gmail.com"  
```
- install vscode and add "Git Graph" extension

### 1.2 - Create repository in Github/Gitea

- configure account in github:  
  choose "Settings" in top-right menu,  
  and change "Profile picture" in "Public profile";  
  choose "Settings" in top-right menu,  
  click "New SSH key" in "SSH and GPG keys",  
  and add public SSH key with name "pub-key-XYZ" for user;  
  in Gitea, language can be changed to English  
  in "User Settings -> Appearance -> Language";  
- create repo with "public" visibility,  
  and "example-", "util-" or "app-" prefix in the name;  
  this allows minimal organisation because no grouping exist in github;  
  create "example-demo_hello" repo as example
- while creating the repo,  
  choose to automatically add a README, gitignore (for Go) and MIT license file;  
- change "Settings" in Gitea repo:  
  keep "main" as name for "Default branch",  
  disable the following "Features": "Wikis", "Issues" and "Projects";  
  change the following "Features":  
  set "Pull request permissions" to "Collaborators only" in "Pull requests";  
  this simplfies the menu for the repo in github;  
- change branching strategy in "Settings" in Gitea repo:  
  only allow "Allow rebase merging" for "Pull Requests";  
  no merge commits nor squash merging are allowed,   
  and head branchis are not deleted after merging pull requests;  
  this greatly simplifies the commit history while maintaining work organisation;  
- after creating git repository in Github,  
  clone git repo locally,  
  add readme content, commit the change locally, and push this commit to Github;  

### 1.3 - Use a Branching Strategy

This strategy with one shared main branch,  
and multiple feature branches for work in progress:
- prevent loosing work results:  
  protect shared branches, like main branch  
- prevent blocking other developers:  
  isolate ongoing work per person in feature branches  
- catches conflicts early:  
  often rebase & test localy before merge into main branch  
- makes change history easy:  
  untangle merges in linear history,  
  by only branching from main branch and merge back after rebase
- improves code quality:  
  code review & feedback before merge by using pull request  
  (called merge request in Gitlab) 
 

## 2 - Prepare go development

### 2.1 - Setup golang development

- download and install Go SDK:  
  see [**[Go - Download and install]**](https://go.dev/doc/install);  
- initialize folder of "example-demo_hello" repo as go module:  
  use the repo name with github-username as prefix for the module name;  
  the init command will create a "go.mod" file in the repo root folder 
``` 
go version  
go mod init RobertTC32/example-demo_hello
```
- because our repo will not only contain go source files,  
  but also documentation, sql scripts, etc,  
  organise these various content categories in the following folder structure:  
  (this is not the folder structure used as standard by the Go community)  
  "src" folder containing go source files,  
  "src-db" folder containing database related files (like sql scripts),  
  "dist" folder containing go compilation results (like executables) to deploy,  
  "dist-db" folder containing sql scripts to deploy for server database migrations,  
  "doc" folder containing documentation,  
  and "temp" folder containing temporary results;  
  create a ".gitkeep" file in the "src", "src-db" and "doc" folders
- add "Go" (from "Go Team at Google") extension in vscode
- implement HelloWorld as console application:  
  create the "main.go" file in "src" folder,  
  and insert the HelloWorld console implementation;  
  compile and run the application
``` 
# "go run" needs "go.mod" file in current folder,  
# and get folder of "main.go" as parameter
go run ./src
go help
```

### 2.2 - Create automation scripts for development

- install "Makefile" software:  
  makefile (taskfile, justfile are newer alternatives)  
  is used as build system to automate development steps;  
  "make" is aleady installed by default on Linux distros (like Ubuntu),  
  and on Windows it is part of the toolset included in the Git installation;  
  see [**[git-scm]**](https://git-scm.com/install/windows)
- add "Makefile Tools" extension in vscode  
- create and use target scripts in the "Makefile." file:  
  implement targets to compile & run ("run" & "all" as default target),  
  build linux executable ("build" target), and clean dist folder ("clean" target);  
  get environment variables from the ".env" file into Makefile
``` 
make run
make build
make clean
make
```


## 3 - Develop first web applications

See [**[todo application which is used as go example]**](https://github.com/RobertTC32/app-demo_todo)  

### 3.1 - Revamp HelloWorld application

- implement HelloWorld as example web application:  
  change "src/main.go" file and use "net/http" standard library;  
  run web application and test using command "curl http://localhost:8080/"
- change HelloWorld implementation:  
  create "hello.html" file in "src/resources" folder;  
  use "html/template" standard library for web gui logic
- setup Tailwind CSS:  
  download from "https://nodejs.org/en" and install node.js (if not installed);  
  see [**[Get started with Tailwind CSS]**](https://tailwindcss.com/docs/installation/tailwind-cli);  
``` 
npm init
npm install tailwindcss @tailwindcss/cli
echo '@import "tailwindcss";' > src/resources/style.tailwindcss.css
# add line in html file: 
#   <link href="public/style.css" rel="stylesheet">
# add lines in in "Makefile." file: 
#   generate-twc: 
#     npx @tailwindcss/cli -i ./src/resources/style.tailwindcss.css -o ./dist/public/style.css
# add line in "scripts" section in "package.json" file: 
#   "dev": "npx @tailwindcss/cli -i ./src/resources/style.tailwindcss.css -o ./dist/public/style.css --watch"
# start npx watch process:
npm run dev 
```
- use "Tailwind CSS" for styling the web page:  
  add link to stylesheet in "hello.html" file;  
  change background color to "bg-blue-200" and margin to "m-4"
- change HelloWorld implementation:  
  add FileServer to handle files in "dist/public" folder 
``` 
(cd ./dist && go run ../src)  
```
- while creating the repo,  
  we have choosen to add a gitignore (for Go) file;  
  in this file we now add "dist/", "dist-db/", "temp/",  
  and "node_modules/" folders to ignore;  
  files in these folders can be re-generated;  
  also add lines in the readme file to explain the requirements
- install extensions in vscode for web development:  
  "Live Server", "REST Client",  
  "Tailwind CSS IntelliSense", "Tailwind Docs", "Tailwind Fold"  


### 3.2 - Develop primitive Todo application

**Setup and learn configuration, logging and database middleware**
- init node.js, install tailwindcss,  
  copy files from demo_hello in demo_todo,  
  and adapt "README.md", ".env.example", ".env"  and "package.json" files
- create "logger.go" file and add logging and usage of environment variables:  
  import standard "log/slog" package for logging;  
  import external "github.com/joho/godotenv" package,  
  for external configuration using environment variables
``` 
   go get github.com/joho/godotenv
``` 
- install docker locally:  
  On Ubuntu:  
  see [**[Install Docker Engine on Ubuntu]**](https://docs.docker.com/engine/install/ubuntu/)  
  see [**[How to Install Docker on an Ubuntu Server]**](https://www.youtube.com/watch?v=Pa-FrV7-DxI)  
  On Windows:  
  see [**[Install Docker Desktop on Windows]**](https://docs.docker.com/desktop/setup/install/windows-install/)  
  see [**[How to install Docker on Windows - 2025]**](https://www.youtube.com/watch?v=740YJZZu7QY)  
- add docker extensions in vscode:  
  "Docker (Extension Pack) -> Container Tools", "Docker DX",  
- install "MariaDB" software:  
  install and start mariadb using docker;  
  see [**[dockerhub - mariadb]**](https://hub.docker.com/search?q=mariadb&badges=official)
``` 
docker run --name mydb --rm -d -p 3306:3306 -e MARIADB_ROOT_PASSWORD=myrootpw -e MARIADB_DATABASE=demo_todo mariadb:12.2 
docker exec -it mydb mariadb -uroot -p
MariaDB> show databases;
MariaDB> use mydatabase;
MariaDB> show tables;
...
MariaDB> exit
``` 
- add extensions in vscode:  
  "SQLTools" (from Matheus Teixeira);
  _  
  click cylinder button in left panel in vscode,  
  click "Add New Connection" button,  
  click "Search VS Code Marketplace" link,  
  and select & install "SQLTools MySQL/MariaDB/TiDB" driver extension;   
  _  
  click cylinder button in left panel in vscode,  
  click "MariaDB" button to select your database driver,  
  enter "MariaDB-Local" as "Connection name", "Server and Port" as "Connect using",  
  "localhost" as "Server Address", "3306" as "Port", "demo_todo" as "Database",  
  "root" as "Username", "Ask on connect" as "Password mode", "30" as "Connection Timeout",  
  click "SAVE CONNECTION", and click "CONNECT NOW";  
  _  
  select the created connection, enter "myrootpw" as "MariaDB-Local password",  
  click ENTER, enter "show tables;" in sql window,  
  end click "run on active connection";
- open the database connection in vscode;  
  create './src-db/resources/create_todos.sql' (sql scripts) with "CREATE TABLE ..." statement,  
  select file content and execute 'run on active connection';  
  create './src-db/resources/insert_todos.sql' (sql scripts) with "INSERT INTO ..." statements,  
  select file content and execute 'run on active connection';  
``` 
# content of "create_todos.sql" file:
use demo_todo;
CREATE TABLE IF NOT EXISTS todos (
    id int unsigned NOT NULL AUTO_INCREMENT,
    title varchar(255) CHARACTER SET utf8 NOT NULL,
    completed boolean NOT NULL DEFAULT FALSE, 
    created_at datetime NOT NULL,
    completed_at datetime DEFAULT NULL,
    PRIMARY KEY (id)
);
# 
# content of "insert_todos.sql" file:
use demo_todo;
DELETE FROM todos;
INSERT INTO todos (title, completed, created_at, completed_at) 
  VALUES ('Bake a cake', TRUE, STR_TO_DATE('2025-02-18 15:44:04', '%Y-%m-%d %H:%i:%s'), STR_TO_DATE('2025-02-18 16:44:04', '%Y-%m-%d %H:%i:%s'));
INSERT INTO todos (title, completed, created_at, completed_at) 
  VALUES ('Feed the cat', FALSE, STR_TO_DATE('2025-02-18 15:55:04', '%Y-%m-%d %H:%i:%s'), NULL);
INSERT INTO todos (title, completed, created_at, completed_at) 
  VALUES ('Take out the trash', FALSE, STR_TO_DATE('2025-02-18 15:57:10', '%Y-%m-%d %H:%i:%s'), NULL);
SELECT * FROM todos;
``` 
- install mariadb driver for go using "database/sql" package:  
  see [**[Go & SQL databases: docker compose setup and first queries step by step tutorial]**](https://www.youtube.com/watch?v=q9uj1CniRYk);  
  "go.mod" and "go.sum" files will be changed during package installation;  
``` 
go get github.com/go-sql-driver/mysql
``` 

**Use configuration, logging and database in go application**

- add new targets in "Makefile" file for database creation:  
  "test-db" (create temp), "clean-db" (clean distro) and "build-db" (create distro)
- create "src/todo.go" file to implement Todo structure and methods 
- create "src/storage.go" file to implement database access and storage for demo_todo;  
  import "database/sql" and "github.com/go-sql-driver/mysql" in this source file;  
  use "Open(...)", "Close()", and "Query(...)" to access the database;  
  the "FindAllTodos()" returns an array of Todo which contains all table rows from db;  
- create "src/resources/list.html" to implement html/template to show all todos in html table;  
  the "main.go" file implements routing and handler,  
  which calls "FindAllTodos()" and "list.html" template;  
  logging level and path of database can be set in ".env" file


**Implement graceful shutdown in go application**

- refactor out all web server functionality from "main.go" in a seperate file;  
  the "Server" struct in this "server.go" file will not only handle server functionality,  
  but also supports a graceful shutdown of the application;  
  when pressing Ctrl+C to stop the running application,  
  the database can now be safely closed;  
  see [**[Graceful Shutdown in Go: Key Patterns you need to know!]**](https://www.youtube.com/watch?v=UPVSeZXBTxI)  
  the code in this video was copied as is and contains advanced golang logic for handling concurrency;  
  we will learn these advanced golang topics in the future


## 4 - Deploy first web applications

### 4.1 - Containerize web applications

- build van container image for each application:  
  create "dockerfile" file with Alpine as base image;  
  add and run the "build-oci" target in the "Makefile" file;  
  in the resulting container image,  
  the application will be installed in "/app/dist" folder as "main.exe",  
  and listening on port "8080";  
  an environment variable "OCI_INT_PORT" is also set to this internal port;  
  the implementation is expecting an external mariadb database,  
  which is identified using environment variables:  
  DB_PORT, DB_HOST, DB_USER, DB_PASSWORD and DB_NAME (like oracle schema);  
- run the containerized HelloWorld and Todo application to verify locally:  
  add metadata in the docker images with labels (eg for versioning);  
  add and run the "run-oci" target in the "Makefile" file;  
``` 
MAKE_TIME = $(shell date +"%FT%H:%M:%SZ")

build-oci: build
	docker build . -t $(OCI_NAME):latest --label "version=${OCI_VERSION}" --label "build=$(MAKE_TIME)"
	docker image ls | grep $(OCI_NAME) 
	docker image inspect $(OCI_NAME):latest

run-oci: 
	docker run  --name myapp --rm -p $(OCI_PORT):$(OCI_INT_PORT) -e OCI_PORT=$(OCI_PORT) \  
    -e DB_PORT=$(DB_PORT) -e DB_HOST=host.docker.internal -e DB_USER=$(DB_USER) -e DB_PASSWORD=$(DB_PASSWORD) -e DB_NAME=$(DB_NAME) \
    -e LOG_LEVEL=debug $(OCI_NAME):latest

``` 
- the todo application must be able to run   
  inside a container (make run-oci) and outside/without container (make run);  
  to determine the situation, an extra function "IsRunningInDockerContainer()" is implemented in "server.go" file,  
  which will be used to set the correct web server port and to show the correct enduser port  

### 4.2 - Deploy web applications

**Saving and loading oci images to/from compressed files**

There are several ways to transfer Docker images between different machines:   
- share your creations on Docker Hub,  
  making them publicly available to the world
- make images available in your home network,  
  using a more private setup and your own local Docker registry
- move docker images manually between machines,  
  which is sometimes necessary in certain situations,  
  like dealing with firewalls or network restriction
  
As beginner I wanted to directly transfer a Docker image between machines without using a registry.
- save an image from your local registry in a file,  
  and make task "save-oci" for this in makefile
``` 
save-oci:
  mkdir -p ./temp
  docker save -o temp/$(OCI_NAME).tar $(OCI_NAME):latest
  du -sh temp/$(OCI_NAME).tar
``` 
- move the file from one system to another (without relying on a registry);  
  I will use "Remote Development Extension Pack" extension in vs code for this
- load a file as image in your local registry,  
  and make task "load-oci" for this in makefile
``` 
load-oci:
  docker load -i temp/$(OCI_NAME).tar
  docker image ls | grep $(OCI_NAME)
``` 

**Transferring files between machines**
- add ssh extensions in vscode:  
  "Remote Development Extension Pack -> WSL & Dev Containers & Remote - SSH & Remote - Tunnels",  
  "Remote SSH: Editing Configuration Files", "Remote Development", "Remote Explorer"  
- create "config." file to easily access servers using remote ssh in vscode:  
  open "Remote-SSH: Open Configuration File..." in Command Panel
``` 
...
Host home-testappserver
  HostName 192.168.0.94
  User myadmin

Host home-prodappserver
  HostName 192.168.0.95
  User myadmin
```

**Deploying HelloWorld and Todo application on application servers**  

- create bash scripts to easily deploy and run containerized application:  
  "ops/demo_hello-deploy.sh" and "ops/demo_hello-run.sh" in "example-demo_hello" project;  
  "ops/mariadb-deploy.sh", "ops/mariadb-run.sh",  
  "ops/demo_todo-deploy.sh" and "ops/demo_todo-run.sh" in "app-demo_todo" project;  
- transfer compressed image files and ops bash scripts to test application server,  
  enable execution of bash scripts using "chmod +x SCRIPT",   
  and execute "demo_hello-deploy.sh" & "mariadb-deploy.sh" & "demo_todo-deploy.sh"
- make ports for applications available in firewall of the host,  
  and start applications using "demo_hello-run.sh" & "mariadb-run.sh" & "demo_todo-run.sh"
- verify that the applications are running via the browser on your desktop;  
  the applications are only accessible in your home network via following urls:  
  "http://home-testappserver.robertthecoder.org:9001" for demo_hello,  
  and "http://home-testappserver.robertthecoder.org:9003" for demo_todo


## 5 - Create documentation and diagrams

### 5.1 - Create development diagrams

- install extensions in vscode for creating diagrams:  
  "Draw.io Integration": add support for editing/viewing draw.io diagrams;  
  "Markdown Preview Mermaid Support", "Mermaid Export": add support for editing/viewing mermaid diagrams;  
  "PlantUML": add support for editing/viewing plantuml diagrams;  
- create an architecture diagram for the Todo and HelloWorld application:  
  create "doc/app-arch.md" file using Mermaid,  
  and generate & copy png on [**[Convert Mermaid Diagrams to High-Quality PNG Images]**](https://www.mermaidonline.live/mermaid-to-png)


### 5.2 - Create development documents

- install editor extensions in vscode:  
  "Prettier - Code formatter": improve text/code formatting;  
  "Markdown All in One", "markdownlint": improve markdown validation and export;   
  "vscode-pdf": add support for viewing pdf files;  
  "YAML": improve support for yaml files;  
- create documentation for the Todo and HelloWorld application:  
  the current document which you are viewing, was created in markdown,  
  and saved in both repos as "doc/golang-beginner-cheatsheet.md" file


---

**In the following "Cheatsheet for golang advanced" part,  
more golang programming topics will be covered,  
and the implementation of the HelloWorld and Todo application will be completed.**

---
