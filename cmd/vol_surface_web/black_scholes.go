package main

import (
	"math"
)

// BlackScholesParams holds the parameters for Black-Scholes calculation
type BlackScholesParams struct {
	SpotPrice    float64 // Current underlying price
	StrikePrice  float64 // Option strike price
	TimeToExpiry float64 // Time to expiry in years
	RiskFreeRate float64 // Risk-free interest rate
	Volatility   float64 // Implied volatility (annualized)
	IsCall       bool    // true for call, false for put
}

// BlackScholesPrice calculates the theoretical option price using Black-Scholes
func BlackScholesPrice(params BlackScholesParams) float64 {
	S := params.SpotPrice
	K := params.StrikePrice
	T := params.TimeToExpiry
	r := params.RiskFreeRate
	σ := params.Volatility

	// Handle edge cases
	if T <= 0 {
		if params.IsCall {
			return math.Max(S-K, 0)
		} else {
			return math.Max(K-S, 0)
		}
	}

	if σ <= 0 {
		return 0
	}

	// Calculate d1 and d2
	d1 := (math.Log(S/K) + (r+0.5*σ*σ)*T) / (σ * math.Sqrt(T))
	d2 := d1 - σ*math.Sqrt(T)

	// Standard normal CDF
	Nd1 := normalCDF(d1)
	Nd2 := normalCDF(d2)
	NminusD1 := normalCDF(-d1)
	NminusD2 := normalCDF(-d2)

	if params.IsCall {
		// Call option price
		return S*Nd1 - K*math.Exp(-r*T)*Nd2
	} else {
		// Put option price
		return K*math.Exp(-r*T)*NminusD2 - S*NminusD1
	}
}

// normalCDF approximates the cumulative distribution function of standard normal distribution
func normalCDF(x float64) float64 {
	// Abramowitz and Stegun approximation
	// Error < 7.5e-8
	a1 := 0.254829592
	a2 := -0.284496736
	a3 := 1.421413741
	a4 := -1.453152027
	a5 := 1.061405429
	p := 0.3275911

	sign := 1.0
	if x < 0 {
		sign = -1.0
	}
	x = math.Abs(x) / math.Sqrt2

	// A&S formula 7.1.26
	t := 1.0 / (1.0 + p*x)
	y := 1.0 - (((((a5*t+a4)*t)+a3)*t+a2)*t+a1)*t*math.Exp(-x*x)

	return 0.5 * (1.0 + sign*y)
}

// ImpliedVolatility calculates implied volatility using Newton-Raphson method
func ImpliedVolatility(marketPrice float64, params BlackScholesParams) float64 {
	const (
		maxIterations = 100
		tolerance     = 1e-6
		minVol        = 0.001  // 0.1%
		maxVol        = 5.0    // 500%
	)

	// Initial guess
	σ := 0.3 // 30% volatility

	for i := 0; i < maxIterations; i++ {
		params.Volatility = σ
		
		// Calculate theoretical price and vega
		theoreticalPrice := BlackScholesPrice(params)
		vega := BlackScholesVega(params)
		
		// Check convergence
		priceDiff := theoreticalPrice - marketPrice
		if math.Abs(priceDiff) < tolerance {
			return σ
		}
		
		// Newton-Raphson update
		if vega > 0 {
			σ = σ - priceDiff/vega
		} else {
			break // Avoid division by zero
		}
		
		// Constrain volatility to reasonable bounds
		if σ < minVol {
			σ = minVol
		}
		if σ > maxVol {
			σ = maxVol
		}
	}
	
	// Return NaN if convergence failed
	return math.NaN()
}

// BlackScholesVega calculates the vega (sensitivity to volatility) of an option
func BlackScholesVega(params BlackScholesParams) float64 {
	S := params.SpotPrice
	K := params.StrikePrice
	T := params.TimeToExpiry
	r := params.RiskFreeRate
	σ := params.Volatility

	if T <= 0 || σ <= 0 {
		return 0
	}

	d1 := (math.Log(S/K) + (r+0.5*σ*σ)*T) / (σ * math.Sqrt(T))
	
	// Vega is the same for calls and puts
	return S * math.Sqrt(T) * normalPDF(d1)
}

// normalPDF calculates the probability density function of standard normal distribution
func normalPDF(x float64) float64 {
	return math.Exp(-0.5*x*x) / math.Sqrt(2*math.Pi)
}

// BlackScholesDelta calculates the delta (sensitivity to underlying price) of an option
func BlackScholesDelta(params BlackScholesParams) float64 {
	S := params.SpotPrice
	K := params.StrikePrice
	T := params.TimeToExpiry
	r := params.RiskFreeRate
	σ := params.Volatility

	if T <= 0 {
		if params.IsCall {
			if S > K {
				return 1.0
			} else {
				return 0.0
			}
		} else {
			if S < K {
				return -1.0
			} else {
				return 0.0
			}
		}
	}

	d1 := (math.Log(S/K) + (r+0.5*σ*σ)*T) / (σ * math.Sqrt(T))

	if params.IsCall {
		return normalCDF(d1)
	} else {
		return normalCDF(d1) - 1.0
	}
}

// BlackScholesGamma calculates the gamma (second derivative w.r.t. underlying price) of an option
func BlackScholesGamma(params BlackScholesParams) float64 {
	S := params.SpotPrice
	K := params.StrikePrice
	T := params.TimeToExpiry
	r := params.RiskFreeRate
	σ := params.Volatility

	if T <= 0 || σ <= 0 {
		return 0
	}

	d1 := (math.Log(S/K) + (r+0.5*σ*σ)*T) / (σ * math.Sqrt(T))
	
	// Gamma is the same for calls and puts
	return normalPDF(d1) / (S * σ * math.Sqrt(T))
}

// BlackScholesTheta calculates the theta (time decay) of an option
func BlackScholesTheta(params BlackScholesParams) float64 {
	S := params.SpotPrice
	K := params.StrikePrice
	T := params.TimeToExpiry
	r := params.RiskFreeRate
	σ := params.Volatility

	if T <= 0 {
		return 0
	}

	d1 := (math.Log(S/K) + (r+0.5*σ*σ)*T) / (σ * math.Sqrt(T))
	d2 := d1 - σ*math.Sqrt(T)

	term1 := -(S * normalPDF(d1) * σ) / (2 * math.Sqrt(T))
	
	if params.IsCall {
		term2 := r * K * math.Exp(-r*T) * normalCDF(d2)
		return term1 - term2
	} else {
		term2 := r * K * math.Exp(-r*T) * normalCDF(-d2)
		return term1 + term2
	}
}