package utils

func CalculateChange(latest float64, past *float64) *float64 {
	if past != nil && *past != 0 {
		change := ((latest - *past) / *past) * 100
		return &change
	}
	return nil
}
