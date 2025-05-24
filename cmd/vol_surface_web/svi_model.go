package main

import (
	"fmt"
	"math"
)

// SVIParameters represents parameters for one expiry slice in the SVI model
type SVIParameters struct {
	A     float64 // ATM variance level
	B     float64 // Slope parameter (controls the angle of the wings)
	Rho   float64 // Rotation parameter (-1 < rho < 1)
	M     float64 // Translation parameter
	Sigma float64 // Curvature parameter (sigma > 0)
	TTM   float64 // Time to maturity for this slice
}

// SVISurface represents a volatility surface using SVI parameterization
type SVISurface struct {
	Expiries   []float64        // Sorted list of expiries
	Parameters []SVIParameters  // SVI parameters for each expiry
	SpotPrice  float64          // Current spot price
	RiskFreeRate float64        // Risk-free rate
}

// computeTotalVariance calculates total implied variance using SVI formula
func (params SVIParameters) computeTotalVariance(logMoneyness float64) float64 {
	k := logMoneyness
	return params.A + params.B*(params.Rho*(k-params.M) + 
		math.Sqrt((k-params.M)*(k-params.M) + params.Sigma*params.Sigma))
}

// computeIV converts total variance to implied volatility
func (params SVIParameters) computeIV(logMoneyness float64) float64 {
	totalVar := params.computeTotalVariance(logMoneyness)
	if totalVar < 0 {
		return 0
	}
	return math.Sqrt(totalVar / params.TTM)
}

// FitSVISlice fits SVI model to one expiry slice with regularization and constraints
func FitSVISlice(strikes []float64, ivs []float64, forward float64, ttm float64) (SVIParameters, error) {
	if len(strikes) != len(ivs) {
		return SVIParameters{}, fmt.Errorf("strikes and ivs must have same length")
	}
	if len(strikes) < 5 {
		return SVIParameters{}, fmt.Errorf("need at least 5 points to fit SVI")
	}

	// Convert to total variance and log-moneyness
	totalVar := make([]float64, len(ivs))
	logMoneyness := make([]float64, len(strikes))
	weights := make([]float64, len(ivs))
	
	for i := range ivs {
		totalVar[i] = ivs[i] * ivs[i] * ttm
		logMoneyness[i] = math.Log(strikes[i] / forward)
		// Weight by inverse of moneyness distance from ATM
		weights[i] = 1.0 / (1.0 + math.Abs(logMoneyness[i]))
	}
	
	// Find ATM index and calculate smart initial guess
	atmIdx := 0
	minDist := math.Abs(logMoneyness[0])
	for i := 1; i < len(logMoneyness); i++ {
		dist := math.Abs(logMoneyness[i])
		if dist < minDist {
			minDist = dist
			atmIdx = i
		}
	}
	
	// Smart initial parameter guess
	atmVar := totalVar[atmIdx]
	
	// Estimate initial B from the slope of variance vs log-moneyness
	leftIdx := atmIdx - 1
	rightIdx := atmIdx + 1
	if leftIdx < 0 {
		leftIdx = 0
	}
	if rightIdx >= len(totalVar) {
		rightIdx = len(totalVar) - 1
	}
	
	initialB := 0.1
	if rightIdx > leftIdx {
		slope := (totalVar[rightIdx] - totalVar[leftIdx]) / (logMoneyness[rightIdx] - logMoneyness[leftIdx])
		initialB = math.Max(0.05, math.Min(0.5, math.Abs(slope)))
	}
	
	params := SVIParameters{
		A:     atmVar * 0.9,
		B:     initialB,
		Rho:   -0.3,
		M:     logMoneyness[atmIdx], // Center at ATM
		Sigma: 0.3,
		TTM:   ttm,
	}
	
	// Optimization with regularization
	learningRate := 0.01
	maxIterations := 2000
	tolerance := 1e-6
	lambda := 0.01 // Reduced regularization parameter
	
	for iter := 0; iter < maxIterations; iter++ {
		// Calculate current error with regularization
		currentError := calculateRegularizedSSE(params, logMoneyness, totalVar, weights, lambda)
		
		// Calculate gradients numerically
		delta := 0.0001
		gradA := (calculateRegularizedSSE(SVIParameters{A: params.A + delta, B: params.B, Rho: params.Rho, M: params.M, Sigma: params.Sigma, TTM: ttm}, logMoneyness, totalVar, weights, lambda) - currentError) / delta
		gradB := (calculateRegularizedSSE(SVIParameters{A: params.A, B: params.B + delta, Rho: params.Rho, M: params.M, Sigma: params.Sigma, TTM: ttm}, logMoneyness, totalVar, weights, lambda) - currentError) / delta
		gradRho := (calculateRegularizedSSE(SVIParameters{A: params.A, B: params.B, Rho: params.Rho + delta, M: params.M, Sigma: params.Sigma, TTM: ttm}, logMoneyness, totalVar, weights, lambda) - currentError) / delta
		gradM := (calculateRegularizedSSE(SVIParameters{A: params.A, B: params.B, Rho: params.Rho, M: params.M + delta, Sigma: params.Sigma, TTM: ttm}, logMoneyness, totalVar, weights, lambda) - currentError) / delta
		gradSigma := (calculateRegularizedSSE(SVIParameters{A: params.A, B: params.B, Rho: params.Rho, M: params.M, Sigma: params.Sigma + delta, TTM: ttm}, logMoneyness, totalVar, weights, lambda) - currentError) / delta
		
		// Update parameters
		params.A -= learningRate * gradA
		params.B -= learningRate * gradB
		params.Rho -= learningRate * gradRho
		params.M -= learningRate * gradM
		params.Sigma -= learningRate * gradSigma
		
		// Apply constraints with better bounds
		// For crypto, ensure minimum variance based on time to expiry
		minVar := 0.3 * 0.3 * params.TTM // 30% minimum annualized vol
		params.A = math.Max(minVar, params.A)
		params.B = math.Max(0.01, math.Min(2.0, params.B)) // Prevent extreme B values
		params.Rho = math.Max(-0.9, math.Min(0.9, params.Rho)) // Keep away from bounds
		params.M = math.Max(-1.0, math.Min(1.0, params.M)) // Reasonable M range
		params.Sigma = math.Max(0.1, math.Min(1.0, params.Sigma)) // Prevent extreme sigma
		
		// Check convergence
		newError := calculateRegularizedSSE(params, logMoneyness, totalVar, weights, lambda)
		if math.Abs(newError-currentError) < tolerance {
			break
		}
		
		// Adaptive learning rate
		if newError > currentError {
			learningRate *= 0.8
		} else if iter % 50 == 0 {
			learningRate *= 0.95
		}
	}
	
	// Verify no-arbitrage conditions
	if !checkSVINoArbitrage(params) {
		params = enforceSVINoArbitrage(params)
	}
	
	return params, nil
}

// calculateSSE calculates sum of squared errors for SVI parameters
func calculateSSE(params SVIParameters, logMoneyness, totalVar []float64) float64 {
	sse := 0.0
	for i := range logMoneyness {
		modelVar := params.computeTotalVariance(logMoneyness[i])
		diff := modelVar - totalVar[i]
		sse += diff * diff
	}
	return sse
}

// calculateRegularizedSSE calculates SSE with regularization penalties
func calculateRegularizedSSE(params SVIParameters, logMoneyness, totalVar, weights []float64, lambda float64) float64 {
	// Data fitting term
	sse := 0.0
	for i := range logMoneyness {
		modelVar := params.computeTotalVariance(logMoneyness[i])
		diff := (modelVar - totalVar[i]) * weights[i]
		sse += diff * diff
	}
	
	// Regularization terms
	reg := 0.0
	
	// Penalize extreme B values
	if params.B < 0.02 {
		reg += lambda * (0.02 - params.B) * (0.02 - params.B)
	}
	
	// Penalize extreme rho
	if math.Abs(params.Rho) > 0.95 {
		reg += lambda * (math.Abs(params.Rho) - 0.95) * (math.Abs(params.Rho) - 0.95)
	}
	
	// Penalize small sigma
	if params.Sigma < 0.05 {
		reg += lambda * (0.05 - params.Sigma) * (0.05 - params.Sigma)
	}
	
	// Butterfly arbitrage penalty
	butterflyPenalty := 0.0
	dk := 0.05
	for k := -0.5; k <= 0.5; k += dk {
		// Check second derivative is positive (convexity)
		w0 := params.computeTotalVariance(k - dk)
		w1 := params.computeTotalVariance(k)
		w2 := params.computeTotalVariance(k + dk)
		
		d2w := (w2 - 2*w1 + w0) / (dk * dk)
		if d2w < -0.001 { // Add small tolerance
			butterflyPenalty += lambda * 10 * d2w * d2w // Reduced penalty
		}
	}
	
	return sse + reg + butterflyPenalty
}

// checkSVINoArbitrage verifies SVI parameters satisfy no-arbitrage conditions
func checkSVINoArbitrage(params SVIParameters) bool {
	a := params.A
	b := params.B
	rho := params.Rho
	sigma := params.Sigma
	
	// Condition 1: a + b*sigma*sqrt(1-rho^2) >= 0
	if a + b*sigma*math.Sqrt(1-rho*rho) < 0 {
		return false
	}
	
	// Condition 2: b >= 0
	if b < 0 {
		return false
	}
	
	// Condition 3: b(1 + |rho|) <= 4
	if b*(1 + math.Abs(rho)) > 4 {
		return false
	}
	
	// Condition 4: sigma > 0
	if sigma <= 0 {
		return false
	}
	
	return true
}

// enforceSVINoArbitrage adjusts parameters to satisfy no-arbitrage conditions
func enforceSVINoArbitrage(params SVIParameters) SVIParameters {
	adjusted := params
	
	// Ensure b > 0
	adjusted.B = math.Max(0.01, adjusted.B)
	
	// Ensure sigma > 0
	adjusted.Sigma = math.Max(0.1, adjusted.Sigma)
	
	// Ensure b(1 + |rho|) <= 4
	maxB := 4.0 / (1 + math.Abs(adjusted.Rho))
	if adjusted.B > maxB {
		adjusted.B = maxB * 0.95
	}
	
	// Ensure a + b*sigma*sqrt(1-rho^2) >= 0
	minA := -adjusted.B * adjusted.Sigma * math.Sqrt(1-adjusted.Rho*adjusted.Rho)
	if adjusted.A < minA {
		adjusted.A = minA + 0.001
	}
	
	return adjusted
}

// GetIV returns implied volatility for a given strike and time to expiry
func (s *SVISurface) GetIV(strike, timeToExpiry float64) float64 {
	// Handle edge cases
	if timeToExpiry <= 0 {
		return 0
	}
	
	// Calculate forward price
	forward := s.SpotPrice * math.Exp(s.RiskFreeRate * timeToExpiry)
	logMoneyness := math.Log(strike / forward)
	
	// Find surrounding expiries
	i := 0
	for i < len(s.Expiries) && s.Expiries[i] < timeToExpiry {
		i++
	}
	
	if i == 0 {
		// Before first expiry - use first slice
		return s.Parameters[0].computeIV(logMoneyness)
	} else if i >= len(s.Expiries) {
		// After last expiry - use last slice
		return s.Parameters[len(s.Parameters)-1].computeIV(logMoneyness)
	}
	
	// Interpolate between two SVI slices
	t1, t2 := s.Expiries[i-1], s.Expiries[i]
	w1 := (t2 - timeToExpiry) / (t2 - t1)
	w2 := 1 - w1
	
	// Interpolate in total variance space
	var1 := s.Parameters[i-1].computeTotalVariance(logMoneyness)
	var2 := s.Parameters[i].computeTotalVariance(logMoneyness)
	
	totalVar := w1*var1 + w2*var2
	
	// Convert back to IV
	if totalVar <= 0 {
		return 0
	}
	return math.Sqrt(totalVar / timeToExpiry)
}

// FitVolatilitySurface fits SVI model to all expiry slices
func FitVolatilitySurface(options []OptionData, spotPrice, riskFreeRate float64) (*SVISurface, error) {
	// Group options by expiry
	expiryMap := make(map[float64][]OptionData)
	for _, opt := range options {
		if !math.IsNaN(opt.ImpliedVolatility) && opt.ImpliedVolatility > 0 {
			ttm := opt.TimeToExpiry
			expiryMap[ttm] = append(expiryMap[ttm], opt)
		}
	}
	
	// Sort expiries
	expiries := make([]float64, 0, len(expiryMap))
	for exp := range expiryMap {
		expiries = append(expiries, exp)
	}
	sortFloats(expiries)
	
	// Fit SVI to each expiry
	surface := &SVISurface{
		Expiries:     expiries,
		Parameters:   make([]SVIParameters, len(expiries)),
		SpotPrice:    spotPrice,
		RiskFreeRate: riskFreeRate,
	}
	
	for i, ttm := range expiries {
		opts := expiryMap[ttm]
		
		// Extract strikes and IVs
		strikes := make([]float64, len(opts))
		ivs := make([]float64, len(opts))
		for j, opt := range opts {
			strikes[j] = opt.Strike
			ivs[j] = opt.ImpliedVolatility
		}
		
		// Calculate forward price
		forward := spotPrice * math.Exp(riskFreeRate * ttm)
		
		// Fit SVI parameters
		params, err := FitSVISlice(strikes, ivs, forward, ttm)
		if err != nil {
			// Use simple interpolation if SVI fitting fails
			params = SVIParameters{
				A:     0.2 * 0.2 * ttm, // 20% vol as default
				B:     0.1,
				Rho:   -0.3,
				M:     0.0,
				Sigma: 0.3,
				TTM:   ttm,
			}
		}
		surface.Parameters[i] = params
	}
	
	return surface, nil
}

// sortFloats sorts a slice of float64 in ascending order
func sortFloats(floats []float64) {
	for i := 0; i < len(floats); i++ {
		for j := i + 1; j < len(floats); j++ {
			if floats[i] > floats[j] {
				floats[i], floats[j] = floats[j], floats[i]
			}
		}
	}
}

// FitVolatilitySurfaceWithConstraints fits SVI with inter-expiry constraints
func FitVolatilitySurfaceWithConstraints(options []OptionData, spotPrice, riskFreeRate float64) (*SVISurface, error) {
	// Group options by expiry
	expiryMap := make(map[float64][]OptionData)
	for _, opt := range options {
		if !math.IsNaN(opt.ImpliedVolatility) && opt.ImpliedVolatility > 0 {
			ttm := opt.TimeToExpiry
			expiryMap[ttm] = append(expiryMap[ttm], opt)
		}
	}
	
	// Sort expiries
	expiries := make([]float64, 0, len(expiryMap))
	for exp := range expiryMap {
		expiries = append(expiries, exp)
	}
	sortFloats(expiries)
	
	// Fit SVI to each expiry with constraints
	surface := &SVISurface{
		Expiries:     expiries,
		Parameters:   make([]SVIParameters, len(expiries)),
		SpotPrice:    spotPrice,
		RiskFreeRate: riskFreeRate,
	}
	
	var prevParams *SVIParameters
	var prevTTM float64
	
	for i, ttm := range expiries {
		opts := expiryMap[ttm]
		
		// Extract strikes and IVs
		strikes := make([]float64, len(opts))
		ivs := make([]float64, len(opts))
		for j, opt := range opts {
			strikes[j] = opt.Strike
			ivs[j] = opt.ImpliedVolatility
		}
		
		// Calculate forward price
		forward := spotPrice * math.Exp(riskFreeRate * ttm)
		
		// Fit SVI parameters with constraints from previous expiry
		params, err := FitSVIWithConstraints(strikes, ivs, forward, ttm, prevParams, prevTTM)
		if err != nil {
			// Use regularized fitting without constraints
			params, err = FitSVISlice(strikes, ivs, forward, ttm)
			if err != nil {
				// Fallback to default parameters
				// For crypto, use higher base volatility
				atmVol := 0.8 // 80% annual vol is reasonable for crypto
				params = SVIParameters{
					A:     atmVol * atmVol * ttm * 0.8, // Slightly below ATM
					B:     0.2,
					Rho:   -0.3,
					M:     0.0,
					Sigma: 0.3,
					TTM:   ttm,
				}
			}
		}
		
		surface.Parameters[i] = params
		prevParams = &params
		prevTTM = ttm
	}
	
	return surface, nil
}

// FitSVIWithConstraints fits SVI with calendar arbitrage constraints
func FitSVIWithConstraints(strikes []float64, ivs []float64, forward float64, 
	ttm float64, prevParams *SVIParameters, prevTTM float64) (SVIParameters, error) {
	
	// Use the regularized fitting as base
	params, err := FitSVISlice(strikes, ivs, forward, ttm)
	if err != nil {
		return params, err
	}
	
	// If we have previous parameters, enforce calendar constraints
	if prevParams != nil && prevTTM > 0 {
		// Check calendar arbitrage at key moneyness levels
		testMoneyness := []float64{0.8, 0.9, 1.0, 1.1, 1.2}
		
		needsAdjustment := false
		for _, moneyness := range testMoneyness {
			k := math.Log(moneyness)
			
			// Current total variance
			currVar := params.computeTotalVariance(k)
			
			// Previous total variance scaled to current time
			prevVar := prevParams.computeTotalVariance(k)
			prevVarScaled := prevVar * (ttm / prevTTM)
			
			// Check if calendar arbitrage exists
			if currVar < prevVarScaled {
				needsAdjustment = true
				break
			}
		}
		
		// If calendar arbitrage detected, adjust parameters
		if needsAdjustment {
			// Scale A parameter to ensure increasing variance
			minA := prevParams.A * (ttm / prevTTM) * 1.05 // 5% buffer
			if params.A < minA {
				params.A = minA
			}
			
			// Smooth parameter transitions
			smoothingFactor := 0.7
			params.B = smoothingFactor*params.B + (1-smoothingFactor)*prevParams.B
			params.Rho = smoothingFactor*params.Rho + (1-smoothingFactor)*prevParams.Rho
			params.Sigma = smoothingFactor*params.Sigma + (1-smoothingFactor)*prevParams.Sigma
		}
	}
	
	return params, nil
}