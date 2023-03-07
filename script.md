# Boundary Demo

Providing remote access to applications and systems requires secure routing 
to the destination and credentials to authenticate the user. Traditionally, 
you achieve this using a Virtual Private Network (VPN) or a Bastion server to 
bridge into the private network. Credentials are generally provided individually, 
created as part of a manual process, and with password rotation on a best-intention 
basis. This is problematic as access is usually too broad, difficult to audit, 
and complex to maintain. 

In a zero-trust world, access is granted from point to point, not to the network 
edge; credentials are unique to the session, and everything is fully auditable. 
HashiCorp Boundary and Vault provide this solution giving you greater control
over access to your systems.

In this talk, Nic will walk you through the steps needed to configure Boundary 
and Vault showing you how to provide secure access to typical cloud-based systems 
like Kubernetes and virtual machines.

## Flow
* Explain problem with existing Bastion setup
* SSH
  - Problem 1: You need Access without Bastion
    - Show how to deploy and configure boundary worker using username password
    - Explain how to use Vault plugin to provide worker auth replacing hardcoded
      details
  - Problem 2: You need SSH creds
    - Show how to configure one-time access for the ssh server using PAM
      to generate credentials
    - Boundary can inject creds, but it needs access to Vault, show how to register
      a boundary worker for Vault
    - Show how to create a credentials store
    - Show how to inject creds
* Database
  - Problem 1: How to access the database
    - Show how to configure and inject dynamic db credentials / separated by role
* Nomad Job (census)
  - Problem 1: How to access workloads running in highly dynamic environments
    the issue is not access but managing targets.
    - Show how to run a boundary Worker in Nomad
    - Start an application on Nomad, show dynamic ports
    - Register a target
    - Re-allocate, watch everything turn to dust port and location has changed
    - Show Census to manage dynamic targets
  - Problem 2: You are using consul service mesh
    - Show pattern where Boundary worker is running as a sidecar