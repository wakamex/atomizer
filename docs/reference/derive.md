# Derive API Documentation
Complete API reference from https://docs.derive.xyz/reference/
Generated on: 2025-05-25 17:28:45

## Table of Contents

### Overview
- [How to earn broker fees](#api-broker)
- [Private Endpoints](#authentication)
- [Create Session Keys](#create-session-keys)
- [Error Codes](#error-codes)
- [Fees 1](#fees-1)
- [Interface Vs Manual](#interface-vs-manual)
- [Json Rpc](#json-rpc)
- [Matching Algorithms](#matching-algorithms)
- [Naming](#naming)
- [Register via API](#on-chain-manage-session-keys)
- [Onboard Manually](#onboard-manually)
- [Onboard Via Interface](#onboard-via-interface)
- [Documentation](#overview)
- [Post_Private Session Keys](#post_private-session-keys)
- [Protocol Constants](#protocol-constants)
- [Rate Limits](#rate-limits)
- [Use-case](#session-keys)

### Public
- [Post_Public Build Register Session Key Tx](#post_public-build-register-session-key-tx)
- [Post_Public Create Subaccount Debug](#post_public-create-subaccount-debug)
- [Post_Public Deposit Debug](#post_public-deposit-debug)
- [Post_Public Deregister Session Key](#post_public-deregister-session-key)
- [Post_Public Execute Quote Debug](#post_public-execute-quote-debug)
- [Post_Public Get All Currencies](#post_public-get-all-currencies)
- [Post_Public Get All Instruments](#post_public-get-all-instruments)
- [Post_Public Get All Points](#post_public-get-all-points)
- [Post_Public Get Currency](#post_public-get-currency)
- [Post_Public Get Descendant Tree](#post_public-get-descendant-tree)
- [Post_Public Get Funding Rate History](#post_public-get-funding-rate-history)
- [Post_Public Get Instrument](#post_public-get-instrument)
- [Post_Public Get Instruments](#post_public-get-instruments)
- [Post_Public Get Interest Rate History](#post_public-get-interest-rate-history)
- [Post_Public Get Invite Code](#post_public-get-invite-code)
- [Post_Public Get Latest Signed Feeds](#post_public-get-latest-signed-feeds)
- [Post_Public Get Liquidation History](#post_public-get-liquidation-history)
- [Post_Public Get Live Incidents](#post_public-get-live-incidents)
- [Post_Public Get Maker Program Scores](#post_public-get-maker-program-scores)
- [Post_Public Get Maker Programs](#post_public-get-maker-programs)
- [Post_Public Get Margin](#post_public-get-margin)
- [Post_Public Get Option Settlement History](#post_public-get-option-settlement-history)
- [Post_Public Get Option Settlement Prices](#post_public-get-option-settlement-prices)
- [Post_Public Get Points](#post_public-get-points)
- [Post_Public Get Points Leaderboard](#post_public-get-points-leaderboard)
- [Post_Public Get Spot Feed History](#post_public-get-spot-feed-history)
- [Post_Public Get Spot Feed History Candles](#post_public-get-spot-feed-history-candles)
- [Post_Public Get Ticker](#post_public-get-ticker)
- [Post_Public Get Time](#post_public-get-time)
- [Post_Public Get Trade History](#post_public-get-trade-history)
- [Post_Public Get Transaction](#post_public-get-transaction)
- [Post_Public Get Tree Roots](#post_public-get-tree-roots)
- [Post_Public Get Vault Balances](#post_public-get-vault-balances)
- [Post_Public Get Vault Share](#post_public-get-vault-share)
- [Post_Public Get Vault Statistics](#post_public-get-vault-statistics)
- [Post_Public Login](#post_public-login)
- [Post_Public Margin Watch](#post_public-margin-watch)
- [Post_Public Register Session Key](#post_public-register-session-key)
- [Post_Public Send Quote Debug](#post_public-send-quote-debug)
- [Post_Public Statistics](#post_public-statistics)
- [Post_Public Validate Invite Code](#post_public-validate-invite-code)
- [Post_Public Withdraw Debug](#post_public-withdraw-debug)
- [Public Build_Register_Session_Key_Tx](#public-build_register_session_key_tx)
- [Public Create_Subaccount_Debug](#public-create_subaccount_debug)
- [Public Deposit_Debug](#public-deposit_debug)
- [Public Deregister_Session_Key](#public-deregister_session_key)
- [Public Execute_Quote_Debug](#public-execute_quote_debug)
- [Public Get_All_Currencies](#public-get_all_currencies)
- [Public Get_All_Instruments](#public-get_all_instruments)
- [Public Get_All_Points](#public-get_all_points)
- [Public Get_Currency](#public-get_currency)
- [Public Get_Descendant_Tree](#public-get_descendant_tree)
- [Public Get_Funding_Rate_History](#public-get_funding_rate_history)
- [Public Get_Instrument](#public-get_instrument)
- [Public Get_Instruments](#public-get_instruments)
- [Public Get_Interest_Rate_History](#public-get_interest_rate_history)
- [Public Get_Invite_Code](#public-get_invite_code)
- [Public Get_Latest_Signed_Feeds](#public-get_latest_signed_feeds)
- [Public Get_Liquidation_History](#public-get_liquidation_history)
- [Public Get_Live_Incidents](#public-get_live_incidents)
- [Public Get_Maker_Program_Scores](#public-get_maker_program_scores)
- [Public Get_Maker_Programs](#public-get_maker_programs)
- [Public Get_Margin](#public-get_margin)
- [Public Get_Option_Settlement_History](#public-get_option_settlement_history)
- [Public Get_Option_Settlement_Prices](#public-get_option_settlement_prices)
- [Public Get_Points](#public-get_points)
- [Public Get_Points_Leaderboard](#public-get_points_leaderboard)
- [Public Get_Spot_Feed_History](#public-get_spot_feed_history)
- [Public Get_Spot_Feed_History_Candles](#public-get_spot_feed_history_candles)
- [Public Get_Ticker](#public-get_ticker)
- [Public Get_Time](#public-get_time)
- [Public Get_Trade_History](#public-get_trade_history)
- [Public Get_Transaction](#public-get_transaction)
- [Public Get_Tree_Roots](#public-get_tree_roots)
- [Public Get_Vault_Balances](#public-get_vault_balances)
- [Public Get_Vault_Share](#public-get_vault_share)
- [Public Get_Vault_Statistics](#public-get_vault_statistics)
- [Public Login](#public-login)
- [Public Margin_Watch](#public-margin_watch)
- [Public Register_Session_Key](#public-register_session_key)
- [Public Send_Quote_Debug](#public-send_quote_debug)
- [Public Statistics](#public-statistics)
- [Public Validate_Invite_Code](#public-validate_invite_code)
- [Public Withdraw_Debug](#public-withdraw_debug)

### Private
- [Post_Private Cancel](#post_private-cancel)
- [Post_Private Cancel All](#post_private-cancel-all)
- [Post_Private Cancel Batch Quotes](#post_private-cancel-batch-quotes)
- [Post_Private Cancel Batch Rfqs](#post_private-cancel-batch-rfqs)
- [Post_Private Cancel By Instrument](#post_private-cancel-by-instrument)
- [Post_Private Cancel By Label](#post_private-cancel-by-label)
- [Post_Private Cancel By Nonce](#post_private-cancel-by-nonce)
- [Post_Private Cancel Quote](#post_private-cancel-quote)
- [Post_Private Cancel Rfq](#post_private-cancel-rfq)
- [Post_Private Cancel Trigger Order](#post_private-cancel-trigger-order)
- [Post_Private Change Subaccount Label](#post_private-change-subaccount-label)
- [Post_Private Create Subaccount](#post_private-create-subaccount)
- [Post_Private Deposit](#post_private-deposit)
- [Post_Private Edit Session Key](#post_private-edit-session-key)
- [Post_Private Execute Quote](#post_private-execute-quote)
- [Post_Private Expired And Cancelled History](#post_private-expired-and-cancelled-history)
- [Post_Private Get Account](#post_private-get-account)
- [Post_Private Get All Portfolios](#post_private-get-all-portfolios)
- [Post_Private Get Collaterals](#post_private-get-collaterals)
- [Post_Private Get Deposit History](#post_private-get-deposit-history)
- [Post_Private Get Erc20 Transfer History](#post_private-get-erc20-transfer-history)
- [Post_Private Get Funding History](#post_private-get-funding-history)
- [Post_Private Get Interest History](#post_private-get-interest-history)
- [Post_Private Get Liquidation History](#post_private-get-liquidation-history)
- [Post_Private Get Liquidator History](#post_private-get-liquidator-history)
- [Post_Private Get Margin](#post_private-get-margin)
- [Post_Private Get Mmp Config](#post_private-get-mmp-config)
- [Post_Private Get Notifications](#post_private-get-notifications)
- [Post_Private Get Open Orders](#post_private-get-open-orders)
- [Post_Private Get Option Settlement History](#post_private-get-option-settlement-history)
- [Post_Private Get Order](#post_private-get-order)
- [Post_Private Get Order History](#post_private-get-order-history)
- [Post_Private Get Orders](#post_private-get-orders)
- [Post_Private Get Positions](#post_private-get-positions)
- [Post_Private Get Quotes](#post_private-get-quotes)
- [Post_Private Get Rfqs](#post_private-get-rfqs)
- [Post_Private Get Subaccount](#post_private-get-subaccount)
- [Post_Private Get Subaccount Value History](#post_private-get-subaccount-value-history)
- [Post_Private Get Subaccounts](#post_private-get-subaccounts)
- [Post_Private Get Trade History](#post_private-get-trade-history)
- [Post_Private Get Withdrawal History](#post_private-get-withdrawal-history)
- [Post_Private Liquidate](#post_private-liquidate)
- [Post_Private Order](#post_private-order)
- [Post_Private Order Debug](#post_private-order-debug)
- [Post_Private Poll Quotes](#post_private-poll-quotes)
- [Post_Private Poll Rfqs](#post_private-poll-rfqs)
- [Post_Private Register Scoped Session Key](#post_private-register-scoped-session-key)
- [Post_Private Replace](#post_private-replace)
- [Post_Private Reset Mmp](#post_private-reset-mmp)
- [Post_Private Rfq Get Best Quote](#post_private-rfq-get-best-quote)
- [Post_Private Send Quote](#post_private-send-quote)
- [Post_Private Send Rfq](#post_private-send-rfq)
- [Post_Private Set Cancel On Disconnect](#post_private-set-cancel-on-disconnect)
- [Post_Private Set Mmp Config](#post_private-set-mmp-config)
- [Post_Private Transfer Erc20](#post_private-transfer-erc20)
- [Post_Private Transfer Position](#post_private-transfer-position)
- [Post_Private Transfer Positions](#post_private-transfer-positions)
- [Post_Private Update Notifications](#post_private-update-notifications)
- [Post_Private Withdraw](#post_private-withdraw)
- [Private Cancel](#private-cancel)
- [Private Cancel All](#private-cancel-all)
- [Private Cancel_Batch_Quotes](#private-cancel_batch_quotes)
- [Private Cancel_Batch_Rfqs](#private-cancel_batch_rfqs)
- [Private Cancel_By_Instrument](#private-cancel_by_instrument)
- [Private Cancel_By_Label](#private-cancel_by_label)
- [Private Cancel_By_Nonce](#private-cancel_by_nonce)
- [Private Cancel_Quote](#private-cancel_quote)
- [Private Cancel_Rfq](#private-cancel_rfq)
- [Private Cancel_Trigger_Order](#private-cancel_trigger_order)
- [Private Change_Session_Key_Label](#private-change_session_key_label)
- [Private Change_Subaccount_Label](#private-change_subaccount_label)
- [Private Create_Subaccount](#private-create_subaccount)
- [Private Deposit](#private-deposit)
- [Private Edit_Session_Key](#private-edit_session_key)
- [Private Execute_Quote](#private-execute_quote)
- [Private Expired_And_Cancelled_History](#private-expired_and_cancelled_history)
- [Private Get_Account](#private-get_account)
- [Private Get_All_Portfolios](#private-get_all_portfolios)
- [Private Get_Collaterals](#private-get_collaterals)
- [Private Get_Deposit_History](#private-get_deposit_history)
- [Private Get_Erc20_Transfer_History](#private-get_erc20_transfer_history)
- [Private Get_Funding_History](#private-get_funding_history)
- [Private Get_Interest_History](#private-get_interest_history)
- [Private Get_Liquidation_History](#private-get_liquidation_history)
- [Private Get_Liquidator_History](#private-get_liquidator_history)
- [Private Get_Margin](#private-get_margin)
- [Private Get_Mmp_Config](#private-get_mmp_config)
- [Private Get_Notifications](#private-get_notifications)
- [Private Get_Open_Orders](#private-get_open_orders)
- [Private Get_Option_Settlement_History](#private-get_option_settlement_history)
- [Private Get_Order](#private-get_order)
- [Private Get_Order_History](#private-get_order_history)
- [Private Get_Orders](#private-get_orders)
- [Private Get_Positions](#private-get_positions)
- [Private Get_Quotes](#private-get_quotes)
- [Private Get_Rfqs](#private-get_rfqs)
- [Private Get_Subaccount](#private-get_subaccount)
- [Private Get_Subaccount_Value_History](#private-get_subaccount_value_history)
- [Private Get_Subaccounts](#private-get_subaccounts)
- [Private Get_Trade_History](#private-get_trade_history)
- [Private Get_Withdrawal_History](#private-get_withdrawal_history)
- [Private Liquidate](#private-liquidate)
- [Private Order](#private-order)
- [Private Order_Debug](#private-order_debug)
- [Private Poll_Quotes](#private-poll_quotes)
- [Private Poll_Rfqs](#private-poll_rfqs)
- [Private Register_Scoped_Session_Key](#private-register_scoped_session_key)
- [Private Replace](#private-replace)
- [Private Reset_Mmp](#private-reset_mmp)
- [Private Rfq_Get_Best_Quote](#private-rfq_get_best_quote)
- [Private Send_Quote](#private-send_quote)
- [Private Send_Rfq](#private-send_rfq)
- [Private Session_Keys](#private-session_keys)
- [Private Set_Cancel_On_Disconnect](#private-set_cancel_on_disconnect)
- [Private Set_Mmp_Config](#private-set_mmp_config)
- [Private Transfer_Erc20](#private-transfer_erc20)
- [Private Transfer_Position](#private-transfer_position)
- [Private Transfer_Positions](#private-transfer_positions)
- [Private Update_Notifications](#private-update_notifications)
- [Private Withdraw](#private-withdraw)

### Websocket
- [Auctions Watch](#auctions-watch)
- [Margin Watch](#margin-watch)
- [Orderbook Instrument_Name Group Depth](#orderbook-instrument_name-group-depth)
- [Spot_Feed Currency](#spot_feed-currency)
- [Subaccount_Id Balances](#subaccount_id-balances)
- [Subaccount_Id Orders](#subaccount_id-orders)
- [Subaccount_Id Quotes](#subaccount_id-quotes)
- [Subaccount_Id Trades](#subaccount_id-trades)
- [Subaccount_Id Trades Tx_Status](#subaccount_id-trades-tx_status)
- [Subscribe](#subscribe)
- [Ticker Instrument_Name Interval](#ticker-instrument_name-interval)
- [Trades Instrument_Name](#trades-instrument_name)
- [Trades Instrument_Type Currency](#trades-instrument_type-currency)
- [Trades Instrument_Type Currency Tx_Status](#trades-instrument_type-currency-tx_status)
- [Unsubscribe](#unsubscribe)
- [Wallet Rfqs](#wallet-rfqs)

### Other
- [Solidity Objects](#create-or-deposit-to-subaccount)
- [Testnet](#deposit-to-lyra-chain)
- [Program Purpose](#institutional-trading-rewards-program)
- [0. Constants & Setup](#liquidation-api)
- [Multiple Subaccounts](#multiple-subaccounts)
- [On Chain Withdraw](#on-chain-withdraw)
- [Maintenance Margin](#open-orders-margin)
- [0. Constants & Setup](#rfq-quoting-and-execution)
- [Rfq Quoting And Execution Javascript Copy](#rfq-quoting-and-execution-javascript-copy)
- [Solidity Objects](#submit-order)
- [Submit Order Javascript Copy](#submit-order-javascript-copy)
- [Transfer](#transfer)
- [Transfer Collateral](#transfer-collateral)
- [Ux Create Or Deposit To Subaccount](#ux-create-or-deposit-to-subaccount)
- [Ux Withdraw](#ux-withdraw)

---


# Overview


## api-broker

**Title:** How to earn broker fees
**URL:** https://docs.derive.xyz/reference/api-broker

Individual traders and community partners can earn "broker" fees when generating trading volumes! This is different to the "referral" rewards earned by sharing invite links to new users via UX.

The "broker" fees scheme is a flexible scheme by which partners can earn fees on individual trades that users make (not only in the initial account creation phase).

# How to earn broker fees

Simply place the "Derive Wallet" address of the wallet in which you'd like to receive rewards in the `referral_code` fields of when sending orders via `private/order`. Note, this `referral_code` can be passed in through both the WebSocket and REST APIs.

For more on integrating with the Derive API see [onboarding](/docs/onboard-via-interface)

# Fee tiers

You can instantly begin earning 10% of the trade fees by using the above method.

For more custom integrations, please fill out this [form](https://docs.google.com/forms/d/e/1FAIpQLSe7FD6NIk6fzT0tOoOsS6jlQP1wP0RUkdiAcxseKjQzmZxy-A/viewform?usp=preview)

---

## authentication

**Title:** Private Endpoints
**URL:** https://docs.derive.xyz/reference/authentication

There are two authentication types, both sign-able via the owner or registered session key wallets.

Please refer to the [Derive Python Action Signing SDK](https://pypi.org/project/derive_action_signing/) for actual examples.

# Private Endpoints

All private endpoints and messages starting with `private/` in both REST and Websocket require authentication.

### Scheme

| Param | Description |
| --- | --- |
| `X-LyraWallet` | The Derive wallet address (not "owner") of account. This is NOT your original EOA, but the smart contract wallet on the Derive Chain. To find it in the website go to Home -> Developers -> "Derive Wallet". |
| `X-LyraTimestamp` | Current UTC timestamp in ms |
| `X-LyraSignature` | Keccak-256 signature (standard ETH signing) of `X-LyraTimestamp` using`X-LyraWallet` or registered `session_key` private key |

The authentication scheme uses Ethereum signatures to validate that the sender of the request is either the owner of the account or a registered session key.

### REST

Add the authentication scheme as headers into any `private/` REST request:

TypeScript

```

let wallet = new ethers.Wallet(process.env.OWNER_PRIVATE_KEY as string, provider);
let timestamp = Date.now() // ensure UTC
let signature = (await wallet.signMessage(timestamp)).toString()

const response = await axios.request<R>({
  "POST",
  "https://api-demo.lyra.finance/private/get_subaccounts",
  {wallet: wallet.address},
  {
  	"X-LyraWallet": wallet.address,
  	"X-LyraTimestamp": timestamp,
  	"X-LyraSignature": signature
  }
});

```

### Websocket

Authenticate your websocket session by sending the below `public/login` message:

TypeScript

```
let wallet = new ethers.Wallet(process.env.OWNER_PRIVATE_KEY as string, provider);
let timestamp = Date.now() // ensure UTC
let signature = await wallet.signMessage(timestamp)).toString()

wsc.send(JSON.stringify({
  method: 'public/login',
  params: {
    "wallet": wallet.address,
    "timestamp": timestamp,
    "signature": signature
  },
  id: 1,
}));

```

Private channels require the above login method to be called before they can be subscribed to. Channels are considered private when they reference `{subaccount_id}` or `{wallet}` parameter.

# Self-custodial Requests

Due to the self-custodial nature of the API, the orderbook cannot force the below user actions without an explicit signature by the user:

1. Post Orders
2. Deposit / Withdrawal
3. Transfer

As part of the request, the **owner or a registered session key** must explicitly sign a payload using the private key of the wallet - which is verified via the respective "module" in `Matching.sol`.

Thus, each self-custodial request requires two auth steps:

- pass endpoint authentication (via REST headers or WS login)
- include a sign of the specific payload as one of the API params

**Note** refer to the docs for each request to check the exact param scheme.

> ## 👍 See [Onboard via Interface](/docs/onboard-via-interface) or [Submit Order](/docs/submit-order) guides for example "self-custodial" requests.

# Session Keys

Session keys can be used to sign private requests instead of the owner wallet. See the [Session Keys](https://v2-docs.lyra.finance/reference/session-keys) section

---

## create-session-keys

**Title:** Create Session Keys
**URL:** https://docs.derive.xyz/reference/create-session-keys

For more information on why you need session keys, refer to the "API Reference" > "Session Keys".

Managing session keys via UX can be easily done by going to the "Account Settings" dropdown and clicking "Developers":

![](https://files.readme.io/38b1026-image.png)

Session Keys in the UX are broken down into two types. Under the hood, both have the same access rights, however each is stored in a separate manner:

- Developer Session Keys: private keys of these session keys are never stored in in the client nor server. Common use case for this is to register the session key via UX and then use the private key to sign requests in your trading scripts.
- Device Session Keys: created when a user deposits and trades via the UX. To enable instant deposits and trades, a random session private key is generated, registered and encrypted in the browser's local storage. You can manually revoke device session keys below.

You may revoke both developer and device session keys at any time.

---

## error-codes

**Title:** Error Codes
**URL:** https://docs.derive.xyz/reference/error-codes

| Code | Message | Description |
| --- | --- | --- |
| 0 | `""` | No error (typically omitted from a successful response) |
| -32000 | `"Rate limit exceeded"` | Check per IP rate limits for non-auth requests or your account details for auth requests |
| -32100 | `"Number of concurrent websocket clients limit exceeded"` | Check per IP max concurrent clients for non-auth requests or your account details for auth requests |
| -32700 | `"Parse error"` | Invalid JSON was received by the server. An error occurred on the server while parsing the JSON text. |
| -32600 | `"Invalid Request"` | The JSON sent is not a valid Request object. |
| -32601 | `"Method not found"` | The method does not exist / is not available. |
| -32602 | `"Invalid params"` | Invalid method parameter(s). |
| -32603 | `"Internal error"` | Internal JSON-RPC error. |
| 9000 | `"Order confirmation timeout"` | Order confirmation timed out but order status is unknown. Please check status of open orders. |
| 9001 | `"Engine confirmation timeout"` | Order confirmation timed out while processing within the matching engine. The order has been dropped |
| 10000 | `"Manager not found"` | The requested margin type does not exist. |
| 10001 | `"Asset is not an ERC20 token"` | The requested asset is not an ERC20 token and therefore cannot be deposited/withdrawn/transferred. |
| 10002 | `"Sender and recipient wallet do not match"` | Transfers can only be made to subaccounts under the same wallet. |
| 10003 | `"Sender and recipient subaccount IDs are the same"` | Transfers can only be made to a different subaccount id. |
| 10004 | `"Multiple currencies not supported"` | Only standard margin accounts may hold assets of multiple currencies. |
| 10005 | `"Maximum number of subaccounts per wallet reached"` | Withdraw any unused subaccounts to add new ones. |
| 10006 | `"Maximum number of session keys per wallet reached"` | Deactivate any unused session keys to add new ones. |
| 10007 | `"Maximum number of assets per subaccount reached"` | Number of assets in a subaccount is limited by the on-chain constraints. |
| 10008 | `"Maximum number of expiries per subaccount reached"` | Number of expiries in a portfolio margin subaccount is limited by the on-chain constraints. |
| 10009 | `"Recipient subaccount ID of the transfer cannot be 0"` | Transfers must be made to registerred non-zero subaccounts. |
| 10010 | `"PMRM only supports USDC asset collateral. Cannot trade spot markets."` | PMRM only supports USDC asset collateral. Cannot trade spot markets. |
| 10011 | `"ERC20 allowance is insufficient"` | ERC20 allowance is insufficient for the requested action. |
| 10012 | `"ERC20 balance is less than transfer amount"` | ERC20 balance is insufficient for the requested action. |
| 10013 | `"There is a pending deposit for this asset"` | There is a pending deposit for this asset within the last 15 minutes. Please wait and try again. |
| 10014 | `"There is a pending withdrawal for this asset"` | There is a pending withdrawal for this asset within the last 15 minutes. Please wait and try again. |
| 10015 | `"PortfolioMargin2 supports multiple collaterals but only options and perps of the same currency"` | PortfolioMargin2 supports multiple collaterals but only options and perps of the same currency |
| 11000 | `"Insufficient funds"` | Insufficient funds to place order. |
| 11002 | `"Order rejected from queue"` | Order did not reach the queue, matching engine might be down, or order was in the queue for too long. |
| 11003 | `"Already cancelled"` | Order is already cancelled. |
| 11004 | `"Already filled"` | Order is already filled. |
| 11005 | `"Already expired"` | Order is already expired. |
| 11006 | `"Does not exist"` | Order does not exist. |
| 11007 | `"Self-crossing disallowed"` | Order was rejected because it crossed with another order placed by the same user. |
| 11008 | `"Post-only reject"` | A post-only order was rejected because it would have matched with an existing order. |
| 11009 | `"Zero liquidity for market or IOC/FOK order"` | A market or an IOC/FOK order was rejected because there was no liquidity within the provided limit price. |
| 11010 | `"Post-only invalid order type"` | A post-only order was rejected because it was not a limit order. |
| 11011 | `"Invalid signature expiry"` | Order expiry is above or below the min/max validity or is beyond expiry (for options only). |
| 11012 | `"Invalid amount"` | Order amount is invalid, see data for details. |
| 11013 | `"Invalid limit price"` | Order limit price is invalid, see data for details. |
| 11014 | `"Fill-or-kill not filled"` | A fill-or-kill order was not filled. |
| 11015 | `"MMP frozen"` | An order was rejected because the market maker protections were triggered. |
| 11016 | `"Already consumed"` | Order is already consumed (Filled/expired/rejected). |
| 11017 | `"Non unique nonce"` | This nonce has already been used, please use a new nonce. |
| 11018 | `"Invalid nonce date"` | First 10 digits of nonce must be a UTC sec timestamp within 1 hour of the true UTC timestamp. |
| 11019 | `"Open orders limit exceeded"` | Too many open orders for this subaccount. |
| 11020 | `"Negative ERC20 balance"` | Wrapped ERC20 balances cannot be negative. |
| 11021 | `"Instrument is not live"` | Instrument has either been deactivated or expired (if an Option) |
| 11022 | `"Reject timestamp exceeded"` | Order was rejected because it reached the engine after the supplied `reject_timestamp`. |
| 11023 | `"Max fee order param is too low"` | Max fee order param must always be >= 2 x max(taker, maker) x spot\_price. If the order crosses the book, it must be >= 2 x max(taker, maker) x spot\_price + base\_fee / fill\_amount. |
| 11024 | `"Reduce only not supported with this time in force"` | Reduce only orders have to be market orders or non-resting limit orders (IOC or FOK). |
| 11025 | `"Reduce only reject"` | A reduce-only order was rejected because it would have increased position size. |
| 11026 | `"Transfer reject"` | A transfer was rejected because it would have increased a subaccount position size. |
| 11027 | `"Subaccount undergoing liquidation"` | A trade or order was rejected because the subaccount is undergoing a liquidation. |
| 11028 | `"Replaced order filled amount does not match expected state."` | New create order was reverted as the filled amount of the old order does not match the expected filled amount. |
| 11050 | `"Trigger order was cancelled between the time worker sent order and engine processed order"` | Trigger order was not placed as it was cancelled right before entering engine |
| 11051 | `"Trigger price must be higher than the current price for stop orders and vice versa for take orders"` | Make sure the trigger price is properly set |
| 11052 | `"Trigger order limit exceeded (separate limit from regular orders)"` | Too many trigger orders for this subaccount. |
| 11053 | `"Index and last-trade trigger price types not supported yet"` | Only mark price is supported for now |
| 11054 | `"Trigger orders cannot replace or be replaced"` | Trigger orders cannot replace or be replaced |
| 11055 | `"Market order limit_price is unfillable at the given trigger price"` | For market trigger orders, make sure limit price is crossable |
| 11100 | `"Leg instruments are not unique"` | Leg provided in the RFQ or Quote parameter must not have repeated instrument names. |
| 11101 | `"RFQ not found"` | RFQ query or cancellation failed because nothing was found with the given filters. |
| 11102 | `"Quote not found"` | Quote query or cancellation failed because nothing was found with the given filters. |
| 11103 | `"Quote leg does not match RFQ leg"` | Legs provided in quote parameters must match the legs in the RFQ. |
| 11104 | `"Requested quote or RFQ is not open"` | Quote submission failed because the RFQ is either expired, filled or cancelled. |
| 11105 | `"Requested quote ID references a different RFQ ID"` | Quote execution failed because the RFQ ID parameter does not match the RFQ ID in the quote. |
| 11106 | `"Invalid RFQ counterparty"` | RFQ submission failed because the counterparty is not authorized to act as an RFQ maker or does not exist. |
| 11107 | `"Quote maker total cost too high"` | Quote submission failed because the maker total cost exceeded price bandwidth. |
| 11200 | `"Auction not ongoing"` | Supplied liquidated subaccount has no ongoing auction. |
| 11201 | `"Open orders not allowed"` | Bidding subaccount is not allowed to have open orders. |
| 11202 | `"Price limit exceeded"` | Supplied bid price limit is too low for this auction. |
| 11203 | `"Last trade ID mismatch"` | Liquidated subaccount has a different last trade ID. |
| 12000 | `"Asset not found"` | Requested asset does not exist. |
| 12001 | `"Instrument not found"` | Requested instrument does not exist. |
| 12002 | `"Currency not found"` | Requested currency does not exist. |
| 12003 | `"USDC does not have asset caps per manager"` | USDC does not have asset caps per manager |
| 13000 | `"Invalid channels"` | All channels in the subscribe request were invalid. |
| 14000 | `"Account not found"` | Requested wallet does not exist. |
| 14001 | `"Subaccount not found"` | Requested subaccount does not belong to a registered wallet. |
| 14002 | `"Subaccount was withdrawn"` | Requested subaccount was withdrawn. Please re-deposit subaccount. |
| 14008 | `"Cannot reduce expiry using registerSessionKey RPC route"` | Must use the deregisterSessionKey RPC route to reduce expiry. |
| 14009 | `"Session key expiry must be > utc_now + 10 min"` | Increase the expiry\_sec. |
| 14010 | `"Session key already registered for this account"` | Requested session key is already registered for this account. |
| 14011 | `"Session key already registered with another account"` | Requested session key is already registered with another account. |
| 14012 | `"Address must be checksummed"` | Requested address must be checksumed |
| 14013 | `"String is not a valid ethereum address"` | String must be a valid ethereum address: e.g. 0xd3cda913deb6f67967b99d67acdfa1712c293601 |
| 14014 | `"Signature invalid for message or transaction"` | Address recovered from message and signature does not match the signer |
| 14015 | `"Transaction count for given wallet does not match provided nonce"` | Ensure the nonce is set correctly |
| 14016 | `"The provided signed raw transaction contains function name that does not match the expected function name"` | Ensure that the right contract abi was used when generating the transaction |
| 14017 | `"The provided signed raw transaction contains contract address that does not match the expected contract address"` | Ensure that the right contract address was used when generating the transaction |
| 14018 | `"The provided signed raw transaction contains function params that do not match any expected function params"` | Ensure that the right contract abi was used when generating the transaction |
| 14019 | `"The provided signed raw transaction contains function param values that do not match the expected values"` | Ensure that the signed function inputs match the ones provided in the request |
| 14020 | `"The X-LyraWallet header does not match the requested subaccount_id or wallet"` | Ensure the X-LyraWallet header is set to match the requested subaccount\_id or wallet |
| 14021 | `"The X-LyraWallet header not provided"` | Ensure the X-LyraWallet header is included in the request |
| 14022 | `"Subscription to a private channel failed"` | A private channel is not authorized for this websocket connection, or it was requested over a public method |
| 14023 | `"Signer in on-chain related request is not wallet owner or registered session key"` | Address of the signer must be a wallet owner or registered session key |
| 14024 | `"Chain ID must match the current roll up chain id"` | Refer to v2-docs for chain id for each deployment |
| 14025 | `"The private request is missing a wallet or subaccount_id param"` | Ensure request params include a subaccount\_id or wallet |
| 14026 | `"Session key not found"` | Requested session key does not exist. |
| 14027 | `"Unauthorized as RFQ maker"` | The account is not authorized to act as an RFQ maker. |
| 14028 | `"Cross currency RFQ not supported"` | RFQs only support same-currency legs at this moment. |
| 14029 | `"Session key IP not whitelisted"` | IP address of the request is not whitelisted for the session key. |
| 14030 | `"Session key expired"` | Session key has expired. |
| 14031 | `"Unauthorized key scope"` | The key scope provided does not meet the minimum required scope. |
| 14032 | `"Scope should not be changed"` | The key scope of the registered session key is not admin. It should not be elevated for security purposes |
| 14033 | `"Account not whitelisted for atomic orders"` | Account not whitelisted for atomic orders while trying to submit an atomic order. |
| 16000 | `"You are in a restricted region that violates our terms of service."` | You are in a restricted region that violates our terms of service. You may withdraw funds any time but deposits, transfers, orders are blocked |
| 16001 | `"Account is disabled due to compliance violations, please contact support to enable it."` | Account is disabled due to compliance violations, please contact support to enable it. |
| 16100 | `"Sentinel authorization is invalid"` | The sentinel authorization is invalid |
| 17000 | `"This accoount does not have a shareable invite code"` | This accoount does not have a shareable invite code |
| 17001 | `"Invalid invite code"` | Invalid invite code |
| 17002 | `"Invite code already registered for this account"` | Invite code already registered for this account |
| 17003 | `"Invite code has no remaining uses"` | Invite code has no remaining uses |
| 17004 | `"Requirement for successful invite registration not met"` | Requirement for successful invite registration not met |
| 17005 | `"Account must register with a valid invite code to be elligible for points"` | Account must register with a valid invite code to be elligible for points |
| 17006 | `"Point program does not exist"` | Point program does not exist |
| 17007 | `"Invalid leaderboard page number"` | Invalid leaderboard page number |
| 18000 | `"Invalid block number"` | Invalid block number |
| 18001 | `"Failed to estimate block number. Please try again later."` | Failed to estimate block number. Please try again later. |
| 18002 | `"The provided smart contract owner does not match the wallet in LightAccountFactory.getAddress()"` | The provided smart contract owner does not match the wallet in LightAccountFactory.getAddress() |
| 18003 | `"Vault ERC20 asset does not exist"` | Vault ERC20 asset does not exist |
| 18004 | `"Vault ERC20 pool does not exist"` | Vault ERC20 pool does not exist |
| 18005 | `"Must add asset to pool before getting balances"` | Must add asset to pool before getting balances |
| 18006 | `"Invalid Swell season. Swell seasons are in the form 'swell_season_X'."` | Invalid Swell season. Swell seasons are in the form 'swell\_season\_X'. |
| 18007 | `"Vault not found"` | Vault not found |
| 19000 | `"Maker program not found"` | Maker program not found |

---

## fees-1

**Title:** Fees 1
**URL:** https://docs.derive.xyz/reference/fees-1

## Orderbook

Fees on Derive differ if you are are maker (i.e. you put out a resting limit order that is filled) or if you are a taker (i.e. you buy or sell against an existing order). The fees also differ by instrument (options vs perps), and are set out in the table below:

| Instrument | Maker | Taker |
| --- | --- | --- |
| Perp | 0.01% x notional volume | $0.1 + 0.03% x notional volume |
| Option | 0.01% x notional volume | $0.5 + 0.03% x notional volume |

Option notional fees are capped at 12.5% of the value of the option.

Note that takers pay an additional `base_fee` of $0.10 (Perps) or $0.5 (Options) per trade regardless of their trade size. This fee is waived for verified market maker accounts.

Examples:

1. Alice buys 2 ETH 2,000 puts using an aggressive order, and the oracle spot price is $2200:  
   `fee = $0.5 + 0.03% * 2 * $2200 = $1.82`
2. Bob opens a 0.1 BTC perp sell limit order and later gets filled by Charlie, with spot at $43,000:  
   `feeBob = 0.01% * 0.1 * $43,000 = $0.43`  
   `feeCharlie = $0.1 + 0.03% * 0.1 * $43,000 = $1.39`

## RFQs

Trades conducted via RFQs get charged the taker notional fee rate to both counterparties (and a base fee to the taker side). Additionally, multi-leg trades **enjoy up to 100% discounts** on the cheaper of the legs based on the below explained set of rules. In summary, for the most common use cases:

- 2-leg option spreads like straddles, verticals, calendars, etc. get charged zero fee on their second leg
- Hedged options (option + perpetual) get charged zero fee on the cheapest between the perp and the option leg

For more complicated trades, the following rules apply. All legs of an RFQ get grouped into `long calls`, `long puts`, `short calls`, `short puts`, `perps`, and total fee is calculated within each group. The full fee is always charged on the most expensive group, no matter how many groups there are. The trades in the remaining groups get a fee discount applied to them in the following order:

1. Cheapest group gets 100% discount
2. Second and third cheapest group gets 50% discount
3. Other groups do not get any further discounts

For example:

1. A call spread has 2 legs, and the two legs belong to different groups (one in `long calls` one in `short calls`). The cheapest leg / group gets a 100% discount.
2. A straddle or a strangle has two legs, and the two legs belong to different groups. The cheapest leg gets a 100% discount.
3. Two long calls bought at different strikes or expiries both fall into the same `long calls` group and therefore have no discount applied to any of them.
4. A risk reversal with a perp hedge has 3 legs each falling in a different group. The cheapest leg gets 100% discount, second cheapest gets 50% discount, and the most expensive leg is paid in full.

### Box Spreads

Additionally, the system recognizes box spreads (a 4-legged trade with a long call and a short put at one strike, and a short call and a long put at another strike, all at the same expiry) as a special strategy with a different fee schedule. A box spread can be thought of as a zero coupon bond paying the notional `(strike_1 - strike_2)`dollars at expiry, and it typically trades at a discount to its notional.

Derive charges a "yield spread" fee for this "bond" equal to `notional x 1% x years_to_expiry`, e.g. for a box with strikes $4000 and $5000 ($1000 notional) and 1 month to expiry, the fee would be `$1000 x 1% x 1/12 = $0.83`. This fee is charged to both maker and taker side, plus a $0.5 base fee is charged to the taker side only.

---

## interface-vs-manual

**Title:** Interface Vs Manual
**URL:** https://docs.derive.xyz/reference/interface-vs-manual

The Derive API can handle both on-chain and pure orderbook requests.

This means that users can onboard onto the API via:

### Interface

- Onboard with no code or API requests (see [Onboard via Interface](/docs/onboard-via-interface))
- Continue to interact with the exchange and protocol via Interface or Session Keys
- A smart-contract wallet is created as the **proxy owner of the account** which the original ETH wallet controls
- Creating account, depositing, withdrawals and session key creation all possible via UX
- Gives ability to monitor open trades / orders / history and even execute trades via UX + API
- For mainnet: [www.derive.xyz](https://www.derive.xyz)
- For testnet: [testnet.derive.xyz](https://testnet.derive.xyz/)

### Manually

- Onboard completely through a combination of calls to on-chain contracts and requests to the API (see [Onboard Manually](/docs/onboard-manually))
- Removes ability to interact with account via Interface.
- Removes the "smart-contract" wallet layer

Refer to either option in [Onboard via Interface](/docs/onboard-via-interface) or [Onboard Manually](/docs/onboard-manually) sections.

---

## json-rpc

**Title:** Json Rpc
**URL:** https://docs.derive.xyz/reference/json-rpc

[JSON-RPC](https://www.jsonrpc.org/specification) is a standard RPC (remote procedure call) protocol widely used both by the exchanges and the native Ethereum ecosystem. Derive API is built on top of this specification with minor adjustments.

This API protocol is transport agnostic: it can be used both over Websockets and HTTP. The request parameters and method names are equivalent over both protocols, however HTTP does not support subscriptions.

## WebSocket

In WebSockets, clients must send messages as JSON objects in the with the following fields:

JSON

```
{
  "id": string,
  "method": string,
  "params": object
}

```

**NOTE** the WebSocket URL is `wss://api-demo.lyra.finance/ws`

## HTTP Post

When used with HTTP POST requests, client must send messages to the endpoint's specified url with the appropriate path, and attach a POST payload that matches the `params` schema. See REST API for examples.

When used with HTTP GET requests, client must send messages over the same base url, and provide the `params` object via url-encoded query parameters. See REST API for examples.

Regardless of the chosen transport protocol, the server will always send back messages in the form:

JSON

```
{
  "id": string,
  "result": object
}

```

if the RPC call was successful, or:

JSON

```
{
  "id": string,
  "error": {
    "code": number,
    "message": string, 
    "data": string
  }
}

```

if the RPC call resulted in an error. [RPC Error Codes](/reference/error-codes) are shared across HTTP and Websocket APIs.

---

## matching-algorithms

**Title:** Matching Algorithms
**URL:** https://docs.derive.xyz/reference/matching-algorithms

Derive supports both FIFO (a.k.a. price/time) and pro-rata matching algorithms, as well as "blends" thereof (e.g. a % of the order being matched FIFO and a % pro-rata).

### FIFO

In FIFO matching, resting orders are ranked by price first, then by the order creation time. Order with the smallest creation timestamp gets top priority amongst all orders of the same price.

### Pro-rata

In pro-rata matching, resting orders are ranked by price first. If several orders have the same price, then the share of the incoming taker order they get is determined pro-rata by the orders' sizes.

For example, suppose Alice and Bob are quoting 10 and 30 contracts respectively at a price of $150, and Charlie sends a taker order of size 20. Alice's share of the totals size is `10 / (10 + 30) = 0.25`, therefore her order gets filled `20 x 0.25 = 5` contracts. Bob's share is `30 / (10 + 30) = 0.75` so he gets filled the remaining `20 x 0.75 = 15`.

### FIFO & pro-rata blend

The above 2 algorithms are building blocks of Derive's FIFO & pro-rata blend. At a high level, the algorithm performs the following 3 steps at every price level:

1. FIFO pass: a certain % or size of an order gets routed through regular FIFO.
2. Pro-rata pass: the remainder gets filled pro-rata, some rounding is applied to the fill share of every participating order.
3. FIFO cleanup pass: the part of the order that still remains unfilled due to rounding performed in (2) gets routed FIFO.

There are 3 parameters governing the degree of this blend (available in `get_instrument` and `get_ticker` payloads and channels)

1. `pro_rata_fraction` determines the maximum % of the order's size that can get filled pro-rata. If this number is zero, the algorithm is equivalent to full FIFO.
2. `fifo_min_allocation` determines the minimum order size threshold that will get routed through FIFO no matter what the parameter in (1) is set to. This adds an incentive for market makers to better the current market, since they are the first to improve the price, small flow will go more towards them. For example, if this value is 5 and `pro_rata_fraction` is 80%, an order of size 10 will have a size of `max(5, 10 x 20%) = 5` routed through FIFO and the remaining 5 through pro-rata.
3. `pro_rata_amount_step` determines the rounding of the fill shares under pro-rata, for example suppose this value is 1, and if Alice's pro-rata share of an order is 25%, and the order is of size 5. Then Alice's unrounded share of the order is 1.25, which gets rounded down to 1. This rounding would happen for every pro-rata fill participant, and the unfilled portions (the "rounding errors") are added up and routed FIFO.

## Algorithms for Products

1. Perpetuals are matched through regular FIFO
2. Options are matched through the blend algorithm with the following parameters:

| Parameter | ETH | BTC |
| --- | --- | --- |
| `pro_rata_fraction` | 0.8 | 0.8 |
| `fifo_min_allocation` | 10 | 1 |
| `pro_rata_amount_step` | 1 | 0.1 |

---

## naming

**Title:** Naming
**URL:** https://docs.derive.xyz/reference/naming

Derive products use the following naming system:

| Product | Template | Examples |
| --- | --- | --- |
| Quote / Base | $TICKER | ETH, BTC, USDC |
| Perpetual | $TICKER-PERP | ETH-PERP |
| Option | $TICKER-$DAY$MONTH$YEAR-$STRIKE-C|P | BTC-20250316-420-C,  BTC-20230916-580-P |
| Spot | $BASE-$QUOTE | ETH-USDC, BTC-USDC |

---

## on-chain-manage-session-keys

**Title:** Register via API
**URL:** https://docs.derive.xyz/reference/on-chain-manage-session-keys

Refer to the "Session Keys" section in the "API Reference" for more information on the nature of session keys.

# Register via API

To avoid the overhead of managing a connection to the Derive Rollup via RPC, you can submit session key registration and deregistration directly through the orderbook. Transaction cost is still paid by the signer.

***NOTE*** Need to wait up to a minute after tx submission for the API to allow session key usage.

TypeScript

```
let wallet = new ethers.Wallet(process.env.OWNER_PRIVATE_KEY as string, provider);
let newSessionKey = ethers.Wallet.createRandom();
let expirySec = Date.now() / 1000 + 3600 // valid for 1hr

// get API to build unsigned tx
const buildTxResponse = await axios.request<R>({
  "POST"
  "https://api-demo.lyra.finance/public/build_register_session_key_tx,
  {
    wallet: wallet.address,
    public_session_key: newSessionKey.address,
    expiry_sec: expirySec, 
    nonce: await wallet.getNonce(),
    gas: ethers.toBigInt(5000000),
	}
});

// API submits tx on-chain
const registerResponse = await axios.request<R>({
  "POST",
  "https://api-demo.lyra.finance/public/register_session_key,
  {
    wallet: wallet.address,
    label: 'my_label',
    public_session_key: newSessionKey.address,
    expiry_sec: expirySec, // 1 hour
    signed_raw_tx: await wallet.signTransaction(
      buildSessionKeyTxResponse.data.result.tx_params
    );
	};
});

```

The same flow can be used with the `public/deregister_session_key` endpoint to delete a session key.

# Register Directly On-chain

TypeScript

```
let wallet = new ethers.Wallet(process.env.OWNER_PRIVATE_KEY as string, provider);
let newSessionKey = ethers.Wallet.createRandom();
let expirySec = Date.now() / 1000 + 3600 // valid for 1hr
let matchingABI = ["function registerSessionKey(address toAllow, uint256 expiry)"] 

const matchingContract = new ethers.Contract(
  process.env.MATCHING_ADDRESS,
  matchingABI,
  provider
)

let tx = await matching.connect(wallet).registerSessionKey.(
  newSessionKey.address, expirySec
)

```

# Deregister Directly On-chain

Once the transaction is submitted, the session key becomes unusable after a **10 minute cooldown period**

TypeScript

```
let wallet = new ethers.Wallet(process.env.OWNER_PRIVATE_KEY as string, provider);
let newSessionKey = ethers.Wallet.createRandom();
let expirySec = Date.now() / 1000 + 3600 // valid for 1hr
let matchingABI = ["function deregisterSessionKey(address sessionKey)"] 

const matchingContract = new ethers.Contract(
  process.env.MATCHING_ADDRESS,
  matchingABI,
  provider
)

let tx = await matching.connect(wallet).deregisterSessionKey.(
  newSessionKey.address
)

```

---

## onboard-manually

**Title:** Onboard Manually
**URL:** https://docs.derive.xyz/reference/onboard-manually

1. [Deposit to Derive Chain](/docs/deposit-to-lyra-chain)
2. [Create or Deposit to Subaccount](/docs/create-or-deposit-to-subaccount)
3. [Manage Session Keys](/docs/on-chain-manage-session-keys)

---

## onboard-via-interface

**Title:** Onboard Via Interface
**URL:** https://docs.derive.xyz/reference/onboard-via-interface

1. [Create Subaccount and Deposit](/docs/ux-create-or-deposit-to-subaccount)
2. [Manage Session Keys](/docs/create-session-keys)
3. [Multiple Subaccounts](/docs/multiple-subaccounts)
4. [Transfer](/docs/transfer)
5. [Withdraw](/docs/ux-withdraw)

---

## overview

**Title:** Documentation
**URL:** https://docs.derive.xyz/reference/overview

Derive is a self-custodial, high performance crypto trading platform for options, perpetuals and spot trading.

The trading platform is made up of three components:

1. **Derive Chain:** A settlement layer for transactions. This is an Optimistic rollup built on the [OP Stack](https://stack.optimism.io/), secured by Ethereum mainnet. Governed by the [Derive DAO](https://gov.lyra.finance/).
2. **Derive Protocol:** A settlement protocol that enables permissionless, self-custodial margin trading of perpetuals, options and spot, deployed to the Derive Chain (formerly Lyra Chain). Governed by the [Derive DAO](https://gov.lyra.finance/).
3. **Derive Exchange:** An orderbook that efficiently matches orders and settles them to the Derive Protocol. Operated by Lyra Trading Co.

The following docs describe Derive technical concepts relating to the Protocol and Exchange. For a deep dive into the Derive Chain, consult the [OP Stack Docs](https://stack.optimism.io/). For a deep dive into the Derive DAO and Derive's governance framework, consult the [Governance Docs](https://gov.lyra.finance). Note that Derive operates independently of the Lyra V1 AMM.

> ## 📘 The Derive Exchange uses a centralized limit order book, but remains self-custodial, and settles trades and liquidations in a trustless manner.

# Documentation

Visit [Documentation](/docs/introduction) for Onboarding Guides and a deep dive into the Derive Protocol's standard margin and portfolio margin rules, as well as an overview of supported products, liquidations and oracles:

- [Interface vs Manual Onboarding](/docs/introduction-1)
- [Supported Products](/docs/supported-products)
- [Standard Margin](/docs/standard-margin)
- [Portfolio Margin](/docs/portfolio-margin)
- [Liquidations](/docs/liquidations)
- [Oracles](/docs/oracles)

# Derive API

The Derive API provides access to the Derive Exchange orderbook which matches orders and settles trades. Derive provides two interfaces to access the API:

### Mainnet

- [JSON-RPC over HTTP](/reference/public) at <https://api.lyra.finance>
- [JSON-RPC over Websocket](/reference/subscribe) at wss://api.lyra.finance/ws

### Testnet

- [JSON-RPC over HTTP](/reference/public) at <https://api-demo.lyra.finance>
- [JSON-RPC over Websocket](/reference/subscribe) at wss://api-demo.lyra.finance/ws

The API v2.0-alpha is available in our testing environment which settles to Derive Chain / Protocol testnet. All examples in this documentation refer to the test environment.

---

## post_private-session-keys

**Title:** Post_Private Session Keys
**URL:** https://docs.derive.xyz/reference/post_private-session-keys

ShellNodeRubyPHPPython

Base URL

https://api.lyra.finance/private/session\_keys

Click `Try It!` to start a request and see the response here!

---

## protocol-constants

**Title:** Protocol Constants
**URL:** https://docs.derive.xyz/reference/protocol-constants

## Rollup RPC Node

| Contract | Mainnet Address | Testnet Address |
| --- | --- | --- |
| RPC Endpoint | <https://rpc.lyra.finance> | [<https://rpc-prod-testnet-0eakp60405.t.conduit.xyz>](https://rpc-prod-testnet-0eakp60405.t.conduit.xyz) |
| Block Explorer | <https://explorer.lyra.finance> | [<https://explorer-prod-testnet-0eakp60405.t.conduit.xyz>](https://explorer-prod-testnet-0eakp60405.t.conduit.xyz) |

## Contracts

| Contract | Mainnet Address | Testnet Address |
| --- | --- | --- |
| Matching.sol | 0xeB8d770ec18DB98Db922E9D83260A585b9F0DeAD | 0x3cc154e220c2197c5337b7Bd13363DD127Bc0C6E |
| SubAccount.sol | 0xE7603DF191D699d8BD9891b821347dbAb889E5a5 | 0xb9ed1cc0c50bca7a391a6819e9cAb466f5501d73 |
| ERC20.sol [USDC: 6 decimals] | 0x6879287835A86F50f784313dBEd5E5cCC5bb8481 | 0xe80F2a02398BBf1ab2C9cc52caD1978159c215BD |
| CashAsset.sol | 0x57B03E14d409ADC7fAb6CFc44b5886CAD2D5f02b | 0x6caf294DaC985ff653d5aE75b4FF8E0A66025928 |
| TradeModule.sol | 0xB8D20c2B7a1Ad2EE33Bc50eF10876eD3035b5e7b | 0x87F2863866D85E3192a35A73b388BD625D83f2be |
| TransferModule.sol | 0x01259207A40925b794C8ac320456F7F6c8FE2636 | 0x0CFC1a4a90741aB242cAfaCD798b409E12e68926 |
| DepositModule.sol | 0x9B3FE5E5a3bcEa5df4E08c41Ce89C4e3Ff01Ace3 | 0x43223Db33AdA0575D2E100829543f8B04A37a1ec |
| WithdrawalModule.sol | 0x9d0E8f5b25384C7310CB8C6aE32C8fbeb645d083 | 0xe850641C5207dc5E9423fB15f89ae6031A05fd92 |
| StandardRiskManager.sol | 0x28c9ddF9A3B29c2E6a561c1BC520954e5A33de5D | 0x28bE681F7bEa6f465cbcA1D25A2125fe7533391C |
| RFQ.sol | 0x9371352CCef6f5b36EfDFE90942fFE622Ab77F1D | 0x4E4DD8Be1e461913D9A5DBC4B830e67a8694ebCa |
| LiquidateAddress.sol | 0x66d23e59DaEEF13904eFA2D4B8658aeD05f59a92 | 0x3e2a570B915fEDAFf6176A261d105A4A68a0EA8D |

## ETH Market

| Contract | Mainnet Address | Testnet Address |
| --- | --- | --- |
| BaseAsset.sol |  | 0x41675b7746AE0E464f2594d258CF399c392A179C |
| ERC20.sol [WETH: 18 decimals] |  | 0x3a34565D81156cF0B1b9bC5f14FD00333bcf6B93 |
| Option.sol | 0x4BB4C3CDc7562f08e9910A0C7D8bB7e108861eB4 | 0xBcB494059969DAaB460E0B5d4f5c2366aab79aa1 |
| Perp.sol | 0xAf65752C4643E25C02F693f9D4FE19cF23a095E3 | 0x010e26422790C6Cb3872330980FAa7628FD20294 |
| PortfolioRiskManager.sol | 0xe7cD9370CdE6C9b5eAbCe8f86d01822d3de205A0 | 0xDF448056d7bf3f9Ca13d713114e17f1B7470DeBF |

## BTC Market

| Contract | Mainnet Address | Testnet Address |
| --- | --- | --- |
| BaseAsset.sol |  | 0x0776C34032C618770ca2Be9eD3a11148128b1A50 |
| ERC20.sol [WBTC: 8 decimals] |  | 0xF1493F3602Ab0fC576375a20D7E4B4714DB4422d |
| Option.sol | 0xd0711b9eBE84b778483709CDe62BacFDBAE13623 | 0xAeB81cbe6b19CeEB0dBE0d230CFFE35Bb40a13a7 |
| Perp.sol | 0xDBa83C0C654DB1cd914FA2710bA743e925B53086 | 0xAFB6Bb95cd70D5367e2C39e9dbEb422B9815339D |
| PortfolioRiskManager.sol | 0x45DA02B9cCF384d7DbDD7b2b13e705BADB43Db0D | 0xbaC0328cd4Af53d52F9266Cdbd5bf46720320A20 |

## Constants

| Constant | Mainnet | Testnet |
| --- | --- | --- |
| DOMAIN\_SEPARATOR (Matching.sol) | 0xd96e5f90797da7ec8dc4e276260c7f3f87fedf68775fbe1ef116e996fc60441b | 0x9bcf4dc06df5d8bf23af818d5716491b995020f377d3b7b64c29ed14e3dd1105 |
| ACTION\_TYPEHASH (Matching.sol) | 0x4d7a9f27c403ff9c0f19bce61d76d82f9aa29f8d6d4b0c5474607d9770d1af17 | 0x4d7a9f27c403ff9c0f19bce61d76d82f9aa29f8d6d4b0c5474607d9770d1af17 |
| CHAIN\_ID | 957 | 901 |

---

## rate-limits

**Title:** Rate Limits
**URL:** https://docs.derive.xyz/reference/rate-limits

The below rate limits have been implemented to safeguard our system. Rate limiters use a "fixed window" algorithm to discretely refill the request allowance every 5 seconds.

Market makers are eligible for higher rate limits. To apply for increased rates, please contact our support team.

| Type | matching | per-instrument matching | non-matching | connections | Burst Multiplier |
| --- | --- | --- | --- | --- | --- |
| Trader | 1 TPS | 1 TPS | 5 TPS | 4x per IP | 5x |
| Market Maker | 500+ TPS | 10+ TPS | 500+ TPS | up to 64x per IP | 5x |

> Burst requests for both REST and WebSockets are refreshed every 5 seconds. E.g a trader can send 5x matching requests in a single burst but must wait 5 seconds before any further requests can be sent.

  

## Matching, Non-Matching, and Custom Requests

The below requests are counted as `matching` and `per-instrument matching` requests:

- `private/order`
- `private/replace` (counted as 1 request)
- `private/cancel`
- `private/cancel_by_nonce`
- `private/cancel_by_instrument`
- `private/cancel_by_label`(if `instrument_name` param is set)

`custom` rate-limited requests:

- `private/cancel_all`- 1 TPS
- `private/cancel_by_label` - 10 TPS (if `instrument_name` param is NOT set)

All requests outside of the above are counted as `non-matching`.

  

## REST

All `non-matching` requests over the REST API are rate limited per IP at a flat 10 TPS with 5x burst. If the limit is crossed,`429 Too Many Requests` response is returned.

  

## WebSocket

To access the above rate limits, all clients **must authorize themselves via the `public/login` route**. Otherwise, a reduced rate limit will be applied. Requests exceeding the rate limit will receive the below response:

JSON

```
{
  id: number
  error: {
    code: -32000,
    message: "Rate limit exceeded",
    data: "Retry after ${cooldown} ms"
  }
}

```

Only via WebSockets, the live remaining rate limits can be checked using the `private/getRateLimits` route.

JavaScript

```
// send getRateLimit request
ws.send(JSON.stringify({
  'public/getRateLimits`,
  {},
  id: 1,
});

// example response
{
  "id": 4,
  "result": {
      "remaining_matching": {
        "remainingPoints": 22,
        "msBeforeNext": 4809,
        "consumedPoints": 3
      },
      "remaining_non_matching": {
        "remainingPoints": 98,
        "msBeforeNext": 4809,
        "consumedPoints": 2
      },
      "remaining_connections": {
        "remainingPoints": 29,
        "msBeforeNext": 8881,
        "consumedPoints": 3
      },
      "remaining_per_instrument": {
      	"ETH-PERP": {
          "remainingPoints": 29,
          "msBeforeNext": 381,
          "consumedPoints": 3
        },
        "ETH-08242024-3200-C": {
          "remainingPoints": 29,
          "msBeforeNext": 381,
          "consumedPoints": 10
        }
      }
}  

```

---

## session-keys

**Title:** Use-case
**URL:** https://docs.derive.xyz/reference/session-keys

A session key is simply an Ethereum wallet. Account owners can give other Ethereum wallets temporary access to their accounts via session keys. Any Ethereum wallet can be registered as a session key.

# Use-case

For large accounts, session keys are a useful way to give other users temporary access to:

1. Sign `private/` requests (note: always pass in the "derive wallet" of the account in to `X-LyraWallet` and not the session key).
2. Due to the self-custodial nature of the API, the orderbook cannot force withdrawals, transfers or orders without an explicit user signature. Session Keys (and the account owner) can sign payloads for these sensitive requests (e.g. orders, withdrawals, deposits).
3. Session Keys can only deposit and withdraw funds to the original account owner
4. Session keys cannot be used to bridge funds
5. When using the UX to on-board (see "UX Guides"), session keys are the only way to programmatically trade / manage your account.

For guides on managing session keys, refer to [Onboard via Interface](/docs/onboard-via-interface) and [Onboard Manually](/docs/onboard-manually) guides.

Please refer to the [Derive Python Action Signing SDK](https://pypi.org/project/derive_action_signing/) for actual examples.

# Scopes

When registering a scoped session key, you have the ability to specify a scope for what that session key can access. For now there are three different scopes for session keys.

1. Admin
   1. This scope gives all permissions to all endpoints. Including trading, depositing/withdrawing, signing orders, and any other API on the system. This scope is applied by default to all session keys that are registered via raw transaction in either the public `register_session_key` endpoint, or the private `register_scoped_session_key` endpoint.
2. Account
   1. This scope can set non-order attributes at an account level. For example, this API can toggle `set_cancel_on_disconnect`, cancel orders, send RFQs, or edit session key attributes.
   2. This scope is not able to deposit, withdraw, trade, or call any other endpoint that requires a `signature` parameter.
   3. This scope includes all permissions from `read_only`.
3. Read only
   1. This scope is responsible for viewing orders, account info, or any other kind of private history. This scope can not edit any attributes of an account, or create any orders.

Each private endpoint is required to inform you of the minimum required scope. For example, If an API requires `account` scope. You can call it with your `admin` or `account` level session keys.

---

# Public


## post_public-build-register-session-key-tx

**Title:** Post_Public Build Register Session Key Tx
**URL:** https://docs.derive.xyz/reference/post_public-build-register-session-key-tx

ShellNodeRubyPHPPython

Base URL

https://api.lyra.finance/public/build\_register\_session\_key\_tx

Click `Try It!` to start a request and see the response here!

---

## post_public-create-subaccount-debug

**Title:** Post_Public Create Subaccount Debug
**URL:** https://docs.derive.xyz/reference/post_public-create-subaccount-debug

ShellNodeRubyPHPPython

Base URL

https://api.lyra.finance/public/create\_subaccount\_debug

Click `Try It!` to start a request and see the response here!

---

## post_public-deposit-debug

**Title:** Post_Public Deposit Debug
**URL:** https://docs.derive.xyz/reference/post_public-deposit-debug

ShellNodeRubyPHPPython

Base URL

https://api.lyra.finance/public/deposit\_debug

Click `Try It!` to start a request and see the response here!

---

## post_public-deregister-session-key

**Title:** Post_Public Deregister Session Key
**URL:** https://docs.derive.xyz/reference/post_public-deregister-session-key

ShellNodeRubyPHPPython

Base URL

https://api.lyra.finance/public/deregister\_session\_key

Click `Try It!` to start a request and see the response here!

---

## post_public-execute-quote-debug

**Title:** Post_Public Execute Quote Debug
**URL:** https://docs.derive.xyz/reference/post_public-execute-quote-debug

ShellNodeRubyPHPPython

Base URL

https://api.lyra.finance/public/execute\_quote\_debug

Click `Try It!` to start a request and see the response here!

---

## post_public-get-all-currencies

**Title:** Post_Public Get All Currencies
**URL:** https://docs.derive.xyz/reference/post_public-get-all-currencies

ShellNodeRubyPHPPython

Base URL

https://api.lyra.finance/public/get\_all\_currencies

Click `Try It!` to start a request and see the response here!

---

## post_public-get-all-instruments

**Title:** Post_Public Get All Instruments
**URL:** https://docs.derive.xyz/reference/post_public-get-all-instruments

ShellNodeRubyPHPPython

Base URL

https://api.lyra.finance/public/get\_all\_instruments

Click `Try It!` to start a request and see the response here!

---

## post_public-get-all-points

**Title:** Post_Public Get All Points
**URL:** https://docs.derive.xyz/reference/post_public-get-all-points

ShellNodeRubyPHPPython

Base URL

https://api.lyra.finance/public/get\_all\_points

Click `Try It!` to start a request and see the response here!

---

## post_public-get-currency

**Title:** Post_Public Get Currency
**URL:** https://docs.derive.xyz/reference/post_public-get-currency

ShellNodeRubyPHPPython

Base URL

https://api.lyra.finance/public/get\_currency

Click `Try It!` to start a request and see the response here!

---

## post_public-get-descendant-tree

**Title:** Post_Public Get Descendant Tree
**URL:** https://docs.derive.xyz/reference/post_public-get-descendant-tree

ShellNodeRubyPHPPython

Base URL

https://api.lyra.finance/public/get\_descendant\_tree

Click `Try It!` to start a request and see the response here!

---

## post_public-get-funding-rate-history

**Title:** Post_Public Get Funding Rate History
**URL:** https://docs.derive.xyz/reference/post_public-get-funding-rate-history

ShellNodeRubyPHPPython

Base URL

https://api.lyra.finance/public/get\_funding\_rate\_history

Click `Try It!` to start a request and see the response here!

---

## post_public-get-instrument

**Title:** Post_Public Get Instrument
**URL:** https://docs.derive.xyz/reference/post_public-get-instrument

ShellNodeRubyPHPPython

Base URL

https://api.lyra.finance/public/get\_instrument

Click `Try It!` to start a request and see the response here!

---

## post_public-get-instruments

**Title:** Post_Public Get Instruments
**URL:** https://docs.derive.xyz/reference/post_public-get-instruments

ShellNodeRubyPHPPython

Base URL

https://api.lyra.finance/public/get\_instruments

Click `Try It!` to start a request and see the response here!

---

## post_public-get-interest-rate-history

**Title:** Post_Public Get Interest Rate History
**URL:** https://docs.derive.xyz/reference/post_public-get-interest-rate-history

ShellNodeRubyPHPPython

Base URL

https://api.lyra.finance/public/get\_interest\_rate\_history

Click `Try It!` to start a request and see the response here!

---

## post_public-get-invite-code

**Title:** Post_Public Get Invite Code
**URL:** https://docs.derive.xyz/reference/post_public-get-invite-code

ShellNodeRubyPHPPython

Base URL

https://api.lyra.finance/public/get\_invite\_code

Click `Try It!` to start a request and see the response here!

---

## post_public-get-latest-signed-feeds

**Title:** Post_Public Get Latest Signed Feeds
**URL:** https://docs.derive.xyz/reference/post_public-get-latest-signed-feeds

ShellNodeRubyPHPPython

Base URL

https://api.lyra.finance/public/get\_latest\_signed\_feeds

Click `Try It!` to start a request and see the response here!

---

## post_public-get-liquidation-history

**Title:** Post_Public Get Liquidation History
**URL:** https://docs.derive.xyz/reference/post_public-get-liquidation-history

ShellNodeRubyPHPPython

Base URL

https://api.lyra.finance/public/get\_liquidation\_history

Click `Try It!` to start a request and see the response here!

---

## post_public-get-live-incidents

**Title:** Post_Public Get Live Incidents
**URL:** https://docs.derive.xyz/reference/post_public-get-live-incidents

ShellNodeRubyPHPPython

Base URL

https://api.lyra.finance/public/get\_live\_incidents

Click `Try It!` to start a request and see the response here!

---

## post_public-get-maker-program-scores

**Title:** Post_Public Get Maker Program Scores
**URL:** https://docs.derive.xyz/reference/post_public-get-maker-program-scores

ShellNodeRubyPHPPython

Base URL

https://api.lyra.finance/public/get\_maker\_program\_scores

Click `Try It!` to start a request and see the response here!

---

## post_public-get-maker-programs

**Title:** Post_Public Get Maker Programs
**URL:** https://docs.derive.xyz/reference/post_public-get-maker-programs

ShellNodeRubyPHPPython

Base URL

https://api.lyra.finance/public/get\_maker\_programs

Click `Try It!` to start a request and see the response here!

---

## post_public-get-margin

**Title:** Post_Public Get Margin
**URL:** https://docs.derive.xyz/reference/post_public-get-margin

ShellNodeRubyPHPPython

Base URL

https://api.lyra.finance/public/get\_margin

Click `Try It!` to start a request and see the response here!

---

## post_public-get-option-settlement-history

**Title:** Post_Public Get Option Settlement History
**URL:** https://docs.derive.xyz/reference/post_public-get-option-settlement-history

ShellNodeRubyPHPPython

Base URL

https://api.lyra.finance/public/get\_option\_settlement\_history

Click `Try It!` to start a request and see the response here!

---

## post_public-get-option-settlement-prices

**Title:** Post_Public Get Option Settlement Prices
**URL:** https://docs.derive.xyz/reference/post_public-get-option-settlement-prices

ShellNodeRubyPHPPython

Base URL

https://api.lyra.finance/public/get\_option\_settlement\_prices

Click `Try It!` to start a request and see the response here!

---

## post_public-get-points

**Title:** Post_Public Get Points
**URL:** https://docs.derive.xyz/reference/post_public-get-points

ShellNodeRubyPHPPython

Base URL

https://api.lyra.finance/public/get\_points

Click `Try It!` to start a request and see the response here!

---

## post_public-get-points-leaderboard

**Title:** Post_Public Get Points Leaderboard
**URL:** https://docs.derive.xyz/reference/post_public-get-points-leaderboard

ShellNodeRubyPHPPython

Base URL

https://api.lyra.finance/public/get\_points\_leaderboard

Click `Try It!` to start a request and see the response here!

---

## post_public-get-spot-feed-history

**Title:** Post_Public Get Spot Feed History
**URL:** https://docs.derive.xyz/reference/post_public-get-spot-feed-history

ShellNodeRubyPHPPython

Base URL

https://api.lyra.finance/public/get\_spot\_feed\_history

Click `Try It!` to start a request and see the response here!

---

## post_public-get-spot-feed-history-candles

**Title:** Post_Public Get Spot Feed History Candles
**URL:** https://docs.derive.xyz/reference/post_public-get-spot-feed-history-candles

ShellNodeRubyPHPPython

Base URL

https://api.lyra.finance/public/get\_spot\_feed\_history\_candles

Click `Try It!` to start a request and see the response here!

---

## post_public-get-ticker

**Title:** Post_Public Get Ticker
**URL:** https://docs.derive.xyz/reference/post_public-get-ticker

ShellNodeRubyPHPPython

Base URL

https://api.lyra.finance/public/get\_ticker

Click `Try It!` to start a request and see the response here!

---

## post_public-get-time

**Title:** Post_Public Get Time
**URL:** https://docs.derive.xyz/reference/post_public-get-time

ShellNodeRubyPHPPython

Base URL

https://api.lyra.finance/public/get\_time

Click `Try It!` to start a request and see the response here!

---

## post_public-get-trade-history

**Title:** Post_Public Get Trade History
**URL:** https://docs.derive.xyz/reference/post_public-get-trade-history

ShellNodeRubyPHPPython

Base URL

https://api.lyra.finance/public/get\_trade\_history

Click `Try It!` to start a request and see the response here!

---

## post_public-get-transaction

**Title:** Post_Public Get Transaction
**URL:** https://docs.derive.xyz/reference/post_public-get-transaction

ShellNodeRubyPHPPython

Base URL

https://api.lyra.finance/public/get\_transaction

Click `Try It!` to start a request and see the response here!

---

## post_public-get-tree-roots

**Title:** Post_Public Get Tree Roots
**URL:** https://docs.derive.xyz/reference/post_public-get-tree-roots

ShellNodeRubyPHPPython

Base URL

https://api.lyra.finance/public/get\_tree\_roots

Click `Try It!` to start a request and see the response here!

---

## post_public-get-vault-balances

**Title:** Post_Public Get Vault Balances
**URL:** https://docs.derive.xyz/reference/post_public-get-vault-balances

ShellNodeRubyPHPPython

Base URL

https://api.lyra.finance/public/get\_vault\_balances

Click `Try It!` to start a request and see the response here!

---

## post_public-get-vault-share

**Title:** Post_Public Get Vault Share
**URL:** https://docs.derive.xyz/reference/post_public-get-vault-share

ShellNodeRubyPHPPython

Base URL

https://api.lyra.finance/public/get\_vault\_share

Click `Try It!` to start a request and see the response here!

---

## post_public-get-vault-statistics

**Title:** Post_Public Get Vault Statistics
**URL:** https://docs.derive.xyz/reference/post_public-get-vault-statistics

ShellNodeRubyPHPPython

Base URL

https://api.lyra.finance/public/get\_vault\_statistics

Click `Try It!` to start a request and see the response here!

---

## post_public-login

**Title:** Post_Public Login
**URL:** https://docs.derive.xyz/reference/post_public-login

ShellNodeRubyPHPPython

Base URL

https://api.lyra.finance/public/login

Click `Try It!` to start a request and see the response here!

---

## post_public-margin-watch

**Title:** Post_Public Margin Watch
**URL:** https://docs.derive.xyz/reference/post_public-margin-watch

ShellNodeRubyPHPPython

Base URL

https://api.lyra.finance/public/margin\_watch

Click `Try It!` to start a request and see the response here!

---

## post_public-register-session-key

**Title:** Post_Public Register Session Key
**URL:** https://docs.derive.xyz/reference/post_public-register-session-key

ShellNodeRubyPHPPython

Base URL

https://api.lyra.finance/public/register\_session\_key

Click `Try It!` to start a request and see the response here!

---

## post_public-send-quote-debug

**Title:** Post_Public Send Quote Debug
**URL:** https://docs.derive.xyz/reference/post_public-send-quote-debug

ShellNodeRubyPHPPython

Base URL

https://api.lyra.finance/public/send\_quote\_debug

Click `Try It!` to start a request and see the response here!

---

## post_public-statistics

**Title:** Post_Public Statistics
**URL:** https://docs.derive.xyz/reference/post_public-statistics

ShellNodeRubyPHPPython

Base URL

https://api.lyra.finance/public/statistics

Click `Try It!` to start a request and see the response here!

---

## post_public-validate-invite-code

**Title:** Post_Public Validate Invite Code
**URL:** https://docs.derive.xyz/reference/post_public-validate-invite-code

ShellNodeRubyPHPPython

Base URL

https://api.lyra.finance/public/validate\_invite\_code

Click `Try It!` to start a request and see the response here!

---

## post_public-withdraw-debug

**Title:** Post_Public Withdraw Debug
**URL:** https://docs.derive.xyz/reference/post_public-withdraw-debug

ShellNodeRubyPHPPython

Base URL

https://api.lyra.finance/public/withdraw\_debug

Click `Try It!` to start a request and see the response here!

---

## public-build_register_session_key_tx

**Title:** Public Build_Register_Session_Key_Tx
**URL:** https://docs.derive.xyz/reference/public-build_register_session_key_tx

### Method Name

#### `public/build_register_session_key_tx`

Build a signable transaction params dictionary.

### Parameters

|  |
| --- |
| **expiry\_sec**integer required  Expiry of the session key |
| **gas**integer required  Gas allowance for transaction. If none, will use estimateGas \* 150% |
| **nonce**integer required  Wallet's transaction count, If none, will use eth.getTransactionCount() |
| **public\_session\_key**string required  Session key in the form of an Ethereum EOA |
| **wallet**string required  Ethereum wallet address of account |

### Response

|  |
| --- |
| **id**string or integer required |
| **result**object required |
| result.**tx\_params**object required  Transaction params in dictionary form, same as `TxParams` in `web3.py` |

### Example

ShellJavaScriptPython

```
{request_example_shell}

```

```
{request_example_javascript}

```

```
{request_example_python}

```

> The above command returns JSON structured like this:

JSON

```
{response_example_json}

```

---

## public-create_subaccount_debug

**Title:** Public Create_Subaccount_Debug
**URL:** https://docs.derive.xyz/reference/public-create_subaccount_debug

### Method Name

#### `public/create_subaccount_debug`

Used for debugging only, do not use in production. Will return the incremental encoded and hashed data.  
  
See guides in Documentation for more.

### Parameters

|  |
| --- |
| **amount**string required  Amount of the asset to deposit |
| **asset\_name**string required  Name of asset to deposit |
| **margin\_type**string required  `PM` (Portfolio Margin) or `SM` (Standard Margin) enum  `PM` `SM` `PM2` |
| **nonce**integer required  Unique nonce defined as  (e.g. 1695836058725001, where 001 is the random number) |
| **signature\_expiry\_sec**integer required  Unix timestamp in seconds. Expiry MUST be >5min from now |
| **signer**string required  Ethereum wallet address that is signing the deposit |
| **wallet**string required  Ethereum wallet address |
| **currency**string  Base currency of the subaccount (only for `PM`) |

### Response

|  |
| --- |
| **id**string or integer required |
| **result**object required |
| result.**action\_hash**string required  Keccak hashed action data |
| result.**encoded\_data**string required  ABI encoded deposit data |
| result.**encoded\_data\_hashed**string required  Keccak hashed encoded\_data |
| result.**typed\_data\_hash**string required  EIP 712 typed data hash |

### Example

ShellJavaScriptPython

```
{request_example_shell}

```

```
{request_example_javascript}

```

```
{request_example_python}

```

> The above command returns JSON structured like this:

JSON

```
{response_example_json}

```

---

## public-deposit_debug

**Title:** Public Deposit_Debug
**URL:** https://docs.derive.xyz/reference/public-deposit_debug

### Method Name

#### `public/deposit_debug`

Used for debugging only, do not use in production. Will return the incremental encoded and hashed data.  
  
See guides in Documentation for more.

### Parameters

|  |
| --- |
| **amount**string required  Amount of the asset to deposit |
| **asset\_name**string required  Name of asset to deposit |
| **nonce**integer required  Unique nonce defined as  (e.g. 1695836058725001, where 001 is the random number) |
| **signature\_expiry\_sec**integer required  Unix timestamp in seconds. Expiry MUST be >5min from now |
| **signer**string required  Ethereum wallet address that is signing the deposit |
| **subaccount\_id**integer required  Subaccount\_id |
| **is\_atomic\_signing**boolean  Used by vaults to determine whether the signature is an EIP-1271 signature |

### Response

|  |
| --- |
| **id**string or integer required |
| **result**object required |
| result.**action\_hash**string required  Keccak hashed action data |
| result.**encoded\_data**string required  ABI encoded deposit data |
| result.**encoded\_data\_hashed**string required  Keccak hashed encoded\_data |
| result.**typed\_data\_hash**string required  EIP 712 typed data hash |

### Example

ShellJavaScriptPython

```
{request_example_shell}

```

```
{request_example_javascript}

```

```
{request_example_python}

```

> The above command returns JSON structured like this:

JSON

```
{response_example_json}

```

---

## public-deregister_session_key

**Title:** Public Deregister_Session_Key
**URL:** https://docs.derive.xyz/reference/public-deregister_session_key

### Method Name

#### `public/deregister_session_key`

Used for de-registering admin scoped keys. For other scopes, use `/edit_session_key`.

### Parameters

|  |
| --- |
| **public\_session\_key**string required  Session key in the form of an Ethereum EOA |
| **signed\_raw\_tx**string required  A signed RLP encoded ETH transaction in form of a hex string (same as `w3.eth.account.sign_transaction(unsigned_tx, private_key).rawTransaction.hex()`) |
| **wallet**string required  Ethereum wallet address of account |

### Response

|  |
| --- |
| **id**string or integer required |
| **result**object required |
| result.**public\_session\_key**string required  Session key in the form of an Ethereum EOA |
| result.**transaction\_id**string required  ID to lookup status of transaction |

### Example

ShellJavaScriptPython

```
{request_example_shell}

```

```
{request_example_javascript}

```

```
{request_example_python}

```

> The above command returns JSON structured like this:

JSON

```
{response_example_json}

```

---

## public-execute_quote_debug

**Title:** Public Execute_Quote_Debug
**URL:** https://docs.derive.xyz/reference/public-execute_quote_debug

### Method Name

#### `public/execute_quote_debug`

Sends a quote in response to an RFQ request.  
The legs supplied in the parameters must exactly match those in the RFQ.

### Parameters

|  |
| --- |
| **direction**string required  Quote direction, `buy` means trading each leg at its direction, `sell` means trading each leg in the opposite direction. enum  `buy` `sell` |
| **max\_fee**string required  Max fee ($ for the full trade). Request will be rejected if the supplied max fee is below the estimated fee for this trade. |
| **nonce**integer required  Unique nonce defined as a concatenated `UTC timestamp in ms` and `random number up to 6 digits` (e.g. 1695836058725001, where 001 is the random number) |
| **quote\_id**string required  Quote ID to execute against |
| **rfq\_id**string required  RFQ ID to execute (must be sent by `subaccount_id`) |
| **signature**string required  Ethereum signature of the quote |
| **signature\_expiry\_sec**integer required  Unix timestamp in seconds. Expiry MUST be at least 310 seconds from now. Once time till signature expiry reaches 300 seconds, the quote will be considered expired. This buffer is meant to ensure the trade can settle on chain in case of a blockchain congestion. |
| **signer**string required  Owner wallet address or registered session key that signed the quote |
| **subaccount\_id**integer required  Subaccount ID |
| **legs**array of objects required  Quote legs |
| legs[].**amount**string required  Amount in units of the base |
| legs[].**direction**string required  Leg direction enum  `buy` `sell` |
| legs[].**instrument\_name**string required  Instrument name |
| legs[].**price**string required  Leg price |
|  |
| **label**string  Optional user-defined label for the quote |

### Response

|  |
| --- |
| **id**string or integer required |
| **result**object required |
| result.**action\_hash**string required  Keccak hashed action data |
| result.**encoded\_data**string required  ABI encoded deposit data |
| result.**encoded\_data\_hashed**string required  Keccak hashed encoded\_data |
| result.**encoded\_legs**string required  ABI encoded legs data |
| result.**legs\_hash**string required  Keccak hashed legs data |
| result.**typed\_data\_hash**string required  EIP 712 typed data hash |

### Example

ShellJavaScriptPython

```
{request_example_shell}

```

```
{request_example_javascript}

```

```
{request_example_python}

```

> The above command returns JSON structured like this:

JSON

```
{response_example_json}

```

---

## public-get_all_currencies

**Title:** Public Get_All_Currencies
**URL:** https://docs.derive.xyz/reference/public-get_all_currencies

### Method Name

#### `public/get_all_currencies`

Get all active currencies with their spot price, spot price 24hrs ago.  
  
For real-time updates, recommend using channels -> ticker or orderbook.

### Parameters

|  |
| --- |

### Response

|  |
| --- |
| **id**string or integer required |
| **result**array of objects required |
| result[].**borrow\_apy**string required  Borrow APY (only for USDC) |
| result[].**currency**string required  Underlying currency of asset (`ETH`, `BTC`, etc) |
| result[].**market\_type**string required  Market type of the currency enum  `ALL` `SRM_BASE_ONLY` `SRM_OPTION_ONLY` `SRM_PERP_ONLY` `CASH` |
| result[].**spot\_price**string required  Spot price of the currency |
| result[].**srm\_im\_discount**string required  Initial Margin discount for given collateral in Standard Manager (e.g. LTV). Only the Standard Manager supports non-USDC collateral |
| result[].**srm\_mm\_discount**string required  Maintenance Margin discount for given collateral in Standard Manager (e.g. liquidation threshold). Only the Standard Manager supports non-USDC collateral |
| result[].**supply\_apy**string required  Supply APY (only for USDC) |
| result[].**total\_borrow**string required  Total collateral borrowed in the protocol (only USDC is borrowable) |
| result[].**total\_supply**string required  Total collateral supplied in the protocol |
| result[].**asset\_cap\_and\_supply\_per\_manager**object required  Current open interest and open interest cap by manager and asset type |
| result[].**instrument\_types**array of strings required  Instrument types supported for the currency |
| result[].**managers**array of objects required  Managers supported for the currency |
| result[].managers[].**address**string required  Address of the manager |
| result[].managers[].**margin\_type**string required  Margin type of the manager enum  `PM` `SM` `PM2` |
| result[].managers[].**currency**string or null  Currency of the manager (only applies to portfolio managers) |
|  |
| result[].**pm2\_collateral\_discounts**array of objects required  Initial and Maintenance Margin discounts for given collateral in PM2 |
| result[].pm2\_collateral\_discounts[].**im\_discount**string required  Initial Margin discount for given collateral in PM2 |
| result[].pm2\_collateral\_discounts[].**manager\_currency**string required  Currency of the manager |
| result[].pm2\_collateral\_discounts[].**mm\_discount**string required  Maintenance Margin discount for given collateral in PM2 |
|  |
| result[].**protocol\_asset\_addresses**object required  Asset addressses of the derive protocol for given currency |
| result[].protocol\_asset\_addresses.**option**string or null  Address of the Derive protocol option contract (none if not supported) |
| result[].protocol\_asset\_addresses.**perp**string or null  Address of the Derive protocol perp contract (none if not supported) |
| result[].protocol\_asset\_addresses.**spot**string or null  Address of the Derive protocol spot contract (none if not supported) |
| result[].protocol\_asset\_addresses.**underlying\_erc20**string or null  Address of the erc20 asset on Derive chain. This is the asset that is deposited into the spot asset |
|  |
| result[].**spot\_price\_24h**string or null  Spot price of the currency 24 hours ago |

### Example

ShellJavaScriptPython

```
{request_example_shell}

```

```
{request_example_javascript}

```

```
{request_example_python}

```

> The above command returns JSON structured like this:

JSON

```
{response_example_json}

```

---

## public-get_all_instruments

**Title:** Public Get_All_Instruments
**URL:** https://docs.derive.xyz/reference/public-get_all_instruments

### Method Name

#### `public/get_all_instruments`

Get a paginated history of all instruments

### Parameters

|  |
| --- |
| **expired**boolean required  If `True`: include expired instruments. |
| **instrument\_type**string required  `erc20`, `option`, or `perp` enum  `erc20` `option` `perp` |
| **currency**string  Underlying currency of asset (`ETH`, `BTC`, etc) |
| **page**integer  Page number of results to return (default 1, returns last if above `num_pages`) |
| **page\_size**integer  Number of results per page (default 100, max 1000) |

### Response

|  |
| --- |
| **id**string or integer required |
| **result**object required |
| result.**instruments**array of objects required  List of instruments |
| result.instruments[].**amount\_step**string required  Minimum valid increment of order amount |
| result.instruments[].**base\_asset\_address**string required  Blockchain address of the base asset |
| result.instruments[].**base\_asset\_sub\_id**string required  Sub ID of the specific base asset as defined in Asset.sol |
| result.instruments[].**base\_currency**string required  Underlying currency of base asset (`ETH`, `BTC`, etc) |
| result.instruments[].**base\_fee**string required  $ base fee added to every taker order |
| result.instruments[].**erc20\_details**object or null required  Details of the erc20 asset (if applicable) |
| result.instruments[].erc20\_details.**decimals**integer required  Number of decimals of the underlying on-chain ERC20 token |
| result.instruments[].erc20\_details.**borrow\_index**string  Latest borrow index as per `CashAsset.sol` implementation |
| result.instruments[].erc20\_details.**supply\_index**string  Latest supply index as per `CashAsset.sol` implementation |
| result.instruments[].erc20\_details.**underlying\_erc20\_address**string  Address of underlying on-chain ERC20 (not V2 asset) |
|  |
| result.instruments[].**fifo\_min\_allocation**string required  Minimum number of contracts that get filled using FIFO. Actual number of contracts that gets filled by FIFO will be the max between this value and (1 - pro\_rata\_fraction) x order\_amount, plus any size leftovers due to rounding. |
| result.instruments[].**instrument\_name**string required  Instrument name |
| result.instruments[].**instrument\_type**string required  `erc20`, `option`, or `perp` enum  `erc20` `option` `perp` |
| result.instruments[].**is\_active**boolean required  If `True`: instrument is tradeable within `activation` and `deactivation` timestamps |
| result.instruments[].**maker\_fee\_rate**string required  Percent of spot price fee rate for makers |
| result.instruments[].**maximum\_amount**string required  Maximum valid amount of contracts / tokens per trade |
| result.instruments[].**minimum\_amount**string required  Minimum valid amount of contracts / tokens per trade |
| result.instruments[].**option\_details**object or null required  Details of the option asset (if applicable) |
| result.instruments[].option\_details.**expiry**integer required  Unix timestamp of expiry date (in seconds) |
| result.instruments[].option\_details.**index**string required  Underlying settlement price index |
| result.instruments[].option\_details.**option\_type**string required   enum  `C` `P` |
| result.instruments[].option\_details.**strike**string required |
| result.instruments[].option\_details.**settlement\_price**string or null  Settlement price of the option |
|  |
| result.instruments[].**perp\_details**object or null required  Details of the perp asset (if applicable) |
| result.instruments[].perp\_details.**aggregate\_funding**string required  Latest aggregated funding as per `PerpAsset.sol` |
| result.instruments[].perp\_details.**funding\_rate**string required  Current hourly funding rate as per `PerpAsset.sol` |
| result.instruments[].perp\_details.**index**string required  Underlying spot price index for funding rate |
| result.instruments[].perp\_details.**max\_rate\_per\_hour**string required  Max rate per hour as per `PerpAsset.sol` |
| result.instruments[].perp\_details.**min\_rate\_per\_hour**string required  Min rate per hour as per `PerpAsset.sol` |
| result.instruments[].perp\_details.**static\_interest\_rate**string required  Static interest rate as per `PerpAsset.sol` |
|  |
| result.instruments[].**pro\_rata\_amount\_step**string required  Pro-rata fill share of every order is rounded down to be a multiple of this number. Leftovers of the order due to rounding are filled FIFO. |
| result.instruments[].**pro\_rata\_fraction**string required  Fraction of order that gets filled using pro-rata matching. If zero, the matching is full FIFO. |
| result.instruments[].**quote\_currency**string required  Quote currency (`USD` for perps, `USDC` for options) |
| result.instruments[].**scheduled\_activation**integer required  Timestamp at which became or will become active (if applicable) |
| result.instruments[].**scheduled\_deactivation**integer required  Scheduled deactivation time for instrument (if applicable) |
| result.instruments[].**taker\_fee\_rate**string required  Percent of spot price fee rate for takers |
| result.instruments[].**tick\_size**string required  Tick size of the instrument, i.e. minimum price increment |
| result.instruments[].**mark\_price\_fee\_rate\_cap**string or null  Percent of option price fee cap, e.g. 12.5%, null if not applicable |
|  |
| result.**pagination**object required  Pagination info |
| result.pagination.**count**integer required  Total number of items, across all pages |
| result.pagination.**num\_pages**integer required  Number of pages |

### Example

ShellJavaScriptPython

```
{request_example_shell}

```

```
{request_example_javascript}

```

```
{request_example_python}

```

> The above command returns JSON structured like this:

JSON

```
{response_example_json}

```

---

## public-get_all_points

**Title:** Public Get_All_Points
**URL:** https://docs.derive.xyz/reference/public-get_all_points

### Method Name

#### `public/get_all_points`

Get all points for all users for a given program. This request is cached in WSGI.

### Parameters

|  |
| --- |
| **program**string required  Program for which to count points |

### Response

|  |
| --- |
| **id**string or integer required |
| **result**object required |
| result.**total\_notional\_volume**string required  Total $ notional volume traded by all users elligible for points |
| result.**total\_users**integer required  Total number of users in the program |
| result.**points**object required  Points earned per category |

### Example

ShellJavaScriptPython

```
{request_example_shell}

```

```
{request_example_javascript}

```

```
{request_example_python}

```

> The above command returns JSON structured like this:

JSON

```
{response_example_json}

```

---

## public-get_currency

**Title:** Public Get_Currency
**URL:** https://docs.derive.xyz/reference/public-get_currency

### Method Name

#### `public/get_currency`

Get currency related risk params, spot price 24hrs ago and lending details for a specific currency.

### Parameters

|  |
| --- |
| **currency**string required  Underlying currency of asset (`ETH`, `BTC`, etc) |

### Response

|  |
| --- |
| **id**string or integer required |
| **result**object required |
| result.**borrow\_apy**string required  Borrow APY (only for USDC) |
| result.**currency**string required  Underlying currency of asset (`ETH`, `BTC`, etc) |
| result.**market\_type**string required  Market type of the currency enum  `ALL` `SRM_BASE_ONLY` `SRM_OPTION_ONLY` `SRM_PERP_ONLY` `CASH` |
| result.**spot\_price**string required  Spot price of the currency |
| result.**srm\_im\_discount**string required  Initial Margin discount for given collateral in Standard Manager (e.g. LTV). Only the Standard Manager supports non-USDC collateral |
| result.**srm\_mm\_discount**string required  Maintenance Margin discount for given collateral in Standard Manager (e.g. liquidation threshold). Only the Standard Manager supports non-USDC collateral |
| result.**supply\_apy**string required  Supply APY (only for USDC) |
| result.**total\_borrow**string required  Total collateral borrowed in the protocol (only USDC is borrowable) |
| result.**total\_supply**string required  Total collateral supplied in the protocol |
| result.**asset\_cap\_and\_supply\_per\_manager**object required  Current open interest and open interest cap by manager and asset type |
| result.**instrument\_types**array of strings required  Instrument types supported for the currency |
| result.**managers**array of objects required  Managers supported for the currency |
| result.managers[].**address**string required  Address of the manager |
| result.managers[].**margin\_type**string required  Margin type of the manager enum  `PM` `SM` `PM2` |
| result.managers[].**currency**string or null  Currency of the manager (only applies to portfolio managers) |
|  |
| result.**pm2\_collateral\_discounts**array of objects required  Initial and Maintenance Margin discounts for given collateral in PM2 |
| result.pm2\_collateral\_discounts[].**im\_discount**string required  Initial Margin discount for given collateral in PM2 |
| result.pm2\_collateral\_discounts[].**manager\_currency**string required  Currency of the manager |
| result.pm2\_collateral\_discounts[].**mm\_discount**string required  Maintenance Margin discount for given collateral in PM2 |
|  |
| result.**protocol\_asset\_addresses**object required  Asset addressses of the derive protocol for given currency |
| result.protocol\_asset\_addresses.**option**string or null  Address of the Derive protocol option contract (none if not supported) |
| result.protocol\_asset\_addresses.**perp**string or null  Address of the Derive protocol perp contract (none if not supported) |
| result.protocol\_asset\_addresses.**spot**string or null  Address of the Derive protocol spot contract (none if not supported) |
| result.protocol\_asset\_addresses.**underlying\_erc20**string or null  Address of the erc20 asset on Derive chain. This is the asset that is deposited into the spot asset |
|  |
| result.**spot\_price\_24h**string or null  Spot price of the currency 24 hours ago |

### Example

ShellJavaScriptPython

```
{request_example_shell}

```

```
{request_example_javascript}

```

```
{request_example_python}

```

> The above command returns JSON structured like this:

JSON

```
{response_example_json}

```

---

## public-get_descendant_tree

**Title:** Public Get_Descendant_Tree
**URL:** https://docs.derive.xyz/reference/public-get_descendant_tree

### Method Name

#### `public/get_descendant_tree`

Returns the tree of descendants given a root wallet up to 2 levels deep.

### Parameters

|  |
| --- |
| **wallet\_or\_invite\_code**string or integer required  Wallet of account owner to get descendants for. |

### Response

|  |
| --- |
| **id**string or integer required |
| **result**object required |
| result.**parent**string required  Wallet address of wallet that referred this account |
| result.**descendants**object required  Nested dict representing a tree of descendants two layers deep |

### Example

ShellJavaScriptPython

```
{request_example_shell}

```

```
{request_example_javascript}

```

```
{request_example_python}

```

> The above command returns JSON structured like this:

JSON

```
{response_example_json}

```

---

## public-get_funding_rate_history

**Title:** Public Get_Funding_Rate_History
**URL:** https://docs.derive.xyz/reference/public-get_funding_rate_history

### Method Name

#### `public/get_funding_rate_history`

Get funding rate history. Start timestamp is restricted to at most 30 days ago.  
End timestamp greater than current time will be truncated to current time.  
Zero start timestamp is allowed and will default to 30 days from the end timestamp.  
  
DB: read replica

### Parameters

|  |
| --- |
| **instrument\_name**string required  Instrument name to return funding history for |
| **end\_timestamp**integer  End timestamp of the event history (default current time) |
| **period**string  Period of the funding rate enum  `900` `3600` `14400` `28800` `86400` |
| **start\_timestamp**integer  Start timestamp of the event history (default 0) |

### Response

|  |
| --- |
| **id**string or integer required |
| **result**object required |
| result.**funding\_rate\_history**array of objects required  List of funding rates |
| result.funding\_rate\_history[].**funding\_rate**string required  Hourly funding rate value at the event time |
| result.funding\_rate\_history[].**timestamp**integer required  Timestamp of the funding rate update (in ms since UNIX epoch) |

### Example

ShellJavaScriptPython

```
{request_example_shell}

```

```
{request_example_javascript}

```

```
{request_example_python}

```

> The above command returns JSON structured like this:

JSON

```
{response_example_json}

```

---

## public-get_instrument

**Title:** Public Get_Instrument
**URL:** https://docs.derive.xyz/reference/public-get_instrument

### Method Name

#### `public/get_instrument`

Get single instrument by asset name

### Parameters

|  |
| --- |
| **instrument\_name**string required  Instrument name |

### Response

|  |
| --- |
| **id**string or integer required |
| **result**object required |
| result.**amount\_step**string required  Minimum valid increment of order amount |
| result.**base\_asset\_address**string required  Blockchain address of the base asset |
| result.**base\_asset\_sub\_id**string required  Sub ID of the specific base asset as defined in Asset.sol |
| result.**base\_currency**string required  Underlying currency of base asset (`ETH`, `BTC`, etc) |
| result.**base\_fee**string required  $ base fee added to every taker order |
| result.**erc20\_details**object or null required  Details of the erc20 asset (if applicable) |
| result.erc20\_details.**decimals**integer required  Number of decimals of the underlying on-chain ERC20 token |
| result.erc20\_details.**borrow\_index**string  Latest borrow index as per `CashAsset.sol` implementation |
| result.erc20\_details.**supply\_index**string  Latest supply index as per `CashAsset.sol` implementation |
| result.erc20\_details.**underlying\_erc20\_address**string  Address of underlying on-chain ERC20 (not V2 asset) |
|  |
| result.**fifo\_min\_allocation**string required  Minimum number of contracts that get filled using FIFO. Actual number of contracts that gets filled by FIFO will be the max between this value and (1 - pro\_rata\_fraction) x order\_amount, plus any size leftovers due to rounding. |
| result.**instrument\_name**string required  Instrument name |
| result.**instrument\_type**string required  `erc20`, `option`, or `perp` enum  `erc20` `option` `perp` |
| result.**is\_active**boolean required  If `True`: instrument is tradeable within `activation` and `deactivation` timestamps |
| result.**maker\_fee\_rate**string required  Percent of spot price fee rate for makers |
| result.**maximum\_amount**string required  Maximum valid amount of contracts / tokens per trade |
| result.**minimum\_amount**string required  Minimum valid amount of contracts / tokens per trade |
| result.**option\_details**object or null required  Details of the option asset (if applicable) |
| result.option\_details.**expiry**integer required  Unix timestamp of expiry date (in seconds) |
| result.option\_details.**index**string required  Underlying settlement price index |
| result.option\_details.**option\_type**string required   enum  `C` `P` |
| result.option\_details.**strike**string required |
| result.option\_details.**settlement\_price**string or null  Settlement price of the option |
|  |
| result.**perp\_details**object or null required  Details of the perp asset (if applicable) |
| result.perp\_details.**aggregate\_funding**string required  Latest aggregated funding as per `PerpAsset.sol` |
| result.perp\_details.**funding\_rate**string required  Current hourly funding rate as per `PerpAsset.sol` |
| result.perp\_details.**index**string required  Underlying spot price index for funding rate |
| result.perp\_details.**max\_rate\_per\_hour**string required  Max rate per hour as per `PerpAsset.sol` |
| result.perp\_details.**min\_rate\_per\_hour**string required  Min rate per hour as per `PerpAsset.sol` |
| result.perp\_details.**static\_interest\_rate**string required  Static interest rate as per `PerpAsset.sol` |
|  |
| result.**pro\_rata\_amount\_step**string required  Pro-rata fill share of every order is rounded down to be a multiple of this number. Leftovers of the order due to rounding are filled FIFO. |
| result.**pro\_rata\_fraction**string required  Fraction of order that gets filled using pro-rata matching. If zero, the matching is full FIFO. |
| result.**quote\_currency**string required  Quote currency (`USD` for perps, `USDC` for options) |
| result.**scheduled\_activation**integer required  Timestamp at which became or will become active (if applicable) |
| result.**scheduled\_deactivation**integer required  Scheduled deactivation time for instrument (if applicable) |
| result.**taker\_fee\_rate**string required  Percent of spot price fee rate for takers |
| result.**tick\_size**string required  Tick size of the instrument, i.e. minimum price increment |
| result.**mark\_price\_fee\_rate\_cap**string or null  Percent of option price fee cap, e.g. 12.5%, null if not applicable |

### Example

ShellJavaScriptPython

```
{request_example_shell}

```

```
{request_example_javascript}

```

```
{request_example_python}

```

> The above command returns JSON structured like this:

JSON

```
{response_example_json}

```

---

## public-get_instruments

**Title:** Public Get_Instruments
**URL:** https://docs.derive.xyz/reference/public-get_instruments

### Method Name

#### `public/get_instruments`

Get all active instruments for a given `currency` and `type`.

### Parameters

|  |
| --- |
| **currency**string required  Underlying currency of asset (`ETH`, `BTC`, etc) |
| **expired**boolean required  If `True`: include expired assets. Note: will soon be capped up to only 1 week in the past. |
| **instrument\_type**string required  `erc20`, `option`, or `perp` enum  `erc20` `option` `perp` |

### Response

|  |
| --- |
| **id**string or integer required |
| **result**array of objects required |
| result[].**amount\_step**string required  Minimum valid increment of order amount |
| result[].**base\_asset\_address**string required  Blockchain address of the base asset |
| result[].**base\_asset\_sub\_id**string required  Sub ID of the specific base asset as defined in Asset.sol |
| result[].**base\_currency**string required  Underlying currency of base asset (`ETH`, `BTC`, etc) |
| result[].**base\_fee**string required  $ base fee added to every taker order |
| result[].**erc20\_details**object or null required  Details of the erc20 asset (if applicable) |
| result[].erc20\_details.**decimals**integer required  Number of decimals of the underlying on-chain ERC20 token |
| result[].erc20\_details.**borrow\_index**string  Latest borrow index as per `CashAsset.sol` implementation |
| result[].erc20\_details.**supply\_index**string  Latest supply index as per `CashAsset.sol` implementation |
| result[].erc20\_details.**underlying\_erc20\_address**string  Address of underlying on-chain ERC20 (not V2 asset) |
|  |
| result[].**fifo\_min\_allocation**string required  Minimum number of contracts that get filled using FIFO. Actual number of contracts that gets filled by FIFO will be the max between this value and (1 - pro\_rata\_fraction) x order\_amount, plus any size leftovers due to rounding. |
| result[].**instrument\_name**string required  Instrument name |
| result[].**instrument\_type**string required  `erc20`, `option`, or `perp` enum  `erc20` `option` `perp` |
| result[].**is\_active**boolean required  If `True`: instrument is tradeable within `activation` and `deactivation` timestamps |
| result[].**maker\_fee\_rate**string required  Percent of spot price fee rate for makers |
| result[].**maximum\_amount**string required  Maximum valid amount of contracts / tokens per trade |
| result[].**minimum\_amount**string required  Minimum valid amount of contracts / tokens per trade |
| result[].**option\_details**object or null required  Details of the option asset (if applicable) |
| result[].option\_details.**expiry**integer required  Unix timestamp of expiry date (in seconds) |
| result[].option\_details.**index**string required  Underlying settlement price index |
| result[].option\_details.**option\_type**string required   enum  `C` `P` |
| result[].option\_details.**strike**string required |
| result[].option\_details.**settlement\_price**string or null  Settlement price of the option |
|  |
| result[].**perp\_details**object or null required  Details of the perp asset (if applicable) |
| result[].perp\_details.**aggregate\_funding**string required  Latest aggregated funding as per `PerpAsset.sol` |
| result[].perp\_details.**funding\_rate**string required  Current hourly funding rate as per `PerpAsset.sol` |
| result[].perp\_details.**index**string required  Underlying spot price index for funding rate |
| result[].perp\_details.**max\_rate\_per\_hour**string required  Max rate per hour as per `PerpAsset.sol` |
| result[].perp\_details.**min\_rate\_per\_hour**string required  Min rate per hour as per `PerpAsset.sol` |
| result[].perp\_details.**static\_interest\_rate**string required  Static interest rate as per `PerpAsset.sol` |
|  |
| result[].**pro\_rata\_amount\_step**string required  Pro-rata fill share of every order is rounded down to be a multiple of this number. Leftovers of the order due to rounding are filled FIFO. |
| result[].**pro\_rata\_fraction**string required  Fraction of order that gets filled using pro-rata matching. If zero, the matching is full FIFO. |
| result[].**quote\_currency**string required  Quote currency (`USD` for perps, `USDC` for options) |
| result[].**scheduled\_activation**integer required  Timestamp at which became or will become active (if applicable) |
| result[].**scheduled\_deactivation**integer required  Scheduled deactivation time for instrument (if applicable) |
| result[].**taker\_fee\_rate**string required  Percent of spot price fee rate for takers |
| result[].**tick\_size**string required  Tick size of the instrument, i.e. minimum price increment |
| result[].**mark\_price\_fee\_rate\_cap**string or null  Percent of option price fee cap, e.g. 12.5%, null if not applicable |

### Example

ShellJavaScriptPython

```
{request_example_shell}

```

```
{request_example_javascript}

```

```
{request_example_python}

```

> The above command returns JSON structured like this:

JSON

```
{response_example_json}

```

---

## public-get_interest_rate_history

**Title:** Public Get_Interest_Rate_History
**URL:** https://docs.derive.xyz/reference/public-get_interest_rate_history

### Method Name

#### `public/get_interest_rate_history`

Get latest USDC interest rate history

### Parameters

|  |
| --- |
| **from\_timestamp\_sec**integer required  From timestamp in seconds |
| **to\_timestamp\_sec**integer required  To timestamp in seconds |
| **page**integer  Page number of results to return (default 1, returns last if above `num_pages`) |
| **page\_size**integer  Number of results per page (default 100, max 1000) |

### Response

|  |
| --- |
| **id**string or integer required |
| **result**object required |
| result.**interest\_rates**array of objects required  List of interest rates, recent first |
| result.interest\_rates[].**block**integer required  Derive Chain block number |
| result.interest\_rates[].**borrow\_apy**string required  Borrow APY |
| result.interest\_rates[].**supply\_apy**string required  Supply APY |
| result.interest\_rates[].**timestamp\_sec**integer required  Timestamp in seconds |
| result.interest\_rates[].**total\_borrow**string required  Total USDC borrowed |
| result.interest\_rates[].**total\_supply**string required  Total USDC supplied |
|  |
| result.**pagination**object required  Pagination info |
| result.pagination.**count**integer required  Total number of items, across all pages |
| result.pagination.**num\_pages**integer required  Number of pages |

### Example

ShellJavaScriptPython

```
{request_example_shell}

```

```
{request_example_javascript}

```

```
{request_example_python}

```

> The above command returns JSON structured like this:

JSON

```
{response_example_json}

```

---

## public-get_invite_code

**Title:** Public Get_Invite_Code
**URL:** https://docs.derive.xyz/reference/public-get_invite_code

### Method Name

#### `public/get_invite_code`

TODO description

### Parameters

|  |
| --- |
| **wallet**string required  Wallet address of the user |

### Response

|  |
| --- |
| **id**string or integer required |
| **result**object required |
| result.**code**string required  Invite code |
| result.**remaining\_uses**integer required  Remaining uses of the invite code. Unlimited use codes will return -1 |

### Example

ShellJavaScriptPython

```
{request_example_shell}

```

```
{request_example_javascript}

```

```
{request_example_python}

```

> The above command returns JSON structured like this:

JSON

```
{response_example_json}

```

---

## public-get_latest_signed_feeds

**Title:** Public Get_Latest_Signed_Feeds
**URL:** https://docs.derive.xyz/reference/public-get_latest_signed_feeds

### Method Name

#### `public/get_latest_signed_feeds`

Get latest signed data feeds

### Parameters

|  |
| --- |
| **currency**string  Currency filter, (defaults to all currencies) |
| **expiry**integer  Expiry filter for options and forward data (defaults to all expiries). Use `0` to get data only for spot and perpetuals |

### Response

|  |
| --- |
| **id**string or integer required |
| **result**object required |
| result.**fwd\_data**object required  currency -> expiry -> latest forward feed data |
| result.**perp\_data**object required  currency -> feed type -> latest perp feed data |
| result.**rate\_data**object required  currency -> expiry -> latest rate feed data |
| result.**spot\_data**object required  currency -> latest spot feed data |
| result.**vol\_data**object required  currency -> expiry -> latest vol feed data |

### Example

ShellJavaScriptPython

```
{request_example_shell}

```

```
{request_example_javascript}

```

```
{request_example_python}

```

> The above command returns JSON structured like this:

JSON

```
{response_example_json}

```

---

## public-get_liquidation_history

**Title:** Public Get_Liquidation_History
**URL:** https://docs.derive.xyz/reference/public-get_liquidation_history

### Method Name

#### `public/get_liquidation_history`

Returns a paginated liquidation history for all subaccounts. Note that the pagination is based on the number of  
raw events that include bids, auction start, and auction end events. This means that the count returned in the  
pagination info will be larger than the total number of auction events. This also means the number of returned  
auctions per page will be smaller than the supplied `page_size`.

### Parameters

|  |
| --- |
| **end\_timestamp**integer  End timestamp of the event history (default current time) |
| **page**integer  Page number of results to return (default 1, returns last if above `num_pages`) |
| **page\_size**integer  Number of results per page (default 100, max 1000) |
| **start\_timestamp**integer  Start timestamp of the event history (default 0) |
| **subaccount\_id**integer  (Optional) Subaccount ID |

### Response

|  |
| --- |
| **id**string or integer required |
| **result**object required |
| result.**auctions**array of objects required  List of auction results |
| result.auctions[].**auction\_id**string required  Unique ID of the auction |
| result.auctions[].**auction\_type**string required  Type of auction enum  `solvent` `insolvent` |
| result.auctions[].**end\_timestamp**integer or null required  Timestamp of the auction end (in ms since UNIX epoch), or `null` if not ended |
| result.auctions[].**fee**string required  Fee paid by the subaccount |
| result.auctions[].**start\_timestamp**integer required  Timestamp of the auction start (in ms since UNIX epoch) |
| result.auctions[].**subaccount\_id**integer required  Liquidated subaccount ID |
| result.auctions[].**tx\_hash**string required  Hash of the transaction that started the auction |
| result.auctions[].**bids**array of objects required  List of auction bid events |
| result.auctions[].bids[].**cash\_received**string required  Cash received by the subaccount for the liquidation. For the liquidated accounts this is the amount the liquidator paid to buy out the percentage of the portfolio. For the liquidator account, this is the amount they received from the security module (if positive) or the amount they paid for the bid (if negative) |
| result.auctions[].bids[].**discount\_pnl**string required  Realized PnL due to liquidating or being liquidated at a discount to mark portfolio value |
| result.auctions[].bids[].**percent\_liquidated**string required  Percent of the subaccount that was liquidated |
| result.auctions[].bids[].**realized\_pnl**string required  Realized PnL of the auction bid, assuming positions are closed at mark price at the time of the liquidation |
| result.auctions[].bids[].**realized\_pnl\_excl\_fees**string required  Realized PnL of the auction bid, excluding fees from total cost basis, assuming positions are closed at mark price at the time of the liquidation |
| result.auctions[].bids[].**timestamp**integer required  Timestamp of the bid (in ms since UNIX epoch) |
| result.auctions[].bids[].**tx\_hash**string required  Hash of the bid transaction |
| result.auctions[].bids[].**amounts\_liquidated**object required  Amounts of each asset that were closed |
| result.auctions[].bids[].**positions\_realized\_pnl**object required  Realized PnL of each position that was closed |
| result.auctions[].bids[].**positions\_realized\_pnl\_excl\_fees**object required  Realized PnL of each position that was closed, excluding fees from total cost basis |
|  |
| result.**pagination**object required  Pagination info |
| result.pagination.**count**integer required  Total number of items, across all pages |
| result.pagination.**num\_pages**integer required  Number of pages |

### Example

ShellJavaScriptPython

```
{request_example_shell}

```

```
{request_example_javascript}

```

```
{request_example_python}

```

> The above command returns JSON structured like this:

JSON

```
{response_example_json}

```

---

## public-get_live_incidents

**Title:** Public Get_Live_Incidents
**URL:** https://docs.derive.xyz/reference/public-get_live_incidents

### Method Name

#### `public/get_live_incidents`

TODO description

### Parameters

|  |
| --- |

### Response

|  |
| --- |
| **id**string or integer required |
| **result**object required |
| result.**incidents**array of objects required  List of ongoing incidents |
| result.incidents[].**creation\_timestamp\_sec**integer required  Timestamp of incident in UTC sec |
| result.incidents[].**label**string required  Incident label |
| result.incidents[].**message**string required  Incident message |
| result.incidents[].**monitor\_type**string required  Incident trigger type enum  `manual` `auto` |
| result.incidents[].**severity**string required  Incident severity enum  `low` `medium` `high` |

### Example

ShellJavaScriptPython

```
{request_example_shell}

```

```
{request_example_javascript}

```

```
{request_example_python}

```

> The above command returns JSON structured like this:

JSON

```
{response_example_json}

```

---

## public-get_maker_program_scores

**Title:** Public Get_Maker_Program_Scores
**URL:** https://docs.derive.xyz/reference/public-get_maker_program_scores

### Method Name

#### `public/get_maker_program_scores`

Get scores breakdown by maker program.

### Parameters

|  |
| --- |
| **epoch\_start\_timestamp**integer required  Start timestamp of the program epoch |
| **program\_name**string required  Program name |

### Response

|  |
| --- |
| **id**string or integer required |
| **result**object required |
| result.**total\_score**string required  Total score across all market makers for the epoch |
| result.**total\_volume**string required  Total volume across all market makers for the epoch |
| result.**program**object required  Program details |
| result.program.**end\_timestamp**integer required  End timestamp of the epoch |
| result.program.**min\_notional**string required  Minimum dollar notional to quote for eligibility |
| result.program.**name**string required  Name of the program |
| result.program.**start\_timestamp**integer required  Start timestamp of the epoch |
| result.program.**asset\_types**array of strings required  List of asset types covered by the program |
| result.program.**currencies**array of strings required  List of currencies covered by the program |
| result.program.**rewards**object required  Rewards for the program as a token -> total reward amount mapping |
|  |
| result.**scores**array of objects required  Scores breakdown of the program by market maker |
| result.scores[].**coverage\_score**string required  Coverag component of the score of the account for this program |
| result.scores[].**holder\_boost**string required  A custom account multiplier for the score due to holding tokens |
| result.scores[].**quality\_score**string required  Quality component of the score of the account for this program |
| result.scores[].**total\_score**string required  Total score of the account for this program |
| result.scores[].**volume**string required  Volume traded by the account for this epoch |
| result.scores[].**volume\_multiplier**string required  Multiplier for the volume traded by the account |
| result.scores[].**wallet**string required  Wallet address of the account |

### Example

ShellJavaScriptPython

```
{request_example_shell}

```

```
{request_example_javascript}

```

```
{request_example_python}

```

> The above command returns JSON structured like this:

JSON

```
{response_example_json}

```

---

## public-get_maker_programs

**Title:** Public Get_Maker_Programs
**URL:** https://docs.derive.xyz/reference/public-get_maker_programs

### Method Name

#### `public/get_maker_programs`

Get all maker programs, including past / historical ones.

### Parameters

|  |
| --- |

### Response

|  |
| --- |
| **id**string or integer required |
| **result**array of objects required |
| result[].**end\_timestamp**integer required  End timestamp of the epoch |
| result[].**min\_notional**string required  Minimum dollar notional to quote for eligibility |
| result[].**name**string required  Name of the program |
| result[].**start\_timestamp**integer required  Start timestamp of the epoch |
| result[].**asset\_types**array of strings required  List of asset types covered by the program |
| result[].**currencies**array of strings required  List of currencies covered by the program |
| result[].**rewards**object required  Rewards for the program as a token -> total reward amount mapping |

### Example

ShellJavaScriptPython

```
{request_example_shell}

```

```
{request_example_javascript}

```

```
{request_example_python}

```

> The above command returns JSON structured like this:

JSON

```
{response_example_json}

```

---

## public-get_margin

**Title:** Public Get_Margin
**URL:** https://docs.derive.xyz/reference/public-get_margin

### Method Name

#### `public/get_margin`

Calculates margin for a given portfolio and (optionally) a simulated state change.  
Does not take into account open orders margin requirements.public/withdraw\_debug

### Parameters

|  |
| --- |
| **margin\_type**string required  `PM` (Portfolio Margin) or `SM` (Standard Margin) enum  `PM` `SM` `PM2` |
| **simulated\_collaterals**array of objects required  List of collaterals in a simulated portfolio |
| simulated\_collaterals[].**amount**string required  Collateral amount to simulate |
| simulated\_collaterals[].**asset\_name**string required  Collateral ERC20 asset name (e.g. ETH, USDC, WSTETH) |
|  |
| **simulated\_positions**array of objects required  List of positions in a simulated portfolio |
| simulated\_positions[].**amount**string required  Position amount to simulate |
| simulated\_positions[].**instrument\_name**string required  Instrument name |
| simulated\_positions[].**entry\_price**string  Only for perps. Entry price to use in the simulation. Mark price is used if not provided. |
|  |
| **market**string  Must be defined for Portfolio Margin |
| **simulated\_collateral\_changes**array of objects  Optional, add collaterals to simulate deposits / withdrawals / spot trades |
| simulated\_collateral\_changes[].**amount**string required  Collateral amount to simulate |
| simulated\_collateral\_changes[].**asset\_name**string required  Collateral ERC20 asset name (e.g. ETH, USDC, WSTETH) |
|  |
| **simulated\_position\_changes**array of objects  Optional, add positions to simulate perp / option trades |
| simulated\_position\_changes[].**amount**string required  Position amount to simulate |
| simulated\_position\_changes[].**instrument\_name**string required  Instrument name |
| simulated\_position\_changes[].**entry\_price**string  Only for perps. Entry price to use in the simulation. Mark price is used if not provided. |

### Response

|  |
| --- |
| **id**string or integer required |
| **result**object required |
| result.**is\_valid\_trade**boolean required  True if trade passes margin requirement |
| result.**post\_initial\_margin**string required  Initial margin requirement post trade |
| result.**post\_maintenance\_margin**string required  Maintenance margin requirement post trade |
| result.**pre\_initial\_margin**string required  Initial margin requirement before trade |
| result.**pre\_maintenance\_margin**string required  Maintenance margin requirement before trade |
| result.**subaccount\_id**integer required  Subaccount\_id |

### Example

ShellJavaScriptPython

```
{request_example_shell}

```

```
{request_example_javascript}

```

```
{request_example_python}

```

> The above command returns JSON structured like this:

JSON

```
{response_example_json}

```

---

## public-get_option_settlement_history

**Title:** Public Get_Option_Settlement_History
**URL:** https://docs.derive.xyz/reference/public-get_option_settlement_history

### Method Name

#### `public/get_option_settlement_history`

Get expired option settlement history for a subaccount

### Parameters

|  |
| --- |
| **page**integer  Page number of results to return (default 1, returns last if above `num_pages`) |
| **page\_size**integer  Number of results per page (default 100, max 1000) |
| **subaccount\_id**integer  Subaccount ID filter (defaults to all if not provided) |

### Response

|  |
| --- |
| **id**string or integer required |
| **result**object required |
| result.**pagination**object required  Pagination info |
| result.pagination.**count**integer required  Total number of items, across all pages |
| result.pagination.**num\_pages**integer required  Number of pages |
|  |
| result.**settlements**array of objects required  List of expired option settlements |
| result.settlements[].**amount**string required  Amount that was settled |
| result.settlements[].**expiry**integer required  Expiry timestamp of the option |
| result.settlements[].**instrument\_name**string required  Instrument name |
| result.settlements[].**option\_settlement\_pnl**string required  USD profit or loss from option settlements calculated as: settlement value - (average cost including fees x amount) |
| result.settlements[].**option\_settlement\_pnl\_excl\_fees**string required  USD profit or loss from option settlements calculated as: settlement value - (average price excluding fees x amount) |
| result.settlements[].**settlement\_price**string required  Price of option settlement |
| result.settlements[].**subaccount\_id**integer required  Subaccount ID of the settlement event |

### Example

ShellJavaScriptPython

```
{request_example_shell}

```

```
{request_example_javascript}

```

```
{request_example_python}

```

> The above command returns JSON structured like this:

JSON

```
{response_example_json}

```

---

## public-get_option_settlement_prices

**Title:** Public Get_Option_Settlement_Prices
**URL:** https://docs.derive.xyz/reference/public-get_option_settlement_prices

### Method Name

#### `public/get_option_settlement_prices`

Get settlement prices by expiry for each currency

### Parameters

|  |
| --- |
| **currency**string required  Currency for which to show expiries |

### Response

|  |
| --- |
| **id**string or integer required |
| **result**object required |
| result.**expiries**array of objects required  List of expiry details |
| result.expiries[].**expiry\_date**string required  Expiry date in `YYYYMMDD` format |
| result.expiries[].**price**string or null required  Settlement price will show None if not yet settled onchain |
| result.expiries[].**utc\_expiry\_sec**integer required  UTC timestamp of expiry |

### Example

ShellJavaScriptPython

```
{request_example_shell}

```

```
{request_example_javascript}

```

```
{request_example_python}

```

> The above command returns JSON structured like this:

JSON

```
{response_example_json}

```

---

## public-get_points

**Title:** Public Get_Points
**URL:** https://docs.derive.xyz/reference/public-get_points

### Method Name

#### `public/get_points`

Get all points for user across all programs. This request is cached in WSGI.

### Parameters

|  |
| --- |
| **program**string required  Program for which to count points |
| **wallet**string required  Wallet address of the user |

### Response

|  |
| --- |
| **id**string or integer required |
| **result**object required |
| result.**flag**string required  Flag user for special treatment |
| result.**parent**string required  Wallet address of referrer |
| result.**percent\_share\_of\_points**string required  Deprecated |
| result.**total\_notional\_volume**string required  Total $ notional volume traded by the user in program |
| result.**total\_users**integer required  Deprecated |
| result.**user\_rank**integer or null required  Deprecated |
| result.**points**object required  Points earned per category |

### Example

ShellJavaScriptPython

```
{request_example_shell}

```

```
{request_example_javascript}

```

```
{request_example_python}

```

> The above command returns JSON structured like this:

JSON

```
{response_example_json}

```

---

## public-get_points_leaderboard

**Title:** Public Get_Points_Leaderboard
**URL:** https://docs.derive.xyz/reference/public-get_points_leaderboard

### Method Name

#### `public/get_points_leaderboard`

Get top 250 users based on points earned. Can paginate where each page contains exactly 250 users.

### Parameters

|  |
| --- |
| **page**integer required  Page number of leaderboard. Each page holds up to 500 users and starts at 1. |
| **program**string required  Program for which to count leaderboard |

### Response

|  |
| --- |
| **id**string or integer required |
| **result**object required |
| result.**pages**integer required  Total number of pages in the leaderboard |
| result.**total\_users**integer required  Total number of users in the program |
| result.**leaderboard**array of objects required  List of up to 500 users in order of highest points |
| result.leaderboard[].**flag**string required  Flag user for special treatment |
| result.leaderboard[].**parent**string required  Wallet address of referrer |
| result.leaderboard[].**percent\_share\_of\_points**string required  Percentage of total points earned by user |
| result.leaderboard[].**points**string required  Total points for the user |
| result.leaderboard[].**rank**integer required  Leaderboard rank of the user |
| result.leaderboard[].**total\_volume**string required  Deprecated |
| result.leaderboard[].**wallet**string required  Wallet address of the user |

### Example

ShellJavaScriptPython

```
{request_example_shell}

```

```
{request_example_javascript}

```

```
{request_example_python}

```

> The above command returns JSON structured like this:

JSON

```
{response_example_json}

```

---

## public-get_spot_feed_history

**Title:** Public Get_Spot_Feed_History
**URL:** https://docs.derive.xyz/reference/public-get_spot_feed_history

### Method Name

#### `public/get_spot_feed_history`

Get spot feed history by currency  
  
DB: read replica

### Parameters

|  |
| --- |
| **currency**string required  Currency |
| **end\_timestamp**integer required  End timestamp |
| **period**integer required  Period |
| **start\_timestamp**integer required  Start timestamp |

### Response

|  |
| --- |
| **id**string or integer required |
| **result**object required |
| result.**currency**string required  Currency |
| result.**spot\_feed\_history**array of objects required  Spot feed history |
| result.spot\_feed\_history[].**price**string required  Spot price |
| result.spot\_feed\_history[].**timestamp**integer required  Timestamp of when the spot price was recored into the database |
| result.spot\_feed\_history[].**timestamp\_bucket**integer required  Timestamp bucket; this value is regularly spaced out with `period` seconds between data points, missing values are forward-filled from earlier data where possible, if no earlier data is available, values are back-filled from the first observed data point |

### Example

ShellJavaScriptPython

```
{request_example_shell}

```

```
{request_example_javascript}

```

```
{request_example_python}

```

> The above command returns JSON structured like this:

JSON

```
{response_example_json}

```

---

## public-get_spot_feed_history_candles

**Title:** Public Get_Spot_Feed_History_Candles
**URL:** https://docs.derive.xyz/reference/public-get_spot_feed_history_candles

### Method Name

#### `public/get_spot_feed_history_candles`

Get spot feed history candles by currency  
  
DB: read replica

### Parameters

|  |
| --- |
| **currency**string required  Currency |
| **end\_timestamp**integer required  End timestamp |
| **period**string required  Period enum  `60` `300` `900` `1800` `3600` `14400` `28800` `86400` `604800` |
| **start\_timestamp**integer required  Start timestamp |

### Response

|  |
| --- |
| **id**string or integer required |
| **result**object required |
| result.**currency**string required  Currency |
| result.**spot\_feed\_history**array of objects required  Spot feed history candles |
| result.spot\_feed\_history[].**close\_price**string required  Close price |
| result.spot\_feed\_history[].**high\_price**string required  High price |
| result.spot\_feed\_history[].**low\_price**string required  Low price |
| result.spot\_feed\_history[].**open\_price**string required  Open price |
| result.spot\_feed\_history[].**price**string required  Spot price |
| result.spot\_feed\_history[].**timestamp**integer required  Timestamp of when the spot price was recored into the database |
| result.spot\_feed\_history[].**timestamp\_bucket**integer required  Timestamp bucket; this value is regularly spaced out with `period` seconds between data points, missing values are forward-filled from earlier data where possible, if no earlier data is available, values are back-filled from the first observed data point |

### Example

ShellJavaScriptPython

```
{request_example_shell}

```

```
{request_example_javascript}

```

```
{request_example_python}

```

> The above command returns JSON structured like this:

JSON

```
{response_example_json}

```

---

## public-get_ticker

**Title:** Public Get_Ticker
**URL:** https://docs.derive.xyz/reference/public-get_ticker

### Method Name

#### `public/get_ticker`

Get ticker information (best bid / ask, instrument contraints, fees info, etc.) for a single instrument

### Parameters

|  |
| --- |
| **instrument\_name**string required  Instrument name |

### Response

|  |
| --- |
| **id**string or integer required |
| **result**object required |
| result.**amount\_step**string required  Minimum valid increment of order amount |
| result.**base\_asset\_address**string required  Blockchain address of the base asset |
| result.**base\_asset\_sub\_id**string required  Sub ID of the specific base asset as defined in Asset.sol |
| result.**base\_currency**string required  Underlying currency of base asset (`ETH`, `BTC`, etc) |
| result.**base\_fee**string required  $ base fee added to every taker order |
| result.**best\_ask\_amount**string required  Amount of contracts / tokens available at best ask price |
| result.**best\_ask\_price**string required  Best ask price |
| result.**best\_bid\_amount**string required  Amount of contracts / tokens available at best bid price |
| result.**best\_bid\_price**string required  Best bid price |
| result.**erc20\_details**object or null required  Details of the erc20 asset (if applicable) |
| result.erc20\_details.**decimals**integer required  Number of decimals of the underlying on-chain ERC20 token |
| result.erc20\_details.**borrow\_index**string  Latest borrow index as per `CashAsset.sol` implementation |
| result.erc20\_details.**supply\_index**string  Latest supply index as per `CashAsset.sol` implementation |
| result.erc20\_details.**underlying\_erc20\_address**string  Address of underlying on-chain ERC20 (not V2 asset) |
|  |
| result.**fifo\_min\_allocation**string required  Minimum number of contracts that get filled using FIFO. Actual number of contracts that gets filled by FIFO will be the max between this value and (1 - pro\_rata\_fraction) x order\_amount, plus any size leftovers due to rounding. |
| result.**five\_percent\_ask\_depth**string required  Total amount of contracts / tokens available at 5 percent above best ask price |
| result.**five\_percent\_bid\_depth**string required  Total amount of contracts / tokens available at 5 percent below best bid price |
| result.**index\_price**string required  Index price |
| result.**instrument\_name**string required  Instrument name |
| result.**instrument\_type**string required  `erc20`, `option`, or `perp` enum  `erc20` `option` `perp` |
| result.**is\_active**boolean required  If `True`: instrument is tradeable within `activation` and `deactivation` timestamps |
| result.**maker\_fee\_rate**string required  Percent of spot price fee rate for makers |
| result.**mark\_price**string required  Mark price |
| result.**max\_price**string required  Maximum price at which an agressive buyer can be matched. Any portion of a market order that would execute above this price will be cancelled. A limit buy order with limit price above this value is treated as post only (i.e. it will be rejected if it would cross any existing resting order). |
| result.**maximum\_amount**string required  Maximum valid amount of contracts / tokens per trade |
| result.**min\_price**string required  Minimum price at which an agressive seller can be matched. Any portion of a market order that would execute below this price will be cancelled. A limit sell order with limit price below this value is treated as post only (i.e. it will be rejected if it would cross any existing resting order). |
| result.**minimum\_amount**string required  Minimum valid amount of contracts / tokens per trade |
| result.**option\_details**object or null required  Details of the option asset (if applicable) |
| result.option\_details.**expiry**integer required  Unix timestamp of expiry date (in seconds) |
| result.option\_details.**index**string required  Underlying settlement price index |
| result.option\_details.**option\_type**string required   enum  `C` `P` |
| result.option\_details.**strike**string required |
| result.option\_details.**settlement\_price**string or null  Settlement price of the option |
|  |
| result.**option\_pricing**object or null required  Greeks, forward price, iv and mark price of the instrument (options only) |
| result.option\_pricing.**ask\_iv**string required  Implied volatility of the current best ask |
| result.option\_pricing.**bid\_iv**string required  Implied volatility of the current best bid |
| result.option\_pricing.**delta**string required  Delta of the option |
| result.option\_pricing.**discount\_factor**string required  Discount factor used to calculate option premium |
| result.option\_pricing.**forward\_price**string required  Forward price used to calculate option premium |
| result.option\_pricing.**gamma**string required  Gamma of the option |
| result.option\_pricing.**iv**string required  Implied volatility of the option |
| result.option\_pricing.**mark\_price**string required  Mark price of the option |
| result.option\_pricing.**rho**string required  Rho of the option |
| result.option\_pricing.**theta**string required  Theta of the option |
| result.option\_pricing.**vega**string required  Vega of the option |
|  |
| result.**perp\_details**object or null required  Details of the perp asset (if applicable) |
| result.perp\_details.**aggregate\_funding**string required  Latest aggregated funding as per `PerpAsset.sol` |
| result.perp\_details.**funding\_rate**string required  Current hourly funding rate as per `PerpAsset.sol` |
| result.perp\_details.**index**string required  Underlying spot price index for funding rate |
| result.perp\_details.**max\_rate\_per\_hour**string required  Max rate per hour as per `PerpAsset.sol` |
| result.perp\_details.**min\_rate\_per\_hour**string required  Min rate per hour as per `PerpAsset.sol` |
| result.perp\_details.**static\_interest\_rate**string required  Static interest rate as per `PerpAsset.sol` |
|  |
| result.**pro\_rata\_amount\_step**string required  Pro-rata fill share of every order is rounded down to be a multiple of this number. Leftovers of the order due to rounding are filled FIFO. |
| result.**pro\_rata\_fraction**string required  Fraction of order that gets filled using pro-rata matching. If zero, the matching is full FIFO. |
| result.**quote\_currency**string required  Quote currency (`USD` for perps, `USDC` for options) |
| result.**scheduled\_activation**integer required  Timestamp at which became or will become active (if applicable) |
| result.**scheduled\_deactivation**integer required  Scheduled deactivation time for instrument (if applicable) |
| result.**taker\_fee\_rate**string required  Percent of spot price fee rate for takers |
| result.**tick\_size**string required  Tick size of the instrument, i.e. minimum price increment |
| result.**timestamp**integer required  Timestamp of the ticker feed snapshot |
| result.**open\_interest**object required  Margin type of subaccount (`PM` (Portfolio Margin) or `SM` (Standard Margin)) -> (current open interest, open interest cap, manager currency) |
| result.**stats**object required  Aggregate trading stats for the last 24 hours |
| result.stats.**contract\_volume**string required  Number of contracts traded during last 24 hours |
| result.stats.**high**string required  Highest trade price during last 24h |
| result.stats.**low**string required  Lowest trade price during last 24h |
| result.stats.**num\_trades**string required  Number of trades during last 24h |
| result.stats.**open\_interest**string required  Current total open interest |
| result.stats.**percent\_change**string required  24-hour price change expressed as a percentage. Options: percent change in vol; Perps: percent change in mark price |
| result.stats.**usd\_change**string required  24-hour price change in USD. |
|  |
| result.**mark\_price\_fee\_rate\_cap**string or null  Percent of option price fee cap, e.g. 12.5%, null if not applicable |

### Example

ShellJavaScriptPython

```
{request_example_shell}

```

```
{request_example_javascript}

```

```
{request_example_python}

```

> The above command returns JSON structured like this:

JSON

```
{response_example_json}

```

---

## public-get_time

**Title:** Public Get_Time
**URL:** https://docs.derive.xyz/reference/public-get_time

### Method Name

#### `public/get_time`

TODO description

### Parameters

|  |
| --- |

### Response

|  |
| --- |
| **id**string or integer required |
| **result**integer required  Current time in milliseconds since UNIX epoch |

### Example

ShellJavaScriptPython

```
{request_example_shell}

```

```
{request_example_javascript}

```

```
{request_example_python}

```

> The above command returns JSON structured like this:

JSON

```
{response_example_json}

```

---

## public-get_trade_history

**Title:** Public Get_Trade_History
**URL:** https://docs.derive.xyz/reference/public-get_trade_history

### Method Name

#### `public/get_trade_history`

Get trade history for a subaccount, with filter parameters.

### Parameters

|  |
| --- |
| **currency**string  Currency to filter by (defaults to all) |
| **from\_timestamp**integer  Earliest timestamp to filter by (in ms since Unix epoch). If not provied, defaults to 0. |
| **instrument\_name**string  Instrument name to filter by (defaults to all) |
| **instrument\_type**string  Instrument type to filter by (defaults to all) enum  `erc20` `option` `perp` |
| **page**integer  Page number of results to return (default 1, returns last if above `num_pages`) |
| **page\_size**integer  Number of results per page (default 100, max 1000) |
| **subaccount\_id**integer  Subaccount ID to filter by |
| **to\_timestamp**integer  Latest timestamp to filter by (in ms since Unix epoch). If not provied, defaults to returning all data up to current time. |
| **trade\_id**string  Trade ID to filter by. If set, all other filters are ignored |
| **tx\_hash**string  On-chain tx hash to filter by. If set, all other filters are ignored |
| **tx\_status**string  Transaction status to filter by (default `settled`). enum  `settled` `reverted` `timed_out` |

### Response

|  |
| --- |
| **id**string or integer required |
| **result**object required |
| result.**pagination**object required  Pagination info |
| result.pagination.**count**integer required  Total number of items, across all pages |
| result.pagination.**num\_pages**integer required  Number of pages |
|  |
| result.**trades**array of objects required  List of trades |
| result.trades[].**direction**string required  Order direction enum  `buy` `sell` |
| result.trades[].**expected\_rebate**string required  Expected rebate for this trade |
| result.trades[].**index\_price**string required  Index price of the underlying at the time of the trade |
| result.trades[].**instrument\_name**string required  Instrument name |
| result.trades[].**liquidity\_role**string required  Role of the user in the trade enum  `maker` `taker` |
| result.trades[].**mark\_price**string required  Mark price of the instrument at the time of the trade |
| result.trades[].**quote\_id**string or null required  Quote ID if the trade was executed via RFQ |
| result.trades[].**realized\_pnl**string required  Realized PnL for this trade |
| result.trades[].**realized\_pnl\_excl\_fees**string required  Realized PnL for this trade using cost accounting that excludes fees |
| result.trades[].**subaccount\_id**integer required  Subaccount ID |
| result.trades[].**timestamp**integer required  Trade timestamp (in ms since Unix epoch) |
| result.trades[].**trade\_amount**string required  Amount filled in this trade |
| result.trades[].**trade\_fee**string required  Fee for this trade |
| result.trades[].**trade\_id**string required  Trade ID |
| result.trades[].**trade\_price**string required  Price at which the trade was filled |
| result.trades[].**tx\_hash**string required  Blockchain transaction hash |
| result.trades[].**tx\_status**string required  Blockchain transaction status enum  `settled` `reverted` `timed_out` |
| result.trades[].**wallet**string required  Wallet address (owner) of the subaccount |

### Example

ShellJavaScriptPython

```
{request_example_shell}

```

```
{request_example_javascript}

```

```
{request_example_python}

```

> The above command returns JSON structured like this:

JSON

```
{response_example_json}

```

---

## public-get_transaction

**Title:** Public Get_Transaction
**URL:** https://docs.derive.xyz/reference/public-get_transaction

### Method Name

#### `public/get_transaction`

Used for getting a transaction by its transaction id

### Parameters

|  |
| --- |
| **transaction\_id**string required  transaction\_id of the transaction to get |

### Response

|  |
| --- |
| **id**string or integer required |
| **result**object required |
| result.**data**string required  Data used to create transaction |
| result.**error\_log**string or null required  Error log if failed tx |
| result.**status**string required  Status of the transaction enum  `requested` `pending` `settled` `reverted` `ignored` `timed_out` |
| result.**transaction\_hash**string or null required  Transaction hash of a pending tx |

### Example

ShellJavaScriptPython

```
{request_example_shell}

```

```
{request_example_javascript}

```

```
{request_example_python}

```

> The above command returns JSON structured like this:

JSON

```
{response_example_json}

```

---

## public-get_tree_roots

**Title:** Public Get_Tree_Roots
**URL:** https://docs.derive.xyz/reference/public-get_tree_roots

### Method Name

#### `public/get_tree_roots`

Returns the roots of a tree from which full tree can be constructed with public/get\_descendant\_tree.

### Parameters

|  |
| --- |

### Response

|  |
| --- |
| **id**string or integer required |
| **result**object required |
| result.**roots**array of strings required  Roots of tree from which whole tree can be constructed |

### Example

ShellJavaScriptPython

```
{request_example_shell}

```

```
{request_example_javascript}

```

```
{request_example_python}

```

> The above command returns JSON structured like this:

JSON

```
{response_example_json}

```

---

## public-get_vault_balances

**Title:** Public Get_Vault_Balances
**URL:** https://docs.derive.xyz/reference/public-get_vault_balances

### Method Name

#### `public/get_vault_balances`

Get all vault assets held by user. Can query by smart contract address or smart contract owner.  
  
Includes VaultERC20Pool balances

### Parameters

|  |
| --- |
| **smart\_contract\_owner**string  If wallet not provided, can query balances by EOA that owns smart contract wallet |
| **wallet**string  Ethereum wallet address of smart contract |

### Response

|  |
| --- |
| **id**string or integer required |
| **result**array of objects required |
| result[].**address**string required |
| result[].**amount**string required |
| result[].**chain\_id**integer required |
| result[].**name**string required |
| result[].**vault\_asset\_type**string required |

### Example

ShellJavaScriptPython

```
{request_example_shell}

```

```
{request_example_javascript}

```

```
{request_example_python}

```

> The above command returns JSON structured like this:

JSON

```
{response_example_json}

```

---

## public-get_vault_share

**Title:** Public Get_Vault_Share
**URL:** https://docs.derive.xyz/reference/public-get_vault_share

### Method Name

#### `public/get_vault_share`

Gets the value of a vault's token against the base currency, underlying currency, and USD for a timestamp range.  
  
The name of the vault from the Vault proxy contract is used to fetch the vault's value.

### Parameters

|  |
| --- |
| **from\_timestamp\_sec**integer required  From timestamp in seconds |
| **to\_timestamp\_sec**integer required  To timestamp in seconds |
| **vault\_name**string required  Name of the vault |
| **page**integer  Page number of results to return (default 1, returns last if above `num_pages`) |
| **page\_size**integer  Number of results per page (default 100, max 1000) |

### Response

|  |
| --- |
| **id**string or integer required |
| **result**object required |
| result.**pagination**object required  Pagination info |
| result.pagination.**count**integer required  Total number of items, across all pages |
| result.pagination.**num\_pages**integer required  Number of pages |
|  |
| result.**vault\_shares**array of objects required  List of vault history shares, recent first |
| result.vault\_shares[].**base\_value**string required  The value of the vault's token against the base currency. Ex: rswETHC vs rswETH |
| result.vault\_shares[].**block\_number**integer required  The Derive chain block number |
| result.vault\_shares[].**block\_timestamp**integer required  Timestamp of the Derive chain block number |
| result.vault\_shares[].**underlying\_value**string or null required  The value of the vault's token against the underlying currency. Ex: rswETHC vs ETH |
| result.vault\_shares[].**usd\_value**string required  The value of the vault's token against USD |

### Example

ShellJavaScriptPython

```
{request_example_shell}

```

```
{request_example_javascript}

```

```
{request_example_python}

```

> The above command returns JSON structured like this:

JSON

```
{response_example_json}

```

---

## public-get_vault_statistics

**Title:** Public Get_Vault_Statistics
**URL:** https://docs.derive.xyz/reference/public-get_vault_statistics

### Method Name

#### `public/get_vault_statistics`

Gets all the latest vault shareRate, totalSupply and TVL values for all vaults.  
  
For data on shares across chains, use public/get\_vault\_assets.

### Parameters

|  |
| --- |

### Response

|  |
| --- |
| **id**string or integer required |
| **result**array of objects required |
| result[].**base\_value**string required  The value of the vault's token against the base currency. Ex: rswETHC vs rswETH |
| result[].**block\_number**integer required  The Derive chain block number |
| result[].**block\_timestamp**integer required  Timestamp of the Derive chain block number |
| result[].**subaccount\_value\_at\_last\_trade**string or null required  Will return None before vault creates subaccount or if no trades found. |
| result[].**total\_supply**string required  Total supply of the vault's token on Derive chain |
| result[].**underlying\_value**string or null required  The value of the vault's token against the underlying currency. Ex: rswETHC vs ETH |
| result[].**usd\_tvl**string required  Total USD TVL of the vault |
| result[].**usd\_value**string required  The value of the vault's token against USD |
| result[].**vault\_name**string required  Name of the vault |

### Example

ShellJavaScriptPython

```
{request_example_shell}

```

```
{request_example_javascript}

```

```
{request_example_python}

```

> The above command returns JSON structured like this:

JSON

```
{response_example_json}

```

---

## public-login

**Title:** Public Login
**URL:** https://docs.derive.xyz/reference/public-login

### Method Name

#### `public/login`

Authenticate a websocket connection. Unavailable via HTTP.

### Parameters

|  |
| --- |
| **signature**string required  Signature of the timestamp, signed with the wallet's private key or a session key |
| **timestamp**string required  Message that was signed, in the form of a timestamp in ms since Unix epoch |
| **wallet**string required  Public key (wallet) of the account |

### Response

|  |
| --- |
| **id**string or integer required |
| **result**array of integers required  List of subaccount IDs that have been authenticated |

### Example

ShellJavaScriptPython

```
{request_example_shell}

```

```
{request_example_javascript}

```

```
{request_example_python}

```

> The above command returns JSON structured like this:

JSON

```
{response_example_json}

```

---

## public-margin_watch

**Title:** Public Margin_Watch
**URL:** https://docs.derive.xyz/reference/public-margin_watch

### Method Name

#### `public/margin_watch`

Calculates MtM and maintenance margin for a given subaccount.

### Parameters

|  |
| --- |
| **subaccount\_id**integer required  Subaccount ID to get margin for. |
| **force\_onchain**boolean  Force the fetching of on-chain balances, default False. |

### Response

|  |
| --- |
| **id**string or integer required |
| **result**object required |
| result.**currency**string required  Currency of subaccount |
| result.**initial\_margin**string required  Total initial margin requirement of all positions and collaterals. |
| result.**maintenance\_margin**string required  Total maintenance margin requirement of all positions and collaterals.If this value falls below zero, the subaccount will be flagged for liquidation. |
| result.**margin\_type**string required  Margin type of subaccount (`PM` (Portfolio Margin) or `SM` (Standard Margin)) enum  `PM` `SM` `PM2` |
| result.**subaccount\_id**integer required  Subaccount\_id |
| result.**subaccount\_value**string required  Total mark-to-market value of all positions and collaterals |
| result.**valuation\_timestamp**integer required  Timestamp (in seconds since epoch) of when margin and MtM were computed. |
| result.**collaterals**array of objects required  All collaterals that count towards margin of subaccount |
| result.collaterals[].**amount**string required  Asset amount of given collateral |
| result.collaterals[].**asset\_name**string required  Asset name |
| result.collaterals[].**asset\_type**string required  Type of asset collateral (currently always `erc20`) enum  `erc20` `option` `perp` |
| result.collaterals[].**initial\_margin**string required  USD value of collateral that contributes to initial margin |
| result.collaterals[].**maintenance\_margin**string required  USD value of collateral that contributes to maintenance margin |
| result.collaterals[].**mark\_price**string required  Current mark price of the asset |
| result.collaterals[].**mark\_value**string required  USD value of the collateral (amount \* mark price) |
|  |
| result.**positions**array of objects required  All active positions of subaccount |
| result.positions[].**amount**string required  Position amount held by subaccount |
| result.positions[].**delta**string required  Asset delta (w.r.t. forward price for options, `1.0` for perps) |
| result.positions[].**gamma**string required  Asset gamma (zero for non-options) |
| result.positions[].**index\_price**string required  Current index (oracle) price for position's currency |
| result.positions[].**initial\_margin**string required  USD initial margin requirement for this position |
| result.positions[].**instrument\_name**string required  Instrument name (same as the base Asset name) |
| result.positions[].**instrument\_type**string required  `erc20`, `option`, or `perp` enum  `erc20` `option` `perp` |
| result.positions[].**liquidation\_price**string or null required  Index price at which position will be liquidated |
| result.positions[].**maintenance\_margin**string required  USD maintenance margin requirement for this position |
| result.positions[].**mark\_price**string required  Current mark price for position's instrument |
| result.positions[].**mark\_value**string required  USD value of the position; this represents how much USD can be recieved by fully closing the position at the current oracle price |
| result.positions[].**theta**string required  Asset theta (zero for non-options) |
| result.positions[].**vega**string required  Asset vega (zero for non-options) |

### Example

ShellJavaScriptPython

```
{request_example_shell}

```

```
{request_example_javascript}

```

```
{request_example_python}

```

> The above command returns JSON structured like this:

JSON

```
{response_example_json}

```

---

## public-register_session_key

**Title:** Public Register_Session_Key
**URL:** https://docs.derive.xyz/reference/public-register_session_key

### Method Name

#### `public/register_session_key`

Register or update expiry of an existing session key.  
Currently, this only supports creating admin level session keys.  
Keys with fewer permissions are registered via `/register_scoped_session_key`  
  
Expiries updated on admin session keys may not happen immediately due to waiting for the onchain transaction to settle.

### Parameters

|  |
| --- |
| **expiry\_sec**integer required  Expiry of the session key |
| **label**string required  Ethereum wallet address |
| **public\_session\_key**string required  Session key in the form of an Ethereum EOA |
| **signed\_raw\_tx**string required  A signed RLP encoded ETH transaction in form of a hex string (same as `w3.eth.account.sign_transaction(unsigned_tx, private_key).rawTransaction.hex()`) |
| **wallet**string required  Ethereum wallet address of account |

### Response

|  |
| --- |
| **id**string or integer required |
| **result**object required |
| result.**label**string required  User-defined session key label |
| result.**public\_session\_key**string required  Session key in the form of an Ethereum EOA |
| result.**transaction\_id**string required  ID to lookup status of transaction |

### Example

ShellJavaScriptPython

```
{request_example_shell}

```

```
{request_example_javascript}

```

```
{request_example_python}

```

> The above command returns JSON structured like this:

JSON

```
{response_example_json}

```

---

## public-send_quote_debug

**Title:** Public Send_Quote_Debug
**URL:** https://docs.derive.xyz/reference/public-send_quote_debug

### Method Name

#### `public/send_quote_debug`

Sends a quote in response to an RFQ request.  
The legs supplied in the parameters must exactly match those in the RFQ.

### Parameters

|  |
| --- |
| **direction**string required  Quote direction, `buy` means trading each leg at its direction, `sell` means trading each leg in the opposite direction. enum  `buy` `sell` |
| **max\_fee**string required  Max fee ($ for the full trade). Request will be rejected if the supplied max fee is below the estimated fee for this trade. |
| **nonce**integer required  Unique nonce defined as a concatenated `UTC timestamp in ms` and `random number up to 6 digits` (e.g. 1695836058725001, where 001 is the random number) |
| **rfq\_id**string required  RFQ ID the quote is for |
| **signature**string required  Ethereum signature of the quote |
| **signature\_expiry\_sec**integer required  Unix timestamp in seconds. Expiry MUST be at least 310 seconds from now. Once time till signature expiry reaches 300 seconds, the quote will be considered expired. This buffer is meant to ensure the trade can settle on chain in case of a blockchain congestion. |
| **signer**string required  Owner wallet address or registered session key that signed the quote |
| **subaccount\_id**integer required  Subaccount ID |
| **legs**array of objects required  Quote legs |
| legs[].**amount**string required  Amount in units of the base |
| legs[].**direction**string required  Leg direction enum  `buy` `sell` |
| legs[].**instrument\_name**string required  Instrument name |
| legs[].**price**string required  Leg price |
|  |
| **label**string  Optional user-defined label for the quote |
| **mmp**boolean  Whether the quote is tagged for market maker protections (default false) |

### Response

|  |
| --- |
| **id**string or integer required |
| **result**object required |
| result.**action\_hash**string required  Keccak hashed action data |
| result.**encoded\_data**string required  ABI encoded deposit data |
| result.**encoded\_data\_hashed**string required  Keccak hashed encoded\_data |
| result.**typed\_data\_hash**string required  EIP 712 typed data hash |

### Example

ShellJavaScriptPython

```
{request_example_shell}

```

```
{request_example_javascript}

```

```
{request_example_python}

```

> The above command returns JSON structured like this:

JSON

```
{response_example_json}

```

---

## public-statistics

**Title:** Public Statistics
**URL:** https://docs.derive.xyz/reference/public-statistics

### Method Name

#### `public/statistics`

Get statistics for a specific instrument or instrument type

### Parameters

|  |
| --- |
| **instrument\_name**string required  Instrument name or 'ALL', 'OPTION', 'PERP', 'SPOT' |
| **currency**string  Currency for stats |
| **end\_time**integer  End time for statistics in ms |

### Response

|  |
| --- |
| **id**string or integer required |
| **result**object required |
| result.**daily\_fees**string required  24h Fees |
| result.**daily\_notional\_volume**string required  24h Notional volume |
| result.**daily\_premium\_volume**string required  24h Premium volume |
| result.**daily\_trades**integer required  24h Trades |
| result.**open\_interest**string required  Open interest |
| result.**total\_fees**string required  Total fees |
| result.**total\_notional\_volume**string required  Total notional volume |
| result.**total\_premium\_volume**string required  Total premium volume |
| result.**total\_trades**integer required  Total trades |

### Example

ShellJavaScriptPython

```
{request_example_shell}

```

```
{request_example_javascript}

```

```
{request_example_python}

```

> The above command returns JSON structured like this:

JSON

```
{response_example_json}

```

---

## public-validate_invite_code

**Title:** Public Validate_Invite_Code
**URL:** https://docs.derive.xyz/reference/public-validate_invite_code

### Method Name

#### `public/validate_invite_code`

Validate if invite is valid and useable

### Parameters

|  |
| --- |
| **code**string or integer required  5 digit, alpha numeric code |

### Response

|  |
| --- |
| **id**string or integer required |
| **result**string required   enum  `ok` |

### Example

ShellJavaScriptPython

```
{request_example_shell}

```

```
{request_example_javascript}

```

```
{request_example_python}

```

> The above command returns JSON structured like this:

JSON

```
{response_example_json}

```

---

## public-withdraw_debug

**Title:** Public Withdraw_Debug
**URL:** https://docs.derive.xyz/reference/public-withdraw_debug

### Method Name

#### `public/withdraw_debug`

Used for debugging only, do not use in production. Will return the incremental encoded and hashed data.  
  
See guides in Documentation for more.

### Parameters

|  |
| --- |
| **amount**string required  Amount of the asset to withdraw |
| **asset\_name**string required  Name of asset to withdraw |
| **nonce**integer required  Unique nonce defined as  (e.g. 1695836058725001, where 001 is the random number) |
| **signature\_expiry\_sec**integer required  Unix timestamp in seconds. Expiry MUST be >5min from now |
| **signer**string required  Ethereum wallet address that is signing the withdraw |
| **subaccount\_id**integer required  Subaccount\_id |
| **is\_atomic\_signing**boolean  Used by vaults to determine whether the signature is an EIP-1271 signature |

### Response

|  |
| --- |
| **id**string or integer required |
| **result**object required |
| result.**action\_hash**string required  Keccak hashed action data |
| result.**encoded\_data**string required  ABI encoded deposit data |
| result.**encoded\_data\_hashed**string required  Keccak hashed encoded\_data |
| result.**typed\_data\_hash**string required  EIP 712 typed data hash |

### Example

ShellJavaScriptPython

```
{request_example_shell}

```

```
{request_example_javascript}

```

```
{request_example_python}

```

> The above command returns JSON structured like this:

JSON

```
{response_example_json}

```

---

# Private


## post_private-cancel

**Title:** Post_Private Cancel
**URL:** https://docs.derive.xyz/reference/post_private-cancel

ShellNodeRubyPHPPython

Base URL

https://api.lyra.finance/private/cancel

Click `Try It!` to start a request and see the response here!

---

## post_private-cancel-all

**Title:** Post_Private Cancel All
**URL:** https://docs.derive.xyz/reference/post_private-cancel-all

ShellNodeRubyPHPPython

Base URL

https://api.lyra.finance/private/cancel\_all

Click `Try It!` to start a request and see the response here!

---

## post_private-cancel-batch-quotes

**Title:** Post_Private Cancel Batch Quotes
**URL:** https://docs.derive.xyz/reference/post_private-cancel-batch-quotes

ShellNodeRubyPHPPython

Base URL

https://api.lyra.finance/private/cancel\_batch\_quotes

Click `Try It!` to start a request and see the response here!

---

## post_private-cancel-batch-rfqs

**Title:** Post_Private Cancel Batch Rfqs
**URL:** https://docs.derive.xyz/reference/post_private-cancel-batch-rfqs

ShellNodeRubyPHPPython

Base URL

https://api.lyra.finance/private/cancel\_batch\_rfqs

Click `Try It!` to start a request and see the response here!

---

## post_private-cancel-by-instrument

**Title:** Post_Private Cancel By Instrument
**URL:** https://docs.derive.xyz/reference/post_private-cancel-by-instrument

ShellNodeRubyPHPPython

Base URL

https://api.lyra.finance/private/cancel\_by\_instrument

Click `Try It!` to start a request and see the response here!

---

## post_private-cancel-by-label

**Title:** Post_Private Cancel By Label
**URL:** https://docs.derive.xyz/reference/post_private-cancel-by-label

ShellNodeRubyPHPPython

Base URL

https://api.lyra.finance/private/cancel\_by\_label

Click `Try It!` to start a request and see the response here!

---

## post_private-cancel-by-nonce

**Title:** Post_Private Cancel By Nonce
**URL:** https://docs.derive.xyz/reference/post_private-cancel-by-nonce

ShellNodeRubyPHPPython

Base URL

https://api.lyra.finance/private/cancel\_by\_nonce

Click `Try It!` to start a request and see the response here!

---

## post_private-cancel-quote

**Title:** Post_Private Cancel Quote
**URL:** https://docs.derive.xyz/reference/post_private-cancel-quote

ShellNodeRubyPHPPython

Base URL

https://api.lyra.finance/private/cancel\_quote

Click `Try It!` to start a request and see the response here!

---

## post_private-cancel-rfq

**Title:** Post_Private Cancel Rfq
**URL:** https://docs.derive.xyz/reference/post_private-cancel-rfq

ShellNodeRubyPHPPython

Base URL

https://api.lyra.finance/private/cancel\_rfq

Click `Try It!` to start a request and see the response here!

---

## post_private-cancel-trigger-order

**Title:** Post_Private Cancel Trigger Order
**URL:** https://docs.derive.xyz/reference/post_private-cancel-trigger-order

ShellNodeRubyPHPPython

Base URL

https://api.lyra.finance/private/cancel\_trigger\_order

Click `Try It!` to start a request and see the response here!

---

## post_private-change-subaccount-label

**Title:** Post_Private Change Subaccount Label
**URL:** https://docs.derive.xyz/reference/post_private-change-subaccount-label

ShellNodeRubyPHPPython

Base URL

https://api.lyra.finance/private/change\_subaccount\_label

Click `Try It!` to start a request and see the response here!

---

## post_private-create-subaccount

**Title:** Post_Private Create Subaccount
**URL:** https://docs.derive.xyz/reference/post_private-create-subaccount

ShellNodeRubyPHPPython

Base URL

https://api.lyra.finance/private/create\_subaccount

Click `Try It!` to start a request and see the response here!

---

## post_private-deposit

**Title:** Post_Private Deposit
**URL:** https://docs.derive.xyz/reference/post_private-deposit

ShellNodeRubyPHPPython

Base URL

https://api.lyra.finance/private/deposit

Click `Try It!` to start a request and see the response here!

---

## post_private-edit-session-key

**Title:** Post_Private Edit Session Key
**URL:** https://docs.derive.xyz/reference/post_private-edit-session-key

ShellNodeRubyPHPPython

Base URL

https://api.lyra.finance/private/edit\_session\_key

Click `Try It!` to start a request and see the response here!

---

## post_private-execute-quote

**Title:** Post_Private Execute Quote
**URL:** https://docs.derive.xyz/reference/post_private-execute-quote

ShellNodeRubyPHPPython

Base URL

https://api.lyra.finance/private/execute\_quote

Click `Try It!` to start a request and see the response here!

---

## post_private-expired-and-cancelled-history

**Title:** Post_Private Expired And Cancelled History
**URL:** https://docs.derive.xyz/reference/post_private-expired-and-cancelled-history

ShellNodeRubyPHPPython

Base URL

https://api.lyra.finance/private/expired\_and\_cancelled\_history

Click `Try It!` to start a request and see the response here!

---

## post_private-get-account

**Title:** Post_Private Get Account
**URL:** https://docs.derive.xyz/reference/post_private-get-account

ShellNodeRubyPHPPython

Base URL

https://api.lyra.finance/private/get\_account

Click `Try It!` to start a request and see the response here!

---

## post_private-get-all-portfolios

**Title:** Post_Private Get All Portfolios
**URL:** https://docs.derive.xyz/reference/post_private-get-all-portfolios

ShellNodeRubyPHPPython

Base URL

https://api.lyra.finance/private/get\_all\_portfolios

Click `Try It!` to start a request and see the response here!

---

## post_private-get-collaterals

**Title:** Post_Private Get Collaterals
**URL:** https://docs.derive.xyz/reference/post_private-get-collaterals

ShellNodeRubyPHPPython

Base URL

https://api.lyra.finance/private/get\_collaterals

Click `Try It!` to start a request and see the response here!

---

## post_private-get-deposit-history

**Title:** Post_Private Get Deposit History
**URL:** https://docs.derive.xyz/reference/post_private-get-deposit-history

ShellNodeRubyPHPPython

Base URL

https://api.lyra.finance/private/get\_deposit\_history

Click `Try It!` to start a request and see the response here!

---

## post_private-get-erc20-transfer-history

**Title:** Post_Private Get Erc20 Transfer History
**URL:** https://docs.derive.xyz/reference/post_private-get-erc20-transfer-history

ShellNodeRubyPHPPython

Base URL

https://api.lyra.finance/private/get\_erc20\_transfer\_history

Click `Try It!` to start a request and see the response here!

---

## post_private-get-funding-history

**Title:** Post_Private Get Funding History
**URL:** https://docs.derive.xyz/reference/post_private-get-funding-history

ShellNodeRubyPHPPython

Base URL

https://api.lyra.finance/private/get\_funding\_history

Click `Try It!` to start a request and see the response here!

---

## post_private-get-interest-history

**Title:** Post_Private Get Interest History
**URL:** https://docs.derive.xyz/reference/post_private-get-interest-history

ShellNodeRubyPHPPython

Base URL

https://api.lyra.finance/private/get\_interest\_history

Click `Try It!` to start a request and see the response here!

---

## post_private-get-liquidation-history

**Title:** Post_Private Get Liquidation History
**URL:** https://docs.derive.xyz/reference/post_private-get-liquidation-history

ShellNodeRubyPHPPython

Base URL

https://api.lyra.finance/private/get\_liquidation\_history

Click `Try It!` to start a request and see the response here!

---

## post_private-get-liquidator-history

**Title:** Post_Private Get Liquidator History
**URL:** https://docs.derive.xyz/reference/post_private-get-liquidator-history

ShellNodeRubyPHPPython

Base URL

https://api.lyra.finance/private/get\_liquidator\_history

Click `Try It!` to start a request and see the response here!

---

## post_private-get-margin

**Title:** Post_Private Get Margin
**URL:** https://docs.derive.xyz/reference/post_private-get-margin

ShellNodeRubyPHPPython

Base URL

https://api.lyra.finance/private/get\_margin

Click `Try It!` to start a request and see the response here!

---

## post_private-get-mmp-config

**Title:** Post_Private Get Mmp Config
**URL:** https://docs.derive.xyz/reference/post_private-get-mmp-config

ShellNodeRubyPHPPython

Base URL

https://api.lyra.finance/private/get\_mmp\_config

Click `Try It!` to start a request and see the response here!

---

## post_private-get-notifications

**Title:** Post_Private Get Notifications
**URL:** https://docs.derive.xyz/reference/post_private-get-notifications

ShellNodeRubyPHPPython

Base URL

https://api.lyra.finance/private/get\_notifications

Click `Try It!` to start a request and see the response here!

---

## post_private-get-open-orders

**Title:** Post_Private Get Open Orders
**URL:** https://docs.derive.xyz/reference/post_private-get-open-orders

ShellNodeRubyPHPPython

Base URL

https://api.lyra.finance/private/get\_open\_orders

Click `Try It!` to start a request and see the response here!

---

## post_private-get-option-settlement-history

**Title:** Post_Private Get Option Settlement History
**URL:** https://docs.derive.xyz/reference/post_private-get-option-settlement-history

ShellNodeRubyPHPPython

Base URL

https://api.lyra.finance/private/get\_option\_settlement\_history

Click `Try It!` to start a request and see the response here!

---

## post_private-get-order

**Title:** Post_Private Get Order
**URL:** https://docs.derive.xyz/reference/post_private-get-order

ShellNodeRubyPHPPython

Base URL

https://api.lyra.finance/private/get\_order

Click `Try It!` to start a request and see the response here!

---

## post_private-get-order-history

**Title:** Post_Private Get Order History
**URL:** https://docs.derive.xyz/reference/post_private-get-order-history

ShellNodeRubyPHPPython

Base URL

https://api.lyra.finance/private/get\_order\_history

Click `Try It!` to start a request and see the response here!

---

## post_private-get-orders

**Title:** Post_Private Get Orders
**URL:** https://docs.derive.xyz/reference/post_private-get-orders

ShellNodeRubyPHPPython

Base URL

https://api.lyra.finance/private/get\_orders

Click `Try It!` to start a request and see the response here!

---

## post_private-get-positions

**Title:** Post_Private Get Positions
**URL:** https://docs.derive.xyz/reference/post_private-get-positions

ShellNodeRubyPHPPython

Base URL

https://api.lyra.finance/private/get\_positions

Click `Try It!` to start a request and see the response here!

---

## post_private-get-quotes

**Title:** Post_Private Get Quotes
**URL:** https://docs.derive.xyz/reference/post_private-get-quotes

ShellNodeRubyPHPPython

Base URL

https://api.lyra.finance/private/get\_quotes

Click `Try It!` to start a request and see the response here!

---

## post_private-get-rfqs

**Title:** Post_Private Get Rfqs
**URL:** https://docs.derive.xyz/reference/post_private-get-rfqs

ShellNodeRubyPHPPython

Base URL

https://api.lyra.finance/private/get\_rfqs

Click `Try It!` to start a request and see the response here!

---

## post_private-get-subaccount

**Title:** Post_Private Get Subaccount
**URL:** https://docs.derive.xyz/reference/post_private-get-subaccount

ShellNodeRubyPHPPython

Base URL

https://api.lyra.finance/private/get\_subaccount

Click `Try It!` to start a request and see the response here!

---

## post_private-get-subaccount-value-history

**Title:** Post_Private Get Subaccount Value History
**URL:** https://docs.derive.xyz/reference/post_private-get-subaccount-value-history

ShellNodeRubyPHPPython

Base URL

https://api.lyra.finance/private/get\_subaccount\_value\_history

Click `Try It!` to start a request and see the response here!

---

## post_private-get-subaccounts

**Title:** Post_Private Get Subaccounts
**URL:** https://docs.derive.xyz/reference/post_private-get-subaccounts

ShellNodeRubyPHPPython

Base URL

https://api.lyra.finance/private/get\_subaccounts

Click `Try It!` to start a request and see the response here!

---

## post_private-get-trade-history

**Title:** Post_Private Get Trade History
**URL:** https://docs.derive.xyz/reference/post_private-get-trade-history

ShellNodeRubyPHPPython

Base URL

https://api.lyra.finance/private/get\_trade\_history

Click `Try It!` to start a request and see the response here!

---

## post_private-get-withdrawal-history

**Title:** Post_Private Get Withdrawal History
**URL:** https://docs.derive.xyz/reference/post_private-get-withdrawal-history

ShellNodeRubyPHPPython

Base URL

https://api.lyra.finance/private/get\_withdrawal\_history

Click `Try It!` to start a request and see the response here!

---

## post_private-liquidate

**Title:** Post_Private Liquidate
**URL:** https://docs.derive.xyz/reference/post_private-liquidate

ShellNodeRubyPHPPython

Base URL

https://api.lyra.finance/private/liquidate

Click `Try It!` to start a request and see the response here!

---

## post_private-order

**Title:** Post_Private Order
**URL:** https://docs.derive.xyz/reference/post_private-order

ShellNodeRubyPHPPython

Base URL

https://api.lyra.finance/private/order

Click `Try It!` to start a request and see the response here!

---

## post_private-order-debug

**Title:** Post_Private Order Debug
**URL:** https://docs.derive.xyz/reference/post_private-order-debug

ShellNodeRubyPHPPython

Base URL

https://api.lyra.finance/private/order\_debug

Click `Try It!` to start a request and see the response here!

---

## post_private-poll-quotes

**Title:** Post_Private Poll Quotes
**URL:** https://docs.derive.xyz/reference/post_private-poll-quotes

ShellNodeRubyPHPPython

Base URL

https://api.lyra.finance/private/poll\_quotes

Click `Try It!` to start a request and see the response here!

---

## post_private-poll-rfqs

**Title:** Post_Private Poll Rfqs
**URL:** https://docs.derive.xyz/reference/post_private-poll-rfqs

ShellNodeRubyPHPPython

Base URL

https://api.lyra.finance/private/poll\_rfqs

Click `Try It!` to start a request and see the response here!

---

## post_private-register-scoped-session-key

**Title:** Post_Private Register Scoped Session Key
**URL:** https://docs.derive.xyz/reference/post_private-register-scoped-session-key

ShellNodeRubyPHPPython

Base URL

https://api.lyra.finance/private/register\_scoped\_session\_key

Click `Try It!` to start a request and see the response here!

---

## post_private-replace

**Title:** Post_Private Replace
**URL:** https://docs.derive.xyz/reference/post_private-replace

ShellNodeRubyPHPPython

Base URL

https://api.lyra.finance/private/replace

Click `Try It!` to start a request and see the response here!

---

## post_private-reset-mmp

**Title:** Post_Private Reset Mmp
**URL:** https://docs.derive.xyz/reference/post_private-reset-mmp

ShellNodeRubyPHPPython

Base URL

https://api.lyra.finance/private/reset\_mmp

Click `Try It!` to start a request and see the response here!

---

## post_private-rfq-get-best-quote

**Title:** Post_Private Rfq Get Best Quote
**URL:** https://docs.derive.xyz/reference/post_private-rfq-get-best-quote

ShellNodeRubyPHPPython

Base URL

https://api.lyra.finance/private/rfq\_get\_best\_quote

Click `Try It!` to start a request and see the response here!

---

## post_private-send-quote

**Title:** Post_Private Send Quote
**URL:** https://docs.derive.xyz/reference/post_private-send-quote

ShellNodeRubyPHPPython

Base URL

https://api.lyra.finance/private/send\_quote

Click `Try It!` to start a request and see the response here!

---

## post_private-send-rfq

**Title:** Post_Private Send Rfq
**URL:** https://docs.derive.xyz/reference/post_private-send-rfq

ShellNodeRubyPHPPython

Base URL

https://api.lyra.finance/private/send\_rfq

Click `Try It!` to start a request and see the response here!

---

## post_private-set-cancel-on-disconnect

**Title:** Post_Private Set Cancel On Disconnect
**URL:** https://docs.derive.xyz/reference/post_private-set-cancel-on-disconnect

ShellNodeRubyPHPPython

Base URL

https://api.lyra.finance/private/set\_cancel\_on\_disconnect

Click `Try It!` to start a request and see the response here!

---

## post_private-set-mmp-config

**Title:** Post_Private Set Mmp Config
**URL:** https://docs.derive.xyz/reference/post_private-set-mmp-config

ShellNodeRubyPHPPython

Base URL

https://api.lyra.finance/private/set\_mmp\_config

Click `Try It!` to start a request and see the response here!

---

## post_private-transfer-erc20

**Title:** Post_Private Transfer Erc20
**URL:** https://docs.derive.xyz/reference/post_private-transfer-erc20

ShellNodeRubyPHPPython

Base URL

https://api.lyra.finance/private/transfer\_erc20

Click `Try It!` to start a request and see the response here!

---

## post_private-transfer-position

**Title:** Post_Private Transfer Position
**URL:** https://docs.derive.xyz/reference/post_private-transfer-position

ShellNodeRubyPHPPython

Base URL

https://api.lyra.finance/private/transfer\_position

Click `Try It!` to start a request and see the response here!

---

## post_private-transfer-positions

**Title:** Post_Private Transfer Positions
**URL:** https://docs.derive.xyz/reference/post_private-transfer-positions

ShellNodeRubyPHPPython

Base URL

https://api.lyra.finance/private/transfer\_positions

Click `Try It!` to start a request and see the response here!

---

## post_private-update-notifications

**Title:** Post_Private Update Notifications
**URL:** https://docs.derive.xyz/reference/post_private-update-notifications

ShellNodeRubyPHPPython

Base URL

https://api.lyra.finance/private/update\_notifications

Click `Try It!` to start a request and see the response here!

---

## post_private-withdraw

**Title:** Post_Private Withdraw
**URL:** https://docs.derive.xyz/reference/post_private-withdraw

ShellNodeRubyPHPPython

Base URL

https://api.lyra.finance/private/withdraw

Click `Try It!` to start a request and see the response here!

---

## private-cancel

**Title:** Private Cancel
**URL:** https://docs.derive.xyz/reference/private-cancel

### Method Name

#### `private/cancel`

Cancel a single order.  
  
Other `private/cancel_*` routes are available through both REST and WebSocket.  
Required minimum session key permission level is `admin`

### Parameters

|  |
| --- |
| **instrument\_name**string required |
| **order\_id**string required |
| **subaccount\_id**integer required |

### Response

|  |
| --- |
| **id**string or integer required |
| **result**object required |
| result.**amount**string required  Order amount in units of the base |
| result.**average\_price**string required  Average fill price |
| result.**cancel\_reason**string required  If cancelled, reason behind order cancellation enum  `user_request` `mmp_trigger` `insufficient_margin` `signed_max_fee_too_low` `cancel_on_disconnect` `ioc_or_market_partial_fill` `session_key_deregistered` `subaccount_withdrawn` `compliance` `trigger_failed` `validation_failed` |
| result.**creation\_timestamp**integer required  Creation timestamp (in ms since Unix epoch) |
| result.**direction**string required  Order direction enum  `buy` `sell` |
| result.**filled\_amount**string required  Total filled amount for the order |
| result.**instrument\_name**string required  Instrument name |
| result.**is\_transfer**boolean required  Whether the order was generated through `private/transfer_position` |
| result.**label**string required  Optional user-defined label for the order |
| result.**last\_update\_timestamp**integer required  Last update timestamp (in ms since Unix epoch) |
| result.**limit\_price**string required  Limit price in quote currency |
| result.**max\_fee**string required  Max fee in units of the quote currency |
| result.**mmp**boolean required  Whether the order is tagged for market maker protections |
| result.**nonce**integer required  Unique nonce defined as  (e.g. 1695836058725001, where 001 is the random number) |
| result.**order\_fee**string required  Total order fee paid so far |
| result.**order\_id**string required  Order ID |
| result.**order\_status**string required  Order status enum  `open` `filled` `cancelled` `expired` `untriggered` |
| result.**order\_type**string required  Order type enum  `limit` `market` |
| result.**quote\_id**string or null required  Quote ID if the trade was executed via RFQ |
| result.**signature**string required  Ethereum signature of the order |
| result.**signature\_expiry\_sec**integer required  Signature expiry timestamp |
| result.**signer**string required  Owner wallet address or registered session key that signed order |
| result.**subaccount\_id**integer required  Subaccount ID |
| result.**time\_in\_force**string required  Time in force enum  `gtc` `post_only` `fok` `ioc` |
| result.**replaced\_order\_id**string or null  If replaced, ID of the order that was replaced |
| result.**trigger\_price**string or null  (Required for trigger orders) Index or Market price to trigger order at |
| result.**trigger\_price\_type**string or null  (Required for trigger orders) Trigger with Index or Mark Price enum  `mark` `index` |
| result.**trigger\_reject\_message**string or null  (Required for trigger orders) Error message if error occured during trigger |
| result.**trigger\_type**string or null  (Required for trigger orders) Stop-loss or Take-profit. enum  `stoploss` `takeprofit` |

### Example

ShellJavaScriptPython

```
{request_example_shell}

```

```
{request_example_javascript}

```

```
{request_example_python}

```

> The above command returns JSON structured like this:

JSON

```
{response_example_json}

```

---

## private-cancel-all

**Title:** Private Cancel All
**URL:** https://docs.derive.xyz/reference/private-cancel-all

### Method Name

#### `private/cancel_all`

Cancel all orders for this instrument.  
Required minimum session key permission level is `admin`

### Parameters

|  |
| --- |
| **subaccount\_id**integer required |

### Response

|  |
| --- |
| **id**string or integer required |
| **result**string required   enum  `ok` |

### Example

ShellJavaScriptPython

```
{request_example_shell}

```

```
{request_example_javascript}

```

```
{request_example_python}

```

> The above command returns JSON structured like this:

JSON

```
{response_example_json}

```

---

## private-cancel_batch_quotes

**Title:** Private Cancel_Batch_Quotes
**URL:** https://docs.derive.xyz/reference/private-cancel_batch_quotes

### Method Name

#### `private/cancel_batch_quotes`

Cancels quotes given optional filters. If no filters are provided, all quotes by the subaccount are cancelled.  
All filters are combined using `AND` logic, so mutually exclusive filters will result in no quotes being cancelled.  
Required minimum session key permission level is `admin`

### Parameters

|  |
| --- |
| **subaccount\_id**integer required  Subaccount ID |
| **label**string  Cancel quotes with this label |
| **nonce**integer  Cancel quote with this nonce |
| **quote\_id**string  Quote ID to cancel |
| **rfq\_id**string  Cancel quotes for this RFQ ID |

### Response

|  |
| --- |
| **id**string or integer required |
| **result**object required |
| result.**cancelled\_ids**array of strings required  Quote IDs of the cancelled quotes |

### Example

ShellJavaScriptPython

```
{request_example_shell}

```

```
{request_example_javascript}

```

```
{request_example_python}

```

> The above command returns JSON structured like this:

JSON

```
{response_example_json}

```

---

## private-cancel_batch_rfqs

**Title:** Private Cancel_Batch_Rfqs
**URL:** https://docs.derive.xyz/reference/private-cancel_batch_rfqs

### Method Name

#### `private/cancel_batch_rfqs`

Cancels RFQs given optional filters.  
If no filters are provided, all RFQs for the subaccount are cancelled.  
All filters are combined using `AND` logic, so mutually exclusive filters will result in no RFQs being cancelled.  
Required minimum session key permission level is `account`

### Parameters

|  |
| --- |
| **subaccount\_id**integer required  Subaccount ID |
| **label**string  Cancel RFQs with this label |
| **nonce**integer  Cancel RFQ with this nonce |
| **rfq\_id**string  RFQ ID to cancel |

### Response

|  |
| --- |
| **id**string or integer required |
| **result**object required |
| result.**cancelled\_ids**array of strings required  RFQ IDs of the cancelled RFQs |

### Example

ShellJavaScriptPython

```
{request_example_shell}

```

```
{request_example_javascript}

```

```
{request_example_python}

```

> The above command returns JSON structured like this:

JSON

```
{response_example_json}

```

---

## private-cancel_by_instrument

**Title:** Private Cancel_By_Instrument
**URL:** https://docs.derive.xyz/reference/private-cancel_by_instrument

### Method Name

#### `private/cancel_by_instrument`

Cancel all orders for this instrument.  
Required minimum session key permission level is `admin`

### Parameters

|  |
| --- |
| **instrument\_name**string required  Cancel all orders for this instrument |
| **subaccount\_id**integer required  Subaccount ID |

### Response

|  |
| --- |
| **id**string or integer required |
| **result**object required |
| result.**cancelled\_orders**integer required  Number of cancelled orders |

### Example

ShellJavaScriptPython

```
{request_example_shell}

```

```
{request_example_javascript}

```

```
{request_example_python}

```

> The above command returns JSON structured like this:

JSON

```
{response_example_json}

```

---

## private-cancel_by_label

**Title:** Private Cancel_By_Label
**URL:** https://docs.derive.xyz/reference/private-cancel_by_label

### Method Name

#### `private/cancel_by_label`

Cancel all open orders for a given subaccount and a given label. If instrument\_name is provided, only orders for that instrument will be cancelled.  
Required minimum session key permission level is `admin`

### Parameters

|  |
| --- |
| **label**string required  Cancel all orders for this label |
| **subaccount\_id**integer required  Subaccount ID |
| **instrument\_name**string  Instrument name. If not provided, all orders for all instruments with the label will be cancelled. If provided, request counts as a regular matching request for ratelimit purposes. |

### Response

|  |
| --- |
| **id**string or integer required |
| **result**object required |
| result.**cancelled\_orders**integer required  Number of cancelled orders |

### Example

ShellJavaScriptPython

```
{request_example_shell}

```

```
{request_example_javascript}

```

```
{request_example_python}

```

> The above command returns JSON structured like this:

JSON

```
{response_example_json}

```

---

## private-cancel_by_nonce

**Title:** Private Cancel_By_Nonce
**URL:** https://docs.derive.xyz/reference/private-cancel_by_nonce

### Method Name

#### `private/cancel_by_nonce`

Cancel a single order by nonce. Uses up that nonce if the order does not exist, so any future orders with that nonce will fail  
Required minimum session key permission level is `admin`

### Parameters

|  |
| --- |
| **instrument\_name**string required  Instrument name |
| **nonce**integer required  Cancel an order with this nonce |
| **subaccount\_id**integer required  Subaccount ID |
| **wallet**string required  Wallet address |

### Response

|  |
| --- |
| **id**string or integer required |
| **result**object required |
| result.**cancelled\_orders**integer required  Number of cancelled orders |

### Example

ShellJavaScriptPython

```
{request_example_shell}

```

```
{request_example_javascript}

```

```
{request_example_python}

```

> The above command returns JSON structured like this:

JSON

```
{response_example_json}

```

---

## private-cancel_quote

**Title:** Private Cancel_Quote
**URL:** https://docs.derive.xyz/reference/private-cancel_quote

### Method Name

#### `private/cancel_quote`

Cancels an open quote.  
Required minimum session key permission level is `admin`

### Parameters

|  |
| --- |
| **quote\_id**string required  Quote ID to cancel |
| **subaccount\_id**integer required  Subaccount ID |

### Response

|  |
| --- |
| **id**string or integer required |
| **result**object required |
| result.**cancel\_reason**string required  Cancel reason, if any enum  `user_request` `insufficient_margin` `signed_max_fee_too_low` `mmp_trigger` `cancel_on_disconnect` `session_key_deregistered` `subaccount_withdrawn` `rfq_no_longer_open` `compliance` |
| result.**creation\_timestamp**integer required  Creation timestamp in ms since Unix epoch |
| result.**direction**string required  Quote direction enum  `buy` `sell` |
| result.**fee**string required  Fee paid for this quote (if executed) |
| result.**is\_transfer**boolean required  Whether the order was generated through `private/transfer_position` |
| result.**label**string required  User-defined label, if any |
| result.**last\_update\_timestamp**integer required  Last update timestamp in ms since Unix epoch |
| result.**legs\_hash**string required  Hash of the legs of the best quote to be signed by the taker. |
| result.**liquidity\_role**string required  Liquidity role enum  `maker` `taker` |
| result.**max\_fee**string required  Signed max fee |
| result.**mmp**boolean required  Whether the quote is tagged for market maker protections (default false) |
| result.**nonce**integer required  Nonce |
| result.**quote\_id**string required  Quote ID |
| result.**rfq\_id**string required  RFQ ID |
| result.**signature**string required  Ethereum signature of the quote |
| result.**signature\_expiry\_sec**integer required  Unix timestamp in seconds |
| result.**signer**string required  Owner wallet address or registered session key that signed the quote |
| result.**status**string required  Status enum  `open` `filled` `cancelled` `expired` |
| result.**subaccount\_id**integer required  Subaccount ID |
| result.**tx\_hash**string or null required  Blockchain transaction hash (only for executed quotes) |
| result.**tx\_status**string or null required  Blockchain transaction status (only for executed quotes) enum  `requested` `pending` `settled` `reverted` `ignored` `timed_out` |
| result.**legs**array of objects required  Quote legs |
| result.legs[].**amount**string required  Amount in units of the base |
| result.legs[].**direction**string required  Leg direction enum  `buy` `sell` |
| result.legs[].**instrument\_name**string required  Instrument name |
| result.legs[].**price**string required  Leg price |

### Example

ShellJavaScriptPython

```
{request_example_shell}

```

```
{request_example_javascript}

```

```
{request_example_python}

```

> The above command returns JSON structured like this:

JSON

```
{response_example_json}

```

---

## private-cancel_rfq

**Title:** Private Cancel_Rfq
**URL:** https://docs.derive.xyz/reference/private-cancel_rfq

### Method Name

#### `private/cancel_rfq`

Cancels a single RFQ by id.  
Required minimum session key permission level is `account`

### Parameters

|  |
| --- |
| **rfq\_id**string required  RFQ ID to cancel |
| **subaccount\_id**integer required  Subaccount ID |

### Response

|  |
| --- |
| **id**string or integer required |
| **result**string required  The result of this method call, `ok` if successful enum  `ok` |

### Example

ShellJavaScriptPython

```
{request_example_shell}

```

```
{request_example_javascript}

```

```
{request_example_python}

```

> The above command returns JSON structured like this:

JSON

```
{response_example_json}

```

---

## private-cancel_trigger_order

**Title:** Private Cancel_Trigger_Order
**URL:** https://docs.derive.xyz/reference/private-cancel_trigger_order

### Method Name

#### `private/cancel_trigger_order`

Cancels a trigger order.  
Required minimum session key permission level is `admin`

### Parameters

|  |
| --- |
| **order\_id**string required |
| **subaccount\_id**integer required |

### Response

|  |
| --- |
| **id**string or integer required |
| **result**object required |
| result.**amount**string required  Order amount in units of the base |
| result.**average\_price**string required  Average fill price |
| result.**cancel\_reason**string required  If cancelled, reason behind order cancellation enum  `user_request` `mmp_trigger` `insufficient_margin` `signed_max_fee_too_low` `cancel_on_disconnect` `ioc_or_market_partial_fill` `session_key_deregistered` `subaccount_withdrawn` `compliance` `trigger_failed` `validation_failed` |
| result.**creation\_timestamp**integer required  Creation timestamp (in ms since Unix epoch) |
| result.**direction**string required  Order direction enum  `buy` `sell` |
| result.**filled\_amount**string required  Total filled amount for the order |
| result.**instrument\_name**string required  Instrument name |
| result.**is\_transfer**boolean required  Whether the order was generated through `private/transfer_position` |
| result.**label**string required  Optional user-defined label for the order |
| result.**last\_update\_timestamp**integer required  Last update timestamp (in ms since Unix epoch) |
| result.**limit\_price**string required  Limit price in quote currency |
| result.**max\_fee**string required  Max fee in units of the quote currency |
| result.**mmp**boolean required  Whether the order is tagged for market maker protections |
| result.**nonce**integer required  Unique nonce defined as  (e.g. 1695836058725001, where 001 is the random number) |
| result.**order\_fee**string required  Total order fee paid so far |
| result.**order\_id**string required  Order ID |
| result.**order\_status**string required  Order status enum  `open` `filled` `cancelled` `expired` `untriggered` |
| result.**order\_type**string required  Order type enum  `limit` `market` |
| result.**quote\_id**string or null required  Quote ID if the trade was executed via RFQ |
| result.**signature**string required  Ethereum signature of the order |
| result.**signature\_expiry\_sec**integer required  Signature expiry timestamp |
| result.**signer**string required  Owner wallet address or registered session key that signed order |
| result.**subaccount\_id**integer required  Subaccount ID |
| result.**time\_in\_force**string required  Time in force enum  `gtc` `post_only` `fok` `ioc` |
| result.**replaced\_order\_id**string or null  If replaced, ID of the order that was replaced |
| result.**trigger\_price**string or null  (Required for trigger orders) Index or Market price to trigger order at |
| result.**trigger\_price\_type**string or null  (Required for trigger orders) Trigger with Index or Mark Price enum  `mark` `index` |
| result.**trigger\_reject\_message**string or null  (Required for trigger orders) Error message if error occured during trigger |
| result.**trigger\_type**string or null  (Required for trigger orders) Stop-loss or Take-profit. enum  `stoploss` `takeprofit` |

### Example

ShellJavaScriptPython

```
{request_example_shell}

```

```
{request_example_javascript}

```

```
{request_example_python}

```

> The above command returns JSON structured like this:

JSON

```
{response_example_json}

```

---

## private-change_session_key_label

**Title:** Private Change_Session_Key_Label
**URL:** https://docs.derive.xyz/reference/private-change_session_key_label

### Method Name

#### `private/change_session_key_label`

TODO description

### Parameters

|  |
| --- |
| **label**string required |

### Response

|  |
| --- |
| **id**string or integer required |
| **result**object required |
| result.**label**string required |

### Example

ShellJavaScriptPython

```
{request_example_shell}

```

```
{request_example_javascript}

```

```
{request_example_python}

```

> The above command returns JSON structured like this:

JSON

```
{response_example_json}

```

---

## private-change_subaccount_label

**Title:** Private Change_Subaccount_Label
**URL:** https://docs.derive.xyz/reference/private-change_subaccount_label

### Method Name

#### `private/change_subaccount_label`

Change a user defined label for given subaccount  
Required minimum session key permission level is `account`

### Parameters

|  |
| --- |
| **label**string required  User defined label |
| **subaccount\_id**integer required  Subaccount\_id |

### Response

|  |
| --- |
| **id**string or integer required |
| **result**object required |
| result.**label**string required  User defined label |
| result.**subaccount\_id**integer required  Subaccount\_id |

### Example

ShellJavaScriptPython

```
{request_example_shell}

```

```
{request_example_javascript}

```

```
{request_example_python}

```

> The above command returns JSON structured like this:

JSON

```
{response_example_json}

```

---

## private-create_subaccount

**Title:** Private Create_Subaccount
**URL:** https://docs.derive.xyz/reference/private-create_subaccount

### Method Name

#### `private/create_subaccount`

Create a new subaccount under a given wallet, and deposit an asset into that subaccount.  
  
See `public/create_subaccount_debug` for debugging invalid signature issues or go to guides in Documentation.  
Required minimum session key permission level is `admin`

### Parameters

|  |
| --- |
| **amount**string required  Amount of the asset to deposit |
| **asset\_name**string required  Name of asset to deposit |
| **margin\_type**string required  `PM` (Portfolio Margin) or `SM` (Standard Margin) enum  `PM` `SM` `PM2` |
| **nonce**integer required  Unique nonce defined as  (e.g. 1695836058725001, where 001 is the random number) |
| **signature**string required  Ethereum signature of the deposit |
| **signature\_expiry\_sec**integer required  Unix timestamp in seconds. Expiry MUST be >5min from now |
| **signer**string required  Ethereum wallet address that is signing the deposit |
| **wallet**string required  Ethereum wallet address |
| **currency**string  Base currency of the subaccount (only for `PM`) |

### Response

|  |
| --- |
| **id**string or integer required |
| **result**object required |
| result.**status**string required  `requested` |
| result.**transaction\_id**string required  Transaction id of the request |

### Example

ShellJavaScriptPython

```
{request_example_shell}

```

```
{request_example_javascript}

```

```
{request_example_python}

```

> The above command returns JSON structured like this:

JSON

```
{response_example_json}

```

---

## private-deposit

**Title:** Private Deposit
**URL:** https://docs.derive.xyz/reference/private-deposit

### Method Name

#### `private/deposit`

Deposit an asset to a subaccount.  
  
See `public/deposit_debug' for debugging invalid signature issues or go to guides in Documentation. Required minimum session key permission level is` admin`

### Parameters

|  |
| --- |
| **amount**string required  Amount of the asset to deposit |
| **asset\_name**string required  Name of asset to deposit |
| **nonce**integer required  Unique nonce defined as  (e.g. 1695836058725001, where 001 is the random number) |
| **signature**string required  Ethereum signature of the deposit |
| **signature\_expiry\_sec**integer required  Unix timestamp in seconds. Expiry MUST be >5min from now |
| **signer**string required  Ethereum wallet address that is signing the deposit |
| **subaccount\_id**integer required  Subaccount\_id |
| **is\_atomic\_signing**boolean  Used by vaults to determine whether the signature is an EIP-1271 signature |

### Response

|  |
| --- |
| **id**string or integer required |
| **result**object required |
| result.**status**string required  `requested` |
| result.**transaction\_id**string required  Transaction id of the deposit |

### Example

ShellJavaScriptPython

```
{request_example_shell}

```

```
{request_example_javascript}

```

```
{request_example_python}

```

> The above command returns JSON structured like this:

JSON

```
{response_example_json}

```

---

## private-edit_session_key

**Title:** Private Edit_Session_Key
**URL:** https://docs.derive.xyz/reference/private-edit_session_key

### Method Name

#### `private/edit_session_key`

Edits session key parameters such as label and IP whitelist.  
For non-admin keys you can also toggle whether to disable a particular key.  
Disabling non-admin keys must be done through /deregister\_session\_key  
Required minimum session key permission level is `account`

### Parameters

|  |
| --- |
| **public\_session\_key**string required  Session key in the form of an Ethereum EOA |
| **wallet**string required  Ethereum wallet address of account |
| **disable**boolean  Flag whether or not to disable to session key. Defaulted to false. Only allowed for non-admin keys. Admin keys must go through `/deregister_session_key` for now. |
| **ip\_whitelist**array of strings  Optional list of whitelisted IPs, an empty list can be supplied to whitelist all IPs |
| **label**string  Optional new label for the session key |

### Response

|  |
| --- |
| **id**string or integer required |
| **result**object required |
| result.**expiry\_sec**integer required  Session key expiry timestamp in sec |
| result.**label**string required  User-defined session key label |
| result.**public\_session\_key**string required  Public session key address (Ethereum EOA) |
| result.**scope**string required  Session key permission level scope |
| result.**ip\_whitelist**array of strings required  List of whitelisted IPs, if empty then any IP is allowed. |

### Example

ShellJavaScriptPython

```
{request_example_shell}

```

```
{request_example_javascript}

```

```
{request_example_python}

```

> The above command returns JSON structured like this:

JSON

```
{response_example_json}

```

---

## private-execute_quote

**Title:** Private Execute_Quote
**URL:** https://docs.derive.xyz/reference/private-execute_quote

### Method Name

#### `private/execute_quote`

Executes a quote.  
Required minimum session key permission level is `admin`

### Parameters

|  |
| --- |
| **direction**string required  Quote direction, `buy` means trading each leg at its direction, `sell` means trading each leg in the opposite direction. enum  `buy` `sell` |
| **max\_fee**string required  Max fee ($ for the full trade). Request will be rejected if the supplied max fee is below the estimated fee for this trade. |
| **nonce**integer required  Unique nonce defined as a concatenated `UTC timestamp in ms` and `random number up to 6 digits` (e.g. 1695836058725001, where 001 is the random number) |
| **quote\_id**string required  Quote ID to execute against |
| **rfq\_id**string required  RFQ ID to execute (must be sent by `subaccount_id`) |
| **signature**string required  Ethereum signature of the quote |
| **signature\_expiry\_sec**integer required  Unix timestamp in seconds. Expiry MUST be at least 310 seconds from now. Once time till signature expiry reaches 300 seconds, the quote will be considered expired. This buffer is meant to ensure the trade can settle on chain in case of a blockchain congestion. |
| **signer**string required  Owner wallet address or registered session key that signed the quote |
| **subaccount\_id**integer required  Subaccount ID |
| **legs**array of objects required  Quote legs |
| legs[].**amount**string required  Amount in units of the base |
| legs[].**direction**string required  Leg direction enum  `buy` `sell` |
| legs[].**instrument\_name**string required  Instrument name |
| legs[].**price**string required  Leg price |
|  |
| **label**string  Optional user-defined label for the quote |

### Response

|  |
| --- |
| **id**string or integer required |
| **result**object required |
| result.**cancel\_reason**string required  Cancel reason, if any enum  `user_request` `insufficient_margin` `signed_max_fee_too_low` `mmp_trigger` `cancel_on_disconnect` `session_key_deregistered` `subaccount_withdrawn` `rfq_no_longer_open` `compliance` |
| result.**creation\_timestamp**integer required  Creation timestamp in ms since Unix epoch |
| result.**direction**string required  Quote direction enum  `buy` `sell` |
| result.**fee**string required  Fee paid for this quote (if executed) |
| result.**is\_transfer**boolean required  Whether the order was generated through `private/transfer_position` |
| result.**label**string required  User-defined label, if any |
| result.**last\_update\_timestamp**integer required  Last update timestamp in ms since Unix epoch |
| result.**legs\_hash**string required  Hash of the legs of the best quote to be signed by the taker. |
| result.**liquidity\_role**string required  Liquidity role enum  `maker` `taker` |
| result.**max\_fee**string required  Signed max fee |
| result.**mmp**boolean required  Whether the quote is tagged for market maker protections (default false) |
| result.**nonce**integer required  Nonce |
| result.**quote\_id**string required  Quote ID |
| result.**rfq\_id**string required  RFQ ID |
| result.**signature**string required  Ethereum signature of the quote |
| result.**signature\_expiry\_sec**integer required  Unix timestamp in seconds |
| result.**signer**string required  Owner wallet address or registered session key that signed the quote |
| result.**status**string required  Status enum  `open` `filled` `cancelled` `expired` |
| result.**subaccount\_id**integer required  Subaccount ID |
| result.**tx\_hash**string or null required  Blockchain transaction hash (only for executed quotes) |
| result.**tx\_status**string or null required  Blockchain transaction status (only for executed quotes) enum  `requested` `pending` `settled` `reverted` `ignored` `timed_out` |
| result.**legs**array of objects required  Quote legs |
| result.legs[].**amount**string required  Amount in units of the base |
| result.legs[].**direction**string required  Leg direction enum  `buy` `sell` |
| result.legs[].**instrument\_name**string required  Instrument name |
| result.legs[].**price**string required  Leg price |

### Example

ShellJavaScriptPython

```
{request_example_shell}

```

```
{request_example_javascript}

```

```
{request_example_python}

```

> The above command returns JSON structured like this:

JSON

```
{response_example_json}

```

---

## private-expired_and_cancelled_history

**Title:** Private Expired_And_Cancelled_History
**URL:** https://docs.derive.xyz/reference/private-expired_and_cancelled_history

### Method Name

#### `private/expired_and_cancelled_history`

Generate a list of URLs to retrieve archived orders  
Required minimum session key permission level is `read_only`

### Parameters

|  |
| --- |
| **end\_timestamp**integer required  End Unix timestamp |
| **expiry**integer required  Expiry of download link in seconds. Maximum of 604800. |
| **start\_timestamp**integer required  Start Unix timestamp |
| **subaccount\_id**integer required  Subaccount to download data for |
| **wallet**string required  Wallet to download data for |

### Response

|  |
| --- |
| **id**string or integer required |
| **result**object required |
| result.**presigned\_urls**array of strings required  List of presigned URLs to the snapshots |

### Example

ShellJavaScriptPython

```
{request_example_shell}

```

```
{request_example_javascript}

```

```
{request_example_python}

```

> The above command returns JSON structured like this:

JSON

```
{response_example_json}

```

---

## private-get_account

**Title:** Private Get_Account
**URL:** https://docs.derive.xyz/reference/private-get_account

### Method Name

#### `private/get_account`

Account details getter  
Required minimum session key permission level is `read_only`

### Parameters

|  |
| --- |
| **wallet**string required  Ethereum wallet address of account |

### Response

|  |
| --- |
| **id**string or integer required |
| **result**object required |
| result.**cancel\_on\_disconnect**boolean required  Whether cancel on disconnect is enabled for the account |
| result.**is\_rfq\_maker**boolean required  Whether account allowed to market make RFQs |
| result.**wallet**string required  Ethereum wallet address |
| result.**websocket\_matching\_tps**integer required  Max transactions per second for matching requests over websocket (see `Rate Limiting` in docs) |
| result.**websocket\_non\_matching\_tps**integer required  Max transactions per second for non-matching requests over websocket (see `Rate Limiting` in docs) |
| result.**websocket\_option\_tps**integer required  Max transactions per second for EACH option instrument over websocket (see `Rate Limiting` in docs) |
| result.**websocket\_perp\_tps**integer required  Max transactions per second for EACH perp instrument over websocket (see `Rate Limiting` in docs) |
| result.**fee\_info**object required  Fee information for the account |
| result.fee\_info.**base\_fee\_discount**string required  Base fee discount |
| result.fee\_info.**option\_maker\_fee**string or null required  Option maker fee - uses default instrument fee rate if None |
| result.fee\_info.**option\_taker\_fee**string or null required  Option taker fee - uses default instrument fee rate if None |
| result.fee\_info.**perp\_maker\_fee**string or null required  Perp maker fee - uses default instrument fee rate if None |
| result.fee\_info.**perp\_taker\_fee**string or null required  Perp taker fee - uses default instrument fee rate if None |
| result.fee\_info.**rfq\_maker\_discount**string required  RFQ maker fee discount |
| result.fee\_info.**rfq\_taker\_discount**string required  RFQ taker fee discount |
| result.fee\_info.**spot\_maker\_fee**string or null required  Spot maker fee - uses default instrument fee rate if None |
| result.fee\_info.**spot\_taker\_fee**string or null required  Spot taker fee - uses default instrument fee rate if None |
|  |
| result.**per\_endpoint\_tps**object required  If a particular endpoint has a different max TPS, it will be specified here |
| result.**subaccount\_ids**array of integers required  List of subaccount\_ids owned by the wallet in `SubAccounts.sol` |

### Example

ShellJavaScriptPython

```
{request_example_shell}

```

```
{request_example_javascript}

```

```
{request_example_python}

```

> The above command returns JSON structured like this:

JSON

```
{response_example_json}

```

---

## private-get_all_portfolios

**Title:** Private Get_All_Portfolios
**URL:** https://docs.derive.xyz/reference/private-get_all_portfolios

### Method Name

#### `private/get_all_portfolios`

Get all portfolios of a wallet  
Required minimum session key permission level is `read_only`

### Parameters

|  |
| --- |
| **wallet**string required  Wallet address |

### Response

|  |
| --- |
| **id**string or integer required |
| **result**array of objects required |
| result[].**collaterals\_initial\_margin**string required  Total initial margin credit contributed by collaterals |
| result[].**collaterals\_maintenance\_margin**string required  Total maintenance margin credit contributed by collaterals |
| result[].**collaterals\_value**string required  Total mark-to-market value of all collaterals |
| result[].**currency**string required  Currency of subaccount |
| result[].**initial\_margin**string required  Total initial margin requirement of all positions and collaterals.Trades will be rejected if this value falls below zero after the trade. |
| result[].**is\_under\_liquidation**boolean required  Whether the subaccount is undergoing a liquidation auction |
| result[].**label**string required  User defined label |
| result[].**maintenance\_margin**string required  Total maintenance margin requirement of all positions and collaterals.If this value falls below zero, the subaccount will be flagged for liquidation. |
| result[].**margin\_type**string required  Margin type of subaccount (`PM` (Portfolio Margin) or `SM` (Standard Margin)) enum  `PM` `SM` `PM2` |
| result[].**open\_orders\_margin**string required  Total margin requirement of all open orders.Orders will be rejected if this value plus initial margin are below zero after the order. |
| result[].**positions\_initial\_margin**string required  Total initial margin requirement of all positions |
| result[].**positions\_maintenance\_margin**string required  Total maintenance margin requirement of all positions |
| result[].**positions\_value**string required  Total mark-to-market value of all positions |
| result[].**projected\_margin\_change**string required  Projected change in maintenance margin requirement between now and projected margin at 8:01 UTC. If this value plus current maintenance margin ise below zero, the account is at risk of being flagged for liquidation right after the upcoming expiry. |
| result[].**subaccount\_id**integer required  Subaccount\_id |
| result[].**subaccount\_value**string required  Total mark-to-market value of all positions and collaterals |
| result[].**collaterals**array of objects required  All collaterals that count towards margin of subaccount |
| result[].collaterals[].**amount**string required  Asset amount of given collateral |
| result[].collaterals[].**amount\_step**string required  Minimum amount step for the collateral |
| result[].collaterals[].**asset\_name**string required  Asset name |
| result[].collaterals[].**asset\_type**string required  Type of asset collateral (currently always `erc20`) enum  `erc20` `option` `perp` |
| result[].collaterals[].**average\_price**string required  Average price of the collateral, 0 for USDC. |
| result[].collaterals[].**average\_price\_excl\_fees**string required  Average price of whole position excluding fees |
| result[].collaterals[].**creation\_timestamp**integer required  Timestamp of when the position was opened (in ms since Unix epoch) |
| result[].collaterals[].**cumulative\_interest**string required  Cumulative interest earned on supplying collateral or paid for borrowing |
| result[].collaterals[].**currency**string required  Underlying currency of asset (`ETH`, `BTC`, etc) |
| result[].collaterals[].**initial\_margin**string required  USD value of collateral that contributes to initial margin |
| result[].collaterals[].**maintenance\_margin**string required  USD value of collateral that contributes to maintenance margin |
| result[].collaterals[].**mark\_price**string required  Current mark price of the asset |
| result[].collaterals[].**mark\_value**string required  USD value of the collateral (amount \* mark price) |
| result[].collaterals[].**open\_orders\_margin**string required  USD margin requirement for all open orders for this asset / instrument |
| result[].collaterals[].**pending\_interest**string required  Portion of interest that has not yet been settled on-chain. This number is added to the portfolio value for margin calculations purposes. |
| result[].collaterals[].**realized\_pnl**string required  Realized trading profit or loss of the collateral, 0 for USDC. |
| result[].collaterals[].**realized\_pnl\_excl\_fees**string required  Realized trading profit or loss of the position excluding fees |
| result[].collaterals[].**total\_fees**string required  Total fees paid for opening and changing the position |
| result[].collaterals[].**unrealized\_pnl**string required  Unrealized trading profit or loss of the collateral, 0 for USDC. |
| result[].collaterals[].**unrealized\_pnl\_excl\_fees**string required  Unrealized trading profit or loss of the position excluding fees |
|  |
| result[].**open\_orders**array of objects required  All open orders of subaccount |
| result[].open\_orders[].**amount**string required  Order amount in units of the base |
| result[].open\_orders[].**average\_price**string required  Average fill price |
| result[].open\_orders[].**cancel\_reason**string required  If cancelled, reason behind order cancellation enum  `user_request` `mmp_trigger` `insufficient_margin` `signed_max_fee_too_low` `cancel_on_disconnect` `ioc_or_market_partial_fill` `session_key_deregistered` `subaccount_withdrawn` `compliance` `trigger_failed` `validation_failed` |
| result[].open\_orders[].**creation\_timestamp**integer required  Creation timestamp (in ms since Unix epoch) |
| result[].open\_orders[].**direction**string required  Order direction enum  `buy` `sell` |
| result[].open\_orders[].**filled\_amount**string required  Total filled amount for the order |
| result[].open\_orders[].**instrument\_name**string required  Instrument name |
| result[].open\_orders[].**is\_transfer**boolean required  Whether the order was generated through `private/transfer_position` |
| result[].open\_orders[].**label**string required  Optional user-defined label for the order |
| result[].open\_orders[].**last\_update\_timestamp**integer required  Last update timestamp (in ms since Unix epoch) |
| result[].open\_orders[].**limit\_price**string required  Limit price in quote currency |
| result[].open\_orders[].**max\_fee**string required  Max fee in units of the quote currency |
| result[].open\_orders[].**mmp**boolean required  Whether the order is tagged for market maker protections |
| result[].open\_orders[].**nonce**integer required  Unique nonce defined as  (e.g. 1695836058725001, where 001 is the random number) |
| result[].open\_orders[].**order\_fee**string required  Total order fee paid so far |
| result[].open\_orders[].**order\_id**string required  Order ID |
| result[].open\_orders[].**order\_status**string required  Order status enum  `open` `filled` `cancelled` `expired` `untriggered` |
| result[].open\_orders[].**order\_type**string required  Order type enum  `limit` `market` |
| result[].open\_orders[].**quote\_id**string or null required  Quote ID if the trade was executed via RFQ |
| result[].open\_orders[].**signature**string required  Ethereum signature of the order |
| result[].open\_orders[].**signature\_expiry\_sec**integer required  Signature expiry timestamp |
| result[].open\_orders[].**signer**string required  Owner wallet address or registered session key that signed order |
| result[].open\_orders[].**subaccount\_id**integer required  Subaccount ID |
| result[].open\_orders[].**time\_in\_force**string required  Time in force enum  `gtc` `post_only` `fok` `ioc` |
| result[].open\_orders[].**replaced\_order\_id**string or null  If replaced, ID of the order that was replaced |
| result[].open\_orders[].**trigger\_price**string or null  (Required for trigger orders) Index or Market price to trigger order at |
| result[].open\_orders[].**trigger\_price\_type**string or null  (Required for trigger orders) Trigger with Index or Mark Price enum  `mark` `index` |
| result[].open\_orders[].**trigger\_reject\_message**string or null  (Required for trigger orders) Error message if error occured during trigger |
| result[].open\_orders[].**trigger\_type**string or null  (Required for trigger orders) Stop-loss or Take-profit. enum  `stoploss` `takeprofit` |
|  |
| result[].**positions**array of objects required  All active positions of subaccount |
| result[].positions[].**amount**string required  Position amount held by subaccount |
| result[].positions[].**amount\_step**string required  Minimum amount step for the position |
| result[].positions[].**average\_price**string required  Average price of whole position |
| result[].positions[].**average\_price\_excl\_fees**string required  Average price of whole position excluding fees |
| result[].positions[].**creation\_timestamp**integer required  Timestamp of when the position was opened (in ms since Unix epoch) |
| result[].positions[].**cumulative\_funding**string required  Cumulative funding for the position (only for perpetuals). |
| result[].positions[].**delta**string required  Asset delta (w.r.t. forward price for options, `1.0` for perps) |
| result[].positions[].**gamma**string required  Asset gamma (zero for non-options) |
| result[].positions[].**index\_price**string required  Current index (oracle) price for position's currency |
| result[].positions[].**initial\_margin**string required  USD initial margin requirement for this position |
| result[].positions[].**instrument\_name**string required  Instrument name (same as the base Asset name) |
| result[].positions[].**instrument\_type**string required  `erc20`, `option`, or `perp` enum  `erc20` `option` `perp` |
| result[].positions[].**leverage**string or null required  Only for perps. Leverage of the position, defined as `abs(notional) / collateral net of options margin` |
| result[].positions[].**liquidation\_price**string or null required  Index price at which position will be liquidated |
| result[].positions[].**maintenance\_margin**string required  USD maintenance margin requirement for this position |
| result[].positions[].**mark\_price**string required  Current mark price for position's instrument |
| result[].positions[].**mark\_value**string required  USD value of the position; this represents how much USD can be recieved by fully closing the position at the current oracle price |
| result[].positions[].**net\_settlements**string required  Net amount of USD from position settlements that has been paid to the user's subaccount. This number is subtracted from the portfolio value for margin calculations purposes. Positive values mean the user has recieved USD from settlements, or is awaiting settlement of USD losses. Negative values mean the user has paid USD for settlements, or is awaiting settlement of USD gains. |
| result[].positions[].**open\_orders\_margin**string required  USD margin requirement for all open orders for this asset / instrument |
| result[].positions[].**pending\_funding**string required  A portion of funding payments that has not yet been settled into cash balance (only for perpetuals). This number is added to the portfolio value for margin calculations purposes. |
| result[].positions[].**realized\_pnl**string required  Realized trading profit or loss of the position. |
| result[].positions[].**realized\_pnl\_excl\_fees**string required  Realized trading profit or loss of the position excluding fees |
| result[].positions[].**theta**string required  Asset theta (zero for non-options) |
| result[].positions[].**total\_fees**string required  Total fees paid for opening and changing the position |
| result[].positions[].**unrealized\_pnl**string required  Unrealized trading profit or loss of the position. |
| result[].positions[].**unrealized\_pnl\_excl\_fees**string required  Unrealized trading profit or loss of the position excluding fees |
| result[].positions[].**vega**string required  Asset vega (zero for non-options) |

### Example

ShellJavaScriptPython

```
{request_example_shell}

```

```
{request_example_javascript}

```

```
{request_example_python}

```

> The above command returns JSON structured like this:

JSON

```
{response_example_json}

```

---

## private-get_collaterals

**Title:** Private Get_Collaterals
**URL:** https://docs.derive.xyz/reference/private-get_collaterals

### Method Name

#### `private/get_collaterals`

Get collaterals of a subaccount  
Required minimum session key permission level is `read_only`

### Parameters

|  |
| --- |
| **subaccount\_id**integer required  Subaccount\_id |

### Response

|  |
| --- |
| **id**string or integer required |
| **result**object required |
| result.**subaccount\_id**integer required  Subaccount\_id |
| result.**collaterals**array of objects required  All collaterals that count towards margin of subaccount |
| result.collaterals[].**amount**string required  Asset amount of given collateral |
| result.collaterals[].**amount\_step**string required  Minimum amount step for the collateral |
| result.collaterals[].**asset\_name**string required  Asset name |
| result.collaterals[].**asset\_type**string required  Type of asset collateral (currently always `erc20`) enum  `erc20` `option` `perp` |
| result.collaterals[].**average\_price**string required  Average price of the collateral, 0 for USDC. |
| result.collaterals[].**average\_price\_excl\_fees**string required  Average price of whole position excluding fees |
| result.collaterals[].**creation\_timestamp**integer required  Timestamp of when the position was opened (in ms since Unix epoch) |
| result.collaterals[].**cumulative\_interest**string required  Cumulative interest earned on supplying collateral or paid for borrowing |
| result.collaterals[].**currency**string required  Underlying currency of asset (`ETH`, `BTC`, etc) |
| result.collaterals[].**initial\_margin**string required  USD value of collateral that contributes to initial margin |
| result.collaterals[].**maintenance\_margin**string required  USD value of collateral that contributes to maintenance margin |
| result.collaterals[].**mark\_price**string required  Current mark price of the asset |
| result.collaterals[].**mark\_value**string required  USD value of the collateral (amount \* mark price) |
| result.collaterals[].**open\_orders\_margin**string required  USD margin requirement for all open orders for this asset / instrument |
| result.collaterals[].**pending\_interest**string required  Portion of interest that has not yet been settled on-chain. This number is added to the portfolio value for margin calculations purposes. |
| result.collaterals[].**realized\_pnl**string required  Realized trading profit or loss of the collateral, 0 for USDC. |
| result.collaterals[].**realized\_pnl\_excl\_fees**string required  Realized trading profit or loss of the position excluding fees |
| result.collaterals[].**total\_fees**string required  Total fees paid for opening and changing the position |
| result.collaterals[].**unrealized\_pnl**string required  Unrealized trading profit or loss of the collateral, 0 for USDC. |
| result.collaterals[].**unrealized\_pnl\_excl\_fees**string required  Unrealized trading profit or loss of the position excluding fees |

### Example

ShellJavaScriptPython

```
{request_example_shell}

```

```
{request_example_javascript}

```

```
{request_example_python}

```

> The above command returns JSON structured like this:

JSON

```
{response_example_json}

```

---

## private-get_deposit_history

**Title:** Private Get_Deposit_History
**URL:** https://docs.derive.xyz/reference/private-get_deposit_history

### Method Name

#### `private/get_deposit_history`

Get subaccount deposit history.  
Required minimum session key permission level is `read_only`

### Parameters

|  |
| --- |
| **subaccount\_id**integer required  Subaccount id |
| **end\_timestamp**integer  End timestamp of the event history (default current time) |
| **start\_timestamp**integer  Start timestamp of the event history (default 0) |

### Response

|  |
| --- |
| **id**string or integer required |
| **result**object required |
| result.**events**array of objects required  List of deposit payments |
| result.events[].**amount**string required  Amount deposited by the subaccount |
| result.events[].**asset**string required  Asset deposited |
| result.events[].**error\_log**object or null required  If failed, error log for reason |
| result.events[].**timestamp**integer required  Timestamp of the deposit (in ms since UNIX epoch) |
| result.events[].**transaction\_id**string required  Transaction ID |
| result.events[].**tx\_hash**string required  Hash of the transaction that deposited the funds |
| result.events[].**tx\_status**string required  Status of the transaction that deposited the funds enum  `requested` `pending` `settled` `reverted` `ignored` `timed_out` |

### Example

ShellJavaScriptPython

```
{request_example_shell}

```

```
{request_example_javascript}

```

```
{request_example_python}

```

> The above command returns JSON structured like this:

JSON

```
{response_example_json}

```

---

## private-get_erc20_transfer_history

**Title:** Private Get_Erc20_Transfer_History
**URL:** https://docs.derive.xyz/reference/private-get_erc20_transfer_history

### Method Name

#### `private/get_erc20_transfer_history`

Get subaccount erc20 transfer history.  
  
Position transfers (e.g. options or perps) are treated as trades. Use `private/get_trade_history` for position transfer history.  
Required minimum session key permission level is `read_only`

### Parameters

|  |
| --- |
| **subaccount\_id**integer required  Subaccount id |
| **end\_timestamp**integer  End timestamp of the event history (default current time) |
| **start\_timestamp**integer  Start timestamp of the event history (default 0) |

### Response

|  |
| --- |
| **id**string or integer required |
| **result**object required |
| result.**events**array of objects required  List of erc20 transfers |
| result.events[].**amount**string required  Amount withdrawn by the subaccount |
| result.events[].**asset**string required  Asset withdrawn |
| result.events[].**counterparty\_subaccount\_id**integer required  Recipient or sender subaccount\_id of transfer |
| result.events[].**is\_outgoing**boolean required  True if the transfer was initiated by the subaccount, False otherwise |
| result.events[].**timestamp**integer required  Timestamp of the transfer (in ms since UNIX epoch) |
| result.events[].**tx\_hash**string required  Hash of the transaction that withdrew the funds |

### Example

ShellJavaScriptPython

```
{request_example_shell}

```

```
{request_example_javascript}

```

```
{request_example_python}

```

> The above command returns JSON structured like this:

JSON

```
{response_example_json}

```

---

## private-get_funding_history

**Title:** Private Get_Funding_History
**URL:** https://docs.derive.xyz/reference/private-get_funding_history

### Method Name

#### `private/get_funding_history`

Get subaccount funding history.  
  
DB: read replica  
Required minimum session key permission level is `read_only`

### Parameters

|  |
| --- |
| **subaccount\_id**integer required  Subaccount id |
| **end\_timestamp**integer  End timestamp of the event history (default current time) |
| **instrument\_name**string  Instrument name (returns history for all perpetuals if not provided) |
| **page**integer  Page number of results to return (default 1, returns last if above `num_pages`) |
| **page\_size**integer  Number of results per page (default 100, max 1000) |
| **start\_timestamp**integer  Start timestamp of the event history (default 0) |

### Response

|  |
| --- |
| **id**string or integer required |
| **result**object required |
| result.**events**array of objects required  List of funding payments |
| result.events[].**funding**string required  Dollar funding paid (if negative) or received (if positive) by the subaccount |
| result.events[].**instrument\_name**string required  Instrument name |
| result.events[].**pnl**string required  Cashflow from the perp PnL settlement |
| result.events[].**timestamp**integer required  Timestamp of the funding payment (in ms since UNIX epoch) |
|  |
| result.**pagination**object required  Pagination information |
| result.pagination.**count**integer required  Total number of items, across all pages |
| result.pagination.**num\_pages**integer required  Number of pages |

### Example

ShellJavaScriptPython

```
{request_example_shell}

```

```
{request_example_javascript}

```

```
{request_example_python}

```

> The above command returns JSON structured like this:

JSON

```
{response_example_json}

```

---

## private-get_interest_history

**Title:** Private Get_Interest_History
**URL:** https://docs.derive.xyz/reference/private-get_interest_history

### Method Name

#### `private/get_interest_history`

Get subaccount interest payment history.  
Required minimum session key permission level is `read_only`

### Parameters

|  |
| --- |
| **subaccount\_id**integer required  Subaccount id |
| **end\_timestamp**integer  End timestamp of the event history (default current time) |
| **start\_timestamp**integer  Start timestamp of the event history (default 0) |

### Response

|  |
| --- |
| **id**string or integer required |
| **result**object required |
| result.**events**array of objects required  List of interest payments |
| result.events[].**interest**string required  Dollar interest paid (if negative) or received (if positive) by the subaccount |
| result.events[].**timestamp**integer required  Timestamp of the interest payment (in ms since UNIX epoch) |

### Example

ShellJavaScriptPython

```
{request_example_shell}

```

```
{request_example_javascript}

```

```
{request_example_python}

```

> The above command returns JSON structured like this:

JSON

```
{response_example_json}

```

---

## private-get_liquidation_history

**Title:** Private Get_Liquidation_History
**URL:** https://docs.derive.xyz/reference/private-get_liquidation_history

### Method Name

#### `private/get_liquidation_history`

Required minimum session key permission level is `read_only`

### Parameters

|  |
| --- |
| **subaccount\_id**integer required  Subaccount id |
| **end\_timestamp**integer  End timestamp of the event history (default current time) |
| **start\_timestamp**integer  Start timestamp of the event history (default 0) |

### Response

|  |
| --- |
| **id**string or integer required |
| **result**array of objects required |
| result[].**auction\_id**string required  Unique ID of the auction |
| result[].**auction\_type**string required  Type of auction enum  `solvent` `insolvent` |
| result[].**end\_timestamp**integer or null required  Timestamp of the auction end (in ms since UNIX epoch), or `null` if not ended |
| result[].**fee**string required  Fee paid by the subaccount |
| result[].**start\_timestamp**integer required  Timestamp of the auction start (in ms since UNIX epoch) |
| result[].**subaccount\_id**integer required  Liquidated subaccount ID |
| result[].**tx\_hash**string required  Hash of the transaction that started the auction |
| result[].**bids**array of objects required  List of auction bid events |
| result[].bids[].**cash\_received**string required  Cash received by the subaccount for the liquidation. For the liquidated accounts this is the amount the liquidator paid to buy out the percentage of the portfolio. For the liquidator account, this is the amount they received from the security module (if positive) or the amount they paid for the bid (if negative) |
| result[].bids[].**discount\_pnl**string required  Realized PnL due to liquidating or being liquidated at a discount to mark portfolio value |
| result[].bids[].**percent\_liquidated**string required  Percent of the subaccount that was liquidated |
| result[].bids[].**realized\_pnl**string required  Realized PnL of the auction bid, assuming positions are closed at mark price at the time of the liquidation |
| result[].bids[].**realized\_pnl\_excl\_fees**string required  Realized PnL of the auction bid, excluding fees from total cost basis, assuming positions are closed at mark price at the time of the liquidation |
| result[].bids[].**timestamp**integer required  Timestamp of the bid (in ms since UNIX epoch) |
| result[].bids[].**tx\_hash**string required  Hash of the bid transaction |
| result[].bids[].**amounts\_liquidated**object required  Amounts of each asset that were closed |
| result[].bids[].**positions\_realized\_pnl**object required  Realized PnL of each position that was closed |
| result[].bids[].**positions\_realized\_pnl\_excl\_fees**object required  Realized PnL of each position that was closed, excluding fees from total cost basis |

### Example

ShellJavaScriptPython

```
{request_example_shell}

```

```
{request_example_javascript}

```

```
{request_example_python}

```

> The above command returns JSON structured like this:

JSON

```
{response_example_json}

```

---

## private-get_liquidator_history

**Title:** Private Get_Liquidator_History
**URL:** https://docs.derive.xyz/reference/private-get_liquidator_history

### Method Name

#### `private/get_liquidator_history`

Returns a paginated history of auctions that the subaccount has participated in as a liquidator.  
Required minimum session key permission level is `read_only`

### Parameters

|  |
| --- |
| **subaccount\_id**integer required  Subaccount ID |
| **end\_timestamp**integer  End timestamp of the event history (default current time) |
| **page**integer  Page number of results to return (default 1, returns last if above `num_pages`) |
| **page\_size**integer  Number of results per page (default 100, max 1000) |
| **start\_timestamp**integer  Start timestamp of the event history (default 0) |

### Response

|  |
| --- |
| **id**string or integer required |
| **result**object required |
| result.**bids**array of objects required  List of auction bid events |
| result.bids[].**cash\_received**string required  Cash received by the subaccount for the liquidation. For the liquidated accounts this is the amount the liquidator paid to buy out the percentage of the portfolio. For the liquidator account, this is the amount they received from the security module (if positive) or the amount they paid for the bid (if negative) |
| result.bids[].**discount\_pnl**string required  Realized PnL due to liquidating or being liquidated at a discount to mark portfolio value |
| result.bids[].**percent\_liquidated**string required  Percent of the subaccount that was liquidated |
| result.bids[].**realized\_pnl**string required  Realized PnL of the auction bid, assuming positions are closed at mark price at the time of the liquidation |
| result.bids[].**realized\_pnl\_excl\_fees**string required  Realized PnL of the auction bid, excluding fees from total cost basis, assuming positions are closed at mark price at the time of the liquidation |
| result.bids[].**timestamp**integer required  Timestamp of the bid (in ms since UNIX epoch) |
| result.bids[].**tx\_hash**string required  Hash of the bid transaction |
| result.bids[].**amounts\_liquidated**object required  Amounts of each asset that were closed |
| result.bids[].**positions\_realized\_pnl**object required  Realized PnL of each position that was closed |
| result.bids[].**positions\_realized\_pnl\_excl\_fees**object required  Realized PnL of each position that was closed, excluding fees from total cost basis |
|  |
| result.**pagination**object required  Pagination info |
| result.pagination.**count**integer required  Total number of items, across all pages |
| result.pagination.**num\_pages**integer required  Number of pages |

### Example

ShellJavaScriptPython

```
{request_example_shell}

```

```
{request_example_javascript}

```

```
{request_example_python}

```

> The above command returns JSON structured like this:

JSON

```
{response_example_json}

```

---

## private-get_margin

**Title:** Private Get_Margin
**URL:** https://docs.derive.xyz/reference/private-get_margin

### Method Name

#### `private/get_margin`

Calculates margin for a given subaccount and (optionally) a simulated state change. Does not take into account  
open orders margin requirements.  
Required minimum session key permission level is `read_only`

### Parameters

|  |
| --- |
| **subaccount\_id**integer required  Subaccount\_id |
| **simulated\_collateral\_changes**array of objects  Optional, add collaterals to simulate deposits / withdrawals / spot trades |
| simulated\_collateral\_changes[].**amount**string required  Collateral amount to simulate |
| simulated\_collateral\_changes[].**asset\_name**string required  Collateral ERC20 asset name (e.g. ETH, USDC, WSTETH) |
|  |
| **simulated\_position\_changes**array of objects  Optional, add positions to simulate perp / option trades |
| simulated\_position\_changes[].**amount**string required  Position amount to simulate |
| simulated\_position\_changes[].**instrument\_name**string required  Instrument name |
| simulated\_position\_changes[].**entry\_price**string  Only for perps. Entry price to use in the simulation. Mark price is used if not provided. |

### Response

|  |
| --- |
| **id**string or integer required |
| **result**object required |
| result.**is\_valid\_trade**boolean required  True if trade passes margin requirement |
| result.**post\_initial\_margin**string required  Initial margin requirement post trade |
| result.**post\_maintenance\_margin**string required  Maintenance margin requirement post trade |
| result.**pre\_initial\_margin**string required  Initial margin requirement before trade |
| result.**pre\_maintenance\_margin**string required  Maintenance margin requirement before trade |
| result.**subaccount\_id**integer required  Subaccount\_id |

### Example

ShellJavaScriptPython

```
{request_example_shell}

```

```
{request_example_javascript}

```

```
{request_example_python}

```

> The above command returns JSON structured like this:

JSON

```
{response_example_json}

```

---

## private-get_mmp_config

**Title:** Private Get_Mmp_Config
**URL:** https://docs.derive.xyz/reference/private-get_mmp_config

### Method Name

#### `private/get_mmp_config`

Get the current mmp config for a subaccount (optionally filtered by currency)  
Required minimum session key permission level is `read_only`

### Parameters

|  |
| --- |
| **subaccount\_id**integer required  Subaccount\_id for which to get the config |
| **currency**string  Currency to get the config for. If not provided, returns all configs for the subaccount |

### Response

|  |
| --- |
| **id**string or integer required |
| **result**array of objects required |
| result[].**currency**string required  Currency of this mmp config |
| result[].**is\_frozen**boolean required  Whether the subaccount is currently frozen |
| result[].**mmp\_frozen\_time**integer required  Time interval in ms setting how long the subaccount is frozen after an mmp trigger, if 0 then a manual reset would be required via private/reset\_mmp |
| result[].**mmp\_interval**integer required  Time interval in ms over which the limits are monotored, if 0 then mmp is disabled |
| result[].**mmp\_unfreeze\_time**integer required  Timestamp in ms after which the subaccount will be unfrozen |
| result[].**subaccount\_id**integer required  Subaccount\_id for which to set the config |
| result[].**mmp\_amount\_limit**string  Maximum total order amount that can be traded within the mmp\_interval across all instruments of the provided currency. The amounts are not netted, so a filled bid of 1 and a filled ask of 2 would count as 3. Default: 0 (no limit) |
| result[].**mmp\_delta\_limit**string  Maximum total delta that can be traded within the mmp\_interval across all instruments of the provided currency. This quantity is netted, so a filled order with +1 delta and a filled order with -2 delta would count as -1 Default: 0 (no limit) |

### Example

ShellJavaScriptPython

```
{request_example_shell}

```

```
{request_example_javascript}

```

```
{request_example_python}

```

> The above command returns JSON structured like this:

JSON

```
{response_example_json}

```

---

## private-get_notifications

**Title:** Private Get_Notifications
**URL:** https://docs.derive.xyz/reference/private-get_notifications

### Method Name

#### `private/get_notifications`

Get the notifications related to a subaccount.  
Required minimum session key permission level is `read_only`

### Parameters

|  |
| --- |
| **page**integer  Page number of results to return |
| **page\_size**integer  Number of results per page (must be between 0-50) |
| **status**string  Status of the notification enum  `unseen` `seen` `hidden` |
| **subaccount\_id**integer  Subaccount\_id (must be set if wallet param is not set) |
| **type**array of strings  List of notification types to filter by |
| **wallet**string  Wallet address (if set, subaccount\_id ignored) |

### Response

|  |
| --- |
| **id**string or integer required |
| **result**object required |
| result.**notifications**array of objects required  Notification response |
| result.notifications[].**event**string required  The specific event leading to the notification. |
| result.notifications[].**id**integer required  The unique identifier for the notification. |
| result.notifications[].**status**string required  The status of the notification, indicating if it has been read, pending, or processed. |
| result.notifications[].**subaccount\_id**integer required  The subaccount\_id associated with the notification. |
| result.notifications[].**timestamp**integer required  The timestamp indicating when the notification was created or triggered. |
| result.notifications[].**event\_details**object required  A JSON-structured dictionary containing detailed data or context about the event. |
| result.notifications[].**transaction\_id**integer or null  The transaction id associated with the notification. |
| result.notifications[].**tx\_hash**string or null  The transaction hash associated with the notification. |
|  |
| result.**pagination**object required  Pagination info |
| result.pagination.**count**integer required  Total number of items, across all pages |
| result.pagination.**num\_pages**integer required  Number of pages |

### Example

ShellJavaScriptPython

```
{request_example_shell}

```

```
{request_example_javascript}

```

```
{request_example_python}

```

> The above command returns JSON structured like this:

JSON

```
{response_example_json}

```

---

## private-get_open_orders

**Title:** Private Get_Open_Orders
**URL:** https://docs.derive.xyz/reference/private-get_open_orders

### Method Name

#### `private/get_open_orders`

Get all open orders of a subacccount  
Required minimum session key permission level is `read_only`

### Parameters

|  |
| --- |
| **subaccount\_id**integer required  Subaccount\_id for which to get open orders |

### Response

|  |
| --- |
| **id**string or integer required |
| **result**object required |
| result.**subaccount\_id**integer required  Subaccount\_id for which to get open orders |
| result.**orders**array of objects required  List of open orders |
| result.orders[].**amount**string required  Order amount in units of the base |
| result.orders[].**average\_price**string required  Average fill price |
| result.orders[].**cancel\_reason**string required  If cancelled, reason behind order cancellation enum  `user_request` `mmp_trigger` `insufficient_margin` `signed_max_fee_too_low` `cancel_on_disconnect` `ioc_or_market_partial_fill` `session_key_deregistered` `subaccount_withdrawn` `compliance` `trigger_failed` `validation_failed` |
| result.orders[].**creation\_timestamp**integer required  Creation timestamp (in ms since Unix epoch) |
| result.orders[].**direction**string required  Order direction enum  `buy` `sell` |
| result.orders[].**filled\_amount**string required  Total filled amount for the order |
| result.orders[].**instrument\_name**string required  Instrument name |
| result.orders[].**is\_transfer**boolean required  Whether the order was generated through `private/transfer_position` |
| result.orders[].**label**string required  Optional user-defined label for the order |
| result.orders[].**last\_update\_timestamp**integer required  Last update timestamp (in ms since Unix epoch) |
| result.orders[].**limit\_price**string required  Limit price in quote currency |
| result.orders[].**max\_fee**string required  Max fee in units of the quote currency |
| result.orders[].**mmp**boolean required  Whether the order is tagged for market maker protections |
| result.orders[].**nonce**integer required  Unique nonce defined as  (e.g. 1695836058725001, where 001 is the random number) |
| result.orders[].**order\_fee**string required  Total order fee paid so far |
| result.orders[].**order\_id**string required  Order ID |
| result.orders[].**order\_status**string required  Order status enum  `open` `filled` `cancelled` `expired` `untriggered` |
| result.orders[].**order\_type**string required  Order type enum  `limit` `market` |
| result.orders[].**quote\_id**string or null required  Quote ID if the trade was executed via RFQ |
| result.orders[].**signature**string required  Ethereum signature of the order |
| result.orders[].**signature\_expiry\_sec**integer required  Signature expiry timestamp |
| result.orders[].**signer**string required  Owner wallet address or registered session key that signed order |
| result.orders[].**subaccount\_id**integer required  Subaccount ID |
| result.orders[].**time\_in\_force**string required  Time in force enum  `gtc` `post_only` `fok` `ioc` |
| result.orders[].**replaced\_order\_id**string or null  If replaced, ID of the order that was replaced |
| result.orders[].**trigger\_price**string or null  (Required for trigger orders) Index or Market price to trigger order at |
| result.orders[].**trigger\_price\_type**string or null  (Required for trigger orders) Trigger with Index or Mark Price enum  `mark` `index` |
| result.orders[].**trigger\_reject\_message**string or null  (Required for trigger orders) Error message if error occured during trigger |
| result.orders[].**trigger\_type**string or null  (Required for trigger orders) Stop-loss or Take-profit. enum  `stoploss` `takeprofit` |

### Example

ShellJavaScriptPython

```
{request_example_shell}

```

```
{request_example_javascript}

```

```
{request_example_python}

```

> The above command returns JSON structured like this:

JSON

```
{response_example_json}

```

---

## private-get_option_settlement_history

**Title:** Private Get_Option_Settlement_History
**URL:** https://docs.derive.xyz/reference/private-get_option_settlement_history

### Method Name

#### `private/get_option_settlement_history`

Get expired option settlement history for a subaccount  
Required minimum session key permission level is `read_only`

### Parameters

|  |
| --- |
| **subaccount\_id**integer required  Subaccount ID for which to get expired option settlement history |

### Response

|  |
| --- |
| **id**string or integer required |
| **result**object required |
| result.**subaccount\_id**integer required  Subaccount\_id for which to get expired option settlement history |
| result.**settlements**array of objects required  List of expired option settlements |
| result.settlements[].**amount**string required  Amount that was settled |
| result.settlements[].**expiry**integer required  Expiry timestamp of the option |
| result.settlements[].**instrument\_name**string required  Instrument name |
| result.settlements[].**option\_settlement\_pnl**string required  USD profit or loss from option settlements calculated as: settlement value - (average cost including fees x amount) |
| result.settlements[].**option\_settlement\_pnl\_excl\_fees**string required  USD profit or loss from option settlements calculated as: settlement value - (average price excluding fees x amount) |
| result.settlements[].**settlement\_price**string required  Price of option settlement |
| result.settlements[].**subaccount\_id**integer required  Subaccount ID of the settlement event |

### Example

ShellJavaScriptPython

```
{request_example_shell}

```

```
{request_example_javascript}

```

```
{request_example_python}

```

> The above command returns JSON structured like this:

JSON

```
{response_example_json}

```

---

## private-get_order

**Title:** Private Get_Order
**URL:** https://docs.derive.xyz/reference/private-get_order

### Method Name

#### `private/get_order`

Get state of an order by order id  
Required minimum session key permission level is `read_only`

### Parameters

|  |
| --- |
| **order\_id**string required  Order ID |
| **subaccount\_id**integer required  Subaccount ID |

### Response

|  |
| --- |
| **id**string or integer required |
| **result**object required |
| result.**amount**string required  Order amount in units of the base |
| result.**average\_price**string required  Average fill price |
| result.**cancel\_reason**string required  If cancelled, reason behind order cancellation enum  `user_request` `mmp_trigger` `insufficient_margin` `signed_max_fee_too_low` `cancel_on_disconnect` `ioc_or_market_partial_fill` `session_key_deregistered` `subaccount_withdrawn` `compliance` `trigger_failed` `validation_failed` |
| result.**creation\_timestamp**integer required  Creation timestamp (in ms since Unix epoch) |
| result.**direction**string required  Order direction enum  `buy` `sell` |
| result.**filled\_amount**string required  Total filled amount for the order |
| result.**instrument\_name**string required  Instrument name |
| result.**is\_transfer**boolean required  Whether the order was generated through `private/transfer_position` |
| result.**label**string required  Optional user-defined label for the order |
| result.**last\_update\_timestamp**integer required  Last update timestamp (in ms since Unix epoch) |
| result.**limit\_price**string required  Limit price in quote currency |
| result.**max\_fee**string required  Max fee in units of the quote currency |
| result.**mmp**boolean required  Whether the order is tagged for market maker protections |
| result.**nonce**integer required  Unique nonce defined as  (e.g. 1695836058725001, where 001 is the random number) |
| result.**order\_fee**string required  Total order fee paid so far |
| result.**order\_id**string required  Order ID |
| result.**order\_status**string required  Order status enum  `open` `filled` `cancelled` `expired` `untriggered` |
| result.**order\_type**string required  Order type enum  `limit` `market` |
| result.**quote\_id**string or null required  Quote ID if the trade was executed via RFQ |
| result.**signature**string required  Ethereum signature of the order |
| result.**signature\_expiry\_sec**integer required  Signature expiry timestamp |
| result.**signer**string required  Owner wallet address or registered session key that signed order |
| result.**subaccount\_id**integer required  Subaccount ID |
| result.**time\_in\_force**string required  Time in force enum  `gtc` `post_only` `fok` `ioc` |
| result.**replaced\_order\_id**string or null  If replaced, ID of the order that was replaced |
| result.**trigger\_price**string or null  (Required for trigger orders) Index or Market price to trigger order at |
| result.**trigger\_price\_type**string or null  (Required for trigger orders) Trigger with Index or Mark Price enum  `mark` `index` |
| result.**trigger\_reject\_message**string or null  (Required for trigger orders) Error message if error occured during trigger |
| result.**trigger\_type**string or null  (Required for trigger orders) Stop-loss or Take-profit. enum  `stoploss` `takeprofit` |

### Example

ShellJavaScriptPython

```
{request_example_shell}

```

```
{request_example_javascript}

```

```
{request_example_python}

```

> The above command returns JSON structured like this:

JSON

```
{response_example_json}

```

---

## private-get_order_history

**Title:** Private Get_Order_History
**URL:** https://docs.derive.xyz/reference/private-get_order_history

### Method Name

#### `private/get_order_history`

Get order history for a subaccount  
Required minimum session key permission level is `read_only`

### Parameters

|  |
| --- |
| **subaccount\_id**integer required  Subaccount\_id for which to get order history |
| **page**integer  Page number of results to return (default 1, returns last if above `num_pages`) |
| **page\_size**integer  Number of results per page (default 100, max 1000) |

### Response

|  |
| --- |
| **id**string or integer required |
| **result**object required |
| result.**subaccount\_id**integer required  Subaccount\_id for which to get open orders |
| result.**orders**array of objects required  List of open orders |
| result.orders[].**amount**string required  Order amount in units of the base |
| result.orders[].**average\_price**string required  Average fill price |
| result.orders[].**cancel\_reason**string required  If cancelled, reason behind order cancellation enum  `user_request` `mmp_trigger` `insufficient_margin` `signed_max_fee_too_low` `cancel_on_disconnect` `ioc_or_market_partial_fill` `session_key_deregistered` `subaccount_withdrawn` `compliance` `trigger_failed` `validation_failed` |
| result.orders[].**creation\_timestamp**integer required  Creation timestamp (in ms since Unix epoch) |
| result.orders[].**direction**string required  Order direction enum  `buy` `sell` |
| result.orders[].**filled\_amount**string required  Total filled amount for the order |
| result.orders[].**instrument\_name**string required  Instrument name |
| result.orders[].**is\_transfer**boolean required  Whether the order was generated through `private/transfer_position` |
| result.orders[].**label**string required  Optional user-defined label for the order |
| result.orders[].**last\_update\_timestamp**integer required  Last update timestamp (in ms since Unix epoch) |
| result.orders[].**limit\_price**string required  Limit price in quote currency |
| result.orders[].**max\_fee**string required  Max fee in units of the quote currency |
| result.orders[].**mmp**boolean required  Whether the order is tagged for market maker protections |
| result.orders[].**nonce**integer required  Unique nonce defined as  (e.g. 1695836058725001, where 001 is the random number) |
| result.orders[].**order\_fee**string required  Total order fee paid so far |
| result.orders[].**order\_id**string required  Order ID |
| result.orders[].**order\_status**string required  Order status enum  `open` `filled` `cancelled` `expired` `untriggered` |
| result.orders[].**order\_type**string required  Order type enum  `limit` `market` |
| result.orders[].**quote\_id**string or null required  Quote ID if the trade was executed via RFQ |
| result.orders[].**signature**string required  Ethereum signature of the order |
| result.orders[].**signature\_expiry\_sec**integer required  Signature expiry timestamp |
| result.orders[].**signer**string required  Owner wallet address or registered session key that signed order |
| result.orders[].**subaccount\_id**integer required  Subaccount ID |
| result.orders[].**time\_in\_force**string required  Time in force enum  `gtc` `post_only` `fok` `ioc` |
| result.orders[].**replaced\_order\_id**string or null  If replaced, ID of the order that was replaced |
| result.orders[].**trigger\_price**string or null  (Required for trigger orders) Index or Market price to trigger order at |
| result.orders[].**trigger\_price\_type**string or null  (Required for trigger orders) Trigger with Index or Mark Price enum  `mark` `index` |
| result.orders[].**trigger\_reject\_message**string or null  (Required for trigger orders) Error message if error occured during trigger |
| result.orders[].**trigger\_type**string or null  (Required for trigger orders) Stop-loss or Take-profit. enum  `stoploss` `takeprofit` |
|  |
| result.**pagination**object required  Pagination info |
| result.pagination.**count**integer required  Total number of items, across all pages |
| result.pagination.**num\_pages**integer required  Number of pages |

### Example

ShellJavaScriptPython

```
{request_example_shell}

```

```
{request_example_javascript}

```

```
{request_example_python}

```

> The above command returns JSON structured like this:

JSON

```
{response_example_json}

```

---

## private-get_orders

**Title:** Private Get_Orders
**URL:** https://docs.derive.xyz/reference/private-get_orders

### Method Name

#### `private/get_orders`

Get orders for a subaccount, with optional filtering.  
Required minimum session key permission level is `read_only`

### Parameters

|  |
| --- |
| **subaccount\_id**integer required  Subaccount\_id for which to get open orders |
| **instrument\_name**string  Filter by instrument name |
| **label**string  Filter by label |
| **page**integer  Page number of results to return (default 1, returns last if above `num_pages`) |
| **page\_size**integer  Number of results per page (default 100, max 1000) |
| **status**string  Filter by order status enum  `open` `filled` `cancelled` `expired` `untriggered` |

### Response

|  |
| --- |
| **id**string or integer required |
| **result**object required |
| result.**subaccount\_id**integer required  Subaccount\_id for which to get open orders |
| result.**orders**array of objects required  List of open orders |
| result.orders[].**amount**string required  Order amount in units of the base |
| result.orders[].**average\_price**string required  Average fill price |
| result.orders[].**cancel\_reason**string required  If cancelled, reason behind order cancellation enum  `user_request` `mmp_trigger` `insufficient_margin` `signed_max_fee_too_low` `cancel_on_disconnect` `ioc_or_market_partial_fill` `session_key_deregistered` `subaccount_withdrawn` `compliance` `trigger_failed` `validation_failed` |
| result.orders[].**creation\_timestamp**integer required  Creation timestamp (in ms since Unix epoch) |
| result.orders[].**direction**string required  Order direction enum  `buy` `sell` |
| result.orders[].**filled\_amount**string required  Total filled amount for the order |
| result.orders[].**instrument\_name**string required  Instrument name |
| result.orders[].**is\_transfer**boolean required  Whether the order was generated through `private/transfer_position` |
| result.orders[].**label**string required  Optional user-defined label for the order |
| result.orders[].**last\_update\_timestamp**integer required  Last update timestamp (in ms since Unix epoch) |
| result.orders[].**limit\_price**string required  Limit price in quote currency |
| result.orders[].**max\_fee**string required  Max fee in units of the quote currency |
| result.orders[].**mmp**boolean required  Whether the order is tagged for market maker protections |
| result.orders[].**nonce**integer required  Unique nonce defined as  (e.g. 1695836058725001, where 001 is the random number) |
| result.orders[].**order\_fee**string required  Total order fee paid so far |
| result.orders[].**order\_id**string required  Order ID |
| result.orders[].**order\_status**string required  Order status enum  `open` `filled` `cancelled` `expired` `untriggered` |
| result.orders[].**order\_type**string required  Order type enum  `limit` `market` |
| result.orders[].**quote\_id**string or null required  Quote ID if the trade was executed via RFQ |
| result.orders[].**signature**string required  Ethereum signature of the order |
| result.orders[].**signature\_expiry\_sec**integer required  Signature expiry timestamp |
| result.orders[].**signer**string required  Owner wallet address or registered session key that signed order |
| result.orders[].**subaccount\_id**integer required  Subaccount ID |
| result.orders[].**time\_in\_force**string required  Time in force enum  `gtc` `post_only` `fok` `ioc` |
| result.orders[].**replaced\_order\_id**string or null  If replaced, ID of the order that was replaced |
| result.orders[].**trigger\_price**string or null  (Required for trigger orders) Index or Market price to trigger order at |
| result.orders[].**trigger\_price\_type**string or null  (Required for trigger orders) Trigger with Index or Mark Price enum  `mark` `index` |
| result.orders[].**trigger\_reject\_message**string or null  (Required for trigger orders) Error message if error occured during trigger |
| result.orders[].**trigger\_type**string or null  (Required for trigger orders) Stop-loss or Take-profit. enum  `stoploss` `takeprofit` |
|  |
| result.**pagination**object required  Pagination info |
| result.pagination.**count**integer required  Total number of items, across all pages |
| result.pagination.**num\_pages**integer required  Number of pages |

### Example

ShellJavaScriptPython

```
{request_example_shell}

```

```
{request_example_javascript}

```

```
{request_example_python}

```

> The above command returns JSON structured like this:

JSON

```
{response_example_json}

```

---

## private-get_positions

**Title:** Private Get_Positions
**URL:** https://docs.derive.xyz/reference/private-get_positions

### Method Name

#### `private/get_positions`

Get active positions of a subaccount  
Required minimum session key permission level is `read_only`

### Parameters

|  |
| --- |
| **subaccount\_id**integer required  Subaccount\_id |

### Response

|  |
| --- |
| **id**string or integer required |
| **result**object required |
| result.**subaccount\_id**integer required  Subaccount\_id |
| result.**positions**array of objects required  All active positions of subaccount |
| result.positions[].**amount**string required  Position amount held by subaccount |
| result.positions[].**amount\_step**string required  Minimum amount step for the position |
| result.positions[].**average\_price**string required  Average price of whole position |
| result.positions[].**average\_price\_excl\_fees**string required  Average price of whole position excluding fees |
| result.positions[].**creation\_timestamp**integer required  Timestamp of when the position was opened (in ms since Unix epoch) |
| result.positions[].**cumulative\_funding**string required  Cumulative funding for the position (only for perpetuals). |
| result.positions[].**delta**string required  Asset delta (w.r.t. forward price for options, `1.0` for perps) |
| result.positions[].**gamma**string required  Asset gamma (zero for non-options) |
| result.positions[].**index\_price**string required  Current index (oracle) price for position's currency |
| result.positions[].**initial\_margin**string required  USD initial margin requirement for this position |
| result.positions[].**instrument\_name**string required  Instrument name (same as the base Asset name) |
| result.positions[].**instrument\_type**string required  `erc20`, `option`, or `perp` enum  `erc20` `option` `perp` |
| result.positions[].**leverage**string or null required  Only for perps. Leverage of the position, defined as `abs(notional) / collateral net of options margin` |
| result.positions[].**liquidation\_price**string or null required  Index price at which position will be liquidated |
| result.positions[].**maintenance\_margin**string required  USD maintenance margin requirement for this position |
| result.positions[].**mark\_price**string required  Current mark price for position's instrument |
| result.positions[].**mark\_value**string required  USD value of the position; this represents how much USD can be recieved by fully closing the position at the current oracle price |
| result.positions[].**net\_settlements**string required  Net amount of USD from position settlements that has been paid to the user's subaccount. This number is subtracted from the portfolio value for margin calculations purposes. Positive values mean the user has recieved USD from settlements, or is awaiting settlement of USD losses. Negative values mean the user has paid USD for settlements, or is awaiting settlement of USD gains. |
| result.positions[].**open\_orders\_margin**string required  USD margin requirement for all open orders for this asset / instrument |
| result.positions[].**pending\_funding**string required  A portion of funding payments that has not yet been settled into cash balance (only for perpetuals). This number is added to the portfolio value for margin calculations purposes. |
| result.positions[].**realized\_pnl**string required  Realized trading profit or loss of the position. |
| result.positions[].**realized\_pnl\_excl\_fees**string required  Realized trading profit or loss of the position excluding fees |
| result.positions[].**theta**string required  Asset theta (zero for non-options) |
| result.positions[].**total\_fees**string required  Total fees paid for opening and changing the position |
| result.positions[].**unrealized\_pnl**string required  Unrealized trading profit or loss of the position. |
| result.positions[].**unrealized\_pnl\_excl\_fees**string required  Unrealized trading profit or loss of the position excluding fees |
| result.positions[].**vega**string required  Asset vega (zero for non-options) |

### Example

ShellJavaScriptPython

```
{request_example_shell}

```

```
{request_example_javascript}

```

```
{request_example_python}

```

> The above command returns JSON structured like this:

JSON

```
{response_example_json}

```

---

## private-get_quotes

**Title:** Private Get_Quotes
**URL:** https://docs.derive.xyz/reference/private-get_quotes

### Method Name

#### `private/get_quotes`

Retrieves a list of quotes matching filter criteria.  
Market makers can use this to get their open quotes, quote history, etc.  
Required minimum session key permission level is `read_only`

### Parameters

|  |
| --- |
| **subaccount\_id**integer required  Subaccount ID for auth purposes, returned data will be scoped to this subaccount. |
| **from\_timestamp**integer  Earliest timestamp to filter by (in ms since Unix epoch). If not provied, defaults to 0. |
| **page**integer  Page number of results to return (default 1, returns last if above `num_pages`) |
| **page\_size**integer  Number of results per page (default 100, max 1000) |
| **quote\_id**string  Quote ID filter, if applicable |
| **rfq\_id**string  RFQ ID filter, if applicable |
| **status**string  Quote status filter, if applicable enum  `open` `filled` `cancelled` `expired` |
| **to\_timestamp**integer  Latest timestamp to filter by (in ms since Unix epoch). If not provied, defaults to returning all data up to current time. |

### Response

|  |
| --- |
| **id**string or integer required |
| **result**object required |
| result.**pagination**object required  Pagination info |
| result.pagination.**count**integer required  Total number of items, across all pages |
| result.pagination.**num\_pages**integer required  Number of pages |
|  |
| result.**quotes**array of objects required  Quotes matching filter criteria |
| result.quotes[].**cancel\_reason**string required  Cancel reason, if any enum  `user_request` `insufficient_margin` `signed_max_fee_too_low` `mmp_trigger` `cancel_on_disconnect` `session_key_deregistered` `subaccount_withdrawn` `rfq_no_longer_open` `compliance` |
| result.quotes[].**creation\_timestamp**integer required  Creation timestamp in ms since Unix epoch |
| result.quotes[].**direction**string required  Quote direction enum  `buy` `sell` |
| result.quotes[].**fee**string required  Fee paid for this quote (if executed) |
| result.quotes[].**is\_transfer**boolean required  Whether the order was generated through `private/transfer_position` |
| result.quotes[].**label**string required  User-defined label, if any |
| result.quotes[].**last\_update\_timestamp**integer required  Last update timestamp in ms since Unix epoch |
| result.quotes[].**legs\_hash**string required  Hash of the legs of the best quote to be signed by the taker. |
| result.quotes[].**liquidity\_role**string required  Liquidity role enum  `maker` `taker` |
| result.quotes[].**max\_fee**string required  Signed max fee |
| result.quotes[].**mmp**boolean required  Whether the quote is tagged for market maker protections (default false) |
| result.quotes[].**nonce**integer required  Nonce |
| result.quotes[].**quote\_id**string required  Quote ID |
| result.quotes[].**rfq\_id**string required  RFQ ID |
| result.quotes[].**signature**string required  Ethereum signature of the quote |
| result.quotes[].**signature\_expiry\_sec**integer required  Unix timestamp in seconds |
| result.quotes[].**signer**string required  Owner wallet address or registered session key that signed the quote |
| result.quotes[].**status**string required  Status enum  `open` `filled` `cancelled` `expired` |
| result.quotes[].**subaccount\_id**integer required  Subaccount ID |
| result.quotes[].**tx\_hash**string or null required  Blockchain transaction hash (only for executed quotes) |
| result.quotes[].**tx\_status**string or null required  Blockchain transaction status (only for executed quotes) enum  `requested` `pending` `settled` `reverted` `ignored` `timed_out` |
| result.quotes[].**legs**array of objects required  Quote legs |
| result.quotes[].legs[].**amount**string required  Amount in units of the base |
| result.quotes[].legs[].**direction**string required  Leg direction enum  `buy` `sell` |
| result.quotes[].legs[].**instrument\_name**string required  Instrument name |
| result.quotes[].legs[].**price**string required  Leg price |

### Example

ShellJavaScriptPython

```
{request_example_shell}

```

```
{request_example_javascript}

```

```
{request_example_python}

```

> The above command returns JSON structured like this:

JSON

```
{response_example_json}

```

---

## private-get_rfqs

**Title:** Private Get_Rfqs
**URL:** https://docs.derive.xyz/reference/private-get_rfqs

### Method Name

#### `private/get_rfqs`

Retrieves a list of RFQs matching filter criteria. Takers can use this to get their open RFQs, RFQ history, etc.  
Required minimum session key permission level is `read_only`

### Parameters

|  |
| --- |
| **subaccount\_id**integer required  Subaccount ID for auth purposes, returned data will be scoped to this subaccount. |
| **from\_timestamp**integer  Earliest `last_update_timestamp` to filter by (in ms since Unix epoch). If not provied, defaults to 0. |
| **page**integer  Page number of results to return (default 1, returns last if above `num_pages`) |
| **page\_size**integer  Number of results per page (default 100, max 1000) |
| **rfq\_id**string  RFQ ID filter, if applicable |
| **status**string  RFQ status filter, if applicable enum  `open` `filled` `cancelled` `expired` |
| **to\_timestamp**integer  Latest `last_update_timestamp` to filter by (in ms since Unix epoch). If not provied, defaults to returning all data up to current time. |

### Response

|  |
| --- |
| **id**string or integer required |
| **result**object required |
| result.**pagination**object required  Pagination info |
| result.pagination.**count**integer required  Total number of items, across all pages |
| result.pagination.**num\_pages**integer required  Number of pages |
|  |
| result.**rfqs**array of objects required  RFQs matching filter criteria |
| result.rfqs[].**ask\_total\_cost**string or null required  Ask total cost for the RFQ implied from orderbook (as `sell`) |
| result.rfqs[].**bid\_total\_cost**string or null required  Bid total cost for the RFQ implied from orderbook (as `buy`) |
| result.rfqs[].**cancel\_reason**string required  Cancel reason, if any enum  `user_request` `insufficient_margin` `signed_max_fee_too_low` `mmp_trigger` `cancel_on_disconnect` `session_key_deregistered` `subaccount_withdrawn` `rfq_no_longer_open` `compliance` |
| result.rfqs[].**counterparties**array of strings or null required  List of requested counterparties, if applicable |
| result.rfqs[].**creation\_timestamp**integer required  Creation timestamp in ms since Unix epoch |
| result.rfqs[].**filled\_direction**string or null required  Direction at which the RFQ was filled (only if filled) enum  `buy` `sell` |
| result.rfqs[].**label**string required  User-defined label, if any |
| result.rfqs[].**last\_update\_timestamp**integer required  Last update timestamp in ms since Unix epoch |
| result.rfqs[].**mark\_total\_cost**string or null required  Mark total cost for the RFQ (assuming `buy` direction) |
| result.rfqs[].**max\_total\_cost**string or null required  Max total cost for the RFQ |
| result.rfqs[].**min\_total\_cost**string or null required  Min total cost for the RFQ |
| result.rfqs[].**rfq\_id**string required  RFQ ID |
| result.rfqs[].**status**string required  Status enum  `open` `filled` `cancelled` `expired` |
| result.rfqs[].**subaccount\_id**integer required  Subaccount ID |
| result.rfqs[].**total\_cost**string or null required  Total cost for the RFQ (only if filled) |
| result.rfqs[].**valid\_until**integer required  RFQ expiry timestamp in ms since Unix epoch |
| result.rfqs[].**legs**array of objects required  RFQ legs |
| result.rfqs[].legs[].**amount**string required  Amount in units of the base |
| result.rfqs[].legs[].**direction**string required  Leg direction enum  `buy` `sell` |
| result.rfqs[].legs[].**instrument\_name**string required  Instrument name |

### Example

ShellJavaScriptPython

```
{request_example_shell}

```

```
{request_example_javascript}

```

```
{request_example_python}

```

> The above command returns JSON structured like this:

JSON

```
{response_example_json}

```

---

## private-get_subaccount

**Title:** Private Get_Subaccount
**URL:** https://docs.derive.xyz/reference/private-get_subaccount

### Method Name

#### `private/get_subaccount`

Get open orders, active positions, and collaterals of a subaccount  
Required minimum session key permission level is `read_only`

### Parameters

|  |
| --- |
| **subaccount\_id**integer required  Subaccount\_id |

### Response

|  |
| --- |
| **id**string or integer required |
| **result**object required |
| result.**collaterals\_initial\_margin**string required  Total initial margin credit contributed by collaterals |
| result.**collaterals\_maintenance\_margin**string required  Total maintenance margin credit contributed by collaterals |
| result.**collaterals\_value**string required  Total mark-to-market value of all collaterals |
| result.**currency**string required  Currency of subaccount |
| result.**initial\_margin**string required  Total initial margin requirement of all positions and collaterals.Trades will be rejected if this value falls below zero after the trade. |
| result.**is\_under\_liquidation**boolean required  Whether the subaccount is undergoing a liquidation auction |
| result.**label**string required  User defined label |
| result.**maintenance\_margin**string required  Total maintenance margin requirement of all positions and collaterals.If this value falls below zero, the subaccount will be flagged for liquidation. |
| result.**margin\_type**string required  Margin type of subaccount (`PM` (Portfolio Margin) or `SM` (Standard Margin)) enum  `PM` `SM` `PM2` |
| result.**open\_orders\_margin**string required  Total margin requirement of all open orders.Orders will be rejected if this value plus initial margin are below zero after the order. |
| result.**positions\_initial\_margin**string required  Total initial margin requirement of all positions |
| result.**positions\_maintenance\_margin**string required  Total maintenance margin requirement of all positions |
| result.**positions\_value**string required  Total mark-to-market value of all positions |
| result.**projected\_margin\_change**string required  Projected change in maintenance margin requirement between now and projected margin at 8:01 UTC. If this value plus current maintenance margin ise below zero, the account is at risk of being flagged for liquidation right after the upcoming expiry. |
| result.**subaccount\_id**integer required  Subaccount\_id |
| result.**subaccount\_value**string required  Total mark-to-market value of all positions and collaterals |
| result.**collaterals**array of objects required  All collaterals that count towards margin of subaccount |
| result.collaterals[].**amount**string required  Asset amount of given collateral |
| result.collaterals[].**amount\_step**string required  Minimum amount step for the collateral |
| result.collaterals[].**asset\_name**string required  Asset name |
| result.collaterals[].**asset\_type**string required  Type of asset collateral (currently always `erc20`) enum  `erc20` `option` `perp` |
| result.collaterals[].**average\_price**string required  Average price of the collateral, 0 for USDC. |
| result.collaterals[].**average\_price\_excl\_fees**string required  Average price of whole position excluding fees |
| result.collaterals[].**creation\_timestamp**integer required  Timestamp of when the position was opened (in ms since Unix epoch) |
| result.collaterals[].**cumulative\_interest**string required  Cumulative interest earned on supplying collateral or paid for borrowing |
| result.collaterals[].**currency**string required  Underlying currency of asset (`ETH`, `BTC`, etc) |
| result.collaterals[].**initial\_margin**string required  USD value of collateral that contributes to initial margin |
| result.collaterals[].**maintenance\_margin**string required  USD value of collateral that contributes to maintenance margin |
| result.collaterals[].**mark\_price**string required  Current mark price of the asset |
| result.collaterals[].**mark\_value**string required  USD value of the collateral (amount \* mark price) |
| result.collaterals[].**open\_orders\_margin**string required  USD margin requirement for all open orders for this asset / instrument |
| result.collaterals[].**pending\_interest**string required  Portion of interest that has not yet been settled on-chain. This number is added to the portfolio value for margin calculations purposes. |
| result.collaterals[].**realized\_pnl**string required  Realized trading profit or loss of the collateral, 0 for USDC. |
| result.collaterals[].**realized\_pnl\_excl\_fees**string required  Realized trading profit or loss of the position excluding fees |
| result.collaterals[].**total\_fees**string required  Total fees paid for opening and changing the position |
| result.collaterals[].**unrealized\_pnl**string required  Unrealized trading profit or loss of the collateral, 0 for USDC. |
| result.collaterals[].**unrealized\_pnl\_excl\_fees**string required  Unrealized trading profit or loss of the position excluding fees |
|  |
| result.**open\_orders**array of objects required  All open orders of subaccount |
| result.open\_orders[].**amount**string required  Order amount in units of the base |
| result.open\_orders[].**average\_price**string required  Average fill price |
| result.open\_orders[].**cancel\_reason**string required  If cancelled, reason behind order cancellation enum  `user_request` `mmp_trigger` `insufficient_margin` `signed_max_fee_too_low` `cancel_on_disconnect` `ioc_or_market_partial_fill` `session_key_deregistered` `subaccount_withdrawn` `compliance` `trigger_failed` `validation_failed` |
| result.open\_orders[].**creation\_timestamp**integer required  Creation timestamp (in ms since Unix epoch) |
| result.open\_orders[].**direction**string required  Order direction enum  `buy` `sell` |
| result.open\_orders[].**filled\_amount**string required  Total filled amount for the order |
| result.open\_orders[].**instrument\_name**string required  Instrument name |
| result.open\_orders[].**is\_transfer**boolean required  Whether the order was generated through `private/transfer_position` |
| result.open\_orders[].**label**string required  Optional user-defined label for the order |
| result.open\_orders[].**last\_update\_timestamp**integer required  Last update timestamp (in ms since Unix epoch) |
| result.open\_orders[].**limit\_price**string required  Limit price in quote currency |
| result.open\_orders[].**max\_fee**string required  Max fee in units of the quote currency |
| result.open\_orders[].**mmp**boolean required  Whether the order is tagged for market maker protections |
| result.open\_orders[].**nonce**integer required  Unique nonce defined as  (e.g. 1695836058725001, where 001 is the random number) |
| result.open\_orders[].**order\_fee**string required  Total order fee paid so far |
| result.open\_orders[].**order\_id**string required  Order ID |
| result.open\_orders[].**order\_status**string required  Order status enum  `open` `filled` `cancelled` `expired` `untriggered` |
| result.open\_orders[].**order\_type**string required  Order type enum  `limit` `market` |
| result.open\_orders[].**quote\_id**string or null required  Quote ID if the trade was executed via RFQ |
| result.open\_orders[].**signature**string required  Ethereum signature of the order |
| result.open\_orders[].**signature\_expiry\_sec**integer required  Signature expiry timestamp |
| result.open\_orders[].**signer**string required  Owner wallet address or registered session key that signed order |
| result.open\_orders[].**subaccount\_id**integer required  Subaccount ID |
| result.open\_orders[].**time\_in\_force**string required  Time in force enum  `gtc` `post_only` `fok` `ioc` |
| result.open\_orders[].**replaced\_order\_id**string or null  If replaced, ID of the order that was replaced |
| result.open\_orders[].**trigger\_price**string or null  (Required for trigger orders) Index or Market price to trigger order at |
| result.open\_orders[].**trigger\_price\_type**string or null  (Required for trigger orders) Trigger with Index or Mark Price enum  `mark` `index` |
| result.open\_orders[].**trigger\_reject\_message**string or null  (Required for trigger orders) Error message if error occured during trigger |
| result.open\_orders[].**trigger\_type**string or null  (Required for trigger orders) Stop-loss or Take-profit. enum  `stoploss` `takeprofit` |
|  |
| result.**positions**array of objects required  All active positions of subaccount |
| result.positions[].**amount**string required  Position amount held by subaccount |
| result.positions[].**amount\_step**string required  Minimum amount step for the position |
| result.positions[].**average\_price**string required  Average price of whole position |
| result.positions[].**average\_price\_excl\_fees**string required  Average price of whole position excluding fees |
| result.positions[].**creation\_timestamp**integer required  Timestamp of when the position was opened (in ms since Unix epoch) |
| result.positions[].**cumulative\_funding**string required  Cumulative funding for the position (only for perpetuals). |
| result.positions[].**delta**string required  Asset delta (w.r.t. forward price for options, `1.0` for perps) |
| result.positions[].**gamma**string required  Asset gamma (zero for non-options) |
| result.positions[].**index\_price**string required  Current index (oracle) price for position's currency |
| result.positions[].**initial\_margin**string required  USD initial margin requirement for this position |
| result.positions[].**instrument\_name**string required  Instrument name (same as the base Asset name) |
| result.positions[].**instrument\_type**string required  `erc20`, `option`, or `perp` enum  `erc20` `option` `perp` |
| result.positions[].**leverage**string or null required  Only for perps. Leverage of the position, defined as `abs(notional) / collateral net of options margin` |
| result.positions[].**liquidation\_price**string or null required  Index price at which position will be liquidated |
| result.positions[].**maintenance\_margin**string required  USD maintenance margin requirement for this position |
| result.positions[].**mark\_price**string required  Current mark price for position's instrument |
| result.positions[].**mark\_value**string required  USD value of the position; this represents how much USD can be recieved by fully closing the position at the current oracle price |
| result.positions[].**net\_settlements**string required  Net amount of USD from position settlements that has been paid to the user's subaccount. This number is subtracted from the portfolio value for margin calculations purposes. Positive values mean the user has recieved USD from settlements, or is awaiting settlement of USD losses. Negative values mean the user has paid USD for settlements, or is awaiting settlement of USD gains. |
| result.positions[].**open\_orders\_margin**string required  USD margin requirement for all open orders for this asset / instrument |
| result.positions[].**pending\_funding**string required  A portion of funding payments that has not yet been settled into cash balance (only for perpetuals). This number is added to the portfolio value for margin calculations purposes. |
| result.positions[].**realized\_pnl**string required  Realized trading profit or loss of the position. |
| result.positions[].**realized\_pnl\_excl\_fees**string required  Realized trading profit or loss of the position excluding fees |
| result.positions[].**theta**string required  Asset theta (zero for non-options) |
| result.positions[].**total\_fees**string required  Total fees paid for opening and changing the position |
| result.positions[].**unrealized\_pnl**string required  Unrealized trading profit or loss of the position. |
| result.positions[].**unrealized\_pnl\_excl\_fees**string required  Unrealized trading profit or loss of the position excluding fees |
| result.positions[].**vega**string required  Asset vega (zero for non-options) |

### Example

ShellJavaScriptPython

```
{request_example_shell}

```

```
{request_example_javascript}

```

```
{request_example_python}

```

> The above command returns JSON structured like this:

JSON

```
{response_example_json}

```

---

## private-get_subaccount_value_history

**Title:** Private Get_Subaccount_Value_History
**URL:** https://docs.derive.xyz/reference/private-get_subaccount_value_history

### Method Name

#### `private/get_subaccount_value_history`

Get the value history of a subaccount  
Required minimum session key permission level is `read_only`

### Parameters

|  |
| --- |
| **end\_timestamp**integer required  End timestamp |
| **period**integer required  Period |
| **start\_timestamp**integer required  Start timestamp |
| **subaccount\_id**integer required  Subaccount\_id |

### Response

|  |
| --- |
| **id**string or integer required |
| **result**object required |
| result.**subaccount\_id**integer required  Subaccount\_id |
| result.**subaccount\_value\_history**array of objects required  Subaccount value history |
| result.subaccount\_value\_history[].**subaccount\_value**string required  Total mark-to-market value of all positions and collaterals |
| result.subaccount\_value\_history[].**timestamp**integer required  Timestamp of when the subaccount value was recorded into the database |

### Example

ShellJavaScriptPython

```
{request_example_shell}

```

```
{request_example_javascript}

```

```
{request_example_python}

```

> The above command returns JSON structured like this:

JSON

```
{response_example_json}

```

---

## private-get_subaccounts

**Title:** Private Get_Subaccounts
**URL:** https://docs.derive.xyz/reference/private-get_subaccounts

### Method Name

#### `private/get_subaccounts`

Get all subaccounts of an account / wallet  
Required minimum session key permission level is `read_only`

### Parameters

|  |
| --- |
| **wallet**string required  Ethereum wallet address of account |

### Response

|  |
| --- |
| **id**string or integer required |
| **result**object required |
| result.**wallet**string required  Ethereum wallet address |
| result.**subaccount\_ids**array of integers required  List of subaccount\_ids owned by the wallet in `SubAccounts.sol` |

### Example

ShellJavaScriptPython

```
{request_example_shell}

```

```
{request_example_javascript}

```

```
{request_example_python}

```

> The above command returns JSON structured like this:

JSON

```
{response_example_json}

```

---

## private-get_trade_history

**Title:** Private Get_Trade_History
**URL:** https://docs.derive.xyz/reference/private-get_trade_history

### Method Name

#### `private/get_trade_history`

Get trade history for a subaccount, with filter parameters.  
Required minimum session key permission level is `read_only`

### Parameters

|  |
| --- |
| **from\_timestamp**integer  Earliest timestamp to filter by (in ms since Unix epoch). If not provied, defaults to 0. |
| **instrument\_name**string  Instrument name to filter by |
| **order\_id**string  Order id to filter by |
| **page**integer  Page number of results to return (default 1, returns last if above `num_pages`) |
| **page\_size**integer  Number of results per page (default 100, max 1000) |
| **quote\_id**string  If supplied, quote id to filter by. Supports either a concrete UUID, or `is_quote` and `is_not_quote` enum |
| **subaccount\_id**integer  Subaccount\_id (must be set if wallet is blank) |
| **to\_timestamp**integer  Latest timestamp to filter by (in ms since Unix epoch). If not provied, defaults to returning all data up to current time. |
| **wallet**string  Wallet address (if set, subaccount\_id ignored) |

### Response

|  |
| --- |
| **id**string or integer required |
| **result**object required |
| result.**subaccount\_id**integer required  Subaccount ID requested, or 0 if not provided |
| result.**pagination**object required  Pagination info |
| result.pagination.**count**integer required  Total number of items, across all pages |
| result.pagination.**num\_pages**integer required  Number of pages |
|  |
| result.**trades**array of objects required  List of trades |
| result.trades[].**direction**string required  Order direction enum  `buy` `sell` |
| result.trades[].**expected\_rebate**string required  Expected rebate for this trade |
| result.trades[].**index\_price**string required  Index price of the underlying at the time of the trade |
| result.trades[].**instrument\_name**string required  Instrument name |
| result.trades[].**is\_transfer**boolean required  Whether the trade was generated through `private/transfer_position` |
| result.trades[].**label**string required  Optional user-defined label for the order |
| result.trades[].**liquidity\_role**string required  Role of the user in the trade enum  `maker` `taker` |
| result.trades[].**mark\_price**string required  Mark price of the instrument at the time of the trade |
| result.trades[].**order\_id**string required  Order ID |
| result.trades[].**quote\_id**string or null required  Quote ID if the trade was executed via RFQ |
| result.trades[].**realized\_pnl**string required  Realized PnL for this trade |
| result.trades[].**realized\_pnl\_excl\_fees**string required  Realized PnL for this trade using cost accounting that excludes fees |
| result.trades[].**subaccount\_id**integer required  Subaccount ID |
| result.trades[].**timestamp**integer required  Trade timestamp (in ms since Unix epoch) |
| result.trades[].**trade\_amount**string required  Amount filled in this trade |
| result.trades[].**trade\_fee**string required  Fee for this trade |
| result.trades[].**trade\_id**string required  Trade ID |
| result.trades[].**trade\_price**string required  Price at which the trade was filled |
| result.trades[].**transaction\_id**string required  The transaction id of the related settlement transaction |
| result.trades[].**tx\_hash**string or null required  Blockchain transaction hash |
| result.trades[].**tx\_status**string required  Blockchain transaction status enum  `requested` `pending` `settled` `reverted` `ignored` `timed_out` |

### Example

ShellJavaScriptPython

```
{request_example_shell}

```

```
{request_example_javascript}

```

```
{request_example_python}

```

> The above command returns JSON structured like this:

JSON

```
{response_example_json}

```

---

## private-get_withdrawal_history

**Title:** Private Get_Withdrawal_History
**URL:** https://docs.derive.xyz/reference/private-get_withdrawal_history

### Method Name

#### `private/get_withdrawal_history`

Get subaccount withdrawal history.  
Required minimum session key permission level is `read_only`

### Parameters

|  |
| --- |
| **subaccount\_id**integer required  Subaccount id |
| **end\_timestamp**integer  End timestamp of the event history (default current time) |
| **start\_timestamp**integer  Start timestamp of the event history (default 0) |

### Response

|  |
| --- |
| **id**string or integer required |
| **result**object required |
| result.**events**array of objects required  List of withdrawals |
| result.events[].**amount**string required  Amount withdrawn by the subaccount |
| result.events[].**asset**string required  Asset withdrawn |
| result.events[].**error\_log**object or null required  If failed, error log for reason |
| result.events[].**timestamp**integer required  Timestamp of the withdrawal (in ms since UNIX epoch) |
| result.events[].**tx\_hash**string required  Hash of the transaction that withdrew the funds |
| result.events[].**tx\_status**string required  Status of the transaction that deposited the funds enum  `requested` `pending` `settled` `reverted` `ignored` `timed_out` |

### Example

ShellJavaScriptPython

```
{request_example_shell}

```

```
{request_example_javascript}

```

```
{request_example_python}

```

> The above command returns JSON structured like this:

JSON

```
{response_example_json}

```

---

## private-liquidate

**Title:** Private Liquidate
**URL:** https://docs.derive.xyz/reference/private-liquidate

### Method Name

#### `private/liquidate`

Liquidates a given subaccount using funds from another subaccount. This endpoint has a few limitations:  
1. If succesful, the RPC will freeze the caller's subaccount until the bid is settled or is reverted on chain.  
2. The caller's subaccount must not have any open orders.  
3. The caller's subaccount must have enough withdrawable cash to cover the bid and the buffer margin requirements.  
Required minimum session key permission level is `admin`

### Parameters

|  |
| --- |
| **cash\_transfer**string required  Amount of cash to transfer to a newly created subaccount for bidding. Must be non-negative. |
| **last\_seen\_trade\_id**integer required  Last seen trade ID for account being liquidated. Not checked if set to 0. |
| **liquidated\_subaccount\_id**integer required  Subaccount ID of the account to be liquidated. |
| **nonce**integer required  Unique nonce defined as  (e.g. 1695836058725001, where 001 is the random number) |
| **percent\_bid**string required  Percent of the liquidated position to bid for. Will bid for the maximum possible percent of the position if set to 1 |
| **price\_limit**string required  Maximum amount of cash to be paid from bidder to liquidated account (supports negative amounts for insolvent auctions). Not checked if set to 0. |
| **signature**string required  Ethereum signature of the order |
| **signature\_expiry\_sec**integer required  Unix timestamp in seconds. Order signature becomes invalid after this time, and the system will cancel the order.Expiry MUST be at least 5 min from now. |
| **signer**string required  Owner wallet address or registered session key that signed order |
| **subaccount\_id**integer required  Subaccount ID owned by wallet, that will be doing the bidding. |

### Response

|  |
| --- |
| **id**string or integer required |
| **result**object required |
| result.**estimated\_bid\_price**string required  Estimated bid price for this liquidation |
| result.**estimated\_discount\_pnl**string required  Estimated profit (increase in the subaccount mark value) if the liquidation is successful. |
| result.**estimated\_percent\_bid**string required  Estimated percent of account the bid will aquire |
| result.**transaction\_id**string required  The transaction id of the related settlement transaction |

### Example

ShellJavaScriptPython

```
{request_example_shell}

```

```
{request_example_javascript}

```

```
{request_example_python}

```

> The above command returns JSON structured like this:

JSON

```
{response_example_json}

```

---

## private-order

**Title:** Private Order
**URL:** https://docs.derive.xyz/reference/private-order

### Method Name

#### `private/order`

Create a new order.  
Required minimum session key permission level is `admin`

### Parameters

|  |
| --- |
| **amount**string required  Order amount in units of the base |
| **direction**string required  Order direction enum  `buy` `sell` |
| **instrument\_name**string required  Instrument name |
| **limit\_price**string required  Limit price in quote currency. This field is still required for market orders because it is a component of the signature. However, market orders will not leave a resting order in the book in case of a partial fill. |
| **max\_fee**string required  Max fee per unit of volume, denominated in units of the quote currency (usually USDC).Order will be rejected if the supplied max fee is below the estimated fee for this order. |
| **nonce**integer required  Unique nonce defined as  (e.g. 1695836058725001, where 001 is the random number).Note, using a random number beyond 3 digits will cause JSON serialization to fail. |
| **signature**string required  Ethereum signature of the order |
| **signature\_expiry\_sec**integer required  Unix timestamp in seconds. Order signature becomes invalid after this time, and the system will cancel the order.Expiry MUST be at least 5 min from now. |
| **signer**string required  Owner wallet address or registered session key that signed order |
| **subaccount\_id**integer required  Subaccount ID |
| **is\_atomic\_signing**boolean  Used by vaults to determine whether the signature is an EIP-1271 signature. |
| **label**string  Optional user-defined label for the order |
| **mmp**boolean  Whether the order is tagged for market maker protections (default false) |
| **order\_type**string  Order type: - `limit`: limit order (default) - `market`: market order, note that limit\_price is still required for market orders, but unfilled order portion will be marked as cancelled enum  `limit` `market` |
| **reduce\_only**boolean  If true, the order will not be able to increase position's size (default false). If the order amount exceeds available position size, the order will be filled up to the position size and the remainder will be cancelled. This flag is only supported for market orders or non-resting limit orders (IOC or FOK) |
| **referral\_code**string  Optional referral code for the order |
| **reject\_timestamp**integer  UTC timestamp in ms, if provided the matching engine will reject the order with an error if `reject_timestamp` < `server_time`. Note that the timestamp must be consistent with the server time: use `public/get_time` method to obtain current server time. |
| **time\_in\_force**string  Time in force behaviour: - `gtc`: good til cancelled (default) - `post_only`: a limit order that will be rejected if it crosses any order in the book, i.e. acts as a taker order - `fok`: fill or kill, will be rejected if it is not fully filled - `ioc`: immediate or cancel, fill at best bid/ask (market) or at limit price (limit), the unfilled portion is cancelled Note that the order will still expire on the `signature_expiry_sec` timestamp. enum  `gtc` `post_only` `fok` `ioc` |
| **trigger\_price**string  (Required for trigger orders) "index" or "mark" price to trigger order at |
| **trigger\_price\_type**string  (Required for trigger orders) Trigger with "mark" price as "index" price type not supported yet. enum  `mark` `index` |
| **trigger\_type**string  (Required for trigger orders) "stoploss" or "takeprofit" enum  `stoploss` `takeprofit` |

### Response

|  |
| --- |
| **id**string or integer required |
| **result**object required |
| result.**order**object required |
| result.order.**amount**string required  Order amount in units of the base |
| result.order.**average\_price**string required  Average fill price |
| result.order.**cancel\_reason**string required  If cancelled, reason behind order cancellation enum  `user_request` `mmp_trigger` `insufficient_margin` `signed_max_fee_too_low` `cancel_on_disconnect` `ioc_or_market_partial_fill` `session_key_deregistered` `subaccount_withdrawn` `compliance` `trigger_failed` `validation_failed` |
| result.order.**creation\_timestamp**integer required  Creation timestamp (in ms since Unix epoch) |
| result.order.**direction**string required  Order direction enum  `buy` `sell` |
| result.order.**filled\_amount**string required  Total filled amount for the order |
| result.order.**instrument\_name**string required  Instrument name |
| result.order.**is\_transfer**boolean required  Whether the order was generated through `private/transfer_position` |
| result.order.**label**string required  Optional user-defined label for the order |
| result.order.**last\_update\_timestamp**integer required  Last update timestamp (in ms since Unix epoch) |
| result.order.**limit\_price**string required  Limit price in quote currency |
| result.order.**max\_fee**string required  Max fee in units of the quote currency |
| result.order.**mmp**boolean required  Whether the order is tagged for market maker protections |
| result.order.**nonce**integer required  Unique nonce defined as  (e.g. 1695836058725001, where 001 is the random number) |
| result.order.**order\_fee**string required  Total order fee paid so far |
| result.order.**order\_id**string required  Order ID |
| result.order.**order\_status**string required  Order status enum  `open` `filled` `cancelled` `expired` `untriggered` |
| result.order.**order\_type**string required  Order type enum  `limit` `market` |
| result.order.**quote\_id**string or null required  Quote ID if the trade was executed via RFQ |
| result.order.**signature**string required  Ethereum signature of the order |
| result.order.**signature\_expiry\_sec**integer required  Signature expiry timestamp |
| result.order.**signer**string required  Owner wallet address or registered session key that signed order |
| result.order.**subaccount\_id**integer required  Subaccount ID |
| result.order.**time\_in\_force**string required  Time in force enum  `gtc` `post_only` `fok` `ioc` |
| result.order.**replaced\_order\_id**string or null  If replaced, ID of the order that was replaced |
| result.order.**trigger\_price**string or null  (Required for trigger orders) Index or Market price to trigger order at |
| result.order.**trigger\_price\_type**string or null  (Required for trigger orders) Trigger with Index or Mark Price enum  `mark` `index` |
| result.order.**trigger\_reject\_message**string or null  (Required for trigger orders) Error message if error occured during trigger |
| result.order.**trigger\_type**string or null  (Required for trigger orders) Stop-loss or Take-profit. enum  `stoploss` `takeprofit` |
|  |
| result.**trades**array of objects required |
| result.trades[].**direction**string required  Order direction enum  `buy` `sell` |
| result.trades[].**expected\_rebate**string required  Expected rebate for this trade |
| result.trades[].**index\_price**string required  Index price of the underlying at the time of the trade |
| result.trades[].**instrument\_name**string required  Instrument name |
| result.trades[].**is\_transfer**boolean required  Whether the trade was generated through `private/transfer_position` |
| result.trades[].**label**string required  Optional user-defined label for the order |
| result.trades[].**liquidity\_role**string required  Role of the user in the trade enum  `maker` `taker` |
| result.trades[].**mark\_price**string required  Mark price of the instrument at the time of the trade |
| result.trades[].**order\_id**string required  Order ID |
| result.trades[].**quote\_id**string or null required  Quote ID if the trade was executed via RFQ |
| result.trades[].**realized\_pnl**string required  Realized PnL for this trade |
| result.trades[].**realized\_pnl\_excl\_fees**string required  Realized PnL for this trade using cost accounting that excludes fees |
| result.trades[].**subaccount\_id**integer required  Subaccount ID |
| result.trades[].**timestamp**integer required  Trade timestamp (in ms since Unix epoch) |
| result.trades[].**trade\_amount**string required  Amount filled in this trade |
| result.trades[].**trade\_fee**string required  Fee for this trade |
| result.trades[].**trade\_id**string required  Trade ID |
| result.trades[].**trade\_price**string required  Price at which the trade was filled |
| result.trades[].**transaction\_id**string required  The transaction id of the related settlement transaction |
| result.trades[].**tx\_hash**string or null required  Blockchain transaction hash |
| result.trades[].**tx\_status**string required  Blockchain transaction status enum  `requested` `pending` `settled` `reverted` `ignored` `timed_out` |

### Example

ShellJavaScriptPython

```
{request_example_shell}

```

```
{request_example_javascript}

```

```
{request_example_python}

```

> The above command returns JSON structured like this:

JSON

```
{response_example_json}

```

---

## private-order_debug

**Title:** Private Order_Debug
**URL:** https://docs.derive.xyz/reference/private-order_debug

### Method Name

#### `private/order_debug`

Debug a new order  
Required minimum session key permission level is `read_only`

### Parameters

|  |
| --- |
| **amount**string required  Order amount in units of the base |
| **direction**string required  Order direction enum  `buy` `sell` |
| **instrument\_name**string required  Instrument name |
| **limit\_price**string required  Limit price in quote currency. This field is still required for market orders because it is a component of the signature. However, market orders will not leave a resting order in the book in case of a partial fill. |
| **max\_fee**string required  Max fee per unit of volume, denominated in units of the quote currency (usually USDC).Order will be rejected if the supplied max fee is below the estimated fee for this order. |
| **nonce**integer required  Unique nonce defined as  (e.g. 1695836058725001, where 001 is the random number).Note, using a random number beyond 3 digits will cause JSON serialization to fail. |
| **signature**string required  Ethereum signature of the order |
| **signature\_expiry\_sec**integer required  Unix timestamp in seconds. Order signature becomes invalid after this time, and the system will cancel the order.Expiry MUST be at least 5 min from now. |
| **signer**string required  Owner wallet address or registered session key that signed order |
| **subaccount\_id**integer required  Subaccount ID |
| **is\_atomic\_signing**boolean  Used by vaults to determine whether the signature is an EIP-1271 signature. |
| **label**string  Optional user-defined label for the order |
| **mmp**boolean  Whether the order is tagged for market maker protections (default false) |
| **order\_type**string  Order type: - `limit`: limit order (default) - `market`: market order, note that limit\_price is still required for market orders, but unfilled order portion will be marked as cancelled enum  `limit` `market` |
| **reduce\_only**boolean  If true, the order will not be able to increase position's size (default false). If the order amount exceeds available position size, the order will be filled up to the position size and the remainder will be cancelled. This flag is only supported for market orders or non-resting limit orders (IOC or FOK) |
| **referral\_code**string  Optional referral code for the order |
| **reject\_timestamp**integer  UTC timestamp in ms, if provided the matching engine will reject the order with an error if `reject_timestamp` < `server_time`. Note that the timestamp must be consistent with the server time: use `public/get_time` method to obtain current server time. |
| **time\_in\_force**string  Time in force behaviour: - `gtc`: good til cancelled (default) - `post_only`: a limit order that will be rejected if it crosses any order in the book, i.e. acts as a taker order - `fok`: fill or kill, will be rejected if it is not fully filled - `ioc`: immediate or cancel, fill at best bid/ask (market) or at limit price (limit), the unfilled portion is cancelled Note that the order will still expire on the `signature_expiry_sec` timestamp. enum  `gtc` `post_only` `fok` `ioc` |
| **trigger\_price**string  (Required for trigger orders) "index" or "mark" price to trigger order at |
| **trigger\_price\_type**string  (Required for trigger orders) Trigger with "mark" price as "index" price type not supported yet. enum  `mark` `index` |
| **trigger\_type**string  (Required for trigger orders) "stoploss" or "takeprofit" enum  `stoploss` `takeprofit` |

### Response

|  |
| --- |
| **id**string or integer required |
| **result**object required |
| result.**action\_hash**string required  Keccak hashed action data |
| result.**encoded\_data**string required  ABI encoded order data |
| result.**encoded\_data\_hashed**string required  Keccak hashed encoded\_data |
| result.**typed\_data\_hash**string required  EIP 712 typed data hash |
| result.**raw\_data**object required  Raw order data |
| result.raw\_data.**expiry**integer required |
| result.raw\_data.**is\_atomic\_signing**boolean required |
| result.raw\_data.**module**string required |
| result.raw\_data.**nonce**integer required |
| result.raw\_data.**owner**string required |
| result.raw\_data.**signature**string required |
| result.raw\_data.**signer**string required |
| result.raw\_data.**subaccount\_id**integer required |
| result.raw\_data.**data**object required |
| result.raw\_data.data.**asset**string required |
| result.raw\_data.data.**desired\_amount**string required |
| result.raw\_data.data.**is\_bid**boolean required |
| result.raw\_data.data.**limit\_price**string required |
| result.raw\_data.data.**recipient\_id**integer required |
| result.raw\_data.data.**sub\_id**integer required |
| result.raw\_data.data.**trade\_id**string required |
| result.raw\_data.data.**worst\_fee**string required |

### Example

ShellJavaScriptPython

```
{request_example_shell}

```

```
{request_example_javascript}

```

```
{request_example_python}

```

> The above command returns JSON structured like this:

JSON

```
{response_example_json}

```

---

## private-poll_quotes

**Title:** Private Poll_Quotes
**URL:** https://docs.derive.xyz/reference/private-poll_quotes

### Method Name

#### `private/poll_quotes`

Retrieves a list of quotes matching filter criteria.  
Takers can use this to poll open quotes that they can fill against their open RFQs.  
Required minimum session key permission level is `read_only`

### Parameters

|  |
| --- |
| **subaccount\_id**integer required  Subaccount ID for auth purposes, returned data will be scoped to this subaccount. |
| **from\_timestamp**integer  Earliest timestamp to filter by (in ms since Unix epoch). If not provied, defaults to 0. |
| **page**integer  Page number of results to return (default 1, returns last if above `num_pages`) |
| **page\_size**integer  Number of results per page (default 100, max 1000) |
| **quote\_id**string  Quote ID filter, if applicable |
| **rfq\_id**string  RFQ ID filter, if applicable |
| **status**string  Quote status filter, if applicable enum  `open` `filled` `cancelled` `expired` |
| **to\_timestamp**integer  Latest timestamp to filter by (in ms since Unix epoch). If not provied, defaults to returning all data up to current time. |

### Response

|  |
| --- |
| **id**string or integer required |
| **result**object required |
| result.**pagination**object required  Pagination info |
| result.pagination.**count**integer required  Total number of items, across all pages |
| result.pagination.**num\_pages**integer required  Number of pages |
|  |
| result.**quotes**array of objects required  Quotes matching filter criteria |
| result.quotes[].**cancel\_reason**string required  Cancel reason, if any enum  `user_request` `insufficient_margin` `signed_max_fee_too_low` `mmp_trigger` `cancel_on_disconnect` `session_key_deregistered` `subaccount_withdrawn` `rfq_no_longer_open` `compliance` |
| result.quotes[].**creation\_timestamp**integer required  Creation timestamp in ms since Unix epoch |
| result.quotes[].**direction**string required  Quote direction enum  `buy` `sell` |
| result.quotes[].**last\_update\_timestamp**integer required  Last update timestamp in ms since Unix epoch |
| result.quotes[].**legs\_hash**string required  Hash of the legs of the best quote to be signed by the taker. |
| result.quotes[].**liquidity\_role**string required  Liquidity role enum  `maker` `taker` |
| result.quotes[].**quote\_id**string required  Quote ID |
| result.quotes[].**rfq\_id**string required  RFQ ID |
| result.quotes[].**status**string required  Status enum  `open` `filled` `cancelled` `expired` |
| result.quotes[].**subaccount\_id**integer required  Subaccount ID |
| result.quotes[].**tx\_hash**string or null required  Blockchain transaction hash (only for executed quotes) |
| result.quotes[].**tx\_status**string or null required  Blockchain transaction status (only for executed quotes) enum  `requested` `pending` `settled` `reverted` `ignored` `timed_out` |
| result.quotes[].**wallet**string required  Wallet address of the quote sender |
| result.quotes[].**legs**array of objects required  Quote legs |
| result.quotes[].legs[].**amount**string required  Amount in units of the base |
| result.quotes[].legs[].**direction**string required  Leg direction enum  `buy` `sell` |
| result.quotes[].legs[].**instrument\_name**string required  Instrument name |
| result.quotes[].legs[].**price**string required  Leg price |

### Example

ShellJavaScriptPython

```
{request_example_shell}

```

```
{request_example_javascript}

```

```
{request_example_python}

```

> The above command returns JSON structured like this:

JSON

```
{response_example_json}

```

---

## private-poll_rfqs

**Title:** Private Poll_Rfqs
**URL:** https://docs.derive.xyz/reference/private-poll_rfqs

### Method Name

#### `private/poll_rfqs`

Retrieves a list of RFQs matching filter criteria. Market makers can use this to poll RFQs directed to them.  
Required minimum session key permission level is `read_only`

### Parameters

|  |
| --- |
| **subaccount\_id**integer required  Subaccount ID for auth purposes, returned data will be scoped to this subaccount. |
| **from\_timestamp**integer  Earliest `last_update_timestamp` to filter by (in ms since Unix epoch). If not provied, defaults to 0. |
| **page**integer  Page number of results to return (default 1, returns last if above `num_pages`) |
| **page\_size**integer  Number of results per page (default 100, max 1000) |
| **rfq\_id**string  RFQ ID filter, if applicable |
| **rfq\_subaccount\_id**integer  Filter returned RFQs by rfq requestor subaccount |
| **status**string  RFQ status filter, if applicable enum  `open` `filled` `cancelled` `expired` |
| **to\_timestamp**integer  Latest `last_update_timestamp` to filter by (in ms since Unix epoch). If not provied, defaults to returning all data up to current time. |

### Response

|  |
| --- |
| **id**string or integer required |
| **result**object required |
| result.**pagination**object required  Pagination info |
| result.pagination.**count**integer required  Total number of items, across all pages |
| result.pagination.**num\_pages**integer required  Number of pages |
|  |
| result.**rfqs**array of objects required  RFQs matching filter criteria |
| result.rfqs[].**cancel\_reason**string required  Cancel reason, if any enum  `user_request` `insufficient_margin` `signed_max_fee_too_low` `mmp_trigger` `cancel_on_disconnect` `session_key_deregistered` `subaccount_withdrawn` `rfq_no_longer_open` `compliance` |
| result.rfqs[].**creation\_timestamp**integer required  Creation timestamp in ms since Unix epoch |
| result.rfqs[].**filled\_direction**string or null required  Direction at which the RFQ was filled (only if filled) enum  `buy` `sell` |
| result.rfqs[].**last\_update\_timestamp**integer required  Last update timestamp in ms since Unix epoch |
| result.rfqs[].**rfq\_id**string required  RFQ ID |
| result.rfqs[].**status**string required  Status enum  `open` `filled` `cancelled` `expired` |
| result.rfqs[].**subaccount\_id**integer required  Subaccount ID |
| result.rfqs[].**total\_cost**string or null required  Total cost for the RFQ (only if filled) |
| result.rfqs[].**valid\_until**integer required  RFQ expiry timestamp in ms since Unix epoch |
| result.rfqs[].**legs**array of objects required  RFQ legs |
| result.rfqs[].legs[].**amount**string required  Amount in units of the base |
| result.rfqs[].legs[].**direction**string required  Leg direction enum  `buy` `sell` |
| result.rfqs[].legs[].**instrument\_name**string required  Instrument name |

### Example

ShellJavaScriptPython

```
{request_example_shell}

```

```
{request_example_javascript}

```

```
{request_example_python}

```

> The above command returns JSON structured like this:

JSON

```
{response_example_json}

```

---

## private-register_scoped_session_key

**Title:** Private Register_Scoped_Session_Key
**URL:** https://docs.derive.xyz/reference/private-register_scoped_session_key

### Method Name

#### `private/register_scoped_session_key`

Registers a new session key bounded to a scope without a transaction attached.  
If you want to register an admin key, you must provide a signed raw transaction.  
Required minimum session key permission level is `account`

### Parameters

|  |
| --- |
| **expiry\_sec**integer required  Expiry of the session key |
| **public\_session\_key**string required  Session key in the form of an Ethereum EOA |
| **wallet**string required  Ethereum wallet address of account |
| **ip\_whitelist**array of strings  List of whitelisted IPs, if empty then any IP is allowed. |
| **label**string  User-defined session key label |
| **scope**string  Scope of the session key. Defaults to READ\_ONLY level permissions.  enum  `admin` `account` `read_only` |
| **signed\_raw\_tx**string  A signed RLP encoded ETH transaction in form of a hex string (same as `w3.eth.account.sign_transaction(unsigned_tx, private_key).rawTransaction.hex()`) Must be included if the scope is ADMIN. |

### Response

|  |
| --- |
| **id**string or integer required |
| **result**object required |
| result.**expiry\_sec**integer required  Session key expiry timestamp in sec |
| result.**ip\_whitelist**array of strings or null required  List of whitelisted IPs, if empty then any IP is allowed. |
| result.**label**string or null required  User-defined session key label |
| result.**public\_session\_key**string required  Session key in the form of an Ethereum EOA |
| result.**scope**string required  Session key permission level scope enum  `admin` `account` `read_only` |
| result.**transaction\_id**string or null required  ID to lookup status of transaction if signed\_raw\_tx is provided |

### Example

ShellJavaScriptPython

```
{request_example_shell}

```

```
{request_example_javascript}

```

```
{request_example_python}

```

> The above command returns JSON structured like this:

JSON

```
{response_example_json}

```

---

## private-replace

**Title:** Private Replace
**URL:** https://docs.derive.xyz/reference/private-replace

### Method Name

#### `private/replace`

Cancel an existing order with nonce or order\_id and create new order with different order\_id in a single RPC call.  
  
If the cancel fails, the new order will not be created.  
If the cancel succeeds but the new order fails, the old order will still be cancelled.  
Required minimum session key permission level is `admin`

### Parameters

|  |
| --- |
| **amount**string required  Order amount in units of the base |
| **direction**string required  Order direction enum  `buy` `sell` |
| **instrument\_name**string required  Instrument name |
| **limit\_price**string required  Limit price in quote currency. This field is still required for market orders because it is a component of the signature. However, market orders will not leave a resting order in the book in case of a partial fill. |
| **max\_fee**string required  Max fee per unit of volume, denominated in units of the quote currency (usually USDC).Order will be rejected if the supplied max fee is below the estimated fee for this order. |
| **nonce**integer required  Unique nonce defined as  (e.g. 1695836058725001, where 001 is the random number).Note, using a random number beyond 3 digits will cause JSON serialization to fail. |
| **signature**string required  Ethereum signature of the order |
| **signature\_expiry\_sec**integer required  Unix timestamp in seconds. Order signature becomes invalid after this time, and the system will cancel the order.Expiry MUST be at least 5 min from now. |
| **signer**string required  Owner wallet address or registered session key that signed order |
| **subaccount\_id**integer required  Subaccount ID |
| **expected\_filled\_amount**string  Optional check to only create new order if old order filled\_amount is equal to this value. |
| **is\_atomic\_signing**boolean  Used by vaults to determine whether the signature is an EIP-1271 signature. |
| **label**string  Optional user-defined label for the order |
| **mmp**boolean  Whether the order is tagged for market maker protections (default false) |
| **nonce\_to\_cancel**integer  Cancel order by nonce (choose either order\_id or nonce). |
| **order\_id\_to\_cancel**string  Cancel order by order\_id (choose either order\_id or nonce). |
| **order\_type**string  Order type: - `limit`: limit order (default) - `market`: market order, note that limit\_price is still required for market orders, but unfilled order portion will be marked as cancelled enum  `limit` `market` |
| **reduce\_only**boolean  If true, the order will not be able to increase position's size (default false). If the order amount exceeds available position size, the order will be filled up to the position size and the remainder will be cancelled. This flag is only supported for market orders or non-resting limit orders (IOC or FOK) |
| **referral\_code**string  Optional referral code for the order |
| **reject\_timestamp**integer  UTC timestamp in ms, if provided the matching engine will reject the order with an error if `reject_timestamp` < `server_time`. Note that the timestamp must be consistent with the server time: use `public/get_time` method to obtain current server time. |
| **time\_in\_force**string  Time in force behaviour: - `gtc`: good til cancelled (default) - `post_only`: a limit order that will be rejected if it crosses any order in the book, i.e. acts as a taker order - `fok`: fill or kill, will be rejected if it is not fully filled - `ioc`: immediate or cancel, fill at best bid/ask (market) or at limit price (limit), the unfilled portion is cancelled Note that the order will still expire on the `signature_expiry_sec` timestamp. enum  `gtc` `post_only` `fok` `ioc` |
| **trigger\_price**string  (Required for trigger orders) "index" or "mark" price to trigger order at |
| **trigger\_price\_type**string  (Required for trigger orders) Trigger with "mark" price as "index" price type not supported yet. enum  `mark` `index` |
| **trigger\_type**string  (Required for trigger orders) "stoploss" or "takeprofit" enum  `stoploss` `takeprofit` |

### Response

|  |
| --- |
| **id**string or integer required |
| **result**object required |
| result.**cancelled\_order**object required  Order that was cancelled |
| result.cancelled\_order.**amount**string required  Order amount in units of the base |
| result.cancelled\_order.**average\_price**string required  Average fill price |
| result.cancelled\_order.**cancel\_reason**string required  If cancelled, reason behind order cancellation enum  `user_request` `mmp_trigger` `insufficient_margin` `signed_max_fee_too_low` `cancel_on_disconnect` `ioc_or_market_partial_fill` `session_key_deregistered` `subaccount_withdrawn` `compliance` `trigger_failed` `validation_failed` |
| result.cancelled\_order.**creation\_timestamp**integer required  Creation timestamp (in ms since Unix epoch) |
| result.cancelled\_order.**direction**string required  Order direction enum  `buy` `sell` |
| result.cancelled\_order.**filled\_amount**string required  Total filled amount for the order |
| result.cancelled\_order.**instrument\_name**string required  Instrument name |
| result.cancelled\_order.**is\_transfer**boolean required  Whether the order was generated through `private/transfer_position` |
| result.cancelled\_order.**label**string required  Optional user-defined label for the order |
| result.cancelled\_order.**last\_update\_timestamp**integer required  Last update timestamp (in ms since Unix epoch) |
| result.cancelled\_order.**limit\_price**string required  Limit price in quote currency |
| result.cancelled\_order.**max\_fee**string required  Max fee in units of the quote currency |
| result.cancelled\_order.**mmp**boolean required  Whether the order is tagged for market maker protections |
| result.cancelled\_order.**nonce**integer required  Unique nonce defined as  (e.g. 1695836058725001, where 001 is the random number) |
| result.cancelled\_order.**order\_fee**string required  Total order fee paid so far |
| result.cancelled\_order.**order\_id**string required  Order ID |
| result.cancelled\_order.**order\_status**string required  Order status enum  `open` `filled` `cancelled` `expired` `untriggered` |
| result.cancelled\_order.**order\_type**string required  Order type enum  `limit` `market` |
| result.cancelled\_order.**quote\_id**string or null required  Quote ID if the trade was executed via RFQ |
| result.cancelled\_order.**signature**string required  Ethereum signature of the order |
| result.cancelled\_order.**signature\_expiry\_sec**integer required  Signature expiry timestamp |
| result.cancelled\_order.**signer**string required  Owner wallet address or registered session key that signed order |
| result.cancelled\_order.**subaccount\_id**integer required  Subaccount ID |
| result.cancelled\_order.**time\_in\_force**string required  Time in force enum  `gtc` `post_only` `fok` `ioc` |
| result.cancelled\_order.**replaced\_order\_id**string or null  If replaced, ID of the order that was replaced |
| result.cancelled\_order.**trigger\_price**string or null  (Required for trigger orders) Index or Market price to trigger order at |
| result.cancelled\_order.**trigger\_price\_type**string or null  (Required for trigger orders) Trigger with Index or Mark Price enum  `mark` `index` |
| result.cancelled\_order.**trigger\_reject\_message**string or null  (Required for trigger orders) Error message if error occured during trigger |
| result.cancelled\_order.**trigger\_type**string or null  (Required for trigger orders) Stop-loss or Take-profit. enum  `stoploss` `takeprofit` |
|  |
| result.**create\_order\_error**object or null  Optional. Returns error during new order creation |
| result.create\_order\_error.**code**integer required |
| result.create\_order\_error.**message**string required |
| result.create\_order\_error.**data**string or null |
|  |
| result.**order**object or null  New order that was created |
| result.order.**amount**string required  Order amount in units of the base |
| result.order.**average\_price**string required  Average fill price |
| result.order.**cancel\_reason**string required  If cancelled, reason behind order cancellation enum  `user_request` `mmp_trigger` `insufficient_margin` `signed_max_fee_too_low` `cancel_on_disconnect` `ioc_or_market_partial_fill` `session_key_deregistered` `subaccount_withdrawn` `compliance` `trigger_failed` `validation_failed` |
| result.order.**creation\_timestamp**integer required  Creation timestamp (in ms since Unix epoch) |
| result.order.**direction**string required  Order direction enum  `buy` `sell` |
| result.order.**filled\_amount**string required  Total filled amount for the order |
| result.order.**instrument\_name**string required  Instrument name |
| result.order.**is\_transfer**boolean required  Whether the order was generated through `private/transfer_position` |
| result.order.**label**string required  Optional user-defined label for the order |
| result.order.**last\_update\_timestamp**integer required  Last update timestamp (in ms since Unix epoch) |
| result.order.**limit\_price**string required  Limit price in quote currency |
| result.order.**max\_fee**string required  Max fee in units of the quote currency |
| result.order.**mmp**boolean required  Whether the order is tagged for market maker protections |
| result.order.**nonce**integer required  Unique nonce defined as  (e.g. 1695836058725001, where 001 is the random number) |
| result.order.**order\_fee**string required  Total order fee paid so far |
| result.order.**order\_id**string required  Order ID |
| result.order.**order\_status**string required  Order status enum  `open` `filled` `cancelled` `expired` `untriggered` |
| result.order.**order\_type**string required  Order type enum  `limit` `market` |
| result.order.**quote\_id**string or null required  Quote ID if the trade was executed via RFQ |
| result.order.**signature**string required  Ethereum signature of the order |
| result.order.**signature\_expiry\_sec**integer required  Signature expiry timestamp |
| result.order.**signer**string required  Owner wallet address or registered session key that signed order |
| result.order.**subaccount\_id**integer required  Subaccount ID |
| result.order.**time\_in\_force**string required  Time in force enum  `gtc` `post_only` `fok` `ioc` |
| result.order.**replaced\_order\_id**string or null  If replaced, ID of the order that was replaced |
| result.order.**trigger\_price**string or null  (Required for trigger orders) Index or Market price to trigger order at |
| result.order.**trigger\_price\_type**string or null  (Required for trigger orders) Trigger with Index or Mark Price enum  `mark` `index` |
| result.order.**trigger\_reject\_message**string or null  (Required for trigger orders) Error message if error occured during trigger |
| result.order.**trigger\_type**string or null  (Required for trigger orders) Stop-loss or Take-profit. enum  `stoploss` `takeprofit` |
|  |
| result.**trades**array of objects or null  List of trades executed by the created order |
| result.trades[].**direction**string required  Order direction enum  `buy` `sell` |
| result.trades[].**expected\_rebate**string required  Expected rebate for this trade |
| result.trades[].**index\_price**string required  Index price of the underlying at the time of the trade |
| result.trades[].**instrument\_name**string required  Instrument name |
| result.trades[].**is\_transfer**boolean required  Whether the trade was generated through `private/transfer_position` |
| result.trades[].**label**string required  Optional user-defined label for the order |
| result.trades[].**liquidity\_role**string required  Role of the user in the trade enum  `maker` `taker` |
| result.trades[].**mark\_price**string required  Mark price of the instrument at the time of the trade |
| result.trades[].**order\_id**string required  Order ID |
| result.trades[].**quote\_id**string or null required  Quote ID if the trade was executed via RFQ |
| result.trades[].**realized\_pnl**string required  Realized PnL for this trade |
| result.trades[].**realized\_pnl\_excl\_fees**string required  Realized PnL for this trade using cost accounting that excludes fees |
| result.trades[].**subaccount\_id**integer required  Subaccount ID |
| result.trades[].**timestamp**integer required  Trade timestamp (in ms since Unix epoch) |
| result.trades[].**trade\_amount**string required  Amount filled in this trade |
| result.trades[].**trade\_fee**string required  Fee for this trade |
| result.trades[].**trade\_id**string required  Trade ID |
| result.trades[].**trade\_price**string required  Price at which the trade was filled |
| result.trades[].**transaction\_id**string required  The transaction id of the related settlement transaction |
| result.trades[].**tx\_hash**string or null required  Blockchain transaction hash |
| result.trades[].**tx\_status**string required  Blockchain transaction status enum  `requested` `pending` `settled` `reverted` `ignored` `timed_out` |

### Example

ShellJavaScriptPython

```
{request_example_shell}

```

```
{request_example_javascript}

```

```
{request_example_python}

```

> The above command returns JSON structured like this:

JSON

```
{response_example_json}

```

---

## private-reset_mmp

**Title:** Private Reset_Mmp
**URL:** https://docs.derive.xyz/reference/private-reset_mmp

### Method Name

#### `private/reset_mmp`

Resets (unfreezes) the mmp state for a subaccount (optionally filtered by currency)  
Required minimum session key permission level is `account`

### Parameters

|  |
| --- |
| **subaccount\_id**integer required  Subaccount\_id for which to reset the mmp state |
| **currency**string  Currency to reset the mmp for. If not provided, resets all configs for the subaccount |

### Response

|  |
| --- |
| **id**string or integer required |
| **result**string required  The result of this method call, `ok` if successful enum  `ok` |

### Example

ShellJavaScriptPython

```
{request_example_shell}

```

```
{request_example_javascript}

```

```
{request_example_python}

```

> The above command returns JSON structured like this:

JSON

```
{response_example_json}

```

---

## private-rfq_get_best_quote

**Title:** Private Rfq_Get_Best_Quote
**URL:** https://docs.derive.xyz/reference/private-rfq_get_best_quote

### Method Name

#### `private/rfq_get_best_quote`

Performs a "dry run" on an RFQ, returning the estimated fee and whether the trade is expected to pass.  
  
Should any exception be raised in the process of evaluating the trade, a standard RPC error will be returned  
with the error details.  
Required minimum session key permission level is `read_only`

### Parameters

|  |
| --- |
| **subaccount\_id**integer required  Subaccount ID |
| **legs**array of objects required  RFQ legs |
| legs[].**amount**string required  Amount in units of the base |
| legs[].**direction**string required  Leg direction enum  `buy` `sell` |
| legs[].**instrument\_name**string required  Instrument name |
|  |
| **counterparties**array of strings  Optional list of market maker account addresses to request quotes from. If not supplied, all market makers who are approved as RFQ makers will be notified. |
| **direction**string  Planned execution direction (default `buy`) enum  `buy` `sell` |
| **label**string  Optional user-defined label for the RFQ |
| **max\_total\_cost**string  An optional max total cost for the RFQ. Only used when the RFQ sender executes as buyer. Polling endpoints and channels will ignore quotes where the total cost across all legs is above this value. Positive values mean the RFQ sender expects to pay $, negative mean the RFQ sender expects to receive $.This field is not disclosed to the market makers. |
| **min\_total\_cost**string  An optional min total cost for the RFQ. Only used when the RFQ sender executes as seller. Polling endpoints and channels will ignore quotes where the total cost across all legs is below this value. Positive values mean the RFQ sender expects to receive $, negative mean the RFQ sender expects to pay $.This field is not disclosed to the market makers. |
| **rfq\_id**string  RFQ ID to get best quote for. If not provided, will return estimates based on mark prices |

### Response

|  |
| --- |
| **id**string or integer required |
| **result**object required |
| result.**best\_quote**object or null required  Best quote for the RFQ (or null if RFQ is not created yet or quotes do not exist). This object should be used to sign a taker quote and call into `execute_quote` RPC. |
| result.best\_quote.**cancel\_reason**string required  Cancel reason, if any enum  `user_request` `insufficient_margin` `signed_max_fee_too_low` `mmp_trigger` `cancel_on_disconnect` `session_key_deregistered` `subaccount_withdrawn` `rfq_no_longer_open` `compliance` |
| result.best\_quote.**creation\_timestamp**integer required  Creation timestamp in ms since Unix epoch |
| result.best\_quote.**direction**string required  Quote direction enum  `buy` `sell` |
| result.best\_quote.**last\_update\_timestamp**integer required  Last update timestamp in ms since Unix epoch |
| result.best\_quote.**legs\_hash**string required  Hash of the legs of the best quote to be signed by the taker. |
| result.best\_quote.**liquidity\_role**string required  Liquidity role enum  `maker` `taker` |
| result.best\_quote.**quote\_id**string required  Quote ID |
| result.best\_quote.**rfq\_id**string required  RFQ ID |
| result.best\_quote.**status**string required  Status enum  `open` `filled` `cancelled` `expired` |
| result.best\_quote.**subaccount\_id**integer required  Subaccount ID |
| result.best\_quote.**tx\_hash**string or null required  Blockchain transaction hash (only for executed quotes) |
| result.best\_quote.**tx\_status**string or null required  Blockchain transaction status (only for executed quotes) enum  `requested` `pending` `settled` `reverted` `ignored` `timed_out` |
| result.best\_quote.**wallet**string required  Wallet address of the quote sender |
| result.best\_quote.**legs**array of objects required  Quote legs |
| result.best\_quote.legs[].**amount**string required  Amount in units of the base |
| result.best\_quote.legs[].**direction**string required  Leg direction enum  `buy` `sell` |
| result.best\_quote.legs[].**instrument\_name**string required  Instrument name |
| result.best\_quote.legs[].**price**string required  Leg price |
|  |
| result.**down\_liquidation\_price**string or null required  Liquidation price if the trade were to be filled and the market moves down. |
| result.**estimated\_fee**string required  An estimate for how much the user will pay in fees ($ for the whole trade). |
| result.**estimated\_realized\_pnl**string required  An estimate for the realized PnL of the trade. |
| result.**estimated\_realized\_pnl\_excl\_fees**string required  An estimate for the realized PnL of the trade. with cost basis calculated without considering fees. |
| result.**estimated\_total\_cost**string required  An estimate for the total $ cost of the trade. |
| result.**invalid\_reason**string or null required  Reason for the RFQ being invalid, if any. enum  `Account is currently under maintenance margin requirements, trading is frozen.` `This order would cause account to fall under maintenance margin requirements.` `Insufficient buying power, only a single risk-reducing open order is allowed.` `Insufficient buying power, consider reducing order size.` `Insufficient buying power, consider reducing order size or canceling other orders.` `Consider canceling other limit orders or using IOC, FOK, or market orders. This order is risk-reducing, but if filled with other open orders, buying power might be insufficient.` `Insufficient buying power.` |
| result.**is\_valid**boolean required  `True` if RFQ is expected to pass margin requirements. |
| result.**post\_initial\_margin**string required  User's hypothetical margin balance if the trade were to get executed. |
| result.**post\_liquidation\_price**string or null required  Liquidation price if the trade were to be filled. If both upside and downside liquidation prices exist, returns the closest one to the current index price. |
| result.**pre\_initial\_margin**string required  User's initial margin balance before the trade. |
| result.**suggested\_max\_fee**string required  Recommended value for `max_fee` of the trade. |
| result.**up\_liquidation\_price**string or null required  Liquidation price if the trade were to be filled and the market moves up. |

### Example

ShellJavaScriptPython

```
{request_example_shell}

```

```
{request_example_javascript}

```

```
{request_example_python}

```

> The above command returns JSON structured like this:

JSON

```
{response_example_json}

```

---

## private-send_quote

**Title:** Private Send_Quote
**URL:** https://docs.derive.xyz/reference/private-send_quote

### Method Name

#### `private/send_quote`

Sends a quote in response to an RFQ request.  
The legs supplied in the parameters must exactly match those in the RFQ.  
Required minimum session key permission level is `admin`

### Parameters

|  |
| --- |
| **direction**string required  Quote direction, `buy` means trading each leg at its direction, `sell` means trading each leg in the opposite direction. enum  `buy` `sell` |
| **max\_fee**string required  Max fee ($ for the full trade). Request will be rejected if the supplied max fee is below the estimated fee for this trade. |
| **nonce**integer required  Unique nonce defined as a concatenated `UTC timestamp in ms` and `random number up to 6 digits` (e.g. 1695836058725001, where 001 is the random number) |
| **rfq\_id**string required  RFQ ID the quote is for |
| **signature**string required  Ethereum signature of the quote |
| **signature\_expiry\_sec**integer required  Unix timestamp in seconds. Expiry MUST be at least 310 seconds from now. Once time till signature expiry reaches 300 seconds, the quote will be considered expired. This buffer is meant to ensure the trade can settle on chain in case of a blockchain congestion. |
| **signer**string required  Owner wallet address or registered session key that signed the quote |
| **subaccount\_id**integer required  Subaccount ID |
| **legs**array of objects required  Quote legs |
| legs[].**amount**string required  Amount in units of the base |
| legs[].**direction**string required  Leg direction enum  `buy` `sell` |
| legs[].**instrument\_name**string required  Instrument name |
| legs[].**price**string required  Leg price |
|  |
| **label**string  Optional user-defined label for the quote |
| **mmp**boolean  Whether the quote is tagged for market maker protections (default false) |

### Response

|  |
| --- |
| **id**string or integer required |
| **result**object required |
| result.**cancel\_reason**string required  Cancel reason, if any enum  `user_request` `insufficient_margin` `signed_max_fee_too_low` `mmp_trigger` `cancel_on_disconnect` `session_key_deregistered` `subaccount_withdrawn` `rfq_no_longer_open` `compliance` |
| result.**creation\_timestamp**integer required  Creation timestamp in ms since Unix epoch |
| result.**direction**string required  Quote direction enum  `buy` `sell` |
| result.**fee**string required  Fee paid for this quote (if executed) |
| result.**is\_transfer**boolean required  Whether the order was generated through `private/transfer_position` |
| result.**label**string required  User-defined label, if any |
| result.**last\_update\_timestamp**integer required  Last update timestamp in ms since Unix epoch |
| result.**legs\_hash**string required  Hash of the legs of the best quote to be signed by the taker. |
| result.**liquidity\_role**string required  Liquidity role enum  `maker` `taker` |
| result.**max\_fee**string required  Signed max fee |
| result.**mmp**boolean required  Whether the quote is tagged for market maker protections (default false) |
| result.**nonce**integer required  Nonce |
| result.**quote\_id**string required  Quote ID |
| result.**rfq\_id**string required  RFQ ID |
| result.**signature**string required  Ethereum signature of the quote |
| result.**signature\_expiry\_sec**integer required  Unix timestamp in seconds |
| result.**signer**string required  Owner wallet address or registered session key that signed the quote |
| result.**status**string required  Status enum  `open` `filled` `cancelled` `expired` |
| result.**subaccount\_id**integer required  Subaccount ID |
| result.**tx\_hash**string or null required  Blockchain transaction hash (only for executed quotes) |
| result.**tx\_status**string or null required  Blockchain transaction status (only for executed quotes) enum  `requested` `pending` `settled` `reverted` `ignored` `timed_out` |
| result.**legs**array of objects required  Quote legs |
| result.legs[].**amount**string required  Amount in units of the base |
| result.legs[].**direction**string required  Leg direction enum  `buy` `sell` |
| result.legs[].**instrument\_name**string required  Instrument name |
| result.legs[].**price**string required  Leg price |

### Example

ShellJavaScriptPython

```
{request_example_shell}

```

```
{request_example_javascript}

```

```
{request_example_python}

```

> The above command returns JSON structured like this:

JSON

```
{response_example_json}

```

---

## private-send_rfq

**Title:** Private Send_Rfq
**URL:** https://docs.derive.xyz/reference/private-send_rfq

### Method Name

#### `private/send_rfq`

Requests two-sided quotes from participating market makers.  
Required minimum session key permission level is `account`

### Parameters

|  |
| --- |
| **subaccount\_id**integer required  Subaccount ID |
| **legs**array of objects required  RFQ legs |
| legs[].**amount**string required  Amount in units of the base |
| legs[].**direction**string required  Leg direction enum  `buy` `sell` |
| legs[].**instrument\_name**string required  Instrument name |
|  |
| **counterparties**array of strings  Optional list of market maker account addresses to request quotes from. If not supplied, all market makers who are approved as RFQ makers will be notified. |
| **label**string  Optional user-defined label for the RFQ |
| **max\_total\_cost**string  An optional max total cost for the RFQ. Only used when the RFQ sender executes as buyer. Polling endpoints and channels will ignore quotes where the total cost across all legs is above this value. Positive values mean the RFQ sender expects to pay $, negative mean the RFQ sender expects to receive $.This field is not disclosed to the market makers. |
| **min\_total\_cost**string  An optional min total cost for the RFQ. Only used when the RFQ sender executes as seller. Polling endpoints and channels will ignore quotes where the total cost across all legs is below this value. Positive values mean the RFQ sender expects to receive $, negative mean the RFQ sender expects to pay $.This field is not disclosed to the market makers. |

### Response

|  |
| --- |
| **id**string or integer required |
| **result**object required |
| result.**ask\_total\_cost**string or null required  Ask total cost for the RFQ implied from orderbook (as `sell`) |
| result.**bid\_total\_cost**string or null required  Bid total cost for the RFQ implied from orderbook (as `buy`) |
| result.**cancel\_reason**string required  Cancel reason, if any enum  `user_request` `insufficient_margin` `signed_max_fee_too_low` `mmp_trigger` `cancel_on_disconnect` `session_key_deregistered` `subaccount_withdrawn` `rfq_no_longer_open` `compliance` |
| result.**counterparties**array of strings or null required  List of requested counterparties, if applicable |
| result.**creation\_timestamp**integer required  Creation timestamp in ms since Unix epoch |
| result.**filled\_direction**string or null required  Direction at which the RFQ was filled (only if filled) enum  `buy` `sell` |
| result.**label**string required  User-defined label, if any |
| result.**last\_update\_timestamp**integer required  Last update timestamp in ms since Unix epoch |
| result.**mark\_total\_cost**string or null required  Mark total cost for the RFQ (assuming `buy` direction) |
| result.**max\_total\_cost**string or null required  Max total cost for the RFQ |
| result.**min\_total\_cost**string or null required  Min total cost for the RFQ |
| result.**rfq\_id**string required  RFQ ID |
| result.**status**string required  Status enum  `open` `filled` `cancelled` `expired` |
| result.**subaccount\_id**integer required  Subaccount ID |
| result.**total\_cost**string or null required  Total cost for the RFQ (only if filled) |
| result.**valid\_until**integer required  RFQ expiry timestamp in ms since Unix epoch |
| result.**legs**array of objects required  RFQ legs |
| result.legs[].**amount**string required  Amount in units of the base |
| result.legs[].**direction**string required  Leg direction enum  `buy` `sell` |
| result.legs[].**instrument\_name**string required  Instrument name |

### Example

ShellJavaScriptPython

```
{request_example_shell}

```

```
{request_example_javascript}

```

```
{request_example_python}

```

> The above command returns JSON structured like this:

JSON

```
{response_example_json}

```

---

## private-session_keys

**Title:** Private Session_Keys
**URL:** https://docs.derive.xyz/reference/private-session_keys

### Method Name

#### `private/session_keys`

Required minimum session key permission level is `read_only`

### Parameters

|  |
| --- |
| **wallet**string required  Ethereum wallet address of account |

### Response

|  |
| --- |
| **id**string or integer required |
| **result**object required |
| result.**public\_session\_keys**array of objects required  List of session keys (includes unactivated and expired keys) |
| result.public\_session\_keys[].**expiry\_sec**integer required  Session key expiry timestamp in sec |
| result.public\_session\_keys[].**label**string required  User-defined session key label |
| result.public\_session\_keys[].**public\_session\_key**string required  Public session key address (Ethereum EOA) |
| result.public\_session\_keys[].**scope**string required  Session key permission level scope |
| result.public\_session\_keys[].**ip\_whitelist**array of strings required  List of whitelisted IPs, if empty then any IP is allowed. |

### Example

ShellJavaScriptPython

```
{request_example_shell}

```

```
{request_example_javascript}

```

```
{request_example_python}

```

> The above command returns JSON structured like this:

JSON

```
{response_example_json}

```

---

## private-set_cancel_on_disconnect

**Title:** Private Set_Cancel_On_Disconnect
**URL:** https://docs.derive.xyz/reference/private-set_cancel_on_disconnect

### Method Name

#### `private/set_cancel_on_disconnect`

Enables cancel on disconnect for the account  
Required minimum session key permission level is `account`

### Parameters

|  |
| --- |
| **enabled**boolean required  Whether to enable or disable cancel on disconnect |
| **wallet**string required  Public key (wallet) of the account |

### Response

|  |
| --- |
| **id**string or integer required |
| **result**string required   enum  `ok` |

### Example

ShellJavaScriptPython

```
{request_example_shell}

```

```
{request_example_javascript}

```

```
{request_example_python}

```

> The above command returns JSON structured like this:

JSON

```
{response_example_json}

```

---

## private-set_mmp_config

**Title:** Private Set_Mmp_Config
**URL:** https://docs.derive.xyz/reference/private-set_mmp_config

### Method Name

#### `private/set_mmp_config`

Set the mmp config for the subaccount and currency  
Required minimum session key permission level is `account`

### Parameters

|  |
| --- |
| **currency**string required  Currency of this mmp config |
| **mmp\_frozen\_time**integer required  Time interval in ms setting how long the subaccount is frozen after an mmp trigger, if 0 then a manual reset would be required via private/reset\_mmp |
| **mmp\_interval**integer required  Time interval in ms over which the limits are monotored, if 0 then mmp is disabled |
| **subaccount\_id**integer required  Subaccount\_id for which to set the config |
| **mmp\_amount\_limit**string  Maximum total order amount that can be traded within the mmp\_interval across all instruments of the provided currency. The amounts are not netted, so a filled bid of 1 and a filled ask of 2 would count as 3. Default: 0 (no limit) |
| **mmp\_delta\_limit**string  Maximum total delta that can be traded within the mmp\_interval across all instruments of the provided currency. This quantity is netted, so a filled order with +1 delta and a filled order with -2 delta would count as -1 Default: 0 (no limit) |

### Response

|  |
| --- |
| **id**string or integer required |
| **result**object required |
| result.**currency**string required  Currency of this mmp config |
| result.**mmp\_frozen\_time**integer required  Time interval in ms setting how long the subaccount is frozen after an mmp trigger, if 0 then a manual reset would be required via private/reset\_mmp |
| result.**mmp\_interval**integer required  Time interval in ms over which the limits are monotored, if 0 then mmp is disabled |
| result.**subaccount\_id**integer required  Subaccount\_id for which to set the config |
| result.**mmp\_amount\_limit**string  Maximum total order amount that can be traded within the mmp\_interval across all instruments of the provided currency. The amounts are not netted, so a filled bid of 1 and a filled ask of 2 would count as 3. Default: 0 (no limit) |
| result.**mmp\_delta\_limit**string  Maximum total delta that can be traded within the mmp\_interval across all instruments of the provided currency. This quantity is netted, so a filled order with +1 delta and a filled order with -2 delta would count as -1 Default: 0 (no limit) |

### Example

ShellJavaScriptPython

```
{request_example_shell}

```

```
{request_example_javascript}

```

```
{request_example_python}

```

> The above command returns JSON structured like this:

JSON

```
{response_example_json}

```

---

## private-transfer_erc20

**Title:** Private Transfer_Erc20
**URL:** https://docs.derive.xyz/reference/private-transfer_erc20

### Method Name

#### `private/transfer_erc20`

Transfer ERC20 assets from one subaccount to another (e.g. USDC or ETH).  
  
For transfering positions (e.g. options or perps), use `private/transfer_position` instead.  
Required minimum session key permission level is `admin`

### Parameters

|  |
| --- |
| **recipient\_subaccount\_id**integer required  Subaccount\_id of the recipient |
| **subaccount\_id**integer required  Subaccount\_id |
| **recipient\_details**object required  Details of the recipient |
| recipient\_details.**nonce**integer required  Unique nonce defined as  (e.g. 1695836058725001, where 001 is the random number) |
| recipient\_details.**signature**string required  Ethereum signature of the transfer |
| recipient\_details.**signature\_expiry\_sec**integer required  Unix timestamp in seconds. Expiry MUST be >5min from now |
| recipient\_details.**signer**string required  Ethereum wallet address that is signing the transfer |
|  |
| **sender\_details**object required  Details of the sender |
| sender\_details.**nonce**integer required  Unique nonce defined as  (e.g. 1695836058725001, where 001 is the random number) |
| sender\_details.**signature**string required  Ethereum signature of the transfer |
| sender\_details.**signature\_expiry\_sec**integer required  Unix timestamp in seconds. Expiry MUST be >5min from now |
| sender\_details.**signer**string required  Ethereum wallet address that is signing the transfer |
|  |
| **transfer**object required  Transfer details |
| transfer.**address**string required  Ethereum address of the asset being transferred |
| transfer.**amount**string required  Amount to transfer |
| transfer.**sub\_id**integer required  Sub ID of the asset being transferred |

### Response

|  |
| --- |
| **id**string or integer required |
| **result**object required |
| result.**status**string required  `requested` |
| result.**transaction\_id**string required  Transaction id of the transfer |

### Example

ShellJavaScriptPython

```
{request_example_shell}

```

```
{request_example_javascript}

```

```
{request_example_python}

```

> The above command returns JSON structured like this:

JSON

```
{response_example_json}

```

---

## private-transfer_position

**Title:** Private Transfer_Position
**URL:** https://docs.derive.xyz/reference/private-transfer_position

### Method Name

#### `private/transfer_position`

Transfers a positions from one subaccount to another, owned by the same wallet.  
  
The transfer is executed as a pair of orders crossing each other.  
The maker order is created first, followed by a taker order crossing it.  
The order amounts, limit prices and instrument name must be the same for both orders.  
Fee is not charged and a zero `max_fee` must be signed.  
The maker order is forcibly considered to be `reduce_only`, meaning it can only reduce the position size.  
  
History: For position transfer history, use the `private/get_trade_history` RPC (not `private/get_erc20_transfer_history`).  
Required minimum session key permission level is `admin`

### Parameters

|  |
| --- |
| **wallet**string required  Public key (wallet) of the account |
| **maker\_params**object required  Maker order parameters and signature. Maximum transfer amount is limited by the size of the maker position. Transfers that increase the maker's position size are not allowed. |
| maker\_params.**amount**string required  Order amount in units of the base |
| maker\_params.**direction**string required  Order direction enum  `buy` `sell` |
| maker\_params.**instrument\_name**string required  Instrument name |
| maker\_params.**limit\_price**string required  Limit price in quote currency. This field is still required for market orders because it is a component of the signature. However, market orders will not leave a resting order in the book in case of a partial fill. |
| maker\_params.**max\_fee**string required  Max fee per unit of volume, denominated in units of the quote currency (usually USDC).Order will be rejected if the supplied max fee is below the estimated fee for this order. |
| maker\_params.**nonce**integer required  Unique nonce defined as  (e.g. 1695836058725001, where 001 is the random number).Note, using a random number beyond 3 digits will cause JSON serialization to fail. |
| maker\_params.**signature**string required  Ethereum signature of the order |
| maker\_params.**signature\_expiry\_sec**integer required  Unix timestamp in seconds. Order signature becomes invalid after this time, and the system will cancel the order.Expiry MUST be at least 5 min from now. |
| maker\_params.**signer**string required  Owner wallet address or registered session key that signed order |
| maker\_params.**subaccount\_id**integer required  Subaccount ID |
|  |
| **taker\_params**object required  Taker order parameters and signature |
| taker\_params.**amount**string required  Order amount in units of the base |
| taker\_params.**direction**string required  Order direction enum  `buy` `sell` |
| taker\_params.**instrument\_name**string required  Instrument name |
| taker\_params.**limit\_price**string required  Limit price in quote currency. This field is still required for market orders because it is a component of the signature. However, market orders will not leave a resting order in the book in case of a partial fill. |
| taker\_params.**max\_fee**string required  Max fee per unit of volume, denominated in units of the quote currency (usually USDC).Order will be rejected if the supplied max fee is below the estimated fee for this order. |
| taker\_params.**nonce**integer required  Unique nonce defined as  (e.g. 1695836058725001, where 001 is the random number).Note, using a random number beyond 3 digits will cause JSON serialization to fail. |
| taker\_params.**signature**string required  Ethereum signature of the order |
| taker\_params.**signature\_expiry\_sec**integer required  Unix timestamp in seconds. Order signature becomes invalid after this time, and the system will cancel the order.Expiry MUST be at least 5 min from now. |
| taker\_params.**signer**string required  Owner wallet address or registered session key that signed order |
| taker\_params.**subaccount\_id**integer required  Subaccount ID |

### Response

|  |
| --- |
| **id**string or integer required |
| **result**object required |
| result.**maker\_order**object required |
| result.maker\_order.**amount**string required  Order amount in units of the base |
| result.maker\_order.**average\_price**string required  Average fill price |
| result.maker\_order.**cancel\_reason**string required  If cancelled, reason behind order cancellation enum  `user_request` `mmp_trigger` `insufficient_margin` `signed_max_fee_too_low` `cancel_on_disconnect` `ioc_or_market_partial_fill` `session_key_deregistered` `subaccount_withdrawn` `compliance` `trigger_failed` `validation_failed` |
| result.maker\_order.**creation\_timestamp**integer required  Creation timestamp (in ms since Unix epoch) |
| result.maker\_order.**direction**string required  Order direction enum  `buy` `sell` |
| result.maker\_order.**filled\_amount**string required  Total filled amount for the order |
| result.maker\_order.**instrument\_name**string required  Instrument name |
| result.maker\_order.**is\_transfer**boolean required  Whether the order was generated through `private/transfer_position` |
| result.maker\_order.**label**string required  Optional user-defined label for the order |
| result.maker\_order.**last\_update\_timestamp**integer required  Last update timestamp (in ms since Unix epoch) |
| result.maker\_order.**limit\_price**string required  Limit price in quote currency |
| result.maker\_order.**max\_fee**string required  Max fee in units of the quote currency |
| result.maker\_order.**mmp**boolean required  Whether the order is tagged for market maker protections |
| result.maker\_order.**nonce**integer required  Unique nonce defined as  (e.g. 1695836058725001, where 001 is the random number) |
| result.maker\_order.**order\_fee**string required  Total order fee paid so far |
| result.maker\_order.**order\_id**string required  Order ID |
| result.maker\_order.**order\_status**string required  Order status enum  `open` `filled` `cancelled` `expired` `untriggered` |
| result.maker\_order.**order\_type**string required  Order type enum  `limit` `market` |
| result.maker\_order.**quote\_id**string or null required  Quote ID if the trade was executed via RFQ |
| result.maker\_order.**signature**string required  Ethereum signature of the order |
| result.maker\_order.**signature\_expiry\_sec**integer required  Signature expiry timestamp |
| result.maker\_order.**signer**string required  Owner wallet address or registered session key that signed order |
| result.maker\_order.**subaccount\_id**integer required  Subaccount ID |
| result.maker\_order.**time\_in\_force**string required  Time in force enum  `gtc` `post_only` `fok` `ioc` |
| result.maker\_order.**replaced\_order\_id**string or null  If replaced, ID of the order that was replaced |
| result.maker\_order.**trigger\_price**string or null  (Required for trigger orders) Index or Market price to trigger order at |
| result.maker\_order.**trigger\_price\_type**string or null  (Required for trigger orders) Trigger with Index or Mark Price enum  `mark` `index` |
| result.maker\_order.**trigger\_reject\_message**string or null  (Required for trigger orders) Error message if error occured during trigger |
| result.maker\_order.**trigger\_type**string or null  (Required for trigger orders) Stop-loss or Take-profit. enum  `stoploss` `takeprofit` |
|  |
| result.**maker\_trade**object required |
| result.maker\_trade.**direction**string required  Order direction enum  `buy` `sell` |
| result.maker\_trade.**expected\_rebate**string required  Expected rebate for this trade |
| result.maker\_trade.**index\_price**string required  Index price of the underlying at the time of the trade |
| result.maker\_trade.**instrument\_name**string required  Instrument name |
| result.maker\_trade.**is\_transfer**boolean required  Whether the trade was generated through `private/transfer_position` |
| result.maker\_trade.**label**string required  Optional user-defined label for the order |
| result.maker\_trade.**liquidity\_role**string required  Role of the user in the trade enum  `maker` `taker` |
| result.maker\_trade.**mark\_price**string required  Mark price of the instrument at the time of the trade |
| result.maker\_trade.**order\_id**string required  Order ID |
| result.maker\_trade.**quote\_id**string or null required  Quote ID if the trade was executed via RFQ |
| result.maker\_trade.**realized\_pnl**string required  Realized PnL for this trade |
| result.maker\_trade.**realized\_pnl\_excl\_fees**string required  Realized PnL for this trade using cost accounting that excludes fees |
| result.maker\_trade.**subaccount\_id**integer required  Subaccount ID |
| result.maker\_trade.**timestamp**integer required  Trade timestamp (in ms since Unix epoch) |
| result.maker\_trade.**trade\_amount**string required  Amount filled in this trade |
| result.maker\_trade.**trade\_fee**string required  Fee for this trade |
| result.maker\_trade.**trade\_id**string required  Trade ID |
| result.maker\_trade.**trade\_price**string required  Price at which the trade was filled |
| result.maker\_trade.**transaction\_id**string required  The transaction id of the related settlement transaction |
| result.maker\_trade.**tx\_hash**string or null required  Blockchain transaction hash |
| result.maker\_trade.**tx\_status**string required  Blockchain transaction status enum  `requested` `pending` `settled` `reverted` `ignored` `timed_out` |
|  |
| result.**taker\_order**object required |
| result.taker\_order.**amount**string required  Order amount in units of the base |
| result.taker\_order.**average\_price**string required  Average fill price |
| result.taker\_order.**cancel\_reason**string required  If cancelled, reason behind order cancellation enum  `user_request` `mmp_trigger` `insufficient_margin` `signed_max_fee_too_low` `cancel_on_disconnect` `ioc_or_market_partial_fill` `session_key_deregistered` `subaccount_withdrawn` `compliance` `trigger_failed` `validation_failed` |
| result.taker\_order.**creation\_timestamp**integer required  Creation timestamp (in ms since Unix epoch) |
| result.taker\_order.**direction**string required  Order direction enum  `buy` `sell` |
| result.taker\_order.**filled\_amount**string required  Total filled amount for the order |
| result.taker\_order.**instrument\_name**string required  Instrument name |
| result.taker\_order.**is\_transfer**boolean required  Whether the order was generated through `private/transfer_position` |
| result.taker\_order.**label**string required  Optional user-defined label for the order |
| result.taker\_order.**last\_update\_timestamp**integer required  Last update timestamp (in ms since Unix epoch) |
| result.taker\_order.**limit\_price**string required  Limit price in quote currency |
| result.taker\_order.**max\_fee**string required  Max fee in units of the quote currency |
| result.taker\_order.**mmp**boolean required  Whether the order is tagged for market maker protections |
| result.taker\_order.**nonce**integer required  Unique nonce defined as  (e.g. 1695836058725001, where 001 is the random number) |
| result.taker\_order.**order\_fee**string required  Total order fee paid so far |
| result.taker\_order.**order\_id**string required  Order ID |
| result.taker\_order.**order\_status**string required  Order status enum  `open` `filled` `cancelled` `expired` `untriggered` |
| result.taker\_order.**order\_type**string required  Order type enum  `limit` `market` |
| result.taker\_order.**quote\_id**string or null required  Quote ID if the trade was executed via RFQ |
| result.taker\_order.**signature**string required  Ethereum signature of the order |
| result.taker\_order.**signature\_expiry\_sec**integer required  Signature expiry timestamp |
| result.taker\_order.**signer**string required  Owner wallet address or registered session key that signed order |
| result.taker\_order.**subaccount\_id**integer required  Subaccount ID |
| result.taker\_order.**time\_in\_force**string required  Time in force enum  `gtc` `post_only` `fok` `ioc` |
| result.taker\_order.**replaced\_order\_id**string or null  If replaced, ID of the order that was replaced |
| result.taker\_order.**trigger\_price**string or null  (Required for trigger orders) Index or Market price to trigger order at |
| result.taker\_order.**trigger\_price\_type**string or null  (Required for trigger orders) Trigger with Index or Mark Price enum  `mark` `index` |
| result.taker\_order.**trigger\_reject\_message**string or null  (Required for trigger orders) Error message if error occured during trigger |
| result.taker\_order.**trigger\_type**string or null  (Required for trigger orders) Stop-loss or Take-profit. enum  `stoploss` `takeprofit` |
|  |
| result.**taker\_trade**object required |
| result.taker\_trade.**direction**string required  Order direction enum  `buy` `sell` |
| result.taker\_trade.**expected\_rebate**string required  Expected rebate for this trade |
| result.taker\_trade.**index\_price**string required  Index price of the underlying at the time of the trade |
| result.taker\_trade.**instrument\_name**string required  Instrument name |
| result.taker\_trade.**is\_transfer**boolean required  Whether the trade was generated through `private/transfer_position` |
| result.taker\_trade.**label**string required  Optional user-defined label for the order |
| result.taker\_trade.**liquidity\_role**string required  Role of the user in the trade enum  `maker` `taker` |
| result.taker\_trade.**mark\_price**string required  Mark price of the instrument at the time of the trade |
| result.taker\_trade.**order\_id**string required  Order ID |
| result.taker\_trade.**quote\_id**string or null required  Quote ID if the trade was executed via RFQ |
| result.taker\_trade.**realized\_pnl**string required  Realized PnL for this trade |
| result.taker\_trade.**realized\_pnl\_excl\_fees**string required  Realized PnL for this trade using cost accounting that excludes fees |
| result.taker\_trade.**subaccount\_id**integer required  Subaccount ID |
| result.taker\_trade.**timestamp**integer required  Trade timestamp (in ms since Unix epoch) |
| result.taker\_trade.**trade\_amount**string required  Amount filled in this trade |
| result.taker\_trade.**trade\_fee**string required  Fee for this trade |
| result.taker\_trade.**trade\_id**string required  Trade ID |
| result.taker\_trade.**trade\_price**string required  Price at which the trade was filled |
| result.taker\_trade.**transaction\_id**string required  The transaction id of the related settlement transaction |
| result.taker\_trade.**tx\_hash**string or null required  Blockchain transaction hash |
| result.taker\_trade.**tx\_status**string required  Blockchain transaction status enum  `requested` `pending` `settled` `reverted` `ignored` `timed_out` |

### Example

ShellJavaScriptPython

```
{request_example_shell}

```

```
{request_example_javascript}

```

```
{request_example_python}

```

> The above command returns JSON structured like this:

JSON

```
{response_example_json}

```

---

## private-transfer_positions

**Title:** Private Transfer_Positions
**URL:** https://docs.derive.xyz/reference/private-transfer_positions

### Method Name

#### `private/transfer_positions`

Transfers multiple positions from one subaccount to another, owned by the same wallet.  
  
The transfer is executed as a an RFQ. A mock RFQ is first created from the taker parameters, followed by a maker quote and a taker execute.  
The leg amounts, prices and instrument name must be the same in both param payloads.  
Fee is not charged and a zero `max_fee` must be signed.  
Every leg in the transfer must be a position reduction for either maker or taker (or both).  
  
History: for position transfer history, use the `private/get_trade_history` RPC (not `private/get_erc20_transfer_history`).  
Required minimum session key permission level is `admin`

### Parameters

|  |
| --- |
| **wallet**string required  Public key (wallet) of the account |
| **maker\_params**object required  Maker quote parameters and signature |
| maker\_params.**direction**string required  Quote direction, `buy` means trading each leg at its direction, `sell` means trading each leg in the opposite direction. enum  `buy` `sell` |
| maker\_params.**max\_fee**string required  Max fee ($ for the full trade). Request will be rejected if the supplied max fee is below the estimated fee for this trade. |
| maker\_params.**nonce**integer required  Unique nonce defined as a concatenated `UTC timestamp in ms` and `random number up to 6 digits` (e.g. 1695836058725001, where 001 is the random number) |
| maker\_params.**signature**string required  Ethereum signature of the quote |
| maker\_params.**signature\_expiry\_sec**integer required  Unix timestamp in seconds. Expiry MUST be at least 310 seconds from now. Once time till signature expiry reaches 300 seconds, the quote will be considered expired. This buffer is meant to ensure the trade can settle on chain in case of a blockchain congestion. |
| maker\_params.**signer**string required  Owner wallet address or registered session key that signed the quote |
| maker\_params.**subaccount\_id**integer required  Subaccount ID |
| maker\_params.**legs**array of objects required  Quote legs |
| maker\_params.legs[].**amount**string required  Amount in units of the base |
| maker\_params.legs[].**direction**string required  Leg direction enum  `buy` `sell` |
| maker\_params.legs[].**instrument\_name**string required  Instrument name |
| maker\_params.legs[].**price**string required  Leg price |
|  |
| **taker\_params**object required  Taker quote execution parameters and signature |
| taker\_params.**direction**string required  Quote direction, `buy` means trading each leg at its direction, `sell` means trading each leg in the opposite direction. enum  `buy` `sell` |
| taker\_params.**max\_fee**string required  Max fee ($ for the full trade). Request will be rejected if the supplied max fee is below the estimated fee for this trade. |
| taker\_params.**nonce**integer required  Unique nonce defined as a concatenated `UTC timestamp in ms` and `random number up to 6 digits` (e.g. 1695836058725001, where 001 is the random number) |
| taker\_params.**signature**string required  Ethereum signature of the quote |
| taker\_params.**signature\_expiry\_sec**integer required  Unix timestamp in seconds. Expiry MUST be at least 310 seconds from now. Once time till signature expiry reaches 300 seconds, the quote will be considered expired. This buffer is meant to ensure the trade can settle on chain in case of a blockchain congestion. |
| taker\_params.**signer**string required  Owner wallet address or registered session key that signed the quote |
| taker\_params.**subaccount\_id**integer required  Subaccount ID |
| taker\_params.**legs**array of objects required  Quote legs |
| taker\_params.legs[].**amount**string required  Amount in units of the base |
| taker\_params.legs[].**direction**string required  Leg direction enum  `buy` `sell` |
| taker\_params.legs[].**instrument\_name**string required  Instrument name |
| taker\_params.legs[].**price**string required  Leg price |

### Response

|  |
| --- |
| **id**string or integer required |
| **result**object required |
| result.**maker\_quote**object required  Created maker-side quote object |
| result.maker\_quote.**cancel\_reason**string required  Cancel reason, if any enum  `user_request` `insufficient_margin` `signed_max_fee_too_low` `mmp_trigger` `cancel_on_disconnect` `session_key_deregistered` `subaccount_withdrawn` `rfq_no_longer_open` `compliance` |
| result.maker\_quote.**creation\_timestamp**integer required  Creation timestamp in ms since Unix epoch |
| result.maker\_quote.**direction**string required  Quote direction enum  `buy` `sell` |
| result.maker\_quote.**fee**string required  Fee paid for this quote (if executed) |
| result.maker\_quote.**is\_transfer**boolean required  Whether the order was generated through `private/transfer_position` |
| result.maker\_quote.**label**string required  User-defined label, if any |
| result.maker\_quote.**last\_update\_timestamp**integer required  Last update timestamp in ms since Unix epoch |
| result.maker\_quote.**legs\_hash**string required  Hash of the legs of the best quote to be signed by the taker. |
| result.maker\_quote.**liquidity\_role**string required  Liquidity role enum  `maker` `taker` |
| result.maker\_quote.**max\_fee**string required  Signed max fee |
| result.maker\_quote.**mmp**boolean required  Whether the quote is tagged for market maker protections (default false) |
| result.maker\_quote.**nonce**integer required  Nonce |
| result.maker\_quote.**quote\_id**string required  Quote ID |
| result.maker\_quote.**rfq\_id**string required  RFQ ID |
| result.maker\_quote.**signature**string required  Ethereum signature of the quote |
| result.maker\_quote.**signature\_expiry\_sec**integer required  Unix timestamp in seconds |
| result.maker\_quote.**signer**string required  Owner wallet address or registered session key that signed the quote |
| result.maker\_quote.**status**string required  Status enum  `open` `filled` `cancelled` `expired` |
| result.maker\_quote.**subaccount\_id**integer required  Subaccount ID |
| result.maker\_quote.**tx\_hash**string or null required  Blockchain transaction hash (only for executed quotes) |
| result.maker\_quote.**tx\_status**string or null required  Blockchain transaction status (only for executed quotes) enum  `requested` `pending` `settled` `reverted` `ignored` `timed_out` |
| result.maker\_quote.**legs**array of objects required  Quote legs |
| result.maker\_quote.legs[].**amount**string required  Amount in units of the base |
| result.maker\_quote.legs[].**direction**string required  Leg direction enum  `buy` `sell` |
| result.maker\_quote.legs[].**instrument\_name**string required  Instrument name |
| result.maker\_quote.legs[].**price**string required  Leg price |
|  |
| result.**taker\_quote**object required  Created taker-side quote object |
| result.taker\_quote.**cancel\_reason**string required  Cancel reason, if any enum  `user_request` `insufficient_margin` `signed_max_fee_too_low` `mmp_trigger` `cancel_on_disconnect` `session_key_deregistered` `subaccount_withdrawn` `rfq_no_longer_open` `compliance` |
| result.taker\_quote.**creation\_timestamp**integer required  Creation timestamp in ms since Unix epoch |
| result.taker\_quote.**direction**string required  Quote direction enum  `buy` `sell` |
| result.taker\_quote.**fee**string required  Fee paid for this quote (if executed) |
| result.taker\_quote.**is\_transfer**boolean required  Whether the order was generated through `private/transfer_position` |
| result.taker\_quote.**label**string required  User-defined label, if any |
| result.taker\_quote.**last\_update\_timestamp**integer required  Last update timestamp in ms since Unix epoch |
| result.taker\_quote.**legs\_hash**string required  Hash of the legs of the best quote to be signed by the taker. |
| result.taker\_quote.**liquidity\_role**string required  Liquidity role enum  `maker` `taker` |
| result.taker\_quote.**max\_fee**string required  Signed max fee |
| result.taker\_quote.**mmp**boolean required  Whether the quote is tagged for market maker protections (default false) |
| result.taker\_quote.**nonce**integer required  Nonce |
| result.taker\_quote.**quote\_id**string required  Quote ID |
| result.taker\_quote.**rfq\_id**string required  RFQ ID |
| result.taker\_quote.**signature**string required  Ethereum signature of the quote |
| result.taker\_quote.**signature\_expiry\_sec**integer required  Unix timestamp in seconds |
| result.taker\_quote.**signer**string required  Owner wallet address or registered session key that signed the quote |
| result.taker\_quote.**status**string required  Status enum  `open` `filled` `cancelled` `expired` |
| result.taker\_quote.**subaccount\_id**integer required  Subaccount ID |
| result.taker\_quote.**tx\_hash**string or null required  Blockchain transaction hash (only for executed quotes) |
| result.taker\_quote.**tx\_status**string or null required  Blockchain transaction status (only for executed quotes) enum  `requested` `pending` `settled` `reverted` `ignored` `timed_out` |
| result.taker\_quote.**legs**array of objects required  Quote legs |
| result.taker\_quote.legs[].**amount**string required  Amount in units of the base |
| result.taker\_quote.legs[].**direction**string required  Leg direction enum  `buy` `sell` |
| result.taker\_quote.legs[].**instrument\_name**string required  Instrument name |
| result.taker\_quote.legs[].**price**string required  Leg price |

### Example

ShellJavaScriptPython

```
{request_example_shell}

```

```
{request_example_javascript}

```

```
{request_example_python}

```

> The above command returns JSON structured like this:

JSON

```
{response_example_json}

```

---

## private-update_notifications

**Title:** Private Update_Notifications
**URL:** https://docs.derive.xyz/reference/private-update_notifications

### Method Name

#### `private/update_notifications`

RPC to mark specified notifications as seen for a given subaccount.  
Required minimum session key permission level is `account`

### Parameters

|  |
| --- |
| **subaccount\_id**integer required  Subaccount\_id |
| **notification\_ids**array of integers required  List of notification IDs to be marked as seen |
| **status**string  Status of the notification enum  `unseen` `seen` `hidden` |

### Response

|  |
| --- |
| **id**string or integer required |
| **result**object required |
| result.**updated\_count**integer required  Number of notifications marked as seen |

### Example

ShellJavaScriptPython

```
{request_example_shell}

```

```
{request_example_javascript}

```

```
{request_example_python}

```

> The above command returns JSON structured like this:

JSON

```
{response_example_json}

```

---

## private-withdraw

**Title:** Private Withdraw
**URL:** https://docs.derive.xyz/reference/private-withdraw

### Method Name

#### `private/withdraw`

Withdraw an asset to wallet.  
  
See `public/withdraw_debug` for debugging invalid signature issues or go to guides in Documentation.  
Required minimum session key permission level is `admin`

### Parameters

|  |
| --- |
| **amount**string required  Amount of the asset to withdraw |
| **asset\_name**string required  Name of asset to withdraw |
| **nonce**integer required  Unique nonce defined as  (e.g. 1695836058725001, where 001 is the random number) |
| **signature**string required  Ethereum signature of the withdraw |
| **signature\_expiry\_sec**integer required  Unix timestamp in seconds. Expiry MUST be >5min from now |
| **signer**string required  Ethereum wallet address that is signing the withdraw |
| **subaccount\_id**integer required  Subaccount\_id |
| **is\_atomic\_signing**boolean  Used by vaults to determine whether the signature is an EIP-1271 signature |

### Response

|  |
| --- |
| **id**string or integer required |
| **result**object required |
| result.**status**string required  `requested` |
| result.**transaction\_id**string required  Transaction id of the withdrawal |

### Example

ShellJavaScriptPython

```
{request_example_shell}

```

```
{request_example_javascript}

```

```
{request_example_python}

```

> The above command returns JSON structured like this:

JSON

```
{response_example_json}

```

---

# Websocket


## auctions-watch

**Title:** Auctions Watch
**URL:** https://docs.derive.xyz/reference/auctions-watch

### Channel Name Schema

#### `auctions.watch`

Subscribe to state of ongoing auctions.

### Channel Parameters

|  |
| --- |

### Notification Data

|  |
| --- |
| **channel**string required  Subscribed channel name |
| **data**array of objects required |
| data[].**details**object or null required  Details of the auction, if ongoing. |
| data[].details.**currency**string or null required  Currency of subaccount (for PM margin type). |
| data[].details.**estimated\_bid\_price**string required  Estimated bid price for this liquidation. This value is not scaled by the bid percent, and instead represents the discounted mark value of the whole subaccount. This value will be negative for insolvent auctions. |
| data[].details.**estimated\_discount\_pnl**string required  Estimated profit relative to `estimated_mtm` if the liquidation is successful, assuming execution at `estimated_percent_bid` and `estimated_bid_price`. |
| data[].details.**estimated\_mtm**string required  Estimated mark-to-market value of the subaccount being auctioned off. This value is not scaled by the bid percent, and instead represents the un-discounted mark value of the whole subaccount. |
| data[].details.**estimated\_percent\_bid**string required  An estimate for the maximum percent of the subaccount that can be liquidated. |
| data[].details.**last\_seen\_trade\_id**integer required  Last trade ID for the account being auctioned off (to use in the `private/liquidate` endpoint). This value is used to ensure that the state of balances reported in `subaccount_balances` is in sync. A trade ID error from `private/liquidate` indicates that the channel is currently out of sync with the on-chain state of the subaccount due to a pending bid. |
| data[].details.**margin\_type**string required  Margin type of the subaccount being auctioned off. It is recommended to bid on subaccounts using the same margin type and currency as to not run into unsupported currency errors or maximum account size limits. enum  `PM` `SM` `PM2` |
| data[].details.**min\_cash\_transfer**string required  Suggested minimum amount of cash to transfer to a newly created subaccount for bidding (to use in the `private/liquidate` endpoint). Any unused cash will get returned back to the original subaccount. If the bidder plans to bid less than the `estimated_percent_bid`, they may scale this value down accordingly. |
| data[].details.**min\_price\_limit**string required  Estimated minimum `price_limit` (to use in the `private/liquidate` endpoint). This is the minimum amount of cash that would be required to buy out the percent of the subaccount. If the bidder plans to bid less than the `estimated_percent_bid`, they may scale this value down accordingly. |
| data[].details.**subaccount\_balances**object required  Current balances of the subaccount being auctioned off. The bidder should expect to receive a percentage of these balances proportional to the `estimated_percent_bid`, and pay `estimated_bid_price * estimated_percent_bid` for them. These balances already include any pending perp settlements and funding payments into the USDC balance. |
|  |
| data[].**state**string required  State of the auction. enum  `ongoing` `ended` |
| data[].**subaccount\_id**integer required  Subaccount ID being auctioned off. |
| data[].**timestamp**integer required  Timestamp of the auction result (in milliseconds since epoch). |

### Example

> *Subscriptions are only available via websockets.*

JavaScriptPython

```
{request_example_javascript}

```

```
{request_example_python}

```

> Notification messages on this channel will look like this:

JSON

```
{response_example_json}

```

---

## margin-watch

**Title:** Margin Watch
**URL:** https://docs.derive.xyz/reference/margin-watch

### Channel Name Schema

#### `margin.watch`

Subscribe to state of margin and MtM of all users.

### Channel Parameters

|  |
| --- |

### Notification Data

|  |
| --- |
| **channel**string required  Subscribed channel name |
| **data**array of objects required |
| data[].**currency**string required  Currency of subaccount |
| data[].**initial\_margin**string required  Total initial margin requirement of all positions and collaterals. |
| data[].**maintenance\_margin**string required  Total maintenance margin requirement of all positions and collaterals.If this value falls below zero, the subaccount will be flagged for liquidation. |
| data[].**margin\_type**string required  Margin type of subaccount (`PM` (Portfolio Margin) or `SM` (Standard Margin)) enum  `PM` `SM` `PM2` |
| data[].**subaccount\_id**integer required  Subaccount\_id |
| data[].**subaccount\_value**string required  Total mark-to-market value of all positions and collaterals |
| data[].**valuation\_timestamp**integer required  Timestamp (in seconds since epoch) of when margin and MtM were computed. |
| data[].**collaterals**array of objects required  All collaterals that count towards margin of subaccount |
| data[].collaterals[].**amount**string required  Asset amount of given collateral |
| data[].collaterals[].**asset\_name**string required  Asset name |
| data[].collaterals[].**asset\_type**string required  Type of asset collateral (currently always `erc20`) enum  `erc20` `option` `perp` |
| data[].collaterals[].**initial\_margin**string required  USD value of collateral that contributes to initial margin |
| data[].collaterals[].**maintenance\_margin**string required  USD value of collateral that contributes to maintenance margin |
| data[].collaterals[].**mark\_price**string required  Current mark price of the asset |
| data[].collaterals[].**mark\_value**string required  USD value of the collateral (amount \* mark price) |
|  |
| data[].**positions**array of objects required  All active positions of subaccount |
| data[].positions[].**amount**string required  Position amount held by subaccount |
| data[].positions[].**delta**string required  Asset delta (w.r.t. forward price for options, `1.0` for perps) |
| data[].positions[].**gamma**string required  Asset gamma (zero for non-options) |
| data[].positions[].**index\_price**string required  Current index (oracle) price for position's currency |
| data[].positions[].**initial\_margin**string required  USD initial margin requirement for this position |
| data[].positions[].**instrument\_name**string required  Instrument name (same as the base Asset name) |
| data[].positions[].**instrument\_type**string required  `erc20`, `option`, or `perp` enum  `erc20` `option` `perp` |
| data[].positions[].**liquidation\_price**string or null required  Index price at which position will be liquidated |
| data[].positions[].**maintenance\_margin**string required  USD maintenance margin requirement for this position |
| data[].positions[].**mark\_price**string required  Current mark price for position's instrument |
| data[].positions[].**mark\_value**string required  USD value of the position; this represents how much USD can be recieved by fully closing the position at the current oracle price |
| data[].positions[].**theta**string required  Asset theta (zero for non-options) |
| data[].positions[].**vega**string required  Asset vega (zero for non-options) |

### Example

> *Subscriptions are only available via websockets.*

JavaScriptPython

```
{request_example_javascript}

```

```
{request_example_python}

```

> Notification messages on this channel will look like this:

JSON

```
{response_example_json}

```

---

## orderbook-instrument_name-group-depth

**Title:** Orderbook Instrument_Name Group Depth
**URL:** https://docs.derive.xyz/reference/orderbook-instrument_name-group-depth

### Channel Name Schema

#### `orderbook.{instrument_name}.{group}.{depth}`

Periodically publishes bids and asks for an instrument.

### Channel Parameters

|  |
| --- |
| **depth**string required  Number of price levels returned enum  `1` `10` `20` `100` |
| **group**string required  Price grouping (rounding) enum  `1` `10` `100` |
| **instrument\_name**string required  Instrument name |

### Notification Data

|  |
| --- |
| **channel**string required  Subscribed channel name |
| **data**object required |
| data.**instrument\_name**string required  Instrument name |
| data.**publish\_id**integer required  Publish ID, incremented for each publish |
| data.**timestamp**integer required  Timestamp of the orderbook snapshot |
| data.**asks**array of arrays required  List of asks as [price, amount] tuples optionally grouped into price buckets |
| data.**bids**array of arrays required  List of bids as [price, amount] tuples optionally grouped into price buckets |

### Example

> *Subscriptions are only available via websockets.*

JavaScriptPython

```
{request_example_javascript}

```

```
{request_example_python}

```

> Notification messages on this channel will look like this:

JSON

```
{response_example_json}

```

---

## spot_feed-currency

**Title:** Spot_Feed Currency
**URL:** https://docs.derive.xyz/reference/spot_feed-currency

### Channel Name Schema

#### `spot_feed.{currency}`

Periodically publishes spot index price by currency.

### Channel Parameters

|  |
| --- |
| **currency**string required  Currency |

### Notification Data

|  |
| --- |
| **channel**string required  Subscribed channel name |
| **data**object required |
| data.**timestamp**integer required  Timestamp of the spot feed snapshot |
| data.**feeds**object required  Spot feed data |

### Example

> *Subscriptions are only available via websockets.*

JavaScriptPython

```
{request_example_javascript}

```

```
{request_example_python}

```

> Notification messages on this channel will look like this:

JSON

```
{response_example_json}

```

---

## subaccount_id-balances

**Title:** Subaccount_Id Balances
**URL:** https://docs.derive.xyz/reference/subaccount_id-balances

### Channel Name Schema

#### `{subaccount_id}.balances`

Subscribe to changes in user's positions for a given subaccount ID.  
  
For perpetuals, additional balance updates are emitted under the name Q-{ccy}-PERP where Q stands for "quote".  
This balance is a proxy for an on-chain state of lastMarkPrice.  
Because of a synchronization lag with the on-chain state, the orderbook instead keeps track of a running total cost of perpetual trades,  
  
For example:  
Q-ETH-PERP balance of $6,600 and an ETH-PERP balance of 3 means the lastMarkPrice state is estimated to be $2,200.

### Channel Parameters

|  |
| --- |
| **subaccount\_id**integer required  Subaccount ID |

### Notification Data

|  |
| --- |
| **channel**string required  Subscribed channel name |
| **data**array of objects required |
| data[].**name**string required  Name of colletaral asset or instrument |
| data[].**new\_balance**string required  Balance after update |
| data[].**previous\_balance**string required  Balance before update |
| data[].**update\_type**string required  Type of transaction enum  `trade` `asset_deposit` `asset_withdrawal` `transfer` `subaccount_deposit` `subaccount_withdrawal` `liquidation` `liquidator` `onchain_drift_fix` `perp_settlement` `option_settlement` `interest_accrual` `onchain_revert` `double_revert` |

### Example

> *Subscriptions are only available via websockets.*

JavaScriptPython

```
{request_example_javascript}

```

```
{request_example_python}

```

> Notification messages on this channel will look like this:

JSON

```
{response_example_json}

```

---

## subaccount_id-orders

**Title:** Subaccount_Id Orders
**URL:** https://docs.derive.xyz/reference/subaccount_id-orders

### Channel Name Schema

#### `{subaccount_id}.orders`

Subscribe to changes in user's orders for a given subaccount ID.

### Channel Parameters

|  |
| --- |
| **subaccount\_id**integer required  Subaccount ID |

### Notification Data

|  |
| --- |
| **channel**string required  Subscribed channel name |
| **data**array of objects required |
| data[].**amount**string required  Order amount in units of the base |
| data[].**average\_price**string required  Average fill price |
| data[].**cancel\_reason**string required  If cancelled, reason behind order cancellation enum  `user_request` `mmp_trigger` `insufficient_margin` `signed_max_fee_too_low` `cancel_on_disconnect` `ioc_or_market_partial_fill` `session_key_deregistered` `subaccount_withdrawn` `compliance` `trigger_failed` `validation_failed` |
| data[].**creation\_timestamp**integer required  Creation timestamp (in ms since Unix epoch) |
| data[].**direction**string required  Order direction enum  `buy` `sell` |
| data[].**filled\_amount**string required  Total filled amount for the order |
| data[].**instrument\_name**string required  Instrument name |
| data[].**is\_transfer**boolean required  Whether the order was generated through `private/transfer_position` |
| data[].**label**string required  Optional user-defined label for the order |
| data[].**last\_update\_timestamp**integer required  Last update timestamp (in ms since Unix epoch) |
| data[].**limit\_price**string required  Limit price in quote currency |
| data[].**max\_fee**string required  Max fee in units of the quote currency |
| data[].**mmp**boolean required  Whether the order is tagged for market maker protections |
| data[].**nonce**integer required  Unique nonce defined as  (e.g. 1695836058725001, where 001 is the random number) |
| data[].**order\_fee**string required  Total order fee paid so far |
| data[].**order\_id**string required  Order ID |
| data[].**order\_status**string required  Order status enum  `open` `filled` `cancelled` `expired` `untriggered` |
| data[].**order\_type**string required  Order type enum  `limit` `market` |
| data[].**quote\_id**string or null required  Quote ID if the trade was executed via RFQ |
| data[].**signature**string required  Ethereum signature of the order |
| data[].**signature\_expiry\_sec**integer required  Signature expiry timestamp |
| data[].**signer**string required  Owner wallet address or registered session key that signed order |
| data[].**subaccount\_id**integer required  Subaccount ID |
| data[].**time\_in\_force**string required  Time in force enum  `gtc` `post_only` `fok` `ioc` |
| data[].**replaced\_order\_id**string or null  If replaced, ID of the order that was replaced |
| data[].**trigger\_price**string or null  (Required for trigger orders) Index or Market price to trigger order at |
| data[].**trigger\_price\_type**string or null  (Required for trigger orders) Trigger with Index or Mark Price enum  `mark` `index` |
| data[].**trigger\_reject\_message**string or null  (Required for trigger orders) Error message if error occured during trigger |
| data[].**trigger\_type**string or null  (Required for trigger orders) Stop-loss or Take-profit. enum  `stoploss` `takeprofit` |

### Example

> *Subscriptions are only available via websockets.*

JavaScriptPython

```
{request_example_javascript}

```

```
{request_example_python}

```

> Notification messages on this channel will look like this:

JSON

```
{response_example_json}

```

---

## subaccount_id-quotes

**Title:** Subaccount_Id Quotes
**URL:** https://docs.derive.xyz/reference/subaccount_id-quotes

### Channel Name Schema

#### `{subaccount_id}.quotes`

Subscribe to quote state for a given subaccount ID.  
This will notify the usser about the state change of the quotes they have sent.

### Channel Parameters

|  |
| --- |
| **subaccount\_id**integer required  Subaccount ID to get quote state updates for |

### Notification Data

|  |
| --- |
| **channel**string required  Subscribed channel name |
| **data**array of objects required |
| data[].**cancel\_reason**string required  Cancel reason, if any enum  `user_request` `insufficient_margin` `signed_max_fee_too_low` `mmp_trigger` `cancel_on_disconnect` `session_key_deregistered` `subaccount_withdrawn` `rfq_no_longer_open` `compliance` |
| data[].**creation\_timestamp**integer required  Creation timestamp in ms since Unix epoch |
| data[].**direction**string required  Quote direction enum  `buy` `sell` |
| data[].**fee**string required  Fee paid for this quote (if executed) |
| data[].**is\_transfer**boolean required  Whether the order was generated through `private/transfer_position` |
| data[].**label**string required  User-defined label, if any |
| data[].**last\_update\_timestamp**integer required  Last update timestamp in ms since Unix epoch |
| data[].**legs\_hash**string required  Hash of the legs of the best quote to be signed by the taker. |
| data[].**liquidity\_role**string required  Liquidity role enum  `maker` `taker` |
| data[].**max\_fee**string required  Signed max fee |
| data[].**mmp**boolean required  Whether the quote is tagged for market maker protections (default false) |
| data[].**nonce**integer required  Nonce |
| data[].**quote\_id**string required  Quote ID |
| data[].**rfq\_id**string required  RFQ ID |
| data[].**signature**string required  Ethereum signature of the quote |
| data[].**signature\_expiry\_sec**integer required  Unix timestamp in seconds |
| data[].**signer**string required  Owner wallet address or registered session key that signed the quote |
| data[].**status**string required  Status enum  `open` `filled` `cancelled` `expired` |
| data[].**subaccount\_id**integer required  Subaccount ID |
| data[].**tx\_hash**string or null required  Blockchain transaction hash (only for executed quotes) |
| data[].**tx\_status**string or null required  Blockchain transaction status (only for executed quotes) enum  `requested` `pending` `settled` `reverted` `ignored` `timed_out` |
| data[].**legs**array of objects required  Quote legs |
| data[].legs[].**amount**string required  Amount in units of the base |
| data[].legs[].**direction**string required  Leg direction enum  `buy` `sell` |
| data[].legs[].**instrument\_name**string required  Instrument name |
| data[].legs[].**price**string required  Leg price |

### Example

> *Subscriptions are only available via websockets.*

JavaScriptPython

```
{request_example_javascript}

```

```
{request_example_python}

```

> Notification messages on this channel will look like this:

JSON

```
{response_example_json}

```

---

## subaccount_id-trades

**Title:** Subaccount_Id Trades
**URL:** https://docs.derive.xyz/reference/subaccount_id-trades

### Channel Name Schema

#### `{subaccount_id}.trades`

Subscribe to user's trades (order executions) for a given subaccount ID.

### Channel Parameters

|  |
| --- |
| **subaccount\_id**integer required  Subaccount ID |

### Notification Data

|  |
| --- |
| **channel**string required  Subscribed channel name |
| **data**array of objects required |
| data[].**direction**string required  Order direction enum  `buy` `sell` |
| data[].**expected\_rebate**string required  Expected rebate for this trade |
| data[].**index\_price**string required  Index price of the underlying at the time of the trade |
| data[].**instrument\_name**string required  Instrument name |
| data[].**is\_transfer**boolean required  Whether the trade was generated through `private/transfer_position` |
| data[].**label**string required  Optional user-defined label for the order |
| data[].**liquidity\_role**string required  Role of the user in the trade enum  `maker` `taker` |
| data[].**mark\_price**string required  Mark price of the instrument at the time of the trade |
| data[].**order\_id**string required  Order ID |
| data[].**quote\_id**string or null required  Quote ID if the trade was executed via RFQ |
| data[].**realized\_pnl**string required  Realized PnL for this trade |
| data[].**realized\_pnl\_excl\_fees**string required  Realized PnL for this trade using cost accounting that excludes fees |
| data[].**subaccount\_id**integer required  Subaccount ID |
| data[].**timestamp**integer required  Trade timestamp (in ms since Unix epoch) |
| data[].**trade\_amount**string required  Amount filled in this trade |
| data[].**trade\_fee**string required  Fee for this trade |
| data[].**trade\_id**string required  Trade ID |
| data[].**trade\_price**string required  Price at which the trade was filled |
| data[].**transaction\_id**string required  The transaction id of the related settlement transaction |
| data[].**tx\_hash**string or null required  Blockchain transaction hash |
| data[].**tx\_status**string required  Blockchain transaction status enum  `requested` `pending` `settled` `reverted` `ignored` `timed_out` |

### Example

> *Subscriptions are only available via websockets.*

JavaScriptPython

```
{request_example_javascript}

```

```
{request_example_python}

```

> Notification messages on this channel will look like this:

JSON

```
{response_example_json}

```

---

## subaccount_id-trades-tx_status

**Title:** Subaccount_Id Trades Tx_Status
**URL:** https://docs.derive.xyz/reference/subaccount_id-trades-tx_status

### Channel Name Schema

#### `{subaccount_id}.trades.{tx_status}`

Subscribe to user's trade settlement for a given subaccount ID.

### Channel Parameters

|  |
| --- |
| **subaccount\_id**integer required  Subaccount ID |
| **tx\_status**string required  Transaction status (`settled` or `reverted`) enum  `settled` `reverted` `timed_out` |

### Notification Data

|  |
| --- |
| **channel**string required  Subscribed channel name |
| **data**array of objects required |
| data[].**direction**string required  Order direction enum  `buy` `sell` |
| data[].**expected\_rebate**string required  Expected rebate for this trade |
| data[].**index\_price**string required  Index price of the underlying at the time of the trade |
| data[].**instrument\_name**string required  Instrument name |
| data[].**is\_transfer**boolean required  Whether the trade was generated through `private/transfer_position` |
| data[].**label**string required  Optional user-defined label for the order |
| data[].**liquidity\_role**string required  Role of the user in the trade enum  `maker` `taker` |
| data[].**mark\_price**string required  Mark price of the instrument at the time of the trade |
| data[].**order\_id**string required  Order ID |
| data[].**quote\_id**string or null required  Quote ID if the trade was executed via RFQ |
| data[].**realized\_pnl**string required  Realized PnL for this trade |
| data[].**realized\_pnl\_excl\_fees**string required  Realized PnL for this trade using cost accounting that excludes fees |
| data[].**subaccount\_id**integer required  Subaccount ID |
| data[].**timestamp**integer required  Trade timestamp (in ms since Unix epoch) |
| data[].**trade\_amount**string required  Amount filled in this trade |
| data[].**trade\_fee**string required  Fee for this trade |
| data[].**trade\_id**string required  Trade ID |
| data[].**trade\_price**string required  Price at which the trade was filled |
| data[].**transaction\_id**string required  The transaction id of the related settlement transaction |
| data[].**tx\_hash**string or null required  Blockchain transaction hash |
| data[].**tx\_status**string required  Blockchain transaction status enum  `requested` `pending` `settled` `reverted` `ignored` `timed_out` |

### Example

> *Subscriptions are only available via websockets.*

JavaScriptPython

```
{request_example_javascript}

```

```
{request_example_python}

```

> Notification messages on this channel will look like this:

JSON

```
{response_example_json}

```

---

## subscribe

**Title:** Subscribe
**URL:** https://docs.derive.xyz/reference/subscribe

### Method Name

#### `subscribe`

Subscribe to a list of channels.

### Parameters

|  |
| --- |
| **channels**array of strings required  A list of channels names to subscribe to |

### Response

|  |
| --- |
| **id**string or integer required |
| **result**object required |
| result.**current\_subscriptions**array of strings required  A list of channels subscribed to after the subscribe operation. |
| result.**status**object required  A mapping of `channel` ⭢ `status`. Successful subscriptions will have status `ok`. Failed subscriptions will contain an error message. |

### Example

ShellJavaScriptPython

```
{request_example_shell}

```

```
{request_example_javascript}

```

```
{request_example_python}

```

> The above command returns JSON structured like this:

JSON

```
{response_example_json}

```

---

## ticker-instrument_name-interval

**Title:** Ticker Instrument_Name Interval
**URL:** https://docs.derive.xyz/reference/ticker-instrument_name-interval

### Channel Name Schema

#### `ticker.{instrument_name}.{interval}`

Periodically publishes ticker info (best bid / ask, instrument contraints, fees, etc.) for a single instrument.

### Channel Parameters

|  |
| --- |
| **instrument\_name**string required  Instrument name |
| **interval**string required  Interval in milliseconds enum  `100` `1000` |

### Notification Data

|  |
| --- |
| **channel**string required  Subscribed channel name |
| **data**object required |
| data.**timestamp**integer required  Timestamp of the ticker feed snapshot |
| data.**instrument\_ticker**object required  Instrument of the ticker feed snapshot |
| data.instrument\_ticker.**amount\_step**string required  Minimum valid increment of order amount |
| data.instrument\_ticker.**base\_asset\_address**string required  Blockchain address of the base asset |
| data.instrument\_ticker.**base\_asset\_sub\_id**string required  Sub ID of the specific base asset as defined in Asset.sol |
| data.instrument\_ticker.**base\_currency**string required  Underlying currency of base asset (`ETH`, `BTC`, etc) |
| data.instrument\_ticker.**base\_fee**string required  $ base fee added to every taker order |
| data.instrument\_ticker.**best\_ask\_amount**string required  Amount of contracts / tokens available at best ask price |
| data.instrument\_ticker.**best\_ask\_price**string required  Best ask price |
| data.instrument\_ticker.**best\_bid\_amount**string required  Amount of contracts / tokens available at best bid price |
| data.instrument\_ticker.**best\_bid\_price**string required  Best bid price |
| data.instrument\_ticker.**erc20\_details**object or null required  Details of the erc20 asset (if applicable) |
| data.instrument\_ticker.erc20\_details.**decimals**integer required  Number of decimals of the underlying on-chain ERC20 token |
| data.instrument\_ticker.erc20\_details.**borrow\_index**string  Latest borrow index as per `CashAsset.sol` implementation |
| data.instrument\_ticker.erc20\_details.**supply\_index**string  Latest supply index as per `CashAsset.sol` implementation |
| data.instrument\_ticker.erc20\_details.**underlying\_erc20\_address**string  Address of underlying on-chain ERC20 (not V2 asset) |
|  |
| data.instrument\_ticker.**fifo\_min\_allocation**string required  Minimum number of contracts that get filled using FIFO. Actual number of contracts that gets filled by FIFO will be the max between this value and (1 - pro\_rata\_fraction) x order\_amount, plus any size leftovers due to rounding. |
| data.instrument\_ticker.**five\_percent\_ask\_depth**string required  Total amount of contracts / tokens available at 5 percent above best ask price |
| data.instrument\_ticker.**five\_percent\_bid\_depth**string required  Total amount of contracts / tokens available at 5 percent below best bid price |
| data.instrument\_ticker.**index\_price**string required  Index price |
| data.instrument\_ticker.**instrument\_name**string required  Instrument name |
| data.instrument\_ticker.**instrument\_type**string required  `erc20`, `option`, or `perp` enum  `erc20` `option` `perp` |
| data.instrument\_ticker.**is\_active**boolean required  If `True`: instrument is tradeable within `activation` and `deactivation` timestamps |
| data.instrument\_ticker.**maker\_fee\_rate**string required  Percent of spot price fee rate for makers |
| data.instrument\_ticker.**mark\_price**string required  Mark price |
| data.instrument\_ticker.**max\_price**string required  Maximum price at which an agressive buyer can be matched. Any portion of a market order that would execute above this price will be cancelled. A limit buy order with limit price above this value is treated as post only (i.e. it will be rejected if it would cross any existing resting order). |
| data.instrument\_ticker.**maximum\_amount**string required  Maximum valid amount of contracts / tokens per trade |
| data.instrument\_ticker.**min\_price**string required  Minimum price at which an agressive seller can be matched. Any portion of a market order that would execute below this price will be cancelled. A limit sell order with limit price below this value is treated as post only (i.e. it will be rejected if it would cross any existing resting order). |
| data.instrument\_ticker.**minimum\_amount**string required  Minimum valid amount of contracts / tokens per trade |
| data.instrument\_ticker.**option\_details**object or null required  Details of the option asset (if applicable) |
| data.instrument\_ticker.option\_details.**expiry**integer required  Unix timestamp of expiry date (in seconds) |
| data.instrument\_ticker.option\_details.**index**string required  Underlying settlement price index |
| data.instrument\_ticker.option\_details.**option\_type**string required   enum  `C` `P` |
| data.instrument\_ticker.option\_details.**strike**string required |
| data.instrument\_ticker.option\_details.**settlement\_price**string or null  Settlement price of the option |
|  |
| data.instrument\_ticker.**option\_pricing**object or null required  Greeks, forward price, iv and mark price of the instrument (options only) |
| data.instrument\_ticker.option\_pricing.**ask\_iv**string required  Implied volatility of the current best ask |
| data.instrument\_ticker.option\_pricing.**bid\_iv**string required  Implied volatility of the current best bid |
| data.instrument\_ticker.option\_pricing.**delta**string required  Delta of the option |
| data.instrument\_ticker.option\_pricing.**discount\_factor**string required  Discount factor used to calculate option premium |
| data.instrument\_ticker.option\_pricing.**forward\_price**string required  Forward price used to calculate option premium |
| data.instrument\_ticker.option\_pricing.**gamma**string required  Gamma of the option |
| data.instrument\_ticker.option\_pricing.**iv**string required  Implied volatility of the option |
| data.instrument\_ticker.option\_pricing.**mark\_price**string required  Mark price of the option |
| data.instrument\_ticker.option\_pricing.**rho**string required  Rho of the option |
| data.instrument\_ticker.option\_pricing.**theta**string required  Theta of the option |
| data.instrument\_ticker.option\_pricing.**vega**string required  Vega of the option |
|  |
| data.instrument\_ticker.**perp\_details**object or null required  Details of the perp asset (if applicable) |
| data.instrument\_ticker.perp\_details.**aggregate\_funding**string required  Latest aggregated funding as per `PerpAsset.sol` |
| data.instrument\_ticker.perp\_details.**funding\_rate**string required  Current hourly funding rate as per `PerpAsset.sol` |
| data.instrument\_ticker.perp\_details.**index**string required  Underlying spot price index for funding rate |
| data.instrument\_ticker.perp\_details.**max\_rate\_per\_hour**string required  Max rate per hour as per `PerpAsset.sol` |
| data.instrument\_ticker.perp\_details.**min\_rate\_per\_hour**string required  Min rate per hour as per `PerpAsset.sol` |
| data.instrument\_ticker.perp\_details.**static\_interest\_rate**string required  Static interest rate as per `PerpAsset.sol` |
|  |
| data.instrument\_ticker.**pro\_rata\_amount\_step**string required  Pro-rata fill share of every order is rounded down to be a multiple of this number. Leftovers of the order due to rounding are filled FIFO. |
| data.instrument\_ticker.**pro\_rata\_fraction**string required  Fraction of order that gets filled using pro-rata matching. If zero, the matching is full FIFO. |
| data.instrument\_ticker.**quote\_currency**string required  Quote currency (`USD` for perps, `USDC` for options) |
| data.instrument\_ticker.**scheduled\_activation**integer required  Timestamp at which became or will become active (if applicable) |
| data.instrument\_ticker.**scheduled\_deactivation**integer required  Scheduled deactivation time for instrument (if applicable) |
| data.instrument\_ticker.**taker\_fee\_rate**string required  Percent of spot price fee rate for takers |
| data.instrument\_ticker.**tick\_size**string required  Tick size of the instrument, i.e. minimum price increment |
| data.instrument\_ticker.**timestamp**integer required  Timestamp of the ticker feed snapshot |
| data.instrument\_ticker.**open\_interest**object required  Margin type of subaccount (`PM` (Portfolio Margin) or `SM` (Standard Margin)) -> (current open interest, open interest cap, manager currency) |
| data.instrument\_ticker.**stats**object required  Aggregate trading stats for the last 24 hours |
| data.instrument\_ticker.stats.**contract\_volume**string required  Number of contracts traded during last 24 hours |
| data.instrument\_ticker.stats.**high**string required  Highest trade price during last 24h |
| data.instrument\_ticker.stats.**low**string required  Lowest trade price during last 24h |
| data.instrument\_ticker.stats.**num\_trades**string required  Number of trades during last 24h |
| data.instrument\_ticker.stats.**open\_interest**string required  Current total open interest |
| data.instrument\_ticker.stats.**percent\_change**string required  24-hour price change expressed as a percentage. Options: percent change in vol; Perps: percent change in mark price |
| data.instrument\_ticker.stats.**usd\_change**string required  24-hour price change in USD. |
|  |
| data.instrument\_ticker.**mark\_price\_fee\_rate\_cap**string or null  Percent of option price fee cap, e.g. 12.5%, null if not applicable |

### Example

> *Subscriptions are only available via websockets.*

JavaScriptPython

```
{request_example_javascript}

```

```
{request_example_python}

```

> Notification messages on this channel will look like this:

JSON

```
{response_example_json}

```

---

## trades-instrument_name

**Title:** Trades Instrument_Name
**URL:** https://docs.derive.xyz/reference/trades-instrument_name

### Channel Name Schema

#### `trades.{instrument_name}`

Subscribe to trades (order executions) for a given instrument name.

### Channel Parameters

|  |
| --- |
| **instrument\_name**string required  Instrument name |

### Notification Data

|  |
| --- |
| **channel**string required  Subscribed channel name |
| **data**array of objects required |
| data[].**direction**string required  Direction of the taker order enum  `buy` `sell` |
| data[].**index\_price**string required  Index price of the underlying at the time of the trade |
| data[].**instrument\_name**string required  Instrument name |
| data[].**mark\_price**string required  Mark price of the instrument at the time of the trade |
| data[].**quote\_id**string or null required  Quote ID if the trade was executed via RFQ |
| data[].**timestamp**integer required  Trade timestamp (in ms since Unix epoch) |
| data[].**trade\_amount**string required  Amount filled in this trade |
| data[].**trade\_id**string required  Trade ID |
| data[].**trade\_price**string required  Price at which the trade was filled |

### Example

> *Subscriptions are only available via websockets.*

JavaScriptPython

```
{request_example_javascript}

```

```
{request_example_python}

```

> Notification messages on this channel will look like this:

JSON

```
{response_example_json}

```

---

## trades-instrument_type-currency

**Title:** Trades Instrument_Type Currency
**URL:** https://docs.derive.xyz/reference/trades-instrument_type-currency

### Channel Name Schema

#### `trades.{instrument_type}.{currency}`

Subscribe to trades (order executions) for a given instrument type and currency.

### Channel Parameters

|  |
| --- |
| **currency**string required  Currency |
| **instrument\_type**string required  Instrument type enum  `erc20` `option` `perp` |

### Notification Data

|  |
| --- |
| **channel**string required  Subscribed channel name |
| **data**array of objects required |
| data[].**direction**string required  Direction of the taker order enum  `buy` `sell` |
| data[].**index\_price**string required  Index price of the underlying at the time of the trade |
| data[].**instrument\_name**string required  Instrument name |
| data[].**mark\_price**string required  Mark price of the instrument at the time of the trade |
| data[].**quote\_id**string or null required  Quote ID if the trade was executed via RFQ |
| data[].**timestamp**integer required  Trade timestamp (in ms since Unix epoch) |
| data[].**trade\_amount**string required  Amount filled in this trade |
| data[].**trade\_id**string required  Trade ID |
| data[].**trade\_price**string required  Price at which the trade was filled |

### Example

> *Subscriptions are only available via websockets.*

JavaScriptPython

```
{request_example_javascript}

```

```
{request_example_python}

```

> Notification messages on this channel will look like this:

JSON

```
{response_example_json}

```

---

## trades-instrument_type-currency-tx_status

**Title:** Trades Instrument_Type Currency Tx_Status
**URL:** https://docs.derive.xyz/reference/trades-instrument_type-currency-tx_status

### Channel Name Schema

#### `trades.{instrument_type}.{currency}.{tx_status}`

Subscribe to the status on on-chain trade settlement events for a given instrument type and currency.

### Channel Parameters

|  |
| --- |
| **currency**string required  Currency |
| **instrument\_type**string required  Instrument type enum  `erc20` `option` `perp` |
| **tx\_status**string required  Transaction status enum  `settled` `reverted` `timed_out` |

### Notification Data

|  |
| --- |
| **channel**string required  Subscribed channel name |
| **data**array of objects required |
| data[].**direction**string required  Order direction enum  `buy` `sell` |
| data[].**expected\_rebate**string required  Expected rebate for this trade |
| data[].**index\_price**string required  Index price of the underlying at the time of the trade |
| data[].**instrument\_name**string required  Instrument name |
| data[].**liquidity\_role**string required  Role of the user in the trade enum  `maker` `taker` |
| data[].**mark\_price**string required  Mark price of the instrument at the time of the trade |
| data[].**quote\_id**string or null required  Quote ID if the trade was executed via RFQ |
| data[].**realized\_pnl**string required  Realized PnL for this trade |
| data[].**realized\_pnl\_excl\_fees**string required  Realized PnL for this trade using cost accounting that excludes fees |
| data[].**subaccount\_id**integer required  Subaccount ID |
| data[].**timestamp**integer required  Trade timestamp (in ms since Unix epoch) |
| data[].**trade\_amount**string required  Amount filled in this trade |
| data[].**trade\_fee**string required  Fee for this trade |
| data[].**trade\_id**string required  Trade ID |
| data[].**trade\_price**string required  Price at which the trade was filled |
| data[].**tx\_hash**string required  Blockchain transaction hash |
| data[].**tx\_status**string required  Blockchain transaction status enum  `settled` `reverted` `timed_out` |
| data[].**wallet**string required  Wallet address (owner) of the subaccount |

### Example

> *Subscriptions are only available via websockets.*

JavaScriptPython

```
{request_example_javascript}

```

```
{request_example_python}

```

> Notification messages on this channel will look like this:

JSON

```
{response_example_json}

```

---

## unsubscribe

**Title:** Unsubscribe
**URL:** https://docs.derive.xyz/reference/unsubscribe

### Method Name

#### `unsubscribe`

Unsubscribe from a list of channels, or all channels.

### Parameters

|  |
| --- |
| **channels**array of strings  A list of channels names to unsubscribe from. If not provided, unsubscribe from all channels. |

### Response

|  |
| --- |
| **id**string or integer required |
| **result**object required |
| result.**remaining\_subscriptions**array of strings required  A list of channels still subscribed to after the unsubscribe operation. |
| result.**status**object required  A mapping of `channel` ⭢ `status`. Successfully unsubscribed channels will have status `ok`. |

### Example

ShellJavaScriptPython

```
{request_example_shell}

```

```
{request_example_javascript}

```

```
{request_example_python}

```

> The above command returns JSON structured like this:

JSON

```
{response_example_json}

```

---

## wallet-rfqs

**Title:** Wallet Rfqs
**URL:** https://docs.derive.xyz/reference/wallet-rfqs

### Channel Name Schema

#### `{wallet}.rfqs`

Subscribe to RFQs directed to a given wallet.

### Channel Parameters

|  |
| --- |
| **wallet**string required  Account (wallet) of RFQ market maker |

### Notification Data

|  |
| --- |
| **channel**string required  Subscribed channel name |
| **data**array of objects required |
| data[].**cancel\_reason**string required  Cancel reason, if any enum  `user_request` `insufficient_margin` `signed_max_fee_too_low` `mmp_trigger` `cancel_on_disconnect` `session_key_deregistered` `subaccount_withdrawn` `rfq_no_longer_open` `compliance` |
| data[].**creation\_timestamp**integer required  Creation timestamp in ms since Unix epoch |
| data[].**filled\_direction**string or null required  Direction at which the RFQ was filled (only if filled) enum  `buy` `sell` |
| data[].**last\_update\_timestamp**integer required  Last update timestamp in ms since Unix epoch |
| data[].**rfq\_id**string required  RFQ ID |
| data[].**status**string required  Status enum  `open` `filled` `cancelled` `expired` |
| data[].**subaccount\_id**integer required  Subaccount ID |
| data[].**total\_cost**string or null required  Total cost for the RFQ (only if filled) |
| data[].**valid\_until**integer required  RFQ expiry timestamp in ms since Unix epoch |
| data[].**legs**array of objects required  RFQ legs |
| data[].legs[].**amount**string required  Amount in units of the base |
| data[].legs[].**direction**string required  Leg direction enum  `buy` `sell` |
| data[].legs[].**instrument\_name**string required  Instrument name |

### Example

> *Subscriptions are only available via websockets.*

JavaScriptPython

```
{request_example_javascript}

```

```
{request_example_python}

```

> Notification messages on this channel will look like this:

JSON

```
{response_example_json}

```

---

# Other


## create-or-deposit-to-subaccount

**Title:** Solidity Objects
**URL:** https://docs.derive.xyz/reference/create-or-deposit-to-subaccount

**NOTE** If you haven't done so yet, use the `public/create_account` to create an account using an ETH EOA or Smart Contract wallet.

Creating subaccount and depositing cash is a "self-custodial" request, a request that is guaranteed to not be alter-able by anyone except you, which means that some extra signing is required.

Creating/Depositing is almost the same process, and you can create a new account and deposit to it within as a single request.

To create or deposit to a Sub Account:

1. Approve USDC for spending on the deposit contract
2. Create a DepositModuleData object and encode it as bytes
3. Create, hash and sign a SignedAction object using the encoded DepositData
4. Send the request to either `/private/deposit` or `/private/create_subaccount`

For debugging, we have provided routes that return all intermediate values using the `public/create_subaccount_debug` and `public/deposit_debug` routes.

---

> ## 👍 If you are struggling to encode data correctly, you can use the public/create\_subaccount\_debug endpoint. The route takes in all raw inputs and returns intermediary outputs shown in the below steps.

## Code Example:

### 1. Approve USDC for spend on the deposit contract

This will allow the Deposit Module in`Matching.sol` to debit the USDC in your wallet.

TypeScript

```
async function approveUSDCForDeposit(wallet: ethers.Wallet, amount: string) {
    const USDCcontract = new ethers.Contract(
      USDC_ADDRESS,
      ["function approve(address _spender, uint256 _value) public returns (bool success)"],
      wallet
    );
    const nonce = await wallet.provider?.getTransactionCount(wallet.address, "pending");
    await USDCcontract.approve(DEPOSIT_MODULE_ADDRESS, ethers.parseUnits(amount, 6), {  
        gasLimit: 1000000,
        nonce: nonce
    });
}

```

**NOTE**: to create an empty SubAccount you can skip this step and use 0 as the deposit amount in the next step

### 2. Encode `DepositModuleData`

In order to ensure that only you can deposit cash, the below creates a `DepositModuleData` object and encodes it as `Bytes` to be used in the `SignedAction` object. `depositAmount` can be zero if you are creating an account.

TypeScript

```
function encodeDepositData(amount: string): Buffer {
    let encoded_data = encoder.encode( // same as "encoded_data" in public/create_subaccount_debug
      ['uint256', 'address', 'address'],
      [
        ethers.parseUnits(amount, 6),
        CASH_ADDRESS,
        STANDARD_RISK_MANAGER_ADDRESS
      ]
    );
    return ethers.keccak256(Buffer.from(encoded_data.slice(2), 'hex')) // same as "encoded_data_hashed" in public/create_subaccount_debug
}

```

### 3. Encode and Sign `SignedAction`

Finally, a `SignedAction` object is created and signed using your private key (or session\_key) which will include the `DepositModuleData` and additional needed info.

> ## 👍 The `subaccoutId` should be 0 if you are creating a new account.

TypeScript

```
function generateSignature(subaccountId: number, encoded_data_hashed: Buffer, expiry: number, nonce: number): string {
    const action_hash = ethers.keccak256(  // same as "action_hash" in public/create_subaccount_debug
        encoder.encode(
            ['bytes32', 'uint256', 'uint256', 'address', 'bytes32', 'uint256', 'address', 'address'],
            [
                ACTION_TYPEHASH,
                subaccountId,
                nonce,
                DEPOSIT_MODULE_ADDRESS,
                encoded_data_hashed,
                expiry, // must be >5 min from now
                wallet.address,
                wallet.address
            ]
        )
    );

    const typed_data_hash = ethers.keccak256( // same as "typed_data_hash" in public/create_subaccount_debug
        Buffer.concat([
            Buffer.from("1901", "hex"),
            Buffer.from(DOMAIN_SEPARATOR.slice(2), "hex"),
            Buffer.from(action_hash.slice(2), "hex"),
        ])
    );

    return wallet.signingKey.sign(typed_data_hash).serialized  
}

```

### Fully working example with raw data

TypeScript

```
import { ethers, Contract } from "ethers";
import axios from 'axios';
import dotenv from 'dotenv';
import { getUTCEpochSec } from "../utils/timer";

dotenv.config();

// Environment variables, double check these in the docs constants section
const PRIVATE_KEY = process.env.OWNER_PRIVATE_KEY as string;
const PROVIDER_URL = 'https://l2-prod-testnet-0eakp60405.t.conduit.xyz'
const USDC_ADDRESS = '0xe80F2a02398BBf1ab2C9cc52caD1978159c215BD'
const DEPOSIT_MODULE_ADDRESS = '0x43223Db33AdA0575D2E100829543f8B04A37a1ec'
const CASH_ADDRESS = '0x6caf294DaC985ff653d5aE75b4FF8E0A66025928'
const ACTION_TYPEHASH = '0x4d7a9f27c403ff9c0f19bce61d76d82f9aa29f8d6d4b0c5474607d9770d1af17'
const STANDARD_RISK_MANAGER_ADDRESS = '0x28bE681F7bEa6f465cbcA1D25A2125fe7533391C' // Use the "PortfolioManager" address if using PM
const DOMAIN_SEPARATOR = '0x9bcf4dc06df5d8bf23af818d5716491b995020f377d3b7b64c29ed14e3dd1105'

// Ethers setup
const PROVIDER = new ethers.JsonRpcProvider(PROVIDER_URL);
const wallet = new ethers.Wallet(PRIVATE_KEY, PROVIDER);
const encoder = ethers.AbiCoder.defaultAbiCoder();  

const depositAmount = "10000";
const subaccountId = 0; // 0 For a new account

async function approveUSDCForDeposit(wallet: ethers.Wallet, amount: string) {
    const USDCcontract = new ethers.Contract(
      USDC_ADDRESS,
      ["function approve(address _spender, uint256 _value) public returns (bool success)"],
      wallet
    );
    const nonce = await wallet.provider?.getTransactionCount(wallet.address, "pending");
    await USDCcontract.approve(DEPOSIT_MODULE_ADDRESS, ethers.parseUnits(amount, 6), {  
        gasLimit: 1000000,
        nonce: nonce
    });
}

function encodeDepositData(amount: string): Buffer {
    let encoded_data = encoder.encode( // same as "encoded_data" in public/create_subaccount_debug
      ['uint256', 'address', 'address'],
      [
        ethers.parseUnits(amount, 6),
        CASH_ADDRESS,
        STANDARD_RISK_MANAGER_ADDRESS
      ]
    );
    return ethers.keccak256(Buffer.from(encoded_data.slice(2), 'hex')) // same as "encoded_data_hashed" in public/create_subaccount_debug
}

function generateSignature(subaccountId: number, encodedData: Buffer, expiry: number, nonce: number): string {
    const action_hash = ethers.keccak256(  // same as "action_hash" in public/create_subaccount_debug
        encoder.encode(
            ['bytes32', 'uint256', 'uint256', 'address', 'bytes32', 'uint256', 'address', 'address'],
            [
                ACTION_TYPEHASH,
                subaccountId,
                nonce,
                DEPOSIT_MODULE_ADDRESS,
                encodedData,
                expiry,
                wallet.address,
                wallet.address
            ]
        )
    );

    const typed_data_hash = ethers.keccak256( // same as "typed_data_hash" in public/create_subaccount_debug
        Buffer.concat([
            Buffer.from("1901", "hex"),
            Buffer.from(DOMAIN_SEPARATOR.slice(2), "hex"),
            Buffer.from(action_hash.slice(2), "hex"),
        ])
    );

    return wallet.signingKey.sign(typed_data_hash).serialized 
}

async function signAuthenticationHeader() {
    const timestamp = Date.now().toString();
    const signature = await wallet.signMessage(timestamp);
    return {
        "X-LyraWallet": wallet.address,
        "X-LyraTimestamp": timestamp,
        "X-LyraSignature": signature
    };
}

async function createSubaccount() {
    // An action nonce is used to prevent replay attacks
		// LYRA nonce format: ${CURRENT UTC MS +/- 1 day}${RANDOM 3 DIGIT NUMBER}
    const nonce = Number(`${Date.now()}${Math.round(Math.random() * 999)}`);
    const expiry = getUTCEpochSec() + 600; // must be >5 min from now

    const encoded_data_hashed = encodeDepositData(depositAmount); // same as "encoded_data_hashed" in public/create_subaccount_debug
    const depositSignature = generateSignature(subaccountId, encoded_data_hashed, expiry, nonce);
    const authHeader = await signAuthenticationHeader();

    await approveUSDCForDeposit(wallet, depositAmount);

    try {
        const response = await axios.request({
            method: "POST",
            url: "https://api-demo.lyra.finance/private/create_subaccount",
            data: {
                margin_type: "SM",
                wallet: wallet.address,
                signer: wallet.address,
                nonce: nonce,
                amount: depositAmount,
                signature: depositSignature,
                signature_expiry_sec: expiry,
                asset_name: 'USDC',
            },
            headers: authHeader,
        });
    
        console.log(JSON.stringify(response.data, null, '\t'));
    } catch (error) {
        console.error("Error depositing to subaccount:", error);
    }
}

createSubaccount();

```

---

## [Optional] Get subaccount id

If you created a new subaccount, you can retrieve the new subaccount id once the transaction has settled

TypeScript

```
async get_subaccounts(){
  let timestamp = Date.now() // ensure UTC
  let signature = await wallet.signMessage(timestamp).toString()

  const response = await axios.request<R>({
    "POST",
    "https://api-demo.lyra.finance/private/get_subaccounts",
    {wallet: wallet.address},
    {
      "X-LyraWallet": wallet.address,
      "X-LyraTimestamp": timestamp,
      "X-LyraSignature": signature
    }
  });
}

```

---

# Solidity Objects

### SignedAction Schema

| Param | Type | Description |
| --- | --- | --- |
| `subaccount_id` | `uint` | User subaccount id for the action (0 for a new subaccounts when depositing) |
| `nonce` | `uint` | Unique nonce defined as <UTC\_timestamp in ms><random\_number\_up\_to\_6\_digits> (e.g. 1695836058725001, where 001 is the random number) |
| `module` | `address` | Deposit module address (see [Protocol Constants](/reference/protocol-constants)) |
| `data` | `bytes` | Encoded module data ("DepositData" for deposits/creates) |
| `expiry` | `uint` | Signature expiry timestamp in sec |
| `owner` | `address` | Wallet address of the account owner |
| `signer` | `address` | Either owner wallet or session key |

### DepositModuleData Schema

| Param | Type | Description |
| --- | --- | --- |
| `amount` | `uint` | Amount to deposit (Can be 0 to create an empty account) |
| `asset` | `address` | Address of the asset being deposited. See [Protocol Constants](/reference/protocol-constants) |
| `managerForNewAccount` | `address` | Use the "StandardManager.sol" address if using SM or "PortfolioManager.sol" address if using PM. See [Protocol Constants](/reference/protocol-constants) |

---

## deposit-to-lyra-chain

**Title:** Testnet
**URL:** https://docs.derive.xyz/reference/deposit-to-lyra-chain

The easiest way to deposit to Derive Chain is by setting up an account via the Interface (see the guides under "Onboard via Interface) as this fully handles bridging / deposits / account setups / withdrawals. **Only use this guide if you'd like to setup your account fully on-chain.**

# Testnet

When onboarding fully on-chain in testnet you will need to reach out via our Discord `v2-support` channel . Note, onboarding via Interface will not require this step as you can mint USDC directly to testnet accounts via the interface.

# Mainnet

Derive Chain does not have a generic bridging interface yet. To deposit ETH and USDC to Derive Chain, you'll need to use etherscan.

> ## ❗️ **Do not deposit to Derive Chain unless you know what you are doing.** Derive uses multiple different bridges for different assets. If you are not sure, ask via the Discord `v2-support` channel. Funds can be lost if done incorrectly.

> ## 🚧 The [exchange frontend](https://www.lyra.finance/) uses smart contract wallets, and isn't suitable for a programmatic account setup that uses a regular EOA and private key.

## Depositing collateral assets (USDC, WETH, WBTC, etc.)

**To deposit collateral assets to Derive Chain, use the socket bridges.** Socket bridges enable fast withdrawals from the Derive chain within certain daily limits. Depending on the asset, deposits and withdrawals are available between Ethereum L1, Optimism and Arbitrum.

Derive uses a custom bridge that uses [Socket](https://socket.tech/) smart contracts and L1-L2 messaging infrastructure. This enables a fast bridge has the option for fast withdrawals that aren't subject to the 7 day challenge period.

For a full list of tokens, bridges and connectors, [see the socket github.](https://github.com/SocketDotTech/socket-plugs/blob/main/deployments/superbridge/prod_lyra_addresses.json) This has been tabulated below:

| Currency | Source chain | Token Address | Bridge address |
| --- | --- | --- | --- |
| USDC | Eth Mainnet | 0xA0b86991c6218b36c1d19D4a2e9Eb0cE3606eB48 | 0x6D303CEE7959f814042D31E0624fB88Ec6fbcC1d |
| USDC.e | Optimism | 0x7F5c764cBc14f9669B88837ca1490cCa17c31607 | 0xBb9CF28Bc1B41c5c7c76Ee1B2722C33eBB8fbD8C |
| USDC.e | Arbitrum | 0xFF970A61A04b1cA14834A43f5dE4533eBDDB5CC8 | 0xFB7B06538d837e4212D72E2A38e6c074F9076E0B |
| USDC | Optimism | 0x0b2C639c533813f4Aa9D7837CAf62653d097Ff85 | 0xDEf0bfBdf7530C75AB3C73f8d2F64d9eaA7aA98e |
| USDC | Arbitrum | 0xaf88d065e77c8cC2239327C5EDb3A432268e5831 | 0x5e027ad442e031424b5a2C0ad6f656662Be32882 |
| WETH | Eth Mainnet | 0xc02aaa39b223fe8d0a0e5c4f27ead9083c756cc2 | 0xD4efe33C66B8CdE33B8896a2126E41e5dB571b7e |
| WETH | Optimism | 0x4200000000000000000000000000000000000006 | 0xdD4c717a69763176d8B7A687728e228597eAB86d |
| WETH | Arbitrum | 0x82af49447d8a07e3bd95bd0d56f35241523fbab1 | 0x8e9f58E6c206CB9C98aBb9F235E0f02D65dFc922 |
| WBTC | Eth Mainnet | 0x3Eec7c855aF33280F1eD38b93059F5aa5862E3ab | 0x3Eec7c855aF33280F1eD38b93059F5aa5862E3ab |
| WBTC | Optimism | 0x68f180fcce6836688e9084f035309e29bf0a2095 | 0xE5967877065f111a556850d8f05b8DaD88edCEc9 |
| WBTC | Arbitrum | 0x2f2a2543b76a4166549f7aab2e75bef0aefc5b0f | 0x3D20c6A2b719129af175E0ff7B1875DEb360896f |

Each of the bridges have their own limits for deposits and withdrawals. For example, the USDC mainnet fast connector is subject to the following global daily limits:

- $10m USDC in deposits
- $1m USDC in withdrawals

Here is the contract interface you can use to deposit USDC via the Socket bridge: <https://etherscan.io/address/0x6d303cee7959f814042d31e0624fb88ec6fbcc1d#writeContract>

Here is a sample USDC deposit via the Socket bridge: <https://etherscan.io/tx/0x69272bbed41fd09f4b50bba6e0e451cc57a19fe81db41ac7819e003cb3088a00>

## Depositing ETH

ETH on the Derive chain is needed for certain transactions that require direct smart contract interaction. If you are manually setting up your account, you will need some small amount of eth. **Do not use this method for wETH to be used as collateral in the protocol.**

You may be able to bridge using the interface provided by the superbridge team here: <https://app.rollbridge.app/lyra>

Alternatively, to manually deposit ETH to Derive Chain, use the standard bridge. The native bridge is[OP Stack's](https://stack.optimism.io/) native bridge for deposits and withdrawals.

Deposits and withdrawals are subject to the following delays:

- Deposits confirmed in 5-10 minutes
- Withdrawals confirmed in 7 days, after the [challenge period](https://docs.optimism.io/builders/dapp-developers/bridging/messaging#understanding-the-challenge-period)

Here is the contract interface you can use to deposit via the native bridge: <https://etherscan.io/address/0x61e44dc0dae6888b5a301887732217d5725b0bff#writeProxyContract>

Here is a sample ETH deposit via the native bridge: <https://etherscan.io/tx/0x1c6b7bb4e060d2e335dfc1b3501d9e778cec1adac80652645f645a6d79daf159>

---

## institutional-trading-rewards-program

**Title:** Program Purpose
**URL:** https://docs.derive.xyz/reference/institutional-trading-rewards-program

*\*The program outlined below and all numbers provided are subject to change.*

# Program Purpose

The purpose of the Program is to support the development of the products listed below by increasing liquidity and volume in the Exchange’s central limit order book and RFQ platform, therefore, benefiting all Participants in the market.

# Program Scope

Options, Spot, and Perpetual Futures trading on the Derive Exchange.

# Eligibility

The Program is open to any firm or individual completing the application and meeting program requirements. All applicants are subject to review and approval by the Exchange. The program is opt-in so all firms must provide all relevant information to be considered for rewards.

If a participant meets the conditions of different fee tiers through market making, and total trading volume, they will enjoy their highest eligible fee tier. Wash trading is strictly prohibited and will result in disqualification from all reward programs.

# Partner Incentives

## Market Maker Rewards:

**60K $OP rewards pool per 28-day epoch** distributed to qualifying Market Makers split pro-rata based on Market Making Score.

- 40K allocated to options (80% orderbook, 20% RFQ)
- 20K allocated to perpetual futures (50% to majors, 50% to Alt markets)

Up to **$500,000 USDC per 28-day epoch exchange rebate program**

## OP Perp Taker Rewards

10k OP rewards **pool per 28-day epoch** distributed to qualifying Takers split pro-rata based on taker volume

- Requirements: > $250M notional OR > 7.5% in TAKER VOLUME

## DRV and OP Public Trading Incentives

- Traders and Market Makers earn points on net taker fees paid
- Currently, 10K OP allocated weekly in public rewards

## DRV Market Maker Incentives

1M DRV rewards pool distributed to qualifying Market Makers split pro-rata based on Market Making Score.

- 500K allocated to options
- 500K allocated to perpetual futures (50% to majors, 50% to Alt markets)
- DRV will be distributed after the airdrop period when the token is launched. 50% of the earned Market Maker DRV rewards are to be paid out immediately when DRV goes live, while the other half is subject to a 6-epoch vesting period. If the address becomes inactive (**earns rewards < 1 % of MM score**) for any epoch during the vesting period, the vested rewards will be subject to forfeiture. All rewards are subject to governance approval.

# Trading Fees

Derive's fee structure consists of:

- A maker/taker fee model
- Fee-based rebate program for market makers and volume program participants

| Fees | Perpetual Futures Maker | Perpetual Futures Taker | Spot Maker | Spot Taker | Options Maker | Options Taker |
| --- | --- | --- | --- | --- | --- | --- |
| ETH | 1bps | 3bps | 15bps | 15bps | 1bps | 3bps |
| BTC | 1bps | 3bps | 15bps | 15bps | 1bps | 3bps |
| ALT | 1bps | 3bps | 15bps | 15bps | 1bps | 3bps |

*\*Derive's matching fees are subject to change.*

# USDC Rewards Pool Overview:

Market Makers add value to the protocol by lowering the cost to trade. The Derive Exchange pays program participants on their Maker Volume, and offers discounts to high-volume Takers, through a rebate program. Trading rebates are distributed in 28-day epochs and outlined in the table below.

**$500,000 USDC rewards pool per 28-day epoch**

- Rebate program
- $250K (50%) Options
- $250K (50%) Perpetual Futures (50% to majors, 50% to Alt markets)
- Participants receive discounted fees atomically
- Negative fees are processed at a rebate at the end of each epoch
- New Market Makers eligible to receive introductory top fee tier for first full epoch of quoting
- Fee tier and rebate determined by Market Maker Ranking and Volume Trading rebates are distributed in 28-day epochs and outlined in the table below:

| MM % Score Share |  | 28-day Volume |  | 28-day Volume Share |  | stDRV Holdings | Perp Maker Fee | Perp Taker Fee | Option Maker Fee | Option Taker Fee | Spot Maker Fee | Spot Taker Fee |
| --- | --- | --- | --- | --- | --- | --- | --- | --- | --- | --- | --- | --- |
| 10% | OR | ≥ 250M | OR | ≥ 10% | AND | ≥ 500,000 stDRV | -0.010% | 0.015% | -0.010% | 0.015% | -0.010% | 0.050% |
| 10% | OR | ≥ 250M | OR | ≥ 10% | OR | ≥ 500,000 stDRV | -0.005% | 0.020% | -0.005% | 0.020% | 0.000% | 0.050% |
| 5% | OR | ≥ 100M | OR | ≥ 5% | OR | ≥ 250,000 stDRV | -0.0025% | 0.025% | 0.0025% | 0.025% | 0.050% | 0.070% |
| 1% | OR | ≥ 25M | OR | ≥ 2% | OR | ≥ 100,000 stDRV | 0.000% | 0.030% | 0.000% | 0.030% | 0.100% | 0.090% |
| All traders |  | < 25M |  |  |  |  | 0.01% | 0.030% | 0.01% | 0.030% | 0.150% | 0.150% |

**RFQ Discounts**

| 28-day RFQ Volume |  | 28-day RFQ Volume Share | RFQ Maker Fee Discount | RFQ Taker Fee Discount |
| --- | --- | --- | --- | --- |
| ≥ 100M | OR | ≥ 5% | 100% | 50% |
| ≥ 25M | OR | ≥ 2% | 25% | 10% |
| < 25M | OR | < 2% | 0% | 0% |

- MM ranking, Volume share, and stDRV holding requirements are by market, i.e. spot, options, perps separately.
- To achieve best fee tiers it requires 500k stDRV per program i.e. 1M required for both options and perps
- Discounts will apply on a percentage basis to fees on spreads and other complex fee logic and calculated pro-rata if $500K rewards is exceeded
- Fees are subject to change
- See [fee documentation](https://docs.lyra.finance/reference/fees-1) for more details

# Obligations:

The central limit order book will be snapshotted randomly, in 15 minutes windows to evaluate Market Maker performance. Optimizations to be added for supporting desired strike ranges and expiries if deemed necessary after launch. For simplification, the only requirement for MMs is to meet the minimum quote size in order for their quotes to be counted for a given snapshot:

To qualify for rewards, Traders must meet the following requirements:

- Market Maker Program:
  - Min quote size:
    - ETH, BTC, SOL, DOGE = $25k
    - Alt markets (BNB, XRP, LINK, AVAX, UNI, TAO, WIF, OP, NEAR, ARB, AAVE, INJ, BONK, TIA, SUI, ENA. PEPE, WORLDCOIN, SEI, EIG, BITCOIN, DEGEN) = $5K
  - Delta Range: only options with a Delta > 1 and < 98 will be included in the scoring snapshots
  - 28-day Notional Trading Volume > $5M
  - Trading Volume Share > 2%

# Additional Incentives:

## Increased Rate Limits

Rate limits on matching engine requests are in place as a safeguard for the exchange's order processing capacity. Participants in the Taker Incentive Program are eligible for the highest level of Matching Engine Requests. Initially, there will be two tiers for rate limiting based on Market Maker and Taker status. The exchange reserves the right to add additional rate limit tiers to be assigned based on a combination of Volume Share and Market Maker Rankings.

![](https://files.readme.io/3b8ed9789693cda59222a4bdcb7e1ead4c08c17e8159f70acda60846c3385580-image.png)

## Advanced Market Maker Protections

Participants in the program may choose to enable the following Market Maker Protections:

- Cancel\_On\_Disconnect
- Trade\_limit - max # of trades per time interval
- Quantity\_limit - max # of instruments per time interval
- Delta\_limit - max # delta per time interval
- Post\_only - order rejects if it would execute on post
- Frozen\_time - auto reset of MMPs
- Manual Reset (if desired)

### How to participate?

- **No commitments required**
- Please provide contact information [here](https://forms.gle/ivDyhdSUGGWPJsJh7)

### Monitoring and Termination of Status

The Exchange reserves the right to remove any Participant from this Program if the Exchange has determined, in good faith, that the Participant consistently and egregiously underperforms the obligations, as determined by the Exchange in its sole discretion. Moreover, the Exchange reserves the right to prohibit the participant and any affiliated entities or individuals from trading, accessing, or participating in any Exchange products and programs for an indefinite period of time, including a prohibition that extends for several years, if the Exchange determines the Participant is found to have engaged in willful misconduct of Exchange Rules.

# Appendix - Market Making Scoring Calculations

**Market Makers are scored on:**

- Market Coverage (40% for Options, 50% for perps)
- Market Quality (40% for Options, 50% for perps)
- RFQ (20% for Options, 0% for perps).

**Market Coverage, Market Quality scores are boosted by:**

- Distance from Best Market Multiplier
- Market Scaling Factor

RFQ Scores are boosted by:

- Distance from Max Cost.
- Time Score
- Market Scaling Factors

**MM Scores are then weighted by:**

- Volume Share

## Distance from Best Market Multipliers

- MM’s score for **Market Coverage and Market Quality** include a series of multipliers based on an order's price distance from best market price
- The further quotes are away from the BBO, the lower the multiplier
- The closer orders are to the BBO, the higher the multiplier
- Low-quality markets have a multiplier of 0

| Distance from BBO Options | Multiplier |
| --- | --- |
| < 0.10% | 5 |
| 0.10% < Order Price < 0.50% | 1.5 |
| 0.50% < Order Price < 1.0% | 1 |
| 1.0% < Order Price < 2.0% | 0.5 |
| > 2% | 0 |

| Distance from BBO Perps | Multiplier |
| --- | --- |
| < 0.0050% | 5 |
| 0.005% < Order Price < 0.01% | 1.5 |
| 0.01% < Order Price < 0.05% | 1 |
| 0.05% < Order Price < 0.10% | 0.5 |
| > 0.1% | 0 |

*\*Weightings and categories are subject to change, see the documentation for the most up to date.*

## Distance from Max Cost

- MM’s RFQ Score includes a series of multipliers based on an RFQ responses price distance from the maximum cost of the order
- The closer quotes are to the Max Cost, the lower the multiplier
- The more competitive the RFQ responses are, the higher the multiplier
- Low-quality markets have a multiplier of 0

| Distance from Max Cost | Multiplier |
| --- | --- |
| < 0% | 0 |
| 0 - 1% | 1 |
| 1 - 3% | 2 |
| 3%+ | 4 |

*\*Weightings and categories are subject to change, see the documentation for the most up to date.*

## Market Scaling Factor

Each market, set of expiries, or group of strikes can have a unique Market Scaling Factor to encourage liquidity. As markets mature, MSF can be set to 0 and new markets will be incentivized. Initially all market scaling factors will be set to 1.

| Market | Market Scaling Factor |
| --- | --- |
| ETH Perps | 1 |
| ETH Weekly Options < 7 DTE | 3 |
| ETH Long-Dated Options > 7 DTE | 1 |
| BTC Perps | 1 |
| BTC Weekly Options < 7 DTE | 3 |
| BTC Long-Dated Options > 7 DTE | 1 |
| SOL Perps | 1 |
| DOGE Perps | 1 |

| Alt Markets | Market Scaling Factor |
| --- | --- |
| BNB | 1 |
| XRP | 1 |
| LINK | 1 |
| AVAX | 1 |
| UNI | 1 |
| TAO | 1 |
| WIF | 1 |
| OP | 1 |
| NEAR | 1 |
| ARB | 1 |
| Aave | 1 |
| INJ | 1 |
| BONK | 1 |
| TIA | 1 |
| SUI | 1 |
| ENA | 1 |
| PEPE | 1 |
| Worldcoin | 1 |
| SEI | 1 |
| EIG | 1 |

*\*Market scaling factors are subject to change, see the documentation for the most up to date.*

## Option & Perpetual Futures Scoring

### Market Coverage 50% (40% post RFQ implementation)

- Time in Market
  - % of the time MMs quotes are on for specified strikes, and expiries. A MM is considered "on" if they are meeting min quoting obligations when the snapshot is taken.
  - √(∑ (# snapshots MM is on for *Distance from Best Market Multiplier) / # of snapshots taken)* Market Scaling Factor.
  - Let:
    - Db = Distance from Best Market Multiplier on the bid side
    - Da = Distance from Best Market Multiplier on the ask side
    - N = # of snapshots taken
    - n = # of snapshots MM is on for
    - F = Market Scaling Factor

![](https://files.readme.io/739ac9b27ff42abce3058148f14cb5f6b8e781afd332c5db448f93c609889d12-Schermafbeelding_2024-11-27_om_00.30.52.png)

### Market Quality 50% (40% post RFQ implementation)

- Book Size
  - Let:
    - Vb = MM quantity bid volume, scaled by its multiplier
    - Va = MM scaled quantity ask volume
    - Ta =Total scaled quantity bid Volume
    - Tb = Total scaled quantity ask Volume
    - Dmax = Maximum Distance from BBO Multiplier (currently 5x)
  - Total MM bid/ask volume relative to exchanges total bid-ask volume taken at each snapshot
  - Sqrt taken to smooth results
  - Scaled volumes = volume scaled by distance from BBO multipliers (I.e. $1000 notional at top of book = $5K)

![](https://files.readme.io/ebe599e1f0b6764d0236526782da862b739dddc4cc5253bf22ca6d38187a1c69-Schermafbeelding_2024-11-27_om_00.34.06.png)

### RFQ Score 20% Options (0% Perps) - Scored Separately

Note: upon launch of RFQ, Market Coverage and Market Quality Scores weighting will each be reduced by 15% and 5% respectively to account for a 20% RFQ allocation.

RFQ is scored as follows:

- For each MM, we consider all quotes for a given RFQ
- For each quote, we compute the `DistanceFromMaxCost` and the `TimeScore`
- E.g. if MM ABC posts 3 different responses to Alice’s RFQ, we compute `DistanceFromMaxCost/TimeScores` for each quote and select the maximum
- If a RFQ is filled, it is scaled by `filledScale = 1.0`, otherwise `0`

  

- *maxj(⋅)* indicates that we're taking the max combined score for a given order (i.e. the order).
  - Let:
    - RFQV = RFQ volume for the given order
    - RFQT =Total RFQ Volume is the total volume (i.e. all other quotes from unique MMs) on this order. Note that each MM’s volume is only counted once.
    - *Fi* = Market Scaling Factor
    - TS = Time Score
    - RFQD = RFQ response’s Distance From Max Cost
    - *fi* = `filledScale` which is 1.0 if the order is filled and 0.1 otherwise

![](https://files.readme.io/5920d002f784994727c4dd584f7af927b93c3e58023d6888a518bc251f711b53-Schermafbeelding_2024-11-27_om_19.55.00.png)

**RFQ Example**

Alice posts an RFQ for 100 x LONG ETH $2900/$3200 call spreads at `t=0` with

- Best bid on $2900 = $24
- Best ask on $3200 = $10
- Order book BBO = $14

At

- t = 0.4s, ABC posts a quote at $13.5 per spread - 100 ETH
- t = 1.5s, XYZ posts a quote at $13 per spread - 100 ETH
- t = 2.5s XYZ posts a fill at $12.5 per spread - 100 ETH

For ABC’s Quote:

- RFQD= ($1400-$1350)/$1400 = 3.5% = 4
- RFQV = $1350
- RFQT = ??
- TimeScore = 2
- f = 0.1

For XYZ’s Quote:

- RFQD= ($1400-$1300)/$1400 = 7.1% = 4
- RFQV = $1300
- RFQT = ??
- TimeScore = 1
- f = 0.1

For XYZ’s Fill:

- RFQD= ($1400-$1250)/$1400 = 10.7% = 4
- RFQV = $1250
- RFQT = ??
- TimeScore = 0.5
- f = 1

ABC’s `RFQScore` = 1/???√∑(1250*1*.1\*Max((4+2),(4+1))

XYZ’s `RFQScore` = √∑(($30,000*5*.5)/($90,000))  *1* 1.2 = 1

ABC gets a 5 DMM multiplier for being filled on their second quote

XYZ Gets a 5 DMM for being filled

## Trading Volume Multiplier

MMs Score is boosted by % volume traded

Volume Weight = 2.0

![](https://files.readme.io/9d3d061a3e8ee639b09cf328992a0b06ef4494c7033309c25b92ecee75a0eacc-Schermafbeelding_2024-11-27_om_01.19.26.png)

## Final Score

![](https://files.readme.io/325f9dfd5b7d64149f62b2395efa5a6b4b43c4cc8c1260c830f41c7bcdac272a-Schermafbeelding_2025-02-25_om_00.07.11.png)

## Example (Options):

|  | Market Coverage | Market Quality | Quote Efficiency | Initial Score | Initial Rank | Volume Boost | Final Score | Final Rank |
| --- | --- | --- | --- | --- | --- | --- | --- | --- |
| Metrics | Time in Market | Book Size | Msg-Volume Ratio | Score |  |  |  |  |
| Weight | 50% | 40% | 10% | 100% |  |  |  |  |
| Score | 0.8 | 0.9 | 0.001 |  |  |  |  |  |
| Results | 0.40 | 0.36 | 0.0001 | 0.761 | 5 | 1.2 | 1.09 | 3 |

  

*\*addresses designated as market maker are not eligible for the referral program.*

---

## liquidation-api

**Title:** 0. Constants & Setup
**URL:** https://docs.derive.xyz/reference/liquidation-api

Liquidations on Derive are permissionless, meaning that anyone can flag an underwater subaccount for liquidations and bid on them. However, maintaining infrastructure for monitoring margin of the subaccounts, query their balances and do bidding on-chain is cumbersome to implement, which is why Derive also provides a liquidations API which works similar to orders and RFQs.

For those interested in manual on-chain liquidations, or those interested in the details on how the auction is implemented, please refer to the full [liquidations guide here](https://docs.lyra.finance/docs/liquidations-1). This guide will focus on how to use the API to monitor live auctions and liquidate flagged users.

The main features of the liquidation auctions are as follows:

- Subaccounts get flagged for liquidation when their maintenance margin falls below zero
- A subacount gets put up for an auction
- Liquidators bid to take over a percentage of the user's subaccount, where the percentage is capped to being "the minimum needed to bring the user's account into a good state"
- The bids take over an equal percentage of every position and collateral the liquidated user is holding, for example a 50% bid on a subaccount holding 2 long ETH perps, 1 short option, 0.5 ETH spot and $2000 USDC would result in the liquidator acquiring 1 long ETH perp, 0.5 short options, 0.25 ETH spot and $1000 USDC in exchange for paying the bid price to the liquidated user
- The auction price starts off at a 5% discount relative to an oracle mark-to-market value of the whole subaccount, then decays quickly to 30% over 15 minutes, followed by a slow decay towards a 100% discount
- If the liquidated user is insolvent, the auction instead starts off at a price of zero (i.e. liquidator can take over the account without paying anything) and then decays down to a negative price (i.e. liquidator being paid by security module to take over the bad account) over 60 minutes

The liquidations flow consists of the following steps:

1. Authentication.
2. Setting up a subscription to `auctions.watch` channel which publishes the state of ongoing auctions.
3. Upon receiving the notification from `auctions.watch`, call a `private/liquidate` RPC endpoint, assuming the current bid price is acceptable for the liquidator. Note that the liquidator's subaccount gets locked until the transaction settles, which can take up to 10 seconds, no orders or other liquidations can be submitted in this locked state.
4. Transaction state can be polled using `public/get_transaction` endpoint (the `private/liquidate` response will contain a `transaction_id` to track), note that the transaction may either end up being `settled` or `reverted`, e.g. if another bid took place shortly before this bid.
5. After the transaction settles, the liquidator will receive an update to their balances over the `{subaccount_id}.balances` channel (if subscribed), and they can check the details of the liquidation using the `private/get_liquidator_history` endpoint.

# 0. Constants & Setup

This examples use the following protocol constants, subaccount IDs, etc.

TypeScript

```
import { ethers } from 'ethers';
import axios from 'axios';
import dotenv from 'dotenv';
import { WebSocket } from 'ws';

dotenv.config();

const PRIVATE_KEY = process.env.SESSION_PRIVATE_KEY as string;
const WS_ADDRESS = 'wss://api-demo.lyra.finance/ws';
const PROVIDER_URL = 'https://rpc-prod-testnet-0eakp60405.t.conduit.xyz/';
const HTTP_ADDRESS = 'https://api-demo.lyra.finance';
const ACTION_TYPEHASH = '0x4d7a9f27c403ff9c0f19bce61d76d82f9aa29f8d6d4b0c5474607d9770d1af17';
const DOMAIN_SEPARATOR = '0x9bcf4dc06df5d8bf23af818d5716491b995020f377d3b7b64c29ed14e3dd1105';
const LIQUIDATE_ADDRESS = '0x3e2a570B915fEDAFf6176A261d105A4A68a0EA8D';

const PROVIDER = new ethers.JsonRpcProvider(PROVIDER_URL);
const SIGNER = new ethers.Wallet(PRIVATE_KEY, PROVIDER);

/// if using UI: "Funding Wallet Address" under https://testnet.lyra.finance/settings
const ACCOUNT = '0x2225F2B33AA18a48EAbb675f918E950878C53BE6'

const ENCODER = ethers.AbiCoder.defaultAbiCoder();

const subaccountIdLiquidator = 36919
const AUCTIONS_CHANNEL = "auctions.watch"

```

# 1. Authentication

The first step is to set up a websocket connection and log into it - see the [Authentication](/reference/authentication) section for more. Note that the subscription to `auctions.watch` is public and does not require authentication. The private RPC calls can be executed over either websockets or HTTP, and this example will be utilizing websockets.

TypeScript

```
async function connectWs(): Promise<WebSocket> {
  return new Promise((resolve, reject) => {
    const ws = new WebSocket(WS_ADDRESS);

    ws.on('open', () => {
      setTimeout(() => resolve(ws), 50);
    });

    ws.on('error', reject);

    ws.on('close', (code: number, reason: Buffer) => {
      if (code && reason.toString()) {
        console.log(`WebSocket closed with code: ${code}`, `Reason: ${reason}`);
      }
    });
  });
};

async function signAuthenticationHeader(): Promise<{[key: string]: string}> {
  const timestamp = Date.now().toString();
  const signature = await SIGNER.signMessage(timestamp);
    return {
      wallet: ACCOUNT,
      timestamp: timestamp,
      signature: signature,
    };
}

async function loginClient(wsc: WebSocket) {
  const rpcId = Math.floor(Math.random() * 10000);
  wsc.on('message', (data: string) => {
    const response = JSON.parse(data);
    if (response.id === rpcId) {
      console.log(`Got login response with id ${rpcId}:`);
      console.log(response);
    }
  });
  const login_request = JSON.stringify({
      method: 'public/login',
      params: await signAuthenticationHeader(),
      id: rpcId
  });
  wsc.send(login_request);
  await new Promise(resolve => setTimeout(resolve, 2000));
}

```

# 2. Subscribe to `auctions.watch`

The `auctions.watch` channel publishes the state ongoing auctions, as well notifies the subscribers when the auction ends.

Sample data from this channel may look like:

JSON

```
[{
  "subaccount_id": 78202,
  "state": "ongoing",
  "timestamp": 1724715992679,
  "details": {  
     currency: null,  
     margin_type: 'SM',  
     min_cash_transfer: '1277.025016',  
     min_price_limit: '1053.257525',  
     last_seen_trade_id: 419716,  
     estimated_percent_bid: '0.175226',  
     estimated_bid_price: '6010.865122',  
     estimated_mtm: '9082.715036',  
     estimated_discount_pnl: '538.266783',  
     subaccount_balances: {  
       "USDC": '9643.945424640872731732',  
       "ETH-20241228-2500-C": '-1',  
       "ETH-PERP": '65'  
     }  
   }
}]

```

The full schema is available in the [API reference](https://docs.lyra.finance/reference/auctions-watch), and is also shown in the example below.

When a subaccount gets flagged for liquidation, a message in the above form will be emitted to subscribers. For as long as the auction is ongoing, messages will keep being re-sent with the updated data (e.g. changes in the subaccount MtM, changes in balances in case bids were made, etc.). Note that there is no guaranteed frequency for these data updates since during time of a high number of auctions, the publisher might slow down.

When the auction ends, a special message will be sent in the below format:

JSON

```
[{
  "subaccount_id": 78202,
  "state": "ended",
  "timestamp": 1724718992679,
  "details": null
}]

```

Please refer to the [API reference](https://docs.lyra.finance/reference/auctions-watch) for the explanation of each field. An informal description is provided below:

- `currency` and `margin_type` are useful as it is preferred to match the liquidated account type when bidding (due to the limits on the # of assets in the subaccounts)
- `min_cash_transfer` is how much USDC the liquidator needs to have to make a bid **and** supply enough additional funds to meet margin requirements of the acquired portfolio
- `min_price_limit` is how much USDC the liquidator needs to pay for the bid (note that this value is typically smaller than `estimated_bid_price` since it is scaled down by the maximum % that can be liquidated)
- `last_seen_trade_id` is how we ensure the liquidated subaccount's state doesn't suddenly change (due to e.g. another bid) - say a liquidator sends a bid with this field being `419716` thinking they are about to get 17.5% of the above portfolio (11.3 perps, 0.175 short options, 1687 USDC), but another bid comes through causing the portfolio to change the last millisecond, then the RPC call will error out *or* the on-chain transaction will return with status `reverted`
- `estimated_percent_bid` is approximately how much of the account can currently be liquidated
- `estimated_bid_price` is the current bid price referencing the whole portfolio (i.e. the discounted mark-to-market value)
- `estimated_mtm` is the current mark-to-market value of the liquidated portfolio
- `estimated_discount_pnl` is roughly how much in $ will the liquidator make upon successful liquidation, on mark-to-market basis, note in the example above this equals `(estimated_mtm - estimated_discount_pnl) * estimated_percent_bid`
- `subaccount_balances` are the current balances of the liquidated subaccount - the liquidator should expect to acquire each of these balances scaled by `estimated_percent_bid` in exchange for paying `estimated_percent_bid * estimated_bid_price`

Below is a simple code snippet that listens to the `auctions.watch` and writes the results to a mapping:

TypeScript

```

export interface AuctionDetailsSchema {
  currency: string | null;
  estimated_bid_price: string;
  estimated_discount_pnl: string;
  estimated_mtm: string;
  estimated_percent_bid: string;
  last_seen_trade_id: number;
  margin_type: "PM" | "SM";
  min_cash_transfer: string;
  min_price_limit: string;
  subaccount_balances: {[k: string]: string};  // asset name (as in get_subaccount) -> decimal balance
}

export type State = "ongoing" | "ended";

export interface AuctionResultSchema {
  details: AuctionDetailsSchema | null;
  state: State;
  subaccount_id: number;
  timestamp: number;
}

const AUCTIONS_STATE: {[subaccount_id: number]: AuctionDetailsSchema} = {}

async function subscribeAuctions(): Promise<WebSocket> {
    const wsc = await connectWs();

    wsc.on('message', (data: string) => {
      const response = JSON.parse(data);
      if (response.params?.channel == AUCTIONS_CHANNEL) {
        const data = response.params.data as AuctionResultSchema[];
        for (const auction of data) {
          if (auction.state === 'ongoing') {
            AUCTIONS_STATE[auction.subaccount_id] = auction.details!;
          } else {
            delete AUCTIONS_STATE[auction.subaccount_id];
          }
        }
      }
    })

    const subscribeRequest = JSON.stringify({
        method: 'subscribe',
        params: {
            channels: [AUCTIONS_CHANNEL],
        },
        id: Math.floor(Math.random() * 10000)
    });
    wsc.send(subscribeRequest);
    return wsc;
}

```

# 3. Sign and call `private/liquidate`

The `private/liquidate` endpoint is similar to deposits, withdrawals, orders and RFQs in that it requires a "signed action". The signature will be verified on-chain by the [liquidation module](https://github.com/lyra-finance/v2-matching/blob/master/src/modules/LiquidateModule.sol).

Under the hood, the module will create a temporary subaccount where some cash would be transferred to from the caller's subaccount. The temporary subaccount will then call the onchain `bid` function and acquire the percentage of the liquidated portfolio. Finally, the temporary subaccount will be "merged" back into the caller's subaccount (i.e. all positions and collaterals will be transferred back to the caller). Note that the module technically makes the merge feature optional, but the API only supports `mergeAccount=true`.

Below is a snippet that performs signing and a call to the API:

TypeScript

```
export async function sendSignedLiquidate(
  wsc: WebSocket,
  subaccountId: number,
  liquidatedSubaccountId: number,
  priceLimit: string,
  percentBid: string,
  cashTransfer: string,
  lastSeenTradeId: number,
): Promise<any> {

  const nonce = Number(`${Date.now()}${Math.round(Math.random() * 999)}`)
  const signatureExpirySec = Math.floor(Date.now() / 1000 + 600)

  // struct being signed:
  // struct LiquidationData {
  //   uint liquidatedAccountId;  // subaccount id to liquidate
  //   uint cashTransfer;  // $ transferred into a temporary subaccount from caller's subaccount
  //   uint percentOfAcc;  // % of the liquidatedAccountId to liquidate
  //   int priceLimit;  // max $ to pay for the liquidated portion
  //   uint lastSeenTradeId;  // validation in case the liquidated account changes (e.g. via someone else's bid)
  //   bool mergeAccount;  // whether to merge the temporary subaccount into the caller's subaccount, must be true
  // }

  const liquidateDataABI = ['uint256', 'uint256', 'uint256', 'int256', 'uint256', 'bool'];
  const liquidateData = [
    liquidatedSubaccountId,
    ethers.parseUnits(cashTransfer, 18),
    ethers.parseUnits(percentBid, 18),
    ethers.parseUnits(priceLimit, 18),
    lastSeenTradeId,
    true, // API only supports merging the liquidated portion into caller's subaccount
  ];
  const liquidationData = ENCODER.encode(liquidateDataABI, liquidateData);
  const hashedLiquidationData = ethers.keccak256(liquidationData);

  const actionHash = ethers.keccak256(
    ENCODER.encode(
      ['bytes32', 'uint256', 'uint256', 'address', 'bytes32', 'uint256', 'address', 'address'],
      [
        ACTION_TYPEHASH,
        subaccountId,
        nonce,
        LIQUIDATE_ADDRESS,
        hashedLiquidationData,
        signatureExpirySec,
        ACCOUNT,
        SIGNER.address
      ]
    )
    );

  const signature = SIGNER.signingKey.sign(
    ethers.keccak256(Buffer.concat([
      Buffer.from("1901", "hex"),
      Buffer.from(DOMAIN_SEPARATOR.slice(2), "hex"),
      Buffer.from(actionHash.slice(2), "hex")
    ]))
  ).serialized;

  const rpcId = Math.floor(Math.random() * 10000);
  let liquidateResp = undefined;
  wsc.on('message', (data: string) => {
    const response = JSON.parse(data);
    if (response.id === rpcId) {
      console.log(`Got liquidate response with id ${rpcId}:`);
      console.log(response);
      liquidateResp = response;
    }
  });

  const params = {
    subaccount_id: subaccountId,
    liquidated_subaccount_id: liquidatedSubaccountId,
    cash_transfer: cashTransfer,
    percent_bid: percentBid,
    price_limit: priceLimit,
    last_seen_trade_id: lastSeenTradeId,
    signature: signature,
    signature_expiry_sec: signatureExpirySec,
    nonce: nonce,
    signer: SIGNER.address,
  }

  console.log(`Sending liquidate request with id ${rpcId}:`);
  console.log(params);

  wsc.send(JSON.stringify({
    method: 'private/liquidate',
    params: params,
    id: rpcId
  }));

  for (let i = 0; i < 10; i++) {
    if (liquidateResp) {
      break;
    }
    await new Promise(resolve => setTimeout(resolve, 1000));
  }

  return liquidateResp;
}

```

Note that the `auctions.watch` channel can help with figuring out values for some of the call parameters such as `cash_transfer`and `last_seen_trade_id`. For example, here's how one could send a liquidation call using the auctions channel output:

TypeScript

```
async function liquidateTest() {
  const wscSub = await subscribeAuctions();
  await new Promise(resolve => setTimeout(resolve, 5000));
  console.log(AUCTIONS_STATE);
  const wsc = await connectWs();
  await loginClient(wsc);
  let liquidateResp: any = undefined;

  for (const subaccount_id in AUCTIONS_STATE) {
    const auction = AUCTIONS_STATE[subaccount_id];
    if (auction.margin_type !== 'SM') {
      continue; // only liquidate same margin type / currency combination as your subaccount
    }

    console.log(`Liquidating subaccount ${subaccount_id} with auction details:`);
    console.log(auction);

    // add small buffer to price limit and cash transfer in case market moves (here 5% of mtm or $10)
    const buffer = Math.max(Math.abs((+auction.estimated_mtm * 0.05)), 10)

    // cashTransfer has to be strictly greater than 0
    const cashTransfer = (+auction.min_cash_transfer == 0 ? 0.1 : +auction.min_cash_transfer + buffer).toFixed(2);

    liquidateResp = await sendSignedLiquidate(
      wsc,
      subaccountIdLiquidator,
      Number(subaccount_id),
      (+auction.min_price_limit + buffer).toFixed(2),
      '1', // liquidate "up to" 100% of the subaccount (actual % may be less, see auction.estimated_percent_bid)
      cashTransfer,
      auction.last_seen_trade_id,
    );
    break; // for the sake of example liquidate just the first auction
  }
}

```

# 4. Get transaction status

The RPC response to `private/liquidate` will contain `transaction_id` that can be polled via `public/get_transaction`.

TypeScript

```
// Also available over websockets
async function getTxStatus(txId: string) {
  const resp = await axios.post(`${HTTP_ADDRESS}/public/get_transaction`, {transaction_id: txId})
  return resp.data.result;
}

```

A sample response looks like:

JSON

```
{
  "result": {
    "status": "settled",
    "transaction_hash": "0x5ee256e742f6d88366f0bcdd76ce991b8785ddb5533a04994f2fd53c8b6e699e",
    "data": {
      "data": {
        "percent_bid": "1",
        "price_limit": "1144.08",
        "cash_transfer": "1367.85",
        "merge_account": true,
        "last_seen_trade_id": 419716,
        "liquidated_subaccount_id": 78202
      },
      "nonce": 1724715980786660,
      "owner": "0x2225F2B33AA18a48EAbb675f918E950878C53BE6",
      "expiry": 1724716280,
      "module": "0x3e2a570B915fEDAFf6176A261d105A4A68a0EA8D",
      "signer": "0xb94dCcaDf0c72E4A472f6ccf07595Ba27B49e033",
      "signature": "0xd328dfd529b975f74821e518007090015fa7c409fe803377ef496b0f3c305f010538d64e8264b16bc631d143df7f36b28016e58faf00240dac27d04adf75c5c71b",
      "subaccount_id": 36919
    },
    "error_log": {}
  },
  "id": "682f6ea4-a801-4306-bbd5-07fadf983fd6"
}

```

If the `status` field is either `settled` or `reverted`, then it is safe to assume that the transaction has been finalized (succeed or failed, respectively), and the liquidator can react accordingly.

Others statuses can be find in the [reference](https://docs.lyra.finance/reference/post_public-get-transaction), most commonly a transaction will be `pending` for a few seconds before finalizing.

# 5. Get balances and history

Only in the event of a successful `settled` transaction will the `{subaccount_id}.balances` [channel](https://docs.lyra.finance/reference/subaccount_id-balances) publish a balance update with the `update_type` of `liquidator`. This channel is currently the best way to keep balances in sync.

Additionally, a `settled` liquidation transaction will be recorded and will be viewable in the `private/get_liquidator_history` endpoint:

TypeScript

```
// get and log liquidator history
async function getLiquidatorHist(wsc: WebSocket) {
  // avaiable as a regular HTTP call as well
  const rpcId = Math.floor(Math.random() * 10000);
  wsc.on('message', (data: string) => {
    const response = JSON.parse(data);
    if (response.id === rpcId) {
      console.log(`Got liquidator history response with id ${rpcId}:`);
      console.log(JSON.stringify(response, null, 2));
    }
  });

  const params = {
    subaccount_id: subaccountIdLiquidator,
  }

  console.log(`Sending liquidator history request with id ${rpcId}:`);
  console.log(params);

  wsc.send(JSON.stringify({
    method: 'private/get_liquidator_history',
    params: params,
    id: rpcId
  }));
}

```

Sample output :

JSON

```
{
  "id": 4922,
  "result": {
    "bids": [
      {
        "timestamp": 1724715992679793,
        "tx_hash": "0x5ee256e742f6d88366f0bcdd76ce991b8785ddb5533a04994f2fd53c8b6e699e",
        "realized_pnl": "180.03751237036511212",
        "realized_pnl_excl_fees": "180.03751237036511212",
        "discount_pnl": "524.1153982323594391345977783203125",
        "percent_liquidated": "0.173325903015437602",
        "cash_received": "-1043.19768569751330486",
        "amounts_liquidated": {
          "ETH-PERP": "11.26618369600344413",
          "ETH-20241228-2500-C": "-0.173325903015437602",
          "USDC": "1674.313180473687227546"
        },
        "positions_realized_pnl": {
          "ETH-PERP": "180.03751237036511212",
        },
        "positions_realized_pnl_excl_fees": {
          "ETH-PERP": "180.03751237036511212",
        }
      }
 ]}}

```

## Putting it all together

TypeScript

```
import { ethers } from 'ethers';
import axios from 'axios';
import dotenv from 'dotenv';
import { WebSocket } from 'ws';

dotenv.config();

const PRIVATE_KEY = process.env.SESSION_PRIVATE_KEY as string;
const WS_ADDRESS = 'wss://api-demo.lyra.finance/ws';
const PROVIDER_URL = 'https://rpc-prod-testnet-0eakp60405.t.conduit.xyz/';
const HTTP_ADDRESS = 'https://api-demo.lyra.finance';
const ACTION_TYPEHASH = '0x4d7a9f27c403ff9c0f19bce61d76d82f9aa29f8d6d4b0c5474607d9770d1af17';
const DOMAIN_SEPARATOR = '0x9bcf4dc06df5d8bf23af818d5716491b995020f377d3b7b64c29ed14e3dd1105';
const LIQUIDATE_ADDRESS = '0x3e2a570B915fEDAFf6176A261d105A4A68a0EA8D';

const PROVIDER = new ethers.JsonRpcProvider(PROVIDER_URL);
const SIGNER = new ethers.Wallet(PRIVATE_KEY, PROVIDER);

/// if using UI: "Funding Wallet Address" under https://testnet.lyra.finance/settings
const ACCOUNT = '0x2225F2B33AA18a48EAbb675f918E950878C53BE6'

const ENCODER = ethers.AbiCoder.defaultAbiCoder();

const subaccountIdLiquidator = 36919
const AUCTIONS_CHANNEL = "auctions.watch"

async function connectWs(): Promise<WebSocket> {
  return new Promise((resolve, reject) => {
    const ws = new WebSocket(WS_ADDRESS);

    ws.on('open', () => {
      setTimeout(() => resolve(ws), 50);
    });

    ws.on('error', reject);

    ws.on('close', (code: number, reason: Buffer) => {
      if (code && reason.toString()) {
        console.log(`WebSocket closed with code: ${code}`, `Reason: ${reason}`);
      }
    });
  });
};

async function signAuthenticationHeader(): Promise<{[key: string]: string}> {
  const timestamp = Date.now().toString();
  const signature = await SIGNER.signMessage(timestamp);
    return {
      wallet: ACCOUNT,
      timestamp: timestamp,
      signature: signature,
    };
}

async function loginClient(wsc: WebSocket) {
  const rpcId = Math.floor(Math.random() * 10000);
  wsc.on('message', (data: string) => {
    const response = JSON.parse(data);
    if (response.id === rpcId) {
      console.log(`Got login response with id ${rpcId}:`);
      console.log(response);
    }
  });
  const login_request = JSON.stringify({
      method: 'public/login',
      params: await signAuthenticationHeader(),
      id: rpcId
  });
  wsc.send(login_request);
  await new Promise(resolve => setTimeout(resolve, 2000));
}

export interface AuctionDetailsSchema {
  currency: string | null;
  estimated_bid_price: string;
  estimated_discount_pnl: string;
  estimated_mtm: string;
  estimated_percent_bid: string;
  last_seen_trade_id: number;
  margin_type: "PM" | "SM";
  min_cash_transfer: string;
  min_price_limit: string;
  subaccount_balances: {[k: string]: string};  // asset name (as in get_subaccount) -> decimal balance
}

export type State = "ongoing" | "ended";

export interface AuctionResultSchema {
  details: AuctionDetailsSchema | null;
  state: State;
  subaccount_id: number;
  timestamp: number;
}

const AUCTIONS_STATE: {[subaccount_id: number]: AuctionDetailsSchema} = {}

async function subscribeAuctions(): Promise<WebSocket> {
    const wsc = await connectWs();

    wsc.on('message', (data: string) => {
      const response = JSON.parse(data);
      if (response.params?.channel == AUCTIONS_CHANNEL) {
        const data = response.params.data as AuctionResultSchema[];
        for (const auction of data) {
          if (auction.state === 'ongoing') {
            AUCTIONS_STATE[auction.subaccount_id] = auction.details!;
          } else {
            delete AUCTIONS_STATE[auction.subaccount_id];
          }
        }
      }
    })

    const subscribeRequest = JSON.stringify({
        method: 'subscribe',
        params: {
            channels: [AUCTIONS_CHANNEL],
        },
        id: Math.floor(Math.random() * 10000)
    });
    wsc.send(subscribeRequest);
    return wsc;
}

export async function sendSignedLiquidate(
  wsc: WebSocket,
  subaccountId: number,
  liquidatedSubaccountId: number,
  priceLimit: string,
  percentBid: string,
  cashTransfer: string,
  lastSeenTradeId: number,
): Promise<any> {

  const nonce = Number(`${Date.now()}${Math.round(Math.random() * 999)}`)
  const signatureExpirySec = Math.floor(Date.now() / 1000 + 600)

  // struct being signed:
  // struct LiquidationData {
  //   uint liquidatedAccountId;  // subaccount id to liquidate
  //   uint cashTransfer;  // $ transferred into a temporary subaccount from caller's subaccount
  //   uint percentOfAcc;  // % of the liquidatedAccountId to liquidate
  //   int priceLimit;  // max $ to pay for the liquidated portion
  //   uint lastSeenTradeId;  // validation in case the liquidated account changes (e.g. via someone else's bid)
  //   bool mergeAccount;  // whether to merge the temporary subaccount into the caller's subaccount, must be true
  // }

  const liquidateDataABI = ['uint256', 'uint256', 'uint256', 'int256', 'uint256', 'bool'];
  const liquidateData = [
    liquidatedSubaccountId,
    ethers.parseUnits(cashTransfer, 18),
    ethers.parseUnits(percentBid, 18),
    ethers.parseUnits(priceLimit, 18),
    lastSeenTradeId,
    true, // API only supports merging the liquidated portion into caller's subaccount
  ];
  const liquidationData = ENCODER.encode(liquidateDataABI, liquidateData);
  const hashedLiquidationData = ethers.keccak256(liquidationData);

  const actionHash = ethers.keccak256(
    ENCODER.encode(
      ['bytes32', 'uint256', 'uint256', 'address', 'bytes32', 'uint256', 'address', 'address'],
      [
        ACTION_TYPEHASH,
        subaccountId,
        nonce,
        LIQUIDATE_ADDRESS,
        hashedLiquidationData,
        signatureExpirySec,
        ACCOUNT,
        SIGNER.address
      ]
    )
    );

  const signature = SIGNER.signingKey.sign(
    ethers.keccak256(Buffer.concat([
      Buffer.from("1901", "hex"),
      Buffer.from(DOMAIN_SEPARATOR.slice(2), "hex"),
      Buffer.from(actionHash.slice(2), "hex")
    ]))
  ).serialized;

  const rpcId = Math.floor(Math.random() * 10000);
  let liquidateResp = undefined;
  wsc.on('message', (data: string) => {
    const response = JSON.parse(data);
    if (response.id === rpcId) {
      console.log(`Got liquidate response with id ${rpcId}:`);
      console.log(response);
      liquidateResp = response;
    }
  });

  const params = {
    subaccount_id: subaccountId,
    liquidated_subaccount_id: liquidatedSubaccountId,
    cash_transfer: cashTransfer,
    percent_bid: percentBid,
    price_limit: priceLimit,
    last_seen_trade_id: lastSeenTradeId,
    signature: signature,
    signature_expiry_sec: signatureExpirySec,
    nonce: nonce,
    signer: SIGNER.address,
  }

  console.log(`Sending liquidate request with id ${rpcId}:`);
  console.log(params);

  wsc.send(JSON.stringify({
    method: 'private/liquidate',
    params: params,
    id: rpcId
  }));

  for (let i = 0; i < 10; i++) {
    if (liquidateResp) {
      break;
    }
    await new Promise(resolve => setTimeout(resolve, 1000));
  }

  return liquidateResp;
}

async function getLiquidatorHist(wsc: WebSocket) {
  // avaiable as a regular HTTP call as well
  const rpcId = Math.floor(Math.random() * 10000);
  wsc.on('message', (data: string) => {
    const response = JSON.parse(data);
    if (response.id === rpcId) {
      console.log(`Got liquidator history response with id ${rpcId}:`);
      console.log(JSON.stringify(response, null, 2));
    }
  });

  const params = {
    subaccount_id: subaccountIdLiquidator,
  }

  console.log(`Sending liquidator history request with id ${rpcId}:`);
  console.log(params);

  wsc.send(JSON.stringify({
    method: 'private/get_liquidator_history',
    params: params,
    id: rpcId
  }));

  return rpcId;
}

async function getTxStatus(txId: string) {
  const resp = await axios.post(`${HTTP_ADDRESS}/public/get_transaction`, {transaction_id: txId})
  return resp.data.result;
}

async function liquidateTest() {
  const wscSub = await subscribeAuctions();
  await new Promise(resolve => setTimeout(resolve, 5000));
  console.log(AUCTIONS_STATE);
  const wsc = await connectWs();
  await loginClient(wsc);
  let liquidateResp: any = undefined;

  for (const subaccount_id in AUCTIONS_STATE) {
    const auction = AUCTIONS_STATE[subaccount_id];
    if (auction.margin_type !== 'SM') {
      continue; // only liquidate same margin type / currency combination as your subaccount
    }

    console.log(`Liquidating subaccount ${subaccount_id} with auction details:`);
    console.log(auction);

    // add small buffer to price limit and cash transfer in case market moves (here 5% of mtm or $10)
    const buffer = Math.max(Math.abs((+auction.estimated_mtm * 0.05)), 10)

    // cashTransfer has to be strictly greater than 0
    const cashTransfer = (+auction.min_cash_transfer == 0 ? 0.1 : +auction.min_cash_transfer + buffer).toFixed(2);

    liquidateResp = await sendSignedLiquidate(
      wsc,
      subaccountIdLiquidator,
      Number(subaccount_id),
      (+auction.min_price_limit + buffer).toFixed(2),
      '1', // liquidate "up to" 100% of the subaccount (actual % may be less, see auction.estimated_percent_bid)
      cashTransfer,
      auction.last_seen_trade_id,
    );
    break; // for the sake of example liquidate just the first auction
  }

  if (liquidateResp === undefined) {
    console.log('Could not get a liquidation result');
    return;
  }

  for (let i = 0; i < 10; i++) {
    let tx = await getTxStatus(liquidateResp.result.transaction_id);
    if (tx.status === 'settled') {
      console.log('Liquidation successful');
      break;
    }
    else {
      console.log(`Liquidation status: ${tx.status}`);
    }
    await new Promise(resolve => setTimeout(resolve, 1000));
  }

  await getLiquidatorHist(wsc);
}

liquidateTest();

```

---

## multiple-subaccounts

**Title:** Multiple Subaccounts
**URL:** https://docs.derive.xyz/reference/multiple-subaccounts

You may create several subaccounts all managed by the same Ethereum Wallet / Smart-contract Wallet. To do so, go to the "Subaccounts" page in the "Account Settings" dropdown:

![](https://files.readme.io/073874b-image.png)

Here you will see your first subaccount with deposited funds. You may now click "Create Subaccount" and choose an account label and subaccount margin type:

![](https://files.readme.io/66dc92e-image.png)

Once created, you can click on "Transfer" to move funds from your first subaccount to the newly created one.

![](https://files.readme.io/a4996db-image.png)

---

## on-chain-withdraw

**Title:** On Chain Withdraw
**URL:** https://docs.derive.xyz/reference/on-chain-withdraw

The withdrawal flow guide is WIP, however it is nearly identical to the deposit flow.

We also have a `public/withdraw_debug` route to help pin point any hashing / signing issues.

---

## open-orders-margin

**Title:** Maintenance Margin
**URL:** https://docs.derive.xyz/reference/open-orders-margin

# Maintenance Margin

The orderbook calculates maintenance margin in accordance with the protocol rules, and makes the values available over various endpoints such as `get_subaccount`. For more details refer to the [standard](https://docs.lyra.finance/docs/standard-margin-1) and [portfolio](https://docs.lyra.finance/docs/portfolio-margin-1) margin sections.

# Initial Margin

Similarly, the orderbook calculates initial margin using protocol rules, and ensures that no trade gets sent for settlement if the margin value after trade would be insufficient. Refer to the aforementioned sections for more detail.

# Open Orders Margin

Limit orders that stay open in the book require that the account has extra margin to cover them if they were to get filled.

The orderbook backend will inspect account's open orders`[order_1, order_2,...]`and find a "worst subset" of these orders, where "worst" is defined as a set of orders that, if filled, leads to the smallest *initial margin* possible. While performing those simulated fills, the backend will take into account the premiums paid or received for option bids and asks, as well as the current positions owned by the account.

For example, suppose the open orders and positions are:

- Orders:`[bid 10 perps @ $1999, ask 100 perps @ $2001, bid 10 1w calls @ $55, bid 5 2w calls @ $75]`
- Positions:`[long 90 perps]`

The backend will try and group the orders by their delta and / or vega sign and arrive at a conclusion that `[bid 10 perps @ $1999, bid 10 calls @ $55, bid 5 2w calls @ $75]` is the worst fill scenario. The open order margin for these orders will then be calculated by finding how much *extra initial margin* would the account require if those orders were to get filled.

For every new open (i.e. non-crossing) order that arrives to the orderbook, the risk engine checks if the sum of current initial margin and the open orders margin is non-negative. In other words, new orders are accepted as long as the account can honour the "worst fill" scenario.

The `private/get_subaccount` endpoint can be used to check which orders have been flagged as "worst subset" and how much open orders margin they require.

## Market Maker Protections (MMPs) and Open Orders Margin

Oftentimes portfolio margin users would be market makers quoting hundreds of assets at the same time. If they have tight MMP limits, it is impossible for them to get filled on all of these quotes simultaneously, so it would be unreasonably capital inefficient to require them to lock margin for very large subsets of orders.

Therefore, for portfolio margin accounts, the process of finding the "worst subset" is constrained by account's market maker protection settings. *If MMP amount limits are enabled*, the "worst subset" of orders would be reduced to an orders subset which can be filled subject to staying within MMP amount limit. Note that the reduced subset cannot be smaller than 2 distinct assets, i.e. the smallest possible open orders margin requirement still enforces that the market maker can honour at least 2 fills on two of the "worst" assets they are quoting.

Using the above example - if MMP amount limit was set to 3, then the worst orders subset would exclude the `bid 5 2w calls @ $75` and would just consist of `[bid 10 perps @ $1999, bid 10 calls @ $55]`, because at least 2 assets have to be fillable by the account.

If the MMP limit was too high (e.g. 30), then the subset would remain unchanged.

Finally note that only MMP *amount* limit supports this capital efficiency improvement, the *delta* limit is ignored. More info on the MMPs can be found in the API reference under `private/set_mmp_config`.

---

## rfq-quoting-and-execution

**Title:** 0. Constants & Setup
**URL:** https://docs.derive.xyz/reference/rfq-quoting-and-execution

Similar to orderbook trading, RFQs are "self-custodial", and they require signed messages to be settled. Those signed messages guarantee that all legs of an RFQ will execute at the specified prices and amounts, as well as that the fee charged by the orderbook does not exceed the signed `max_fee`.

Unlike orderbook trading, makers and takers follow different rules and sign slightly different messages in order to complete an RFQ. The full flow is below:

1. **[Taker & Maker]** Authentication
2. **[Taker]** Send RFQ
3. **[Maker]** Listen or poll for RFQs
4. **[Maker]** In response to an RFQ, sign and send a quote
5. **[Taker]** Poll for the Quotes (market makers' replies to RFQs) and pick the best one
6. **[Taker]** Sign an execute message for the selected quote

NOTE: to greatly simplify signatures, you may use our [Derive Python Signing SDK](https://pypi.org/project/derive_action_signing/) (see the Python RFQ example on the next page).

# 0. Constants & Setup

This examples use the following protocol constants, subaccount IDs, leg instruments, etc.

TypeScript

```
import { ethers } from 'ethers';
import axios from 'axios';
import dotenv from 'dotenv';

dotenv.config();

const PRIVATE_KEY = process.env.OWNER_PRIVATE_KEY as string;
const PROVIDER_URL = 'https://l2-prod-testnet-0eakp60405.t.conduit.xyz';
const HTTP_ADDRESS = 'https://api-demo.lyra.finance';
const ACTION_TYPEHASH = '0x4d7a9f27c403ff9c0f19bce61d76d82f9aa29f8d6d4b0c5474607d9770d1af17';
const DOMAIN_SEPARATOR = '0x9bcf4dc06df5d8bf23af818d5716491b995020f377d3b7b64c29ed14e3dd1105';
const OPTION_ADDRESS = '0xBcB494059969DAaB460E0B5d4f5c2366aab79aa1';
const RFQ_ADDRESS = '0x4E4DD8Be1e461913D9A5DBC4B830e67a8694ebCa'

const PROVIDER = new ethers.JsonRpcProvider(PROVIDER_URL);
const wallet = new ethers.Wallet(PRIVATE_KEY, PROVIDER);
const encoder = ethers.AbiCoder.defaultAbiCoder();

const subaccount_id_rfq = 23525
const subaccount_id_maker = 8

const LEG_1_NAME = 'ETH-20240329-2400-C'
const LEG_2_NAME = 'ETH-20240329-2600-C'

// can retreive with public/get_instrument
const LEGS_TO_SUB_ID: any = {
  'ETH-20240329-2400-C': '39614082287924319838483674368',
  'ETH-20240329-2600-C': '39614082373823665758483674368'
}

```

# 1. Authentication (both makers and takers)

In this guide, we'll use REST API for all examples / requests. As such, auth is done over headers as described in the [Authentication](/reference/authentication) section:

TypeScript

```
async function signAuthenticationHeader() {
  const timestamp = Date.now().toString();
  const signature = await wallet.signMessage(timestamp);
  return {
      "X-LyraWallet": wallet.address,
      "X-LyraTimestamp": timestamp,
      "X-LyraSignature": signature
  };
}

```

# 2. Send RFQ

Takers send RFQs, which do *not* specify the direction of execution.

TypeScript

```
function createRfqObject(): object {
  const rfq = {
    subaccount_id: subaccount_id_rfq,
    // NOTE: legs MUST be sorted by instrument_name where sorting key is instrument_name
    legs: [
      {
        instrument_name: LEG_1_NAME,
        amount: '3',
        direction: 'buy'
      },
      {
        instrument_name: LEG_2_NAME,
        amount: '3',
        direction: 'sell'
      }
    ],
  };
  return rfq;
}

async function sendRfq(rfq: object) {
  const authHeader = await signAuthenticationHeader();
  const resp = await axios.post(`${HTTP_ADDRESS}/private/send_rfq`, rfq, {headers: authHeader})
  return resp.data.result;
}

```

# 3. Listen or poll for RFQs

Note that market maker wallets must be approved by the support team in order to get access to the maker API. To get live RFQs, one can use either the polling endpoint (`poll_rfqs`) or the `{wallet}.rfqs` channel. Below example shows uses the `poll_rfqs` endpoint.

TypeScript

```
// NOTE: types defined in this example are just for illustration and are not robust,
// use the docs to get more info such as allowed enum values, etc.
type RfqLeg = {
  instrument_name: string,
  amount: string,
  direction: 'buy' | 'sell'
}

type RfqResponse = {
  subaccount_id: number,
  creation_timestamp: number,
  last_update_timestamp: number,
  status: string,
  cancel_reason: string,
  rfq_id: string,
  valid_until: number,
  legs: Array<RfqLeg>
}

async function pollRfq() : Promise<RfqResponse> {
  // account owner of the subaccount_id must be approved to act as RFQ maker
  // can also use {wallet}.rfqs channel to listen for RFQs (same response format)
  const authHeader = await signAuthenticationHeader();
  const resp = await axios.post(`${HTTP_ADDRESS}/private/poll_rfqs`, {subaccount_id: subaccount_id_maker, status: 'open'}, {headers: authHeader})
  // for the sake of example just return the latest RFQ
  return resp.data.result.rfqs[0]
}

```

# 4. In response to an RFQ, sign and send a quote

When the execution occurs, the `RfqModule.sol` contract will validate the maker and taker signatures, effectively ensuring that the two parties "agreed" on all of the leg names, amounts and prices.

Quotes can be sent in either `buy` or `sell` direction. The `buy` quote will execute the legs in the same direction as the legs' definition (e.g. a `buy` quote on a `sell` call option leg will be executed as a short call). A `sell` quote flips the direction of every leg in the RFQ. Note that the quote direction affects signature logic, because the contracts work with *signed* leg amounts.

TypeScript

```
type QuoteLeg = {
  instrument_name: string,
  amount: string,
  direction: 'buy' | 'sell',
  price: string
}

type EncodedLeg = [string, string, ethers.BigNumberish, ethers.BigNumberish]

function encodePricedLegs(legs: Array<QuoteLeg>, direction: 'buy' | 'sell'): Array<EncodedLeg> {
  const dirSign = BigInt(direction === 'buy' ? 1 : -1);
  const encoded_legs : Array<EncodedLeg> = legs.map((leg) => {
    const subid = LEGS_TO_SUB_ID[leg.instrument_name];
    const legSign = BigInt(leg.direction === 'buy' ? 1 : -1);
    const signedAmount = ethers.parseUnits(leg.amount, 18) * legSign * dirSign;
    return [OPTION_ADDRESS, subid, ethers.parseUnits(leg.price, 18), signedAmount];
  });
  return encoded_legs;
}

function encodeQuoteData(encoded_legs: Array<EncodedLeg>, max_fee: string): string {
  const rfqData = [ethers.parseUnits(max_fee, 18), encoded_legs];
  const QuoteDataABI = ['(uint,(address,uint,uint,int)[])'];
  const encodedData = encoder.encode(QuoteDataABI, [rfqData]);
  const hashedData = ethers.keccak256(Buffer.from(encodedData.slice(2), 'hex'));
  return hashedData;
}

function signAction(action: any, actionData: string) {
  const action_hash = ethers.keccak256(
    encoder.encode(
      ['bytes32', 'uint256', 'uint256', 'address', 'bytes32', 'uint256', 'address', 'address'],
      [
        ACTION_TYPEHASH,
        action.subaccount_id,
        action.nonce,
        RFQ_ADDRESS,
        actionData,
        action.signature_expiry_sec,
        wallet.address,
        action.signer
      ]
    )
  );
  action.signature = wallet.signingKey.sign(
    ethers.keccak256(Buffer.concat([
      Buffer.from("1901", "hex"),
      Buffer.from(DOMAIN_SEPARATOR.slice(2), "hex"),
      Buffer.from(action_hash.slice(2), "hex")
    ]))
  ).serialized;
}

function signQuote(quote: any) {
  const encoded_legs = encodePricedLegs(quote.legs, quote.direction);
  const quoteData = encodeQuoteData(encoded_legs, quote.max_fee);
  signAction(quote, quoteData)
}

function createQuoteObject(rfq_response: RfqResponse, direction: 'buy' | 'sell') : object {
  const pricedLegs: Array<any> = rfq_response.legs;
  pricedLegs[0].price = direction == 'buy' ? '160' : '180';
  pricedLegs[1].price = direction == 'buy' ? '70' : '50';
  return {
    subaccount_id: subaccount_id_maker,
    rfq_id: rfq_response.rfq_id,
    legs: pricedLegs,
    direction: direction,
    max_fee: '10',
    nonce: Number(`${Date.now()}${Math.round(Math.random() * 999)}`),
    signer: wallet.address,
    signature_expiry_sec: Math.floor(Date.now() / 1000 + 350),
    signature: "filled_in_below"
  };
}

async function sendQuote(rfq_response: RfqResponse, direction: 'buy' | 'sell') {
  const quote = createQuoteObject(rfq_response, direction);
  signQuote(quote);
  const authHeader = await signAuthenticationHeader();
  const resp = await axios.post(`${HTTP_ADDRESS}/private/send_quote`, quote, {headers: authHeader})
  return resp.data.result;
}

```

# 5. Poll for the Quotes and pick the best one

Takers can poll the quotes, and use the polled object's fields to sign an execute message in the next step.

TypeScript

```
type QuoteResultPublicSchema = {
  cancel_reason: string;
  creation_timestamp: number;
  direction: 'buy' | 'sell';
  last_update_timestamp: number;
  legs: Array<QuoteLeg>;
  legs_hash: string;
  liquidity_role: 'maker' | 'taker';
  quote_id: string;
  rfq_id: string;
  status: string;
  subaccount_id: number;
  tx_hash: string | null;
  tx_status: string | null;
}

async function pollQuotes(rfq_id: string): Promise<Array<QuoteResultPublicSchema>> {
  const authHeader = await signAuthenticationHeader();
  const resp = await axios.post(`${HTTP_ADDRESS}/private/poll_quotes`, {subaccount_id: subaccount_id_rfq, rfq_id: rfq_id, status: 'open'}, {headers: authHeader})
  return resp.data.result.quotes;
}

```

# 6. Sign an execute message for the selected quote

Signing an execute message is very similar to the quote signing, except the type signatures differ a little:

- Market makers sign `{'max_fee': uint, legs: EncodedLeg[]}`
- Takers sign `{'max_fee': uint, legs_hash: bytes32}`

The `legs_hash` is simply a keccak256-hashed array of the same legs as what market makers sign in their quote. The smart contract ensures that the two parties agreed on the leg amounts / prices etc. by hashing maker's array of legs and comparing it to the `legs_hash`.

TypeScript

```
function encodeExecuteData(encoded_legs: Array<EncodedLeg>, max_fee: string): string {
  const encoder = ethers.AbiCoder.defaultAbiCoder();
  const orderHashABI = ['(address,uint,uint,int)[]'];
  const orderHash = ethers.keccak256(Buffer.from(encoder.encode(orderHashABI, [encoded_legs]).slice(2), 'hex'));
  const ExectuteDataABI = ['bytes32', 'uint'];
  const encodedData = encoder.encode(ExectuteDataABI, [orderHash, ethers.parseUnits(max_fee, 18)]);
  const hashedData = ethers.keccak256(Buffer.from(encodedData.slice(2), 'hex'));
  return hashedData;
}

function signExecute(execute: any) {
  const encoded_legs = encodePricedLegs(execute.legs, execute.direction === 'buy' ? 'sell' : 'buy');
  const executeData = encodeExecuteData(encoded_legs, execute.max_fee);
  signAction(execute, executeData)
}

function createExecuteObject(quote: QuoteResultPublicSchema) : object {
  return {
    subaccount_id: subaccount_id_rfq,
    quote_id: quote.quote_id,
    rfq_id: quote.rfq_id,
    direction: quote.direction === 'buy' ? 'sell' : 'buy',
    max_fee: '10',
    nonce: Number(`${Date.now()}${Math.round(Math.random() * 999)}`),
    signer: wallet.address,
    signature_expiry_sec: Math.floor(Date.now() / 1000 + 350),
    legs: quote.legs,
    signature: "filled_in_below"
  }
}

async function sendExecute(quote: QuoteResultPublicSchema) {
  const execute = createExecuteObject(quote);
  signExecute(execute);
  const authHeader = await signAuthenticationHeader();
  const resp = await axios.post(`${HTTP_ADDRESS}/private/execute_quote`, execute, {headers: authHeader})
  return resp.data.result;
}

```

## Putting it all together

Below is an example of the end-to-end RFQ flow, from creating RFQs, to signing maker quotes, to executing them. For illustration purposes the same account is used (the account owns two different subaccounts).

TypeScript

```
import { ethers } from 'ethers';
import axios from 'axios';
import dotenv from 'dotenv';

dotenv.config();

const PRIVATE_KEY = process.env.OWNER_PRIVATE_KEY as string;
const PROVIDER_URL = 'https://l2-prod-testnet-0eakp60405.t.conduit.xyz';
const HTTP_ADDRESS = 'https://api-demo.lyra.finance';
const ACTION_TYPEHASH = '0x4d7a9f27c403ff9c0f19bce61d76d82f9aa29f8d6d4b0c5474607d9770d1af17';
const DOMAIN_SEPARATOR = '0x9bcf4dc06df5d8bf23af818d5716491b995020f377d3b7b64c29ed14e3dd1105';
const OPTION_ADDRESS = '0xBcB494059969DAaB460E0B5d4f5c2366aab79aa1';
const RFQ_ADDRESS = '0x4E4DD8Be1e461913D9A5DBC4B830e67a8694ebCa'

const PROVIDER = new ethers.JsonRpcProvider(PROVIDER_URL);
const wallet = new ethers.Wallet(PRIVATE_KEY, PROVIDER);
const encoder = ethers.AbiCoder.defaultAbiCoder();

const subaccount_id_rfq = 23525
const subaccount_id_maker = 8

const LEG_1_NAME = 'ETH-20240329-2400-C'
const LEG_2_NAME = 'ETH-20240329-2600-C'

// can retreive with public/get_instrument
const LEGS_TO_SUB_ID: any = {
  'ETH-20240329-2400-C': '39614082287924319838483674368',
  'ETH-20240329-2600-C': '39614082373823665758483674368'
}

async function signAuthenticationHeader() {
  const timestamp = Date.now().toString();
  const signature = await wallet.signMessage(timestamp);
  return {
      "X-LyraWallet": wallet.address,
      "X-LyraTimestamp": timestamp,
      "X-LyraSignature": signature
  };
}

// Schemas

type RfqLeg = {
  instrument_name: string,
  amount: string,
  direction: 'buy' | 'sell'
}

type QuoteLeg = {
  instrument_name: string,
  amount: string,
  direction: 'buy' | 'sell',
  price: string
}

type RfqResponse = {
  subaccount_id: number,
  creation_timestamp: number,
  last_update_timestamp: number,
  status: string,
  cancel_reason: string,
  rfq_id: string,
  valid_until: number,
  legs: Array<RfqLeg>
}

type QuoteResultPublicSchema = {
  cancel_reason: string;
  creation_timestamp: number;
  direction: 'buy' | 'sell';
  last_update_timestamp: number;
  legs: Array<QuoteLeg>;
  legs_hash: string;
  liquidity_role: 'maker' | 'taker';
  quote_id: string;
  rfq_id: string;
  status: string;
  subaccount_id: number;
  tx_hash: string | null;
  tx_status: string | null;
}

function createRfqObject(): object {
  const rfq = {
    subaccount_id: subaccount_id_rfq,
    // NOTE: legs MUST be sorted by instrument_name where sorting key is instrument_name
    legs: [
      {
        instrument_name: LEG_1_NAME,
        amount: '3',
        direction: 'buy'
      },
      {
        instrument_name: LEG_2_NAME,
        amount: '3',
        direction: 'sell'
      }
    ],
  };
  return rfq;
}

function createQuoteObject(rfq_response: RfqResponse, direction: 'buy' | 'sell') : object {
  const pricedLegs: Array<any> = rfq_response.legs;
  pricedLegs[0].price = direction == 'buy' ? '160' : '180';
  pricedLegs[1].price = direction == 'buy' ? '70' : '50';
  return {
    subaccount_id: subaccount_id_maker,
    rfq_id: rfq_response.rfq_id,
    legs: pricedLegs,
    direction: direction,
    max_fee: '10',
    nonce: Number(`${Date.now()}${Math.round(Math.random() * 999)}`),
    signer: wallet.address,
    signature_expiry_sec: Math.floor(Date.now() / 1000 + 350),
    signature: "filled_in_below"
  };
}

function createExecuteObject(quote: QuoteResultPublicSchema) : object {
  return {
    subaccount_id: subaccount_id_rfq,
    quote_id: quote.quote_id,
    rfq_id: quote.rfq_id,
    direction: quote.direction === 'buy' ? 'sell' : 'buy',
    max_fee: '10',
    nonce: Number(`${Date.now()}${Math.round(Math.random() * 999)}`),
    signer: wallet.address,
    signature_expiry_sec: Math.floor(Date.now() / 1000 + 350),
    legs: quote.legs,
    signature: "filled_in_below"
  }
}

// Getters / Polling

async function pollRfq() : Promise<RfqResponse> {
  // account owner of the subaccount_id must be approved to act as RFQ maker
  // can also use {wallet}.rfqs channel to listen for RFQs (same response format)
  const authHeader = await signAuthenticationHeader();
  const resp = await axios.post(`${HTTP_ADDRESS}/private/poll_rfqs`, {subaccount_id: subaccount_id_maker, status: 'open'}, {headers: authHeader})
  console.log(`found ${resp.data.result.rfqs.length} RFQs`)
  return resp.data.result.rfqs[0]
}

async function pollQuotes(rfq_id: string): Promise<Array<QuoteResultPublicSchema>> {
  const authHeader = await signAuthenticationHeader();
  const resp = await axios.post(`${HTTP_ADDRESS}/private/poll_quotes`, {subaccount_id: subaccount_id_rfq, rfq_id: rfq_id, status: 'open'}, {headers: authHeader})
  return resp.data.result.quotes;
}

// Signatures and Encoding

function signAction(action: any, actionData: string) {
  const action_hash = ethers.keccak256(
    encoder.encode(
      ['bytes32', 'uint256', 'uint256', 'address', 'bytes32', 'uint256', 'address', 'address'],
      [
        ACTION_TYPEHASH,
        action.subaccount_id,
        action.nonce,
        RFQ_ADDRESS,
        actionData,
        action.signature_expiry_sec,
        wallet.address,
        action.signer
      ]
    )
  );
  action.signature = wallet.signingKey.sign(
    ethers.keccak256(Buffer.concat([
      Buffer.from("1901", "hex"),
      Buffer.from(DOMAIN_SEPARATOR.slice(2), "hex"),
      Buffer.from(action_hash.slice(2), "hex")
    ]))
  ).serialized;
}

type EncodedLeg = [string, string, ethers.BigNumberish, ethers.BigNumberish]

function encodePricedLegs(legs: Array<QuoteLeg>, direction: 'buy' | 'sell'): Array<EncodedLeg> {
  const dirSign = BigInt(direction === 'buy' ? 1 : -1);
  const encoded_legs : Array<EncodedLeg> = legs.map((leg) => {
    const subid = LEGS_TO_SUB_ID[leg.instrument_name];
    const legSign = BigInt(leg.direction === 'buy' ? 1 : -1);
    const signedAmount = ethers.parseUnits(leg.amount, 18) * legSign * dirSign;
    return [OPTION_ADDRESS, subid, ethers.parseUnits(leg.price, 18), signedAmount];
  });
  return encoded_legs;
}

function encodeQuoteData(encoded_legs: Array<EncodedLeg>, max_fee: string): string {
  const rfqData = [ethers.parseUnits(max_fee, 18), encoded_legs];
  const QuoteDataABI = ['(uint,(address,uint,uint,int)[])'];
  const encodedData = encoder.encode(QuoteDataABI, [rfqData]);
  const hashedData = ethers.keccak256(Buffer.from(encodedData.slice(2), 'hex'));
  return hashedData;
}

function signQuote(quote: any) {
  const encoded_legs = encodePricedLegs(quote.legs, quote.direction);
  const quoteData = encodeQuoteData(encoded_legs, quote.max_fee);
  signAction(quote, quoteData)
}

function encodeExecuteData(encoded_legs: Array<EncodedLeg>, max_fee: string): string {
  const encoder = ethers.AbiCoder.defaultAbiCoder();
  const orderHashABI = ['(address,uint,uint,int)[]'];
  const orderHash = ethers.keccak256(Buffer.from(encoder.encode(orderHashABI, [encoded_legs]).slice(2), 'hex'));
  const ExectuteDataABI = ['bytes32', 'uint'];
  const encodedData = encoder.encode(ExectuteDataABI, [orderHash, ethers.parseUnits(max_fee, 18)]);
  const hashedData = ethers.keccak256(Buffer.from(encodedData.slice(2), 'hex'));
  return hashedData;
}

function signExecute(execute: any) {
  const encoded_legs = encodePricedLegs(execute.legs, execute.direction === 'buy' ? 'sell' : 'buy');
  const executeData = encodeExecuteData(encoded_legs, execute.max_fee);
  signAction(execute, executeData)
}

// Send API

async function sendRfq(rfq: object) {
  const authHeader = await signAuthenticationHeader();
  const resp = await axios.post(`${HTTP_ADDRESS}/private/send_rfq`, rfq, {headers: authHeader})
  return resp.data.result;
}

async function sendQuote(rfq_response: RfqResponse, direction: 'buy' | 'sell') {
  const quote = createQuoteObject(rfq_response, direction);
  signQuote(quote);
  const authHeader = await signAuthenticationHeader();
  const resp = await axios.post(`${HTTP_ADDRESS}/private/send_quote`, quote, {headers: authHeader})
  return resp.data.result;
}

async function sendExecute(quote: QuoteResultPublicSchema) {
  const execute = createExecuteObject(quote);
  signExecute(execute);
  const authHeader = await signAuthenticationHeader();
  const resp = await axios.post(`${HTTP_ADDRESS}/private/execute_quote`, execute, {headers: authHeader})
  return resp.data.result;
}

// Helpers to check if the RFQ is filled

async function getSubaccount(subaccount_id: number) {
  const resp = await axios.post(`${HTTP_ADDRESS}/private/get_subaccount`, {subaccount_id: subaccount_id}, {headers: await signAuthenticationHeader()})
  return resp.data.result;
}

async function getFilledQuotes() {
  const authHeader = await signAuthenticationHeader();
  const resp = await axios.post(`${HTTP_ADDRESS}/private/get_quotes`, {subaccount_id: subaccount_id_rfq, status: 'filled'}, {headers: authHeader})
  return resp.data.result;
}

async function completeRfq() {
    await sendRfq(createRfqObject())
    const rfq_response = await pollRfq();
    console.log(rfq_response);

    const buy_response = await sendQuote(rfq_response, 'buy');
    console.log(buy_response);

    const sell_response = await sendQuote(rfq_response, 'sell');
    console.log(sell_response);

    const quotes = await pollQuotes(rfq_response.rfq_id);
    console.log(quotes);

    const buyQuote = quotes.find((quote) => quote.direction === 'buy') as QuoteResultPublicSchema;
    console.log(buyQuote);

    const sellQuote = quotes.find((quote) => quote.direction === 'sell') as QuoteResultPublicSchema;
    console.log(sellQuote);

    const executeAsSeller = await sendExecute(buyQuote);
    console.log(executeAsSeller);

    // NOTE optionally execute as buyer instead, but only one side can execute
    // const executeAsBuyer = await sendExecute(sellQuote);
    // console.log(executeAsBuyer);

    console.log(await getSubaccount(subaccount_id_rfq));
    console.log(await getFilledQuotes());
}

completeRfq();

```

---

## rfq-quoting-and-execution-javascript-copy

**Title:** Rfq Quoting And Execution Javascript Copy
**URL:** https://docs.derive.xyz/reference/rfq-quoting-and-execution-javascript-copy

Similar to orderbook trading, RFQs are "self-custodial", and they require signed messages to be settled. Those signed messages guarantee that all legs of an RFQ will execute at the specified prices and amounts, as well as that the fee charged by the orderbook does not exceed the signed `max_fee`.

Unlike orderbook trading, makers and takers follow different rules and sign slightly different messages in order to complete an RFQ. The full flow is below:

1. **[Taker & Maker]** Authentication
2. **[Taker]** Send RFQ
3. **[Maker]** Listen or poll for RFQs
4. **[Maker]** In response to an RFQ, sign and send a quote
5. **[Taker]** Poll for the Quotes (market makers' replies to RFQs) and pick the best one
6. **[Taker]** Sign an execute message for the selected quote

The [Derive Python Action Signing SDK](https://pypi.org/project/derive_action_signing/) can be used to help with signing the self-custodial actions as part of the above flow.

For Taker actions (steps 1, 2, 5, 6) refer to the [RFQ execute example](https://github.com/derivexyz/v2-action-signing-python/blob/master/examples/rfq_execute.py)

For Maker actions (steps 1, 3, 4) refer to the [RFQ quote example](https://github.com/derivexyz/v2-action-signing-python/blob/master/examples/rfq_quote.py)

---

## submit-order

**Title:** Solidity Objects
**URL:** https://docs.derive.xyz/reference/submit-order

Order submission is a "self-custodial" request, a request that is guaranteed to not be alter-able by anyone except you, which means that it must past both:

1. orderbook authentication (steps 1-2)
2. on-chain signature verification (steps 3-4)

> ## 👍 If you are struggling to encode data correctly, you can use the public/order\_debug endpoint. The route takes in all raw inputs and returns intermediary outputs shown in the below steps.

## 1. Authenticate

The first step is to login via WebSocket - see the [Authentication](/reference/authentication) section for more:

TypeScript

```
async function signAuthenticationHeader(): Promise<{[key: string]: string}> {
  const timestamp = Date.now().toString();
  const signature = await wallet.signMessage(timestamp);  
    return {
      wallet: wallet.address,
      timestamp: timestamp,
      signature: signature,
    };
}

const connectWs = async (): Promise<WebSocket> => {
    return new Promise((resolve, reject) => {
        const ws = new WebSocket(WS_ADDRESS);

        ws.on('open', () => {
            setTimeout(() => resolve(ws), 50);
        });

        ws.on('error', reject);

        ws.on('close', (code: number, reason: Buffer) => {
            if (code && reason.toString()) {
                console.log(`WebSocket closed with code: ${code}`, `Reason: ${reason}`);
            }
        });
    });
};

async function loginClient(wsc: WebSocket) {
    const login_request = JSON.stringify({
        method: 'public/login',
        params: await signAuthenticationHeader(),
        id: Math.floor(Math.random() * 10000)
    });
    wsc.send(login_request);
    await new Promise(resolve => setTimeout(resolve, 2000));
}

```

## 2. Define

See the WebSocket API reference for `private/order` on more param documentation.

The same order can also be sent through REST, see the REST API reference for `private/order` for more info.

TypeScript

```
function defineOrder(): any {
    return {
        instrument_name: OPTION_NAME,
        subaccount_id: subaccount_id,
        direction: "buy",
        limit_price: 310,
        amount: 1,
        signature_expiry_sec: Math.floor(Date.now() / 1000 + 600), // must be >5min from now
        max_fee: "0.01",
        nonce: Number(`${Date.now()}${Math.round(Math.random() * 999)}`), // LYRA nonce format: ${CURRENT UTC MS +/- 1 day}${RANDOM 3 DIGIT NUMBER}
        signer: wallet.address,
        order_type: "limit",
        mmp: false,
        signature: "filled_in_below"
    };
}

```

## 3. Sign

When a fill occurs, this signature will be verified by the on-chain `Matching.sol` contract to ensure that you approved this trade.

TypeScript

```

function encodeTradeData(order: any): string {
  let encoded_data = encoder.encode( // same as "encoded_data" in public/order_debug
  
    ['address', 'uint', 'int', 'int', 'uint', 'uint', 'bool'],
    [
      ASSET_ADDRESS, 
      OPTION_SUB_ID, 
      ethers.parseUnits(order.limit_price.toString(), 18), 
      ethers.parseUnits(order.amount.toString(), 18), 
      ethers.parseUnits(order.max_fee.toString(), 18), 
      order.subaccount_id, order.direction === 'buy'
    ]
  );
  return ethers.keccak256(Buffer.from(encoded_data.slice(2), 'hex')) // same as "encoded_data_hashed" in public/order_debug
}

async function signOrder(order: any) {
    const tradeModuleData = encodeTradeData(order)

    const action_hash = ethers.keccak256(
        encoder.encode(
          ['bytes32', 'uint256', 'uint256', 'address', 'bytes32', 'uint256', 'address', 'address'], 
          [
            ACTION_TYPEHASH, 
            order.subaccount_id, 
            order.nonce, 
            TRADE_MODULE_ADDRESS, 
            tradeModuleData, 
            order.signature_expiry_sec, 
            wallet.address, 
            order.signer
          ]
        )
    ); // same as "action_hash" in public/order_debug

    order.signature = wallet.signingKey.sign(
        ethers.keccak256(Buffer.concat([
          Buffer.from("1901", "hex"), 
          Buffer.from(DOMAIN_SEPARATOR.slice(2), "hex"), 
          Buffer.from(action_hash.slice(2), "hex")
        ]))  // same as "typed_data_hash" in public/order_debug
    ).serialized;
}

```

## 4. Send

You will most likely have more involved listeners, but for example purposes a built-in listener is added into the submitOrder function.

TypeScript

```
async function submitOrder(order: any, ws: WebSocket) {
    return new Promise((resolve, reject) => {
        const id = Math.floor(Math.random() * 1000);
        ws.send(JSON.stringify({
            method: 'private/order',
            params: order,
            id: id
        }));

        ws.on('message', (message: string) => {
            const msg = JSON.parse(message);
            if (msg.id === id) {
                console.log('Got order response:', msg);
                resolve(msg);
            }
        });
    });
}

```

## Putting it all together

TypeScript

```
import { ethers } from "ethers";
import { WebSocket } from 'ws';
import dotenv from 'dotenv';

dotenv.config();

const PRIVATE_KEY = process.env.OWNER_PRIVATE_KEY as string;
const PROVIDER_URL = 'https://l2-prod-testnet-0eakp60405.t.conduit.xyz';
const WS_ADDRESS = 'wss://api-demo.lyra.finance/ws';
const ACTION_TYPEHASH = '0x4d7a9f27c403ff9c0f19bce61d76d82f9aa29f8d6d4b0c5474607d9770d1af17';
const DOMAIN_SEPARATOR = '0x9bcf4dc06df5d8bf23af818d5716491b995020f377d3b7b64c29ed14e3dd1105';
const ASSET_ADDRESS = '0xBcB494059969DAaB460E0B5d4f5c2366aab79aa1';
const TRADE_MODULE_ADDRESS = '0x87F2863866D85E3192a35A73b388BD625D83f2be';

const PROVIDER = new ethers.JsonRpcProvider(PROVIDER_URL);
const wallet = new ethers.Wallet(PRIVATE_KEY, PROVIDER);
const encoder = ethers.AbiCoder.defaultAbiCoder();
const subaccount_id = 9

const OPTION_NAME = 'ETH-20231027-1500-P'
const OPTION_SUB_ID = '644245094401698393600' // can retreive with public/get_instrument

async function signAuthenticationHeader(): Promise<{[key: string]: string}> {
  const timestamp = Date.now().toString();
  const signature = await wallet.signMessage(timestamp);  
    return {
      wallet: wallet.address,
      timestamp: timestamp,
      signature: signature,
    };
}

const connectWs = async (): Promise<WebSocket> => {
    return new Promise((resolve, reject) => {
        const ws = new WebSocket(WS_ADDRESS);

        ws.on('open', () => {
            setTimeout(() => resolve(ws), 50);
        });

        ws.on('error', reject);

        ws.on('close', (code: number, reason: Buffer) => {
            if (code && reason.toString()) {
                console.log(`WebSocket closed with code: ${code}`, `Reason: ${reason}`);
            }
        });
    });
};

async function loginClient(wsc: WebSocket) {
    const login_request = JSON.stringify({
        method: 'public/login',
        params: await signAuthenticationHeader(),
        id: Math.floor(Math.random() * 10000)
    });
    wsc.send(login_request);
    await new Promise(resolve => setTimeout(resolve, 2000));
}

function defineOrder(): any {
    return {
        instrument_name: OPTION_NAME,
        subaccount_id: subaccount_id,
        direction: "buy",
        limit_price: 310,
        amount: 1,
        signature_expiry_sec: Math.floor(Date.now() / 1000 + 600), // must be >5min from now
        max_fee: "0.01",
        nonce: Number(`${Date.now()}${Math.round(Math.random() * 999)}`), // LYRA nonce format: ${CURRENT UTC MS +/- 1 day}${RANDOM 3 DIGIT NUMBER}
        signer: wallet.address,
        order_type: "limit",
        mmp: false,
        signature: "filled_in_below"
    };
}

function encodeTradeData(order: any): string {
  let encoded_data = encoder.encode( // same as "encoded_data" in public/order_debug
    ['address', 'uint', 'int', 'int', 'uint', 'uint', 'bool'],
    [
      ASSET_ADDRESS, 
      OPTION_SUB_ID, 
      ethers.parseUnits(order.limit_price.toString(), 18), 
      ethers.parseUnits(order.amount.toString(), 18), 
      ethers.parseUnits(order.max_fee.toString(), 18), 
      order.subaccount_id, order.direction === 'buy']
    );
  return ethers.keccak256(Buffer.from(encoded_data.slice(2), 'hex')) // same as "encoded_data_hashed" in public/order_debug
}

async function signOrder(order: any) {
    const tradeModuleData = encodeTradeData(order)

    const action_hash = ethers.keccak256(
        encoder.encode(
          ['bytes32', 'uint256', 'uint256', 'address', 'bytes32', 'uint256', 'address', 'address'], 
          [
            ACTION_TYPEHASH, 
            order.subaccount_id, 
            order.nonce, 
            TRADE_MODULE_ADDRESS, 
            tradeModuleData, 
            order.signature_expiry_sec, 
            wallet.address, 
            order.signer
          ]
        )
    ); // same as "action_hash" in public/order_debug

    order.signature = wallet.signingKey.sign(
        ethers.keccak256(Buffer.concat([
          Buffer.from("1901", "hex"), 
          Buffer.from(DOMAIN_SEPARATOR.slice(2), "hex"), 
          Buffer.from(action_hash.slice(2), "hex")
        ]))  // same as "typed_data_hash" in public/order_debug
    ).serialized;
}

async function submitOrder(order: any, ws: WebSocket) {
    return new Promise((resolve, reject) => {
        const id = Math.floor(Math.random() * 1000);
        ws.send(JSON.stringify({
            method: 'private/order',
            params: order,
            id: id
        }));

        ws.on('message', (message: string) => {
            const msg = JSON.parse(message);
            if (msg.id === id) {
                console.log('Got order response:', msg);
                resolve(msg);
            }
        });
    });
}

async function completeOrder() {
    const ws = await connectWs();
    await loginClient(ws);
    const order = defineOrder();
    await signOrder(order);
    await submitOrder(order, ws);
}

completeOrder();

```

# Solidity Objects

### SignedAction Schema

| Param | Type | Description |
| --- | --- | --- |
| `subaccount_id` | `uint` | User subaccount id for the action (0 for a new subaccounts when depositing) |
| `nonce` | `uint` | Unique nonce defined as <UTC\_timestamp in ms><random\_number\_up\_to\_6\_digits> (e.g. 1695836058725001, where 001 is the random number) |
| `module` | `address` | Deposit module address (see [Protocol Constants](/reference/protocol-constants)) |
| `data` | `bytes` | Encoded module data ("TradeModuleData" for orders) |
| `expiry` | `uint` | Signature expiry timestamp in sec |
| `owner` | `address` | Wallet address of the account owner |
| `signer` | `address` | Either owner wallet or session key |

### TradeModuleData Schema

| Param | Type | Description |
| --- | --- | --- |
| `asset` | `address` | Get with `public/get_instrument` (base\_asset\_address) |
| `subId` | `uint` | Sub ID of the asset (Get from public/get\_instrument endpoint) |
| `amount` | `int` | Max amount willing to trade |
| `max_fee` | `uint` | max fee |
| `recipient_id` | `uint` | User subaccount id |
| `isBid` | `bool` | Bid or Ask |

---

## submit-order-javascript-copy

**Title:** Submit Order Javascript Copy
**URL:** https://docs.derive.xyz/reference/submit-order-javascript-copy

Order submission is a "self-custodial" request, a request that is guaranteed to not be alter-able by anyone except you, which means that it must past both:

1. orderbook authentication
2. on-chain signature verification

The [Derive Python Action Signing SDK](https://pypi.org/project/derive_action_signing/) can be used perform both steps. You may use it for orders as well as other self-custodial requests (e.g. deposits, withdrawals, and etc).

Refer to the [full order example](https://github.com/derivexyz/v2-action-signing-python/blob/master/examples/order.py) in the SDK repo.

---

## transfer

**Title:** Transfer
**URL:** https://docs.derive.xyz/reference/transfer

Click on "Account Settings" and choose "Subaccounts" to see a list of your subacounts.

![](https://files.readme.io/f52c3c8-image.png)

Click on "Transfer" to transfer funds:

![](https://files.readme.io/7899c3a-image.png)

---

## transfer-collateral

**Title:** Transfer Collateral
**URL:** https://docs.derive.xyz/reference/transfer-collateral

Note there are 3 primary methods of transferring assets on Derive:

1. Transfer Collateral via `private/transfer_ec20` (e.g. USDC, ETH, BTC and etc): see [Transfer Collateral Example](https://github.com/derivexyz/v2-action-signing-python/blob/master/examples/transfer_erc20.py)
2. Transfer a single position via `private/transfer_position` (e.g. ETH-PERP or ETH option): example is WIP but same effect can be achieved by #3
3. Transfer multiple positions via `private/transfer_positions`: e.g. ETH\_PERP + BTC option: see [Transfer Multiple Positions Example](https://github.com/derivexyz/v2-action-signing-python/blob/master/examples/transfer_positions.py)

The above examples use the [Derive Python Action Signing SDK](https://pypi.org/project/derive_action_signing/) to greatly simplify signing the self-custodial actions and authentication.

---

## ux-create-or-deposit-to-subaccount

**Title:** Ux Create Or Deposit To Subaccount
**URL:** https://docs.derive.xyz/reference/ux-create-or-deposit-to-subaccount

This is the preferred / no-code method for integrating with the Derive exchange.

You may use the user interface to:

- Bridge from Mainnet / OP / Arbitrum
- Deposit / withdraw funds into / from the exchange
- Transfer funds between subaccounts
- Create several subaccounts with different margin types
- Monitor and manage positions / open orders via UX
- Manage session keys

If you'd like to complete these steps fully on-chain refer to [Onboard Manually](/docs/onboard-manually).

## Step 1: Connect Wallet

Load the [www.derive.xyz](http://www.derive.xyz) website and follow the "connect wallet" flow:

![](https://files.readme.io/2c3867c-image.png)

You may use "Metamask" if you'd like to onboard via a hardware wallet.

## Step 2: Launch "Getting Started" Flow

Enter the "Developer" page by clicking on the "Account Settings" drop down on the top right of the page. You should see the below page.

Follow the flow to create a subaccount, mint USDC and create your first session key.

![](https://files.readme.io/7f4bbcb-image.png)

Refer to the other guides in [Onboard via Interface](/docs/onboard-via-interface) section for other actions.

## Smart-contract Wallets

When onboarding via the UX, Derive creates a smart-contract wallet wrapper around your original Ethereum Wallet. Your wallet still has full control over all actions, however the all funds are owned by the smart contract wallet.

This means when you view transactions on etherscan, transfers / fills / deposits will all appear to happen to this Smart-contract wallet address.

![](https://files.readme.io/d62c1df-image.png)

You can go to the "Account Settings" dropdown and click on "Account" to see:

- Owner Address: Original Ethereum wallet used to create the account
- Wallet Address: Smart-contract wallet (still fully controlled by the original Ethereum wallet)

---

## ux-withdraw

**Title:** Ux Withdraw
**URL:** https://docs.derive.xyz/reference/ux-withdraw

Choose the "Withdraw" button in the "Account Settings" dropdown:

![](https://files.readme.io/8e128cc-image.png)

You may choose to withdraw to Ethereum, Optimism or Arbitrum.

---
