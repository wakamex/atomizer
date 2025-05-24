package main

import (
	"fmt"
	"math"
)

// SSVIParameters represents the Surface SVI parameterization
// SSVI ensures calendar arbitrage-free surfaces by construction
type SSVIParameters struct {
	Theta float64 // ATM variance term structure parameter (θ)
	Rho   float64 // Correlation parameter (ρ)
	Phi   float64 // Vol of vol parameter (φ)
}

// SSVISurface represents a volatility surface using SSVI parameterization
type SSVISurface struct {
	Parameters   SSVIParameters
	SpotPrice    float64
	RiskFreeRate float64
}

// ComputeTotalVariance calculates total implied variance using SSVI formula
// w(k,t) = θt/2 * [1 + ρφk + √((φk + ρ)² + (1 - ρ²))]
func (s *SSVISurface) ComputeTotalVariance(logMoneyness float64, timeToExpiry float64) float64 {
	k := logMoneyness
	t := timeToExpiry
	theta := s.Parameters.Theta
	rho := s.Parameters.Rho
	phi := s.Parameters.Phi
	
	// SSVI formula: w(k,t) = θt/2 * [1 + ρφk + √((φk + ρ)² + (1 - ρ²))]
	// Note: The sqrt term should be: sqrt((phi*k + rho)^2 + (1 - rho^2))
	term1 := phi*k + rho
	term2 := math.Sqrt(term1*term1 + (1 - rho*rho))
	
	totalVar := (theta * t / 2.0) * (1 + rho*phi*k + term2)
	
	// Debug - log some values to understand the issue
	if k > -0.1 && k < 0.1 && t > 0.09 && t < 0.11 {
		fmt.Printf("SSVI Debug: k=%.4f, t=%.4f, theta=%.4f, rho=%.4f, phi=%.4f -> term2=%.4f, totalVar=%.4f\n", 
			k, t, theta, rho, phi, term2, totalVar)
	}
	
	return totalVar
}

// GetIV returns implied volatility for a given strike and time to expiry
func (s *SSVISurface) GetIV(strike, timeToExpiry float64) float64 {
	if timeToExpiry <= 0 {
		return 0
	}
	
	// Calculate forward price and log-moneyness
	forward := s.SpotPrice * math.Exp(s.RiskFreeRate * timeToExpiry)
	logMoneyness := math.Log(strike / forward)
	
	// Get total variance
	totalVar := s.ComputeTotalVariance(logMoneyness, timeToExpiry)
	
	// Convert to IV
	if totalVar <= 0 {
		return 0
	}
	
	return math.Sqrt(totalVar / timeToExpiry)
}

// FitSSVISurface fits SSVI model to all option data simultaneously
func FitSSVISurface(options []OptionData, spotPrice, riskFreeRate float64) (*SSVISurface, error) {
	// Filter and prepare data
	validOptions := []OptionData{}
	for _, opt := range options {
		if !math.IsNaN(opt.ImpliedVolatility) && opt.ImpliedVolatility > 0 && 
		   opt.TimeToExpiry > 0 && opt.Moneyness > 0.5 && opt.Moneyness < 2.0 {
			validOptions = append(validOptions, opt)
		}
	}
	
	if len(validOptions) < 10 {
		return nil, fmt.Errorf("insufficient valid options for SSVI fitting: %d", len(validOptions))
	}
	
	// Initial parameter guess based on ATM volatility
	atmVol := estimateATMVol(validOptions)
	
	// Analyze term structure to set better initial parameters
	shortTermVol, longTermVol := analyzeTermStructure(validOptions)
	
	// Calculate term structure characteristics
	// volRatio > 1 means inverted (backwardation), < 1 means normal (contango)
	volRatio := shortTermVol / longTermVol
	termStructureAdjustment := 1.0 + 0.2*math.Abs(volRatio-1.0) // Boost params when term structure is steep
	
	// Adjust initial parameters based on market observations
	params := SSVIParameters{
		Theta: atmVol * atmVol * termStructureAdjustment,  // Higher theta when term structure is steep
		Rho:   -0.4 - 0.1*math.Max(0, volRatio-1.0),      // Stronger skew in backwardation
		Phi:   1.0 + 0.2*math.Abs(volRatio-1.0),          // More smile when term structure is steep
	}
	
	fmt.Printf("SSVI: ATM vol: %.2f%%, Short-term vol: %.2f%%, Long-term vol: %.2f%%\n", 
		atmVol*100, shortTermVol*100, longTermVol*100)
	fmt.Printf("SSVI: Term structure ratio: %.3f, Initial params: Theta=%.4f, Rho=%.4f, Phi=%.4f\n", 
		volRatio, params.Theta, params.Rho, params.Phi)
	
	// Optimization using gradient descent with constraints
	// Use adaptive learning rate based on parameter sensitivity
	learningRate := 0.005  // More conservative for stability
	maxIterations := 8000  // More iterations for better convergence
	tolerance := 1e-7      // Tighter tolerance
	prevError := math.MaxFloat64
	
	fmt.Printf("SSVI Initial params: Theta=%.4f, Rho=%.4f, Phi=%.4f\n", params.Theta, params.Rho, params.Phi)
	
	for iter := 0; iter < maxIterations; iter++ {
		// Calculate current error and gradients
		currentError, gradTheta, gradRho, gradPhi := calculateSSVIGradients(
			params, validOptions, spotPrice, riskFreeRate)
		
		// Log progress periodically
		if iter%500 == 0 {
			fmt.Printf("SSVI Iter %d: Error=%.6f, Theta=%.4f, Rho=%.4f, Phi=%.4f\n", 
				iter, currentError, params.Theta, params.Rho, params.Phi)
		}
		
		// Update parameters with gradient descent
		params.Theta -= learningRate * gradTheta
		params.Rho -= learningRate * gradRho
		params.Phi -= learningRate * gradPhi
		
		// Apply constraints (prevent Theta from going too low)
		params.Theta = math.Max(atmVol*atmVol*0.8, math.Min(10.0, params.Theta)) // Prevent too low theta
		params.Rho = math.Max(-0.95, math.Min(0.95, params.Rho))
		params.Phi = math.Max(0.8, math.Min(2.0, params.Phi)) // Higher minimum phi
		
		// Check convergence based on change in error
		if iter > 0 && math.Abs(currentError-prevError) < tolerance {
			fmt.Printf("SSVI converged at iter %d\n", iter)
			break
		}
		prevError = currentError
		
		// Adaptive learning rate
		if iter%100 == 0 && iter > 0 {
			learningRate *= 0.98
		}
	}
	
	// Verify no-arbitrage conditions
	if !checkSSVINoArbitrage(params) {
		// Apply corrections to ensure no-arbitrage
		params = enforceSSVINoArbitrage(params)
	}
	
	surface := &SSVISurface{
		Parameters:   params,
		SpotPrice:    spotPrice,
		RiskFreeRate: riskFreeRate,
	}
	
	return surface, nil
}

// calculateSSVIGradients computes error and gradients for SSVI parameters
func calculateSSVIGradients(params SSVIParameters, options []OptionData, 
	spotPrice, riskFreeRate float64) (error, gradTheta, gradRho, gradPhi float64) {
	
	delta := 0.0001
	
	// Calculate base error
	error = calculateSSVIError(params, options, spotPrice, riskFreeRate)
	
	// Numerical gradients
	paramsTheta := params
	paramsTheta.Theta += delta
	errorTheta := calculateSSVIError(paramsTheta, options, spotPrice, riskFreeRate)
	gradTheta = (errorTheta - error) / delta
	
	paramsRho := params
	paramsRho.Rho += delta
	errorRho := calculateSSVIError(paramsRho, options, spotPrice, riskFreeRate)
	gradRho = (errorRho - error) / delta
	
	paramsPhi := params
	paramsPhi.Phi += delta
	errorPhi := calculateSSVIError(paramsPhi, options, spotPrice, riskFreeRate)
	gradPhi = (errorPhi - error) / delta
	
	return error, gradTheta, gradRho, gradPhi
}

// calculateSSVIError computes the fitting error for SSVI parameters
func calculateSSVIError(params SSVIParameters, options []OptionData, 
	spotPrice, riskFreeRate float64) float64 {
	
	sse := 0.0
	surface := &SSVISurface{
		Parameters:   params,
		SpotPrice:    spotPrice,
		RiskFreeRate: riskFreeRate,
	}
	
	for _, opt := range options {
		forward := spotPrice * math.Exp(riskFreeRate * opt.TimeToExpiry)
		logMoneyness := math.Log(opt.Strike / forward)
		
		// Market IV
		marketIV := opt.ImpliedVolatility
		
		// Model IV
		modelVar := surface.ComputeTotalVariance(logMoneyness, opt.TimeToExpiry)
		modelIV := math.Sqrt(modelVar / opt.TimeToExpiry)
		
		// Use IV difference instead of variance difference for more stable fitting
		weight := 1.0
		if !math.IsNaN(opt.Vega) && opt.Vega > 0 {
			weight = 1.0 / math.Sqrt(opt.Vega) // Weight by inverse of vega
		}
		
		diff := (marketIV - modelIV) * weight
		sse += diff * diff
	}
	
	// Add regularization to prevent extreme parameters
	reg := 0.0
	
	// Penalize extreme phi values
	if params.Phi > 1.5 {
		reg += 100 * (params.Phi - 1.5) * (params.Phi - 1.5)
	}
	
	// Penalize extreme correlation
	if math.Abs(params.Rho) > 0.8 {
		reg += 100 * (math.Abs(params.Rho) - 0.8) * (math.Abs(params.Rho) - 0.8)
	}
	
	return sse + reg
}

// estimateATMVol estimates at-the-money volatility from options data
func estimateATMVol(options []OptionData) float64 {
	atmVols := []float64{}
	
	for _, opt := range options {
		// Find near-ATM options
		if opt.Moneyness > 0.95 && opt.Moneyness < 1.05 {
			atmVols = append(atmVols, opt.ImpliedVolatility)
		}
	}
	
	if len(atmVols) == 0 {
		// Fallback to all options average
		sum := 0.0
		for _, opt := range options {
			sum += opt.ImpliedVolatility
		}
		return sum / float64(len(options))
	}
	
	// Return median ATM vol
	return median(atmVols)
}

// analyzeTermStructure analyzes short-term vs long-term volatility patterns
func analyzeTermStructure(options []OptionData) (shortTermVol, longTermVol float64) {
	shortTerm := []float64{}
	longTerm := []float64{}
	
	for _, opt := range options {
		// Focus on near-ATM options for term structure
		if opt.Moneyness > 0.85 && opt.Moneyness < 1.15 {
			if opt.TimeToExpiry < 0.1 { // Less than ~36 days
				shortTerm = append(shortTerm, opt.ImpliedVolatility)
			} else if opt.TimeToExpiry > 0.3 { // More than ~110 days
				longTerm = append(longTerm, opt.ImpliedVolatility)
			}
		}
	}
	
	// Calculate averages with market-based defaults
	shortTermVol = 0.7 // Default
	if len(shortTerm) > 0 {
		shortTermVol = median(shortTerm)
	}
	
	longTermVol = 0.65 // Default
	if len(longTerm) > 0 {
		longTermVol = median(longTerm)
	}
	
	// Don't force any relationship - let the data speak
	// Term structure can be normal (short < long) or inverted (short > long)
	
	return shortTermVol, longTermVol
}

// checkSSVINoArbitrage verifies SSVI parameters satisfy no-arbitrage conditions
func checkSSVINoArbitrage(params SSVIParameters) bool {
	theta := params.Theta
	rho := params.Rho
	phi := params.Phi
	
	// Basic parameter bounds
	if theta <= 0 || math.Abs(rho) >= 1 || phi <= 0 {
		return false
	}
	
	// Check Gatheral's conditions for SSVI
	// Condition 1: φ(1 + |ρ|) ≤ 2
	if phi*(1+math.Abs(rho)) > 2 {
		return false
	}
	
	// Condition 2: For calendar spread arbitrage-free
	// The function g(k) = θt/2 * [1 + ρφk + √((φk + ρ)² + (1 - ρ²))]
	// must be increasing in t for all k
	// This is automatically satisfied by construction in SSVI
	
	return true
}

// enforceSSVINoArbitrage adjusts parameters to satisfy no-arbitrage conditions
func enforceSSVINoArbitrage(params SSVIParameters) SSVIParameters {
	adjusted := params
	
	// Ensure basic bounds
	adjusted.Theta = math.Max(0.001, adjusted.Theta)
	adjusted.Rho = math.Max(-0.99, math.Min(0.99, adjusted.Rho))
	adjusted.Phi = math.Max(0.01, adjusted.Phi)
	
	// Enforce Gatheral's condition: φ(1 + |ρ|) ≤ 2
	maxPhi := 2.0 / (1 + math.Abs(adjusted.Rho))
	if adjusted.Phi > maxPhi {
		adjusted.Phi = maxPhi * 0.95 // Stay slightly below the bound
	}
	
	return adjusted
}

// median calculates the median of a slice of floats
func median(values []float64) float64 {
	n := len(values)
	if n == 0 {
		return 0
	}
	
	// Sort values
	sorted := make([]float64, n)
	copy(sorted, values)
	for i := 0; i < n; i++ {
		for j := i + 1; j < n; j++ {
			if sorted[i] > sorted[j] {
				sorted[i], sorted[j] = sorted[j], sorted[i]
			}
		}
	}
	
	if n%2 == 0 {
		return (sorted[n/2-1] + sorted[n/2]) / 2
	}
	return sorted[n/2]
}