package climodel

import (
	currency "github.com/0chain/system_test/internal/currency"
	"time"
)

type Provider int

type Timestamp int64

const (
	ProviderMiner Provider = iota + 1
	ProviderSharder
	ProviderBlobber
	ProviderValidator
	ProviderAuthorizer
)

var providerString = []string{"unknown", "miner", "sharder", "blobber", "validator", "authorizer"}

func (p Provider) String() string {
	return providerString[p]
}

type PoolStatus int

const (
	Active PoolStatus = iota
	Pending
	Inactive
	Unstaking
	Deleting
	Deleted
)

var poolString = []string{"active", "pending", "inactive", "unstaking", "deleting"}

func (p PoolStatus) String() string {
	return poolString[p]
}

type Wallet struct {
	ClientID            string `json:"client_id"`
	ClientPublicKey     string `json:"client_public_key"`
	EncryptionPublicKey string `json:"encryption_public_key"`
}

type Allocation struct {
	ID             string    `json:"id"`
	Tx             string    `json:"tx"`
	Name           string    `json:"name"`
	ExpirationDate int64     `json:"expiration_date"`
	DataShards     int       `json:"data_shards"`
	ParityShards   int       `json:"parity_shards"`
	Size           int64     `json:"size"`
	Owner          string    `json:"owner_id"`
	OwnerPublicKey string    `json:"owner_public_key"`
	Payer          string    `json:"payer_id"`
	Blobbers       []Blobber `json:"blobbers"`
	// Stats          *AllocationStats          `json:"stats"`
	TimeUnit    time.Duration `json:"time_unit"`
	IsImmutable bool          `json:"is_immutable"`

	WritePool int64 `json:"write_pool"`

	// BlobberDetails contains real terms used for the allocation.
	// If the allocation has updated, then terms calculated using
	// weighted average values.
	BlobberDetails []*BlobberAllocation `json:"blobber_details"`

	// ReadPriceRange is requested reading prices range.
	ReadPriceRange PriceRange `json:"read_price_range"`

	// WritePriceRange is requested writing prices range.
	WritePriceRange PriceRange `json:"write_price_range"`

	ChallengeCompletionTime time.Duration `json:"challenge_completion_time"`

	StartTime         int64    `json:"start_time"`
	Finalized         bool     `json:"finalized,omitempty"`
	Canceled          bool     `json:"canceled,omitempty"`
	MovedToChallenge  int64    `json:"moved_to_challenge,omitempty"`
	MovedBack         int64    `json:"moved_back,omitempty"`
	MovedToValidators int64    `json:"moved_to_validators,omitempty"`
	Curators          []string `json:"curators"`
}

type AllocationFile struct {
	Name string `json:"name"`
	Path string `json:"path"`
	Type string `json:"type"`
	Size int    `json:"size"`
	Hash string `json:"hash"`
}
type Blobber struct {
	ID      string `json:"id"`
	Baseurl string `json:"url"`
}

type ReadPoolInfo struct {
	Balance int64 `json:"balance"`
}

type ListFileResult struct {
	Name            string    `json:"name"`
	Path            string    `json:"path"`
	Type            string    `json:"type"`
	Size            int64     `json:"size"`
	Hash            string    `json:"hash"`
	Mimetype        string    `json:"mimetype"`
	NumBlocks       int       `json:"num_blocks"`
	LookupHash      string    `json:"lookup_hash"`
	EncryptionKey   string    `json:"encryption_key"`
	ActualSize      int64     `json:"actual_size"`
	ActualNumBlocks int       `json:"actual_num_blocks"`
	CreatedAt       Timestamp `json:"created_at"`
	UpdatedAt       Timestamp `json:"updated_at"`
}

type Terms struct {
	Read_price         int64         `json:"read_price"`
	Write_price        int64         `json:"write_price"`
	Min_lock_demand    float64       `json:"min_lock_demand"`
	Max_offer_duration time.Duration `json:"max_offer_duration"`
}

type Settings struct {
	Delegate_wallet string  `json:"delegate_wallet"`
	Min_stake       int     `json:"min_stake"`
	Max_stake       int     `json:"max_stake"`
	Num_delegates   int     `json:"num_delegates"`
	Service_charge  float64 `json:"service_charge"`
}

type BlobberInfo struct {
	Id                  string   `json:"id"`
	Url                 string   `json:"url"`
	Capacity            int      `json:"capacity"`
	Last_health_check   int      `json:"last_health_check"`
	Allocated           int      `json:"allocated"`
	Terms               Terms    `json:"terms"`
	Stake_pool_settings Settings `json:"stake_pool_settings"`
}

type ChallengePoolInfo struct {
	Id         string `json:"id"`
	Balance    int64  `json:"balance"`
	StartTime  int64  `json:"start_time"`
	Expiration int64  `json:"expiration"`
	Finalized  bool   `json:"finalized"`
}

type FileMetaResult struct {
	Name            string          `json:"Name"`
	Path            string          `json:"Path"`
	Type            string          `json:"Type"`
	Size            int64           `json:"Size"`
	ActualFileSize  int64           `json:"ActualFileSize"`
	LookupHash      string          `json:"LookupHash"`
	Hash            string          `json:"Hash"`
	MimeType        string          `json:"MimeType"`
	ActualNumBlocks int             `json:"ActualNumBlocks"`
	EncryptedKey    string          `json:"EncryptedKey"`
	CommitMetaTxns  []CommitMetaTxn `json:"CommitMetaTxns"`
	Collaborators   []Collaborator  `json:"Collaborators"`
}

type CommitMetaTxn struct {
	RefID     int64  `json:"ref_id"`
	TxnID     string `json:"txn_id"`
	CreatedAt string `json:"created_at"`
}

type Collaborator struct {
	RefID     int64  `json:"ref_id"`
	ClientID  string `json:"client_id"`
	CreatedAt string `json:"created_at"`
}

type CommitResponse struct {
	//FIXME: POSSIBLE ISSUE: json-tags are not available for commit response

	TxnID    string `json:"TxnID"`
	MetaData struct {
		Name            string          `json:"Name"`
		Type            string          `json:"Type"`
		Path            string          `json:"Path"`
		LookupHash      string          `json:"LookupHash"`
		Hash            string          `json:"Hash"`
		MimeType        string          `json:"MimeType"`
		EncryptedKey    string          `json:"EncryptedKey"`
		Size            int64           `json:"Size"`
		ActualFileSize  int64           `json:"ActualFileSize"`
		ActualNumBlocks int             `json:"ActualNumBlocks"`
		CommitMetaTxns  []CommitMetaTxn `json:"CommitMetaTxns"`
		Collaborators   []Collaborator  `json:"Collaborators"`
	} `json:"MetaData"`
}

type PriceRange struct {
	Min int64 `json:"min"`
	Max int64 `json:"max"`
}

type BlobberAllocation struct {
	BlobberID       string `json:"blobber_id"`
	Size            int64  `json:"size"`
	Terms           Terms  `json:"terms"`
	MinLockDemand   int64  `json:"min_lock_demand"`
	Spent           int64  `json:"spent"`
	Penalty         int64  `json:"penalty"`
	ReadReward      int64  `json:"read_reward"`
	Returned        int64  `json:"returned"`
	ChallengeReward int64  `json:"challenge_reward"`
	FinalReward     int64  `json:"final_reward"`
}

type StakePoolInfo struct {
	ID           string                      `json:"pool_id"`      // pool ID
	Balance      int64                       `json:"balance"`      // total balance
	Unstake      int64                       `json:"unstake"`      // total unstake amount
	Free         int64                       `json:"free"`         // free staked space
	Capacity     int64                       `json:"capacity"`     // blobber bid
	WritePrice   int64                       `json:"write_price"`  // its write price
	OffersTotal  int64                       `json:"offers_total"` //
	UnstakeTotal int64                       `json:"unstake_total"`
	Delegate     []StakePoolDelegatePoolInfo `json:"delegate"`
	Penalty      int64                       `json:"penalty"` // total for all
	Rewards      int64                       `json:"rewards"`
	Settings     StakePoolSettings           `json:"settings"`
}

type StakePoolDelegatePoolInfo struct {
	ID         string `json:"id"`          // blobber ID
	Balance    int64  `json:"balance"`     // current balance
	DelegateID string `json:"delegate_id"` // wallet
	Rewards    int64  `json:"rewards"`     // current
	UnStake    bool   `json:"unstake"`     // want to unstake

	TotalReward  int64  `json:"total_reward"`
	TotalPenalty int64  `json:"total_penalty"`
	Status       string `json:"status"`
	RoundCreated int64  `json:"round_created"`
}

type StakePoolSettings struct {
	// DelegateWallet for pool owner.
	DelegateWallet string `json:"delegate_wallet"`
	// MinStake allowed.
	MinStake currency.Coin `json:"min_stake"`
	// MaxStake allowed.
	MaxStake currency.Coin `json:"max_stake"`
	// MaxNumDelegates maximum allowed.
	MaxNumDelegates int `json:"num_delegates"`
	// ServiceCharge is blobber service charge.
	ServiceCharge float64 `json:"service_charge"`
}

type NodeList struct {
	Nodes []Node `json:"Nodes"`
}

type DelegatePool struct {
	Balance      int64  `json:"balance"`
	Reward       int64  `json:"reward"`
	Status       int    `json:"status"`
	RoundCreated int64  `json:"round_created"` // used for cool down
	DelegateID   string `json:"delegate_id"`
}

type StakePool struct {
	Pools    map[string]*DelegatePool `json:"pools"`
	Reward   int64                    `json:"rewards"`
	Settings StakePoolSettings        `json:"settings"`
	Minter   int                      `json:"minter"`
}

type Node struct {
	SimpleNode  `json:"simple_miner"`
	StakePool   `json:"stake_pool"`
	Round       int64 `json:"round"`
	TotalReward int64 `json:"total_reward"`
}

type SimpleNode struct {
	ID         string      `json:"id"`
	N2NHost    string      `json:"n2n_host"`
	Host       string      `json:"host"`
	Port       int         `json:"port"`
	PublicKey  string      `json:"public_key"`
	ShortName  string      `json:"short_name"`
	BuildTag   string      `json:"build_tag"`
	TotalStake int64       `json:"total_stake"`
	Stat       interface{} `json:"stat"`
}

type Sharder struct {
	ID           string `json:"id"`
	Version      string `json:"version"`
	CreationDate int64  `json:"creation_date"`
	PublicKey    string `json:"public_key"`
	N2NHost      string `json:"n2n_host"`
	Host         string `json:"host"`
	Port         int    `json:"port"`
	Path         string `json:"path"`
	Type         int    `json:"type"`
	Description  string `json:"description"`
	SetIndex     int    `json:"set_index"`
	Status       int    `json:"status"`
	Info         struct {
		BuildTag                string `json:"build_tag"`
		StateMissingNodes       int    `json:"state_missing_nodes"`
		MinersMedianNetworkTime int64  `json:"miners_median_network_time"`
		AvgBlockTxns            int    `json:"avg_block_txns"`
	} `json:"info"`
}

type FileStats struct {
	Name                string    `json:"name"`
	Size                int64     `json:"size"`
	PathHash            string    `json:"path_hash"`
	Path                string    `json:"path"`
	NumOfBlocks         int64     `json:"num_of_blocks"`
	NumOfUpdates        int64     `json:"num_of_updates"`
	NumOfBlockDownloads int64     `json:"num_of_block_downloads"`
	NumOfChallenges     int64     `json:"num_of_failed_challenges"`
	LastChallengeTxn    string    `json:"last_challenge_txn"`
	WriteMarkerTxn      string    `json:"write_marker_txn"`
	BlobberID           string    `json:"blobber_id"`
	BlobberURL          string    `json:"blobber_url"`
	BlockchainAware     bool      `json:"blockchain_aware"`
	CreatedAt           time.Time `json:"CreatedAt"`
}

type BlobberDetails struct {
	ID                string            `json:"id"`
	BaseURL           string            `json:"url"`
	Terms             Terms             `json:"terms"`
	Capacity          int64             `json:"capacity"`
	Allocated         int64             `json:"allocated"`
	LastHealthCheck   int64             `json:"last_health_check"`
	PublicKey         string            `json:"-"`
	StakePoolSettings StakePoolSettings `json:"stake_pool_settings"`
}

type Validator struct {
	ID             string  `json:"validator_id"`
	BaseURL        string  `json:"url"`
	PublicKey      string  `json:"-"`
	DelegateWallet string  `json:"delegate_wallet"`
	MinStake       int64   `json:"min_stake"`
	MaxStake       int64   `json:"max_stake"`
	NumDelegates   int     `json:"num_delegates"`
	ServiceCharge  float64 `json:"service_charge"`
	TotalStake     int64   `json:"stake"`
}

type FileDiff struct {
	Op   string `json:"operation"`
	Path string `json:"path"`
	Type string `json:"type"`
}

type FreeStorageMarker struct {
	Assigner   string  `json:"assigner,omitempty"`
	Recipient  string  `json:"recipient"`
	FreeTokens float64 `json:"free_tokens"`
	Timestamp  int64   `json:"timestamp"`
	Signature  string  `json:"signature,omitempty"`
}

type WalletFile struct {
	ClientID    string    `json:"client_id"`
	ClientKey   string    `json:"client_key"`
	Keys        []KeyPair `json:"keys"`
	Mnemonic    string    `json:"mnemonics"`
	Version     string    `json:"version"`
	DateCreated string    `json:"date_created"`
}

type KeyPair struct {
	PublicKey  string `json:"public_key"`
	PrivateKey string `json:"private_key"`
}

type Miner struct {
	ID                string      `json:"id"`
	N2NHost           string      `json:"n2n_host"`
	Host              string      `json:"host"`
	Port              int         `json:"port"`
	PublicKey         string      `json:"public_key"`
	ShortName         string      `json:"short_name"`
	BuildTag          string      `json:"build_tag"`
	TotalStake        int         `json:"total_stake"`
	DelegateWallet    string      `json:"delegate_wallet"`
	ServiceCharge     float64     `json:"service_charge"`
	NumberOfDelegates int         `json:"number_of_delegates"`
	MinStake          int64       `json:"min_stake"`
	MaxStake          int64       `json:"max_stake"`
	Stat              interface{} `json:"stat"`
}

type MinerSCNodes struct {
	Nodes []Node `json:"Nodes"`
}

type MinerSCDelegatePoolInfo struct {
	ID         string `json:"id"`
	Balance    int64  `json:"balance"`
	Reward     int64  `json:"reward"`      // uncollected reread
	RewardPaid int64  `json:"reward_paid"` // total reward all time
	Status     string `json:"status"`
}

type LockConfig struct {
	ID               string           `json:"ID"`
	SimpleGlobalNode SimpleGlobalNode `json:"simple_global_node"`
	MinLockPeriod    int64            `json:"min_lock_period"`
}

type SimpleGlobalNode struct {
	MaxMint     int64   `json:"max_mint"`
	TotalMinted int64   `json:"total_minted"`
	MinLock     int64   `json:"min_lock"`
	Apr         float64 `json:"apr"`
	OwnerId     string  `json:"owner_id"`
}

type MinerSCUserPoolsInfo struct {
	Pools map[string][]*MinerSCDelegatePoolInfo `json:"pools"`
}

type PoolStats struct {
	DelegateID   string `json:"delegate_id"`
	High         int64  `json:"high"` // } interests and rewards
	Low          int64  `json:"low"`  // }
	InterestPaid int64  `json:"interest_paid"`
	RewardPaid   int64  `json:"reward_paid"`
	NumRounds    int64  `json:"number_rounds"`
	Status       string `json:"status"`
}

type TokenPool struct {
	ID      string `json:"id"`
	Balance int64  `json:"balance"`
}

type ZCNLockingPool struct {
	TokenPool `json:"pool"`
}

type SendTransaction struct {
	Status string `json:"status"`
	Txn    string `json:"tx"`
	Nonce  string `json:"nonce"`
}
