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

// FitSVISlice fits SVI model to one expiry slice using simple gradient descent
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
	
	for i := range ivs {
		totalVar[i] = ivs[i] * ivs[i] * ttm
		logMoneyness[i] = math.Log(strikes[i] / forward)
	}
	
	// Find ATM index
	atmIdx := 0
	minDist := math.Abs(logMoneyness[0])
	for i := 1; i < len(logMoneyness); i++ {
		dist := math.Abs(logMoneyness[i])
		if dist < minDist {
			minDist = dist
			atmIdx = i
		}
	}
	
	// Initial parameter guess
	atmVar := totalVar[atmIdx]
	params := SVIParameters{
		A:     atmVar * 0.8,
		B:     0.1,
		Rho:   -0.3,
		M:     0.0,
		Sigma: 0.3,
		TTM:   ttm,
	}
	
	// Simple gradient descent optimization
	learningRate := 0.01
	maxIterations := 1000
	tolerance := 1e-6
	
	for iter := 0; iter < maxIterations; iter++ {
		// Calculate current error
		currentError := calculateSSE(params, logMoneyness, totalVar)
		
		// Calculate gradients numerically
		delta := 0.0001
		gradA := (calculateSSE(SVIParameters{A: params.A + delta, B: params.B, Rho: params.Rho, M: params.M, Sigma: params.Sigma, TTM: ttm}, logMoneyness, totalVar) - currentError) / delta
		gradB := (calculateSSE(SVIParameters{A: params.A, B: params.B + delta, Rho: params.Rho, M: params.M, Sigma: params.Sigma, TTM: ttm}, logMoneyness, totalVar) - currentError) / delta
		gradRho := (calculateSSE(SVIParameters{A: params.A, B: params.B, Rho: params.Rho + delta, M: params.M, Sigma: params.Sigma, TTM: ttm}, logMoneyness, totalVar) - currentError) / delta
		gradM := (calculateSSE(SVIParameters{A: params.A, B: params.B, Rho: params.Rho, M: params.M + delta, Sigma: params.Sigma, TTM: ttm}, logMoneyness, totalVar) - currentError) / delta
		gradSigma := (calculateSSE(SVIParameters{A: params.A, B: params.B, Rho: params.Rho, M: params.M, Sigma: params.Sigma + delta, TTM: ttm}, logMoneyness, totalVar) - currentError) / delta
		
		// Update parameters
		params.A -= learningRate * gradA
		params.B -= learningRate * gradB
		params.Rho -= learningRate * gradRho
		params.M -= learningRate * gradM
		params.Sigma -= learningRate * gradSigma
		
		// Apply constraints
		params.A = math.Max(0.0001, params.A)
		params.B = math.Max(0.0001, params.B)
		params.Rho = math.Max(-0.999, math.Min(0.999, params.Rho))
		params.Sigma = math.Max(0.0001, params.Sigma)
		
		// Check convergence
		newError := calculateSSE(params, logMoneyness, totalVar)
		if math.Abs(newError-currentError) < tolerance {
			break
		}
		
		// Adaptive learning rate
		if newError > currentError {
			learningRate *= 0.5
		}
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