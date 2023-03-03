---
title: "HashiCorp Boundary and Nomad"
author: "Nic Jackson"
slug: "boundary_nomad"
---

## Consul UI

To access the Consul UI point your browser at:

`http://1-consul-server.container.shipyard.run:8500`

## Boundary UI

The Boundary UI can be accessed at the following location:

`http://boundary.container.shipyard.run:9200`

### Default Login Credentials
* Username: nicj
* Password: password

## Default Auth Method and Scope

A default auth method, organization scope, and project scope is created automatically,
these identifiers can be read from the following files:

* org_id - ./files/org_id
* project_id - ./files/project_id
* auth_method_id - ./files/auth_method_id

## Nomad UI

Nomad is running on a random port, to deterimine the Nomad port
use the following command.

```shell
shipyard output NOMAD_ADDR
```
