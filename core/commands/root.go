package commands

import (
	"errors"
	"io"
	"strings"

	oldcmds "github.com/ipfs/go-ipfs/commands"
	lgc "github.com/ipfs/go-ipfs/commands/legacy"
	dag "github.com/ipfs/go-ipfs/core/commands/dag"
	e "github.com/ipfs/go-ipfs/core/commands/e"
	name "github.com/ipfs/go-ipfs/core/commands/name"
	ocmd "github.com/ipfs/go-ipfs/core/commands/object"
	unixfs "github.com/ipfs/go-ipfs/core/commands/unixfs"
	cidenc "gx/ipfs/QmNWQygwYxgz3QzXG2ytTkrHkZ4HnnSh94ASox3JjktFcR/go-cidutil/cidenc"

	"gx/ipfs/QmPTfgFTo9PFr1PvPKyKoeMgBvYPh6cX3aDP7DHKVbnCbi/go-ipfs-cmds"
	logging "gx/ipfs/QmRREK2CAZ5Re2Bd9zZFG6FeYDppUWt5cMgsoUEp3ktgSr/go-log"
	"gx/ipfs/QmSP88ryZkHSRn1fnngAaV2Vcn63WUJzAavnRM9CVdU1Ky/go-ipfs-cmdkit"
	mbase "gx/ipfs/QmekxXDhCxCJRNuzmHreuaT3BsuJcsjcXWNrtV9C8DRHtd/go-multibase"
)

var log = logging.Logger("core/commands")

var ErrNotOnline = errors.New("this command must be run in online mode. Try running 'ipfs daemon' first")

const (
	ApiOption          = "api"
	CidBaseOption      = "cid-base"
	UpgradeCidV0Option = "upgrade-cidv0"
)

var Root = &cmds.Command{
	Helptext: cmdkit.HelpText{
		Tagline:  "Global p2p merkle-dag filesystem.",
		Synopsis: "ipfs [--config=<config> | -c] [--debug=<debug> | -D] [--help=<help>] [-h=<h>] [--local=<local> | -L] [--api=<api>] <command> ...",
		Subcommands: `
BASIC COMMANDS
  init          Initialize ipfs local configuration
  add <path>    Add a file to IPFS
  cat <ref>     Show IPFS object data
  get <ref>     Download IPFS objects
  ls <ref>      List links from an object
  refs <ref>    List hashes of links from an object

DATA STRUCTURE COMMANDS
  block         Interact with raw blocks in the datastore
  object        Interact with raw dag nodes
  files         Interact with objects as if they were a unix filesystem
  dag           Interact with IPLD documents (experimental)

ADVANCED COMMANDS
  daemon        Start a long-running daemon process
  mount         Mount an IPFS read-only mountpoint
  resolve       Resolve any type of name
  name          Publish and resolve IPNS names
  key           Create and list IPNS name keypairs
  dns           Resolve DNS links
  pin           Pin objects to local storage
  repo          Manipulate the IPFS repository
  stats         Various operational stats
  p2p           Libp2p stream mounting
  filestore     Manage the filestore (experimental)

NETWORK COMMANDS
  id            Show info about IPFS peers
  bootstrap     Add or remove bootstrap peers
  swarm         Manage connections to the p2p network
  dht           Query the DHT for values or peers
  ping          Measure the latency of a connection
  diag          Print diagnostics

TOOL COMMANDS
  config        Manage configuration
  version       Show ipfs version information
  update        Download and apply go-ipfs updates
  commands      List all available commands

Use 'ipfs <command> --help' to learn more about each command.

ipfs uses a repository in the local file system. By default, the repo is
located at ~/.ipfs. To change the repo location, set the $IPFS_PATH
environment variable:

  export IPFS_PATH=/path/to/ipfsrepo

EXIT STATUS

The CLI will exit with one of the following values:

0     Successful execution.
1     Failed executions.
`,
	},
	Options: []cmdkit.Option{
		cmdkit.StringOption("config", "c", "Path to the configuration file to use."),
		cmdkit.BoolOption("debug", "D", "Operate in debug mode."),
		cmdkit.BoolOption("help", "Show the full command help text."),
		cmdkit.BoolOption("h", "Show a short version of the command help text."),
		cmdkit.BoolOption("local", "L", "Run the command locally, instead of using the daemon."),
		cmdkit.StringOption(ApiOption, "Use a specific API instance (defaults to /ip4/127.0.0.1/tcp/5001)"),
		cmdkit.StringOption(CidBaseOption, "mbase", "Multi-base to use to encode version 1 CIDs in output."),
		cmdkit.BoolOption(UpgradeCidV0Option, "Upgrade CID version 0 to version 1 in output."),

		// global options, added to every command
		cmds.OptionEncodingType,
		cmds.OptionStreamChannels,
		cmds.OptionTimeout,
	},
}

// commandsDaemonCmd is the "ipfs commands" command for daemon
var CommandsDaemonCmd = CommandsCmd(Root)

var rootSubcommands = map[string]*cmds.Command{
	"add":       AddCmd,
	"bitswap":   BitswapCmd,
	"block":     BlockCmd,
	"cat":       CatCmd,
	"commands":  CommandsDaemonCmd,
	"files":     FilesCmd,
	"filestore": FileStoreCmd,
	"get":       GetCmd,
	"pubsub":    PubsubCmd,
	"repo":      RepoCmd,
	"stats":     StatsCmd,
	"bootstrap": lgc.NewCommand(BootstrapCmd),
	"config":    lgc.NewCommand(ConfigCmd),
	"dag":       lgc.NewCommand(dag.DagCmd),
	"dht":       lgc.NewCommand(DhtCmd),
	"diag":      lgc.NewCommand(DiagCmd),
	"dns":       lgc.NewCommand(DNSCmd),
	"id":        lgc.NewCommand(IDCmd),
	"key":       KeyCmd,
	"log":       lgc.NewCommand(LogCmd),
	"ls":        lgc.NewCommand(LsCmd),
	"mount":     lgc.NewCommand(MountCmd),
	"name":      name.NameCmd,
	"object":    ocmd.ObjectCmd,
	"pin":       lgc.NewCommand(PinCmd),
	"ping":      lgc.NewCommand(PingCmd),
	"p2p":       lgc.NewCommand(P2PCmd),
	"refs":      lgc.NewCommand(RefsCmd),
	"resolve":   lgc.NewCommand(ResolveCmd),
	"swarm":     lgc.NewCommand(SwarmCmd),
	"tar":       lgc.NewCommand(TarCmd),
	"file":      lgc.NewCommand(unixfs.UnixFSCmd),
	"update":    lgc.NewCommand(ExternalBinary()),
	"urlstore":  urlStoreCmd,
	"version":   lgc.NewCommand(VersionCmd),
	"shutdown":  daemonShutdownCmd,
}

// RootRO is the readonly version of Root
var RootRO = &cmds.Command{}

var CommandsDaemonROCmd = CommandsCmd(RootRO)

var RefsROCmd = &oldcmds.Command{}

var rootROSubcommands = map[string]*cmds.Command{
	"commands": CommandsDaemonROCmd,
	"cat":      CatCmd,
	"block": &cmds.Command{
		Subcommands: map[string]*cmds.Command{
			"stat": blockStatCmd,
			"get":  blockGetCmd,
		},
	},
	"get": GetCmd,
	"dns": lgc.NewCommand(DNSCmd),
	"ls":  lgc.NewCommand(LsCmd),
	"name": &cmds.Command{
		Subcommands: map[string]*cmds.Command{
			"resolve": name.IpnsCmd,
		},
	},
	"object": lgc.NewCommand(&oldcmds.Command{
		Subcommands: map[string]*oldcmds.Command{
			"data":  ocmd.ObjectDataCmd,
			"links": ocmd.ObjectLinksCmd,
			"get":   ocmd.ObjectGetCmd,
			"stat":  ocmd.ObjectStatCmd,
		},
	}),
	"dag": lgc.NewCommand(&oldcmds.Command{
		Subcommands: map[string]*oldcmds.Command{
			"get":     dag.DagGetCmd,
			"resolve": dag.DagResolveCmd,
		},
	}),
	"resolve": lgc.NewCommand(ResolveCmd),
	"version": lgc.NewCommand(VersionCmd),
}

func init() {
	Root.ProcessHelp()
	*RootRO = *Root

	// sanitize readonly refs command
	*RefsROCmd = *RefsCmd
	RefsROCmd.Subcommands = map[string]*oldcmds.Command{}

	// this was in the big map definition above before,
	// but if we leave it there lgc.NewCommand will be executed
	// before the value is updated (:/sanitize readonly refs command/)
	rootROSubcommands["refs"] = lgc.NewCommand(RefsROCmd)

	Root.Subcommands = rootSubcommands

	RootRO.Subcommands = rootROSubcommands
}

type MessageOutput struct {
	Message string
}

func MessageTextMarshaler(res oldcmds.Response) (io.Reader, error) {
	v, err := unwrapOutput(res.Output())
	if err != nil {
		return nil, err
	}

	out, ok := v.(*MessageOutput)
	if !ok {
		return nil, e.TypeErr(out, v)
	}

	return strings.NewReader(out.Message), nil
}

type CidBaseHandler struct {
	base           string
	upgrade        bool
	upgradeDefined bool
	args           []string
	enc            *cidenc.Encoder
}

func NewCidBaseHandler(req *cmds.Request) *CidBaseHandler {
	h := &CidBaseHandler{}
	h.base, _ = req.Options[CidBaseOption].(string)
	h.upgrade, h.upgradeDefined = req.Options[UpgradeCidV0Option].(bool)
	h.args = req.Arguments
	return h
}

func NewCidBaseHandlerLegacy(req oldcmds.Request) *CidBaseHandler {
	h := &CidBaseHandler{}
	h.base, _, _ = req.Option(CidBaseOption).String()
	h.upgrade, h.upgradeDefined, _ = req.Option(UpgradeCidV0Option).Bool()
	h.args = req.Arguments()
	return h
}

func (h *CidBaseHandler) UseGlobal() *CidBaseHandler {
	h.enc = &cidenc.Default
	return h
}

func (h *CidBaseHandler) Proc() (*CidBaseHandler, error) {
	var e cidenc.Encoder = cidenc.Default
	if h.base != "" {
		var err error
		e.Base, err = mbase.EncoderByName(h.base)
		if err != nil {
			return h, err
		}
	}

	e.Upgrade = h.upgrade
	if h.base != "" && !h.upgradeDefined {
		e.Upgrade = true
	}

	if h.enc == nil {
		h.enc = &cidenc.Encoder{}
	}
	*h.enc = e
	return h, nil
}

func (h *CidBaseHandler) Encoder() cidenc.Encoder {
	return *h.enc
}

func (h *CidBaseHandler) EncoderFromPath(p string) cidenc.Encoder {
	if h.base == "" {
		enc, _ := cidenc.FromPath(*h.enc, p)
		return enc
	} else {
		return *h.enc
	}
}

func (h *CidBaseHandler) EncoderWithOverride() cidenc.Interface {
	if h.base == "" {
		enc := cidenc.NewOverride(*h.enc)
		enc.Add(h.args...)
		return enc
	} else {
		return *h.enc
	}
}
