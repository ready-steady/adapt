package global

func score(strategy Strategy, state *State, ni, no uint) []float64 {
	nn := uint(len(state.Counts))
	scores := make([]float64, nn)
	for i, o := uint(0), uint(0); i < nn; i++ {
		count := state.Counts[i]
		element := Element{
			Lindex:       state.Lindices[i*ni : (i+1)*ni],
			Indices:      state.Indices[o*ni : (o+count)*ni],
			Volumes:      state.Volumes[o:(o + count)],
			Observations: state.Observations[o*no : (o+count)*no],
			Surpluses:    state.Surpluses[o*no : (o+count)*no],
		}
		scores[i] = strategy.Score(&element)
		o += count
	}
	return scores
}
