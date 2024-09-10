package exec

import (
	"errors"
	"github.com/chainreactors/malice-network/client/command/common"
	"github.com/chainreactors/malice-network/client/core/intermediate/builtin"
	"github.com/chainreactors/malice-network/client/repl"
	"github.com/chainreactors/malice-network/helper/consts"
	"github.com/chainreactors/malice-network/helper/helper"
	"github.com/chainreactors/malice-network/proto/client/clientpb"
	"github.com/chainreactors/malice-network/proto/implant/implantpb"
	"github.com/chainreactors/malice-network/proto/services/clientrpc"
	"github.com/spf13/cobra"
	"google.golang.org/protobuf/proto"
)

// ExecuteExeCmd - Execute PE on sacrifice process
func ExecuteExeCmd(cmd *cobra.Command, con *repl.Console) {
	path, args, output, timeout, arch, process := common.ParseFullBinaryParams(cmd)
	sac, _ := common.ParseSacrifice(cmd)
	task, err := ExecExe(con.Rpc, con.GetInteractive(), path, args, output, timeout, arch, process, sac)
	if err != nil {
		con.Log.Errorf("Execute EXE error: %v", err)
		return
	}
	session := con.GetInteractive()
	con.AddCallback(task, func(msg proto.Message) {
		resp, _ := builtin.ParseAssembly(msg.(*implantpb.Spite))
		session.Log.Console(resp)
	})
}

func ExecExe(rpc clientrpc.MaliceRPCClient, sess *repl.Session, pePath string, args []string, output bool, timeout int, arch string, process string, sac *implantpb.SacrificeProcess) (*clientpb.Task, error) {
	if arch == "" {
		arch = sess.Os.Arch
	}
	binary, err := common.NewBinary(consts.ModuleExecuteExe, pePath, args, output, timeout, arch, process, sac)
	if err != nil {
		return nil, err
	}
	if helper.CheckPEType(binary.Bin) != consts.EXEFile {
		return nil, errors.New("the file is not a EXE file")
	}
	task, err := rpc.ExecuteEXE(sess.Context(), binary)
	if err != nil {
		return nil, err
	}
	return task, nil
}

// InlineExeCmd - Execute PE in current process
func InlineExeCmd(cmd *cobra.Command, con *repl.Console) {
	session := con.GetInteractive()
	path, args, output, timeout, arch, process := common.ParseFullBinaryParams(cmd)
	task, err := InlineExe(con.Rpc, session, path, args, output, timeout, arch, process)
	if err != nil {
		con.Log.Errorf("Execute ESE error: %v", err)
		return
	}
	con.AddCallback(task, func(msg proto.Message) {
		resp, _ := builtin.ParseAssembly(msg.(*implantpb.Spite))
		session.Log.Console(resp)
	})
}

func InlineExe(rpc clientrpc.MaliceRPCClient, sess *repl.Session, path string, args []string,
	output bool, timeout int, arch string, process string) (*clientpb.Task, error) {
	if arch == "" {
		arch = sess.Os.Arch
	}
	binary, err := common.NewBinary(consts.ModuleAliasInlineExe, path, args, output, timeout, arch, process, nil)
	if err != nil {
		return nil, err
	}
	if helper.CheckPEType(binary.Bin) != consts.EXEFile {
		return nil, errors.New("the file is not a PE file")

	}
	task, err := rpc.ExecuteEXE(sess.Context(), binary)
	if err != nil {
		return nil, err
	}
	return task, nil
}
