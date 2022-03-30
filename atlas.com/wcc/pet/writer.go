package pet

import "atlas-wcc/socket/response"

func WriteForEachPet(w *response.Writer, ps []Model, pe func(w *response.Writer, p Model), pne func(w *response.Writer)) {
	for i := 0; i < 3; i++ {
		if ps != nil && len(ps) > i {
			pe(w, ps[i])
		} else {
			pne(w)
		}
	}
}

func WritePetItemId(w *response.Writer, p Model) {
	w.WriteInt(p.ItemId())
}

func WriteEmptyPetItemId(w *response.Writer) {
	w.WriteInt(0)
}
