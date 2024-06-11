/* Copyright 2017-2020 Victor Penso, Matteo Dessalvi, Joeri Hermans

This program is free software: you can redistribute it and/or modify
it under the terms of the GNU General Public License as published by
the Free Software Foundation, either version 3 of the License, or
(at your option) any later version.

This program is distributed in the hope that it will be useful,
but WITHOUT ANY WARRANTY; without even the implied warranty of
MERCHANTABILITY or FITNESS FOR A PARTICULAR PURPOSE.  See the
GNU General Public License for more details.

You should have received a copy of the GNU General Public License
along with this program.  If not, see <http://www.gnu.org/licenses/>. */

package main

import (
	"flag"
	"github.com/prometheus/client_golang/prometheus"
	"github.com/prometheus/client_golang/prometheus/promhttp"
	"log"
	"net/http"
	"prometheus-slurm-exporter/account"
	"prometheus-slurm-exporter/allocation"
	"prometheus-slurm-exporter/compute"
	"prometheus-slurm-exporter/schedule"
	"strconv"
)

func init() {
	// Metrics have to be registered to be exposed
	prometheus.MustRegister(account.NewAccountsCollector())      // from accounts.go
	prometheus.MustRegister(compute.NewCPUsCollector())          // from compute.go
	prometheus.MustRegister(allocation.NewNodesCollector())      // from nodes.go
	prometheus.MustRegister(allocation.NewNodeCollector())       // from node.go
	prometheus.MustRegister(allocation.NewPartitionsCollector()) // from partitions.go
	prometheus.MustRegister(schedule.NewQueueCollector())        // from queue.go
	prometheus.MustRegister(schedule.NewSchedulerCollector())    // from scheduler.go
	prometheus.MustRegister(allocation.NewFairShareCollector())  // from sshare.go
	prometheus.MustRegister(account.NewUsersCollector())         // from users.go
}

var listenAddress = flag.String(
	"listen-address",
	":8080",
	"The address to listen on for HTTP requests.")

var gpuAcct = flag.Bool(
	"gpus-acct",
	false,
	"Enable GPUs accounting")

func main() {
	flag.Parse()

	// Turn on GPUs accounting only if the corresponding command line option is set to true.
	if *gpuAcct {
		prometheus.MustRegister(compute.NewGPUsCollector()) // from gpus.go
	}

	// The Handler function provides a default handler to expose metrics
	// via an HTTP server. "/metrics" is the usual endpoint for that.

	log.Println("Starting Server: %s", *listenAddress)
	log.Println("GPUs Accounting: %b", strconv.FormatBool(*gpuAcct))
	http.Handle("/metrics", promhttp.Handler())
	log.Fatal(http.ListenAndServe(*listenAddress, nil))
	//http.ListenAndServe(*listenAddress, nil)
}
