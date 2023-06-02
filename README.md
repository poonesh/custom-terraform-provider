
# High Level Explanation

This repo contains a basic example of implementing a custom terraform provider. When run, it will simply populate a postgres database depending on the provided terraform resources. This example consists of 4 parts: a server, a client, and a provider, all written in golang, and the example terraform code which uses the provider.

## Server

The server is a simple api which accepts simple REST api requests and inserts, deletes or modifies rows in a database. 

## Client

The client is a simple client which sends requests to the server.

## Provider

The provider works alongside terraform core to parse and interpret terraform code, and call the client methods appropritately. Under the hood, terraform core also manages the state file.

## Example Terraform Code

The example .tf files show how we can now write new terraform resources, which will be interpreted by the provider.


# Installation and Running

## 1. Create Database

Start by [installing postgres](https://www.postgresql.org/download/) and ensure it is running.
With brew you can do this by
```
sudo brew services start postgres
```
The following command will create the foods table used by the server.
```
cd server
psql postgres -f foods.sql;
cd ..
```

## 2. Run the server

Open `server/main.go` and edit the user and password for your postgres database.

To run the server, 

```
cd server
go mod init
go mod vendor
go run main.go
cd ..
```


## 3. Build and install the provider

If running on an OS other than MacOS, you will need to change the OS_ARCH variable in provider/Makefile.

Then, to install the provider, run

```
cd provider
make
cd ..
```

This will build and install the provider in a place terraform can find it.


## 4. Run the terraform code

Now that the provider is installed, we can run the terraform code in the examples directory. The provided example will create a few resources, each which will add a row to the database table created in step 1. Commenting the resource out and re-running the terraform will cause these rows to be deleted.

```
cd examples
terraform init
terraform apply
cd ..
```
