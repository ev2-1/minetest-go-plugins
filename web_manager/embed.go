package main

import (
	"embed"

	"strings"
)

//go:embed html/*
//go:embed js/*
//go:embed css/*
var files embed.FS

//go:embed packets.txt
var packetfile string

var pkts = "pkts " + strings.ReplaceAll(packetfile, "\n", ",")
