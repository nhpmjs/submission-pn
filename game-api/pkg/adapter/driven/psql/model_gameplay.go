package psql

func (gp *GamePlay) nextUser() {
	if gp.isLastUser() {
		// next frame
		gp.CurrentFrame = gp.CurrentFrame + 1
	}

	gp.CurrentUser = (gp.CurrentUser + 1) % len(gp.Participants)
	gp.CurrentRoll = 1
}

func (gp *GamePlay) isLastUser() bool {
	return gp.CurrentUser == len(gp.Participants)-1
}

func (gp *GamePlay) nextRoll() {
	if gp.CurrentFrame < 10 {
		if gp.CurrentRoll == 2 {
			if gp.isLastUser() {
				// next frame
				gp.CurrentFrame = gp.CurrentFrame + 1
			}
			// next user
			gp.CurrentUser = (gp.CurrentUser + 1) % len(gp.Participants)
		}

		gp.CurrentRoll = increment(gp.CurrentRoll)
	} else {
		gp.CurrentRoll = gp.CurrentRoll + 1
	}
}
