package types

import (
	"database/sql"
	"database/sql/driver"
	"encoding/json"
	"fmt"
	"html/template"
	"math/big"
	"strconv"
	"time"

	itypes "github.com/gobitfly/eth-rewards/types"
	"github.com/lib/pq"
	"github.com/pkg/errors"
	"github.com/shopspring/decimal"
	"github.com/theQRL/go-zond/accounts/abi"
	"github.com/theQRL/go-zond/common"
	geth_types "github.com/theQRL/go-zond/core/types"
)

// PageData is a struct to hold web page data
type PageData struct {
	Active             string
	AdConfigurations   []*AdConfig
	Meta               *Meta
	ShowSyncingMessage bool
	// User                  *User
	Data                  interface{}
	Version               string
	Year                  int
	ChainSlotsPerEpoch    uint64
	ChainSecondsPerSlot   uint64
	ChainGenesisTimestamp uint64
	CurrentEpoch          uint64
	LatestFinalizedEpoch  uint64
	CurrentSlot           uint64
	FinalizationDelay     uint64
	Mainnet               bool
	DepositContract       string
	Rates                 *Rates
	InfoBanner            *template.HTML
	// IsUserClientUpdated   func(uint64) bool
	ChainConfig         ClChainConfig
	Lang                string
	Debug               bool
	DebugTemplates      []string
	DebugSession        map[string]interface{}
	GasNow              *GasNowPageData
	GlobalNotification  template.HTML
	AvailableCurrencies []string
	MainMenuItems       []MainMenuItem
	TermsOfServiceUrl   string
	PrivacyPolicyUrl    string
}

type MainMenuItem struct {
	Label        string
	Path         string
	IsActive     bool
	HasBigGroups bool // if HasBigGroups is set to true then the NavigationGroups will be ordered horizontally and their Label will be shown
	Groups       []NavigationGroup
}

type NavigationGroup struct {
	Label string // only used for "BigGroups"
	Links []NavigationLink
}

type NavigationLink struct {
	Label         string
	Path          string
	CustomIcon    string
	Icon          string
	IsHidden      bool
	IsHighlighted bool
}

type Rates struct {
	TickerCurrency                    string                `json:"tickerCurrency"`
	TickerCurrencySymbol              string                `json:"tickerCurrencySymbol"`
	SelectedCurrency                  string                `json:"selectedCurrency"`
	SelectedCurrencySymbol            string                `json:"selectedCurrencySymbol"`
	MainCurrency                      string                `json:"mainCurrency"`
	MainCurrencySymbol                string                `json:"mainCurrencySymbol"`
	MainCurrencyPrice                 float64               `json:"mainCurrencyPrice"`
	MainCurrencyPriceFormatted        template.HTML         `json:"mainCurrencyPriceFormatted"`
	MainCurrencyPriceKFormatted       template.HTML         `json:"mainCurrencyKFormatted"`
	MainCurrencyTickerPrice           float64               `json:"mainCurrencyTickerPrice"`
	MainCurrencyTickerPriceFormatted  template.HTML         `json:"mainCurrencyTickerPriceFormatted"`
	MainCurrencyTickerPriceKFormatted template.HTML         `json:"mainCurrencyTickerPriceKFormatted"`
	ElCurrency                        string                `json:"elCurrency"`
	ElCurrencySymbol                  string                `json:"elCurrencySymbol"`
	ElCurrencyPrice                   float64               `json:"elCurrencyPrice"`
	ElCurrencyPriceFormatted          template.HTML         `json:"elCurrencyPriceFormatted"`
	ElCurrencyPriceKFormatted         template.HTML         `json:"elCurrencyKFormatted"`
	ClCurrency                        string                `json:"clCurrency"`
	ClCurrencySymbol                  string                `json:"clCurrencySymbol"`
	ClCurrencyPrice                   float64               `json:"clCurrencyPrice"`
	ClCurrencyPriceFormatted          template.HTML         `json:"clCurrencyPriceFormatted"`
	ClCurrencyPriceKFormatted         template.HTML         `json:"clCurrencyKFormatted"`
	MainCurrencyPrices                map[string]RatesPrice `json:"mainCurrencyTickerPrices"`
}

type RatesPrice struct {
	Symbol     string        `json:"symbol"`
	RoundPrice uint64        `json:"roundPrice"`
	TruncPrice template.HTML `json:"truncPrice"`
}

// Meta is a struct to hold metadata about the page
type Meta struct {
	Title       string
	Description string
	Path        string
	Tlabel1     string
	Tdata1      string
	Tlabel2     string
	Tdata2      string
	GATag       string
	NoTrack     bool
	Templates   string
}

// LatestState is a struct to hold data for the banner
type LatestState struct {
	LastProposedSlot      uint64 `json:"lastProposedSlot"`
	CurrentSlot           uint64 `json:"currentSlot"`
	CurrentEpoch          uint64 `json:"currentEpoch"`
	CurrentFinalizedEpoch uint64 `json:"currentFinalizedEpoch"`
	FinalityDelay         uint64 `json:"finalityDelay"`
	IsSyncing             bool   `json:"syncing"`
	Rates                 *Rates `json:"rates"`
}

type Stats struct {
	TopDepositors                  *[]StatsTopDepositors
	InvalidDepositCount            *uint64 `db:"count"`
	UniqueValidatorCount           *uint64 `db:"count"`
	TotalValidatorCount            *uint64 `db:"count"`
	ActiveValidatorCount           *uint64 `db:"count"`
	PendingValidatorCount          *uint64 `db:"count"`
	ValidatorChurnLimit            *uint64
	ValidatorActivationChurnLimit  *uint64
	LatestValidatorWithdrawalIndex *uint64 `db:"index"`
	WithdrawableValidatorCount     *uint64 `db:"count"`
	// WithdrawableAmount             *uint64 `db:"amount"`
	PendingDilithiumChangeValidatorCount *uint64 `db:"count"`
	NonWithdrawableCount                 *uint64 `db:"count"`
	TotalAmountWithdrawn                 *uint64 `db:"amount"`
	WithdrawalCount                      *uint64 `db:"count"`
	TotalAmountDeposited                 *uint64 `db:"amount"`
	DilithiumChangeCount                 *uint64 `db:"count"`
}

type StatsTopDepositors struct {
	Address      string `db:"from_address"`
	DepositCount uint64 `db:"count"`
}

// IndexPageData is a struct to hold info for the main web page
type IndexPageData struct {
	NetworkName               string `json:"networkName"`
	DepositContract           string `json:"depositContract"`
	ShowSyncingMessage        bool
	CurrentEpoch              uint64                 `json:"current_epoch"`
	CurrentFinalizedEpoch     uint64                 `json:"current_finalized_epoch"`
	CurrentSlot               uint64                 `json:"current_slot"`
	ScheduledCount            uint8                  `json:"scheduled_count"`
	FinalityDelay             uint64                 `json:"finality_delay"`
	ActiveValidators          uint64                 `json:"active_validators"`
	EnteringValidators        uint64                 `json:"entering_validators"`
	ExitingValidators         uint64                 `json:"exiting_validators"`
	StakedEther               string                 `json:"staked_ether"`
	AverageBalance            string                 `json:"average_balance"`
	DepositedTotal            float64                `json:"deposit_total"`
	DepositThreshold          float64                `json:"deposit_threshold"`
	ValidatorsRemaining       float64                `json:"validators_remaining"`
	NetworkStartTs            int64                  `json:"network_start_ts"`
	MinGenesisTime            int64                  `json:"minGenesisTime"`
	Blocks                    []*IndexPageDataBlocks `json:"blocks"`
	Epochs                    []*IndexPageDataEpochs `json:"epochs"`
	StakedEtherChartData      [][]float64            `json:"staked_ether_chart_data"`
	ActiveValidatorsChartData [][]float64            `json:"active_validators_chart_data"`
	Title                     template.HTML          `json:"title"`
	Subtitle                  template.HTML          `json:"subtitle"`
	Genesis                   bool                   `json:"genesis"`
	GenesisPeriod             bool                   `json:"genesis_period"`
	Mainnet                   bool                   `json:"mainnet"`
	DepositChart              *ChartsPageDataChart
	DepositDistribution       *ChartsPageDataChart
	Countdown                 interface{}
	SlotVizData               *SlotVizPageData `json:"slotVizData"`
	ClCurrency                string
	ElCurrency                string
	ValidatorsPerEpoch        uint64
	ValidatorsPerDay          uint64
	NewDepositProcessAfter    string
}

type SlotVizPageData struct {
	Epochs   []*SlotVizEpochs
	Selector string
	Config   ExplorerConfigurationKeyMap
}

type IndexPageDataEpochs struct {
	Epoch                            uint64        `json:"epoch"`
	Ts                               time.Time     `json:"ts"`
	Finalized                        bool          `json:"finalized"`
	FinalizedFormatted               template.HTML `json:"finalized_formatted"`
	EligibleEther                    uint64        `json:"eligibleether"`
	EligibleEtherFormatted           template.HTML `json:"eligibleether_formatted"`
	GlobalParticipationRate          float64       `json:"globalparticipationrate"`
	GlobalParticipationRateFormatted template.HTML `json:"globalparticipationrate_formatted"`
	VotedEther                       uint64        `json:"votedether"`
	VotedEtherFormatted              template.HTML `json:"votedether_formatted"`
}

// IndexPageDataBlocks is a struct to hold detail data for the main web page
type IndexPageDataBlocks struct {
	Epoch                uint64        `json:"epoch"`
	Slot                 uint64        `json:"slot"`
	Ts                   time.Time     `json:"ts"`
	Proposer             uint64        `db:"proposer" json:"proposer"`
	ProposerFormatted    template.HTML `json:"proposer_formatted"`
	BlockRoot            []byte        `db:"blockroot" json:"block_root"`
	BlockRootFormatted   string        `json:"block_root_formatted"`
	ParentRoot           []byte        `db:"parentroot" json:"parent_root"`
	Attestations         uint64        `db:"attestationscount" json:"attestations"`
	Deposits             uint64        `db:"depositscount" json:"deposits"`
	Withdrawals          uint64        `db:"withdrawalcount" json:"withdrawals"`
	Exits                uint64        `db:"voluntaryexitscount" json:"exits"`
	Proposerslashings    uint64        `db:"proposerslashingscount" json:"proposerslashings"`
	Attesterslashings    uint64        `db:"attesterslashingscount" json:"attesterslashings"`
	SyncAggParticipation float64       `db:"syncaggregate_participation" json:"sync_aggregate_participation"`
	Status               uint64        `db:"status" json:"status"`
	StatusFormatted      template.HTML `json:"status_formatted"`
	Votes                uint64        `db:"votes" json:"votes"`
	Graffiti             []byte        `db:"graffiti"`
	ProposerName         string        `db:"name"`
	ExecutionBlockNumber int           `db:"exec_block_number" json:"exec_block_number"`
}

// IndexPageEpochHistory is a struct to hold the epoch history for the main web page
type IndexPageEpochHistory struct {
	Epoch                   uint64 `db:"epoch"`
	ValidatorsCount         uint64 `db:"validatorscount"`
	EligibleEther           uint64 `db:"eligibleether"`
	Finalized               bool   `db:"finalized"`
	AverageValidatorBalance uint64 `db:"averagevalidatorbalance"`
}

// IndexPageDataBlocks is a struct to hold detail data for the main web page
type BlocksPageDataBlocks struct {
	TotalCount           uint64        `db:"total_count"`
	Epoch                uint64        `json:"epoch"`
	Slot                 uint64        `json:"slot"`
	Ts                   time.Time     `json:"ts"`
	Proposer             uint64        `db:"proposer" json:"proposer"`
	ProposerFormatted    template.HTML `json:"proposer_formatted"`
	BlockRoot            []byte        `db:"blockroot" json:"block_root"`
	BlockRootFormatted   string        `json:"block_root_formatted"`
	ParentRoot           []byte        `db:"parentroot" json:"parent_root"`
	Attestations         uint64        `db:"attestationscount" json:"attestations"`
	Deposits             uint64        `db:"depositscount" json:"deposits"`
	Withdrawals          uint64        `db:"withdrawalcount" json:"withdrawals"`
	Exits                uint64        `db:"voluntaryexitscount" json:"exits"`
	Proposerslashings    uint64        `db:"proposerslashingscount" json:"proposerslashings"`
	Attesterslashings    uint64        `db:"attesterslashingscount" json:"attesterslashings"`
	SyncAggParticipation float64       `db:"syncaggregate_participation" json:"sync_aggregate_participation"`
	Status               uint64        `db:"status" json:"status"`
	StatusFormatted      template.HTML `json:"status_formatted"`
	Votes                uint64        `db:"votes" json:"votes"`
	Graffiti             []byte        `db:"graffiti"`
	ProposerName         string        `db:"name"`
}

// ValidatorsPageData is a struct to hold data about the validators page
type ValidatorsPageData struct {
	TotalCount           uint64
	DepositedCount       uint64
	PendingCount         uint64
	ActiveCount          uint64
	ActiveOnlineCount    uint64
	ActiveOfflineCount   uint64
	SlashingCount        uint64
	SlashingOnlineCount  uint64
	SlashingOfflineCount uint64
	Slashed              uint64
	ExitingCount         uint64
	ExitingOnlineCount   uint64
	ExitingOfflineCount  uint64
	ExitedCount          uint64
	VoluntaryExitsCount  uint64
	UnknownCount         uint64
	Validators           []*ValidatorsData
}

// ValidatorsData is a struct to hold data about validators
type ValidatorsData struct {
	TotalCount                 uint64 `db:"total_count"`
	Epoch                      uint64 `db:"epoch"`
	PublicKey                  []byte `db:"pubkey"`
	ValidatorIndex             uint64 `db:"validatorindex"`
	WithdrawableEpoch          uint64 `db:"withdrawableepoch"`
	CurrentBalance             uint64 `db:"balance"`
	EffectiveBalance           uint64 `db:"effectivebalance"`
	Slashed                    bool   `db:"slashed"`
	ActivationEligibilityEpoch uint64 `db:"activationeligibilityepoch"`
	ActivationEpoch            uint64 `db:"activationepoch"`
	ExitEpoch                  uint64 `db:"exitepoch"`
	LastAttestationSlot        int64
	Name                       string  `db:"name"`
	State                      string  `db:"state"`
	MissedProposals            uint64  `db:"missedproposals"`
	ExecutedProposals          uint64  `db:"executedproposals"`
	MissedAttestations         uint64  `db:"missedattestations"`
	ExecutedAttestations       uint64  `db:"executedattestations"`
	Performance7d              int64   `db:"performance7d"`
	DepositAddress             *[]byte `db:"from_address"`
}

// ValidatorPageData is a struct to hold data for the validators page
type ValidatorPageData struct {
	Epoch                           uint64 `db:"epoch"`
	ValidatorIndex                  uint64 `db:"validatorindex"`
	PublicKey                       []byte `db:"pubkey"`
	WithdrawableEpoch               uint64 `db:"withdrawableepoch"`
	WithdrawCredentials             []byte `db:"withdrawalcredentials"`
	CurrentBalance                  uint64 `db:"balance"`
	BalanceActivation               uint64 `db:"balanceactivation"`
	EffectiveBalance                uint64 `db:"effectivebalance"`
	Slashed                         bool   `db:"slashed"`
	SlashedBy                       uint64
	SlashedAt                       uint64
	SlashedFor                      string
	ActivationEligibilityEpoch      uint64 `db:"activationeligibilityepoch"`
	ActivationEpoch                 uint64 `db:"activationepoch"`
	ExitEpoch                       uint64 `db:"exitepoch"`
	Index                           uint64 `db:"index"`
	LastAttestationSlot             uint64
	Name                            string         `db:"name"`
	Pool                            string         `db:"pool"`
	Tags                            pq.StringArray `db:"tags"`
	WithdrawableTs                  time.Time
	ActivationEligibilityTs         time.Time
	ActivationTs                    time.Time
	ExitTs                          time.Time
	Status                          string `db:"status"`
	AttestationsCount               uint64
	ExecutedAttestationsCount       uint64
	MissedAttestationsCount         uint64
	UnmissedAttestationsPercentage  float64 // missed/(executed+orphaned)
	DepositsCount                   uint64
	WithdrawalCount                 uint64
	SlashingsCount                  uint64
	PendingCount                    uint64
	SyncCount                       uint64 // amount of sync committees the validator was (and is) part of
	SlotsPerSyncCommittee           uint64
	FutureDutiesEpoch               uint64
	SlotsDoneInCurrentSyncCommittee uint64
	ScheduledSyncCountSlots         uint64
	ParticipatedSyncCountSlots      uint64
	MissedSyncCountSlots            uint64
	OrphanedSyncCountSlots          uint64
	UnmissedSyncPercentage          float64 // participated/(participated+missed)
	Income                          *ValidatorEarnings
	IncomeToday                     ClEl          `json:"incomeToday"`
	Income1d                        ClEl          `json:"income1d"`
	Income7d                        ClEl          `json:"income7d"`
	Income31d                       ClEl          `json:"income31d"`
	IncomeTotal                     ClEl          `json:"incomeTotal"`
	IncomeTotalFormatted            template.HTML `json:"incomeTotalFormatted"`
	Apr7d                           ClElFloat64   `json:"apr7d"`
	Apr31d                          ClElFloat64   `json:"apr31d"`
	Apr365d                         ClElFloat64   `json:"apr365d"`
	SyncLuck                        float64
	SyncEstimate                    *time.Time
	AvgSyncInterval                 *time.Duration
	Rank7d                          int64 `db:"rank7d"`
	RankCount                       int64 `db:"rank_count"`
	RankPercentage                  float64
	IncomeHistoryChartData          []*ChartDataPoint
	ExecutionIncomeHistoryData      []*ChartDataPoint
	Deposits                        *ValidatorDeposits
	Eth1DepositAddress              []byte
	FlashMessage                    string
	SubscriptionFlash               []interface{}
	// User                                     *User
	AttestationInclusionEffectiveness        float64
	CsrfField                                template.HTML
	NetworkStats                             *IndexPageData
	ChurnRate                                uint64
	QueuePosition                            uint64
	EstimatedActivationTs                    time.Time
	EstimatedActivationEpoch                 uint64
	InclusionDelay                           int64
	CurrentAttestationStreak                 uint64
	LongestAttestationStreak                 uint64
	ShowMultipleWithdrawalCredentialsWarning bool
	DilithiumChange                          *DilithiumChange
	IsWithdrawableAddress                    bool
	EstimatedNextWithdrawal                  template.HTML
	// AddValidatorWatchlistModal               *AddValidatorWatchlistModal
	NextWithdrawalRow [][]interface{}
	ValidatorProposalData
}

type ValidatorStatsTablePageData struct {
	ValidatorIndex uint64
	Rows           []*ValidatorStatsTableRow
	Currency       string
}

type ValidatorStatsTableRow struct {
	ValidatorIndex         uint64
	Day                    int64         `db:"day"`
	StartBalance           sql.NullInt64 `db:"start_balance"`
	EndBalance             sql.NullInt64 `db:"end_balance"`
	Income                 int64         `db:"cl_rewards_gwei"`
	IncomeExchangeRate     float64       `db:"-"`
	IncomeExchangeCurrency string        `db:"-"`
	IncomeExchanged        float64       `db:"-"`
	MinBalance             sql.NullInt64 `db:"min_balance"`
	MaxBalance             sql.NullInt64 `db:"max_balance"`
	StartEffectiveBalance  sql.NullInt64 `db:"start_effective_balance"`
	EndEffectiveBalance    sql.NullInt64 `db:"end_effective_balance"`
	MinEffectiveBalance    sql.NullInt64 `db:"min_effective_balance"`
	MaxEffectiveBalance    sql.NullInt64 `db:"max_effective_balance"`
	MissedAttestations     sql.NullInt64 `db:"missed_attestations"`
	ProposedBlocks         sql.NullInt64 `db:"proposed_blocks"`
	MissedBlocks           sql.NullInt64 `db:"missed_blocks"`
	OrphanedBlocks         sql.NullInt64 `db:"orphaned_blocks"`
	AttesterSlashings      sql.NullInt64 `db:"attester_slashings"`
	ProposerSlashings      sql.NullInt64 `db:"proposer_slashings"`
	Deposits               sql.NullInt64 `db:"deposits"`
	DepositsAmount         sql.NullInt64 `db:"deposits_amount"`
	ParticipatedSync       sql.NullInt64 `db:"participated_sync"`
	MissedSync             sql.NullInt64 `db:"missed_sync"`
	OrphanedSync           sql.NullInt64 `db:"orphaned_sync"`
}

type ChartDataPoint struct {
	X     float64 `json:"x"`
	Y     float64 `json:"y"`
	Color string  `json:"color"`
}

// ValidatorRank is a struct for validator rank data
type ValidatorRank struct {
	Rank int64 `db:"rank" json:"rank"`
}

// DailyProposalCount is a struct for the daily proposal count data
type DailyProposalCount struct {
	Day      int64
	Proposed uint
	Missed   uint
	Orphaned uint
}

// ValidatorBalanceHistory is a struct for the validator balance history data
type ValidatorBalanceHistory struct {
	Day              uint64 `db:"day"`
	Balance          uint64 `db:"balance"`
	EffectiveBalance uint64 `db:"effectivebalance"`
}

// ValidatorBalanceHistory is a struct for the validator income history data
type ValidatorIncomeHistory struct {
	Day              int64         `db:"day"` // day can be -1 which is pre-genesis
	ClRewards        int64         `db:"cl_rewards_gwei"`
	EndBalance       sql.NullInt64 `db:"end_balance"`
	StartBalance     sql.NullInt64 `db:"start_balance"`
	DepositAmount    sql.NullInt64 `db:"deposits_amount"`
	WithdrawalAmount sql.NullInt64 `db:"withdrawals_amount"`
}

type ValidatorBalanceHistoryChartData struct {
	Epoch   uint64
	Balance uint64
}

// ValidatorBalance is a struct for the validator balance data
type ValidatorBalance struct {
	Epoch            uint64 `db:"epoch"`
	Balance          uint64 `db:"balance"`
	EffectiveBalance uint64 `db:"effectivebalance"`
	Index            uint64 `db:"validatorindex"`
	PublicKey        []byte `db:"pubkey"`
}

// ValidatorPerformance is a struct for the validator performance data
type ValidatorPerformance struct {
	Rank            uint64 `db:"rank"`
	Index           uint64 `db:"validatorindex"`
	PublicKey       []byte `db:"pubkey"`
	Name            string `db:"name"`
	Balance         uint64 `db:"balance"`
	Performance1d   int64  `db:"performance1d"`
	Performance7d   int64  `db:"performance7d"`
	Performance31d  int64  `db:"performance31d"`
	Performance365d int64  `db:"performance365d"`
	Rank7d          int64  `db:"rank7d"`
	TotalCount      uint64 `db:"total_count"`
}

// ValidatorAttestation is a struct for the validators attestations data
type ValidatorAttestation struct {
	Index          uint64
	Epoch          uint64 `db:"epoch"`
	AttesterSlot   uint64 `db:"attesterslot"`
	CommitteeIndex uint64 `db:"committeeindex"`
	Status         uint64 `db:"status"`
	InclusionSlot  uint64 `db:"inclusionslot"`
	Delay          int64  `db:"delay"`
	// EarliestInclusionSlot uint64 `db:"earliestinclusionslot"`
}

// ValidatorSyncParticipation hold information about sync-participation of a validator
type ValidatorSyncParticipation struct {
	Period uint64 `db:"period"`
	Slot   uint64 `db:"slot"`
	Status uint64 `db:"status"`
}

// type AvgInclusionDistance struct {
// 	InclusionSlot         uint64 `db:"inclusionslot"`
// 	EarliestInclusionSlot uint64 `db:"earliestinclusionslot"`
// }

// VisPageData is a struct to hold the visualizations page data
type VisPageData struct {
	ChartData  []*VisChartData
	StartEpoch uint64
	EndEpoch   uint64
}

// VisChartData is a struct to hold the visualizations chart data
type VisChartData struct {
	Slot       uint64 `db:"slot" json:"-"`
	BlockRoot  []byte `db:"blockroot" json:"-"`
	ParentRoot []byte `db:"parentroot" json:"-"`

	Proposer uint64 `db:"proposer" json:"proposer"`

	Number    uint64   `json:"number"`
	Timestamp uint64   `json:"timestamp"`
	Hash      string   `json:"hash"`
	Parents   []string `json:"parents"`
}

// VisVotesPageData is a struct for the visualization votes page data
type VisVotesPageData struct {
	ChartData []*VotesVisChartData
}

// VotesVisChartData is a struct for the visualization chart data
type VotesVisChartData struct {
	Slot       uint64        `db:"slot" json:"slot"`
	BlockRoot  string        `db:"blockroot" json:"blockRoot"`
	ParentRoot string        `db:"parentroot" json:"parentRoot"`
	Validators pq.Int64Array `db:"validators" json:"validators"`
}

// BlockPageData is a struct block data used in the block/slot page
type BlockPageData struct {
	Epoch                  uint64  `db:"epoch"`
	EpochFinalized         bool    `db:"epoch_finalized"`
	PrevEpochFinalized     bool    `db:"prev_epoch_finalized"`
	EpochParticipationRate float64 `db:"epoch_participation_rate"`
	Ts                     time.Time
	NextSlot               uint64
	PreviousSlot           uint64
	Proposer               uint64  `db:"proposer"`
	BlockRoot              []byte  `db:"blockroot"`
	ParentRoot             []byte  `db:"parentroot"`
	StateRoot              []byte  `db:"stateroot"`
	Signature              []byte  `db:"signature"`
	RandaoReveal           []byte  `db:"randaoreveal"`
	Graffiti               []byte  `db:"graffiti"`
	ProposerName           string  `db:"name"`
	Eth1dataDepositroot    []byte  `db:"eth1data_depositroot"`
	Eth1dataDepositcount   uint64  `db:"eth1data_depositcount"`
	Eth1dataBlockhash      []byte  `db:"eth1data_blockhash"`
	SyncAggregateBits      []byte  `db:"syncaggregate_bits"`
	SyncAggregateSignature []byte  `db:"syncaggregate_signature"`
	SyncAggParticipation   float64 `db:"syncaggregate_participation"`
	ProposerSlashingsCount uint64  `db:"proposerslashingscount"`
	AttesterSlashingsCount uint64  `db:"attesterslashingscount"`
	AttestationsCount      uint64  `db:"attestationscount"`
	DepositsCount          uint64  `db:"depositscount"`
	WithdrawalCount        uint64  `db:"withdrawalcount"`
	DilithiumChangeCount   uint64  `db:"dilithium_change_count"`
	VoluntaryExitscount    uint64  `db:"voluntaryexitscount"`
	SlashingsCount         uint64
	VotesCount             uint64
	VotingValidatorsCount  uint64
	Mainnet                bool

	ExecParentHash        []byte        `db:"exec_parent_hash"`
	ExecFeeRecipient      []byte        `db:"exec_fee_recipient"`
	ExecStateRoot         []byte        `db:"exec_state_root"`
	ExecReceiptsRoot      []byte        `db:"exec_receipts_root"`
	ExecLogsBloom         []byte        `db:"exec_logs_bloom"`
	ExecRandom            []byte        `db:"exec_random"`
	ExecGasLimit          sql.NullInt64 `db:"exec_gas_limit"`
	ExecGasUsed           sql.NullInt64 `db:"exec_gas_used"`
	ExecTimestamp         sql.NullInt64 `db:"exec_timestamp"`
	ExecTime              time.Time
	ExecExtraData         []byte        `db:"exec_extra_data"`
	ExecBaseFeePerGas     sql.NullInt64 `db:"exec_base_fee_per_gas"`
	ExecBlockHash         []byte        `db:"exec_block_hash"`
	ExecTransactionsCount uint64        `db:"exec_transactions_count"`

	Transactions []*BlockPageTransaction

	Withdrawals []*Withdrawals

	ExecutionData *Eth1BlockPageData

	Attestations      []*BlockPageAttestation // Attestations included in this block
	VoluntaryExits    []*BlockPageVoluntaryExits
	Votes             []*BlockVote // Attestations that voted for that block
	AttesterSlashings []*BlockPageAttesterSlashing
	ProposerSlashings []*BlockPageProposerSlashing
	SyncCommittee     []uint64 // TODO: Setting it to contain the validator index

	Tags TagMetadataSlice `db:"tags"`
	ValidatorProposalInfo
}

func (u *BlockPageData) MarshalJSON() ([]byte, error) {
	type Alias BlockPageData
	return json.Marshal(&struct {
		BlockRoot string
		Ts        int64
		*Alias
	}{
		BlockRoot: fmt.Sprintf("%x", u.BlockRoot),
		Ts:        u.Ts.Unix(),
		Alias:     (*Alias)(u),
	})
}

// BlockVote stores a vote for a given block
type BlockVote struct {
	Validator      uint64 `db:"validator"`
	IncludedIn     uint64 `db:"included_in"`
	CommitteeIndex uint64 `db:"committee_index"`
}

// BlockPageMinMaxSlot is a struct to hold min/max slot data
type BlockPageMinMaxSlot struct {
	MinSlot uint64
	MaxSlot uint64
}

// BlockPageTransaction is a struct to hold execution transactions on the block page
type BlockPageTransaction struct {
	BlockSlot    uint64 `db:"block_slot"`
	BlockIndex   uint64 `db:"block_index"`
	TxHash       []byte `db:"txhash"`
	AccountNonce uint64 `db:"nonce"`
	// big endian
	Price       []byte `db:"gas_price"`
	PricePretty string
	GasLimit    uint64 `db:"gas_limit"`
	Sender      []byte `db:"sender"`
	Recipient   []byte `db:"recipient"`
	// big endian
	Amount       []byte `db:"amount"`
	AmountPretty string
	Payload      []byte `db:"payload"`

	// TODO: transaction type

	MaxPriorityFeePerGas uint64 `db:"max_priority_fee_per_gas"`
	MaxFeePerGas         uint64 `db:"max_fee_per_gas"`
}

// BlockPageAttestation is a struct to hold attestations on the block page
type BlockPageAttestation struct {
	BlockSlot       uint64        `db:"block_slot"`
	BlockIndex      uint64        `db:"block_index"`
	AggregationBits []byte        `db:"aggregationbits"`
	Validators      pq.Int64Array `db:"validators"`
	Signature       []byte        `db:"signature"`
	Slot            uint64        `db:"slot"`
	CommitteeIndex  uint64        `db:"committeeindex"`
	BeaconBlockRoot []byte        `db:"beaconblockroot"`
	SourceEpoch     uint64        `db:"source_epoch"`
	SourceRoot      []byte        `db:"source_root"`
	TargetEpoch     uint64        `db:"target_epoch"`
	TargetRoot      []byte        `db:"target_root"`
}

// BlockPageDeposit is a struct to hold data for deposits on the block page
type BlockPageDeposit struct {
	PublicKey             []byte `db:"publickey"`
	WithdrawalCredentials []byte `db:"withdrawalcredentials"`
	Amount                uint64 `db:"amount"`
	Signature             []byte `db:"signature"`
}

// BlockPageVoluntaryExits is a struct to hold data for voluntary exits on the block page
type BlockPageVoluntaryExits struct {
	ValidatorIndex uint64 `db:"validatorindex"`
	Signature      []byte `db:"signature"`
}

// BlockPageAttesterSlashing is a struct to hold data for attester slashings on the block page
type BlockPageAttesterSlashing struct {
	BlockSlot                   uint64        `db:"block_slot"`
	BlockIndex                  uint64        `db:"block_index"`
	Attestation1Indices         pq.Int64Array `db:"attestation1_indices"`
	Attestation1Signature       []byte        `db:"attestation1_signature"`
	Attestation1Slot            uint64        `db:"attestation1_slot"`
	Attestation1Index           uint64        `db:"attestation1_index"`
	Attestation1BeaconBlockRoot []byte        `db:"attestation1_beaconblockroot"`
	Attestation1SourceEpoch     uint64        `db:"attestation1_source_epoch"`
	Attestation1SourceRoot      []byte        `db:"attestation1_source_root"`
	Attestation1TargetEpoch     uint64        `db:"attestation1_target_epoch"`
	Attestation1TargetRoot      []byte        `db:"attestation1_target_root"`
	Attestation2Indices         pq.Int64Array `db:"attestation2_indices"`
	Attestation2Signature       []byte        `db:"attestation2_signature"`
	Attestation2Slot            uint64        `db:"attestation2_slot"`
	Attestation2Index           uint64        `db:"attestation2_index"`
	Attestation2BeaconBlockRoot []byte        `db:"attestation2_beaconblockroot"`
	Attestation2SourceEpoch     uint64        `db:"attestation2_source_epoch"`
	Attestation2SourceRoot      []byte        `db:"attestation2_source_root"`
	Attestation2TargetEpoch     uint64        `db:"attestation2_target_epoch"`
	Attestation2TargetRoot      []byte        `db:"attestation2_target_root"`
	SlashedValidators           []int64
}

// BlockPageProposerSlashing is a struct to hold data for proposer slashings on the block page
type BlockPageProposerSlashing struct {
	BlockSlot         uint64 `db:"block_slot"`
	BlockIndex        uint64 `db:"block_index"`
	BlockRoot         []byte `db:"block_root" json:"block_root"`
	ProposerIndex     uint64 `db:"proposerindex"`
	Header1Slot       uint64 `db:"header1_slot"`
	Header1ParentRoot []byte `db:"header1_parentroot"`
	Header1StateRoot  []byte `db:"header1_stateroot"`
	Header1BodyRoot   []byte `db:"header1_bodyroot"`
	Header1Signature  []byte `db:"header1_signature"`
	Header2Slot       uint64 `db:"header2_slot"`
	Header2ParentRoot []byte `db:"header2_parentroot"`
	Header2StateRoot  []byte `db:"header2_stateroot"`
	Header2BodyRoot   []byte `db:"header2_bodyroot"`
	Header2Signature  []byte `db:"header2_signature"`
}

// DataTableResponse is a struct to hold data for data table responses
type DataTableResponse struct {
	Draw            uint64          `json:"draw"`
	RecordsTotal    uint64          `json:"recordsTotal"`
	RecordsFiltered uint64          `json:"recordsFiltered"`
	Data            [][]interface{} `json:"data"`
	PageLength      uint64          `json:"pageLength"`
	DisplayStart    uint64          `json:"displayStart"`
	PagingToken     string          `json:"pagingToken"`
}

// EpochsPageData is a struct to hold epoch data for the epochs page
type EpochsPageData struct {
	Epoch                   uint64  `db:"epoch"`
	BlocksCount             uint64  `db:"blockscount"`
	ProposerSlashingsCount  uint64  `db:"proposerslashingscount"`
	AttesterSlashingsCount  uint64  `db:"attesterslashingscount"`
	AttestationsCount       uint64  `db:"attestationscount"`
	DepositsCount           uint64  `db:"depositscount"`
	WithdrawalCount         uint64  `db:"withdrawalcount"`
	VoluntaryExitsCount     uint64  `db:"voluntaryexitscount"`
	ValidatorsCount         uint64  `db:"validatorscount"`
	AverageValidatorBalance uint64  `db:"averagevalidatorbalance"`
	Finalized               bool    `db:"finalized"`
	EligibleEther           uint64  `db:"eligibleether"`
	GlobalParticipationRate float64 `db:"globalparticipationrate"`
	VotedEther              uint64  `db:"votedether"`
}

// EpochPageData is a struct to hold detailed epoch data for the epoch page
type EpochPageData struct {
	Epoch                   uint64        `db:"epoch"`
	BlocksCount             uint64        `db:"blockscount"`
	ProposerSlashingsCount  uint64        `db:"proposerslashingscount"`
	AttesterSlashingsCount  uint64        `db:"attesterslashingscount"`
	AttestationsCount       uint64        `db:"attestationscount"`
	DepositsCount           uint64        `db:"depositscount"`
	WithdrawalCount         uint64        `db:"withdrawalcount"`
	DepositTotal            uint64        `db:"deposittotal"`
	WithdrawalTotal         template.HTML `db:"withdrawaltotal"`
	VoluntaryExitsCount     uint64        `db:"voluntaryexitscount"`
	ValidatorsCount         uint64        `db:"validatorscount"`
	AverageValidatorBalance uint64        `db:"averagevalidatorbalance"`
	Finalized               bool          `db:"finalized"`
	EligibleEther           uint64        `db:"eligibleether"`
	GlobalParticipationRate float64       `db:"globalparticipationrate"`
	VotedEther              uint64        `db:"votedether"`

	Blocks []*IndexPageDataBlocks

	SyncParticipationRate float64
	Ts                    time.Time
	NextEpoch             uint64
	PreviousEpoch         uint64
	ProposedCount         uint64
	MissedCount           uint64
	ScheduledCount        uint64
	OrphanedCount         uint64
}

// EpochPageMinMaxSlot is a struct for the min/max epoch data
type EpochPageMinMaxSlot struct {
	MinEpoch uint64
	MaxEpoch uint64
}

// SearchAheadEpochsResult is a struct to hold the search ahead epochs results
type SearchAheadEpochsResult []struct {
	Epoch string `db:"epoch" json:"epoch,omitempty"`
}

// SearchAheadSlotsResult is a struct to hold the search ahead slots results
type SearchAheadSlotsResult []struct {
	Slot string `db:"slot" json:"slot,omitempty"`
	Root string `db:"blockroot" json:"blockroot,omitempty"`
}

// SearchAheadBlockssResult is a struct to hold the search ahead block results
type SearchAheadBlocksResult []struct {
	Block string `json:"block,omitempty"`
	Hash  string `json:"hash,omitempty"`
}

type SearchAheadTransactionsResult []struct {
	Slot   string `db:"slot" json:"slot,omitempty"`
	TxHash string `db:"txhash" json:"txhash,omitempty"`
}

// SearchAheadGraffitiResult is a struct to hold the search ahead blocks results with a given graffiti
type SearchAheadGraffitiResult []struct {
	Graffiti string `db:"graffiti" json:"graffiti,omitempty"`
	Count    string `db:"count" json:"count,omitempty"`
}

// SearchAheadEth1Result is a struct to hold the search ahead eth1 results
type SearchAheadEth1Result []struct {
	Publickey   string `db:"publickey" json:"publickey,omitempty"`
	Eth1Address string `db:"from_address" json:"address,omitempty"`
}

// SearchAheadValidatorsResult is a struct to hold the search ahead validators results
type SearchAheadValidatorsResult []struct {
	Index  string `db:"index" json:"index,omitempty"`
	Pubkey string `db:"pubkey" json:"pubkey,omitempty"`
}

// SearchAheadPubkeyResult is a struct to hold the search ahead public key results
type SearchAheadPubkeyResult []struct {
	Pubkey string `db:"pubkey" json:"pubkey,omitempty"`
}

// GenericChartData is a struct to hold chart data
type GenericChartData struct {
	IsNormalChart                   bool
	ShowGapHider                    bool
	XAxisLabelsFormatter            template.JS
	TooltipFormatter                template.JS
	TooltipShared                   bool
	TooltipUseHTML                  bool
	TooltipSplit                    bool
	TooltipFollowPointer            bool
	PlotOptionsSeriesEventsClick    template.JS
	PlotOptionsPie                  template.JS
	DataLabelsEnabled               bool
	DataLabelsFormatter             template.JS
	PlotOptionsSeriesCursor         string
	Title                           string                    `json:"title"`
	Subtitle                        string                    `json:"subtitle"`
	XAxisTitle                      string                    `json:"x_axis_title"`
	YAxisTitle                      string                    `json:"y_axis_title"`
	Type                            string                    `json:"type"`
	StackingMode                    string                    `json:"stacking_mode"`
	ColumnDataGroupingApproximation string                    // "average", "averages", "open", "high", "low", "close" and "sum"
	Series                          []*GenericChartDataSeries `json:"series"`
	Drilldown                       interface{}               `json:"drilldown"`
	Footer                          string                    `json:"footer"`
}

type SeriesDataItem struct {
	Name string `json:"name"`
	Y    uint64 `json:"y"`
}

// GenericChartDataSeries is a struct to hold chart series data
type GenericChartDataSeries struct {
	Name  string      `json:"name"`
	Data  interface{} `json:"data"`
	Stack string      `json:"stack,omitempty"`
	Type  string      `json:"type,omitempty"`
	Color string      `json:"color,omitempty"`
}

// ChartsPageDataChart is a struct to hold a chart for the charts-page
type ChartsPageDataChart struct {
	Order  int
	Path   string
	Data   *GenericChartData
	Height int
}

// ChartsPageData is a struct to hold charts for the charts-page and a disclaimer
type ChartsPageData struct {
	ChartsPageDataCharts []ChartsPageDataChart
	Disclaimer           string
}
type HeatmapData struct {
	// BalanceHistory DashboardValidatorBalanceHistory `json:"balance_history"`
	// Earnings       ValidatorEarnings                `json:"earnings"`
	// Validators     [][]interface{}                  `json:"validators"`
	Csrf           string `json:"csrf"`
	ValidatorLimit int    `json:"valLimit"`
	Epochs         []uint64
	Validators     []uint64
	IncomeData     [][3]int64
	MinIncome      int64
	MaxIncome      int64
}

// DashboardData is a struct to hold data for the dashboard-page
type DashboardData struct {
	Csrf           string `json:"csrf"`
	ValidatorLimit int    `json:"valLimit"`
}

// DashboardValidatorBalanceHistory is a struct to hold data for the balance-history on the dashboard-page
type DashboardValidatorBalanceHistory struct {
	Epoch            uint64  `db:"epoch"`
	Balance          uint64  `db:"balance"`
	EffectiveBalance uint64  `db:"effectivebalance"`
	ValidatorCount   float64 `db:"validatorcount"`
}

// ValidatorEarnings is a struct to hold the earnings of one or multiple validators

type ClEl struct {
	El    decimal.Decimal
	Cl    decimal.Decimal
	Total decimal.Decimal
}

type ClElInt64 struct {
	El    int64
	Cl    int64
	Total float64
}

type ClElFloat64 struct {
	El    float64
	Cl    float64
	Total float64 // in ClCurrency if El and Cl have different currencies
}

type ValidatorProposalData struct {
	Proposals                [][]uint64
	BlocksCount              uint64
	ScheduledBlocksCount     uint64
	LastScheduledSlot        uint64
	MissedBlocksCount        uint64
	OrphanedBlocksCount      uint64
	ProposedBlocksCount      uint64
	UnmissedBlocksPercentage float64 // missed/(executed+orphaned+scheduled)
	ProposalLuck             float64
	AvgSlotInterval          *time.Duration
	ProposalEstimate         *time.Time
}

type ValidatorEarnings struct {
	Income1d                ClEl          `json:"income1d"`
	Income7d                ClEl          `json:"income7d"`
	Income31d               ClEl          `json:"income31d"`
	IncomeToday             ClEl          `json:"incomeToday"`
	IncomeTotal             ClEl          `json:"incomeTotal"`
	Apr7d                   ClElFloat64   `json:"apr"`
	Apr31d                  ClElFloat64   `json:"apr31d"`
	Apr365d                 ClElFloat64   `json:"apr365d"`
	TotalDeposits           int64         `json:"totalDeposits"`
	TotalWithdrawals        uint64        `json:"totalWithdrawals"`
	EarningsInPeriodBalance int64         `json:"earningsInPeriodBalance"`
	EarningsInPeriod        int64         `json:"earningsInPeriod"`
	EpochStart              int64         `json:"epochStart"`
	EpochEnd                int64         `json:"epochEnd"`
	LastDayFormatted        template.HTML `json:"lastDayFormatted"`
	LastWeekFormatted       template.HTML `json:"lastWeekFormatted"`
	LastMonthFormatted      template.HTML `json:"lastMonthFormatted"`
	TotalFormatted          template.HTML `json:"totalFormatted"`
	TotalChangeFormatted    template.HTML `json:"totalChangeFormatted"`
	TotalBalance            template.HTML `json:"totalBalance"`
	ProposalData            ValidatorProposalData
}

// ValidatorAttestationSlashing is a struct to hold data of an attestation-slashing
type ValidatorAttestationSlashing struct {
	Epoch                  uint64        `db:"epoch" json:"epoch,omitempty"`
	Slot                   uint64        `db:"slot" json:"slot,omitempty"`
	Proposer               uint64        `db:"proposer" json:"proposer,omitempty"`
	Attestestation1Indices pq.Int64Array `db:"attestation1_indices" json:"attestation1_indices,omitempty"`
	Attestestation2Indices pq.Int64Array `db:"attestation2_indices" json:"attestation2_indices,omitempty"`
}

type ValidatorProposerSlashing struct {
	Epoch         uint64 `db:"epoch" json:"epoch,omitempty"`
	Slot          uint64 `db:"slot" json:"slot,omitempty"`
	Proposer      uint64 `db:"proposer" json:"proposer,omitempty"`
	ProposerIndex uint64 `db:"proposerindex" json:"proposer_index,omitempty"`
}

type ValidatorHistory struct {
	Epoch             uint64                       `db:"epoch" json:"epoch,omitempty"`
	BalanceChange     sql.NullInt64                `db:"balancechange" json:"balance_change,omitempty"`
	AttesterSlot      sql.NullInt64                `db:"attestatation_attesterslot" json:"attester_slot,omitempty"`
	InclusionSlot     sql.NullInt64                `db:"attestation_inclusionslot" json:"inclusion_slot,omitempty"`
	AttestationStatus uint64                       `db:"attestation_status" json:"attestation_status,omitempty"`
	ProposalStatus    sql.NullInt64                `db:"proposal_status" json:"proposal_status,omitempty"`
	ProposalSlot      sql.NullInt64                `db:"proposal_slot" json:"proposal_slot,omitempty"`
	IncomeDetails     *itypes.ValidatorEpochIncome `db:"-" json:"income_details,omitempty"`
	WithdrawalStatus  sql.NullInt64                `db:"withdrawal_status" json:"withdrawal_status,omitempty"`
	WithdrawalSlot    sql.NullInt64                `db:"withdrawal_slot" json:"withdrawal_slot,omitempty"`
}

type ValidatorSlashing struct {
	Epoch                  uint64        `db:"epoch" json:"epoch,omitempty"`
	Slot                   uint64        `db:"slot" json:"slot,omitempty"`
	Proposer               uint64        `db:"proposer" json:"proposer,omitempty"`
	SlashedValidator       *uint64       `db:"slashedvalidator" json:"slashed_validator,omitempty"`
	Attestestation1Indices pq.Int64Array `db:"attestation1_indices" json:"attestation1_indices,omitempty"`
	Attestestation2Indices pq.Int64Array `db:"attestation2_indices" json:"attestation2_indices,omitempty"`
	Type                   string        `db:"type" json:"type"`
}

type StakingCalculatorPageData struct {
	BestValidatorBalanceHistory *[]ValidatorBalanceHistory
	WatchlistBalanceHistory     [][]interface{}
	TotalStaked                 uint64
	// EtherscanApiBaseUrl         string
}

type DepositsPageData struct {
	*Stats
	DepositContract string
	DepositChart    *ChartsPageDataChart
}

type EthOneDepositLeaderBoardPageData struct {
	DepositContract string
}

// EpochsPageData is a struct to hold epoch data for the epochs page
type EthOneDepositsData struct {
	TxHash                []byte    `db:"tx_hash"`
	TxInput               []byte    `db:"tx_input"`
	TxIndex               uint64    `db:"tx_index"`
	BlockNumber           uint64    `db:"block_number"`
	BlockTs               time.Time `db:"block_ts"`
	FromAddress           []byte    `db:"from_address"`
	PublicKey             []byte    `db:"publickey"`
	WithdrawalCredentials []byte    `db:"withdrawal_credentials"`
	Amount                uint64    `db:"amount"`
	Signature             []byte    `db:"signature"`
	MerkletreeIndex       []byte    `db:"merkletree_index"`
	State                 string    `db:"state"`
	ValidSignature        bool      `db:"valid_signature"`
}

type EthOneDepositLeaderboardData struct {
	FromAddress        []byte `db:"from_address"`
	Amount             uint64 `db:"amount"`
	ValidCount         uint64 `db:"validcount"`
	InvalidCount       uint64 `db:"invalidcount"`
	TotalCount         uint64 `db:"totalcount"`
	PendingCount       uint64 `db:"pendingcount"`
	SlashedCount       uint64 `db:"slashedcount"`
	ActiveCount        uint64 `db:"activecount"`
	VoluntaryExitCount uint64 `db:"voluntary_exit_count"`
}

type EthTwoDepositData struct {
	BlockSlot             uint64 `db:"block_slot"`
	BlockIndex            uint64 `db:"block_index"`
	Proof                 []byte `db:"proof"`
	Publickey             []byte `db:"publickey"`
	ValidatorIndex        uint64 `db:"validatorindex"`
	Withdrawalcredentials []byte `db:"withdrawalcredentials"`
	Amount                uint64 `db:"amount"`
	Signature             []byte `db:"signature"`
}

type ValidatorDeposits struct {
	Eth1Deposits      []Eth1Deposit
	LastEth1DepositTs int64
	Eth2Deposits      []Eth2Deposit
}

type MyCryptoSignature struct {
	Address string `json:"address"`
	Msg     string `json:"msg"`
	Sig     string `json:"sig"`
	Version string `json:"version"`
}

type FilterSubscription struct {
	User     uint64
	PriceIds []string
}

type PasswordResetNotAllowedError struct{}

func (e *PasswordResetNotAllowedError) Error() string {
	return "password reset not allowed for this account"
}

type RateLimitError struct {
	TimeLeft time.Duration
}

func (e *RateLimitError) Error() string {
	return fmt.Sprintf("rate limit has been exceeded, %v left", e.TimeLeft)
}

type Empty struct {
}

// GoogleRecaptchaResponse ...
type GoogleRecaptchaResponse struct {
	Success            bool     `json:"success"`
	ChallengeTimestamp string   `json:"challenge_ts"`
	Hostname           string   `json:"hostname"`
	ErrorCodes         []string `json:"error-codes"`
	Score              float32  `json:"score,omitempty"`
	Action             string   `json:"action,omitempty"`
}

type Price struct {
	TS  time.Time `db:"ts"`
	EUR float64   `db:"eur"`
	USD float64   `db:"usd"`
	GBP float64   `db:"gbp"`
	CAD float64   `db:"cad"`
	JPY float64   `db:"jpy"`
	CNY float64   `db:"cny"`
	RUB float64   `db:"rub"`
	AUD float64   `db:"aud"`
}

type ApiStatistics struct {
	Daily      *int `db:"daily"`
	Monthly    *int `db:"monthly"`
	MaxDaily   *int
	MaxMonthly *int
}

type AdConfigurationPageData struct {
	Configurations []*AdConfig
	CsrfField      template.HTML
	New            AdConfig
	TemplateNames  []string
}

type DataTableSaveState struct {
	Key     string                      `json:"key"`
	Time    uint64                      `json:"time"`   // Time stamp of when the object was created
	Start   uint64                      `json:"start"`  // Display start point
	Length  uint64                      `json:"length"` // Page length
	Order   [][]string                  `json:"order"`  // 2D array of column ordering information (see `order` option)
	Search  DataTableSaveStateSearch    `json:"search"`
	Columns []DataTableSaveStateColumns `json:"columns"`
}

func (e *DataTableSaveState) Scan(value interface{}) error {
	b, ok := value.([]byte)
	if !ok {
		return errors.New("type assertion to []byte failed")
	}

	return json.Unmarshal(b, &e)
}

func (a DataTableSaveState) Value() (driver.Value, error) {
	return json.Marshal(a)
}

type DataTableSaveStateOrder struct {
}

type DataTableSaveStateSearch struct {
	Search          string `json:"search"`          // Search term
	Regex           bool   `json:"regex"`           // Indicate if the search term should be treated as regex or not
	Smart           bool   `json:"smart"`           // Flag to enable DataTables smart search
	CaseInsensitive bool   `json:"caseInsensitive"` // Case insensitive flag
}

type DataTableSaveStateColumns struct {
	Visible bool                     `json:"visible"`
	Search  DataTableSaveStateSearch `json:"search"`
}

type Eth1AddressPageData struct {
	Address            string `json:"address"`
	EnsName            string `json:"ensName"`
	IsContract         bool
	QRCode             string `json:"qr_code_base64"`
	QRCodeInverse      string
	Metadata           *Eth1AddressMetadata
	WithdrawalsSummary template.HTML
	BlocksMinedTable   *DataTableResponse
	TransactionsTable  *DataTableResponse
	InternalTxnsTable  *DataTableResponse
	Erc20Table         *DataTableResponse
	Erc721Table        *DataTableResponse
	Erc1155Table       *DataTableResponse
	WithdrawalsTable   *DataTableResponse
	EtherValue         template.HTML
	Tabs               []Eth1AddressPageTabs
}

type ContractInteractionType uint8

const (
	CONTRACT_NONE        ContractInteractionType = 0
	CONTRACT_CREATION    ContractInteractionType = 1
	CONTRACT_PRESENT     ContractInteractionType = 2
	CONTRACT_DESTRUCTION ContractInteractionType = 3
)

type Eth1AddressPageTabs struct {
	Id   string
	Href string
	Text string
	Data *DataTableResponse
}

type Eth1AddressMetadata struct {
	Balances                []*Eth1AddressBalance
	ERC20TokenLimit         uint64
	ERC20TokenLimitExceeded bool
	ERC20                   *ERC20Metadata
	Name                    string
	Tags                    []template.HTML
	EthBalance              *Eth1AddressBalance
}

type Eth1AddressBalance struct {
	Address  []byte
	Token    []byte
	Balance  []byte
	Metadata *ERC20Metadata
}

type ERC20TokenPrice struct {
	Token       []byte
	Price       []byte
	TotalSupply []byte
}

type ERC20Metadata struct {
	Decimals     []byte
	Symbol       string
	Name         string
	Description  string
	Logo         []byte
	LogoFormat   string
	TotalSupply  []byte
	OfficialSite string
	Price        []byte
}

func (metadata ERC20Metadata) MarshalBinary() ([]byte, error) {
	return json.Marshal(metadata)
}

func (metadata ERC20Metadata) UnmarshalBinary(data []byte) error {
	return json.Unmarshal(data, &metadata)
}

type ContractMetadata struct {
	Name    string
	ABI     *abi.ABI `msgpack:"-" json:"-"`
	ABIJson []byte
}

type Eth1TokenPageData struct {
	Token            string `json:"token"`
	Address          string `json:"address"`
	QRCode           string `json:"qr_code_base64"`
	QRCodeInverse    string
	Metadata         *ERC20Metadata
	Balance          *Eth1AddressBalance
	Holders          template.HTML `json:"holders"`
	Transfers        template.HTML `json:"transfers"`
	Price            template.HTML `json:"price"`
	Supply           template.HTML `json:"supply"`
	MarketCap        template.HTML `json:"marketCap"`
	DilutedMarketCap template.HTML `json:"dilutedMarketCap"`
	Decimals         template.HTML `json:"decimals"`
	Contract         template.HTML `json:"contract"`
	WebSite          template.HTML `json:"website"`
	SocialProfiles   template.HTML `json:"socialProfiles"`
	TransfersTable   *DataTableResponse
	HoldersTable     *DataTableResponse
}

type ITransaction struct {
	From      template.HTML
	To        template.HTML
	Amount    template.HTML
	TracePath template.HTML
	Advanced  bool
	Gas       struct {
		Limit uint64
		// Usage    uint64
		// UsedPerc float64
	}
}

type Transfer struct {
	From   template.HTML
	To     template.HTML
	Amount template.HTML
	Token  template.HTML
}

type DepositContractInteraction struct {
	ValidatorPubkey []byte
	WithdrawalCreds []byte
	Amount          []byte
}

type EpochInfo struct {
	Finalized     bool    `db:"finalized"`
	Participation float64 `db:"globalparticipationrate"`
}

type Eth1TxData struct {
	From         common.Address
	To           *common.Address
	InternalTxns []ITransaction
	FromName     string
	ToName       string
	Gas          struct {
		BlockBaseFee   []byte
		MaxFee         []byte
		MaxPriorityFee []byte
		Used           uint64
		UsedPerc       float64
		Limit          uint64
		TxFee          []byte
		EffectiveFee   []byte
	}
	Epoch                       EpochInfo
	TypeFormatted               string
	Type                        uint8
	Nonce                       uint64
	TxnPosition                 uint
	Hash                        common.Hash
	Value                       []byte
	Receipt                     *geth_types.Receipt
	ErrorMsg                    string
	BlockNumber                 int64
	Timestamp                   time.Time
	IsPending                   bool
	TargetIsContract            bool
	IsContractCreation          bool
	CallData                    string
	Method                      string
	Events                      []*Eth1EventData
	Transfers                   []*Transfer
	DepositContractInteractions []DepositContractInteraction
	CurrentEtherPrice           template.HTML
	HistoricalEtherPrice        template.HTML
}

type Eth1EventData struct {
	Address     common.Address
	Name        string
	Topics      []common.Hash
	Data        []byte
	DecodedData map[string]Eth1DecodedEventData
}

type Eth1DecodedEventData struct {
	Type    string
	Value   string
	Raw     string
	Address common.Address
}

type SourcifyContractMetadata struct {
	Compiler struct {
		Version string `json:"version"`
	} `json:"compiler"`
	Language string `json:"language"`
	Output   struct {
		Abi []struct {
			Anonymous bool `json:"anonymous"`
			Inputs    []struct {
				Indexed      bool   `json:"indexed"`
				InternalType string `json:"internalType"`
				Name         string `json:"name"`
				Type         string `json:"type"`
			} `json:"inputs"`
			Name    string `json:"name"`
			Outputs []struct {
				InternalType string `json:"internalType"`
				Name         string `json:"name"`
				Type         string `json:"type"`
			} `json:"outputs"`
			StateMutability string `json:"stateMutability"`
			Type            string `json:"type"`
		} `json:"abi"`
	} `json:"output"`
	Settings struct {
		CompilationTarget struct {
			Browser_Stakehavens_sol string `json:"browser/Stakehavens.sol"`
		} `json:"compilationTarget"`
		EvmVersion string   `json:"evmVersion"`
		Libraries  struct{} `json:"libraries"`
		Metadata   struct {
			BytecodeHash string `json:"bytecodeHash"`
		} `json:"metadata"`
		Optimizer struct {
			Enabled bool  `json:"enabled"`
			Runs    int64 `json:"runs"`
		} `json:"optimizer"`
		Remappings []interface{} `json:"remappings"`
	} `json:"settings"`
	Sources struct {
		Browser_Stakehavens_sol struct {
			Keccak256 string   `json:"keccak256"`
			Urls      []string `json:"urls"`
		} `json:"browser/Stakehavens.sol"`
	} `json:"sources"`
	Version int64 `json:"version"`
}

type EtherscanContractMetadata struct {
	Message string `json:"message"`
	Result  []struct {
		Abi                  string `json:"ABI"`
		CompilerVersion      string `json:"CompilerVersion"`
		ConstructorArguments string `json:"ConstructorArguments"`
		ContractName         string `json:"ContractName"`
		EVMVersion           string `json:"EVMVersion"`
		Implementation       string `json:"Implementation"`
		Library              string `json:"Library"`
		LicenseType          string `json:"LicenseType"`
		OptimizationUsed     string `json:"OptimizationUsed"`
		Proxy                string `json:"Proxy"`
		Runs                 string `json:"Runs"`
		SourceCode           string `json:"SourceCode"`
		SwarmSource          string `json:"SwarmSource"`
	} `json:"result"`
	Status string `json:"status"`
}

type Eth1BlockPageData struct {
	Number          uint64
	PreviousBlock   uint64
	NextBlock       uint64
	TxCount         uint64
	WithdrawalCount uint64
	Hash            string
	ParentHash      string
	MinerAddress    string
	MinerFormatted  template.HTML
	Reward          *big.Int
	TxFees          *big.Int
	GasUsage        template.HTML
	GasLimit        uint64
	LowestGasPrice  *big.Int
	Ts              time.Time
	BaseFeePerGas   *big.Int
	BurnedFees      *big.Int
	BurnedTxFees    *big.Int
	Extra           string
	Txs             []Eth1BlockPageTransaction
	State           string
}

type Eth1BlockPageTransaction struct {
	Hash          string
	HashFormatted template.HTML
	From          string
	FromFormatted template.HTML
	To            string
	ToFormatted   template.HTML
	Value         *big.Int
	Fee           *big.Int
	GasPrice      *big.Int
	Method        string
}

type SlotVizSlots struct {
	BlockRoot []byte
	Epoch     uint64
	Slot      uint64
	Status    string `json:"status"`
	Active    bool   `json:"active"`
}
type SlotVizEpochs struct {
	Epoch          uint64          `json:"epoch"`
	Finalized      bool            `json:"finalized"`
	Justified      bool            `json:"justified"`
	Justifying     bool            `json:"justifying"`
	Particicpation float64         `json:"participation"`
	Slots          []*SlotVizSlots `json:"slots"`
}

type BurnPageDataBlock struct {
	Number        int64     `json:"number"`
	Hash          string    `json:"hash"`
	GasTarget     int64     `json:"gas_target" db:"gaslimit"`
	GasUsed       int64     `json:"gas_used" db:"gasused"`
	Rewards       float64   `json:"mining_reward" db:"miningreward"`
	Txn           int       `json:"tx_count" db:"tx_count"`
	Age           time.Time `json:"time" db:"time"`
	BaseFeePerGas float64   `json:"base_fee_per_gas" db:"basefeepergas"`
	BurnedFees    float64   `json:"burned_fees" db:"burnedfees"`
}

type BurnPageData struct {
	TotalBurned      float64              `json:"total_burned"`
	Blocks           []*BurnPageDataBlock `json:"blocks"`
	BaseFeeTrend     int                  `json:"base_fee_trend"`
	BurnRate1h       float64              `json:"burn_rate_1_h"`
	BurnRate24h      float64              `json:"burn_rate_24_h"`
	BlockUtilization float64              `json:"block_utilization"`
	Emission         float64              `json:"emission"`
	Price            float64              `json:"price_usd"`
	Currency         string               `json:"currency"`
}

type CorrelationDataResponse struct {
	Status  string      `json:"status"`
	Data    interface{} `json:"data"`
	Message string      `json:"message"`
}

type CorrelationData struct {
	Indicator string  `db:"indicator" json:"indicator,omitempty"`
	Time      float64 `db:"time" json:"time,omitempty"`
	Value     float64 `db:"value" json:"value,omitempty"`
}

type EthStoreStatistics struct {
	EffectiveBalances         [][]float64
	TotalRewards              [][]float64
	APRs                      [][]float64
	YesterdayRewards          float64
	YesterdayEffectiveBalance float64
	ProjectedAPR              float64
	YesterdayTs               int64
	StartEpoch                uint64
}

type DilithiumChange struct {
	Slot            uint64 `db:"slot" json:"slot,omitempty"`
	BlockRoot       []byte `db:"block_rot" json:"blockroot,omitempty"`
	Validatorindex  uint64 `db:"validatorindex" json:"validatorindex,omitempty"`
	DilithiumPubkey []byte `db:"pubkey" json:"pubkey,omitempty"`
	Address         []byte `db:"address" json:"address,omitempty"`
	Signature       []byte `db:"signature" json:"signature,omitempty"`
}

type ValidatorsDilithiumChange struct {
	Slot                     uint64 `db:"slot" json:"slot,omitempty"`
	BlockRoot                []byte `db:"block_root" json:"blockroot,omitempty"`
	Validatorindex           uint64 `db:"validatorindex" json:"validatorindex,omitempty"`
	DilithiumPubkey          []byte `db:"pubkey" json:"pubkey,omitempty"`
	Address                  []byte `db:"address" json:"address,omitempty"`
	Signature                []byte `db:"signature" json:"signature,omitempty"`
	WithdrawalCredentialsOld []byte `db:"withdrawalcredentials" json:"withdrawalcredentials,omitempty"`
}

// AdConfig is a struct to hold the configuration for one specific ad banner placement
type AdConfig struct {
	Id              string `db:"id"`
	TemplateId      string `db:"template_id"`
	JQuerySelector  string `db:"jquery_selector"`
	InsertMode      string `db:"insert_mode"`
	RefreshInterval uint64 `db:"refresh_interval"`
	Enabled         bool   `db:"enabled"`
	ForAllUsers     bool   `db:"for_all_users"`
	BannerId        uint64 `db:"banner_id"`
	HtmlContent     string `db:"html_content"`
}

type ExplorerConfigurationCategory string
type ExplorerConfigurationKey string
type ExplorerConfigValue struct {
	Value    string `db:"value"`
	DataType string `db:"data_type"`
}
type ExplorerConfig struct {
	Category ExplorerConfigurationCategory `db:"category"`
	Key      ExplorerConfigurationKey      `db:"key"`
	ExplorerConfigValue
}
type ExplorerConfigurationKeyMap map[ExplorerConfigurationKey]ExplorerConfigValue
type ExplorerConfigurationMap map[ExplorerConfigurationCategory]ExplorerConfigurationKeyMap

func (configMap ExplorerConfigurationMap) GetConfigValue(category ExplorerConfigurationCategory, configKey ExplorerConfigurationKey) (ExplorerConfigValue, error) {
	configValue := ExplorerConfigValue{}
	keyMap, ok := configMap[category]
	if ok {
		configValue, ok = keyMap[configKey]
		if ok {
			return configValue, nil
		}
	}
	return configValue, fmt.Errorf("config value for %s %s not found", category, configKey)
}

func (configMap ExplorerConfigurationMap) GetUInt64Value(category ExplorerConfigurationCategory, configKey ExplorerConfigurationKey) (uint64, error) {
	configValue, err := configMap.GetConfigValue(category, configKey)
	if err == nil {
		if configValue.DataType != "int" {
			return 0, fmt.Errorf("wrong data type for %s %s, got %s, expected %s", category, configKey, configValue.DataType, "int")
		} else {
			return strconv.ParseUint(configValue.Value, 10, 64)
		}
	}
	return 0, err
}

func (configMap ExplorerConfigurationMap) GetStringValue(category ExplorerConfigurationCategory, configKey ExplorerConfigurationKey) (string, error) {
	configValue, err := configMap.GetConfigValue(category, configKey)
	return configValue.Value, err
}

type WithdrawalsPageData struct {
	Stats           *Stats
	WithdrawalChart *ChartsPageDataChart
}

type WithdrawalStats struct {
	WithdrawalsCount                   uint64
	WithdrawalsTotal                   uint64
	DilithiumChangeCount               uint64
	ValidatorsWithDilithiumCredentials uint64
}

type ChangeWithdrawalCredentialsPageData struct {
	FlashMessage string
	CsrfField    template.HTML
	RecaptchaKey string
}

type BroadcastPageData struct {
	Stats        *Stats
	FlashMessage string
	CaptchaId    string
	CsrfField    template.HTML
	RecaptchaKey string
}

type BroadcastStatusPageData struct {
	Job          *NodeJob
	JobTypeLabel string
	JobTitle     string
	JobJson      string
	Validators   *[]NodeJobValidatorInfo
}

type ValidatorIncomePerformance struct {
	ClIncomeWei1d    decimal.Decimal `db:"cl_performance_wei_1d"`
	ClIncomeWei7d    decimal.Decimal `db:"cl_performance_wei_7d"`
	ClIncomeWei31d   decimal.Decimal `db:"cl_performance_wei_31d"`
	ClIncomeWei365d  decimal.Decimal `db:"cl_performance_wei_365d"`
	ClIncomeWeiTotal decimal.Decimal `db:"cl_performance_wei_total"`
	ElIncomeWei1d    decimal.Decimal `db:"el_performance_wei_1d"`
	ElIncomeWei7d    decimal.Decimal `db:"el_performance_wei_7d"`
	ElIncomeWei31d   decimal.Decimal `db:"el_performance_wei_31d"`
	ElIncomeWei365d  decimal.Decimal `db:"el_performance_wei_365d"`
	ElIncomeWeiTotal decimal.Decimal `db:"el_performance_wei_total"`
}

type ValidatorProposalInfo struct {
	Slot            uint64        `db:"slot"`
	Status          uint64        `db:"status"`
	ExecBlockNumber sql.NullInt64 `db:"exec_block_number"`
}
