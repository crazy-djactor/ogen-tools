package main

import (
	"archive/tar"
	"bufio"
	"compress/gzip"
	"context"
	"crypto/tls"
	"encoding/json"
	"flag"
	"fmt"
	"github.com/go-git/go-git/v5"
	"github.com/go-git/go-git/v5/plumbing"
	"github.com/multiformats/go-multiaddr"
	"github.com/olympus-protocol/ogen-tools/launcher/config"
	"github.com/olympus-protocol/ogen/bls"
	ogenconf "github.com/olympus-protocol/ogen/config"
	"github.com/olympus-protocol/ogen/params"
	"github.com/olympus-protocol/ogen/primitives"
	"github.com/olympus-protocol/ogen/proto"
	"github.com/sethvargo/go-password/password"
	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials"
	"io"
	"io/ioutil"
	"log"
	"math/rand"
	"net/http"
	"os"
	"os/exec"
	"path"
	"path/filepath"
	"strconv"
	"strings"
	"sync"
	"time"
)

// RPCClient represents an RPC connection to a server.
type RPCClient struct {
	address   string
	conn      *grpc.ClientConn
	Available bool
	Network   proto.NetworkClient
	Utils     proto.UtilsClient
}

func init() {
	// We initialize ogen bls module with testnet params
	err := bls.Initialize(params.TestNet)
	if err != nil {
		panic(err)
	}
}

var datadir = "./data/"

var ogenSubFolderPrefix = "ogen-node-"

var genesisTime = time.Unix(time.Now().Unix()+120, 0)

var premineAccount = bls.RandKey()

func main() {
	log.Println("Loading Configuration")
	c := loadConfig()

	log.Println("Creating Folder Structure")
	err := folders(c)
	if err != nil {
		log.Fatal(err)
	}

	if c.Source {
		log.Println("Building Ogen")
		err = buildOgen(c)
		if err != nil {
			log.Fatal(err)
		}
	} else {
		log.Println("Downloading Ogen")
		err = downloadOgen()
		if err != nil {
			log.Fatal(err)
		}
	}

	log.Println("Generate Validators")
	var validators []primitives.ValidatorInitialization
	var wg sync.WaitGroup
	for i := 1; i <= c.Nodes; i++ {
		wg.Add(1)
		go func(index int, wg *sync.WaitGroup) {
			defer wg.Done()
			v := genValidators(index, c.Password, c.Validators)
			validators = append(validators, v...)
		}(i, &wg)
	}
	wg.Wait()

	premineAccountString, err := premineAccount.PublicKey().ToAccount()
	if err != nil {
		log.Fatal(err)
	}

	chain := ogenconf.ChainFile{
		Validators:         validators,
		GenesisTime:        uint64(genesisTime.Unix()),
		InitialConnections: nil,
		PremineAddress:     premineAccountString,
	}

	log.Println("Generating and copying chain file")
	err = generateChainFile(c, chain)
	if err != nil {
		log.Fatal(err)
	}

	log.Printf("Starting %v ogen instances with %v validators \n", c.Nodes, c.Nodes*c.Validators)
	var iwg sync.WaitGroup
	iwg.Add(1)

	go func() {
		err = runInstances(c, &iwg)
		if err != nil {
			log.Fatal(err)
		}
	}()

	log.Println("Connecting nodes and start proposers")

	_, extMa := startChain(c)

	log.Println("Creating chain file for external usage")

	var initialConnections []string
	for _, ma := range extMa {
		initialConnections = append(initialConnections, ma.String())
	}
	chain.InitialConnections = initialConnections

	marshal, err := json.Marshal(&chain)
	if err != nil {
		log.Println("Error:" + err.Error())
	}
	err = ioutil.WriteFile("./chain.json", marshal, 0777)
	if err != nil {
		log.Println("Error:" + err.Error())
	}

	log.Println("Network ready!")
	premineWif, err := premineAccount.ToWIF()
	if err != nil {
		log.Fatal(err)
	}
	log.Printf("The premine private key is: %s. This key is not attached to any node. \n", premineWif)

	iwg.Wait()
}

func loadConfig() config.Config {
	var pass, externalHost, branch string
	var nodes, validators int
	var debug, source bool

	flag.StringVar(&pass, "password", "", "Password for keystore and wallet")
	flag.StringVar(&externalHost, "host", "127.0.0.1", "IP of the external host to use on chain file")
	flag.IntVar(&nodes, "nodes", 5, "Setup the amount of nodes the testnet (minimum of 5 nodes)")
	flag.IntVar(&validators, "validators", 32, "Define the amount of validators per node (default 32 nodes)")
	//flag.BoolVar(&debug, "debug", false, "Use this flag to start nodes on debug mode")
	flag.BoolVar(&debug, "debug", true, "Use this flag to start nodes on debug mode")
	flag.BoolVar(&source, "source", false, "Use this flag to build from source")
	flag.StringVar(&branch, "branch", "master", "When using the source you can specify a branch to build from")
	flag.Parse()

	if pass == "" {
		pass, _ = password.Generate(32, 10, 0, false, false)
	}

	c := config.Config{
		Password:     pass,
		Nodes:        nodes,
		Validators:   validators,
		ExternalHost: externalHost,
		Debug:        debug,
		Source:       source,
		Branch:       branch,
	}
	return c
}

func folders(c config.Config) error {

	_ = os.RemoveAll(datadir)

	err := os.Mkdir(datadir, 0777)
	if err != nil {
		return err
	}

	for i := 1; i <= c.Nodes; i++ {
		err := os.Mkdir(path.Join(datadir, ogenSubFolderPrefix+strconv.Itoa(i)), 0777)
		if err != nil {
			return err
		}
	}

	return nil
}

func buildOgen(c config.Config) error {
	_ = os.RemoveAll("./bin")
	_ = os.RemoveAll("./ogen")

clone:
	_, err := git.PlainClone("./ogen", false, &git.CloneOptions{
		URL:           "https://github.com/olympus-protocol/ogen",
		Progress:      os.Stdout,
		ReferenceName: plumbing.NewBranchReferenceName(c.Branch),
	})
	if err != nil {
		// If the repo already exists, delete it and clone again
		if err == git.ErrRepositoryAlreadyExists {
			_ = os.RemoveAll("./ogen")
			goto clone
		}
		return err
	}
	cmd := exec.Command("scripts/build.sh")
	cmd.Dir = "./ogen"
	cmd.Stdin = os.Stdin
	cmd.Stdout = os.Stdout
	err = cmd.Run()
	if err != nil {
		return err
	}

	err = os.Mkdir("./bin", 0777)
	if err != nil {
		return err
	}

	err = os.Rename("./ogen/ogen", "./bin/ogen")
	if err != nil {
		return err
	}

	err = os.Chmod("./bin/ogen", 0777)
	if err != nil {
		return err
	}

	return nil
}

func RemoveContents(dir string) error {
	d, err := os.Open(dir)
	if err != nil {
		return err
	}
	defer d.Close()
	names, err := d.Readdirnames(-1)
	if err != nil {
		return err
	}
	for _, name := range names {
		err = os.RemoveAll(filepath.Join(dir, name))
		if err != nil {
			return err
		}
	}
	return nil
}

func downloadOgen() error {
	_ = os.RemoveAll("./bin")

	file := "https://public.oly.tech/olympus/ogen-release/ogen-0.0.1-linux-amd64.tar.gz"
	resp, err := http.Get(file)
	if err != nil {
		return err
	}

	defer resp.Body.Close()

	_ = os.Mkdir("./bin", 0777)

	err = extractTar(resp.Body)
	if err != nil {
		return err
	}

	err = os.Rename("./ogen-0.0.1/ogen", "./bin/ogen")
	if err != nil {
		return err
	}

	//err = os.Remove("./ogen-0.0.1")
	err = RemoveContents("./ogen-0.0.1")
	if err != nil {
		return err
	}

	err = os.Chmod("./bin/ogen", 0777)
	if err != nil {
		return err
	}

	return nil
}

func extractTar(stream io.Reader) error {
	log.Println("Extracting Ogen")

	uncompressedStream, err := gzip.NewReader(stream)
	if err != nil {
		return err
	}

	tarReader := tar.NewReader(uncompressedStream)

	for true {

		header, err := tarReader.Next()

		if err == io.EOF {
			break
		}

		if err != nil {
			return err
		}

		switch header.Typeflag {
		case tar.TypeDir:
			if err := os.Mkdir(header.Name, 0755); err != nil {
				return err
			}

		case tar.TypeReg:
			outFile, err := os.Create(header.Name)
			if err != nil {
				return err
			}
			if _, err := io.Copy(outFile, tarReader); err != nil {
				return err
			}
			err = outFile.Close()
			if err != nil {
				return err
			}

		default:
			return err
		}
	}
	return nil
}

func genValidators(index int, password string, amount int) []primitives.ValidatorInitialization {

	var v []primitives.ValidatorInitialization
	dataDirAbsPath, err := filepath.Abs(path.Join(datadir, ogenSubFolderPrefix+strconv.Itoa(index)))
	if err != nil {
		log.Println("Error: " + err.Error())
		return nil
	}

	payee, err := genWallet(dataDirAbsPath, password)
	if err != nil {
		log.Println("Error: " + err.Error())
		return nil
	}

	cmd := exec.Command("./ogen", "--datadir="+dataDirAbsPath, "generate", strconv.Itoa(amount))

	p, err := filepath.Abs("bin/")
	if err != nil {
		log.Println("Error: " + err.Error())
		return nil
	}

	cmd.Dir = p

	stdout, err := cmd.StdoutPipe()
	if err != nil {
		return nil
	}

	cmd.Stderr = cmd.Stdout

	done := make(chan struct{})

	scanner := bufio.NewScanner(stdout)

	go func() {
		for scanner.Scan() {
			line := scanner.Text()
			data := strings.Split(line, ":")
			if data[0] == "Public Key" {
				val := primitives.ValidatorInitialization{
					PubKey:       strings.Replace(data[1], " ", "", -1),
					PayeeAddress: payee,
				}
				v = append(v, val)
			}
		}

		done <- struct{}{}
	}()

	err = cmd.Start()
	if err != nil {
		log.Println("Error:" + err.Error())
		return nil
	}

	<-done

	err = cmd.Wait()
	if err != nil {
		log.Println("Error:" + err.Error())
		return nil
	}

	return v
}

func genWallet(dataDir string, password string) (string, error) {
	cmd := exec.Command("./ogen", "--datadir="+dataDir, "wallet", "validators", "testnet", password)

	p, err := filepath.Abs("bin/")
	if err != nil {
		return "", err
	}
	cmd.Dir = p

	stdout, err := cmd.StdoutPipe()
	if err != nil {
		return "", err
	}

	cmd.Stderr = cmd.Stdout

	done := make(chan struct{})

	scanner := bufio.NewScanner(stdout)
	var payee string
	go func() {

		for scanner.Scan() {
			line := scanner.Text()
			data := strings.Split(line, ":")
			payee = strings.Replace(data[1], " ", "", -1)
		}

		done <- struct{}{}

	}()

	err = cmd.Start()
	if err != nil {
		log.Println("Error:" + err.Error())
		return "", err
	}

	<-done

	err = cmd.Wait()
	if err != nil {
		log.Println("Error:" + err.Error())
		return "", err
	}
	return payee, nil
}

func generateChainFile(conf config.Config, c ogenconf.ChainFile) error {
	marshal, err := json.Marshal(c)
	if err != nil {
		return err
	}
	for i := 1; i <= conf.Nodes; i++ {
		err = ioutil.WriteFile(path.Join(datadir, ogenSubFolderPrefix+strconv.Itoa(i), "chain.json"), marshal, 0777)
		if err != nil {
			return err
		}
	}
	return nil
}

func runInstances(c config.Config, gwg *sync.WaitGroup) error {
	var wg sync.WaitGroup

	for i := 1; i <= c.Nodes; i++ {
		wg.Add(1)

		go func(index int, wg *sync.WaitGroup) {
			defer wg.Done()

			dataDirAbsPath, err := filepath.Abs(path.Join(datadir, ogenSubFolderPrefix+strconv.Itoa(index)))
			if err != nil {
				log.Println("Error: " + err.Error())
				return
			}

			var cmd *exec.Cmd
			if c.Debug {
				cmd = exec.Command("./ogen", "--debug", "--datadir="+dataDirAbsPath, "--log_file", "--port="+strconv.Itoa(24000+index), "--rpc_port="+strconv.Itoa(25000+index))
			} else {
				cmd = exec.Command("./ogen", "--datadir="+dataDirAbsPath, "--log_file", "--port="+strconv.Itoa(24000+index), "--rpc_port="+strconv.Itoa(25000+index))
			}

			p, err := filepath.Abs("bin/")
			if err != nil {
				log.Println("Error: " + err.Error())
				return
			}

			cmd.Dir = p

			stdout, err := cmd.StdoutPipe()
			if err != nil {
				return
			}

			cmd.Stderr = cmd.Stdout

			done := make(chan struct{})

			scanner := bufio.NewScanner(stdout)

			go func() {

				for scanner.Scan() {
					line := scanner.Text()
					fmt.Print(line)
				}

				done <- struct{}{}

			}()

			err = cmd.Start()
			if err != nil {
				log.Println("Error:" + err.Error())
				return
			}

			<-done

			err = cmd.Wait()
			if err != nil {
				log.Println("Error:" + err.Error())
				return
			}

		}(i, &wg)
	}
	wg.Wait()
	gwg.Done()
	return nil
}

func startChain(c config.Config) (local []multiaddr.Multiaddr, external []multiaddr.Multiaddr) {
	var peerAddr, externalAddr []multiaddr.Multiaddr

	var fwg sync.WaitGroup
	fwg.Add(c.Nodes)

	// Get all node IDs
	for i := 1; i <= c.Nodes; i++ {
		go func(wg *sync.WaitGroup, index int) {
			defer wg.Done()
			client := rpcClient(index)
		get:
			netInfo, err := client.Network.GetNetworkInfo(context.Background(), &proto.Empty{})
			if err != nil {
				goto get
			}
			maL, err := multiaddr.NewMultiaddr("/ip4/127.0.0.1/tcp/" + strconv.Itoa(24000+index) + "/p2p/" + netInfo.ID)
			if err != nil {
				return
			}
			maE, err := multiaddr.NewMultiaddr("/ip4/" + c.ExternalHost + "/tcp/" + strconv.Itoa(24000+index) + "/p2p/" + netInfo.ID)
			if err != nil {
				return
			}
			peerAddr = append(peerAddr, maL)
			externalAddr = append(externalAddr, maE)
			_ = client.Close()
		}(&fwg, i)
	}

	fwg.Wait()

	// Connect nodes between them

	var swg sync.WaitGroup

	swg.Add(c.Nodes)
	for i := 1; i <= c.Nodes; i++ {
		go func(index int, wg *sync.WaitGroup) {
			defer wg.Done()
			client := rpcClient(index)
			peers := getRandomPeers(peerAddr)
			for _, p := range peers {
				_, err := client.Network.AddPeer(context.Background(), &proto.IP{Host: p.String()})
				if err != nil {
					continue
				}
			}
			_ = client.Close()
		}(i, &swg)
	}

	swg.Wait()

	// Start the block proposers

	var twg sync.WaitGroup
	twg.Add(c.Nodes)
	for i := 1; i <= c.Nodes; i++ {
		go func(index int, wg *sync.WaitGroup) {
			defer wg.Done()
			client := rpcClient(index)
			_, err := client.Utils.StartProposer(context.Background(), &proto.Empty{})
			if err != nil {
				return
			}
		}(i, &twg)

	}

	twg.Wait()

	return peerAddr, externalAddr
}

func getRandomPeers(peers []multiaddr.Multiaddr) []multiaddr.Multiaddr {
	randPeers := make([]multiaddr.Multiaddr, 8)
	for i := range randPeers {
		r := rand.Intn(len(peers))
		randPeers[i] = peers[r]
	}
	return randPeers
}

// NewRPCClient creates a new RPC client.
func rpcClient(nodeNum int) *RPCClient {
	c := &tls.Config{
		InsecureSkipVerify: true,
	}
	conn, err := grpc.Dial("127.0.0.1:"+strconv.Itoa(25000+nodeNum), grpc.WithTransportCredentials(credentials.NewTLS(c)))
	if err != nil {
		panic("unable to connect to rpc server")
	}
	client := &RPCClient{
		conn:    conn,
		address: "127.0.0.1:" + strconv.Itoa(25000+nodeNum),
		Network: proto.NewNetworkClient(conn),
		Utils:   proto.NewUtilsClient(conn),
	}
	return client
}

func (rpc *RPCClient) Close() error {
	return rpc.conn.Close()
}
