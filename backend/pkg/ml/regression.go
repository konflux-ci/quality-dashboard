package ml

import "math"

// LinearRegression represents a simple linear regression model.
type LinearRegression struct {
	Slope     float64
	Intercept float64
}

// Fit performs linear regression on the given data.
func (lr *LinearRegression) Fit(x, y []float64) {
	if len(x) != len(y) || len(x) == 0 {
		panic("Input data length mismatch or empty data")
	}

	// Calculate the mean of x and y
	meanX := mean(x)
	meanY := mean(y)

	// Calculate the slope (m) and intercept (b)
	numerator, denominator := 0.0, 0.0
	for i := 0; i < len(x); i++ {
		numerator += (x[i] - meanX) * (y[i] - meanY)
		denominator += math.Pow(x[i]-meanX, 2)
	}

	lr.Slope = numerator / denominator
	lr.Intercept = meanY - lr.Slope*meanX
}

// Predict predicts the dependent variable for a given independent variable.
func (lr *LinearRegression) Predict(x float64) float64 {
	return lr.Slope*x + lr.Intercept
}

// mean calculates the mean of a slice of float64 values.
func mean(values []float64) float64 {
	sum := 0.0
	for _, v := range values {
		sum += v
	}
	return sum / float64(len(values))
}
