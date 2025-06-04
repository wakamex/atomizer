# Deribit Integration for Options Analysis Tool

The options analysis tool now supports fetching and analyzing options from both Derive/Lyra and Deribit exchanges.

## Usage

### Command Line

Use the `--exchange=deribit` flag to analyze Deribit options:

```bash
# Analyze Deribit options
./analyze_options.sh query --exchange=deribit

# Query second-nearest expiry for Deribit
./analyze_options.sh query 2 --exchange=deribit

# Export Deribit ETH calls
./analyze_options.sh export --exchange=deribit
```

### Environment Variables

- `DERIBIT_TEST_MODE=true` - Use Deribit testnet instead of mainnet
- `ANALYZE_BTC=true` - Analyze BTC options instead of ETH (Deribit only)

### Examples

```bash
# Analyze BTC options on Deribit
ANALYZE_BTC=true ./analyze_options.sh query --exchange=deribit

# Use Deribit testnet
DERIBIT_TEST_MODE=true ./analyze_options.sh stats --exchange=deribit

# Query second-nearest BTC expiry on testnet
DERIBIT_TEST_MODE=true ANALYZE_BTC=true ./analyze_options.sh query 2 --exchange=deribit
```

## Features

- Fetches all active options from Deribit (BTC and ETH)
- Converts Deribit instrument format to common format
- Fetches real-time ticker data including Greeks
- Supports both mainnet and testnet
- Async ticker fetching for improved performance

## Data Mapping

| Deribit Field | Common Field |
|---------------|--------------|
| instrument_name | InstrumentName |
| base_currency | BaseCurrency |
| strike | OptionDetails.Strike |
| option_type | OptionDetails.OptionType |
| expiration_timestamp | OptionDetails.Expiry |
| best_bid_price | BestBidPrice |
| best_ask_price | BestAskPrice |
| mark_price | MarkPrice |
| index_price | IndexPrice |
| greeks.delta | OptionPricing.Delta |
| greeks.gamma | OptionPricing.Gamma |
| greeks.theta | OptionPricing.Theta |
| greeks.vega | OptionPricing.Vega |
| stats.volume | Stats.ContractVolume |

## Notes

- Deribit options expire at 08:00 UTC
- Instrument format: `BTC-27DEC24-100000-C` (Currency-Expiry-Strike-Type)
- No authentication required for public market data
- Rate limiting is handled automatically with concurrent request limits