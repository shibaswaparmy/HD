package cli

import (
	"github.com/cosmos/cosmos-sdk/codec"
	"github.com/spf13/cobra"
	"github.com/cosmos/cosmos-sdk/client/context"
	"github.com/maticnetwork/heimdall/checkpoint"
	"github.com/maticnetwork/heimdall/helper"
	"time"
	"github.com/spf13/viper"
	"github.com/ethereum/go-ethereum/common"
	"strconv"
	"fmt"
)

// send checkpoint transaction
func GetSendCheckpointTx(cdc *codec.Codec) *cobra.Command  {
	cmd:=&cobra.Command{
		Use:   "send-checkpoint",
		Short: "send checkpoint to tendermint and ethereum chain ",
		RunE: func(cmd *cobra.Command, args []string) error {
			cliCtx := context.NewCLIContext().WithCodec(cdc)

			ProposerStr := viper.GetString(FlagProposerAddress)
			if viper.GetString(FlagProposerAddress) ==""{
				 return fmt.Errorf("Proposer address cannot be empty")
			}
			StartBlockStr := viper.GetString(FlagStartBlock)
			if StartBlockStr == "" {
				return fmt.Errorf("Start block cannot be empty")
			}

			EndBlockStr := viper.GetString(FlagEndBlock)
			if EndBlockStr == ""{
				return fmt.Errorf("End block cannot be empty")
			}

			RootHashStr := viper.GetString(FlagRootHash)
			if RootHashStr == ""{
				return fmt.Errorf("Root Hash cannot be empty")
			}

			Proposer:= common.HexToAddress(ProposerStr)

			StartBlock,err := strconv.ParseUint(StartBlockStr,10,64)
			if err!=nil{
				return err
			}

			EndBlock,err:= strconv.ParseUint(EndBlockStr,10,64)
			if err!=nil{
				return err
			}

			RootHash:=common.HexToHash(RootHashStr)

			msg := checkpoint.NewMsgCheckpointBlock(
				Proposer,
				StartBlock,
				EndBlock,
				RootHash,
				uint64(time.Now().Unix()),
			)

			return helper.CreateAndSendTx(msg,cliCtx)
		},
	}
	cmd.Flags().StringP(FlagProposerAddress,"p",helper.GetPubKey().Address().String(),"--proposer=0xAddress")
	cmd.Flags().StringP(FlagStartBlock,"s","start nlocl","proposer address")
	cmd.Flags().StringP(FlagEndBlock,"e","start nlocl","proposer address")
	cmd.Flags().StringP(FlagRootHash,"r","start nlocl","proposer address")
	cmd.MarkFlagRequired(FlagStartBlock)
	cmd.MarkFlagRequired(FlagEndBlock)
	cmd.MarkFlagRequired(FlagRootHash)

	return cmd
}


// send checkpoint ack transaction
func GetCheckpointACKTx(cdc *codec.Codec) *cobra.Command  {
	cmd:=&cobra.Command{
		Use:   "send-ack",
		Short: "send acknowledgement for checkpoint in buffer",
		RunE: func(cmd *cobra.Command, args []string) error {
			cliCtx := context.NewCLIContext().WithCodec(cdc)

			HeaderBlockStr:= viper.GetString(FlagHeaderNumber)

			HeaderBlock,err := strconv.ParseUint(HeaderBlockStr,10,64)
			if err!=nil{
				return err
			}

			msg := checkpoint.NewMsgCheckpointAck(HeaderBlock)

			return helper.CreateAndSendTx(msg,cliCtx)
		},
	}

	cmd.MarkFlagRequired(FlagHeaderNumber)
	return cmd
}

// send no-ack transaction
func GetCheckpointNoACKTx(cdc *codec.Codec) *cobra.Command{
	cmd:=&cobra.Command{
		Use:   "send-NoACK",
		Short: "send no-acknowledgement for last proposer",
		RunE: func(cmd *cobra.Command, args []string) error {
			cliCtx := context.NewCLIContext().WithCodec(cdc)

			msg := checkpoint.NewMsgCheckpointNoAck(uint64(time.Now().Unix()))

			return helper.CreateAndSendTx(msg,cliCtx)
		},
	}
	return cmd
}

