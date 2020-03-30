package queryables

import "net/http"

// Collection of queryable info
type Collection []*Info

func (coll *Collection) Read(r *http.Request) map[string]interface{} {
	res := make(map[string]interface{})
	for _, i := range *coll {
		key := i.DaoKey
		val := i.Value(r)

		// skip if value is nil
		if val == nil {
			continue
		}

		if i.Transform != nil {
			key, val = i.Transform(i.DaoKey, val)
		}

		res[key] = val
	}

	return res
}
