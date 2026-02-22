
# Autonomics  
A proof of concept of the Autonomic Water Network Architecture.

---

# Compilation

Build the project using:

```bash
go build .\main.go
```

This compiles both the type checker and the interpreter.

Tested with:

- Go version: `go1.25.0`
- Platform: `windows/amd64`

---

# Running

Execute:

```bash
./main.exe
```

This runs the WDN simulator together with the homeostasis and reconfiguration scenarios.

---

# Resources

## Core Files

- `main.sess`  
  Runs the program by deploying the `iCPSDLproxy`, `eventManager`, and `agentDeployment` processes.

- `wdnSimulator`  
  Implements the `sessions` processes that simulate the water distribution network.

- `homeostasis.sess`  
  Implements the `sessions` processes running the homeostasis control loop.

- `reconf.sess`  
  Implements the `sessions` processes running the reconfiguration control loop.

---

## Example Folder

The `example` folder contains configuration files used to build the iCPS-DL knowledge graph:

- `wdn.dss`  
  Schema for describing water distribution networks in iCPS-DL.

- `simple.dss`  
  Description of the water distribution network use case.

- `agents.dss`  
  Repository containing MPST agent descriptions.

- `faulty.dss`  
  Description of the water distribution network without sensor `s6` (fault at sensor `s6`).

- `faulty2.dss`  
  Description of the water distribution network without device `dev2` (fault at device `dev2`).