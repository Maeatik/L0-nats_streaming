package main

import (
	"flag"
	"log"
	"os"

	_ "github.com/lib/pq"
	"github.com/nats-io/nats.go"
)

func usage() {
	log.Printf("Usage: nats-pub [-s server] [-creds file] [-nkey file] [-tlscert file] [-tlskey file] [-tlscacert file] <subject> <msg>\n")
	flag.PrintDefaults()
}

func showUsageAndExit(exitcode int) {
	usage()
	os.Exit(exitcode)
}

func main() {


	var urls = flag.String("s", nats.DefaultURL, "The nats server URLs (separated by comma)")
	var userCreds = flag.String("creds", "", "User Credentials File")
	var nkeyFile = flag.String("nkey", "", "NKey Seed File")
	var tlsClientCert = flag.String("tlscert", "", "TLS client certificate file")
	var tlsClientKey = flag.String("tlskey", "", "Private key file for client certificate")
	var tlsCACert = flag.String("tlscacert", "", "CA certificate to verify peer against")
	var reply = flag.String("reply", "", "Sets a specific reply subject")
	var showHelp = flag.Bool("h", false, "Show help message")
	//var modelShow = flag.String("model", "", "Show simple model database record")
	log.SetFlags(0)
	flag.Usage = usage
	flag.Parse()

	if *showHelp {
		showUsageAndExit(0)
	}

	args := flag.Args()
	if len(args) != 2 {
		showUsageAndExit(1)
	}

	opts := []nats.Option{nats.Name("NATS Sample Publisher")}

	//if *modelShow != "" {
	//	opts = append(opts, nats.Token(string(getModel(urls))))
	//}
	if *userCreds != "" && *nkeyFile != "" {
		log.Fatal("specify -seed or -creds")
	}

	// Use UserCredentials
	if *userCreds != "" {
		opts = append(opts, nats.UserCredentials(*userCreds))
	}



	// Use TLS client authentication
	if *tlsClientCert != "" && *tlsClientKey != "" {
		opts = append(opts, nats.ClientCert(*tlsClientCert, *tlsClientKey))
	}

	// Use specific CA certificate
	if *tlsCACert != "" {
		opts = append(opts, nats.RootCAs(*tlsCACert))
	}

	// Use Nkey authentication.
	if *nkeyFile != "" {
		opt, err := nats.NkeyOptionFromSeed(*nkeyFile)
		if err != nil {
			log.Fatal(err)
		}
		opts = append(opts, opt)
	}


	sc, err := nats.Connect(*urls, opts...)
	if err != nil {
		log.Fatal(err)
	}
	defer sc.Close()

	subj, msg := args[0], []byte(args[1])
	if reply != nil && *reply != "" {
		sc.PublishRequest(subj, *reply, msg)
	} else {
		sc.Publish(subj, msg)
	}

	sc.Flush()

	if err := sc.LastError(); err != nil {
		log.Fatal(err)
	} else {
		log.Printf("Published [%s] : '%s'\n", subj, msg)
	}
}
