# Azure DevOps Pipelines Agent Registrator

This is a simple CLI helper tool that registers fake/offline Azure Pipelines Agents against the Azure Pipelines API
(see https://learn.microsoft.com/en-us/rest/api/azure/devops/distributedtask/agents/add?view=azure-devops-rest-7.0).
It sets the _system_ capability `ExtraAgentContainers` to a specific string

It was created as a companion tool for https://github.com/MShekow/azure-pipelines-k8s-agent-scaler to register fake
agents, so that a follow-up job (that demands specific capabilities) has the chance to start. This works around the
Azure DevOps Pipelines limitation that jobs with demands for which there is no _registered_ agent will immediately be
cancelled.

## Usage

Example for Linux/UNIX:

```bash
./agent-registrator \
  -organization-url https://dev.azure.com/foobar \
  -pool-name your-azure-devops-pool-name \
  -pat <Azure DevOps Personal Access Token with 'Agent Pools Read&Manage' permission> \
  -agent-name-prefix dummy-agent \
  -extra-agent-contaners 'ubuntu,registry.hub.docker.com/library/ubuntu:22.04,250m,64Mi'
```

You can download the release from the _GitHub Releases_ page of this project.
