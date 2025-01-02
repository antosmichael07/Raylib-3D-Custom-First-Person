package rlfp

// Gets the distance between 2 points without rooting the result to make it faster
//
// #1 argument a: float32 - the x position of the first point
//
// #2 argument b: float32 - the y position of the first point
//
// #3 argument c: float32 - the x position of the second point
//
// #4 argument d: float32 - the y position of the second point
//
// #1 return: float32 - the distance between the 2 points squared
func getDistance(a, b, c, d float32) float32 {
	return (a-c)*(a-c) + (b-d)*(b-d)
}
