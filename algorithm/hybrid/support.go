package hybrid

func score(target Target, state *state, ni, no uint) []float64 {
	nn := uint(len(state.Counts))
	scores := make([]float64, 0, no)
	for i, offset := uint(0), uint(0); i < nn; i++ {
		count := state.Counts[i]
		element := Element{
			Lindex:    state.Lindices[i*ni : (i+1)*ni],
			Indices:   state.Indices[offset*ni : (offset+count)*ni],
			Volumes:   state.Volumes[offset:(offset + count)],
			Values:    state.Observations[offset*no : (offset+count)*no],
			Surpluses: state.Surpluses[offset*no : (offset+count)*no],
		}
		scores = append(scores, target.Score(&element)...)
		offset += count
	}
	return scores
}
