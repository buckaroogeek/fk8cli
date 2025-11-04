# `fk8cli` - Fedora Kubernetes Utility

Golang based utility to install Kubernetes via `kubeadm` on Fedora machines.

`fk8cli` is a command line client packaged in an rpm (COPR initially) for
installation via dnf. The client is installed on each node in the cluster
and executed via remote shell. It can also function well in shell scripts.

Inspired by C. Collicutt's [go-install-kubernetes](https://github.com/ccollicutt/go-install-kubernetes)
utility. I package Kubernetes, CRI-O and a few other related components
for Fedora. In order to test a Kubernetes rpm I usually follow steps outlined
on [Fedora Quick Docs](https://docs.fedoraproject.org/en-US/quick-docs/using-kubernetes-kubeadm/).
I also have an ansible role and playbooks for more automation. Yet I have
not been fully satisfied with these approachs - time consuming and/or
requiring more software on the VM than I was happy with. In searching for
alternatives I ran across C. Collicutt's bash script and then the golang
implementation. This seems to be a useful approach for my purposes and
may be useful to other Fedora users that want to explore basic Kubernetes.

## Scope

Install Kubernetes and CRI-O on a Fedora machine using rpms provided by
the distribution. Initialize a node on using the
appropriate configuration for a control plane or worker node (or a
node with both Kubernetes roles).

## Stack

### Current stack

Initial stack consists of:

* Fedora machine (vm or bare metal).
* Kubernetes rpms from the Fedora repository.
* `kubeadm` (kubernetesx.xx-client rpm) to initialize the cluster.
* CRI-O as the CRI implementation.
* `crun` for the container runtime.
* `flannel` for CNI.

### Future enhancements

* Containerd.
* Calico or other CNI.

## Desired behavior

1. Configure a machine as a node
1. node is the object to configure
1. settings package defines a struct of node
1. is the node a control-plane, worker or both?
1. version of kubernetes and cri-o to install
1. future - what CNI or CRI to use?
1. future - what runtime to use?
1. app configuration: log file, verbosity

## Execution flow

1. Execute: `fkutil [flags] [optional k8s version]`
1. Process flags
1. ~~Look for `fkutil.yaml` (local directory, location via flag)~~
1. ~~Process optional yaml file~~
1. Execute as configured via flags~~/yaml~~
1. Call appInit - initialize fk8cli based on flags
1. Call criInstall - install CRI (current cri-o)
1. Call k8sInstall - installs k8s rpms
1. Call nodeInit - initializes node
1. Call appClose - any clean up

## Command line options

Command line: fk8cli [flags] k8s-version

E.G. fk8cli -c 1.34

1. -c Configure as control-plane, not exclusive (default)
1. -w Configure as worker, not exclusive
1. ~~-s Configure as single (both control-plane and worker)~~
1. -d Configure as dual (both control-plane and worker)
1. -q Quiet - do not echo log file content to standard out
1. ~~-k value --k8s=value    Set Kubernetes version (e.g. 1.33)~~
1. ~~-f value --file=value   Set yaml file location (path or URL)~~
1. ~~version                 fk8cli Version~~
1. -h --help               help including current version
1. k8s-version   a Semver string with or without v or V prefix
