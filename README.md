# Faster Incident Response (Faster IR, FaIR)

As a junior minion at [Fluency Security](https://www.fluencysecurity.com/ "Fluency Corporation's Homepage") (yes, the interns), I aim to improve the Fluency Security product while having fun and being goofy.

One feature we want to implement was an improved version of certsocietegenerale's [Fast Incident Response (FIR)](https://github.com/certsocietegenerale/FIR).

Therefore, we introduce to you, the improved version of FIR, Faster Incident Response, or FaIR.

## Purpose

Faster Incident Response is to be a stand alone incident response tracking system.  With the use of Security, Orcheteration and Response (SOAR) systems, tracking of incidents will no longer be kept in the Security Information and Event Management (SIEM) system.

## Installation

To install FaIR:

* Download and install Git. Installation instructions for each operating system are found [here](https://git-scm.com/downloads "Download Git")
* Install the following:
  * Golang (also known as Go). This is the language used to program this application. For instructions, visit Golang's website for [download instructions](https://golang.org/dl/) and [installation instructions](https://golang.org/doc/install).
  * Gin-Gonic's Gin web framework. Use this command per their [instructions](https://github.com/gin-gonic/gin#installation): `go get -u github.com/gin-gonic/gin`
  * MongoDB, the system's current database for storing all objects. Please refer to [MongoDB's website](https://docs.mongodb.com/manual/installation/) for proper instructions for your system.
  * Redis (for user authorizations, this will almost definitely be migrated to MongoDB)
* Clone (download) the project to your system by running `git clone https://github.com/SecurityDo/fsb.git`

To start FaIR, open two terminal windows:

* In the first window, run the `mongod` command to start the MongoDB daemon
* In the second window
  * Run `go build` to build the project into a binary
  * Run `.fir` to execute the binary
* Visit localhost:8080, and your local version of FaIR is now running!

## Development

This project aims to use a correct Open RESTful API, using common tools, such as Golang and MongoDB.
Development will involving unit testing using Golang's native unit testing framework.
In the end, this project will allow other SIEM projects to quickly integrate solutions with a low chance of dependency- and platform-specific build issues.

### Status

The project is in the early stages of development. The current objectives are:

* [ ] User accounts: Authorizations
* [ ] Web interface: incorporate reactivity using Vue.js (NOTE: this task might take some time)
* [ ] Unit Testing

Features already implemented include:

* [x] Definitions of the base schemas for the database
* [x] RESTful API to the object of that scheme (current)
* [x] Implement basic CRUD functionalities
* [x] Authentication: salting passwords
* [x] Add a (very basic) site

### Project Structure

The project uses the following directory structure for development:

* Application (main.go) runs from the root directory
* Database models reside in `models`
* Gin routes reside in `routes`
* The code that bridges the route to the database, including all CRUD operations, are kept in `controllers`
* Basic functionalities directly involving the database are put in `mongo`
* Unit testing scripts are to be placed in `tests`.  Tests are to end in `_test.go`

### Data Models

Data models are used to store different types of objects -- incidents, users, artifacts, files -- in the database, which is currently a MongoDB database. An online, well-documented API has yet to be added.
A documented API of a similar product (called simple incident response) can be accessed on Postman [here](https://documenter.getpostman.com/view/1117493/SVSHs9hS?version=latest).

#### Incidents

This is the json for an incident object that gets used by /api/incidents:

```go
{
  "detection": string,
  "actor": string,
  "plan": string,
  "fileset": []string,
  "datecreated": int64, // automatically tracked
  "lastmodified": int64, // automatically tracked
  "isstarred": bool,
  "subject": string,
  "description": string,
  "severity": int,
  "isincident": bool,
  "ismajor": bool,
  "status": string,
  "confidentiality": int,
  "category": string,
  "openedby": string,
  "concernedbusinesslines": []string
}
```

#### Users

This is the json for an incident object that gets used by /api/users:

```go
{
  "username": string,
  "password": string,
  "email": string,
  "url": string,
  "groups": []string
}
```

#### Artifacts

This is the json for an incident object that gets used by /api/artifacts:

```go
{
  "name": string,
  "artifact": struct {
    "description": string,
    "datecreated": int64, // automatically tracked
    "lastmodified": int64 // automatically tracked
  }
}
```

#### Files

This is the json for an incident object that gets used by /api/files:

```go
{
  "datecreated": int64, // automatically tracked
  "lastmodified": int64, // automatically tracked
  "filepath": string
}
```
