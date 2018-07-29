# Ansible Inventory for Hetzner Cloud
[![CircleCI](https://circleci.com/gh/artheus/hcloud-ansible-inv/tree/master.svg?style=svg)](https://circleci.com/gh/artheus/hcloud-ansible-inv/tree/master)
[![Go Report Card](https://goreportcard.com/badge/github.com/artheus/hcloud-ansible-inv)](https://goreportcard.com/report/github.com/artheus/hcloud-ansible-inv)
[![GoDoc](https://godoc.org/github.com/artheus/hcloud-ansible-inv?status.svg)](https://godoc.org/github.com/artheus/hcloud-ansible-inv)

Automate your [Hetzner Cloud](https://www.hetzner.de/cloud) instances by using a dynamic inventory script for [Ansible](https://github.com/ansible/ansible).

![See it in action](https://github.com/thannaske/hcloud-ansible-inv/raw/master/example.png)

## Getting Started
See [Getting Started](https://github.com/thannaske/hetzner-cloud-ansible-inventory/wiki/Getting-Started) in the repository's wiki. Here you will find always up-to-date installation instructions as well as remarks concerning the configuration and usage of the inventory script.

## Usage
You are able to use the within your Ansible commands using the `-i` flag.

`HETZNER_CLOUD_KEY=example ansible -i hcloud-ansible-inv all -m ping`

This command should execute the Ansible ping module and should return a pong for each server you are running at Hetzner Cloud.
Please consult [Ansible's documentation](http://docs.ansible.com) for further resources concerning the usage of Ansible itself.

## License
This project is open source (MIT License). For more information see [LICENSE](https://github.com/thannaske/hcloud-ansible-inv/blob/master/LICENSE).

## Acknowledgements
This project is using the [Hetzner Cloud API Client](https://github.com/hetznercloud/hcloud-go) and [jeffail's Gabs](https://github.com/Jeffail/gabs) (painless JSON processing).
